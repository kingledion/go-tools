package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChildren(t *testing.T) {

	node1 := &node{primary: 1}
	node2 := &node{primary: 2}

	tests := map[string]struct {
		n        Node
		expChild []Node
	}{
		"nil child array": {
			n:        &node{},
			expChild: nil,
		},
		"empty child array": {
			n:        &node{children: []Node{}},
			expChild: []Node{},
		},
		"success": {
			n:        &node{children: []Node{node1, node2}},
			expChild: []Node{node1, node2},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotChild := tt.n.GetChildren()

			assert.Equal(t, tt.expChild, gotChild)
		})
	}
}

func TestGetID(t *testing.T) {

	tests := map[string]struct {
		n     Node
		expID uint64
	}{
		"success": {
			n:     &node{primary: 1},
			expID: 1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotID := tt.n.GetID()

			assert.Equal(t, tt.expID, gotID)
		})
	}
}

func TestNodeAdd(t *testing.T) {

	node1 := &node{primary: 1}
	node2 := &node{primary: 2}
	node3 := &node{primary: 3}

	tests := map[string]struct {
		n        Node
		argNode  []Node
		expChild []Node
	}{
		"nil child array": {
			n:        &node{},
			argNode:  []Node{node1},
			expChild: []Node{node1},
		},
		"empty child array": {
			n:        &node{children: []Node{}},
			argNode:  []Node{node1},
			expChild: []Node{node1},
		},
		"add none": {
			n:        &node{children: []Node{node1}},
			argNode:  nil,
			expChild: []Node{node1},
		},
		"non-empty child array": {
			n:        &node{children: []Node{node1}},
			argNode:  []Node{node2, node3},
			expChild: []Node{node1, node2, node3},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.n.Add(tt.argNode...)

			assert.Equal(t, tt.expChild, tt.n.GetChildren())
		})
	}
}

func TestIsParent(t *testing.T) {

	tests := map[string]struct {
		n       Node
		argID   uint
		expBool bool
	}{
		"true": {
			n:       &node{parentID: 1},
			argID:   1,
			expBool: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotBool := tt.n.IsParent(tt.argID)

			assert.Equal(t, tt.expBool, gotBool)
		})
	}
}

func TestGetParent(t *testing.T) {

	node1 := &node{primary: 1}

	tests := map[string]struct {
		n       Node
		expNode Node
	}{
		"nil parent": {
			n:       &node{},
			expNode: nil,
		},
		"success": {
			n:       &node{parent: node1},
			expNode: node1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotNode := tt.n.GetParent()

			assert.Equal(t, tt.expNode, gotNode)
		})
	}
}
