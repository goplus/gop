package qlang

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"qlang.io/qlang.spec.v1"
	"qiniupkg.com/text/tpl.v1/interpreter.util"

	qip "qiniupkg.com/text/tpl.v1/interpreter"
)

var (
	ErrFunctionWithoutBody        = errors.New("function without body")
	ErrMultiAssignExprMustBeSlice = errors.New("expression of multi assignment must return slice")
)

var (
	ErrReturn = &qip.RuntimeError{Err: errors.New("return")}
)

// -----------------------------------------------------------------------------

func (stk *Stack) BaseFrame() int {

	return len(*stk)
}

func (stk *Stack) SetFrame(n int) {

	*stk = (*stk)[:n]
}

func PushName(stk *Stack, name string) {

	stk.Push(name)
}

// -----------------------------------------------------------------------------

func (p *Interpreter) Clear() {

	p.stk.SetFrame(p.base)
}

func (p *Interpreter) Unset(name string) {

	delete(p.vars, name)
}

func (p *Interpreter) Inc(name string) error {

	vars, val := p.getVar(name)
	val = qlang.Inc(val)
	vars[name] = val
	p.stk.Push(val)
	return nil
}

func (p *Interpreter) Dec(name string) error {

	vars, val := p.getVar(name)
	val = qlang.Dec(val)
	vars[name] = val
	p.stk.Push(val)
	return nil
}

type externVar struct {
	vars map[string]interface{}
}

func (p *Interpreter) getVar(name string) (vars map[string]interface{}, val interface{}) {

	vars = p.vars
	val, ok := vars[name]
	if !ok {
		panic(fmt.Sprintf("variable `%s` not found", name))
	}
	if e, ok := val.(externVar); ok {
		vars = e.vars
		val = vars[name]
	}
	return
}

func (p *Interpreter) AddAssign(name string) error {

	vars, val := p.getVar(name)
	v, ok := p.stk.Pop()
	if !ok {
		return ErrAssignWithoutVal
	}
	val = qlang.Add(val, v)
	vars[name] = val
	p.stk.Push(val)
	return nil
}

func (p *Interpreter) SubAssign(name string) error {

	vars, val := p.getVar(name)
	v, ok := p.stk.Pop()
	if !ok {
		return ErrAssignWithoutVal
	}
	val = qlang.Sub(val, v)
	vars[name] = val
	p.stk.Push(val)
	return nil
}

func (p *Interpreter) MulAssign(name string) error {

	vars, val := p.getVar(name)
	v, ok := p.stk.Pop()
	if !ok {
		return ErrAssignWithoutVal
	}
	val = qlang.Mul(val, v)
	vars[name] = val
	p.stk.Push(val)
	return nil
}

func (p *Interpreter) QuoAssign(name string) error {

	vars, val := p.getVar(name)
	v, ok := p.stk.Pop()
	if !ok {
		return ErrAssignWithoutVal
	}
	val = qlang.Quo(val, v)
	vars[name] = val
	p.stk.Push(val)
	return nil
}

func (p *Interpreter) ModAssign(name string) error {

	vars, val := p.getVar(name)
	v, ok := p.stk.Pop()
	if !ok {
		return ErrAssignWithoutVal
	}
	val = qlang.Mod(val, v)
	vars[name] = val
	p.stk.Push(val)
	return nil
}

func (p *Interpreter) MultiAssign() error {

	stk := p.stk
	val, ok := stk.Pop()
	if !ok {
		return ErrAssignWithoutVal
	}

	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Slice {
		return ErrMultiAssignExprMustBeSlice
	}

	arity := stk.popArity() + 1
	n := v.Len()
	if arity != n {
		return fmt.Errorf("multi assignment error: require %d variables, but we got %d", n, arity)
	}

	names := stk.popFnArgs(arity)
	for i := 0; i < n; i++ {
		name := names[i]
		vars := p.getVars(name)
		vars[name] = v.Index(i).Interface()
	}

	p.stk.Push(val)
	return nil
}

func (p *Interpreter) Assign(name string) error {

	vars := p.getVars(name)
	if v, ok := p.stk.Pop(); ok {
		vars[name] = v
		p.stk.Push(v)
		return nil
	}
	return ErrAssignWithoutVal
}

func (p *Interpreter) getVars(name string) (vars map[string]interface{}) {

	vars = p.vars
	if val, ok := vars[name]; ok { // 变量已经存在
		if e, ok := val.(externVar); ok {
			vars = e.vars
		}
		return
	}

	if name[0] != '_' { // NOTE: '_' 开头的变量是私有的，不可继承
		for t := p.parent; t != nil; t = t.parent {
			if _, ok := t.vars[name]; ok {
				panic(fmt.Sprintf("variable `%s` exists in extern function", name))
			}
		}
	}
	return
}

func (p *Interpreter) Ref(name string) {

	p.stk.Push(p.getRef(name))
}

func (p *Interpreter) getRef(name string) interface{} {

	if name == "unset" {
		return p.Unset
	}

	val, ok := p.vars[name]
	if ok {
		if e, ok1 := val.(externVar); ok1 {
			val = e.vars[name]
		}
	} else {
		if name[0] != '_' { // NOTE: '_' 开头的变量是私有的，不可继承
			for t := p.parent; t != nil; t = t.parent {
				if val, ok = t.vars[name]; ok {
					e, ok1 := val.(externVar)
					if ok1 {
						val = e.vars[name]
					} else {
						e = externVar{t.vars}
					}
					p.vars[name] = e // 缓存访问过的变量
					goto lzDone
				}
			}
		}
		if val, ok = qlang.Fntable[name]; !ok {
			panic("symbol not found: " + name)
		}
	}

lzDone:
	return val
}

func (p *Interpreter) Function(e interpreter.Engine) error {

	stk := p.stk
	fnb, ok := stk.Pop()
	if !ok {
		return ErrFunctionWithoutBody
	}

	variadic := stk.popArity()
	arity := stk.popArity()
	args := stk.popFnArgs(arity)
	fn := &Function{
		Args:     args,
		fnb:      fnb,
		engine:   e,
		parent:   p,
		Variadic: variadic != 0,
	}
	stk.Push(fn)
	return nil
}

func (p *Interpreter) Return(e interpreter.Engine) {

	arity := p.stk.popArity()
	if arity == 0 {
		p.ret = nil
	} else {
		if v, ok := p.stk.Pop(); ok {
			p.ret = v
		}
	}
	if p.parent != nil {
		panic(ErrReturn) // 利用 panic 来实现 return (正常退出)
	}

	p.ExecDefers(e)

	if p.ret == nil {
		os.Exit(0)
	}
	if v, ok := p.ret.(int); ok {
		os.Exit(v)
	}
	panic("must return `int` for main function")
}

func (p *Interpreter) ExecDefers(e interpreter.Engine) {

	for d := p.defers; d != nil; d = d.next {
		if err := e.EvalCode(p, "expr", d.src); err != nil {
			panic(err)
		}
	}
	p.defers = nil
}

func (p *Interpreter) Defer() {

	src, ok := p.stk.Pop()
	if !ok {
		panic("unreachable")
	}

	p.defers = &Defer{next: p.defers, src: src}
	p.stk.Push(nil)
}

func (p *Interpreter) Recover() {

	if parent := p.parent; parent != nil {
		e := parent.Recov
		parent.Recov = nil
		p.stk.Push(e)
	} else {
		p.stk.Push(nil)
	}
}

// -----------------------------------------------------------------------------

type Function struct {
	Args     []string
	fnb      interface{}
	engine   interpreter.Engine
	parent   *Interpreter
	Cls      *Class
	Variadic bool
}

func (p *Function) Call(args ...interface{}) interface{} {

	n := len(p.Args)
	if p.Variadic {
		if len(args) < n-1 {
			panic(fmt.Sprintf("function requires >= %d arguments, but we got %d", n-1, len(args)))
		}
	} else {
		if len(args) != n {
			panic(fmt.Sprintf("function requires %d arguments, but we got %d", n, len(args)))
		}
	}

	if p.fnb == nil {
		return nil
	}

	parent := p.parent
	stk := parent.stk
	vars := make(map[string]interface{})

	base := stk.BaseFrame()
	calc := &Interpreter{
		vars:   vars,
		parent: parent,
		stk:    stk,
		base:   base,
	}

	if p.Variadic {
		for i := 0; i < n-1; i++ {
			vars[p.Args[i]] = args[i]
		}
		vars[p.Args[n-1]] = args[n-1:]
	} else {
		for i, arg := range args {
			vars[p.Args[i]] = arg
		}
	}

	err := func() error {
		defer func() {
			if e := recover(); e != ErrReturn { // 正常 return 导致，见 (*Interpreter).Return 函数
				calc.Recov = e
			}
			calc.ExecDefers(p.engine)
			stk.SetFrame(base)
			if calc.Recov != nil {
				panic(calc.Recov)
			}
		}()
		return p.engine.EvalCode(calc, "doc", p.fnb)
	}()
	if err != nil {
		panic(err)
	}
	return calc.ret
}

// -----------------------------------------------------------------------------

func (p *Interpreter) Main(e interpreter.Engine) {

	fnb, _ := p.stk.Pop()
	fn := &Function{
		fnb:      fnb,
		engine:   e,
		parent:   p,
	}
	ret := fn.Call()

	code := 0
	if ret != nil {
		code = ret.(int)
	}
	os.Exit(code)
}

// -----------------------------------------------------------------------------

