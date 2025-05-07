package main

import (
	"fmt"
	"github.com/qiniu/x/gop/osx"
	"io"
)

var r io.Reader

func main() {
	for _gop_it := osx.Lines(r).Gop_Enum(); ; {
		var _gop_ok bool
		line, _gop_ok := _gop_it.Next()
		if !_gop_ok {
			break
		}
		fmt.Println(line)
	}
}
