package retry

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type example struct {
	i int
}

func (e *example) Process() bool {
	if rand.Intn(3) < 1 {
		fmt.Println("Fail on", e.i)
		return true
	}
	fmt.Println("Success on", e.i)
	return false
}

func (e *example) Fail() {
	fmt.Println("Abandoning", e.i)
}

func TestRetry(t *testing.T) {

	fmt.Println("Starting test")

	rand.Seed(time.Now().UnixNano())

	// setup

	capacity := uint(10)

	genelem := func() *example {
		return &example{rand.Int()}
	}

	count := 100

	// execution

	r := New(capacity)

	for i := 0; i < count; i++ {
		e := genelem()
		fmt.Println("New element", e.i)
		r.Push(e)
	}

	r.Close()

	fmt.Println("Waiting")

	r.Wait()

}
