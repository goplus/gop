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
	"math/big"
	"reflect"
	"strings"
	"unsafe"

	"github.com/goplus/gop/ast/astutil"
	"github.com/goplus/gop/constant"
	"github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/token"
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

type shiftValue struct {
	x           *constVal
	r           exec.Reserved
	fnCheckType func(typ reflect.Type)
}

func (p *shiftValue) checkType(t reflect.Type) {
	p.fnCheckType(t)
}

func (p *shiftValue) bound(t reflect.Type) {
	v := boundConst(p.x, t)
	p.x.v = v
	p.x.kind = t.Kind()
}

func (p *shiftValue) Kind() iKind {
	return p.x.kind
}

func (p *shiftValue) Type() reflect.Type {
	return boundType(p.x)
}

func (p *shiftValue) NumValues() int {
	return 1
}

func (p *shiftValue) Value(i int) iValue {
	return p
}

type goValue struct {
	t reflect.Type
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

func v2cv(x interface{}) interface{} {
	v := reflect.ValueOf(x)
	kind := v.Kind()
	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return constant.MakeInt64(v.Int())
	case kind >= reflect.Uint && kind <= reflect.Uintptr:
		return constant.MakeUint64(v.Uint())
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return constant.MakeFloat64(v.Float())
	case kind >= reflect.Complex64 && kind <= reflect.Complex128:
		c := v.Complex()
		r := constant.MakeFloat64(real(c))
		i := constant.MakeFloat64(imag(c))
		return constant.BinaryOp(r, token.ADD, constant.MakeImag(i))
	}
	return x
}

type constVal struct {
	v       interface{}
	kind    iKind
	reserve exec.Reserved
}

func newConstVal(v interface{}, kind iKind) *constVal {
	return &constVal{v: v2cv(v), kind: kind, reserve: exec.InvalidReserved}
}

func (p *constVal) Kind() iKind {
	return p.kind
}

func (p *constVal) Type() reflect.Type {
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
		return reflect.Int
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

func unaryOp(xop token.Token, op exec.Operator, x *constVal) *constVal {
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
	if cv, ok := x.v.(constant.Value); ok && !isConstBound(xkind) {
		v := constant.UnaryOp(xop, cv, 0)
		return &constVal{kind: xkind, v: v, reserve: -1}
	}
	t := exec.TypeFromKind(kindReal)
	vx := boundConst(x, t)
	v := CallBuiltinOp(kindReal, op, vx)
	if op == exec.OpNeg && xkind >= reflect.Uint && xkind <= reflect.Uintptr {
		log.Panicf("constant -%v overflows %v", vx, xkind)
	}
	return newConstVal(v, xkind)
}

func binaryOp(xop token.Token, op exec.Operator, x, y *constVal) *constVal {
	i := op.GetInfo()
	xkind := x.kind
	ykind := y.kind
	var kind, kindReal astutil.ConstKind
	if op == exec.OpLsh || op == exec.OpRsh {
		kind, kindReal = xkind, realKindOf(xkind)
	} else if isConstBound(xkind) {
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
	cx, xok := x.v.(constant.Value)
	cy, yok := y.v.(constant.Value)
	if xok && yok {
		var v constant.Value
		switch xop {
		case token.SHL, token.SHR:
			var s uint
			if !isConstBound(y.kind) {
				i := boundConst(y, exec.TyInt).(int)
				if i < 0 {
					log.Panicf("invalid negative shift count: %v", i)
				}
				s = uint(i)
			} else {
				s = boundConst(y, exec.TyUint).(uint)
			}
			v = constant.Shift(cx, xop, s)
		case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
			if isBoundNumberType(xkind) || isBoundNumberType(ykind) {
				goto bound
			}
			b := constant.Compare(cx, xop, cy)
			return &constVal{kind: kind, v: b, reserve: -1}
		default:
			if (kind >= exec.Int && kind <= exec.Uintptr) || kind == exec.BigInt {
				if !isBoundNumberType(xkind) {
					cx = extractUnboundInt(cx, xkind)
				}
				if !isBoundNumberType(ykind) {
					cy = extractUnboundInt(cy, ykind)
				}
			}
			v = constant.BinaryOp(cx, xop, cy)
		}
		if kind == exec.ConstUnboundInt {
			v = extractUnboundInt(v, kind)
		}
		cv := &constVal{kind: kind, v: v, reserve: -1}
		if isBoundNumberType(kind) {
			return newConstVal(boundConst(cv, t), kind)
		}
		return cv
	}
bound:
	vx := boundConst(x, t)
	vy := boundConst(y, t)
	v := CallBuiltinOp(kindReal, op, vx, vy)
	return &constVal{kind: kind, v: v, reserve: -1}
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

func intKind(kind reflect.Kind) reflect.Kind {
	if kind == reflect.Interface {
		return []reflect.Kind{reflect.Int32, reflect.Int64}[ptrSizeIndex]
	}
	return kind
}

func floatKind(kind reflect.Kind) reflect.Kind {
	if kind == reflect.Interface {
		return reflect.Float64
	}
	return kind
}

func complexKind(kind reflect.Kind) reflect.Kind {
	if kind == reflect.Interface {
		return reflect.Complex128
	}
	return kind
}

func extractUnboundInt(cv constant.Value, kind exec.Kind) constant.Value {
	switch val := constant.Val(cv).(type) {
	case *big.Rat:
		if kind == exec.ConstUnboundInt {
			s := val.FloatString(1)
			pos := strings.Index(s, ".")
			cv = constant.MakeFromLiteral(s[:pos], token.INT, 0)
		} else if kind == exec.ConstUnboundFloat {
			if !val.IsInt() {
				log.Panicf("constant %v truncated to integer", cv)
			}
			cv = constant.ToInt(cv)
		}
	case *big.Float:
		if kind == exec.ConstUnboundInt {
			s := val.String()
			pos := strings.Index(s, ".")
			cv = constant.MakeFromLiteral(s[:pos], token.INT, 0)
		} else if kind == exec.ConstUnboundFloat {
			if !val.IsInt() {
				log.Panicf("constant %v truncated to integer", cv)
			}
			cv = constant.ToInt(cv)
			if cv.Kind() == constant.Unknown {
				log.Panicf("integer too large")
			}
		}
	}
	return cv
}

func constantValue(cv constant.Value, ckind exec.Kind, kind reflect.Kind) interface{} {
	if cv.Kind() == constant.Complex {
		cr := constant.Real(cv)
		ci := constant.Imag(cv)
		if kind == reflect.Interface || kind == reflect.Complex64 || kind == reflect.Complex128 {
			kind = complexKind(kind)
			if doesOverflow(cr, kind) {
				log.Panicf("constant %v overflows %v", cv, kind)
			}
			if doesOverflow(ci, kind) {
				log.Panicf("constant %v overflows %v", cv, kind)
			}
			r, _ := constant.Float64Val(cr)
			i, _ := constant.Float64Val(ci)
			return complex(r, i)
		}
		i, _ := constant.Float64Val(ci)
		if i != 0 {
			log.Panicf("constant %v truncated to %v", cv, kind)
		}
		cv = cr
	}
	v := constant.Val(cv)
	switch val := v.(type) {
	case *big.Rat:
		if ckind == exec.ConstUnboundInt {
			s := val.FloatString(1)
			pos := strings.Index(s, ".")
			cv = constant.MakeFromLiteral(s[:pos], token.INT, 0)
			v = constant.Val(cv)
		} else if kind >= exec.Int && kind <= exec.Uintptr || kind == exec.BigInt {
			if val.IsInt() {
				v = val.Num()
			} else if ckind == exec.ConstUnboundFloat {
				log.Panicf("constant %v truncated to integer", cv)
			} else {
				s := val.FloatString(1)
				pos := strings.Index(s, ".")
				cv = constant.MakeFromLiteral(s[:pos], token.INT, 0)
				v = constant.Val(cv)
			}
		}
	case *big.Float:
		if ckind == exec.ConstUnboundInt {
			s := val.String()
			pos := strings.Index(s, ".")
			cv = constant.MakeFromLiteral(s[:pos], token.INT, 0)
			v = constant.Val(cv)
		} else if kind >= exec.Int && kind <= exec.Uintptr || kind == exec.BigInt {
			if !val.IsInt() {
				log.Panicf("constant %v truncated to integer", cv)
			}
			v, _ = val.Int(nil)
			cv = constant.Make(v)
		}
	}
	switch val := v.(type) {
	case int64:
		kind = intKind(kind)
		if doesOverflow(cv, kind) {
			log.Panicf("constant %v overflows %v", cv, kind)
		}
	case *big.Int:
		kind = intKind(kind)
		if kind >= exec.Int && kind <= exec.Uintptr {
			if doesOverflow(cv, kind) {
				log.Panicf("constant %v overflows %v", cv, kind)
			}
			v = val.Int64()
		} else if kind >= exec.Float32 && kind <= exec.Float64 {
			if doesOverflow(cv, kind) {
				log.Panicf("constant %v overflows %v", cv, kind)
			}
			v, _ = new(big.Float).SetInt(val).Float64()
		}
	case *big.Float:
		kind = floatKind(kind)
		if kind >= exec.Float32 && kind <= exec.Complex128 {
			if doesOverflow(cv, kind) {
				log.Panicf("constant %v overflows %v", cv, kind)
			}
			v, _ = val.Float64()
		}
	case *big.Rat:
		kind = floatKind(kind)
		if kind >= exec.Float32 && kind <= exec.Complex128 {
			if doesOverflow(cv, kind) {
				log.Panicf("constant %v overflows %v", cv, kind)
			}
			v, _ = val.Float64()
		}
	}
	return v
}

func boundConst(x *constVal, t reflect.Type) interface{} {
	v := x.v
	kind := kindOf(t)
	if v == nil {
		if kind >= reflect.Chan && kind <= reflect.Slice {
			return reflect.Zero(t).Interface()
		} else if kind == reflect.UnsafePointer {
			return reflect.ValueOf(unsafe.Pointer(nil)).Interface()
		}
		log.Panicln("boundConst: can't convert nil into", t)
	}
	if cv, ok := v.(constant.Value); ok {
		v = constantValue(cv, x.kind, kind)
	}
	sval := reflect.ValueOf(v)
	st := sval.Type()
	if t == st {
		return v
	}
	if kind == reflect.Interface {
		if sval.Kind() == reflect.Int64 {
			return int(sval.Int())
		}
	} else if kind == reflect.Complex128 || kind == reflect.Complex64 {
		if skind := sval.Kind(); skind >= reflect.Int && skind <= reflect.Float64 {
			fval := sval.Convert(exec.TyFloat64).Float()
			if kind == reflect.Complex64 {
				return complex(float32(fval), 0)
			} else {
				return complex(fval, 0)
			}
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
