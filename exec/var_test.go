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

	code := NewBuilder(nil).
		Push(5).
		Push("32").
		CallGoFunv(sprint, 2).
		StoreVar(0, 0). // x = sprint(5, "32")
		Push("78").
		LoadVar(0, 0).
		CallGoFun(strcat).
		StoreVar(0, 1). // y = strcat("78", x)
		Resolve()

	tvars := Struct([]StructField{
		{Name: "X", Type: TyString},
		{Name: "Y", Type: TyString},
	})
	ctx := NewContextEx(code, tvars, nil)
	ctx.Exec(0, code.Len())
	if v := ctx.GetVar(1); v != "78532" {
		t.Fatal("y != 78532, ret =", v)
	}
}

func TestParentCtx(t *testing.T) {
	sprint, ok := I.FindVariadicFunc("Sprint")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindVariadicFunc failed: Sprintf/strcat")
	}

	code := NewBuilder(nil).
		Push(5).
		Push("32").
		CallGoFunv(sprint, 2).
		StoreVar(0, 0). // x = sprint(5, "32")
		LoadVar(2, 0).
		LoadVar(0, 0).
		CallGoFun(strcat).
		StoreVar(2, 0). // z = strcat(z, x)
		Resolve()

	tvars1 := Struct([]StructField{
		{Name: "Z", Type: TyString},
	})
	tvars2 := Struct([]StructField{
		{Name: "X", Type: TyString},
		{Name: "Y", Type: TyString},
	})
	pp := NewContextEx(code, tvars1, nil)
	pp.SetVar(0, "78")

	pp2 := NewContextEx(code, nil, pp)
	ctx := NewContextEx(code, tvars2, pp2)
	ctx.Exec(0, code.Len())
	if v := pp.GetVar(0); v != "78532" {
		t.Fatal("y != 78532, ret =", v)
	}
}

func TestAddrVar(t *testing.T) {
	sprint, ok := I.FindVariadicFunc("Sprint")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindVariadicFunc failed: Sprintf/strcat")
	}

	code := NewBuilder(nil).
		Push(5).
		Push("32").
		CallGoFunv(sprint, 2).
		StoreVar(0, 0). // x = sprint(5, "32")
		AddrVar(2, 0).
		LoadVar(2, 0).
		AddrVar(0, 0).
		AddrOp(String, OpAddrVal).
		CallGoFun(strcat).
		AddrOp(String, OpAssign). // z = strcat(z, x)
		Resolve()

	tvars1 := Struct([]StructField{
		{Name: "Z", Type: TyString},
	})
	tvars2 := Struct([]StructField{
		{Name: "X", Type: TyString},
		{Name: "Y", Type: TyString},
	})
	pp := NewContextEx(code, tvars1, nil)
	pp.SetVar(0, "78")

	pp2 := NewContextEx(code, nil, pp)
	ctx := NewContextEx(code, tvars2, pp2)
	ctx.Exec(0, code.Len())
	if v := pp.GetVar(0); v != "78532" {
		t.Fatal("y != 78532, ret =", v)
	}
}

// -----------------------------------------------------------------------------
