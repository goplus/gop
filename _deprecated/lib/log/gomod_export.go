// Package log provide Go+ "log" package, as "log" package in Go.
package log

import (
	io "io"
	log "log"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execFatal(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	log.Fatal(args...)
	p.PopN(arity)
}

func execFatalf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	log.Fatalf(args[0].(string), args[1:]...)
	p.PopN(arity)
}

func execFatalln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	log.Fatalln(args...)
	p.PopN(arity)
}

func execFlags(_ int, p *gop.Context) {
	ret0 := log.Flags()
	p.Ret(0, ret0)
}

func toType0(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execmLoggerSetOutput(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*log.Logger).SetOutput(toType0(args[1]))
	p.PopN(2)
}

func execmLoggerOutput(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*log.Logger).Output(args[1].(int), args[2].(string))
	p.Ret(3, ret0)
}

func execmLoggerPrintf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	args[0].(*log.Logger).Printf(args[1].(string), args[2:]...)
	p.PopN(arity)
}

func execmLoggerPrint(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	args[0].(*log.Logger).Print(args[1:]...)
	p.PopN(arity)
}

func execmLoggerPrintln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	args[0].(*log.Logger).Println(args[1:]...)
	p.PopN(arity)
}

func execmLoggerFatal(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	args[0].(*log.Logger).Fatal(args[1:]...)
	p.PopN(arity)
}

func execmLoggerFatalf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	args[0].(*log.Logger).Fatalf(args[1].(string), args[2:]...)
	p.PopN(arity)
}

func execmLoggerFatalln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	args[0].(*log.Logger).Fatalln(args[1:]...)
	p.PopN(arity)
}

func execmLoggerPanic(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	args[0].(*log.Logger).Panic(args[1:]...)
	p.PopN(arity)
}

func execmLoggerPanicf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	args[0].(*log.Logger).Panicf(args[1].(string), args[2:]...)
	p.PopN(arity)
}

func execmLoggerPanicln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	args[0].(*log.Logger).Panicln(args[1:]...)
	p.PopN(arity)
}

func execmLoggerFlags(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*log.Logger).Flags()
	p.Ret(1, ret0)
}

func execmLoggerSetFlags(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*log.Logger).SetFlags(args[1].(int))
	p.PopN(2)
}

func execmLoggerPrefix(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*log.Logger).Prefix()
	p.Ret(1, ret0)
}

func execmLoggerSetPrefix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*log.Logger).SetPrefix(args[1].(string))
	p.PopN(2)
}

func execmLoggerWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*log.Logger).Writer()
	p.Ret(1, ret0)
}

func execNew(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := log.New(toType0(args[0]), args[1].(string), args[2].(int))
	p.Ret(3, ret0)
}

func execOutput(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := log.Output(args[0].(int), args[1].(string))
	p.Ret(2, ret0)
}

func execPanic(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	log.Panic(args...)
	p.PopN(arity)
}

func execPanicf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	log.Panicf(args[0].(string), args[1:]...)
	p.PopN(arity)
}

func execPanicln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	log.Panicln(args...)
	p.PopN(arity)
}

func execPrefix(_ int, p *gop.Context) {
	ret0 := log.Prefix()
	p.Ret(0, ret0)
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

func execSetFlags(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	log.SetFlags(args[0].(int))
	p.PopN(1)
}

func execSetOutput(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	log.SetOutput(toType0(args[0]))
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

// I is a Go package instance.
var I = gop.NewGoPackage("log")

func init() {
	I.RegisterFuncs(
		I.Func("Flags", log.Flags, execFlags),
		I.Func("(*Logger).SetOutput", (*log.Logger).SetOutput, execmLoggerSetOutput),
		I.Func("(*Logger).Output", (*log.Logger).Output, execmLoggerOutput),
		I.Func("(*Logger).Flags", (*log.Logger).Flags, execmLoggerFlags),
		I.Func("(*Logger).SetFlags", (*log.Logger).SetFlags, execmLoggerSetFlags),
		I.Func("(*Logger).Prefix", (*log.Logger).Prefix, execmLoggerPrefix),
		I.Func("(*Logger).SetPrefix", (*log.Logger).SetPrefix, execmLoggerSetPrefix),
		I.Func("(*Logger).Writer", (*log.Logger).Writer, execmLoggerWriter),
		I.Func("New", log.New, execNew),
		I.Func("Output", log.Output, execOutput),
		I.Func("Prefix", log.Prefix, execPrefix),
		I.Func("SetFlags", log.SetFlags, execSetFlags),
		I.Func("SetOutput", log.SetOutput, execSetOutput),
		I.Func("SetPrefix", log.SetPrefix, execSetPrefix),
		I.Func("Writer", log.Writer, execWriter),
	)
	I.RegisterFuncvs(
		I.Funcv("Fatal", log.Fatal, execFatal),
		I.Funcv("Fatalf", log.Fatalf, execFatalf),
		I.Funcv("Fatalln", log.Fatalln, execFatalln),
		I.Funcv("(*Logger).Printf", (*log.Logger).Printf, execmLoggerPrintf),
		I.Funcv("(*Logger).Print", (*log.Logger).Print, execmLoggerPrint),
		I.Funcv("(*Logger).Println", (*log.Logger).Println, execmLoggerPrintln),
		I.Funcv("(*Logger).Fatal", (*log.Logger).Fatal, execmLoggerFatal),
		I.Funcv("(*Logger).Fatalf", (*log.Logger).Fatalf, execmLoggerFatalf),
		I.Funcv("(*Logger).Fatalln", (*log.Logger).Fatalln, execmLoggerFatalln),
		I.Funcv("(*Logger).Panic", (*log.Logger).Panic, execmLoggerPanic),
		I.Funcv("(*Logger).Panicf", (*log.Logger).Panicf, execmLoggerPanicf),
		I.Funcv("(*Logger).Panicln", (*log.Logger).Panicln, execmLoggerPanicln),
		I.Funcv("Panic", log.Panic, execPanic),
		I.Funcv("Panicf", log.Panicf, execPanicf),
		I.Funcv("Panicln", log.Panicln, execPanicln),
		I.Funcv("Print", log.Print, execPrint),
		I.Funcv("Printf", log.Printf, execPrintf),
		I.Funcv("Println", log.Println, execPrintln),
	)
	I.RegisterTypes(
		I.Type("Logger", reflect.TypeOf((*log.Logger)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("LUTC", qspec.ConstUnboundInt, log.LUTC),
		I.Const("Ldate", qspec.ConstUnboundInt, log.Ldate),
		I.Const("Llongfile", qspec.ConstUnboundInt, log.Llongfile),
		I.Const("Lmicroseconds", qspec.ConstUnboundInt, log.Lmicroseconds),
		I.Const("Lmsgprefix", qspec.ConstUnboundInt, log.Lmsgprefix),
		I.Const("Lshortfile", qspec.ConstUnboundInt, log.Lshortfile),
		I.Const("LstdFlags", qspec.ConstUnboundInt, log.LstdFlags),
		I.Const("Ltime", qspec.ConstUnboundInt, log.Ltime),
	)
}
