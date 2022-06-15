package tree

import (
	"github.com/phf/go-queue/queue"
)

// BFS implements a breadth first search through a tree. For each node reached
// in the traversal, all child nodes are added to a queue of future nodes to search.
// Thus all children of the current node will be serch before any children of
// those children are searched.
func (t *Tree) BFS() <-chan Node {
	search := make(chan Node)
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

	return search
}

func bfs(q *queue.Queue, search chan<- Node) bool {

	current := q.PopFront()
	switch c := current.(type) {
	case Node:
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
