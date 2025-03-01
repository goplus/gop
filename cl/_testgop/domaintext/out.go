package main

import (
	"github.com/goplus/gop/tpl"
	"github.com/qiniu/x/errors"
)

func main() {
	cl := func() (_gop_ret tpl.Compiler) {
		var _gop_err error
		_gop_ret, _gop_err = tpl.New(`expr = INT % ("+" | "-")`)
		if _gop_err != nil {
			_gop_err = errors.NewFrame(_gop_err, "tpl`expr = INT % (\"+\" | \"-\")`", "/Users/xushiwei/work/gop/cl/_testgop/domaintext/in.gop", 1, "main.main")
			panic(_gop_err)
		}
		return
	}()
	cl.ParseExpr("1+2", nil)
}
