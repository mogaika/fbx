package fbx

import (
	"fmt"
	"reflect"
	"strings"
)

type Node struct {
	Name       string
	Nodes      []*Node
	Properties []interface{}
}

func (n *Node) AddNode(node *Node) *Node {
	n.Nodes = append(n.Nodes, node)
	return n
}

func (n *Node) AddProperty(prop interface{}) *Node {
	n.Properties = append(n.Properties, prop)
	return n
}

func (n *Node) AddProperties(props ...interface{}) *Node {
	n.Properties = append(n.Properties, props...)
	return n
}

func (n *Node) sprint(sb *strings.Builder, depth int) {
	tab := func(amount int) {
		for i := 0; i < amount; i++ {
			sb.WriteString("|   ")
		}
	}

	if depth >= 0 {
		tab(depth)
		sb.WriteString(fmt.Sprintf("node %q\n", n.Name))
	}
	for _, property := range n.Properties {
		tab(depth + 1)
		sb.WriteString(fmt.Sprintf("- property: (%s) %+#v\n", reflect.TypeOf(property).String(), property))
	}
	for _, node := range n.Nodes {
		node.sprint(sb, depth+1)
	}
}

func (n *Node) SPrint() string {
	var sb strings.Builder
	n.sprint(&sb, 0)
	return sb.String()
}
