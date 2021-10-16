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
	"fmt"
	"github.com/qiniu/x/log"
	"os"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gox"
)

// Cmd - gop install
var Cmd = &base.Command{
	UsageLine: "gop install [-v] <GopPackages>",
	Short:     "Build Go+ files and install target to GOBIN",
}

var (
	flag        = &Cmd.Flag
	flagVerbose = flag.Bool("v", false, "print verbose information")
)

func init() {
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
	base.GenGoForBuild(dir, recursive, func() { fmt.Fprintln(os.Stderr, "GenGo failed, stop installing") })
	base.RunGoCmd(dir, "install", args...)
}

// -----------------------------------------------------------------------------
