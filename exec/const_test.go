package exec

import (
	"testing"
)

// -----------------------------------------------------------------------------

func TestConst1(t *testing.T) {

	code := NewBuilder(nil).
		Push(1 << 32).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 1<<32 {
		t.Fatal("1<<32 != 1<<32, ret =", v)
	}
}

func TestConst2(t *testing.T) {

	code := NewBuilder(nil).
		Push(uint64(1 << 32)).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != uint64(1<<32) {
		t.Fatal("1<<32 != 1<<32, ret =", v)
	}
}

func TestConst3(t *testing.T) {

	code := NewBuilder(nil).
		Push(uint32(1 << 30)).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != uint32(1<<30) {
		t.Fatal("1<<30 != 1<<30, ret =", v)
	}
}

func TestConst4(t *testing.T) {

	code := NewBuilder(nil).
		Push(int32(1 << 30)).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != int32(1<<30) {
		t.Fatal("1<<30 != 1<<30, ret =", v)
	}
}

func TestConst5(t *testing.T) {

	code := NewBuilder(nil).
		Push(uint(1 << 12)).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != uint(1<<12) {
		t.Fatal("1<<12 != 1<<12, ret =", v)
	}
}

// -----------------------------------------------------------------------------
