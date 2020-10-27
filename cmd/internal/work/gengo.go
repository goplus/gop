package work

import (
	"errors"
	"fmt"
	"go/format"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/exec/golang"
	"github.com/goplus/gop/parser"
)

func GenGo(dir, toDir string) error {
	if strings.HasPrefix(dir, "_") {
		return nil
	}
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ReadDir failed:", err)
		return err
	}
	var isPkg bool
	for _, fi := range fis {
		if fi.IsDir() {
			pkgDir := path.Join(dir, fi.Name())
			err = GenGo(pkgDir, toDir)
			if err != nil {
				return err
			}
			continue
		}
		if strings.HasSuffix(fi.Name(), ".gop") {
			isPkg = true
		}
	}
	if isPkg {
		fmt.Printf("Compiling %s ...\n", dir)
		isPkg, err = genGopkg(dir, toDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] %v\n\n", err)
			return err
		}
	}
	return nil
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
	data, err := format.Source(b)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dir+"/gop_autogen.go", data, 0666)
}

func genGopkg(pkgDir, tmpDir string) (mainPkg bool, err error) {
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
	pkgDir, _ = filepath.Abs(pkgDir)
	pkgs, err := parser.ParseDir(fset, pkgDir, nil, 0)
	if err != nil {
		return
	}
	if len(pkgs) != 1 {
		return false, fmt.Errorf("too many packages (%d) in the same directory", len(pkgs))
	}

	pkg := getPkg(pkgs)
	b := golang.NewBuilder(pkg.Name, nil, fset)
	_, err = cl.NewPackage(b.Interface(), pkg, fset, cl.PkgActClAll)
	if err != nil {
		return
	}
	code := b.Resolve()
	return pkg.Name == "main", saveGoFile(tmpDir, code)
}

func getPkg(pkgs map[string]*ast.Package) *ast.Package {
	for _, pkg := range pkgs {
		return pkg
	}
	return nil
}
