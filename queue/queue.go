package queue

// type Queue struct {
// 	internal []interface{}
// 	in       int // points to empty spot for next push
// 	out      int // points to element for next pop
// }

// func New() *Queue {
// 	return &Queue{internal: make([]interface{}, 1)}
// }

// func (q *Queue) upsize() {

// }

// func (q *Queue) downsize() {

// }

// func (q *Queue) dec(i int) int {
// 	return (i - 1) & (len(q.internal) - 1)
// }

// // Push
// func (q *Queue) Push(e interface{}) {

// 	q.upsize()
// 	q.in = q.dec(q.in)
// 	q.internal[q.in] = e

// }

// func (q *Queue) Pop() interface{} {
// 	o := q.internal[q.out]
// 	if o == nil {
// 		return nil
// 	}
// 	q.internal[q.out] = nil
// 	q.out = q.dec(q.out)
// 	q.downsize()
// 	return o
// }
