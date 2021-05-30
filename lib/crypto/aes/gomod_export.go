// Package aes provide Go+ "crypto/aes" package, as "crypto/aes" package in Go.
package aes

import (
	aes "crypto/aes"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmKeySizeErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(aes.KeySizeError).Error()
	p.Ret(1, ret0)
}

func execNewCipher(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := aes.NewCipher(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/aes")

func init() {
	I.RegisterFuncs(
		I.Func("(KeySizeError).Error", (aes.KeySizeError).Error, execmKeySizeErrorError),
		I.Func("NewCipher", aes.NewCipher, execNewCipher),
	)
	I.RegisterTypes(
		I.Type("KeySizeError", reflect.TypeOf((*aes.KeySizeError)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("BlockSize", qspec.ConstUnboundInt, aes.BlockSize),
	)
}
