// Package rsa provide Go+ "crypto/rsa" package, as "crypto/rsa" package in Go.
package rsa

import (
	crypto "crypto"
	rsa "crypto/rsa"
	hash "hash"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func toType0(v interface{}) hash.Hash {
	if v == nil {
		return nil
	}
	return v.(hash.Hash)
}

func toType1(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execDecryptOAEP(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0, ret1 := rsa.DecryptOAEP(toType0(args[0]), toType1(args[1]), args[2].(*rsa.PrivateKey), args[3].([]byte), args[4].([]byte))
	p.Ret(5, ret0, ret1)
}

func execDecryptPKCS1v15(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := rsa.DecryptPKCS1v15(toType1(args[0]), args[1].(*rsa.PrivateKey), args[2].([]byte))
	p.Ret(3, ret0, ret1)
}

func execDecryptPKCS1v15SessionKey(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := rsa.DecryptPKCS1v15SessionKey(toType1(args[0]), args[1].(*rsa.PrivateKey), args[2].([]byte), args[3].([]byte))
	p.Ret(4, ret0)
}

func execEncryptOAEP(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0, ret1 := rsa.EncryptOAEP(toType0(args[0]), toType1(args[1]), args[2].(*rsa.PublicKey), args[3].([]byte), args[4].([]byte))
	p.Ret(5, ret0, ret1)
}

func execEncryptPKCS1v15(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := rsa.EncryptPKCS1v15(toType1(args[0]), args[1].(*rsa.PublicKey), args[2].([]byte))
	p.Ret(3, ret0, ret1)
}

func execGenerateKey(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := rsa.GenerateKey(toType1(args[0]), args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execGenerateMultiPrimeKey(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := rsa.GenerateMultiPrimeKey(toType1(args[0]), args[1].(int), args[2].(int))
	p.Ret(3, ret0, ret1)
}

func execmPSSOptionsHashFunc(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rsa.PSSOptions).HashFunc()
	p.Ret(1, ret0)
}

func execmPrivateKeyPublic(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rsa.PrivateKey).Public()
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
	ret0, ret1 := args[0].(*rsa.PrivateKey).Sign(toType1(args[1]), args[2].([]byte), toType2(args[3]))
	p.Ret(4, ret0, ret1)
}

func execmPrivateKeyDecrypt(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*rsa.PrivateKey).Decrypt(toType1(args[1]), args[2].([]byte), args[3])
	p.Ret(4, ret0, ret1)
}

func execmPrivateKeyValidate(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rsa.PrivateKey).Validate()
	p.Ret(1, ret0)
}

func execmPrivateKeyPrecompute(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*rsa.PrivateKey).Precompute()
	p.PopN(1)
}

func execmPublicKeySize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rsa.PublicKey).Size()
	p.Ret(1, ret0)
}

func execSignPKCS1v15(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := rsa.SignPKCS1v15(toType1(args[0]), args[1].(*rsa.PrivateKey), args[2].(crypto.Hash), args[3].([]byte))
	p.Ret(4, ret0, ret1)
}

func execSignPSS(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0, ret1 := rsa.SignPSS(toType1(args[0]), args[1].(*rsa.PrivateKey), args[2].(crypto.Hash), args[3].([]byte), args[4].(*rsa.PSSOptions))
	p.Ret(5, ret0, ret1)
}

func execVerifyPKCS1v15(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := rsa.VerifyPKCS1v15(args[0].(*rsa.PublicKey), args[1].(crypto.Hash), args[2].([]byte), args[3].([]byte))
	p.Ret(4, ret0)
}

func execVerifyPSS(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0 := rsa.VerifyPSS(args[0].(*rsa.PublicKey), args[1].(crypto.Hash), args[2].([]byte), args[3].([]byte), args[4].(*rsa.PSSOptions))
	p.Ret(5, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/rsa")

func init() {
	I.RegisterFuncs(
		I.Func("DecryptOAEP", rsa.DecryptOAEP, execDecryptOAEP),
		I.Func("DecryptPKCS1v15", rsa.DecryptPKCS1v15, execDecryptPKCS1v15),
		I.Func("DecryptPKCS1v15SessionKey", rsa.DecryptPKCS1v15SessionKey, execDecryptPKCS1v15SessionKey),
		I.Func("EncryptOAEP", rsa.EncryptOAEP, execEncryptOAEP),
		I.Func("EncryptPKCS1v15", rsa.EncryptPKCS1v15, execEncryptPKCS1v15),
		I.Func("GenerateKey", rsa.GenerateKey, execGenerateKey),
		I.Func("GenerateMultiPrimeKey", rsa.GenerateMultiPrimeKey, execGenerateMultiPrimeKey),
		I.Func("(*PSSOptions).HashFunc", (*rsa.PSSOptions).HashFunc, execmPSSOptionsHashFunc),
		I.Func("(*PrivateKey).Public", (*rsa.PrivateKey).Public, execmPrivateKeyPublic),
		I.Func("(*PrivateKey).Sign", (*rsa.PrivateKey).Sign, execmPrivateKeySign),
		I.Func("(*PrivateKey).Decrypt", (*rsa.PrivateKey).Decrypt, execmPrivateKeyDecrypt),
		I.Func("(*PrivateKey).Validate", (*rsa.PrivateKey).Validate, execmPrivateKeyValidate),
		I.Func("(*PrivateKey).Precompute", (*rsa.PrivateKey).Precompute, execmPrivateKeyPrecompute),
		I.Func("(*PublicKey).Size", (*rsa.PublicKey).Size, execmPublicKeySize),
		I.Func("SignPKCS1v15", rsa.SignPKCS1v15, execSignPKCS1v15),
		I.Func("SignPSS", rsa.SignPSS, execSignPSS),
		I.Func("VerifyPKCS1v15", rsa.VerifyPKCS1v15, execVerifyPKCS1v15),
		I.Func("VerifyPSS", rsa.VerifyPSS, execVerifyPSS),
	)
	I.RegisterVars(
		I.Var("ErrDecryption", &rsa.ErrDecryption),
		I.Var("ErrMessageTooLong", &rsa.ErrMessageTooLong),
		I.Var("ErrVerification", &rsa.ErrVerification),
	)
	I.RegisterTypes(
		I.Type("CRTValue", reflect.TypeOf((*rsa.CRTValue)(nil)).Elem()),
		I.Type("OAEPOptions", reflect.TypeOf((*rsa.OAEPOptions)(nil)).Elem()),
		I.Type("PKCS1v15DecryptOptions", reflect.TypeOf((*rsa.PKCS1v15DecryptOptions)(nil)).Elem()),
		I.Type("PSSOptions", reflect.TypeOf((*rsa.PSSOptions)(nil)).Elem()),
		I.Type("PrecomputedValues", reflect.TypeOf((*rsa.PrecomputedValues)(nil)).Elem()),
		I.Type("PrivateKey", reflect.TypeOf((*rsa.PrivateKey)(nil)).Elem()),
		I.Type("PublicKey", reflect.TypeOf((*rsa.PublicKey)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("PSSSaltLengthAuto", qspec.ConstUnboundInt, rsa.PSSSaltLengthAuto),
		I.Const("PSSSaltLengthEqualsHash", qspec.ConstUnboundInt, rsa.PSSSaltLengthEqualsHash),
	)
}
