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
