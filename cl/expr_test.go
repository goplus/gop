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
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/qiniu/goplus/ast/asttest"
	"github.com/qiniu/goplus/parser"
	"github.com/qiniu/goplus/token"

	exec "github.com/qiniu/goplus/exec/bytecode"
)

// -----------------------------------------------------------------------------

var fsTestUnbound = asttest.NewSingleFileFS("/foo", "bar.gop", `
	println("Hello " + "qiniu:", 123, 4.5, 7i)
`)

func TestUnbound(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestUnbound, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
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
	if v := ctx.Get(-2); v != int(28) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestTypeCast = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := []byte("hello")
	x
`)

func TestTypeCast(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestTypeCast, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Get(-1); !bytes.Equal(v.([]byte), []byte("hello")) {
		t.Fatal("v:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestMake = asttest.NewSingleFileFS("/foo", "bar.gop", `
	a := make([]int, 0, 4)
	a = append(a, 1, 2, 3)
	println(a)
`)

func TestMake(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestMake, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
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
	if v := ctx.Get(-2); v != int(8) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestMake2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	a := make([]int, 0, 4)
	a = append(a, [1, 2, 3]...)
	println(a)
`)

func TestMake2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestMake2, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
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
	if v := ctx.Get(-2); v != int(8) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestMake3 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	n := 4
	a := make(map[string]interface{}, uint16(n))
	println(a)
`)

func TestMake3(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestMake3, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
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
	if v := ctx.Get(-2); v != int(6) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestMake4 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	import "reflect"

	a := make(chan *func(), uint16(4))
	println(reflect.TypeOf(a))
`)

func TestMake4(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestMake4, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
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
	if v := ctx.Get(-2); v != int(13) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestOperator = asttest.NewSingleFileFS("/foo", "bar.gop", `
	println("Hello", 123 * 4.5, 1 + 7i)
`)

func TestOperator(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestOperator, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
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
	if v := ctx.Get(-2); v != int(19) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestVar = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := 123.1
	println("Hello", x)
`)

func TestVar(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestVar, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
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
	if v := ctx.Get(-2); v != int(12) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestVarOp = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := 123.1
	y := 1 + x
	println("Hello", y + 10)
	n, err := println("Hello", y + 10)
	println("ret:", n << 1, err)
`)

func TestVarOp(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestVarOp, "/foo", nil, 0)
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

var fsTestGoPackage = asttest.NewSingleFileFS("/foo", "bar.gop", `
	import "fmt"
	import gostrings "strings"

	x := gostrings.NewReplacer("?", "!").Replace("hello, world???")
	fmt.Println("x: " + x)
`)

func TestGoPackage(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestGoPackage, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(19) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestSlice = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := []float64{1, 2.3, 3.6}
	println("x:", x)
`)

func TestSlice(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestSlice, "/foo", nil, 0)
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

var fsTestSlice2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := []float64{1, 2: 3.4, 5}
	println("x:", x)
`)

func TestSlice2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestSlice2, "/foo", nil, 0)
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

var fsTestArray = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := [4]float64{1, 2.3, 3.6}
	println("x:", x)

	y := [...]float64{1, 2.3, 3.6}
	println("y:", y)
`)

func TestArray(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestArray, "/foo", nil, 0)
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

var fsTestArray2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := [...]float64{1, 3: 3.4, 5}
	x[1] = 217
	println("x:", x, "x[1]:", x[1])
`)

func TestArray2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestArray2, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(30) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestMap = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := map[string]float64{"Hello": 1, "xsw": 3.4}
	println("x:", x)
`)

func TestMap(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestMap, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(24) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestMapLit = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := {"Hello": 1, "xsw": 3.4}
	println("x:", x)
`)

func TestMapLit(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestMapLit, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(24) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestMapLit2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := {"Hello": 1, "xsw": "3.4"}
	println("x:", x)

	println("empty map:", {})
`)

func TestMapLit2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestMapLit2, "/foo", nil, 0)
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

var fsTestMapIdx = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := {"Hello": 1, "xsw": "3.4"}
	y := {1: "glang", 5: "Hi"}
	i := 1
	q := "Q"
	key := "xsw"
	x["xsw"], y[i] = 3.1415926, q
	println("x:", x, "y:", y)
	println("x[key]:", x[key], "y[1]:", y[1])
`)

func TestMapIdx(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestMapIdx, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(26) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsSliceLit = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := [1, 3.4]
	println("x:", x)

	y := [1]
	println("y:", y)

	z := [1+2i, "xsw"]
	println("z:", z)

	println("empty slice:", [])
`)

func TestSliceLit(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsSliceLit, "/foo", nil, 0)
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

var fsSliceIdx = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := [1, 3.4, 17]
	n, m := 1, uint16(0)
	x[1] = 32.7
	x[m] = 36.86
	println("x:", x[2], x[m], x[n])
`)

func TestSliceIdx(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsSliceIdx, "/foo", nil, 0)
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

var fsListComprehension = asttest.NewSingleFileFS("/foo", "bar.gop", `
	y := [i+x for i, x <- [1, 2, 3, 4]]
	println("y:", y)
`)

func TestListComprehension(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsListComprehension, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(13) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsListComprehension2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	y := [i+x for i, x <- {3: 1, 5: 2, 7: 3, 11: 4}]
	println("y:", y)
`)

func TestListComprehension2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsListComprehension2, "/foo", nil, 0)
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

var fsListComprehensionFilter = asttest.NewSingleFileFS("/foo", "bar.gop", `
	y := [i+x for i, x <- {3: 1, 5: 2, 7: 3, 11: 4}, x % 2 == 1]
	println("y:", y)
`)

func TestListComprehensionFilter(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsListComprehensionFilter, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(10) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsMapComprehension = asttest.NewSingleFileFS("/foo", "bar.gop", `
	y := {x: i for i, x <- [3, 5, 7, 11, 13]}
	println("y:", y)
`)

func TestMapComprehension(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsMapComprehension, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(30) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsMapComprehensionFilter = asttest.NewSingleFileFS("/foo", "bar.gop", `
	y := {x: i for i, x <- [3, 5, 7, 11, 13], i % 2 == 1}
	println("y:", y)
`)

func TestMapComprehensionFilter(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsMapComprehensionFilter, "/foo", nil, 0)
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

var fsMapComprehension2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	y := {v: k for k, v <- {"Hello": "xsw", "Hi": "glang"}}
	println("y:", y)
`)

func TestMapComprehension2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsMapComprehension2, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(27) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsMapComprehension3 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	println({x: i for i, x <- [3, 5, 7, 11, 13]})
	println({x: i for i, x <- [3, 5, 7, 11, 13]})
`)

func TestMapComprehension3(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsMapComprehension3, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(27) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsMapComprehension4 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	arr := [1, 2, 3, 4, 5, 6]
	x := [[a, b] for a <- arr, a < b for b <- arr, b > 2]
	println("x:", x)
`)

func TestMapComprehension4(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsMapComprehension4, "/foo", nil, 0)
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
	if v := ctx.Get(-2); v != int(89) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestErrWrapExpr = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := println("Hello qiniu")!
	x
`)

func TestErrWrapExpr(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestErrWrapExpr, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Get(-1); v != int(12) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestErrWrapExpr2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
import (
	"strconv"
)

func add(x, y string) (int, error) {
	return strconv.Atoi(x)? + strconv.Atoi(y)?, nil
}

x := add("100", "23")!
x
`)

func TestErrWrapExpr2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestErrWrapExpr2, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Get(-1); v != int(123) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestRational = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := 3/4r + 5/7r
	x
`)

func TestRational(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestRational, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Get(-1); v.(*big.Rat).Cmp(big.NewRat(41, 28)) != 0 {
		t.Fatal("v:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestRational2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	a := 3/4r
	x := a + 5/7r
	x
`)

func TestRational2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestRational2, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Get(-1); v.(*big.Rat).Cmp(big.NewRat(41, 28)) != 0 {
		t.Fatal("v:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestRational3 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	y := 3.14159265358979323846264338327950288419716939937510582097494459r
	y *= 2
	y
`)

func TestRational3(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestRational3, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	y, _ := new(big.Float).SetString(
		"3.14159265358979323846264338327950288419716939937510582097494459")
	y.Mul(y, big.NewFloat(2))
	if v := ctx.Get(-1); v.(*big.Float).Cmp(y) != 0 {
		t.Fatal("v:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestRational4 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	a := 3/4r
	b := 5/7r
	if a > b {
		a = a + 1
	}
	a
`)

func TestRational4(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestRational4, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Get(-1); v.(*big.Rat).Cmp(big.NewRat(7, 4)) != 0 {
		t.Fatal("v:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestRational5 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := 1/3r + 1r*2r
	x
`)

func TestRational5(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestRational5, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Get(-1); v.(*big.Rat).Cmp(big.NewRat(7, 3)) != 0 {
		t.Fatal("v:", v)
	}
}

// -----------------------------------------------------------------------------

var testDeleteClauses = map[string]testData{
	"delete_int_key": {`
					m:={1:1,2:2}
					delete(m,1)
					println(m)
					delete(m,3)
					println(m)
					delete(m,2)
					println(m)
					`, []string{"map[2:2]", "map[2:2]", "map[]"}},
	"delete_string_key": {`
					m:={"hello":1,"Go+":2}
					delete(m,"hello")
					println(m)
					delete(m,"hi")
					println(m)
					delete(m,"Go+")
					println(m)
					`, []string{"map[Go+:2]", "map[Go+:2]", "map[]"}},
	"delete_var_string_key": {`
					m:={"hello":1,"Go+":2}
					delete(m,"hello")
					println(m)
					a:="hi"
					delete(m,a)
					println(m)
					arr:=["Go+"]
					delete(m,arr[0])
					println(m)
					`, []string{"map[Go+:2]", "map[Go+:2]", "map[]"}},
	"delete_var_map_string_key": {`
					ma:=[{"hello":1,"Go+":2}]
					delete(ma[0],"hello")
					println(ma[0])
					a:="hi"
					delete(ma[0],a)
					println(ma[0])
					arr:=["Go+"]
					delete(ma[0],arr[0])
					println(ma[0])
					`, []string{"map[Go+:2]", "map[Go+:2]", "map[]"}},
	"delete_no_key_panic": {`
					m:={"hello":1,"Go+":2}
					delete(m)
					`, []string{"_panic"}},
	"delete_multi_key_panic": {`
					m:={"hello":1,"Go+":2}
					delete(m,"hi","hi")
					`, []string{"_panic"}},
	"delete_not_map_panic": {`
					m:=[1,2,3]
					delete(m,1)
					`, []string{"_panic"}},
}

func TestDelete(t *testing.T) {
	for name, clause := range testDeleteClauses {
		testSingleStmt(name, t, asttest.NewSingleFileFS("/foo", "bar.gop", clause.clause), clause.wants)
	}
}

// -----------------------------------------------------------------------------
