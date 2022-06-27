package tree

import (
	"github.com/phf/go-queue/queue"
)

// TraversalType determines the order in which an operation is performed on a tree.
type TraversalType int

const (
	// TraverseBreadthFirst traverses the tree breadth first. For any node
	// currently being traversed, all children of this node are traversed
	// before any of the children's children are visited.
	TraverseBreadthFirst TraversalType = iota
	// TraverseDepthFrist traverses the tree depth first. For any node
	// currently being traversed, all descendents of any child will be
	// traversed before any subsequent children of the current node are
	// visited.
	TraverseDepthFirst
)

// Traverse visits each node of a tree in a specified order, returning
// those nodes to an iterator-like chennel.
//
// This function takes as an argumentt a TraversalType which defines the
// order of traversal. All node in the tree are traversed in this order.
// The Nodes traversed are pushed to an unbuffered channel and must be
// consumed by caller.
//
// If a tree is modified after the traversal has begun, any node that is
// added after its correct place in traversal order will not be visited, nor
// will any of its children.
func (t *Tree[T]) Traverse(trvsl TraversalType) <-chan Node[T] {
	search := make(chan Node[T])

	switch trvsl {
	case TraverseBreadthFirst:
		q := queue.New()
		q.PushBack(t.root)
		go func() {
			for {
				if bfs(q, search) {
					close(search)
					break
				}
			}
		}()
	}

	return search

}

func bfs[T any](q *queue.Queue, search chan<- Node[T]) bool {

	current := q.PopFront()
	switch c := current.(type) {
	case Node[T]:
		for _, c := range c.GetChildren() {
			q.PushBack(c)
		}
		search <- c
		return false
	case nil:
		return true
	default:
		// Should be unreachable...
		return true
	}

}
