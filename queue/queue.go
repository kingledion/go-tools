package queue

import "sync"

type queue struct {
	q         chan interface{}
	capacity  int
	expanding *sync.Mutex
}

func new() *queue {

	return &queue{
		q:         make(chan interface{}, 2),
		capacity:  2,
		expanding: &sync.Mutex{},
	}
}

func (q *queue) push(elem interface{}) {
	select {
	case q.q <- elem:
	default:
		// resize the channel
		q.capacity = q.capacity * 2
		oldChan := q.q
		q.q = make(chan interface{}, q.capacity)
		for v := range *oldChan {
			q.q <- v
		}

	}
}

func (q *queue) pop() interface{} {
	return <-q.q
}
