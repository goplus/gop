package exec

import (
	"errors"
	"fmt"
	"reflect"

	"qlang.io/qlang.spec.v1"
)

var (
	ErrAssignWithoutVal           = errors.New("variable assign without value")
	ErrMultiAssignExprMustBeSlice = errors.New("expression of multi assignment must return slice")
)

// -----------------------------------------------------------------------------
// Unset

type iUnset struct {
	name string
}

func (p *iUnset) Exec(stk *Stack, ctx *Context) {
	delete(ctx.vars, p.name)
}

func Unset(name string) Instr {
	return &iUnset{name}
}

// -----------------------------------------------------------------------------
// Assign

type iAssign struct {
	name string
}

func (p *iAssign) Exec(stk *Stack, ctx *Context) {

	vars := ctx.getVars(p.name)
	if v, ok := stk.Pop(); ok {
		vars[p.name] = v
	} else {
		panic(ErrAssignWithoutVal)
	}
}

type externVar struct {
	vars map[string]interface{}
}

func (p *Context) getVars(name string) (vars map[string]interface{}) {

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

func Assign(name string) Instr {
	return &iAssign{name}
}

// -----------------------------------------------------------------------------

type iMultiAssignFromSlice struct {
	names []string
}

func (p *iMultiAssignFromSlice) Exec(stk *Stack, ctx *Context) {

	val, ok := stk.Pop()
	if !ok {
		panic(ErrAssignWithoutVal)
	}

	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Slice {
		panic(ErrMultiAssignExprMustBeSlice)
	}

	n := v.Len()
	arity := len(p.names)
	if arity != n {
		panic(fmt.Errorf("multi assignment error: require %d variables, but we got %d", n, arity))
	}

	for i, name := range p.names {
		vars := ctx.getVars(name)
		vars[name] = v.Index(i).Interface()
	}
}

func MultiAssignFromSlice(names []string) Instr {
	return &iMultiAssignFromSlice{names}
}

// -----------------------------------------------------------------------------

type iMultiAssign struct {
	names []string
}

func (p *iMultiAssign) Exec(stk *Stack, ctx *Context) {

	n := len(p.names)
	for i := n; i > 0; {
		val, ok := stk.Pop()
		if !ok {
			panic(ErrAssignWithoutVal)
		}
		i--
		name := p.names[i]
		vars := ctx.getVars(name)
		vars[name] = val
	}
}

func MultiAssign(names []string) Instr {
	return &iMultiAssign{names}
}

// -----------------------------------------------------------------------------
// AddAssign/SubAssign/MulAssign/QuoAssign/ModAssign/Inc/Dec

type iOpAssign struct {
	name string
	op   func(a, b interface{}) interface{}
}

type iOp1Assign struct {
	name string
	op   func(a interface{}) interface{}
}

func (p *iOpAssign) Exec(stk *Stack, ctx *Context) {

	vars, val := ctx.getVar(p.name)
	v, ok := stk.Pop()
	if !ok {
		panic(ErrAssignWithoutVal)
	}
	val = p.op(val, v)
	vars[p.name] = val
}

func (p *iOp1Assign) Exec(stk *Stack, ctx *Context) {

	vars, val := ctx.getVar(p.name)
	val = p.op(val)
	vars[p.name] = val
}

func (p *Context) getVar(name string) (vars map[string]interface{}, val interface{}) {

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

func OpAssign(name string, op func(a, b interface{}) interface{}) Instr {
	return &iOpAssign{name, op}
}

func AddAssign(name string) Instr {
	return &iOpAssign{name, qlang.Add}
}

func SubAssign(name string) Instr {
	return &iOpAssign{name, qlang.Sub}
}

func MulAssign(name string) Instr {
	return &iOpAssign{name, qlang.Mul}
}

func QuoAssign(name string) Instr {
	return &iOpAssign{name, qlang.Quo}
}

func ModAssign(name string) Instr {
	return &iOpAssign{name, qlang.Mod}
}

func Inc(name string) Instr {
	return &iOp1Assign{name, qlang.Inc}
}

func Dec(name string) Instr {
	return &iOp1Assign{name, qlang.Dec}
}

// -----------------------------------------------------------------------------

type iRef struct {
	name string
}

func (p *iRef) Exec(stk *Stack, ctx *Context) {

	stk.Push(ctx.getRef(p.name))
}

func (p *Context) getRef(name string) interface{} {

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

func Ref(name string) Instr {
	return &iRef{name}
}

// -----------------------------------------------------------------------------

