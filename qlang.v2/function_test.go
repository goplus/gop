package qlang_test

import (
	"testing"

	_ "qlang.io/qlang/builtin"
	"qlang.io/qlang.v2/qlang"
)

// -----------------------------------------------------------------------------

const testReturnCode = `

x = fn {
	return 3
}
`

func TestReturn(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testReturnCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testReturnCode2 = `

x, y = fn {
	return 3, 2
}
`

func TestReturn2(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testReturnCode2), "")
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

const testAnonymFnCode = `

x = 1

main {
	x; x = 2
	println("x:", x)
}

x = 3
println("x:", x)
`

func TestAnonymFn(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testAnonymFnCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testAnonymFnCode2 = `

x = 1
fn { x; x = 2 }
x = 3
`

func TestAnonymFn2(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testAnonymFnCode2), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testDeferCode = `

defer fn { x; x = 2 }
x = 1
`

func TestDefer(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testDeferCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 2 {
		t.Fatal("x != 2, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testDeferCode2 = `

defer fn() { x; x = 2 }()
x = 1
`

func TestDefer2(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testDeferCode2), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 2 {
		t.Fatal("x != 2, x =", v)
	}
}

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
}

// -----------------------------------------------------------------------------

