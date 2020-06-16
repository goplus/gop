GoPlus - The Go+ language for data science
========

[![LICENSE](https://img.shields.io/github/license/qiniu/goplus.svg)](https://github.com/qiniu/goplus/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/qiniu/goplus.png?branch=master)](https://travis-ci.org/qiniu/goplus)
[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/goplus)](https://goreportcard.com/report/github.com/qiniu/goplus)
[![GitHub release](https://img.shields.io/github/v/tag/qiniu/goplus.svg?label=release)](https://github.com/qiniu/goplus/releases)
[![Coverage Status](https://codecov.io/gh/qiniu/goplus/branch/master/graph/badge.svg)](https://codecov.io/gh/qiniu/goplus)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://godoc.org/github.com/qiniu/goplus)


## Summary about Go+

What are mainly impressions about Go+?

- A static typed language.
- Fully compatible with [the Go language](https://github.com/golang/go).
- Script-like style, and more readable code for data science than Go.

For example, the following is a legal Go+ source code:

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

For example, we  support `list comprehension`, which make data processing easier.

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

All Go features (not including `cgo`) will be supported.

* See [supported the Go language features](https://github.com/qiniu/goplus/wiki/Supported-Go-features).

**All Go packages (even these packages use `cgo`) can be imported by Go+.**

```go
import (
    "fmt"
    "strings"
)

x := strings.NewReplacer("?", "!").Replace("hello, world???")
fmt.Println("x:", x)
```

Be interested in how it works? See [Dive into Go+](https://github.com/qiniu/goplus/wiki/Dive-into-Goplus).

**Also, all Go+ packages can be converted into Go packages, and then be imported by Go.**

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

	"github.com/qiniu/goplus/tutorial/14-Using-goplus-in-Go/foo"
)

func main() {
	rmap := foo.ReverseMap(map[string]int{"Hi": 1, "Hello": 2})
	fmt.Println(rmap)
}
```

How to compile this exmaple?

```bash
gop go tutorial/ # Convert all Go+ packages in tutorial/ into Go packages
go install ./...
```

Or:

```bash
gop install ./... # Convert Go+ packages and go install ./...
```

Go [tutorial/14-Using-goplus-in-Go](https://github.com/qiniu/goplus/tree/master/tutorial/14-Using-goplus-in-Go) to get the source code.

Note: The `gop` command isn't provided currently (in alpha stage). Instead, we provide `qrun` and `qgo` commands.


## How to build

Current version: [![GitHub release](https://img.shields.io/github/v/tag/qiniu/goplus.svg?label=)](https://github.com/qiniu/goplus/releases)

```bash
go get github.com/qiniu/goplus@vX.X.XX
#or: git clone git@github.com:qiniu/goplus.git
cd goplus
go install -v ./...
```

## Tutorials

See https://github.com/qiniu/goplus/tree/master/tutorial


## Go+ features

### Bytecode vs. Go code

Go+ suppports bytecode backend or generating Go code.

When we use `gop go` or `gop install` command, it generates Go code to covert Go+ package into Go packages.

When we use `gop run` command, it doesn't call `go run` command. It generates bytecode to execute.

### Commands

```bash
gop go <gopSrcDir> # Convert all Go+ packages under <gopSrcDir> into Go packages, recursively
gop run <gopSrcDir> # Running <gopSrcDir> as a Go+ main package
gop run <gopSrcFile> # Running <gopSrcFile> as a Go+ script
gop install ./... # Convert all Go+ packages under ./ and go install ./...
```

The `gop` command isn't provided currently (in alpha stage). Instead, we provide the following commands:

```bash
qrun <gopSrcDir> # gop run <gopSrcDir>
qrun -asm <gopSrcDir> # generates `asm` code of Go+ bytecode backend
qrun -quiet <gopSrcDir> # don't generate any compiling stage log
qrun -debug <gopSrcDir> # print debug information
qrun -prof <gopSrcDir> # do profile and generate profile report
qgo <gopSrcDir> # gop go <gopSrcDir>
qgo -test <gopSrcDir>
```

Note:

* `qgo -test <gopSrcDir>` converts Go+ packages into Go packages, and for every package, it call `go run <gopPkgDir>/gop_autogen.go` and `qrun -quiet <gopPkgDir>` to compare their outputs. If their outputs aren't equal, the test case fails.

### Rational number: bigint, bigrat, bigfloat

We introduce rational number as native Go+ types. We use -r suffix as a rational constant. For example, (1r << 200) means a big int whose value is equal to 2<sup>200</sup>. And 4/5r means the rational constant 4/5.

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

We reinvent error handling specification in Go+. We call them `ErrWrap expressions`:

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

How these `ErrWrap expressions` work? See [Error Handling](https://github.com/qiniu/goplus/wiki/Error-Handling) for more information.

### Go features

All Go features (not including `cgo`) will be supported.

* See [supported the Go language features](https://github.com/qiniu/goplus/wiki/Supported-Go-features).

## Contributing

The Go+ project welcomes all contributors. We appreciate your help!

Here are [list of Go+ Contributors](https://github.com/qiniu/goplus/wiki/Goplus-Contributors). We award an email account (XXX@goplus.org) for every contributor. And we suggest you commit code by using this email account:

```bash
git config --global user.email XXX@goplus.org
```

What does `a contributor of Go+` means? He must meet one of the following conditions: 

* At least one pull request of a full feature implemention.
* At least three pull requests of feature enhancements.
* At least ten pull requests of any kind issues.
