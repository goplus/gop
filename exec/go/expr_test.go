package exec

import (
	"go/token"
	"io/ioutil"
	"os"
	"testing"

	"github.com/qiniu/qlang/v6/cl"
	"github.com/qiniu/qlang/v6/parser"
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
	return ioutil.WriteFile(dir+"/qlang_autogen.go", b, 0666)
}

// -----------------------------------------------------------------------------

func testGenGo(t *testing.T, pkgDir string) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgDir, nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := NewBuilder(nil, fset)
	_, err = cl.NewPackage(b.Interface(), bar, fset)
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
		testGenGo(t, dir+"/"+fi.Name())
	}
}

// -----------------------------------------------------------------------------
