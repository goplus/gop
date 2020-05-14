package exec

import (
	"reflect"

	"github.com/qiniu/x/log"
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
		log.Panicln("pushInt failed: invalid kind -", kind)
	}
	stk.Push(val)
}

func pushInt32(stk *Context, kind reflect.Kind, v int32) {
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
		log.Panicln("pushInt failed: invalid kind -", kind)
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
		log.Panicln("pushUint failed: invalid kind -", kind)
	}
	stk.Push(val)
}

func pushUint32(stk *Context, kind reflect.Kind, v uint32) {
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
		log.Panicln("pushUint failed: invalid kind -", kind)
	}
	stk.Push(val)
}

// -----------------------------------------------------------------------------

var valSpecs = []interface{}{
	false,
	true,
	nil,
}

func execPushValSpec(i Instr, stk *Context) {
	stk.Push(valSpecs[i&bitsOperand])
}

func execPushStringR(i Instr, stk *Context) {
	v := stk.code.stringConsts[i&bitsOperand]
	stk.Push(v)
}

func execPushIntR(i Instr, stk *Context) {
	v := stk.code.intConsts[i&bitsOpIntOperand]
	kind := reflect.Int + reflect.Kind((i>>bitsOpIntShift)&7)
	pushInt(stk, kind, v)
}

func execPushInt(i Instr, stk *Context) {
	v := int32(i) << bitsOpInt >> bitsOpInt
	kind := reflect.Int + reflect.Kind((i>>bitsOpIntShift)&7)
	pushInt32(stk, kind, v)
}

func execPushUintR(i Instr, stk *Context) {
	v := stk.code.uintConsts[i&bitsOpIntOperand]
	kind := reflect.Uint + reflect.Kind((i>>bitsOpIntShift)&7)
	pushUint(stk, kind, v)
}

func execPushUint(i Instr, stk *Context) {
	v := i & bitsOpIntOperand
	kind := reflect.Uint + reflect.Kind((i>>bitsOpIntShift)&7)
	pushUint32(stk, kind, v)
}

func execPushFloatR(i Instr, stk *Context) {
	v := stk.code.valConsts[i&bitsOpFloatOperand]
	stk.Push(v)
}

func execPop(i Instr, stk *Context) {
	n := len(stk.data) - int(i&bitsOperand)
	stk.data = stk.data[:n]
}

// -----------------------------------------------------------------------------

// Push instr
func (p *Builder) pushInstr(val interface{}, off int) (i Instr) {
	if val == nil {
		return iPushNil
	}
	v := reflect.ValueOf(val)
	kind := v.Kind()
	if kind >= reflect.Int && kind <= reflect.Int64 {
		iv := v.Int()
		ivStore := int64(int32(iv) << bitsOpInt >> bitsOpInt)
		if iv != ivStore {
			code := p.code
			i = (opPushIntR << bitsOpShift) | (uint32(kind-reflect.Int) << bitsOpIntShift) | uint32(len(code.intConsts))
			code.intConsts = append(code.intConsts, iv)
			return
		}
		i = (opPushInt << bitsOpShift) | (uint32(kind-reflect.Int) << bitsOpIntShift) | (uint32(iv) & bitsOpIntOperand)
	} else if kind >= reflect.Uint && kind <= reflect.Uintptr {
		iv := v.Uint()
		if iv != (iv & bitsOpIntOperand) {
			code := p.code
			i = (opPushUintR << bitsOpShift) | (uint32(kind-reflect.Uint) << bitsOpIntShift) | uint32(len(code.uintConsts))
			code.uintConsts = append(code.uintConsts, iv)
			return
		}
		i = (opPushUint << bitsOpShift) | (uint32(kind-reflect.Uint) << bitsOpIntShift) | (uint32(iv) & bitsOpIntOperand)
	} else if kind == reflect.Bool {
		if val.(bool) {
			i = iPushTrue
		} else {
			i = iPushFalse
		}
	} else if kind == reflect.String {
		code := p.code
		i = (opPushStringR << bitsOpShift) | uint32(len(code.stringConsts))
		code.stringConsts = append(code.stringConsts, val.(string))
	} else if kind >= reflect.Float32 && kind <= reflect.Complex128 {
		code := p.code
		i = (opPushFloatR << bitsOpShift) | (uint32(kind-reflect.Float32) << bitsOpFloatShift) | uint32(len(code.valConsts))
		code.valConsts = append(code.valConsts, val)
	} else {
		log.Panicln("Push failed: unsupported type:", reflect.TypeOf(val), "-", val)
	}
	return
}

// Push instr
func (p *Builder) Push(val interface{}) *Builder {
	code := p.code
	i := p.pushInstr(val, len(code.data))
	code.data = append(code.data, i)
	return p
}

// Push instr
func (p Reserved) Push(b *Builder, val interface{}) {
	i := b.pushInstr(val, int(p))
	b.code.data[p] = i
}

// Pop instr
func (p *Builder) Pop(n int) *Builder {
	i := (opPop << bitsOpShift) | uint32(n)
	p.code.data = append(p.code.data, i)
	return p
}

// -----------------------------------------------------------------------------
