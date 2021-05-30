// Package rpc provide Go+ "net/rpc" package, as "net/rpc" package in Go.
package rpc

import (
	io "io"
	net "net"
	http "net/http"
	rpc "net/rpc"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func toType0(v interface{}) net.Listener {
	if v == nil {
		return nil
	}
	return v.(net.Listener)
}

func execAccept(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	rpc.Accept(toType0(args[0]))
	p.PopN(1)
}

func execmClientClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rpc.Client).Close()
	p.Ret(1, ret0)
}

func execmClientGo(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0 := args[0].(*rpc.Client).Go(args[1].(string), args[2], args[3], args[4].(chan *rpc.Call))
	p.Ret(5, ret0)
}

func execmClientCall(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*rpc.Client).Call(args[1].(string), args[2], args[3])
	p.Ret(4, ret0)
}

func execiClientCodecClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(rpc.ClientCodec).Close()
	p.Ret(1, ret0)
}

func execiClientCodecReadResponseBody(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(rpc.ClientCodec).ReadResponseBody(args[1])
	p.Ret(2, ret0)
}

func execiClientCodecReadResponseHeader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(rpc.ClientCodec).ReadResponseHeader(args[1].(*rpc.Response))
	p.Ret(2, ret0)
}

func execiClientCodecWriteRequest(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(rpc.ClientCodec).WriteRequest(args[1].(*rpc.Request), args[2])
	p.Ret(3, ret0)
}

func execDial(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := rpc.Dial(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execDialHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := rpc.DialHTTP(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execDialHTTPPath(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := rpc.DialHTTPPath(args[0].(string), args[1].(string), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execHandleHTTP(_ int, p *gop.Context) {
	rpc.HandleHTTP()
	p.PopN(0)
}

func toType1(v interface{}) io.ReadWriteCloser {
	if v == nil {
		return nil
	}
	return v.(io.ReadWriteCloser)
}

func execNewClient(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := rpc.NewClient(toType1(args[0]))
	p.Ret(1, ret0)
}

func toType2(v interface{}) rpc.ClientCodec {
	if v == nil {
		return nil
	}
	return v.(rpc.ClientCodec)
}

func execNewClientWithCodec(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := rpc.NewClientWithCodec(toType2(args[0]))
	p.Ret(1, ret0)
}

func execNewServer(_ int, p *gop.Context) {
	ret0 := rpc.NewServer()
	p.Ret(0, ret0)
}

func execRegister(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := rpc.Register(args[0])
	p.Ret(1, ret0)
}

func execRegisterName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := rpc.RegisterName(args[0].(string), args[1])
	p.Ret(2, ret0)
}

func toType3(v interface{}) rpc.ServerCodec {
	if v == nil {
		return nil
	}
	return v.(rpc.ServerCodec)
}

func execServeCodec(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	rpc.ServeCodec(toType3(args[0]))
	p.PopN(1)
}

func execServeConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	rpc.ServeConn(toType1(args[0]))
	p.PopN(1)
}

func execServeRequest(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := rpc.ServeRequest(toType3(args[0]))
	p.Ret(1, ret0)
}

func execmServerRegister(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*rpc.Server).Register(args[1])
	p.Ret(2, ret0)
}

func execmServerRegisterName(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*rpc.Server).RegisterName(args[1].(string), args[2])
	p.Ret(3, ret0)
}

func execmServerServeConn(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*rpc.Server).ServeConn(toType1(args[1]))
	p.PopN(2)
}

func execmServerServeCodec(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*rpc.Server).ServeCodec(toType3(args[1]))
	p.PopN(2)
}

func execmServerServeRequest(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*rpc.Server).ServeRequest(toType3(args[1]))
	p.Ret(2, ret0)
}

func execmServerAccept(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*rpc.Server).Accept(toType0(args[1]))
	p.PopN(2)
}

func toType4(v interface{}) http.ResponseWriter {
	if v == nil {
		return nil
	}
	return v.(http.ResponseWriter)
}

func execmServerServeHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*rpc.Server).ServeHTTP(toType4(args[1]), args[2].(*http.Request))
	p.PopN(3)
}

func execmServerHandleHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*rpc.Server).HandleHTTP(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execiServerCodecClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(rpc.ServerCodec).Close()
	p.Ret(1, ret0)
}

func execiServerCodecReadRequestBody(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(rpc.ServerCodec).ReadRequestBody(args[1])
	p.Ret(2, ret0)
}

func execiServerCodecReadRequestHeader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(rpc.ServerCodec).ReadRequestHeader(args[1].(*rpc.Request))
	p.Ret(2, ret0)
}

func execiServerCodecWriteResponse(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(rpc.ServerCodec).WriteResponse(args[1].(*rpc.Response), args[2])
	p.Ret(3, ret0)
}

func execmServerErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(rpc.ServerError).Error()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net/rpc")

func init() {
	I.RegisterFuncs(
		I.Func("Accept", rpc.Accept, execAccept),
		I.Func("(*Client).Close", (*rpc.Client).Close, execmClientClose),
		I.Func("(*Client).Go", (*rpc.Client).Go, execmClientGo),
		I.Func("(*Client).Call", (*rpc.Client).Call, execmClientCall),
		I.Func("(ClientCodec).Close", (rpc.ClientCodec).Close, execiClientCodecClose),
		I.Func("(ClientCodec).ReadResponseBody", (rpc.ClientCodec).ReadResponseBody, execiClientCodecReadResponseBody),
		I.Func("(ClientCodec).ReadResponseHeader", (rpc.ClientCodec).ReadResponseHeader, execiClientCodecReadResponseHeader),
		I.Func("(ClientCodec).WriteRequest", (rpc.ClientCodec).WriteRequest, execiClientCodecWriteRequest),
		I.Func("Dial", rpc.Dial, execDial),
		I.Func("DialHTTP", rpc.DialHTTP, execDialHTTP),
		I.Func("DialHTTPPath", rpc.DialHTTPPath, execDialHTTPPath),
		I.Func("HandleHTTP", rpc.HandleHTTP, execHandleHTTP),
		I.Func("NewClient", rpc.NewClient, execNewClient),
		I.Func("NewClientWithCodec", rpc.NewClientWithCodec, execNewClientWithCodec),
		I.Func("NewServer", rpc.NewServer, execNewServer),
		I.Func("Register", rpc.Register, execRegister),
		I.Func("RegisterName", rpc.RegisterName, execRegisterName),
		I.Func("ServeCodec", rpc.ServeCodec, execServeCodec),
		I.Func("ServeConn", rpc.ServeConn, execServeConn),
		I.Func("ServeRequest", rpc.ServeRequest, execServeRequest),
		I.Func("(*Server).Register", (*rpc.Server).Register, execmServerRegister),
		I.Func("(*Server).RegisterName", (*rpc.Server).RegisterName, execmServerRegisterName),
		I.Func("(*Server).ServeConn", (*rpc.Server).ServeConn, execmServerServeConn),
		I.Func("(*Server).ServeCodec", (*rpc.Server).ServeCodec, execmServerServeCodec),
		I.Func("(*Server).ServeRequest", (*rpc.Server).ServeRequest, execmServerServeRequest),
		I.Func("(*Server).Accept", (*rpc.Server).Accept, execmServerAccept),
		I.Func("(*Server).ServeHTTP", (*rpc.Server).ServeHTTP, execmServerServeHTTP),
		I.Func("(*Server).HandleHTTP", (*rpc.Server).HandleHTTP, execmServerHandleHTTP),
		I.Func("(ServerCodec).Close", (rpc.ServerCodec).Close, execiServerCodecClose),
		I.Func("(ServerCodec).ReadRequestBody", (rpc.ServerCodec).ReadRequestBody, execiServerCodecReadRequestBody),
		I.Func("(ServerCodec).ReadRequestHeader", (rpc.ServerCodec).ReadRequestHeader, execiServerCodecReadRequestHeader),
		I.Func("(ServerCodec).WriteResponse", (rpc.ServerCodec).WriteResponse, execiServerCodecWriteResponse),
		I.Func("(ServerError).Error", (rpc.ServerError).Error, execmServerErrorError),
	)
	I.RegisterVars(
		I.Var("DefaultServer", &rpc.DefaultServer),
		I.Var("ErrShutdown", &rpc.ErrShutdown),
	)
	I.RegisterTypes(
		I.Type("Call", reflect.TypeOf((*rpc.Call)(nil)).Elem()),
		I.Type("Client", reflect.TypeOf((*rpc.Client)(nil)).Elem()),
		I.Type("ClientCodec", reflect.TypeOf((*rpc.ClientCodec)(nil)).Elem()),
		I.Type("Request", reflect.TypeOf((*rpc.Request)(nil)).Elem()),
		I.Type("Response", reflect.TypeOf((*rpc.Response)(nil)).Elem()),
		I.Type("Server", reflect.TypeOf((*rpc.Server)(nil)).Elem()),
		I.Type("ServerCodec", reflect.TypeOf((*rpc.ServerCodec)(nil)).Elem()),
		I.Type("ServerError", reflect.TypeOf((*rpc.ServerError)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("DefaultDebugPath", qspec.ConstBoundString, rpc.DefaultDebugPath),
		I.Const("DefaultRPCPath", qspec.ConstBoundString, rpc.DefaultRPCPath),
	)
}
