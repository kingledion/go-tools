/*
Package tree implements a simple tree that can be serialized and deserialized
from a byte-based storage, such as a file or cloud storage.

A tree is here defined as a graph having three properties:
  - a single root node with no inbound edges
  - all non-root nodes have exactly one inbound edge (parent)
  - any node may have any number of outbound edges (children)
The nodes of the tree are assumed to have a primary identifier or key by
which parent and child relationships can be defined.

The implemented tree data structure consists of a pointer to the root node
and an index allowing fast lookup of nodes by their primary key. Each node
holds the primary key of itself and its parent, as well as a pointer to its
parent and an array of pointers to its children. All nodes also store an
arbitrary set of node data, which can be any structure. Children of a node are
not ordered, and therefore are traversed in an order determined by the order
in which they are added to the tree.

This package includes tree traversal algorithms for breadth-first and depth-
first search.
*/
package tree

import (
	"encoding/json"
	"fmt"
	"io"
)

// Tree is a data structure representing a tree. It contains a pointer to
// a root node and an index of primary keys implemented as a hash map.
type Tree struct {
	root    Node
	primary *index
}

// Empty creates and returns an empty tree. The empty tree has a nil pointer
// to its root node and an empty node index.
func Empty() *Tree {
	return &Tree{
		primary: &index{},
	}
}

// Root returns the root node of a tree. If the tree has no nodes, this
// function returns nil.
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
func (t *Tree) Add(nodeID uint, parentID uint, data any) (added bool, exists bool) {
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
			child.setParent(parent)
			parent.AddChildren(child)
		}
	}

	// add to primary index
	t.primary.insert(nodeID, child)

	added = true
	return
}

func (t *Tree) reroot(newHead Node) {
	t.root.setParent(newHead)
	newHead.AddChildren(t.root)
	t.root = newHead
}

// Merge another tree (passed in the argument) into the target tree (passed as the
// subject of this method call). All data from the other tree is added to the target
// tree, if a relationship can be found between the two trees. A relationship is
// established if and only if the parent of the head of the other tree is found
// in the target tree.
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
		other.root.setParent(f)

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
// The argument TraversalType will determine the traversal order in which
// the tree is serialized. TraversalType does not matter for deserialization;
// the internal metadata of the nodes will create the shape of the tree when
// it is deserialized, not the order in which the nodes are serialized
// to storage.
//
// The associated data of each node is serialized with it. This data may be
// set the the caller and may not be serializable. If the associated data
// cannot be serialzied using the json package, then this function will
// throw an error. Once any error is thrown, serialization stops.
//
// The serialization is implemented into a goroutine which will populate the
// ReadCloser return value as elements are consumed from it by the caller.
// the <-chan error exists to pass any serialization error back from the
// encoding goroutine.
func (t *Tree) Serialize(trvsl TraversalType) (io.ReadCloser, <-chan error) {
	reader, writer := io.Pipe()
	errchan := make(chan error)

	go func() {
		encoder := json.NewEncoder(writer)
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

// Deserialize decodes a data stream into a tree.
//
// Decode is validated for data streams encoded via the [`Serialize`]
// method on [`Tree`]. While encoding is implemented with the json packagae,
// There is no guarantee that it will deserialize data encoded in any other way.
//
// The argument ReadCloser is a stream with data from a serialized tree. If any
// node of the tree fails to deserialize, this function will abord and return an
// error.
func Deserialize(stream io.ReadCloser) (*Tree, error) {
	decoder := json.NewDecoder(stream)
	var n node
	var t *Tree = Empty()
	for {

		err := decoder.Decode(&n)
		if err == io.EOF {
			//log.Printf("deserialize - end of file")
			return t, nil
		}

		if err != nil {
			//log.Printf("deserialize - error: %s", err)
			return nil, fmt.Errorf("error deserializing: %w", err)
		}

		//log.Printf("deserialize - adding %d %d %+v", n.GetID(), n.GetParentID(), n.GetData())
		t.Add(n.GetID(), n.GetParentID(), n.GetData())

	}

}
