package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBFS(t *testing.T) {

	tests := map[string]struct {
		tree      func() *Tree
		traversal TraversalType
		expSearch []uint
	}{
		"success": {
			tree: func() *Tree {
				node6 := &node{Primary: 6}
				node5 := &node{Primary: 5}
				node4 := &node{Primary: 4}
				node3 := &node{Primary: 3, children: []Node{node4, node5}}
				node2 := &node{Primary: 2, children: []Node{node6}}
				node1 := &node{Primary: 1, children: []Node{node2, node3}}
				return &Tree{root: node1}
			},
			traversal: TraverseBreadthFirst,
			expSearch: []uint{1, 2, 3, 6, 4, 5},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			i := 0
			for g := range tt.tree().Traverse(tt.traversal) {
				assert.Equal(t, tt.expSearch[i], g.GetID())
				i = i + 1
			}

		})
	}
}
