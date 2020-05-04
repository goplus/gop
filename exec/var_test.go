package exec

import (
	"testing"
)

// -----------------------------------------------------------------------------

func TestVar(t *testing.T) {
	sprint, ok := I.FindVariadicFunc("Sprint")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindVariadicFunc failed: Sprintf/strcat")
	}

	x := MakeAddr(0, 0)
	y := MakeAddr(0, 1)
	code := NewBuilder(nil).
		Push(5).
		Push("32").
		CallGoFunv(sprint, 2).
		StoreVar(x). // x = sprint(5, "32")
		Push("78").
		LoadVar(x).
		CallGoFun(strcat).
		StoreVar(y). // y = strcat("78", x)
		Resolve()

	tvars := Struct([]StructField{
		{Name: "X", Type: TyString},
		{Name: "Y", Type: TyString},
	})
	ctx := NewContextEx(code, tvars, nil)
	ctx.Exec(0, code.Len())
	if v := ctx.getVar(1); v != "78532" {
		t.Fatal("y != 78532, ret =", v)
	}
}

func TestParentCtx(t *testing.T) {
	sprint, ok := I.FindVariadicFunc("Sprint")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindVariadicFunc failed: Sprintf/strcat")
	}

	x := MakeAddr(0, 0)
	z := MakeAddr(2, 0)
	code := NewBuilder(nil).
		Push(5).
		Push("32").
		CallGoFunv(sprint, 2).
		StoreVar(x). // x = sprint(5, "32")
		LoadVar(z).
		LoadVar(x).
		CallGoFun(strcat).
		StoreVar(z). // z = strcat(z, x)
		Resolve()

	tvars1 := Struct([]StructField{
		{Name: "Z", Type: TyString},
	})
	tvars2 := Struct([]StructField{
		{Name: "X", Type: TyString},
		{Name: "Y", Type: TyString},
	})
	pp := NewContextEx(code, tvars1, nil)
	pp.FastSetVar(0, "78")

	pp2 := NewContextEx(code, nil, pp)

	ctx := NewContextEx(code, tvars2, pp2)
	ctx.Exec(0, code.Len())
	if v := pp.FastGetVar(0); v != "78532" {
		t.Fatal("y != 78532, ret =", v)
	}
}

func TestAddrVar(t *testing.T) {
	sprint, ok := I.FindVariadicFunc("Sprint")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindVariadicFunc failed: Sprintf/strcat")
	}

	x := MakeAddr(0, 0)
	z := MakeAddr(2, 0)
	code := NewBuilder(nil).
		Push(5).
		Push("32").
		CallGoFunv(sprint, 2).
		StoreVar(x). // x = sprint(5, "32")
		LoadVar(z).
		AddrVar(x). // &x
		AddrOp(String, OpAddrVal).
		CallGoFun(strcat).
		AddrVar(z).
		AddrOp(String, OpAssign). // z = strcat(z, *&x)
		Resolve()

	tvars1 := Struct([]StructField{
		{Name: "Z", Type: TyString},
	})
	tvars2 := Struct([]StructField{
		{Name: "X", Type: TyString},
		{Name: "Y", Type: TyString},
	})
	pp := NewContextEx(code, tvars1, nil)

	pp2 := NewContextEx(code, nil, pp)
	ctx := NewContextEx(code, tvars2, pp2)
	ctx.FastSetVar(z, "78")
	ctx.Exec(0, code.Len())
	if v := ctx.FastGetVar(z); v != "78532" {
		t.Fatal("y != 78532, ret =", v)
	}
}

// -----------------------------------------------------------------------------
