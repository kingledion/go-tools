package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {

	node1 := &node{primary: 1}

	var tests = map[string]struct {
		tree *Tree
		exp  Node
	}{
		"nil root": {
			tree: Empty(),
			exp:  nil,
		},
		"non-nil root": {
			tree: &Tree{
				root: node1,
			},
			exp: node1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.tree.Root()
			assert.Equal(t, tt.exp, got)
		})
	}
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
		t.Run(name, func(t *testing.T) {
			tree := tt.prep()
			gotAdded, gotExists := tree.Add(tt.add.nodeID, tt.add.parentID, "")

			assert.Equal(t, tt.expAdded, gotAdded)
			assert.Equal(t, tt.expExists, gotExists)
		})
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
		t.Run(name, func(t *testing.T) {
			tree := Empty()
			for _, input := range tt.adds {
				tree.Add(input.nodeID, input.parentID, "")
			}

			assert.Equal(t, tt.expBFC, bfc([]Node{tree.root}, []uint64{}))
			assert.Equal(t, tt.expDFC, dfc(tree.root, []uint64{}))
		})

	}
}
