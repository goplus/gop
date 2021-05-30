// Package httptrace provide Go+ "net/http/httptrace" package, as "net/http/httptrace" package in Go.
package httptrace

import (
	context "context"
	httptrace "net/http/httptrace"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func toType0(v interface{}) context.Context {
	if v == nil {
		return nil
	}
	return v.(context.Context)
}

func execContextClientTrace(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := httptrace.ContextClientTrace(toType0(args[0]))
	p.Ret(1, ret0)
}

func execWithClientTrace(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := httptrace.WithClientTrace(toType0(args[0]), args[1].(*httptrace.ClientTrace))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net/http/httptrace")

func init() {
	I.RegisterFuncs(
		I.Func("ContextClientTrace", httptrace.ContextClientTrace, execContextClientTrace),
		I.Func("WithClientTrace", httptrace.WithClientTrace, execWithClientTrace),
	)
	I.RegisterTypes(
		I.Type("ClientTrace", reflect.TypeOf((*httptrace.ClientTrace)(nil)).Elem()),
		I.Type("DNSDoneInfo", reflect.TypeOf((*httptrace.DNSDoneInfo)(nil)).Elem()),
		I.Type("DNSStartInfo", reflect.TypeOf((*httptrace.DNSStartInfo)(nil)).Elem()),
		I.Type("GotConnInfo", reflect.TypeOf((*httptrace.GotConnInfo)(nil)).Elem()),
		I.Type("WroteRequestInfo", reflect.TypeOf((*httptrace.WroteRequestInfo)(nil)).Elem()),
	)
}
