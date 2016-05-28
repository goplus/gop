package exec

import (
	"qlang.io/qlang.spec.v1"
)

// -----------------------------------------------------------------------------

type tchan int

func (p tchan) OptimizableGetArity() int {

	return 1
}

func (p tchan) Exec(stk *Stack, ctx *Context) {

	n := len(stk.data) - 1
	stk.data[n] = qlang.ChanOf(stk.data[n])
}

// Chan is an instruction that returns chan T.
//
var Chan Instr = tchan(0)

// -----------------------------------------------------------------------------

type slice int

func (p slice) OptimizableGetArity() int {

	return 1
}

func (p slice) Exec(stk *Stack, ctx *Context) {

	n := len(stk.data) - 1
	stk.data[n] = qlang.Slice(stk.data[n])
}

// Slice is an instruction that returns []T.
//
var Slice Instr = slice(0)

// -----------------------------------------------------------------------------

type sliceFrom int
type sliceFromTy int

var nilVarSlice = make([]interface{}, 0)

func (p sliceFrom) OptimizableGetArity() int {

	return int(p)
}

func (p sliceFromTy) OptimizableGetArity() int {

	return int(p)
}

func (p sliceFrom) Exec(stk *Stack, ctx *Context) {

	if p == 0 {
		stk.data = append(stk.data, nilVarSlice)
		return
	}
	n := len(stk.data) - int(p)
	stk.data[n] = qlang.SliceFrom(stk.data[n:]...)
	stk.data = stk.data[:n+1]
}

func (p sliceFromTy) Exec(stk *Stack, ctx *Context) {

	n := len(stk.data) - int(p)
	stk.data[n] = qlang.SliceFromTy(stk.data[n:]...)
	stk.data = stk.data[:n+1]
}

// SliceFrom is an instruction that creates slice in [a1, a2, ...] form.
//
func SliceFrom(arity int) Instr {

	return sliceFrom(arity)
}

// SliceFromTy is an instruction that creates slice in []T{a1, a2, ...} form.
//
func SliceFromTy(arity int) Instr {

	return sliceFromTy(arity)
}

// -----------------------------------------------------------------------------
// StructInit

type iStructInit int

func (p iStructInit) OptimizableGetArity() int {

	return int(p)
}

func (p iStructInit) Exec(stk *Stack, ctx *Context) {

	n := len(stk.data) - int(p)
	stk.data[n] = qlang.StructInit(stk.data[n:]...)
	stk.data = stk.data[:n+1]
}

// StructInit returns a StructInit instruction that means `&StructType{name1: expr1, name2: expr2, ...}`.
//
func StructInit(arity int) Instr {
	return iStructInit(arity)
}

// -----------------------------------------------------------------------------
// MapInit

type iMapInit int

func (p iMapInit) OptimizableGetArity() int {

	return int(p)
}

func (p iMapInit) Exec(stk *Stack, ctx *Context) {

	n := len(stk.data) - int(p)
	stk.data[n] = qlang.MapInit(stk.data[n:]...)
	stk.data = stk.data[:n+1]
}

// MapInit returns a MapInit instruction that means `map[key]elem{key1: expr1, key2: expr2, ...}`.
//
func MapInit(arity int) Instr {
	return iMapInit(arity)
}

// -----------------------------------------------------------------------------

type tmap int

func (p tmap) OptimizableGetArity() int {

	return 2
}

func (p tmap) Exec(stk *Stack, ctx *Context) {

	n := len(stk.data) - 1
	stk.data[n-1] = qlang.Map(stk.data[n-1], stk.data[n])
	stk.data = stk.data[:n]
}

// Map is an instruction that returns `map[key]elem`.
//
var Map Instr = tmap(0)

// -----------------------------------------------------------------------------
