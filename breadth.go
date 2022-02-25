package tree

// An interator implementing a breadth-first search
type BFSArray []Node

func (t *Tree) BFS() BFSArray {
	iter := BFSArray{}
	q := []Node{t.head}
	return bfs(q, iter)

}

func bfs(q []Node, iter BFSArray) BFSArray {

	if len(q) == 0 {
		return iter
	}
	current := q[0]
	iter = append(iter, current)
	q = append(q[1:], current.GetChildren()...)
	return bfs(q, iter)
}
