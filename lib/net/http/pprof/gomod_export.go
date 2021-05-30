// Package pprof provide Go+ "net/http/pprof" package, as "net/http/pprof" package in Go.
package pprof

import (
	http "net/http"
	pprof "net/http/pprof"

	gop "github.com/goplus/gop"
)

func toType0(v interface{}) http.ResponseWriter {
	if v == nil {
		return nil
	}
	return v.(http.ResponseWriter)
}

func execCmdline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	pprof.Cmdline(toType0(args[0]), args[1].(*http.Request))
	p.PopN(2)
}

func execHandler(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := pprof.Handler(args[0].(string))
	p.Ret(1, ret0)
}

func execIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	pprof.Index(toType0(args[0]), args[1].(*http.Request))
	p.PopN(2)
}

func execProfile(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	pprof.Profile(toType0(args[0]), args[1].(*http.Request))
	p.PopN(2)
}

func execSymbol(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	pprof.Symbol(toType0(args[0]), args[1].(*http.Request))
	p.PopN(2)
}

func execTrace(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	pprof.Trace(toType0(args[0]), args[1].(*http.Request))
	p.PopN(2)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net/http/pprof")

func init() {
	I.RegisterFuncs(
		I.Func("Cmdline", pprof.Cmdline, execCmdline),
		I.Func("Handler", pprof.Handler, execHandler),
		I.Func("Index", pprof.Index, execIndex),
		I.Func("Profile", pprof.Profile, execProfile),
		I.Func("Symbol", pprof.Symbol, execSymbol),
		I.Func("Trace", pprof.Trace, execTrace),
	)
}
