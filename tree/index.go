package tree

import (
	"log"
	"reflect"
)

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

func (idx *index[T]) findByVal(val T) []Node[T]{
	if idx == nil { // do we need an error check here?
		log.Println("Attempting to insert in an undefined index")
		return nil
	}

	// The reflection could be removed, by adding a type constraint, but that
	// would not be perfect for this labrary.
	valData := reflect.ValueOf(val)
	if !valData.Comparable(){
		log.Println("The value type is incomparable")
	}

	// Maybe there is a better way to search for an object in an unsorted list,
	// but i have no idea how as of right now. Besides O(n) is not that bad.
	var nl []Node[T]
	for _, node := range *idx{
		nodeData := reflect.ValueOf(node.GetData())
		if valData.Equal(nodeData){
			nl = append(nl, node)
		}
	}	
	if len(nl) == 0{
		log.Println("Not found a node with given value")
		return nil
	}
	return nl
}
