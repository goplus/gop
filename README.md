## The Go+ language for engineering, STEM education, and data science

[![Build Status](https://github.com/goplus/gop/actions/workflows/go.yml/badge.svg)](https://github.com/goplus/gop/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/gop)](https://goreportcard.com/report/github.com/goplus/gop)
[![Coverage Status](https://codecov.io/gh/goplus/gop/branch/main/graph/badge.svg)](https://codecov.io/gh/goplus/gop)
[![GitHub release](https://img.shields.io/github/v/tag/goplus/gop.svg?label=release)](https://github.com/goplus/gop/releases)
[![Tutorials](https://img.shields.io/badge/tutorial-Go+-blue.svg)](https://tutorial.goplus.org/)
[![Playground](https://img.shields.io/badge/playground-Go+-blue.svg)](https://play.goplus.org/)
[![VSCode](https://img.shields.io/badge/vscode-Go+-teal.svg)](https://github.com/gopcode/vscode-goplus)
[![Readme](https://img.shields.io/badge/README-中文-teal.svg)](https://github.com/goplus/gop/blob/main/README_zh.md)

* **For engineering**: working in the simplest language that can be mastered by children.
* **For STEM education**: studying an engineering language that can be used to work in the future.
* **For data science**: communicating with engineers in the same language.

## How to install

For now, we suggest you install Go+ from source code.

```bash
git clone https://github.com/goplus/gop.git
cd gop

# On mac/linux run:
./all.bash
# On Windows run:
all.bat
```

Actually, `all.bash` and `all.bat` will use `go run cmd/make.go` underneath.


## Code style (important)

* https://tutorial.goplus.org/hello-world


## Summary about Go+

What are mainly impressions about Go+?

- A static typed language.
- Fully compatible with [the Go language](https://github.com/golang/go).
- Script-like style, and more readable code than Go.

For example, the following is legal Go+ source code:

```coffee
println [1, 2, 3.4]
```

How do we do this in the Go language?

```go
package main

import "fmt"

func main() {
    fmt.Println([]float64{1, 2, 3.4})
}
```

Of course, we don't only do less-typing things.

For example, we support [list comprehension](https://en.wikipedia.org/wiki/List_comprehension), which makes data processing easier.

```go
println [x*x for x <- 1:6:2] // output: [1 9 25]

mapData := {"Hi": 1, "Hello": 2, "Go+": 3}
reversedMap := {v: k for k, v <- mapData}
println reversedMap // output: map[1:Hi 2:Hello 3:Go+]
```

We will keep Go+ simple. This is why we call it Go+, not Go++.

Less is exponentially more.

It's for Go, and it's also for Go+.


## Tutorials

* https://tutorial.goplus.org/


## Playground

Go+ Playground based on Docker (recommended):
* https://play.goplus.org/

Go+ Playground based on GopherJS (currently only available in v0.7.x):
* https://jsplay.goplus.org/


## Compatibility with Go

All Go features will be supported (including partially support `cgo`, see [below](#bytecode-vs-go-code)).

**All Go packages (even these packages use `cgo`) can be imported by Go+.**

```coffee
import (
    "fmt"
    "strings"
)

x := strings.NewReplacer("?", "!").Replace("hello, world???")
fmt.Println "x:", x
```

**And all Go+ packages can also be imported in Go programs. What you need to do is just using `gop` command instead of `go`.**

First, let's make a directory named `14-Using-goplus-in-Go`.

Then write a Go+ package named [foo](https://github.com/goplus/tutorial/tree/main/14-Using-goplus-in-Go/foo) in it:

```go
package foo

func ReverseMap(m map[string]int) map[int]string {
    return {v: k for k, v <- m}
}
```

Then use it in a Go package [14-Using-goplus-in-Go/gomain](https://github.com/goplus/tutorial/tree/main/14-Using-goplus-in-Go/gomain):

```go
package main

import (
    "fmt"

    "github.com/goplus/tutorial/14-Using-goplus-in-Go/foo"
)

func main() {
    rmap := foo.ReverseMap(map[string]int{"Hi": 1, "Hello": 2})
    fmt.Println(rmap)
}
```

How to build this example? You can use:

```bash
gop install -v ./...
```

Go [github.com/goplus/tutorial/14-Using-goplus-in-Go](https://github.com/goplus/tutorial/tree/main/14-Using-goplus-in-Go) to get the source code.


## Bytecode vs. Go code

Go+ supports bytecode backend and Go code generation.

When we use `gop` command, it generates Go code to covert Go+ package into Go packages.

```bash
gop run     # Run a Go+ program
gop install # Build Go+ files and install target to GOBIN
gop build   # Build Go+ files
gop test    # Test Go+ packages
gop fmt     # Format Go+ packages
gop clean   # Clean all Go+ auto generated files
gop go      # Convert Go+ packages into Go packages
```

When we use [`igop`](https://github.com/goplus/igop) command, it generates bytecode to execute.

```bash
igop  # Run a Go+ program
```

In bytecode mode, Go+ doesn't support `cgo`. However, in Go-code-generation mode, Go+ fully supports `cgo`.


## Go+ features

### Rational number: bigint, bigrat, bigfloat

We introduce the rational number as native Go+ types. We use suffix `r` to denote rational literals. For example, (1r << 200) means a big int whose value is equal to 2<sup>200</sup>. And 4/5r means the rational constant 4/5.

```go
var a bigint = 1r << 65  // bigint, large than int64
var b bigrat = 4/5r      // bigrat
c := b - 1/3r + 3 * 1/2r // bigrat
println a, b, c

var x *big.Int = 1r << 65 // (1r << 65) is untyped bigint, and can be assigned to *big.Int
var y *big.Rat = 4/5r
println x, y
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

### Lambda expression

```go
func plot(fn func(x float64) float64) {
    // ...
}

func plot2(fn func(x float64) (float64, float64)) {
    // ...
}

plot x => x * x           // plot(func(x float64) float64 { return x * x })
plot2 x => (x * x, x + x) // plot2(func(x float64) (float64, float64) { return x * x, x + x })
```

### Deduce struct type

```go
type Config struct {
    Dir   string
    Level int
}

func foo(conf *Config) {
    // ...
}

foo {Dir: "/foo/bar", Level: 1}
```

Here `foo {Dir: "/foo/bar", Level: 1}` is equivalent to `foo(&Config{Dir: "/foo/bar", Level: 1})`. However, you can't replace `foo(&Config{"/foo/bar", 1})` with `foo {"/foo/bar", 1}`, because it is confusing to consider `{"/foo/bar", 1}` as a struct literal.

You also can omit struct types in a return statement. For example:

```go
type Result struct {
    Text string
}

func foo() *Result {
    return {Text: "Hi, Go+"} // return &Result{Text: "Hi, Go+"}
}
```


### List comprehension

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

### Select data from a collection

```go
type student struct {
    name  string
    score int
}

students := [student{"Ken", 90}, student{"Jason", 80}, student{"Lily", 85}]

unknownScore, ok := {x.score for x <- students, x.name == "Unknown"}
jasonScore := {x.score for x <- students, x.name == "Jason"}

println unknownScore, ok // output: 0 false
println jasonScore // output: 80
```

### Check if data exists in a collection

```go
type student struct {
    name  string
    score int
}

students := [student{"Ken", 90}, student{"Jason", 80}, student{"Lily", 85}]

hasJason := {for x <- students, x.name == "Jason"} // is any student named Jason?
hasFailed := {for x <- students, x.score < 60}     // is any student failed?
```

### For loop

```go
sum := 0
for x <- [1, 3, 5, 7, 11, 13, 17], x > 3 {
    sum += x
}
```


### Range expression (`start:end:step`)

```go
for i <- :10 {
    println i
}

for i := range :10:2 {
    println i
}

for i := range 1:10:3 {
    println i
}

for range :10 {
    println "Range expression"
}
```


### For range of UDT

```go
type Foo struct {
}

// Gop_Enum(proc func(val ValType)) or:
// Gop_Enum(proc func(key KeyType, val ValType))
func (p *Foo) Gop_Enum(proc func(key int, val string)) {
    // ...
}

foo := &Foo{}
for k, v := range foo {
    println k, v
}

for k, v <- foo {
    println k, v
}

println {v: k for k, v <- foo}
```

**Note: you can't use break/continue or return statements in for range of udt.Gop_Enum(callback).**


### For range of UDT2

```go
type FooIter struct {
}

// (Iterator) Next() (val ValType, ok bool) or:
// (Iterator) Next() (key KeyType, val ValType, ok bool)
func (p *FooIter) Next() (key int, val string, ok bool) {
    // ...
}

type Foo struct {
}

// Gop_Enum() Iterator
func (p *Foo) Gop_Enum() *FooIter {
    // ...
}

foo := &Foo{}
for k, v := range foo {
    println k, v
}

for k, v <- foo {
    println k, v
}

println {v: k for k, v <- foo}
```

### Overload operators

```go
import "math/big"

type MyBigInt struct {
    *big.Int
}

func Int(v *big.Int) MyBigInt {
    return MyBigInt{v}
}

func (a MyBigInt) + (b MyBigInt) MyBigInt { // binary operator
    return MyBigInt{new(big.Int).Add(a.Int, b.Int)}
}

func (a MyBigInt) += (b MyBigInt) {
    a.Int.Add(a.Int, b.Int)
}

func -(a MyBigInt) MyBigInt { // unary operator
    return MyBigInt{new(big.Int).Neg(a.Int)}
}

a := Int(1r)
a += Int(2r)
println a + Int(3r)
println -a
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

println `add("100", "23"):`, add("100", "23")!

sum, err := add("10", "abc")
println `add("10", "abc"):`, sum, err

println `addSafe("10", "abc"):`, addSafe("10", "abc")
```

The output of this example is:

```
add("100", "23"): 123
add("10", "abc"): 0 strconv.Atoi: parsing "abc": invalid syntax

===> errors stack:
main.add("10", "abc")
    /Users/xsw/tutorial/15-ErrWrap/err_wrap.gop:6 strconv.Atoi(y)?

addSafe("10", "abc"): 10
```

Compared to corresponding Go code, It is clear and more readable.

And the most interesting thing is, the return error contains the full error stack. When we got an error, it is very easy to position what the root cause is.

How these `ErrWrap expressions` work? See [Error Handling](https://github.com/goplus/gop/wiki/Error-Handling) for more information.


### Auto property

Let's see an example written in Go+:

```go
import "gop/ast/goptest"

doc := goptest.New(`... Go+ code ...`)!

println doc.Any().FuncDecl().Name()
```

In many languages, there is a concept named `property` who has `get` and `set` methods.

Suppose we have `get property`, the above example will be:

```go
import "gop/ast/goptest"

doc := goptest.New(`... Go+ code ...`)!

println doc.any.funcDecl.name
```

In Go+, we introduce a concept named `auto property`. It is a `get property`, but is implemented automatically. If we have a method named `Bar()`, then we will have a `get property` named `bar` at the same time.


### Unix shebang

You can use Go+ programs as shell scripts now. For example:

```go
#!/usr/bin/env -S gop run

println "Hello, Go+"

println 1r << 129
println 1/3r + 2/7r*2

arr := [1, 3, 5, 7, 11, 13, 17, 19]
println arr
println [x*x for x <- arr, x > 3]

m := {"Hi": 1, "Go+": 2}
println m
println {v: k for k, v <- m}
println [k for k, _ <- m]
println [v for v <- m]
```

Go [20-Unix-Shebang/shebang](https://github.com/goplus/tutorial/blob/main/20-Unix-Shebang/shebang) to get the source code.


### Go features

All Go features (including partially support `cgo`) will be supported. In bytecode mode, Go+ doesn't support `cgo`. However, in Go-code-generation mode, Go+ fully supports `cgo`.


## IDE Plugins

* vscode: https://github.com/gopcode/vscode-goplus


## Contributing

The Go+ project welcomes all contributors. We appreciate your help!

Here are [list of Go+ Contributors](https://github.com/goplus/gop/wiki/Goplus-Contributors). We award an email account (XXX@goplus.org) for every contributor. And we suggest you commit code by using this email account:

```bash
git config --global user.email XXX@goplus.org
```

If you did this, remember to add your `XXX@goplus.org` email to https://github.com/settings/emails.

What does `a contributor to Go+` mean? You must meet one of the following conditions:

* At least one pull request of a full-featured implemention.
* At least three pull requests of feature enhancements.
* At least ten pull requests of any kind issues.

Where can you start?

* [![Issues](https://img.shields.io/badge/ISSUEs-Go+-blue.svg)](https://github.com/goplus/gop/issues)
* [![Issues](https://img.shields.io/badge/ISSUEs-NumGo+-blue.svg)](https://github.com/numgoplus/ng/issues)
* [![Issues](https://img.shields.io/badge/ISSUEs-PandasGo+-blue.svg)](https://github.com/goplus/pandas/issues)
* [![Issues](https://img.shields.io/badge/ISSUEs-vscode%20Go+-blue.svg)](https://github.com/gopcode/vscode-goplus/issues)
* [![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/goplus/gop)](https://www.tickgit.com/browse?repo=github.com/goplus/gop)
