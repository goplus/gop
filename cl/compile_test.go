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
	"bytes"
	"testing"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/parser/parsertest"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

func gopClTest(t *testing.T, gopcode, expected string) {
	fset := token.NewFileSet()
	fs := parsertest.NewSingleFileFS("/foo", "bar.gop", gopcode)
	pkgs, err := parser.ParseFSDir(fset, fs, "/foo", nil, 0)
	if err != nil {
		t.Fatal("ParseFSDir:", err)
	}
	bar := pkgs["main"]
	pkg, err := cl.NewPackage("", bar, fset, nil)
	if err != nil {
		t.Fatal("NewPackage:", err)
	}
	var b bytes.Buffer
	err = gox.WriteTo(&b, pkg)
	if err != nil {
		t.Fatal("gox.WriteTo failed:", err)
	}
	result := b.String()
	if result != expected {
		t.Fatalf("\nResult:\n%s\nExpected:%s\n", result, expected)
	}
}

func TestImport(t *testing.T) {
	gopClTest(t, `import "fmt"

func main() {
}`, `package main

import fmt "fmt"

func main() {
}
`)
}

func TestFunc(t *testing.T) {
	gopClTest(t, `import "fmt"

func foo(fmt string, args ...interface{}) {
}

func main() {
}`, `package main

import fmt "fmt"

func foo(fmt string, args ...interface {
}) {
}
func main() {
}
`)
}

func TestUnnamedMainFunc(t *testing.T) {
	gopClTest(t, `i := 1`, `package main

func main() {
	var i int
	i = 1
}
`)
}
