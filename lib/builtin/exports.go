package builtin

import (
	"fmt"
	"io"

	"github.com/qiniu/qlang/spec"
)

// -----------------------------------------------------------------------------

func execPrint(arity uint32, p *spec.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Print(args...)
	p.Ret(arity, n, err)
}

func execPrintf(arity uint32, p *spec.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Printf(args[0].(string), args[1:]...)
	p.Ret(arity, n, err)
}

func execPrintln(arity uint32, p *spec.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Println(args...)
	p.Ret(arity, n, err)
}

func execFprintln(arity uint32, p *spec.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Fprintln(args[0].(io.Writer), args[1:]...)
	p.Ret(arity, n, err)
}

// -----------------------------------------------------------------------------

// I is a Go package instance.
var I = spec.NewPackage("")

func init() {
	I.RegisterVariadicFuncs(
		I.Func("print", fmt.Print, execPrint),
		I.Func("printf", fmt.Printf, execPrintf),
		I.Func("println", fmt.Println, execPrintln),
		I.Func("fprintln", fmt.Fprintln, execFprintln),
	)
}

// -----------------------------------------------------------------------------
