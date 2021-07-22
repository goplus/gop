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
	"testing"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/parser/parsertest"
)

func sourceErrorTest(t *testing.T, msg, src string) {
	fs := parsertest.NewSingleFileFS("/foo", "bar.gop", src)
	pkgs, err := parser.ParseFSDir(gblFset, fs, "/foo", nil, 0)
	if err != nil {
		t.Fatal("ParseFSDir:", err)
	}
	conf := *baseConf.Ensure()
	conf.NoFileLine = false
	conf.TargetDir = "/foo"
	bar := pkgs["main"]
	_, err = cl.NewPackage("", bar, &conf)
	if err == nil {
		t.Fatal("no error?")
	}
	if ret := err.Error(); ret != msg {
		t.Fatalf("\nError: \"%s\"\nExpected: \"%s\"\n", ret, msg)
	}
}

func _TestErrNewVar(t *testing.T) {
	sourceErrorTest(t,
		``, `
a := 1
a := "Hi"
`)
}

func TestErrSliceLit(t *testing.T) {
	sourceErrorTest(t,
		`./bar.gop:1 cannot use "Hi"+"!" (type untyped string) as type int in slice literal`,
		`a := []int{"Hi"+"!"}`)
}

func TestErrMapLit(t *testing.T) {
	sourceErrorTest(t,
		`./bar.gop:2 cannot use 1+2 (type untyped int) as type string in map key
./bar.gop:3 cannot use "Go" + "+" (type untyped string) as type int in map value`,
		`
a := map[string]int{1+2: 2}
b := map[string]int{"Hi": "Go" + "+"}
`)
}

func TestErrLabel(t *testing.T) {
	sourceErrorTest(t,
		`./bar.gop:4 label foo already defined at ./bar.gop:2
./bar.gop:2 label foo defined and not used`,
		`
foo:
	i := 1
foo:
	i++
`)
	sourceErrorTest(t,
		`./bar.gop:1 label foo is not defined`,
		`goto foo`)
}
