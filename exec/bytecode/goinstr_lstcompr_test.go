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
	"testing"

	"github.com/goplus/gop/exec.spec"
)

func TestLargeSlice(t *testing.T) {
	b := newBuilder()
	ret := []string{}
	for i := 0; i < bitsFuncvArityMax+1; i++ {
		b.Push("32")
		ret = append(ret, "32")
	}
	code := b.
		MakeArray(reflect.SliceOf(TyString), bitsFuncvArityMax+1).
		MakeArray(reflect.SliceOf(TyString), -1).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, ret) {
		t.Fatal("32 times(1024) mkslice != `32` times(1024) slice, ret =", v)
	}
}

func TestLargeArray(t *testing.T) {
	b := newBuilder()
	ret := [bitsFuncvArityMax + 1]string{}
	for i := 0; i < bitsFuncvArityMax+1; i++ {
		b.Push("32")
		ret[i] = "32"
	}
	code := b.
		MakeArray(reflect.ArrayOf(bitsFuncvArityMax+1, TyString), bitsFuncvArityMax+1).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, ret) {
		t.Fatal("32 times(1024) mkslice != `32` times(1024) slice, ret =", v)
	}
}

func TestMap(t *testing.T) {
	code := newBuilder().
		Push("Hello").
		Push(3.2).
		Push("xsw").
		Push(1.0).
		MakeMap(reflect.MapOf(TyString, TyFloat64), 2).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, map[string]float64{"Hello": 3.2, "xsw": 1.0}) {
		t.Fatal("expected: {`Hello`: 3.2, `xsw`: 1.0}, ret =", v)
	}
}

func TestMapIndex(t *testing.T) {
	m := NewVar(reflect.MapOf(TyString, TyFloat64), "")
	code := newBuilder().
		DefineVar(m).
		Push("Hello").
		Push(3.2).
		Push("xsw").
		Push(1.0).
		MakeMap(reflect.MapOf(TyString, TyFloat64), 2).
		StoreVar(m).
		LoadVar(m).
		Push("go+").
		MapIndex().
		LoadVar(m).
		Push("xsw").
		MapIndex().
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Get(-1); v != 1.0 {
		t.Fatal("{`Hello`: 3.2, `xsw`: 1.0}[`xsw`] != 1.0, ret =", v)
	}
	if v := ctx.Get(-2); v != 0.0 {
		t.Fatal("{`Hello`: 3.2, `xsw`: 1.0}[`go+`] != 1.0, ret =", v)
	}
}

func TestSetMapIndex(t *testing.T) {
	a := NewVar(reflect.MapOf(TyString, TyFloat64), "")
	code := newBuilder().
		DefineVar(a).
		Push(2.0).
		Push("Hello").
		Push(3.2).
		Push("xsw").
		Push(1.0).
		MakeMap(reflect.MapOf(TyString, TyFloat64), 2).
		StoreVar(a).
		LoadVar(a).
		Push("xsw").
		SetMapIndex().
		LoadVar(a).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, map[string]float64{"Hello": 3.2, "xsw": 2.0}) {
		t.Fatal("expected: {`Hello`: 3.2, `xsw`: 2.0}, ret =", v)
	}
}

func TestMapComprehension(t *testing.T) {
	typData := reflect.MapOf(TyString, TyInt)
	key := NewVar(TyString, "k")
	val := NewVar(TyInt, "v")
	f := defaultImpl.NewForPhrase(typData).(*ForPhrase)
	c := defaultImpl.NewComprehension(reflect.MapOf(TyInt, TyString)).(*Comprehension)
	code := newBuilder().
		MapComprehension(c).
		Push("Hello").
		Push(3).
		Push("xsw").
		Push(1).
		MakeMap(typData, 2).
		ForPhrase(f, key, val, true).
		LoadVar(val).
		LoadVar(key).
		EndForPhrase(f).
		EndComprehension(c).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, map[int]string{3: "Hello", 1: "xsw"}) {
		t.Fatal(`expected: {3: "Hello", 1: "xsw"}, ret =`, v)
	}
}

func TestMapComprehensionFilter(t *testing.T) {
	typData := reflect.MapOf(TyString, TyInt)
	key := NewVar(TyString, "k")
	val := NewVar(TyInt, "v")
	f := NewForPhrase(typData)
	c := NewComprehension(reflect.MapOf(TyInt, TyString))
	code := newBuilder().
		MapComprehension(c).
		Push("Hello").
		Push(3).
		Push("xsw").
		Push(1).
		MakeMap(typData, 2).
		ForPhrase(f, key, val).
		LoadVar(val).
		Push(2).
		BuiltinOp(Int, OpLE).
		FilterForPhrase(f).
		LoadVar(val).
		LoadVar(key).
		EndForPhrase(f).
		EndComprehension(c).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, map[int]string{1: "xsw"}) {
		t.Fatal(`expected: {1: "xsw"}, ret =`, v)
	}
}

func TestListComprehension(t *testing.T) {
	typData := reflect.ArrayOf(4, TyInt)
	x := NewVar(TyInt, "x")
	f := NewForPhrase(typData)
	c := NewComprehension(reflect.SliceOf(TyInt))
	code := newBuilder().
		ListComprehension(c).
		Push(1).
		Push(3).
		Push(5).
		Push(7).
		MakeArray(typData, 4).
		ForPhrase(f, nil, x).
		LoadVar(x).
		LoadVar(x).
		BuiltinOp(Int, OpMul).
		EndForPhrase(f).
		EndComprehension(c).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, []int{1, 9, 25, 49}) {
		t.Fatal(`expected: [1, 9, 25, 49], ret =`, v)
	}
}

func TestListComprehensionFilter(t *testing.T) {
	typData := reflect.ArrayOf(4, TyInt)
	x := NewVar(TyInt, "x")
	f := NewForPhrase(typData)
	c := NewComprehension(reflect.SliceOf(TyInt))
	code := newBuilder().
		ListComprehension(c).
		Push(1).
		Push(3).
		Push(5).
		Push(7).
		MakeArray(typData, 4).
		ForPhrase(f, nil, x).
		LoadVar(x).
		Push(3).
		BuiltinOp(Int, OpGT). // x > 3
		FilterForPhrase(f).
		LoadVar(x).
		LoadVar(x).
		BuiltinOp(Int, OpMul).
		EndForPhrase(f).
		EndComprehension(c).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, []int{25, 49}) {
		t.Fatal(`expected: [25, 49], ret =`, v)
	}
}

func TestMapComprehension2(t *testing.T) {
	typData := reflect.SliceOf(TyInt)
	i := NewVar(TyInt, "i")
	x := NewVar(TyInt, "x")
	f := NewForPhrase(typData)
	c := NewComprehension(reflect.MapOf(TyInt, TyInt))
	code := newBuilder().
		MapComprehension(c).
		Push(1).
		Push(3).
		Push(5).
		Push(7).
		MakeArray(typData, 4).
		ForPhrase(f, i, x).
		LoadVar(x).
		LoadVar(x).
		BuiltinOp(Int, OpMul).
		LoadVar(i).
		EndForPhrase(f).
		EndComprehension(c).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, map[int]int{1: 0, 9: 1, 25: 2, 49: 3}) {
		t.Fatal(`expected: {1: 0, 9: 1, 25: 2, 49: 3}, ret =`, v)
	}
}

func TestListComprehensionEx(t *testing.T) {
	typData := reflect.SliceOf(TyInt)
	a := NewVar(TyInt, "a")
	b := NewVar(TyInt, "b")
	fa := NewForPhrase(typData)
	fb := NewForPhrase(typData)
	c := NewComprehension(typData)
	code := newBuilder().
		ListComprehension(c).
		Push(5).
		Push(6).
		Push(7).
		MakeArray(typData, 3).
		ForPhrase(fb, nil, b).
		Push(1).
		Push(2).
		Push(3).
		Push(4).
		MakeArray(typData, 4).
		ForPhrase(fa, nil, a).
		LoadVar(a).
		Push(1).
		BuiltinOp(Int, OpGT). // a > 1
		FilterForPhrase(fa).
		LoadVar(a).
		LoadVar(b).
		BuiltinOp(Int, OpMul).
		EndForPhrase(fa).
		EndForPhrase(fb).
		EndComprehension(c).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, []int{10, 15, 20, 12, 18, 24, 14, 21, 28}) {
		t.Fatal(`expected: [10, 15, 20, 12, 18, 24, 14, 21, 28], ret =`, v)
	}
}

func TestZero(t *testing.T) {
	code := newBuilder().
		Zero(TyFloat64).
		Push(3.2).
		BuiltinOp(Float64, OpAdd).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 3.2 {
		t.Fatal("0 + 3.2 != 3.2, ret =", v)
	}
}

func TestIndex(t *testing.T) {
	code := newBuilder().
		Push(3.2).
		Push(1.2).
		Push(2.4).
		MakeArray(reflect.SliceOf(TyFloat64), 3).
		Index(1).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 1.2 {
		t.Fatal("[3.2, 1.2, 2.4][1] != 1.2, ret:", v)
	}
}

func TestIndex2(t *testing.T) {
	code := newBuilder().
		Push(3.2).
		Push(1.2).
		Push(2.4).
		MakeArray(reflect.SliceOf(TyFloat64), 3).
		Push(2).
		Index(-1).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 2.4 {
		t.Fatal("[3.2, 1.2, 2.4][2] != 2.4, ret:", v)
	}
}

func TestAddrIndex(t *testing.T) {
	a := NewVar(reflect.SliceOf(TyFloat64), "")
	code := newBuilder().
		DefineVar(a).
		Push(0.7).
		Push(3.2).
		Push(1.2).
		Push(2.4).
		MakeArray(reflect.SliceOf(TyFloat64), 3).
		StoreVar(a).
		LoadVar(a).
		Push(2).
		SetIndex(-1).
		LoadVar(a).
		AddrIndex(2).
		AddrOp(Float64, OpAddrVal).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 0.7 {
		t.Fatal("[3.2, 1.2, 0.7], ret:", v)
	}
}

func TestAddrLargeIndex(t *testing.T) {
	a := NewVar(reflect.SliceOf(TyFloat64), "")
	code := newBuilder().
		DefineVar(a).
		Push(bitsOpIndexOperand+1).
		Make(reflect.SliceOf(TyFloat64), 1).
		StoreVar(a).
		Push(1.7).
		LoadVar(a).
		SetIndex(bitsOpIndexOperand).
		LoadVar(a).
		AddrIndex(bitsOpIndexOperand).
		AddrOp(Float64, OpAddrVal).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 1.7 {
		t.Fatal("v != 1.7, ret:", v)
	}
}

func TestSetIndex(t *testing.T) {
	a := NewVar(reflect.SliceOf(TyFloat64), "")
	code := newBuilder().
		DefineVar(a).
		Push(0.7).
		Push(3.2).
		Push(1.2).
		Push(2.4).
		MakeArray(reflect.SliceOf(TyFloat64), 3).
		StoreVar(a).
		LoadVar(a).
		Push(2).
		SetIndex(-1).
		LoadVar(a).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, []float64{3.2, 1.2, 0.7}) {
		t.Fatal("[3.2, 1.2, 0.7], ret:", v)
	}
}

func TestSetLargeIndex(t *testing.T) {
	a := NewVar(reflect.SliceOf(TyFloat64), "")
	code := newBuilder().
		DefineVar(a).
		Push(bitsOpIndexOperand+1).
		Make(reflect.SliceOf(TyFloat64), 1).
		StoreVar(a).
		Push(1.7).
		LoadVar(a).
		SetIndex(bitsOpIndexOperand).
		LoadVar(a).
		Index(bitsOpIndexOperand).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 1.7 {
		t.Fatal("v != 1.7, ret:", v)
	}
}

func TestSlice(t *testing.T) {
	code := newBuilder().
		Push(3.2).
		Push(1.2).
		Push(2.4).
		MakeArray(reflect.SliceOf(TyFloat64), 3).
		Slice(0, 2).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, []float64{3.2, 1.2}) {
		t.Fatal("[3.2, 1.2, 2.4][0:2] != [3.2, 1.2], ret:", v)
	}
}

func TestSliceLarge(t *testing.T) {
	a := NewVar(reflect.SliceOf(TyFloat64), "")
	code := newBuilder().
		DefineVar(a).
		Push(SliceConstIndexLast+1).
		Make(reflect.SliceOf(TyFloat64), 1).
		StoreVar(a).
		Push(1.7).
		LoadVar(a).
		SetIndex(SliceConstIndexLast).
		LoadVar(a).
		Slice(SliceConstIndexLast, SliceConstIndexLast+1).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, []float64{1.7}) {
		t.Fatal("ret != [1.7], ret:", v)
	}
}

func TestSlice2(t *testing.T) {
	code := newBuilder().
		Push(3.2).
		Push(1.2).
		Push(2.4).
		MakeArray(reflect.SliceOf(TyFloat64), 3).
		Push(1).
		Slice(SliceDefaultIndex, -1).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, []float64{3.2}) {
		t.Fatal("[3.2, 1.2, 2.4][:1] != [3.2], ret:", v)
	}
}

func TestSlice3(t *testing.T) {
	code := newBuilder().
		Push(3.2).
		Push(1.2).
		Push(2.4).
		MakeArray(reflect.SliceOf(TyFloat64), 3).
		Push(1).
		Slice(-1, SliceDefaultIndex).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, []float64{1.2, 2.4}) {
		t.Fatal("[3.2, 1.2, 2.4][1:] != [1.2, 2.4], ret:", v)
	}
}

func TestSlice4(t *testing.T) {
	code := newBuilder().
		Push(3.2).
		Push(1.2).
		Push(2.4).
		MakeArray(reflect.SliceOf(TyFloat64), 3).
		Push(1).
		Slice3(SliceDefaultIndex, -1, 2).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, []float64{3.2}) || reflect.ValueOf(v).Cap() != 2 {
		t.Fatal("[3.2, 1.2, 2.4][:1] != [3.2], ret:", v)
	}
}

func TestAppend(t *testing.T) {
	code := newBuilder().
		Zero(reflect.SliceOf(TyFloat64)).
		Push(3.2).
		Push(1.2).
		Push(2.4).
		Append(TyFloat64, 4).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, []float64{3.2, 1.2, 2.4}) {
		t.Fatal("ret != [3.2, 1.2, 2.4], ret:", v)
	}
}

func TestAppend2(t *testing.T) {
	sliceTy := reflect.SliceOf(TyFloat64)
	code := newBuilder().
		Zero(sliceTy).
		Push(3.2).
		Push(1.2).
		Push(2.4).
		MakeArray(sliceTy, 3).
		Append(TyFloat64, -1).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, []float64{3.2, 1.2, 2.4}) {
		t.Fatal("ret != [3.2, 1.2, 2.4], ret:", v)
	}
}

func TestAppend3(t *testing.T) {
	sliceTy := reflect.SliceOf(TyFloat64)
	code := newBuilder().
		New(sliceTy).
		AddrOp(reflect.Slice, exec.OpAddrVal).
		Push(3.2).
		Push(1.2).
		Push(2.4).
		MakeArray(sliceTy, 3).
		Append(TyFloat64, -1).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, []float64{3.2, 1.2, 2.4}) {
		t.Fatal("ret != [3.2, 1.2, 2.4], ret:", v)
	}
}

func TestMakeSlice(t *testing.T) {
	code := newBuilder().
		Push(2).
		Make(reflect.SliceOf(TyFloat64), 1).
		Push(3.2).
		Push(1.2).
		Push(2.4).
		Append(TyFloat64, 4).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, []float64{0, 0, 3.2, 1.2, 2.4}) {
		t.Fatal("ret != [0, 0, 3.2, 1.2, 2.4], ret:", v)
	}
}

func TestMakeSlice2(t *testing.T) {
	code := newBuilder().
		Push(2).
		Push(4).
		Make(reflect.SliceOf(TyFloat64), 2).
		Push(3.2).
		Push(1.2).
		Push(2.4).
		Append(TyFloat64, 4).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, []float64{0, 0, 3.2, 1.2, 2.4}) {
		t.Fatal("ret != [0, 0, 3.2, 1.2, 2.4], ret:", v)
	}
}

func TestMakeMap(t *testing.T) {
	code := newBuilder().
		Make(reflect.MapOf(TyInt, TyFloat64), 0).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, map[int]float64{}) {
		t.Fatal("ret != {}, ret:", v)
	}
}

func TestMakeMap2(t *testing.T) {
	code := newBuilder().
		Push(2).
		Make(reflect.MapOf(TyInt, TyFloat64), 1).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, map[int]float64{}) {
		t.Fatal("ret != {}, ret:", v)
	}
}

func TestMakeChan(t *testing.T) {
	typ := reflect.ChanOf(reflect.BothDir, TyInt)
	code := newBuilder().
		Make(typ, 0).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); reflect.TypeOf(v) != typ {
		t.Fatal("ret != chan int, ret:", v)
	}
}

func TestMakeChan2(t *testing.T) {
	typ := reflect.ChanOf(reflect.BothDir, TyInt)
	code := newBuilder().
		Push(2).
		Make(typ, 1).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); reflect.TypeOf(v) != typ {
		t.Fatal("ret != chan int, ret:", v)
	}
}

func TestTypeCast(t *testing.T) {
	code := newBuilder().
		Push(byte('5')).
		TypeCast(TyUint8, TyString).
		Push("6").
		BuiltinOp(String, OpAdd).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "56" {
		t.Fatal("`5` `6` add != `56`, ret =", v)
	}
}

func TestDelete(t *testing.T) {
	tyMap := reflect.MapOf(TyString, TyInt)
	m := NewVar(tyMap, "")
	code := newBuilder().
		DefineVar(m).
		Push("Hello").
		Push(1).
		Push("Go+").
		Push(2).
		MakeMap(tyMap, 2).
		StoreVar(m).
		LoadVar(m).
		Push("Hello").
		GoBuiltin(tyMap, GobDelete).
		LoadVar(m).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, map[string]int{"Go+": 2}) {
		t.Fatal("expected: {`Go+`: 2}, ret =", v)
	}
}

func TestCopy(t *testing.T) {
	a := NewVar(reflect.SliceOf(TyInt), "")
	b := NewVar(reflect.SliceOf(TyInt), "")
	code := newBuilder().
		DefineVar(a).
		DefineVar(b).
		Push(1).
		Push(2).
		Push(3).
		MakeArray(reflect.SliceOf(TyInt), 3).
		StoreVar(a).
		Push(11).
		MakeArray(reflect.SliceOf(TyInt), 1).
		StoreVar(b).
		LoadVar(b).
		LoadVar(a).
		GoBuiltin(nil, GobCopy).
		LoadVar(b).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	arr := ctx.Pop().([]int)
	n := ctx.Pop().(int)
	if n != 1 || reflect.DeepEqual(arr, []int{11, 2, 3}) {
		t.Fatal("copy failed")
	}
}

func TestCopy2(t *testing.T) {
	a := NewVar(reflect.SliceOf(TyByte), "")
	code := newBuilder().
		DefineVar(a).
		Push(byte(96)).
		Push(byte(97)).
		MakeArray(reflect.SliceOf(TyByte), 2).
		StoreVar(a).
		LoadVar(a).
		Push("hello").
		GoBuiltin(nil, GobCopy).
		LoadVar(a).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	arr := ctx.Pop().([]byte)
	n := ctx.Pop().(int)
	if n != 2 || string(arr) != "he" {
		t.Fatal("copy failed")
	}
}
