// Package path provide Go+ "path" package, as "path" package in Go.
package path

import (
	path "path"

	gop "github.com/goplus/gop"
)

func execBase(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := path.Base(args[0].(string))
	p.Ret(1, ret0)
}

func execClean(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := path.Clean(args[0].(string))
	p.Ret(1, ret0)
}

func execDir(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := path.Dir(args[0].(string))
	p.Ret(1, ret0)
}

func execExt(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := path.Ext(args[0].(string))
	p.Ret(1, ret0)
}

func execIsAbs(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := path.IsAbs(args[0].(string))
	p.Ret(1, ret0)
}

func execJoin(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := path.Join(gop.ToStrings(args)...)
	p.Ret(arity, ret0)
}

func execMatch(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := path.Match(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execSplit(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := path.Split(args[0].(string))
	p.Ret(1, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("path")

func init() {
	I.RegisterFuncs(
		I.Func("Base", path.Base, execBase),
		I.Func("Clean", path.Clean, execClean),
		I.Func("Dir", path.Dir, execDir),
		I.Func("Ext", path.Ext, execExt),
		I.Func("IsAbs", path.IsAbs, execIsAbs),
		I.Func("Match", path.Match, execMatch),
		I.Func("Split", path.Split, execSplit),
	)
	I.RegisterFuncvs(
		I.Funcv("Join", path.Join, execJoin),
	)
	I.RegisterVars(
		I.Var("ErrBadPattern", &path.ErrBadPattern),
	)
}
