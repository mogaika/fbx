package fbx

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io/ioutil"
	"reflect"

	"github.com/mogaika/binrw"
	"github.com/pkg/errors"
)

type DataType rune

const TYPE_BOOL DataType = 'C'
const TYPE_INT16 DataType = 'Y'
const TYPE_INT32 DataType = 'I'
const TYPE_INT64 DataType = 'L'
const TYPE_FLOAT32 DataType = 'F'
const TYPE_FLOAT64 DataType = 'D'

const TYPE_ARRAY_FLOAT32 DataType = 'f'
const TYPE_ARRAY_FLOAT64 DataType = 'd'
const TYPE_ARRAY_INT64 DataType = 'l'
const TYPE_ARRAY_INT32 DataType = 'i'
const TYPE_ARRAY_BOOL DataType = 'b'

const TYPE_STRING DataType = 'S'
const TYPE_RAW DataType = 'R'

func (dt DataType) IsArray() bool {
	switch dt {
	case TYPE_ARRAY_FLOAT64, TYPE_ARRAY_FLOAT32,
		TYPE_ARRAY_INT64, TYPE_ARRAY_INT32,
		TYPE_ARRAY_BOOL:
		return true
	}
	return false
}

func (dt DataType) IsSpecial() bool {
	switch dt {
	case TYPE_STRING, TYPE_RAW:
		return true
	}
	return false
}

func dataTypeFromInterface(v interface{}) (dt DataType, err error) {
	switch v.(type) {
	case bool:
		dt = TYPE_BOOL
	case int16:
		dt = TYPE_INT16
	case int32:
		dt = TYPE_INT32
	case int64:
		dt = TYPE_INT64
	case float32:
		dt = TYPE_FLOAT32
	case float64:
		dt = TYPE_FLOAT64
	case []bool:
		dt = TYPE_ARRAY_BOOL
	case []int32:
		dt = TYPE_ARRAY_INT32
	case []int64:
		dt = TYPE_ARRAY_INT64
	case []float32:
		dt = TYPE_ARRAY_FLOAT32
	case []float64:
		dt = TYPE_ARRAY_FLOAT64
	case []byte:
		dt = TYPE_RAW
	case string:
		dt = TYPE_STRING
	default:
		err = errors.Errorf("Can't find data type for type %v", reflect.TypeOf(v))
	}
	return dt, err
}

func (dt DataType) read(r *binrw.Reader) (interface{}, error) {
	if dt.IsArray() {
		arrayElementsCount := int(r.ReadU32())
		encoding := r.ReadU32()
		compressedLength := r.ReadU32()

		arrayData := r.ReadBuf(int(compressedLength))

		switch encoding {
		case 0:
		case 1:
			compressedReader, err := zlib.NewReader(bytes.NewReader(arrayData))
			if err != nil {
				return nil, errors.Wrapf(err, "Unable to init zlib")
			}

			if uncompressedData, err := ioutil.ReadAll(compressedReader); err != nil {
				return nil, errors.Wrapf(err, "Unable to unpack zlib")
			} else {
				arrayData = uncompressedData
			}
		default:
			return nil, errors.Errorf("Unknown array encoding %v", encoding)
		}

		if dt == TYPE_ARRAY_BOOL {
			bools := make([]bool, arrayElementsCount)
			for i := range bools {
				bools[i] = arrayData[i]&1 != 0
			}
			return bools, r.Error()
		}

		var arr interface{}
		switch dt {
		case TYPE_ARRAY_INT32:
			arr = make([]int32, arrayElementsCount)
		case TYPE_ARRAY_INT64:
			arr = make([]int64, arrayElementsCount)
		case TYPE_ARRAY_FLOAT32:
			arr = make([]float32, arrayElementsCount)
		case TYPE_ARRAY_FLOAT64:
			arr = make([]float64, arrayElementsCount)
		default:
			return nil, errors.Errorf("Unknown array type %q", string(dt))
		}

		if err := binary.Read(bytes.NewReader(arrayData), binary.LittleEndian, arr); err != nil {
			return nil, errors.Wrapf(err, "Can't read array of property type %q", string(dt))
		} else {
			return arr, r.Error()
		}
	} else if dt.IsSpecial() {
		data := r.ReadBuf(int(r.ReadU32()))
		switch dt {
		case TYPE_STRING:
			return string(data), r.Error()
		case TYPE_RAW:
			return data, r.Error()
		}
	} else {
		switch dt {
		case TYPE_BOOL:
			return r.ReadU8()&1 != 0, r.Error()
		case TYPE_INT16:
			return r.ReadI16(), r.Error()
		case TYPE_INT32:
			return r.ReadI32(), r.Error()
		case TYPE_INT64:
			return r.ReadI64(), r.Error()
		case TYPE_FLOAT32:
			return r.ReadF32(), r.Error()
		case TYPE_FLOAT64:
			return r.ReadF64(), r.Error()
		}
	}
	return nil, errors.Errorf("Invalid type %q for reading", string(dt))
}

func (dt DataType) write(w *binrw.Writer, v interface{}) error {
	if dt.IsArray() {
		vType := reflect.ValueOf(v)
		arrayLength := vType.Len()

		var data []byte
		if dt == TYPE_ARRAY_BOOL {
			data = make([]byte, arrayLength)
			for i, elV := range v.([]bool) {
				if elV {
					data[i] = 0
				} else {
					data[i] = 1
				}
			}
		} else {
			var buf bytes.Buffer
			if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
				return errors.Wrapf(err, "Can't write array of property type %v", vType)
			}
			data = buf.Bytes()
		}

		var compressedBuffer bytes.Buffer
		compressedWriter := zlib.NewWriter(&compressedBuffer)

		if _, err := compressedWriter.Write(data); err != nil {
			return errors.Wrapf(err, "Unable to pack zlib")
		}

		if err := compressedWriter.Close(); err != nil {
			return errors.Wrapf(err, "Unable to close zlib")
		}

		w.WriteU32(uint32(arrayLength))
		if compressedBuffer.Len() < len(data) {
			w.WriteU32(1)
			data = compressedBuffer.Bytes()
		} else {
			w.WriteU32(0)
		}
		w.WriteU32(uint32(len(data)))
		w.Write(data)
	} else if dt.IsSpecial() {
		var data []byte
		switch dt {
		case TYPE_STRING:
			data = []byte(v.(string))
		case TYPE_RAW:
			data = v.([]byte)
		}
		w.WriteU32(uint32(len(data)))
		w.Write(data)
	} else {
		switch dt {
		case TYPE_BOOL:
			if v.(bool) {
				w.WriteU8(1)
			} else {
				w.WriteU8(0)
			}
		case TYPE_INT16:
			w.WriteI16(v.(int16))
		case TYPE_INT32:
			w.WriteI32(v.(int32))
		case TYPE_INT64:
			w.WriteI64(v.(int64))
		case TYPE_FLOAT32:
			w.WriteF32(v.(float32))
		case TYPE_FLOAT64:
			w.WriteF64(v.(float64))
		default:
			return errors.Errorf("Invalid type %q for writing", string(dt))
		}
	}
	return w.Error()
}
