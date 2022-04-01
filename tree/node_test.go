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
		expID uint
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

func TestNodeAddChildren(t *testing.T) {

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
		"add nil": {
			n:        &node{children: []Node{node1}},
			argNode:  nil,
			expChild: []Node{node1},
		},
		"add empty array": {
			n:        &node{children: []Node{node1}},
			argNode:  []Node{},
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
			tt.n.AddChildren(tt.argNode...)

			assert.Equal(t, tt.expChild, tt.n.GetChildren())
		})
	}
}

func TestGetParentID(t *testing.T) {

	tests := map[string]struct {
		n   Node
		exp uint
	}{
		"trivial": {
			n:   &node{parentID: 1},
			exp: 1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.n.GetParentID()

			assert.Equal(t, tt.exp, got)
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
