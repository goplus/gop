package main

import (
	"testing"
)

const fibCode = `
fib = fn() {
	a, b = [0, 1]
	f = fn() {
		a, b = [b, a+b]
		return a
	}
	return f
}
f = fib()
println(f(), f(), f(), f(), f())
`

func _TestCompileAndRun(t *testing.T) {
	req := Request{
		Body: fibCode,
	}
	res, err := compileAndRun(&req)
	if err != nil {
		t.Fatal("compileAndRun test error")
	}
	if res.Events[0].Message != "1 1 2 3 5\n" {
		t.Fatal("compileAndRun res test error")
	}
}
