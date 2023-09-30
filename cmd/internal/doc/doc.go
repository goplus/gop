/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
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

package doc

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"syscall"

	"github.com/goplus/gop"
	"github.com/goplus/gop/cl/outline"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/x/gopenv"
	"github.com/goplus/gop/x/gopprojs"
)

// -----------------------------------------------------------------------------

// gop doc
var Cmd = &base.Command{
	UsageLine: "gop doc [-all] [pkgPath]",
	Short:     "Show documentation for package or symbol",
}

var (
	flag = &Cmd.Flag
	all  = flag.Bool("all", false, "show all the documentation for the package.")
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

	proj, _, err := gopprojs.ParseOne(pattern...)
	if err != nil {
		log.Panicln("gopprojs.ParseOne:", err)
	}

	gopEnv := gopenv.Get()
	conf := &gop.Config{Gop: gopEnv}
	outlinePkg(proj, conf)
}

func outlinePkg(proj gopprojs.Proj, conf *gop.Config) {
	var obj string
	var out outline.Package
	var err error
	switch v := proj.(type) {
	case *gopprojs.DirProj:
		obj = v.Dir
		out, err = gop.Outline(obj, conf)
	case *gopprojs.PkgPathProj:
		obj = v.Path
		out, err = gop.OutlinePkgPath("", v.Path, conf, true)
	default:
		log.Panicln("`gop doc` doesn't support", reflect.TypeOf(v))
	}
	if err == syscall.ENOENT {
		fmt.Fprintf(os.Stderr, "gop doc %v: not found\n", obj)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		_ = out.Outline(*all)
		return
	}
}

// -----------------------------------------------------------------------------
