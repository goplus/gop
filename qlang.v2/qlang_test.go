package qlang_test

import (
	"reflect"
	"testing"

	"qlang.io/qlang.v2/qlang"
)

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
