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

func NewNode(name string, properties ...interface{}) *Node {
	return &Node{
		Name:       name,
		Properties: properties,
		Nodes:      make([]*Node, 0),
	}
}

func (n *Node) AddNode(node *Node) *Node {
	n.Nodes = append(n.Nodes, node)
	return n
}

func (n *Node) AddNodes(nodes ...*Node) *Node {
	n.Nodes = append(n.Nodes, nodes...)
	return n
}

func (n *Node) AddNewNode(name string, createParams ...interface{}) *Node {
	return n.AddNode(NewNode(name, createParams))
}

func (n *Node) AddProperty(property interface{}) *Node {
	n.Properties = append(n.Properties, property)
	return n
}

func (n *Node) AddProperties(properties ...interface{}) *Node {
	n.Properties = append(n.Properties, properties...)
	return n
}

func (n *Node) GetNode(name string) *Node {
	for _, node := range n.Nodes {
		if node.Name == name {
			return node
		}
	}
	return nil
}

func (n *Node) GetNodes(name string) []*Node {
	nodes := make([]*Node, 0)
	for _, node := range n.Nodes {
		if node.Name == name {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (n *Node) GetOrAddNode(node *Node) *Node {
	result := n.GetNode(node.Name)
	if result == nil {
		n.AddNode(node)
		return node
	} else {
		return result
	}
}

func (n *Node) GetOrAddNewNode(name string, createParams ...interface{}) *Node {
	result := n.GetNode(name)
	if result == nil {
		node := NewNode(name, createParams)
		n.AddNode(node)
		return node
	} else {
		return result
	}
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
