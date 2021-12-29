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

package deps

import (
	"fmt"
	"log"
	"sort"

	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/x/gopmod"
)

// -----------------------------------------------------------------------------

// Cmd - gop deps
var Cmd = &base.Command{
	UsageLine: "gop deps [-r -v] [package]",
	Short:     "Show dependencies of a package or module",
}

var (
	flag      = &Cmd.Flag
	recursive = flag.Bool("r", false, "get dependencies recursively.")
	_         = flag.Bool("v", false, "print verbose information.")
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	err := flag.Parse(args)
	if err != nil {
		log.Fatalln("parse input arguments failed:", err)
	}
	var pkgPath string
	narg := flag.NArg()
	if narg < 1 {
		pkgPath = "."
	} else {
		pkgPath = flag.Arg(0)
	}
	getDeps(pkgPath, *recursive)
}

func getDeps(pkgPath string, recursive bool) {
	mod, err := gopmod.Load(pkgPath)
	check(err)
	check(mod.RegisterClasses())
	imports, err := mod.Imports(pkgPath, recursive)
	check(err)
	sort.Strings(imports)
	for _, imp := range imports {
		fmt.Println(imp)
	}
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// -----------------------------------------------------------------------------
