/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package main

import (
	"fmt"
	"os"

	"github.com/qiniu/qlang/v6/cl"
	"github.com/qiniu/qlang/v6/exec"
	"github.com/qiniu/qlang/v6/parser"
	"github.com/qiniu/qlang/v6/token"
	"github.com/qiniu/x/log"

	_ "github.com/qiniu/qlang/v6/lib/builtin"
	_ "github.com/qiniu/qlang/v6/lib/fmt"
	_ "github.com/qiniu/qlang/v6/lib/reflect"
	_ "github.com/qiniu/qlang/v6/lib/strings"
)

// -----------------------------------------------------------------------------

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: qrun <qlangSrcDir>")
		return
	}
	log.SetFlags(log.Ldefault &^ log.LstdFlags)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, os.Args[1], nil, 0)
	if err != nil {
		log.Fatalln("ParseDir failed:", err)
	}
	cl.CallBuiltinOp = exec.CallBuiltinOp

	b := exec.NewBuilder(nil)
	_, err = cl.NewPackage(b.Interface(), pkgs["main"], fset)
	if err != nil {
		log.Fatalln("cl.NewPackage failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
}

// -----------------------------------------------------------------------------
