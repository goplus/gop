package exec

import (
	"testing"
	"qlang.io/lib/builtin"
)

// -----------------------------------------------------------------------------

func TestMul(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(
		Push(2),
		Push(3.0),
		Call(builtin.Mul),
	)

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 6.0) {
		t.Fatal("mul != 6")
	}
}

func TestInvalidMax(t *testing.T) {

	defer func() {
		err := recover()
		if err != ErrArityRequired {
			t.Fatal("code.New:", err)
		}
	}()

	New(
		Push(2),
		Push(3.0),
		Push(5),
		Call(builtin.Max),
	)
}

func TestMax(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(
		Push(2),
		Push(3.0),
		Push(5),
		Call(builtin.Max, 3),
	)

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 5.0) {
		t.Fatal("max != 5")
	}
}

func TestCallFn(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(
		Push(builtin.Max),
		Push(2),
		Push(3.0),
		Push(5),
		CallFn(3),
	)

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 5.0) {
		t.Fatal("max != 5")
	}
}

func TestCallFnv(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(
		Push(builtin.Max),
		Push(2),
		Push(3.0),
		Push(5),
		Call(builtin.SliceFrom, 3),
		CallFnv(1),
	)

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 5.0) {
		t.Fatal("max != 5")
	}
}

// -----------------------------------------------------------------------------

