// Package flate provide Go+ "compress/flate" package, as "compress/flate" package in Go.
package flate

import (
	flate "compress/flate"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmCorruptInputErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(flate.CorruptInputError).Error()
	p.Ret(1, ret0)
}

func execmInternalErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(flate.InternalError).Error()
	p.Ret(1, ret0)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewReader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := flate.NewReader(toType0(args[0]))
	p.Ret(1, ret0)
}

func execNewReaderDict(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := flate.NewReaderDict(toType0(args[0]), args[1].([]byte))
	p.Ret(2, ret0)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execNewWriter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := flate.NewWriter(toType1(args[0]), args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execNewWriterDict(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := flate.NewWriterDict(toType1(args[0]), args[1].(int), args[2].([]byte))
	p.Ret(3, ret0, ret1)
}

func execmReadErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*flate.ReadError).Error()
	p.Ret(1, ret0)
}

func execiReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(flate.Reader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiReaderReadByte(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(flate.Reader).ReadByte()
	p.Ret(1, ret0, ret1)
}

func execiResetterReset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(flate.Resetter).Reset(toType0(args[1]), args[2].([]byte))
	p.Ret(3, ret0)
}

func execmWriteErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*flate.WriteError).Error()
	p.Ret(1, ret0)
}

func execmWriterWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*flate.Writer).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmWriterFlush(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*flate.Writer).Flush()
	p.Ret(1, ret0)
}

func execmWriterClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*flate.Writer).Close()
	p.Ret(1, ret0)
}

func execmWriterReset(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*flate.Writer).Reset(toType1(args[1]))
	p.PopN(2)
}

// I is a Go package instance.
var I = gop.NewGoPackage("compress/flate")

func init() {
	I.RegisterFuncs(
		I.Func("(CorruptInputError).Error", (flate.CorruptInputError).Error, execmCorruptInputErrorError),
		I.Func("(InternalError).Error", (flate.InternalError).Error, execmInternalErrorError),
		I.Func("NewReader", flate.NewReader, execNewReader),
		I.Func("NewReaderDict", flate.NewReaderDict, execNewReaderDict),
		I.Func("NewWriter", flate.NewWriter, execNewWriter),
		I.Func("NewWriterDict", flate.NewWriterDict, execNewWriterDict),
		I.Func("(*ReadError).Error", (*flate.ReadError).Error, execmReadErrorError),
		I.Func("(Reader).Read", (flate.Reader).Read, execiReaderRead),
		I.Func("(Reader).ReadByte", (flate.Reader).ReadByte, execiReaderReadByte),
		I.Func("(Resetter).Reset", (flate.Resetter).Reset, execiResetterReset),
		I.Func("(*WriteError).Error", (*flate.WriteError).Error, execmWriteErrorError),
		I.Func("(*Writer).Write", (*flate.Writer).Write, execmWriterWrite),
		I.Func("(*Writer).Flush", (*flate.Writer).Flush, execmWriterFlush),
		I.Func("(*Writer).Close", (*flate.Writer).Close, execmWriterClose),
		I.Func("(*Writer).Reset", (*flate.Writer).Reset, execmWriterReset),
	)
	I.RegisterTypes(
		I.Type("CorruptInputError", reflect.TypeOf((*flate.CorruptInputError)(nil)).Elem()),
		I.Type("InternalError", reflect.TypeOf((*flate.InternalError)(nil)).Elem()),
		I.Type("ReadError", reflect.TypeOf((*flate.ReadError)(nil)).Elem()),
		I.Type("Reader", reflect.TypeOf((*flate.Reader)(nil)).Elem()),
		I.Type("Resetter", reflect.TypeOf((*flate.Resetter)(nil)).Elem()),
		I.Type("WriteError", reflect.TypeOf((*flate.WriteError)(nil)).Elem()),
		I.Type("Writer", reflect.TypeOf((*flate.Writer)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("BestCompression", qspec.ConstUnboundInt, flate.BestCompression),
		I.Const("BestSpeed", qspec.ConstUnboundInt, flate.BestSpeed),
		I.Const("DefaultCompression", qspec.ConstUnboundInt, flate.DefaultCompression),
		I.Const("HuffmanOnly", qspec.ConstUnboundInt, flate.HuffmanOnly),
		I.Const("NoCompression", qspec.ConstUnboundInt, flate.NoCompression),
	)
}
