// Package sha512 provide Go+ "crypto/sha512" package, as "crypto/sha512" package in Go.
package sha512

import (
	sha512 "crypto/sha512"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execNew(_ int, p *gop.Context) {
	ret0 := sha512.New()
	p.Ret(0, ret0)
}

func execNew384(_ int, p *gop.Context) {
	ret0 := sha512.New384()
	p.Ret(0, ret0)
}

func execNew512_224(_ int, p *gop.Context) {
	ret0 := sha512.New512_224()
	p.Ret(0, ret0)
}

func execNew512_256(_ int, p *gop.Context) {
	ret0 := sha512.New512_256()
	p.Ret(0, ret0)
}

func execSum384(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := sha512.Sum384(args[0].([]byte))
	p.Ret(1, ret0)
}

func execSum512(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := sha512.Sum512(args[0].([]byte))
	p.Ret(1, ret0)
}

func execSum512_224(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := sha512.Sum512_224(args[0].([]byte))
	p.Ret(1, ret0)
}

func execSum512_256(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := sha512.Sum512_256(args[0].([]byte))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/sha512")

func init() {
	I.RegisterFuncs(
		I.Func("New", sha512.New, execNew),
		I.Func("New384", sha512.New384, execNew384),
		I.Func("New512_224", sha512.New512_224, execNew512_224),
		I.Func("New512_256", sha512.New512_256, execNew512_256),
		I.Func("Sum384", sha512.Sum384, execSum384),
		I.Func("Sum512", sha512.Sum512, execSum512),
		I.Func("Sum512_224", sha512.Sum512_224, execSum512_224),
		I.Func("Sum512_256", sha512.Sum512_256, execSum512_256),
	)
	I.RegisterConsts(
		I.Const("BlockSize", qspec.ConstUnboundInt, sha512.BlockSize),
		I.Const("Size", qspec.ConstUnboundInt, sha512.Size),
		I.Const("Size224", qspec.ConstUnboundInt, sha512.Size224),
		I.Const("Size256", qspec.ConstUnboundInt, sha512.Size256),
		I.Const("Size384", qspec.ConstUnboundInt, sha512.Size384),
	)
}
