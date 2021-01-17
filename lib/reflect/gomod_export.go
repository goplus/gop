// Package reflect provide Go+ "reflect" package, as "reflect" package in Go.
package reflect

import (
	reflect "reflect"
	unsafe "unsafe"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
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

func execmChanDirString(_ int, p *gop.Context) {
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

func execmKindString(_ int, p *gop.Context) {
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

func execmMapIterKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*reflect.MapIter).Key()
	p.Ret(1, ret0)
}

func execmMapIterValue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*reflect.MapIter).Value()
	p.Ret(1, ret0)
}

func execmMapIterNext(_ int, p *gop.Context) {
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

func execmStructTagGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.StructTag).Get(args[1].(string))
	p.Ret(2, ret0)
}

func execmStructTagLookup(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(reflect.StructTag).Lookup(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execSwapper(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflect.Swapper(args[0])
	p.Ret(1, ret0)
}

func execiTypeAlign(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).Align()
	p.Ret(1, ret0)
}

func execiTypeAssignableTo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Type).AssignableTo(toType0(args[1]))
	p.Ret(2, ret0)
}

func execiTypeBits(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).Bits()
	p.Ret(1, ret0)
}

func execiTypeChanDir(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).ChanDir()
	p.Ret(1, ret0)
}

func execiTypeComparable(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).Comparable()
	p.Ret(1, ret0)
}

func execiTypeConvertibleTo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Type).ConvertibleTo(toType0(args[1]))
	p.Ret(2, ret0)
}

func execiTypeElem(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).Elem()
	p.Ret(1, ret0)
}

func execiTypeField(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Type).Field(args[1].(int))
	p.Ret(2, ret0)
}

func execiTypeFieldAlign(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).FieldAlign()
	p.Ret(1, ret0)
}

func execiTypeFieldByIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Type).FieldByIndex(args[1].([]int))
	p.Ret(2, ret0)
}

func execiTypeFieldByName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(reflect.Type).FieldByName(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiTypeFieldByNameFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(reflect.Type).FieldByNameFunc(args[1].(func(string) bool))
	p.Ret(2, ret0, ret1)
}

func execiTypeImplements(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Type).Implements(toType0(args[1]))
	p.Ret(2, ret0)
}

func execiTypeIn(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Type).In(args[1].(int))
	p.Ret(2, ret0)
}

func execiTypeIsVariadic(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).IsVariadic()
	p.Ret(1, ret0)
}

func execiTypeKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).Key()
	p.Ret(1, ret0)
}

func execiTypeKind(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).Kind()
	p.Ret(1, ret0)
}

func execiTypeLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).Len()
	p.Ret(1, ret0)
}

func execiTypeMethod(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Type).Method(args[1].(int))
	p.Ret(2, ret0)
}

func execiTypeMethodByName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(reflect.Type).MethodByName(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiTypeName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).Name()
	p.Ret(1, ret0)
}

func execiTypeNumField(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).NumField()
	p.Ret(1, ret0)
}

func execiTypeNumIn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).NumIn()
	p.Ret(1, ret0)
}

func execiTypeNumMethod(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).NumMethod()
	p.Ret(1, ret0)
}

func execiTypeNumOut(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).NumOut()
	p.Ret(1, ret0)
}

func execiTypeOut(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Type).Out(args[1].(int))
	p.Ret(2, ret0)
}

func execiTypePkgPath(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).PkgPath()
	p.Ret(1, ret0)
}

func execiTypeSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).Size()
	p.Ret(1, ret0)
}

func execiTypeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Type).String()
	p.Ret(1, ret0)
}

func execTypeOf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflect.TypeOf(args[0])
	p.Ret(1, ret0)
}

func execmValueAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Addr()
	p.Ret(1, ret0)
}

func execmValueBool(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Bool()
	p.Ret(1, ret0)
}

func execmValueBytes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Bytes()
	p.Ret(1, ret0)
}

func execmValueCanAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).CanAddr()
	p.Ret(1, ret0)
}

func execmValueCanSet(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).CanSet()
	p.Ret(1, ret0)
}

func execmValueCall(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).Call(args[1].([]reflect.Value))
	p.Ret(2, ret0)
}

func execmValueCallSlice(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).CallSlice(args[1].([]reflect.Value))
	p.Ret(2, ret0)
}

func execmValueCap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Cap()
	p.Ret(1, ret0)
}

func execmValueClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(reflect.Value).Close()
	p.PopN(1)
}

func execmValueComplex(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Complex()
	p.Ret(1, ret0)
}

func execmValueElem(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Elem()
	p.Ret(1, ret0)
}

func execmValueField(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).Field(args[1].(int))
	p.Ret(2, ret0)
}

func execmValueFieldByIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).FieldByIndex(args[1].([]int))
	p.Ret(2, ret0)
}

func execmValueFieldByName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).FieldByName(args[1].(string))
	p.Ret(2, ret0)
}

func execmValueFieldByNameFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).FieldByNameFunc(args[1].(func(string) bool))
	p.Ret(2, ret0)
}

func execmValueFloat(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Float()
	p.Ret(1, ret0)
}

func execmValueIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).Index(args[1].(int))
	p.Ret(2, ret0)
}

func execmValueInt(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Int()
	p.Ret(1, ret0)
}

func execmValueCanInterface(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).CanInterface()
	p.Ret(1, ret0)
}

func execmValueInterface(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Interface()
	p.Ret(1, ret0)
}

func execmValueInterfaceData(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).InterfaceData()
	p.Ret(1, ret0)
}

func execmValueIsNil(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).IsNil()
	p.Ret(1, ret0)
}

func execmValueIsValid(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).IsValid()
	p.Ret(1, ret0)
}

func execmValueIsZero(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).IsZero()
	p.Ret(1, ret0)
}

func execmValueKind(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Kind()
	p.Ret(1, ret0)
}

func execmValueLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Len()
	p.Ret(1, ret0)
}

func execmValueMapIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).MapIndex(args[1].(reflect.Value))
	p.Ret(2, ret0)
}

func execmValueMapKeys(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).MapKeys()
	p.Ret(1, ret0)
}

func execmValueMapRange(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).MapRange()
	p.Ret(1, ret0)
}

func execmValueMethod(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).Method(args[1].(int))
	p.Ret(2, ret0)
}

func execmValueNumMethod(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).NumMethod()
	p.Ret(1, ret0)
}

func execmValueMethodByName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).MethodByName(args[1].(string))
	p.Ret(2, ret0)
}

func execmValueNumField(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).NumField()
	p.Ret(1, ret0)
}

func execmValueOverflowComplex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).OverflowComplex(args[1].(complex128))
	p.Ret(2, ret0)
}

func execmValueOverflowFloat(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).OverflowFloat(args[1].(float64))
	p.Ret(2, ret0)
}

func execmValueOverflowInt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).OverflowInt(args[1].(int64))
	p.Ret(2, ret0)
}

func execmValueOverflowUint(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).OverflowUint(args[1].(uint64))
	p.Ret(2, ret0)
}

func execmValuePointer(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Pointer()
	p.Ret(1, ret0)
}

func execmValueRecv(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(reflect.Value).Recv()
	p.Ret(1, ret0, ret1)
}

func execmValueSend(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).Send(args[1].(reflect.Value))
	p.PopN(2)
}

func execmValueSet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).Set(args[1].(reflect.Value))
	p.PopN(2)
}

func execmValueSetBool(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetBool(args[1].(bool))
	p.PopN(2)
}

func execmValueSetBytes(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetBytes(args[1].([]byte))
	p.PopN(2)
}

func execmValueSetComplex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetComplex(args[1].(complex128))
	p.PopN(2)
}

func execmValueSetFloat(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetFloat(args[1].(float64))
	p.PopN(2)
}

func execmValueSetInt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetInt(args[1].(int64))
	p.PopN(2)
}

func execmValueSetLen(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetLen(args[1].(int))
	p.PopN(2)
}

func execmValueSetCap(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetCap(args[1].(int))
	p.PopN(2)
}

func execmValueSetMapIndex(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(reflect.Value).SetMapIndex(args[1].(reflect.Value), args[2].(reflect.Value))
	p.PopN(3)
}

func execmValueSetUint(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetUint(args[1].(uint64))
	p.PopN(2)
}

func execmValueSetPointer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetPointer(args[1].(unsafe.Pointer))
	p.PopN(2)
}

func execmValueSetString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(reflect.Value).SetString(args[1].(string))
	p.PopN(2)
}

func execmValueSlice(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(reflect.Value).Slice(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmValueSlice3(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(reflect.Value).Slice3(args[1].(int), args[2].(int), args[3].(int))
	p.Ret(4, ret0)
}

func execmValueString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).String()
	p.Ret(1, ret0)
}

func execmValueTryRecv(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(reflect.Value).TryRecv()
	p.Ret(1, ret0, ret1)
}

func execmValueTrySend(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).TrySend(args[1].(reflect.Value))
	p.Ret(2, ret0)
}

func execmValueType(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Type()
	p.Ret(1, ret0)
}

func execmValueUint(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).Uint()
	p.Ret(1, ret0)
}

func execmValueUnsafeAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(reflect.Value).UnsafeAddr()
	p.Ret(1, ret0)
}

func execmValueConvert(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(reflect.Value).Convert(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmValueErrorError(_ int, p *gop.Context) {
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
		I.Func("(ChanDir).String", (reflect.ChanDir).String, execmChanDirString),
		I.Func("ChanOf", reflect.ChanOf, execChanOf),
		I.Func("Copy", reflect.Copy, execCopy),
		I.Func("DeepEqual", reflect.DeepEqual, execDeepEqual),
		I.Func("FuncOf", reflect.FuncOf, execFuncOf),
		I.Func("Indirect", reflect.Indirect, execIndirect),
		I.Func("(Kind).String", (reflect.Kind).String, execmKindString),
		I.Func("MakeChan", reflect.MakeChan, execMakeChan),
		I.Func("MakeFunc", reflect.MakeFunc, execMakeFunc),
		I.Func("MakeMap", reflect.MakeMap, execMakeMap),
		I.Func("MakeMapWithSize", reflect.MakeMapWithSize, execMakeMapWithSize),
		I.Func("MakeSlice", reflect.MakeSlice, execMakeSlice),
		I.Func("(*MapIter).Key", (*reflect.MapIter).Key, execmMapIterKey),
		I.Func("(*MapIter).Value", (*reflect.MapIter).Value, execmMapIterValue),
		I.Func("(*MapIter).Next", (*reflect.MapIter).Next, execmMapIterNext),
		I.Func("MapOf", reflect.MapOf, execMapOf),
		I.Func("New", reflect.New, execNew),
		I.Func("NewAt", reflect.NewAt, execNewAt),
		I.Func("PtrTo", reflect.PtrTo, execPtrTo),
		I.Func("Select", reflect.Select, execSelect),
		I.Func("SliceOf", reflect.SliceOf, execSliceOf),
		I.Func("StructOf", reflect.StructOf, execStructOf),
		I.Func("(StructTag).Get", (reflect.StructTag).Get, execmStructTagGet),
		I.Func("(StructTag).Lookup", (reflect.StructTag).Lookup, execmStructTagLookup),
		I.Func("Swapper", reflect.Swapper, execSwapper),
		I.Func("(Type).Align", (reflect.Type).Align, execiTypeAlign),
		I.Func("(Type).AssignableTo", (reflect.Type).AssignableTo, execiTypeAssignableTo),
		I.Func("(Type).Bits", (reflect.Type).Bits, execiTypeBits),
		I.Func("(Type).ChanDir", (reflect.Type).ChanDir, execiTypeChanDir),
		I.Func("(Type).Comparable", (reflect.Type).Comparable, execiTypeComparable),
		I.Func("(Type).ConvertibleTo", (reflect.Type).ConvertibleTo, execiTypeConvertibleTo),
		I.Func("(Type).Elem", (reflect.Type).Elem, execiTypeElem),
		I.Func("(Type).Field", (reflect.Type).Field, execiTypeField),
		I.Func("(Type).FieldAlign", (reflect.Type).FieldAlign, execiTypeFieldAlign),
		I.Func("(Type).FieldByIndex", (reflect.Type).FieldByIndex, execiTypeFieldByIndex),
		I.Func("(Type).FieldByName", (reflect.Type).FieldByName, execiTypeFieldByName),
		I.Func("(Type).FieldByNameFunc", (reflect.Type).FieldByNameFunc, execiTypeFieldByNameFunc),
		I.Func("(Type).Implements", (reflect.Type).Implements, execiTypeImplements),
		I.Func("(Type).In", (reflect.Type).In, execiTypeIn),
		I.Func("(Type).IsVariadic", (reflect.Type).IsVariadic, execiTypeIsVariadic),
		I.Func("(Type).Key", (reflect.Type).Key, execiTypeKey),
		I.Func("(Type).Kind", (reflect.Type).Kind, execiTypeKind),
		I.Func("(Type).Len", (reflect.Type).Len, execiTypeLen),
		I.Func("(Type).Method", (reflect.Type).Method, execiTypeMethod),
		I.Func("(Type).MethodByName", (reflect.Type).MethodByName, execiTypeMethodByName),
		I.Func("(Type).Name", (reflect.Type).Name, execiTypeName),
		I.Func("(Type).NumField", (reflect.Type).NumField, execiTypeNumField),
		I.Func("(Type).NumIn", (reflect.Type).NumIn, execiTypeNumIn),
		I.Func("(Type).NumMethod", (reflect.Type).NumMethod, execiTypeNumMethod),
		I.Func("(Type).NumOut", (reflect.Type).NumOut, execiTypeNumOut),
		I.Func("(Type).Out", (reflect.Type).Out, execiTypeOut),
		I.Func("(Type).PkgPath", (reflect.Type).PkgPath, execiTypePkgPath),
		I.Func("(Type).Size", (reflect.Type).Size, execiTypeSize),
		I.Func("(Type).String", (reflect.Type).String, execiTypeString),
		I.Func("TypeOf", reflect.TypeOf, execTypeOf),
		I.Func("(Value).Addr", (reflect.Value).Addr, execmValueAddr),
		I.Func("(Value).Bool", (reflect.Value).Bool, execmValueBool),
		I.Func("(Value).Bytes", (reflect.Value).Bytes, execmValueBytes),
		I.Func("(Value).CanAddr", (reflect.Value).CanAddr, execmValueCanAddr),
		I.Func("(Value).CanSet", (reflect.Value).CanSet, execmValueCanSet),
		I.Func("(Value).Call", (reflect.Value).Call, execmValueCall),
		I.Func("(Value).CallSlice", (reflect.Value).CallSlice, execmValueCallSlice),
		I.Func("(Value).Cap", (reflect.Value).Cap, execmValueCap),
		I.Func("(Value).Close", (reflect.Value).Close, execmValueClose),
		I.Func("(Value).Complex", (reflect.Value).Complex, execmValueComplex),
		I.Func("(Value).Elem", (reflect.Value).Elem, execmValueElem),
		I.Func("(Value).Field", (reflect.Value).Field, execmValueField),
		I.Func("(Value).FieldByIndex", (reflect.Value).FieldByIndex, execmValueFieldByIndex),
		I.Func("(Value).FieldByName", (reflect.Value).FieldByName, execmValueFieldByName),
		I.Func("(Value).FieldByNameFunc", (reflect.Value).FieldByNameFunc, execmValueFieldByNameFunc),
		I.Func("(Value).Float", (reflect.Value).Float, execmValueFloat),
		I.Func("(Value).Index", (reflect.Value).Index, execmValueIndex),
		I.Func("(Value).Int", (reflect.Value).Int, execmValueInt),
		I.Func("(Value).CanInterface", (reflect.Value).CanInterface, execmValueCanInterface),
		I.Func("(Value).Interface", (reflect.Value).Interface, execmValueInterface),
		I.Func("(Value).InterfaceData", (reflect.Value).InterfaceData, execmValueInterfaceData),
		I.Func("(Value).IsNil", (reflect.Value).IsNil, execmValueIsNil),
		I.Func("(Value).IsValid", (reflect.Value).IsValid, execmValueIsValid),
		I.Func("(Value).IsZero", (reflect.Value).IsZero, execmValueIsZero),
		I.Func("(Value).Kind", (reflect.Value).Kind, execmValueKind),
		I.Func("(Value).Len", (reflect.Value).Len, execmValueLen),
		I.Func("(Value).MapIndex", (reflect.Value).MapIndex, execmValueMapIndex),
		I.Func("(Value).MapKeys", (reflect.Value).MapKeys, execmValueMapKeys),
		I.Func("(Value).MapRange", (reflect.Value).MapRange, execmValueMapRange),
		I.Func("(Value).Method", (reflect.Value).Method, execmValueMethod),
		I.Func("(Value).NumMethod", (reflect.Value).NumMethod, execmValueNumMethod),
		I.Func("(Value).MethodByName", (reflect.Value).MethodByName, execmValueMethodByName),
		I.Func("(Value).NumField", (reflect.Value).NumField, execmValueNumField),
		I.Func("(Value).OverflowComplex", (reflect.Value).OverflowComplex, execmValueOverflowComplex),
		I.Func("(Value).OverflowFloat", (reflect.Value).OverflowFloat, execmValueOverflowFloat),
		I.Func("(Value).OverflowInt", (reflect.Value).OverflowInt, execmValueOverflowInt),
		I.Func("(Value).OverflowUint", (reflect.Value).OverflowUint, execmValueOverflowUint),
		I.Func("(Value).Pointer", (reflect.Value).Pointer, execmValuePointer),
		I.Func("(Value).Recv", (reflect.Value).Recv, execmValueRecv),
		I.Func("(Value).Send", (reflect.Value).Send, execmValueSend),
		I.Func("(Value).Set", (reflect.Value).Set, execmValueSet),
		I.Func("(Value).SetBool", (reflect.Value).SetBool, execmValueSetBool),
		I.Func("(Value).SetBytes", (reflect.Value).SetBytes, execmValueSetBytes),
		I.Func("(Value).SetComplex", (reflect.Value).SetComplex, execmValueSetComplex),
		I.Func("(Value).SetFloat", (reflect.Value).SetFloat, execmValueSetFloat),
		I.Func("(Value).SetInt", (reflect.Value).SetInt, execmValueSetInt),
		I.Func("(Value).SetLen", (reflect.Value).SetLen, execmValueSetLen),
		I.Func("(Value).SetCap", (reflect.Value).SetCap, execmValueSetCap),
		I.Func("(Value).SetMapIndex", (reflect.Value).SetMapIndex, execmValueSetMapIndex),
		I.Func("(Value).SetUint", (reflect.Value).SetUint, execmValueSetUint),
		I.Func("(Value).SetPointer", (reflect.Value).SetPointer, execmValueSetPointer),
		I.Func("(Value).SetString", (reflect.Value).SetString, execmValueSetString),
		I.Func("(Value).Slice", (reflect.Value).Slice, execmValueSlice),
		I.Func("(Value).Slice3", (reflect.Value).Slice3, execmValueSlice3),
		I.Func("(Value).String", (reflect.Value).String, execmValueString),
		I.Func("(Value).TryRecv", (reflect.Value).TryRecv, execmValueTryRecv),
		I.Func("(Value).TrySend", (reflect.Value).TrySend, execmValueTrySend),
		I.Func("(Value).Type", (reflect.Value).Type, execmValueType),
		I.Func("(Value).Uint", (reflect.Value).Uint, execmValueUint),
		I.Func("(Value).UnsafeAddr", (reflect.Value).UnsafeAddr, execmValueUnsafeAddr),
		I.Func("(Value).Convert", (reflect.Value).Convert, execmValueConvert),
		I.Func("(*ValueError).Error", (*reflect.ValueError).Error, execmValueErrorError),
		I.Func("ValueOf", reflect.ValueOf, execValueOf),
		I.Func("Zero", reflect.Zero, execZero),
	)
	I.RegisterFuncvs(
		I.Funcv("Append", reflect.Append, execAppend),
	)
	I.RegisterTypes(
		I.Type("ChanDir", reflect.TypeOf((*reflect.ChanDir)(nil)).Elem()),
		I.Type("Kind", reflect.TypeOf((*reflect.Kind)(nil)).Elem()),
		I.Type("MapIter", reflect.TypeOf((*reflect.MapIter)(nil)).Elem()),
		I.Type("Method", reflect.TypeOf((*reflect.Method)(nil)).Elem()),
		I.Type("SelectCase", reflect.TypeOf((*reflect.SelectCase)(nil)).Elem()),
		I.Type("SelectDir", reflect.TypeOf((*reflect.SelectDir)(nil)).Elem()),
		I.Type("SliceHeader", reflect.TypeOf((*reflect.SliceHeader)(nil)).Elem()),
		I.Type("StringHeader", reflect.TypeOf((*reflect.StringHeader)(nil)).Elem()),
		I.Type("StructField", reflect.TypeOf((*reflect.StructField)(nil)).Elem()),
		I.Type("StructTag", reflect.TypeOf((*reflect.StructTag)(nil)).Elem()),
		I.Type("Type", reflect.TypeOf((*reflect.Type)(nil)).Elem()),
		I.Type("Value", reflect.TypeOf((*reflect.Value)(nil)).Elem()),
		I.Type("ValueError", reflect.TypeOf((*reflect.ValueError)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("Array", qspec.Uint, reflect.Array),
		I.Const("Bool", qspec.Uint, reflect.Bool),
		I.Const("BothDir", qspec.Int, reflect.BothDir),
		I.Const("Chan", qspec.Uint, reflect.Chan),
		I.Const("Complex128", qspec.Uint, reflect.Complex128),
		I.Const("Complex64", qspec.Uint, reflect.Complex64),
		I.Const("Float32", qspec.Uint, reflect.Float32),
		I.Const("Float64", qspec.Uint, reflect.Float64),
		I.Const("Func", qspec.Uint, reflect.Func),
		I.Const("Int", qspec.Uint, reflect.Int),
		I.Const("Int16", qspec.Uint, reflect.Int16),
		I.Const("Int32", qspec.Uint, reflect.Int32),
		I.Const("Int64", qspec.Uint, reflect.Int64),
		I.Const("Int8", qspec.Uint, reflect.Int8),
		I.Const("Interface", qspec.Uint, reflect.Interface),
		I.Const("Invalid", qspec.Uint, reflect.Invalid),
		I.Const("Map", qspec.Uint, reflect.Map),
		I.Const("Ptr", qspec.Uint, reflect.Ptr),
		I.Const("RecvDir", qspec.Int, reflect.RecvDir),
		I.Const("SelectDefault", qspec.Int, reflect.SelectDefault),
		I.Const("SelectRecv", qspec.Int, reflect.SelectRecv),
		I.Const("SelectSend", qspec.Int, reflect.SelectSend),
		I.Const("SendDir", qspec.Int, reflect.SendDir),
		I.Const("Slice", qspec.Uint, reflect.Slice),
		I.Const("String", qspec.Uint, reflect.String),
		I.Const("Struct", qspec.Uint, reflect.Struct),
		I.Const("Uint", qspec.Uint, reflect.Uint),
		I.Const("Uint16", qspec.Uint, reflect.Uint16),
		I.Const("Uint32", qspec.Uint, reflect.Uint32),
		I.Const("Uint64", qspec.Uint, reflect.Uint64),
		I.Const("Uint8", qspec.Uint, reflect.Uint8),
		I.Const("Uintptr", qspec.Uint, reflect.Uintptr),
		I.Const("UnsafePointer", qspec.Uint, reflect.UnsafePointer),
	)
}
