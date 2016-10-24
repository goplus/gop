package exec

import (
	"fmt"
	"reflect"

	qlang "qlang.io/spec"
)

// -----------------------------------------------------------------------------
// AssignEx

type iAssignEx int

func (p iAssignEx) Exec(stk *Stack, ctx *Context) {

	v, ok1 := stk.Pop()
	k, ok2 := stk.Pop()
	if !ok1 || !ok2 {
		panic(ErrAssignWithoutVal)
	}

	doAssign(k, v, ctx)
}

func doAssign(k, v interface{}, ctx *Context) {

	switch t := k.(type) {
	case *variable:
		varRef := ctx.getVarRef(t.Name)
		*varRef = v
	case *qlang.DataIndex:
		qlang.SetIndex(t.Data, t.Index, v)
	default:
		panic("invalid assignment statement")
	}
}

const (
	symbolNameBits = 20
	symbolIndexMax = 1 << symbolNameBits
)

func (p *Context) getVarRef(name int) (varRef *interface{}) {

	if name < symbolIndexMax {
		return p.FastRefVar(name)
	}
	scope := name >> symbolNameBits
	for scope > 0 {
		p = p.parent
		scope--
	}
	return p.FastRefVar(name & (symbolIndexMax - 1))
}

// AssignEx is qlang AssignEx instruction.
//
var AssignEx Instr = iAssignEx(0)

// -----------------------------------------------------------------------------

type iMultiAssignFromSliceEx int

func (p iMultiAssignFromSliceEx) Exec(stk *Stack, ctx *Context) {

	arity := int(p)
	args := stk.PopNArgs(arity + 1)

	v := reflect.ValueOf(args[arity])
	if v.Kind() != reflect.Slice {
		panic(ErrMultiAssignExprMustBeSlice)
	}

	n := v.Len()
	if arity != n {
		panic(fmt.Errorf("multi assignment error: require %d variables, but we got %d", n, arity))
	}

	for i := 0; i < arity; i++ {
		doAssign(args[i], v.Index(i).Interface(), ctx)
	}
}

// MultiAssignFromSliceEx returns a MultiAssignFromSliceEx instruction.
//
func MultiAssignFromSliceEx(arity int) Instr {
	return iMultiAssignFromSliceEx(arity)
}

// -----------------------------------------------------------------------------

type iMultiAssignEx int

func (p iMultiAssignEx) Exec(stk *Stack, ctx *Context) {

	arity := int(p)
	args := stk.PopNArgs(arity << 1)

	for i := 0; i < arity; i++ {
		doAssign(args[i], args[arity+i], ctx)
	}
}

// MultiAssignEx returns a MultiAssignEx instruction.
//
func MultiAssignEx(arity int) Instr {
	return iMultiAssignEx(arity)
}

// -----------------------------------------------------------------------------
// AddAssign/SubAssign/MulAssign/QuoAssign/ModAssign/Inc/Dec

func execOpAssignEx(p func(a, b interface{}) interface{}, stk *Stack, ctx *Context) {

	v, ok1 := stk.Pop()
	k, ok2 := stk.Pop()
	if !ok1 || !ok2 {
		panic(ErrAssignWithoutVal)
	}

	switch t := k.(type) {
	case *variable:
		name := t.Name
		varRef := ctx.getVarRef(name)
		*varRef = p(*varRef, v)
	case *qlang.DataIndex:
		val := qlang.Get(t.Data, t.Index)
		val = p(val, v)
		qlang.SetIndex(t.Data, t.Index, val)
	default:
		panic("invalid op assignment statement")
	}
}

func execOp1AssignEx(p func(a interface{}) interface{}, stk *Stack, ctx *Context) {

	k, ok := stk.Pop()
	if !ok {
		panic(ErrAssignWithoutVal)
	}

	switch t := k.(type) {
	case *variable:
		name := t.Name
		varRef := ctx.getVarRef(name)
		*varRef = p(*varRef)
	case *qlang.DataIndex:
		val := qlang.Get(t.Data, t.Index)
		val = p(val)
		qlang.SetIndex(t.Data, t.Index, val)
	default:
		panic("invalid op1 assignment statement")
	}
}

// -----------------------------------------------------------------------------

type iAddAssignEx int
type iSubAssignEx int
type iMulAssignEx int
type iQuoAssignEx int
type iModAssignEx int
type iXorAssignEx int
type iBitAndAssignEx int
type iBitOrAssignEx int
type iAndNotAssignEx int
type iLshrAssignEx int
type iRshrAssignEx int
type iIncEx int
type iDecEx int

func (p iAddAssignEx) Exec(stk *Stack, ctx *Context) {
	execOpAssignEx(qlang.Add, stk, ctx)
}

func (p iSubAssignEx) Exec(stk *Stack, ctx *Context) {
	execOpAssignEx(qlang.Sub, stk, ctx)
}

func (p iMulAssignEx) Exec(stk *Stack, ctx *Context) {
	execOpAssignEx(qlang.Mul, stk, ctx)
}

func (p iQuoAssignEx) Exec(stk *Stack, ctx *Context) {
	execOpAssignEx(qlang.Quo, stk, ctx)
}

func (p iModAssignEx) Exec(stk *Stack, ctx *Context) {
	execOpAssignEx(qlang.Mod, stk, ctx)
}

func (p iXorAssignEx) Exec(stk *Stack, ctx *Context) {
	execOpAssignEx(qlang.Xor, stk, ctx)
}

func (p iBitAndAssignEx) Exec(stk *Stack, ctx *Context) {
	execOpAssignEx(qlang.BitAnd, stk, ctx)
}

func (p iBitOrAssignEx) Exec(stk *Stack, ctx *Context) {
	execOpAssignEx(qlang.BitOr, stk, ctx)
}

func (p iAndNotAssignEx) Exec(stk *Stack, ctx *Context) {
	execOpAssignEx(qlang.AndNot, stk, ctx)
}

func (p iLshrAssignEx) Exec(stk *Stack, ctx *Context) {
	execOpAssignEx(qlang.Lshr, stk, ctx)
}

func (p iRshrAssignEx) Exec(stk *Stack, ctx *Context) {
	execOpAssignEx(qlang.Rshr, stk, ctx)
}

func (p iIncEx) Exec(stk *Stack, ctx *Context) {
	execOp1AssignEx(qlang.Inc, stk, ctx)
}

func (p iDecEx) Exec(stk *Stack, ctx *Context) {
	execOp1AssignEx(qlang.Dec, stk, ctx)
}

var (
	AddAssignEx    Instr = iAddAssignEx(0)
	SubAssignEx    Instr = iSubAssignEx(0)
	MulAssignEx    Instr = iMulAssignEx(0)
	QuoAssignEx    Instr = iQuoAssignEx(0)
	ModAssignEx    Instr = iModAssignEx(0)
	XorAssignEx    Instr = iXorAssignEx(0)
	BitAndAssignEx Instr = iBitAndAssignEx(0)
	BitOrAssignEx  Instr = iBitOrAssignEx(0)
	AndNotAssignEx Instr = iAndNotAssignEx(0)
	LshrAssignEx   Instr = iLshrAssignEx(0)
	RshrAssignEx   Instr = iRshrAssignEx(0)
	IncEx          Instr = iIncEx(0)
	DecEx          Instr = iDecEx(0)
)

// -----------------------------------------------------------------------------
