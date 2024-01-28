/*
 * Copyright (c) 2024 The GoPlus Authors (goplus.org). All rights reserved.
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

package cl_test

import (
	"testing"
)

func TestOverloadOp(t *testing.T) {
	gopClTest(t, `
type foo struct {
}

func (a *foo) + (b *foo) *foo {
	println("a + b")
	return &foo{}
}

func (a foo) - (b foo) foo {
	println("a - b")
	return foo{}
}

func -(a foo) {
	println("-a")
}

func ++(a foo) {
	println("a++")
}

func (a foo) != (b foo) bool{
	println("a!=b")
	return true
}

var a, b foo
var c = a - b
var d = -a       // TODO: -a have no return value!
var e = a!=b
`, `package main

import "fmt"

type foo struct {
}

func (a *foo) Gop_Add(b *foo) *foo {
	fmt.Println("a + b")
	return &foo{}
}
func (a foo) Gop_Sub(b foo) foo {
	fmt.Println("a - b")
	return foo{}
}
func (a foo) Gop_NE(b foo) bool {
	fmt.Println("a!=b")
	return true
}
func (a foo) Gop_Neg() {
	fmt.Println("-a")
}
func (a foo) Gop_Inc() {
	fmt.Println("a++")
}

var a, b foo
var c = (foo).Gop_Sub(a, b)
var d = a.Gop_Neg()
var e = (foo).Gop_NE(a, b)
`)
}

func TestOverloadOp2(t *testing.T) {
	gopClTest(t, `
type foo struct {
}

func (a foo) mulInt(b int) (ret foo) {
	return
}

func (a foo) mulFoo(b foo) (ret foo) {
	return
}

func intMulFoo(a int, b foo) (ret foo) {
	return
}

func (foo).* = (
	(foo).mulInt
	(foo).mulFoo
	intMulFoo
)

var a, b foo

println a * 10
println a * b
println 10 * a
`, `package main

import "fmt"

type foo struct {
}

const Gopo__foo__Gop_Mul = ".mulInt,.mulFoo,intMulFoo"

func (a foo) mulInt(b int) (ret foo) {
	return
}
func (a foo) mulFoo(b foo) (ret foo) {
	return
}
func intMulFoo(a int, b foo) (ret foo) {
	return
}

var a, b foo

func main() {
	fmt.Println((foo).mulInt(a, 10))
	fmt.Println((foo).mulFoo(a, b))
	fmt.Println(intMulFoo(10, a))
}
`)
}

func TestOverloadMethod(t *testing.T) {
	gopClTest(t, `
type foo struct {
}

func (a *foo) mulInt(b int) *foo {
	println "mulInt"
	return a
}

func (a *foo) mulFoo(b *foo) *foo {
	println "mulFoo"
	return a
}

func (foo).mul = (
	(foo).mulInt
	(foo).mulFoo
)

var a, b foo
var c = a.mul(100)
var d = a.mul(c)
`, `package main

import "fmt"

type foo struct {
}

const Gopo_foo_mul = ".mulInt,.mulFoo"

func (a *foo) mulInt(b int) *foo {
	fmt.Println("mulInt")
	return a
}
func (a *foo) mulFoo(b *foo) *foo {
	fmt.Println("mulFoo")
	return a
}

var a, b foo
var c = a.mulInt(100)
var d = a.mulFoo(c)
`)
}

func TestOverloadFunc(t *testing.T) {
	gopClTest(t, `
func add = (
	func(a, b int) int {
		return a + b
	}
	func(a, b string) string {
		return a + b
	}
)

println add(100, 7)
println add("Hello", "World")
`, `package main

import "fmt"

func add__0(a int, b int) int {
	return a + b
}
func add__1(a string, b string) string {
	return a + b
}
func main() {
	fmt.Println(add__0(100, 7))
	fmt.Println(add__1("Hello", "World"))
}
`)
}

func TestOverloadFunc2(t *testing.T) {
	gopClTest(t, `
func mulInt(a, b int) int {
	return a * b
}

func mulFloat(a, b float64) float64 {
	return a * b
}

func mul = (
	mulInt
	mulFloat
)

println mul(100, 7)
println mul(1.2, 3.14)
`, `package main

import "fmt"

const Gopo_mul = "mulInt,mulFloat"

func mulInt(a int, b int) int {
	return a * b
}
func mulFloat(a float64, b float64) float64 {
	return a * b
}
func main() {
	fmt.Println(mulInt(100, 7))
	fmt.Println(mulFloat(1.2, 3.14))
}
`)
}

func TestStringLitBasic(t *testing.T) {
	gopClTest(t, `echo "$$"`, `package main

import "fmt"

func main() {
	fmt.Println("$")
}
`)
}

func TestStringLitVar(t *testing.T) {
	gopClTest(t, `
x := 1
println "Hi, " + "a${x}b"`, `package main

import (
	"fmt"
	"strconv"
)

func main() {
	x := 1
	fmt.Println("Hi, " + ("a" + strconv.Itoa(x) + "b"))
}
`)
}

func TestFileOpen(t *testing.T) {
	gopClTest(t, `
for line <- open("foo.txt")! {
	println line
}
`, `package main

import (
	"fmt"
	"github.com/goplus/gop/builtin/iox"
	"github.com/qiniu/x/errors"
	"os"
)

func main() {
	for _gop_it := iox.EnumLines(func() (_gop_ret *os.File) {
		var _gop_err error
		_gop_ret, _gop_err = os.Open("foo.txt")
		if _gop_err != nil {
			_gop_err = errors.NewFrame(_gop_err, "open(\"foo.txt\")", "/foo/bar.gop", 2, "main.main")
			panic(_gop_err)
		}
		return
	}()); ; {
		var _gop_ok bool
		line, _gop_ok := _gop_it.Next()
		if !_gop_ok {
			break
		}
		fmt.Println(line)
	}
}
`)
}

func TestFileEnumLines(t *testing.T) {
	gopClTest(t, `
import "os"

for line <- os.Stdin {
	println line
}
`, `package main

import (
	"fmt"
	"github.com/goplus/gop/builtin/iox"
	"os"
)

func main() {
	for _gop_it := iox.EnumLines(os.Stdin); ; {
		var _gop_ok bool
		line, _gop_ok := _gop_it.Next()
		if !_gop_ok {
			break
		}
		fmt.Println(line)
	}
}
`)
}

func TestIoxLines(t *testing.T) {
	gopClTest(t, `
import "io"

var r io.Reader

for line <- lines(r) {
	println line
}
`, `package main

import (
	"fmt"
	"github.com/goplus/gop/builtin/iox"
	"io"
)

var r io.Reader

func main() {
	for _gop_it := iox.Lines(r).Gop_Enum(); ; {
		var _gop_ok bool
		line, _gop_ok := _gop_it.Next()
		if !_gop_ok {
			break
		}
		fmt.Println(line)
	}
}
`)
}

func TestMixedGo(t *testing.T) {
	gopMixedClTest(t, "main", `package main

import "strconv"

const n = 10

func f(v int) string {
	return strconv.Itoa(v)
}

type foo struct {
	v int
}

func (a foo) _() {
}

func (a foo) Str() string {
	return f(a.v)
}

func (a *foo) Bar() int {
	return 0
}

type foo2 = foo
type foo3 foo2
`, `
var a [n]int
var b string = f(n)
var c foo2
var d int = c.v
var e = foo3{}
var x string = c.str
`, `package main

var a [10]int
var b string = f(n)
var c foo
var d int = c.v
var e = foo3{}
var x string = c.Str()
`, true)
	gopMixedClTest(t, "main", `package main
type Point struct {
	X int
	Y int
}
`, `
type T struct{}
println(&T{},&Point{10,20})
`, `package main

import "fmt"

type T struct {
}

func main() {
	fmt.Println(&T{}, &Point{10, 20})
}
`, false)
}

func Test_RangeExpressionIf_Issue1243(t *testing.T) {
	gopClTest(t, `
for i <- :10, i%3 == 0 {
	println i
}`, `package main

import "fmt"

func main() {
	for i := 0; i < 10; i += 1 {
		if i%3 == 0 {
			fmt.Println(i)
		}
	}
}
`)
}
