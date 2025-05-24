package main

import (
	"fmt"
	"os"

	"github.com/qiniu/x/osx"
)

func main() {
	for _gop_it := osx.EnumLines(os.Stdin); ; {
		var _gop_ok bool
		line, _gop_ok := _gop_it.Next()
		if !_gop_ok {
			break
		}
		fmt.Println(line)
	}
}
