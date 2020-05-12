package main

import (
	"fmt"
	"os"

	"github.com/qiniu/qlang/cl"
	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/qlang/parser"
	"github.com/qiniu/qlang/token"
	"github.com/qiniu/x/log"

	_ "github.com/qiniu/qlang/lib/builtin"
	_ "github.com/qiniu/qlang/lib/fmt"
	_ "github.com/qiniu/qlang/lib/strings"
)

// -----------------------------------------------------------------------------

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: qrun <qlangSrcDir>")
		return
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, os.Args[1], nil, 0)
	if err != nil {
		log.Fatalln("ParseDir failed:", err)
	}

	b := exec.NewBuilder(nil)
	_, err = cl.NewPackage(b, pkgs["main"])
	if err != nil {
		log.Fatalln("cl.NewPackage failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
}

// -----------------------------------------------------------------------------
