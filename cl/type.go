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

type iFuncType interface {
	In(i int) reflect.Type
	Out(i int) reflect.Type
	NumIn() int
	NumOut() int
	IsVariadic() bool
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
	tfn iFuncType
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

func (p *goFunc) Proto() iFuncType {
	return reflect.TypeOf(p.v.This)
}

func (p *goFunc) TypeOf() iType {
	return &goType{t: reflect.TypeOf(p.v.This)}
}

func (p *goFunc) TypesOut() *retType {
	return &retType{tfn: p.Proto()}
}

func newGoFunc(addr uint32, kind exec.SymbolKind) *goFunc {
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

var (
	typEmptyInterface = reflect.TypeOf((*interface{})(nil)).Elem()
)

type constVal struct {
	v       interface{}
	kind    astutil.ConstKind
	reserve exec.Reserved
}

func (p *constVal) TypeOf() iType {
	if astutil.IsConstBound(p.kind) {
		return &goType{t: reflect.TypeOf(p.v)}
	}
	return unboundType(p.kind)
}

func (p *constVal) bound(t reflect.Type, b *exec.Builder) {
	if p.reserve == -1 { // bounded
		if p.kind != t.Kind() {
			if t == typEmptyInterface {
				return
			}
			log.Fatalln("function call with invalid argument type: requires", t, ", but got", p.kind)
		}
		return
	}
	var v interface{}
	if t == typEmptyInterface {
		switch nv := p.v.(type) {
		case int64:
			v = int(nv)
		case uint64:
			v = uint(nv)
		default:
			v = p.v
		}
	} else {
		if nv, ok := astutil.ConstBound(p.v, t); ok {
			v = nv
		} else {
			log.Fatalln("function call with invalid argument type: requires", t, ", but got", reflect.TypeOf(p.v))
		}
	}
	p.reserve.Push(b, v)
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

var (
	// ErrFuncArgNoReturnValue error.
	ErrFuncArgNoReturnValue = errors.New("function argument expression doesn't have return value")
	// ErrFuncArgCantBeMultiValue error.
	ErrFuncArgCantBeMultiValue = errors.New("function argument expression can't be multi values")
)

func checkFuncCall(tfn iFuncType, args []interface{}, b *exec.Builder) (arity int) {
	narg := tfn.NumIn()
	variadic := tfn.IsVariadic()
	if variadic {
		narg--
	}
	if len(args) == 1 {
		targ := args[0].(iValue).TypeOf()
		n := targ.NumValues()
		if n != 1 { // TODO
			return n
		}
	}
	for idx, arg := range args {
		var treq reflect.Type
		if variadic && idx >= narg {
			treq = tfn.In(narg).Elem()
		} else {
			treq = tfn.In(idx)
		}
		if v, ok := arg.(*constVal); ok {
			v.bound(treq, b)
		}
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
