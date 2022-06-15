package tree

// Minimal implementaiton of a breadth and depth first search for testing
// if a tree matches what is expected. Yields primary keys of all nodes
// in a tree in order of selected traversal.

// returns the primary keys of a tree in order of a breadth first search
func bfc(q []Node, iter []uint) []uint {
	iter = append(iter, q[0].GetID())
	q = append(q[1:], q[0].GetChildren()...)
	if len(q) == 0 {
		return iter
	}
	return bfc(q, iter)
}

// returns the primary keys of a tree in order of a breadth first search
func dfc(n Node, iter []uint) []uint {
	iter = append(iter, n.GetID())
	for _, c := range n.GetChildren() {
		iter = dfc(c, iter)
	}
	return iter
}
