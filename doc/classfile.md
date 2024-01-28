Go+ Classfiles
=====

Rob Pike once said that if he could only introduce one feature to Go, he would choose to introduce `interface` instead of `goroutine`. `classfile` is as important to Go+ as `interface` is to Go.

In the design philosophy of Go+, we do not recommend `DSL` (Domain Specific Language). But `SDF` (Specific Domain Friendliness) is very important. The Go+ philosophy about `SDF` is:

```
Don't define a language for specific domain.
Abstract domain knowledge for it.
```

Go+ introduces `classfile` to abstract domain knowledge.


### classfile: Unit Test

Go+ has built-in support for a classfile to simplify unit testing. This classfile has the file suffix `_test.gox`.

Suppose you have a function named `foo`:

```go
func foo(v int) int {
	return v * 2
}
```

Then you can create a `foo_test.gox` file to test it (see [unit-test/foo_test.gox](testdata/unit-test/foo_test.gox)):

```go
if v := foo(50); v == 100 {
	t.log "test foo(50) ok"
} else {
	t.error "foo() ret: ${v}"
}

t.run "foo -10", t => {
	if foo(-10) != -20 {
		t.fatal "foo(-10) != -20"
	}
}

t.run "foo 0", t => {
	if foo(0) != 0 {
		t.fatal "foo(0) != 0"
	}
}
```

You don't need to define a series of `TestXXX` functions like Go, just write your test code directly.

If you want to run a subtest case, use `t.run`.


### classfile: HTTP Web Framework

