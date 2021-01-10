// Package fmt provide Go+ "fmt" package, as "fmt" package in Go.
package fmt

import (
	fmt "fmt"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execErrorf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := fmt.Errorf(args[0].(string), args[1:]...)
	p.Ret(arity, ret0)
}

func toType0(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execFprint(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Fprint(toType0(args[0]), args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func execFprintf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Fprintf(toType0(args[0]), args[1].(string), args[2:]...)
	p.Ret(arity, ret0, ret1)
}

func execFprintln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Fprintln(toType0(args[0]), args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func toType1(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execFscan(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Fscan(toType1(args[0]), args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func execFscanf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Fscanf(toType1(args[0]), args[1].(string), args[2:]...)
	p.Ret(arity, ret0, ret1)
}

func execFscanln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Fscanln(toType1(args[0]), args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func execPrint(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Print(args...)
	p.Ret(arity, ret0, ret1)
}

func execPrintf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Printf(args[0].(string), args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func execPrintln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Println(args...)
	p.Ret(arity, ret0, ret1)
}

func execScan(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Scan(args...)
	p.Ret(arity, ret0, ret1)
}

func execScanf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Scanf(args[0].(string), args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func execScanln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Scanln(args...)
	p.Ret(arity, ret0, ret1)
}

func execSprint(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := fmt.Sprint(args...)
	p.Ret(arity, ret0)
}

func execSprintf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := fmt.Sprintf(args[0].(string), args[1:]...)
	p.Ret(arity, ret0)
}

func execSprintln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := fmt.Sprintln(args...)
	p.Ret(arity, ret0)
}

func execSscan(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Sscan(args[0].(string), args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func execSscanf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Sscanf(args[0].(string), args[1].(string), args[2:]...)
	p.Ret(arity, ret0, ret1)
}

func execSscanln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Sscanln(args[0].(string), args[1:]...)
	p.Ret(arity, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("fmt")

func init() {
	I.RegisterFuncvs(
		I.Funcv("Errorf", fmt.Errorf, execErrorf),
		I.Funcv("Fprint", fmt.Fprint, execFprint),
		I.Funcv("Fprintf", fmt.Fprintf, execFprintf),
		I.Funcv("Fprintln", fmt.Fprintln, execFprintln),
		I.Funcv("Fscan", fmt.Fscan, execFscan),
		I.Funcv("Fscanf", fmt.Fscanf, execFscanf),
		I.Funcv("Fscanln", fmt.Fscanln, execFscanln),
		I.Funcv("Print", fmt.Print, execPrint),
		I.Funcv("Printf", fmt.Printf, execPrintf),
		I.Funcv("Println", fmt.Println, execPrintln),
		I.Funcv("Scan", fmt.Scan, execScan),
		I.Funcv("Scanf", fmt.Scanf, execScanf),
		I.Funcv("Scanln", fmt.Scanln, execScanln),
		I.Funcv("Sprint", fmt.Sprint, execSprint),
		I.Funcv("Sprintf", fmt.Sprintf, execSprintf),
		I.Funcv("Sprintln", fmt.Sprintln, execSprintln),
		I.Funcv("Sscan", fmt.Sscan, execSscan),
		I.Funcv("Sscanf", fmt.Sscanf, execSscanf),
		I.Funcv("Sscanln", fmt.Sscanln, execSscanln),
	)
	I.RegisterTypes(
		I.Type("Formatter", reflect.TypeOf((*fmt.Formatter)(nil)).Elem()),
		I.Type("GoStringer", reflect.TypeOf((*fmt.GoStringer)(nil)).Elem()),
		I.Type("ScanState", reflect.TypeOf((*fmt.ScanState)(nil)).Elem()),
		I.Type("Scanner", reflect.TypeOf((*fmt.Scanner)(nil)).Elem()),
		I.Type("State", reflect.TypeOf((*fmt.State)(nil)).Elem()),
		I.Type("Stringer", reflect.TypeOf((*fmt.Stringer)(nil)).Elem()),
	)
}
