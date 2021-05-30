// Package dsa provide Go+ "crypto/dsa" package, as "crypto/dsa" package in Go.
package dsa

import (
	dsa "crypto/dsa"
	io "io"
	big "math/big"
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
	args := p.GetArgs(2)
	ret0 := dsa.GenerateKey(args[0].(*dsa.PrivateKey), toType0(args[1]))
	p.Ret(2, ret0)
}

func execGenerateParameters(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := dsa.GenerateParameters(args[0].(*dsa.Parameters), toType0(args[1]), args[2].(dsa.ParameterSizes))
	p.Ret(3, ret0)
}

func execSign(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1, ret2 := dsa.Sign(toType0(args[0]), args[1].(*dsa.PrivateKey), args[2].([]byte))
	p.Ret(3, ret0, ret1, ret2)
}

func execVerify(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := dsa.Verify(args[0].(*dsa.PublicKey), args[1].([]byte), args[2].(*big.Int), args[3].(*big.Int))
	p.Ret(4, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/dsa")

func init() {
	I.RegisterFuncs(
		I.Func("GenerateKey", dsa.GenerateKey, execGenerateKey),
		I.Func("GenerateParameters", dsa.GenerateParameters, execGenerateParameters),
		I.Func("Sign", dsa.Sign, execSign),
		I.Func("Verify", dsa.Verify, execVerify),
	)
	I.RegisterVars(
		I.Var("ErrInvalidPublicKey", &dsa.ErrInvalidPublicKey),
	)
	I.RegisterTypes(
		I.Type("ParameterSizes", reflect.TypeOf((*dsa.ParameterSizes)(nil)).Elem()),
		I.Type("Parameters", reflect.TypeOf((*dsa.Parameters)(nil)).Elem()),
		I.Type("PrivateKey", reflect.TypeOf((*dsa.PrivateKey)(nil)).Elem()),
		I.Type("PublicKey", reflect.TypeOf((*dsa.PublicKey)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("L1024N160", qspec.Int, dsa.L1024N160),
		I.Const("L2048N224", qspec.Int, dsa.L2048N224),
		I.Const("L2048N256", qspec.Int, dsa.L2048N256),
		I.Const("L3072N256", qspec.Int, dsa.L3072N256),
	)
}
