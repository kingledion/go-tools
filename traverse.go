package tree

// BFSArray is an array of Nodes ordered as a breadth first search
type BFSArray []Node

// BFS creates an array of all nodes in the tree in order of a breadth first
// search.
func (t *Tree) BFS() BFSArray {
	iter := BFSArray{}
	q := []Node{t.root}
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
