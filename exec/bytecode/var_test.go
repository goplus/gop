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
	"reflect"
	"testing"
)

// -----------------------------------------------------------------------------

func TestVar(t *testing.T) {
	sprint, ok := I.FindFuncv("Sprint")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindFunc failed: Sprintf/strcat")
	}

	x := defaultImpl.NewVar(TyString, "x").(*Var)
	y := NewVar(TyString, "y")
	b := newBuilder()
	code := b.
		DefineVar(x, y).
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2, 2).
		StoreVar(x). // x = sprint(5, "32")
		Push("78").
		LoadVar(x).
		CallGoFunc(strcat, 2).
		StoreVar(y). // y = strcat("78", x)
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.getVar(1); v != "78532" {
		t.Fatal("y != 78532, ret =", v)
	}
	_ = x.Type()
	_ = b.InCurrentCtx(x)
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
	code := newBuilder().
		DefineVar(z).
		DefineFunc(newFunc("", 2).Args()).
		DefineVar(x, y).
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2, 2).
		StoreVar(x). // x = sprint(5, "32")
		LoadVar(z).
		LoadVar(x).
		CallGoFunc(strcat, 2).
		StoreVar(z). // z = strcat(z, x)
		Resolve()

	ctx := NewContext(code)
	ctx.SetVar(z, "78")

	p0 := ctx.getScope(true)
	old := ctx.switchScope(p0, newVarManager())

	p1 := ctx.getScope(true)
	ctx.switchScope(p1, newVarManager(x, y))

	ctx.Exec(0, code.Len())
	ctx.restoreScope(old)
	if v := ctx.GetVar(z); v != "78532" {
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
	code := newBuilder().
		DefineVar(z).
		DefineFunc(newFunc("", 2).Args()).
		DefineVar(x, y).
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2, 2).
		StoreVar(x). // x = sprint(5, "32")
		LoadVar(z).
		AddrVar(x). // &x
		AddrOp(String, OpAddrVal).
		CallGoFunc(strcat, 2).
		AddrVar(z).
		AddrOp(String, OpAssign). // z = strcat(z, *&x)
		Resolve()

	ctx := NewContext(code)
	ctx.SetVar(z, "78")

	p0 := ctx.getScope(true)
	old := ctx.switchScope(p0, newVarManager())

	p1 := ctx.getScope(true)
	ctx.switchScope(p1, newVarManager(x, y))

	ctx.Exec(0, code.Len())
	ctx.restoreScope(old)
	if v := ctx.GetVar(z); v != "78532" {
		t.Fatal("y != 78532, ret =", v)
	}
}

func TestContext_CloneSetVarScope(t *testing.T) {
	sprint, ok := I.FindFuncv("Sprint")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindFunc failed: Sprintf/strcat")
	}

	x := defaultImpl.NewVar(TyString, "x").(*Var)
	y := NewVar(TyString, "y")
	b := newBuilder()
	code := b.
		DefineVar(x, y).
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2, 2).
		StoreVar(x). // x = sprint(5, "32")
		Push("78").
		LoadVar(x).
		CallGoFunc(strcat, 2).
		StoreVar(y). // y = strcat("78", x)
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())

	z := defaultImpl.NewVar(TyString, "z").(*Var)
	h := NewVar(TyString, "h")
	b1 := newBuilder()
	code = b1.
		DefineVar(z, h).
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2, 2).
		StoreVar(z).
		Push("78").
		LoadVar(z).
		CallGoFunc(strcat, 2).
		StoreVar(h).
		Resolve()

	ctxNew := NewContext(code)
	ctx.CloneSetVarScope(ctxNew)
	if ctxNew.GetVar(y) != ctx.GetVar(y) {
		t.Fatal("clone set var scop err")
	}
}

func TestContext_CloneSetVarScope2(t *testing.T) {
	var m func()
	x := NewVar(TyString, "x")
	foo := newFunc("", 1)
	ret := NewVar(TyString, "ret")
	fn := NewVar(reflect.TypeOf(m), "foo")
	code := newBuilder().
		DefineVar(ret).
		DefineVar(x).
		DefineVar(fn).
		Push("001").
		StoreVar(x).
		GoClosure(foo).
		StoreVar(fn).
		LoadVar(fn).
		CallGoClosure(0, 0, false).
		Return(-1).
		DefineFunc(foo.Args()).
		LoadVar(x).
		StoreVar(ret).
		EndFunc(foo).
		Resolve()

	setPackage(foo, code)

	ctx := NewContext(code)
	ip := ctx.Exec(0, code.Len())

	x2 := NewVar(TyString, "x")
	ret2 := NewVar(TyString, "ret")
	fn2 := NewVar(reflect.TypeOf(m), "foo")

	code2 := newBuilder().
		DefineVar(ret2).
		DefineVar(x2).
		DefineVar(fn2).
		Push("001").
		StoreVar(x).
		GoClosure(foo).
		StoreVar(fn2).
		LoadVar(fn2).
		CallGoClosure(0, 0, false).
		Push("002").
		StoreVar(x).
		LoadVar(fn2).
		CallGoClosure(0, 0, false).
		Return(-1).
		DefineFunc(foo.Args()).
		LoadVar(x).
		StoreVar(ret).
		EndFunc(foo).
		Resolve()

	setPackage(foo, code2)

	ctxNew := NewContext(code2)
	ctx.CloneSetVarScope(ctxNew)
	ctxNew.Exec(ip-1, code2.Len())

	if v := ctxNew.GetVar(ret2); v != "002" {
		t.Fatal("clone set var scop closure err", v)
	}
}

func TestContext_UpdateCode(t *testing.T) {
	var m func()
	x := NewVar(TyString, "x")
	foo := newFunc("", 1)
	ret := NewVar(TyString, "ret")
	fn := NewVar(reflect.TypeOf(m), "foo")
	code := newBuilder().
		DefineVar(ret).
		DefineVar(x).
		DefineVar(fn).
		Push("001").
		StoreVar(x).
		GoClosure(foo).
		StoreVar(fn).
		LoadVar(fn).
		CallGoClosure(0, 0, false).
		Return(-1).
		DefineFunc(foo.Args()).
		LoadVar(x).
		StoreVar(ret).
		EndFunc(foo).
		Resolve()

	setPackage(foo, code)

	ctx := NewContext(code)
	ip := ctx.Exec(0, code.Len())

	x2 := NewVar(TyString, "x")
	ret2 := NewVar(TyString, "ret")
	fn2 := NewVar(reflect.TypeOf(m), "foo")

	code2 := newBuilder().
		DefineVar(ret2).
		DefineVar(x2).
		DefineVar(fn2).
		Push("001").
		StoreVar(x).
		GoClosure(foo).
		StoreVar(fn2).
		LoadVar(fn2).
		CallGoClosure(0, 0, false).
		Push("002").
		StoreVar(x).
		LoadVar(fn2).
		CallGoClosure(0, 0, false).
		Return(-1).
		DefineFunc(foo.Args()).
		LoadVar(x).
		StoreVar(ret).
		EndFunc(foo).
		Resolve()

	setPackage(foo, code2)

	ctx.UpdateCode(code2)
	ctx.Exec(ip-1, code2.Len())

	if v := ctx.GetVar(ret); v != "002" {
		t.Fatal("update code closure err", v)
	}
}

// -----------------------------------------------------------------------------
