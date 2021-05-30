// Package signal provide Go+ "os/signal" package, as "os/signal" package in Go.
package signal

import (
	os "os"
	signal "os/signal"

	gop "github.com/goplus/gop"
)

func toType0(v interface{}) os.Signal {
	if v == nil {
		return nil
	}
	return v.(os.Signal)
}

func toSlice0(args []interface{}) []os.Signal {
	ret := make([]os.Signal, len(args))
	for i, arg := range args {
		ret[i] = toType0(arg)
	}
	return ret
}

func execIgnore(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	signal.Ignore(toSlice0(args)...)
	p.PopN(arity)
}

func execIgnored(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := signal.Ignored(toType0(args[0]))
	p.Ret(1, ret0)
}

func execNotify(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	signal.Notify(args[0].(chan<- os.Signal), toSlice0(args[1:])...)
	p.PopN(arity)
}

func execReset(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	signal.Reset(toSlice0(args)...)
	p.PopN(arity)
}

func execStop(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	signal.Stop(args[0].(chan<- os.Signal))
	p.PopN(1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("os/signal")

func init() {
	I.RegisterFuncs(
		I.Func("Ignored", signal.Ignored, execIgnored),
		I.Func("Stop", signal.Stop, execStop),
	)
	I.RegisterFuncvs(
		I.Funcv("Ignore", signal.Ignore, execIgnore),
		I.Funcv("Notify", signal.Notify, execNotify),
		I.Funcv("Reset", signal.Reset, execReset),
	)
}
