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
	"os"
	"testing"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/parser/parsertest"
	"github.com/goplus/gop/scanner"
)

func codeErrorTest(t *testing.T, msg, src string) {
	fs := parsertest.NewSingleFileFS("/foo", "bar.gop", src)
	pkgs, err := parser.ParseFSDir(gblFset, fs, "/foo", nil, 0)
	if err != nil {
		scanner.PrintError(os.Stderr, err)
		t.Fatal("parser.ParseFSDir failed")
	}
	conf := *baseConf.Ensure()
	conf.NoFileLine = false
	conf.WorkingDir = "/foo"
	conf.TargetDir = "/foo"
	bar := pkgs["main"]
	_, err = cl.NewPackage("", bar, &conf)
	if err == nil {
		t.Fatal("no error?")
	}
	if ret := err.Error(); ret != msg {
		t.Fatalf("\nError: \"%s\"\nExpected: \"%s\"\n", ret, msg)
	}
}

func TestErrErrWrap(t *testing.T) {
	codeErrorTest(t,
		"./bar.gop:2:2: undefined: a", `func main() {
	a!
}
`)
}

func TestErrVar(t *testing.T) {
	codeErrorTest(t,
		"./bar.gop:6:5: assignment mismatch: 1 variables but fmt.Println returns 2 values", `import "fmt"

func main() {
}

var a = fmt.Println(1)
`)
	codeErrorTest(t,
		"./bar.gop:4:5: assignment mismatch: 1 variables but 2 values", `func main() {
}

var a = 1, 2
`)
	codeErrorTest(t,
		"./bar.gop:2:2: undefined: foo", `func main() {
	foo.x = 1
}
`)
	codeErrorTest(t,
		"./bar.gop:2:2: use of builtin len not in function call", `func main() {
	len.x = 1
}
`)
	codeErrorTest(t,
		"./bar.gop:2:10: undefined: foo", `func main() {
	println(foo.x)
}
`)
	codeErrorTest(t,
		"./bar.gop:2:9: use of builtin len not in function call", `func main() {
println(len.x)
}
`)
	codeErrorTest(t,
		"./bar.gop:2:10: undefined: foo", `func main() {
	println(foo)
}
`)
	codeErrorTest(t,
		"./bar.gop:3:20: use of builtin len not in function call", `package main

func foo(v map[int]len) {
}
`)
	codeErrorTest(t,
		"./bar.gop:5:20: bar is not a type", `package main

var bar = 1

func foo(v map[int]bar) {
}
`)
	codeErrorTest(t,
		"./bar.gop:2:6: use of builtin len not in function call", `func main() {
	new(len)
}
`)
	codeErrorTest(t,
		"./bar.gop:2:2: undefined: foo", `func main() {
	foo = 1
}
`)
	codeErrorTest(t,
		"./bar.gop:2:9: cannot use _ as value", `func main() {
	foo := _
}
`)
	codeErrorTest(t,
		"./bar.gop:2:9: use of builtin len not in function call", `func main() {
	foo := len
}
`)
	codeErrorTest(t,
		"./bar.gop:2:2: use of builtin println not in function call", `func main() {
	println = "hello"
}
`)
}

func TestErrImport(t *testing.T) {
	codeErrorTest(t,
		"./bar.gop:5:2: cannot refer to unexported name os.undefined", `
import "os"

func foo() {
	os.undefined
}`)
	codeErrorTest(t,
		"./bar.gop:5:2: undefined: os.UndefinedObject", `
import "os"

func foo() {
	os.UndefinedObject
}`)
	codeErrorTest(t,
		"./bar.gop:2:13: undefined: testing", `
func foo(t *testing.T) {
}`)
	codeErrorTest(t,
		"./bar.gop:4:12: testing.Verbose is not a type", `
import "testing"

func foo(t testing.Verbose) {
}`)
}

func TestErrConst(t *testing.T) {
	codeErrorTest(t,
		"./bar.gop:3:7: a redeclared in this block\n\tprevious declaration at ./bar.gop:2:5", `
var a int
const a = 1
`)
	codeErrorTest(t,
		"./bar.gop:4:2: missing value in const declaration", `
const (
	a = iota
	b, c
)
`)
}

func TestErrNewVar(t *testing.T) {
	codeErrorTest(t,
		"./bar.gop:3:5: a redeclared in this block\n\tprevious declaration at ./bar.gop:2:5", `
var a int
var a string
`)
}

func TestErrDefineVar(t *testing.T) {
	codeErrorTest(t, "./bar.gop:3:1: no new variables on left side of :=\n"+
		"./bar.gop:3:6: cannot use \"Hi\" (type untyped string) as type int in assignment", `
a := 1
a := "Hi"
`)
}

func TestErrAssign(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:8:1: assignment mismatch: 1 variables but bar returns 2 values`, `

func bar() (n int, err error) {
	return
}

x := 1
x = bar()
`)
	codeErrorTest(t,
		`./bar.gop:4:1: assignment mismatch: 1 variables but 2 values`, `

x := 1
x = 1, "Hi"
`)
}

func TestErrReturn(t *testing.T) {
	codeErrorTest(t,
		"./bar.gop:4:2: too few arguments to return\n\thave (untyped int)\n\twant (int, error)", `

func foo() (int, error) {
	return 1
}
`)
	codeErrorTest(t,
		"./bar.gop:4:2: too many arguments to return\n\thave (untyped int, untyped int, untyped string)\n\twant (int, error)", `

func foo() (int, error) {
	return 1, 2, "Hi"
}
`)
	codeErrorTest(t,
		`./bar.gop:4:12: cannot use "Hi" (type untyped string) as type error in return argument`, `

func foo() (int, error) {
	return 1, "Hi"
}
`)
	codeErrorTest(t,
		"./bar.gop:8:2: too few arguments to return\n\thave (byte)\n\twant (int, error)", `

func bar() (v byte) {
	return
}

func foo() (int, error) {
	return bar()
}
`)
	codeErrorTest(t,
		"./bar.gop:8:2: too many arguments to return\n\thave (n int, err error)\n\twant (v byte)", `

func bar() (n int, err error) {
	return
}

func foo() (v byte) {
	return bar()
}
`)
	codeErrorTest(t,
		`./bar.gop:8:2: cannot use byte value as type error in return argument`, `

func bar() (n int, v byte) {
	return
}

func foo() (int, error) {
	return bar()
}
`)
	codeErrorTest(t,
		"./bar.gop:4:2: not enough arguments to return\n\thave ()\n\twant (byte)", `

func foo() byte {
	return
}
`)
}

func TestErrForRange(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:4:8: cannot assign type string to a (type int) in range`, `
a := 1
var b []string
for _, a = range b {
}
`)
}

func TestErrInitFunc(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:2:1: func init must have no arguments and no return values`, `
func init(v byte) {
}
`)
}

func TestErrRecv(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:5:9: invalid receiver type a (a is a pointer type)`, `

type a *int

func (p a) foo() {
}
`)
	codeErrorTest(t,
		`./bar.gop:2:9: invalid receiver type error (error is an interface type)`, `
func (p error) foo() {
}
`)
	codeErrorTest(t,
		`./bar.gop:2:9: invalid receiver type []byte ([]byte is not a defined type)`, `
func (p []byte) foo() {
}
`)
	codeErrorTest(t,
		`./bar.gop:2:10: invalid receiver type []byte ([]byte is not a defined type)`, `
func (p *[]byte) foo() {
}
`)
}

func TestErrStructLit(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:3:39: too many values in struct{x int; y string}{...}`, `
x := 1
a := struct{x int; y string}{1, "Hi", 2}
`)
	codeErrorTest(t,
		`./bar.gop:3:30: too few values in struct{x int; y string}{...}`, `
x := 1
a := struct{x int; y string}{1}
`)
	codeErrorTest(t,
		`./bar.gop:3:33: cannot use x (type int) as type string in value of field y`, `
x := 1
a := struct{x int; y string}{1, x}
`)
}

func TestErrArray(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:3:8: non-constant array bound n`, `
var n int
var a [n]int
`)
}

func TestErrArrayLit(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:3:14: cannot use a as index which must be non-negative integer constant`,
		`
a := "Hi"
b := [10]int{a: 1}
`)
	codeErrorTest(t,
		`./bar.gop:3:20: array index 10 out of bounds [0:10]`,
		`
a := "Hi"
b := [10]int{9: 1, 3}
`)
	codeErrorTest(t,
		`./bar.gop:3:16: array index 1 out of bounds [0:1]`,
		`
a := "Hi"
b := [1]int{1, 2}
`)
	codeErrorTest(t,
		`./bar.gop:3:14: array index 12 (value 12) out of bounds [0:10]`,
		`
a := "Hi"
b := [10]int{12: 2}
`)
	codeErrorTest(t,
		`./bar.gop:3:14: cannot use a+"!" (type string) as type int in array literal`,
		`
a := "Hi"
b := [10]int{a+"!"}
`)
	codeErrorTest(t,
		`./bar.gop:3:17: cannot use a (type string) as type int in array literal`,
		`
a := "Hi"
b := [10]int{2: a}
`)
}

func TestErrSliceLit(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:3:12: cannot use a as index which must be non-negative integer constant`,
		`
a := "Hi"
b := []int{a: 1}
`)
	codeErrorTest(t,
		`./bar.gop:3:12: cannot use a (type string) as type int in slice literal`,
		`
a := "Hi"
b := []int{a}
`)
	codeErrorTest(t,
		`./bar.gop:3:15: cannot use a (type string) as type int in slice literal`,
		`
a := "Hi"
b := []int{2: a}
`)
}

func TestErrMapLit(t *testing.T) {
	codeErrorTest(t, // TODO: first column need correct
		`./bar.gop:2:34: cannot use 1+2 (type untyped int) as type string in map key
./bar.gop:3:27: cannot use "Go" + "+" (type untyped string) as type int in map value`,
		`
a := map[string]int{1+2: 2}
b := map[string]int{"Hi": "Go" + "+"}
`)
}

func TestErrSlice(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:4:6: cannot slice a (type *byte)`,
		`
var a *byte
x := 1
b := a[x:2]
`)
	codeErrorTest(t,
		`./bar.gop:3:6: cannot slice a (type bool)`,
		`
a := true
b := a[1:2]
`)
	codeErrorTest(t,
		`./bar.gop:3:6: invalid operation a[1:2:5] (3-index slice of string)`,
		`
a := "Hi"
b := a[1:2:5]
`)
}

func TestErrIndex(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:3:10: assignment mismatch: 2 variables but 1 values`,
		`
a := "Hi"
b, ok := a[1]
`)
	codeErrorTest(t,
		`./bar.gop:3:6: invalid operation: a[1] (type bool does not support indexing)`,
		`
a := true
b := a[1]
`)
}

func TestErrIndexRef(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:3:1: cannot assign to a[1] (strings are immutable)`,
		`
a := "Hi"
a[1] = 'e'
`)
}

func TestErrStar(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:3:2: invalid indirect of a (type string)`,
		`
a := "Hi"
*a = 'e'
`)
	codeErrorTest(t,
		`./bar.gop:3:7: invalid indirect of a (type string)`,
		`
a := "Hi"
b := *a
`)
}

func TestErrMember(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:3:6: a.x undefined (type string has no field or method x)`,
		`
a := "Hello"
b := a.x
`)
}

func TestErrMemberRef(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:3:1: a.x undefined (type string has no field or method x)`,
		`
a := "Hello"
a.x = 1
`)
	codeErrorTest(t,
		`./bar.gop:5:1: a.x undefined (type aaa has no field or method x)`,
		`
type aaa byte

a := aaa(0)
a.x = 1
`)
	codeErrorTest(t,
		`./bar.gop:5:1: a.z undefined (type aaa has no field or method z)`,
		`
type aaa struct {x int; y string}

a := aaa{}
a.z = 1
`)
	codeErrorTest(t,
		`./bar.gop:3:1: a.z undefined (type struct{x int; y string} has no field or method z)`,
		`
a := struct{x int; y string}{}
a.z = 1
`)
}

func TestErrLabel(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:4:1: label foo already defined at ./bar.gop:2:1
./bar.gop:2:1: label foo defined and not used`,
		`x := 1
foo:
	i := 1
foo:
	i++
`)
	codeErrorTest(t,
		`./bar.gop:2:6: label foo is not defined`,
		`x := 1
goto foo`)
}
