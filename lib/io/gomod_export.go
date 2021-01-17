// Package io provide Go+ "io" package, as "io" package in Go.
package io

import (
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execiByteReaderReadByte(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(io.ByteReader).ReadByte()
	p.Ret(1, ret0, ret1)
}

func execiByteScannerUnreadByte(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(io.ByteScanner).UnreadByte()
	p.Ret(1, ret0)
}

func execiByteWriterWriteByte(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(io.ByteWriter).WriteByte(args[1].(byte))
	p.Ret(2, ret0)
}

func execiCloserClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(io.Closer).Close()
	p.Ret(1, ret0)
}

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

func execmLimitedReaderRead(_ int, p *gop.Context) {
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

func execmPipeReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*io.PipeReader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmPipeReaderClose(_ int, p *gop.Context) {
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

func execmPipeReaderCloseWithError(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*io.PipeReader).CloseWithError(toType3(args[1]))
	p.Ret(2, ret0)
}

func execmPipeWriterWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*io.PipeWriter).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmPipeWriterClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*io.PipeWriter).Close()
	p.Ret(1, ret0)
}

func execmPipeWriterCloseWithError(_ int, p *gop.Context) {
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

func execiReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(io.Reader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiReaderAtReadAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(io.ReaderAt).ReadAt(args[1].([]byte), args[2].(int64))
	p.Ret(3, ret0, ret1)
}

func execiReaderFromReadFrom(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(io.ReaderFrom).ReadFrom(toType1(args[1]))
	p.Ret(2, ret0, ret1)
}

func execiRuneReaderReadRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(io.RuneReader).ReadRune()
	p.Ret(1, ret0, ret1, ret2)
}

func execiRuneScannerUnreadRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(io.RuneScanner).UnreadRune()
	p.Ret(1, ret0)
}

func execmSectionReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*io.SectionReader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmSectionReaderSeek(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*io.SectionReader).Seek(args[1].(int64), args[2].(int))
	p.Ret(3, ret0, ret1)
}

func execmSectionReaderReadAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*io.SectionReader).ReadAt(args[1].([]byte), args[2].(int64))
	p.Ret(3, ret0, ret1)
}

func execmSectionReaderSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*io.SectionReader).Size()
	p.Ret(1, ret0)
}

func execiSeekerSeek(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(io.Seeker).Seek(args[1].(int64), args[2].(int))
	p.Ret(3, ret0, ret1)
}

func execiStringWriterWriteString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(io.StringWriter).WriteString(args[1].(string))
	p.Ret(2, ret0, ret1)
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

func execiWriterWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(io.Writer).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiWriterAtWriteAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(io.WriterAt).WriteAt(args[1].([]byte), args[2].(int64))
	p.Ret(3, ret0, ret1)
}

func execiWriterToWriteTo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(io.WriterTo).WriteTo(toType0(args[1]))
	p.Ret(2, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("io")

func init() {
	I.RegisterFuncs(
		I.Func("(ByteReader).ReadByte", (io.ByteReader).ReadByte, execiByteReaderReadByte),
		I.Func("(ByteScanner).ReadByte", (io.ByteReader).ReadByte, execiByteReaderReadByte),
		I.Func("(ByteScanner).UnreadByte", (io.ByteScanner).UnreadByte, execiByteScannerUnreadByte),
		I.Func("(ByteWriter).WriteByte", (io.ByteWriter).WriteByte, execiByteWriterWriteByte),
		I.Func("(Closer).Close", (io.Closer).Close, execiCloserClose),
		I.Func("Copy", io.Copy, execCopy),
		I.Func("CopyBuffer", io.CopyBuffer, execCopyBuffer),
		I.Func("CopyN", io.CopyN, execCopyN),
		I.Func("LimitReader", io.LimitReader, execLimitReader),
		I.Func("(*LimitedReader).Read", (*io.LimitedReader).Read, execmLimitedReaderRead),
		I.Func("NewSectionReader", io.NewSectionReader, execNewSectionReader),
		I.Func("Pipe", io.Pipe, execPipe),
		I.Func("(*PipeReader).Read", (*io.PipeReader).Read, execmPipeReaderRead),
		I.Func("(*PipeReader).Close", (*io.PipeReader).Close, execmPipeReaderClose),
		I.Func("(*PipeReader).CloseWithError", (*io.PipeReader).CloseWithError, execmPipeReaderCloseWithError),
		I.Func("(*PipeWriter).Write", (*io.PipeWriter).Write, execmPipeWriterWrite),
		I.Func("(*PipeWriter).Close", (*io.PipeWriter).Close, execmPipeWriterClose),
		I.Func("(*PipeWriter).CloseWithError", (*io.PipeWriter).CloseWithError, execmPipeWriterCloseWithError),
		I.Func("ReadAtLeast", io.ReadAtLeast, execReadAtLeast),
		I.Func("(ReadCloser).Close", (io.Closer).Close, execiCloserClose),
		I.Func("(ReadCloser).Read", (io.Reader).Read, execiReaderRead),
		I.Func("ReadFull", io.ReadFull, execReadFull),
		I.Func("(ReadSeeker).Read", (io.Reader).Read, execiReaderRead),
		I.Func("(ReadSeeker).Seek", (io.Seeker).Seek, execiSeekerSeek),
		I.Func("(ReadWriteCloser).Close", (io.Closer).Close, execiCloserClose),
		I.Func("(ReadWriteCloser).Read", (io.Reader).Read, execiReaderRead),
		I.Func("(ReadWriteCloser).Write", (io.Writer).Write, execiWriterWrite),
		I.Func("(ReadWriteSeeker).Read", (io.Reader).Read, execiReaderRead),
		I.Func("(ReadWriteSeeker).Seek", (io.Seeker).Seek, execiSeekerSeek),
		I.Func("(ReadWriteSeeker).Write", (io.Writer).Write, execiWriterWrite),
		I.Func("(ReadWriter).Read", (io.Reader).Read, execiReaderRead),
		I.Func("(ReadWriter).Write", (io.Writer).Write, execiWriterWrite),
		I.Func("(Reader).Read", (io.Reader).Read, execiReaderRead),
		I.Func("(ReaderAt).ReadAt", (io.ReaderAt).ReadAt, execiReaderAtReadAt),
		I.Func("(ReaderFrom).ReadFrom", (io.ReaderFrom).ReadFrom, execiReaderFromReadFrom),
		I.Func("(RuneReader).ReadRune", (io.RuneReader).ReadRune, execiRuneReaderReadRune),
		I.Func("(RuneScanner).ReadRune", (io.RuneReader).ReadRune, execiRuneReaderReadRune),
		I.Func("(RuneScanner).UnreadRune", (io.RuneScanner).UnreadRune, execiRuneScannerUnreadRune),
		I.Func("(*SectionReader).Read", (*io.SectionReader).Read, execmSectionReaderRead),
		I.Func("(*SectionReader).Seek", (*io.SectionReader).Seek, execmSectionReaderSeek),
		I.Func("(*SectionReader).ReadAt", (*io.SectionReader).ReadAt, execmSectionReaderReadAt),
		I.Func("(*SectionReader).Size", (*io.SectionReader).Size, execmSectionReaderSize),
		I.Func("(Seeker).Seek", (io.Seeker).Seek, execiSeekerSeek),
		I.Func("(StringWriter).WriteString", (io.StringWriter).WriteString, execiStringWriterWriteString),
		I.Func("TeeReader", io.TeeReader, execTeeReader),
		I.Func("(WriteCloser).Close", (io.Closer).Close, execiCloserClose),
		I.Func("(WriteCloser).Write", (io.Writer).Write, execiWriterWrite),
		I.Func("(WriteSeeker).Seek", (io.Seeker).Seek, execiSeekerSeek),
		I.Func("(WriteSeeker).Write", (io.Writer).Write, execiWriterWrite),
		I.Func("WriteString", io.WriteString, execWriteString),
		I.Func("(Writer).Write", (io.Writer).Write, execiWriterWrite),
		I.Func("(WriterAt).WriteAt", (io.WriterAt).WriteAt, execiWriterAtWriteAt),
		I.Func("(WriterTo).WriteTo", (io.WriterTo).WriteTo, execiWriterToWriteTo),
	)
	I.RegisterFuncvs(
		I.Funcv("MultiReader", io.MultiReader, execMultiReader),
		I.Funcv("MultiWriter", io.MultiWriter, execMultiWriter),
	)
	I.RegisterVars(
		I.Var("EOF", &io.EOF),
		I.Var("ErrClosedPipe", &io.ErrClosedPipe),
		I.Var("ErrNoProgress", &io.ErrNoProgress),
		I.Var("ErrShortBuffer", &io.ErrShortBuffer),
		I.Var("ErrShortWrite", &io.ErrShortWrite),
		I.Var("ErrUnexpectedEOF", &io.ErrUnexpectedEOF),
	)
	I.RegisterTypes(
		I.Type("ByteReader", reflect.TypeOf((*io.ByteReader)(nil)).Elem()),
		I.Type("ByteScanner", reflect.TypeOf((*io.ByteScanner)(nil)).Elem()),
		I.Type("ByteWriter", reflect.TypeOf((*io.ByteWriter)(nil)).Elem()),
		I.Type("Closer", reflect.TypeOf((*io.Closer)(nil)).Elem()),
		I.Type("LimitedReader", reflect.TypeOf((*io.LimitedReader)(nil)).Elem()),
		I.Type("PipeReader", reflect.TypeOf((*io.PipeReader)(nil)).Elem()),
		I.Type("PipeWriter", reflect.TypeOf((*io.PipeWriter)(nil)).Elem()),
		I.Type("ReadCloser", reflect.TypeOf((*io.ReadCloser)(nil)).Elem()),
		I.Type("ReadSeeker", reflect.TypeOf((*io.ReadSeeker)(nil)).Elem()),
		I.Type("ReadWriteCloser", reflect.TypeOf((*io.ReadWriteCloser)(nil)).Elem()),
		I.Type("ReadWriteSeeker", reflect.TypeOf((*io.ReadWriteSeeker)(nil)).Elem()),
		I.Type("ReadWriter", reflect.TypeOf((*io.ReadWriter)(nil)).Elem()),
		I.Type("Reader", reflect.TypeOf((*io.Reader)(nil)).Elem()),
		I.Type("ReaderAt", reflect.TypeOf((*io.ReaderAt)(nil)).Elem()),
		I.Type("ReaderFrom", reflect.TypeOf((*io.ReaderFrom)(nil)).Elem()),
		I.Type("RuneReader", reflect.TypeOf((*io.RuneReader)(nil)).Elem()),
		I.Type("RuneScanner", reflect.TypeOf((*io.RuneScanner)(nil)).Elem()),
		I.Type("SectionReader", reflect.TypeOf((*io.SectionReader)(nil)).Elem()),
		I.Type("Seeker", reflect.TypeOf((*io.Seeker)(nil)).Elem()),
		I.Type("StringWriter", reflect.TypeOf((*io.StringWriter)(nil)).Elem()),
		I.Type("WriteCloser", reflect.TypeOf((*io.WriteCloser)(nil)).Elem()),
		I.Type("WriteSeeker", reflect.TypeOf((*io.WriteSeeker)(nil)).Elem()),
		I.Type("Writer", reflect.TypeOf((*io.Writer)(nil)).Elem()),
		I.Type("WriterAt", reflect.TypeOf((*io.WriterAt)(nil)).Elem()),
		I.Type("WriterTo", reflect.TypeOf((*io.WriterTo)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("SeekCurrent", qspec.ConstUnboundInt, io.SeekCurrent),
		I.Const("SeekEnd", qspec.ConstUnboundInt, io.SeekEnd),
		I.Const("SeekStart", qspec.ConstUnboundInt, io.SeekStart),
	)
}
