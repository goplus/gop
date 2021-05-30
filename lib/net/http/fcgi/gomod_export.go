// Package fcgi provide Go+ "net/http/fcgi" package, as "net/http/fcgi" package in Go.
package fcgi

import (
	net "net"
	http "net/http"
	fcgi "net/http/fcgi"

	gop "github.com/goplus/gop"
)

func execProcessEnv(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := fcgi.ProcessEnv(args[0].(*http.Request))
	p.Ret(1, ret0)
}

func toType0(v interface{}) net.Listener {
	if v == nil {
		return nil
	}
	return v.(net.Listener)
}

func toType1(v interface{}) http.Handler {
	if v == nil {
		return nil
	}
	return v.(http.Handler)
}

func execServe(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := fcgi.Serve(toType0(args[0]), toType1(args[1]))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net/http/fcgi")

func init() {
	I.RegisterFuncs(
		I.Func("ProcessEnv", fcgi.ProcessEnv, execProcessEnv),
		I.Func("Serve", fcgi.Serve, execServe),
	)
	I.RegisterVars(
		I.Var("ErrConnClosed", &fcgi.ErrConnClosed),
		I.Var("ErrRequestAborted", &fcgi.ErrRequestAborted),
	)
}
