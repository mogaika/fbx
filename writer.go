package fbx

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"os"
	"reflect"

	"github.com/mogaika/binrw"
	"github.com/pkg/errors"
)

func writeNode(w *binrw.Writer, n *Node) error {
	posForRawEnd := w.Offset()
	w.Skip(4 + 4 + 4)

	w.WriteU8(uint8(len(n.Name)))
	w.Write([]byte(n.Name))

	offBeforeProperty := w.Offset()
	for iProperty, property := range n.Properties {
		if err := writeProperty(w, property); err != nil {
			return errors.Wrapf(err, "Error writing property index %v", iProperty)
		}
	}
	offAfterProperty := w.Offset()

	if n.Nodes != nil && len(n.Nodes) != 0 {
		if err := writeNodes(w, n.Nodes); err != nil {
			return errors.Wrapf(err, "Error writing nodes of node %q", n.Name)
		}
	}

	endOffset := w.Offset()

	w.Seek(posForRawEnd, os.SEEK_SET)
	w.WriteU32(uint32(endOffset))
	w.WriteU32(uint32(len(n.Properties)))
	w.WriteU32(uint32(offAfterProperty - offBeforeProperty))
	w.Seek(endOffset, os.SEEK_SET)

	return w.Error()
}

func writeNodes(w *binrw.Writer, nodes []*Node) error {
	for _, n := range nodes {
		if err := writeNode(w, n); err != nil {
			return errors.Wrapf(err, "Error writing node")
		}
	}
	w.Write(RAW_NULL_ENTRY)
	return w.Error()
}

func writeProperty(w *binrw.Writer, v interface{}) error {
	if dt, err := dataTypeFromInterface(v); err != nil {
		return err
	} else {
		w.WriteU8(uint8(dt))
		return writePropertyData(w, dt, v)
	}
}
func writePropertyData(w *binrw.Writer, dt DataType, v interface{}) error {
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

		if len(data) >= 1024 {
			compressedWriter := zlib.NewWriter(&compressedBuffer)

			if _, err := compressedWriter.Write(data); err != nil {
				return errors.Wrapf(err, "Unable to pack zlib")
			}

			if err := compressedWriter.Close(); err != nil {
				return errors.Wrapf(err, "Unable to close zlib")
			}
		}

		w.WriteU32(uint32(arrayLength))
		if compressedBuffer.Len() != 0 && compressedBuffer.Len() < len(data) {
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

func Write(sourceW io.WriteSeeker, f *FBX) error {
	w := binrw.NewWriter(sourceW, binary.LittleEndian)

	w.Write(RAW_MAGIC)
	w.WriteU16(0x1a)
	w.WriteU32(uint32(f.Version))

	return writeNodes(w, f.Root.Nodes)
}
