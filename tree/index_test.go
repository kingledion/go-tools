package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexFind(t *testing.T) {

	node1 := &node[int]{primary: 1}
	node2 := &node[int]{primary: 2}

	tests := map[string]struct {
		index   index[int]
		argID   uint
		expNode Node[int]
	}{
		"nil index": {
			index:   nil,
			argID:   1,
			expNode: nil,
		},
		"not in index": {
			index: index[int]{
				1: node1,
				2: node2,
			},
			argID:   3,
			expNode: nil,
		},
		"success": {
			index: index[int]{
				1: node1,
				2: node2,
			},
			argID:   2,
			expNode: node2,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			gotNode := tt.index.find(tt.argID)
			assert.Equal(t, tt.expNode, gotNode)
		})
	}
}


func TestIndexFindByValue(t *testing.T){
	tree := Empty[int]()
	tree.Add(1, 0, 1)  //   1   
	tree.Add(2, 1, 2)  //  / \  // visual representation
	tree.Add(3, 1, 3)  // 2   3
	tree.Add(4, 3, 3)  //      \
	                   //       3
	// node1, _ := tree.Find(1)
	node2, _ := tree.Find(2)
	node3, _ := tree.Find(3)
	node4, _ := tree.Find(4)
	
	tests1 := map[string]struct{
		tr *Tree[int]
		val int
		expNode Node[int]
		okVal bool
	}{
		"empty tree": {
			tr: Empty[int](),
			val: 1,
			expNode: nil,
			okVal: false,
		},
		"not found": {
			tr: tree,
			val: 5,
			expNode: nil,
			okVal: false,
		},
		"1 nodes": {
			tr: tree,
			val: 2,
			expNode: node2,
			okVal: true,
		},
		"2 and more nodes": {
			tr: tree,
			val: 3,
			expNode: node3,
			okVal: false,
		},
	}

	tests2 := map[string]struct{
		tr *Tree[int]
		val int
		expNodes []Node[int]
		okVal bool
	}{
		
		"empty tree": {
			tr: Empty[int](),
			val: 1,
			expNodes: nil,
			okVal: false,
		},
		"not found": {
			tr: tree,
			val: 5,
			expNodes: nil,
			okVal: false,
		},
		"1 nodes": {
			tr: tree,
			val: 2,
			expNodes: []Node[int]{node2},
			okVal: true,
		},
		"2 and more nodes": {
			tr: tree,
			val: 3,
			expNodes: []Node[int]{node3, node4},
			okVal: false,
		},
	}

	for name, tt := range tests1 {
		t.Run(name, func(t *testing.T) {
			gotNode, _ := tt.tr.FindByValue(tt.val)
			assert.Equal(t, tt.expNode, gotNode)
		})
	}


	for name, tt := range tests2 {
		t.Run(name, func(t *testing.T) {
			gotNodes, _ := tt.tr.FindMultipleByValue(tt.val)
			assert.Equal(t, tt.expNodes, []Node[int](gotNodes))
		})
	}
}
