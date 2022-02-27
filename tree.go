/*
Package tree implements a simple tree that can be built from and stored to a
row based data format, such as a relational database or csv file.

A tree is here defined as a graph having three properties:
  - a single root node with no inbound edges
  - all non-root nodes have exactly one inbound edge (parent)
  - any node may have any number of outbound edges (children)
The nodes of the tree are assumed to have a primary identifier or key by
which parent and child relationships can be defined.

The implemented tree data structure consists of a pointer to the root node
and an index allowing fast lookup of nodes by their primary key. Each node
holds the primary key of itself and its parent, as well as an array of
pointers to its children. Children of a node are not ordered, and therefore
are traversed in an order determined by the order in which they are added
to the tree.

This package includes a breadth first search algorithm.
*/
package tree

// Tree is a data structure representing a tree. It contains a pointer to
// a root node and an index of primary keys.
type Tree struct {
	root    Node
	primary *index
}

// Empty creates and returns an empty tree. The empty tree has a nil pointer
// to its root node.
func Empty() *Tree {
	return &Tree{
		primary: &index{},
	}
}

// Root returns the root node of a tree. If the tree is empty, it returns
// nil.
func (t *Tree) Root() Node {
	return t.root
}

// Add inserts an element into a tree as a node. This function returns two
// boolean values:
//   - added - indicates that the element was successfully inserted into
//     the tree
//   - exists - indicates that the element's primary was already found
//     in the tree
//
// There are two ways that insertion can fail. If the element's parent is not
// found, then the element is not inserted; added and exists will both be
// false. If the element's primary key already exists in the index, then
// added will be false and exists will be true.
// If the element is added as expected, then added will be true and exists
// will be false.
//
// If the element to be added has a primary key that matches the parent key
// of the root node, the tree will be re-rooted by adding this element as the
// new root. If there is a cyclical reference when attempting to re-root, i.e. the
// parent of the existing root is the new node and the parent of the new node is
// the exiting rool, the element will fail to add.
func (t *Tree) Add(nodeID uint64, parentID uint64, data interface{}) (added bool, exists bool) {
	child := &node{primary: nodeID, parentID: parentID, data: data}

	// Return false if this element has already been added
	if t.primary.find(nodeID) != nil {
		exists = true
		return
	}

	if t.root == nil { // always insert the first element
		t.root = child
	} else {

		parent := t.primary.find(parentID)
		if parent == nil {
			if t.root.IsParent(nodeID) { // parent does not exist but incoming node is parent of root
				t.reroot(child)
			} else { // parent does not exist, do not add
				return
			}
		} else {
			if t.root.IsParent(nodeID) { // parent exists, but incoming node causes cycle
				return
			}
			// parent exists, add
			child.addParent(parent)
			parent.Add(child)
		}
	}

	// add to primary index
	t.primary.insert(nodeID, child)

	added = true
	return
}

func (t *Tree) reroot(newHead Node) {
	t.root.addParent(newHead)
	newHead.Add(t.root)
	t.root = newHead
}
