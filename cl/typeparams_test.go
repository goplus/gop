//go:build go1.18
// +build go1.18

package cl_test

import (
	"go/scanner"
	"os"
	"runtime"
	"testing"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/parser/fsx/memfs"
)

func TestTypeParams(t *testing.T) {
	gopMixedClTest(t, "main", `package main

type Data[X, Y any] struct {
	v X
}

func (p *Data[T, R]) foo() {
	fmt.Println(p.v)
}
`, `
v := Data[int, float64]{1}
v.foo()
`, `package main

func main() {
	v := Data[int, float64]{1}
	v.foo()
}
`)
}

func TestTypeParamsFunc(t *testing.T) {
	gopMixedClTest(t, "main", `package main

type Number interface {
	~int | ~uint | (float64)
}

func Sum[T Number](vec []T) T {
	var sum T
	for _, elt := range vec {
		sum = sum + elt
	}
	return sum
}

func At[T interface{ ~[]E }, E any](x T, i int) E {
	return x[i]
}

func Loader[T1 any, T2 any](p1 T1, p2 T2) T1 {
	return p1
}

func Add[T1 any, T2 ~int|~uint](v1 T1, v2 ...T2) (sum T2) {
	println(v1)
	for _, v := range v2 {
		sum += v
	}
	return sum
}

type Int []int
var MyInts = Int{1,2,3,4}
`, `
_ = At[[]int]
_ = At[[]int,int]
_ = Sum[int]
_ = Loader[*int,int]
_ = Add[string,int]
var s1 int = Sum([1, 2, 3])
var s2 int = Sum[int]([1, 2, 3])
var v1 int = At([1, 2, 3], 1)
var v2 int = At[[]int]([1, 2, 3], 1)
var v3 int = At[[]int, int]([1, 2, 3], 1)
var n1 int = Add("hello", 1, 2, 3)
var n2 int = Add[string]("hello", 1, 2, 3)
var n3 int = Add[string, int]("hello", 1, 2, 3)
var n4 int = Add("hello", [1, 2, 3]...)
var n5 int = Add("hello", [1, 2, 3]...)
var n6 int = Add[string]("hello", MyInts...)
var n7 int = Add[string, int]("hello", [1, 2, 3]...)
var p1 *int = Loader[*int](nil, 1)
var p2 *int = Loader[*int, int](nil, 1)
var fn1 func(p1 *int, p2 int) *int
fn1 = Loader[*int, int]
fn1(nil, 1)
`, `package main

func main() {
	_ = At[[]int]
	_ = At[[]int, int]
	_ = Sum[int]
	_ = Loader[*int, int]
	_ = Add[string, int]
	var s1 int = Sum([]int{1, 2, 3})
	var s2 int = Sum[int]([]int{1, 2, 3})
	var v1 int = At([]int{1, 2, 3}, 1)
	var v2 int = At[[]int]([]int{1, 2, 3}, 1)
	var v3 int = At[[]int, int]([]int{1, 2, 3}, 1)
	var n1 int = Add("hello", 1, 2, 3)
	var n2 int = Add[string]("hello", 1, 2, 3)
	var n3 int = Add[string, int]("hello", 1, 2, 3)
	var n4 int = Add("hello", []int{1, 2, 3}...)
	var n5 int = Add("hello", []int{1, 2, 3}...)
	var n6 int = Add[string]("hello", MyInts...)
	var n7 int = Add[string, int]("hello", []int{1, 2, 3}...)
	var p1 *int = Loader[*int](nil, 1)
	var p2 *int = Loader[*int, int](nil, 1)
	var fn1 func(p1 *int, p2 int) *int
	fn1 = Loader[*int, int]
	fn1(nil, 1)
}
`)
}

func TestTypeParamsType(t *testing.T) {
	gopMixedClTest(t, "main", `package main
type Data[T any] struct {
	v T
}
func(p *Data[T]) Set(v T) {
	p.v = v
}
func(p *(Data[T1])) Set2(v T1) {
	p.v = v
}
type sliceOf[E any] interface {
	~[]E
}
type Slice[S sliceOf[T], T any] struct {
	Data S
}
func (p *Slice[S, T]) Append(t ...T) S {
	p.Data = append(p.Data, t...)
	return p.Data
}
func (p *Slice[S1, T1]) Append2(t ...T1) S1 {
	p.Data = append(p.Data, t...)
	return p.Data
}
type (
	DataInt = Data[int]
	SliceInt = Slice[[]int,int]
)
`, `
type DataString = Data[string]
type SliceString = Slice[[]string,string]
println(DataInt{1}.v)
println(DataString{"hello"}.v)
println(Data[int]{100}.v)
println(Data[string]{"hello"}.v)
println(Data[struct{X int;Y int}]{}.v.X)

v1 := SliceInt{}
v2 := SliceString{}
v3 := Slice[[]int,int]{}
v3.Append([1,2,3,4]...)
v3.Append2([1,2,3,4]...)
`, `package main

import "fmt"

type DataString = Data[string]
type SliceString = Slice[[]string, string]

func main() {
	fmt.Println(Data[int]{1}.v)
	fmt.Println(Data[string]{"hello"}.v)
	fmt.Println(Data[int]{100}.v)
	fmt.Println(Data[string]{"hello"}.v)
	fmt.Println(Data[struct {
		X int
		Y int
	}]{}.v.X)
	v1 := Slice[[]int, int]{}
	v2 := Slice[[]string, string]{}
	v3 := Slice[[]int, int]{}
	v3.Append([]int{1, 2, 3, 4}...)
	v3.Append2([]int{1, 2, 3, 4}...)
}
`)
}

func TestTypeParamsComparable(t *testing.T) {
	gopMixedClTest(t, "main", `package main
// Index returns the index of x in s, or -1 if not found.
func Index[T comparable](s []T, x T) int {
	for i, v := range s {
		// v and x are type T, which has the comparable
		// constraint, so we can use == here.
		if v == x {
			return i
		}
	}
	return -1
}

var IndexInt = Index[int]
`, `
v1 := IndexInt([1,2,3,4],1)
v2 := Index([1,2,3,4],1)
v3 := Index[int]([1,2,3,4],1)
v4 := Index(["a","b","c","d"],"b")
`, `package main

func main() {
	v1 := IndexInt([]int{1, 2, 3, 4}, 1)
	v2 := Index([]int{1, 2, 3, 4}, 1)
	v3 := Index[int]([]int{1, 2, 3, 4}, 1)
	v4 := Index([]string{"a", "b", "c", "d"}, "b")
}
`)
}

func mixedErrorTest(t *testing.T, msg, gocode, gopcode string) {
	mixedErrorTestEx(t, "main", msg, gocode, gopcode)
}

func mixedErrorTestEx(t *testing.T, pkgname, msg, gocode, gopcode string) {
	fs := memfs.TwoFiles("/foo", "a.go", gocode, "b.gop", gopcode)
	pkgs, err := parser.ParseFSDir(gblFset, fs, "/foo", parser.Config{})
	if err != nil {
		scanner.PrintError(os.Stderr, err)
		t.Fatal("parser.ParseFSDir failed")
	}
	conf := *gblConf
	conf.NoFileLine = false
	conf.RelativeBase = "/foo"
	conf.WorkingDir = "/foo"
	bar := pkgs[pkgname]
	_, err = cl.NewPackage("", bar, &conf)
	if err == nil {
		t.Fatal("no error?")
	}
	if ret := err.Error(); ret != msg {
		t.Fatalf("\nError: \"%s\"\nExpected: \"%s\"\n", ret, msg)
	}
}

func TestTypeParamsErrorInstantiate(t *testing.T) {
	var msg string
	switch runtime.Version()[:6] {
	case "go1.18":
		msg = `b.gop:2:1: uint does not implement Number`
	case "go1.19":
		msg = `b.gop:2:1: uint does not implement Number (uint missing in ~int | float64)`
	default:
		msg = `b.gop:2:1: uint does not satisfy Number (uint missing in ~int | float64)`
	}

	mixedErrorTest(t, msg, `
package main

type Number interface {
	~int | float64
}

func Sum[T Number](vec []T) T {
	var sum T
	for _, elt := range vec {
		sum = sum + elt
	}
	return sum
}

var	SumInt = Sum[int]
`, `
Sum[uint]
`)
}

func TestTypeParamsErrorMatch(t *testing.T) {
	var msg string
	switch runtime.Version()[:6] {
	case "go1.18", "go1.19":
		msg = `b.gop:2:5: T does not match ~[]E`
	case "go1.20":
		msg = `b.gop:2:5: int does not match ~[]E`
	default:
		msg = `b.gop:2:5: T (type int) does not satisfy interface{interface{~[]E}}`
	}
	mixedErrorTest(t, msg, `
package main

func At[T interface{ ~[]E }, E any](x T, i int) E {
	return x[i]
}

var	AtInt = At[[]int]
`, `
_ = At[int]
`)
}

func TestTypeParamsErrInferFunc(t *testing.T) {
	mixedErrorTest(t, `b.gop:2:5: cannot infer T2 (/foo/a.go:4:21)`, `
package main

func Loader[T1 any, T2 any](p1 T1, p2 T2) T1 {
	return p1
}
`, `
_ = Loader[int]
`)
}

func TestTypeParamsErrArgumentsParameters1(t *testing.T) {
	mixedErrorTest(t, `b.gop:2:7: got 1 type arguments but Data[T1, T2 interface{}] has 2 type parameters`, `
package main

type Data[T1 any, T2 any] struct {
	v1 T1
	v2 T2
}
`, `
var v Data[int]
`)
}

func TestTypeParamsErrArgumentsParameters2(t *testing.T) {
	mixedErrorTest(t, `b.gop:2:7: got 3 type arguments but Data[T1, T2 interface{}] has 2 type parameters`, `
package main

type Data[T1 any, T2 any] struct {
	v1 T1
	v2 T2
}
`, `
var v Data[int,int,int]
`)
}

func TestTypeParamsErrArgumentsParameters3(t *testing.T) {
	mixedErrorTest(t, `b.gop:2:1: got 3 type arguments but func[T1, T2 interface{}](t1 T1, t2 T2) has 2 type parameters`, `
package main

func Test[T1 any, T2 any](t1 T1, t2 T2) {
	println(t1,t2)
}
`, `
Test[int,int,int](1,2)
`)
}

func TestTypeParamsErrCallArguments1(t *testing.T) {
	mixedErrorTest(t, `b.gop:2:1: not enough arguments in call to Test
	have (untyped int)
	want (T1, T2)`, `
package main

func Test[T1 any, T2 any](t1 T1, t2 T2) {
	println(t1,t2)
}
`, `
Test(1)
`)
}

func TestTypeParamsErrCallArguments2(t *testing.T) {
	mixedErrorTest(t, `b.gop:2:1: too many arguments in call to Test
	have (untyped int, untyped int, untyped int)
	want (T1, T2)`, `
package main

func Test[T1 any, T2 any](t1 T1, t2 T2) {
	println(t1,t2)
}
`, `
Test(1,2,3)
`)
}

func TestTypeParamsErrCallArguments3(t *testing.T) {
	mixedErrorTest(t, `b.gop:2:1: too many arguments in call to Test
	have (untyped int, untyped int)
	want ()`, `
package main

func Test[T1 any, T2 any]() {
	var t1 T1
	var t2 T2
	println(t1,t2)
}
`, `
Test(1,2)
`)
}

func TestTypeParamsErrCallVariadicArguments1(t *testing.T) {
	mixedErrorTest(t, `b.gop:2:1: not enough arguments in call to Add
	have ()
	want (T1, ...T2)`, `
package main

func Add[T1 any, T2 ~int|~uint](v1 T1, v2 ...T2) (sum T2) {
	println(v1)
	for _, v := range v2 {
		sum += v
	}
	return sum
}
`, `
Add()
`)
}

func TestTypeParamsErrCallVariadicArguments2(t *testing.T) {
	mixedErrorTest(t, `b.gop:2:1: cannot infer T2 (a.go:4:18)`, `
package main

func Add[T1 any, T2 ~int|~uint](v1 T1, v2 ...T2) (sum T2) {
	println(v1)
	for _, v := range v2 {
		sum += v
	}
	return sum
}
`, `
Add(1)
`)
}

func TestTypeParamsRecvTypeError1(t *testing.T) {
	mixedErrorTest(t, `a.go:7:9: cannot use generic type Data[T interface{}] without instantiation`, `
package main

type Data[T any] struct {
	v T
}
func(p *Data) Test() {
}
`, `
Data[int]{}.Test()
`)
}

func TestTypeParamsRecvTypeError2(t *testing.T) {
	mixedErrorTest(t, `a.go:7:9: got 2 arguments but 1 type parameters`, `
package main

type Data[T any] struct {
	v T
}
func(p *Data[T1,T2]) Test() {
}
`, `
Data[int]{}.Test()
`)
}

func TestTypeParamsRecvTypeError3(t *testing.T) {
	mixedErrorTest(t, `a.go:8:9: got 1 type parameter, but receiver base type declares 2`, `
package main

type Data[T1 any, T2 any] struct {
	v1 T1
	v2 T2
}
func(p *Data[T1]) Test() {
}
`, `
Data[int,int]{}.Test()
`)
}

func TestGenericTypeWithoutInst1(t *testing.T) {
	mixedErrorTest(t, `a.go:8:9: cannot use generic type Data[T1, T2 interface{}] without instantiation`, `
package main

type Data[T1 any, T2 any] struct {
	v1 T1
	v2 T2
}
func(p *Data) Test() {
}
`, `
var v Data[int,int]
`)
}

func TestGenericTypeWithoutInst2(t *testing.T) {
	mixedErrorTest(t, `a.go:10:2: cannot use generic type Data[T1, T2 interface{}] without instantiation`, `
package main

type Data[T1 any, T2 any] struct {
	v1 T1
	v2 T2
}

type My[T any] struct {
	Data
}
`, `
var v My[int]
`)
}

func TestGenericTypeWithoutInst3(t *testing.T) {
	mixedErrorTest(t, `a.go:10:2: cannot use generic type Data[T1, T2 interface{}] without instantiation`, `
package main

type Data[T1 any, T2 any] struct {
	v1 T1
	v2 T2
}

type My struct {
	Data
}
`, `
var v My
`)
}

func TestGenericTypeWithoutInst4(t *testing.T) {
	mixedErrorTest(t, `a.go:10:15: cannot use generic type Data[T1, T2 interface{}] without instantiation`, `
package main

type Data[T1 any, T2 any] struct {
	v1 T1
	v2 T2
}

type My struct {
	v map[string]Data
}
`, `
var v My
`)
}

func TestGenericTypeWithoutInst5(t *testing.T) {
	mixedErrorTest(t, `b.gop:2:7: cannot use generic type Data[T1, T2 interface{}] without instantiation`, `
package main

type Data[T1 any, T2 any] struct {
	v1 T1
	v2 T2
}
`, `
var v Data
`)
}

func TestGenericTypeWithoutInst6(t *testing.T) {
	mixedErrorTest(t, `b.gop:2:8: cannot use generic type Data[T1, T2 interface{}] without instantiation`, `
package main

type Data[T1 any, T2 any] struct {
	v1 T1
	v2 T2
}
`, `
type T Data
`)
}

func TestGenericTypeWithoutInst7(t *testing.T) {
	mixedErrorTest(t, `b.gop:3:2: cannot use generic type Data[T1, T2 interface{}] without instantiation`, `
package main

type Data[T1 any, T2 any] struct {
	v1 T1
	v2 T2
}
`, `
type My struct {
	Data
}
`)
}

func TestGenericTypeWithoutInst8(t *testing.T) {
	mixedErrorTest(t, `b.gop:2:23: cannot use generic type Data[T1, T2 interface{}] without instantiation`, `
package main

type Data[T1 any, T2 any] struct {
	v1 T1
	v2 T2
}
`, `
func test(v1 int, v2 *Data) {
}
`)
}

func TestGenericTypeCompositeLit(t *testing.T) {
	gopMixedClTest(t, "main", `package main
type A[T any] struct {
	m T
}

type B[T any] struct {
	n A[T]
}

`, `
var a [2]int
if 0 == a[1] {
	println "world"
}
println B[int]{}.n
if 0 < (B[int]{}).n.m {
}
`, `package main

import "fmt"

var a [2]int

func main() {
	if 0 == a[1] {
		fmt.Println("world")
	}
	fmt.Println(B[int]{}.n)
	if 0 < (B[int]{}).n.m {
	}
}
`)
}
