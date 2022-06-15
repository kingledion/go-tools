package queue

import "sync"

type Queue[T any] struct {
	q        chan T
	capacity int
	pl       *sync.Mutex
}

func (q *Queue[T]) Empty() *Queue[T] {

	return &Queue[T]{
		q:        make(chan T, 2),
		capacity: 2,
		pl:       &sync.Mutex{},
	}
}

func (q *Queue[T]) push(elem T) {
	q.pl.Lock()
	select {
	case q.q <- elem:
	default:
		// resize the channel
		q.capacity = q.capacity * 2
		oldChan := q.q
		q.q = make(chan T, q.capacity)
		for v := range oldChan {
			q.q <- v
		}

	}
	q.pl.Unlock()
}

func (q *Queue[T]) pop() T {
	return <-q.q
}
