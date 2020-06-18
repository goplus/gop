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
	"reflect"
	"strings"
	"testing"

	"github.com/qiniu/goplus/ast/asttest"
	qspec "github.com/qiniu/goplus/exec.spec"
	exec "github.com/qiniu/goplus/exec/bytecode"
	"github.com/qiniu/goplus/parser"
	"github.com/qiniu/goplus/token"
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
		{"MinInt64", reflect.Int64, int64(math.MinInt64)},
		{"MaxUint64", reflect.Uint64, uint64(math.MaxUint64)},
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

type testGoVarInfo struct {
	Name  string
	Addr  interface{}
	Store interface{}
}

func TestPkgGoVar(t *testing.T) {
	var I = exec.NewGoPackage("pkg_test_var")
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
	infos := []testGoVarInfo{
		{"True", &v1, false},
		{"False", &v2, true},
		{"A", &v3, 'B'},
		{"String1", &v4, "Info2"},
		{"String2", &v5, "Inf3"},
		{"Int1", &v6, 100},
		{"Int2", &v7, 200},
		{"Ar1", &v8, []int{200, 300}},
		{"Ar2", &v9, []string{"03", "04"}},
		{"M1", &v10, nil},
	}

	var vars []exec.GoVarInfo
	for _, info := range infos {
		vars = append(vars, I.Var(info.Name, info.Addr))
	}
	I.RegisterVars(vars...)

	var testSource string
	testSource = `package main

import (
	pkg "pkg_test_var"
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
	// for i, info := range infos {
	// 	s := fmt.Sprintf("pkg.%v = %v\n", info.Name, info.Store)
	// 	testSource += s
	// 	name = fmt.Sprintf("id_store%v := pkg.%v", i, info.Name)
	// 	retList2 = append(retList, name)
	// }
	// ret list
	testSource += strings.Join(retList, "\n")
	// ret list
	// testSource += strings.Join(retList2, "\n")

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
		if v := ctx.Get(i - n); reflect.DeepEqual(v, reflect.ValueOf(info.Addr).Elem().Interface()) {
			t.Fatal(info.Name, v, reflect.ValueOf(info.Addr).Elem().Interface())
		}
	}
}
