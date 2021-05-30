// Package hmac provide Go+ "crypto/hmac" package, as "crypto/hmac" package in Go.
package hmac

import (
	hmac "crypto/hmac"
	hash "hash"

	gop "github.com/goplus/gop"
)

func execEqual(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := hmac.Equal(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execNew(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := hmac.New(args[0].(func() hash.Hash), args[1].([]byte))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/hmac")

func init() {
	I.RegisterFuncs(
		I.Func("Equal", hmac.Equal, execEqual),
		I.Func("New", hmac.New, execNew),
	)
}
