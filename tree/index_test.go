package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexFind(t *testing.T) {

	node1 := &node{Primary: 1}
	node2 := &node{Primary: 2}

	tests := map[string]struct {
		index   index
		argID   uint
		expNode Node
	}{
		"nil index": {
			index:   nil,
			argID:   1,
			expNode: nil,
		},
		"not in index": {
			index: index{
				1: node1,
				2: node2,
			},
			argID:   3,
			expNode: nil,
		},
		"success": {
			index: index{
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
