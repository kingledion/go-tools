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

	assert.Equal(t, &exp, got, "Root pointers should be equal")

}
