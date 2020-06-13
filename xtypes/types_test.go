package xtypes

import (
	"errors"
	"go/types"
	"reflect"
	"testing"
)

var (
	pkg  = types.NewPackage("github.com/qiniu/foo", "foo")
	tFlt = reflect.TypeOf((*Flt)(nil)).Elem()
)

type finder struct{}
type Flt float64

func (p *finder) FindGoType(pkgPath string, namedType string) (reflect.Type, bool) {
	if pkgPath == "github.com/qiniu/foo" {
		if namedType == "Flt" {
			return tFlt, true
		}
	}
	return nil, false
}

func TestTypes(t *testing.T) {
	tFloat64 := basicTypes[types.Float64]
	tyFloat64 := types.Typ[types.Float64]
	if ty, err := ToType(tyFloat64, nil); err != nil || ty != tFloat64 {
		t.Fatal("TestTypes Float64 failed:", ty, err)
	}

	tyUntypedFloat := types.Typ[types.UntypedFloat]
	if _, err := ToType(tyUntypedFloat, nil); err != ErrUntyped {
		t.Fatal("TestTypes UntypedFloat failed:", err)
	}

	tyFloat64Ptr := types.NewPointer(tyFloat64)
	if ty, err := ToType(tyFloat64Ptr, nil); err != nil || ty != reflect.PtrTo(tFloat64) {
		t.Fatal("TestTypes Pointer failed:", ty, err)
	}

	tyUntypedFloatPtr := types.NewPointer(tyUntypedFloat)
	if _, err := ToType(tyUntypedFloatPtr, nil); !errors.Is(err, ErrUntyped) {
		t.Fatal("TestTypes UntypedFloat Pointer failed:", err)
	}

	tyFloat64Slice := types.NewSlice(tyFloat64)
	if ty, err := ToType(tyFloat64Slice, nil); err != nil || ty != reflect.SliceOf(tFloat64) {
		t.Fatal("TestTypes Slice failed:", ty, err)
	}

	tyUntypedFloatSlice := types.NewSlice(tyUntypedFloat)
	if _, err := ToType(tyUntypedFloatSlice, nil); !errors.Is(err, ErrUntyped) {
		t.Fatal("TestTypes UntypedFloat Slice failed:", err)
	}

	tyFloat64Array := types.NewArray(tyFloat64, 10)
	if ty, err := ToType(tyFloat64Array, nil); err != nil || ty != reflect.ArrayOf(10, tFloat64) {
		t.Fatal("TestTypes Array failed:", ty, err)
	}

	tyFloat64Array2 := types.NewArray(tyFloat64, -1)
	if _, err := ToType(tyFloat64Array2, nil); !errors.Is(err, ErrUnknownArrayLen) {
		t.Fatal("TestTypes Array failed:", err)
	}

	tyUntypedFloatArray := types.NewArray(tyUntypedFloat, 10)
	if _, err := ToType(tyUntypedFloatArray, nil); !errors.Is(err, ErrUntyped) {
		t.Fatal("TestTypes UntypedFloat Array failed:", err)
	}

	tyFloat64Map := types.NewMap(tyFloat64, tyFloat64)
	if ty, err := ToType(tyFloat64Map, nil); err != nil || ty != reflect.MapOf(tFloat64, tFloat64) {
		t.Fatal("TestTypes Map failed:", ty, err)
	}

	tyUntypedFloatMap := types.NewMap(tyUntypedFloat, tyFloat64)
	if _, err := ToType(tyUntypedFloatMap, nil); !errors.Is(err, ErrUntyped) {
		t.Fatal("TestTypes UntypedFloat Map failed:", err)
	}

	tyUntypedFloatMap2 := types.NewMap(tyFloat64, tyUntypedFloat)
	if _, err := ToType(tyUntypedFloatMap2, nil); !errors.Is(err, ErrUntyped) {
		t.Fatal("TestTypes UntypedFloat Map failed:", err)
	}

	tyFloat64Chan := types.NewChan(types.SendRecv, tyFloat64)
	if ty, err := ToType(tyFloat64Chan, nil); err != nil || ty != reflect.ChanOf(reflect.BothDir, tFloat64) {
		t.Fatal("TestTypes Chan failed:", ty, err)
	}

	tyFloat64Chan0 := types.NewChan(types.RecvOnly, tyFloat64)
	if ty, err := ToType(tyFloat64Chan0, nil); err != nil || ty != reflect.ChanOf(reflect.RecvDir, tFloat64) {
		t.Fatal("TestTypes Chan failed:", ty, err)
	}

	tyFloat64Chan1 := types.NewChan(types.SendOnly, tyFloat64)
	if ty, err := ToType(tyFloat64Chan1, nil); err != nil || ty != reflect.ChanOf(reflect.SendDir, tFloat64) {
		t.Fatal("TestTypes Chan failed:", ty, err)
	}

	tyFloat64Chan2 := types.NewChan(types.SendOnly, tyUntypedFloat)
	if _, err := ToType(tyFloat64Chan2, nil); !errors.Is(err, ErrUntyped) {
		t.Fatal("TestTypes Chan failed:", err)
	}

	tFlds := []reflect.StructField{
		{Name: "x", PkgPath: pkg.Path(), Type: tFloat64},
	}
	fld1 := types.NewVar(0, pkg, "x", tyFloat64)
	flds := []*types.Var{fld1}
	tyBar := types.NewStruct(flds, nil)
	if ty, err := ToType(tyBar, nil); err != nil || ty != reflect.StructOf(tFlds) {
		t.Fatal("TestTypes Struct failed:", ty, err)
	}

	fld2 := types.NewVar(0, pkg, "x", tyUntypedFloat)
	flds2 := []*types.Var{fld2}
	tyBar2 := types.NewStruct(flds2, nil)
	if _, err := ToType(tyBar2, nil); !errors.Is(err, ErrUntyped) {
		t.Fatal("TestTypes Struct failed:", err)
	}

	tFlds3 := []reflect.StructField{
		{Name: "X", Type: tFloat64},
	}
	fld3 := types.NewVar(0, pkg, "X", tyFloat64)
	flds3 := []*types.Var{fld3}
	tyBar3 := types.NewStruct(flds3, nil)
	if ty, err := ToType(tyBar3, nil); err != nil || ty != reflect.StructOf(tFlds3) {
		t.Fatal("TestTypes Struct failed:", ty, err)
	}

	tFlds4 := []reflect.StructField{
		{Name: "Float64", Type: tFloat64, Anonymous: true},
	}
	fld4 := types.NewField(0, pkg, "Float64", tyFloat64, true)
	flds4 := []*types.Var{fld4}
	tyBar4 := types.NewStruct(flds4, nil)
	if ty, err := ToType(tyBar4, nil); err != nil || ty != reflect.StructOf(tFlds4) {
		t.Fatal("TestTypes Struct failed:", ty, err)
	}

	barArgs := types.NewTuple(types.NewParam(0, pkg, "v", tyFloat64))
	sigBar := types.NewSignature(nil, barArgs, nil, false)
	fnBar := types.NewFunc(0, pkg, "Bar", sigBar)
	methods := []*types.Func{fnBar}
	tyIntf := types.NewInterface(methods, nil)
	if ty, err := ToType(tyIntf, nil); err != nil || ty != tyEmptyInterface {
		t.Fatal("TestTypes Interface failed:", ty, err)
	}

	gf := new(finder)
	fltName := types.NewTypeName(0, pkg, "Flt", nil)
	tyFlt := types.NewNamed(fltName, tyFloat64, nil)
	if ty, err := ToType(tyFlt, gf); err != nil || ty != tFlt {
		t.Fatal("TestTypes Named failed:", ty, err)
	}

	fltName2 := types.NewTypeName(0, pkg, "Flt2", nil)
	tyFlt2 := types.NewNamed(fltName2, tyFloat64, nil)
	if ty, err := ToType(tyFlt2, gf); err != nil || ty != tFloat64 {
		t.Fatal("TestTypes Named failed:", ty, err)
	}

	fltName3 := types.NewTypeName(0, pkg, "Flt3", nil)
	tyFlt3 := types.NewNamed(fltName3, tyUntypedFloat, nil)
	if _, err := ToType(tyFlt3, gf); !errors.Is(err, ErrUntyped) {
		t.Fatal("TestTypes Named failed:", err)
	}
}
