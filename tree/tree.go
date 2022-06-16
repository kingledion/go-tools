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

import (
	"encoding/gob"
	"fmt"
	"io"
)

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
// parent of the existing root is the new node and the parent of the new node
// exists elsewhere in the tree, the element will fail to add.
//
// Do not set a primaryID to zero, as this value should be reserved for the
// case where a node has no parent.
func (t *Tree) Add(nodeID uint, parentID uint, data interface{}) (added bool, exists bool) {
	child := &node{Primary: nodeID, ParentID: parentID, Data: data}

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
			if t.root.GetParentID() == nodeID { // parent does not exist but incoming node is parent of root
				t.reroot(child)
			} else { // parent does not exist, do not add
				return
			}
		} else {
			if t.root.GetParentID() == nodeID { // parent exists, but incoming node causes cycle
				return
			}
			// parent exists, add
			child.SetParent(parent)
			parent.AddChildren(child)
		}
	}

	// add to primary index
	t.primary.insert(nodeID, child)

	added = true
	return
}

func (t *Tree) reroot(newHead Node) {
	t.root.SetParent(newHead)
	newHead.AddChildren(t.root)
	t.root = newHead
}

// Merge another tree (passed in the argument) into the target tree (passed as the
// subject of this method call). All data from the other tree is added to the target
// tree.
//
// If the merge is successful, returns true, otherwise return false. The merge can
// fail if there are duplicate primary keys between the two trees. The merge
// can also fail if the parent of the head of the other tree is not found in the
// target tree.
func (t *Tree) Merge(other *Tree) bool {

	if other == nil {
		return false
	}

	headParent := other.root.GetParentID()

	f := t.primary.find(headParent)
	if f != nil {

		// check for duplicate primary ids
		for k := range *other.primary {
			if t.primary.find(k) != nil {
				return false
			}
		}

		f.AddChildren(other.root)
		other.root.SetParent(f)

		// copy other index to new tree
		for k, n := range *other.primary {
			t.primary.insert(k, n)
		}
		return true
	}

	return false

}

// Find looks up a node by its primary key. If the node is found, then
// ok is true and a Node is returned. If the node is not found, then
// ok is false an a nil pointer is returned.
func (t *Tree) Find(id uint) (n Node, ok bool) {
	f := t.primary.find(id)
	if f == nil {
		return
	}
	return f, true
}

// FindParents finds the list of all parent nodes between a target node and the
// root of a tree. The node is identified by its primary key. If the primary
// key cannot be found in the tree, then ok is false and an empty array is returned.
// If the target node is found in the tree, then ok is true and the parent nodes
// are returned in an array. If the target node is the root of the tree, the
// parent nodes array is empty.
//
// The parent nodes array is ordered from immediate parent first to tree root
// last.
func (t *Tree) FindParents(id uint) (parents []Node, ok bool) {

	f := t.primary.find(id)
	if f == nil {
		return
	}

	for n := f.GetParent(); n != nil; n = n.GetParent() {
		parents = append(parents, n)
	}

	return parents, true
}

// Serialize encodes the tree as a byte stream.
//
// The argument io.Writer will recieve the serialized bytes. The argument
// TraversalType will determine the traversal order in which the tree is
// serialized. TraversalType does not have to be remembered for
// deserialization; the internal metadata of the nodes directs how the
// tree is rebuilt, not the order in which the nodes are serialized
// to storage.
func (t *Tree) Serialize(trvsl TraversalType) (io.ReadCloser, <-chan error) {
	reader, writer := io.Pipe()
	errchan := make(chan error)

	go func() {
		encoder := gob.NewEncoder(writer)
		for n := range t.Traverse(trvsl) {
			err := encoder.Encode(n)
			if err != nil {
				errchan <- err
				writer.Close()
				return
			}

		}

		close(errchan)
		writer.Close()
	}()

	return reader, errchan
}

// Deserialize decodes a data source into a tree.
//
// Elements from the data store are read one by one as decoded gobs.
// These elements are resolved into nodes, and then inserted into the tree.
func Deserialize(stream io.ReadCloser) (*Tree, error) {
	decoder := gob.NewDecoder(stream)
	var n Node
	var t *Tree = Empty()
	for {

		err := decoder.Decode(n)
		if err == io.EOF {
			return t, nil
		}

		if err != nil {
			return nil, fmt.Errorf("error deserializing: %w", err)
		}

		t.Add(n.GetID(), n.GetParentID(), n.GetData())

	}

}
