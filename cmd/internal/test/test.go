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

// Package test implements the “gop test” command.
package test

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/goplus/gogen"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/tool"
	"github.com/goplus/gop/x/gocmd"
	"github.com/goplus/gop/x/gopprojs"
)

// gop test
var Cmd = &base.Command{
	UsageLine: "gop test [-debug] [packages]",
	Short:     "Test Go+ packages",
}

var (
	flag      = &Cmd.Flag
	flagDebug = flag.Bool("debug", false, "print debug information")
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	pass := PassTestFlags(cmd)
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

	if *flagDebug {
		gogen.SetDebug(gogen.DbgFlagAll &^ gogen.DbgFlagComments)
		cl.SetDebug(cl.DbgFlagAll)
		cl.SetDisableRecover(true)
	}

	conf, err := tool.NewDefaultConf(".", 0, pass.Tags())
	if err != nil {
		log.Panicln("tool.NewDefaultConf:", err)
	}
	defer conf.UpdateCache()

	confCmd := conf.NewGoCmdConf()
	confCmd.Flags = pass.Args
	for _, proj := range projs {
		test(proj, conf, confCmd)
	}
}

func test(proj gopprojs.Proj, conf *tool.Config, test *gocmd.TestConfig) {
	const flags = tool.GenFlagPrompt
	var obj string
	var err error
	switch v := proj.(type) {
	case *gopprojs.DirProj:
		obj = v.Dir
		err = tool.TestDir(obj, conf, test, flags)
	case *gopprojs.PkgPathProj:
		obj = v.Path
		err = tool.TestPkgPath("", v.Path, conf, test, flags)
	case *gopprojs.FilesProj:
		err = tool.TestFiles(v.Files, conf, test)
	default:
		log.Panicln("`gop test` doesn't support", reflect.TypeOf(v))
	}
	if tool.NotFound(err) {
		fmt.Fprintf(os.Stderr, "gop test %v: not found\n", obj)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		return
	}
	os.Exit(1)
}

// -----------------------------------------------------------------------------
