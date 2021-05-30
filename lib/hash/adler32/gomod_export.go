// Package adler32 provide Go+ "hash/adler32" package, as "hash/adler32" package in Go.
package adler32

import (
	adler32 "hash/adler32"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execChecksum(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := adler32.Checksum(args[0].([]byte))
	p.Ret(1, ret0)
}

func execNew(_ int, p *gop.Context) {
	ret0 := adler32.New()
	p.Ret(0, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("hash/adler32")

func init() {
	I.RegisterFuncs(
		I.Func("Checksum", adler32.Checksum, execChecksum),
		I.Func("New", adler32.New, execNew),
	)
	I.RegisterConsts(
		I.Const("Size", qspec.ConstUnboundInt, adler32.Size),
	)
}
