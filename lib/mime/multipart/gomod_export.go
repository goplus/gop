// Package multipart provide Go+ "mime/multipart" package, as "mime/multipart" package in Go.
package multipart

import (
	io "io"
	multipart "mime/multipart"
	textproto "net/textproto"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execiFileClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(multipart.File).Close()
	p.Ret(1, ret0)
}

func execiFileRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(multipart.File).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiFileReadAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(multipart.File).ReadAt(args[1].([]byte), args[2].(int64))
	p.Ret(3, ret0, ret1)
}

func execiFileSeek(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(multipart.File).Seek(args[1].(int64), args[2].(int))
	p.Ret(3, ret0, ret1)
}

func execmFileHeaderOpen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*multipart.FileHeader).Open()
	p.Ret(1, ret0, ret1)
}

func execmFormRemoveAll(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*multipart.Form).RemoveAll()
	p.Ret(1, ret0)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewReader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := multipart.NewReader(toType0(args[0]), args[1].(string))
	p.Ret(2, ret0)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execNewWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := multipart.NewWriter(toType1(args[0]))
	p.Ret(1, ret0)
}

func execmPartFormName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*multipart.Part).FormName()
	p.Ret(1, ret0)
}

func execmPartFileName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*multipart.Part).FileName()
	p.Ret(1, ret0)
}

func execmPartRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*multipart.Part).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmPartClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*multipart.Part).Close()
	p.Ret(1, ret0)
}

func execmReaderReadForm(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*multipart.Reader).ReadForm(args[1].(int64))
	p.Ret(2, ret0, ret1)
}

func execmReaderNextPart(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*multipart.Reader).NextPart()
	p.Ret(1, ret0, ret1)
}

func execmReaderNextRawPart(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*multipart.Reader).NextRawPart()
	p.Ret(1, ret0, ret1)
}

func execmWriterBoundary(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*multipart.Writer).Boundary()
	p.Ret(1, ret0)
}

func execmWriterSetBoundary(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*multipart.Writer).SetBoundary(args[1].(string))
	p.Ret(2, ret0)
}

func execmWriterFormDataContentType(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*multipart.Writer).FormDataContentType()
	p.Ret(1, ret0)
}

func execmWriterCreatePart(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*multipart.Writer).CreatePart(args[1].(textproto.MIMEHeader))
	p.Ret(2, ret0, ret1)
}

func execmWriterCreateFormFile(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*multipart.Writer).CreateFormFile(args[1].(string), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execmWriterCreateFormField(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*multipart.Writer).CreateFormField(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmWriterWriteField(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*multipart.Writer).WriteField(args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execmWriterClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*multipart.Writer).Close()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("mime/multipart")

func init() {
	I.RegisterFuncs(
		I.Func("(File).Close", (multipart.File).Close, execiFileClose),
		I.Func("(File).Read", (multipart.File).Read, execiFileRead),
		I.Func("(File).ReadAt", (multipart.File).ReadAt, execiFileReadAt),
		I.Func("(File).Seek", (multipart.File).Seek, execiFileSeek),
		I.Func("(*FileHeader).Open", (*multipart.FileHeader).Open, execmFileHeaderOpen),
		I.Func("(*Form).RemoveAll", (*multipart.Form).RemoveAll, execmFormRemoveAll),
		I.Func("NewReader", multipart.NewReader, execNewReader),
		I.Func("NewWriter", multipart.NewWriter, execNewWriter),
		I.Func("(*Part).FormName", (*multipart.Part).FormName, execmPartFormName),
		I.Func("(*Part).FileName", (*multipart.Part).FileName, execmPartFileName),
		I.Func("(*Part).Read", (*multipart.Part).Read, execmPartRead),
		I.Func("(*Part).Close", (*multipart.Part).Close, execmPartClose),
		I.Func("(*Reader).ReadForm", (*multipart.Reader).ReadForm, execmReaderReadForm),
		I.Func("(*Reader).NextPart", (*multipart.Reader).NextPart, execmReaderNextPart),
		I.Func("(*Reader).NextRawPart", (*multipart.Reader).NextRawPart, execmReaderNextRawPart),
		I.Func("(*Writer).Boundary", (*multipart.Writer).Boundary, execmWriterBoundary),
		I.Func("(*Writer).SetBoundary", (*multipart.Writer).SetBoundary, execmWriterSetBoundary),
		I.Func("(*Writer).FormDataContentType", (*multipart.Writer).FormDataContentType, execmWriterFormDataContentType),
		I.Func("(*Writer).CreatePart", (*multipart.Writer).CreatePart, execmWriterCreatePart),
		I.Func("(*Writer).CreateFormFile", (*multipart.Writer).CreateFormFile, execmWriterCreateFormFile),
		I.Func("(*Writer).CreateFormField", (*multipart.Writer).CreateFormField, execmWriterCreateFormField),
		I.Func("(*Writer).WriteField", (*multipart.Writer).WriteField, execmWriterWriteField),
		I.Func("(*Writer).Close", (*multipart.Writer).Close, execmWriterClose),
	)
	I.RegisterVars(
		I.Var("ErrMessageTooLarge", &multipart.ErrMessageTooLarge),
	)
	I.RegisterTypes(
		I.Type("File", reflect.TypeOf((*multipart.File)(nil)).Elem()),
		I.Type("FileHeader", reflect.TypeOf((*multipart.FileHeader)(nil)).Elem()),
		I.Type("Form", reflect.TypeOf((*multipart.Form)(nil)).Elem()),
		I.Type("Part", reflect.TypeOf((*multipart.Part)(nil)).Elem()),
		I.Type("Reader", reflect.TypeOf((*multipart.Reader)(nil)).Elem()),
		I.Type("Writer", reflect.TypeOf((*multipart.Writer)(nil)).Elem()),
	)
}
