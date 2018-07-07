package vm

import (
	"sync"
)

// Stack is a basic stack implimentation
type Stack struct {
	data    []*Pointer
	pointer int
	// Although every thread has its own stack, vm's main thread still can be accessed by other threads.
	// This is why we need a lock in stack
	// TODO: Find a way to fix this instead of put lock on every stack.
	sync.RWMutex
}

// StackFrame is a segment of a stack. Each method call will start a new stack frame.
type StackFrame struct {
	stack *Stack
	nargs int

	pushed bool
}

// Push an element into the frame
func (f *StackFrame) Push(p *Pointer) {
	f.stack.Push(p)
	f.pushed = true
}

// Pop an element from the stack, if it exists
func (f *StackFrame) Pop() *Pointer {
	if f.NArgs() <= 0 {
		return nil
	}

	return f.stack.Pop()
}

// Peek at an element of the frame
func (f *StackFrame) Peek(i int) *Pointer {
	var p *Pointer
	if f.NArgs()-1 < i {
		p = f.stack.Peek(f.stack.pointer - i)
	}
	return p
}

// NArgs returns the current size of the frame
func (f *StackFrame) NArgs() int {
	return f.nargs
}

// Set a value at a given index in the stack. TODO: Maybe we should be checking for size before we do this.
func (s *Stack) Set(index int, pointer *Pointer) {
	s.Lock()

	s.data[index] = pointer

	s.Unlock()
}

// Push an element to the top of the stack
func (s *Stack) Push(v *Pointer) {
	s.Lock()

	if len(s.data) <= s.pointer {
		s.data = append(s.data, v)
	} else {
		s.data[s.pointer] = v
	}

	s.pointer++
	s.Unlock()
}

// Pop an element off the top of the stack
func (s *Stack) Pop() *Pointer {
	s.Lock()

	if len(s.data) < 1 {
		panic("Nothing to pop!")
	}

	if s.pointer < 0 {
		panic("SP is not normal!")
	}

	if s.pointer > 0 {
		s.pointer--
	}

	v := s.data[s.pointer]
	s.data[s.pointer] = nil
	s.Unlock()
	return v
}

// Peek at an element in the stack
func (s *Stack) Peek(i int) *Pointer {
	var r *Pointer
	s.RLock()
	r = s.data[i]
	s.RUnlock()
	return r
}

func (s *Stack) top() *Pointer {
	var r *Pointer
	s.RLock()

	if len(s.data) == 0 {
		r = nil
	} else if s.pointer > 0 {
		r = s.data[s.pointer-1]
	} else {
		r = s.data[0]
	}

	s.RUnlock()

	return r
}
