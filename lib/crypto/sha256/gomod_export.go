// Package sha256 provide Go+ "crypto/sha256" package, as "crypto/sha256" package in Go.
package sha256

import (
	sha256 "crypto/sha256"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execNew(_ int, p *gop.Context) {
	ret0 := sha256.New()
	p.Ret(0, ret0)
}

func execNew224(_ int, p *gop.Context) {
	ret0 := sha256.New224()
	p.Ret(0, ret0)
}

func execSum224(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := sha256.Sum224(args[0].([]byte))
	p.Ret(1, ret0)
}

func execSum256(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := sha256.Sum256(args[0].([]byte))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/sha256")

func init() {
	I.RegisterFuncs(
		I.Func("New", sha256.New, execNew),
		I.Func("New224", sha256.New224, execNew224),
		I.Func("Sum224", sha256.Sum224, execSum224),
		I.Func("Sum256", sha256.Sum256, execSum256),
	)
	I.RegisterConsts(
		I.Const("BlockSize", qspec.ConstUnboundInt, sha256.BlockSize),
		I.Const("Size", qspec.ConstUnboundInt, sha256.Size),
		I.Const("Size224", qspec.ConstUnboundInt, sha256.Size224),
	)
}
