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

// Package install implements the ``gop install'' command.
package install

import (
	"github.com/goplus/gop/cmd/internal/base"
)

// Cmd - gop install
var Cmd = &base.Command{
	UsageLine: "gop install [-v] <gopSrcDir|gopSrcFile>",
	Short:     "Build go+ files and install target to GOBIN",
}

var (
	flag        = &Cmd.Flag
	flagVerbose bool
)

func init() {
	Cmd.Run = runCmd
	flag.BoolVar(&flagVerbose, "v", false, "print the names of packages as they are compiled.")
}

func runCmd(cmd *base.Command, args []string) {
	panic("TODO: go install not impl")
	/*
		flag.Parse(args)

		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Fail to build: %v", err)
		}

		paths := flag.Args()
		if len(paths) == 0 {
			paths = append(paths, dir)
		}

		log.SetFlags(log.Ldefault &^ log.LstdFlags)

		fset := token.NewFileSet()
		pkgs, errs := work.LoadPackages(fset, paths)
		if len(errs) > 0 {
			log.Fatalf("load packages error: %v\n", errs)
		}
		if len(pkgs) == 0 {
			fmt.Println("no Go+ files in ", paths)
		}
		for _, pkg := range pkgs {
			err := work.GenGoPkg(fset, pkg.Pkg, pkg.Dir)
			if err != nil {
				log.Fatalf("generate go package error: %v\n", err)
			}
			target := filepath.Join(work.GOPBIN(), pkg.Target)
			if pkg.IsDir {
				err = work.GoInstall(pkg.Dir)
			} else {
				err = work.GoBuild(pkg.Dir, target, "gop_autogen.go")
			}
			if err != nil {
				log.Fatalf("go build error: %v\n", err)
			}
			if flagVerbose && pkg.Name == "main" {
				fmt.Println(target)
			}
		}
	*/
}
