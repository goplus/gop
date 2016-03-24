package qlang_test

import (
	"testing"

	_ "qlang.io/qlang/builtin"
	"qlang.io/qlang.v2/qlang"
)

// -----------------------------------------------------------------------------

const testClosureCode1 = `

test = fn() {
	g = 1
	return fn() {
		return g
	}
}

x = 0

fn() {
	f = test()
	g = 2
	x; x = f()
	println("f:", x)
}()
`

func TestClosure1(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testClosureCode1), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 1 {
		t.Fatal("x != 1, x =", v)
	}
	return
}

// -----------------------------------------------------------------------------

const testClosureCode2 = `

test = fn() {
	g = 1
	return class {
		fn f() {
			return g
		}
	}
}

x = 0

fn() {
	Foo = test()
	foo = new Foo
	g = 2
	x; x = foo.f()
	println("foo.f:", x)
}()
`

func TestClosure2(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testClosureCode2), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 1 {
		t.Fatal("x != 1, x =", v)
	}
	return
}

// -----------------------------------------------------------------------------

