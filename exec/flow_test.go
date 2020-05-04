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
		Push(50).
		Push(6).
		BuiltinOp(Int, OpDiv).
		Jmp(label2).
		Label(label1).
		Push(5).
		Push(2).
		BuiltinOp(Int, OpMod).
		Label(label2).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 8 {
		t.Fatal("50 6 div != 8, ret =", v)
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
		Push(5.0).
		Push(2.0).
		BuiltinOp(Float64, OpMul).
		Label(label2).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 10.0 {
		t.Fatal("5.0 2.0 mul != 10.0, ret =", v)
	}
}

// -----------------------------------------------------------------------------

func TestCase1(t *testing.T) {
	done := NewLabel()
	label1 := NewLabel()
	label2 := NewLabel()
	code := NewBuilder(nil).
		Push(1).
		Push(1).
		CaseNE(label1).
		Push(5).
		Push(6).
		BuiltinOp(Int, OpAdd).
		Jmp(done).
		Label(label1).
		Push(2).
		CaseNE(label2).
		Push(5).
		Push(2).
		BuiltinOp(Int, OpMod).
		Jmp(done).
		Label(label2).
		Default().
		Push(100).
		Label(done).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 11 {
		t.Fatal("5 6 add != 11, ret =", v)
	}
}

func TestCase2(t *testing.T) {
	done := NewLabel()
	label1 := NewLabel()
	label2 := NewLabel()
	code := NewBuilder(nil).
		Push(2).
		Push(1).
		CaseNE(label1).
		Push(5).
		Push(6).
		BuiltinOp(Int, OpMul).
		Jmp(done).
		Label(label1).
		Push(2).
		CaseNE(label2).
		Push(5).
		Push(2).
		BuiltinOp(Int, OpSub).
		Jmp(done).
		Label(label2).
		Default().
		Push(100).
		Label(done).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 3 {
		t.Fatal("5 2 sub != 3, ret =", v)
	}
}

func TestDefault(t *testing.T) {
	done := NewLabel()
	label1 := NewLabel()
	label2 := NewLabel()
	code := NewBuilder(nil).
		Push(3).
		Push(1).
		CaseNE(label1).
		Push(5).
		Push(6).
		BuiltinOp(Int, OpMul).
		Jmp(done).
		Label(label1).
		Push(2).
		CaseNE(label2).
		Push(5).
		Push(2).
		BuiltinOp(Int, OpMod).
		Jmp(done).
		Label(label2).
		Default().
		Push(100).
		Label(done).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 100 {
		t.Fatal("100 != 100, ret =", v)
	}
}

// -----------------------------------------------------------------------------
