package queue

type queue struct {
	q chan interface{}
}

func new(capacity int) *queue {
	return &queue{
		q: make(chan interface{}, capacity),
	}
}

func (q *queue) push(elem interface{}) {
	q.q <- elem
}

func (q *queue) pop() interface{} {
	return <-q.q
}
