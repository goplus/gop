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

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/x/gopproj"
	"github.com/goplus/gop/x/gopprojs"
	"github.com/goplus/gox"
	"github.com/qiniu/x/log"
)

// gop run
var Cmd = &base.Command{
	UsageLine: "gop run [-asm -quiet -debug -nr -gop -prof] package [arguments...]",
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

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	err := flag.Parse(args)
	if err != nil {
		log.Fatalln("parse input arguments failed:", err)
	}
	if flag.NArg() < 1 {
		cmd.Usage(os.Stderr)
	}
	gopRun(flag.Args())
}

func gopRun(args []string) {
	proj, args, err := gopprojs.ParseOne(args...)
	if err != nil {
		log.Fatalln(err)
	}

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

	flags := 0
	if *flagGop {
		flags = gopproj.FlagGoAsGoPlus
	}
	var ctx = gopproj.New("")
	goProj, err := ctx.OpenProject(flags, proj)
	if err != nil {
		fmt.Fprintln(os.Stderr, "OpenProject failed:", err)
		os.Exit(1)
	}
	goProj.ExecArgs = args
	goProj.FlagNRINC = *flagNorun
	goProj.FlagRTOE = *flagRTOE
	if goProj.FlagRTOE {
		goProj.UseDefaultCtx = true
	}
	cmd := ctx.GoCommand("run", goProj)
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
