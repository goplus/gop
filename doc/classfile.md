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


### yap: Yet Another Go/Go+ HTTP Web Framework

This classfile has the file suffix `_yap.gox`.

Before using `yap`, you need to add it to `go.mod` by using `go get`:

```sh
go get github.com/goplus/yap@latest
```

Find `require github.com/goplus/yap` statement in `go.mod` and add `//gop:class` at the end of the line:

```go.mod
require github.com/goplus/yap v0.7.1 //gop:class
```

#### Router and Parameters

demo in Go+ classfile ([hello_yap.gox](https://github.com/goplus/yap/blob/v0.7.2/demo/classfile_hello/hello_yap.gox)):

```go
get "/p/:id", ctx => {
	ctx.json {
		"id": ctx.param("id"),
	}
}
handle "/", ctx => {
	ctx.html `<html><body>Hello, <a href="/p/123">Yap</a>!</body></html>`
}

run ":8080"
```

#### Static files

Static files server demo in Go+ classfile ([staticfile_yap.gox](https://github.com/goplus/yap/blob/v0.7.2/demo/classfile_static/staticfile_yap.gox)):

```go
static "/foo", FS("public")
static "/"

run ":8080"
```

#### YAP Template

demo in Go+ classfile ([blog_yap.gox](https://github.com/goplus/yap/blob/v0.7.2/demo/classfile_blog/blog_yap.gox), [article_yap.html](https://github.com/goplus/yap/blob/v0.7.2/demo/classfile_blog/yap/article_yap.html)):

```go
get "/p/:id", ctx => {
	ctx.yap "article", {
		"id": ctx.param("id"),
	}
}

run ":8080"
```

### yaptest: HTTP Unit Test Framework

This classfile has the file suffix `_ytest.gox`.

```go
host "https://example.com", "http://localhost:8080"
testauth := oauth2("...")

run "urlWithVar", => {
	id := "123"
	get "https://example.com/p/${id}"
	ret
	echo "code:", resp.code
	echo "body:", resp.body
}

run "matchWithVar", => {
	code := Var(int)
	id := "123"
	get "https://example.com/p/${id}"
	ret code
	echo "code:", code
	match code, 200
}

run "postWithAuth", => {
	id := "123"
	title := "title"
	author := "author"
	post "https://example.com/p/${id}"
	auth testauth
	json {
		"title":  title,
		"author": author,
	}
	ret 200 // match resp.code, 200
	echo "body:", resp.body
}

run "mathJsonObject", => {
	title := Var(string)
	author := Var(string)
	id := "123"
	get "https://example.com/p/${id}"
	ret 200
	json {
		"title":  title,
		"author": author,
	}
	echo "title:", title
	echo "author:", author
}
```

### spx: A Go+ 2D Game Engine for STEM education

This classfile has the file suffix `.spx`.
