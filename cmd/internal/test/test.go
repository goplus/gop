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
	"os/exec"
	"strings"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/gengo"
	"github.com/goplus/gop/cmd/internal/base"
)

// Cmd - gop install
var Cmd = &base.Command{
	UsageLine: "gop test [-v] <GopPackages>",
	Short:     "Test Go+ packages",
}

var (
	flag = &Cmd.Flag
	_    = flag.Bool("v", false, "print the names of packages as they are compiled.")
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
	var exitCode int
	var recursive bool
	var dir = flag.Arg(0)
	if strings.HasSuffix(dir, "/...") {
		dir = dir[:len(dir)-4]
		recursive = true
	}

	runner := new(gengo.Runner)
	runner.SetAfter(func(p *gengo.Runner, dir string, flags int) error {
		errs := p.ResetErrors()
		if errs != nil {
			for _, err := range errs {
				fmt.Fprintln(os.Stderr, err)
			}
			fmt.Fprintln(os.Stderr)
		}
		return nil
	})
	runner.GenGo(dir, recursive, &cl.Config{CacheLoadPkgs: true})
	goCmd(dir, "test", args...)
	os.Exit(exitCode)
}

func goCmd(dir string, op string, args ...string) error {
	opwargs := make([]string, len(args)+1)
	opwargs[0] = op
	copy(opwargs[1:], args)
	cmd := exec.Command("go", opwargs...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

// -----------------------------------------------------------------------------
