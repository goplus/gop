package exec

import (
	"testing"
)

// -----------------------------------------------------------------------------

func TestIf1(t *testing.T) {

	label1 := NewLabel()
	label2 := NewLabel()
	code := NewBuilder(nil).
		Push(true).
		JmpIfFalse(label1).
		Push(5).
		Push(6).
		BuiltinOp(Int, OpMul).
		Jmp(label2).
		Label(label1).
		Push(5).
		Push(2).
		BuiltinOp(Int, OpMod).
		Label(label2).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v, ok := ctx.Pop(); !(ok && v == 30) {
		t.Fatal("5 6 mul != 30, ret =", v)
	}
}

func TestIf2(t *testing.T) {

	label1 := NewLabel()
	label2 := NewLabel()
	code := NewBuilder(nil).
		Push(false).
		JmpIfFalse(label1).
		Push(5).
		Push(6).
		BuiltinOp(Int, OpMul).
		Jmp(label2).
		Label(label1).
		Push(5).
		Push(2).
		BuiltinOp(Int, OpMod).
		Label(label2).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v, ok := ctx.Pop(); !(ok && v == 1) {
		t.Fatal("5 2 mod != 1, ret =", v)
	}
}

// -----------------------------------------------------------------------------
