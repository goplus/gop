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
	var I = exec.NewGoPackage("pkg")
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
	"pkg"
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
