package exec

import (
	"testing"
)

// -----------------------------------------------------------------------------

func NewContext() *Context {

	return NewContextEx(nil)
}

func TestCodeLine(t *testing.T) {

	code := New(Nil, Nil, Nil)
	code.CodeLine("", 1)
	code.Block(Nil, Nil, Nil)
	code.CodeLine("", 2)

	if _, line := code.Line(0); line != 1 {
		t.Fatal("code.Line(0):", line)
	}

	if _, line := code.Line(2); line != 1 {
		t.Fatal("code.Line(2):", line)
	}

	if _, line := code.Line(3); line != 2 {
		t.Fatal("code.Line(3):", line)
	}

	if _, line := code.Line(4); line != 2 {
		t.Fatal("code.Line(4):", line)
	}

	if _, line := code.Line(5); line != 2 {
		t.Fatal("code.Line(6):", line)
	}

	if _, line := code.Line(6); line != 0 {
		t.Fatal("code.Line(6):", line)
	}
}

// -----------------------------------------------------------------------------
