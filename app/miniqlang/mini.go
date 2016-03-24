package main

import (
	"fmt"
	"os"

	"qlang.io/qlang.v1/qlang"
	_ "qlang.io/qlang/builtin" // 导入 builtin 包
)

// -----------------------------------------------------------------------------

func main() {

	lang, err := qlang.New(nil) // 参数 nil 也可以改为 qlang.InsertSemis
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = lang.SafeEval(`"str" + 2`)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	v, _ := lang.Ret()
	fmt.Println(v)
}

// -----------------------------------------------------------------------------

