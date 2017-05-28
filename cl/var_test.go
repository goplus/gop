package qlang_test

import (
	"testing"

	"qlang.io/cl/qlang"
	_ "qlang.io/lib/builtin"
	ql "qlang.io/spec"
)

// -----------------------------------------------------------------------------

const testGetCode = `

a = [1, 2, 3]
x = a[2]
b = "123"
y = b[2]
c = {"a": 1, "b": 2, "c": 3}
z = c.a
`

func TestGet(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testGetCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
	if v, ok := lang.GetVar("y"); !ok || v != byte('3') {
		t.Fatal("y != '3', y =", v)
	}
	if v, ok := lang.GetVar("z"); !ok || v != 1 {
		t.Fatal("z != 1, z =", v)
	}
}

// -----------------------------------------------------------------------------

const testGetUndefinedCode = `

c = {}
x = c["a"]["b"]
`

func TestGetUndefined(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testGetUndefinedCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v := lang.Var("x"); v != ql.Undefined {
		t.Fatal("x != undefined, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testIncCode = `

a = [1, 100]
a[0]++
x = a[0]
x++
`

func TestInc(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testIncCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testMulAssignCode = `

t = {"a": 2, "b": 2}
t["a"] *= 2
t.a *= 2
x = t.a
x *= 2
`

func TestMulAssign(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testMulAssignCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 16 {
		t.Fatal("x != 16, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testSwapCode = `

t = {"a": 2, "b": 3}
t.a, t.b = t.b, t.a
x, y = t.a, t.b
`

func TestSwap(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testSwapCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
	if v, ok := lang.GetVar("y"); !ok || v != 2 {
		t.Fatal("y != 2, y =", v)
	}
}

// -----------------------------------------------------------------------------

const testMultiAssignCode = `

x, y = 3, 2
`

func TestMultiAssign(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testMultiAssignCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
	if v, ok := lang.GetVar("y"); !ok || v != 2 {
		t.Fatal("y != 2, y =", v)
	}
}

// -----------------------------------------------------------------------------

const testMultiAssignCode2 = `

x, y = [3, 2]
`

func TestMultiAssign2(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testMultiAssignCode2), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
	if v, ok := lang.GetVar("y"); !ok || v != 2 {
		t.Fatal("y != 2, y =", v)
	}
}

// -----------------------------------------------------------------------------

const testErrMultiAssignCode = `

x, y = 3, 2, 1
`

func TestErrMultiAssign(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testErrMultiAssignCode), "")
	if err == nil {
		t.Fatal("qlang.SafeExec?")
	}
	if err.Error() != "line 3: runtime error: argument count of multi assignment doesn't match" {
		t.Fatal("testErrMultiAssignCode:", err)
	}
}

// -----------------------------------------------------------------------------

const testMultiAssignExCode = `

a = [1, 2]
b = {"a": 1, "b": "hello"}
a[1], b.b = 3, 2
x, y = a[1], b.b
`

func TestMultiAssignEx(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testMultiAssignExCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
	if v, ok := lang.GetVar("y"); !ok || v != 2 {
		t.Fatal("y != 2, y =", v)
	}
}

// -----------------------------------------------------------------------------

const testAssignExCode = `

Foo = class {
	fn _init(v) {
		this.foo = v
	}

	fn f() {
		return this.foo
	}
}

a = new Foo(3)
x = a.f()
`

func TestAssignEx(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testAssignExCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testClassAddAssignCode = `

Foo = class {
	fn _init() {
		this.foo = 1
	}

	fn f() {
		this.foo += 2
	}
}

a = new Foo
a.f()
x = a.foo
`

func TestClassAddAssign(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testClassAddAssignCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
}

// -----------------------------------------------------------------------------
