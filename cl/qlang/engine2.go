package qlang

import (
	"errors"
	"fmt"
	"reflect"

	"qiniupkg.com/text/tpl.v1/interpreter"

	qcl "qlang.io/cl"
	"qlang.io/exec"
	qlang "qlang.io/spec"
)

// Options represent interpreter options.
//
type Options interpreter.Options

var (
	// InsertSemis is interpreter options that means to insert semis smartly.
	InsertSemis = (*Options)(interpreter.InsertSemis)
)

// SetReadFile sets the `ReadFile` function.
//
func SetReadFile(fn func(file string) ([]byte, error)) {

	qcl.ReadFile = fn
}

// SetFindEntry sets the `FindEntry` function.
//
func SetFindEntry(fn func(file string, libs []string) (string, error)) {

	qcl.FindEntry = fn
}

// SetOnPop sets OnPop callback.
//
func SetOnPop(fn func(v interface{})) {

	exec.OnPop = fn
}

// SetDumpCode sets dump code mode:
//	"1" - dump code with rem instruction.
//	"2" - dump code without rem instruction.
//	else - don't dump code.
//
func SetDumpCode(dumpCode string) {

	switch dumpCode {
	case "true", "1":
		qcl.DumpCode = 1
	case "2":
		qcl.DumpCode = 2
	default:
		qcl.DumpCode = 0
	}
}

// Debug sets dump code mode to "1" for debug.
//
func Debug(fn func()) {

	SetDumpCode("1")
	defer SetDumpCode("0")
	fn()
}

// -----------------------------------------------------------------------------

// A Qlang represents the qlang compiler and executor.
//
type Qlang struct {
	*exec.Context
	cl *qcl.Compiler
}

// New returns a new qlang instance.
//
func New() *Qlang {

	cl := qcl.New()
	stk := exec.NewStack()
	ctx := exec.NewContextEx(cl.GlobalSymbols())
	ctx.Stack = stk
	ctx.Code = cl.Code()
	return &Qlang{ctx, cl}
}

// SetLibs sets lib paths for searching modules.
//
func (p *Qlang) SetLibs(libs string) {

	p.cl.SetLibs(libs)
}

// Cl compiles a source code.
//
func (p *Qlang) Cl(codeText []byte, fname string) (end int, err error) {

	end = p.cl.Cl(codeText, fname)
	p.cl.Done()
	p.ResizeVars()
	return
}

// SafeCl compiles a source code, without panic (will convert panic into an error).
//
func (p *Qlang) SafeCl(codeText []byte, fname string) (end int, err error) {

	defer func() {
		if e := recover(); e != nil {
			switch v := e.(type) {
			case string:
				err = errors.New(v)
			case error:
				err = v
			default:
				panic(e)
			}
		}
	}()

	return p.Cl(codeText, fname)
}

// Exec compiles and executes a source code.
//
func (p *Qlang) Exec(codeText []byte, fname string) (err error) {

	code := p.cl.Code()
	start := code.Len()
	end, err := p.Cl(codeText, fname)
	if err != nil {
		return
	}

	if qcl.DumpCode != 0 {
		code.Dump(start)
	}

	p.ExecBlock(start, end, p.cl.GlobalSymbols())
	return
}

// Eval compiles and executes a source code.
//
func (p *Qlang) Eval(expr string) (err error) {

	return p.Exec([]byte(expr), "")
}

// SafeExec compiles and executes a source code, without panic (will convert panic into an error).
//
func (p *Qlang) SafeExec(code []byte, fname string) (err error) {

	defer func() {
		if e := recover(); e != nil {
			switch v := e.(type) {
			case string:
				err = errors.New(v)
			case error:
				err = v
			default:
				panic(e)
			}
		}
	}()

	err = p.Exec(code, fname)
	return
}

// SafeEval compiles and executes a source code, without panic (will convert panic into an error).
//
func (p *Qlang) SafeEval(expr string) (err error) {

	return p.SafeExec([]byte(expr), "")
}

// InjectMethods injects some methods into a class.
// `pcls` can be a `*exec.Class` object or a `string` typed class name.
//
func (p *Qlang) InjectMethods(pcls interface{}, code []byte) (err error) {

	var cls *exec.Class
	switch v := pcls.(type) {
	case *exec.Class:
		cls = v
	case string:
		val, ok := p.GetVar(v)
		if !ok {
			return fmt.Errorf("class `%s` not exists", v)
		}
		if cls, ok = val.(*exec.Class); !ok {
			return fmt.Errorf("var `%s` not a class", v)
		}
	default:
		return fmt.Errorf("invalid cls argument type: %v", reflect.TypeOf(pcls))
	}
	err = p.cl.InjectMethods(cls, code)
	p.ResizeVars()
	return
}

// Import imports a module written in Go.
//
func Import(mod string, table map[string]interface{}) {

	qlang.Import(mod, table)
}

// SetAutoCall is reserved for internal use.
//
func SetAutoCall(t reflect.Type) {

	qlang.SetAutoCall(t)
}

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"new": New,
	"New": New,
}

// -----------------------------------------------------------------------------
