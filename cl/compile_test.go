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

package cl

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/ast/asttest"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/qiniu/x/log"

	exec "github.com/goplus/gop/exec/bytecode"
	_ "github.com/goplus/gop/lib"
)

func init() {
	log.SetFlags(log.Ldefault &^ log.LstdFlags)
	log.SetOutputLevel(log.Ldebug)
	CallBuiltinOp = exec.CallBuiltinOp
}

func newPackage(
	out *exec.Builder, pkg *ast.Package, fset *token.FileSet) (p *Package, noExecCtx bool, err error) {
	p, ctx, err := newPackageEx(out, pkg, fset)
	if err != nil {
		return
	}
	entry, _ := ctx.findFunc("main")
	noExecCtx = isNoExecCtx(ctx, entry.body)
	return
}

func newPackageEx(
	out *exec.Builder, pkg *ast.Package, fset *token.FileSet) (p *Package, ctx *blockCtx, err error) {
	b := out.Interface()
	if p, err = NewPackage(b, pkg, fset, PkgActClMain); err != nil {
		return
	}
	ctxPkg := newPkgCtx(b, pkg, fset)
	ctx = newGblBlockCtx(ctxPkg)
	ctx.syms = p.syms
	return
}

// -----------------------------------------------------------------------------

var fsTestBasic = asttest.NewSingleFileFS("/foo", "bar.gop", `
	println("Hello", "xsw", "- nice to meet you!")
	println("Hello, world!")
	return
`)

func TestBasic(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestBasic, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, bctx, err := newPackageEx(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("results:", ctx.Get(-2), ctx.Get(-1))
	if v := ctx.Get(-1); v != nil {
		t.Fatal("error:", v)
	}
	if v := ctx.Get(-2); v != int(14) {
		t.Fatal("n:", v)
	}

	_, _ = NewPackage(nil, nil, nil, 0)
	e := newError(nil, "cannot slice a (type *%v)", "[]int")
	_ = e.Error()
	ev := &execVar{exec.NewVar(reflect.TypeOf(0), "")}
	_ = ev.getType()
	_ = ev.inCurrentCtx(bctx)
	sv := new(stackVar)
	_ = sv.getType()
	_ = sv.inCurrentCtx(nil)
	logError(bctx, nil, "unknown error")
	entry, _ := bctx.findFunc("main")
	defer func() {
		recover()
		defer func() {
			recover()
		}()
		logIllTypeMapIndexPanic(bctx, entry.body, reflect.TypeOf(0), reflect.TypeOf(1.2))
	}()
	logNonIntegerIdxPanic(bctx, entry.body, reflect.String)
}

// -----------------------------------------------------------------------------

var fsTestBasic2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	arr := [...]float64{1, 2}
	slice := make([]float64, 0, 32)
	title := "Hello,world!2020-05-27"
	s := title[0:len(title)-len("2006-01-02")]
	slice = slice[:3:10]
	println(s, len(s), s[0], s[:2], arr[:1], len(arr), cap(arr), slice, len(slice), cap(slice))
`)

func TestBasic2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestBasic2, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, noExecCtx, err := newPackage(b, bar, fset)
	if err != nil || !noExecCtx {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("results:", ctx.Get(-2), ctx.Get(-1))
	if v := ctx.Get(-1); v != nil {
		t.Fatal("error:", v)
	}
	if v := ctx.Get(-2); v != int(43) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestMainPkgNoMain = asttest.NewSingleFileFS("/foo", "bar.gop", `
package main

func ReverseMap(m map[string]int) map[int]string {
    return {v: k for k, v <- m}
}
`)

func TestMainPkgNoMain(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestMainPkgNoMain, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b.Interface(), bar, fset, PkgActClMain)
	if err != ErrMainFuncNotFound {
		t.Fatal("NewPackage failed:", err)
	}
}

// -----------------------------------------------------------------------------

var fsTestPkg = asttest.NewSingleFileFS("/foo", "bar.gop", `
package foo

func ReverseMap(m map[string]int) map[int]string {
    return {v: k for k, v <- m}
}
`)

func TestPkg(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestPkg, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["foo"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b.Interface(), bar, fset, PkgActClMain)
	if err != ErrNotAMainPackage {
		t.Fatal("NewPackage failed:", err)
	}
	pkg, err := NewPackage(b.Interface(), bar, fset, PkgActNone)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}

	kind, sym, ok := pkg.Find("ReverseMap")
	if !ok {
		t.Fatal("pkg.Find failed: ReverseMap not found")
	}
	if kind != SymFunc {
		t.Fatal("pkg.Find failed: kind != SymFunc")
	}

	f := sym.(*FuncDecl).Compile()
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Push(map[string]int{"Hi": 1, "Hello": 2})
	ctx.Call(f)
	if v := ctx.Get(-1); !reflect.DeepEqual(v, map[int]string{1: "Hi", 2: "Hello"}) {
		t.Fatal("ReverseMap failed: ret =", v)
	}
}

func TestPkg2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestPkg, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["foo"]
	b := exec.NewBuilder(nil)
	pkg, err := NewPackage(b.Interface(), bar, fset, PkgActClAll)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	kind, sym, ok := pkg.Find("ReverseMap")
	if !ok || kind != SymFunc {
		t.Fatal("pkg.Find failed: ReverseMap func not found")
	}
	f := sym.(*FuncDecl).Get()
	ctx := exec.NewContext(code)
	ctx.Push(map[string]int{"Hi": 1, "Hello": 2})
	ctx.Call(f)
	if v := ctx.Get(-1); !reflect.DeepEqual(v, map[int]string{1: "Hi", 2: "Hello"}) {
		t.Fatal("ReverseMap failed: ret =", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestType = asttest.NewSingleFileFS("/foo", "bar.gop", `
println(struct {
	A int
	B string
}{1, "Hello"})
println(&struct {
	A int
	B string
}{1, "Hello"})
`)

func TestStruct(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestType, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b.Interface(), bar, fset, PkgActClMain)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Get(-1); v != nil {
		t.Fatal("error:", v)
	}
	if v := ctx.Get(-2); v != int(11) {
		t.Fatal("n:", v)
	}
	if v := ctx.Get(-3); v != nil {
		t.Fatal("error:", v)
	}
	if v := ctx.Get(-4); v != int(10) {
		t.Fatal("n:", v)
	}
}
