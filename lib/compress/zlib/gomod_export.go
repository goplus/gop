// Package zlib provide Go+ "compress/zlib" package, as "compress/zlib" package in Go.
package zlib

import (
	zlib "compress/zlib"
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
	ret0, ret1 := zlib.NewReader(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

func execNewReaderDict(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := zlib.NewReaderDict(toType0(args[0]), args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execNewWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := zlib.NewWriter(toType1(args[0]))
	p.Ret(1, ret0)
}

func execNewWriterLevel(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := zlib.NewWriterLevel(toType1(args[0]), args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execNewWriterLevelDict(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := zlib.NewWriterLevelDict(toType1(args[0]), args[1].(int), args[2].([]byte))
	p.Ret(3, ret0, ret1)
}

func execiResetterReset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(zlib.Resetter).Reset(toType0(args[1]), args[2].([]byte))
	p.Ret(3, ret0)
}

func execmWriterReset(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*zlib.Writer).Reset(toType1(args[1]))
	p.PopN(2)
}

func execmWriterWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*zlib.Writer).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmWriterFlush(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*zlib.Writer).Flush()
	p.Ret(1, ret0)
}

func execmWriterClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*zlib.Writer).Close()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("compress/zlib")

func init() {
	I.RegisterFuncs(
		I.Func("NewReader", zlib.NewReader, execNewReader),
		I.Func("NewReaderDict", zlib.NewReaderDict, execNewReaderDict),
		I.Func("NewWriter", zlib.NewWriter, execNewWriter),
		I.Func("NewWriterLevel", zlib.NewWriterLevel, execNewWriterLevel),
		I.Func("NewWriterLevelDict", zlib.NewWriterLevelDict, execNewWriterLevelDict),
		I.Func("(Resetter).Reset", (zlib.Resetter).Reset, execiResetterReset),
		I.Func("(*Writer).Reset", (*zlib.Writer).Reset, execmWriterReset),
		I.Func("(*Writer).Write", (*zlib.Writer).Write, execmWriterWrite),
		I.Func("(*Writer).Flush", (*zlib.Writer).Flush, execmWriterFlush),
		I.Func("(*Writer).Close", (*zlib.Writer).Close, execmWriterClose),
	)
	I.RegisterVars(
		I.Var("ErrChecksum", &zlib.ErrChecksum),
		I.Var("ErrDictionary", &zlib.ErrDictionary),
		I.Var("ErrHeader", &zlib.ErrHeader),
	)
	I.RegisterTypes(
		I.Type("Resetter", reflect.TypeOf((*zlib.Resetter)(nil)).Elem()),
		I.Type("Writer", reflect.TypeOf((*zlib.Writer)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("BestCompression", qspec.ConstUnboundInt, zlib.BestCompression),
		I.Const("BestSpeed", qspec.ConstUnboundInt, zlib.BestSpeed),
		I.Const("DefaultCompression", qspec.ConstUnboundInt, zlib.DefaultCompression),
		I.Const("HuffmanOnly", qspec.ConstUnboundInt, zlib.HuffmanOnly),
		I.Const("NoCompression", qspec.ConstUnboundInt, zlib.NoCompression),
	)
}
