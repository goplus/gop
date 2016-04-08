package qlang_test

import (
	"testing"

	_ "qlang.io/qlang/builtin"
	"qlang.io/qlang.v2/qlang"
)

// -----------------------------------------------------------------------------

func TestIf(t *testing.T) {

	lang, err := qlang.New(nil)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeEval(`fn { 0 }; if true { 1; 2; 3 } else { 4; 5 }`)
	if err != nil {
		t.Fatal("qlang.SafeEval:", err)
	}
	if v, ok := lang.Ret(); !ok || v != 3 {
		t.Fatal("ret =", v)
	}
}

// -----------------------------------------------------------------------------

