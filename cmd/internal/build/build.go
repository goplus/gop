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
	"log"
	"os"
	"path/filepath"
	"reflect"
	"syscall"

	"github.com/goplus/gop"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/x/gocmd"
	"github.com/goplus/gop/x/gopenv"
	"github.com/goplus/gop/x/gopprojs"
	"github.com/goplus/gox"
)

// gop build
var Cmd = &base.Command{
	UsageLine: "gop build [-v -o output] [packages]",
	Short:     "Build Go+ files",
}

var (
	flagVerbose = flag.Bool("v", false, "print verbose information")
	flagOutput  = flag.String("o", "", "gop build output file")
	flag        = &Cmd.Flag
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(_ *base.Command, args []string) {
	err := flag.Parse(base.SkipSwitches(args, flag))
	if err != nil {
		log.Panicln("parse input arguments failed:", err)
	}

	if *flagVerbose {
		gox.SetDebug(gox.DbgFlagAll &^ gox.DbgFlagComments)
		cl.SetDebug(cl.DbgFlagAll)
		cl.SetDisableRecover(true)
	}

	sargs := flag.Args()
	if len(sargs) == 0 {
		sargs = []string{"."}
	}
	proj, sargs, err := gopprojs.ParseOne(sargs...)

	if err != nil {
		log.Panicln(err)
	}
	if len(sargs) != 0 {
		log.Panicln("too many arguments:", sargs)
	}

	gopEnv := gopenv.Get()
	conf := &gop.Config{Gop: gopEnv}
	confCmd := &gocmd.BuildConfig{Gop: gopEnv}
	if *flagOutput != "" {
		output, err := filepath.Abs(*flagOutput)
		if err != nil {
			log.Panicln(err)
		}
		confCmd.Flags = []string{"-o", output}
	}
	confCmd.Flags = append(confCmd.Flags, args...)
	build(proj, conf, confCmd)
}

func build(proj gopprojs.Proj, conf *gop.Config, build *gocmd.BuildConfig) {
	var obj string
	var err error
	switch v := proj.(type) {
	case *gopprojs.DirProj:
		obj = v.Dir
		err = gop.BuildDir(obj, conf, build)
	case *gopprojs.PkgPathProj:
		obj = v.Path
		err = gop.BuildPkgPath("", v.Path, conf, build)
	case *gopprojs.FilesProj:
		err = gop.BuildFiles(v.Files, conf, build)
	default:
		log.Panicln("`gop build` doesn't support", reflect.TypeOf(v))
	}
	if err == syscall.ENOENT {
		fmt.Fprintf(os.Stderr, "gop build %v: not found\n", obj)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		return
	}
	os.Exit(1)
}

// -----------------------------------------------------------------------------
