package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {

	exp := node{}
	tree := Tree{
		root: &exp,
	}

	got := tree.Root()

	assert.Equal(t, &exp, got, "root pointers should be equal")

}

func TestAdd(t *testing.T) {

	type addInput struct {
		nodeID   uint64
		parentID uint64
	}

	var tests = map[string]struct {
		adds   []addInput
		expBFC []uint64
		expDFC []uint64
	}{
		"three level parent-child": {
			adds: []addInput{
				{1, 0},
				{2, 1},
				{3, 2},
			},
			expBFC: []uint64{1, 2, 3},
			expDFC: []uint64{1, 2, 3},
		},
	}

	for _, tt := range tests {
		tree := Empty()
		for _, input := range tt.adds {
			tree.Add(input.nodeID, input.parentID, "")
		}

		assert.Equal(t, tt.expBFC, bfc([]Node{tree.root}, []uint64{}), "breadth-first search should match expected")
		assert.Equal(t, tt.expDFC, dfc(tree.root, []uint64{}), "depth-first search should match expected")

	}
}
