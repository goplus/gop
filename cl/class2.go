package qlang

import (
	"errors"

	ipt "qiniupkg.com/text/tpl.v1/interpreter"
	"qiniupkg.com/text/tpl.v1/interpreter.util"
	"qlang.io/exec"
)

// -----------------------------------------------------------------------------

type functionInfo struct {
	args     []string // args[0] => function name
	fnb      interface{}
	variadic bool
}

// -----------------------------------------------------------------------------

func (p *Compiler) memberFuncDecl() {

	fnb, _ := p.gstk.Pop()
	variadic := p.popArity()
	arity := p.popArity()
	args := p.gstk.PopFnArgs(arity + 1)
	fn := &functionInfo{
		args:     args,
		fnb:      fnb,
		variadic: variadic != 0,
	}
	p.gstk.Push(fn)
}

func (p *Compiler) addMethods(cls *exec.Class, e interpreter.Engine, parent *funcCtx, members []interface{}) {

	for _, val := range members {
		v := val.(*functionInfo)
		name := v.args[0]
		v.args[0] = "this"
		start, end, symtbl := p.clFunc(e, "doc", v.fnb, parent, v.args)
		fn := exec.NewFunction(cls, start, end, symtbl, v.args, v.variadic)
		fn.Parent = cls.Ctx
		cls.Fns[name] = fn
	}
}

func (p *Compiler) fnClass(e interpreter.Engine) {

	arity := p.popArity()
	members := p.gstk.PopNArgs(arity)
	instr := p.code.Reserve()
	fnctx := p.fnctx
	p.exits = append(p.exits, func() {
		cls := exec.IClass()
		p.addMethods(cls, e, fnctx, members)
		instr.Set(cls)
	})
}

// InjectMethods injects some methods into a class.
//
func (p *Compiler) InjectMethods(cls *exec.Class, code []byte) (err error) {

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

	engine, err := ipt.New(p, p.Opts)
	if err != nil {
		return
	}
	p.ipt = engine

	src, err := engine.Tokenize(code, "")
	if err != nil {
		return
	}

	p.cl(engine, "methods", src)
	arity := p.popArity()
	members := p.gstk.PopNArgs(arity)
	p.addMethods(cls, engine, p.fnctx, members)
	p.Done()
	return
}

func (p *Compiler) fnNew() {

	nArgs := p.popArity()
	if nArgs != 0 {
		nArgs = p.popArity()
	}
	p.code.Block(exec.INew(nArgs))
}

func (p *Compiler) memberRef(name string) {

	p.code.Block(exec.MemberRef(name))
}

// -----------------------------------------------------------------------------
