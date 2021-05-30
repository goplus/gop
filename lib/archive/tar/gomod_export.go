// Package tar provide Go+ "archive/tar" package, as "archive/tar" package in Go.
package tar

import (
	tar "archive/tar"
	io "io"
	os "os"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func toType0(v interface{}) os.FileInfo {
	if v == nil {
		return nil
	}
	return v.(os.FileInfo)
}

func execFileInfoHeader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := tar.FileInfoHeader(toType0(args[0]), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmFormatString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(tar.Format).String()
	p.Ret(1, ret0)
}

func execmHeaderFileInfo(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*tar.Header).FileInfo()
	p.Ret(1, ret0)
}

func toType1(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewReader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := tar.NewReader(toType1(args[0]))
	p.Ret(1, ret0)
}

func toType2(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execNewWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := tar.NewWriter(toType2(args[0]))
	p.Ret(1, ret0)
}

func execmReaderNext(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*tar.Reader).Next()
	p.Ret(1, ret0, ret1)
}

func execmReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*tar.Reader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmWriterFlush(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*tar.Writer).Flush()
	p.Ret(1, ret0)
}

func execmWriterWriteHeader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*tar.Writer).WriteHeader(args[1].(*tar.Header))
	p.Ret(2, ret0)
}

func execmWriterWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*tar.Writer).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmWriterClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*tar.Writer).Close()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("archive/tar")

func init() {
	I.RegisterFuncs(
		I.Func("FileInfoHeader", tar.FileInfoHeader, execFileInfoHeader),
		I.Func("(Format).String", (tar.Format).String, execmFormatString),
		I.Func("(*Header).FileInfo", (*tar.Header).FileInfo, execmHeaderFileInfo),
		I.Func("NewReader", tar.NewReader, execNewReader),
		I.Func("NewWriter", tar.NewWriter, execNewWriter),
		I.Func("(*Reader).Next", (*tar.Reader).Next, execmReaderNext),
		I.Func("(*Reader).Read", (*tar.Reader).Read, execmReaderRead),
		I.Func("(*Writer).Flush", (*tar.Writer).Flush, execmWriterFlush),
		I.Func("(*Writer).WriteHeader", (*tar.Writer).WriteHeader, execmWriterWriteHeader),
		I.Func("(*Writer).Write", (*tar.Writer).Write, execmWriterWrite),
		I.Func("(*Writer).Close", (*tar.Writer).Close, execmWriterClose),
	)
	I.RegisterVars(
		I.Var("ErrFieldTooLong", &tar.ErrFieldTooLong),
		I.Var("ErrHeader", &tar.ErrHeader),
		I.Var("ErrWriteAfterClose", &tar.ErrWriteAfterClose),
		I.Var("ErrWriteTooLong", &tar.ErrWriteTooLong),
	)
	I.RegisterTypes(
		I.Type("Format", reflect.TypeOf((*tar.Format)(nil)).Elem()),
		I.Type("Header", reflect.TypeOf((*tar.Header)(nil)).Elem()),
		I.Type("Reader", reflect.TypeOf((*tar.Reader)(nil)).Elem()),
		I.Type("Writer", reflect.TypeOf((*tar.Writer)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("FormatGNU", qspec.Int, tar.FormatGNU),
		I.Const("FormatPAX", qspec.Int, tar.FormatPAX),
		I.Const("FormatUSTAR", qspec.Int, tar.FormatUSTAR),
		I.Const("FormatUnknown", qspec.Int, tar.FormatUnknown),
		I.Const("TypeBlock", qspec.ConstBoundRune, tar.TypeBlock),
		I.Const("TypeChar", qspec.ConstBoundRune, tar.TypeChar),
		I.Const("TypeCont", qspec.ConstBoundRune, tar.TypeCont),
		I.Const("TypeDir", qspec.ConstBoundRune, tar.TypeDir),
		I.Const("TypeFifo", qspec.ConstBoundRune, tar.TypeFifo),
		I.Const("TypeGNULongLink", qspec.ConstBoundRune, tar.TypeGNULongLink),
		I.Const("TypeGNULongName", qspec.ConstBoundRune, tar.TypeGNULongName),
		I.Const("TypeGNUSparse", qspec.ConstBoundRune, tar.TypeGNUSparse),
		I.Const("TypeLink", qspec.ConstBoundRune, tar.TypeLink),
		I.Const("TypeReg", qspec.ConstBoundRune, tar.TypeReg),
		I.Const("TypeRegA", qspec.ConstBoundRune, tar.TypeRegA),
		I.Const("TypeSymlink", qspec.ConstBoundRune, tar.TypeSymlink),
		I.Const("TypeXGlobalHeader", qspec.ConstBoundRune, tar.TypeXGlobalHeader),
		I.Const("TypeXHeader", qspec.ConstBoundRune, tar.TypeXHeader),
	)
}
