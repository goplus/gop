package cl

import (
	"errors"
	"go/ast"
	"reflect"

	"github.com/qiniu/qlang/ast/astutil"
	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/x/log"
)

// iType represents a qlang type.
//  - *goType
//  - *retType
type iType interface {
	TypeValue(i int) iType
	NumValues() int
}

// iValue represents a qlang value.
//  - *goFunc
//  - *constVal
//  - *exprResult
type iValue interface {
	TypeOf() iType
}

// -----------------------------------------------------------------------------

type goType struct {
	t reflect.Type
}

func (p *goType) TypeValue(i int) iType {
	return p
}

func (p *goType) NumValues() int {
	return 1
}

// -----------------------------------------------------------------------------

type unboundType astutil.ConstKind

func (p unboundType) TypeValue(i int) iType {
	return p
}

func (p unboundType) NumValues() int {
	return 1
}

// -----------------------------------------------------------------------------

type retType struct {
	tfn reflect.Type
}

func (p *retType) TypeValue(i int) iType {
	return &goType{t: p.tfn.Out(i)}
}

func (p *retType) NumValues() int {
	return p.tfn.NumOut()
}

// -----------------------------------------------------------------------------

type goFunc struct {
	v    *exec.GoFuncInfo
	addr uint32
	kind exec.SymbolKind
}

func (p *goFunc) Proto() reflect.Type {
	return reflect.TypeOf(p.v.This)
}

func (p *goFunc) TypeOf() iType {
	return &goType{t: p.Proto()}
}

func (p *goFunc) TypesOut() *retType {
	return &retType{tfn: p.Proto()}
}

func getGoFunc(addr uint32, kind exec.SymbolKind) *goFunc {
	var fi *exec.GoFuncInfo
	switch kind {
	case exec.SymbolFunc:
		fi = exec.GoFuncAddr(addr).GetInfo()
	case exec.SymbolVariadicFunc:
		fi = exec.GoVariadicFuncAddr(addr).GetInfo()
	default:
		log.Fatalln("getGoFunc: unknown -", kind, addr)
	}
	return &goFunc{v: fi, addr: addr, kind: kind}
}

// -----------------------------------------------------------------------------

type constVal struct {
	v    interface{}
	kind astutil.ConstKind
}

func (p *constVal) TypeOf() iType {
	if astutil.IsConstBound(p.kind) {
		return &goType{t: reflect.TypeOf(p.v)}
	}
	return unboundType(p.kind)
}

// -----------------------------------------------------------------------------

type exprResult struct {
	results iType
	expr    ast.Expr
}

func (p *exprResult) TypeOf() iType {
	return p.results
}

// -----------------------------------------------------------------------------

type iFunc interface {
	In(i int) reflect.Type
	NumIn() int
	IsVariadic() bool
}

var (
	// ErrFuncArgNoReturnValue error.
	ErrFuncArgNoReturnValue = errors.New("function argument expression doesn't have return value")
	// ErrFuncArgCantBeMultiValue error.
	ErrFuncArgCantBeMultiValue = errors.New("function argument expression can't be multi values")
)

func checkFuncCall(tfn iFunc, args []interface{}) (arity int) {
	if len(args) == 1 {
		targ := args[0].(iValue).TypeOf()
		return targ.NumValues()
	}
	for _, arg := range args {
		targ := arg.(iValue).TypeOf()
		n := targ.NumValues()
		if n != 1 {
			if n == 0 {
				log.Fatalln("checkFuncCall:", ErrFuncArgNoReturnValue)
			} else {
				log.Fatalln("checkFuncCall:", ErrFuncArgCantBeMultiValue)
			}
		}
	}
	return len(args)
}

// -----------------------------------------------------------------------------
