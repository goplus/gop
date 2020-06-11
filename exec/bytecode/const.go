/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

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

package bytecode

import (
	"reflect"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

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

func execPushInt(i Instr, stk *Context) {
	v := int32(i) << bitsOpInt >> bitsOpInt
	kind := reflect.Int + reflect.Kind((i>>bitsOpIntShift)&7)
	pushInt32(stk, kind, v)
}

func execPushUint(i Instr, stk *Context) {
	v := i & bitsOpIntOperand
	kind := reflect.Uint + reflect.Kind((i>>bitsOpIntShift)&7)
	pushUint32(stk, kind, v)
}

func execPushConstR(i Instr, stk *Context) {
	v := stk.code.valConsts[i&bitsOperand]
	stk.Push(v)
}

func execPop(i Instr, stk *Context) {
	n := len(stk.data) - int(i&bitsOperand)
	stk.data = stk.data[:n]
}

// -----------------------------------------------------------------------------

// Push instr
func (p *Builder) pushInstr(val interface{}) (i Instr) {
	if val == nil {
		return iPushNil
	}
	v := reflect.ValueOf(val)
	kind := v.Kind()
	if kind >= reflect.Int && kind <= reflect.Int64 {
		iv := v.Int()
		ivStore := int64(int32(iv) << bitsOpInt >> bitsOpInt)
		if iv == ivStore {
			i = (opPushInt << bitsOpShift) | (uint32(kind-reflect.Int) << bitsOpIntShift) | (uint32(iv) & bitsOpIntOperand)
			return
		}
	} else if kind >= reflect.Uint && kind <= reflect.Uintptr {
		iv := v.Uint()
		if iv == (iv & bitsOpIntOperand) {
			i = (opPushUint << bitsOpShift) | (uint32(kind-reflect.Uint) << bitsOpIntShift) | (uint32(iv) & bitsOpIntOperand)
			return
		}
	} else if kind == reflect.Bool {
		if val.(bool) {
			return iPushTrue
		}
		return iPushFalse
	} else if kind != reflect.String && !(kind >= reflect.Float32 && kind <= reflect.Complex128) {
		if !(kind == reflect.Ptr && v.Type().Elem().PkgPath() == "math/big") {
			log.Panicln("Push failed: unsupported type:", reflect.TypeOf(val), "-", val)
		}
	}
	code := p.code
	i = (opPushConstR << bitsOpShift) | uint32(len(code.valConsts))
	code.valConsts = append(code.valConsts, val)
	return
}

// Push instr
func (p *Builder) Push(val interface{}) *Builder {
	p.code.data = append(p.code.data, p.pushInstr(val))
	return p
}

// ReservedAsPush sets Reserved as Push(v)
func (p *Builder) ReservedAsPush(r Reserved, val interface{}) {
	p.code.data[r] = p.pushInstr(val)
}

// Pop instr
func (p *Builder) Pop(n int) *Builder {
	i := (opPop << bitsOpShift) | uint32(n)
	p.code.data = append(p.code.data, i)
	return p
}

// -----------------------------------------------------------------------------
