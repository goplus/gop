package cltest

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/qiniu/x/log"

	exec "github.com/goplus/gop/exec/bytecode"
	_ "github.com/goplus/gop/lib" // libraries
)

// -----------------------------------------------------------------------------

func getPkg(pkgs map[string]*ast.Package) *ast.Package {
	for _, pkg := range pkgs {
		return pkg
	}
	return nil
}

func testFrom(t *testing.T, pkgDir, sel, exclude string) {
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
	exec.NewContext(code).Run()
}

// FromTestdata - run test cases from a directory
func FromTestdata(t *testing.T, dir, sel, exclude string) {
	cl.CallBuiltinOp = exec.CallBuiltinOp
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		testFrom(t, dir+"/"+fi.Name(), sel, exclude)
	}
}

// -----------------------------------------------------------------------------
