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
	"fmt"
	"reflect"
	"testing"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

func init() {
	log.SetOutputLevel(log.Ldebug)
	if opCallGoFunc != SymbolFunc || opCallGoFuncv != SymbolFuncv {
		panic("opCallGoFunc != SymbolFunc || opCallGoFuncv != SymbolFuncv")
	}
}

func Strcat(a, b string) string {
	return a + b
}

func execStrcat(arity int, p *Context) {
	args := p.GetArgs(2)
	ret := Strcat(args[0].(string), args[1].(string))
	p.Ret(2, ret)
}

func execSprint(arity int, p *Context) {
	args := p.GetArgs(arity)
	s := fmt.Sprint(args...)
	p.Ret(arity, s)
}

func execSprintf(arity int, p *Context) {
	args := p.GetArgs(arity)
	s := fmt.Sprintf(args[0].(string), args[1:]...)
	p.Ret(arity, s)
}

func execErrorf(arity int, p *Context) {
	args := p.GetArgs(arity)
	s := fmt.Errorf(args[0].(string), args[1:]...)
	p.Ret(arity, s)
}

// I is a Go package instance.
var I = NewGoPackage("foo")

func init() {
	I.RegisterFuncvs(
		I.Funcv("Sprint", fmt.Sprint, execSprint),
		I.Funcv("Sprintf", fmt.Sprintf, execSprintf),
		I.Funcv("Errorf", fmt.Errorf, execErrorf),
	)
	I.RegisterFuncs(
		I.Func("strcat", Strcat, execStrcat),
	)
	I.RegisterVars(
		I.Var("x", new(int)),
		I.Var("y", new(int)),
	)
	I.RegisterConsts(
		I.Const("true", reflect.Bool, true),
		I.Const("false", reflect.Bool, false),
		I.Const("nil", ConstUnboundPtr, nil),
	)
	I.RegisterTypes(
		I.Rtype(reflect.TypeOf((*Context)(nil))),
		I.Rtype(reflect.TypeOf((*Code)(nil))),
		I.Rtype(reflect.TypeOf((*Stack)(nil))),
		I.Type("rune", TyRune),
	)
	_ = I.PkgPath()
}

func TestVarAndConst(t *testing.T) {
	if ci, ok := I.FindConst("true"); !ok || ci.Value != true {
		t.Fatal("FindConst failed:", ci.Value)
	}
	if ci, ok := I.FindConst("nil"); !ok || ci.Kind != ConstUnboundPtr {
		t.Fatal("FindConst failed:", ci.Kind)
	}
	if addr, ok := I.FindVar("x"); !ok || addr != 0 {
		t.Fatal("FindVar failed:", addr)
	}
}

func TestGoVar(t *testing.T) {
	sprint, ok := I.FindFuncv("Sprint")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindFunc failed: Sprintf/strcat")
	}

	pkg := NewGoPackage("pkg")

	pkg.RegisterVars(
		I.Var("x", new(string)),
		I.Var("y", new(string)),
		I.Var("i", new(int)),
		I.Var("j", new(int)),
	)

	x, ok := pkg.FindVar("x")
	if !ok {
		t.Fatal("FindVar failed:", x)
	}
	y, ok := pkg.FindVar("y")
	if !ok {
		t.Fatal("FindVar failed:", y)
	}
	i, ok := pkg.FindVar("i")
	if !ok {
		t.Fatal("FindVar failed:", i)
	}
	j, ok := pkg.FindVar("j")
	if !ok {
		t.Fatal("FindVar failed:", j)
	}

	b := newBuilder()
	code := b.
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2, 2).
		StoreGoVar(x). // x = sprint(5, "32")
		Push("78").
		LoadGoVar(x).
		CallGoFunc(strcat, 2).
		StoreGoVar(y). // y = strcat("78", x)
		Push(5).
		StoreGoVar(i).
		LoadGoVar(i).
		AddrGoVar(j).
		AddrOp(Int, OpAssign).
		Push(5).
		AddrGoVar(j).
		AddrOp(Int, OpAddAssign).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())

	if v := reflect.ValueOf(govars[x].Addr).Elem().String(); v != "532" {
		t.Fatal("x != 532, x = ", v)
	}
	if v := reflect.ValueOf(govars[y].Addr).Elem().String(); v != "78532" {
		t.Fatal("y != 78532, x =", v)
	}
	if v := reflect.ValueOf(govars[i].Addr).Elem().Int(); v != 5 {
		t.Fatal("i !=5, x =", i)
	}
	if v := reflect.ValueOf(govars[j].Addr).Elem().Int(); v != 10 {
		t.Fatal("j !=10, x =", j)
	}
}

type testBaseInfo struct {
	Info string
}

type testPoint struct {
	testBaseInfo
	X  int
	Y  int
	Ar [5]testBaseInfo
}

func TestGoField(t *testing.T) {
	pkg := NewGoPackage("pkg_field")

	v := testPoint{}
	v.Info = "Info"
	pkg.RegisterVars(
		I.Var("pt", &v),
	)

	i, ok := pkg.FindVar("pt")
	if !ok {
		t.Fatal("FindVar failed:", i)
	}

	x := NewVar(TyInt, "x")
	b := newBuilder()
	code := b.
		DefineVar(x).
		Push("hello").
		AddrGoVar(i).
		StoreGoField([]int{0, 0}).
		Push(-1).
		StoreVar(x).
		LoadVar(x).
		AddrGoVar(i).
		StoreGoField([]int{1}).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())

	if v.Info != "hello" || v.X != -1 {
		t.Fatal("pt", v)
	}
}

func TestSprint(t *testing.T) {
	sprint, ok := I.FindFuncv("Sprint")
	if !ok {
		t.Fatal("FindFuncv failed: Sprint")
	}

	code := newBuilder().
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2, 2).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "532" {
		t.Fatal("5 `32` sprint != 532, ret =", v)
	}
}

func TestSprintf(t *testing.T) {
	sprintf, ok := I.FindFuncv("Sprintf")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindFunc failed: Sprintf/strcat")
	}
	_ = defaultImpl.GetGoFuncType(strcat)
	_ = defaultImpl.GetGoFuncvType(sprintf)

	code := newBuilder().
		Push("Hello, %v, %d, %s").
		Push(1.3).
		Push(1).
		Push("x").
		Push("sw").
		CallGoFunc(strcat, 2).
		CallGoFuncv(sprintf, 4, 4).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "Hello, 1.3, 1, xsw" {
		t.Fatal("format 1.3 1 `x` `sw` strcat sprintf != `Hello, 1.3, 1, xsw`, ret =", v)
	}
}

func TestLargeArity(t *testing.T) {
	sprint, kind, ok := defaultImpl.FindGoPackage("foo").Find("Sprint")
	if !ok || kind != SymbolFuncv {
		t.Fatal("Find failed: Sprint")
	}

	b := newBuilder()
	ret := ""
	for i := 0; i < bitsFuncvArityMax+1; i++ {
		b.Push("32")
		ret += "32"
	}
	code := b.
		CallGoFuncv(GoFuncvAddr(sprint), bitsFuncvArityMax+1, bitsFuncvArityMax+1).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != ret {
		t.Fatal("32 times(1024) sprint != `32` times(1024), ret =", v)
	}
}

func TestType(t *testing.T) {
	typ, ok := I.FindType("Context")
	if !ok {
		t.Fatal("FindType failed: Context not found")
	}
	fmt.Println(typ)

	typ, ok = FindGoPackage("foo").FindType("rune")
	if !ok || typ != TyRune {
		t.Fatal("FindType failed: rune not found")
	}
	fmt.Println(typ)
}
