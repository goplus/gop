// Package zip provide Go+ "archive/zip" package, as "archive/zip" package in Go.
package zip

import (
	zip "archive/zip"
	io "io"
	os "os"
	reflect "reflect"
	time "time"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmFileDataOffset(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*zip.File).DataOffset()
	p.Ret(1, ret0, ret1)
}

func execmFileOpen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*zip.File).Open()
	p.Ret(1, ret0, ret1)
}

func execmFileHeaderFileInfo(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*zip.FileHeader).FileInfo()
	p.Ret(1, ret0)
}

func execmFileHeaderModTime(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*zip.FileHeader).ModTime()
	p.Ret(1, ret0)
}

func execmFileHeaderSetModTime(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*zip.FileHeader).SetModTime(args[1].(time.Time))
	p.PopN(2)
}

func execmFileHeaderMode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*zip.FileHeader).Mode()
	p.Ret(1, ret0)
}

func execmFileHeaderSetMode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*zip.FileHeader).SetMode(args[1].(os.FileMode))
	p.PopN(2)
}

func toType0(v interface{}) os.FileInfo {
	if v == nil {
		return nil
	}
	return v.(os.FileInfo)
}

func execFileInfoHeader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := zip.FileInfoHeader(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

func toType1(v interface{}) io.ReaderAt {
	if v == nil {
		return nil
	}
	return v.(io.ReaderAt)
}

func execNewReader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := zip.NewReader(toType1(args[0]), args[1].(int64))
	p.Ret(2, ret0, ret1)
}

func toType2(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execNewWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := zip.NewWriter(toType2(args[0]))
	p.Ret(1, ret0)
}

func execOpenReader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := zip.OpenReader(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execmReadCloserClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*zip.ReadCloser).Close()
	p.Ret(1, ret0)
}

func execmReaderRegisterDecompressor(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*zip.Reader).RegisterDecompressor(args[1].(uint16), args[2].(zip.Decompressor))
	p.PopN(3)
}

func execRegisterCompressor(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	zip.RegisterCompressor(args[0].(uint16), args[1].(zip.Compressor))
	p.PopN(2)
}

func execRegisterDecompressor(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	zip.RegisterDecompressor(args[0].(uint16), args[1].(zip.Decompressor))
	p.PopN(2)
}

func execmWriterSetOffset(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*zip.Writer).SetOffset(args[1].(int64))
	p.PopN(2)
}

func execmWriterFlush(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*zip.Writer).Flush()
	p.Ret(1, ret0)
}

func execmWriterSetComment(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*zip.Writer).SetComment(args[1].(string))
	p.Ret(2, ret0)
}

func execmWriterClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*zip.Writer).Close()
	p.Ret(1, ret0)
}

func execmWriterCreate(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*zip.Writer).Create(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmWriterCreateHeader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*zip.Writer).CreateHeader(args[1].(*zip.FileHeader))
	p.Ret(2, ret0, ret1)
}

func execmWriterRegisterCompressor(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*zip.Writer).RegisterCompressor(args[1].(uint16), args[2].(zip.Compressor))
	p.PopN(3)
}

// I is a Go package instance.
var I = gop.NewGoPackage("archive/zip")

func init() {
	I.RegisterFuncs(
		I.Func("(*File).DataOffset", (*zip.File).DataOffset, execmFileDataOffset),
		I.Func("(*File).Open", (*zip.File).Open, execmFileOpen),
		I.Func("(*FileHeader).FileInfo", (*zip.FileHeader).FileInfo, execmFileHeaderFileInfo),
		I.Func("(*FileHeader).ModTime", (*zip.FileHeader).ModTime, execmFileHeaderModTime),
		I.Func("(*FileHeader).SetModTime", (*zip.FileHeader).SetModTime, execmFileHeaderSetModTime),
		I.Func("(*FileHeader).Mode", (*zip.FileHeader).Mode, execmFileHeaderMode),
		I.Func("(*FileHeader).SetMode", (*zip.FileHeader).SetMode, execmFileHeaderSetMode),
		I.Func("FileInfoHeader", zip.FileInfoHeader, execFileInfoHeader),
		I.Func("NewReader", zip.NewReader, execNewReader),
		I.Func("NewWriter", zip.NewWriter, execNewWriter),
		I.Func("OpenReader", zip.OpenReader, execOpenReader),
		I.Func("(*ReadCloser).Close", (*zip.ReadCloser).Close, execmReadCloserClose),
		I.Func("(*Reader).RegisterDecompressor", (*zip.Reader).RegisterDecompressor, execmReaderRegisterDecompressor),
		I.Func("RegisterCompressor", zip.RegisterCompressor, execRegisterCompressor),
		I.Func("RegisterDecompressor", zip.RegisterDecompressor, execRegisterDecompressor),
		I.Func("(*Writer).SetOffset", (*zip.Writer).SetOffset, execmWriterSetOffset),
		I.Func("(*Writer).Flush", (*zip.Writer).Flush, execmWriterFlush),
		I.Func("(*Writer).SetComment", (*zip.Writer).SetComment, execmWriterSetComment),
		I.Func("(*Writer).Close", (*zip.Writer).Close, execmWriterClose),
		I.Func("(*Writer).Create", (*zip.Writer).Create, execmWriterCreate),
		I.Func("(*Writer).CreateHeader", (*zip.Writer).CreateHeader, execmWriterCreateHeader),
		I.Func("(*Writer).RegisterCompressor", (*zip.Writer).RegisterCompressor, execmWriterRegisterCompressor),
	)
	I.RegisterVars(
		I.Var("ErrAlgorithm", &zip.ErrAlgorithm),
		I.Var("ErrChecksum", &zip.ErrChecksum),
		I.Var("ErrFormat", &zip.ErrFormat),
	)
	I.RegisterTypes(
		I.Type("Compressor", reflect.TypeOf((*zip.Compressor)(nil)).Elem()),
		I.Type("Decompressor", reflect.TypeOf((*zip.Decompressor)(nil)).Elem()),
		I.Type("File", reflect.TypeOf((*zip.File)(nil)).Elem()),
		I.Type("FileHeader", reflect.TypeOf((*zip.FileHeader)(nil)).Elem()),
		I.Type("ReadCloser", reflect.TypeOf((*zip.ReadCloser)(nil)).Elem()),
		I.Type("Reader", reflect.TypeOf((*zip.Reader)(nil)).Elem()),
		I.Type("Writer", reflect.TypeOf((*zip.Writer)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("Deflate", qspec.Uint16, zip.Deflate),
		I.Const("Store", qspec.Uint16, zip.Store),
	)
}
