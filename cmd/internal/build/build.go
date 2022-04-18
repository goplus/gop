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

// Package build implements the ``gop build'' command.
package build

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/qiniu/x/log"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/x/gopproj"
	"github.com/goplus/gop/x/gopprojs"
	"github.com/goplus/gox"
)

// gop build
var Cmd = &base.Command{
	UsageLine: "gop build [-v -o output] [packages]",
	Short:     "Build Go+ files",
}

var (
	flagVerbose = flag.Bool("v", false, "print verbose information.")
	flagOutput  = flag.String("o", "a.out", "gop build output file.")
	flag        = &Cmd.Flag
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(_ *base.Command, args []string) {
	err := flag.Parse(args)
	if err != nil {
		log.Fatalln("parse input arguments failed:", err)
	}

	if *flagVerbose {
		gox.SetDebug(gox.DbgFlagAll &^ gox.DbgFlagComments)
		cl.SetDebug(cl.DbgFlagAll)
		cl.SetDisableRecover(true)
	}

	args = flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}
	gopBuild(args)
}

func gopBuild(args []string) {
	proj, args, err := gopprojs.ParseOne(args...)
	if err != nil {
		log.Fatalln(err)
	}

	flags := 0
	ctx, goProj, err := gopproj.OpenProject(flags, proj)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	goProj.BuildArgs = []string{"-o", *flagOutput}
	cmd := ctx.GoCommand("build", goProj)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	err = cmd.Run()
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			os.Exit(e.ExitCode())
		default:
			log.Fatalln(err)
		}
	}
}

// -----------------------------------------------------------------------------
