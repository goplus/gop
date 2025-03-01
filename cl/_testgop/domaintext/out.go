package main

import "github.com/goplus/gop/tpl"

func main() {
	cl, err := tpl.New(`expr = INT % ("+" | "-")`)
	cl.ParseExpr("1+2", nil)
	_ = err
}
