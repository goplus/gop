package exec

import (
	"errors"

	qlang "qlang.io/spec"
)

var (
	// ErrAssignWithoutVal is returned when variable assign without value
	ErrAssignWithoutVal = errors.New("variable assign without value")

	// ErrMultiAssignExprMustBeSlice is returned when expression of multi assignment must be a slice
	ErrMultiAssignExprMustBeSlice = errors.New("expression of multi assignment must be a slice")
)

// -----------------------------------------------------------------------------

type iRef struct {
	name int
}

func (p *iRef) Exec(stk *Stack, ctx *Context) {

	stk.Push(ctx.getRef(p.name))
}

func (p *iRef) ToVar() Instr {

	return &iVar{p.name}
}

func (p *Context) getRef(name int) interface{} {

	if name < symbolIndexMax {
		return p.FastGetVar(name)
	}
	scope := name >> symbolNameBits
	for scope > 0 {
		p = p.parent
		scope--
	}
	return p.FastGetVar(name & (symbolIndexMax - 1))
}

// SymbolIndex returns symbol index value.
//
func SymbolIndex(id, scope int) int {
	return id | (scope << symbolNameBits)
}

// Ref returns an instruction that refers a variable.
//
func Ref(name int) Instr {
	return &iRef{name}
}

// -----------------------------------------------------------------------------
// Var

type iVar struct {
	name int
}

func (p *iVar) Exec(stk *Stack, ctx *Context) {

	stk.Push(&variable{Name: p.name})
}

// Var returns a Var instruction.
//
func Var(name int) Instr {
	return &iVar{name}
}

// -----------------------------------------------------------------------------

type iGet int

func (p iGet) Exec(stk *Stack, ctx *Context) {

	k, ok1 := stk.Pop()
	o, ok2 := stk.Pop()
	if !ok1 || !ok2 {
		panic("unexpected to call `Get` instruction")
	}
	stk.Push(qlang.Get(o, k))
}

func (p iGet) ToVar() Instr {

	return GetVar
}

// Get is the Get instruction.
//
var Get Instr = iGet(0)

// -----------------------------------------------------------------------------

type iGetVar int

func (p iGetVar) Exec(stk *Stack, ctx *Context) {

	k, ok1 := stk.Pop()
	o, ok2 := stk.Pop()
	if !ok1 || !ok2 {
		panic("unexpected to call `GetVar` instruction")
	}
	stk.Push(&qlang.DataIndex{Data: o, Index: k})
}

// GetVar is the Get instruction.
//
var GetVar Instr = iGetVar(0)

// -----------------------------------------------------------------------------

type iGfnRef struct {
	val   interface{}
	toVar func() Instr
}

func (p *iGfnRef) Exec(stk *Stack, ctx *Context) {

	stk.Push(p.val)
}

func (p *iGfnRef) ToVar() Instr {

	return p.toVar()
}

// GfnRef returns an instruction that refers a fntable item.
//
func GfnRef(v interface{}, toVar func() Instr) Instr {
	return &iGfnRef{v, toVar}
}

// -----------------------------------------------------------------------------
