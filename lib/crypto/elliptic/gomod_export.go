// Package elliptic provide Go+ "crypto/elliptic" package, as "crypto/elliptic" package in Go.
package elliptic

import (
	elliptic "crypto/elliptic"
	io "io"
	big "math/big"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execiCurveAdd(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0, ret1 := args[0].(elliptic.Curve).Add(args[1].(*big.Int), args[2].(*big.Int), args[3].(*big.Int), args[4].(*big.Int))
	p.Ret(5, ret0, ret1)
}

func execiCurveDouble(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(elliptic.Curve).Double(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0, ret1)
}

func execiCurveIsOnCurve(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(elliptic.Curve).IsOnCurve(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execiCurveParams(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(elliptic.Curve).Params()
	p.Ret(1, ret0)
}

func execiCurveScalarBaseMult(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(elliptic.Curve).ScalarBaseMult(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiCurveScalarMult(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(elliptic.Curve).ScalarMult(args[1].(*big.Int), args[2].(*big.Int), args[3].([]byte))
	p.Ret(4, ret0, ret1)
}

func execmCurveParamsParams(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*elliptic.CurveParams).Params()
	p.Ret(1, ret0)
}

func execmCurveParamsIsOnCurve(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*elliptic.CurveParams).IsOnCurve(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmCurveParamsAdd(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0, ret1 := args[0].(*elliptic.CurveParams).Add(args[1].(*big.Int), args[2].(*big.Int), args[3].(*big.Int), args[4].(*big.Int))
	p.Ret(5, ret0, ret1)
}

func execmCurveParamsDouble(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*elliptic.CurveParams).Double(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0, ret1)
}

func execmCurveParamsScalarMult(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*elliptic.CurveParams).ScalarMult(args[1].(*big.Int), args[2].(*big.Int), args[3].([]byte))
	p.Ret(4, ret0, ret1)
}

func execmCurveParamsScalarBaseMult(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*elliptic.CurveParams).ScalarBaseMult(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

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
	ret0, ret1, ret2, ret3 := elliptic.GenerateKey(toType0(args[0]), toType1(args[1]))
	p.Ret(2, ret0, ret1, ret2, ret3)
}

func execMarshal(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := elliptic.Marshal(toType0(args[0]), args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execP224(_ int, p *gop.Context) {
	ret0 := elliptic.P224()
	p.Ret(0, ret0)
}

func execP256(_ int, p *gop.Context) {
	ret0 := elliptic.P256()
	p.Ret(0, ret0)
}

func execP384(_ int, p *gop.Context) {
	ret0 := elliptic.P384()
	p.Ret(0, ret0)
}

func execP521(_ int, p *gop.Context) {
	ret0 := elliptic.P521()
	p.Ret(0, ret0)
}

func execUnmarshal(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := elliptic.Unmarshal(toType0(args[0]), args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/elliptic")

func init() {
	I.RegisterFuncs(
		I.Func("(Curve).Add", (elliptic.Curve).Add, execiCurveAdd),
		I.Func("(Curve).Double", (elliptic.Curve).Double, execiCurveDouble),
		I.Func("(Curve).IsOnCurve", (elliptic.Curve).IsOnCurve, execiCurveIsOnCurve),
		I.Func("(Curve).Params", (elliptic.Curve).Params, execiCurveParams),
		I.Func("(Curve).ScalarBaseMult", (elliptic.Curve).ScalarBaseMult, execiCurveScalarBaseMult),
		I.Func("(Curve).ScalarMult", (elliptic.Curve).ScalarMult, execiCurveScalarMult),
		I.Func("(*CurveParams).Params", (*elliptic.CurveParams).Params, execmCurveParamsParams),
		I.Func("(*CurveParams).IsOnCurve", (*elliptic.CurveParams).IsOnCurve, execmCurveParamsIsOnCurve),
		I.Func("(*CurveParams).Add", (*elliptic.CurveParams).Add, execmCurveParamsAdd),
		I.Func("(*CurveParams).Double", (*elliptic.CurveParams).Double, execmCurveParamsDouble),
		I.Func("(*CurveParams).ScalarMult", (*elliptic.CurveParams).ScalarMult, execmCurveParamsScalarMult),
		I.Func("(*CurveParams).ScalarBaseMult", (*elliptic.CurveParams).ScalarBaseMult, execmCurveParamsScalarBaseMult),
		I.Func("GenerateKey", elliptic.GenerateKey, execGenerateKey),
		I.Func("Marshal", elliptic.Marshal, execMarshal),
		I.Func("P224", elliptic.P224, execP224),
		I.Func("P256", elliptic.P256, execP256),
		I.Func("P384", elliptic.P384, execP384),
		I.Func("P521", elliptic.P521, execP521),
		I.Func("Unmarshal", elliptic.Unmarshal, execUnmarshal),
	)
	I.RegisterTypes(
		I.Type("Curve", reflect.TypeOf((*elliptic.Curve)(nil)).Elem()),
		I.Type("CurveParams", reflect.TypeOf((*elliptic.CurveParams)(nil)).Elem()),
	)
}
