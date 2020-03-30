package main

import (
	"fmt"
	"os"

	"github.com/qiniu/qlang"
	_ "github.com/qiniu/qlang/lib/builtin" // 导入 builtin 包
)

// -----------------------------------------------------------------------------

const scriptCode = `
	x = 1 + 2
`

func main() {

	lang := qlang.New()
	err := lang.SafeExec([]byte(scriptCode), "")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	v, _ := lang.Var("x")
	fmt.Println("x:", v)
}

// -----------------------------------------------------------------------------
