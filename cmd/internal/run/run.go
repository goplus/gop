/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package run implements the ``gop run'' command.
package run

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/qiniu/x/log"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/scanner"
	"github.com/goplus/gop/token"
	"github.com/goplus/gop/x/mod/modload"
	"github.com/goplus/gox"
)

// -----------------------------------------------------------------------------

// Cmd - gop run
var Cmd = &base.Command{
	UsageLine: "gop run [-asm -quiet -debug -nr -gop -prof] <gopSrcDir|gopSrcFile>",
	Short:     "Run a Go+ program",
}

var (
	flag        = &Cmd.Flag
	flagAsm     = flag.Bool("asm", false, "generates `asm` code of Go+ bytecode backend")
	flagVerbose = flag.Bool("v", false, "print verbose information")
	flagQuiet   = flag.Bool("quiet", false, "don't generate any compiling stage log")
	flagDebug   = flag.Bool("debug", false, "set log level to debug")
	flagNorun   = flag.Bool("nr", false, "don't run if no change")
	flagRTOE    = flag.Bool("rtoe", false, "remove tempfile on error")
	flagGop     = flag.Bool("gop", false, "parse a .go file as a .gop file")
	flagProf    = flag.Bool("prof", false, "do profile and generate profile report")
)

const (
	parserMode = parser.ParseComments
)

func init() {
	Cmd.Run = runCmd
}

func saveGoFile(gofile string, pkg *gox.Package) error {
	dir := filepath.Dir(gofile)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	return gox.WriteFile(gofile, pkg, false)
}

func runCmd(cmd *base.Command, args []string) {
	err := flag.Parse(args)
	if err != nil {
		log.Fatalln("parse input arguments failed:", err)
	}
	if flag.NArg() < 1 {
		cmd.Usage(os.Stderr)
	}
	args = flag.Args()[1:]

	if *flagQuiet {
		log.SetOutputLevel(0x7000)
	} else if *flagDebug {
		log.SetOutputLevel(log.Ldebug)
		gox.SetDebug(gox.DbgFlagAll)
		cl.SetDebug(cl.DbgFlagAll)
	}
	if *flagVerbose {
		gox.SetDebug(gox.DbgFlagAll &^ gox.DbgFlagComments)
		cl.SetDebug(cl.DbgFlagAll)
		cl.SetDisableRecover(true)
	} else if *flagAsm {
		gox.SetDebug(gox.DbgFlagInstruction)
	}
	if *flagProf {
		panic("TODO: profile not impl")
	}

	fset := token.NewFileSet()
	src, _ := filepath.Abs(flag.Arg(0))
	fi, err := os.Stat(src)
	if err != nil {
		log.Fatalln("input arg check failed:", err)
	}
	isDir := fi.IsDir()

	var isDirty bool
	var srcDir, gofile string
	var pkgs map[string]*ast.Package
	if isDir {
		srcDir = src
		gofile = src + "/gop_autogen.go"
		modload.UpdateGoMod(srcDir)
		isDirty = true // TODO: check if code changed
		if isDirty {
			pkgs, err = parser.ParseDir(fset, src, nil, parserMode)
		} else if *flagNorun {
			return
		}
	} else {
		gopRun(src, args...)
		return
	}
	if err != nil {
		scanner.PrintError(os.Stderr, err)
		os.Exit(10)
	}

	if isDirty {
		mainPkg, ok := pkgs["main"]
		if !ok {
			if len(pkgs) == 0 && isDir { // not a Go+ package, try runGoPkg
				runGoPkg(src, args, true)
				return
			}
			fmt.Fprintln(os.Stderr, "TODO: not a main package")
			os.Exit(12)
		}

		modDir := findGoModDir(srcDir)
		conf := &cl.Config{Dir: modDir, TargetDir: srcDir, Fset: fset}
		out, err := cl.NewPackage("", mainPkg, conf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(11)
		}
		err = saveGoFile(gofile, out)
		if err != nil {
			log.Fatalln("saveGoFile failed:", err)
		}
	}

	goRun(gofile, args)
	if *flagProf {
		panic("TODO: profile not impl")
	}
}

func goRun(file string, args []string) {
	goArgs := make([]string, len(args)+2)
	goArgs[0] = "run"
	goArgs[1] = file
	copy(goArgs[2:], args)
	cmd := exec.Command("go", goArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			os.Exit(e.ExitCode())
		default:
			log.Fatalln("go run failed:", err)
		}
	}
}

func runGoPkg(src string, args []string, doRun bool) {
	if doRun {
		goRun(src+"/.", args)
	}
}

// -----------------------------------------------------------------------------
