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

const (

	// Number of buttons in lock
	NumButtons = 5

	// Maximum number of buttons that can be pressed in a single key press
	MaxAtOnce = 2
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
// pressing a single button e.g "3" or pressing multiple buttons at once,
// e.g "25". The button parameter in the methods for this type is zero
// based. That is 0 means the "1" button. 1 means the "2" button etc.
type KeyPress [NumButtons]bool

// SingleKeyPress returns a key press involving a single button.
func SingleKeyPress(button int) KeyPress {
	var result KeyPress
	result[button] = true
	return result
}

// Add adds an additional button to returned KeyPress while leaving k
// unchanged.
func (k KeyPress) Add(button int) KeyPress {
	k[button] = true
	return k
}

// Highest returns the 0 based index of the highest button pressed in k.
func (k KeyPress) Highest() int {
	for i := NumButtons - 1; i >= 0; i-- {
		if k[i] {
			return i
		}
	}
	panic("Highest called on zero KeyPress")
}

// Len returns the number of simultaneously pressed buttons in k.
func (k KeyPress) Len() int {
	var result int
	for i := 0; i < NumButtons; i++ {
		if k[i] {
			result++
		}
	}
	return result
}

// String returns the string representation of k. e.g "3" or "25" where
// the numerals are one based.
func (k KeyPress) String() string {
	buffer := make([]byte, 0, NumButtons)
	for i := 0; i < NumButtons; i++ {
		if k[i] {
			buffer = append(buffer, '1'+byte(i))
		}
	}
	return string(buffer)
}

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
		parts[i] = ks[i].String()
	}
	return strings.Join(parts, "-")
}

// Lock represents the state of a 5 button lock. The zero value of Lock
// represents a 5 button lock with no buttons pushed.
type Lock [NumButtons]bool

// NextPresses returns all the legal next key presses depending on the state
// of this lock.
func (l Lock) NextPresses() iter.Seq[KeyPress] {
	return func(yield func(KeyPress) bool) {
		queue := NewQueue[KeyPress]()
		for i := 0; i < NumButtons; i++ {
			if !l[i] {
				queue.Enqueue(SingleKeyPress(i))
			}
		}
		for !queue.IsEmpty() {
			press := queue.Dequeue()
			if !yield(press) {
				return
			}
			if press.Len() == MaxAtOnce {
				continue
			}
			for i := press.Highest() + 1; i < NumButtons; i++ {
				if !l[i] {
					queue.Enqueue(press.Add(i))
				}
			}
		}
	}
}

// Apply applies the KeyPress kp to l and returns the resulting
// lock while leaving l unchanged.
func (l Lock) Apply(kp KeyPress) Lock {
	for i := 0; i < NumButtons; i++ {
		if kp[i] {
			l[i] = true
		}
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
