package qlang

import (
	"errors"
	"reflect"

	"qiniupkg.com/text/tpl.v1/interpreter"

	"qlang.io/exec.v2"
	"qlang.io/qlang.spec.v1"
	qlangv2 "qlang.io/qlang.v2"
)

type Options interpreter.Options

var (
	InsertSemis = (*Options)(interpreter.InsertSemis)
)

func SetReadFile(fn func(file string) ([]byte, error)) {

	qlangv2.ReadFile = fn
}

func SetFindEntry(fn func(file string, libs []string) (string, error)) {

	qlangv2.FindEntry = fn
}

func SetOnPop(fn func(v interface{})) {

	exec.OnPop = fn
}

func SetDumpCode(dumpCode string) {

	switch dumpCode {
	case "true", "1":
		qlangv2.DumpCode = 1
	case "2":
		qlangv2.DumpCode = 2
	default:
		qlangv2.DumpCode = 0
	}
}

func Debug(fn func()) {

	SetDumpCode("1")
	defer SetDumpCode("0")
	fn()
}

// -----------------------------------------------------------------------------

type Qlang struct {
	*exec.Context
	cl *qlangv2.Compiler
}

func New(options *Options) (lang *Qlang, err error) {

	cl := qlangv2.New()
	cl.Opts = (*interpreter.Options)(options)
	stk := exec.NewStack()
	ctx := exec.NewContext()
	ctx.Stack = stk
	ctx.Code = cl.Code()
	return &Qlang{ctx, cl}, nil
}

func (p *Qlang) SetLibs(libs string) {

	p.cl.SetLibs(libs)
}

func (p *Qlang) Cl(codeText []byte, fname string) (end int, err error) {

	end = p.cl.Cl(codeText, fname)
	p.cl.Done()
	return
}

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

func (p *Qlang) Exec(codeText []byte, fname string) (err error) {

	code := p.cl.Code()
	start := code.Len()
	end, err := p.Cl(codeText, fname)
	if err != nil {
		return
	}

	if qlangv2.DumpCode != 0 {
		code.Dump(start)
	}

	p.ExecBlock(start, end)
	return
}

func (p *Qlang) Eval(expr string) (err error) {

	return p.Exec([]byte(expr), "")
}

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

func (p *Qlang) SafeEval(expr string) (err error) {

	return p.SafeExec([]byte(expr), "")
}

func Import(mod string, table map[string]interface{}) {

	qlang.Import(mod, table)
}

func SetAutoCall(t reflect.Type) {

	qlang.SetAutoCall(t)
}

// -----------------------------------------------------------------------------

