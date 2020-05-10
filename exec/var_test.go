package exec

import (
	"testing"
)

// -----------------------------------------------------------------------------

func TestVar(t *testing.T) {
	sprint, ok := I.FindFuncv("Sprint")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindFunc failed: Sprintf/strcat")
	}

	x := NewVar(TyString, "x")
	y := NewVar(TyString, "y")
	code := NewBuilder(nil).
		DefineVar(x, y).
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2).
		StoreVar(x). // x = sprint(5, "32")
		Push("78").
		LoadVar(x).
		CallGoFunc(strcat).
		StoreVar(y). // y = strcat("78", x)
		Resolve()

	ctx := NewContext(code, x, y)
	ctx.Exec(0, code.Len())
	if v := ctx.getVar(1); v != "78532" {
		t.Fatal("y != 78532, ret =", v)
	}
}

func TestParentCtx(t *testing.T) {
	sprint, ok := I.FindFuncv("Sprint")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindFunc failed: Sprintf/strcat")
	}

	x := NewVar(TyString, "x")
	y := NewVar(TyString, "y")
	z := NewVar(TyString, "z")
	code := NewBuilder(nil).
		DefineVar(z).
		SetNestDepth(2).
		DefineVar(x, y).
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2).
		StoreVar(x). // x = sprint(5, "32")
		LoadVar(z).
		LoadVar(x).
		CallGoFunc(strcat).
		StoreVar(z). // z = strcat(z, x)
		Resolve()

	p1 := NewContext(code, z)
	p1.SetVar(z, "78")

	p2 := NewContextEx(p1, p1.Stack, p1.code)
	ctx := NewContextEx(p2, p2.Stack, p2.code, x, y)
	ctx.Exec(0, code.Len())
	if v := p1.GetVar(z); v != "78532" {
		t.Fatal("z != 78532, ret =", v)
	}
}

func TestAddrVar(t *testing.T) {
	sprint, ok := I.FindFuncv("Sprint")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindFunc failed: Sprintf/strcat")
	}

	x := NewVar(TyString, "x")
	y := NewVar(TyString, "y")
	z := NewVar(TyString, "z")
	code := NewBuilder(nil).
		DefineVar(z).
		SetNestDepth(2).
		DefineVar(x, y).
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2).
		StoreVar(x). // x = sprint(5, "32")
		LoadVar(z).
		AddrVar(x). // &x
		AddrOp(String, OpAddrVal).
		CallGoFunc(strcat).
		AddrVar(z).
		AddrOp(String, OpAssign). // z = strcat(z, *&x)
		Resolve()

	p1 := NewContext(code, z)
	p1.SetVar(z, "78")

	p2 := NewContextEx(p1, p1.Stack, p1.code)
	ctx := NewContextEx(p2, p2.Stack, p2.code, x, y)
	ctx.Exec(0, code.Len())
	if v := p1.GetVar(z); v != "78532" {
		t.Fatal("y != 78532, ret =", v)
	}
}

// -----------------------------------------------------------------------------
