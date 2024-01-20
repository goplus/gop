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
	"github.com/goplus/mod/modfile"
)

var (
	goxConf = getGoxConf()
)

func getGoxConf() *gox.Config {
	fset := token.NewFileSet()
	imp := packages.NewImporter(fset)
	return &gox.Config{Fset: fset, Importer: imp}
}

func TestErrStringLit(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("TestErrStringLit: no panic?")
		}
	}()
	compileStringLitEx(nil, nil, &ast.BasicLit{
		Value: "Hello",
		Extra: &ast.StringLitEx{
			Parts: []any{1},
		},
	})
}

func TestNodeInterp(t *testing.T) {
	ni := &nodeInterp{}
	if v := ni.Caller(&ast.Ident{}); v != "the function call" {
		t.Fatal("TestNodeInterp:", v)
	}
}

func TestMarkAutogen(t *testing.T) {
	old := noMarkAutogen
	noMarkAutogen = false

	NewPackage("", &ast.Package{Files: map[string]*ast.File{
		"main.t2gmx": {IsProj: true},
	}}, &Config{
		LookupClass: lookupClassErr,
	})

	noMarkAutogen = old
}

func TestClassNameAndExt(t *testing.T) {
	name, ext := ClassNameAndExt("/foo/bar.abc_yap.gox")
	if name != "bar" || ext != "_yap.gox" {
		t.Fatal("classNameAndExt:", name, ext)
	}
}

func TestErrMultiStarRecv(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("TestErrMultiStarRecv: no panic?")
		}
	}()
	getRecvType(&ast.StarExpr{
		X: &ast.StarExpr{},
	})
}

func TestErrAssign(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("TestErrAssign: no panic?")
		}
	}()
	ctx := &blockCtx{}
	compileAssignStmt(ctx, &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent("foo"),
				Sel: ast.NewIdent("bar"),
			},
		},
	})
}

func TestErrPanicToRecv(t *testing.T) {
	ctx := &blockCtx{
		tlookup: &typeParamLookup{
			[]*types.TypeParam{
				types.NewTypeParam(types.NewTypeName(0, nil, "t", nil), nil),
			},
		},
	}
	recv := &ast.FieldList{
		List: []*ast.Field{
			{Type: &ast.SelectorExpr{}},
		},
	}
	func() {
		defer func() {
			if e := recover(); e == nil {
				t.Fatal("TestErrPanicToRecv: no panic?")
			}
		}()
		toRecv(ctx, recv)
	}()
}

func TestCompileErrWrapExpr(t *testing.T) {
	defer func() {
		if e := recover(); e != "TODO: can't use expr? in global" {
			t.Fatal("TestCompileErrWrapExpr failed")
		}
	}()
	pkg := gox.NewPackage("", "foo", goxConf)
	ctx := &blockCtx{pkg: pkg, cb: pkg.CB()}
	compileErrWrapExpr(ctx, &ast.ErrWrapExpr{Tok: token.QUESTION}, 0)
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
	if !isCommandWithoutArgs(
		&ast.SelectorExpr{
			X:   &ast.SelectorExpr{X: ast.NewIdent("foo"), Sel: ast.NewIdent("bar")},
			Sel: ast.NewIdent("val"),
		}) {
		t.Fatal("TestCanAutoCall failed")
	}
}

func TestClRangeStmt(t *testing.T) {
	ctx := &blockCtx{
		cb: &gox.CodeBuilder{},
	}
	stmt := &ast.RangeStmt{
		Tok:  token.DEFINE,
		X:    &ast.SliceLit{},
		Body: &ast.BlockStmt{},
	}
	compileRangeStmt(ctx, stmt)
	stmt.Tok = token.ASSIGN
	stmt.Value = &ast.Ident{Name: "_"}
	compileRangeStmt(ctx, stmt)
}

// -----------------------------------------------------------------------------

func TestGetStringConst(t *testing.T) {
	spx := &gox.PkgRef{Types: types.NewPackage("", "foo")}
	if v := getStringConst(spx, "unknown"); v != "" {
		t.Fatal("getStringConst:", v)
	}
}

func TestSpxRef(t *testing.T) {
	defer func() {
		if e := recover(); !isError(e, "foo.bar not found") {
			t.Fatal("TestSpxRef:", e)
		}
	}()
	pkg := &gox.PkgRef{
		Types: types.NewPackage("foo", "foo"),
	}
	spxRef(pkg, "bar")
}

func isError(e interface{}, msg string) bool {
	if e != nil {
		if err, ok := e.(error); ok {
			return err.Error() == msg
		}
		if err, ok := e.(string); ok {
			return err == msg
		}
	}
	return false
}

func TestGmxProject(t *testing.T) {
	pkg := gox.NewPackage("", "foo", goxConf)
	ctx := &pkgCtx{
		projs:   make(map[string]*gmxProject),
		classes: make(map[*ast.File]gmxClass),
	}
	gmx := loadClass(ctx, pkg, "main.t2gmx", &ast.File{IsProj: true}, &Config{
		LookupClass: lookupClass,
	})
	scheds := gmx.getScheds(pkg.CB())
	if len(scheds) != 2 || scheds[0] == nil || scheds[0] != scheds[1] {
		t.Fatal("TestGmxProject failed")
	}
	gmx.hasScheds = false
	if gmx.getScheds(nil) != nil {
		t.Fatal("TestGmxProject failed: hasScheds?")
	}

	/* _, err := NewPackage("", &ast.Package{Files: map[string]*ast.File{
		"main.t2gmx": {
			IsProj: true,
		},
	}}, &Config{
		LookupClass: lookupClassErr,
	})
	if e := err.Error(); e != `github.com/goplus/gop/cl/internal/libc.Game not found` {
		t.Fatal("newGmx:", e)
	} */

	func() {
		defer func() {
			if e := recover(); e != "TODO: class not found" {
				t.Fatal("TestGmxProject failed:", e)
			}
		}()
		loadClass(nil, pkg, "main.abcx", &ast.File{IsProj: true}, &Config{
			LookupClass: lookupClass,
		})
	}()
	func() {
		defer func() {
			if e := recover(); e != "TODO: multiple project files found" {
				t.Fatal("TestGmxProject failed:", e)
			}
		}()
		loadClass(ctx, pkg, "main.t2gmx", &ast.File{IsProj: true}, &Config{
			LookupClass: lookupClass,
		})
	}()
}

func TestSpxLookup(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("TestSpxLookup failed: no error?")
		}
	}()
	spxLookup(nil, "foo")
}

func lookupClass(ext string) (c *modfile.Project, ok bool) {
	switch ext {
	case ".t2gmx", ".t2spx":
		return &modfile.Project{
			Ext: ".t2gmx", Class: "Game",
			Works:    []*modfile.Class{{Ext: ".t2spx", Class: "Sprite"}},
			PkgPaths: []string{"github.com/goplus/gop/cl/internal/spx2"}}, true
	}
	return
}

func lookupClassErr(ext string) (c *modfile.Project, ok bool) {
	switch ext {
	case ".t2gmx", ".t2spx":
		return &modfile.Project{
			Ext: ".t2gmx", Class: "Game",
			Works:    []*modfile.Class{{Ext: ".t2spx", Class: "Sprite"}},
			PkgPaths: []string{"github.com/goplus/gop/cl/internal/libc"}}, true
	}
	return
}

func TestGetGoFile(t *testing.T) {
	if f := genGoFile("a_test.gop", true); f != testingGoFile {
		t.Fatal("TestGetGoFile:", f)
	}
	if f := genGoFile("a_test.gop", false); f != skippingGoFile {
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

func TestErrLoadImport(t *testing.T) {
	testPanic(t, "-: unknownpkg not found or not a valid C package (c2go.a.pub file not found).\n", func() {
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
