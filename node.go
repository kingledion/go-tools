package tree

import "fmt"

type Node interface {
	GetChildren() []Node
	IsParent(id uint64) bool
	GetID() uint64
	Add(...Node)
	GetParent() Node
	AddParent(Node)
}

type node struct {
	primary  uint64
	parentID uint64
	parent   Node
	data     interface{}
	children []Node
}

func (n *node) GetChildren() []Node {
	return n.children
}

func (n *node) GetID() uint64 {
	return n.primary
}

func (n *node) Add(children ...Node) {
	n.children = append(n.children, children[:]...)
}

func (n *node) IsParent(id uint64) bool {
	return n.parentID == id
}

func (n *node) GetParent() Node {
	return n.parent
}

func (n *node) AddParent(parent Node) {
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
