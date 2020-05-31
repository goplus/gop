package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/qiniu/qlang/v6/cl"
	"github.com/qiniu/qlang/v6/exec/golang"
	"github.com/qiniu/qlang/v6/parser"
	"github.com/qiniu/qlang/v6/token"

	exec "github.com/qiniu/qlang/v6/exec/bytecode"
	_ "github.com/qiniu/qlang/v6/lib/builtin"
	_ "github.com/qiniu/qlang/v6/lib/fmt"
	_ "github.com/qiniu/qlang/v6/lib/reflect"
	_ "github.com/qiniu/qlang/v6/lib/strings"
)

func saveGoFile(dir string, code *golang.Code) error {
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

func genGo(pkgDir string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch v := e.(type) {
			case string:
				err = errors.New(v)
			case error:
				err = v
			default:
				panic(e)
			}
		}
	}()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgDir, nil, 0)
	if err != nil {
		return
	}
	if len(pkgs) != 1 {
		return fmt.Errorf("too many package in the same directory")
	}

	bar := pkgs["main"]
	b := golang.NewBuilder(nil, fset)
	_, err = cl.NewPackage(b.Interface(), bar, fset)
	if err != nil {
		return
	}
	code := b.Resolve()
	return saveGoFile(pkgDir, code)
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: qgo <qlangSrcDir>")
		return
	}
	dir := os.Args[1]
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ReadDir failed:", err)
		return
	}
	cl.CallBuiltinOp = exec.CallBuiltinOp
	for _, fi := range fis {
		pkgDir := path.Join(dir, fi.Name())
		fmt.Printf("Compiling %s ...\n", pkgDir)
		if err = genGo(pkgDir); err != nil {
			fmt.Fprintf(os.Stderr, "# %v\n[ERROR] %v\n\n", pkgDir, err)
		}
	}
}
