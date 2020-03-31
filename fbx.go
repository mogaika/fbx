package fbx

import (
	"strings"
)

type FBX struct {
	Version int
	Root    *Node
}

const RAW_MAGIC_SIZE = 0x15
const RAW_HEADER_SIZE = 0x19
const RAW_NULL_ENTRY_SIZE = 0xb

// "Kaydara FBX Binary  \x00"
var RAW_MAGIC []byte = []byte{
	0x4B, 0x61, 0x79, 0x64, 0x61, 0x72, 0x61, 0x20,
	0x46, 0x42, 0x58, 0x20, 0x42, 0x69, 0x6E, 0x61,
	0x72, 0x79, 0x20, 0x20, 0x00}

var RAW_NULL_ENTRY []byte = make([]byte, RAW_NULL_ENTRY_SIZE)

func (f *FBX) SPrint() string {
	var sb strings.Builder
	f.Root.sprint(&sb, -1)
	return sb.String()
}
