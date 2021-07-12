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

package cl_test

import (
	"bytes"
	"testing"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/parser/parsertest"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

var (
	gblFset     *token.FileSet
	gblLoadPkgs gox.LoadPkgsFunc
)

func init() {
	gox.SetDebug(true)
	gblFset = token.NewFileSet()
	gblLoadPkgs = gox.NewLoadPkgsCached(nil)
}

func gopClTest(t *testing.T, gopcode, expected string) {
	fset := gblFset
	fs := parsertest.NewSingleFileFS("/foo", "bar.gop", gopcode)
	pkgs, err := parser.ParseFSDir(fset, fs, "/foo", nil, 0)
	if err != nil {
		t.Fatal("ParseFSDir:", err)
	}
	bar := pkgs["main"]
	pkg, err := cl.NewPackage("", bar, fset, &cl.Config{
		LoadPkgs: gblLoadPkgs,
	})
	if err != nil {
		t.Fatal("NewPackage:", err)
	}
	var b bytes.Buffer
	err = gox.WriteTo(&b, pkg)
	if err != nil {
		t.Fatal("gox.WriteTo failed:", err)
	}
	result := b.String()
	if result != expected {
		t.Fatalf("\nResult:\n%s\nExpected:\n%s\n", result, expected)
	}
}

func TestCompositeLit(t *testing.T) {
	gopClTest(t, `
x := []float64{1, 3.4, 5}
y := map[string]int{"Hello": 1, "Go+": 5}
z := [...]int{1, 3, 5}
a := {"Hello": 1, "Go+": 5.1}
`, `package main

func main() {
	x := []float64{1, 3.4, 5}
	y := map[string]int{"Hello": 1, "Go+": 5}
	z := [3]int{1, 3, 5}
	a := map[string]float64{"Hello": 1, "Go+": 5.1}
}
`)
}

func TestSliceLit(t *testing.T) {
	gopClTest(t, `
x := [1, 3.4, 5]
y := [1]
z := []
`, `package main

func main() {
	x := []float64{1, 3.4, 5}
	y := []int{1}
	z := []interface {
	}{}
}
`)
}

func TestKeyValModeLit(t *testing.T) {
	gopClTest(t, `
a := [...]float64{1, 3: 3.4, 5}
b := []float64{2: 1.2, 3, 6: 4.5}
`, `package main

func main() {
	a := [5]float64{1, 3: 3.4, 5}
	b := []float64{2: 1.2, 3, 6: 4.5}
}
`)
}

func TestFor(t *testing.T) {
	gopClTest(t, `
a := [1, 3.4, 5]
for i := 0; i < 3; i=i+1 {
	println(i)
}
`, `package main

import fmt "fmt"

func main() {
	a := []float64{1, 3.4, 5}
	for i := 0; i < 3; i = i + 1 {
		fmt.Println(i)
	}
}
`)
}

func TestRangeStmt(t *testing.T) {
	gopClTest(t, `
a := [1, 3.4, 5]
for _, x := range a {
	println(x)
}
`, `package main

import fmt "fmt"

func main() {
	a := []float64{1, 3.4, 5}
	for _, x := range a {
		fmt.Println(x)
	}
}
`)
}

func TestListComprehension(t *testing.T) {
	gopClTest(t, `
a := [1, 3.4, 5]
b := [x*x for x <- a]
`, `package main

func main() {
	a := []float64{1, 3.4, 5}
	b := func() (_gop_ret []float64) {
		for _, x := range a {
			_gop_ret = append(_gop_ret, x*x)
		}
		return
	}()
}
`)
}

func TestListComprehensionMultiLevel(t *testing.T) {
	gopClTest(t, `
arr := [1, 2, 3, 4.1, 5, 6]
x := [[a, b] for a <- arr, a < b for b <- arr, b > 2]
println("x:", x)
`, `package main

import fmt "fmt"

func main() {
	arr := []float64{1, 2, 3, 4.1, 5, 6}
	x := func() (_gop_ret [][]float64) {
		for _, b := range arr {
			if b > 2 {
				for _, a := range arr {
					if a < b {
						_gop_ret = append(_gop_ret, []float64{a, b})
					}
				}
			}
		}
		return
	}()
	fmt.Println("x:", x)
}
`)
}

func TestImport(t *testing.T) {
	gopClTest(t, `import "fmt"

func main() {
	fmt.Println("Hi")
}`, `package main

import fmt "fmt"

func main() {
	fmt.Println("Hi")
}
`)
}

func TestImportUnused(t *testing.T) {
	gopClTest(t, `import "fmt"

func main() {
}`, `package main

func main() {
}
`)
}

func TestAnonymousImport(t *testing.T) {
	gopClTest(t, `println("Hello")
printf("Hello Go+\n")
`, `package main

import fmt "fmt"

func main() {
	fmt.Println("Hello")
	fmt.Printf("Hello Go+\n")
}
`)
}

func TestMemberVal(t *testing.T) {
	gopClTest(t, `import "strings"

x := strings.NewReplacer("?", "!").Replace("hello, world???")
println("x:", x)
`, `package main

import (
	strings "strings"
	fmt "fmt"
)

func main() {
	x := strings.NewReplacer("?", "!").Replace("hello, world???")
	fmt.Println("x:", x)
}
`)
}

func TestIf(t *testing.T) {
	gopClTest(t, `x := 0
if t := false; t {
	x = 3
} else {
	x = 5
}
println("x:", x)
`, `package main

import fmt "fmt"

func main() {
	x := 0
	if t := false; t {
		x = 3
	} else {
		x = 5
	}
	fmt.Println("x:", x)
}
`)
}

func TestSwitch(t *testing.T) {
	gopClTest(t, `x := 0
switch s := "Hello"; s {
default:
	x = 7
case "world", "hi":
	x = 5
case "xsw":
	x = 3
}
println("x:", x)

v := "Hello"
switch {
case v == "xsw":
	x = 3
case v == "hi", v == "world":
	x = 9
default:
	x = 11
}
println("x:", x)
`, `package main

import fmt "fmt"

func main() {
	x := 0
	switch s := "Hello"; s {
	default:
		x = 7
	case "world", "hi":
		x = 5
	case "xsw":
		x = 3
	}
	fmt.Println("x:", x)
	v := "Hello"
	switch {
	case v == "xsw":
		x = 3
	case v == "hi", v == "world":
		x = 9
	default:
		x = 11
	}
	fmt.Println("x:", x)
}
`)
}

func TestSwitchFallthrough(t *testing.T) {
	gopClTest(t, `v := "Hello"
switch v {
case "Hello":
	println(v)
	fallthrough
case "hi":
	println(v)
	fallthrough
default:
	println(v)
}
`, `package main

import fmt "fmt"

func main() {
	v := "Hello"
	switch v {
	case "Hello":
		fmt.Println(v)
		fallthrough
	case "hi":
		fmt.Println(v)
		fallthrough
	default:
		fmt.Println(v)
	}
}
`)
}

func TestReturn(t *testing.T) {
	gopClTest(t, `
func foo(format string, args ...interface{}) (int, error) {
	return println(format, args...)
}

func main() {
}
`, `package main

import fmt "fmt"

func foo(format string, args ...interface {
}) ( int,  error) {
	return fmt.Println(format, args...)
}
func main() {
}
`)
}

func TestReturnExpr(t *testing.T) {
	gopClTest(t, `
func foo(format string, args ...interface{}) (int, error) {
	return 0, nil
}

func main() {
}
`, `package main

func foo(format string, args ...interface {
}) ( int,  error) {
	return 0, nil
}
func main() {
}
`)
}

func TestClosure(t *testing.T) {
	gopClTest(t, `import "fmt"

func(v string) {
	fmt.Println(v)
}("Hello")
`, `package main

import fmt "fmt"

func main() {
	func(v string) {
		fmt.Println(v)
	}("Hello")
}
`)
}

func TestFunc(t *testing.T) {
	gopClTest(t, `func foo(format string, args ...interface{}) {
}

func main() {
}`, `package main

func foo(format string, args ...interface {
}) {
}
func main() {
}
`)
}

func TestUnnamedMainFunc(t *testing.T) {
	gopClTest(t, `i := 1`, `package main

func main() {
	i := 1
}
`)
}

func TestFuncAsParam(t *testing.T) {
	gopClTest(t, `import "fmt"

func bar(foo func(string, ...interface{}) (int, error)) {
	foo("Hello, %v!\n", "Go+")
}

bar(fmt.Printf)
`, `package main

import fmt "fmt"

func bar(foo func( string,  ...interface {
}) ( int,  error)) {
	foo("Hello, %v!\n", "Go+")
}
func main() {
	bar(fmt.Printf)
}
`)
}

func TestFuncAsParam2(t *testing.T) {
	gopClTest(t, `import (
	"fmt"
	"strings"
)

func foo(x string) string {
	return strings.NewReplacer("?", "!").Replace(x)
}

func printf(format string, args ...interface{}) (n int, err error) {
	n, err = fmt.Printf(format, args...)
	return
}

func bar(foo func(string, ...interface{}) (int, error)) {
	foo("Hello, %v!\n", "Go+")
}

bar(printf)
fmt.Println(foo("Hello, world???"))
fmt.Println(printf("Hello, %v\n", "Go+"))
`, `package main

import (
	fmt "fmt"
	strings "strings"
)

func foo(x string) ( string) {
	return strings.NewReplacer("?", "!").Replace(x)
}
func printf(format string, args ...interface {
}) (n int, err error) {
	n, err = fmt.Printf(format, args...)
	return
}
func bar(foo func( string,  ...interface {
}) ( int,  error)) {
	foo("Hello, %v!\n", "Go+")
}
func main() {
	bar(printf)
	fmt.Println(foo("Hello, world???"))
	fmt.Println(printf("Hello, %v\n", "Go+"))
}
`)
}

func TestFuncCall(t *testing.T) {
	gopClTest(t, `import "fmt"

fmt.Println("Hello")`, `package main

import fmt "fmt"

func main() {
	fmt.Println("Hello")
}
`)
}

func TestFuncCallEllipsis(t *testing.T) {
	gopClTest(t, `import "fmt"

func foo(args ...interface{}) {
	fmt.Println(args...)
}

func main() {
}`, `package main

import fmt "fmt"

func foo(args ...interface {
}) {
	fmt.Println(args...)
}
func main() {
}
`)
}

func TestFuncCallCodeOrder(t *testing.T) {
	gopClTest(t, `import "fmt"

func main() {
	foo("Hello", 123)
}

func foo(args ...interface{}) {
	fmt.Println(args...)
}
`, `package main

import fmt "fmt"

func main() {
	foo("Hello", 123)
}
func foo(args ...interface {
}) {
	fmt.Println(args...)
}
`)
}

func TestInterfaceMethods(t *testing.T) {
	gopClTest(t, `package main

func foo(v ...interface { Bar() }) {
}

func main() {
}`, `package main

func foo(v ...interface {
	Bar()
}) {
}
func main() {
}
`)
}

func TestAssignUnderscore(t *testing.T) {
	gopClTest(t, `import log "fmt"

_, err := log.Println("Hello")
`, `package main

import fmt "fmt"

func main() {
	_, err := fmt.Println("Hello")
}
`)
}

func TestOperator(t *testing.T) {
	gopClTest(t, `
a := "Hi"
b := a + "!"
c := 13
d := -c
`, `package main

func main() {
	a := "Hi"
	b := a + "!"
	c := 13
	d := -c
}
`)
}
