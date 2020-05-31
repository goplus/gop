package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/qiniu/qlang/v6/cl"
	"github.com/qiniu/qlang/v6/exec/golang"
	"github.com/qiniu/qlang/v6/parser"
	"github.com/qiniu/qlang/v6/token"
	"github.com/qiniu/x/log"

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

func genGopkg(pkgDir string) (err error) {
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

func genGo(dir string) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ReadDir failed:", err)
		return
	}
	var isPkg bool
	for _, fi := range fis {
		if fi.IsDir() {
			pkgDir := path.Join(dir, fi.Name())
			genGo(pkgDir)
			continue
		}
		if strings.HasSuffix(fi.Name(), ".ql") {
			isPkg = true
		}
	}
	if isPkg {
		fmt.Printf("Compiling %s ...\n", dir)
		err = genGopkg(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "# %v\n[ERROR] %v\n\n", dir, err)
		}
	}
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: qgo <qlangSrcDir>")
		return
	}
	dir := os.Args[1]
	cl.CallBuiltinOp = exec.CallBuiltinOp
	log.SetFlags(log.Ldefault &^ log.LstdFlags)
	genGo(dir)
}
