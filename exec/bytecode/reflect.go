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
	"unsafe"

	"github.com/goplus/gop/ast"
)

type Value struct {
	typ *unsafe.Pointer
	ptr unsafe.Pointer
	flag
}

type flag uintptr

const (
	flagKindWidth        = 5 // there are 27 kinds
	flagKindMask    flag = 1<<flagKindWidth - 1
	flagStickyRO    flag = 1 << 5
	flagEmbedRO     flag = 1 << 6
	flagIndir       flag = 1 << 7
	flagAddr        flag = 1 << 8
	flagMethod      flag = 1 << 9
	flagMethodShift      = 10
	flagRO          flag = flagStickyRO | flagEmbedRO
)

func Field(v reflect.Value, i int) reflect.Value {
	sf := v.Type().Field(i)
	if !ast.IsExported(sf.Name) {
		vf := v.Field(i)
		(*Value)(unsafe.Pointer(&vf)).flag &= ^flagRO
		return vf
	}
	return v.Field(i)
}

func FieldByIndex(v reflect.Value, index []int) reflect.Value {
	sf := v.Type().FieldByIndex(index)
	if !ast.IsExported(sf.Name) {
		vf := v.FieldByIndex(index)
		(*Value)(unsafe.Pointer(&vf)).flag &= ^flagRO
		return vf
	}
	return v.FieldByIndex(index)
}
