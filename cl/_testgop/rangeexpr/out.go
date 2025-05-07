package main

import (
	"fmt"
	"github.com/qiniu/x/gop"
)

func main() {
	fmt.Println(func() (_gop_ret []int) {
		for _gop_it := gop.NewRange__0(0, 3, 1).Gop_Enum(); ; {
			var _gop_ok bool
			x, _gop_ok := _gop_it.Next()
			if !_gop_ok {
				break
			}
			_gop_ret = append(_gop_ret, x)
		}
		return
	}())
}
