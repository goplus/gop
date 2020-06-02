package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/qiniu/qlang/v6/ast"
	"github.com/qiniu/qlang/v6/cl"
	"github.com/qiniu/qlang/v6/exec/golang"
	"github.com/qiniu/qlang/v6/parser"
	"github.com/qiniu/qlang/v6/token"
	"github.com/qiniu/x/log"

	"github.com/qiniu/qlang/v6/exec/bytecode"
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

	pkg := getPkg(pkgs)
	b := golang.NewBuilder(pkg.Name, nil, fset)
	_, err = cl.NewPackage(b.Interface(), pkg, fset, cl.PkgActClAll)
	if err != nil {
		return
	}
	code := b.Resolve()
	return saveGoFile(pkgDir, code)
}

func getPkg(pkgs map[string]*ast.Package) *ast.Package {
	for _, pkg := range pkgs {
		return pkg
	}
	return nil
}

// -----------------------------------------------------------------------------

func testPkg(dir string) {
	cmd1 := exec.Command("go", "run", path.Join(dir, "qlang_autogen.go"))
	gorun, err := cmd1.CombinedOutput()
	if err != nil {
		os.Stderr.Write(gorun)
		fmt.Fprintf(os.Stderr, "[ERROR] `%v` failed: %v\n", cmd1, err)
		return
	}
	cmd2 := exec.Command("qrun", "-quiet", dir) // -quiet: don't generate any log
	qrun, err := cmd2.CombinedOutput()
	if err != nil {
		os.Stderr.Write(qrun)
		fmt.Fprintf(os.Stderr, "[ERROR] `%v` failed: %v\n", cmd2, err)
		return
	}
	if !bytes.Equal(gorun, qrun) {
		fmt.Fprintf(os.Stderr, "[ERROR] Output has differences!\n")
		fmt.Fprintf(os.Stderr, ">>> Output of `%v`:\n", cmd1)
		os.Stderr.Write(gorun)
		fmt.Fprintf(os.Stderr, "\n>>> Output of `%v`:\n", cmd2)
		os.Stderr.Write(qrun)
	}
}

// -----------------------------------------------------------------------------

func genGo(dir string, test bool) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ReadDir failed:", err)
		return
	}
	var isPkg bool
	for _, fi := range fis {
		if fi.IsDir() {
			pkgDir := path.Join(dir, fi.Name())
			genGo(pkgDir, test)
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
			fmt.Fprintf(os.Stderr, "[ERROR] %v\n\n", err)
		} else if test {
			fmt.Printf("Testing %s ...\n", dir)
			testPkg(dir)
		}
	}
}

// -----------------------------------------------------------------------------

var (
	flagTest = flag.Bool("test", false, "test qlang package")
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: qgo [-test] <qlangSrcDir>")
		flag.PrintDefaults()
		return
	}
	dir := flag.Arg(0)
	cl.CallBuiltinOp = bytecode.CallBuiltinOp
	log.SetFlags(log.Ldefault &^ log.LstdFlags)
	genGo(dir, *flagTest)
}

// -----------------------------------------------------------------------------
