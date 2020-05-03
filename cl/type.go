package cl

import (
	"errors"
	"reflect"

	"github.com/qiniu/qlang/ast/astutil"
	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/qlang/token"
	"github.com/qiniu/x/log"
)

// iType represents a qlang type.
//  - *goType
//  - *retType
//  - unboundType
type iType interface {
	TypeValue(i int) iType
	NumValues() int
	Kind() reflect.Kind
}

// iValue represents a qlang value.
//  - *goFunc
//  - *constVal
//  - *exprResult
//  - *operatorResult
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

func (p *goType) Kind() reflect.Kind {
	return p.t.Kind()
}

func (p *goType) TypeValue(i int) iType {
	return p
}

func (p *goType) NumValues() int {
	return 1
}

// -----------------------------------------------------------------------------

type unboundType astutil.ConstKind

func (p unboundType) Kind() reflect.Kind {
	return reflect.Kind(p)
}

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

func (p *retType) Kind() reflect.Kind {
	return reflect.Invalid
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

func binaryOp(op exec.Operator, x, y *constVal) *constVal {
	i := op.GetInfo()
	xkind := x.kind
	ykind := y.kind
	var kind, kindReal astutil.ConstKind
	if astutil.IsConstBound(xkind) {
		kind, kindReal = xkind, xkind
	} else if astutil.IsConstBound(ykind) {
		kind, kindReal = ykind, ykind
	} else if xkind < ykind {
		kind, kindReal = ykind, realKindOf(ykind)
	} else {
		kind, kindReal = xkind, realKindOf(xkind)
	}
	if (i.InFirst & (1 << kindReal)) == 0 {
		log.Fatalln("binaryOp failed: invalid first argument type.")
	}
	t := exec.TypeFromKind(kindReal)
	vx, xok := astutil.ConstBound(x.v, t)
	vy, yok := astutil.ConstBound(y.v, t)
	if !xok || !yok {
		log.Fatalln("binaryOp failed: invalid argument type -", t)
	}
	v := exec.CallBuiltinOp(kindReal, op, vx, vy)
	return &constVal{kind: kind, v: v, reserve: -1}
}

func realKindOf(kind astutil.ConstKind) reflect.Kind {
	switch kind {
	case astutil.ConstUnboundInt:
		return reflect.Int64
	case astutil.ConstUnboundFloat:
		return reflect.Float64
	case astutil.ConstUnboundComplex:
		return reflect.Complex128
	default:
		return kind
	}
}

// -----------------------------------------------------------------------------

type exprResult struct {
	results iType
}

func (p *exprResult) TypeOf() iType {
	return p.results
}

// -----------------------------------------------------------------------------

type operatorResult struct {
	kind exec.Kind
	op   exec.Operator
}

func (p *operatorResult) TypeOf() iType {
	op := p.op.GetInfo()
	var kind exec.Kind
	if op.Out == exec.SameAsFirst {
		kind = p.kind
	} else {
		kind = op.Out
	}
	return &goType{t: exec.TypeFromKind(kind)}
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

func checkBinaryOp(kind exec.Kind, op exec.Operator, x, y interface{}, b *exec.Builder) {
	if xcons, xok := x.(*constVal); xok {
		if xcons.reserve != -1 {
			xv, ok := astutil.ConstBound(xcons.v, exec.TypeFromKind(kind))
			if !ok {
				log.Fatalln("checkBinaryOp: invalid operator", kind, "argument type.")
			}
			xcons.reserve.Push(b, xv)
		}
	}
	if ycons, yok := y.(*constVal); yok {
		i := op.GetInfo()
		if i.InSecond != (1 << exec.SameAsFirst) {
			if (uint64(ycons.kind) & i.InSecond) == 0 {
				log.Fatalln("checkBinaryOp: invalid operator", kind, "argument type.")
			}
			kind = ycons.kind
		}
		if ycons.reserve != -1 {
			yv, ok := astutil.ConstBound(ycons.v, exec.TypeFromKind(kind))
			if !ok {
				log.Fatalln("checkBinaryOp: invalid operator", kind, "argument type.")
			}
			ycons.reserve.Push(b, yv)
		}
	}
}

func inferBinaryOp(tok token.Token, x, y interface{}) (kind exec.Kind, op exec.Operator) {
	tx := x.(iValue).TypeOf()
	ty := y.(iValue).TypeOf()
	op = binaryOps[tok]
	if tx.NumValues() != 1 || ty.NumValues() != 1 {
		log.Fatalln("inferBinaryOp failed: argument isn't an expr.")
	}
	kind = tx.TypeValue(0).Kind()
	if !astutil.IsConstBound(kind) {
		kind = ty.TypeValue(0).Kind()
		if !astutil.IsConstBound(kind) {
			log.Fatalln("inferBinaryOp failed: expect x, y aren't const values either.")
		}
	}
	return
}

// -----------------------------------------------------------------------------
