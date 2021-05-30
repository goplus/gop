// Package ecdsa provide Go+ "crypto/ecdsa" package, as "crypto/ecdsa" package in Go.
package ecdsa

import (
	crypto "crypto"
	ecdsa "crypto/ecdsa"
	elliptic "crypto/elliptic"
	io "io"
	big "math/big"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func toType0(v interface{}) elliptic.Curve {
	if v == nil {
		return nil
	}
	return v.(elliptic.Curve)
}

func toType1(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execGenerateKey(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := ecdsa.GenerateKey(toType0(args[0]), toType1(args[1]))
	p.Ret(2, ret0, ret1)
}

func execmPrivateKeyPublic(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*ecdsa.PrivateKey).Public()
	p.Ret(1, ret0)
}

func toType2(v interface{}) crypto.SignerOpts {
	if v == nil {
		return nil
	}
	return v.(crypto.SignerOpts)
}

func execmPrivateKeySign(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*ecdsa.PrivateKey).Sign(toType1(args[1]), args[2].([]byte), toType2(args[3]))
	p.Ret(4, ret0, ret1)
}

func execSign(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1, ret2 := ecdsa.Sign(toType1(args[0]), args[1].(*ecdsa.PrivateKey), args[2].([]byte))
	p.Ret(3, ret0, ret1, ret2)
}

func execVerify(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := ecdsa.Verify(args[0].(*ecdsa.PublicKey), args[1].([]byte), args[2].(*big.Int), args[3].(*big.Int))
	p.Ret(4, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/ecdsa")

func init() {
	I.RegisterFuncs(
		I.Func("GenerateKey", ecdsa.GenerateKey, execGenerateKey),
		I.Func("(*PrivateKey).Public", (*ecdsa.PrivateKey).Public, execmPrivateKeyPublic),
		I.Func("(*PrivateKey).Sign", (*ecdsa.PrivateKey).Sign, execmPrivateKeySign),
		I.Func("Sign", ecdsa.Sign, execSign),
		I.Func("Verify", ecdsa.Verify, execVerify),
	)
	I.RegisterTypes(
		I.Type("PrivateKey", reflect.TypeOf((*ecdsa.PrivateKey)(nil)).Elem()),
		I.Type("PublicKey", reflect.TypeOf((*ecdsa.PublicKey)(nil)).Elem()),
	)
}
