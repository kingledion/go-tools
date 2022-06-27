package tree

import "log"

type index[T any] map[uint]Node[T]

func (idx *index[T]) find(id uint) Node[T] {
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

func (idx *index[T]) insert(id uint, node Node[T]) bool {
	if idx == nil { // do we need an error check here?
		log.Println("Attempting to insert in an undefined index")
		return false
	}
	m := *idx
	m[id] = node
	return true
}
