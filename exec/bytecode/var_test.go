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
	old := ctx.switchScope(p0, newVarManager(), nil)

	p1 := ctx.getScope(true)
	ctx.switchScope(p1, newVarManager(x, y), nil)

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
	old := ctx.switchScope(p0, newVarManager(), nil)

	p1 := ctx.getScope(true)
	ctx.switchScope(p1, newVarManager(x, y), nil)

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

// -----------------------------------------------------------------------------
