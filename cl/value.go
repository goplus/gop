package cl

import (
	"reflect"

	"github.com/qiniu/qlang/ast/astutil"
	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/x/log"
)

type iKind = astutil.ConstKind

// iValue represents a qlang value(s).
//  - *goFunc
//  - *goValue
//  - *nonValue
//  - *constVal
//  - *funcResult
type iValue interface {
	Type() reflect.Type
	Kind() iKind
	Value(i int) iValue
	NumValues() int
}

// -----------------------------------------------------------------------------

type goValue struct {
	t reflect.Type
}

func (p *goValue) Kind() iKind {
	return p.t.Kind()
}

func (p *goValue) Type() reflect.Type {
	return p.t
}

func (p *goValue) NumValues() int {
	return 1
}

func (p *goValue) Value(i int) iValue {
	return p
}

// -----------------------------------------------------------------------------

type nonValue struct {
	v interface{} // *exec.GoPackage, type, etc.
}

func (p *nonValue) Kind() iKind {
	return reflect.Invalid
}

func (p *nonValue) Type() reflect.Type {
	return nil
}

func (p *nonValue) NumValues() int {
	return 0
}

func (p *nonValue) Value(i int) iValue {
	return p
}

// -----------------------------------------------------------------------------

type funcResults struct {
	tfn reflect.Type
}

func (p *funcResults) Kind() iKind {
	panic("don't call me")
}

func (p *funcResults) Type() reflect.Type {
	panic("don't call me")
}

func (p *funcResults) NumValues() int {
	return p.tfn.NumOut()
}

func (p *funcResults) Value(i int) iValue {
	return &goValue{t: p.tfn.Out(i)}
}

func newFuncResults(tfn reflect.Type) iValue {
	if tfn.NumOut() == 1 {
		return &goValue{t: tfn.Out(0)}
	}
	return &funcResults{tfn: tfn}
}

// -----------------------------------------------------------------------------

type qlFunc funcDecl

func newQlFunc(f *funcDecl) *qlFunc {
	return (*qlFunc)(f)
}

func (p *qlFunc) Kind() iKind {
	return reflect.Func
}

func (p *qlFunc) Type() reflect.Type {
	return ((*funcDecl)(p)).typeOf()
}

func (p *qlFunc) NumValues() int {
	return 1
}

func (p *qlFunc) Value(i int) iValue {
	return p
}

func (p *qlFunc) Results() iValue {
	return newFuncResults(p.Type())
}

func (p *qlFunc) Proto() iFuncType {
	return p.Type()
}

// -----------------------------------------------------------------------------

type goFunc struct {
	v        *exec.GoFuncInfo
	addr     uint32
	kind     exec.SymbolKind
	isMethod int // 0 - global func, 1 - method
}

func newGoFunc(addr uint32, kind exec.SymbolKind, isMethod int) *goFunc {
	var fi *exec.GoFuncInfo
	switch kind {
	case exec.SymbolFunc:
		fi = exec.GoFuncAddr(addr).GetInfo()
	case exec.SymbolFuncv:
		fi = exec.GoFuncvAddr(addr).GetInfo()
	default:
		log.Panicln("getGoFunc: unknown -", kind, addr)
	}
	return &goFunc{v: fi, addr: addr, kind: kind, isMethod: isMethod}
}

func (p *goFunc) Kind() iKind {
	return reflect.Func
}

func (p *goFunc) Type() reflect.Type {
	return reflect.TypeOf(p.v.This)
}

func (p *goFunc) NumValues() int {
	return 1
}

func (p *goFunc) Value(i int) iValue {
	return p
}

func (p *goFunc) Results() iValue {
	return newFuncResults(reflect.TypeOf(p.v.This))
}

func (p *goFunc) Proto() iFuncType {
	return reflect.TypeOf(p.v.This)
}

// -----------------------------------------------------------------------------

type constVal struct {
	v       interface{}
	kind    iKind
	reserve exec.Reserved
}

func (p *constVal) Kind() iKind {
	return p.kind
}

func (p *constVal) Type() reflect.Type {
	if astutil.IsConstBound(p.kind) {
		return exec.TypeFromKind(p.kind)
	}
	panic("don't call constVal.TypeOf: unbounded")
}

func (p *constVal) NumValues() int {
	return 1
}

func (p *constVal) Value(i int) iValue {
	return p
}

func (p *constVal) boundKind() reflect.Kind {
	if astutil.IsConstBound(p.kind) {
		return p.kind
	}
	switch p.kind {
	case astutil.ConstUnboundInt:
		if _, ok := p.v.(int64); ok {
			return reflect.Int
		}
		return reflect.Uint
	case astutil.ConstUnboundFloat:
		return reflect.Float64
	case astutil.ConstUnboundComplex:
		return reflect.Complex128
	}
	log.Panicln("boundKind: unexpected type kind -", p.kind)
	return reflect.Invalid
}

func (p *constVal) boundType() reflect.Type {
	return exec.TypeFromKind(p.boundKind())
}

func boundType(in iValue) reflect.Type {
	if v, ok := in.(*constVal); ok {
		return v.boundType()
	}
	return in.Type()
}

func (p *constVal) bound(t reflect.Type, b *exec.Builder) {
	if p.reserve == -1 { // bounded
		if p.kind != t.Kind() {
			if t == exec.TyEmptyInterface {
				return
			}
			log.Panicln("function call with invalid argument type: requires", t, ", but got", p.kind)
		}
		return
	}
	v, ok := boundConst(p.v, t)
	if !ok {
		log.Panicln("function call with invalid argument type: requires", t, ", but got", reflect.TypeOf(p.v))
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
		log.Panicln("binaryOp failed: invalid first argument type.")
	}
	t := exec.TypeFromKind(kindReal)
	vx, xok := boundConst(x.v, t)
	vy, yok := boundConst(y.v, t)
	if !xok || !yok {
		log.Panicln("binaryOp failed: invalid argument type -", t)
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
