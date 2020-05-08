package exec

import (
	"fmt"
	"reflect"
	"testing"
)

// -----------------------------------------------------------------------------

func TestFunc(t *testing.T) {
	strcat, ok := I.FindFunc("strcat")
	if !ok {
		t.Fatal("FindFunc failed: strcat")
	}
	fmt.Println("strcat:", strcat.GetInfo())

	foo := NewFunc("foo")
	code := NewBuilder(nil).
		Push(nil).
		Push("x").
		Push("sw").
		CallFunc(foo).
		Return().
		DefineFunc(
			foo.Return(TyString).
				Args(TyString, TyString)).
		Load(-2).
		Load(-1).
		CallGoFunc(strcat).
		Store(-3).
		EndFunc(foo).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "xsw" {
		t.Fatal("`x` `sw` foo != `xsw`, ret =", v)
	}
}

func TestFuncv(t *testing.T) {
	sprintf, ok := I.FindFuncv("Sprintf")
	if !ok {
		t.Fatal("FindFunc failed: Sprintf")
	}
	fmt.Println("sprintf:", sprintf.GetInfo())

	tyInterfaceSlice := reflect.SliceOf(TyEmptyInterface)

	foo := NewFunc("foo")
	format := NewVar(TyString, "format")
	args := NewVar(tyInterfaceSlice, "args")
	code := NewBuilder(nil).
		Push(nil).
		Push("Hello, %v, %d, %s").
		Push(1.3).
		Push(1).
		Push("xsw").
		CallFuncv(foo, 4).
		Return().
		DefineFunc(
			foo.Return(TyString).
				Vargs(TyString, tyInterfaceSlice).
				DefineVar(format, args)).
		Load(-2).
		StoreVar(format).
		Load(-1).
		StoreVar(args).
		LoadVar(format).
		LoadVar(args).
		CallGoFuncv(sprintf, -1). // sprintf(format, args...)
		Store(-3).
		EndFunc(foo).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "Hello, 1.3, 1, xsw" {
		t.Fatal("format 1.3 1 `xsw` sprintf != `Hello, 1.3, 1, xsw`, ret =", v)
	}
}

func TestFuncLargeArity(t *testing.T) {
	sprint, kind, ok := FindGoPackage("").Find("Sprint")
	if !ok || kind != SymbolFuncv {
		t.Fatal("Find failed: Sprint")
	}

	tyStringSlice := reflect.SliceOf(TyString)

	foo := NewFunc("foo")
	bar := NewFunc("bar")
	b := NewBuilder(nil).
		Push(nil)
	ret := ""
	for i := 0; i < bitsFuncvArityMax+1; i++ {
		b.Push("32")
		ret += "32"
	}
	code := b.
		CallFuncv(foo, bitsFuncvArityMax+1).
		Return().
		DefineFunc(
			bar.Return(TyString).
				Vargs(tyStringSlice)).
		Load(-1).
		CallGoFuncv(GoFuncvAddr(sprint), -1).
		Store(-2).
		EndFunc(bar).
		DefineFunc(
			foo.Return(TyString).
				Vargs(tyStringSlice)).
		Push(nil).
		Load(-1).
		CallFuncv(bar, -1).
		Store(-2).
		EndFunc(foo).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != ret {
		t.Fatal("32 times(1024) sprint != `32` times(1024), ret =", v)
	}
}

// -----------------------------------------------------------------------------
