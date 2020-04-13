package fbx

import (
	"fmt"
	"log"
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

var RAW_NULL_FOOTER_MAGIC []byte = []byte{
	0xf8, 0x5a, 0x8c, 0x6a, 0xde, 0xf5, 0xd9, 0x7e,
	0xec, 0xe9, 0xc, 0xe3, 0x75, 0x8f, 0x29, 0xb,
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

func (f *FBX) PrintConnectionsTree(id int64) { f.printConnectionsTree(id, 0) }
func (f *FBX) printConnectionsTree(id int64, tab int) {
	connections := f.Root.GetNode("Connections")
	objects := f.Root.GetNode("Objects")

	if tab > 30 {
		panic(tab)
	}
	sTab := func(i int) string {
		s := ""
		for j := 0; j < i; j++ {
			s += "|  "
		}
		return s
	}
	for _, conn := range connections.Nodes {
		if conn.Properties[2].(int64) != id {
			continue
		}
		childId := conn.Properties[1].(int64)
		for _, object := range objects.Nodes {
			if object.Properties[0].(int64) == childId {
				log.Printf("%s %v %q %q",
					sTab(tab), childId, object.Properties[1], object.Properties[2])
				f.printConnectionsTree(childId, tab+1)
			}
		}
	}
}

func (f *FBX) PrintConnectionsList(id int64) { f.printConnectionsList(id, 0) }
func (f *FBX) printConnectionsList(id int64, tab int) {
	connections := f.Root.GetNode("Connections")
	objects := f.Root.GetNode("Objects")

	getObject := func(id int64) *Node {
		for _, object := range objects.Nodes {
			if object.Properties[0].(int64) == id {
				return object
			}
		}
		return nil
	}
	for _, object := range connections.Nodes {
		child := getObject(object.Properties[1].(int64))
		parent := getObject(object.Properties[2].(int64))
		if parent != nil {
			log.Printf("%s c %09d p %09d  c %q %q   p %q %q",
				object.Properties[0].(string), object.Properties[1].(int64), object.Properties[2].(int64),
				child.Properties[1], child.Properties[2], parent.Properties[1], parent.Properties[2])
		} else {
			log.Printf("%s c %09d p %09d  c %q %q | root ",
				object.Properties[0].(string), object.Properties[1].(int64), object.Properties[2].(int64),
				child.Properties[1], child.Properties[2])
		}
	}
}
