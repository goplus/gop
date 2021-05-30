// Package httputil provide Go+ "net/http/httputil" package, as "net/http/httputil" package in Go.
package httputil

import (
	bufio "bufio"
	io "io"
	net "net"
	http "net/http"
	httputil "net/http/httputil"
	url "net/url"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execiBufferPoolGet(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(httputil.BufferPool).Get()
	p.Ret(1, ret0)
}

func execiBufferPoolPut(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(httputil.BufferPool).Put(args[1].([]byte))
	p.PopN(2)
}

func execmClientConnHijack(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*httputil.ClientConn).Hijack()
	p.Ret(1, ret0, ret1)
}

func execmClientConnClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*httputil.ClientConn).Close()
	p.Ret(1, ret0)
}

func execmClientConnWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*httputil.ClientConn).Write(args[1].(*http.Request))
	p.Ret(2, ret0)
}

func execmClientConnPending(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*httputil.ClientConn).Pending()
	p.Ret(1, ret0)
}

func execmClientConnRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*httputil.ClientConn).Read(args[1].(*http.Request))
	p.Ret(2, ret0, ret1)
}

func execmClientConnDo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*httputil.ClientConn).Do(args[1].(*http.Request))
	p.Ret(2, ret0, ret1)
}

func execDumpRequest(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := httputil.DumpRequest(args[0].(*http.Request), args[1].(bool))
	p.Ret(2, ret0, ret1)
}

func execDumpRequestOut(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := httputil.DumpRequestOut(args[0].(*http.Request), args[1].(bool))
	p.Ret(2, ret0, ret1)
}

func execDumpResponse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := httputil.DumpResponse(args[0].(*http.Response), args[1].(bool))
	p.Ret(2, ret0, ret1)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewChunkedReader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := httputil.NewChunkedReader(toType0(args[0]))
	p.Ret(1, ret0)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execNewChunkedWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := httputil.NewChunkedWriter(toType1(args[0]))
	p.Ret(1, ret0)
}

func toType2(v interface{}) net.Conn {
	if v == nil {
		return nil
	}
	return v.(net.Conn)
}

func execNewClientConn(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := httputil.NewClientConn(toType2(args[0]), args[1].(*bufio.Reader))
	p.Ret(2, ret0)
}

func execNewProxyClientConn(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := httputil.NewProxyClientConn(toType2(args[0]), args[1].(*bufio.Reader))
	p.Ret(2, ret0)
}

func execNewServerConn(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := httputil.NewServerConn(toType2(args[0]), args[1].(*bufio.Reader))
	p.Ret(2, ret0)
}

func execNewSingleHostReverseProxy(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := httputil.NewSingleHostReverseProxy(args[0].(*url.URL))
	p.Ret(1, ret0)
}

func toType3(v interface{}) http.ResponseWriter {
	if v == nil {
		return nil
	}
	return v.(http.ResponseWriter)
}

func execmReverseProxyServeHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*httputil.ReverseProxy).ServeHTTP(toType3(args[1]), args[2].(*http.Request))
	p.PopN(3)
}

func execmServerConnHijack(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*httputil.ServerConn).Hijack()
	p.Ret(1, ret0, ret1)
}

func execmServerConnClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*httputil.ServerConn).Close()
	p.Ret(1, ret0)
}

func execmServerConnRead(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*httputil.ServerConn).Read()
	p.Ret(1, ret0, ret1)
}

func execmServerConnPending(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*httputil.ServerConn).Pending()
	p.Ret(1, ret0)
}

func execmServerConnWrite(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*httputil.ServerConn).Write(args[1].(*http.Request), args[2].(*http.Response))
	p.Ret(3, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net/http/httputil")

func init() {
	I.RegisterFuncs(
		I.Func("(BufferPool).Get", (httputil.BufferPool).Get, execiBufferPoolGet),
		I.Func("(BufferPool).Put", (httputil.BufferPool).Put, execiBufferPoolPut),
		I.Func("(*ClientConn).Hijack", (*httputil.ClientConn).Hijack, execmClientConnHijack),
		I.Func("(*ClientConn).Close", (*httputil.ClientConn).Close, execmClientConnClose),
		I.Func("(*ClientConn).Write", (*httputil.ClientConn).Write, execmClientConnWrite),
		I.Func("(*ClientConn).Pending", (*httputil.ClientConn).Pending, execmClientConnPending),
		I.Func("(*ClientConn).Read", (*httputil.ClientConn).Read, execmClientConnRead),
		I.Func("(*ClientConn).Do", (*httputil.ClientConn).Do, execmClientConnDo),
		I.Func("DumpRequest", httputil.DumpRequest, execDumpRequest),
		I.Func("DumpRequestOut", httputil.DumpRequestOut, execDumpRequestOut),
		I.Func("DumpResponse", httputil.DumpResponse, execDumpResponse),
		I.Func("NewChunkedReader", httputil.NewChunkedReader, execNewChunkedReader),
		I.Func("NewChunkedWriter", httputil.NewChunkedWriter, execNewChunkedWriter),
		I.Func("NewClientConn", httputil.NewClientConn, execNewClientConn),
		I.Func("NewProxyClientConn", httputil.NewProxyClientConn, execNewProxyClientConn),
		I.Func("NewServerConn", httputil.NewServerConn, execNewServerConn),
		I.Func("NewSingleHostReverseProxy", httputil.NewSingleHostReverseProxy, execNewSingleHostReverseProxy),
		I.Func("(*ReverseProxy).ServeHTTP", (*httputil.ReverseProxy).ServeHTTP, execmReverseProxyServeHTTP),
		I.Func("(*ServerConn).Hijack", (*httputil.ServerConn).Hijack, execmServerConnHijack),
		I.Func("(*ServerConn).Close", (*httputil.ServerConn).Close, execmServerConnClose),
		I.Func("(*ServerConn).Read", (*httputil.ServerConn).Read, execmServerConnRead),
		I.Func("(*ServerConn).Pending", (*httputil.ServerConn).Pending, execmServerConnPending),
		I.Func("(*ServerConn).Write", (*httputil.ServerConn).Write, execmServerConnWrite),
	)
	I.RegisterVars(
		I.Var("ErrClosed", &httputil.ErrClosed),
		I.Var("ErrLineTooLong", &httputil.ErrLineTooLong),
		I.Var("ErrPersistEOF", &httputil.ErrPersistEOF),
		I.Var("ErrPipeline", &httputil.ErrPipeline),
	)
	I.RegisterTypes(
		I.Type("BufferPool", reflect.TypeOf((*httputil.BufferPool)(nil)).Elem()),
		I.Type("ClientConn", reflect.TypeOf((*httputil.ClientConn)(nil)).Elem()),
		I.Type("ReverseProxy", reflect.TypeOf((*httputil.ReverseProxy)(nil)).Elem()),
		I.Type("ServerConn", reflect.TypeOf((*httputil.ServerConn)(nil)).Elem()),
	)
}
