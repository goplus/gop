/*
 * Copyright (c) 2024 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package doc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestCheckTypeMethod(t *testing.T) {
	if ret := checkTypeMethod("_Foo_a"); ret.typ != "" || ret.name != "Foo_a" {
		t.Fatal("checkTypeMethod:", ret)
	}
	if ret := checkTypeMethod("Foo_a"); ret.typ != "Foo" || ret.name != "a" {
		t.Fatal("checkTypeMethod:", ret)
	}
}

func TestIsGopPackage(t *testing.T) {
	if isGopPackage(&doc.Package{}) {
		t.Fatal("isGopPackage: true?")
	}
}

func TestDocRecv(t *testing.T) {
	if _, ok := docRecv(&ast.Field{}); ok {
		t.Fatal("docRecv: ok?")
	}
}

func TestGopoMethod(t *testing.T) {
	testPkg(t, `
package foo

const GopPackage = true

const Gopo__T__Gop_Add = ".AddInt,AddString"

type T int

// AddInt doc
func (p T) AddInt(b int) *T {}

// AddString doc
func AddString(this *T, b string) *T {}
`, `== Type T ==
- Func AddString -
Doc: AddString doc

func AddString(this *T, b string) *T
- Method AddInt -
Recv: T
Doc: AddInt doc

func (p T) AddInt(b int) *T
- Method Gop_Add -
Recv: T
Doc: AddInt doc

func (p T) Gop_Add(b int) *T
- Method Gop_Add -
Recv: *T
Doc: AddString doc

func (this *T) Gop_Add(b string) *T
`)
}

func TestGopoFn(t *testing.T) {
	testPkg(t, `
package foo

const GopPackage = true

const Gopo_Add = "AddInt,,AddString"

// AddInt doc
func AddInt(a, b int) int {}

// AddString doc
func AddString(a, b string) string {}
`, `== Func Add ==
Doc: AddInt doc

func Add(a, b int) int
== Func Add ==
Doc: AddString doc

func Add(a, b string) string
== Func AddInt ==
Doc: AddInt doc

func AddInt(a, b int) int
== Func AddString ==
Doc: AddString doc

func AddString(a, b string) string
`)
}

func TestOverloadFn(t *testing.T) {
	testPkg(t, `
package foo

const GopPackage = true

// Bar doc
func Bar__0() {
}
`, `== Func Bar ==
Doc: Bar doc

func Bar()
`)
}

func TestOverloadMethod(t *testing.T) {
	testPkg(t, `
package foo

const GopPackage = true

type T int

// Bar doc 1
func (p *T) Bar__0() {
}

// Bar doc 2
func (t T) Bar__1() {
}
`, `== Type T ==
- Method Bar -
Recv: *T
Doc: Bar doc 1

func (p *T) Bar()
- Method Bar -
Recv: T
Doc: Bar doc 2

func (t T) Bar()
`)
}

func printVal(parts []string, format string, val any) []string {
	return append(parts, fmt.Sprintf(format, val))
}

func printFuncDecl(parts []string, fset *token.FileSet, decl *ast.FuncDecl) []string {
	var b bytes.Buffer
	if e := format.Node(&b, fset, decl); e != nil {
		panic(e)
	}
	return append(parts, b.String())
}

func printFunc(parts []string, fset *token.FileSet, format string, fn *doc.Func) []string {
	parts = printVal(parts, format, fn.Name)
	if fn.Recv != "" {
		parts = printVal(parts, "Recv: %s", fn.Recv)
	}
	parts = printVal(parts, "Doc: %s", fn.Doc)
	parts = printFuncDecl(parts, fset, fn.Decl)
	return parts
}

func printFuncs(parts []string, fset *token.FileSet, fns []*doc.Func) []string {
	for _, fn := range fns {
		parts = printFunc(parts, fset, "== Func %s ==", fn)
	}
	return parts
}

func printType(parts []string, fset *token.FileSet, typ *doc.Type) []string {
	parts = append(parts, fmt.Sprintf("== Type %s ==", typ.Name))
	for _, fn := range typ.Funcs {
		parts = printFunc(parts, fset, "- Func %s -", fn)
	}
	for _, fn := range typ.Methods {
		parts = printFunc(parts, fset, "- Method %s -", fn)
	}
	return parts
}

func printTypes(parts []string, fset *token.FileSet, types []*doc.Type) []string {
	for _, typ := range types {
		parts = printType(parts, fset, typ)
	}
	return parts
}

func printPkg(fset *token.FileSet, in *doc.Package) string {
	var parts []string
	parts = printFuncs(parts, fset, in.Funcs)
	parts = printTypes(parts, fset, in.Types)
	return strings.Join(append(parts, ""), "\n")
}

func testPkg(t *testing.T, in, expected string) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "foo.go", in, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	pkg, err := doc.NewFromFiles(fset, []*ast.File{f}, "foo")
	if err != nil {
		t.Fatal(err)
	}
	pkg = Transform(pkg)
	if ret := printPkg(fset, pkg); ret != expected {
		t.Fatalf("got:\n%s\nexpected:\n%s\n", ret, expected)
	}
}
