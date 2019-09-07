package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/qiniu/qlang"
	_ "github.com/qiniu/qlang/lib/builtin" // 导入 builtin 包
)

// -----------------------------------------------------------------------------

var strings_Exports = map[string]interface{}{
	"replacer": strings.NewReplacer,
}

func main() {

	qlang.Import("strings",	strings_Exports) // 导入一个自定义的包，叫 strings（和标准库同名）
	ql := qlang.New()

	err := ql.SafeEval(`x = strings.replacer("?", "!").replace("hello, world???")`)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("x:", ql.Var("x")) // 输出 x: hello, world!!!
}

// -----------------------------------------------------------------------------
