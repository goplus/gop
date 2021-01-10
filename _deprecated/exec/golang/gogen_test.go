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

package golang

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/qiniu/x/log"
)

func saveGoFile(dir string, code *Code) error {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	b, err := code.Bytes(nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dir+"/gop_autogen.go", b, 0666)
}

// -----------------------------------------------------------------------------

func getPkg(pkgs map[string]*ast.Package) *ast.Package {
	for _, pkg := range pkgs {
		return pkg
	}
	return nil
}

func testGenGo(t *testing.T, pkgDir, sel, exclude string) {
	if sel != "" && !strings.Contains(pkgDir, sel) {
		return
	}
	if exclude != "" && strings.Contains(pkgDir, exclude) {
		return
	}
	log.Debug("Compiling", pkgDir)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgDir, nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseDir failed:", err, len(pkgs))
	}

	bar := getPkg(pkgs)
	b := NewBuilder(bar.Name, nil, fset)
	_, err = cl.NewPackage(b.Interface(), bar, fset, cl.PkgActClAll)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()
	err = saveGoFile(pkgDir, code)
	if err != nil {
		t.Fatal("saveGoFile failed:", err)
	}
}

func TestGenGofile(t *testing.T) {
	sel, exclude := "", ""
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir += "/testdata"
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		testGenGo(t, dir+"/"+fi.Name(), sel, exclude)
	}
}

func TestGenGofile2(t *testing.T) {
	sel, exclude := "", ""
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir += "/testdata2"
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		testGenGo(t, dir+"/"+fi.Name(), sel, exclude)
	}
}

// -----------------------------------------------------------------------------
