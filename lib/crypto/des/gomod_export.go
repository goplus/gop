// Package des provide Go+ "crypto/des" package, as "crypto/des" package in Go.
package des

import (
	des "crypto/des"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmKeySizeErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(des.KeySizeError).Error()
	p.Ret(1, ret0)
}

func execNewCipher(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := des.NewCipher(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execNewTripleDESCipher(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := des.NewTripleDESCipher(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/des")

func init() {
	I.RegisterFuncs(
		I.Func("(KeySizeError).Error", (des.KeySizeError).Error, execmKeySizeErrorError),
		I.Func("NewCipher", des.NewCipher, execNewCipher),
		I.Func("NewTripleDESCipher", des.NewTripleDESCipher, execNewTripleDESCipher),
	)
	I.RegisterTypes(
		I.Type("KeySizeError", reflect.TypeOf((*des.KeySizeError)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("BlockSize", qspec.ConstUnboundInt, des.BlockSize),
	)
}
