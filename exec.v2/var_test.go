package exec

import (
	"testing"
	"qlang.io/qlang/builtin"
)

// -----------------------------------------------------------------------------

func TestAssign(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(
		Push(2),
		Push(3.0),
		Call(builtin.Mul),
		Assign("a"),
	)

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 0 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := ctx.Var("a"); !(ok && v == 6.0) {
		t.Fatal("a != 6")
	}
}

func TestOpAssign(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(
		Push(10),
		Assign("a"),
		Push(2),
		Push(3.0),
		Call(builtin.Mul),
		OpAssign("a", builtin.Mul),
	)

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 0 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := ctx.Var("a"); !(ok && v == 60.0) {
		t.Fatal("a != 6")
	}
}

// -----------------------------------------------------------------------------

