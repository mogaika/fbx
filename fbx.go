package fbx

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type FBX struct {
	Version int
	Root    Node
}

func NewFBX(version int) *FBX {
	return &FBX{
		Version: version,
		Root:    *NewNode(""),
	}
}

func (f *FBX) SPrint() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Binary version: %v\n", f.Version))
	f.Root.sprint(&sb, -1)
	return sb.String()
}

const RAW_MAGIC_SIZE = 0x15
const RAW_HEADER_SIZE = 0x19
const RAW_NULL_ENTRY_SIZE = 0xd
const RAW_NULL_FOOTER_1_SIZE = 7
const RAW_NULL_FOOTER_2_SIZE = 0x88

// "Kaydara FBX Binary  \x00"
var RAW_MAGIC []byte = []byte{
	0x4B, 0x61, 0x79, 0x64, 0x61, 0x72, 0x61, 0x20,
	0x46, 0x42, 0x58, 0x20, 0x42, 0x69, 0x6E, 0x61,
	0x72, 0x79, 0x20, 0x20, 0x00}

var RAW_NULL_ENTRY []byte = make([]byte, RAW_NULL_ENTRY_SIZE)

var RAW_FOOTER_SOURCE []byte = []byte{
	0x58, 0xAB, 0xA9, 0xF0, 0x6C, 0xA2, 0xD8, 0x3F,
	0x4D, 0x47, 0x49, 0xA3, 0xB4, 0xB2, 0xE7, 0x3D,
}

var RAW_FOOTER_KEY []byte = []byte{
	0xE2, 0x4F, 0x7B, 0x5F, 0xCD, 0xE4, 0xC8, 0x6D,
	0xDB, 0xD8, 0xFB, 0xD7, 0x40, 0x58, 0xC6, 0x78,
}

func generateFooter(f *FBX) ([]byte, error) {
	headerExtension := f.Root.GetNode("FBXHeaderExtension")
	if headerExtension == nil {
		return nil, errors.New("FBXHeaderExtension not found")
	}
	creationTimestamp := headerExtension.GetNode("CreationTimeStamp")
	if creationTimestamp == nil {
		return nil, errors.New("CreationTimeStamp not found")
	}

	getTimePartOrZero := func(nodeName string) int32 {
		if n := creationTimestamp.GetNode(nodeName); n != nil {
			if len(n.Properties) > 0 {
				return n.Properties[0].(int32)
			}
		}
		return 0
	}

	year := getTimePartOrZero("Year")
	month := getTimePartOrZero("Month")
	day := getTimePartOrZero("Day")
	hour := getTimePartOrZero("Hour")
	minute := getTimePartOrZero("Minute")
	second := getTimePartOrZero("Second")
	millisecond := getTimePartOrZero("Millisecond")

	encode := func(buf []byte, in []byte) {
		needle := byte(64)
		for i := range buf {
			buf[i] ^= needle ^ in[i]
			needle = buf[i]
		}
	}

	timeEncoded := []byte(fmt.Sprintf("%02d%02d%02d%02d%02d%04d%02d",
		second, month, hour, day, millisecond/10, year, minute))

	buf := make([]byte, 16)
	copy(buf, RAW_FOOTER_SOURCE)
	encode(buf, timeEncoded)
	encode(buf, RAW_FOOTER_KEY)
	encode(buf, timeEncoded)
	return buf, nil
}
