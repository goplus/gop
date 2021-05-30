// Package bzip2 provide Go+ "compress/bzip2" package, as "compress/bzip2" package in Go.
package bzip2

import (
	bzip2 "compress/bzip2"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewReader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bzip2.NewReader(toType0(args[0]))
	p.Ret(1, ret0)
}

func execmStructuralErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(bzip2.StructuralError).Error()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("compress/bzip2")

func init() {
	I.RegisterFuncs(
		I.Func("NewReader", bzip2.NewReader, execNewReader),
		I.Func("(StructuralError).Error", (bzip2.StructuralError).Error, execmStructuralErrorError),
	)
	I.RegisterTypes(
		I.Type("StructuralError", reflect.TypeOf((*bzip2.StructuralError)(nil)).Elem()),
	)
}
