/*
 * Copyright (c) 2022 The GoPlus Authors (goplus.org). All rights reserved.
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

package build_test

import (
	"fmt"
	"testing"

	"github.com/goplus/gop/x/build"
)

var (
	ctx = build.Default()
)

func gopClTest(t *testing.T, gopcode, expected string) {
	gopClTestEx(t, "main.gop", gopcode, expected)
}

func gopClTestEx(t *testing.T, filename string, gopcode, expected string) {
	data, err := ctx.BuildFile(filename, gopcode)
	if err != nil {
		t.Fatalf("build gop error: %v", err)
	}
	if string(data) != expected {
		fmt.Println("build gop error:")
		fmt.Println(string(data))
		t.Fail()
	}
}

func TestGop(t *testing.T) {
	gopClTest(t, `
println "Go+"
`, `package main

import fmt "fmt"

func main() {
//line main.gop:2
	fmt.Println("Go+")
}
`)
}

func TestGox(t *testing.T) {
	gopClTestEx(t, "Rect.gox", `
println "Go+"
`, `package main

import fmt "fmt"

type Rect struct {
}

func (this *Rect) Main() {
//line Rect.gox:2
	fmt.Println("Go+")
}
func main() {
}
`)
	gopClTestEx(t, "Rect.gox", `
var (
	Buffer
	v int
)
type Buffer struct {
	buf []byte
}
println "Go+"
`, `package main

import fmt "fmt"

type Buffer struct {
	buf []byte
}
type Rect struct {
	Buffer
	v int
}

func (this *Rect) Main() {
//line Rect.gox:9
	fmt.Println("Go+")
}
func main() {
}
`)
	gopClTestEx(t, "Rect.gox", `
var (
	*Buffer
	v int
)
type Buffer struct {
	buf []byte
}
println "Go+"
`, `package main

import fmt "fmt"

type Buffer struct {
	buf []byte
}
type Rect struct {
	*Buffer
	v int
}

func (this *Rect) Main() {
//line Rect.gox:9
	fmt.Println("Go+")
}
func main() {
}
`)
	gopClTestEx(t, "Rect.gox", `
import "bytes"
var (
	*bytes.Buffer
	v int
)
println "Go+"
`, `package main

import (
	fmt "fmt"
	bytes "bytes"
)

type Rect struct {
	*bytes.Buffer
	v int
}

func (this *Rect) Main() {
//line Rect.gox:7
	fmt.Println("Go+")
}
func main() {
}
`)
	gopClTestEx(t, "Rect.gox", `
import "bytes"
var (
	bytes.Buffer
	v int
)
println "Go+"
`, `package main

import (
	fmt "fmt"
	bytes "bytes"
)

type Rect struct {
	bytes.Buffer
	v int
}

func (this *Rect) Main() {
//line Rect.gox:7
	fmt.Println("Go+")
}
func main() {
}
`)
}

func TestBig(t *testing.T) {
	gopClTest(t, `
a := 1/2r
println a+1/2r
`, `package main

import (
	fmt "fmt"
	ng "github.com/goplus/gop/builtin/ng"
	big "math/big"
)

func main() {
//line main.gop:2
	a := ng.Bigrat_Init__2(big.NewRat(1, 2))
//line main.gop:3
	fmt.Println(a.Gop_Add(ng.Bigrat_Init__2(big.NewRat(1, 2))))
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
	fmt "fmt"
	iox "github.com/goplus/gop/builtin/iox"
	io "io"
)

var r io.Reader

func main() {
//line main.gop:6
	for _gop_it := iox.Lines(r).Gop_Enum(); ; {
		var _gop_ok bool
		line, _gop_ok := _gop_it.Next()
		if !_gop_ok {
			break
		}
//line main.gop:7
		fmt.Println(line)
	}
}
`)
}

func TestErrorWrap(t *testing.T) {
	gopClTest(t, `
import (
    "strconv"
)

func add(x, y string) (int, error) {
    return strconv.Atoi(x)? + strconv.Atoi(y)?, nil
}

func addSafe(x, y string) int {
    return strconv.Atoi(x)?:0 + strconv.Atoi(y)?:0
}

println add("100", "23")!

sum, err := add("10", "abc")
println sum, err

println addSafe("10", "abc")
`, `package main

import (
	fmt "fmt"
	strconv "strconv"
	errors "github.com/qiniu/x/errors"
)

func add(x string, y string) (int, error) {
//line main.gop:7
	var _autoGo_1 int
//line main.gop:7
	{
//line main.gop:7
		var _gop_err error
//line main.gop:7
		_autoGo_1, _gop_err = strconv.Atoi(x)
//line main.gop:7
		if _gop_err != nil {
//line main.gop:7
			_gop_err = errors.NewFrame(_gop_err, "strconv.Atoi(x)", "main.gop", 7, "main.add")
//line main.gop:7
			return 0, _gop_err
		}
//line main.gop:7
		goto _autoGo_2
	_autoGo_2:
//line main.gop:7
	}
//line main.gop:7
	var _autoGo_3 int
//line main.gop:7
	{
//line main.gop:7
		var _gop_err error
//line main.gop:7
		_autoGo_3, _gop_err = strconv.Atoi(y)
//line main.gop:7
		if _gop_err != nil {
//line main.gop:7
			_gop_err = errors.NewFrame(_gop_err, "strconv.Atoi(y)", "main.gop", 7, "main.add")
//line main.gop:7
			return 0, _gop_err
		}
//line main.gop:7
		goto _autoGo_4
	_autoGo_4:
//line main.gop:7
	}
//line main.gop:7
	return _autoGo_1 + _autoGo_3, nil
}
func addSafe(x string, y string) int {
//line main.gop:11
	return func() (_gop_ret int) {
//line main.gop:11
		var _gop_err error
//line main.gop:11
		_gop_ret, _gop_err = strconv.Atoi(x)
//line main.gop:11
		if _gop_err != nil {
//line main.gop:11
			return 0
		}
//line main.gop:11
		return
	}() + func() (_gop_ret int) {
//line main.gop:11
		var _gop_err error
//line main.gop:11
		_gop_ret, _gop_err = strconv.Atoi(y)
//line main.gop:11
		if _gop_err != nil {
//line main.gop:11
			return 0
		}
//line main.gop:11
		return
	}()
}
func main() {
//line main.gop:14
	fmt.Println(func() (_gop_ret int) {
//line main.gop:14
		var _gop_err error
//line main.gop:14
		_gop_ret, _gop_err = add("100", "23")
//line main.gop:14
		if _gop_err != nil {
//line main.gop:14
			_gop_err = errors.NewFrame(_gop_err, "add(\"100\", \"23\")", "main.gop", 14, "main.main")
//line main.gop:14
			panic(_gop_err)
		}
//line main.gop:14
		return
	}())
//line main.gop:16
	sum, err := add("10", "abc")
//line main.gop:17
	fmt.Println(sum, err)
//line main.gop:19
	fmt.Println(addSafe("10", "abc"))
}
`)
}

func init() {
	build.RegisterClassFileType(".tspx", "MyGame", []*build.Class{
		{Ext: ".tspx", Class: "Sprite"},
	}, "github.com/goplus/gop/cl/internal/spx")
}

func TestSpx(t *testing.T) {
	gopClTestEx(t, "main.tspx", `println "hi"`, `package main

import (
	fmt "fmt"
	spx "github.com/goplus/gop/cl/internal/spx"
)

type MyGame struct {
	spx.MyGame
}

func (this *MyGame) MainEntry() {
//line main.tspx:1
	fmt.Println("hi")
}
func main() {
	spx.Gopt_MyGame_Main(new(MyGame))
}
`)
	gopClTestEx(t, "Cat.tspx", `println "hi"`, `package main

import (
	fmt "fmt"
	spx "github.com/goplus/gop/cl/internal/spx"
)

type Cat struct {
	spx.Sprite
}

func (this *Cat) Main() {
//line Cat.tspx:1
	fmt.Println("hi")
}
func main() {
}
`)
}
