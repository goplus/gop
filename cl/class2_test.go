package qlang_test

import (
	"testing"

	"qlang.io/cl/qlang"
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

	lang := qlang.New()
	err := lang.SafeExec([]byte(testInjectCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
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
	if v, ok := lang.GetVar("x2"); !ok || v != 4 {
		t.Fatal("x2 != 4, x2 =", v)
	}
	if v, ok := lang.GetVar("x3"); !ok || v != 11 {
		t.Fatal("x3 != 11, x3 =", v)
	}
}

// -----------------------------------------------------------------------------

type FooMemberRef struct {
	X int
}

func (p *FooMemberRef) Bar() int {
	return p.X
}

func (p FooMemberRef) Bar2() int {
	return p.X
}

func fooMemberPtr() *FooMemberRef {
	return &FooMemberRef{3}
}

func fooCall(v FooMemberRef) int {
	return v.X
}

const testMemberRefCode = `

foo = fooMemberPtr()
x = foo.x
y = foo.bar2()
z = foo.bar()
ret = fooCall(foo)
`

func TestMemberRef(t *testing.T) {

	lang := qlang.New()
	qlang.Import("", map[string]interface{}{
		"fooMemberPtr": fooMemberPtr,
		"fooCall":      fooCall,
	})

	err := lang.SafeExec([]byte(testMemberRefCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
	if v, ok := lang.GetVar("y"); !ok || v != 3 {
		t.Fatal("y != 3, y =", v)
	}
	if v, ok := lang.GetVar("z"); !ok || v != 3 {
		t.Fatal("z != 3, z =", v)
	}
	if v, ok := lang.GetVar("ret"); !ok || v != 3 {
		t.Fatal("ret != 3, ret =", v)
	}
}

// -----------------------------------------------------------------------------

const testNewClassCode = `

t = new class {
	fn f() {
		return 2
	}
}
x = t.f()
`

func TestNewClass(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testNewClassCode), "")
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 2 {
		t.Fatal("x != 2, x =", v)
	}
}

// -----------------------------------------------------------------------------
