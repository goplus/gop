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

// Package bytecode implements a bytecode backend for the Go+ language.
package bytecode

import (
	"reflect"
)

// StoreVal instr
func (p *Builder) StoreVal(val interface{}) *Builder {
	p.code.data = append(p.code.data, opStoreVal<<bitsOpShift|uint32(len(p.code.vals)))
	p.code.vals = append(p.code.vals, val)
	return p
}

func execStoreVal(i Instr, stk *Context) {
	v := stk.code.vals[i&bitsOperand]
	stk.Push(v)
}

func (p *Builder) SetField() *Builder {
	p.code.data = append(p.code.data, opSetField<<bitsOpShift)
	return p
}

func execSetField(i Instr, stk *Context) {
	args := stk.GetArgs(3)
	d := args[0]
	p := args[1]
	field := args[2]
	v := reflect.ValueOf(p)
	f := field.(string)

	if v.Kind() == reflect.Ptr {
		v.Elem().FieldByName(f).Set(reflect.ValueOf(d))
	} else {
		t := reflect.TypeOf(p)
		v2 := reflect.New(t).Elem()
		v2.Set(v)
		v2.FieldByName(f).Set(reflect.ValueOf(d))
		v = v2
	}

	// t := reflect.TypeOf(p)
	// if t.Kind() == reflect.Ptr {
	// 	t = t.Elem()
	// }

	// v2 := reflect.New(t).Elem()

	// for i := 0; i < t.NumField(); i++ {
	// 	field := t.Field(i)
	// 	if field.Name == f {
	// 		v2.Field(i).Set(reflect.ValueOf(d))
	// 	} else {
	// 		v2.Field(i).Set(v.Elem().Field(i))
	// 	}
	// }

	// v = v.FieldByName(f)

	// v.Set(reflect.ValueOf(d))

	stk.Push(v.Interface())
}

func (p *Builder) CallField() *Builder {
	p.code.data = append(p.code.data, opCallField<<bitsOpShift)
	return p
}

func execCallField(i Instr, stk *Context) {
	args := stk.GetArgs(2)
	p := args[0]
	field := args[1]
	v := reflect.ValueOf(p)
	f := field.(string)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	v = reflect.Indirect(v).FieldByName(f)
	stk.Push(v.Interface())
}
