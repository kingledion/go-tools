package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetID(t *testing.T) {

	tests := map[string]struct {
		n     Node
		expID uint
	}{
		"trivial": {
			n:     &node{Primary: 1},
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

func TestGetParentID(t *testing.T) {

	tests := map[string]struct {
		n   Node
		exp uint
	}{
		"trivial": {
			n:   &node{ParentID: 1},
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

func TestGetChildren(t *testing.T) {

	node1 := &node{Primary: 1}
	node2 := &node{Primary: 2}

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

func TestGetParent(t *testing.T) {

	node1 := &node{Primary: 1}

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

func TestNodeAddChildren(t *testing.T) {

	node1 := &node{Primary: 1}
	node2 := &node{Primary: 2}
	node3 := &node{Primary: 3}

	tests := map[string]struct {
		n        Node
		argNodes []Node
		expChild []Node
	}{
		"nil child array": {
			n:        &node{},
			argNodes: []Node{node1},
			expChild: []Node{node1},
		},
		"empty child array": {
			n:        &node{children: []Node{}},
			argNodes: []Node{node1},
			expChild: []Node{node1},
		},
		"add nil": {
			n:        &node{children: []Node{node1}},
			argNodes: nil,
			expChild: []Node{node1},
		},
		"add empty array": {
			n:        &node{children: []Node{node1}},
			argNodes: []Node{},
			expChild: []Node{node1},
		},
		"non-empty child array": {
			n:        &node{children: []Node{node1}},
			argNodes: []Node{node2, node3},
			expChild: []Node{node1, node2, node3},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.n.AddChildren(tt.argNodes...)

			assert.Equal(t, tt.expChild, tt.n.GetChildren())
		})
	}
}

func TestReplaceChildren(t *testing.T) {

	node1 := &node{Primary: 1}
	node2 := &node{Primary: 2}
	node3 := &node{Primary: 3}

	tests := map[string]struct {
		n        Node
		argNodes []Node
		expChild []Node
	}{
		"nil child array": {
			n:        &node{},
			argNodes: []Node{node1},
			expChild: []Node{node1},
		},
		"empty child array": {
			n:        &node{children: []Node{}},
			argNodes: []Node{node1},
			expChild: []Node{node1},
		},
		"use nil": {
			n:        &node{children: []Node{node1}},
			argNodes: nil,
			expChild: []Node{},
		},
		"use empty array": {
			n:        &node{children: []Node{node1}},
			argNodes: []Node{},
			expChild: []Node{},
		},
		"non-empty replacement array": {
			n:        &node{children: []Node{node1}},
			argNodes: []Node{node2, node3},
			expChild: []Node{node2, node3},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.n.ReplaceChildren(tt.argNodes...)

			assert.Equal(t, tt.expChild, tt.n.GetChildren())
		})
	}

}

func TestSetParent(t *testing.T) {

	node1 := &node{Primary: 1}
	node2 := &node{Primary: 2}

	tests := map[string]struct {
		n           Node
		argParent   Node
		expParent   Node
		expParentID uint
	}{
		"set nil parent": {
			n:           &node{Primary: 1},
			argParent:   nil,
			expParent:   nil,
			expParentID: 0,
		},
		"set circular ref parent": {
			n:           node1,
			argParent:   node1,
			expParent:   nil,
			expParentID: 0,
		},
		"success": {
			n:           &node{Primary: 1},
			argParent:   node2,
			expParent:   node2,
			expParentID: 2,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			n := tt.n
			n.setParent(tt.argParent)

			assert.Equal(t, tt.expParent, n.GetParent())
			assert.Equal(t, tt.expParentID, n.GetParentID())
		})
	}

}
