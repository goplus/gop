package qlang_test

import (
	"reflect"
	"strings"
	"testing"

	"qlang.io/cl/qlang"
	_ "qlang.io/lib/builtin"
)

// -----------------------------------------------------------------------------

const testReturnCode = `

x = fn {
	return 3
}
`

func TestReturn(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testReturnCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
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

	lang := qlang.New()
	err := lang.SafeExec([]byte(testReturnCode2), "")
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

	lang := qlang.New()
	err := lang.SafeExec([]byte(testAnonymFnCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
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

	lang := qlang.New()
	err := lang.SafeExec([]byte(testAnonymFnCode2), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testDeferCode = `

defer fn { x; x = 2 }
x = 1
`

func TestDefer(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testDeferCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 2 {
		t.Fatal("x != 2, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testDeferCode2 = `

defer fn() { x; x = 2 }()
x = 1
`

func TestDefer2(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testDeferCode2), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 2 {
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

	lang := qlang.New()
	err := lang.SafeExec([]byte(testClosureCode1), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 1 {
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

	lang := qlang.New()
	err := lang.SafeExec([]byte(testClosureCode2), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 1 {
		t.Fatal("x != 1, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testUint32Code = `

a = uint32(3)
x = foo(a)
`

func fooUint32(v uint32) uint32 {
	return v * 2
}

func TestUint32(t *testing.T) {

	lang := qlang.New()
	lang.SetVar("foo", fooUint32)

	err := lang.SafeExec([]byte(testUint32Code), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != uint32(6) {
		t.Fatal("x != 6, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testTypeOfCode = `

x = type(uint32(3))
`

func TestTypeOf(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testTypeOfCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != reflect.TypeOf(uint32(1)) {
		t.Fatal("x != uint32, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testCallbackCode = `

x = strings.indexFunc("Hello!!!", fn(c) { return c == '!' })
`

func TestCallback(t *testing.T) {

	lang := qlang.New()
	qlang.Import("strings", map[string]interface{}{
		"indexFunc": strings.IndexFunc,
	})

	err := lang.SafeExec([]byte(testCallbackCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 5 {
		t.Fatal("x != 5, x =", v)
	}
}

// -----------------------------------------------------------------------------

func forEach(values []int, fn func(v int)) {
	for _, v := range values {
		fn(v)
	}
}

const testCallback2Code = `

x = 0
forEach([1, 2, 2], fn(v) { x; x += v })
`

func TestCallback2(t *testing.T) {

	lang := qlang.New()
	qlang.Import("", map[string]interface{}{
		"forEach": forEach,
	})

	err := lang.SafeExec([]byte(testCallback2Code), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 5 {
		t.Fatal("x != 5, x =", v)
	}
}

// -----------------------------------------------------------------------------

func toVars(values []int, fn func(v int) interface{}) []interface{} {
	ret := make([]interface{}, len(values))
	for i, v := range values {
		ret[i] = fn(v)
	}
	return ret
}

const testCallback3Code = `

ret = toVars([1, 2, 3], fn(v) { return v+1 })
`

func TestCallback3(t *testing.T) {

	lang := qlang.New()
	qlang.Import("", map[string]interface{}{
		"toVars": toVars,
	})

	err := lang.SafeExec([]byte(testCallback3Code), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	ret, ok := lang.GetVar("ret")
	if !ok {
		t.Fatal("ret not found")
	}
	if v, ok := ret.([]interface{}); !ok || len(v) != 3 {
		t.Fatal("ret != [2, 3, 4], ret =", v)
	}
}

// -----------------------------------------------------------------------------
