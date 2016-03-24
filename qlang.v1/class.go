package qlang

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"qlang.io/qlang.spec.v1"
	"qiniupkg.com/text/tpl.v1/interpreter.util"
)

var (
	ErrNewWithoutClassName   = errors.New("new object without class name")
	ErrNewObjectWithNotClass = errors.New("can't new object: not a class")
)

// -----------------------------------------------------------------------------

type FunctionInfo struct {
	args     []string // args[0] => function name
	fnb      interface{}
	variadic bool
}

type Class struct {
	Fns    map[string]*Function
	engine interpreter.Engine
	parent *Interpreter
}

type Object struct {
	vars map[string]interface{}
	Cls  *Class
}

func (p *Object) SetVar(name string, val interface{}) {

	if _, ok := p.Cls.Fns[name]; ok {
		panic("set failed: class already have a method named " + name)
	}
	p.vars[name] = val
}

func SetMemberVar(m interface{}, args ...interface{}) {

	if v, ok := m.(*Object); ok {
		for i := 0; i < len(args); i += 2 {
			v.SetVar(args[i].(string), args[i+1])
		}
		return
	}
	panic(fmt.Sprintf("type `%v` doesn't support `set` operator", reflect.TypeOf(m)))
}

func init() {
	qlang.Set = SetMemberVar
}

// -----------------------------------------------------------------------------

type thisDeref struct {
	this *Object
	fn   *Function
}

func (p *thisDeref) Call(a ...interface{}) interface{} {

	args := make([]interface{}, len(a)+1)
	args[0] = p.this
	for i, v := range a {
		args[i+1] = v
	}
	return p.fn.Call(args...)
}

// -----------------------------------------------------------------------------

func (p *Interpreter) newClass(e interpreter.Engine, m []interface{}) *Class {

	fns := make(map[string]*Function)
	cls := &Class{Fns: fns, engine: e, parent: p}
	for _, val := range m {
		v := val.(*FunctionInfo)
		name := v.args[0]
		v.args[0] = "this"
		fn := &Function{
			Args:     v.args,
			fnb:      v.fnb,
			engine:   e,
			parent:   p,
			Cls:      cls,
			Variadic: v.variadic,
		}
		fns[name] = fn
	}
	return cls
}

func (p *Interpreter) MemberFuncDecl() error {

	stk := p.stk
	fnb, ok := stk.Pop()
	if !ok {
		return ErrFunctionWithoutBody
	}

	variadic := stk.popArity()
	arity := stk.popArity()
	args := stk.popFnArgs(arity + 1)
	fn := &FunctionInfo{
		args:     args,
		fnb:      fnb,
		variadic: variadic != 0,
	}
	stk.Push(fn)
	return nil
}

func (p *Interpreter) Class(e interpreter.Engine) {

	stk := p.stk
	arity := stk.popArity()
	args := stk.popNArgs(arity)
	cls := p.newClass(e, args)
	stk.Push(cls)
}

func (p *Interpreter) New() error {

	var args []interface{}

	stk := p.stk
	hasArgs := stk.popArity()
	if hasArgs != 0 {
		nArgs := stk.popArity()
		args = stk.popNArgs(nArgs)
	}

	if v, ok := stk.Pop(); ok {
		if cls, ok := v.(*Class); ok {
			obj := &Object{
				vars: make(map[string]interface{}),
				Cls:  cls,
			}
			if init, ok := cls.Fns["_init"]; ok { // 构造函数
				closure := &thisDeref{
					this: obj,
					fn:   init,
				}
				closure.Call(args...)
			}
			stk.Push(obj)
			return nil
		}
		return ErrNewObjectWithNotClass
	}
	return ErrNewWithoutClassName
}

// -----------------------------------------------------------------------------

var (
	typeObjectPtr = reflect.TypeOf((*Object)(nil))
	typeClassPtr  = reflect.TypeOf((*Class)(nil))
)

func (p *Interpreter) MemberRef(name string) error {

	v, ok := p.stk.Pop()
	if !ok {
		return ErrRefWithoutObject
	}

	t := reflect.TypeOf(v)
	switch t {
	case typeObjectPtr:
		o := v.(*Object)
		val, ok := o.vars[name]
		if !ok {
			if fn, ok := o.Cls.Fns[name]; ok {
				t := &thisDeref{
					this: o,
					fn:   fn,
				}
				val = t
			} else {
				return fmt.Errorf("object doesn't has member `%s`", name)
			}
		}
		p.stk.Push(val)
		return nil
	case typeClassPtr:
		o := v.(*Class)
		val, ok := o.Fns[name]
		if !ok {
			return fmt.Errorf("class doesn't has method `%s`", name)
		}
		p.stk.Push(val)
		return nil
	}

	obj := reflect.ValueOf(v)
	switch {
	case obj.Kind() == reflect.Map:
		m := obj.MapIndex(reflect.ValueOf(name))
		if m.IsValid() {
			p.stk.Push(m.Interface())
		} else {
			return fmt.Errorf("member `%s` not found", name)
		}
	default:
		name = strings.Title(name)
		m := obj.MethodByName(name)
		if m.IsValid() {
			if qlang.AutoCall[t] && m.Type().NumIn() == 0 {
				out := m.Call(nil)
				p.stk.PushRet(out)
				return nil
			}
		} else {
			m = reflect.Indirect(obj).FieldByName(name)
			if !m.IsValid() {
				return fmt.Errorf("type `%v` doesn't has member `%s`", obj.Type(), name)
			}
		}
		p.stk.Push(m.Interface())
	}
	return nil
}

// -----------------------------------------------------------------------------
