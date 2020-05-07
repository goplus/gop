package builtin

import (
	"fmt"
	"io"

	qlang "github.com/qiniu/qlang/spec"
)

// -----------------------------------------------------------------------------

// QexecPrint instr
func QexecPrint(arity uint32, p *qlang.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Print(args...)
	p.Ret(arity, n, err)
}

// QexecPrintf instr
func QexecPrintf(arity uint32, p *qlang.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Printf(args[0].(string), args[1:]...)
	p.Ret(arity, n, err)
}

// QexecPrintln instr
func QexecPrintln(arity uint32, p *qlang.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Println(args...)
	p.Ret(arity, n, err)
}

// QexecFprintln instr
func QexecFprintln(arity uint32, p *qlang.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Fprintln(args[0].(io.Writer), args[1:]...)
	p.Ret(arity, n, err)
}

// -----------------------------------------------------------------------------

// I is a Go package instance.
var I = qlang.NewGoPackage("")

func init() {
	I.RegisterVariadicFuncs(
		I.Func("print", fmt.Print, QexecPrint),
		I.Func("printf", fmt.Printf, QexecPrintf),
		I.Func("println", fmt.Println, QexecPrintln),
		I.Func("fprintln", fmt.Fprintln, QexecFprintln),
	)
}

// -----------------------------------------------------------------------------
