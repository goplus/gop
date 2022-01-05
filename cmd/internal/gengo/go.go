/*
 * Copyright (c) 2021-2021 The GoPlus Authors (goplus.org). All rights reserved.
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

// Package gengo implements the ``gop go'' command.
package gengo

import (
	"os"

	"github.com/qiniu/x/log"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/x/gengo"
	"github.com/goplus/gox"
)

// gop go
var Cmd = &base.Command{
	UsageLine: "gop go [-v -n] [packages]",
	Short:     "Convert Go+ packages into Go packages",
}

var (
	flagDontRun = flag.Bool("n", false, "print the processing flow but do not run them.")
	flagVerbose = flag.Bool("v", false, "print verbose information.")
	flag        = &Cmd.Flag
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	err := flag.Parse(args)
	if err != nil {
		log.Fatalln("parse input arguments failed:", err)
	}
	pattern := flag.Args()
	if len(pattern) == 0 {
		pattern = []string{"."}
	}

	if *flagVerbose {
		gox.SetDebug(gox.DbgFlagAll &^ gox.DbgFlagComments)
		cl.SetDebug(cl.DbgFlagAll)
		cl.SetDisableRecover(true)
	}
	if !gengo.GenGo(gengo.Config{}, *flagDontRun, ".", pattern...) {
		os.Exit(1)
	}
}

// -----------------------------------------------------------------------------
