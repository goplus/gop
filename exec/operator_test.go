package exec

import (
	"testing"
)

// -----------------------------------------------------------------------------

func TestMul(t *testing.T) {

	code := NewBuilder(nil).
		Push(5).
		Push(6).
		BuiltinOp(Int, OpMul).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v, ok := ctx.Pop(); !(ok && v == 30) {
		t.Fatal("5 6 mul != 30, ret =", v)
	}
}

func TestStrcat(t *testing.T) {

	code := NewBuilder(nil).
		Push("5").
		Push("6").
		BuiltinOp(String, OpAdd).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v, ok := ctx.Pop(); !(ok && v == "56") {
		t.Fatal("`5` `6` add != `56`, ret =", v)
	}
}

// -----------------------------------------------------------------------------
