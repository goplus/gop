GoPlus - The Go+ language for data science
========

[![Build Status](https://travis-ci.org/goplus/gop.png?branch=master)](https://travis-ci.org/goplus/gop)
[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/gop)](https://goreportcard.com/report/github.com/goplus/gop)
[![GitHub release](https://img.shields.io/github/v/tag/goplus/gop.svg?label=release)](https://github.com/goplus/gop/releases)
[![Coverage Status](https://codecov.io/gh/goplus/gop/branch/master/graph/badge.svg)](https://codecov.io/gh/goplus/gop)
[![REPL](https://img.shields.io/badge/repl-iGo+-blue.svg)](https://repl.goplus.org/)
[![Playground](https://img.shields.io/badge/playground-Go+-blue.svg)](https://play.goplus.org/)
[![VSCode](https://img.shields.io/badge/vscode-Go+-teal.svg)](https://github.com/gopcode/vscode-goplus)
[![GoDoc](https://img.shields.io/badge/godoc-reference-teal.svg)](https://pkg.go.dev/mod/github.com/goplus/gop)


## Summary about Go+

What are mainly impressions about Go+?

- A static typed language.
- Fully compatible with [the Go language](https://github.com/golang/go).
- Script-like style, and more readable code for data science than Go.

For example, the following is legal Go+ source code:

```go
a := [1, 2, 3.4]
println(a)
```

How do we do this in the Go language?

```go
package main

func main() {
    a := []float64{1, 2, 3.4}
    println(a)
}
```

Of course, we don't only do less-typing things.

For example, we support [list comprehension](https://en.wikipedia.org/wiki/List_comprehension), which makes data processing easier.

```go
a := [1, 3, 5, 7, 11]
b := [x*x for x <- a, x > 3]
println(b) // output: [25 49 121]

mapData := {"Hi": 1, "Hello": 2, "Go+": 3}
reversedMap := {v: k for k, v <- mapData}
println(reversedMap) // output: map[1:Hi 2:Hello 3:Go+]
```

We will keep Go+ simple. This is why we call it Go+, not Go++.

Less is exponentially more.

It's for Go, and it's also for Go+.


## Compatibility with Go

All Go features will be supported (including partially support `cgo`, see [below](#bytecode-vs-go-code)).

* See [supported the Go language features](https://github.com/goplus/gop/wiki/Supported-Go-features).

**All Go packages (even these packages use `cgo`) can be imported by Go+.**

```go
import (
    "fmt"
    "strings"
)

x := strings.NewReplacer("?", "!").Replace("hello, world???")
fmt.Println("x:", x)
```

Be interested in how it works? See [Import Go packages in Go+ programs](https://github.com/goplus/gop/wiki/Import-Go-packages-in-GoPlus-programs).

**Also, all Go+ packages can be converted into Go packages and imported in Go programs.**

First, let's make a directory named `tutorial/14-Using-goplus-in-Go`.

Then write a Go+ package named `foo` in it:

```go
package foo

func ReverseMap(m map[string]int) map[int]string {
    return {v: k for k, v <- m}
}
```

Then use it in a Go package:

```go
package main

import (
	"fmt"

	"github.com/goplus/gop/tutorial/14-Using-goplus-in-Go/foo"
)

func main() {
	rmap := foo.ReverseMap(map[string]int{"Hi": 1, "Hello": 2})
	fmt.Println(rmap)
}
```

How to compile this example?

```bash
gop go tutorial/ # Convert all Go+ packages in tutorial/ into Go packages
go install ./...
```

Or:

```bash
gop install ./... # Convert Go+ packages and go install ./...
```

Go [tutorial/14-Using-goplus-in-Go](https://github.com/goplus/gop/tree/master/tutorial/14-Using-goplus-in-Go) to get the source code.


## Playground/REPL

Go+ REPL based on GopherJS/WASM:
* https://repl.goplus.org/

Go+ Playground based on Docker:
* https://play.goplus.org/

Go+ Playground based on GopherJS:
* https://goplusjs.github.io/play/


## Tutorials

See https://github.com/goplus/gop/tree/master/tutorial


## How to build

```bash
git clone git@github.com:goplus/gop.git
cd gop
go install -v ./...
```

## Go+ features

### Bytecode vs. Go code

Go+ supports bytecode backend and Go code generation.

When we use `gop go` or `gop install` command, it generates Go code to covert Go+ package into Go packages.

When we use `gop run` command, it doesn't call `go run` command. It generates bytecode to execute.

In bytecode mode, Go+ doesn't support `cgo`. However, in Go-code-generation mode, Go+ fully supports `cgo`.


### Commands

```bash
gop run         # Run a Go+ program
gop repl        # Run Go+ in REPL/Console mode
gop go [-test]  # Convert Go+ packages into Go packages. If -test specified, it tests related packages.
gop fmt         # Format Go+ packages
gop export      # Export Go packages for Go+ programs
```

See https://github.com/goplus/gop/wiki/Commands for details.

Note:

* `gop go -test <gopSrcDir>` converts Go+ packages into Go packages, and for every package, it call `go run <gopPkgDir>/gop_autogen.go` and `gop run -quiet <gopPkgDir>` to compare their outputs. If their outputs aren't equal, the test case fails.


### Rational number: bigint, bigrat, bigfloat

We introduce the rational number as native Go+ types. We use suffix `r` to denote rational literals. For example, (1r << 200) means a big int whose value is equal to 2<sup>200</sup>. And 4/5r means the rational constant 4/5.

```go
a := 1r << 65   // bigint, large than int64
b := 4/5r       // bigrat
c := b - 1/3r + 3 * 1/2r
println(a, b, c)
```

### Map literal

```go
x := {"Hello": 1, "xsw": 3.4} // map[string]float64
y := {"Hello": 1, "xsw": "Go+"} // map[string]interface{}
z := {"Hello": 1, "xsw": 3} // map[string]int
empty := {} // map[string]interface{}
```

### Slice literal

```go
x := [1, 3.4] // []float64
y := [1] // []int
z := [1+2i, "xsw"] // []interface{}
a := [1, 3.4, 3+4i] // []complex128
b := [5+6i] // []complex128
c := ["xsw", 3] // []interface{}
empty := [] // []interface{}
```

### List/Map comprehension

```go
a := [x*x for x <- [1, 3, 5, 7, 11]]
b := [x*x for x <- [1, 3, 5, 7, 11], x > 3]
c := [i+v for i, v <- [1, 3, 5, 7, 11], i%2 == 1]
d := [k+","+s for k, s <- {"Hello": "xsw", "Hi": "Go+"}]

arr := [1, 2, 3, 4, 5, 6]
e := [[a, b] for a <- arr, a < b for b <- arr, b > 2]

x := {x: i for i, x <- [1, 3, 5, 7, 11]}
y := {x: i for i, x <- [1, 3, 5, 7, 11], i%2 == 1}
z := {v: k for k, v <- {1: "Hello", 3: "Hi", 5: "xsw", 7: "Go+"}, k > 3}
```

### For loop

```go
sum := 0
for x <- [1, 3, 5, 7, 11, 13, 17], x > 3 {
    sum += x
}
```

### Error handling

We reinvent the error handling specification in Go+. We call them `ErrWrap expressions`:

```go
expr! // panic if err
expr? // return if err
expr?:defval // use defval if err
```

How to use them? Here is an example:

```go
import (
	"strconv"
)

func add(x, y string) (int, error) {
	return strconv.Atoi(x)? + strconv.Atoi(y)?, nil
}

func addSafe(x, y string) int {
	return strconv.Atoi(x)?:0 + strconv.Atoi(y)?:0
}

println(`add("100", "23"):`, add("100", "23")!)

sum, err := add("10", "abc")
println(`add("10", "abc"):`, sum, err)

println(`addSafe("10", "abc"):`, addSafe("10", "abc"))
```

The output of this example is:

```
add("100", "23"): 123
add("10", "abc"): 0 strconv.Atoi: parsing "abc": invalid syntax

===> errors stack:
main.add("10", "abc")
	/Users/xsw/goplus/tutorial/15-ErrWrap/err_wrap.gop:6 strconv.Atoi(y)?

addSafe("10", "abc"): 10
```

Compared to corresponding Go code, It is clear and more readable.

And the most interesting thing is, the return error contains the full error stack. When we got an error, it is very easy to position what the root cause is.

How these `ErrWrap expressions` work? See [Error Handling](https://github.com/goplus/gop/wiki/Error-Handling) for more information.


### Auto property

Let's see an example written in Go+:

```go
import "github.com/goplus/gop/ast/goptest"

doc := goptest.New(`... Go+ code ...`)!

println(doc.Any().FuncDecl().Name())
```

In many languages, there is a concept named `property` who has `get` and `set` methods.

Suppose we have `get property`, the above example will be:

```go
import "github.com/goplus/gop/ast/goptest"

doc := goptest.New(`... Go+ code ...`)!

println(doc.any.funcDecl.name)
```

In Go+, we introduce a concept named `auto property`. It is a `get property`, but is implemented automatically. If we have a method named `Bar()`, then we will have a `get property` named `bar` at the same time.


### Unix shebang

You can use Go+ programs as shell scripts now. For example:

```go
#!/usr/bin/env qrun

println("Hello, Go+")

println(1r << 129)
println(1/3r + 2/7r*2)

arr := [1, 3, 5, 7, 11, 13, 17, 19]
println(arr)
println([x*x for x <- arr, x > 3])

m := {"Hi": 1, "Go+": 2}
println(m)
println({v: k for k, v <- m})
println([k for k, _ <- m])
println([v for v <- m])
```

Go [tutorial/20-Unix-Shebang/shebang](https://github.com/goplus/gop/blob/master/tutorial/20-Unix-Shebang/shebang) to get the source code.


### Go features

All Go features (not including `cgo`) will be supported.

* See [supported the Go language features](https://github.com/goplus/gop/wiki/Supported-Go-features).


## IDE Plugins

* vscode: https://github.com/gopcode/vscode-goplus


## Contributing

The Go+ project welcomes all contributors. We appreciate your help!

Here are [list of Go+ Contributors](https://github.com/goplus/gop/wiki/Goplus-Contributors). We award an email account (XXX@goplus.org) for every contributor. And we suggest you commit code by using this email account:

```bash
git config --global user.email XXX@goplus.org
```

If you did this, remember to add your `XXX@goplus.org` email to https://github.com/settings/emails.

What does `a contributor to Go+` mean? He must meet one of the following conditions:

* At least one pull request of a full-featured implemention.
* At least three pull requests of feature enhancements.
* At least ten pull requests of any kind issues.

Where can you start?

* [![Issues](https://img.shields.io/badge/ISSUEs-Go+-blue.svg)](https://github.com/goplus/gop/issues)
* [![Issues](https://img.shields.io/badge/ISSUEs-NumGo+-blue.svg)](https://github.com/numgoplus/ng/issues)
* [![Issues](https://img.shields.io/badge/ISSUEs-PandasGo+-blue.svg)](https://github.com/goplus/pandas/issues)
* [![Issues](https://img.shields.io/badge/ISSUEs-vscode%20Go+-blue.svg)](https://github.com/gopcode/vscode-goplus/issues)
* [![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/goplus/gop)](https://www.tickgit.com/browse?repo=github.com/goplus/gop)
