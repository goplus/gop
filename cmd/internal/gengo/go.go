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

// Package gengo implements the “gop go” command.
package gengo

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/goplus/gop"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/x/gopprojs"
	"github.com/goplus/gox"
)

// gop go
var Cmd = &base.Command{
	UsageLine: "gop go [-v] [packages]",
	Short:     "Convert Go+ packages into Go packages",
}

var (
	flag          = &Cmd.Flag
	flagVerbose   = flag.Bool("v", false, "print verbose information.")
	flagCheckMode = flag.Bool("t", false, "check mode, no generate gop_autogen.go")
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	err := flag.Parse(args)
	if err != nil {
		log.Panicln("parse input arguments failed:", err)
	}
	pattern := flag.Args()
	if len(pattern) == 0 {
		pattern = []string{"."}
	}

	projs, err := gopprojs.ParseAll(pattern...)
	if err != nil {
		log.Panicln("gopprojs.ParseAll:", err)
	}

	if *flagVerbose {
		gox.SetDebug(gox.DbgFlagAll &^ gox.DbgFlagComments)
		cl.SetDebug(cl.DbgFlagAll)
		cl.SetDisableRecover(true)
	}

	for _, proj := range projs {
		switch v := proj.(type) {
		case *gopprojs.DirProj:
			_, _, err = gop.GenGoEx(v.Dir, nil, true, *flagCheckMode, true)
		case *gopprojs.PkgPathProj:
			_, _, err = gop.GenGoPkgPathEx("", v.Path, nil, true, *flagCheckMode, true)
		default:
			log.Panicln("`gop go` doesn't support", reflect.TypeOf(v))
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "GenGo failed.")
			os.Exit(1)
		}
	}
}

// -----------------------------------------------------------------------------
