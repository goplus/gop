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

package printer_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/format"
	"github.com/goplus/gop/printer"
	"github.com/goplus/gop/token"
)

func init() {
	printer.SetDebug(printer.DbgFlagAll)
}

func TestNoPkgDecl(t *testing.T) {
	var dst bytes.Buffer
	fset := token.NewFileSet()
	if err := format.Node(&dst, fset, &ast.File{
		Name:      &ast.Ident{Name: "main"},
		NoPkgDecl: true,
	}); err != nil {
		t.Fatal("format.Node failed:", err)
	}
	if dst.String() != "\n" {
		t.Fatal("TestNoPkgDecl:", dst.String())
	}
}

func TestFuncs(t *testing.T) {
	var dst bytes.Buffer
	fset := token.NewFileSet()
	if err := format.Node(&dst, fset, &ast.File{
		Name: &ast.Ident{Name: "main"},
		Decls: []ast.Decl{
			&ast.FuncDecl{
				Type: &ast.FuncType{Params: &ast.FieldList{}},
				Name: &ast.Ident{Name: "foo"},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{&printer.NewlineStmt{}},
				},
			},
			&ast.FuncDecl{
				Type: &ast.FuncType{Params: &ast.FieldList{}},
				Name: &ast.Ident{Name: "bar"},
				Body: &ast.BlockStmt{},
			},
		},
		NoPkgDecl: true,
	}); err != nil {
		t.Fatal("format.Node failed:", err)
	}
	if dst.String() != `func foo() {

}

func bar() {
}
` {
		t.Fatal("TestNoPkgDecl:", dst.String())
	}
}

func diffBytes(t *testing.T, dst, src []byte) {
	line := 1
	offs := 0 // line offset
	for i := 0; i < len(dst) && i < len(src); i++ {
		d := dst[i]
		s := src[i]
		if d != s {
			t.Errorf("dst:%d: %s\n", line, dst[offs:])
			t.Errorf("src:%d: %s\n", line, src[offs:])
			return
		}
		if s == '\n' {
			line++
			offs = i + 1
		}
	}
	if len(dst) != len(src) {
		t.Errorf("len(dst) = %d, len(src) = %d\nsrc = %q", len(dst), len(src), src)
	}
}

func testFrom(t *testing.T, fpath string) {
	src, err := ioutil.ReadFile(fpath)
	if err != nil {
		t.Fatal(err)
	}

	res, err := format.Source(src)
	if err != nil {
		t.Fatal("Source failed:", err)
	}

	diffBytes(t, res, src)
}

func TestFromTutorial(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = filepath.Join(dir, "../tutorial")
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".gop" {
			t.Log(path)
			testFrom(t, path)
		}
		return nil
	})
}

func TestFromParse(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = filepath.Join(dir, "../parser/_testdata")
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := info.Name()
		if name == "cmd.gop" { // skip this file
			return nil
		}
		if !info.IsDir() && filepath.Ext(name) == ".gop" {
			t.Log(path)
			testFrom(t, path)
		}
		return nil
	})
}
