package qlang_test

import (
	"testing"

	_ "qlang.io/qlang/builtin"
	"qlang.io/qlang.v2/qlang"
)

// -----------------------------------------------------------------------------

func TestEval(t *testing.T) {

	lang, err := qlang.New(nil)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeEval(`1 + 2`)
	if err != nil {
		t.Fatal("qlang.SafeEval:", err)
	}
	if v, ok := lang.Ret(); !ok || v != 3 {
		t.Fatal("ret =", v)
	}
}

// -----------------------------------------------------------------------------

