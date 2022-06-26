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

// Package install implements the ``gop install'' command.
package install

import (
	"fmt"
	"log"
	"os"
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

// gop install
var Cmd = &base.Command{
	UsageLine: "gop install [-v] [packages]",
	Short:     "Build Go+ files and install target to GOBIN",
}

var (
	flag        = &Cmd.Flag
	flagVerbose = flag.Bool("v", false, "print verbose information")
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

	gopEnv := gopenv.Get()
	conf := &gop.Config{Gop: gopEnv}
	confCmd := &gocmd.Config{Gop: gopEnv}
	confCmd.Flags = pass.Args
	for _, proj := range projs {
		install(proj, conf, confCmd)
	}
}

func install(proj gopprojs.Proj, conf *gop.Config, install *gocmd.InstallConfig) {
	var obj string
	var err error
	switch v := proj.(type) {
	case *gopprojs.DirProj:
		obj = v.Dir
		err = gop.InstallDir(obj, conf, install)
	case *gopprojs.PkgPathProj:
		obj = v.Path
		err = gop.InstallPkgPath("", v.Path, conf, install)
	case *gopprojs.FilesProj:
		err = gop.InstallFiles(v.Files, conf, install)
	default:
		log.Panicln("`gop install` doesn't support", reflect.TypeOf(v))
	}
	if err == syscall.ENOENT {
		fmt.Fprintf(os.Stderr, "gop install %v: not found\n", obj)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		return
	}
	os.Exit(1)
}

// -----------------------------------------------------------------------------
