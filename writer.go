package fbx

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/mogaika/binrw"
	"github.com/pkg/errors"
)

func writeProperty(w *binrw.Writer, v interface{}) error {
	if dt, err := dataTypeFromInterface(v); err != nil {
		return err
	} else {
		w.WriteU8(uint8(dt))
		return dt.write(w, v)
	}
}

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

func Write(sourceW io.WriteSeeker, f *FBX) error {
	w := binrw.NewWriter(sourceW, binary.LittleEndian)

	w.Write(RAW_MAGIC)
	w.WriteU16(0x1a)
	w.WriteU32(uint32(f.Version))

	return writeNodes(w, f.Root.Nodes)
}
