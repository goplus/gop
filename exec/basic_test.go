package exec

import (
	"testing"
	"qlang.io/lib/builtin"
)

// -----------------------------------------------------------------------------

func TestSubSlice0(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(
		Push(builtin.Max),
		Push(2),
		Push(3.0),
		Push(5),
		Push(6.1),
		Call(builtin.SliceFrom, 4),
		Op3(builtin.SubSlice, false, false),
		CallFnv(1),
	)

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 6.1) {
		t.Fatal("max != 6.1")
	}
}

func TestSubSlice1(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(
		Push(builtin.Max),
		Push(2),
		Push(3.0),
		Push(5),
		Push(6.1),
		Call(builtin.SliceFrom, 4),
		Push(2),
		Op3(builtin.SubSlice, true, false),
		CallFnv(1),
	)

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 6.1) {
		t.Fatal("max != 6.1")
	}
}

func TestSubSlice2(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(
		Push(builtin.Max),
		Push(2),
		Push(3.0),
		Push(5),
		Push(6.1),
		Call(builtin.SliceFrom, 4),
		Push(2),
		Op3(builtin.SubSlice, false, true),
		CallFnv(1),
	)

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 3.0) {
		t.Fatal("max != 3.0")
	}
}

func TestSubSlice3(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(
		Push(builtin.Max),
		Push(2),
		Push(3.0),
		Push(5),
		Push(6.1),
		Call(builtin.SliceFrom, 4),
		Push(2),
		Push(3),
		Op3(builtin.SubSlice, true, true),
		CallFnv(1),
	)

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 5.0) {
		t.Fatal("max != 5.0")
	}
}

// -----------------------------------------------------------------------------

func TestOr1(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(Push(true))
	reserved := code.Reserve()
	code.Block(
		Push(5),
		Push(6),
		Call(builtin.Mul),
	)
	reserved.Set(Or(code.Len() - reserved.Next()))

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == true) {
		t.Fatal("or != true:", v)
	}
}

func TestOr2(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(Push(false))
	reserved := code.Reserve()
	code.Block(
		Push(5),
		Push(6),
		Call(builtin.Mul),
	)
	reserved.Set(Or(code.Len() - reserved.Next()))

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 30) {
		t.Fatal("or != 30:", v)
	}
}

func TestAnd1(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(Push(false))
	reserved := code.Reserve()
	code.Block(
		Push(5),
		Push(6),
		Call(builtin.Mul),
	)
	reserved.Set(And(code.Len() - reserved.Next()))

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == false) {
		t.Fatal("and != false:", v)
	}
}

func TestAnd2(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(Push(true))
	reserved := code.Reserve()
	code.Block(
		Push(5),
		Push(6),
		Call(builtin.Mul),
	)
	reserved.Set(And(code.Len() - reserved.Next()))

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 30) {
		t.Fatal("and != 30:", v)
	}
}

// -----------------------------------------------------------------------------

func TestIf1(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(Push(true))
	reserved1 := code.Reserve()
	code.Block(
		Push(5),
		Push(6),
		Call(builtin.Mul),
	)
	reserved2 := code.Reserve()
	code.Block(
		Push(5),
		Push(2),
		Call(builtin.Mod),
	)
	reserved1.Set(JmpIfFalse(reserved2.Delta(reserved1)))
	reserved2.Set(Jmp(code.Len() - reserved2.Next()))

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 30) {
		t.Fatal("if != 30:", v)
	}
}

func TestIf2(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(Push(false))
	reserved1 := code.Reserve()
	code.Block(
		Push(5),
		Push(6),
		Call(builtin.Mul),
	)
	reserved2 := code.Reserve()
	code.Block(
		Push(5),
		Push(2),
		Call(builtin.Mod),
	)
	reserved1.Set(JmpIfFalse(reserved2.Delta(reserved1)))
	reserved2.Set(Jmp(code.Len() - reserved2.Next()))

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 1) {
		t.Fatal("if != 1:", v)
	}
}

// -----------------------------------------------------------------------------

func TestCase1(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(Push(1))

	code.Block(Push(1))
	case1 := code.Reserve()
	code.Block(
		Push(5),
		Push(6),
		Call(builtin.Mul),
	)
	jmp1 := code.Reserve()
	case1.Set(Case(code.Len() - case1.Next()))

	code.Block(Push(2))
	case2 := code.Reserve()
	code.Block(
		Push(5),
		Push(2),
		Call(builtin.Mod),
	)
	jmp2 := code.Reserve()
	case2.Set(Case(code.Len() - case2.Next()))

	code.Block(
		Default,
		Push(100),
	)
	jmp1.Set(Jmp(code.Len() - jmp1.Next()))
	jmp2.Set(Jmp(code.Len() - jmp2.Next()))

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 30) {
		t.Fatal("case != 30:", v)
	}
}

func TestCase2(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(Push(2))

	code.Block(Push(1))
	case1 := code.Reserve()
	code.Block(
		Push(5),
		Push(6),
		Call(builtin.Mul),
	)
	jmp1 := code.Reserve()
	case1.Set(Case(code.Len() - case1.Next()))

	code.Block(Push(2))
	case2 := code.Reserve()
	code.Block(
		Push(5),
		Push(2),
		Call(builtin.Mod),
	)
	jmp2 := code.Reserve()
	case2.Set(Case(code.Len() - case2.Next()))

	code.Block(
		Default,
		Push(100),
	)
	jmp1.Set(Jmp(code.Len() - jmp1.Next()))
	jmp2.Set(Jmp(code.Len() - jmp2.Next()))

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 1) {
		t.Fatal("case != 1:", v)
	}
}

func TestCase3(t *testing.T) {

	ctx := NewContext()
	stk := NewStack()

	code := New(Push(3))

	code.Block(Push(1))
	case1 := code.Reserve()
	code.Block(
		Push(5),
		Push(6),
		Call(builtin.Mul),
	)
	jmp1 := code.Reserve()
	case1.Set(Case(code.Len() - case1.Next()))

	code.Block(Push(2))
	case2 := code.Reserve()
	code.Block(
		Push(5),
		Push(2),
		Call(builtin.Mod),
	)
	jmp2 := code.Reserve()
	case2.Set(Case(code.Len() - case2.Next()))

	code.Block(
		Default,
		Push(100),
	)
	jmp1.Set(Jmp(code.Len() - jmp1.Next()))
	jmp2.Set(Jmp(code.Len() - jmp2.Next()))

	code.Exec(0, code.Len(), stk, ctx)
	if stk.BaseFrame() != 1 {
		t.Fatal("code.Exec failed")
	}
	if v, ok := stk.Pop(); !(ok && v == 100) {
		t.Fatal("case != 100:", v)
	}
}

// -----------------------------------------------------------------------------

