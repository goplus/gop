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

// Package install implements the ``gop install'' command.
package install

import (
	"os"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/cmd/internal/work"
	"github.com/goplus/gop/exec/bytecode"
	"github.com/qiniu/x/log"
)

var (
	exitCode = 0
)

// Cmd - gop go
var Cmd = &base.Command{
	UsageLine: "gop install [-v] <gopSrcDir>",
	Short:     "Install compiles and installs the packages named by the import paths.",
}

var (
	flag = &Cmd.Flag
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
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Fail to build: %v", err)
	}

	cl.CallBuiltinOp = bytecode.CallBuiltinOp
	log.SetFlags(log.Ldefault &^ log.LstdFlags)
	err = runInstall(args, dir)
	if err != nil {
		exitCode = -1
	}
	os.Exit(exitCode)
}

// -----------------------------------------------------------------------------

func runInstall(args []string, wd string) error {
	gopInstall, err := work.NewInstall("", args, wd)
	if err != nil {
		log.Fatalf("Fail to install: %v", err)
		return err
	}

	err = gopInstall.Install()
	if err != nil {
		return err
	}
	return nil
}
