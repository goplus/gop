package qlang_test

import (
	"reflect"
	"testing"

	"qlang.io/qlang.v2/qlang"
	"qlang.io/qlang/math"
)

// -----------------------------------------------------------------------------

const testCastFloatCode = `

x = sin(0)
`

func TestCastFloat(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testCastFloatCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 0.0 {
		t.Fatal("x != 0.0, x =", v)
	}
}

func init() {
	qlang.Import("", math.Exports)
}

// -----------------------------------------------------------------------------

const testTestByteSliceCode = `

a = ['a', 'b']
x = type(a)
`

func TestByteSlice(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testTestByteSliceCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != reflect.TypeOf([]byte(nil)) {
		t.Fatal("x != []byte, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testNilSliceCode = `

a = []
a = append(a, 1)
x = type(a)
y = a[0]
`

func TestNilSlice(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testNilSliceCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != reflect.TypeOf([]interface{}(nil)) {
		t.Fatal("x != []interface{}, x =", v)
	}
	if v, ok := lang.Var("y"); !ok || v != 1 {
		t.Fatal("y != 1, y =", v)
	}
}

// -----------------------------------------------------------------------------

type fooMode uint

func castFooMode(mode fooMode) int {
	return int(mode)
}

const testCastFooModeCode = `

y = castFooMode(1)
`

func TestCastFooMode(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}
	qlang.Import("", map[string]interface{}{
		"castFooMode": castFooMode,
	})

	err = lang.SafeExec([]byte(testCastFooModeCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("y"); !ok || v != 1 {
		t.Fatal("y != 1, y =", v)
	}
}

// -----------------------------------------------------------------------------
