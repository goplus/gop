// Package crypto provide Go+ "crypto" package, as "crypto" package in Go.
package crypto

import (
	crypto "crypto"
	hash "hash"
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

func execiDecrypterDecrypt(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(crypto.Decrypter).Decrypt(toType0(args[1]), args[2].([]byte), args[3])
	p.Ret(4, ret0, ret1)
}

func execiDecrypterPublic(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(crypto.Decrypter).Public()
	p.Ret(1, ret0)
}

func execmHashHashFunc(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(crypto.Hash).HashFunc()
	p.Ret(1, ret0)
}

func execmHashSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(crypto.Hash).Size()
	p.Ret(1, ret0)
}

func execmHashNew(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(crypto.Hash).New()
	p.Ret(1, ret0)
}

func execmHashAvailable(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(crypto.Hash).Available()
	p.Ret(1, ret0)
}

func execRegisterHash(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	crypto.RegisterHash(args[0].(crypto.Hash), args[1].(func() hash.Hash))
	p.PopN(2)
}

func execiSignerPublic(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(crypto.Signer).Public()
	p.Ret(1, ret0)
}

func toType1(v interface{}) crypto.SignerOpts {
	if v == nil {
		return nil
	}
	return v.(crypto.SignerOpts)
}

func execiSignerSign(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(crypto.Signer).Sign(toType0(args[1]), args[2].([]byte), toType1(args[3]))
	p.Ret(4, ret0, ret1)
}

func execiSignerOptsHashFunc(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(crypto.SignerOpts).HashFunc()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto")

func init() {
	I.RegisterFuncs(
		I.Func("(Decrypter).Decrypt", (crypto.Decrypter).Decrypt, execiDecrypterDecrypt),
		I.Func("(Decrypter).Public", (crypto.Decrypter).Public, execiDecrypterPublic),
		I.Func("(Hash).HashFunc", (crypto.Hash).HashFunc, execmHashHashFunc),
		I.Func("(Hash).Size", (crypto.Hash).Size, execmHashSize),
		I.Func("(Hash).New", (crypto.Hash).New, execmHashNew),
		I.Func("(Hash).Available", (crypto.Hash).Available, execmHashAvailable),
		I.Func("RegisterHash", crypto.RegisterHash, execRegisterHash),
		I.Func("(Signer).Public", (crypto.Signer).Public, execiSignerPublic),
		I.Func("(Signer).Sign", (crypto.Signer).Sign, execiSignerSign),
		I.Func("(SignerOpts).HashFunc", (crypto.SignerOpts).HashFunc, execiSignerOptsHashFunc),
	)
	I.RegisterTypes(
		I.Type("Decrypter", reflect.TypeOf((*crypto.Decrypter)(nil)).Elem()),
		I.Type("DecrypterOpts", reflect.TypeOf((*crypto.DecrypterOpts)(nil)).Elem()),
		I.Type("Hash", reflect.TypeOf((*crypto.Hash)(nil)).Elem()),
		I.Type("PrivateKey", reflect.TypeOf((*crypto.PrivateKey)(nil)).Elem()),
		I.Type("PublicKey", reflect.TypeOf((*crypto.PublicKey)(nil)).Elem()),
		I.Type("Signer", reflect.TypeOf((*crypto.Signer)(nil)).Elem()),
		I.Type("SignerOpts", reflect.TypeOf((*crypto.SignerOpts)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("BLAKE2b_256", qspec.Uint, crypto.BLAKE2b_256),
		I.Const("BLAKE2b_384", qspec.Uint, crypto.BLAKE2b_384),
		I.Const("BLAKE2b_512", qspec.Uint, crypto.BLAKE2b_512),
		I.Const("BLAKE2s_256", qspec.Uint, crypto.BLAKE2s_256),
		I.Const("MD4", qspec.Uint, crypto.MD4),
		I.Const("MD5", qspec.Uint, crypto.MD5),
		I.Const("MD5SHA1", qspec.Uint, crypto.MD5SHA1),
		I.Const("RIPEMD160", qspec.Uint, crypto.RIPEMD160),
		I.Const("SHA1", qspec.Uint, crypto.SHA1),
		I.Const("SHA224", qspec.Uint, crypto.SHA224),
		I.Const("SHA256", qspec.Uint, crypto.SHA256),
		I.Const("SHA384", qspec.Uint, crypto.SHA384),
		I.Const("SHA3_224", qspec.Uint, crypto.SHA3_224),
		I.Const("SHA3_256", qspec.Uint, crypto.SHA3_256),
		I.Const("SHA3_384", qspec.Uint, crypto.SHA3_384),
		I.Const("SHA3_512", qspec.Uint, crypto.SHA3_512),
		I.Const("SHA512", qspec.Uint, crypto.SHA512),
		I.Const("SHA512_224", qspec.Uint, crypto.SHA512_224),
		I.Const("SHA512_256", qspec.Uint, crypto.SHA512_256),
	)
}
