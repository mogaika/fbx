package fbx

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"

	"github.com/mogaika/binrw"
	"github.com/pkg/errors"
)

func readNode(r *binrw.Reader) (*Node, error) {
	n := &Node{
		Nodes: make([]*Node, 0),
	}

	// curNodeOffset := r.Offset()

	rawEndOffset := r.ReadU32()
	rawNumProperties := r.ReadU32()
	rawPropertyListLen := r.ReadU32()
	rawNameLen := r.ReadU8()

	n.Name = string(r.ReadBuf(int(rawNameLen)))
	// log.Printf("reading node %q offset 0x%x end 0x%x", n.Name, curNodeOffset, rawEndOffset)

	offBeforeProperty := r.Offset()

	n.Properties = make([]interface{}, rawNumProperties)
	for iProperty := range n.Properties {
		if property, err := readProperty(r); err != nil {
			return nil, errors.Wrapf(err, "Error reading property index %v", iProperty)
		} else {
			n.Properties[iProperty] = property
		}
	}

	r.Seek(offBeforeProperty+int64(rawPropertyListLen), os.SEEK_SET)

	if r.Offset() < int64(rawEndOffset)-RAW_NULL_ENTRY_SIZE {
		if nodes, err := readNodes(r); err != nil {
			return nil, errors.Wrapf(err, "Error reading nodes of node %q", n.Name)
		} else {
			n.Nodes = nodes
		}
	}

	r.Seek(int64(rawEndOffset), os.SEEK_SET)
	return n, r.Error()
}

func readNodes(r *binrw.Reader) ([]*Node, error) {
	nodes := make([]*Node, 0)
	for {
		nullBuf := r.PeekBuf(RAW_NULL_ENTRY_SIZE)
		if bytes.Compare(nullBuf, RAW_NULL_ENTRY) == 0 {
			r.Skip(RAW_NULL_ENTRY_SIZE)
			return nodes, r.Error()
		}

		if subNode, err := readNode(r); err != nil {
			return nil, errors.Wrapf(err, "Failed to read node")
		} else {
			nodes = append(nodes, subNode)
		}
	}
}

func readProperty(r *binrw.Reader) (interface{}, error) {
	return readPropertyData(r, DataType(r.ReadU8()))
}

func readPropertyData(r *binrw.Reader, dt DataType) (interface{}, error) {
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

func Read(sourceR io.ReadSeeker) (*FBX, error) {
	r := binrw.NewReader(sourceR, binary.LittleEndian)

	magic := r.ReadBuf(RAW_MAGIC_SIZE)
	if bytes.Compare(magic, RAW_MAGIC) != 0 {
		return nil, errors.Errorf("Invalid header magic: %s", hex.EncodeToString(magic))
	}

	r.Skip(2) // unknown, usually 0x1a

	f := &FBX{
		Version: int(r.ReadU32()),
	}

	if nodes, err := readNodes(r); err != nil {
		return nil, errors.Wrapf(err, "Unable to parse root node")
	} else {
		f.Root.Nodes = nodes
	}

	return f, nil
}
