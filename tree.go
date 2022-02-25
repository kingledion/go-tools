package tree

import "fmt"

// Tree represents a tree structure
type Tree struct {
	head    Node
	primary *index
}

// Empty returns an empty tree
func Empty() Tree {
	return Tree{
		primary: &index{},
	}
}

func (t *Tree) Head() Node {
	return t.head
}

// Add inserts an element into a tree. If the element's parent is not found,
// then the element is not inserted and false is returned. If the element
// is added successfully, true is returned.
func (t *Tree) Add(nodeID uint64, parentID uint64, data interface{}) (bool, bool) {
	child := &node{primary: nodeID, parentID: parentID, data: data}

	// Return false if this element has already been added
	if t.primary.find(nodeID) != nil {
		return false, true
	}

	if t.head == nil { // always insert the first element
		t.head = child
	} else {

		// return false if the parent does not exist and this element is
		// not a missing parent
		parent := t.primary.find(parentID)
		if parent == nil {

			if t.head.IsParent(nodeID) {
				t.Rehead(child)
			} else {
				return false, false
			}
		} else {
			child.AddParent(parent)
			parent.Add(child)
		}
	}

	// add to primary index
	t.primary.insert(nodeID, child)

	return true, false
}

func (t *Tree) Rehead(newHead Node) {
	t.head.AddParent(newHead)
	newHead.Add(t.head)
	t.head = newHead
}

func (t *Tree) PrintData() {
	for _, n := range t.BFS() {
		fmt.Printf("%+v\n", n)
	}
}
