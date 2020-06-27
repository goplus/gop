// Package io provide Go+ "io" package, as "io" package in Go.
package io

import (
	io "io"

	gop "github.com/qiniu/goplus/gop"
)

func toType0(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func toType1(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execCopy(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := io.Copy(toType0(args[0]), toType1(args[1]))
	p.Ret(2, ret0, ret1)
}

func execCopyBuffer(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := io.CopyBuffer(toType0(args[0]), toType1(args[1]), args[2].([]byte))
	p.Ret(3, ret0, ret1)
}

func execCopyN(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := io.CopyN(toType0(args[0]), toType1(args[1]), args[2].(int64))
	p.Ret(3, ret0, ret1)
}

func execLimitReader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := io.LimitReader(toType1(args[0]), args[1].(int64))
	p.Ret(2, ret0)
}

func execLimitedReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*io.LimitedReader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func toSlice0(args []interface{}) []io.Reader {
	ret := make([]io.Reader, len(args))
	for i, arg := range args {
		ret[i] = toType1(arg)
	}
	return ret
}

func execMultiReader(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := io.MultiReader(toSlice0(args)...)
	p.Ret(arity, ret0)
}

func toSlice1(args []interface{}) []io.Writer {
	ret := make([]io.Writer, len(args))
	for i, arg := range args {
		ret[i] = toType0(arg)
	}
	return ret
}

func execMultiWriter(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := io.MultiWriter(toSlice1(args)...)
	p.Ret(arity, ret0)
}

func toType2(v interface{}) io.ReaderAt {
	if v == nil {
		return nil
	}
	return v.(io.ReaderAt)
}

func execNewSectionReader(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := io.NewSectionReader(toType2(args[0]), args[1].(int64), args[2].(int64))
	p.Ret(3, ret0)
}

func execPipe(_ int, p *gop.Context) {
	ret0, ret1 := io.Pipe()
	p.Ret(0, ret0, ret1)
}

func execPipeReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*io.PipeReader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execPipeReaderClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*io.PipeReader).Close()
	p.Ret(1, ret0)
}

func toType3(v interface{}) error {
	if v == nil {
		return nil
	}
	return v.(error)
}

func execPipeReaderCloseWithError(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*io.PipeReader).CloseWithError(toType3(args[1]))
	p.Ret(2, ret0)
}

func execPipeWriterWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*io.PipeWriter).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execPipeWriterClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*io.PipeWriter).Close()
	p.Ret(1, ret0)
}

func execPipeWriterCloseWithError(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*io.PipeWriter).CloseWithError(toType3(args[1]))
	p.Ret(2, ret0)
}

func execReadAtLeast(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := io.ReadAtLeast(toType1(args[0]), args[1].([]byte), args[2].(int))
	p.Ret(3, ret0, ret1)
}

func execReadFull(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := io.ReadFull(toType1(args[0]), args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execSectionReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*io.SectionReader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execSectionReaderSeek(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*io.SectionReader).Seek(args[1].(int64), args[2].(int))
	p.Ret(3, ret0, ret1)
}

func execSectionReaderReadAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*io.SectionReader).ReadAt(args[1].([]byte), args[2].(int64))
	p.Ret(3, ret0, ret1)
}

func execSectionReaderSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*io.SectionReader).Size()
	p.Ret(1, ret0)
}

func execTeeReader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := io.TeeReader(toType1(args[0]), toType0(args[1]))
	p.Ret(2, ret0)
}

func execWriteString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := io.WriteString(toType0(args[0]), args[1].(string))
	p.Ret(2, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("io")

func init() {
	I.RegisterFuncs(
		I.Func("Copy", io.Copy, execCopy),
		I.Func("CopyBuffer", io.CopyBuffer, execCopyBuffer),
		I.Func("CopyN", io.CopyN, execCopyN),
		I.Func("LimitReader", io.LimitReader, execLimitReader),
		I.Func("(*LimitedReader).Read", (*io.LimitedReader).Read, execLimitedReaderRead),
		I.Func("NewSectionReader", io.NewSectionReader, execNewSectionReader),
		I.Func("Pipe", io.Pipe, execPipe),
		I.Func("(*PipeReader).Read", (*io.PipeReader).Read, execPipeReaderRead),
		I.Func("(*PipeReader).Close", (*io.PipeReader).Close, execPipeReaderClose),
		I.Func("(*PipeReader).CloseWithError", (*io.PipeReader).CloseWithError, execPipeReaderCloseWithError),
		I.Func("(*PipeWriter).Write", (*io.PipeWriter).Write, execPipeWriterWrite),
		I.Func("(*PipeWriter).Close", (*io.PipeWriter).Close, execPipeWriterClose),
		I.Func("(*PipeWriter).CloseWithError", (*io.PipeWriter).CloseWithError, execPipeWriterCloseWithError),
		I.Func("ReadAtLeast", io.ReadAtLeast, execReadAtLeast),
		I.Func("ReadFull", io.ReadFull, execReadFull),
		I.Func("(*SectionReader).Read", (*io.SectionReader).Read, execSectionReaderRead),
		I.Func("(*SectionReader).Seek", (*io.SectionReader).Seek, execSectionReaderSeek),
		I.Func("(*SectionReader).ReadAt", (*io.SectionReader).ReadAt, execSectionReaderReadAt),
		I.Func("(*SectionReader).Size", (*io.SectionReader).Size, execSectionReaderSize),
		I.Func("TeeReader", io.TeeReader, execTeeReader),
		I.Func("WriteString", io.WriteString, execWriteString),
	)
	I.RegisterFuncvs(
		I.Funcv("MultiReader", io.MultiReader, execMultiReader),
		I.Funcv("MultiWriter", io.MultiWriter, execMultiWriter),
	)
}
