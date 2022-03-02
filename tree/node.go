package tree

import "fmt"

// Node is the interface for a node within this tree. This package has an
// internal node representation, but any custom structure meeting this
// interface can be substituted.
//
// Any implementation of Node should exists as a pointer to a structure, or
// else there will be serious performance ramifications.
type Node interface {
	// GetChildren returns an array of all children on this node.
	GetChildren() []Node
	// IsParent returns a boolean true if the id parameter matches the parent
	// of this node; and returns false otherwise.
	IsParent(id uint) bool
	// GetID returns the primary key of this node.
	GetID() uint
	// Add adds a list of nodes as children of this node.
	Add(...Node)
	// GetParent returns the parent node of this node.
	GetParent() Node

	SetParent(n Node)
}

type node struct {
	primary  uint
	parentID uint
	parent   Node
	data     interface{}
	children []Node
}

func (n *node) GetChildren() []Node {
	return n.children
}

func (n *node) GetID() uint {
	return n.primary
}

func (n *node) Add(children ...Node) {
	if n.children == nil {
		n.children = []Node{}
	}
	n.children = append(n.children, children[:]...)
}

func (n *node) IsParent(id uint) bool {
	return n.parentID == id
}

func (n *node) GetParent() Node {
	return n.parent
}

func (n *node) SetParent(parent Node) {
	n.parent = parent
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
