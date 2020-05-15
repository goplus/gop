package strings

import (
	"strings"

	qlang "github.com/qiniu/qlang/v6/spec"
)

// -----------------------------------------------------------------------------

func execNewReplacer(arity uint32, p *qlang.Context) {
	args := p.GetArgs(arity)
	repl := strings.NewReplacer(qlang.ToStrings(args)...)
	p.Ret(arity, repl)
}

func execReplacerReplace(zero uint32, p *qlang.Context) {
	args := p.GetArgs(2)
	ret := args[0].(*strings.Replacer).Replace(args[1].(string))
	p.Ret(2, ret)
}

// -----------------------------------------------------------------------------

// I is a Go package instance.
var I = qlang.NewGoPackage("strings")

func init() {
	I.RegisterFuncvs(
		I.Funcv("NewReplacer", strings.NewReplacer, execNewReplacer),
	)
	I.RegisterFuncs(
		I.Func("(*Replacer).Replace", (*strings.Replacer).Replace, execReplacerReplace),
	)
}

// -----------------------------------------------------------------------------
