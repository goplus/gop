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

package gopget

import (
	"fmt"
	"log"
	"os"

	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/x/mod/modfetch"
	"github.com/goplus/gop/x/mod/modload"
	"golang.org/x/mod/module"
)

// -----------------------------------------------------------------------------

// Cmd - gop build
var Cmd = &base.Command{
	UsageLine: "gop get [-v] [packages]",
	Short:     "Version prints the build information for Gop executables",
}

var (
	flag = &Cmd.Flag
	_    = flag.Bool("v", false, "print verbose information.")
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	err := flag.Parse(args)
	if err != nil {
		log.Fatalln("parse input arguments failed:", err)
	}
	narg := flag.NArg()
	if narg < 1 {
		log.Fatalln("TODO: not impl")
	}
	for i := 0; i < narg; i++ {
		get(flag.Arg(i))
	}
}

func get(pkgPath string) {
	mod, err := modload.Load(".")
	check(err)
	check(mod.UpdateGoMod(true))
	modPath, _ := splitPkgPath(pkgPath)
	inited := false
	_, err = modfetch.Get(modPath, func(act string, modVer module.Version) {
		if !inited {
			genGo(modVer)
			check(mod.AddRequire(modVer.Path, modVer.Version))
			check(mod.Save())
			fmt.Fprintf(os.Stderr, "gop get: added %s %s\n", modVer.Path, modVer.Version)
			check(mod.UpdateGoMod(false))
			inited = true
		}
	})
	check(err)
}

func genGo(mod module.Version) {
	// TODO: not impl
}

func splitPkgPath(pkgPath string) (modPathWithVer string, pkgPathNoVer string) {
	return pkgPath, pkgPath
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// -----------------------------------------------------------------------------
