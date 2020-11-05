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

// Package build implements the ``gop build'' command.
package build

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/cmd/internal/work"
	"github.com/goplus/gop/exec/bytecode"
	"github.com/qiniu/x/log"
)

var (
	exitCode = 0
)

// -----------------------------------------------------------------------------

// Cmd - gop go
var Cmd = &base.Command{
	UsageLine: "gop build [-v] [-o output] <gopSrcDir|gopSrcFile>",
	Short:     "Build for all go+ files and execute go build command",
}

var (
	flagBuildOutput string
	flagVerbose     bool
	flag            = &Cmd.Flag
)

func init() {
	flag.StringVar(&flagBuildOutput, "o", "", "go build output file")
	flag.BoolVar(&flagVerbose, "v", false, "print the names of packages as they are compiled.")
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	err := flag.Parse(args)
	if err != nil {
		cmd.Usage(os.Stderr)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Fail to build: %v", err)
	}

	paths := flag.Args()
	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	cl.CallBuiltinOp = bytecode.CallBuiltinOp
	log.SetFlags(log.Ldefault &^ log.LstdFlags)

	fset := token.NewFileSet()
	pkgs, errs := work.LoadPackages(fset, paths)
	if len(errs) > 0 {
		log.Fatalf("load packages error: %v\n", errs)
	}
	for _, pkg := range pkgs {
		err := work.GenGoPkg(fset, pkg.Pkg, pkg.Dir)
		if err != nil {
			log.Fatalf("generate go package error: %v\n", err)
		}
		var target string
		if flagBuildOutput != "" {
			target = filepath.Join(dir, flagBuildOutput)
		} else {
			target = pkg.Target
		}
		err = work.GoBuild(pkg.Dir, target)
		if err != nil {
			log.Fatalf("go build error: %v\n", err)
		}
		if flagVerbose {
			fmt.Printf("gop build %v\n", target)
		}
	}
}
