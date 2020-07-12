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
	"testing"

	"github.com/goplus/gop/cl/cltest"
)

// -----------------------------------------------------------------------------

func TestUnbound(t *testing.T) {
	cltest.Expect(t,
		`println("Hello " + "qiniu:", 123, 4.5, 7i)`,
		"Hello qiniu: 123 4.5 (0+7i)\n",
	)
}

func TestPanic(t *testing.T) {
	cltest.Expect(t,
		`panic("Helo")`, "", "Helo",
	)
}

func TestTypeCast(t *testing.T) {
	cltest.Call(t, `
	x := []byte("hello")
	x
	`).Equal([]byte("hello"))
}

func TestMake(t *testing.T) {
	cltest.Expect(t, `
	    a := make([]int, 0, 4)
		a = append(a, 1, 2, 3)
		println(a)`,
		"[1 2 3]\n",
	)
}

func TestMake2(t *testing.T) {
	cltest.Expect(t, `
		a := make([]int, 0, 4)
		a = append(a, [1, 2, 3]...)
		println(a)`,
		"[1 2 3]\n",
	)
}

func TestMake3(t *testing.T) {
	cltest.Expect(t, `
		n := 4
		a := make(map[string]interface{}, uint16(n))
		println(a)`,
		"map[]\n",
	)
}

func TestMake4(t *testing.T) {
	cltest.Expect(t, `
		import "reflect"

		a := make(chan *func(), uint16(4))
		println(reflect.TypeOf(a))`,
		"chan *func()\n")
}

func TestOperator(t *testing.T) {
	cltest.Expect(t, `
		println("Hello", 123 * 4.5, 1 + 7i)`,
		"Hello 553.5 (1+7i)\n")
}

func TestVar(t *testing.T) {
	cltest.Expect(t, `
		x := 123.1
		println("Hello", x)`,
		"Hello 123.1\n")
}

func TestVarOp(t *testing.T) {
	cltest.Expect(t, `
		x := 123.1
		y := 1 + x
		println("Hello", y + 10)
		n, err := println("Hello", y + 10)
		println("ret:", n << 1, err)`,
		"Hello 134.1\nHello 134.1\nret: 24 <nil>\n")
}

func TestGoPackage(t *testing.T) {
	cltest.Expect(t, `
		import "fmt"
		import gostrings "strings"

		x := gostrings.NewReplacer("?", "!").Replace("hello, world???")
		fmt.Println("x: " + x)`,
		"x: hello, world!!!\n")
}

func TestSlice(t *testing.T) {
	cltest.Expect(t, `
		x := []float64{1, 2.3, 3.6}
		println("x:", x)`,
		"x: [1 2.3 3.6]\n")
}

func TestSlice2(t *testing.T) {
	cltest.Expect(t, `
		x := []float64{1, 2: 3.4, 5}
		println("x:", x)`,
		"x: [1 0 3.4 5]\n")
}
func TestArray(t *testing.T) {
	cltest.Expect(t, `
		x := [4]float64{1, 2.3, 3.6}
		println("x:", x)

		y := [...]float64{1, 2.3, 3.6}
		println("y:", y)`, "x: &[2.3 3.6 0 0]\ny: &[1 2.3 3.6]\n")
}

func TestArray2(t *testing.T) {
	cltest.Expect(t, `
		x := [...]float64{1, 3: 3.4, 5}
		x[1] = 217
		println("x:", x, "x[1]:", x[1])`,
		"x: &[1 217 0 3.4 5] x[1]: 217\n")
}

func TestMap(t *testing.T) {
 	cltest.Expect(t, `
	x := map[string]float64{"Hello": 1, "xsw": 3.4}
	println("x:", x)`,
		"x: map[Hello:1 xsw:3.4]\n")
}

func TestMapLit(t *testing.T) {
	cltest.Expect(t, `
		x := {"Hello": 1, "xsw": 3.4}
		println("x:", x)`, 
		"x: map[Hello:1 xsw:3.4]\n")
}

func TestMapLit2(t *testing.T) {
	cltest.Expect(t, `
		x := {"Hello": 1, "xsw": "3.4"}
		println("x:", x)

		println("empty map:", {})`, 
		"x: map[Hello:1 xsw:3.4]\nempty map: map[]\n")
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
		println("x[key]:", x[key], "y[1]:", y[1])`, 
		"x: map[Hello:1 xsw:3.1415926] y: map[1:Q 5:Hi]\nx[key]: 3.1415926 y[1]: Q\n")
}

func TestSliceLit(t *testing.T) {
	cltest.Expect(t, `
		x := [1, 3.4]
		println("x:", x)

		y := [1]
		println("y:", y)

		z := [1+2i, "xsw"]
		println("z:", z)

		println("empty slice:", [])`, 
		"x: [1 3.4]\ny: [1]\nz: [(1+2i) xsw]\nempty slice: []\n")
}

func TestSliceIdx(t *testing.T) {
	cltest.Expect(t, `
	x := [1, 3.4, 17]
	n, m := 1, uint16(0)
	x[1] = 32.7
	x[m] = 36.86
	println("x:", x[2], x[m], x[n])`, 
	"x: 17 36.86 32.7\n")
}

func TestListComprehension(t *testing.T) {
	cltest.Expect(t, `
		y := [i+x for i, x <- [1, 2, 3, 4]]
		println("y:", y)`, 
		"y: [1 3 5 7]\n")
}

func TestListComprehension2(t *testing.T) {
	cltest.Expect(t, `
		y := [i+x for i, x <- {3: 1, 5: 2, 7: 3, 11: 4}]
		println("y:", y)`, 
		"y: [7 10 15 4]\n")
}

func TestListComprehensionFilter(t *testing.T) {
	cltest.Expect(t, `
		y := [i+x for i, x <- {3: 1, 5: 2, 7: 3, 11: 4}, x % 2 == 1]
		println("y:", y)`, 
		"y: [4 10]\n")
}

func TestMapComprehension(t *testing.T) {
	cltest.Expect(t, `
		y := {x: i for i, x <- [3, 5, 7, 11, 13]}
		println("y:", y)`, 
		"y: map[3:0 5:1 7:2 11:3 13:4]\n")
}

func TestMapComprehensionFilter(t *testing.T) {
	cltest.Expect(t, `
		y := {x: i for i, x <- [3, 5, 7, 11, 13], i % 2 == 1}
		println("y:", y)`, 
		"y: map[5:1 11:3]\n")
}

func TestMapComprehension2(t *testing.T) {
	cltest.Expect(t, `
		y := {v: k for k, v <- {"Hello": "xsw", "Hi": "glang"}}
		println("y:", y)`, 
		"y: map[glang:Hi xsw:Hello]\n")
}

func TestMapComprehension3(t *testing.T) {
	cltest.Expect(t, `
		println({x: i for i, x <- [3, 5, 7, 11, 13]})
		println({x: i for i, x <- [3, 5, 7, 11, 13]})`, 
		"map[3:0 5:1 7:2 11:3 13:4]\nmap[3:0 5:1 7:2 11:3 13:4]\n")
}

func TestMapComprehension4(t *testing.T) {
	cltest.Expect(t, `
		arr := [1, 2, 3, 4, 5, 6]
		x := [[a, b] for a <- arr, a < b for b <- arr, b > 2]
		println("x:", x)`, 
		"x: [[1 3] [2 3] [1 4] [2 4] [3 4] [1 5] [2 5] [3 5] [4 5] [1 6] [2 6] [3 6] [4 6] [5 6]]\n")
}

func TestErrWrapExpr(t *testing.T) {
	cltest.Expect(t, `
		x := println("Hello qiniu")!
		x`, 
		"Hello qiniu\n")
}

func TestErrWrapExpr2(t *testing.T) {
	cltest.Expect(t, `
		import (
			"strconv"
		)
	
		func add(x, y string) (int, error) {
			return strconv.Atoi(x)? + strconv.Atoi(y)?, nil
		}
	
		x := add("100", "23")!
		x`, 
		"")
}

func TestRational(t *testing.T) {
	cltest.Expect(t, `
		x := 3/4r + 5/7r
		x`, 
		"")
}

func TestRational2(t *testing.T) {
	cltest.Expect(t, `
		a := 3/4r
		x := a + 5/7r
		x`, "")
}

func TestRational3(t *testing.T) {
	cltest.Expect(t, `
		y := 3.14159265358979323846264338327950288419716939937510582097494459r
		y *= 2
		y`, "")
}

func  TestRational4(t *testing.T)  {
	cltest.Expect(t, `
		a := 3/4r
		b := 5/7r
		if a > b {
			a = a + 1
		}
		a`, "")
}

func TestRational5(t *testing.T)  {
	cltest.Expect(t, `
		x := 1/3r + 1r*2r
		x`, "")
}
