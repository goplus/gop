package cl

import (
	"reflect"

	"github.com/qiniu/qlang/ast/astutil"
	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/x/log"
)

// iValue represents a qlang value(s).
//  - *goFunc
//  - *goValue
//  - *constVal
//  - *funcResult
//  - *operatorResult
type iValue interface {
	TypeOf() iType
	Value(i int) iValue
	NumValues() int
}

// -----------------------------------------------------------------------------
/*
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
*/
// -----------------------------------------------------------------------------

type goValue struct {
	t reflect.Type
}

func (p *goValue) TypeOf() iType {
	return p.t
}

func (p *goValue) NumValues() int {
	return 1
}

func (p *goValue) Value(i int) iValue {
	return p
}

// -----------------------------------------------------------------------------

type funcResult struct {
	tfn reflect.Type
}

func (p *funcResult) TypeOf() iType {
	panic("don't call me")
}

func (p *funcResult) NumValues() int {
	return p.tfn.NumOut()
}

func (p *funcResult) Value(i int) iValue {
	return &goValue{t: p.tfn.Out(i)}
}

// -----------------------------------------------------------------------------

type goFunc struct {
	v    *exec.GoFuncInfo
	addr uint32
	kind exec.SymbolKind
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

func (p *goFunc) TypeOf() iType {
	return reflect.TypeOf(p.v.This)
}

func (p *goFunc) NumValues() int {
	return 1
}

func (p *goFunc) Value(i int) iValue {
	return p
}

func (p *goFunc) Results() *funcResult {
	return &funcResult{tfn: reflect.TypeOf(p.v.This)}
}

func (p *goFunc) Proto() iFuncType {
	return reflect.TypeOf(p.v.This)
}

// -----------------------------------------------------------------------------

type constVal struct {
	v       interface{}
	kind    astutil.ConstKind
	reserve exec.Reserved
}

func (p *constVal) TypeOf() iType {
	if astutil.IsConstBound(p.kind) {
		return exec.TypeFromKind(p.kind)
	}
	return unboundType(p.kind)
}

func (p *constVal) NumValues() int {
	return 1
}

func (p *constVal) Value(i int) iValue {
	return p
}

func (p *constVal) bound(t reflect.Type, b *exec.Builder) {
	if p.reserve == -1 { // bounded
		if p.kind != t.Kind() {
			if t == exec.TyEmptyInterface {
				return
			}
			log.Fatalln("function call with invalid argument type: requires", t, ", but got", p.kind)
		}
		return
	}
	v, ok := boundConst(p.v, t)
	if !ok {
		log.Fatalln("function call with invalid argument type: requires", t, ", but got", reflect.TypeOf(p.v))
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
	vx, xok := boundConst(x.v, t)
	vy, yok := boundConst(y.v, t)
	if !xok || !yok {
		log.Fatalln("binaryOp failed: invalid argument type -", t)
	}
	v := exec.CallBuiltinOp(kindReal, op, vx, vy)
	return &constVal{kind: kind, v: v, reserve: -1}
}

func boundConst(v interface{}, t reflect.Type) (ret interface{}, ok bool) {
	if t == exec.TyEmptyInterface {
		switch nv := v.(type) {
		case int64:
			return int(nv), true
		case uint64:
			return uint(nv), true
		}
		return v, true
	}
	nkind := t.Kind()
	nv := reflect.New(t).Elem()
	if nkind >= reflect.Int && nkind <= reflect.Int64 {
		switch ov := v.(type) {
		case int64:
			nv.SetInt(ov)
		case uint64:
			nv.SetInt(int64(ov))
		default:
			return nil, false
		}
	} else if nkind >= reflect.Uint && nkind <= reflect.Uintptr {
		switch ov := v.(type) {
		case int64:
			nv.SetUint(uint64(ov))
		case uint64:
			nv.SetUint(ov)
		default:
			return nil, false
		}
	} else if nkind == reflect.Float64 || nkind == reflect.Float32 {
		switch ov := v.(type) {
		case float64:
			nv.SetFloat(ov)
		case int64:
			nv.SetFloat(float64(ov))
		case uint64:
			nv.SetFloat(float64(ov))
		default:
			return nil, false
		}
	} else if nkind == reflect.Complex128 || nkind == reflect.Complex64 {
		switch ov := v.(type) {
		case complex128:
			nv.SetComplex(ov)
		case float64:
			nv.SetComplex(complex(float64(ov), 0))
		case int64:
			nv.SetComplex(complex(float64(ov), 0))
		case uint64:
			nv.SetComplex(complex(float64(ov), 0))
		default:
			return nil, false
		}
	} else {
		return nil, false
	}
	return nv.Interface(), true
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
