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

package cl

import (
	"fmt"
	"os"
	"testing"

	"github.com/qiniu/goplus/ast/asttest"
	"github.com/qiniu/goplus/parser"
	"github.com/qiniu/goplus/token"

	exec "github.com/qiniu/goplus/exec/bytecode"
)

// -----------------------------------------------------------------------------

var fsTestAssign = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x, y := 123, "Hello"
	x
	y
`)

func TestAssign(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestAssign, "/foo", nil, 0)
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
	fmt.Println("x, y:", ctx.Get(-2), ctx.Get(-1))
	if v := ctx.Get(-2); v != 123 {
		t.Fatal("x:", v)
	}
	if v := ctx.Get(-1); v != "Hello" {
		t.Fatal("y:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestSwif = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := 0
	t := "Hello"
	switch {
	case t == "xsw":
		x = 3
	case t == "Hello", t == "world":
		x = 5
	default:
		x = 7
	}
	x
`)

func TestSwitchIf(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestSwif, "/foo", nil, 0)
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
	fmt.Println("x:", ctx.Get(-1))
	if v := ctx.Get(-1); v != 5 {
		t.Fatal("x:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestSwif2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := 0
	t := "Hello"
	switch {
	case t == "xsw":
		x = 3
	case t == "hi", t == "world":
		x = 11
	default:
		x = 7
	}
	x
`)

func TestSwitchIfDefault(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestSwif2, "/foo", nil, 0)
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

	code.Dump(os.Stdout)

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("x:", ctx.Get(-1))
	if v := ctx.Get(-1); v != 7 {
		t.Fatal("x:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestSw = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := 0
	switch t := "Hello"; t {
	case "xsw":
		x = 3
	case "Hello", "world":
		x = 5
	default:
		x= 7
	}
	x
`)

func TestSwitch(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestSw, "/foo", nil, 0)
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
	fmt.Println("x:", ctx.Get(-1))
	if v := ctx.Get(-1); v != 5 {
		t.Fatal("x:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestSw2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := 0
	t := "Hello"
	switch t {
	}
	switch t {
	case "world", "Hello":
		x = 5
	case "xsw":
		x = 3
	}
	x
`)

func TestSwitch2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestSw2, "/foo", nil, 0)
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
	fmt.Println("x:", ctx.Get(-1))
	if v := ctx.Get(-1); v != 5 {
		t.Fatal("x:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestSw3 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := 0
	t := "Hello"
	switch t {
	default:
		x = 7
	case "world", "hi":
		x = 5
	case "xsw":
		x = 3
	}
	x
`)

func TestDefault(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestSw3, "/foo", nil, 0)
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
	fmt.Println("x:", ctx.Get(-1))
	if v := ctx.Get(-1); v != 7 {
		t.Fatal("x:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestIf = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := 0
	if t := false; t {
		x = 3
	} else {
		x = 5
	}
	x
`)

func TestIf(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestIf, "/foo", nil, 0)
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
	fmt.Println("x:", ctx.Get(-1))
	if v := ctx.Get(-1); v != 5 {
		t.Fatal("x:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestIf2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := 5
	if true {
		x = 3
	}
	x
`)

func TestIf2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestIf2, "/foo", nil, 0)
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
	fmt.Println("x:", ctx.Get(-1))
	if v := ctx.Get(-1); v != 3 {
		t.Fatal("x:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestReturn = asttest.NewSingleFileFS("/foo", "bar.gop", `
	import (
		"fmt"
		"strings"
	)

	func foo(x string) string {
		return strings.NewReplacer("?", "!").Replace(x)
	}

	fmt.Println(foo("Hello, world???"))
`)

func TestReturn(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestReturn, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(16) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestReturn2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	func max(a, b int) int {
		if a < b {
			a = b
		}
		return a
	}

	println("max(23,345):", max(23,345))
`)

func TestReturn2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestReturn2, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(17) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestFunc = asttest.NewSingleFileFS("/foo", "bar.gop", `
	import "fmt"

	func foo(x string) (n int, err error) {
		n, err = fmt.Println("x: " + x)
		return
	}

	foo("Hello, world!")
`)

func TestFunc(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestFunc, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(17) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestFuncv = asttest.NewSingleFileFS("/foo", "bar.gop", `
	import "fmt"

	func foo(format string, args ...interface{}) (n int, err error) {
		n, err = printf(format, args...)
		return
	}

	func bar(foo func(string, ...interface{}) (int, error)) {
		foo("Hello, %v!\n", "glang")
	}

	bar(foo)
	println(foo("Hello, %v!\n", 123))
`)

func TestFuncv(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestFuncv, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(9) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestClosure = asttest.NewSingleFileFS("/foo", "bar.gop", `
	import "fmt"

	foo := func(prompt string) (n int, err error) {
		n, err = fmt.Println(prompt + x)
		return
	}

	x := "Hello, world!"
	foo("x: ")
`)

func TestClosure(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestClosure, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, noExecCtx, err := newPackage(b, bar, fset)
	if err != nil || noExecCtx {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("results:", ctx.Get(-2), ctx.Get(-1))
	if v := ctx.Get(-1); v != nil {
		t.Fatal("error:", v)
	}
	if v := ctx.Get(-2); v != int(17) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestClosurev = asttest.NewSingleFileFS("/foo", "bar.gop", `
	import "fmt"

	foo := func(format string, args ...interface{}) (n int, err error) {
		n, err = fmt.Printf(format, args...)
		return
	}

	foo("Hello, %v!\n", "xsw")
`)

func TestClosurev(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestClosurev, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, noExecCtx, err := newPackage(b, bar, fset)
	if err != nil || noExecCtx {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("results:", ctx.Get(-2), ctx.Get(-1))
	if v := ctx.Get(-1); v != nil {
		t.Fatal("error:", v)
	}
	if v := ctx.Get(-2); v != int(12) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestForPhraseStmt = asttest.NewSingleFileFS("/foo", "bar.gop", `
	sum := 0
	for x <- [1, 3, 5, 7], x < 7 {
		sum += x
	}
	println("sum(1,3,5):", sum)
`)

func TestForPhraseStmt(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestForPhraseStmt, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(14) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestForPhraseStmt2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	fns := make([]func() int, 3)
	for i, x <- [3, 15, 777] {
		v := x
		fns[i] = func() int {
			return v
		}
	}
	println("values:", fns[0](), fns[1](), fns[2]())
`)

func TestForPhraseStmt2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestForPhraseStmt2, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, noExecCtx, err := newPackage(b, bar, fset)
	if err != nil || noExecCtx {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("results:", ctx.Get(-2), ctx.Get(-1))
	if v := ctx.Get(-1); v != nil {
		t.Fatal("error:", v)
	}
	if v := ctx.Get(-2); v != 17 {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestForRangeStmt = asttest.NewSingleFileFS("/foo", "bar.gop", `
	sum := 0
	for x := range [1, 3, 5, 7] {
		if x < 7 {
			sum += x
		}
	}
	println("sum(1,3,5):", sum)
`)

func _TestForRangeStmt(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestForRangeStmt, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(15) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestForIncDecStmt = asttest.NewSingleFileFS("/foo", "bar.gop", `
	a,b:=10,2
	{a--;a--;a--}
	{b++;b++;b++}
	println(a,b,a*b)
`)

func TestForIncDecStmt(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestForIncDecStmt, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(7) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------
