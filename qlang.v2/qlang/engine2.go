package qlang

import (
	"errors"
	"path/filepath"
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

// -----------------------------------------------------------------------------

type Qlang struct {
	*exec.Context
	cl      *qlangv2.Compiler
	stk     *exec.Stack
	options *Options
}

func New(options *Options) (lang *Qlang, err error) {

	cl := qlangv2.New()
	stk := exec.NewStack()
	ctx := exec.NewContext()
	ctx.Stack = stk
	ctx.Code = cl.Code()
	return &Qlang{ctx, cl, stk, options}, nil
}

func (p *Qlang) Ret() (v interface{}, ok bool) {

	v, ok = p.stk.Pop()
	p.stk.SetFrame(0)
	return
}

func (p *Qlang) Cl(codeText []byte, fname string) (end int, err error) {

	cl := p.cl
	engine, err := interpreter.New(cl, (*interpreter.Options)(p.options))
	if err != nil {
		return
	}
	vars := cl.Vars()
	vars["__dir__"] = filepath.Dir(fname)
	vars["__file__"] = fname
	err = engine.MatchExactly(codeText, fname)
	if err != nil {
		return
	}
	end = cl.Code().Len()
	cl.Done()
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

	code.Exec(start, end, p.stk, p.Context)
	return
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

