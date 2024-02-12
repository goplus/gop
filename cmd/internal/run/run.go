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

// Package run implements the “gop run” command.
package run

import (
	"fmt"
	"os"
	"reflect"

	"github.com/goplus/gop"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/x/gocmd"
	"github.com/goplus/gop/x/gopprojs"
	"github.com/goplus/gox"
	"github.com/qiniu/x/log"
)

// gop run
var Cmd = &base.Command{
	UsageLine: "gop run [-nc -asm -quiet -debug -prof] package [arguments...]",
	Short:     "Run a Go+ program",
}

var (
	flag        = &Cmd.Flag
	flagAsm     = flag.Bool("asm", false, "generates `asm` code of Go+ bytecode backend")
	flagDebug   = flag.Bool("debug", false, "print debug information")
	flagQuiet   = flag.Bool("quiet", false, "don't generate any compiling stage log")
	flagNoChdir = flag.Bool("nc", false, "don't change dir (only for `gop run pkgPath`)")
	flagProf    = flag.Bool("prof", false, "do profile and generate profile report")
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	pass := base.PassBuildFlags(cmd)
	err := flag.Parse(args)
	if err != nil {
		log.Fatalln("parse input arguments failed:", err)
	}
	if flag.NArg() < 1 {
		cmd.Usage(os.Stderr)
	}

	proj, args, err := gopprojs.ParseOne(flag.Args()...)
	if err != nil {
		log.Fatalln(err)
	}

	if *flagQuiet {
		log.SetOutputLevel(0x7000)
	} else if *flagDebug {
		gox.SetDebug(gox.DbgFlagAll &^ gox.DbgFlagComments)
		cl.SetDebug(cl.DbgFlagAll)
		cl.SetDisableRecover(true)
	} else if *flagAsm {
		gox.SetDebug(gox.DbgFlagInstruction)
	}

	if *flagProf {
		panic("TODO: profile not impl")
	}

	noChdir := *flagNoChdir
	conf, err := gop.NewDefaultConf(".")
	if err != nil {
		log.Panicln("gop.NewDefaultConf:", err)
	}
	confCmd := &gocmd.Config{Gop: conf.Gop}
	confCmd.Flags = pass.Args
	run(proj, args, !noChdir, conf, confCmd)
}

func run(proj gopprojs.Proj, args []string, chDir bool, conf *gop.Config, run *gocmd.RunConfig) {
	var obj string
	var err error
	switch v := proj.(type) {
	case *gopprojs.DirProj:
		obj = v.Dir
		err = gop.RunDir(obj, args, conf, run)
	case *gopprojs.PkgPathProj:
		obj = v.Path
		err = gop.RunPkgPath(v.Path, args, chDir, conf, run)
	case *gopprojs.FilesProj:
		err = gop.RunFiles("", v.Files, args, conf, run)
	default:
		log.Panicln("`gop run` doesn't support", reflect.TypeOf(v))
	}
	if gop.NotFound(err) {
		fmt.Fprintf(os.Stderr, "gop run %v: not found\n", obj)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		return
	}
	os.Exit(1)
}

// -----------------------------------------------------------------------------
