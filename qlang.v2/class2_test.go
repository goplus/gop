package qlang_test

import (
	"testing"

	"qlang.io/qlang.v2/qlang"
)

// -----------------------------------------------------------------------------

const testInjectCode = `

y = 3
typ = class {
	fn foo() {
		return y
	}
}

a = new typ
x = a.foo()
`

const testInjectMethodsCode = `

fn bar(v) {
    return y + v
}

fn setbar(v) {
    y; y = v
}
`

const testInjectCheckCode = `

x2 = a.bar(1)

a.setbar(10)
x3 = a.bar(1)
`

func TestInject(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testInjectCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}

	err = lang.InjectMethods("typ", []byte(testInjectMethodsCode))
	if err != nil {
		t.Fatal("InjectMethods failed:", err)
	}

	err = lang.SafeEval(testInjectCheckCode)
	if err != nil {
		t.Fatal("SafeEval failed:", err)
	}
	if v, ok := lang.Var("x2"); !ok || v != 4 {
		t.Fatal("x2 != 4, x2 =", v)
	}
	if v, ok := lang.Var("x3"); !ok || v != 11 {
		t.Fatal("x3 != 11, x3 =", v)
	}
}

// -----------------------------------------------------------------------------
