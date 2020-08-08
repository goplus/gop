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
package cl_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/goplus/gop/cl/cltest"
)

// -----------------------------------------------------------------------------

func TestNew(t *testing.T) {
	cltest.Expect(t, `
		a := new([2]int)
		println("a:", a)
		`,
		"a: &[0 0]\n",
	)
	cltest.Expect(t, `
		println(new())
		`,
		"",
		"missing argument to new\n",
	)
	cltest.Expect(t, `
		println(new(int, float64))
		`,
		"",
		"too many arguments to new(int)\n",
	)
}

func TestNew2(t *testing.T) {
	cltest.Expect(t, `
		a := new([2]int)
		a[0] = 2
		println("a:", a[0])
		`,
		"a: 2\n",
	)
	cltest.Expect(t, `
		a := new([2]float64)
		a[0] = 1.1
		println("a:", a[0])
		`,
		"a: 1.1\n",
	)
	cltest.Expect(t, `
		a := new([2]string)
		a[0] = "gop"
		println("a:", a[0])
		`,
		"a: gop\n",
	)
}

func TestBadIndex(t *testing.T) {
	cltest.Expect(t, `
		a := new(int)
		println(a[0])
		`,
		"",
		nil,
	)
	cltest.Expect(t, `
		a := new(int)
		a[0] = 2
		`,
		"",
		nil,
	)
}

// -------------------`----------------------------------------------------------

func TestAutoProperty(t *testing.T) {
	script := `
		import "io"

		func New() (*Bar, error) {
			return nil, io.EOF
		}

		bar, err := New()
		if err != nil {
			log.Println(err)
		}
	`
	gopcode := `
		import (
			"github.com/goplus/gop/ast/goptest"
		)

		script := %s

		doc := goptest.New(script)!
		println(doc.any.funcDecl.name)
	`
	cltest.Expect(t,
		fmt.Sprintf(gopcode, "`"+script+"`"),
		"[New main]\n",
	)
}

func TestAutoProperty2(t *testing.T) {
	cltest.Expect(t, `
		import "bytes"
		import "os"

		b := bytes.newBuffer(nil)
		b.writeString("Hello, ")
		b.writeString("Go+")
		println(b.string)
		`,
		"Hello, Go+\n",
	)
	cltest.Expect(t, `
		import "bytes"
		import "os"

		b := bytes.newBuffer(nil)
		b.writeString("Hello, ")
		b.writeString("Go+")
		println(b.string2)
		`,
		"",
		nil, // panic
	)
	cltest.Expect(t, `
		import "bytes"

		bytes.newBuf()
		`,
		"",
		nil, // panic
	)
}

func TestUnbound(t *testing.T) {
	cltest.Expect(t, `
		println("Hello " + "qiniu:", 123, 4.5, 7i)
		`,
		"Hello qiniu: 123 4.5 (0+7i)\n",
	)
}

func TestUnboundInt(t *testing.T) {
	cltest.Expect(t, `
	import "reflect"
	printf("%T",100)
	`,
		"int",
	)
	cltest.Expect(t, `
	import "reflect"
	printf("%T",-100)
	`,
		"int",
	)
}

func TestOverflowsInt(t *testing.T) {
	cltest.Expect(t, `
	println(9223372036854775807)
	`,
		"9223372036854775807\n",
	)
	cltest.Expect(t, `
	println(-9223372036854775808)
	`,
		"-9223372036854775808\n",
	)
	cltest.Expect(t, `
	println(9223372036854775808)
	`,
		"",
		nil,
	)
}

func TestOpLAndLOr(t *testing.T) {
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
func bar() bool {
	println("bar")
	return true
}

func fake() bool {
	println("fake")
	return false
}

if foo() || bar() {
}
println("---")
if foo() && bar() {
}
println("---")
if fake() && bar() {
}
	`, "foo\n---\nfoo\nbar\n---\nfake\n")
}

func TestOpLAndLOr2(t *testing.T) {
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
func bar() bool {
	println("bar")
	return true
}

func fake() bool {
	println("fake")
	return true
}

if foo() && bar() && fake() {
}
	`, "foo\nbar\nfake\n")
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
func bar() bool {
	println("bar")
	return false
}

func fake() bool {
	println("fake")
	return true
}

if foo() && bar() && fake() {
}
	`, "foo\nbar\n")
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return false
}
func bar() bool {
	println("bar")
	return true
}

func fake() bool {
	println("fake")
	return true
}

if foo() || bar() || fake() {
}
	`, "foo\nbar\n")
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
func bar() bool {
	println("bar")
	return true
}

func fake() bool {
	println("fake")
	return true
}

if foo() || bar() || fake() {
}
	`, "foo\n")
}

func TestOpLAndLOr3(t *testing.T) {
	cltest.Expect(t, `
func foo() int {
	println("foo")
	return 0
}
func bar() bool {
	println("bar")
	return true
}
if foo() || bar() {
}
	`, "", nil)
	cltest.Expect(t, `
func foo() int {
	println("foo")
	return 0
}
func bar() bool {
	println("bar")
	return true
}
if foo() && bar() {
}
	`, "", nil)
}

func TestOpLAndLOr4(t *testing.T) {
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
if true || foo() {
}
	`, "")
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
if false || foo() {
}
	`, "foo\n")
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
if true && foo() {
}
	`, "foo\n")
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
if false && foo() {
}
	`, "")
}

func TestPanic(t *testing.T) {
	cltest.Expect(t,
		`panic("Helo")`,
		"",
		"Helo", // panicMsg
	)
}

func TestTypeCast(t *testing.T) {
	cltest.Call(t, `
	x := []byte("hello")
	x
	`).Equal([]byte("hello"))
}

func TestAppendErr(t *testing.T) {
	cltest.Expect(t, `
		append()
		`,
		"",
		"append: argument count not enough\n",
	)
	cltest.Expect(t, `
		x := 1
		append(x, 2)
		`,
		"",
		"append: first argument not a slice\n",
	)
	cltest.Expect(t, `
		defer append([]int{1}, 2)
		`,
		"",
		"defer discards result of append([]int{1}, 2)\n",
	)
}

func TestLenErr(t *testing.T) {
	cltest.Expect(t, `
		len()
		`,
		"",
		"missing argument to len: len()\n",
	)
	cltest.Expect(t, `
		len("a", "b")
		`,
		"",
		`too many arguments to len: len("a", "b")`+"\n",
	)
}

func TestMake(t *testing.T) {
	cltest.Expect(t, `
		make()
		`,
		"",
		"missing argument to make: make()\n",
	)
	cltest.Expect(t, `
		a := make([]int, 0, 4)
		a = append(a, 1, 2, 3)
		println(a)
		`,
		"[1 2 3]\n",
	)
	cltest.Expect(t, `
		a := make([]int, 0, 4)
		a = append(a, [1, 2, 3]...)
		println(a)
		`,
		"[1 2 3]\n",
	)
	cltest.Expect(t, `
		n := 4
		a := make(map[string]interface{}, uint16(n))
		println(a)
		`,
		"map[]\n",
	)
	cltest.Expect(t, `
		import "reflect"

		a := make(chan *func(), uint16(4))
		println(reflect.TypeOf(a))
		`,
		"chan *func()\n",
	)
}

func TestOperator(t *testing.T) {
	cltest.Expect(t, `
		println("Hello", 123 * 4.5, 1 + 7i)
		`,
		"Hello 553.5 (1+7i)\n")
}

func TestVar(t *testing.T) {
	cltest.Expect(t, `
		x := 123.1
		println("Hello", x)
		`,
		"Hello 123.1\n")
}

func TestVarOp(t *testing.T) {
	cltest.Expect(t, `
		x := 123.1
		y := 1 + x
		println("Hello", y + 10)
		n, err := println("Hello", y + 10)
		println("ret:", n << 1, err)
		`,
		"Hello 134.1\nHello 134.1\nret: 24 <nil>\n",
	)
}

func TestGoPackage(t *testing.T) {
	cltest.Expect(t, `
		import "fmt"
		import gostrings "strings"

		x := gostrings.NewReplacer("?", "!").Replace("hello, world???")
		fmt.Println("x: " + x)
		`,
		"x: hello, world!!!\n",
	)
}

func TestSlice(t *testing.T) {
	cltest.Expect(t, `
		x := []float64{1, 2.3, 3.6}
		println("x:", x)
		`,
		"x: [1 2.3 3.6]\n",
	)
	cltest.Expect(t, `
		x := []float64{1, 2: 3.4, 5}
		println("x:", x)
		`,
		"x: [1 0 3.4 5]\n",
	)
}

func TestArray(t *testing.T) {
	cltest.Expect(t, `
		x := [4]float64{1, 2.3, 3.6}
		println("x:", x)

		y := [...]float64{1, 2.3, 3.6}
		println("y:", y)
		`,
		"x: [2.3 3.6 0 0]\ny: [1 2.3 3.6]\n",
	)
	cltest.Expect(t, `
		x := [...]float64{1, 3: 3.4, 5}
		x[1] = 217
		println("x:", x, "x[1]:", x[1])
		`,
		"x: [1 217 0 3.4 5] x[1]: 217\n",
	)
	cltest.Expect(t, `
		x := [...]float64{1, 2.3, 3, 4}
		x[2] = 3.1
		println("x[1:]:", x[1:])
		println(len(x))
	`,
		"x[1:]: [2.3 3.1 4]\n4\n")
}

func TestMap(t *testing.T) {
	cltest.Expect(t, `
		x := map[string]float64{"Hello": 1, "xsw": 3.4}
		println("x:", x)
		`,
		"x: map[Hello:1 xsw:3.4]\n")
}

func TestMapLit(t *testing.T) {
	cltest.Expect(t, `
		x := {"Hello": 1, "xsw": 3.4}
		println("x:", x)
		`,
		"x: map[Hello:1 xsw:3.4]\n",
	)
	cltest.Expect(t, `
		x := {"Hello": 1, "xsw": "3.4"}
		println("x:", x)

		println("empty map:", {})
		`,
		"x: map[Hello:1 xsw:3.4]\nempty map: map[]\n",
	)
}

func TestMapIdx(t *testing.T) {
	cltest.Expect(t, `
		x := {"Hello": 1, "xsw": "3.4"}
		y := {1: "glang", 5: "Hi"}
		i := 1
		q := "Q"
		key := "xsw"
		x["xsw"], y[i] = 3.1415926, q
		println("x:", x, "y:", y)
		println("x[key]:", x[key], "y[1]:", y[1])
		`,
		"x: map[Hello:1 xsw:3.1415926] y: map[1:Q 5:Hi]\nx[key]: 3.1415926 y[1]: Q\n",
	)
}

func TestSliceLit(t *testing.T) {
	cltest.Expect(t, `
		x := [1, 3.4]
		println("x:", x)

		y := [1]
		println("y:", y)

		z := [1+2i, "xsw"]
		println("z:", z)

		println("empty slice:", [])
		`,
		"x: [1 3.4]\ny: [1]\nz: [(1+2i) xsw]\nempty slice: []\n")
}

func TestSliceIdx(t *testing.T) {
	cltest.Expect(t, `
		x := [1, 3.4, 17]
		n, m := 1, uint16(0)
		x[1] = 32.7
		x[m] = 36.86
		println("x:", x[2], x[m], x[n])
		`,
		"x: 17 36.86 32.7\n")
}

func TestListComprehension(t *testing.T) {
	cltest.Expect(t, `
		y := [i+x for i, x <- [1, 2, 3, 4]]
		println("y:", y)
		`,
		"y: [1 3 5 7]\n")
	cltest.Call(t, `
		y := [i+x for i, x <- {3: 1, 5: 2, 7: 3, 11: 4}]
		println("y:", y)
		`, -2).Equal(15)
	cltest.Call(t, `
		y := [i+x for i, x <- {3: 1, 5: 2, 7: 3, 11: 4}, x % 2 == 1]
		println("y:", y)
		`, -2).Equal(10)
}

func TestMapComprehension(t *testing.T) {
	cltest.Expect(t, `
		y := {x: i for i, x <- [3, 5, 7, 11, 13]}
		println("y:", y)
		`,
		"y: map[3:0 5:1 7:2 11:3 13:4]\n",
	)
	cltest.Expect(t, `
		y := {x: i for i, x <- [3, 5, 7, 11, 13], i % 2 == 1}
		println("y:", y)
		`,
		"y: map[5:1 11:3]\n",
	)
	cltest.Expect(t, `
		y := {v: k for k, v <- {"Hello": "xsw", "Hi": "glang"}}
		println("y:", y)
		`,
		"y: map[glang:Hi xsw:Hello]\n",
	)
	cltest.Expect(t, `
		println({x: i for i, x <- [3, 5, 7, 11, 13]})
		println({x: i for i, x <- [3, 5, 7, 11, 13]})
		`,
		"map[3:0 5:1 7:2 11:3 13:4]\nmap[3:0 5:1 7:2 11:3 13:4]\n",
	)
	cltest.Expect(t, `
		arr := [1, 2, 3, 4, 5, 6]
		x := [[a, b] for a <- arr, a < b for b <- arr, b > 2]
		println("x:", x)
		`,
		"x: [[1 3] [2 3] [1 4] [2 4] [3 4] [1 5] [2 5] [3 5] [4 5] [1 6] [2 6] [3 6] [4 6] [5 6]]\n")
}

func TestErrWrapExpr(t *testing.T) {
	cltest.Call(t, `
		x := println("Hello qiniu")!
		x
		`).Equal(12)
	cltest.Call(t, `
		import (
			"strconv"
		)
	
		func add(x, y string) (int, error) {
			return strconv.Atoi(x)? + strconv.Atoi(y)?, nil
		}
	
		x := add("100", "23")!
		x
		`).Equal(123)
}

func TestRational(t *testing.T) {
	cltest.Call(t, `
		x := 3/4r + 5/7r
		x
	`).Equal(big.NewRat(41, 28))
	cltest.Call(t, `
		a := 3/4r
		x := a + 5/7r
		x
	`).Equal(big.NewRat(41, 28))
	y, _ := new(big.Float).SetString(
		"3.14159265358979323846264338327950288419716939937510582097494459")
	y.Mul(y, big.NewFloat(2))
	cltest.Call(t, `
		y := 3.14159265358979323846264338327950288419716939937510582097494459r
		y *= 2
		y
	`).Equal(y)
	cltest.Call(t, `
		a := 3/4r
		b := 5/7r
		if a > b {
			a = a + 1
		}
		a
	`).Equal(big.NewRat(7, 4))
	cltest.Call(t, `
		x := 1/3r + 1r*2r
		x
	`).Equal(big.NewRat(7, 3))
}

type testData struct {
	clause string
	want   string
	panic  bool
}

var testDeleteClauses = map[string]testData{
	"delete_int_key": {`
					m:={1:1,2:2}
					delete(m,1)
					println(m)
					delete(m,3)
					println(m)
					delete(m,2)
					println(m)
					`, "map[2:2]\nmap[2:2]\nmap[]\n", false},
	"delete_string_key": {`
					m:={"hello":1,"Go+":2}
					delete(m,"hello")
					println(m)
					delete(m,"hi")
					println(m)
					delete(m,"Go+")
					println(m)
					`, "map[Go+:2]\nmap[Go+:2]\nmap[]\n", false},
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
					`, "map[Go+:2]\nmap[Go+:2]\nmap[]\n", false},
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
					`, "map[Go+:2]\nmap[Go+:2]\nmap[]\n", false},
	"delete_no_key_panic": {`
					m:={"hello":1,"Go+":2}
					delete(m)
					`, "", true},
	"delete_multi_key_panic": {`
					m:={"hello":1,"Go+":2}
					delete(m,"hi","hi")
					`, "", true},
	"delete_not_map_panic": {`
					m:=[1,2,3]
					delete(m,1)
					`, "", true},
}

func TestDelete(t *testing.T) {
	testScripts(t, "TestDelete", testDeleteClauses)
}

// -----------------------------------------------------------------------------

var testCopyClauses = map[string]testData{
	"copy_int": {`
					a:=[1,2,3]
					b:=[4,5,6]
					n:=copy(b,a)
					println(n)
					println(b)
					`, "3\n[1 2 3]\n", false},
	"copy_string": {`
					a:=["hello"]
					b:=["hi"]
					n:=copy(b,a)
					println(n)
					println(b)
					`, "1\n[hello]\n", false},
	"copy_byte_string": {`
					a:=[byte(65),byte(66),byte(67)]
					println(string(a))
					n:=copy(a,"abc")
					println(n)
					println(a)
					println(string(a))
					`, "ABC\n3\n[97 98 99]\nabc\n", false},
	"copy_first_not_slice_panic": {`
					a:=1
					b:=[1,2,3]
					copy(a,b)
					println(a)
					`, "", true},
	"copy_second_not_slice_panic": {`
					a:=1
					b:=[1,2,3]
					copy(b,a)
					println(b)
					`, "", true},
	"copy_one_args_panic": {`
					a:=[1,2,3]
					copy(a)
					println(a)
					`, "", true},
	"copy_multi_args_panic": {`
					a:=[1,2,3]
					copy(a,a,a)
					println(a)
					`, "", true},
	"copy_string_panic": {`
					a:=[65,66,67]
					copy(a,"abc")
					println(a)
					`, "", true},
	"copy_different_type_panic": {`
					a:=[65,66,67]
					b:=[1.2,1.5,1.7]
					copy(b,a)
					copy(b,a)
					println(b)
					`, "", true},
	"copy_with_operation": {`
					a:=[65,66,67]
					b:=[1]
					println(copy(a,b)+copy(b,a)==2)
					`, "true\n", false},
}

func TestCopy(t *testing.T) {
	testScripts(t, "TestCopy", testCopyClauses)
}

func testScripts(t *testing.T, testName string, scripts map[string]testData) {
	for name, script := range scripts {
		t.Log("Run " + testName + "---" + name)
		var panicMsg []interface{}
		if script.panic {
			panicMsg = append(panicMsg, nil)
		}
		cltest.Expect(t, script.clause, script.want, panicMsg...)
	}
}

// -----------------------------------------------------------------------------
