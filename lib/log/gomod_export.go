// Package strconv provide Go+ "log" package, as "log" package in Go.
package log

import (
	io "io"
	log "log"

	qspec "github.com/goplus/gop/exec.spec"
	gop "github.com/goplus/gop/gop"
)

func execFlags(_ int, p *gop.Context) {
	ret0 := log.Flags()
	p.Ret(0, ret0)
}

func execOutput(arity int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := log.Output(args[0].(int), args[1].(string))
	p.Ret(2, ret0)
}

func execPrefix(_ int, p *gop.Context) {
	ret0 := log.Prefix()
	p.Ret(0, ret0)
}

func execSetFlags(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	log.SetFlags(args[0].(int))
	p.PopN(1)
}

func execSetOutput(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	log.SetOutput(args[0].(io.Writer))
	p.PopN(1)
}

func execSetPrefix(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	log.SetPrefix(args[0].(string))
	p.PopN(1)
}

func execWriter(_ int, p *gop.Context) {
	ret0 := log.Writer()
	p.Ret(0, ret0)
}

func execNew(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := log.New(args[0].(io.Writer), args[1].(string), args[2].(int))
	p.Ret(3, ret0)
}

func execFatal(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	defer p.PopN(arity)
	log.Fatal(args...)
}

func execFatalf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	defer p.PopN(arity)
	log.Fatalf(args[0].(string), args[1:]...)
}

func execFatalln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	defer p.PopN(arity)
	log.Fatalln(args...)
}

func execPanic(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	defer p.PopN(arity)
	log.Panic(args...)
}

func execPanicf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	defer p.PopN(arity)
	log.Panicf(args[0].(string), args[1:]...)
}

func execPanicln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	defer p.PopN(arity)
	log.Panicln(args...)
}

func execPrint(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	log.Print(args...)
	p.PopN(arity)
}

func execPrintf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	log.Printf(args[0].(string), args[1:]...)
	p.PopN(arity)
}

func execPrintln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	log.Println(args...)
	p.PopN(arity)
}

func execLoggerFlags(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*log.Logger).Flags()
	p.Ret(1, ret0)
}

func execLoggerOutput(arity int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*log.Logger).Output(args[1].(int), args[2].(string))
	p.Ret(3, ret0)
}

func execLoggerPrefix(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*log.Logger).Prefix()
	p.Ret(1, ret0)
}

func execLoggerSetFlags(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*log.Logger).SetFlags(args[1].(int))
	p.PopN(2)
}

func execLoggerSetOutput(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*log.Logger).SetOutput(args[1].(io.Writer))
	p.PopN(2)
}

func execLoggerSetPrefix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*log.Logger).SetPrefix(args[1].(string))
	p.PopN(2)
}

func execLoggerWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*log.Logger).Writer()
	p.Ret(1, ret0)
}

func execLoggerFatal(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	defer p.PopN(arity)
	args[0].(*log.Logger).Fatal(args[1:]...)
}

func execLoggerFatalf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	defer p.PopN(arity)
	args[0].(*log.Logger).Fatalf(args[1].(string), args[2:]...)
}

func execLoggerFatalln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	defer p.PopN(arity)
	args[0].(*log.Logger).Fatalln(args[1:]...)
}

func execLoggerPanic(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	defer p.PopN(arity)
	args[0].(*log.Logger).Panic(args[1:]...)
}

func execLoggerPanicf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	defer p.PopN(arity)
	args[0].(*log.Logger).Panicf(args[1].(string), args[2:]...)
}

func execLoggerPanicln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	defer p.PopN(arity)
	args[0].(*log.Logger).Panicln(args[1:]...)
}

func execLoggerPrint(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	args[0].(*log.Logger).Print(args[1:]...)
	p.PopN(arity)
}

func execLoggerPrintf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	args[0].(*log.Logger).Printf(args[1].(string), args[2:]...)
	p.PopN(arity)
}

func execLoggerPrintln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	args[0].(*log.Logger).Println(args[1:]...)
	p.PopN(arity)
}

// I is a Go package instance.
var I = gop.NewGoPackage("log")

func init() {
	I.RegisterFuncs(
		I.Func("Flags", log.Flags, execFlags),
		I.Func("Output", log.Output, execOutput),
		I.Func("Prefix", log.Prefix, execPrefix),
		I.Func("SetFlags", log.SetFlags, execSetFlags),
		I.Func("SetOutput", log.SetOutput, execSetOutput),
		I.Func("SetPrefix", log.SetPrefix, execSetPrefix),
		I.Func("Writer", log.Writer, execWriter),
		I.Func("New", log.New, execNew),
		I.Func("(*Logger).Flags", (*log.Logger).Flags, execLoggerFlags),
		I.Func("(*Logger).Output", (*log.Logger).Output, execLoggerOutput),
		I.Func("(*Logger).Prefix", (*log.Logger).Prefix, execLoggerPrefix),
		I.Func("(*Logger).SetFlags", (*log.Logger).SetFlags, execLoggerSetFlags),
		I.Func("(*Logger).SetOutput", (*log.Logger).SetOutput, execLoggerSetOutput),
		I.Func("(*Logger).SetPrefix", (*log.Logger).SetPrefix, execLoggerSetPrefix),
		I.Func("(*Logger).Writer", (*log.Logger).Writer, execLoggerWriter),
	)
	I.RegisterFuncvs(
		I.Funcv("Fatal", log.Fatal, execFatal),
		I.Funcv("Fatalf", log.Fatalf, execFatalf),
		I.Funcv("Fatalln", log.Fatalln, execFatalln),
		I.Funcv("Panic", log.Panic, execPanic),
		I.Funcv("Panicf", log.Panicf, execPanicf),
		I.Funcv("Panicln", log.Panicln, execPanicln),
		I.Funcv("Print", log.Print, execPrint),
		I.Funcv("Printf", log.Printf, execPrintf),
		I.Funcv("Println", log.Println, execPrintln),
		I.Funcv("(*Logger).Fatal", (*log.Logger).Fatal, execLoggerFatal),
		I.Funcv("(*Logger).Fatalf", (*log.Logger).Fatalf, execLoggerFatalf),
		I.Funcv("(*Logger).Fatalln", (*log.Logger).Fatalln, execLoggerFatalln),
		I.Funcv("(*Logger).Panic", (*log.Logger).Panic, execLoggerPanic),
		I.Funcv("(*Logger).Panicf", (*log.Logger).Panicf, execLoggerPanicf),
		I.Funcv("(*Logger).Panicln", (*log.Logger).Panicln, execLoggerPanicln),
		I.Funcv("(*Logger).Print", (*log.Logger).Print, execLoggerPrint),
		I.Funcv("(*Logger).Printf", (*log.Logger).Printf, execLoggerPrintf),
		I.Funcv("(*Logger).Println", (*log.Logger).Println, execLoggerPrintln),
	)
	I.RegisterConsts(
		I.Const("Ldate", qspec.Int, log.Ldate),
		I.Const("Ltime", qspec.Int, log.Ltime),
		I.Const("Lmicroseconds", qspec.Int, log.Lmicroseconds),
		I.Const("Llongfile", qspec.Int, log.Llongfile),
		I.Const("Lshortfile", qspec.Int, log.Lshortfile),
		I.Const("LUTC", qspec.Int, log.LUTC),
		I.Const("Lmsgprefix", qspec.Int, log.Lmsgprefix),
		I.Const("LstdFlags", qspec.Int, log.LstdFlags),
	)
}
