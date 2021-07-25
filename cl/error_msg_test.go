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

func codeErrorTest(t *testing.T, msg, src string) {
	fs := parsertest.NewSingleFileFS("/foo", "bar.gop", src)
	pkgs, err := parser.ParseFSDir(gblFset, fs, "/foo", nil, 0)
	if err != nil {
		t.Fatal("ParseFSDir:", err)
	}
	conf := *baseConf.Ensure()
	conf.NoFileLine = false
	conf.WorkingDir = "/foo"
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

func TestErrInitFunc(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:2:1 func init must have no arguments and no return values`, `
func init(v byte) {
}
`)
}

func TestErrRecv(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:5:9 invalid receiver type a (a is a pointer type)`, `

type a *int

func (p a) foo() {
}
`)
	codeErrorTest(t,
		`./bar.gop:2:9 invalid receiver type error (error is an interface type)`, `
func (p error) foo() {
}
`)
	codeErrorTest(t,
		`./bar.gop:2:9 invalid receiver type []byte ([]byte is not a defined type)`, `
func (p []byte) foo() {
}
`)
	codeErrorTest(t,
		`./bar.gop:2:10 invalid receiver type []byte ([]byte is not a defined type)`, `
func (p *[]byte) foo() {
}
`)
}

func _TestErrNewVar(t *testing.T) {
	codeErrorTest(t,
		``, `
a := 1
a := "Hi"
`)
}

func TestErrSliceLit(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:3:12 cannot use a (type string) as type int in slice literal`,
		`
a := "Hi"
b := []int{a}
`)
}

func TestErrMapLit(t *testing.T) {
	codeErrorTest(t, // TODO: first column need correct
		`./bar.gop:2:34 cannot use 1+2 (type untyped int) as type string in map key
./bar.gop:3:27 cannot use "Go" + "+" (type untyped string) as type int in map value`,
		`
a := map[string]int{1+2: 2}
b := map[string]int{"Hi": "Go" + "+"}
`)
}

func TestErrMember(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:3:6 a.x undefined (type string has no field or method x)`,
		`
a := "Hello"
b := a.x
`)
}

func TestErrLabel(t *testing.T) {
	codeErrorTest(t,
		`./bar.gop:4:1 label foo already defined at ./bar.gop:2:1
./bar.gop:2:1 label foo defined and not used`,
		`x := 1
foo:
	i := 1
foo:
	i++
`)
	codeErrorTest(t,
		`./bar.gop:2:6 label foo is not defined`,
		`x := 1
goto foo`)
}
