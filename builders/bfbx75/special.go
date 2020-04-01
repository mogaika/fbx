package bfbx75

import (
	"github.com/mogaika/fbx"
)

func C(t string, id1 int64, id2 int64, v ...interface{}) *fbx.Node {
	params := make([]interface{}, 3, 3+len(v))
	params[0] = t
	params[1] = id1
	params[2] = id2
	return fbx.NewNode("P", append(params, v...)...)
}

func P(a1, a2, a3, a4 string, v ...interface{}) *fbx.Node {
	params := make([]interface{}, 4, 4+len(v))
	params[0] = a1
	params[1] = a2
	params[2] = a3
	params[4] = a4
	return fbx.NewNode("P", append(params, v...)...)
}
