package exec

import (
)

// -----------------------------------------------------------------------------

// A Stack represents a FILO container.
//
type Stack struct {
	data []interface{}
}

// NewStack returns a new Stack.
//
func NewStack() *Stack {

	data := make([]interface{}, 0, 64)
	return &Stack{data}
}

// Push pushs a value into this stack.
//
func (p *Stack) Push(v interface{}) {

	p.data = append(p.data, v)
}

// Top returns the last pushed value, if it exists.
//
func (p *Stack) Top() (v interface{}, ok bool) {

	n := len(p.data)
	if n > 0 {
		v, ok = p.data[n-1], true
	}
	return
}

// Pop pops a value from this stack.
//
func (p *Stack) Pop() (v interface{}, ok bool) {

	n := len(p.data)
	if n > 0 {
		v, ok = p.data[n-1], true
		p.data = p.data[:n-1]
	}
	return
}

// -----------------------------------------------------------------------------

// A Context represents the context of an executor.
//
type Context struct {
	Stack  *Stack
	Code   *Code
	parent *Context
}

// -----------------------------------------------------------------------------
