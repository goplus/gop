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

package bytecode

import (
	"reflect"

	"github.com/goplus/reflectx"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// A StructField describes a single field in a struct.
type StructField = reflect.StructField

// A StructInfo represents a struct information.
type StructInfo struct {
	Fields []StructField
	typ    reflect.Type
}

// Struct creates a new struct.
func Struct(fields []StructField) *StructInfo {
	return &StructInfo{
		Fields: fields,
		typ:    reflect.StructOf(fields),
	}
}

// Type returns the struct type.
func (p *StructInfo) Type() reflect.Type {
	return p.typ
}

// -----------------------------------------------------------------------------

func (p *varManager) makeVarsContext(ctx *Context) (vars []reflect.Value) {
	size := len(p.vlist)
	if size > 0 {
		vars = make([]reflect.Value, size)
		for i := 0; i < size; i++ {
			vars[i] = reflect.New(p.vlist[i].typ)
		}
	}
	return
}

func (ctx *varScope) addrVar(idx uint32) interface{} {
	return ctx.vars[int(idx)].Interface()
}

func (ctx *varScope) getVar(idx uint32) interface{} {
	return ctx.vars[int(idx)].Elem().Interface()
}

func (ctx *varScope) setVar(idx uint32, v interface{}) {
	x := ctx.vars[int(idx)].Elem()
	setValue(x, v)
}

func (ctx *varScope) setValue(idx uint32, v reflect.Value) {
	ctx.vars[int(idx)].Elem().Set(v)
}

func setValue(x reflect.Value, v interface{}) {
	if v != nil {
		x.Set(reflect.ValueOf(v))
		return
	}
	x.Set(reflect.Zero(x.Type()))
}

func getValueOf(v interface{}, t reflect.Type) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(t)
}

func getKeyOf(v interface{}, t reflect.Type) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(t.Key())
}

func getElementOf(v interface{}, t reflect.Type) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(t.Elem())
}

func getArgOf(v interface{}, t reflect.Type, i int) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(t.In(i))
}

func getRetOf(v interface{}, fun *FuncInfo, i int) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(fun.vlist[i].typ)
}

// -----------------------------------------------------------------------------

// Struct instr
func (p *Builder) Struct(typ reflect.Type, arity int) *Builder {
	if arity < 0 {
		log.Panicln("Struct failed: can't be variadic.")
	} else if arity >= bitsFuncvArityMax {
		p.Push(arity - bitsFuncvArityMax)
		arity = bitsFuncvArityMax
	}
	i := (opStruct << bitsOpShift) | (uint32(arity) << bitsOpCallFuncvShift) | p.newType(typ)
	p.code.data = append(p.code.data, i)
	return p
}

func execStruct(i Instr, p *Context) {
	typ := getType(i&bitsOpCallFuncvOperand, p)
	arity := int((i >> bitsOpCallFuncvShift) & bitsFuncvArityOperand)
	if arity == bitsFuncvArityMax {
		arity = p.Pop().(int) + bitsFuncvArityMax
	}
	makeStruct(typ, arity, p)
}

func makeStruct(typStruct reflect.Type, arity int, p *Context) {
	n := arity << 1
	args := p.GetArgs(n)

	var ptr bool
	for typStruct.Kind() == reflect.Ptr {
		ptr = true
		typStruct = typStruct.Elem()
	}
	v := reflect.New(typStruct).Elem()
	for i := 0; i < n; i += 2 {
		index := args[i].(int)
		reflectx.Field(v, index).Set(reflect.ValueOf(args[i+1]))
	}
	if ptr {
		p.Ret(n, v.Addr().Interface())
	} else {
		p.Ret(n, v.Interface())
	}
}
