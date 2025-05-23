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

package list

/*
import (
	"fmt"
	"log"

	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/x/gopmod"
)

// -----------------------------------------------------------------------------

// gop list
var Cmd = &base.Command{
	UsageLine: "gop list [-json] [packages]",
	Short:     "List packages or modules",
}

var (
	flag = &Cmd.Flag
	_    = flag.Bool("json", false, "printing in JSON format.")
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

	pkgPaths, err := gopmod.List(pattern...)
	check(err)
	for _, pkgPath := range pkgPaths {
		fmt.Println(pkgPath)
	}
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
*/

// -----------------------------------------------------------------------------
