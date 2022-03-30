package retry

import (
	"fmt"
	"sync"
	"time"
)

// Specifies an element of the retry queue. The element must have methods that process
// and record failure to process. The processing method will return a boolean representing
// whether processing should be tried again.
type Element interface {
	// Process this element. Return true if the process shoud be tried again.
	Process() bool
	// Record failure to process.
	Fail()
}

// Retry is a queue with an automatic retry capability. Elements added to the queue are
// processed and retried until the input is closed and
type Retry struct {
	input chan Element
	retry chan Element
	wait  sync.WaitGroup
}

// Build an empty retry queue
func New(capacity uint) *Retry {
	r := &Retry{
		input: make(chan Element),
		retry: make(chan Element, capacity),
	}
	r.run()
	return r
}

// Add a new element to the retry queue
func (r *Retry) Push(e Element) {
	r.input <- e
}

// Close channel, indicating that no further input will be passed
func (r *Retry) Close() {
	close(r.input)
}

func (r *Retry) Wait() {
	r.wait.Wait()
}

// run starts the processing of elements
func (r *Retry) run() {

	inputTimeout := make(chan bool, 1)

	// input goroutine
	go func() {
		for elem := range r.input {
			tryagain := elem.Process()
			if tryagain {
				r.retry <- elem
			}
		}
		inputTimeout <- true
	}()

	// retry goroutine
	go func() {

		r.wait.Add(1)

		cntr := 0
		rlen := 0
		inputDone := false
		retryDone := false

		for elem := range r.retry {

			// if done, process all remaining elements
			if retryDone {
				elem.Fail()
				continue
			}

			// process current element
			tryagain := elem.Process()
			if tryagain {
				r.retry <- elem
			}

			// if the input channel is closed, start monitoring for shutdown conditions
			select {
			case <-inputTimeout:
				inputDone = true
				rlen = len(r.retry)
			case <-time.After(1 * time.Second):
				continue
			default:
			}

			// if input is done, loop through all elements once then shutdown
			if inputDone {
				fmt.Printf("Input closed, loop %d\n", cntr)
				cntr = cntr + 1

				if cntr > rlen {
					retryDone = true
					close(r.retry)
				}

			}

		}

		r.wait.Done()

	}()

}
