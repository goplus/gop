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

package cl

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/goplus/gop/ast/asttest"
	qspec "github.com/goplus/gop/exec.spec"
	exec "github.com/goplus/gop/exec/bytecode"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
)

type testConstInfo struct {
	Name  string
	Kind  reflect.Kind
	Value interface{}
}

func TestPkgConst(t *testing.T) {
	var I = exec.NewGoPackage("pkg_test_const")
	infos := []testConstInfo{
		{"True", reflect.Bool, true},
		{"False", reflect.Bool, false},
		{"A", qspec.ConstBoundRune, 'A'},
		{"String1", qspec.ConstBoundString, "Info"},
		{"String2", qspec.ConstBoundString, "信息"},
		{"Int1", qspec.ConstUnboundInt, -1024},
		{"Int2", qspec.ConstUnboundInt, 1024},
		{"Int3", qspec.ConstUnboundInt, -10000000},
		{"Int4", qspec.ConstUnboundInt, 10000000},
		// {"MinInt64", qspec.ConstUnboundInt, math.MinInt64},
		// {"MaxInt64", qspec.ConstUnboundInt, math.MaxInt64},
		{"Pi", qspec.ConstUnboundFloat, math.Pi},
		{"Complex", qspec.ConstUnboundComplex, 1 + 2i},
	}

	var consts []exec.GoConstInfo
	for _, info := range infos {
		consts = append(consts, I.Const(info.Name, info.Kind, info.Value))
	}
	I.RegisterConsts(consts...)

	var testSource string
	testSource = `package main

import (
	pkg "pkg_test_const"
)

`

	// make println
	for _, info := range infos {
		testSource += fmt.Sprintf("println(pkg.%v)\n", info.Name)
	}
	// make ret
	var retList []string
	var name string
	for i, info := range infos {
		if !isConstBound(info.Kind) {
			name = fmt.Sprintf("ret%v", i)
			testSource += fmt.Sprintf("%v := pkg.%v\n", name, info.Name)
		} else {
			name = "pkg." + info.Name
		}
		retList = append(retList, name)
	}
	// ret list
	testSource += strings.Join(retList, "\n")
	testSource += ""

	fsTestPkgConst := asttest.NewSingleFileFS("/foo", "bar.gop", testSource)
	t.Log(testSource)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestPkgConst, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())

	n := len(infos)
	for i, info := range infos {
		if v := ctx.Get(i - n); v != info.Value {
			t.Fatal(info.Name, v, info.Value)
		}
	}
}

type testLoadGoVarInfo struct {
	Name string
	Addr interface{}
}

func TestPkgLoadGoVar(t *testing.T) {
	var I = exec.NewGoPackage("pkg_test_var_load")
	v1 := true
	v2 := false
	v3 := rune('A')
	v4 := "Info"
	v5 := "信息"
	v6 := -100
	v7 := uint32(100)
	v8 := []int{100, 200}
	v9 := []string{"01", "02"}
	v10 := make(map[int]string)
	v10[1] = "01"
	v10[2] = "02"
	infos := []testLoadGoVarInfo{
		{"True", &v1},
		{"False", &v2},
		{"A", &v3},
		{"String1", &v4},
		{"String2", &v5},
		{"Int1", &v6},
		{"Int2", &v7},
		{"Ar1", &v8},
		{"Ar2", &v9},
		{"M1", &v10},
	}

	var vars []exec.GoVarInfo
	for _, info := range infos {
		vars = append(vars, I.Var(info.Name, info.Addr))
	}
	I.RegisterVars(vars...)

	var testSource string
	testSource = `package main

import (
	pkg "pkg_test_var_load"
)

`

	// make println
	for _, info := range infos {
		testSource += fmt.Sprintf("println(pkg.%v)\n", info.Name)
	}
	// make ret
	var retList []string
	// var retList2 []string
	var name string
	for i, info := range infos {
		name = fmt.Sprintf("id%v := pkg.%v", i, info.Name)
		retList = append(retList, name)
	}

	testSource += strings.Join(retList, "\n")

	fsTestPkgVar := asttest.NewSingleFileFS("/foo", "bar.gop", testSource)
	t.Log(testSource)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestPkgVar, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())

	n := len(infos)
	for i, info := range infos {
		if v := ctx.Get(i - n); reflect.DeepEqual(v, reflect.ValueOf(info.Addr).Elem().Interface()) {
			t.Fatal(info.Name, v, reflect.ValueOf(info.Addr).Elem().Interface())
		}
	}
}

type testStoreGoVarInfo struct {
	Name  string
	Addr  interface{}
	Store interface{}
	Gop   string
}

func TestPkgStoreGoVar(t *testing.T) {
	var I = exec.NewGoPackage("pkg_test_store_var")
	v1 := true
	v2 := false
	v3 := 'A'
	v4 := "Info"
	v5 := "信息"
	v6 := -100
	v7 := uint32(100)
	v8 := []int{100, 200}
	v9 := []string{"01", "02"}
	v10 := make(map[int]string)
	v10[1] = "01"
	v10[2] = "02"
	v10_1 := make(map[int]string)
	v10_1[1] = "02"
	v10_1[3] = "03"
	infos := []testStoreGoVarInfo{
		{"True", &v1, false, "false"},
		{"False", &v2, true, "true"},
		{"A", &v3, 'B', "'B'"},
		{"String1", &v4, "Info2", `"Info2"`},
		{"String2", &v5, "Inf3", `"Inf3"`},
		{"Int1", &v6, 100, "100"},
		{"Int2", &v7, uint32(200), "200"},
		{"Ar1", &v8, []int{200, 300}, "[200,300]"},
		{"Ar2", &v9, []string{"03", "04"}, `["03","04"]`},
		{"M1", &v10, v10_1, `{1:"02",3:"03"}`},
	}

	var vars []exec.GoVarInfo
	for _, info := range infos {
		vars = append(vars, I.Var(info.Name, info.Addr))
	}
	I.RegisterVars(vars...)

	var testSource string
	testSource = `package main

import (
	pkg "pkg_test_store_var"
)

`

	for _, info := range infos {
		testSource += fmt.Sprintf("pkg.%v = %v\n", info.Name, info.Gop)
	}

	// make println
	for _, info := range infos {
		testSource += fmt.Sprintf("println(pkg.%v)\n", info.Name)
	}

	fsTestPkgVar := asttest.NewSingleFileFS("/foo", "bar.gop", testSource)
	t.Log(testSource)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestPkgVar, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())

	for _, info := range infos {
		v := reflect.ValueOf(info.Addr).Elem().Interface()
		if !reflect.DeepEqual(v, info.Store) {
			t.Fatalf("%v, %v(%T), %v(%T)\n", info.Name, v, v, info.Store, info.Store)
		}
	}
}

func TestPkgGoVarMap(t *testing.T) {
	var I = exec.NewGoPackage("pkg_test_var_map")
	m1 := make(map[int]string)
	m1[0] = "hello"
	m1[1] = "world"
	m1_1 := make(map[int]string)
	m1_1[0] = "001"
	m1_1[1] = "002"

	m2 := make(map[string]int)
	m2["0"] = 100
	m2["1"] = 200
	m2_1 := make(map[string]int)
	m2_1["0"] = -200
	m2_1["1"] = -100
	infos := []testStoreGoVarInfo{
		{"M1", &m1, m1_1, `pkg.M1[0],pkg.M1[1]="001","002"`},
		{"M2", &m2, m2_1, `pkg.M2["0"],pkg.M2["1"]=-200,-100`},
	}

	var vars []exec.GoVarInfo
	for _, info := range infos {
		vars = append(vars, I.Var(info.Name, info.Addr))
	}
	I.RegisterVars(vars...)

	var testSource string
	testSource = `package main

import (
	pkg "pkg_test_var_map"
)

`

	for _, info := range infos {
		testSource += fmt.Sprintf("%v\n", info.Gop)
	}

	// make println
	for _, info := range infos {
		testSource += fmt.Sprintf("println(pkg.%v)\n", info.Name)
	}

	fsTestPkgVar := asttest.NewSingleFileFS("/foo", "bar.gop", testSource)
	t.Log(testSource)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestPkgVar, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())

	for _, info := range infos {
		v := reflect.ValueOf(info.Addr).Elem().Interface()
		if !reflect.DeepEqual(v, info.Store) {
			t.Fatalf("%v, %v(%T), %v(%T)\n", info.Name, v, v, info.Store, info.Store)
		}
	}
}

func TestPkgGoVarSlice(t *testing.T) {
	var I = exec.NewGoPackage("pkg_test_var_slice")

	var s1 []string
	s1 = append(s1, "hello")
	s1 = append(s1, "world")
	var s1_1 []string
	s1_1 = append(s1_1, "001")
	s1_1 = append(s1_1, "002")

	infos := []testStoreGoVarInfo{
		{"S1", &s1, s1_1, `pkg.S1[0],pkg.S1[1]="001","002"`},
	}

	var vars []exec.GoVarInfo
	for _, info := range infos {
		vars = append(vars, I.Var(info.Name, info.Addr))
	}
	I.RegisterVars(vars...)

	var testSource string
	testSource = `package main

import (
	pkg "pkg_test_var_slice"
)

`

	for _, info := range infos {
		testSource += fmt.Sprintf("%v\n", info.Gop)
	}

	// make println
	for _, info := range infos {
		testSource += fmt.Sprintf("println(pkg.%v)\n", info.Name)
	}

	fsTestPkgVar := asttest.NewSingleFileFS("/foo", "bar.gop", testSource)
	t.Log(testSource)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestPkgVar, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())

	for _, info := range infos {
		v := reflect.ValueOf(info.Addr).Elem().Interface()
		if !reflect.DeepEqual(v, info.Store) {
			t.Fatalf("%v, %v(%T), %v(%T)\n", info.Name, v, v, info.Store, info.Store)
		}
	}
}

func TestPkgGoVarArray(t *testing.T) {
	var I = exec.NewGoPackage("pkg_test_var_array")
	var ar1 [2]string
	ar1[0] = "hello"
	ar1[1] = "world"
	var ar1_1 [2]string
	ar1_1[0] = "001"
	ar1_1[1] = "002"

	infos := []testStoreGoVarInfo{
		{"A1", &ar1, ar1_1, `pkg.A1[0],pkg.A1[1]="001","002"`},
	}

	var vars []exec.GoVarInfo
	for _, info := range infos {
		vars = append(vars, I.Var(info.Name, info.Addr))
	}
	I.RegisterVars(vars...)

	var testSource string
	testSource = `package main

import (
	pkg "pkg_test_var_array"
)

	A0 := pkg.A1
	A0[0], A0[1] = "003", "004"
	println(A0)

`

	for _, info := range infos {
		testSource += fmt.Sprintf("%v\n", info.Gop)
	}

	// make println
	for _, info := range infos {
		testSource += fmt.Sprintf("println(pkg.%v)\n", info.Name)
	}

	fsTestPkgVar := asttest.NewSingleFileFS("/foo", "bar.gop", testSource)
	t.Log(testSource)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestPkgVar, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())

	for _, info := range infos {
		v := reflect.ValueOf(info.Addr).Elem().Interface()
		if !reflect.DeepEqual(v, info.Store) {
			t.Fatalf("%v, %v(%T), %v(%T)\n", info.Name, v, v, info.Store, info.Store)
		}
	}
}

type Tpoint struct {
	X int
	Y int
}

type Trect struct {
	Min Tpoint
	Max Tpoint
}

type Tfieldinfo struct {
	V1  bool
	V2  rune
	V3  string
	V4  int
	V5  []int
	V6  []string
	V7  map[int]string
	V8  Trect
	V9  *Trect
	V10 []Trect
	V11 []*Trect
	V12 [2]Trect
	V13 [2]*Trect
}

func tNewRect(x1 int, y1 int, x2 int, y2 int) *Trect {
	return &Trect{Tpoint{x1, y1}, Tpoint{x2, y2}}
}

func tMakeRect(x1 int, y1 int, x2 int, y2 int) Trect {
	return Trect{Tpoint{x1, y1}, Tpoint{x2, y2}}
}

func tNewPoint(x int, y int) *Tpoint {
	return &Tpoint{x, y}
}

func tMakePoint(x int, y int) Tpoint {
	return Tpoint{x, y}
}

func execTestNewPoint(_ int, p *exec.Context) {
	args := p.GetArgs(2)
	ret0 := tNewPoint(args[0].(int), args[1].(int))
	p.Ret(2, ret0)
}

func execTestMakePoint(_ int, p *exec.Context) {
	args := p.GetArgs(2)
	ret0 := tMakePoint(args[0].(int), args[1].(int))
	p.Ret(2, ret0)
}

func execTestNewRect(_ int, p *exec.Context) {
	args := p.GetArgs(4)
	ret0 := tNewRect(args[0].(int), args[1].(int), args[2].(int), args[3].(int))
	p.Ret(4, ret0)
}

func execTestMakeRect(_ int, p *exec.Context) {
	args := p.GetArgs(4)
	ret0 := tMakeRect(args[0].(int), args[1].(int), args[2].(int), args[3].(int))
	p.Ret(4, ret0)
}

func TestPkgField(t *testing.T) {
	var I = exec.NewGoPackage("pkg_test_field")

	m := make(map[int]string)
	m[1] = "hello"
	m[2] = "world"

	var ar [13]interface{}
	ar[0] = true
	ar[1] = 'A'
	ar[2] = "Info"
	ar[3] = -100
	ar[4] = []int{100, 200}
	ar[5] = []string{"hello", "world"}
	ar[6] = m
	ar[7] = tMakeRect(10, 20, 100, 200)
	ar[8] = tNewRect(10, 20, 100, 200)
	ar[9] = []Trect{tMakeRect(10, 20, 30, 40), tMakeRect(50, 60, 70, 80)}
	ar[10] = []*Trect{tNewRect(10, 20, 30, 40), tNewRect(50, 60, 70, 80)}
	ar[11] = [2]Trect{tMakeRect(10, 20, 30, 40), tMakeRect(50, 60, 70, 80)}
	ar[12] = [2]*Trect{tNewRect(10, 20, 30, 40), tNewRect(50, 60, 70, 80)}

	info := &Tfieldinfo{}
	info.V7 = make(map[int]string)
	v := reflect.ValueOf(info).Elem()
	for i := 0; i < v.NumField(); i++ {
		v.Field(i).Set(reflect.ValueOf(ar[i]))
	}

	var out [13]interface{}
	var sum int
	I.RegisterVars(
		I.Var("Info", &info),
		I.Var("Out", &out),
		I.Var("Sum", &sum),
	)
	I.RegisterFuncs(
		I.Func("NewRect", tNewRect, execTestNewRect),
		I.Func("NewPoint", tNewPoint, execTestNewPoint),
		I.Func("MakeRect", tMakeRect, execTestMakeRect),
		I.Func("MakePoint", tMakePoint, execTestMakePoint),
	)

	var testSource string
	testSource = `package main

import (
	pkg "pkg_test_field"
)

`
	const nField = len(ar)

	// make set
	for i := 0; i < nField; i++ {
		testSource += fmt.Sprintf("pkg.Info.V%v = pkg.Info.V%v\n", i+1, i+1)
	}

	// make println
	for i := 0; i < nField; i++ {
		testSource += fmt.Sprintf("println(\"pkg.Info.V%v = \", pkg.Info.V%v)\n", i+1, i+1)
	}

	for i := 0; i < nField; i++ {
		testSource += fmt.Sprintf("pkg.Out[%v] = pkg.Info.V%v\n", i, i+1)
	}

	testSource += `
	pkg.Sum = 1+pkg.Info.V4+pkg.NewPoint(1,2).X+pkg.NewRect(10,20,30,40).Max.Y+pkg.Info.V8.Min.X-pkg.Info.V9.Max.Y
	pkg.Info.V5[0] = -100
	println("pkg.Info.V5", pkg.Info.V5)
	println("pkg.Info.V11[0]",pkg.Info.V11[0])
	pkg.Info.V7[0] = "hello"
	println("pkg.Info.V7[0]",pkg.Info.V7[0])
	pkg.Info.V11[0] = nil
	println("pkg.Info.V11",pkg.Info.V11)
	pkg.Info.V11[1].Min.X = -101
	println("pkg.Info.V11[1]",pkg.Info.V11[1])
	println("pkg.Info.V11[1].Min",pkg.Info.V11[1].Min)
	println("pkg.Info.V11[1].Min.X",pkg.Info.V11[1].Min.X)
	pkg.Info.V13[0] = nil
	println("pkg.Info.V13",pkg.Info.V13)
	pkg.Info.V13[1].Min = pkg.MakePoint(-105,-106)
	println("pkg.Info.V13[1].Min",pkg.Info.V13[1].Min)
	pkg.Info.V13[1].Max.Y = -102
	println("pkg.Info.V13[1].Max.Y",pkg.Info.V13[1].Max.Y)
	`

	fsTestPkgVar := asttest.NewSingleFileFS("/foo", "bar.gop", testSource)
	t.Log(testSource)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestPkgVar, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()
	code.Dump(os.Stdout)

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())

	for i := 0; i < nField; i++ {
		if !reflect.DeepEqual(ar[i], out[i]) {
			t.Fatal(i, ar[i], out[i])
		}
	}
	if sum != -248 {
		t.Fatal("binary expr check fail", sum)
	}
	if info.V5[0] != -100 {
		t.Fatal("V5", info.V5)
	}
	if info.V11[0] != nil || info.V11[1].Min.X != -101 {
		t.Fatal("V11", info.V11)
	}
	if info.V13[0] != nil || info.V13[1].Min.X != -105 || info.V13[1].Max.Y != -102 {
		t.Fatal("V13", info.V13[0], info.V13[1])
	}
}

func TestPkgType(t *testing.T) {
	var I = exec.NewGoPackage("pkg_test_type")
	I.RegisterTypes(
		I.Type("Rect", reflect.TypeOf((*Trect)(nil)).Elem()),
		I.Type("Point", reflect.TypeOf((*Tpoint)(nil)).Elem()),
	)

	var testSource = `
	import pkg "pkg_test_type"

	pkg.Rect{pkg.Point{1,2},pkg.Point{3,4}}
	pkg.Rect{pkg.Point{X:5,Y:6},pkg.Point{X:7,Y:8}}
	&pkg.Rect{pkg.Point{1,2},pkg.Point{3,4}}
	&pkg.Rect{pkg.Point{X:5,Y:6},pkg.Point{X:7,Y:8}}
	pkg.Point{1,2}
	pkg.Point{X:3,Y:4}
	&pkg.Point{5,6}
	&pkg.Point{X:7,Y:8}
	`

	fsTestPkgVar := asttest.NewSingleFileFS("/foo", "bar.gop", testSource)
	t.Log(testSource)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestPkgVar, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()
	code.Dump(os.Stdout)

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Get(-8); v != nil && !reflect.DeepEqual(v, tMakeRect(1, 2, 3, 4)) {
		t.Fatal("check rect", v)
	}
	if v := ctx.Get(-7); v != nil && !reflect.DeepEqual(v, tMakeRect(5, 6, 7, 8)) {
		t.Fatal("check rect", v)
	}
	if v := ctx.Get(-6); v != nil && !reflect.DeepEqual(v, tNewRect(1, 2, 3, 4)) {
		t.Fatal("check rect", v)
	}
	if v := ctx.Get(-5); v != nil && !reflect.DeepEqual(v, tNewRect(5, 6, 7, 8)) {
		t.Fatal("check rect", v)
	}
	if v := ctx.Get(-4); v != nil && !reflect.DeepEqual(v, Tpoint{1, 2}) {
		t.Fatal("check point", v)
	}
	if v := ctx.Get(-3); v != nil && !reflect.DeepEqual(v, Tpoint{X: 3, Y: 4}) {
		t.Fatal("check point", v)
	}
	if v := ctx.Get(-2); v != nil && !reflect.DeepEqual(v, &Tpoint{5, 6}) {
		t.Fatal("check point", v)
	}
	if v := ctx.Get(-1); v != nil && !reflect.DeepEqual(v, &Tpoint{X: 7, Y: 8}) {
		t.Fatal("check point", v)
	}

}

func TestPkgTakeAddr(t *testing.T) {
	var I = exec.NewGoPackage("pkg_test_takeaddr")
	rc := tMakeRect(10, 20, 100, 200)
	ar := [2]Trect{tMakeRect(1, 2, 10, 20), tMakeRect(10, 20, 100, 200)}
	slice := []Trect{tMakeRect(1, 2, 10, 20), tMakeRect(10, 20, 100, 200)}
	m := make(map[int]Trect)
	m[1] = tMakeRect(1, 2, 10, 20)
	m[2] = tMakeRect(10, 20, 100, 200)
	I.RegisterVars(
		I.Var("RC", &rc),
		I.Var("Ar", &ar),
		I.Var("Slice", &slice),
		I.Var("M", &m),
	)
	var testSource = `
	import pkg "pkg_test_takeaddr"

	&pkg.M
	&pkg.Ar
	&pkg.Ar[0]
	&pkg.Ar[0].Min
	&pkg.Ar[1].Max.Y
	&pkg.Slice
	&pkg.Slice[0]
	&pkg.Slice[0].Min
	&pkg.Slice[1].Max.Y
	&pkg.RC
	&pkg.RC.Min
	&pkg.RC.Max.Y
	`

	fsTestPkgVar := asttest.NewSingleFileFS("/foo", "bar.gop", testSource)
	t.Log(testSource)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestPkgVar, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()
	code.Dump(os.Stdout)

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Get(-12); v != &m {
		t.Fatal("takeAddr &m", v)
	}
	if v := ctx.Get(-11); v != &ar {
		t.Fatal("takeAddr &ar", v)
	}
	if v := ctx.Get(-10); v != &ar[0] {
		t.Fatal("takeAddr &ar[0]", v)
	}
	if v := ctx.Get(-9); v != &ar[0].Min {
		t.Fatal("takeAddr &ar[0].Min", v)
	}
	if v := ctx.Get(-8); v != &ar[1].Max.Y {
		t.Fatal("takeAddr &ar[1].Max.Y", v)
	}
	if v := ctx.Get(-7); v != &slice {
		t.Fatal("takeAddr &slice", v)
	}
	if v := ctx.Get(-6); v != &slice[0] {
		t.Fatal("takeAddr &slice[0]", v)
	}
	if v := ctx.Get(-5); v != &slice[0].Min {
		t.Fatal("takeAddr &slice[0].Min", v)
	}
	if v := ctx.Get(-4); v != &slice[1].Max.Y {
		t.Fatal("takeAddr &slice[1].Max.Y", v)
	}
	if v := ctx.Get(-3); v != &rc {
		t.Fatal("takeAddr &rc", v)
	}
	if v := ctx.Get(-2); v != &rc.Min {
		t.Fatal("takeAddr &rc.Min", v)
	}
	if v := ctx.Get(-1); v != &rc.Max.Y {
		t.Fatal("takeAddr &rc.Max.Y", v)
	}
}

func TestPkgFieldBadSet(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("must panic")
		} else {
			t.Log("panic info", r)
		}
	}()

	var I = exec.NewGoPackage("pkg_test_field_bad")

	I.RegisterFuncs(
		I.Func("NewRect", tNewRect, execTestNewRect),
		I.Func("NewPoint", tNewPoint, execTestNewPoint),
		I.Func("MakeRect", tMakeRect, execTestMakeRect),
		I.Func("MakePoint", tMakePoint, execTestMakePoint),
	)

	var testSource string
	testSource = `package main

import (
	pkg "pkg_test_field_bad"
)

pkg.MakeRect(10,20,100,200).Min.X = 10

`
	fsTestPkgVar := asttest.NewSingleFileFS("/foo", "bar.gop", testSource)
	t.Log(testSource)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestPkgVar, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, _, err = newPackage(b, bar, fset)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()
	code.Dump(os.Stdout)

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
}
