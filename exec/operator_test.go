package exec

import (
	"testing"
)

// -----------------------------------------------------------------------------

func TestInt(t *testing.T) {
	code := NewBuilder(nil).
		Push(int(5)).
		Push(int(6)).
		BuiltinOp(Int, OpMul).
		Push(int(36)).
		BuiltinOp(Int, OpMod).
		Push(int(7)).
		BuiltinOp(Int, OpDiv).
		Push(int(1)).
		BuiltinOp(Int, OpAdd).
		Push(int(2)).
		BuiltinOp(Int, OpSub).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != int(3) {
		t.Fatal("5 6 mul 36 mod 7 div 1 add 2 sub != 3, ret =", v)
	}
}

func TestInt64(t *testing.T) {
	code := NewBuilder(nil).
		Push(int64(5)).
		Push(int64(6)).
		BuiltinOp(Int64, OpMul).
		Push(int64(36)).
		BuiltinOp(Int64, OpMod).
		Push(int64(7)).
		BuiltinOp(Int64, OpDiv).
		Push(int64(1)).
		BuiltinOp(Int64, OpAdd).
		Push(int64(2)).
		BuiltinOp(Int64, OpSub).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != int64(3) {
		t.Fatal("5 6 mul 36 mod 7 div 1 add 2 sub != 3, ret =", v)
	}
}

func TestInt32(t *testing.T) {
	code := NewBuilder(nil).
		Push(int32(5)).
		Push(int32(6)).
		BuiltinOp(Int32, OpMul).
		Push(int32(36)).
		BuiltinOp(Int32, OpMod).
		Push(int32(7)).
		BuiltinOp(Int32, OpDiv).
		Push(int32(1)).
		BuiltinOp(Int32, OpAdd).
		Push(int32(2)).
		BuiltinOp(Int32, OpSub).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != int32(3) {
		t.Fatal("5 6 mul 36 mod 7 div 1 add 2 sub != 3, ret =", v)
	}
}

func TestInt16(t *testing.T) {
	code := NewBuilder(nil).
		Push(int16(5)).
		Push(int16(6)).
		BuiltinOp(Int16, OpMul).
		Push(int16(36)).
		BuiltinOp(Int16, OpMod).
		Push(int16(7)).
		BuiltinOp(Int16, OpDiv).
		Push(int16(1)).
		BuiltinOp(Int16, OpAdd).
		Push(int16(2)).
		BuiltinOp(Int16, OpSub).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != int16(3) {
		t.Fatal("5 6 mul 36 mod 7 div 1 add 2 sub != 3, ret =", v)
	}
}

func TestInt8(t *testing.T) {
	code := NewBuilder(nil).
		Push(int8(5)).
		Push(int8(6)).
		BuiltinOp(Int8, OpMul).
		Push(int8(36)).
		BuiltinOp(Int8, OpMod).
		Push(int8(7)).
		BuiltinOp(Int8, OpDiv).
		Push(int8(1)).
		BuiltinOp(Int8, OpAdd).
		Push(int8(2)).
		BuiltinOp(Int8, OpSub).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != int8(3) {
		t.Fatal("5 6 mul 36 mod 7 div 1 add 2 sub != 3, ret =", v)
	}
}

// -----------------------------------------------------------------------------

func TestUint(t *testing.T) {
	code := NewBuilder(nil).
		Push(uint(5)).
		Push(uint(6)).
		BuiltinOp(Uint, OpMul).
		Push(uint(36)).
		BuiltinOp(Uint, OpMod).
		Push(uint(7)).
		BuiltinOp(Uint, OpDiv).
		Push(uint(1)).
		BuiltinOp(Uint, OpAdd).
		Push(uint(2)).
		BuiltinOp(Uint, OpSub).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != uint(3) {
		t.Fatal("5 6 mul 36 mod 7 div 1 add 2 sub != 3, ret =", v)
	}
}

func TestUintptr(t *testing.T) {
	code := NewBuilder(nil).
		Push(uintptr(5)).
		Push(uintptr(6)).
		BuiltinOp(Uintptr, OpMul).
		Push(uintptr(36)).
		BuiltinOp(Uintptr, OpMod).
		Push(uintptr(7)).
		BuiltinOp(Uintptr, OpDiv).
		Push(uintptr(1)).
		BuiltinOp(Uintptr, OpAdd).
		Push(uintptr(2)).
		BuiltinOp(Uintptr, OpSub).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != uintptr(3) {
		t.Fatal("5 6 mul 36 mod 7 div 1 add 2 sub != 3, ret =", v)
	}
}

func TestUint64(t *testing.T) {
	code := NewBuilder(nil).
		Push(uint64(5)).
		Push(uint64(6)).
		BuiltinOp(Uint64, OpMul).
		Push(uint64(36)).
		BuiltinOp(Uint64, OpMod).
		Push(uint64(7)).
		BuiltinOp(Uint64, OpDiv).
		Push(uint64(1)).
		BuiltinOp(Uint64, OpAdd).
		Push(uint64(2)).
		BuiltinOp(Uint64, OpSub).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != uint64(3) {
		t.Fatal("5 6 mul 36 mod 7 div 1 add 2 sub != 3, ret =", v)
	}
}

func TestUint32(t *testing.T) {
	code := NewBuilder(nil).
		Push(uint32(5)).
		Push(uint32(6)).
		BuiltinOp(Uint32, OpMul).
		Push(uint32(36)).
		BuiltinOp(Uint32, OpMod).
		Push(uint32(7)).
		BuiltinOp(Uint32, OpDiv).
		Push(uint32(1)).
		BuiltinOp(Uint32, OpAdd).
		Push(uint32(2)).
		BuiltinOp(Uint32, OpSub).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != uint32(3) {
		t.Fatal("5 6 mul 36 mod 7 div 1 add 2 sub != 3, ret =", v)
	}
}

func TestUint16(t *testing.T) {
	code := NewBuilder(nil).
		Push(uint16(5)).
		Push(uint16(6)).
		BuiltinOp(Uint16, OpMul).
		Push(uint16(36)).
		BuiltinOp(Uint16, OpMod).
		Push(uint16(7)).
		BuiltinOp(Uint16, OpDiv).
		Push(uint16(1)).
		BuiltinOp(Uint16, OpAdd).
		Push(uint16(2)).
		BuiltinOp(Uint16, OpSub).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != uint16(3) {
		t.Fatal("5 6 mul 36 mod 7 div 1 add 2 sub != 3, ret =", v)
	}
}

func TestUint8(t *testing.T) {
	code := NewBuilder(nil).
		Push(uint8(5)).
		Push(uint8(6)).
		BuiltinOp(Uint8, OpMul).
		Push(uint8(36)).
		BuiltinOp(Uint8, OpMod).
		Push(uint8(7)).
		BuiltinOp(Uint8, OpDiv).
		Push(uint8(1)).
		BuiltinOp(Uint8, OpAdd).
		Push(uint8(2)).
		BuiltinOp(Uint8, OpSub).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != uint8(3) {
		t.Fatal("5 6 mul 36 mod 7 div 1 add 2 sub != 3, ret =", v)
	}
}

// -----------------------------------------------------------------------------

func TestFloat64(t *testing.T) {
	code := NewBuilder(nil).
		Push(5.0).
		Push(2.0).
		BuiltinOp(Float64, OpDiv).
		Push(4.0).
		BuiltinOp(Float64, OpMul).
		Push(7.0).
		BuiltinOp(Float64, OpSub).
		Push(1.0).
		BuiltinOp(Float64, OpAdd).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 4.0 {
		t.Fatal("5.0 2.0 div 4.0 mul 7.0 sub 1.0 add != 4.0, ret =", v)
	}
}

func TestFloat32(t *testing.T) {
	code := NewBuilder(nil).
		Push(float32(5.0)).
		Push(float32(2.0)).
		BuiltinOp(Float32, OpAdd).
		Push(float32(1.0)).
		BuiltinOp(Float32, OpSub).
		Push(float32(2.0)).
		BuiltinOp(Float32, OpDiv).
		Push(float32(4.0)).
		BuiltinOp(Float32, OpMul).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != float32(12.0) {
		t.Fatal("5.0 2.0 add 1.0 sub 2.0 div 4.0 mul != 12.0, ret =", v)
	}
}

func TestComplex64(t *testing.T) {
	code := NewBuilder(nil).
		Push(complex64(5.0)).
		Push(complex64(2.0)).
		BuiltinOp(Complex64, OpAdd).
		Push(complex64(1.0)).
		BuiltinOp(Complex64, OpSub).
		Push(complex64(2.0)).
		BuiltinOp(Complex64, OpDiv).
		Push(complex64(4.0)).
		BuiltinOp(Complex64, OpMul).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != complex64(12.0) {
		t.Fatal("5.0 2.0 add 1.0 sub 2.0 div 4.0 mul != 12.0, ret =", v)
	}
}

func TestComplex128(t *testing.T) {
	code := NewBuilder(nil).
		Push(complex128(5.0)).
		Push(complex128(2.0)).
		BuiltinOp(Complex128, OpAdd).
		Push(complex128(1.0)).
		BuiltinOp(Complex128, OpSub).
		Push(complex128(2.0)).
		BuiltinOp(Complex128, OpDiv).
		Push(complex128(4.0)).
		BuiltinOp(Complex128, OpMul).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != complex128(12.0) {
		t.Fatal("5.0 2.0 add 1.0 sub 2.0 div 4.0 mul != 12.0, ret =", v)
	}
}

// -----------------------------------------------------------------------------

func TestStrcat(t *testing.T) {
	code := NewBuilder(nil).
		Push("5").
		Push("6").
		BuiltinOp(String, OpAdd).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "56" {
		t.Fatal("`5` `6` add != `56`, ret =", v)
	}
}

func TestTypeCast(t *testing.T) {
	code := NewBuilder(nil).
		Push(byte('5')).
		TypeCast(Uint8, String).
		Push("6").
		BuiltinOp(String, OpAdd).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "56" {
		t.Fatal("`5` `6` add != `56`, ret =", v)
	}
}

// -----------------------------------------------------------------------------

func TestCallBuiltinOp(t *testing.T) {
	ret := CallBuiltinOp(String, OpAdd, "5", "6")
	if ret.(string) != "56" {
		t.Fatal("CallBuiltinOp failed: ret =", ret)
	}
	ret = CallBuiltinOp(Int, OpAdd, 5, 6)
	if ret.(int) != 11 {
		t.Fatal("CallBuiltinOp failed: ret =", ret)
	}
}

// -----------------------------------------------------------------------------
