package fbx

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"os"

	"github.com/mogaika/binrw"
	"github.com/pkg/errors"
)

func readProperty(r *binrw.Reader) (interface{}, error) {
	return DataType(r.ReadU8()).read(r)
}

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

func Read(sourceR io.ReadSeeker) (*FBX, error) {
	r := binrw.NewReader(sourceR, binary.LittleEndian)

	magic := r.ReadBuf(RAW_MAGIC_SIZE)
	if bytes.Compare(magic, RAW_MAGIC) != 0 {
		return nil, errors.Errorf("Invalid header magic: %s", hex.EncodeToString(magic))
	}

	r.Skip(2) // unknown

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
