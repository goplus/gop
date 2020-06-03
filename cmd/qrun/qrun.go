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
	"flag"
	"fmt"
	"os"

	"github.com/qiniu/goplus/cl"
	"github.com/qiniu/goplus/parser"
	"github.com/qiniu/goplus/token"
	"github.com/qiniu/x/log"

	exec "github.com/qiniu/goplus/exec/bytecode"
	_ "github.com/qiniu/goplus/lib/builtin"
	_ "github.com/qiniu/goplus/lib/fmt"
	_ "github.com/qiniu/goplus/lib/reflect"
	_ "github.com/qiniu/goplus/lib/strings"
)

// -----------------------------------------------------------------------------

var (
	flagAsm   = flag.Bool("asm", false, "generate asm code")
	flagQuiet = flag.Bool("quiet", false, "don't generate any log")
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: qrun [-asm] [-quiet] <gopSrcDir>")
		flag.PrintDefaults()
		return
	}

	log.SetFlags(log.Ldefault &^ log.LstdFlags)
	if *flagQuiet {
		log.SetOutputLevel(0x7000)
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, flag.Arg(0), nil, 0)
	if err != nil {
		log.Fatalln("ParseDir failed:", err)
	}
	cl.CallBuiltinOp = exec.CallBuiltinOp

	b := exec.NewBuilder(nil)
	_, err = cl.NewPackage(b.Interface(), pkgs["main"], fset, cl.PkgActClMain)
	if err != nil {
		log.Fatalln("cl.NewPackage failed:", err)
	}
	code := b.Resolve()
	if *flagAsm {
		code.Dump(os.Stdout)
		return
	}
	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
}

// -----------------------------------------------------------------------------
