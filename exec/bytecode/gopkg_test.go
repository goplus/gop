/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

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

// I is a Go package instance.
var I = NewGoPackage("")

func init() {
	I.RegisterFuncvs(
		I.Funcv("Sprint", fmt.Sprint, execSprint),
		I.Funcv("Sprintf", fmt.Sprintf, execSprintf),
	)
	I.RegisterFuncs(
		I.Func("strcat", Strcat, execStrcat),
	)
	I.RegisterVars(
		I.Var("x", new(int)),
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

func TestSprint(t *testing.T) {
	sprint, ok := I.FindFuncv("Sprint")
	if !ok {
		t.Fatal("FindFuncv failed: Sprint")
	}

	code := newBuilder().
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2).
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
		CallGoFunc(strcat).
		CallGoFuncv(sprintf, 4).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "Hello, 1.3, 1, xsw" {
		t.Fatal("format 1.3 1 `x` `sw` strcat sprintf != `Hello, 1.3, 1, xsw`, ret =", v)
	}
}

func TestLargeArity(t *testing.T) {
	sprint, kind, ok := defaultImpl.FindGoPackage("").Find("Sprint")
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
		CallGoFuncv(GoFuncvAddr(sprint), bitsFuncvArityMax+1).
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

	typ, ok = FindGoPackage("").FindType("rune")
	if !ok || typ != TyRune {
		t.Fatal("FindType failed: rune not found")
	}
	fmt.Println(typ)
}
