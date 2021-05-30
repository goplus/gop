// Package sha1 provide Go+ "crypto/sha1" package, as "crypto/sha1" package in Go.
package sha1

import (
	sha1 "crypto/sha1"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execNew(_ int, p *gop.Context) {
	ret0 := sha1.New()
	p.Ret(0, ret0)
}

func execSum(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := sha1.Sum(args[0].([]byte))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/sha1")

func init() {
	I.RegisterFuncs(
		I.Func("New", sha1.New, execNew),
		I.Func("Sum", sha1.Sum, execSum),
	)
	I.RegisterConsts(
		I.Const("BlockSize", qspec.ConstUnboundInt, sha1.BlockSize),
		I.Const("Size", qspec.ConstUnboundInt, sha1.Size),
	)
}
