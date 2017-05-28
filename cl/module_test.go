package qlang_test

import (
	"testing"

	"qlang.io/cl/qlang"
	_ "qlang.io/lib/builtin"
)

// -----------------------------------------------------------------------------

const scriptA = `

println("in script A")

a = 1
b = 2

foo = fn(a) {
	println("in func foo:", a, b)
}

export b, foo
`

const scriptB = `

include "a.ql"

println("in script B")
foo(a)
`

const scriptC = `

import "a"
import "a" as g

b = 3

g.b = 4
println("in script C:", a.b, g.b)
a.foo(b)
`

func TestInclude(t *testing.T) {

	lang := qlang.New()
	qlang.SetReadFile(func(file string) ([]byte, error) {
		return []byte(scriptA), nil
	})

	err := lang.SafeExec([]byte(scriptB), "b.ql")
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := lang.GetVar("b"); !ok || v != 2 {
		t.Fatal("b != 2, b =", v)
	}
}

func TestImport(t *testing.T) {

	lang := qlang.New()
	qlang.SetFindEntry(func(file string, libs []string) (string, error) {
		return file, nil
	})

	qlang.SetReadFile(func(file string) ([]byte, error) {
		return []byte(scriptA), nil
	})

	err := lang.SafeExec([]byte(scriptC), "c.ql")
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := lang.GetVar("a"); !ok || v.(map[string]interface{})["b"] != 4 {
		t.Fatal("a.b != 4, a.b =", v)
	}
}

// -----------------------------------------------------------------------------

const scriptA1 = `

defer fn() {
	x; x = 2
}()

x = 1
export x
`

const scriptB1 = `

import "a"

println("a.x:", a.x)
`

func TestModuleDefer(t *testing.T) {

	lang := qlang.New()
	qlang.SetFindEntry(func(file string, libs []string) (string, error) {
		return file, nil
	})

	qlang.SetReadFile(func(file string) ([]byte, error) {
		return []byte(scriptA1), nil
	})

	err := lang.SafeExec([]byte(scriptB1), "b.ql")
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := lang.GetVar("a"); !ok || v.(map[string]interface{})["x"] != 2 {
		t.Fatal("x != 2, x =", v)
	}
}

// -----------------------------------------------------------------------------

const scriptA2 = `

Foo = class {
	fn f() {
		return 2
	}
}

export Foo
`

const scriptB2 = `

import "a"

t = new a.Foo
x = t.f()
`

func TestModuleClass(t *testing.T) {

	lang := qlang.New()
	qlang.SetFindEntry(func(file string, libs []string) (string, error) {
		return file, nil
	})

	qlang.SetReadFile(func(file string) ([]byte, error) {
		return []byte(scriptA2), nil
	})

	err := lang.SafeExec([]byte(scriptB2), "b.ql")
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 2 {
		t.Fatal("x != 2, x =", v)
	}
}

// -----------------------------------------------------------------------------

const scriptPi = `pi = 2`

const scriptIncPi = `

include "pi.ql"
b = pi
`

func TestIncludePi(t *testing.T) {

	lang := qlang.New()
	qlang.SetReadFile(func(file string) ([]byte, error) {
		return []byte(scriptPi), nil
	})

	err := lang.Exec([]byte(scriptIncPi), "incPi.ql")
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := lang.GetVar("b"); !ok || v != 2 {
		t.Fatal("b != 2, b =", v)
	}
}

// -----------------------------------------------------------------------------
