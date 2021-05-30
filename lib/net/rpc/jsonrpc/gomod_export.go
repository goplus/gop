// Package jsonrpc provide Go+ "net/rpc/jsonrpc" package, as "net/rpc/jsonrpc" package in Go.
package jsonrpc

import (
	io "io"
	jsonrpc "net/rpc/jsonrpc"

	gop "github.com/goplus/gop"
)

func execDial(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := jsonrpc.Dial(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func toType0(v interface{}) io.ReadWriteCloser {
	if v == nil {
		return nil
	}
	return v.(io.ReadWriteCloser)
}

func execNewClient(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := jsonrpc.NewClient(toType0(args[0]))
	p.Ret(1, ret0)
}

func execNewClientCodec(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := jsonrpc.NewClientCodec(toType0(args[0]))
	p.Ret(1, ret0)
}

func execNewServerCodec(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := jsonrpc.NewServerCodec(toType0(args[0]))
	p.Ret(1, ret0)
}

func execServeConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	jsonrpc.ServeConn(toType0(args[0]))
	p.PopN(1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net/rpc/jsonrpc")

func init() {
	I.RegisterFuncs(
		I.Func("Dial", jsonrpc.Dial, execDial),
		I.Func("NewClient", jsonrpc.NewClient, execNewClient),
		I.Func("NewClientCodec", jsonrpc.NewClientCodec, execNewClientCodec),
		I.Func("NewServerCodec", jsonrpc.NewServerCodec, execNewServerCodec),
		I.Func("ServeConn", jsonrpc.ServeConn, execServeConn),
	)
}
