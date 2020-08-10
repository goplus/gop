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

var testAssign = `
	x, y := 123, "Hello"
	x
	y
`

func TestAssign(t *testing.T) {
	cltest.Call(t, testAssign).Equal("Hello")
	cltest.Call(t, testAssign, -2).Equal(123)
}

// -----------------------------------------------------------------------------

func TestSwitch(t *testing.T) {
	testScripts(t, "TestSwitch", testSwitchIfScripts)
}

var testSwitchIfScripts = map[string]testData{
	"switch_into_case":            {testSwif, "5\n", false},
	"switch_into_default":         {testSwif2, "7\n", false},
	"switch_into_case_with_cond":  {testSw, "5\n", false},
	"switch_into_case_with_empty": {testSw2, "5\n", false},
	"switch_with_no_case":         {testSw3, "7\n", false},
}

var testSwif = `
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
	println(x)
`

var testSwif2 = `
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
	println(x)
`

var testSw = `
	x := 0
	switch t := "Hello"; t {
	case "xsw":
		x = 3
	case "Hello", "world":
		x = 5
	default:
		x= 7
	}
	println(x)

`

var testSw2 = `
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
	println(x)
x
`

var testSw3 = `
	x := 0
	t := "Hello"
	switch t {
	default:
		x = 7
	}
	println(x)
`

// -----------------------------------------------------------------------------

var testIfScripts = map[string]testData{
	"if_with_else":    {testIf, "5\n", false},
	"if_without_else": {testIf2, "3\n", false},
}

var testIf = `
	x := 0
	if t := false; t {
		x = 3
	} else {
		x = 5
	}
	println(x)
`

var testIf2 = `
	x := 5
	if true {
		x = 3
	}
		println(x)
`

func TestIf(t *testing.T) {
	testScripts(t, "TestIf", testIfScripts)
}

// -----------------------------------------------------------------------------

var testReturn = `
	import (
		"fmt"
		"strings"
	)

	func foo(x string) string {
		return strings.NewReplacer("?", "!").Replace(x)
	}

	fmt.Println(foo("Hello, world???"))
`

func TestReturn(t *testing.T) {
	cltest.Expect(t, testReturn, "Hello, world!!!\n")
}

// -----------------------------------------------------------------------------

var testReturn2 = `
	func max(a, b int) int {
		if a < b {
			a = b
		}
		return a
	}

	println("max(23,345):", max(23,345))
`

func TestReturn2(t *testing.T) {
	cltest.Expect(t, testReturn2, "max(23,345): 345\n")
}

// -----------------------------------------------------------------------------

var testFunc = `
	import "fmt"

	func foo(x string) (n int, err error) {
		n, err = fmt.Println("x: " + x)
		return
	}

	foo("Hello, world!")
`

func TestFunc(t *testing.T) {
	cltest.Expect(t, testFunc, "x: Hello, world!\n")
}

// -----------------------------------------------------------------------------

var testFuncv = `
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
`

func TestFuncv(t *testing.T) {
	cltest.Expect(t, testFuncv, "Hello, glang!\n"+"Hello, 123!\n"+"12 <nil>\n")
}

// -----------------------------------------------------------------------------

var testClosure = `
	import "fmt"

	foo := func(prompt string) (n int, err error) {
		n, err = fmt.Println(prompt + x)
		return
	}

	x := "Hello, world!"
	foo("x: ")
`

func TestClosure(t *testing.T) {
	cltest.Expect(t, testClosure, "x: Hello, world!\n")
}

// -----------------------------------------------------------------------------

var testClosurev = `
	import "fmt"

	foo := func(format string, args ...interface{}) (n int, err error) {
		n, err = fmt.Printf(format, args...)
		return
	}

	foo("Hello, %v!\n", "xsw")
`

func TestClosurev(t *testing.T) {
	cltest.Expect(t, testClosurev, "Hello, xsw!\n")

}

// -----------------------------------------------------------------------------

var testForPhraseStmt = `
	sum := 0
	for x <- [1, 3, 5, 7, 11, 13, 17], x > 3 {
		sum += x
	}
	sum
`

func TestForPhraseStmt(t *testing.T) {
	cltest.Call(t, testForPhraseStmt).Equal(53)
}

// -----------------------------------------------------------------------------

var testForPhraseStmt2 = `
	sum := 0
	for x <- [1, 3, 5, 7, 11, 13, 17] {
		if x > 3 {
			sum += x
		}
	}
	sum
`

func TestForPhraseStmt2(t *testing.T) {
	cltest.Call(t, testForPhraseStmt2).Equal(53)
}

// -----------------------------------------------------------------------------

var testForPhraseStmt3 = `
	fns := make([]func() int, 3)
	for i, x <- [3, 15, 777] {
		v := x
		fns[i] = func() int {
			return v
		}
	}
	println("values:", fns[0](), fns[1](), fns[2]())
`

func TestForPhraseStmt3(t *testing.T) {
	cltest.Expect(t, testForPhraseStmt3, "values: 3 15 777\n")
}

// -----------------------------------------------------------------------------

var testForRangeClauses = map[string]testData{
	"no_kv_range_list": {`sum:=0
					for range [1,3,5,7] {
						sum++
					}
					println(sum)
					`, "4\n", false},
	"no_kv_range_map": {`sum:=0
					for range {1:1,2:2,3:3} {
						sum++
					}
					println(sum)
					`, "3\n", false},
	"only_k_range_list": {`sum:=0
					for k :=range [1,3,5,7]{
						sum+=k
					}
					println(sum)
					`, "6\n", false},
	"only_k_range_map": {`sum:=0
					for k :=range {1:1,2:4,3:8,4:16}{
						sum+=k
					}
					println(sum)
					`, "10\n", false},
	"only_v_range_list": {`sum:=0
					for _,v :=range [1,3,5,7]{
						sum+=v
					}
					println(sum)
					`, "16\n", false},
	"only_v_range_map": {`sum:=0
					for _,v :=range {1:1,2:4,3:8,4:16}{
						sum+=v
					}
					println(sum)
					`, "29\n", false},
	"both_kv_range_list": {`sum:=0
					for k,v:=range [1,3,5,7]{
						// 0*1+1*3+2*5+3*7
						sum+=k*v
					}
					println(sum)
					`, "34\n", false},
	"both_kv_range_map": {`sum:=0
					m:={1:2,2:4,3:8}
					for k,v:=range m {
						//1*2+2*4+3*8=34
						sum+=k*v
					}
					println(sum)
					`, "34\n", false},
	"both_kv_assign_simple_range": {` sum:=0
					k,v:=0,0
					for k,v=range [1,2,3,4,5]{
						sum+=k+v
					}
					println(k)
					println(v)
					println(sum)
					`, "4\n5\n25\n", false},
	"both_kv_assign_range_list": {` sum:=0
					m:={1:2,2:4,3:8}
					arr:=[11,22]
					for m[1],m[2]=range arr{
					    sum+=m[1]+m[2]
					}
					println(m[1])
					println(m[2])
					println(sum)
					`, "1\n22\n34\n", false},
	"both_kv_assign_range_map": {` sum:=0
					m:={3:8}
					arr:=[11,22]
					for arr[0],arr[1]=range m{
					    sum+=arr[0]+arr[1]
					}
					println(arr[0])
					println(arr[1])
					println(sum)
					`, "3\n8\n11\n", false},
	"only_v_assign_range": {` sum:=0
					m:={3:8}
					arr:=[11,22]
					for _,arr[1]=range m{
					    sum+=arr[0]+arr[1]
					}
					println(arr[0])
					println(arr[1])
					println(sum)
					`, "11\n8\n19\n", false},
	"only_k_assign_range": {` sum:=0
					m:={3:8}
					arr:=[11,22]
					for arr[0],_=range m{
					    sum+=arr[0]+arr[1]
					}
					println(arr[0])
					println(arr[1])
					println(sum)
					`, "3\n22\n25\n", false},
	"none_kv_assign_range": {` sum:=0
					m:={3:8}
					arr:=[11,22]
					for _,_=range m{
					    sum+=arr[0]+arr[1]
					}
					println(arr[0])
					println(arr[1])
					println(sum)
					`, "11\n22\n33\n", false},
}

func TestRangeStmt(t *testing.T) {
	testScripts(t, "TestRangeStmt", testForRangeClauses)
}

// -----------------------------------------------------------------------------

var testRangeStmt2 = `
	sum := 0
	for _, x := range [1, 3, 5, 7, 11, 13, 17] {
		if x > 3 {
			sum += x
		}
	}
	sum
`

func TestRangeStmt2(t *testing.T) {
	cltest.Call(t, testRangeStmt2).Equal(53)
}

// -----------------------------------------------------------------------------

var testForStmt = `
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
`

func _TestForStmt(t *testing.T) {
	cltest.Expect(t, testForStmt, "values: 3 15 777\n")
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
					`, "16\n", false},
	"for_with_cond_post": {`
					sum := 0
					arr := [1,3,5,7]
					i := 0
					for ; i < len(arr); i+=2 {
						sum+=arr[i]
					}
					println(sum)
					`, "6\n", false},
	"for_with_cond": {`
					arr := [1,3,5,7]
					i := 0
					sum := 0
					for ; i < len(arr) && i < 2; {
						sum+=arr[i]
						i++
					}
					println(sum)
					`, "4\n", false},
	"for_with_init_cond": {`
					arr := [1,3,5,7]
					sum := 0
					for i:=0; i < len(arr); {
						sum+=arr[i]
						i++
					}
					println(sum)
					`, "16\n", false},
	"for_with_continue": {`
					arr := [1,3,5,7]
					sum := 0
					for i:=0; i < len(arr);i++ {
						if arr[i]<5{
							continue
						}
						sum+=arr[i]
					}
					println(sum)
					`, "12\n", false},
	"for_with_break": {`
					arr := [1,3,5,7]
					sum := 0
					for i:=0; i < len(arr);i++ {
						if arr[i]>5{
							break
						}
						sum+=arr[i]
					}
					println(sum)
					`, "9\n", false},
	"for_with_break_label": {`
					arr := [1,3,5,7]
					sum := 0
					L:
					for i:=0; i < len(arr);i++ {
						for j:=0;j<len(arr);j++{
							if j>2{
								break L
							}
							sum+=arr[i]+arr[j]
						}
					}
					println(sum)
					`, "12\n", false}, // (1+1)+(1+3)+(1+5)
	"for_with_continue_label": {`
					arr := [1,3,5,7]
					sum := 0
					L:
					for i:=0; i < len(arr);i++ {
						for j:=0;j<len(arr);j++{
							if j>1{
								continue L
							}
							sum+=arr[i]+arr[j]
						}
					}
					println(sum)
					`, "48\n", false}, // (1+3+5+7)*2+(1+3)*4
	"for_with_continue_break": {`
					arr := [1,3,5,7]
					sum := 0
					for i:=0; i < len(arr);i++ {
						if arr[i]>5{
							break
						}
						for j:=0;j<len(arr);j++{
							if arr[j]<5{
								continue
							}
							sum+=arr[j]
						}
					}
					println(sum)
					`, "36\n", false},
	"for_with_continue_break_continue": {`
					arr := [1,3,5,7]
					sum := 0
					L1:
					for i:=0; i < len(arr);i++ {
						if arr[i]>5{
							break
						}
						for j:=i;j<len(arr);j++{
							if arr[j]<5{
								continue L1
							}
							sum+=arr[j]
						}
					}
					println(sum)
					`, "12\n", false},
	"for_with_continue_panic": {`
					arr := [1,3,5,7]
					sum := 0
					for i:=0; i < len(arr);i++ {
					}
					continue
					println(sum)
					`, "", true},
	"for_with_continue_no_label_panic": {`
					arr := [1,3,5,7]
					sum := 0
					for i:=0; i < len(arr);i++ {
						continue L
					}
					println(sum)
					`, "", true},
	"for_with_break_panic": {`
					arr := [1,3,5,7]
					sum := 0
					for i:=0; i < len(arr);i++ {
					}
					break
					println(sum)
					`, "", true},
	"for_with_break_label_panic": {`
					arr := [1,3,5,7]
					L:
					sum := 0
					for i:=0; i < len(arr);i++ {
						break L
					}
					println(sum)
					`, "", true},
	"for_with_continue_wrong_label_panic": {`
					arr := [1,3,5,7]
					L:
					sum := 0
					for i:=0; i < len(arr);i++ {
						continue L
					}
					println(sum)
					`, "", true},
	"for_with_many_labels": {`
					arr := [1,3,5,7]
					sum := 0
					L:
					for i:=0; i < len(arr);i++ {
						if arr[i]<7{
								continue L
						}
						L1:
						for j:=0;j<len(arr);j++{
							if arr[j]>1{
								break L1
							}
							sum+=arr[i]+arr[j]
						}
					}
					println(sum)
					`, "8\n", false},
	"for_with_many_labels_break": {`
					arr := [1,3,5,7]
					sum := 0
					L:
					for i:=0; i < len(arr);i++ {
						if arr[i]<7{
								continue L
						}
						L1:
						for j:=0;j<len(arr);j++{
							if arr[j]>3{
								break L
							}
							sum+=arr[i]+arr[j]
						}
					}
					println(sum)
					`, "18\n", false},
}

func TestNormalForStmt(t *testing.T) {
	testScripts(t, "TestNormalForStmt", testNormalForClauses)
}

var testForIncDecStmt = `
	a,b:=10,2
	{a--;a--;a--}
	{b++;b++;b++}
	println(a,b,a*b)
`

func TestForIncDecStmt(t *testing.T) {
	cltest.Expect(t, testForIncDecStmt, "7 5 35\n")
}

// -----------------------------------------------------------------------------

var testSwitchBranchClauses = map[string]testData{
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
					`, "0\n1\n7\n", false},
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
					`, "0\n1\n", false},
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
					`, "0\n1\n", false},
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
					`, "0\n1\n7\n", false},
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
					`, "0\n1\n", false},
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
					`, "", true},
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
					`, "", true},
	"switch_break": {`
					x:=0
					y:=2
					switch x {
					case 0:
						if y>0{
							println(y)
							break
						}
						println(x)
					default:
						x=7
						println(x)
					}
					`, "2\n", false},
	"switch_break_label": {`
					x:=0
					y:=2
					L:
					switch x {
					case 0:
						if y>0{
							println(y)
							break L
						}
						println(x)
					default:
						x=7
						println(x)
					}
					`, "2\n", false},
	"switch_for_continue_label": {`
					x:=0
					y:=2
					L:
					for i:=0;i<5;i++{
						switch i {
						case 0:
							if y>0{
								println(y)
								continue L
							}
							println(x)
						case 1:
							println(x)
							x++
							continue L
						case 2:
							println(x)
							x++
							break
						case 3:
							println(x)
							break L
						case 4:
							println(x)
						default:
							x=7
							println(x)
						}
					}
					`, "2\n0\n1\n2\n", false},
}

func TestSwitchBranchStmt(t *testing.T) {
	testScripts(t, "TestSwitchBranchStmt", testSwitchBranchClauses)
}

// -----------------------------------------------------------------------------

var testGotoLabelClauses = map[string]testData{
	"goto_before_label": {`
					goto L
					println("before")
					L:
					println("over")
					`, "over\n", false},
	"goto_after_label": {`
					i:=0
					L:
						if i < 3 {
							println(i)
							i++
							goto L
						}
					println("over")
					`, "0\n1\n2\nover\n", false},
	"goto_multi_labels": {`
					i:=0
					L:
						if i < 3  {
						goto L1
							println(i)
						L1:
							println(i)
							i++
							if i==4{
								goto L3
							}
							goto L
						}
					L3:
					println("over")
					L4:
					`, "0\n1\n2\nover\n", false},
	"goto_nil_label": {`
					goto;
					println("over")
					`, "", true},
	"goto_not_define_label": {`
					goto L
					println("over")
					`, "", true},
	"goto_illegal_block": {`
					goto L
					{
						L:
						println("L")
					}
					`, "", true},
	"goto_redefine_block": {`
					{
						L:
						println("L")
					}
					{
						L:
						println("L")
					}
					goto L
					`, "", true},
}

func TestGotoLabelStmt(t *testing.T) {
	testScripts(t, "TestGotoLabelStmt", testGotoLabelClauses)
}

// -----------------------------------------------------------------------------

var testRangeStmtWithBranchClauses = map[string]testData{
	"range_with_continue": {clause: `
					arr := [1,3,5,7]
					sum := 0
					for i,_:=range arr {
						if arr[i]<5{
							continue
						}
						sum+=arr[i]
					}
					println(sum)
					`, want: "12\n"},
	"range_with_break": {clause: `
					arr := [1,3,5,7]
					sum := 0
					for i,_:=range arr {
						if arr[i]>5{
							break
						}
						sum+=arr[i]
					}
					println(sum)
					`, want: "9\n"},

	"range_with_continue_break": {clause: `
					arr := [1,3,5,7]
					sum := 0
					for i,_:=range arr {
						if arr[i]>5{
							break
						}
						for j:=0;j<len(arr);j++{
							if arr[j]<5{
								continue
							}
							sum+=arr[j]
						}
					}
					println(sum)
					`, want: "36\n"},
	"range_with_return_value": {clause: `
					func foo() int{
						arr := [1,2,3,4]
						sum := 0
						for i, _ <- arr {
							if arr[i]>0{
								return 1
							}
							sum+=arr[i]
						}
						return 0
					}
					println(foo())
					`, want: "1\n"},
	"range_with_only_return": {clause: `
					func foo() {
						arr := [1,2,3,4]
						sum := 0
						for i, _ <- arr {
							if arr[i]>0{
								println("1")
								return 
							}
							sum+=arr[i]
						}
						println("0")
						return 
					}
					foo()
					`, want: "1\n"},
}

func TestRangeStmtWithBranch(t *testing.T) {
	testScripts(t, "TestRangeStmtWithBranch", testRangeStmtWithBranchClauses)
}

// -----------------------------------------------------------------------------

var testForPhraseWithBranchClauses = map[string]testData{
	"for_phrase_with_continue": {clause: `
					arr := [1,3,5,7]
					sum := 0
					for i, _ <- arr {
						if arr[i]<5{
							continue
						}
						sum+=arr[i]
					}
					println(sum)
					`, want: "12\n"},
	"for_phrase_with_break": {clause: `
					arr := [1,3,5,7]
					sum := 0
					for i, _ <- arr {
						if arr[i]>5{
							break
						}
						sum+=arr[i]
					}
					println(sum)
					`, want: "9\n"},
	"for_phrase_with_continue_break": {clause: `
					arr := [1,3,5,7]
					sum := 0
					for i, _ <- arr {
						if arr[i]>5{
							break
						}
						for j:=0;j<len(arr);j++{
							if arr[j]<5{
								continue
							}
							sum+=arr[j]
						}
					}
					println(sum)
					`, want: "36\n"},
	"for_phrase_with_return_value": {clause: `
					func foo() int{
						arr := [1,2,3,4]
						sum := 0
						for i, _ <- arr {
							if arr[i]>0{
								return 1
							}
							sum+=arr[i]
						}
						return 0
					}
					println(foo())
					`, want: "1\n"},
	"for_phrase_with_only_return": {clause: `
					func foo() {
						arr := [1,2,3,4]
						sum := 0
						for i, _ <- arr {
							if arr[i]>0{
								println("1")
								return 
							}
							sum+=arr[i]
						}
						println("0")
						return 
					}
					foo()
					`, want: "1\n"},
}

func TestForPhraseWithBranch(t *testing.T) {
	testScripts(t, "TestForPhraseWithBranch", testForPhraseWithBranchClauses)
}

// -----------------------------------------------------------------------------

var testRangeMapWithBranchClauses = map[string]testData{
	"map_for_with_continue": {clause: `
					arr := {1:1,2:3,3:5,4:7}
					sum := 0
					for i, _ <- arr {
						if arr[i]<5{
							continue
						}
						sum+=arr[i]
					}
					println(sum)
					`, want: "12\n"},
	"map_for_with_break": {clause: `
					arr := {1:1,2:3,3:5,4:7}
					sum := 0
					for i, _ <- arr {
						if arr[i]>0{
							break
						}
						sum+=arr[i]
					}
					println(sum)
					`, want: "0\n"},

	"map_for_range_continue": {clause: `
					arr := {1:1,2:3,3:5,4:7}
					sum := 0
					for i, _ <- arr {
						if arr[i]<5{
							continue
						}
						sum+=arr[i]
					}
					println(sum)
					`, want: "12\n"},
	"map_for_range_break": {clause: `
					arr := {1:1,2:3,3:5,4:7}
					sum := 0
					for i, _ <- arr {
						if arr[i]>0{
							break
						}
						sum+=arr[i]
					}
					println(sum)
					`, want: "0\n"},
	"map_for_range_return_value": {clause: `
					func foo() int{
						arr := {1:1,2:3,3:5,4:7}
						sum := 0
						for i, _ <- arr {
							if arr[i]>0{
								return 1
							}
							sum+=arr[i]
						}
						return 0
					}
					println(foo())
					`, want: "1\n"},
	"map_for_range_only_return": {clause: `
					func foo() {
						arr := {1:1,2:3,3:5,4:7}
						sum := 0
						for i, _ <- arr {
							if arr[i]>0{
								println("1")
								return 
							}
							sum+=arr[i]
						}
						println("0")
						return 
					}
					foo()
					`, want: "1\n"},
}

func TestMapForPhraseWithBranch(t *testing.T) {
	testScripts(t, "TestMapForPhraseWithBranch", testRangeMapWithBranchClauses)
}

// -----------------------------------------------------------------------------

var testDeferClauses = map[string]testData{
	"multi_defer": {clause: `
	func test() {
		defer println("Hello, test defer!")
		println("Hello, test!")
	}
	
	defer println("Hello, defer1!")
	defer println("Hello, defer2!")
	defer println("Hello, defer3!")
	defer test()
	println("Hello, world!")
		`, want: "Hello, world!\nHello, test!\nHello, test defer!\nHello, defer3!\nHello, defer2!\nHello, defer1!\n"},
	"multi_defer_args": {clause: `
	func test(i int) {
		defer println("Hello, test defer!")
		println("Hello, test!",i)
	}

	defer println("Hello, defer1!")
	defer println("Hello, defer2!")
	defer println("Hello, defer3!")
	defer test(1)
	println("Hello, world!")
		`, want: "Hello, world!\nHello, test! 1\nHello, test defer!\nHello, defer3!\nHello, defer2!\nHello, defer1!\n"},
	"multi_defer_goval": {clause: `
	import "fmt"
	import gostrings "strings"
	func test(i int) {
		defer println("Hello, test defer!")
		println("Hello, test!",i)
	}

    defer println(gostrings.NewReplacer("?", "!").Replace("hello, world???"))
	defer println("Hello, defer1!")
	defer println("Hello, defer2!")
	defer println("Hello, defer3!")
	defer test(1)
	println("Hello, world!")
		`, want: "Hello, world!\nHello, test! 1\nHello, test defer!\nHello, defer3!\nHello, defer2!\nHello, defer1!\nhello, world!!!\n"},

	"unnamed_func": {clause: `
		f := func(i int) {
			defer println(i)
			println("hello")
		}
		defer f(1+1)
		println("hello1")
		`, want: "hello1\nhello\n2\n"},
	"unnamed_func2": {clause: `
		f2:= func(i int,err error) {
			defer println(i)
			println(err)
			println("hello")
		}
		defer f2(println("hello world"))
		println("hello1")
		`, want: "hello world\nhello1\n<nil>\nhello\n12\n"},
	"unnamed_func3": {clause: `
		import (
			"fmt"
		)
		
		myprint := func(format string, a ...interface{}) {
			format = "this is test print: " + format
			fmt.Printf(format, a...)
		}
		
		defer myprint("hello %s\n", "defer")
		myprint("hello %s\n", "world")
		`, want: "this is test print: hello world\nthis is test print: hello defer\n"},
	"unnamed_func6": {clause: `
	import (
		"fmt"
	)

	myprint := func() {
		fmt.Printf("hello defer\n")
	}

	defer myprint()
	printf("hello %s\n", "world")
		`, want: "hello world\nhello defer\n"},
	"unnamed_func7": {clause: `
	import (
		"fmt"
	)
	
	myprint := func() {
		fmt.Printf("hello world\n")
	}
	
	defer fmt.Println("hello defer")
	myprint()
		`, want: "hello world\nhello defer\n"},
}

func TestDeferStmt(t *testing.T) {
	testScripts(t, "TestDeferStmt", testDeferClauses)
}

func TestDeferCopy(t *testing.T) {
	cltest.Expect(t, `
		a := [1, 2, 3]
		b := [4, 5, 6]
		defer func() {
			println(b)
		}()
		defer copy(b, a)
		`,
		"[1 2 3]\n",
	)
}

func TestDefer1(t *testing.T) {
	cltest.Expect(t, `
		defer println(println("hello world"))
		println("test defer")
		`,
		"hello world\ntest defer\n12 <nil>\n",
	)
}

func TestDeferInFunc(t *testing.T) {
	cltest.Expect(t, `
		func f() {
			defer println(println("hello world"))
			println("test defer")
		}
		defer println("hello main")
		f()
		f()
		`,
		"hello world\ntest defer\n12 <nil>\n"+
			"hello world\ntest defer\n12 <nil>\n"+
			"hello main\n",
	)
}

func TestDeferInClosure(t *testing.T) {
	cltest.Expect(t, `
		defer func() {
			defer println(println("hello world"))
			println("test defer")
		}()
		`,
		"hello world\ntest defer\n12 <nil>\n",
	)
}

func TestDeferInFor(t *testing.T) {
	cltest.Expect(t, `
		for i:=0;i<10;i++ {
			defer println(i)
		}
		println("test defer")
		`,
		"test defer\n9\n8\n7\n6\n5\n4\n3\n2\n1\n0\n",
	)
}

func TestDefer4(t *testing.T) {
	cltest.Expect(t, `
		import "strings"

		f2 := func(s string, s2 string) {
			defer func() {
				println(s)
				println(s2)
			}()
			s = strings.Replace(s, "?", "!", -1)
			s2 = strings.Replace(s2, "?", "!", -1)
		}
		s := "hello world???"
		s2 := "hello world????"
		defer f2(s, s2)

		println(s)
		println(s2)
		`,
		"hello world???\nhello world????\nhello world!!!\nhello world!!!!\n",
	)
}

func TestDefer5(t *testing.T) {
	cltest.Expect(t, `
		import "fmt"

		defer func(format string, a ...interface{}) {
			format = format
			fmt.Printf(format, a...)
		}("hello %s\n", "defer")

		printf("hello %s\n", "world")
		`,
		"hello world\nhello defer\n",
	)
}

func TestDefer6(t *testing.T) {
	cltest.Expect(t, `
		func f() (x int) {
			defer func() {
				x = 3
			}()
			return 1
		}
		println(f())
		`,
		"3\n",
	)
}

func TestDefer7(t *testing.T) {
	cltest.Expect(t, `
		func h() (x int) {
			for i <- [3, 2, 1] {
				v := i
				defer func() {
					x = v
				}()
			}
			return
		}
		println(h())
		`,
		"3\n",
	)
}

// -----------------------------------------------------------------------------

func TestGo(t *testing.T) {
	cltest.Expect(t, `
		import "time"

		n := 1
		go func() {
			println("Hello, goroutine!")
			n++
			println(n)
		}()
		time.Sleep(1e8)
		`,
		"Hello, goroutine!\n2\n",
	)
}

// -----------------------------------------------------------------------------

var testVarScopeClauses = map[string]testData{
	"variable_redefinition_#issue304": {`
						a := []float64{1, 2, 3.4}
						println(a)
						
						a := []float64{1, 2, 3.4}
						println(a)
					`, "", true},
	"variable_scope_bug_#issue303": {`
						i:=0
						{
						    i:=0
						    i++
						    println("inner is",i)
						}
						println("outer is",i)
					`, "inner is 1\nouter is 0\n", false},
}

func TestVarScopeStmt(t *testing.T) {
	testScripts(t, "TestVarScopeStmt", testVarScopeClauses)
}
