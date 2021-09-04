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

package base

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
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

// RunGoCmd executes `go` command tools.
func RunGoCmd(dir string, op string, args ...string) {
	opwargs := make([]string, len(args)+1)
	opwargs[0] = op
	copy(opwargs[1:], args)
	cmd := exec.Command("go", opwargs...)
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
