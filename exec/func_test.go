package exec

import (
	"fmt"
	"os"
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

	foo := NewFunc("foo", 1)
	ret := NewVar(TyString, "1")
	code := NewBuilder(nil).
		Push("x").
		Push("sw").
		CallFunc(foo).
		Return(-1).
		DefineFunc(
			foo.Return(ret).
				Args(TyString, TyString)).
		Load(-2).
		Load(-1).
		CallGoFunc(strcat).
		StoreVar(ret).
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

	foo := NewFunc("foo", 1)
	bar := NewFunc("bar", 1)
	ret1 := NewVar(TyString, "1")
	ret2 := NewVar(TyString, "1")
	format := NewVar(TyString, "format")
	args := NewVar(tyInterfaceSlice, "args")
	code := NewBuilder(nil).
		Push("Hello, %v, %d, %s").
		Push(1.3).
		Push(1).
		Push("xsw").
		CallFuncv(bar, 4).
		Return(-1).
		DefineFunc(
			foo.Return(ret1).
				Vargs(TyString, tyInterfaceSlice)).
		DefineVar(format, args).
		Load(-2).
		StoreVar(format).
		Load(-1).
		StoreVar(args).
		LoadVar(format).
		LoadVar(args).
		CallGoFuncv(sprintf, -1). // sprintf(format, args...)
		StoreVar(ret1).
		EndFunc(foo).
		DefineFunc(
			bar.Return(ret2).
				Vargs(TyString, tyInterfaceSlice)).
		Load(-2).
		Load(-1).
		CallFuncv(foo, -1). // foo(format, args...)
		Return(2).
		EndFunc(bar).
		Resolve()

	code.Dump(os.Stdout)

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

	foo := NewFunc("foo", 1)
	bar := NewFunc("bar", 1)
	ret1 := NewVar(TyString, "1")
	ret2 := NewVar(TyString, "1")
	b := NewBuilder(nil)
	ret := ""
	for i := 0; i < bitsFuncvArityMax+1; i++ {
		b.Push("32")
		ret += "32"
	}
	code := b.
		CallFuncv(foo, bitsFuncvArityMax+1).
		Return(-1).
		DefineFunc(
			bar.Return(ret1).
				Vargs(tyStringSlice)).
		Load(-1).
		CallGoFuncv(GoFuncvAddr(sprint), -1).
		StoreVar(ret1).
		EndFunc(bar).
		DefineFunc(
			foo.Return(ret2).
				Vargs(tyStringSlice)).
		Load(-1).
		CallFuncv(bar, -1).
		StoreVar(ret2).
		EndFunc(foo).
		Resolve()

	if bar.IsTypeValid() {
		fmt.Println("func bar:", bar.Type())
	}

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != ret {
		t.Fatal("32 times(1024) sprint != `32` times(1024), ret =", v)
	}
}

func TestClosure(t *testing.T) {
	strcat, ok := I.FindFunc("strcat")
	if !ok {
		t.Fatal("FindFunc failed: strcat")
	}
	fmt.Println("strcat:", strcat.GetInfo())

	foo := NewFunc("foo", 1)
	ret := NewVar(TyString, "1")
	code := NewBuilder(nil).
		Push("x").
		Push("sw").
		Closure(foo).
		CallClosure(2).
		Return(-1).
		DefineFunc(
			foo.Return(ret).
				Args(TyString, TyString)).
		Load(-2).
		Load(-1).
		CallGoFunc(strcat).
		StoreVar(ret).
		EndFunc(foo).
		Resolve()

	code.Dump(os.Stdout)

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "xsw" {
		t.Fatal("`x` `sw` foo != `xsw`, ret =", v)
	}
}

func TestClosure2(t *testing.T) {
	sprintf, ok := I.FindFuncv("Sprintf")
	if !ok {
		t.Fatal("FindFunc failed: Sprintf")
	}
	fmt.Println("sprintf:", sprintf.GetInfo())

	tyInterfaceSlice := reflect.SliceOf(TyEmptyInterface)

	foo := NewFunc("foo", 2)
	bar := NewFunc("bar", 1)
	ret1 := NewVar(TyString, "1")
	ret2 := NewVar(TyString, "1")
	code := NewBuilder(nil).
		Push("Hello, %v, %d, %s").
		Push(1.3).
		Push(1).
		Push("xsw").
		CallFuncv(bar, 4).
		Return(-1).
		DefineFunc(
			foo.Return(ret1).
				Vargs(TyString, tyInterfaceSlice)).
		Load(-2).
		Load(-1).
		CallGoFuncv(sprintf, -1). // sprintf(format, args...)
		StoreVar(ret1).
		EndFunc(foo).
		DefineFunc(
			bar.Return(ret2).
				Vargs(TyString, tyInterfaceSlice)).
		Load(-2).
		Load(-1).
		Closure(foo).
		CallClosure(-1). // foo(format, args...)
		Return(2).
		EndFunc(bar).
		Resolve()

	code.Dump(os.Stdout)

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "Hello, 1.3, 1, xsw" {
		t.Fatal("format 1.3 1 `xsw` sprintf != `Hello, 1.3, 1, xsw`, ret =", v)
	}
}

func TestGoClosure(t *testing.T) {
	sprintf, ok := I.FindFuncv("Sprintf")
	if !ok {
		t.Fatal("FindFunc failed: Sprintf")
	}
	fmt.Println("sprintf:", sprintf.GetInfo())

	tyInterfaceSlice := reflect.SliceOf(TyEmptyInterface)

	foo := NewFunc("foo", 2)
	bar := NewFunc("bar", 1)
	ret1 := NewVar(TyString, "1")
	ret2 := NewVar(TyString, "1")
	code := NewBuilder(nil).
		Push("Hello, %v, %d, %s").
		Push(1.3).
		Push(1).
		Push("xsw").
		GoClosure(bar).
		CallGoClosure(4).
		Return(-1).
		DefineFunc(
			foo.Return(ret1).
				Vargs(TyString, tyInterfaceSlice)).
		Load(-2).
		Load(-1).
		CallGoFuncv(sprintf, -1). // sprintf(format, args...)
		StoreVar(ret1).
		EndFunc(foo).
		DefineFunc(
			bar.Return(ret2).
				Vargs(TyString, tyInterfaceSlice)).
		Load(-2).
		Load(-1).
		GoClosure(foo).
		CallGoClosure(-1). // foo(format, args...)
		Return(2).
		EndFunc(bar).
		Resolve()

	code.Dump(os.Stdout)

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "Hello, 1.3, 1, xsw" {
		t.Fatal("format 1.3 1 `xsw` sprintf != `Hello, 1.3, 1, xsw`, ret =", v)
	}
}

// -----------------------------------------------------------------------------
