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

// Package gengo implements the ``gop go'' command.
package gengo

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/exec/bytecode"
	"github.com/goplus/gop/exec/golang"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/qiniu/x/log"
)

var (
	exitCode = 0
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
	data, err := format.Source(b)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dir+"/gop_autogen.go", data, 0666)
}

func genGopkg(pkgDir string) (mainPkg bool, err error) {
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
	return pkg.Name == "main", saveGoFile(pkgDir, code)
}

func getPkg(pkgs map[string]*ast.Package) *ast.Package {
	for _, pkg := range pkgs {
		return pkg
	}
	return nil
}

// -----------------------------------------------------------------------------

func testPkg(dir string) {
	cmd1 := exec.Command("go", "run", path.Join(dir, "gop_autogen.go"))
	gorun, err := cmd1.CombinedOutput()
	if err != nil {
		os.Stderr.Write(gorun)
		fmt.Fprintf(os.Stderr, "[ERROR] `%v` failed: %v\n", cmd1, err)
		exitCode = -1
		return
	}
	cmd2 := exec.Command("qrun", "-quiet", dir) // -quiet: don't generate any log
	qrun, err := cmd2.CombinedOutput()
	if err != nil {
		os.Stderr.Write(qrun)
		fmt.Fprintf(os.Stderr, "[ERROR] `%v` failed: %v\n", cmd2, err)
		exitCode = -1
		return
	}
	if !bytes.Equal(gorun, qrun) {
		fmt.Fprintf(os.Stderr, "[ERROR] Output has differences!\n")
		fmt.Fprintf(os.Stderr, ">>> Output of `%v`:\n", cmd1)
		os.Stderr.Write(gorun)
		fmt.Fprintf(os.Stderr, "\n>>> Output of `%v`:\n", cmd2)
		os.Stderr.Write(qrun)
		exitCode = -1
	}
}

// -----------------------------------------------------------------------------

func genGo(dir string, test bool) {
	if strings.HasPrefix(dir, "_") {
		return
	}
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ReadDir failed:", err)
		exitCode = -1
		return
	}
	var isPkg bool
	for _, fi := range fis {
		if fi.IsDir() {
			pkgDir := path.Join(dir, fi.Name())
			genGo(pkgDir, test)
			continue
		}
		if strings.HasSuffix(fi.Name(), ".gop") {
			isPkg = true
		}
	}
	if isPkg {
		fmt.Printf("Compiling %s ...\n", dir)
		isPkg, err = genGopkg(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] %v\n\n", err)
			exitCode = -1
		} else if isPkg && test {
			fmt.Printf("Testing %s ...\n", dir)
			testPkg(dir)
		}
	}
}

// -----------------------------------------------------------------------------

// Cmd - gop go
var Cmd = &base.Command{
	UsageLine: "gop go [-test] <gopSrcDir>",
	Short:     "Convert Go+ packages into Go packages",
}

var (
	flag     = &Cmd.Flag
	flagTest = flag.Bool("test", false, "test Go+ package")
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	flag.Parse(args)
	if flag.NArg() < 1 {
		cmd.Usage(os.Stderr)
		return
	}
	dir := flag.Arg(0)
	cl.CallBuiltinOp = bytecode.CallBuiltinOp
	log.SetFlags(log.Ldefault &^ log.LstdFlags)
	genGo(dir, *flagTest)
	os.Exit(exitCode)
}

// -----------------------------------------------------------------------------
