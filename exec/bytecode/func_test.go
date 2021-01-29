/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package bytecode

import (
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

	foo := newFunc("foo", 1)
	ret := NewVar(TyString, "1")
	code := newBuilder().
		Push("x").
		Push("sw").
		CallFunc(foo, 2).
		Return(-1).
		DefineFunc(
			foo.Return(ret).
				Args(TyString, TyString)).
		Load(foo, -2).
		Load(foo, -1).
		CallGoFunc(strcat, 2).
		StoreVar(ret).
		EndFunc(foo).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "xsw" {
		t.Fatal("`x` `sw` foo != `xsw`, ret =", v)
	}
	_ = foo.Name()
	_ = foo.NumIn()
	_ = foo.NumOut()
	_ = foo.Out(0)
	_ = foo.IsUnnamedOut()
	_ = foo.IsVariadic()
	_ = foo.Type()

	ctx.Push("g")
	ctx.Push("lang")
	ctx.Call(foo)
	if v := checkPop(ctx); v != "glang" {
		t.Fatal("`g` `lang` foo != `glang`, ret =", v)
	}
}

func TestFuncv(t *testing.T) {
	sprintf, ok := I.FindFuncv("Sprintf")
	if !ok {
		t.Fatal("FindFunc failed: Sprintf")
	}

	tyInterfaceSlice := reflect.SliceOf(TyEmptyInterface)

	foo := NewFunc("foo", 1)
	_ = foo.IsUnnamedOut()

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
		Load(foo, -2).
		StoreVar(format).
		Load(foo, -1).
		StoreVar(args).
		LoadVar(format).
		LoadVar(args).
		CallGoFuncv(sprintf, -1). // sprintf(format, args...)
		StoreVar(ret1).
		EndFunc(foo).
		DefineFunc(
			bar.Return(ret2).
				Vargs(TyString, tyInterfaceSlice)).
		Load(bar, -2).
		Load(bar, -1).
		CallFuncv(foo, -1). // foo(format, args...)
		Push("123").
		Store(bar, -2). // format = "123"
		Return(2).
		EndFunc(bar).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "Hello, 1.3, 1, xsw" {
		t.Fatal("format 1.3 1 `xsw` sprintf != `Hello, 1.3, 1, xsw`, ret =", v)
	}
	_ = foo.NumOut()
	_ = foo.Name()

	ctx.Push("Hello, %v, %d, %s")
	ctx.Push([]interface{}{1.3, 1, "xsw"})
	ctx.Call((*iFuncInfo)(bar))
	if v := checkPop(ctx); v != "Hello, 1.3, 1, xsw" {
		t.Fatal("format 1.3 1 `xsw` sprintf != `Hello, 1.3, 1, xsw`, ret =", v)
	}
}

func TestFuncLargeArity(t *testing.T) {
	sprint, kind, ok := FindGoPackage("foo").Find("Sprint")
	if !ok || kind != SymbolFuncv {
		t.Fatal("Find failed: Sprint")
	}

	tyStringSlice := reflect.SliceOf(TyString)

	foo := newFunc("foo", 1)
	bar := newFunc("bar", 1)
	ret1 := NewVar(TyString, "1")
	ret2 := NewVar(TyString, "1")
	b := newBuilder()
	ret := ""
	for i := 0; i < bitsFuncvArityMax+1; i++ {
		b.Push("32")
		ret += "32"
	}
	code := b.
		CallFuncv(foo, bitsFuncvArityMax+1, bitsFuncvArityMax+1).
		Return(-1).
		DefineFunc(
			bar.Return(ret1).
				Vargs(tyStringSlice)).
		Load(bar, -1).
		CallGoFuncv(GoFuncvAddr(sprint), 1, -1).
		StoreVar(ret1).
		EndFunc(bar).
		DefineFunc(
			foo.Return(ret2).
				Vargs(tyStringSlice)).
		Load(foo, -1).
		CallFuncv(bar, 1, -1).
		StoreVar(ret2).
		EndFunc(foo).
		Resolve()

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

	foo := newFunc("foo", 1)
	ret := NewVar(TyString, "1")
	code := newBuilder().
		Push("x").
		Push("sw").
		Closure(foo).
		CallClosure(2, 2, false).
		Return(-1).
		DefineFunc(
			foo.Return(ret).
				Args(TyString, TyString)).
		Load(foo, -2).
		Load(foo, -1).
		CallGoFunc(strcat, 2).
		StoreVar(ret).
		EndFunc(foo).
		Resolve()

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

	tyInterfaceSlice := reflect.SliceOf(TyEmptyInterface)

	foo := newFunc("foo", 2)
	bar := newFunc("bar", 1)
	ret1 := NewVar(TyString, "1")
	ret2 := NewVar(TyString, "1")
	code := newBuilder().
		Push("Hello, %v, %d, %s").
		Push(1.3).
		Push(1).
		Push("xsw").
		Closure(bar).
		CallClosure(4, 4, false).
		Return(-1).
		DefineFunc(
			foo.Return(ret1).
				Vargs(TyString, tyInterfaceSlice)).
		Load(foo, -2).
		Load(foo, -1).
		CallGoFuncv(sprintf, 2, -1). // sprintf(format, args...)
		StoreVar(ret1).
		EndFunc(foo).
		DefineFunc(
			bar.Return(ret2).
				Vargs(TyString, tyInterfaceSlice)).
		Load(bar, -2).
		Load(bar, -1).
		Closure(foo).
		CallClosure(2, -1, true). // foo(format, args...)
		Return(2).
		EndFunc(bar).
		Resolve()

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

	tyInterfaceSlice := reflect.SliceOf(TyEmptyInterface)

	foo := newFunc("foo", 2)
	bar := newFunc("bar", 1)
	ret1 := NewVar(TyString, "1")
	ret2 := NewVar(TyString, "1")
	code := newBuilder().
		Push("Hello, %v, %d, %s").
		Push(1.3).
		Push(1).
		Push("xsw").
		GoClosure(bar).
		CallGoClosure(4, 4, false).
		Return(-1).
		DefineFunc(
			foo.Return(ret1).
				Vargs(TyString, tyInterfaceSlice)).
		Load(foo, -2).
		Load(foo, -1).
		CallGoFuncv(sprintf, 2, -1). // sprintf(format, args...)
		StoreVar(ret1).
		EndFunc(foo).
		DefineFunc(
			bar.Return(ret2).
				Vargs(TyString, tyInterfaceSlice)).
		Load(bar, -2).
		Load(bar, -1).
		GoClosure(foo).
		CallGoClosure(2, -1, true). // foo(format, args...)
		Return(2).
		EndFunc(bar).
		Resolve()

	setPackage(foo, code)
	setPackage(bar, code)

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "Hello, 1.3, 1, xsw" {
		t.Fatal("format 1.3 1 `xsw` sprintf != `Hello, 1.3, 1, xsw`, ret =", v)
	}
	code.(*Code).Dump(os.Stdout) // for code coverage
	ProfileReport()
}

// -----------------------------------------------------------------------------
