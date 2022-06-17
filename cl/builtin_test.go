/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cl

import (
	"errors"
	"go/types"
	"testing"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
	"github.com/goplus/gox/cpackages"
	"github.com/goplus/gox/packages"
	"github.com/goplus/mod/gopmod"
)

var (
	goxConf = getGoxConf()
)

func getGoxConf() *gox.Config {
	fset := token.NewFileSet()
	imp := packages.NewImporter(fset)
	return &gox.Config{Fset: fset, Importer: imp}
}

func TestCompileErrWrapExpr(t *testing.T) {
	defer func() {
		if e := recover(); e != "TODO: can't use expr? in global" {
			t.Fatal("TestCompileErrWrapExpr failed")
		}
	}()
	pkg := gox.NewPackage("", "foo", goxConf)
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
	pkg := gox.NewPackage("", "foo", goxConf)
	gmx := newGmx(nil, pkg, "main.t2gmx", &Config{
		LookupClass: lookupClass,
	})
	scheds := gmx.getScheds(pkg.CB())
	if len(scheds) != 2 || scheds[0] == nil || scheds[0] != scheds[1] {
		t.Fatal("TestGmxSettings failed")
	}
	gmx.hasScheds = false
	if gmx.getScheds(nil) != nil {
		t.Fatal("TestGmxSettings failed: hasScheds?")
	}
}

func TestSpxLookup(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("TestSpxLookup failed: no error?")
		}
	}()
	spxLookup(nil, "foo")
}

func lookupClass(ext string) (c *gopmod.Class, ok bool) {
	switch ext {
	case ".t2gmx", ".t2spx":
		return &gopmod.Class{
			ProjExt: ".t2gmx", WorkExt: ".t2spx",
			PkgPaths: []string{"github.com/goplus/gop/cl/internal/spx2"}}, true
	}
	return
}

func TestGetGoFile(t *testing.T) {
	if f := getGoFile("a_test.gop", true); f != testingGoFile {
		t.Fatal("TestGetGoFile:", f)
	}
	if f := getGoFile("a_test.gop", false); f != skippingGoFile {
		t.Fatal("TestGetGoFile:", f)
	}
}

func TestC2goBase(t *testing.T) {
	if c2goBase("") != "github.com/goplus/" {
		t.Fatal("c2goBase failed")
	}
}

func TestErrNewType(t *testing.T) {
	testPanic(t, `bar redeclared in this block
	previous declaration at <TODO>
`, func() {
		pkg := types.NewPackage("", "foo")
		newType(pkg, token.NoPos, "bar")
		newType(pkg, token.NoPos, "bar")
	})
}

func TestErrDeclFunc(t *testing.T) {
	testPanic(t, "invalid receiver type **byte (**byte is not a defined type)\n", func() {
		pkg := gox.NewPackage("", "foo", goxConf)
		recv := pkg.NewParam(token.NoPos, "p", types.NewPointer(types.NewPointer(gox.TyByte)))
		declFunc(&blockCtx{pkg: pkg}, recv, &ast.FuncDecl{
			Name: &ast.Ident{Name: "m"},
			Type: &ast.FuncType{Params: &ast.FieldList{}},
		})
	})
}

func TestErrLoadImport(t *testing.T) {
	testPanic(t, ".: unknownpkg not found or not a valid C package (c2go.a.pub file not found).\n", func() {
		pkg := &pkgCtx{
			nodeInterp: &nodeInterp{
				fset: token.NewFileSet(),
			},
			cpkgs: cpackages.NewImporter(
				&cpackages.Config{LookupPub: func(pkgPath string) (pubfile string, err error) {
					return "", errors.New("not found")
				}})}
		ctx := &blockCtx{pkgCtx: pkg}
		spec := &ast.ImportSpec{
			Path: &ast.BasicLit{Kind: token.STRING, Value: `"C/unknownpkg"`},
		}
		loadImport(ctx, spec)
		panic(ctx.errs[0].Error())
	})
}

func TestErrCompileBasicLit(t *testing.T) {
	testPanic(t, "compileBasicLit: invalid syntax\n", func() {
		ctx := &blockCtx{cb: new(gox.CodeBuilder)}
		compileBasicLit(ctx, &ast.BasicLit{Kind: token.CSTRING, Value: `\\x`})
	})
}

func testPanic(t *testing.T, panicMsg string, doPanic func()) {
	t.Run(panicMsg, func(t *testing.T) {
		defer func() {
			if e := recover(); e == nil {
				t.Fatal("testPanic: no error?")
			} else if msg := e.(string); msg != panicMsg {
				t.Fatalf("\nResult:\n%s\nExpected Panic:\n%s\n", msg, panicMsg)
			}
		}()
		doPanic()
	})
}

// -----------------------------------------------------------------------------
