// Package cmplx provide Go+ "math/cmplx" package, as "math/cmplx" package in Go.
package cmplx

import (
	cmplx "math/cmplx"

	gop "github.com/goplus/gop"
)

func execAbs(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Abs(args[0].(complex128))
	p.Ret(1, ret0)
}

func execAcos(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Acos(args[0].(complex128))
	p.Ret(1, ret0)
}

func execAcosh(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Acosh(args[0].(complex128))
	p.Ret(1, ret0)
}

func execAsin(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Asin(args[0].(complex128))
	p.Ret(1, ret0)
}

func execAsinh(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Asinh(args[0].(complex128))
	p.Ret(1, ret0)
}

func execAtan(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Atan(args[0].(complex128))
	p.Ret(1, ret0)
}

func execAtanh(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Atanh(args[0].(complex128))
	p.Ret(1, ret0)
}

func execConj(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Conj(args[0].(complex128))
	p.Ret(1, ret0)
}

func execCos(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Cos(args[0].(complex128))
	p.Ret(1, ret0)
}

func execCosh(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Cosh(args[0].(complex128))
	p.Ret(1, ret0)
}

func execCot(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Cot(args[0].(complex128))
	p.Ret(1, ret0)
}

func execExp(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Exp(args[0].(complex128))
	p.Ret(1, ret0)
}

func execInf(_ int, p *gop.Context) {
	ret0 := cmplx.Inf()
	p.Ret(0, ret0)
}

func execIsInf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.IsInf(args[0].(complex128))
	p.Ret(1, ret0)
}

func execIsNaN(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.IsNaN(args[0].(complex128))
	p.Ret(1, ret0)
}

func execLog(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Log(args[0].(complex128))
	p.Ret(1, ret0)
}

func execLog10(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Log10(args[0].(complex128))
	p.Ret(1, ret0)
}

func execNaN(_ int, p *gop.Context) {
	ret0 := cmplx.NaN()
	p.Ret(0, ret0)
}

func execPhase(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Phase(args[0].(complex128))
	p.Ret(1, ret0)
}

func execPolar(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := cmplx.Polar(args[0].(complex128))
	p.Ret(1, ret0, ret1)
}

func execPow(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := cmplx.Pow(args[0].(complex128), args[1].(complex128))
	p.Ret(2, ret0)
}

func execRect(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := cmplx.Rect(args[0].(float64), args[1].(float64))
	p.Ret(2, ret0)
}

func execSin(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Sin(args[0].(complex128))
	p.Ret(1, ret0)
}

func execSinh(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Sinh(args[0].(complex128))
	p.Ret(1, ret0)
}

func execSqrt(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Sqrt(args[0].(complex128))
	p.Ret(1, ret0)
}

func execTan(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Tan(args[0].(complex128))
	p.Ret(1, ret0)
}

func execTanh(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cmplx.Tanh(args[0].(complex128))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("math/cmplx")

func init() {
	I.RegisterFuncs(
		I.Func("Abs", cmplx.Abs, execAbs),
		I.Func("Acos", cmplx.Acos, execAcos),
		I.Func("Acosh", cmplx.Acosh, execAcosh),
		I.Func("Asin", cmplx.Asin, execAsin),
		I.Func("Asinh", cmplx.Asinh, execAsinh),
		I.Func("Atan", cmplx.Atan, execAtan),
		I.Func("Atanh", cmplx.Atanh, execAtanh),
		I.Func("Conj", cmplx.Conj, execConj),
		I.Func("Cos", cmplx.Cos, execCos),
		I.Func("Cosh", cmplx.Cosh, execCosh),
		I.Func("Cot", cmplx.Cot, execCot),
		I.Func("Exp", cmplx.Exp, execExp),
		I.Func("Inf", cmplx.Inf, execInf),
		I.Func("IsInf", cmplx.IsInf, execIsInf),
		I.Func("IsNaN", cmplx.IsNaN, execIsNaN),
		I.Func("Log", cmplx.Log, execLog),
		I.Func("Log10", cmplx.Log10, execLog10),
		I.Func("NaN", cmplx.NaN, execNaN),
		I.Func("Phase", cmplx.Phase, execPhase),
		I.Func("Polar", cmplx.Polar, execPolar),
		I.Func("Pow", cmplx.Pow, execPow),
		I.Func("Rect", cmplx.Rect, execRect),
		I.Func("Sin", cmplx.Sin, execSin),
		I.Func("Sinh", cmplx.Sinh, execSinh),
		I.Func("Sqrt", cmplx.Sqrt, execSqrt),
		I.Func("Tan", cmplx.Tan, execTan),
		I.Func("Tanh", cmplx.Tanh, execTanh),
	)
}
