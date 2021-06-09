/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package cl

import (
	"math"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/goplus/gop/ast/astutil"
	"github.com/goplus/gop/exec.spec"
	"github.com/qiniu/x/log"
)

type iKind = astutil.ConstKind

// iValue represents a Go+ value(s).
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

func isBool(v iValue) bool {
	return v.NumValues() == 1 && v.Type() == exec.TyBool
}

// -----------------------------------------------------------------------------

type lshValue struct {
	x           *constVal
	r           exec.Reserved
	fnCheckType func(typ reflect.Type)
}

func (p *lshValue) checkType(t reflect.Type) {
	p.fnCheckType(t)
}

func (p *lshValue) bound(t reflect.Type) {
	v := boundConst(p.x, t)
	p.x.v = v
	p.x.kind = t.Kind()
}

func (p *lshValue) Kind() iKind {
	return p.x.kind
}

func (p *lshValue) Type() reflect.Type {
	return boundType(p.x)
}

func (p *lshValue) NumValues() int {
	return 1
}

func (p *lshValue) Value(i int) iValue {
	return p
}

type goValue struct {
	t reflect.Type
	c *constVal
}

func (p *goValue) Kind() iKind {
	return kindOf(p.t)
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
	v interface{} // *exec.GoPackage, goInstr, iType, etc.
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

type wrapValue struct {
	x iValue
}

func (p *wrapValue) Type() reflect.Type {
	if p.x.NumValues() != 2 {
		panic("don't call me")
	}
	return p.x.Value(0).Type()
}

func (p *wrapValue) Kind() iKind {
	if p.x.NumValues() != 2 {
		panic("don't call me")
	}
	return p.x.Value(0).Kind()
}

func (p *wrapValue) NumValues() int {
	return p.x.NumValues() - 1
}

func (p *wrapValue) Value(i int) iValue {
	return p.x.Value(i)
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

func (p *qlFunc) FuncInfo() exec.FuncInfo {
	return ((*funcDecl)(p)).Get()
}

func (p *qlFunc) Kind() iKind {
	return reflect.Func
}

func (p *qlFunc) Type() reflect.Type {
	return ((*funcDecl)(p)).Type()
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
	t        reflect.Type
	addr     uint32
	kind     exec.SymbolKind
	isMethod int // 0 - global func, 1 - method
}

func newGoFunc(addr uint32, kind exec.SymbolKind, isMethod int, ctx *blockCtx) *goFunc {
	var t reflect.Type
	switch kind {
	case exec.SymbolFunc:
		t = ctx.GetGoFuncType(exec.GoFuncAddr(addr))
	case exec.SymbolFuncv:
		t = ctx.GetGoFuncvType(exec.GoFuncvAddr(addr))
	default:
		log.Panicln("getGoFunc: unknown -", kind, addr)
	}
	return &goFunc{t: t, addr: addr, kind: kind, isMethod: isMethod}
}

func (p *goFunc) Kind() iKind {
	return reflect.Func
}

func (p *goFunc) Type() reflect.Type {
	return p.t
}

func (p *goFunc) NumValues() int {
	return 1
}

func (p *goFunc) Value(i int) iValue {
	return p
}

func (p *goFunc) Results() iValue {
	return newFuncResults(p.t)
}

func (p *goFunc) Proto() iFuncType {
	return p.t
}

// -----------------------------------------------------------------------------

// isConstBound checks a const is bound or not.
func isConstBound(kind astutil.ConstKind) bool {
	return astutil.IsConstBound(kind)
}

type constVal struct {
	v       interface{}
	kind    iKind
	reserve exec.Reserved
	typed   reflect.Type
}

func newConstVal(v interface{}, kind iKind) *constVal {
	return &constVal{v: v, kind: kind, reserve: exec.InvalidReserved}
}

func (c *constVal) IsTyped() bool {
	return c.typed != nil && c.typed.PkgPath() != ""
}

func (p *constVal) Kind() iKind {
	return p.kind
}

func (p *constVal) Type() reflect.Type {
	if p.IsTyped() {
		return p.typed
	}
	if isConstBound(p.kind) {
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
	if isConstBound(p.kind) {
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
		if v.IsTyped() {
			return v.typed
		}
		return v.boundType()
	}
	return in.Type()
}

func (p *constVal) bound(t reflect.Type, b exec.Builder) {
	kind := t.Kind()
	if p.reserve == exec.InvalidReserved { // bounded
		if p.kind != kind {
			if t == exec.TyEmptyInterface {
				return
			}
			log.Panicln("function call with invalid argument type: requires", t, ", but got", p.kind)
		}
		return
	}
	if kind == reflect.Interface && astutil.IsConstBound(p.kind) {
		t = exec.TypeFromKind(p.kind)
	}
	v := boundConst(p, t)
	p.v, p.kind = v, kind
	p.reserve.Push(b, v)
}

func unaryOp(op exec.Operator, x *constVal) *constVal {
	i := op.GetInfo()
	xkind := x.kind
	var kindReal astutil.ConstKind
	if isConstBound(xkind) {
		kindReal = xkind
	} else {
		kindReal = realKindOf(xkind)
	}
	if (i.InFirst & (1 << kindReal)) == 0 {
		log.Panicln("unaryOp failed: invalid argument type.")
	}
	t := exec.TypeFromKind(kindReal)
	vx := boundConstCheck(x, t, false)
	v := CallBuiltinOp(kindReal, op, vx)
	c := &constVal{kind: xkind, v: v, reserve: -1}
	c.typed = x.typed
	return c
}

func binaryOp(op exec.Operator, x, y *constVal) *constVal {
	i := op.GetInfo()
	xkind := x.kind
	ykind := y.kind
	var kind, kindReal astutil.ConstKind
	if isConstBound(xkind) {
		kind, kindReal = xkind, xkind
	} else if isConstBound(ykind) {
		kind, kindReal = ykind, ykind
	} else if xkind < ykind {
		kind, kindReal = ykind, realKindOf(ykind)
	} else {
		kind, kindReal = xkind, realKindOf(xkind)
	}
	if (i.InFirst & (1 << kindReal)) == 0 {
		if kindReal != exec.BigInt && op != exec.OpQuo {
			log.Panicln("binaryOp failed: invalid first argument type -", i, kindReal)
		}
		kind = exec.BigRat
	} else if i.Out != exec.SameAsFirst {
		kind = i.Out
	}
	t := exec.TypeFromKind(kindReal)
	vx := boundConstCheck(x, t, false)
	vy := boundConstCheck(y, t, false)
	v := CallBuiltinOp(kindReal, op, vx, vy)
	c := &constVal{kind: kind, v: v, reserve: -1}
	if !(op >= exec.OpLT && op <= exec.OpNENil) {
		c.typed = x.typed
	}
	return c
}

func kindOf(t reflect.Type) exec.Kind {
	kind := t.Kind()
	if kind == reflect.Ptr {
		switch t {
		case exec.TyBigRat:
			return exec.BigRat
		case exec.TyBigInt:
			return exec.BigInt
		case exec.TyBigFloat:
			return exec.BigFloat
		}
	}
	return kind
}

const (
	intSize = strconv.IntSize
)

func isOverflowsIntByInt64(v int64, intSize int) bool {
	if intSize == 32 {
		return v < int64(math.MinInt32) || v > int64(math.MaxInt32)
	}
	return false
}

func isOverflowsIntByUint64(v uint64, intSize int) bool {
	if intSize == 32 {
		return v > uint64(math.MaxInt32)
	}
	return v > uint64(math.MaxInt64)
}

func boundConst(c *constVal, t reflect.Type) interface{} {
	return boundConstCheck(c, t, true)
}

func boundConstCheck(c *constVal, t reflect.Type, chkTyped bool) interface{} {
	kind := kindOf(t)
	if c.v == nil {
		if kind >= reflect.Chan && kind <= reflect.Slice {
			return reflect.Zero(t).Interface()
		} else if kind == reflect.UnsafePointer {
			return reflect.ValueOf(unsafe.Pointer(nil)).Interface()
		}
		log.Panicln("boundConst: can't convert nil into", t)
	}
	sval := reflect.ValueOf(c.v)
	if chkTyped && c.IsTyped() {
		return sval.Convert(c.typed).Interface()
	}
	st := sval.Type()
	if t == st {
		return c.v
	}
	if kind == reflect.Interface {
		skind := sval.Kind()
		if skind == reflect.Uint64 {
			v := sval.Uint()
			if isOverflowsIntByUint64(v, intSize) {
				log.Panicf("constant %v overflows int\n", v)
			}
			return int(v)
		} else if skind == reflect.Int64 {
			v := sval.Int()
			if isOverflowsIntByInt64(v, intSize) {
				log.Panicf("constant %v overflows int\n", v)
			}
			return int(v)
		}
	} else if kind == reflect.Complex128 || kind == reflect.Complex64 {
		if skind := sval.Kind(); skind >= reflect.Int && skind <= reflect.Float64 {
			fval := sval.Convert(exec.TyFloat64).Float()
			return complex(fval, 0)
		}
	} else if kind >= exec.BigInt {
		val := reflect.New(t.Elem())
		skind := kindOf(st)
		switch {
		case skind >= reflect.Int && skind <= reflect.Int64:
			sval = sval.Convert(exec.TyInt64)
			val.MethodByName("SetInt64").Call([]reflect.Value{sval})
		case skind >= reflect.Uint && skind <= reflect.Uintptr:
			sval = sval.Convert(exec.TyUint64)
			val.MethodByName("SetUint64").Call([]reflect.Value{sval})
		case skind >= reflect.Float32 && skind <= reflect.Float64:
			sval = sval.Convert(exec.TyFloat64)
			val.MethodByName("SetFloat64").Call([]reflect.Value{sval})
		case skind == exec.BigInt:
			val.MethodByName("SetInt").Call([]reflect.Value{sval})
		case skind == exec.BigFloat:
			val.MethodByName("SetRat").Call([]reflect.Value{sval})
		default:
			log.Panicln("boundConst: convert type failed -", skind)
		}
		return val.Interface()
	}
	return sval.Convert(t).Interface()
}

func constIsConvertible(v interface{}, t reflect.Type) bool {
	styp := reflect.TypeOf(v)
	skind := styp.Kind()
	switch kind := t.Kind(); kind {
	case reflect.String:
		return skind == reflect.String
	case reflect.Complex128, reflect.Complex64:
		return skind >= reflect.Int && skind <= reflect.Complex128
	}
	return styp.ConvertibleTo(t)
}

func realKindOf(kind astutil.ConstKind) reflect.Kind {
	switch kind {
	case astutil.ConstUnboundInt:
		return reflect.Int64
	case astutil.ConstUnboundFloat:
		return reflect.Float64
	case astutil.ConstUnboundComplex:
		return reflect.Complex128
	case astutil.ConstUnboundPtr:
		return reflect.UnsafePointer
	default:
		return kind
	}
}

func boundElementType(elts []interface{}, base, max, step int) reflect.Type {
	var tBound reflect.Type
	var kindUnbound iKind
	for i := base; i < max; i += step {
		e := elts[i].(iValue)
		if e.NumValues() != 1 {
			log.Panicln("boundElementType: unexpected - multiple return values.")
		}
		kind := e.Kind()
		if !isConstBound(kind) { // unbound
			if kindUnbound < kind {
				kindUnbound = kind
			}
		} else {
			if t := e.Type(); tBound != t {
				if tBound != nil { // mismatched type
					return nil
				}
				tBound = t
			}
		}
	}
	if tBound != nil {
		for i := base; i < max; i += step {
			if e, ok := elts[i].(*constVal); ok {
				if !constIsConvertible(e.v, tBound) { // mismatched type
					return nil
				}
			}
		}
		return tBound
	}
	var kindBound iKind
	for i := base; i < max; i += step {
		if e, ok := elts[i].(*constVal); ok && e.kind == kindUnbound {
			kind := e.boundKind()
			if kind != kindBound {
				if kindBound != 0 { // mismatched type
					return nil
				}
				kindBound = kind
			}
		}
	}
	return exec.TypeFromKind(kindBound)
}

// -----------------------------------------------------------------------------
