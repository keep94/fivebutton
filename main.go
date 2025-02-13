// This program prints out all 936 combinations of a 5 button lock and
// assumes that at most 2 buttons can be pressed simultaneously. This program
// prints out the shorter combination sequences first.
package main

import (
	"fmt"
	"iter"
	"strings"

	"github.com/keep94/itertools"
)

type node[T any] struct {
	value T
	next  *node[T]
}

// Queue represents a FIFO queue.
type Queue[T any] struct {
	front *node[T]
	back  *node[T]
}

func NewQueue[T any]() *Queue[T] {
	n := &node[T]{}
	return &Queue[T]{front: n, back: n}
}

// IsEmpty returns true if q is empty.
func (q *Queue[T]) IsEmpty() bool {
	return (q.front == q.back)
}

// Enqueue adds a new value to the end of q.
func (q *Queue[T]) Enqueue(value T) {
	q.back.value = value
	n := &node[T]{}
	q.back.next = n
	q.back = n
}

// Dequeue pops the first value off the beginning of q.
func (q *Queue[T]) Dequeue() T {
	if q.IsEmpty() {
		panic("Queue already empty")
	}
	result := q.front.value
	q.front = q.front.next
	return result
}

// KeyPress represents a key press of a 5 button lock. A KeyPress includes
// pressing a single button e.g "3" or pressing 2 buttons at once, e.g "25"
type KeyPress string

// KeySequence represents an ordered sequence of key presses on a 5 button
// lock.
type KeySequence []KeyPress

// Append appends a KeyPress to the end of ks and returns the resulting
// KeySequence leaving ks unchanged.
func (ks KeySequence) Append(k KeyPress) KeySequence {
	result := make(KeySequence, 0, len(ks)+1)
	result = append(result, ks...)
	return append(result, k)
}

// String converts ks to a string e.g "5-12-4"
func (ks KeySequence) String() string {
	parts := make([]string, len(ks))
	for i := range ks {
		parts[i] = string(ks[i])
	}
	return strings.Join(parts, "-")
}

const NumKeys = 5

// Lock represents the state of a 5 button lock. The zero value of Lock
// represents a 5 button lock with no buttons pushed.
type Lock [NumKeys]bool

// NextPresses returns all the legal next key presses depending on the state
// of this lock.
func (l Lock) NextPresses() iter.Seq[KeyPress] {
	return func(yield func(KeyPress) bool) {
		sb := make([]byte, 0, 2)
		for i := 0; i < NumKeys; i++ {
			if !l[i] {
				sb = sb[:0]
				sb = append(sb, '0'+byte(i+1))
				if !yield(KeyPress(sb)) {
					return
				}
				for j := i + 1; j < NumKeys; j++ {
					if !l[j] {
						sb = sb[:1]
						sb = append(sb, '0'+byte(j+1))
						if !yield(KeyPress(sb)) {
							return
						}
					}
				}
			}
		}
	}
}

// Apply applies the KeyPress kp to l and returns the resulting
// lock while leaving l unchanged.
func (l Lock) Apply(kp KeyPress) Lock {
	for i := 0; i < len(kp); i++ {
		posit := int(kp[i]-'0') - 1
		l[posit] = true
	}
	return l
}

// State contains the state of the lock and the key presses done so far.
type State struct {
	Lock Lock
	Seq  KeySequence
}

// GetCombinations returns all the combinations of a 5 button lock with the
// shorter combination sequences coming first.
func Combinations() iter.Seq[KeySequence] {
	return func(yield func(KeySequence) bool) {
		queue := NewQueue[State]()
		queue.Enqueue(State{})
		for !queue.IsEmpty() {
			state := queue.Dequeue()
			if !yield(state.Seq) {
				return
			}
			for kp := range state.Lock.NextPresses() {
				lock := state.Lock.Apply(kp)
				seq := state.Seq.Append(kp)
				queue.Enqueue(State{Lock: lock, Seq: seq})
			}
		}
	}
}

func main() {
	for i, ks := range itertools.Enumerate(Combinations()) {
		fmt.Println(i+1, ks)
	}
}
