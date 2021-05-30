// Package md5 provide Go+ "crypto/md5" package, as "crypto/md5" package in Go.
package md5

import (
	md5 "crypto/md5"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execNew(_ int, p *gop.Context) {
	ret0 := md5.New()
	p.Ret(0, ret0)
}

func execSum(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := md5.Sum(args[0].([]byte))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/md5")

func init() {
	I.RegisterFuncs(
		I.Func("New", md5.New, execNew),
		I.Func("Sum", md5.Sum, execSum),
	)
	I.RegisterConsts(
		I.Const("BlockSize", qspec.ConstUnboundInt, md5.BlockSize),
		I.Const("Size", qspec.ConstUnboundInt, md5.Size),
	)
}
