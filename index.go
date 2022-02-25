package tree

import "log"

type index map[uint64]Node

func (idx *index) find(id uint64) Node {
	if idx == nil { // do we need an error check here?
		log.Println("Attempting to find in an undefined index")
		return nil
	}
	m := *idx
	val, exists := m[id]
	if !exists {
		return nil
	}
	return val
}

func (idx *index) insert(id uint64, node Node) bool {
	if idx == nil { // do we need an error check here?
		log.Println("Attempting to insert in an undefined index")
		return false
	}
	m := *idx
	m[id] = node
	return true
}
