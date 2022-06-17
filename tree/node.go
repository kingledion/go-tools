package tree

import (
	"fmt"
)

// Node is the interface for a node within a tree.
//
// This package implementes this interface with an internal struct
// representation. This interface cannot be implemented by other packages
// due to unexported fields.
//
// The properties of a Node are that it has a uint primary key, a parent
// Node and an array of child Nodes. Instantiating a Node through the
// Tree.Add() function allows the node to be created with a parentID. This
// represents the primary key ID of a parent that may not be added to the
// tree yet. This value can be specified to indentify parents added to the
// tree after their children.
//
// The node also has arbitrary data attached to it; any struct may be inserted
// as the data. However, because the data is not specified, this trait contains
// no functionality to modify the data in place. The data should be implemented
// as a pointer to avoid performance ramifications.
//
// The associated data attached to a node must be serializable
// using the json encoding. If either node or arbitrary data is not, the
// Tree.Serialize() function will throw serialization errors when called.
// There is no serialization checking at the time that a Node is added to a tree
// or that its data is updated.
//
// When this node is serialized, only the id and parent id, along with
// associated data are serialized. The internal pointers to parent and children
// are recreated with the nodes are deserialized.
type Node interface {
	// GetID returns the primary key of this node.
	GetID() uint
	// GetParentID returns the primary key of this node's parent.
	GetParentID() uint

	// GetChildren returns an array of pointers to all children of this node.
	GetChildren() []Node
	// GetParent returns a pointer to the parent node of this node.
	GetParent() Node

	// AddChildren adds a list of Nodes as children of this node.
	AddChildren(...Node)
	// ReplaceChildren replaces the current list of children with a new list of
	// Nodes.
	ReplaceChildren(...Node)

	setParent(n Node)

	// GetData retruns this node's internal data.
	GetData() interface{}
	// SetData replaces this nodes data with the argument. The argument may be any
	// type, but must be serializable via json.
	//
	// This function does not attempt to test json encoding when the data is set;
	// any error with encoding will only occur when the data is serialized
	// to a repository.
	SetData(interface{})
}

type node struct {
	Primary  uint
	ParentID uint
	parent   Node
	Data     interface{}
	children []Node
}

func (n *node) GetID() uint {
	return n.Primary
}

func (n *node) GetParentID() uint {
	return n.ParentID
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

func (n *node) setParent(parent Node) {
	n.parent = parent
	n.ParentID = parent.GetID()
}

func (n *node) GetData() interface{} {
	return n.Data
}

func (n *node) SetData(newdata interface{}) {
	n.Data = newdata
}

func (n *node) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(f, "{primary: %d parentID: %d data:%+v children:[", n.Primary, n.ParentID, n.Data)
		for i, n := range n.children {
			if i != 0 {
				fmt.Fprint(f, " ")
			}
			fmt.Fprintf(f, "%d", n.GetID())
		}
		fmt.Fprint(f, "]}")
	}
}
