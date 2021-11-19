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

package base

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/gengo"
)

// SkipSwitches skips all switches and returns non-switch arguments.
func SkipSwitches(args []string, f *flag.FlagSet) []string {
	out := make([]string, 0, len(args))
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if f.Lookup(arg[1:]) == nil { // flag not found
				continue
			}
		}
		out = append(out, arg)
	}
	return out
}

// GetBuildDir Get build directory from arguments
func GetBuildDir(args []string) (dir string, recursive bool) {
	if len(args) == 0 {
		args = []string{"."}
	}
	dir = args[0]
	if strings.HasSuffix(dir, "/...") {
		dir = dir[:len(dir)-4]
		recursive = true
	}
	if fi, err := os.Stat(dir); err == nil {
		if fi.IsDir() {
			return
		}
		return filepath.Dir(dir), recursive
	}
	return
}

// GenGoForBuild Generate go code before building or installing, and cache pkgs if success
func GenGoForBuild(dir string, recursive bool, errorHandle func()) {
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
	runner.GenGo(dir, recursive, baseConf.Ensure())
	if hasError {
		errorHandle()
		os.Exit(1)
	}
	baseConf.PkgsLoader.Save()
}

// RunGoCmd executes `go` command tools.
func RunGoCmd(dir string, op string, args ...string) {
	cmd := exec.Command("go", append([]string{op}, args...)...)
	cmd.Dir = dir
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
			log.Fatalln("RunGoCmd failed:", err)
		}
	}
}
