// Package reflect provide Go+ "reflect" package, as "reflect" package in Go.
package reflect

import (
	reflect "reflect"
	unsafe "unsafe"

	gop "github.com/qiniu/goplus/gop"
)

func toSlice0(args []interface{}) []reflect.Value {
	ret := make([]reflect.Value, len(args))
	for i, arg := range args {
		ret[i] = arg.(reflect.Value)
	}
	return ret
}

func execAppend(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := reflect.Append(args[0].(reflect.Value), toSlice0(args[1:])...)
	p.Ret(arity, ret0)
}

func execAppendSlice(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflect.AppendSlice(args[0].(reflect.Value), args[1].(reflect.Value))
	p.Ret(2, ret0)
}

func toType0(v interface{}) reflect.Type {
	if v == nil {
		return nil
	}
	return v.(reflect.Type)
}

func execArrayOf(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflect.ArrayOf(args[0].(int), toType0(args[1]))
	p.Ret(2, ret0)
}

func execChanDirString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.ChanDir).String()
	p.Ret(1, ret0)
}

func execChanOf(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflect.ChanOf(args[0].(reflect.ChanDir), toType0(args[1]))
	p.Ret(2, ret0)
}

func execCopy(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflect.Copy(args[0].(reflect.Value), args[1].(reflect.Value))
	p.Ret(2, ret0)
}

func execDeepEqual(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflect.DeepEqual(args[0], args[1])
	p.Ret(2, ret0)
}

func execFuncOf(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := reflect.FuncOf(args[0].([]reflect.Type), args[1].([]reflect.Type), args[2].(bool))
	p.Ret(3, ret0)
}

func execIndirect(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflect.Indirect(args[0].(reflect.Value))
	p.Ret(1, ret0)
}

func execKindString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Kind).String()
	p.Ret(1, ret0)
}

func execMakeChan(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflect.MakeChan(toType0(args[0]), args[1].(int))
	p.Ret(2, ret0)
}

func execMakeFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflect.MakeFunc(toType0(args[0]), args[1].(func(args []reflect.Value) (results []reflect.Value)))
	p.Ret(2, ret0)
}

func execMakeMap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflect.MakeMap(toType0(args[0]))
	p.Ret(1, ret0)
}

func execMakeMapWithSize(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflect.MakeMapWithSize(toType0(args[0]), args[1].(int))
	p.Ret(2, ret0)
}

func execMakeSlice(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := reflect.MakeSlice(toType0(args[0]), args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execMapIterKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*reflect.MapIter).Key()
	p.Ret(1, ret0)
}

func execMapIterValue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*reflect.MapIter).Value()
	p.Ret(1, ret0)
}

func execMapIterNext(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*reflect.MapIter).Next()
	p.Ret(1, ret0)
}

func execMapOf(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflect.MapOf(toType0(args[0]), toType0(args[1]))
	p.Ret(2, ret0)
}

func execNew(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflect.New(toType0(args[0]))
	p.Ret(1, ret0)
}

func execNewAt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflect.NewAt(toType0(args[0]), args[1].(unsafe.Pointer))
	p.Ret(2, ret0)
}

func execPtrTo(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflect.PtrTo(toType0(args[0]))
	p.Ret(1, ret0)
}

func execSelect(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := reflect.Select(args[0].([]reflect.SelectCase))
	p.Ret(1, ret0, ret1, ret2)
}

func execSliceOf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflect.SliceOf(toType0(args[0]))
	p.Ret(1, ret0)
}

func execStructOf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflect.StructOf(args[0].([]reflect.StructField))
	p.Ret(1, ret0)
}

func execStructTagGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.StructTag).Get(args[1].(string))
	p.Ret(2, ret0)
}

func execStructTagLookup(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(reflect.StructTag).Lookup(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execSwapper(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflect.Swapper(args[0])
	p.Ret(1, ret0)
}

func execTypeOf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflect.TypeOf(args[0])
	p.Ret(1, ret0)
}

func execValueAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Addr()
	p.Ret(1, ret0)
}

func execValueBool(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Bool()
	p.Ret(1, ret0)
}

func execValueBytes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Bytes()
	p.Ret(1, ret0)
}

func execValueCanAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).CanAddr()
	p.Ret(1, ret0)
}

func execValueCanSet(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).CanSet()
	p.Ret(1, ret0)
}

func execValueCall(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).Call(args[1].([]reflect.Value))
	p.Ret(2, ret0)
}

func execValueCallSlice(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).CallSlice(args[1].([]reflect.Value))
	p.Ret(2, ret0)
}

func execValueCap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Cap()
	p.Ret(1, ret0)
}

func execValueClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(reflect.Value).Close()
	p.PopN(1)
}

func execValueComplex(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Complex()
	p.Ret(1, ret0)
}

func execValueElem(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Elem()
	p.Ret(1, ret0)
}

func execValueField(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).Field(args[1].(int))
	p.Ret(2, ret0)
}

func execValueFieldByIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).FieldByIndex(args[1].([]int))
	p.Ret(2, ret0)
}

func execValueFieldByName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).FieldByName(args[1].(string))
	p.Ret(2, ret0)
}

func execValueFieldByNameFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).FieldByNameFunc(args[1].(func(string) bool))
	p.Ret(2, ret0)
}

func execValueFloat(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Float()
	p.Ret(1, ret0)
}

func execValueIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).Index(args[1].(int))
	p.Ret(2, ret0)
}

func execValueInt(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Int()
	p.Ret(1, ret0)
}

func execValueCanInterface(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).CanInterface()
	p.Ret(1, ret0)
}

func execValueInterface(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Interface()
	p.Ret(1, ret0)
}

func execValueInterfaceData(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).InterfaceData()
	p.Ret(1, ret0)
}

func execValueIsNil(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).IsNil()
	p.Ret(1, ret0)
}

func execValueIsValid(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).IsValid()
	p.Ret(1, ret0)
}

func execValueIsZero(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).IsZero()
	p.Ret(1, ret0)
}

func execValueKind(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Kind()
	p.Ret(1, ret0)
}

func execValueLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Len()
	p.Ret(1, ret0)
}

func execValueMapIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).MapIndex(args[1].(reflect.Value))
	p.Ret(2, ret0)
}

func execValueMapKeys(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).MapKeys()
	p.Ret(1, ret0)
}

func execValueMapRange(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).MapRange()
	p.Ret(1, ret0)
}

func execValueMethod(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).Method(args[1].(int))
	p.Ret(2, ret0)
}

func execValueNumMethod(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).NumMethod()
	p.Ret(1, ret0)
}

func execValueMethodByName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).MethodByName(args[1].(string))
	p.Ret(2, ret0)
}

func execValueNumField(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).NumField()
	p.Ret(1, ret0)
}

func execValueOverflowComplex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).OverflowComplex(args[1].(complex128))
	p.Ret(2, ret0)
}

func execValueOverflowFloat(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).OverflowFloat(args[1].(float64))
	p.Ret(2, ret0)
}

func execValueOverflowInt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).OverflowInt(args[1].(int64))
	p.Ret(2, ret0)
}

func execValueOverflowUint(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).OverflowUint(args[1].(uint64))
	p.Ret(2, ret0)
}

func execValuePointer(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Pointer()
	p.Ret(1, ret0)
}

func execValueRecv(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(reflect.Value).Recv()
	p.Ret(1, ret0, ret1)
}

func execValueSend(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).Send(args[1].(reflect.Value))
	p.PopN(2)
}

func execValueSet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).Set(args[1].(reflect.Value))
	p.PopN(2)
}

func execValueSetBool(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetBool(args[1].(bool))
	p.PopN(2)
}

func execValueSetBytes(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetBytes(args[1].([]byte))
	p.PopN(2)
}

func execValueSetComplex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetComplex(args[1].(complex128))
	p.PopN(2)
}

func execValueSetFloat(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetFloat(args[1].(float64))
	p.PopN(2)
}

func execValueSetInt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetInt(args[1].(int64))
	p.PopN(2)
}

func execValueSetLen(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetLen(args[1].(int))
	p.PopN(2)
}

func execValueSetCap(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetCap(args[1].(int))
	p.PopN(2)
}

func execValueSetMapIndex(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(reflect.Value).SetMapIndex(args[1].(reflect.Value), args[2].(reflect.Value))
	p.PopN(3)
}

func execValueSetUint(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetUint(args[1].(uint64))
	p.PopN(2)
}

func execValueSetPointer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetPointer(args[1].(unsafe.Pointer))
	p.PopN(2)
}

func execValueSetString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetString(args[1].(string))
	p.PopN(2)
}

func execValueSlice(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(reflect.Value).Slice(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execValueSlice3(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(reflect.Value).Slice3(args[1].(int), args[2].(int), args[3].(int))
	p.Ret(4, ret0)
}

func execValueString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).String()
	p.Ret(1, ret0)
}

func execValueTryRecv(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(reflect.Value).TryRecv()
	p.Ret(1, ret0, ret1)
}

func execValueTrySend(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).TrySend(args[1].(reflect.Value))
	p.Ret(2, ret0)
}

func execValueType(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Type()
	p.Ret(1, ret0)
}

func execValueUint(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Uint()
	p.Ret(1, ret0)
}

func execValueUnsafeAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).UnsafeAddr()
	p.Ret(1, ret0)
}

func execValueConvert(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).Convert(toType0(args[1]))
	p.Ret(2, ret0)
}

func execValueErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*reflect.ValueError).Error()
	p.Ret(1, ret0)
}

func execValueOf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflect.ValueOf(args[0])
	p.Ret(1, ret0)
}

func execZero(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflect.Zero(toType0(args[0]))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("reflect")

func init() {
	I.RegisterFuncs(
		I.Func("AppendSlice", reflect.AppendSlice, execAppendSlice),
		I.Func("ArrayOf", reflect.ArrayOf, execArrayOf),
		I.Func("(ChanDir).String", (reflect.ChanDir).String, execChanDirString),
		I.Func("ChanOf", reflect.ChanOf, execChanOf),
		I.Func("Copy", reflect.Copy, execCopy),
		I.Func("DeepEqual", reflect.DeepEqual, execDeepEqual),
		I.Func("FuncOf", reflect.FuncOf, execFuncOf),
		I.Func("Indirect", reflect.Indirect, execIndirect),
		I.Func("(Kind).String", (reflect.Kind).String, execKindString),
		I.Func("MakeChan", reflect.MakeChan, execMakeChan),
		I.Func("MakeFunc", reflect.MakeFunc, execMakeFunc),
		I.Func("MakeMap", reflect.MakeMap, execMakeMap),
		I.Func("MakeMapWithSize", reflect.MakeMapWithSize, execMakeMapWithSize),
		I.Func("MakeSlice", reflect.MakeSlice, execMakeSlice),
		I.Func("(*MapIter).Key", (*reflect.MapIter).Key, execMapIterKey),
		I.Func("(*MapIter).Value", (*reflect.MapIter).Value, execMapIterValue),
		I.Func("(*MapIter).Next", (*reflect.MapIter).Next, execMapIterNext),
		I.Func("MapOf", reflect.MapOf, execMapOf),
		I.Func("New", reflect.New, execNew),
		I.Func("NewAt", reflect.NewAt, execNewAt),
		I.Func("PtrTo", reflect.PtrTo, execPtrTo),
		I.Func("Select", reflect.Select, execSelect),
		I.Func("SliceOf", reflect.SliceOf, execSliceOf),
		I.Func("StructOf", reflect.StructOf, execStructOf),
		I.Func("(StructTag).Get", (reflect.StructTag).Get, execStructTagGet),
		I.Func("(StructTag).Lookup", (reflect.StructTag).Lookup, execStructTagLookup),
		I.Func("Swapper", reflect.Swapper, execSwapper),
		I.Func("TypeOf", reflect.TypeOf, execTypeOf),
		I.Func("(Value).Addr", (reflect.Value).Addr, execValueAddr),
		I.Func("(Value).Bool", (reflect.Value).Bool, execValueBool),
		I.Func("(Value).Bytes", (reflect.Value).Bytes, execValueBytes),
		I.Func("(Value).CanAddr", (reflect.Value).CanAddr, execValueCanAddr),
		I.Func("(Value).CanSet", (reflect.Value).CanSet, execValueCanSet),
		I.Func("(Value).Call", (reflect.Value).Call, execValueCall),
		I.Func("(Value).CallSlice", (reflect.Value).CallSlice, execValueCallSlice),
		I.Func("(Value).Cap", (reflect.Value).Cap, execValueCap),
		I.Func("(Value).Close", (reflect.Value).Close, execValueClose),
		I.Func("(Value).Complex", (reflect.Value).Complex, execValueComplex),
		I.Func("(Value).Elem", (reflect.Value).Elem, execValueElem),
		I.Func("(Value).Field", (reflect.Value).Field, execValueField),
		I.Func("(Value).FieldByIndex", (reflect.Value).FieldByIndex, execValueFieldByIndex),
		I.Func("(Value).FieldByName", (reflect.Value).FieldByName, execValueFieldByName),
		I.Func("(Value).FieldByNameFunc", (reflect.Value).FieldByNameFunc, execValueFieldByNameFunc),
		I.Func("(Value).Float", (reflect.Value).Float, execValueFloat),
		I.Func("(Value).Index", (reflect.Value).Index, execValueIndex),
		I.Func("(Value).Int", (reflect.Value).Int, execValueInt),
		I.Func("(Value).CanInterface", (reflect.Value).CanInterface, execValueCanInterface),
		I.Func("(Value).Interface", (reflect.Value).Interface, execValueInterface),
		I.Func("(Value).InterfaceData", (reflect.Value).InterfaceData, execValueInterfaceData),
		I.Func("(Value).IsNil", (reflect.Value).IsNil, execValueIsNil),
		I.Func("(Value).IsValid", (reflect.Value).IsValid, execValueIsValid),
		I.Func("(Value).IsZero", (reflect.Value).IsZero, execValueIsZero),
		I.Func("(Value).Kind", (reflect.Value).Kind, execValueKind),
		I.Func("(Value).Len", (reflect.Value).Len, execValueLen),
		I.Func("(Value).MapIndex", (reflect.Value).MapIndex, execValueMapIndex),
		I.Func("(Value).MapKeys", (reflect.Value).MapKeys, execValueMapKeys),
		I.Func("(Value).MapRange", (reflect.Value).MapRange, execValueMapRange),
		I.Func("(Value).Method", (reflect.Value).Method, execValueMethod),
		I.Func("(Value).NumMethod", (reflect.Value).NumMethod, execValueNumMethod),
		I.Func("(Value).MethodByName", (reflect.Value).MethodByName, execValueMethodByName),
		I.Func("(Value).NumField", (reflect.Value).NumField, execValueNumField),
		I.Func("(Value).OverflowComplex", (reflect.Value).OverflowComplex, execValueOverflowComplex),
		I.Func("(Value).OverflowFloat", (reflect.Value).OverflowFloat, execValueOverflowFloat),
		I.Func("(Value).OverflowInt", (reflect.Value).OverflowInt, execValueOverflowInt),
		I.Func("(Value).OverflowUint", (reflect.Value).OverflowUint, execValueOverflowUint),
		I.Func("(Value).Pointer", (reflect.Value).Pointer, execValuePointer),
		I.Func("(Value).Recv", (reflect.Value).Recv, execValueRecv),
		I.Func("(Value).Send", (reflect.Value).Send, execValueSend),
		I.Func("(Value).Set", (reflect.Value).Set, execValueSet),
		I.Func("(Value).SetBool", (reflect.Value).SetBool, execValueSetBool),
		I.Func("(Value).SetBytes", (reflect.Value).SetBytes, execValueSetBytes),
		I.Func("(Value).SetComplex", (reflect.Value).SetComplex, execValueSetComplex),
		I.Func("(Value).SetFloat", (reflect.Value).SetFloat, execValueSetFloat),
		I.Func("(Value).SetInt", (reflect.Value).SetInt, execValueSetInt),
		I.Func("(Value).SetLen", (reflect.Value).SetLen, execValueSetLen),
		I.Func("(Value).SetCap", (reflect.Value).SetCap, execValueSetCap),
		I.Func("(Value).SetMapIndex", (reflect.Value).SetMapIndex, execValueSetMapIndex),
		I.Func("(Value).SetUint", (reflect.Value).SetUint, execValueSetUint),
		I.Func("(Value).SetPointer", (reflect.Value).SetPointer, execValueSetPointer),
		I.Func("(Value).SetString", (reflect.Value).SetString, execValueSetString),
		I.Func("(Value).Slice", (reflect.Value).Slice, execValueSlice),
		I.Func("(Value).Slice3", (reflect.Value).Slice3, execValueSlice3),
		I.Func("(Value).String", (reflect.Value).String, execValueString),
		I.Func("(Value).TryRecv", (reflect.Value).TryRecv, execValueTryRecv),
		I.Func("(Value).TrySend", (reflect.Value).TrySend, execValueTrySend),
		I.Func("(Value).Type", (reflect.Value).Type, execValueType),
		I.Func("(Value).Uint", (reflect.Value).Uint, execValueUint),
		I.Func("(Value).UnsafeAddr", (reflect.Value).UnsafeAddr, execValueUnsafeAddr),
		I.Func("(Value).Convert", (reflect.Value).Convert, execValueConvert),
		I.Func("(*ValueError).Error", (*reflect.ValueError).Error, execValueErrorError),
		I.Func("ValueOf", reflect.ValueOf, execValueOf),
		I.Func("Zero", reflect.Zero, execZero),
	)
	I.RegisterFuncvs(
		I.Funcv("Append", reflect.Append, execAppend),
	)
}
