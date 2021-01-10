// Package goptest provide Go+ "github.com/goplus/gop/ast/goptest" package, as "github.com/goplus/gop/ast/goptest" package in Go.
package goptest

import (
	gop "github.com/goplus/gop"
	goptest "github.com/goplus/gop/ast/goptest"
)

func execNew(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := goptest.New(args[0].(string))
	p.Ret(1, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/goplus/gop/ast/goptest")

func init() {
	I.RegisterFuncs(
		I.Func("New", goptest.New, execNew),
	)
}
