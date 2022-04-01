package tree

import "fmt"

// Node is the interface for a node within this tree. This package has an
// internal node representation, but any custom structure meeting this
// interface can be substituted.
//
// Any implementation of Node should exists as a pointer to a structure, or
// else there will be serious performance ramifications.
//
// The properties of a node are that it has a uint primary key, a parent
// Node and an array of child Nodes. Instantiating a Node through the
// Tree.Add() function allows the node to be created with a parentID. This
// represents the primary key ID of a parent that may not be added to the
// tree yet. This value can be specified to indentify parents added to the
// tree after their children.
type Node interface {
	// GetID returns the primary key of this node.
	GetID() uint
	// GetParentID returns the primary key of a node's parent
	GetParentID() uint

	// GetChildren returns an array of all children on this node.
	GetChildren() []Node
	// GetParent returns the parent node of this node.
	GetParent() Node

	// Add adds a list of nodes as children of this node.
	AddChildren(...Node)
	// ReplaceChildren replaces the current list of children with a new list
	ReplaceChildren(...Node)
	// SetParent sets this node's parent to be the argument Node
	SetParent(n Node)

	// GetData retruns the node's data.
	GetData() interface{}
	// ReplaceData replaces nodes data with the argument.
	ReplaceData(interface{})
}

type node struct {
	primary  uint
	parentID uint
	parent   Node
	data     interface{}
	children []Node
}

func (n *node) GetID() uint {
	return n.primary
}

func (n *node) GetParentID() uint {
	return n.parentID
}

func (n *node) GetChildren() []Node {
	return n.children
}

func (n *node) GetParent() Node {
	return n.parent
}

func (n *node) AddChildren(children ...Node) {
	if n.children == nil {
		n.children = []Node{}
	}
	n.children = append(n.children, children[:]...)
}

func (n *node) ReplaceChildren(children ...Node) {
	n.children = []Node{}
	n.AddChildren(children...)
}

func (n *node) SetParent(parent Node) {
	n.parent = parent
}

func (n *node) GetData() interface{} {
	return n.data
}

func (n *node) ReplaceData(newData interface{}) {
	n.data = newData
}

func (n *node) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(f, "{primary: %d parentID: %d data:%+v children:[", n.primary, n.parentID, n.data)
		for i, n := range n.children {
			if i != 0 {
				fmt.Fprint(f, " ")
			}
			fmt.Fprintf(f, "%d", n.GetID())
		}
		fmt.Fprint(f, "]}")
	}
}
