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
	"switch_into_case":            {testSwif, []string{"5"}},
	"switch_into_default":         {fsTestSwif2, []string{"7"}},
	"switch_into_case_with_cond":  {fsTestSw, []string{"5"}},
	"switch_into_case_with_empty": {fsTestSw2, []string{"5"}},
	"switch_with_no_case":         {fsTestSw3, []string{"7"}},
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

var fsTestSwif2 = `
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

var fsTestSw = `
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

var fsTestSw2 = `
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

var fsTestSw3 = `
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
	"if_with_else":    {fsTestIf, []string{"5"}},
	"if_without_else": {fsTestIf2, []string{"3"}},
}

var fsTestIf = `
	x := 0
	if t := false; t {
		x = 3
	} else {
		x = 5
	}
	println(x)
`

var fsTestIf2 = `
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

var fsTestReturn = `
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
	cltest.Expect(t, fsTestReturn, "Hello, world!!!\n")
}

// -----------------------------------------------------------------------------

var fsTestReturn2 = `
	func max(a, b int) int {
		if a < b {
			a = b
		}
		return a
	}

	println("max(23,345):", max(23,345))
`

func TestReturn2(t *testing.T) {
	cltest.Expect(t, fsTestReturn2, "max(23,345): 345\n")
}

// -----------------------------------------------------------------------------

var fsTestFunc = `
	import "fmt"

	func foo(x string) (n int, err error) {
		n, err = fmt.Println("x: " + x)
		return
	}

	foo("Hello, world!")
`

func TestFunc(t *testing.T) {
	cltest.Expect(t, fsTestFunc, "x: Hello, world!\n")
}

// -----------------------------------------------------------------------------

var fsTestFuncv = `
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
	cltest.Expect(t, fsTestFuncv, "Hello, glang!\n"+"Hello, 123!\n"+"12 <nil>\n")
}

// -----------------------------------------------------------------------------

var fsTestClosure = `
	import "fmt"

	foo := func(prompt string) (n int, err error) {
		n, err = fmt.Println(prompt + x)
		return
	}

	x := "Hello, world!"
	foo("x: ")
`

func TestClosure(t *testing.T) {
	cltest.Expect(t, fsTestClosure, "x: Hello, world!\n")
}

// -----------------------------------------------------------------------------

var fsTestClosurev = `
	import "fmt"

	foo := func(format string, args ...interface{}) (n int, err error) {
		n, err = fmt.Printf(format, args...)
		return
	}

	foo("Hello, %v!\n", "xsw")
`

func TestClosurev(t *testing.T) {
	cltest.Expect(t, fsTestClosurev, "Hello, xsw!\n")

}

// -----------------------------------------------------------------------------

var fsTestForPhraseStmt = `
	sum := 0
	for x <- [1, 3, 5, 7, 11, 13, 17], x > 3 {
		sum += x
	}
	sum
`

func TestForPhraseStmt(t *testing.T) {
	cltest.Call(t, fsTestForPhraseStmt).Equal(53)
}

// -----------------------------------------------------------------------------

var fsTestForPhraseStmt2 = `
	sum := 0
	for x <- [1, 3, 5, 7, 11, 13, 17] {
		if x > 3 {
			sum += x
		}
	}
	sum
`

func TestForPhraseStmt2(t *testing.T) {
	cltest.Call(t, fsTestForPhraseStmt2).Equal(53)
}

// -----------------------------------------------------------------------------

var fsTestForPhraseStmt3 = `
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
	cltest.Expect(t, fsTestForPhraseStmt3, "values: 3 15 777\n")
}

// -----------------------------------------------------------------------------

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
	testScripts(t, "TestRangeStmt", testForRangeClauses)
}

// -----------------------------------------------------------------------------

var fsTestRangeStmt2 = `
	sum := 0
	for _, x := range [1, 3, 5, 7, 11, 13, 17] {
		if x > 3 {
			sum += x
		}
	}
	sum
`

func TestRangeStmt2(t *testing.T) {
	cltest.Call(t, fsTestRangeStmt2).Equal(53)
}

// -----------------------------------------------------------------------------

var fsTestForStmt = `
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
	cltest.Expect(t, fsTestForStmt, "values: 3 15 777\n")
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
					`, []string{"12"}},
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
					`, []string{"9"}},
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
					`, []string{"12"}}, // (1+1)+(1+3)+(1+5)
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
					`, []string{"48"}}, // (1+3+5+7)*2+(1+3)*4
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
					`, []string{"36"}},
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
					`, []string{"12"}},
	"for_with_continue_panic": {`
					arr := [1,3,5,7]
					sum := 0
					for i:=0; i < len(arr);i++ {
					}
					continue
					println(sum)
					`, []string{"_panic"}},
	"for_with_continue_no_label_panic": {`
					arr := [1,3,5,7]
					sum := 0
					for i:=0; i < len(arr);i++ {
						continue L
					}
					println(sum)
					`, []string{"_panic"}},
	"for_with_break_panic": {`
					arr := [1,3,5,7]
					sum := 0
					for i:=0; i < len(arr);i++ {
					}
					break
					println(sum)
					`, []string{"_panic"}},
	"for_with_break_label_panic": {`
					arr := [1,3,5,7]
					L:
					sum := 0
					for i:=0; i < len(arr);i++ {
						break L
					}
					println(sum)
					`, []string{"_panic"}},
	"for_with_continue_wrong_label_panic": {`
					arr := [1,3,5,7]
					L:
					sum := 0
					for i:=0; i < len(arr);i++ {
						continue L
					}
					println(sum)
					`, []string{"_panic"}},
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
					`, []string{"8"}},
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
					`, []string{"18"}},
}

func TestNormalForStmt(t *testing.T) {
	testScripts(t, "TestNormalForStmt", testNormalForClauses)
}

var fsTestForIncDecStmt = `
	a,b:=10,2
	{a--;a--;a--}
	{b++;b++;b++}
	println(a,b,a*b)
`

func TestForIncDecStmt(t *testing.T) {
	cltest.Expect(t, fsTestForIncDecStmt, "7 5 35\n")
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
					`, []string{"_panic"}},
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
					`, []string{"_panic"}},
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
					`, []string{"2"}},
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
					`, []string{"2"}},
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

					`, []string{"2", "0", "1", "2"}},
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
					`, []string{"over"}},
	"goto_after_label": {`
					i:=0
					L:
						if i < 3 {
							println(i)
							i++
							goto L
						}
					println("over")
					`, []string{"0", "1", "2", "over"}},
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
					`, []string{"0", "1", "2", "over"}},
	"goto_nil_label": {`
					goto;
					println("over")
					`, []string{"_panic"}},
	"goto_not_define_label": {`
					goto L
					println("over")
					`, []string{"_panic"}},
	"goto_illegal_block": {`
					goto L
					{
						L:
						println("L")
					}
					`, []string{"_panic"}},
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
					`, []string{"_panic"}},
}

func TestGotoLabelStmt(t *testing.T) {
	testScripts(t, "TestGotoLabelStmt", testGotoLabelClauses)
}

// -----------------------------------------------------------------------------

var testRangeStmtWithBranchClauses = map[string]testData{
	"range_with_continue": {`
					arr := [1,3,5,7]
					sum := 0
					for i,_:=range arr {
						if arr[i]<5{
							continue
						}
						sum+=arr[i]
					}
					println(sum)
					`, []string{"12"}},
	"range_with_break": {`
					arr := [1,3,5,7]
					sum := 0
					for i,_:=range arr {
						if arr[i]>5{
							break
						}
						sum+=arr[i]
					}
					println(sum)
					`, []string{"9"}},
	"range_with_break_label": {`
					arr := [1,3,5,7]
					sum := 0
					L:
					for i,_:=range arr {
						for j:=0;j<len(arr);j++{
							if j>2{
								break L
							}
							sum+=arr[i]+arr[j]
						}
					}
					println(sum)
					`, []string{"12"}}, // (1+1)+(1+3)+(1+5)
	"range_with_continue_label": {`
					arr := [1,3,5,7]
					sum := 0
					L:
					for i,_:=range arr {
						for j:=0;j<len(arr);j++{
							if j>1{
								continue L
							}
							sum+=arr[i]+arr[j]
						}
					}
					println(sum)
					`, []string{"48"}}, // (1+3+5+7)*2+(1+3)*4
	"range_with_continue_break": {`
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
					`, []string{"36"}},
	"range_with_continue_break_continue": {`
					arr := [1,3,5,7]
					sum := 0
					L1:
					for i,_:=range arr {
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
					`, []string{"12"}},
	"range_with_continue_panic": {`
					arr := [1,3,5,7]
					sum := 0
					for i,_:=range arr {
					}
					continue
					println(sum)
					`, []string{"_panic"}},
	"range_with_continue_no_label_panic": {`
					arr := [1,3,5,7]
					sum := 0
					for i,_:=range arr {
						continue L
					}
					println(sum)
					`, []string{"_panic"}},
	"range_with_break_panic": {`
					arr := [1,3,5,7]
					sum := 0
					for i,_:=range arr {
					}
					break
					println(sum)
					`, []string{"_panic"}},
	"range_with_break_label_panic": {`
					arr := [1,3,5,7]
					L:
					sum := 0
					for i,_:=range arr {
						break L
					}
					println(sum)
					`, []string{"_panic"}},
	"range_with_continue_wrong_label_panic": {`
					arr := [1,3,5,7]
					L:
					sum := 0
					for i,_:=range arr {
						continue L
					}
					println(sum)
					`, []string{"_panic"}},
	"range_with_many_labels": {`
					arr := [1,3,5,7]
					sum := 0
					L:
					for i,_:=range arr {
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
					`, []string{"8"}},
	"range_with_many_labels_break": {`
					arr := [1,3,5,7]
					sum := 0
					L:
					for i,_:=range arr {
						if arr[i]<7{
								continue L
						}
						L1:
						for j,_:=range arr {
							if arr[j]>3{
								break L
							}
							sum+=arr[i]+arr[j]
						}
					}
					println(sum)
					`, []string{"18"}},
}

func TestRangeStmtWithBranch(t *testing.T) {
	testScripts(t, "TestRangeStmtWithBranch", testRangeStmtWithBranchClauses)
}

// -----------------------------------------------------------------------------

var testForPhraseWithBranchClauses = map[string]testData{
	"for_phrase_with_continue": {`
					arr := [1,3,5,7]
					sum := 0
					for i, _ <- arr {
						if arr[i]<5{
							continue
						}
						sum+=arr[i]
					}
					println(sum)
					`, []string{"12"}},
	"for_phrase_with_break": {`
					arr := [1,3,5,7]
					sum := 0
					for i, _ <- arr {
						if arr[i]>5{
							break
						}
						sum+=arr[i]
					}
					println(sum)
					`, []string{"9"}},
	"for_phrase_with_break_label": {`
					arr := [1,3,5,7]
					sum := 0
					L:
					for i, _ <- arr {
						for j:=0;j<len(arr);j++{
							if j>2{
								break L
							}
							sum+=arr[i]+arr[j]
						}
					}
					println(sum)
					`, []string{"12"}}, // (1+1)+(1+3)+(1+5)
	"for_phrase_with_continue_label": {`
					arr := [1,3,5,7]
					sum := 0
					L:
					for i, _ <- arr {
						for j:=0;j<len(arr);j++{
							if j>1{
								continue L
							}
							sum+=arr[i]+arr[j]
						}
					}
					println(sum)
					`, []string{"48"}}, // (1+3+5+7)*2+(1+3)*4
	"for_phrase_with_continue_break": {`
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
					`, []string{"36"}},
	"for_phrase_with_continue_break_continue": {`
					arr := [1,3,5,7]
					sum := 0
					L1:
					for i, _ <- arr {
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
					`, []string{"12"}},
	"for_phrase_with_continue_panic": {`
					arr := [1,3,5,7]
					sum := 0
					for i, _ <- arr {
					}
					continue
					println(sum)
					`, []string{"_panic"}},
	"for_phrase_with_continue_no_label_panic": {`
					arr := [1,3,5,7]
					sum := 0
					for i, _ <- arr {
						continue L
					}
					println(sum)
					`, []string{"_panic"}},
	"for_phrase_with_break_panic": {`
					arr := [1,3,5,7]
					sum := 0
					for i, _ <- arr {
					}
					break
					println(sum)
					`, []string{"_panic"}},
	"for_phrase_with_break_label_panic": {`
					arr := [1,3,5,7]
					L:
					sum := 0
					for i, _ <- arr {
						break L
					}
					println(sum)
					`, []string{"_panic"}},
	"for_phrase_with_continue_wrong_label_panic": {`
					arr := [1,3,5,7]
					L:
					sum := 0
					for i, _ <- arr {
						continue L
					}
					println(sum)
					`, []string{"_panic"}},
	"for_phrase_with_many_labels": {`
					arr := [1,3,5,7]
					sum := 0
					L:
					for i, _ <- arr {
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
					`, []string{"8"}},
	"for_phrase_with_many_labels_break": {`
					arr := [1,3,5,7]
					sum := 0
					L:
					for i, _ <- arr {
						if arr[i]<7{
								continue L
						}
						L1:
						for j, _ <- arr {
							if arr[j]>3{
								break L
							}
							sum+=arr[i]+arr[j]
						}
					}
					println(sum)
					`, []string{"18"}},
}

func TestForPhraseWithBranch(t *testing.T) {
	testScripts(t, "TestForPhraseWithBranch", testForPhraseWithBranchClauses)
}
