/*
 * Copyright (c) 2021 The XGo Authors (xgo.dev). All rights reserved.
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

/*
import (
	"fmt"
	"log"

	"github.com/goplus/xgo/cmd/internal/base"
	"github.com/goplus/xgo/x/gopmod"
)

// -----------------------------------------------------------------------------

// gop deps
var Cmd = &base.Command{
	UsageLine: "gop deps [-v] [package]",
	Short:     "Show dependencies of a package or module",
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

	var dir string
	narg := flag.NArg()
	if narg < 1 {
		dir = "."
	} else {
		dir = flag.Arg(0)
	}

	imports, err := gopmod.Imports(dir)
	check(err)
	for _, imp := range imports {
		fmt.Println(imp)
	}
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
*/

// -----------------------------------------------------------------------------
