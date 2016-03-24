package qlang

import (
	"errors"
	"reflect"

	"qiniupkg.com/text/tpl.v1/interpreter"
	"qlang.io/qlang.spec.v1"
	qlangv1 "qlang.io/qlang.v1"
)

type Options interpreter.Options

var (
	InsertSemis = (*Options)(interpreter.InsertSemis)
)

// -----------------------------------------------------------------------------

type Qlang struct {
	*qlangv1.Interpreter
	*interpreter.Engine
}

func New(options *Options) (lang Qlang, err error) {

	calc := qlangv1.New()
	engine, err := interpreter.New(calc, (*interpreter.Options)(options))
	if err != nil {
		return
	}
	return Qlang{calc, engine}, nil
}

func (p Qlang) Exec(code []byte, fname string) (err error) {

	return p.MatchExactly(code, fname)
}

func (p Qlang) SafeExec(code []byte, fname string) (err error) {

	defer func() {
		p.Recov = recover()
		p.ExecDefers(p.Engine)
		if e := p.Recov; e != nil {
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

	err = p.MatchExactly(code, fname)
	return
}

func (p Qlang) SafeEval(expr string) (err error) {

	return p.SafeExec([]byte(expr), "")
}

func Import(mod string, table map[string]interface{}) {

	qlang.Import(mod, table)
}

func SetAutoCall(t reflect.Type) {

	qlang.SetAutoCall(t)
}

// -----------------------------------------------------------------------------

