package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBFS(t *testing.T) {

	tests := map[string]struct {
		tree      func() *Tree
		expSearch []uint
	}{
		"success": {
			tree: func() *Tree {
				node6 := &node{primary: 6}
				node5 := &node{primary: 5}
				node4 := &node{primary: 4}
				node3 := &node{primary: 3, children: []Node{node4, node5}}
				node2 := &node{primary: 2, children: []Node{node6}}
				node1 := &node{primary: 1, children: []Node{node2, node3}}
				return &Tree{root: node1}
			},
			expSearch: []uint{1, 2, 3, 6, 4, 5},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			i := 0
			for g := range tt.tree().BFS() {
				assert.Equal(t, tt.expSearch[i], g.GetID())
				i = i + 1
			}

		})
	}
}
