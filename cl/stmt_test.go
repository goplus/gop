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
	"os"
	"testing"

	"github.com/qiniu/goplus/ast/asttest"
	"github.com/qiniu/goplus/gop"
	"github.com/qiniu/goplus/parser"
	"github.com/qiniu/goplus/token"

	exec "github.com/qiniu/goplus/exec/bytecode"
	libbuiltin "github.com/qiniu/goplus/lib/builtin"
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
	for x <- [1, 3, 5, 7, 11, 13, 17], x > 3 {
		sum += x
	}
	sum
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
	if v := ctx.Get(-1); v != 53 {
		t.Fatal("v:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestForPhraseStmt2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	sum := 0
	for x <- [1, 3, 5, 7, 11, 13, 17] {
		if x > 3 {
			sum += x
		}
	}
	sum
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
	if err != nil || !noExecCtx {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Get(-1); v != 53 {
		t.Fatal("v:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestForPhraseStmt3 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	fns := make([]func() int, 3)
	for i, x <- [3, 15, 777] {
		v := x
		fns[i] = func() int {
			return v
		}
	}
	println("values:", fns[0](), fns[1](), fns[2]())
`)

func TestForPhraseStmt3(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestForPhraseStmt3, "/foo", nil, 0)
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

type testData struct {
	clause string
	wants  []string
}

var testForRangeClauses = map[string]testData{
	"no_kv_range_list": {`sum:=0
					for range [1,3,5,7] {
						sum++
					}
					println(sum)
					`, []string{"4"}},
	"no_kv_range_map": {`sum:=0
					for range {1:1,2:2,3:3} {
						sum++
					}
					println(sum)
					`, []string{"3"}},
	"only_k_range_list": {`sum:=0
					for k :=range [1,3,5,7]{
						sum+=k
					}
					println(sum)
					`, []string{"6"}},
	"only_k_range_map": {`sum:=0
					for k :=range {1:1,2:4,3:8,4:16}{
						sum+=k
					}
					println(sum)
					`, []string{"10"}},
	"only_v_range_list": {`sum:=0
					for _,v :=range [1,3,5,7]{
						sum+=v
					}
					println(sum)
					`, []string{"16"}},
	"only_v_range_map": {`sum:=0
					for _,v :=range {1:1,2:4,3:8,4:16}{
						sum+=v
					}
					println(sum)
					`, []string{"29"}},
	"both_kv_range_list": {`sum:=0
					for k,v:=range [1,3,5,7]{
						// 0*1+1*3+2*5+3*7
						sum+=k*v
					}
					println(sum)
					`, []string{"34"}},
	"both_kv_range_map": {`sum:=0
					m:={1:2,2:4,3:8}
					for k,v:=range m { 
						//1*2+2*4+3*8=34
						sum+=k*v
					}
					println(sum)
					`, []string{"34"}},
	"both_kv_assign_simple_range": {` sum:=0
					k,v:=0,0
					for k,v=range [1,2,3,4,5]{
						sum+=k+v
					}
					println(k)
					println(v)
					println(sum)
					`, []string{"4", "5", "25"}},
	"both_kv_assign_range_list": {` sum:=0
					m:={1:2,2:4,3:8}
					arr:=[11,22]
					for m[1],m[2]=range arr{
					    sum+=m[1]+m[2]
					}
					println(m[1])
					println(m[2])
					println(sum)
					`, []string{"1", "22", "34"}},
	"both_kv_assign_range_map": {` sum:=0
					m:={3:8}
					arr:=[11,22]
					for arr[0],arr[1]=range m{
					    sum+=arr[0]+arr[1]
					}
					println(arr[0])
					println(arr[1])
					println(sum)
					`, []string{"3", "8", "11"}},
	"only_v_assign_range": {` sum:=0
					m:={3:8}
					arr:=[11,22]
					for _,arr[1]=range m{
					    sum+=arr[0]+arr[1]
					}
					println(arr[0])
					println(arr[1])
					println(sum)
					`, []string{"11", "8", "19"}},
	"only_k_assign_range": {` sum:=0
					m:={3:8}
					arr:=[11,22]
					for arr[0],_=range m{
					    sum+=arr[0]+arr[1]
					}
					println(arr[0])
					println(arr[1])
					println(sum)
					`, []string{"3", "22", "25"}},
	"none_kv_assign_range": {` sum:=0
					m:={3:8}
					arr:=[11,22]
					for _,_=range m{
					    sum+=arr[0]+arr[1]
					}
					println(arr[0])
					println(arr[1])
					println(sum)
					`, []string{"11", "22", "33"}},
}

func TestRangeStmt(t *testing.T) {
	for name, clause := range testForRangeClauses {
		testForRangeStmt(name, t, asttest.NewSingleFileFS("/foo", "bar.gop", clause.clause), clause.wants)
	}
}

// -----------------------------------------------------------------------------

var fsTestRangeStmt2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	sum := 0
	for _, x := range [1, 3, 5, 7, 11, 13, 17] {
		if x > 3 {
			sum += x
		}
	}
	sum
`)

func TestRangeStmt2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestRangeStmt2, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, noExecCtx, err := newPackage(b, bar, fset)
	if err != nil || !noExecCtx {
		t.Fatal("Compile failed:", err, "noExecCtx:", noExecCtx)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Get(-1); v != 53 {
		t.Fatal("v:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestForStmt = asttest.NewSingleFileFS("/foo", "bar.gop", `
	fns := make([]func() int, 3)
	arr := [3, 15, 777]
	sum := 0
	for i := 0; i < len(arr); i++ {
		v := arr[i]
		fns[i] = func() int {
			return v
		}
	}
	println("values:", fns[0](), fns[1](), fns[2]())
`)

func _TestForStmt(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestForStmt, "/foo", nil, 0)
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

var testNormalForClauses = map[string]testData{
	"for_with_init_cond_post": {`
					sum := 0
					arr := [1,3,5,7]
					for i := 0; i < len(arr); i++ {
						sum+=arr[i]
					}
					println(sum)
					`, []string{"16"}},
	"for_with_cond_post": {`
					sum := 0
					arr := [1,3,5,7]
					i := 0
					for ; i < len(arr); i+=2 {
						sum+=arr[i]
					}
					println(sum)
					`, []string{"6"}},
	"for_with_cond": {`
					arr := [1,3,5,7]
					i := 0
					sum := 0
					for ; i < len(arr) && i < 2; {
						sum+=arr[i]
						i++
					}
					println(sum)
					`, []string{"4"}},
	"for_with_init_cond": {`
					arr := [1,3,5,7]
					sum := 0
					for i:=0; i < len(arr); {
						sum+=arr[i]
						i++
					}
					println(sum)
					`, []string{"16"}},
}

func TestNormalForStmt(t *testing.T) {
	for name, clause := range testNormalForClauses {
		testForRangeStmt(name, t, asttest.NewSingleFileFS("/foo", "bar.gop", clause.clause), clause.wants)
	}
}

func testForRangeStmt(name string, t *testing.T, fs *asttest.MemFS, wants []string) {
	var results []string
	selfPrintln := func(arity int, p *gop.Context) {
		args := p.GetArgs(arity)
		results = append(results, fmt.Sprintln(args...))
		n, err := fmt.Println(args...)
		p.Ret(arity, n, err)
	}
	libbuiltin.I.RegisterFuncvs(libbuiltin.I.Funcv("println", fmt.Print, selfPrintln))

	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fs, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal(name+" : ParseFSDir failed:", err, len(pkgs))
	}
	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, noExecCtx, err := newPackage(b, bar, fset)
	if err != nil || !noExecCtx {
		t.Fatal(name+" :Compile failed:", err)
	}
	code := b.Resolve()
	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Get(-1); v != nil {
		t.Fatal(name+" :error:", v)
	}
	if v := ctx.Get(-2); v != len(results[len(results)-1]) {
		t.Fatal(name+" :n:", v)
	}
	if len(wants) != len(results) {
		t.Fatal(name+" exec fail , wants", wants, ",actually", results)
	}
	for i := 0; i < len(wants); i++ {
		if wants[i]+"\n" != results[i] {
			t.Fatal(name+" exec fail", i, "result wants", wants[i], ",actually", results[i])
		}
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

var testFallthroughClauses = map[string]testData{
	"switch_all_fallthrough": {`
					x:=0
					switch x {
					case 0:
						println(x)
						fallthrough
					case 1:
						x++
						println(x)
						fallthrough
					default:
						x=7
						println(x)
					}
					`, []string{"0", "1", "7"}},
	"switch_one_fallthrough": {`
					x:=0
					switch x {
					case 0,1,2:
						println(x)
						fallthrough
					case 3:
						x++
						println(x)
					default:
						x=7
						println(x)
					}
					`, []string{"0", "1"}},
	"switch__fallthrough": {`
					x:=0
					switch x {
					case 0:
						println(x)
						fallthrough
					case 1:
						x++
						println(x)
					default:
						x=7
						println(x)
					}
					`, []string{"0", "1"}},
	"switch_no_tag_fallthrough": {`
					x:=0
					switch {
					case x==0:
						println(x)
						fallthrough
					case x==1:
						x++
						println(x)
						fallthrough
					default:
						x=7
						println(x)
					}
					`, []string{"0", "1", "7"}},
	"switch_no_tag_one_fallthrough": {`
					x:=0
					switch x {
					case 0:
						println(x)
						fallthrough
					case 1:
						x++
						println(x)
					default:
						x=7
						println(x)
					}
					`, []string{"0", "1"}},
	"switch_fallthrough_panic": {`
					x:=0
					switch x {
					case 0:
						println(x)
						fallthrough
					default:
						x=7
						println(x)
					case 1:
						x++
						println(x)
					fallthrough
					}
					`, []string{"0", "1"}},
	"switch_fallthrough_out_panic": {`
					x:=0
					switch x {
					case 0:
						println(x)
						fallthrough
					default:
						x=7
						println(x)
					case 1:
						x++
						println(x)
					}
					fallthrough
					`, []string{"0", "1"}},
}

func TestFallthroughStmt(t *testing.T) {
	for name, clause := range testFallthroughClauses {
		func() {
			defer func() {
				if r := recover(); r != nil {
				}
			}()
			testForRangeStmt(name, t, asttest.NewSingleFileFS("/foo", "bar.gop", clause.clause), clause.wants)
		}()
	}
}
