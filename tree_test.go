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

type addInput struct {
	nodeID   uint64
	parentID uint64
}

func TestAdd(t *testing.T) {

	var tests = map[string]struct {
		prep      func() *Tree
		add       addInput
		expAdded  bool
		expExists bool
	}{
		"primary exists": {
			prep: func() *Tree {
				n := &node{primary: 1}
				return &Tree{root: n, primary: &index{1: n}}
			},
			add:       addInput{1, 0},
			expAdded:  false,
			expExists: true,
		},
		"root is nil": {
			prep: func() *Tree {
				return &Tree{primary: &index{}}
			},
			add:       addInput{1, 0},
			expAdded:  true,
			expExists: false,
		},
		"re-root": {
			prep: func() *Tree {
				n := &node{primary: 1, parentID: 2}
				return &Tree{root: n, primary: &index{1: n}}
			},
			add:       addInput{2, 3},
			expAdded:  true,
			expExists: false,
		},
		"re-root with cycle": {
			prep: func() *Tree {
				n := &node{primary: 1, parentID: 2}
				return &Tree{root: n, primary: &index{1: n}}
			},
			add:       addInput{2, 1},
			expAdded:  false,
			expExists: false,
		},
		"parent does not exist": {
			prep: func() *Tree {
				n := &node{primary: 1}
				return &Tree{root: n, primary: &index{1: n}}
			},
			add:       addInput{2, 3},
			expAdded:  false,
			expExists: false,
		},
		"added": {
			prep: func() *Tree {
				n := &node{primary: 1}
				return &Tree{root: n, primary: &index{1: n}}
			},
			add:       addInput{2, 1},
			expAdded:  true,
			expExists: false,
		},
	}

	for name, tt := range tests {
		tree := tt.prep()
		gotAdded, gotExists := tree.Add(tt.add.nodeID, tt.add.parentID, "")

		assert.Equal(t, tt.expAdded, gotAdded, "%s: bool added return value does not match", name)
		assert.Equal(t, tt.expExists, gotExists, "%s: bool exists return value does not match", name)

	}

}

func TestAddResults(t *testing.T) {

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
		"three level multi-children": {
			adds: []addInput{
				{1, 0},
				{2, 1},
				{3, 2},
				{4, 1},
			},
			expBFC: []uint64{1, 2, 4, 3},
			expDFC: []uint64{1, 2, 3, 4},
		},
		"re-root with a new subtree": {
			adds: []addInput{
				{1, 2},
				{3, 1},
				{2, 0},
				{4, 2},
			},
			expBFC: []uint64{2, 1, 4, 3},
			expDFC: []uint64{2, 1, 3, 4},
		},
		"failed inserts": {
			adds: []addInput{
				{1, 0},
				{2, 1},
				{3, 2},
				{2, 1},
				{4, 5},
			},
			expBFC: []uint64{1, 2, 3},
			expDFC: []uint64{1, 2, 3},
		},
	}

	for name, tt := range tests {
		tree := Empty()
		for _, input := range tt.adds {
			tree.Add(input.nodeID, input.parentID, "")
		}

		assert.Equal(t, tt.expBFC, bfc([]Node{tree.root}, []uint64{}), "%s: breadth-first search does not match", name)
		assert.Equal(t, tt.expDFC, dfc(tree.root, []uint64{}), "%s: depth-first search does not match", name)

	}
}
