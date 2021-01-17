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

func toType0(v interface{}) fmt.State {
	if v == nil {
		return nil
	}
	return v.(fmt.State)
}

func execiFormatterFormat(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(fmt.Formatter).Format(toType0(args[1]), args[2].(rune))
	p.PopN(3)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execFprint(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Fprint(toType1(args[0]), args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func execFprintf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Fprintf(toType1(args[0]), args[1].(string), args[2:]...)
	p.Ret(arity, ret0, ret1)
}

func execFprintln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Fprintln(toType1(args[0]), args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func toType2(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execFscan(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Fscan(toType2(args[0]), args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func execFscanf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Fscanf(toType2(args[0]), args[1].(string), args[2:]...)
	p.Ret(arity, ret0, ret1)
}

func execFscanln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := fmt.Fscanln(toType2(args[0]), args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func execiGoStringerGoString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fmt.GoStringer).GoString()
	p.Ret(1, ret0)
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

func execiScanStateRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(fmt.ScanState).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiScanStateReadRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(fmt.ScanState).ReadRune()
	p.Ret(1, ret0, ret1, ret2)
}

func execiScanStateSkipSpace(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(fmt.ScanState).SkipSpace()
	p.PopN(1)
}

func execiScanStateToken(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(fmt.ScanState).Token(args[1].(bool), args[2].(func(rune) bool))
	p.Ret(3, ret0, ret1)
}

func execiScanStateUnreadRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fmt.ScanState).UnreadRune()
	p.Ret(1, ret0)
}

func execiScanStateWidth(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(fmt.ScanState).Width()
	p.Ret(1, ret0, ret1)
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

func toType3(v interface{}) fmt.ScanState {
	if v == nil {
		return nil
	}
	return v.(fmt.ScanState)
}

func execiScannerScan(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(fmt.Scanner).Scan(toType3(args[1]), args[2].(rune))
	p.Ret(3, ret0)
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

func execiStateFlag(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(fmt.State).Flag(args[1].(int))
	p.Ret(2, ret0)
}

func execiStatePrecision(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(fmt.State).Precision()
	p.Ret(1, ret0, ret1)
}

func execiStateWidth(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(fmt.State).Width()
	p.Ret(1, ret0, ret1)
}

func execiStateWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(fmt.State).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiStringerString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fmt.Stringer).String()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("fmt")

func init() {
	I.RegisterFuncs(
		I.Func("(Formatter).Format", (fmt.Formatter).Format, execiFormatterFormat),
		I.Func("(GoStringer).GoString", (fmt.GoStringer).GoString, execiGoStringerGoString),
		I.Func("(ScanState).Read", (fmt.ScanState).Read, execiScanStateRead),
		I.Func("(ScanState).ReadRune", (fmt.ScanState).ReadRune, execiScanStateReadRune),
		I.Func("(ScanState).SkipSpace", (fmt.ScanState).SkipSpace, execiScanStateSkipSpace),
		I.Func("(ScanState).Token", (fmt.ScanState).Token, execiScanStateToken),
		I.Func("(ScanState).UnreadRune", (fmt.ScanState).UnreadRune, execiScanStateUnreadRune),
		I.Func("(ScanState).Width", (fmt.ScanState).Width, execiScanStateWidth),
		I.Func("(Scanner).Scan", (fmt.Scanner).Scan, execiScannerScan),
		I.Func("(State).Flag", (fmt.State).Flag, execiStateFlag),
		I.Func("(State).Precision", (fmt.State).Precision, execiStatePrecision),
		I.Func("(State).Width", (fmt.State).Width, execiStateWidth),
		I.Func("(State).Write", (fmt.State).Write, execiStateWrite),
		I.Func("(Stringer).String", (fmt.Stringer).String, execiStringerString),
	)
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
