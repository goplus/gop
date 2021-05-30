// Package lzw provide Go+ "compress/lzw" package, as "compress/lzw" package in Go.
package lzw

import (
	lzw "compress/lzw"
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

func execNewReader(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := lzw.NewReader(toType0(args[0]), args[1].(lzw.Order), args[2].(int))
	p.Ret(3, ret0)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execNewWriter(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := lzw.NewWriter(toType1(args[0]), args[1].(lzw.Order), args[2].(int))
	p.Ret(3, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("compress/lzw")

func init() {
	I.RegisterFuncs(
		I.Func("NewReader", lzw.NewReader, execNewReader),
		I.Func("NewWriter", lzw.NewWriter, execNewWriter),
	)
	I.RegisterTypes(
		I.Type("Order", reflect.TypeOf((*lzw.Order)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("LSB", qspec.Int, lzw.LSB),
		I.Const("MSB", qspec.Int, lzw.MSB),
	)
}
