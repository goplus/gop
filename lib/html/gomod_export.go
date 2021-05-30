// Package html provide Go+ "html" package, as "html" package in Go.
package html

import (
	html "html"

	gop "github.com/goplus/gop"
)

func execEscapeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := html.EscapeString(args[0].(string))
	p.Ret(1, ret0)
}

func execUnescapeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := html.UnescapeString(args[0].(string))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("html")

func init() {
	I.RegisterFuncs(
		I.Func("EscapeString", html.EscapeString, execEscapeString),
		I.Func("UnescapeString", html.UnescapeString, execUnescapeString),
	)
}
