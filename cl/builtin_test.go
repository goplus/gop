/*
 Copyright 2021 The GoPlus Authors (goplus.org)

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

package cl

import (
	"go/types"
	"testing"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

func TestCompileErrWrapExpr(t *testing.T) {
	defer func() {
		if e := recover(); e != "TODO: can't use expr? in global" {
			t.Fatal("TestCompileErrWrapExpr failed")
		}
	}()
	pkg := gox.NewPackage("", "foo", nil)
	ctx := &blockCtx{pkg: pkg, cb: pkg.CB()}
	compileErrWrapExpr(ctx, &ast.ErrWrapExpr{Tok: token.QUESTION})
}

func TestToString(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("toString: no error?")
		}
	}()
	toString(&ast.BasicLit{Kind: token.INT, Value: "1"})
}

func TestGetTypeName(t *testing.T) {
	if getTypeName(types.Typ[types.Int]) != "int" {
		t.Fatal("getTypeName int failed")
	}
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("getTypeName: no error?")
		}
	}()
	getTypeName(types.NewSlice(types.Typ[types.Int]))
}

func TestIsFunc(t *testing.T) {
	if isFunc(nil) {
		t.Fatal("TestIsFunc failed")
	}
}

func TestHandleRecover(t *testing.T) {
	var ctx pkgCtx
	ctx.handleRecover("hello")
	if !(len(ctx.errs) == 1 && ctx.errs[0].Error() == "hello") {
		t.Fatal("TestHandleRecover failed:", ctx.errs)
	}
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("TestHandleRecover failed: no error?")
		}
	}()
	ctx.handleRecover(100)
}

func TestCanAutoCall(t *testing.T) {
	if !canAutoCall(
		&ast.SelectorExpr{
			X:   &ast.SelectorExpr{X: ast.NewIdent("foo"), Sel: ast.NewIdent("bar")},
			Sel: ast.NewIdent("val"),
		}) {
		t.Fatal("TestCanAutoCall failed")
	}
}

// -----------------------------------------------------------------------------

func TestGmxSettings(t *testing.T) {
	pkg := gox.NewPackage("", "foo", nil)
	gmx := newGmx(pkg, "index.t2gmx")
	if len(gmx.scheds) != 2 || gmx.scheds[0] == nil || gmx.scheds[0] != gmx.scheds[1] {
		t.Fatal("TestGmxSettings failed")
	}
}

func init() {
	RegisterClassFileType(".t2gmx", ".t2spx", "github.com/goplus/gop/cl/internal/spx2")
}

// -----------------------------------------------------------------------------
