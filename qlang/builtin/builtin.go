package builtin

import (
	"fmt"
	"reflect"
	"strings"

	qlang "qlang.io/spec"
	"qlang.io/spec/types"
)

// -----------------------------------------------------------------------------

// Panic panics with v.
//
func Panic(v interface{}) {
	panic(v)
}

// Panicf panics with sprintf(format, args...).
//
func Panicf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

// -----------------------------------------------------------------------------

var (
	zeroVal reflect.Value
)

// Mkmap makes a new map object.
//
func Mkmap(typ interface{}, n ...int) interface{} {

	return reflect.MakeMap(types.Reflect(typ)).Interface()
}

// MapOf makes a map type.
//
func MapOf(key, val interface{}) interface{} {

	return reflect.MapOf(types.Reflect(key), types.Reflect(val))
}

// MapFrom creates a map from args.
//
func MapFrom(args ...interface{}) interface{} {

	n := len(args)
	if (n & 1) != 0 {
		panic("please use `mapFrom(key1, val1, key2, val2, ...)`")
	}
	if n == 0 {
		return make(map[string]interface{})
	}

	switch kindOf2Args(args, 0) {
	case reflect.String:
		switch kindOf2Args(args, 1) {
		case reflect.String:
			ret := make(map[string]string, n>>1)
			for i := 0; i < n; i += 2 {
				ret[args[i].(string)] = args[i+1].(string)
			}
			return ret
		case reflect.Int:
			ret := make(map[string]int, n>>1)
			for i := 0; i < n; i += 2 {
				ret[args[i].(string)] = asInt(args[i+1])
			}
			return ret
		case reflect.Float64:
			ret := make(map[string]float64, n>>1)
			for i := 0; i < n; i += 2 {
				ret[args[i].(string)] = asFloat(args[i+1])
			}
			return ret
		default:
			ret := make(map[string]interface{}, n>>1)
			for i := 0; i < n; i += 2 {
				if t := args[i+1]; t != qlang.Undefined {
					ret[args[i].(string)] = t
				}
			}
			return ret
		}
	case reflect.Int:
		switch kindOf2Args(args, 1) {
		case reflect.String:
			ret := make(map[int]string, n>>1)
			for i := 0; i < n; i += 2 {
				ret[asInt(args[i])] = args[i+1].(string)
			}
			return ret
		case reflect.Int:
			ret := make(map[int]int, n>>1)
			for i := 0; i < n; i += 2 {
				ret[asInt(args[i])] = asInt(args[i+1])
			}
			return ret
		case reflect.Float64:
			ret := make(map[int]float64, n>>1)
			for i := 0; i < n; i += 2 {
				ret[asInt(args[i])] = asFloat(args[i+1])
			}
			return ret
		default:
			ret := make(map[int]interface{}, n>>1)
			for i := 0; i < n; i += 2 {
				if t := args[i+1]; t != qlang.Undefined {
					ret[asInt(args[i])] = t
				}
			}
			return ret
		}
	default:
		panic("mapFrom: key type only support `string`, `int` now")
	}
}

// Delete deletes a key from map object.
//
func Delete(m interface{}, key interface{}) {

	reflect.ValueOf(m).SetMapIndex(reflect.ValueOf(key), zeroVal)
}

// Set sets (index, value) pairs to an object. object can be a slice, an array, a map or a user-defined class.
//
func Set(m interface{}, args ...interface{}) {

	n := len(args)
	if (n & 1) != 0 {
		panic("call with invalid argument count: please use `set(obj, member1, val1, ...)")
	}

	o := reflect.ValueOf(m)
	switch o.Kind() {
	case reflect.Slice, reflect.Array:
		telem := reflect.TypeOf(m).Elem()
		for i := 0; i < n; i += 2 {
			val := autoConvert(telem, args[i+1])
			o.Index(args[i].(int)).Set(val)
		}
	case reflect.Map:
		setMapMember(o, args...)
	default:
		setMember(m, args...)
	}
}

// SetIndex sets a (index, value) pair to an object. object can be a slice, an array, a map or a user-defined class.
//
func SetIndex(m, key, v interface{}) {

	o := reflect.ValueOf(m)
	switch o.Kind() {
	case reflect.Map:
		var val reflect.Value
		if v == qlang.Undefined {
			val = zeroVal
		} else {
			val = autoConvert(o.Type().Elem(), v)
		}
		o.SetMapIndex(reflect.ValueOf(key), val)
	case reflect.Slice, reflect.Array:
		if idx, ok := key.(int); ok {
			o.Index(idx).Set(reflect.ValueOf(v))
		} else {
			panic("slice index isn't an integer value")
		}
	default:
		setMember(m, key, v)
	}
}

type varSetter interface {
	SetVar(name string, val interface{})
}

func setMember(m interface{}, args ...interface{}) {

	if v, ok := m.(varSetter); ok {
		for i := 0; i < len(args); i += 2 {
			v.SetVar(args[i].(string), args[i+1])
		}
		return
	}

	o := reflect.ValueOf(m)
	if o.Kind() == reflect.Ptr {
		o = o.Elem()
		if o.Kind() == reflect.Struct {
			setStructMember(o, args...)
			return
		}
	}
	panic(fmt.Sprintf("type `%v` doesn't support `set` operator", reflect.TypeOf(m)))
}

func setStructMember(o reflect.Value, args ...interface{}) {

	for i := 0; i < len(args); i += 2 {
		key := args[i].(string)
		field := o.FieldByName(strings.Title(key))
		if !field.IsValid() {
			panic(fmt.Sprintf("struct `%v` doesn't has member `%v`", o.Type(), key))
		}
		field.Set(reflect.ValueOf(args[i+1]))
	}
}

func setMapMember(o reflect.Value, args ...interface{}) {

	var val reflect.Value
	telem := o.Type().Elem()
	for i := 0; i < len(args); i += 2 {
		key := reflect.ValueOf(args[i])
		t := args[i+1]
		if t == qlang.Undefined {
			val = zeroVal
		} else {
			val = autoConvert(telem, t)
		}
		o.SetMapIndex(key, val)
	}
}

// Get gets a value from an object. object can be a slice, an array, a map or a user-defined class.
//
func Get(m interface{}, key interface{}) interface{} {

	o := reflect.ValueOf(m)
	switch o.Kind() {
	case reflect.Map:
		v := o.MapIndex(reflect.ValueOf(key))
		if v.IsValid() {
			return v.Interface()
		}
		return qlang.Undefined
	case reflect.Slice, reflect.String, reflect.Array:
		return o.Index(key.(int)).Interface()
	case reflect.Int: // undefined?
		return qlang.Undefined
	default:
		return qlang.GetEx(m, key)
	}
}

// GetVar returns a member variable of an object. object can be a slice, an array, a map or a user-defined class.
//
func GetVar(m interface{}, key interface{}) interface{} {

	return &qlang.DataIndex{Data: m, Index: key}
}

// Len returns length of a collection object. object can be a slice, an array, a map, a string or a chan.
//
func Len(a interface{}) int {

	if a == nil {
		return 0
	}
	if ch, ok := a.(*qlang.Chan); ok {
		return ch.Data.Len()
	}
	return reflect.ValueOf(a).Len()
}

// Cap returns capacity of a collection object. object can be a slice, an array or a chan.
//
func Cap(a interface{}) int {

	if a == nil {
		return 0
	}
	if ch, ok := a.(*qlang.Chan); ok {
		return ch.Data.Cap()
	}
	return reflect.ValueOf(a).Cap()
}

// SubSlice returns a[i:j]. if i == nil it returns a[:j]. if j == nil it returns a[i:].
//
func SubSlice(a, i, j interface{}) interface{} {

	var va = reflect.ValueOf(a)
	var i1, j1 int
	if i != nil {
		i1 = asInt(i)
	}
	if j != nil {
		j1 = asInt(j)
	} else {
		j1 = va.Len()
	}
	return va.Slice(i1, j1).Interface()
}

// Copy does copy(a, b).
//
func Copy(a, b interface{}) int {

	return reflect.Copy(reflect.ValueOf(a), reflect.ValueOf(b))
}

// Append does append(a, vals...)
//
func Append(a interface{}, vals ...interface{}) interface{} {

	switch arr := a.(type) {
	case []int:
		return appendInts(arr, vals...)
	case []interface{}:
		return append(arr, vals...)
	case []string:
		return appendStrings(arr, vals...)
	case []byte:
		return appendBytes(arr, vals...)
	case []float64:
		return appendFloats(arr, vals...)
	}

	va := reflect.ValueOf(a)
	telem := va.Type().Elem()
	x := make([]reflect.Value, len(vals))
	for i, v := range vals {
		x[i] = autoConvert(telem, v)
	}
	return reflect.Append(va, x...).Interface()
}

func autoConvert(telem reflect.Type, v interface{}) reflect.Value {

	if v == nil {
		return reflect.Zero(telem)
	}

	val := reflect.ValueOf(v)
	if telem != reflect.TypeOf(v) {
		val = qlang.AutoConvert(val, telem)
	}
	return val
}

func appendFloats(a []float64, vals ...interface{}) interface{} {

	for _, v := range vals {
		switch val := v.(type) {
		case float64:
			a = append(a, val)
		case int:
			a = append(a, float64(val))
		case float32:
			a = append(a, float64(val))
		default:
			panic("unsupported: []float64 append " + reflect.TypeOf(v).String())
		}
	}
	return a
}

func appendInts(a []int, vals ...interface{}) interface{} {

	for _, v := range vals {
		switch val := v.(type) {
		case int:
			a = append(a, val)
		default:
			panic("unsupported: []int append " + reflect.TypeOf(v).String())
		}
	}
	return a
}

func appendBytes(a []byte, vals ...interface{}) interface{} {

	for _, v := range vals {
		switch val := v.(type) {
		case byte:
			a = append(a, val)
		case int:
			a = append(a, byte(val))
		default:
			panic("unsupported: []byte append " + reflect.TypeOf(v).String())
		}
	}
	return a
}

func appendStrings(a []string, vals ...interface{}) interface{} {

	for _, v := range vals {
		switch val := v.(type) {
		case string:
			a = append(a, val)
		default:
			panic("unsupported: []string append " + reflect.TypeOf(v).String())
		}
	}
	return a
}

// Mkslice returns a new slice.
//
func Mkslice(typ interface{}, args ...interface{}) interface{} {

	n, cap := 0, 0
	if len(args) == 1 {
		if v, ok := args[0].(int); ok {
			n, cap = v, v
		} else {
			panic("second param type of func `slice` must be `int`")
		}
	} else if len(args) > 1 {
		if v, ok := args[0].(int); ok {
			n = v
		} else {
			panic("2nd param type of func `slice` must be `int`")
		}
		if v, ok := args[1].(int); ok {
			cap = v
		} else {
			panic("3rd param type of func `slice` must be `int`")
		}
	}
	typSlice := reflect.SliceOf(types.Reflect(typ))
	return reflect.MakeSlice(typSlice, n, cap).Interface()
}

// SliceFrom creates a slice from [a1, a2, ...].
//
func SliceFrom(args ...interface{}) interface{} {

	n := len(args)
	if n == 0 {
		return []interface{}(nil)
	}

	switch kindOfArgs(args) {
	case reflect.Int:
		return appendInts(make([]int, 0, n), args...)
	case reflect.Float64:
		return appendFloats(make([]float64, 0, n), args...)
	case reflect.String:
		return appendStrings(make([]string, 0, n), args...)
	case reflect.Uint8:
		return appendBytes(make([]byte, 0, n), args...)
	default:
		return append(make([]interface{}, 0, n), args...)
	}
}

// SliceFromTy creates a slice from `[]T{a1, a2, ...}`.
//
func SliceFromTy(args ...interface{}) interface{} {

	got, ok := args[0].(qlang.GoTyper)
	if !ok {
		panic(fmt.Sprintf("`%v` is not a qlang type", args[0]))
	}
	t := got.GoType()
	n := len(args)
	ret := reflect.MakeSlice(reflect.SliceOf(t), 0, n-1).Interface()
	return Append(ret, args[1:]...)
}

// SliceOf makes a slice type.
//
func SliceOf(typ interface{}) interface{} {

	return reflect.SliceOf(types.Reflect(typ))
}

// StructInit creates a struct object from `structInit(structType, member1, val1, ...)`.
//
func StructInit(args ...interface{}) interface{} {

	if (len(args) & 1) != 1 {
		panic("call with invalid argument count: please use `structInit(structType, member1, val1, ...)")
	}

	got, ok := args[0].(qlang.GoTyper)
	if !ok {
		panic(fmt.Sprintf("`%v` is not a qlang type", args[0]))
	}
	t := got.GoType()
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("`%v` is not a struct type", args[0]))
	}
	ret := reflect.New(t)
	setStructMember(ret.Elem(), args[1:]...)
	return ret.Interface()
}

// MapInit creates a map object from `mapInit(mapType, member1, val1, ...)`.
//
func MapInit(args ...interface{}) interface{} {

	if (len(args) & 1) != 1 {
		panic("call with invalid argument count: please use `mapInit(mapType, member1, val1, ...)")
	}

	got, ok := args[0].(qlang.GoTyper)
	if !ok {
		panic(fmt.Sprintf("`%v` is not a qlang type", args[0]))
	}
	t := got.GoType()
	if t.Kind() != reflect.Map {
		panic(fmt.Sprintf("`%v` is not a map type", args[0]))
	}
	ret := reflect.MakeMap(t)
	setMapMember(ret, args[1:]...)
	return ret.Interface()
}

// -----------------------------------------------------------------------------
