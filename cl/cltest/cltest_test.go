/*
 Copyright 2020 Goplus Team (goplus.org)

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

package cltest

import (
	"fmt"
	"io/ioutil"
	"os"
	shell "os/exec"
	"path"
	"strings"
	"testing"

	"github.com/qiniu/goplus/ast"
	"github.com/qiniu/goplus/cl"
	"github.com/qiniu/goplus/gop"
	libbuiltin "github.com/qiniu/goplus/lib/builtin"
	libfmt "github.com/qiniu/goplus/lib/fmt"
	"github.com/qiniu/goplus/parser"
	"github.com/qiniu/goplus/token"
	"github.com/qiniu/x/log"

	exec "github.com/qiniu/goplus/exec/bytecode"
	"github.com/qiniu/goplus/exec/golang"
	_ "github.com/qiniu/goplus/lib" // libraries
)

// -----------------------------------------------------------------------------

func TestFromTestdata(t *testing.T) {
	sel, exclude := "", ""
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = path.Join(dir, "./testdata")
	fromTestdata(t, dir, sel, exclude)
}

// FromTestdata - run test cases from a directory
func fromTestdata(t *testing.T, dir, sel, exclude string) {
	cl.CallBuiltinOp = exec.CallBuiltinOp
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		dirName := dir + "/" + fi.Name()
		testGenGo(t, dirName, sel, exclude)
		testFrom(t, dirName, sel, exclude)
	}
}

func getPkg(pkgs map[string]*ast.Package) *ast.Package {
	for _, pkg := range pkgs {
		return pkg
	}
	return nil
}

func testFrom(t *testing.T, pkgDir, sel, exclude string) {
	var bs = &strings.Builder{}
	selfPrintln := func(arity int, p *gop.Context) {
		args := p.GetArgs(arity)
		n, err := fmt.Fprintln(bs, args...)
		p.Ret(arity, n, err)
	}
	selfPrintf := func(arity int, p *gop.Context) {
		args := p.GetArgs(arity)
		n, err := fmt.Fprintf(bs, args[0].(string), args[1:]...)
		p.Ret(arity, n, err)
	}
	libbuiltin.I.RegisterFuncvs(libbuiltin.I.Funcv("println", fmt.Println, selfPrintln))
	libfmt.I.RegisterFuncvs(libfmt.I.Funcv("Printf", fmt.Printf, selfPrintf))
	libfmt.I.RegisterFuncvs(libfmt.I.Funcv("Println", fmt.Println, selfPrintln))

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
	b := exec.NewBuilder(nil)
	_, err = cl.NewPackage(b.Interface(), bar, fset, cl.PkgActClMain)
	if err != nil {
		if err == cl.ErrNotAMainPackage {
			return
		}
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	out, _ := shell.Command("go", "run", pkgDir).CombinedOutput()
	if bs.String() != string(out) {
		t.Log(pkgDir, "run fail")
		t.Log("byte code result", bs.String())
		t.Log("golang code result", string(out))
		t.Fail()
	}
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
	b := golang.NewBuilder(bar.Name, nil, fset)
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
func saveGoFile(dir string, code *golang.Code) error {
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
