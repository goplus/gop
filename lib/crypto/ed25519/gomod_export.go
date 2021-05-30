// Package ed25519 provide Go+ "crypto/ed25519" package, as "crypto/ed25519" package in Go.
package ed25519

import (
	crypto "crypto"
	ed25519 "crypto/ed25519"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execGenerateKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := ed25519.GenerateKey(toType0(args[0]))
	p.Ret(1, ret0, ret1, ret2)
}

func execNewKeyFromSeed(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := ed25519.NewKeyFromSeed(args[0].([]byte))
	p.Ret(1, ret0)
}

func execmPrivateKeyPublic(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(ed25519.PrivateKey).Public()
	p.Ret(1, ret0)
}

func execmPrivateKeySeed(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(ed25519.PrivateKey).Seed()
	p.Ret(1, ret0)
}

func toType1(v interface{}) crypto.SignerOpts {
	if v == nil {
		return nil
	}
	return v.(crypto.SignerOpts)
}

func execmPrivateKeySign(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(ed25519.PrivateKey).Sign(toType0(args[1]), args[2].([]byte), toType1(args[3]))
	p.Ret(4, ret0, ret1)
}

func execSign(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := ed25519.Sign(args[0].(ed25519.PrivateKey), args[1].([]byte))
	p.Ret(2, ret0)
}

func execVerify(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := ed25519.Verify(args[0].(ed25519.PublicKey), args[1].([]byte), args[2].([]byte))
	p.Ret(3, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/ed25519")

func init() {
	I.RegisterFuncs(
		I.Func("GenerateKey", ed25519.GenerateKey, execGenerateKey),
		I.Func("NewKeyFromSeed", ed25519.NewKeyFromSeed, execNewKeyFromSeed),
		I.Func("(PrivateKey).Public", (ed25519.PrivateKey).Public, execmPrivateKeyPublic),
		I.Func("(PrivateKey).Seed", (ed25519.PrivateKey).Seed, execmPrivateKeySeed),
		I.Func("(PrivateKey).Sign", (ed25519.PrivateKey).Sign, execmPrivateKeySign),
		I.Func("Sign", ed25519.Sign, execSign),
		I.Func("Verify", ed25519.Verify, execVerify),
	)
	I.RegisterTypes(
		I.Type("PrivateKey", reflect.TypeOf((*ed25519.PrivateKey)(nil)).Elem()),
		I.Type("PublicKey", reflect.TypeOf((*ed25519.PublicKey)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("PrivateKeySize", qspec.ConstUnboundInt, ed25519.PrivateKeySize),
		I.Const("PublicKeySize", qspec.ConstUnboundInt, ed25519.PublicKeySize),
		I.Const("SeedSize", qspec.ConstUnboundInt, ed25519.SeedSize),
		I.Const("SignatureSize", qspec.ConstUnboundInt, ed25519.SignatureSize),
	)
}
