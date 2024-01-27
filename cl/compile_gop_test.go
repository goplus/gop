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
