package exec

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	qlang "qlang.io/spec"
)

var (
	// ErrNewWithoutClassName is returned when new doesn't specify a class.
	ErrNewWithoutClassName = errors.New("new object without class name")

	// ErrNewObjectWithNotType is returned when new T but T is not a type.
	ErrNewObjectWithNotType = errors.New("can't new object: not a type")

	// ErrRefWithoutObject is returned when refer member without specified an object.
	ErrRefWithoutObject = errors.New("reference without object")
)

// -----------------------------------------------------------------------------

// A Class represents a qlang class.
//
type Class struct {
	Fns map[string]*Function
	Ctx *Context
}

// Exec is required by interface Instr.
//
func (p *Class) Exec(stk *Stack, ctx *Context) {

	p.Ctx = ctx
	for _, f := range p.Fns {
		f.Parent = ctx
	}
	stk.Push(p)
}

// GoType returns the underlying go type. required by `qlang type` spec.
//
func (p *Class) GoType() reflect.Type {

	return typeIntf
}

// NewInstance creates a new instance of a qlang type. required by `qlang type` spec.
//
func (p *Class) NewInstance(stk *Stack, args ...interface{}) interface{} {

	return p.New(stk, args...)
}

// New creates a new instance of this class.
//
func (p *Class) New(stk *Stack, args ...interface{}) *Object {

	obj := &Object{
		vars: make(map[string]interface{}),
		Cls:  p,
	}
	if init, ok := p.Fns["_init"]; ok { // 构造函数
		closure := &Method{
			this: obj,
			fn:   init,
		}
		closure.Call(stk, args...)
	} else if len(args) > 0 {
		panic("constructor `_init` not found")
	}
	return obj
}

// IClass returns a Class instruction.
//
func IClass() *Class {

	fns := make(map[string]*Function)
	return &Class{Fns: fns}
}

// -----------------------------------------------------------------------------

// A Object represents a qlang object.
//
type Object struct {
	vars map[string]interface{}
	Cls  *Class
}

// Vars returns all variables (not including methods) of this object.
//
func (p *Object) Vars() map[string]interface{} {
	return p.vars
}

// SetVar sets the value of a qlang object's member.
//
func (p *Object) SetVar(name string, val interface{}) {

	if _, ok := p.Cls.Fns[name]; ok {
		panic("set failed: class already have a method named " + name)
	}
	p.vars[name] = val
}

// Member returns a member of this object.
//
func (p *Object) Member(name string) interface{} {

	if val, ok := p.vars[name]; ok {
		return val
	}
	if fn, ok := p.Cls.Fns[name]; ok {
		return &Method{
			this: p,
			fn:   fn,
		}
	}
	panic(fmt.Errorf("object doesn't has member `%s`", name))
}

// GetMemberVar implements get(object, key).
//
func GetMemberVar(m interface{}, key interface{}) interface{} {

	if v, ok := m.(*Object); ok {
		if name, ok := key.(string); ok {
			return v.Member(name)
		}
		panic("get(object, member): member should be `string` type")
	}
	panic(fmt.Sprintf("type `%v` doesn't support `get` operator", reflect.TypeOf(m)))
}

func init() {
	qlang.GetEx = GetMemberVar
}

// -----------------------------------------------------------------------------

// A Method represents a method of an Object.
//
type Method struct {
	this *Object
	fn   *Function
}

// Call calls this method with arguments.
//
func (p *Method) Call(stk *Stack, a ...interface{}) interface{} {

	args := make([]interface{}, len(a)+1)
	args[0] = p.this
	for i, v := range a {
		args[i+1] = v
	}
	return p.fn.Call(stk, args...)
}

// -----------------------------------------------------------------------------

type newTyper interface {
	NewInstance(args ...interface{}) interface{}
}

type newTyperEx interface {
	NewInstance(stk *Stack, args ...interface{}) interface{}
}

type iNew int

func (nArgs iNew) Exec(stk *Stack, ctx *Context) {

	var args []interface{}

	if nArgs != 0 {
		args = stk.PopNArgs(int(nArgs))
	}

	if v, ok := stk.Pop(); ok {
		if cls, ok := v.(newTyperEx); ok {
			obj := cls.NewInstance(stk, args...)
			stk.Push(obj)
			return
		}
		if cls, ok := v.(newTyper); ok {
			obj := cls.NewInstance(args...)
			stk.Push(obj)
			return
		}
		panic(ErrNewObjectWithNotType)
	}
	panic(ErrNewWithoutClassName)
}

// INew returns a New instruction.
//
func INew(nArgs int) Instr {
	return iNew(nArgs)
}

// -----------------------------------------------------------------------------

type iMemberRef struct {
	name string
}

var (
	typeObjectPtr = reflect.TypeOf((*Object)(nil))
	typeClassPtr  = reflect.TypeOf((*Class)(nil))
)

func (p *iMemberRef) Exec(stk *Stack, ctx *Context) {

	v, ok := stk.Pop()
	if !ok {
		panic(ErrRefWithoutObject)
	}

	name := p.name
	t := reflect.TypeOf(v)
	switch t {
	case typeObjectPtr:
		val := v.(*Object).Member(name)
		stk.Push(val)
		return
	case typeClassPtr:
		o := v.(*Class)
		val, ok := o.Fns[name]
		if !ok {
			panic(fmt.Errorf("class doesn't has method `%s`", name))
		}
		stk.Push(val)
		return
	}

	obj := reflect.ValueOf(v)
	kind := obj.Kind()
	switch {
	case kind == reflect.Map:
		m := obj.MapIndex(reflect.ValueOf(name))
		if m.IsValid() {
			stk.Push(m.Interface())
		} else {
			panic(fmt.Errorf("member `%s` not found", name))
		}
	default:
		name = strings.Title(name)
		m := obj.MethodByName(name)
		if !m.IsValid() {
			if kind == reflect.Ptr {
				obj = reflect.Indirect(obj)
				m = obj.FieldByName(name)
			} else if obj.CanAddr() {
				obj = obj.Addr()
				m = obj.MethodByName(name)
			}
			if !m.IsValid() {
				panic(fmt.Errorf("type `%v` doesn't has member `%s`", reflect.TypeOf(v), name))
			}
			t = obj.Type()
		}
		if qlang.AutoCall[t] && m.Type().NumIn() == 0 {
			out := m.Call(nil)
			stk.PushRet(out)
			return
		}
		stk.Push(m.Interface())
	}
}

func (p *iMemberRef) ToVar() Instr {
	return &iMemberVar{p.name}
}

// MemberRef returns a MemberRef instruction.
//
func MemberRef(name string) Instr {
	return &iMemberRef{name}
}

// -----------------------------------------------------------------------------
// MemberVar

type iMemberVar struct {
	name string
}

func (p *iMemberVar) Exec(stk *Stack, ctx *Context) {

	v, ok := stk.Pop()
	if !ok {
		panic(ErrRefWithoutObject)
	}

	stk.Push(&qlang.DataIndex{Data: v, Index: p.name})
}

// MemberVar returns a MemberVar instruction.
//
func MemberVar(name string) Instr {
	return &iMemberVar{name}
}

// -----------------------------------------------------------------------------
