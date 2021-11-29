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

	"github.com/qiniu/x/log"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/cmd/internal/modload"
	"github.com/goplus/gox"
)

// -----------------------------------------------------------------------------

// Cmd - gop build
var Cmd = &base.Command{
	UsageLine: "gop build [-v] [-o output] <gopSrcDir|gopSrcFile>",
	Short:     "Build Go+ files",
}

var (
	flagBuildOutput string
	flagVerbose     = flag.Bool("v", false, "print verbose information")
	flag            = &Cmd.Flag
)

func init() {
	flag.StringVar(&flagBuildOutput, "o", "", "gop build output file")
	Cmd.Run = runCmd
}

func runCmd(_ *base.Command, args []string) {
	err := flag.Parse(base.SkipSwitches(args, flag))
	if err != nil {
		log.Fatalln("parse input arguments failed:", err)
	}
	ssargs := flag.Args()
	dir, recursive := base.GetBuildDir(ssargs)

	if *flagVerbose {
		gox.SetDebug(gox.DbgFlagAll &^ gox.DbgFlagComments)
		cl.SetDebug(cl.DbgFlagAll)
		cl.SetDisableRecover(true)
	}
	modload.Load()
	base.GenGoForBuild(dir, recursive, func() { fmt.Fprintln(os.Stderr, "GenGo failed, stop building") })
	base.RunGoCmd(dir, "build", args...)
}

// -----------------------------------------------------------------------------
