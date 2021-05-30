// Package textproto provide Go+ "net/textproto" package, as "net/textproto" package in Go.
package textproto

import (
	bufio "bufio"
	io "io"
	textproto "net/textproto"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execCanonicalMIMEHeaderKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := textproto.CanonicalMIMEHeaderKey(args[0].(string))
	p.Ret(1, ret0)
}

func execmConnClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*textproto.Conn).Close()
	p.Ret(1, ret0)
}

func execmConnCmd(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*textproto.Conn).Cmd(args[1].(string), args[2:]...)
	p.Ret(arity, ret0, ret1)
}

func execDial(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := textproto.Dial(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*textproto.Error).Error()
	p.Ret(1, ret0)
}

func execmMIMEHeaderAdd(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(textproto.MIMEHeader).Add(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execmMIMEHeaderSet(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(textproto.MIMEHeader).Set(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execmMIMEHeaderGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(textproto.MIMEHeader).Get(args[1].(string))
	p.Ret(2, ret0)
}

func execmMIMEHeaderValues(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(textproto.MIMEHeader).Values(args[1].(string))
	p.Ret(2, ret0)
}

func execmMIMEHeaderDel(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(textproto.MIMEHeader).Del(args[1].(string))
	p.PopN(2)
}

func toType0(v interface{}) io.ReadWriteCloser {
	if v == nil {
		return nil
	}
	return v.(io.ReadWriteCloser)
}

func execNewConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := textproto.NewConn(toType0(args[0]))
	p.Ret(1, ret0)
}

func execNewReader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := textproto.NewReader(args[0].(*bufio.Reader))
	p.Ret(1, ret0)
}

func execNewWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := textproto.NewWriter(args[0].(*bufio.Writer))
	p.Ret(1, ret0)
}

func execmPipelineNext(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*textproto.Pipeline).Next()
	p.Ret(1, ret0)
}

func execmPipelineStartRequest(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*textproto.Pipeline).StartRequest(args[1].(uint))
	p.PopN(2)
}

func execmPipelineEndRequest(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*textproto.Pipeline).EndRequest(args[1].(uint))
	p.PopN(2)
}

func execmPipelineStartResponse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*textproto.Pipeline).StartResponse(args[1].(uint))
	p.PopN(2)
}

func execmPipelineEndResponse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*textproto.Pipeline).EndResponse(args[1].(uint))
	p.PopN(2)
}

func execmProtocolErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(textproto.ProtocolError).Error()
	p.Ret(1, ret0)
}

func execmReaderReadLine(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*textproto.Reader).ReadLine()
	p.Ret(1, ret0, ret1)
}

func execmReaderReadLineBytes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*textproto.Reader).ReadLineBytes()
	p.Ret(1, ret0, ret1)
}

func execmReaderReadContinuedLine(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*textproto.Reader).ReadContinuedLine()
	p.Ret(1, ret0, ret1)
}

func execmReaderReadContinuedLineBytes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*textproto.Reader).ReadContinuedLineBytes()
	p.Ret(1, ret0, ret1)
}

func execmReaderReadCodeLine(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(*textproto.Reader).ReadCodeLine(args[1].(int))
	p.Ret(2, ret0, ret1, ret2)
}

func execmReaderReadResponse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(*textproto.Reader).ReadResponse(args[1].(int))
	p.Ret(2, ret0, ret1, ret2)
}

func execmReaderDotReader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*textproto.Reader).DotReader()
	p.Ret(1, ret0)
}

func execmReaderReadDotBytes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*textproto.Reader).ReadDotBytes()
	p.Ret(1, ret0, ret1)
}

func execmReaderReadDotLines(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*textproto.Reader).ReadDotLines()
	p.Ret(1, ret0, ret1)
}

func execmReaderReadMIMEHeader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*textproto.Reader).ReadMIMEHeader()
	p.Ret(1, ret0, ret1)
}

func execTrimBytes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := textproto.TrimBytes(args[0].([]byte))
	p.Ret(1, ret0)
}

func execTrimString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := textproto.TrimString(args[0].(string))
	p.Ret(1, ret0)
}

func execmWriterPrintfLine(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*textproto.Writer).PrintfLine(args[1].(string), args[2:]...)
	p.Ret(arity, ret0)
}

func execmWriterDotWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*textproto.Writer).DotWriter()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net/textproto")

func init() {
	I.RegisterFuncs(
		I.Func("CanonicalMIMEHeaderKey", textproto.CanonicalMIMEHeaderKey, execCanonicalMIMEHeaderKey),
		I.Func("(*Conn).Close", (*textproto.Conn).Close, execmConnClose),
		I.Func("Dial", textproto.Dial, execDial),
		I.Func("(*Error).Error", (*textproto.Error).Error, execmErrorError),
		I.Func("(MIMEHeader).Add", (textproto.MIMEHeader).Add, execmMIMEHeaderAdd),
		I.Func("(MIMEHeader).Set", (textproto.MIMEHeader).Set, execmMIMEHeaderSet),
		I.Func("(MIMEHeader).Get", (textproto.MIMEHeader).Get, execmMIMEHeaderGet),
		I.Func("(MIMEHeader).Values", (textproto.MIMEHeader).Values, execmMIMEHeaderValues),
		I.Func("(MIMEHeader).Del", (textproto.MIMEHeader).Del, execmMIMEHeaderDel),
		I.Func("NewConn", textproto.NewConn, execNewConn),
		I.Func("NewReader", textproto.NewReader, execNewReader),
		I.Func("NewWriter", textproto.NewWriter, execNewWriter),
		I.Func("(*Pipeline).Next", (*textproto.Pipeline).Next, execmPipelineNext),
		I.Func("(*Pipeline).StartRequest", (*textproto.Pipeline).StartRequest, execmPipelineStartRequest),
		I.Func("(*Pipeline).EndRequest", (*textproto.Pipeline).EndRequest, execmPipelineEndRequest),
		I.Func("(*Pipeline).StartResponse", (*textproto.Pipeline).StartResponse, execmPipelineStartResponse),
		I.Func("(*Pipeline).EndResponse", (*textproto.Pipeline).EndResponse, execmPipelineEndResponse),
		I.Func("(ProtocolError).Error", (textproto.ProtocolError).Error, execmProtocolErrorError),
		I.Func("(*Reader).ReadLine", (*textproto.Reader).ReadLine, execmReaderReadLine),
		I.Func("(*Reader).ReadLineBytes", (*textproto.Reader).ReadLineBytes, execmReaderReadLineBytes),
		I.Func("(*Reader).ReadContinuedLine", (*textproto.Reader).ReadContinuedLine, execmReaderReadContinuedLine),
		I.Func("(*Reader).ReadContinuedLineBytes", (*textproto.Reader).ReadContinuedLineBytes, execmReaderReadContinuedLineBytes),
		I.Func("(*Reader).ReadCodeLine", (*textproto.Reader).ReadCodeLine, execmReaderReadCodeLine),
		I.Func("(*Reader).ReadResponse", (*textproto.Reader).ReadResponse, execmReaderReadResponse),
		I.Func("(*Reader).DotReader", (*textproto.Reader).DotReader, execmReaderDotReader),
		I.Func("(*Reader).ReadDotBytes", (*textproto.Reader).ReadDotBytes, execmReaderReadDotBytes),
		I.Func("(*Reader).ReadDotLines", (*textproto.Reader).ReadDotLines, execmReaderReadDotLines),
		I.Func("(*Reader).ReadMIMEHeader", (*textproto.Reader).ReadMIMEHeader, execmReaderReadMIMEHeader),
		I.Func("TrimBytes", textproto.TrimBytes, execTrimBytes),
		I.Func("TrimString", textproto.TrimString, execTrimString),
		I.Func("(*Writer).DotWriter", (*textproto.Writer).DotWriter, execmWriterDotWriter),
	)
	I.RegisterFuncvs(
		I.Funcv("(*Conn).Cmd", (*textproto.Conn).Cmd, execmConnCmd),
		I.Funcv("(*Writer).PrintfLine", (*textproto.Writer).PrintfLine, execmWriterPrintfLine),
	)
	I.RegisterTypes(
		I.Type("Conn", reflect.TypeOf((*textproto.Conn)(nil)).Elem()),
		I.Type("Error", reflect.TypeOf((*textproto.Error)(nil)).Elem()),
		I.Type("MIMEHeader", reflect.TypeOf((*textproto.MIMEHeader)(nil)).Elem()),
		I.Type("Pipeline", reflect.TypeOf((*textproto.Pipeline)(nil)).Elem()),
		I.Type("ProtocolError", reflect.TypeOf((*textproto.ProtocolError)(nil)).Elem()),
		I.Type("Reader", reflect.TypeOf((*textproto.Reader)(nil)).Elem()),
		I.Type("Writer", reflect.TypeOf((*textproto.Writer)(nil)).Elem()),
	)
}
