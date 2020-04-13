package bfbx73

import (
	"github.com/mogaika/fbx"
)

func C(t string, childId int64, parentId int64, v ...interface{}) *fbx.Node {
	params := make([]interface{}, 3, 3+len(v))
	params[0] = t
	params[1] = childId
	params[2] = parentId
	return fbx.NewNode("C", append(params, v...)...)
}

func P(a1, a2, a3, a4 string, v ...interface{}) *fbx.Node {
	params := make([]interface{}, 4, 4+len(v))
	params[0] = a1
	params[1] = a2
	params[2] = a3
	params[3] = a4
	return fbx.NewNode("P", append(params, v...)...)
}
