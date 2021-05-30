// Package httptest provide Go+ "net/http/httptest" package, as "net/http/httptest" package in Go.
package httptest

import (
	io "io"
	http "net/http"
	httptest "net/http/httptest"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execNewRecorder(_ int, p *gop.Context) {
	ret0 := httptest.NewRecorder()
	p.Ret(0, ret0)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewRequest(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := httptest.NewRequest(args[0].(string), args[1].(string), toType0(args[2]))
	p.Ret(3, ret0)
}

func toType1(v interface{}) http.Handler {
	if v == nil {
		return nil
	}
	return v.(http.Handler)
}

func execNewServer(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := httptest.NewServer(toType1(args[0]))
	p.Ret(1, ret0)
}

func execNewTLSServer(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := httptest.NewTLSServer(toType1(args[0]))
	p.Ret(1, ret0)
}

func execNewUnstartedServer(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := httptest.NewUnstartedServer(toType1(args[0]))
	p.Ret(1, ret0)
}

func execmResponseRecorderHeader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*httptest.ResponseRecorder).Header()
	p.Ret(1, ret0)
}

func execmResponseRecorderWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*httptest.ResponseRecorder).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmResponseRecorderWriteString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*httptest.ResponseRecorder).WriteString(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmResponseRecorderWriteHeader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*httptest.ResponseRecorder).WriteHeader(args[1].(int))
	p.PopN(2)
}

func execmResponseRecorderFlush(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*httptest.ResponseRecorder).Flush()
	p.PopN(1)
}

func execmResponseRecorderResult(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*httptest.ResponseRecorder).Result()
	p.Ret(1, ret0)
}

func execmServerStart(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*httptest.Server).Start()
	p.PopN(1)
}

func execmServerStartTLS(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*httptest.Server).StartTLS()
	p.PopN(1)
}

func execmServerClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*httptest.Server).Close()
	p.PopN(1)
}

func execmServerCloseClientConnections(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*httptest.Server).CloseClientConnections()
	p.PopN(1)
}

func execmServerCertificate(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*httptest.Server).Certificate()
	p.Ret(1, ret0)
}

func execmServerClient(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*httptest.Server).Client()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net/http/httptest")

func init() {
	I.RegisterFuncs(
		I.Func("NewRecorder", httptest.NewRecorder, execNewRecorder),
		I.Func("NewRequest", httptest.NewRequest, execNewRequest),
		I.Func("NewServer", httptest.NewServer, execNewServer),
		I.Func("NewTLSServer", httptest.NewTLSServer, execNewTLSServer),
		I.Func("NewUnstartedServer", httptest.NewUnstartedServer, execNewUnstartedServer),
		I.Func("(*ResponseRecorder).Header", (*httptest.ResponseRecorder).Header, execmResponseRecorderHeader),
		I.Func("(*ResponseRecorder).Write", (*httptest.ResponseRecorder).Write, execmResponseRecorderWrite),
		I.Func("(*ResponseRecorder).WriteString", (*httptest.ResponseRecorder).WriteString, execmResponseRecorderWriteString),
		I.Func("(*ResponseRecorder).WriteHeader", (*httptest.ResponseRecorder).WriteHeader, execmResponseRecorderWriteHeader),
		I.Func("(*ResponseRecorder).Flush", (*httptest.ResponseRecorder).Flush, execmResponseRecorderFlush),
		I.Func("(*ResponseRecorder).Result", (*httptest.ResponseRecorder).Result, execmResponseRecorderResult),
		I.Func("(*Server).Start", (*httptest.Server).Start, execmServerStart),
		I.Func("(*Server).StartTLS", (*httptest.Server).StartTLS, execmServerStartTLS),
		I.Func("(*Server).Close", (*httptest.Server).Close, execmServerClose),
		I.Func("(*Server).CloseClientConnections", (*httptest.Server).CloseClientConnections, execmServerCloseClientConnections),
		I.Func("(*Server).Certificate", (*httptest.Server).Certificate, execmServerCertificate),
		I.Func("(*Server).Client", (*httptest.Server).Client, execmServerClient),
	)
	I.RegisterTypes(
		I.Type("ResponseRecorder", reflect.TypeOf((*httptest.ResponseRecorder)(nil)).Elem()),
		I.Type("Server", reflect.TypeOf((*httptest.Server)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("DefaultRemoteAddr", qspec.ConstBoundString, httptest.DefaultRemoteAddr),
	)
}
