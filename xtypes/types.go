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

// Package xtypes provides `go/types` extended utilities. for example,
// converting `types.Type` into `reflect.Type`.
package xtypes

import (
	"errors"
	"fmt"
	"go/types"
	"reflect"
	"syscall"
	"unsafe"
)

// GoTypeFinder interface
type GoTypeFinder interface {
	FindGoType(pkgPath string, namedType string) (reflect.Type, bool)
}

// -----------------------------------------------------------------------------

var basicTypes = [...]reflect.Type{
	types.Bool:          reflect.TypeOf(false),
	types.Int:           reflect.TypeOf(0),
	types.Int8:          reflect.TypeOf(int8(0)),
	types.Int16:         reflect.TypeOf(int16(0)),
	types.Int32:         reflect.TypeOf(int32(0)),
	types.Int64:         reflect.TypeOf(int64(0)),
	types.Uint:          reflect.TypeOf(uint(0)),
	types.Uint8:         reflect.TypeOf(uint8(0)),
	types.Uint16:        reflect.TypeOf(uint16(0)),
	types.Uint32:        reflect.TypeOf(uint32(0)),
	types.Uint64:        reflect.TypeOf(uint64(0)),
	types.Uintptr:       reflect.TypeOf(uintptr(0)),
	types.Float32:       reflect.TypeOf(float32(0)),
	types.Float64:       reflect.TypeOf(float64(0)),
	types.Complex64:     reflect.TypeOf(complex64(0)),
	types.Complex128:    reflect.TypeOf(complex128(0)),
	types.String:        reflect.TypeOf(""),
	types.UnsafePointer: reflect.TypeOf(unsafe.Pointer(nil)),
}

var (
	tyEmptyInterface = reflect.TypeOf((*interface{})(nil)).Elem()
)

var (
	// ErrUntyped error
	ErrUntyped = errors.New("untyped type")
	// ErrUnknownArrayLen error
	ErrUnknownArrayLen = errors.New("unknown array length")
)

// ToType converts a types.Type to reflect.Type.
func ToType(typ types.Type, ctx GoTypeFinder) (reflect.Type, error) {
	switch t := typ.(type) {
	case *types.Basic:
		if kind := t.Kind(); kind >= types.Bool && kind <= types.UnsafePointer {
			return basicTypes[kind], nil
		}
		return nil, ErrUntyped
	case *types.Pointer:
		elem, err := ToType(t.Elem(), ctx)
		if err != nil {
			return nil, fmt.Errorf("unknown pointer elem type - %w", err)
		}
		return reflect.PtrTo(elem), nil
	case *types.Slice:
		elem, err := ToType(t.Elem(), ctx)
		if err != nil {
			return nil, fmt.Errorf("unknown slice elem type - %w", err)
		}
		return reflect.SliceOf(elem), nil
	case *types.Array:
		elem, err := ToType(t.Elem(), ctx)
		if err != nil {
			return nil, fmt.Errorf("unknown array elem type - %w", err)
		}
		n := t.Len()
		if n < 0 {
			return nil, ErrUnknownArrayLen
		}
		return reflect.ArrayOf(int(n), elem), nil
	case *types.Map:
		key, err := ToType(t.Key(), ctx)
		if err != nil {
			return nil, fmt.Errorf("unknown map key type - %w", err)
		}
		elem, err := ToType(t.Elem(), ctx)
		if err != nil {
			return nil, fmt.Errorf("unknown map elem type - %w", err)
		}
		return reflect.MapOf(key, elem), nil
	case *types.Chan:
		elem, err := ToType(t.Elem(), ctx)
		if err != nil {
			return nil, fmt.Errorf("unknown chan elem type - %w", err)
		}
		return reflect.ChanOf(toChanDir(t.Dir()), elem), nil
	case *types.Struct:
		return toStructType(t, ctx)
	case *types.Named:
		return toNamedType(t, ctx)
	case *types.Interface:
		return tyEmptyInterface, nil
	}
	return nil, syscall.EINVAL
}

func toChanDir(d types.ChanDir) reflect.ChanDir {
	switch d {
	case types.SendRecv:
		return reflect.BothDir
	case types.SendOnly:
		return reflect.SendDir
	case types.RecvOnly:
		return reflect.RecvDir
	}
	return 0
}

// toStructType converts a types.Struct to reflect.Type.
func toStructType(t *types.Struct, ctx GoTypeFinder) (typ reflect.Type, err error) {
	n := t.NumFields()
	flds := make([]reflect.StructField, n)
	for i := 0; i < n; i++ {
		v := t.Field(i)
		if flds[i], err = toStructField(v, t.Tag(i), ctx); err != nil {
			return
		}
	}
	return reflect.StructOf(flds), nil
}

func toStructField(v *types.Var, tag string, ctx GoTypeFinder) (fld reflect.StructField, err error) {
	name := v.Name()
	typ, err := ToType(v.Type(), ctx)
	if err != nil {
		err = fmt.Errorf("unknown struct field `%s` type - %w", name, err)
		return
	}
	fld = reflect.StructField{
		Name:      name,
		Type:      typ,
		Tag:       reflect.StructTag(tag),
		Anonymous: v.Anonymous(),
	}
	if !fld.Anonymous && !v.Exported() {
		fld.PkgPath = v.Pkg().Path()
	}
	return
}

func toNamedType(t *types.Named, ctx GoTypeFinder) (reflect.Type, error) {
	name := t.Obj()
	if t, ok := ctx.FindGoType(name.Pkg().Path(), name.Name()); ok {
		return t, nil
	}
	// So, this type is defined in Go+, so in fact it is same as the underlying type.
	typ, err := ToType(t.Underlying(), ctx)
	if err != nil {
		return nil, fmt.Errorf("named type `%s` - %w", name.Name(), err)
	}
	return typ, nil
}

// -----------------------------------------------------------------------------
