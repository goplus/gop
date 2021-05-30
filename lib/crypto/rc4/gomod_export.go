// Package rc4 provide Go+ "crypto/rc4" package, as "crypto/rc4" package in Go.
package rc4

import (
	rc4 "crypto/rc4"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execmCipherReset(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*rc4.Cipher).Reset()
	p.PopN(1)
}

func execmCipherXORKeyStream(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*rc4.Cipher).XORKeyStream(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func execmKeySizeErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(rc4.KeySizeError).Error()
	p.Ret(1, ret0)
}

func execNewCipher(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := rc4.NewCipher(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/rc4")

func init() {
	I.RegisterFuncs(
		I.Func("(*Cipher).Reset", (*rc4.Cipher).Reset, execmCipherReset),
		I.Func("(*Cipher).XORKeyStream", (*rc4.Cipher).XORKeyStream, execmCipherXORKeyStream),
		I.Func("(KeySizeError).Error", (rc4.KeySizeError).Error, execmKeySizeErrorError),
		I.Func("NewCipher", rc4.NewCipher, execNewCipher),
	)
	I.RegisterTypes(
		I.Type("Cipher", reflect.TypeOf((*rc4.Cipher)(nil)).Elem()),
		I.Type("KeySizeError", reflect.TypeOf((*rc4.KeySizeError)(nil)).Elem()),
	)
}
