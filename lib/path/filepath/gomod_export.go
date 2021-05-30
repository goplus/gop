// Package filepath provide Go+ "path/filepath" package, as "path/filepath" package in Go.
package filepath

import (
	filepath "path/filepath"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execAbs(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := filepath.Abs(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execBase(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := filepath.Base(args[0].(string))
	p.Ret(1, ret0)
}

func execClean(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := filepath.Clean(args[0].(string))
	p.Ret(1, ret0)
}

func execDir(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := filepath.Dir(args[0].(string))
	p.Ret(1, ret0)
}

func execEvalSymlinks(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := filepath.EvalSymlinks(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execExt(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := filepath.Ext(args[0].(string))
	p.Ret(1, ret0)
}

func execFromSlash(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := filepath.FromSlash(args[0].(string))
	p.Ret(1, ret0)
}

func execGlob(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := filepath.Glob(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execHasPrefix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := filepath.HasPrefix(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execIsAbs(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := filepath.IsAbs(args[0].(string))
	p.Ret(1, ret0)
}

func execJoin(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := filepath.Join(gop.ToStrings(args)...)
	p.Ret(arity, ret0)
}

func execMatch(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := filepath.Match(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execRel(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := filepath.Rel(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execSplit(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := filepath.Split(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execSplitList(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := filepath.SplitList(args[0].(string))
	p.Ret(1, ret0)
}

func execToSlash(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := filepath.ToSlash(args[0].(string))
	p.Ret(1, ret0)
}

func execVolumeName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := filepath.VolumeName(args[0].(string))
	p.Ret(1, ret0)
}

func execWalk(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := filepath.Walk(args[0].(string), args[1].(filepath.WalkFunc))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("path/filepath")

func init() {
	I.RegisterFuncs(
		I.Func("Abs", filepath.Abs, execAbs),
		I.Func("Base", filepath.Base, execBase),
		I.Func("Clean", filepath.Clean, execClean),
		I.Func("Dir", filepath.Dir, execDir),
		I.Func("EvalSymlinks", filepath.EvalSymlinks, execEvalSymlinks),
		I.Func("Ext", filepath.Ext, execExt),
		I.Func("FromSlash", filepath.FromSlash, execFromSlash),
		I.Func("Glob", filepath.Glob, execGlob),
		I.Func("HasPrefix", filepath.HasPrefix, execHasPrefix),
		I.Func("IsAbs", filepath.IsAbs, execIsAbs),
		I.Func("Match", filepath.Match, execMatch),
		I.Func("Rel", filepath.Rel, execRel),
		I.Func("Split", filepath.Split, execSplit),
		I.Func("SplitList", filepath.SplitList, execSplitList),
		I.Func("ToSlash", filepath.ToSlash, execToSlash),
		I.Func("VolumeName", filepath.VolumeName, execVolumeName),
		I.Func("Walk", filepath.Walk, execWalk),
	)
	I.RegisterFuncvs(
		I.Funcv("Join", filepath.Join, execJoin),
	)
	I.RegisterVars(
		I.Var("ErrBadPattern", &filepath.ErrBadPattern),
		I.Var("SkipDir", &filepath.SkipDir),
	)
	I.RegisterTypes(
		I.Type("WalkFunc", reflect.TypeOf((*filepath.WalkFunc)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("ListSeparator", qspec.ConstBoundRune, filepath.ListSeparator),
		I.Const("Separator", qspec.ConstBoundRune, filepath.Separator),
	)
}
