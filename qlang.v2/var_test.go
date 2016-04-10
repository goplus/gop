package qlang_test

import (
	"testing"

	_ "qlang.io/qlang/builtin"
	"qlang.io/qlang.v2/qlang"
)

// -----------------------------------------------------------------------------

const testMultiAssignCode = `

x, y = 3, 2
`

func TestMultiAssign(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testMultiAssignCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
	if v, ok := lang.Var("y"); !ok || v != 2 {
		t.Fatal("y != 2, y =", v)
	}
}

// -----------------------------------------------------------------------------

const testMultiAssignCode2 = `

x, y = [3, 2]
`

func TestMultiAssign2(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testMultiAssignCode2), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
	if v, ok := lang.Var("y"); !ok || v != 2 {
		t.Fatal("y != 2, y =", v)
	}
}

// -----------------------------------------------------------------------------

const testErrMultiAssignCode = `

x, y = 3, 2, 1
`

func TestErrMultiAssign(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testErrMultiAssignCode), "")
	if err == nil {
		t.Fatal("qlang.SafeExec?")
	}
	if err.Error() != "line 3: runtime error: argument count of multi assignment doesn't match" {
		t.Fatal("testErrMultiAssignCode:", err)
	}
}

// -----------------------------------------------------------------------------

