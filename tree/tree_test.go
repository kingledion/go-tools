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
	nodeID   uint
	parentID uint
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
		expBFC []uint
		expDFC []uint
	}{
		"three level parent-child": {
			adds: []addInput{
				{1, 0},
				{2, 1},
				{3, 2},
			},
			expBFC: []uint{1, 2, 3},
			expDFC: []uint{1, 2, 3},
		},
		"three level multi-children": {
			adds: []addInput{
				{1, 0},
				{2, 1},
				{3, 2},
				{4, 1},
			},
			expBFC: []uint{1, 2, 4, 3},
			expDFC: []uint{1, 2, 3, 4},
		},
		"re-root with a new subtree": {
			adds: []addInput{
				{1, 2},
				{3, 1},
				{2, 0},
				{4, 2},
			},
			expBFC: []uint{2, 1, 4, 3},
			expDFC: []uint{2, 1, 3, 4},
		},
		"failed inserts": {
			adds: []addInput{
				{1, 0},
				{2, 1},
				{3, 2},
				{2, 1},
				{4, 5},
			},
			expBFC: []uint{1, 2, 3},
			expDFC: []uint{1, 2, 3},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tree := Empty()
			for _, input := range tt.adds {
				tree.Add(input.nodeID, input.parentID, "")
			}

			assert.Equal(t, tt.expBFC, bfc([]Node{tree.root}, []uint{}))
			assert.Equal(t, tt.expDFC, dfc(tree.root, []uint{}))

			for _, key := range tt.expBFC {
				k := tree.primary.find(key)
				if assert.NotNil(t, k, "Expceted value for %d not to be nil", key) {
					assert.Equal(t, key, k.GetID())
				}
			}
		})

	}
}

func TestFind(t *testing.T) {

	var tests = map[string]struct {
		prep      func() *Tree
		argID     uint
		expNodeID uint
		expOK     bool
	}{
		"primary does not exist": {
			prep: func() *Tree {
				n := &node{primary: 1}
				t := &Tree{root: n, primary: &index{1: n}}
				t.Add(2, 1, "")
				return t
			},
			argID:     3,
			expNodeID: 0,
			expOK:     false,
		},
		"primary exists - branch end": {
			prep: func() *Tree {
				n := &node{primary: 1}
				t := &Tree{root: n, primary: &index{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 1, "")
				return t
			},
			argID:     3,
			expNodeID: 3,
			expOK:     true,
		},
		"primary exists - mid tree": {
			prep: func() *Tree {
				n := &node{primary: 1}
				t := &Tree{root: n, primary: &index{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 1, "")
				return t
			},
			argID:     2,
			expNodeID: 2,
			expOK:     true,
		},
		"primary exists - root": {
			prep: func() *Tree {
				n := &node{primary: 1}
				t := &Tree{root: n, primary: &index{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 1, "")
				return t
			},
			argID:     1,
			expNodeID: 1,
			expOK:     true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tree := tt.prep()
			gotNode, gotOK := tree.Find(tt.argID)

			var gotNodeID uint = 0
			if gotNode != nil {
				gotNodeID = gotNode.GetID()
			}

			assert.Equal(t, tt.expOK, gotOK)
			assert.Equal(t, tt.expNodeID, gotNodeID)
		})
	}
}

func TestFindParents(t *testing.T) {

	var tests = map[string]struct {
		prep       func() *Tree
		argID      uint
		expNodeIDs []uint
		expOK      bool
	}{
		"primary does not exist": {
			prep: func() *Tree {
				n := &node{primary: 1}
				t := &Tree{root: n, primary: &index{1: n}}
				t.Add(2, 1, "")
				return t
			},
			argID:      3,
			expNodeIDs: []uint{},
			expOK:      false,
		},
		"primary exists - branch end": {
			prep: func() *Tree {
				n := &node{primary: 1}
				t := &Tree{root: n, primary: &index{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 1, "")
				return t
			},
			argID:      3,
			expNodeIDs: []uint{2, 1},
			expOK:      true,
		},
		"primary exists - mid tree": {
			prep: func() *Tree {
				n := &node{primary: 1}
				t := &Tree{root: n, primary: &index{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 1, "")
				return t
			},
			argID:      2,
			expNodeIDs: []uint{1},
			expOK:      true,
		},
		"primary exists - root": {
			prep: func() *Tree {
				n := &node{primary: 1}
				t := &Tree{root: n, primary: &index{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 1, "")
				return t
			},
			argID:      1,
			expNodeIDs: []uint{},
			expOK:      true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tree := tt.prep()
			gotNodes, gotOK := tree.FindParents(tt.argID)

			var gotNodeIDs = make([]uint, len(gotNodes))
			if gotNodes != nil {
				for i, n := range gotNodes {
					gotNodeIDs[i] = n.GetID()
				}
			}

			assert.Equal(t, tt.expOK, gotOK)
			assert.Equal(t, tt.expNodeIDs, gotNodeIDs)
		})
	}
}

func TestMerge(t *testing.T) {

	var tests = map[string]struct {
		prepRoot  func() *Tree
		prepOther func() *Tree
		expOK     bool
		expBFC    []uint
		expDFC    []uint
	}{
		"other parent not in tree": {
			prepRoot: func() *Tree {
				n := &node{primary: 1}
				t := &Tree{root: n, primary: &index{1: n}}
				t.Add(2, 1, "")
				return t
			},
			prepOther: func() *Tree {
				n := &node{primary: 3}
				t := &Tree{root: n, primary: &index{3: n}}
				t.Add(4, 3, "")
				return t
			},
			expOK:  false,
			expBFC: []uint{1, 2},
			expDFC: []uint{1, 2},
		},
		"dulicate keys": {
			prepRoot: func() *Tree {
				n := &node{primary: 1}
				t := &Tree{root: n, primary: &index{1: n}}
				t.Add(2, 1, "")
				return t
			},
			prepOther: func() *Tree {
				t := Empty()
				t.Add(3, 1, "")
				t.Add(2, 3, "")
				return t
			},
			expOK:  false,
			expBFC: []uint{1, 2},
			expDFC: []uint{1, 2},
		},
		"merged - branch end": {
			prepRoot: func() *Tree {
				n := &node{primary: 1}
				t := &Tree{root: n, primary: &index{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 2, "")
				t.Add(5, 1, "")
				return t
			},
			prepOther: func() *Tree {
				t := Empty()
				t.Add(6, 5, "")
				t.Add(7, 6, "")
				return t
			},
			expOK:  true,
			expBFC: []uint{1, 2, 5, 3, 4, 6, 7},
			expDFC: []uint{1, 2, 3, 4, 5, 6, 7},
		},
		"merged - mid tree": {
			prepRoot: func() *Tree {
				n := &node{primary: 1}
				t := &Tree{root: n, primary: &index{1: n}}
				t.Add(2, 1, "")
				t.Add(3, 2, "")
				t.Add(4, 2, "")
				t.Add(5, 1, "")
				return t
			},
			prepOther: func() *Tree {
				t := Empty()
				t.Add(6, 1, "")
				t.Add(7, 6, "")
				return t
			},
			expOK:  true,
			expBFC: []uint{1, 2, 5, 6, 3, 4, 7},
			expDFC: []uint{1, 2, 3, 4, 5, 6, 7},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tree := tt.prepRoot()
			other := tt.prepOther()
			gotOK := tree.Merge(other)

			assert.Equal(t, tt.expOK, gotOK)

			assert.Equal(t, tt.expBFC, bfc([]Node{tree.root}, []uint{}))
			assert.Equal(t, tt.expDFC, dfc(tree.root, []uint{}))

			for _, key := range tt.expBFC {
				k := tree.primary.find(key)
				if assert.NotNil(t, k, "Expceted value for %d not to be nil", key) {
					assert.Equal(t, key, k.GetID())
				}
			}
		})
	}
}
