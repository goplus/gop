package exec

import (
	"reflect"
)

// -----------------------------------------------------------------------------

func pushInt(stk *Context, kind reflect.Kind, v int64) {
	var val interface{}
	switch kind {
	case reflect.Int:
		val = int(v)
	case reflect.Int64:
		val = int64(v)
	case reflect.Int32:
		val = int32(v)
	case reflect.Int16:
		val = int16(v)
	case reflect.Int8:
		val = int8(v)
	default:
		panic("pushInt failed: invalid kind")
	}
	stk.Push(val)
}

func pushUint(stk *Context, kind reflect.Kind, v uint64) {
	var val interface{}
	switch kind {
	case reflect.Uint:
		val = uint(v)
	case reflect.Uint64:
		val = uint64(v)
	case reflect.Uint32:
		val = uint32(v)
	case reflect.Uint8:
		val = uint8(v)
	case reflect.Uint16:
		val = uint16(v)
	case reflect.Uintptr:
		val = uintptr(v)
	default:
		panic("pushUint failed: invalid kind")
	}
	stk.Push(val)
}

// -----------------------------------------------------------------------------

var valSpecs = []interface{}{
	false,
	true,
}

func execPushValSpec(i Instr, stk *Context) {
	stk.Push(valSpecs[i&bitsOperand])
}

func execPushStringR(i Instr, stk *Context) {
	v := stk.Code.stringConsts[i&bitsOperand]
	stk.Push(v)
}

func execPushIntR(i Instr, stk *Context) {
	v := stk.Code.intConsts[i&bitsOpIntOperand]
	kind := reflect.Int + reflect.Kind((i>>bitsOpIntOperand)&7)
	pushInt(stk, kind, v)
}

func execPushInt(i Instr, stk *Context) {
	v := int64(int32(i) << bitsOpIntShift >> bitsOpIntShift)
	kind := reflect.Int + reflect.Kind((i>>bitsOpIntOperand)&7)
	pushInt(stk, kind, v)
}

func execPushUintR(i Instr, stk *Context) {
	v := stk.Code.uintConsts[i&bitsOpIntOperand]
	kind := reflect.Uint + reflect.Kind((i>>bitsOpIntOperand)&7)
	pushUint(stk, kind, v)
}

func execPushUint(i Instr, stk *Context) {
	v := uint64(i & bitsOpIntOperand)
	kind := reflect.Uint + reflect.Kind((i>>bitsOpIntOperand)&7)
	pushUint(stk, kind, v)
}

func execPushFloatR(i Instr, stk *Context) {
	v := stk.Code.valConsts[i&bitsOpFloatOperand]
	stk.Push(v)
}

func execPushFloat(i Instr, stk *Context) {
	panic("execPushFloat: not impl")
}

// -----------------------------------------------------------------------------

func resolveConsts(p *Builder, code *Code) {
	var i Instr
	for val, vu := range p.valConsts {
		switch vu.op {
		case opPushStringR:
			i = (opPushStringR << bitsOpShift) | uint32(len(code.stringConsts))
			code.stringConsts = append(code.stringConsts, val.(string))
		case opPushIntR:
			v := reflect.ValueOf(val)
			kind := v.Kind()
			i = (opPushIntR << bitsOpShift) | (uint32(kind-reflect.Int) << bitsOpIntShift) | uint32(len(code.intConsts))
			code.intConsts = append(code.intConsts, v.Int())
		case opPushUintR:
			v := reflect.ValueOf(val)
			kind := v.Kind()
			i = (opPushUintR << bitsOpShift) | (uint32(kind-reflect.Uint) << bitsOpIntShift) | uint32(len(code.uintConsts))
			code.uintConsts = append(code.uintConsts, v.Uint())
		case opPushFloatR:
			v := reflect.ValueOf(val)
			kind := v.Kind()
			i = (opPushFloatR << bitsOpShift) | (uint32(kind-reflect.Float32) << bitsOpFloatShift) | uint32(len(code.valConsts))
			code.valConsts = append(code.valConsts, val)
		default:
			panic("Resolve failed: unknown type")
		}
		for _, off := range vu.offs {
			code.data[off] = i
		}
	}
}

// -----------------------------------------------------------------------------

func (ctx *Builder) pushUnresolved(op Instr, val interface{}) *Builder {
	vu, ok := ctx.valConsts[val]
	if !ok {
		vu = &valUnresolved{op: op}
		ctx.valConsts[val] = vu
	}
	p := ctx.code
	vu.offs = append(vu.offs, len(p.data))
	p.data = append(p.data, iPushUnresolved)
	return ctx
}

// Push instr
func (ctx *Builder) Push(val interface{}) *Builder {
	var i Instr
	v := reflect.ValueOf(val)
	kind := v.Kind()
	if kind >= reflect.Int && kind <= reflect.Int64 {
		iv := v.Int()
		ivStore := int64(int32(iv) << bitsOpIntShift >> bitsOpIntShift)
		if iv != ivStore {
			return ctx.pushUnresolved(opPushIntR, val)
		}
		i = (opPushInt << bitsOpShift) | (uint32(kind-reflect.Int) << bitsOpIntShift) | (uint32(iv) & bitsOpIntOperand)
	} else if kind >= reflect.Uint && kind <= reflect.Uintptr {
		iv := v.Uint()
		if iv != (iv & bitsOpIntOperand) {
			return ctx.pushUnresolved(opPushUintR, val)
		}
		i = (opPushUint << bitsOpShift) | (uint32(kind-reflect.Uint) << bitsOpIntShift) | (uint32(iv) & bitsOpIntOperand)
	} else if kind == reflect.Bool {
		if val.(bool) {
			i = iPushTrue
		} else {
			i = iPushFalse
		}
	} else if kind == reflect.String {
		return ctx.pushUnresolved(opPushStringR, val)
	} else if kind >= reflect.Float32 && kind <= reflect.Complex128 {
		return ctx.pushUnresolved(opPushFloatR, val)
	} else {
		panic("Push failed: unsupported type")
	}
	ctx.code.data = append(ctx.code.data, i)
	return ctx
}

// -----------------------------------------------------------------------------
