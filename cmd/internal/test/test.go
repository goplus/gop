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

// Package test implements the ``gop test'' command.
package test

import (
	"fmt"
	"os"
	"strings"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/gengo"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gox"
)

// Cmd - gop install
var Cmd = &base.Command{
	UsageLine: "gop test [-v] <GopPackages>",
	Short:     "Test Go+ packages",
}

var (
	flag        = &Cmd.Flag
	flagVerbose = flag.Bool("v", false, "print verbose information")
	flagRebuild = flag.Bool("a", false, "force rebuilding of packages that are already up-to-date")
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	flag.Parse(base.SkipSwitches(args, flag))
	ssargs := flag.Args()
	if len(ssargs) == 0 {
		ssargs = []string{"."}
	}
	var recursive bool
	var dir = ssargs[0]
	if strings.HasSuffix(dir, "/...") {
		dir = dir[:len(dir)-4]
		recursive = true
	}

	if *flagVerbose {
		gox.SetDebug(gox.DbgFlagAll &^ gox.DbgFlagComments)
		cl.SetDebug(cl.DbgFlagAll)
		cl.SetDisableRecover(true)
	}
	hasError := false
	runner := new(gengo.Runner)
	runner.SetAfter(func(p *gengo.Runner, dir string, flags int) error {
		errs := p.ResetErrors()
		if errs != nil {
			hasError = true
			for _, err := range errs {
				fmt.Fprintln(os.Stderr, err)
			}
			fmt.Fprintln(os.Stderr)
		}
		return nil
	})
	baseConf := &cl.Config{PersistLoadPkgs: true}
	runner.GenGo(dir, recursive, *flagRebuild, baseConf.Ensure())
	if hasError {
		os.Exit(1)
	}
	baseConf.PkgsLoader.Save()
	args = removeRebuild(args)
	base.RunGoCmd(dir, "test", args...)
}

func removeRebuild(args []string) (r []string) {
	for _, a := range args {
		if a == "-a" {
			continue
		}
		r = append(r, a)
	}
	return
}

// -----------------------------------------------------------------------------
