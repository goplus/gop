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
	"os"
	"sync"
	"syscall"
	"testing"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/gengo"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/parser/parsertest"
	"github.com/goplus/gop/scanner"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
	"golang.org/x/tools/go/packages"
)

var (
	gblFset  *token.FileSet
	baseConf *cl.Config
)

func init() {
	gox.SetDebug(gox.DbgFlagAll)
	cl.SetDebug(cl.DbgFlagAll)
	gblFset = token.NewFileSet()
	baseConf = &cl.Config{
		Fset:          gblFset,
		ModPath:       "github.com/goplus/gop",
		ModRootDir:    ".",
		GenGoPkg:      new(gengo.Runner).GenGoPkg,
		CacheLoadPkgs: true,
		NoFileLine:    true,
	}
}

func gopClTest(t *testing.T, gopcode, expected string, cachefile ...string) {
	cl.SetDisableRecover(true)
	defer cl.SetDisableRecover(false)

	fs := parsertest.NewSingleFileFS("/foo", "bar.gop", gopcode)
	pkgs, err := parser.ParseFSDir(gblFset, fs, "/foo", nil, 0)
	if err != nil {
		scanner.PrintError(os.Stderr, err)
		t.Fatal("ParseFSDir:", err)
	}
	conf := *baseConf.Ensure()
	if cachefile != nil {
		copy := *baseConf
		copy.PkgsLoader = nil
		copy.CacheFile = cachefile[0]
		copy.PersistLoadPkgs = true
		conf = *copy.Ensure()
	}
	bar := pkgs["main"]
	pkg, err := cl.NewPackage("", bar, &conf)
	if err != nil {
		t.Fatal("NewPackage:", err)
	}
	var b bytes.Buffer
	err = gox.WriteTo(&b, pkg, false)
	if err != nil {
		t.Fatal("gox.WriteTo failed:", err)
	}
	result := b.String()
	if result != expected {
		t.Fatalf("\nResult:\n%s\nExpected:\n%s\n", result, expected)
	}
	if cachefile != nil {
		if err = conf.PkgsLoader.Save(); err != nil {
			t.Fatal("PkgsLoader.Save failed:", err)
		}
	}
}

func TestEmptyPkgsLoader(t *testing.T) {
	l := &cl.PkgsLoader{}
	if l.Save() != nil {
		t.Fatal("PkgsLoader.Save failed")
	}
	if l.GenGoPkgs(nil, nil) != syscall.ENOENT {
		t.Fatal("PkgsLoader.GenGoPkgs failed")
	}
	if _, err := cl.GetModulePath("/dir-not-exists/go.mod"); err == nil {
		t.Fatal("GetModulePath failed")
	}
}

func TestNilConf(t *testing.T) {
	fs := parsertest.NewSingleFileFS("/foo", "bar.gop", `println("Hi")`)
	pkgs, err := parser.ParseFSDir(gblFset, fs, "/foo", nil, 0)
	if err != nil {
		t.Fatal("ParseFSDir:", err)
	}
	cl.NewPackage("", pkgs["main"], nil)
}

func TestInitFunc(t *testing.T) {
	gopClTest(t, `

func init() {}
func init() {}
`, `package main

func init() {
}
func init() {
}
`)
}

func TestChanRecvIssue789(t *testing.T) {
	gopClTest(t, `
func foo(ch chan int) (int, bool) {
	x, ok := (<-ch)
	return x, ok
}
`, `package main

func foo(ch chan int) (int, bool) {
	x, ok := <-ch
	return x, ok
}
`)
}

func TestNamedChanCloseIssue790(t *testing.T) {
	gopClTest(t, `
type XChan chan int

func foo(ch XChan) {
	close(ch)
}
`, `package main

type XChan chan int

func foo(ch XChan) {
	close(ch)
}
`)
}

func TestUntypedFloatIssue788(t *testing.T) {
	gopClTest(t, `
func foo(v int) bool {
    return v > 1.1e5
}
`, `package main

func foo(v int) bool {
	return v > 1.1e5
}
`)
}

func TestUnderscoreConstAndVar(t *testing.T) {
	gopClTest(t, `
const (
	c0 = 1 << iota
	_
	_
	_
	c4
)

func i() int {
	return 23
}

var (
	_ = i()
	_ = i()
)
`, `package main

const (
	c0 = 1 << iota
	_
	_
	_
	c4
)

func i() int {
	return 23
}

var _ = i()
var _ = i()
`)
}

func TestUnderscoreFuncAndMethod(t *testing.T) {
	gopClTest(t, `
func _() {
}

type T struct {
	_, _, _ int
}

func (T) _() {
}

func (T) _() {
}
`, `package main

type T struct {
	_ int
	_ int
	_ int
}

func (T) _() {
}
func (T) _() {
}
func _() {
}
`)
}

func TestErrWrapIssue772(t *testing.T) {
	gopClTest(t, `
package main

func t() (int,int,error){
	return 0, 0, nil
}

func main() {
	a, b := t()!
	println(a, b)
}`, `package main

import fmt "fmt"

func t() (int, int, error) {
	return 0, 0, nil
}
func main() {
	a, b := func() (_gop_ret int, _gop_ret2 int) {
		var _gop_err error
		_gop_ret, _gop_ret2, _gop_err = t()
		if _gop_err != nil {
			panic(_gop_err)
		}
		return
	}()
	fmt.Println(a, b)
}
`)
}

func TestErrWrapIssue778(t *testing.T) {
	gopClTest(t, `
package main

func t() error {
	return nil
}

func main() {
	t()!
}`, `package main

func t() error {
	return nil
}
func main() {
	func() {
		var _gop_err error
		_gop_err = t()
		if _gop_err != nil {
			panic(_gop_err)
		}
		return
	}()
}
`)
}

func TestIssue774(t *testing.T) {
	gopClTest(t, `
package main

import "fmt"

func main() {
	var a AA = &A{str: "hello"}
	fmt.Println(a.(*A))
}

type AA interface {
	String() string
}

type A struct {
	str string
}

func (a *A) String() string {
	return a.str
}
`, `package main

import fmt "fmt"

func main() {
	var a AA = &A{str: "hello"}
	fmt.Println(a.(*A))
}

type AA interface {
	String() string
}
type A struct {
	str string
}

func (a *A) String() string {
	return a.str
}
`)
	gopClTest(t, `
package main

import "fmt"

func main() {
	a := get()
	fmt.Println(a.(*A))
}

func get()AA{
	var a AA
	return a
}

type AA interface {
	String() string
}

type A struct {
	str string
}

func (a *A) String() string {
	return a.str
}
`, `package main

import fmt "fmt"

func main() {
	a := get()
	fmt.Println(a.(*A))
}

type AA interface {
	String() string
}

func get() AA {
	var a AA
	return a
}

type A struct {
	str string
}

func (a *A) String() string {
	return a.str
}
`)
}

func TestBlockStmt(t *testing.T) {
	gopClTest(t, `
package main

func main() {
	{
		type T int
		t := T(100)
		println(t)
	}
	{
		type T string
		t := "hello"
		println(t)
	}
}
`, `package main

import fmt "fmt"

func main() {
	{
		type T int
		t := T(100)
		fmt.Println(t)
	}
	{
		type T string
		t := "hello"
		fmt.Println(t)
	}
}
`)
}

func TestVarAfterMain(t *testing.T) {
	gopClTest(t, `
package main

func main() {
	println(i)
}

var i int
`, `package main

import fmt "fmt"

func main() {
	fmt.Println(i)
}

var i int
`)
	gopClTest(t, `
package main

func f(v float64) float64 {
	return v
}
func main() {
	sink = f(100)
}

var sink float64
`, `package main

func f(v float64) float64 {
	return v
}
func main() {
	sink = f(100)
}

var sink float64
`)
}

func TestVarAfterMain2(t *testing.T) {
	gopClTest(t, `
package main

func main() {
	println(i)
}

var i = 100
`, `package main

import fmt "fmt"

func main() {
	fmt.Println(i)
}

var i = 100
`)
}

func TestVarInMain(t *testing.T) {
	gopClTest(t, `
package main

func main() {
	v := []uint64{2, 3, 5}
	var n = len(v)
	println(n)
}`, `package main

import fmt "fmt"

func main() {
	v := []uint64{2, 3, 5}
	var n = len(v)
	fmt.Println(n)
}
`)
}

func TestSelect(t *testing.T) {
	gopClTest(t, `

func consume(xchg chan int) {
	select {
	case c := <-xchg:
		println(c)
	case xchg <- 1:
		println("send ok")
	default:
		println(0)
	}
}
`, `package main

import fmt "fmt"

func consume(xchg chan int) {
	select {
	case c := <-xchg:
		fmt.Println(c)
	case xchg <- 1:
		fmt.Println("send ok")
	default:
		fmt.Println(0)
	}
}
`)
}

func TestTypeSwitch(t *testing.T) {
	gopClTest(t, `

func bar(p *interface{}) {
}

func foo(v interface{}) {
	switch t := v.(type) {
	case int, string:
		bar(&v)
	case bool:
		var x bool = t
	default:
		bar(nil)
	}
}
`, `package main

func bar(p *interface {
}) {
}
func foo(v interface {
}) {
	switch t := v.(type) {
	case int, string:
		bar(&v)
	case bool:
		var x bool = t
	default:
		bar(nil)
	}
}
`)
}

func TestTypeSwitch2(t *testing.T) {
	gopClTest(t, `

func bar(p *interface{}) {
}

func foo(v interface{}) {
	switch bar(nil); v.(type) {
	case int, string:
		bar(&v)
	}
}
`, `package main

func bar(p *interface {
}) {
}
func foo(v interface {
}) {
	switch bar(nil); v.(type) {
	case int, string:
		bar(&v)
	}
}
`)
}

func TestTypeAssert(t *testing.T) {
	gopClTest(t, `

func foo(v interface{}) {
	x := v.(int)
	y, ok := v.(string)
}
`, `package main

func foo(v interface {
}) {
	x := v.(int)
	y, ok := v.(string)
}
`)
}

func TestInterface(t *testing.T) {
	gopClTest(t, `

type Shape interface {
	Area() float64
}

func foo(shape Shape) {
	shape.Area()
}
`, `package main

type Shape interface {
	Area() float64
}

func foo(shape Shape) {
	shape.Area()
}
`)
}

func TestInterfaceExample(t *testing.T) {
	gopClTest(t, `
type Shape interface {
	Area() float64
}

type Rect struct {
	x, y, w, h float64
}

func (p *Rect) Area() float64 {
	return p.w * p.h
}

type Circle struct {
	x, y, r float64
}

func (p *Circle) Area() float64 {
	return 3.14 * p.r * p.r
}

func Area(shapes ...Shape) float64 {
	s := 0.0
	for shape <- shapes {
		s += shape.Area()
	}
	return s
}

func main() {
	rect := &Rect{0, 0, 2, 5}
	circle := &Circle{0, 0, 3}
	println(Area(circle, rect))
}
`, `package main

import fmt "fmt"

type Shape interface {
	Area() float64
}
type Rect struct {
	x float64
	y float64
	w float64
	h float64
}

func (p *Rect) Area() float64 {
	return p.w * p.h
}

type Circle struct {
	x float64
	y float64
	r float64
}

func (p *Circle) Area() float64 {
	return 3.14 * p.r * p.r
}
func Area(shapes ...Shape) float64 {
	s := 0.0
	for _, shape := range shapes {
		s += shape.Area()
	}
	return s
}
func main() {
	rect := &Rect{0, 0, 2, 5}
	circle := &Circle{0, 0, 3}
	fmt.Println(Area(circle, rect))
}
`)
}

func TestEmbeddField(t *testing.T) {
	gopClTest(t, `import "math/big"

type BigInt struct {
	*big.Int
}`, `package main

import big "math/big"

type BigInt struct {
	*big.Int
}
`)
}

func TestAutoProperty(t *testing.T) {
	gopClTest(t, `import "github.com/goplus/gop/ast/goptest"

func foo(script string) {
	doc := goptest.New(script)!

	println(doc.any.funcDecl.name)
	println(doc.any.importSpec.name)
}
`, `package main

import (
	fmt "fmt"
	goptest "github.com/goplus/gop/ast/goptest"
	gopq "github.com/goplus/gop/ast/gopq"
)

func foo(script string) {
	doc := func() (_gop_ret gopq.NodeSet) {
		var _gop_err error
		_gop_ret, _gop_err = goptest.New(script)
		if _gop_err != nil {
			panic(_gop_err)
		}
		return
	}()
	fmt.Println(doc.Any().FuncDecl().Name())
	fmt.Println(doc.Any().ImportSpec().Name())
}
`)
}

func TestErrWrap(t *testing.T) {
	gopClTest(t, `
import "strconv"

func add(x, y string) (int, error) {
	return strconv.Atoi(x)? + strconv.Atoi(y)?, nil
}
`, `package main

import strconv "strconv"

func add(x string, y string) (int, error) {
	var _autoGop_1 int
	{
		var _gop_err error
		_autoGop_1, _gop_err = strconv.Atoi(x)
		if _gop_err != nil {
			return 0, _gop_err
		}
		goto _autoGop_2
	_autoGop_2:
	}
	var _autoGop_3 int
	{
		var _gop_err error
		_autoGop_3, _gop_err = strconv.Atoi(y)
		if _gop_err != nil {
			return 0, _gop_err
		}
		goto _autoGop_4
	_autoGop_4:
	}
	return _autoGop_1 + _autoGop_3, nil
}
`)
}

func TestErrWrapDefVal(t *testing.T) {
	gopClTest(t, `
import "strconv"

func addSafe(x, y string) int {
	return strconv.Atoi(x)?:0 + strconv.Atoi(y)?:0
}
`, `package main

import strconv "strconv"

func addSafe(x string, y string) int {
	return func() (_gop_ret int) {
		var _gop_err error
		_gop_ret, _gop_err = strconv.Atoi(x)
		if _gop_err != nil {
			return 0
		}
		return
	}() + func() (_gop_ret int) {
		var _gop_err error
		_gop_ret, _gop_err = strconv.Atoi(y)
		if _gop_err != nil {
			return 0
		}
		return
	}()
}
`)
}

func TestErrWrapPanic(t *testing.T) {
	gopClTest(t, `
var ret int = println("Hi")!
`, `package main

import fmt "fmt"

var ret int = func() (_gop_ret int) {
	var _gop_err error
	_gop_ret, _gop_err = fmt.Println("Hi")
	if _gop_err != nil {
		panic(_gop_err)
	}
	return
}()
`)
}

func TestMakeAndNew(t *testing.T) {
	gopClTest(t, `
var a *int = new(int)
var b map[string]int = make(map[string]int)
var c []byte = make([]byte, 0, 2)
`, `package main

var a *int = new(int)
var b map[string]int = make(map[string]int)
var c []byte = make([]byte, 0, 2)
`)
}

func TestVarDecl(t *testing.T) {
	gopClTest(t, `
var a int
var x, y = 1, "Hi"
`, `package main

var a int
var x, y = 1, "Hi"
`)
}

func TestBigIntAdd(t *testing.T) {
	gopClTest(t, `
var x, y bigint
var z bigint = x + y
`, `package main

import builtin "github.com/goplus/gop/builtin"

var x, y builtin.Gop_bigint
var z builtin.Gop_bigint = x.Gop_Add(y)
`)
}

func TestBigIntLit(t *testing.T) {
	gopClTest(t, `
var x = 1r
`, `package main

import (
	builtin "github.com/goplus/gop/builtin"
	big "math/big"
)

var x = builtin.Gop_bigint_Init__1(big.NewInt(1))
`)
}

func TestBigRatLit(t *testing.T) {
	gopClTest(t, `
var x = 1/2r
`, `package main

import (
	builtin "github.com/goplus/gop/builtin"
	big "math/big"
)

var x = builtin.Gop_bigrat_Init__2(big.NewRat(1, 2))
`)
}

func TestBigRatLitAdd(t *testing.T) {
	gopClTest(t, `
var x = 3 + 1/2r
`, `package main

import (
	builtin "github.com/goplus/gop/builtin"
	big "math/big"
)

var x = builtin.Gop_bigrat_Init__2(big.NewRat(7, 2))
`)
}

func TestBigRatAdd(t *testing.T) {
	gox.SetDebug(gox.DbgFlagAll)
	gopClTest(t, `
var x = 3 + 1/2r
var y = x + 100
var z = 100 + y
`, `package main

import (
	builtin "github.com/goplus/gop/builtin"
	big "math/big"
)

var x = builtin.Gop_bigrat_Init__2(big.NewRat(7, 2))
var y = x.Gop_Add(builtin.Gop_bigrat_Init__0(100))
var z = builtin.Gop_bigrat_Init__0(100) + y
`)
}

func TestTypeConv(t *testing.T) {
	gopClTest(t, `
var a = (*struct{})(nil)
var b = interface{}(nil)
var c = (func())(nil)
var x uint32 = uint32(0)
var y *uint32 = (*uint32)(nil)
`, `package main

var a = (*struct {
})(nil)
var b = interface {
}(nil)
var c = (func())(nil)
var x uint32 = uint32(0)
var y *uint32 = (*uint32)(nil)
`)
}

func TestStar(t *testing.T) {
	gopClTest(t, `
var x *uint32 = (*uint32)(nil)
var y uint32 = *x
`, `package main

var x *uint32 = (*uint32)(nil)
var y uint32 = *x
`)
}

func TestSend(t *testing.T) {
	gopClTest(t, `
var x chan bool
x <- true
`, `package main

var x chan bool

func main() {
	x <- true
}
`)
}

func TestIncDec(t *testing.T) {
	gopClTest(t, `
var x uint32
x++
`, `package main

var x uint32

func main() {
	x++
}
`)
}

func TestAssignOp(t *testing.T) {
	gopClTest(t, `
var x uint32
x += 3
`, `package main

var x uint32

func main() {
	x += 3
}
`)
}

func TestBigIntAssignOp(t *testing.T) {
	gopClTest(t, `
var x bigint
x += 3
`, `package main

import builtin "github.com/goplus/gop/builtin"

var x builtin.Gop_bigint

func main() {
	x.Gop_AddAssign(builtin.Gop_bigint_Init__0(3))
}
`)
}

func TestBigIntAssignOp2(t *testing.T) {
	gopClTest(t, `
x := 3r
x *= 2
`, `package main

import (
	builtin "github.com/goplus/gop/builtin"
	big "math/big"
)

func main() {
	x := builtin.Gop_bigint_Init__1(big.NewInt(3))
	x.Gop_MulAssign(builtin.Gop_bigint_Init__0(2))
}
`)
}

func TestBigIntAssignOp3(t *testing.T) {
	gopClTest(t, `
x := 3r
x *= 2r
`, `package main

import (
	builtin "github.com/goplus/gop/builtin"
	big "math/big"
)

func main() {
	x := builtin.Gop_bigint_Init__1(big.NewInt(3))
	x.Gop_MulAssign(builtin.Gop_bigint_Init__1(big.NewInt(2)))
}
`)
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
	z := [...]int{1, 3, 5}
	a := map[string]float64{"Hello": 1, "Go+": 5.1}
}
`)
}

func TestCompositeLit2(t *testing.T) {
	gopClTest(t, `
type foo struct {
	A int
}

x := []*struct{a int}{
	{1}, {3}, {5},
}
y := map[foo]struct{a string}{
	{1}: {"Hi"},
}
z := [...]foo{
	{1}, {3}, {5},
}
`, `package main

type foo struct {
	A int
}

func main() {
	x := []*struct {
		a int
	}{&struct {
		a int
	}{1}, &struct {
		a int
	}{3}, &struct {
		a int
	}{5}}
	y := map[foo]struct {
		a string
	}{foo{1}: struct {
		a string
	}{"Hi"}}
	z := [...]foo{foo{1}, foo{3}, foo{5}}
}
`)
}

// deduce struct type as parameters of a function call
func TestCompositeLit3(t *testing.T) {
	gopClTest(t, `
type Config struct {
	A int
}

func foo(conf *Config) {
}

func bar(conf ...Config) {
}

foo({A: 1})
bar({A: 2})
`, `package main

type Config struct {
	A int
}

func foo(conf *Config) {
}
func bar(conf ...Config) {
}
func main() {
	foo(&Config{A: 1})
	bar(Config{A: 2})
}
`)
}

// deduce struct type as results of a function call
func TestCompositeLit4(t *testing.T) {
	gopClTest(t, `
type Result struct {
	A int
}

func foo() *Result {
	return {A: 1}
}
`, `package main

type Result struct {
	A int
}

func foo() *Result {
	return &Result{A: 1}
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

func TestChan(t *testing.T) {
	gopClTest(t, `
a := make(chan int, 10)
a <- 3
var b int = <-a
x, ok := <-a
`, `package main

func main() {
	a := make(chan int, 10)
	a <- 3
	var b int = <-a
	x, ok := <-a
}
`)
}

func TestKeyValModeLit(t *testing.T) {
	gopClTest(t, `
a := [...]float64{1, 3: 3.4, 5}
b := []float64{2: 1.2, 3, 6: 4.5}
`, `package main

func main() {
	a := [...]float64{1, 3: 3.4, 5}
	b := []float64{2: 1.2, 3, 6: 4.5}
}
`)
}

func TestStructLit(t *testing.T) {
	gopClTest(t, `
type foo struct {
	A int
	B string "tag1:123"
}

a := struct {
	A int
	B string "tag1:123"
}{1, "Hello"}

b := foo{1, "Hello"}
c := foo{B: "Hi"}
`, `package main

type foo struct {
	A int
	B string "tag1:123"
}

func main() {
	a := struct {
		A int
		B string "tag1:123"
	}{1, "Hello"}
	b := foo{1, "Hello"}
	c := foo{B: "Hi"}
}
`)
}

func TestStructType(t *testing.T) {
	gopClTest(t, `
type bar = foo

type foo struct {
	p *bar
	A int
	B string "tag1:123"
}

func main() {
	type a struct {
		p *a
	}
	type b = a
}
`, `package main

type foo struct {
	p *foo
	A int
	B string "tag1:123"
}
type bar = foo

func main() {
	type a struct {
		p *a
	}
	type b = a
}
`)
}

func TestDeferGo(t *testing.T) {
	gopClTest(t, `
go println("Hi")
defer println("Go+")
`, `package main

import fmt "fmt"

func main() {
	go fmt.Println("Hi")
	defer fmt.Println("Go+")
}
`)
}

func TestFor(t *testing.T) {
	gopClTest(t, `
a := [1, 3.4, 5]
for i := 0; i < 3; i=i+1 {
	println(i)
}
for {
	println("loop")
}
`, `package main

import fmt "fmt"

func main() {
	a := []float64{1, 3.4, 5}
	for i := 0; i < 3; i = i + 1 {
		fmt.Println(i)
	}
	for {
		fmt.Println("loop")
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
for i, x := range a {
	println(i, x)
}

var i int
var x float64
for _, x = range a {
	println(i, x)
}
for i, x = range a {
	println(i, x)
}
for range a {
	println("Hi")
}
`, `package main

import fmt "fmt"

func main() {
	a := []float64{1, 3.4, 5}
	for _, x := range a {
		fmt.Println(x)
	}
	for i, x := range a {
		fmt.Println(i, x)
	}
	var i int
	var x float64
	for _, x = range a {
		fmt.Println(i, x)
	}
	for i, x = range a {
		fmt.Println(i, x)
	}
	for range a {
		fmt.Println("Hi")
	}
}
`)
}

func TestRangeStmtUDT(t *testing.T) {
	gopClTest(t, `
type foo struct {
}

func (p *foo) Gop_Enum(c func(key int, val string)) {
}

for k, v := range new(foo) {
	println(k, v)
}
`, `package main

import fmt "fmt"

type foo struct {
}

func (p *foo) Gop_Enum(c func(key int, val string)) {
}
func main() {
	new(foo).Gop_Enum(func(k int, v string) {
		fmt.Println(k, v)
	})
}
`)
}

func TestForPhraseUDT(t *testing.T) {
	gopClTest(t, `
type foo struct {
}

func (p *foo) Gop_Enum(c func(val string)) {
}

for v <- new(foo) {
	println(v)
}
`, `package main

import fmt "fmt"

type foo struct {
}

func (p *foo) Gop_Enum(c func(val string)) {
}
func main() {
	new(foo).Gop_Enum(func(v string) {
		fmt.Println(v)
	})
}
`)
}

func TestForPhraseUDT2(t *testing.T) {
	gopClTest(t, `
type fooIter struct {
}

func (p fooIter) Next() (key string, val int, ok bool) {
	return
}

type foo struct {
}

func (p *foo) Gop_Enum() fooIter {
}

for k, v <- new(foo) {
	println(k, v)
}
`, `package main

import fmt "fmt"

type fooIter struct {
}

func (p fooIter) Next() (key string, val int, ok bool) {
	return
}

type foo struct {
}

func (p *foo) Gop_Enum() fooIter {
}
func main() {
	for _gop_it := new(foo).Gop_Enum(); ; {
		var _gop_ok bool
		k, v, _gop_ok := _gop_it.Next()
		if !_gop_ok {
			break
		}
		fmt.Println(k, v)
	}
}
`)
}

func TestForPhraseUDT3(t *testing.T) {
	gopClTest(t, `
type foo struct {
}

func (p *foo) Gop_Enum(c func(val string)) {
}

println([v for v <- new(foo)])
`, `package main

import fmt "fmt"

type foo struct {
}

func (p *foo) Gop_Enum(c func(val string)) {
}
func main() {
	fmt.Println(func() (_gop_ret []string) {
		new(foo).Gop_Enum(func(v string) {
			_gop_ret = append(_gop_ret, v)
		})
		return
	}())
}
`)
}

func TestForPhraseUDT4(t *testing.T) {
	gopClTest(t, `
type fooIter struct {
	data *foo
	idx  int
}

func (p *fooIter) Next() (key int, val string, ok bool) {
	if p.idx < len(p.data.key) {
		key, val, ok = p.data.key[p.idx], p.data.val[p.idx], true
		p.idx++
	}
	return
}

type foo struct {
	key []int
	val []string
}

func newFoo() *foo {
	return &foo{key: [3, 7], val: ["Hi", "Go+"]}
}

func (p *foo) Gop_Enum() *fooIter {
	return &fooIter{data: p}
}

for k, v <- newFoo() {
	println(k, v)
}
`, `package main

import fmt "fmt"

type fooIter struct {
	data *foo
	idx  int
}
type foo struct {
	key []int
	val []string
}

func (p *fooIter) Next() (key int, val string, ok bool) {
	if p.idx < len(p.data.key) {
		key, val, ok = p.data.key[p.idx], p.data.val[p.idx], true
		p.idx++
	}
	return
}
func (p *foo) Gop_Enum() *fooIter {
	return &fooIter{data: p}
}
func newFoo() *foo {
	return &foo{key: []int{3, 7}, val: []string{"Hi", "Go+"}}
}
func main() {
	for _gop_it := newFoo().Gop_Enum(); ; {
		var _gop_ok bool
		k, v, _gop_ok := _gop_it.Next()
		if !_gop_ok {
			break
		}
		fmt.Println(k, v)
	}
}
`)
}

func TestForPhrase(t *testing.T) {
	gopClTest(t, `
sum := 0
for x <- [1, 3, 5, 7, 11, 13, 17], x > 3 {
	sum = sum + x
}
for i, x <- [1, 3, 5, 7, 11, 13, 17] {
	sum = sum + i*x
}
println("sum(5,7,11,13,17):", sum)
`, `package main

import fmt "fmt"

func main() {
	sum := 0
	for _, x := range []int{1, 3, 5, 7, 11, 13, 17} {
		if x > 3 {
			sum = sum + x
		}
	}
	for i, x := range []int{1, 3, 5, 7, 11, 13, 17} {
		sum = sum + i*x
	}
	fmt.Println("sum(5,7,11,13,17):", sum)
}
`)
}

func TestMapComprehension(t *testing.T) {
	gopClTest(t, `
y := {x: i for i, x <- ["1", "3", "5", "7", "11"]}
`, `package main

func main() {
	y := func() (_gop_ret map[string]int) {
		_gop_ret = map[string]int{}
		for i, x := range []string{"1", "3", "5", "7", "11"} {
			_gop_ret[x] = i
		}
		return
	}()
}
`)
}

func TestMapComprehensionCond(t *testing.T) {
	gopClTest(t, `
z := {v: k for k, v <- {"Hello": 1, "Hi": 3, "xsw": 5, "Go+": 7}, v > 3}
`, `package main

func main() {
	z := func() (_gop_ret map[int]string) {
		_gop_ret = map[int]string{}
		for k, v := range map[string]int{"Hello": 1, "Hi": 3, "xsw": 5, "Go+": 7} {
			if v > 3 {
				_gop_ret[v] = k
			}
		}
		return
	}()
}
`)
}

func TestMapComprehensionCond2(t *testing.T) {
	gopClTest(t, `
z := {t: k for k, v <- {"Hello": 1, "Hi": 3, "xsw": 5, "Go+": 7}, t := v; t > 3}
`, `package main

func main() {
	z := func() (_gop_ret map[int]string) {
		_gop_ret = map[int]string{}
		for k, v := range map[string]int{"Hello": 1, "Hi": 3, "xsw": 5, "Go+": 7} {
			if t := v; t > 3 {
				_gop_ret[t] = k
			}
		}
		return
	}()
}
`)
}

func TestExistsComprehension(t *testing.T) {
	gopClTest(t, `
hasFive := {for x <- ["1", "3", "5", "7", "11"], x == "5"}
`, `package main

func main() {
	hasFive := func() (_gop_ok bool) {
		for _, x := range []string{"1", "3", "5", "7", "11"} {
			if x == "5" {
				return true
			}
		}
		return
	}()
}
`)
}

func TestSelectComprehension(t *testing.T) {
	gopClTest(t, `
y := {i for i, x <- ["1", "3", "5", "7", "11"], x == "5"}
`, `package main

func main() {
	y := func() (_gop_ret int) {
		for i, x := range []string{"1", "3", "5", "7", "11"} {
			if x == "5" {
				return i
			}
		}
		return
	}()
}
`)
}

func TestSelectComprehensionTwoValue(t *testing.T) {
	gopClTest(t, `
y, ok := {i for i, x <- ["1", "3", "5", "7", "11"], x == "5"}
`, `package main

func main() {
	y, ok := func() (_gop_ret int, _gop_ok bool) {
		for i, x := range []string{"1", "3", "5", "7", "11"} {
			if x == "5" {
				return i, true
			}
		}
		return
	}()
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

func TestSliceGet(t *testing.T) {
	gopClTest(t, `
a := [1, 3, 5, 7, 9]
b := a[:3]
c := a[1:]
d := a[1:2:3]
e := "Hello, Go+"[7:]
`, `package main

func main() {
	a := []int{1, 3, 5, 7, 9}
	b := a[:3]
	c := a[1:]
	d := a[1:2:3]
	e := "Hello, Go+"[7:]
}
`)
}

func TestIndexGetTwoValue(t *testing.T) {
	gopClTest(t, `
a := {"Hello": 1, "Hi": 3, "xsw": 5, "Go+": 7}
x, ok := a["Hi"]
y := a["Go+"]
`, `package main

func main() {
	a := map[string]int{"Hello": 1, "Hi": 3, "xsw": 5, "Go+": 7}
	x, ok := a["Hi"]
	y := a["Go+"]
}
`)
}

func TestIndexGet(t *testing.T) {
	gopClTest(t, `
a := [1, 3.4, 5]
b := a[1]
`, `package main

func main() {
	a := []float64{1, 3.4, 5}
	b := a[1]
}
`)
}

func TestIndexRef(t *testing.T) {
	gopClTest(t, `
a := [1, 3.4, 5]
a[1] = 2.1
`, `package main

func main() {
	a := []float64{1, 3.4, 5}
	a[1] = 2.1
}
`)
}

func TestIndexArrayPtrIssue784(t *testing.T) {
	gopClTest(t, `
type intArr [2]int

func foo(a *intArr) {
	a[1] = 10
}
`, `package main

type intArr [2]int

func foo(a *intArr) {
	a[1] = 10
}
`)
}

func TestMemberVal(t *testing.T) {
	gopClTest(t, `import "strings"

x := strings.NewReplacer("?", "!").Replace("hello, world???")
println("x:", x)
`, `package main

import (
	fmt "fmt"
	strings "strings"
)

func main() {
	x := strings.NewReplacer("?", "!").Replace("hello, world???")
	fmt.Println("x:", x)
}
`)
}

func TestNamedPtrMemberIssue786(t *testing.T) {
	gopClTest(t, `
type foo struct {
	req int
}

type pfoo *foo

func bar(p pfoo) {
	println(p.req)
}
`, `package main

import fmt "fmt"

type foo struct {
	req int
}
type pfoo *foo

func bar(p pfoo) {
	fmt.Println(p.req)
}
`)
}

func TestMember(t *testing.T) {
	gopClTest(t, `

import "flag"

a := &struct {
	A int
	B string
}{1, "Hello"}

x := a.A
a.B = "Hi"

flag.Usage = nil
`, `package main

import flag "flag"

func main() {
	a := &struct {
		A int
		B string
	}{1, "Hello"}
	x := a.A
	a.B = "Hi"
	flag.Usage = nil
}
`)
}

func TestElem(t *testing.T) {
	gopClTest(t, `

func foo(a *int, b int) {
	b = *a
	*a = b
}
`, `package main

func foo(a *int, b int) {
	b = *a
	*a = b
}
`)
}

func TestMethod(t *testing.T) {
	gopClTest(t, `
type M int

func (m M) Foo() {
	println("foo", m)
}

func (M) Bar() {
	println("bar")
}
`, `package main

import fmt "fmt"

type M int

func (m M) Foo() {
	fmt.Println("foo", m)
}
func (M) Bar() {
	fmt.Println("bar")
}
`)
}

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

var a, b foo
var c = a - b
var d = -a
`, `package main

import fmt "fmt"

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
func (a foo) Gop_Neg() {
	fmt.Println("-a")
}
func (a foo) Gop_Inc() {
	fmt.Println("a++")
}

var a, b foo
var c = a.Gop_Sub(b)
var d = a.Gop_Neg()
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

func TestImportForceUsed(t *testing.T) {
	gopClTest(t, `import _ "fmt"

func main() {
}`, `package main

import _ "fmt"

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

func TestVarAndConst(t *testing.T) {
	gopClTest(t, `
const (
	i = 1
	x float64 = 1
)
var j int = i
`, `package main

const (
	i         = 1
	x float64 = 1
)

var j int = i
`)
}

func TestDeclStmt(t *testing.T) {
	gopClTest(t, `import "fmt"

func main() {
	const (
		i = 1
		x float64 = 1
	)
	var j int = i
	fmt.Println("Hi")
}`, `package main

import fmt "fmt"

func main() {
	const (
		i         = 1
		x float64 = 1
	)
	var j int = i
	fmt.Println("Hi")
}
`)
}

func TestIf(t *testing.T) {
	gopClTest(t, `x := 0
if t := false; t {
	x = 3
} else if !t {
	x = 5
} else {
	x = 7
}
println("x:", x)
`, `package main

import fmt "fmt"

func main() {
	x := 0
	if t := false; t {
		x = 3
	} else if !t {
		x = 5
	} else {
		x = 7
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

func TestBranchStmt(t *testing.T) {
	gopClTest(t, `
	a := [1, 3.4, 5]
label:
	for i := 0; i < 3; i=i+1 {
		println(i)
		break
		break label
		continue
		continue label
		goto label
	}
`, `package main

import fmt "fmt"

func main() {
	a := []float64{1, 3.4, 5}
label:
	for i := 0; i < 3; i = i + 1 {
		fmt.Println(i)
		break
		break label
		continue
		continue label
		goto label
	}
}
`)
}

func TestReturn(t *testing.T) {
	gopClTest(t, `
func foo(format string, args ...interface{}) (int, error) {
	return printf(format, args...)
}

func main() {
}
`, `package main

import fmt "fmt"

func foo(format string, args ...interface {
}) (int, error) {
	return fmt.Printf(format, args...)
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
}) (int, error) {
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
	gopClTest(t, `func foo(format string, a [10]int, args ...interface{}) {
}

func main() {
}`, `package main

func foo(format string, a [10]int, args ...interface {
}) {
}
func main() {
}
`)
}

func TestLambdaExpr(t *testing.T) {
	gopClTest(t, `
func Map(c []float64, t func(float64) float64) {
	// ...
}

func Map2(c []float64, t func(float64) (float64, float64)) {
	// ...
}

Map([1.2, 3.5, 6], x => x * x)
Map2([1.2, 3.5, 6], x => (x * x, x + x))
`, `package main

func Map(c []float64, t func(float64) float64) {
}
func Map2(c []float64, t func(float64) (float64, float64)) {
}
func main() {
	Map([]float64{1.2, 3.5, 6}, func(x float64) float64 {
		return x * x
	})
	Map2([]float64{1.2, 3.5, 6}, func(x float64) (float64, float64) {
		return x * x, x + x
	})
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

func bar(foo func(string, ...interface {
}) (int, error)) {
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

func foo(x string) string {
	return strings.NewReplacer("?", "!").Replace(x)
}
func printf(format string, args ...interface {
}) (n int, err error) {
	n, err = fmt.Printf(format, args...)
	return
}
func bar(foo func(string, ...interface {
}) (int, error)) {
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

var (
	autogen sync.Mutex
)

func removeAutogenFiles() {
	os.Remove("../tutorial/14-Using-goplus-in-Go/foo/gop_autogen.go")
	os.Remove("../tutorial/14-Using-goplus-in-Go/foo/gop_autogen_test.go")
	os.Remove("../tutorial/14-Using-goplus-in-Go/foo/gop_autogen2_test.go")
}

func TestImportGopPkg(t *testing.T) {
	autogen.Lock()
	defer autogen.Unlock()

	removeAutogenFiles()
	gopClTest(t, `import "github.com/goplus/gop/tutorial/14-Using-goplus-in-Go/foo"

rmap := foo.ReverseMap(map[string]int{"Hi": 1, "Hello": 2})
println(rmap)
`, `package main

import (
	fmt "fmt"
	foo "github.com/goplus/gop/tutorial/14-Using-goplus-in-Go/foo"
)

func main() {
	rmap := foo.ReverseMap(map[string]int{"Hi": 1, "Hello": 2})
	fmt.Println(rmap)
}
`)
}

// bugfix (only depends order of testing functions)
// vet: open tutorial/14-Using-goplus-in-Go/foo/gop_autogen.go: no such file or directory
func TestGopkgDep(t *testing.T) {
	autogen.Lock()
	defer autogen.Unlock()

	removeAutogenFiles()
	const (
		loadTypes = packages.NeedImports | packages.NeedDeps | packages.NeedTypes
		loadModes = loadTypes | packages.NeedName | packages.NeedModule
	)
	loadConf := &packages.Config{Mode: loadModes, Fset: gblFset}
	pkgs, err := baseConf.Ensure().PkgsLoader.Load(
		loadConf, "github.com/goplus/gop/tutorial/14-Using-goplus-in-Go/gomain")
	if err != nil {
		t.Fatal("PkgsLoader.Load failed:", err)
	}
	for _, err = range pkgs[0].Errors {
		t.Fatal("PkgsLoader.Load failed:", err)
	}
	if pkgs[0].Name != "main" {
		t.Fatal("pkg name:", pkgs[0].Name)
	}
}

func TestCallDep(t *testing.T) {
	const (
		cachedir  = "../.gop"
		cachefile = cachedir + "/gop.cache"
	)
	os.Remove(cachefile)
	defer func() {
		os.Remove(cachefile)
		os.Remove(cachedir)
	}()
	for i := 0; i < 2; i++ {
		gopClTest(t, `
import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	ret := New()
	expected := Result{}
	if reflect.DeepEqual(ret, expected) {
		t.Fatal("Test failed:", ret, expected)
	}
}

type Repo struct {
	Title string
}

func newRepo() Repo {
	return {Title: "Hi"}
}

type Result struct {
	Repo Repo
}

func New() Result {
	repo := newRepo()
	return {Repo: repo}
}
`, `package main

import (
	reflect "reflect"
	testing "testing"
)

func TestNew(t *testing.T) {
	ret := New()
	expected := Result{}
	if reflect.DeepEqual(ret, expected) {
		t.Fatal("Test failed:", ret, expected)
	}
}

type Result struct {
	Repo Repo
}

func New() Result {
	repo := newRepo()
	return Result{Repo: repo}
}

type Repo struct {
	Title string
}

func newRepo() Repo {
	return Repo{Title: "Hi"}
}
`, "")
	}
}
