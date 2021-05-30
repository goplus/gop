// Package gzip provide Go+ "compress/gzip" package, as "compress/gzip" package in Go.
package gzip

import (
	gzip "compress/gzip"
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
	args := p.GetArgs(1)
	ret0, ret1 := gzip.NewReader(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execNewWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := gzip.NewWriter(toType1(args[0]))
	p.Ret(1, ret0)
}

func execNewWriterLevel(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := gzip.NewWriterLevel(toType1(args[0]), args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execmReaderReset(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*gzip.Reader).Reset(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmReaderMultistream(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*gzip.Reader).Multistream(args[1].(bool))
	p.PopN(2)
}

func execmReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*gzip.Reader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmReaderClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*gzip.Reader).Close()
	p.Ret(1, ret0)
}

func execmWriterReset(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*gzip.Writer).Reset(toType1(args[1]))
	p.PopN(2)
}

func execmWriterWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*gzip.Writer).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmWriterFlush(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*gzip.Writer).Flush()
	p.Ret(1, ret0)
}

func execmWriterClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*gzip.Writer).Close()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("compress/gzip")

func init() {
	I.RegisterFuncs(
		I.Func("NewReader", gzip.NewReader, execNewReader),
		I.Func("NewWriter", gzip.NewWriter, execNewWriter),
		I.Func("NewWriterLevel", gzip.NewWriterLevel, execNewWriterLevel),
		I.Func("(*Reader).Reset", (*gzip.Reader).Reset, execmReaderReset),
		I.Func("(*Reader).Multistream", (*gzip.Reader).Multistream, execmReaderMultistream),
		I.Func("(*Reader).Read", (*gzip.Reader).Read, execmReaderRead),
		I.Func("(*Reader).Close", (*gzip.Reader).Close, execmReaderClose),
		I.Func("(*Writer).Reset", (*gzip.Writer).Reset, execmWriterReset),
		I.Func("(*Writer).Write", (*gzip.Writer).Write, execmWriterWrite),
		I.Func("(*Writer).Flush", (*gzip.Writer).Flush, execmWriterFlush),
		I.Func("(*Writer).Close", (*gzip.Writer).Close, execmWriterClose),
	)
	I.RegisterVars(
		I.Var("ErrChecksum", &gzip.ErrChecksum),
		I.Var("ErrHeader", &gzip.ErrHeader),
	)
	I.RegisterTypes(
		I.Type("Header", reflect.TypeOf((*gzip.Header)(nil)).Elem()),
		I.Type("Reader", reflect.TypeOf((*gzip.Reader)(nil)).Elem()),
		I.Type("Writer", reflect.TypeOf((*gzip.Writer)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("BestCompression", qspec.ConstUnboundInt, gzip.BestCompression),
		I.Const("BestSpeed", qspec.ConstUnboundInt, gzip.BestSpeed),
		I.Const("DefaultCompression", qspec.ConstUnboundInt, gzip.DefaultCompression),
		I.Const("HuffmanOnly", qspec.ConstUnboundInt, gzip.HuffmanOnly),
		I.Const("NoCompression", qspec.ConstUnboundInt, gzip.NoCompression),
	)
}
