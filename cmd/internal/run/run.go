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

// Package run implements the ``gop run'' command.
package run

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// Cmd - gop run
var Cmd = &base.Command{
	UsageLine: "gop run [-asm -quiet -debug -prof] <gopSrcDir|gopSrcFile>",
	Short:     "Run a Go+ program",
}

var (
	flag        = &Cmd.Flag
	flagAsm     = flag.Bool("asm", false, "generates `asm` code of Go+ bytecode backend")
	flagVerbose = flag.Bool("v", false, "print verbose information")
	flagQuiet   = flag.Bool("quiet", false, "don't generate any compiling stage log")
	flagDebug   = flag.Bool("debug", false, "set log level to debug")
	flagProf    = flag.Bool("prof", false, "do profile and generate profile report")
)

func init() {
	Cmd.Run = runCmd
}

func saveGoFile(gofile string, pkg *gox.Package) error {
	dir := filepath.Dir(gofile)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	return gox.WriteFile(gofile, pkg)
}

func findGoModFile(dir string) (string, error) {
	modfile, err := cl.FindGoModFile(dir)
	if err != nil {
		home := os.Getenv("HOME")
		modfile = home + "/gop/go.mod"
		if fi, e := os.Lstat(modfile); e == nil && !fi.IsDir() {
			return modfile, nil
		}
		modfile = home + "/goplus/go.mod"
		if fi, e := os.Lstat(modfile); e == nil && !fi.IsDir() {
			return modfile, nil
		}
	}
	return modfile, err
}

func findGoModDir(dir string) string {
	modfile, err := findGoModFile(dir)
	if err != nil {
		log.Fatalln("findGoModFile:", err)
	}
	return filepath.Dir(modfile)
}

func runCmd(cmd *base.Command, args []string) {
	flag.Parse(args)
	if flag.NArg() < 1 {
		cmd.Usage(os.Stderr)
	}

	if *flagQuiet {
		log.SetOutputLevel(0x7000)
	} else if *flagDebug {
		log.SetOutputLevel(log.Ldebug)
		gox.SetDebug(gox.DbgFlagAll)
	}
	if *flagVerbose {
		gox.SetDebug(gox.DbgFlagAll &^ gox.DbgFlagComments)
	} else if *flagAsm {
		gox.SetDebug(gox.DbgFlagInstruction)
	}
	if *flagProf {
		panic("TODO: profile not impl")
	}
	fset := token.NewFileSet()

	target, _ := filepath.Abs(flag.Arg(0))
	isDir, err := IsDir(target)
	if err != nil {
		log.Fatalln("input arg check failed:", err)
	}
	var dir string
	var pkgs map[string]*ast.Package
	if isDir {
		dir = target
		pkgs, err = parser.ParseDir(fset, target, nil, 0)
	} else {
		dir = filepath.Dir(target)
		pkgs, err = parser.Parse(fset, target, nil, 0)
	}
	if err != nil {
		log.Fatalln("parser.Parse failed:", err)
	}

	conf := &cl.Config{Dir: findGoModDir(dir)}
	out, err := cl.NewPackage("", pkgs["main"], fset, conf)
	if err != nil {
		log.Fatalln("cl.NewPackage failed:", err)
	}
	var gofile string
	if isDir {
		gofile = target + "/gop_autogen.go"
	} else {
		dir, file := filepath.Split(target)
		gofile = filepath.Join(dir, ".gop", file+".go")
	}
	err = saveGoFile(gofile, out)
	if err != nil {
		log.Fatalln("saveGoFile failed:", err)
	}
	err = goRun(gofile)
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			os.Stderr.Write(e.Stderr)
		default:
			log.Fatalln("go run failed:", err)
		}
	}
	if *flagProf {
		panic("TODO: profile not impl")
	}
}

// IsDir checks a target path is dir or not.
func IsDir(target string) (bool, error) {
	fi, err := os.Stat(target)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

func goRun(target string) error {
	cmd := exec.Command("go", "run", target)
	cmd.Dir, _ = filepath.Split(target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

// -----------------------------------------------------------------------------
