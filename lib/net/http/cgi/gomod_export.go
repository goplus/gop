// Package cgi provide Go+ "net/http/cgi" package, as "net/http/cgi" package in Go.
package cgi

import (
	http "net/http"
	cgi "net/http/cgi"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func toType0(v interface{}) http.ResponseWriter {
	if v == nil {
		return nil
	}
	return v.(http.ResponseWriter)
}

func execmHandlerServeHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*cgi.Handler).ServeHTTP(toType0(args[1]), args[2].(*http.Request))
	p.PopN(3)
}

func execRequest(_ int, p *gop.Context) {
	ret0, ret1 := cgi.Request()
	p.Ret(0, ret0, ret1)
}

func execRequestFromMap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := cgi.RequestFromMap(args[0].(map[string]string))
	p.Ret(1, ret0, ret1)
}

func toType1(v interface{}) http.Handler {
	if v == nil {
		return nil
	}
	return v.(http.Handler)
}

func execServe(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := cgi.Serve(toType1(args[0]))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net/http/cgi")

func init() {
	I.RegisterFuncs(
		I.Func("(*Handler).ServeHTTP", (*cgi.Handler).ServeHTTP, execmHandlerServeHTTP),
		I.Func("Request", cgi.Request, execRequest),
		I.Func("RequestFromMap", cgi.RequestFromMap, execRequestFromMap),
		I.Func("Serve", cgi.Serve, execServe),
	)
	I.RegisterTypes(
		I.Type("Handler", reflect.TypeOf((*cgi.Handler)(nil)).Elem()),
	)
}
