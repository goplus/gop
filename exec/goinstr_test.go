package exec

import (
	"reflect"
	"testing"
)

func TestLargeSlice(t *testing.T) {
	b := NewBuilder(nil)
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
	b := NewBuilder(nil)
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
	code := NewBuilder(nil).
		Push("Hello").
		Push(3.2).
		Push("xsw").
		Push(1.0).
		MakeMap(reflect.MapOf(TyString, TyFloat64), 2).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, map[string]float64{"Hello": 3.2, "xsw": 1.0}) {
		t.Fatal("expected: {`Hello`: 3.2, `xsw`: 1}, ret =", v)
	}
}

func TestMapComprehension(t *testing.T) {
	typData := reflect.MapOf(TyString, TyInt)
	key := NewVar(TyString, "k")
	val := NewVar(TyInt, "v")
	f := NewForPhrase(key, val, typData)
	c := NewComprehension(reflect.MapOf(TyInt, TyString))
	code := NewBuilder(nil).
		MapComprehension(c).
		Push("Hello").
		Push(3).
		Push("xsw").
		Push(1).
		MakeMap(typData, 2).
		ForPhrase(f).
		DefineVar(key, val).
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
	f := NewForPhrase(key, val, typData)
	c := NewComprehension(reflect.MapOf(TyInt, TyString))
	code := NewBuilder(nil).
		MapComprehension(c).
		Push("Hello").
		Push(3).
		Push("xsw").
		Push(1).
		MakeMap(typData, 2).
		ForPhrase(f).
		DefineVar(key, val).
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
	f := NewForPhrase(nil, x, typData)
	c := NewComprehension(reflect.SliceOf(TyInt))
	code := NewBuilder(nil).
		ListComprehension(c).
		Push(1).
		Push(3).
		Push(5).
		Push(7).
		MakeArray(typData, 4).
		ForPhrase(f).
		DefineVar(x).
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
	f := NewForPhrase(nil, x, typData)
	c := NewComprehension(reflect.SliceOf(TyInt))
	code := NewBuilder(nil).
		ListComprehension(c).
		Push(1).
		Push(3).
		Push(5).
		Push(7).
		MakeArray(typData, 4).
		ForPhrase(f).
		DefineVar(x).
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
	f := NewForPhrase(i, x, typData)
	c := NewComprehension(reflect.MapOf(TyInt, TyInt))
	code := NewBuilder(nil).
		MapComprehension(c).
		Push(1).
		Push(3).
		Push(5).
		Push(7).
		MakeArray(typData, 4).
		ForPhrase(f).
		DefineVar(i, x).
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
	fa := NewForPhrase(nil, a, typData)
	fb := NewForPhrase(nil, b, typData)
	c := NewComprehension(typData)
	code := NewBuilder(nil).
		ListComprehension(c).
		Push(5).
		Push(6).
		Push(7).
		MakeArray(typData, 3).
		ForPhrase(fb).
		DefineVar(b).
		Push(1).
		Push(2).
		Push(3).
		Push(4).
		MakeArray(typData, 4).
		ForPhrase(fa).
		DefineVar(a).
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
	code := NewBuilder(nil).
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

func TestAppend(t *testing.T) {
	code := NewBuilder(nil).
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
	code := NewBuilder(nil).
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

func TestMakeSlice(t *testing.T) {
	code := NewBuilder(nil).
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
	code := NewBuilder(nil).
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
	code := NewBuilder(nil).
		Make(reflect.MapOf(TyInt, TyFloat64), 0).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); !reflect.DeepEqual(v, map[int]float64{}) {
		t.Fatal("ret != {}, ret:", v)
	}
}

func TestMakeMap2(t *testing.T) {
	code := NewBuilder(nil).
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
	code := NewBuilder(nil).
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
	code := NewBuilder(nil).
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
	code := NewBuilder(nil).
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
