Go+ Quick Start
======

Our vision is to **enable everyone to create production-level applications**.

#### Easy to learn

* Simple and easy to understand
* Smaller syntax set than Python in best practices

#### Ready for large projects

* Derived from Go and easy to build large projects from its good engineering foundation

The Go+ programming language is designed for engineering, STEM education, and data science.

* **For engineering**: working in the simplest language that can be mastered by children.
* **For STEM education**: studying an engineering language that can be used for work in the future.
* **For data science**: communicating with engineers in the same language.

## How to install

Note: Requires go1.19 or later

### on Windows

```sh
winget install goplus.gop
```

### on Debian/Ubuntu

```sh
sudo bash -c ' echo "deb [trusted=yes] https://pkgs.goplus.org/apt/ /" > /etc/apt/sources.list.d/goplus.list'
sudo apt update
sudo apt install gop
```

### on RedHat/CentOS/Fedora

```sh
sudo bash -c 'echo -e "[goplus]\nname=Go+ Repo\nbaseurl=https://pkgs.goplus.org/yum/\nenabled=1\ngpgcheck=0" > /etc/yum.repos.d/goplus.repo'
sudo yum install gop
```

### on macOS/Linux (Homebrew)

Install via [brew](https://brew.sh/)

```sh
$ brew install goplus
```

### from source code

```bash
git clone https://github.com/goplus/gop.git
cd gop

# On mac/linux run:
./all.bash
# On Windows run:
all.bat
```

Actually, `all.bash` and `all.bat` will use `go run cmd/make.go` underneath.


## Running in Go+ playground

If you don't want install Go+, you can write your Go+ programs in Go+ playground. This is the fastest way to experience Go+.

* Go+ playground based on Docker: https://play.goplus.org/
* Go+ playground based on GopherJS: https://jsplay.goplus.org/

And you can share your Go+ code with your friends.
Here is my `Hello world` program:
* https://play.goplus.org/p/AAh_gQAKAZR.


## Table of Contents

<table>
    <tr><td width=33% valign=top>

* [Hello world](#hello-world)
* [Running a project folder](#running-a-project-folder-with-several-files)
* [Comments](#comments)
* [Variables](#variables)
* [Go+ types](#go-types)
    * [Strings](#strings)
    * [Numbers](#numbers)
    * [Slices](#slices)
    * [Maps](#maps)
* [Module imports](#module-imports)

</td><td width=33% valign=top>

* [Statements & expressions](#statements--expressions)
    * [If..else](#ifelse)
    * [For loop](#for-loop)
    * [Error handling](#error-handling)
* [Functions](#functions)
    * [Returning multiple values](#returning-multiple-values)
    * [Variadic parameters](#variadic-parameters)
    * [Higher order functions](#higher-order-functions)
    * [Lambda expressions](#lambda-expressions)
* [Structs](#structs)

</td><td valign=top>

* [Go/Go+ hybrid programming](#gogo-hybrid-programming)
    * [Run Go+ in watch mode](#run-go-in-watch-mode)
* [Calling C from Go+](#calling-c-from-go)
* [Data processing](#data-processing)
    * [Rational numbers](#rational-numbers)
    * [List comprehension](#list-comprehension)
    * [Select data from a collection](#select-data-from-a-collection)
    * [Check if data exists in a collection](#check-if-data-exists-in-a-collection)
* [Unix shebang](#unix-shebang)
* [Compatibility with Go](#compatibility-with-go)

</td></tr>
</table>


## Hello World

Different from the function call style of most languages, Go+ recommends command style code:

```go
println "Hello world"
```

Save this snippet into a file named `hello.gop`. Now do: `gop run hello.gop`.

Congratulations - you just wrote and executed your first Go+ program!

You can compile a program without execution with `gop build hello.gop`.
See `gop help` for all supported commands.

[`println`](#println) is one of the few [built-in functions](#builtin-functions).
It prints the value passed to it to standard output.

To emphasize our preference for command style, we introduce `echo` as an alias for `println`:

```coffee
echo "Hello world"
```

See https://tutorial.goplus.org/hello-world for more details.

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


## Running a project folder with several files

Suppose you have a folder with several .gop files in it, and you want 
to compile them all into one program. Just do: `gop run .`.

Passing parameters also works, so you can do:
`gop run . --yourparams some_other_stuff`.

Your program can then use the CLI parameters like this:

```go
import "os"

println os.Args
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


## Comments

```go
# This is a single line comment.

// This is a single line comment.

/*
This is a multiline comment.
*/
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


## Variables

```go
name := "Bob"
age := 20
largeNumber := int128(1 << 65)
println name, age
println largeNumber
```

Variables are declared and initialized with `:=`.

The variable's type is inferred from the value on the right hand side.
To choose a different type, use type conversion:
the expression `T(v)` converts the value `v` to the
type `T`.

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>

### Initialization vs. assignment

Note the (important) difference between `:=` and `=`.
`:=` is used for declaring and initializing, `=` is used for assigning.

```go failcompile
age = 21
```

This code will not compile, because the variable `age` is not declared.
All variables need to be declared in Go+.

```go
age := 21
```

The values of multiple variables can be changed in one line.
In this way, their values can be swapped without an intermediary variable.

```go
a, b := 0, 1
a, b = b, a
println a, b // 1, 0
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


## Go+ Types

### Primitive types

```go ignore
bool

int8    int16   int32   int    int64    int128
uint8   uint16  uint32  uint   uint64   uint128

uintptr // similar to C's size_t

byte // alias for uint8
rune // alias for int32, represents a Unicode code point

string

float32 float64

complex64 complex128

bigint bigrat

unsafe.Pointer // similar to C's void*

any // alias for Go's interface{}
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>

### Strings

```go
name := "Bob"
println name.len  // 3
println name[0]   // 66
println name[1:3] // ob
println name[:2]  // Bo
println name[2:]  // b

// or using octal escape `\###` notation where `#` is an octal digit
println "\141a"   // aa

// Unicode can be specified directly as `\u####` where # is a hex digit
// and will be converted internally to its UTF-8 representation
println "\u2605"  // ★
```

String values are immutable. You cannot mutate elements:

```go failcompile
s := "hello 🌎"
s[0] = `H` // not allowed
```

Note that indexing a string will produce a `byte`, not a `rune` nor another `string`.

Strings can be easily converted to integers:

```go
s := "12"
a, err := s.int
b := s.int! // will panic if s isn't a valid integer
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


#### String operators

```go
name := "Bob"
bobby := name + "by" // + is used to concatenate strings
println bobby // Bobby

s := "Hello "
s += "world"
println s // Hello world
```

Most Go+ operators must have values of the same type on both sides. You cannot concatenate an
integer to a string:

```go failcompile
age := 10
println "age = " + age // not allowed
```

We have to either convert `age` to a `string`:

```go
age := 10
println "age = " + age.string
```

However, you can replace `age.string` to `"${age}"`:

```go
age := 10
println "age = ${age}"
```

Here is a more complex example of `${expr}`:

```go
host := "example.com"
page := 0
limit := 20
println "https://${host}/items?page=${page+1}&limit=${limit}" // https://example.com/items?page=1&limit=20
println "$$" // $
```


<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


### Runes

A `rune` represents a single Unicode character and is an alias for `int32`.

```go
rocket := '🚀'
println rocket         // 128640
println string(rocket) // 🚀
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


### Numbers

```go
a := 123
```

This will assign the value of 123 to `a`. By default `a` will have the
type `int`.

You can also use hexadecimal, binary or octal notation for integer literals:

```go
a := 0x7B
b := 0b01111011
c := 0o173
```

All of these will be assigned the same value, 123. They will all have type
`int`, no matter what notation you used.

Go+ also supports writing numbers with `_` as separator:

```go
num := 1_000_000 // same as 1000000
```

If you want a different type of integer, you can use casting:

```go
a := int64(123)
b := uint8(12)
c := int128(12345)
```

Assigning floating point numbers works the same way:

```go
f1 := 1.0
f2 := float32(3.14)
```

If you do not specify the type explicitly, by default float literals will have the type of `float64`.

Float literals can also be declared as a power of ten:

```go
f0 := 42e1   // 420
f1 := 123e-2 // 1.23
f2 := 456e+2 // 45600
```

Go+ has built-in support for [rational numbers](#rational-numbers):

```go
a := 1r << 200  // suffix `r` means `rational`
b := bigint(1 << 200)
```

And you can cast bool to number types (this is NOT supported in Go):

```go
println int(true)       // 1
println float64(true)   // 1
println complex64(true) // (1+0i)
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


### Slices

A slice is a collection of data elements of the same type. A slice literal is a
list of expressions surrounded by square brackets. An individual element can be
accessed using an *index* expression. Indexes start from `0`:

```go
nums := [1, 2, 3]
println nums      // [1 2 3]
println nums.len  // 3
println nums[0]   // 1
println nums[1:3] // [2 3]
println nums[:2]  // [1 2]
println nums[2:]  // [3]

nums[1] = 5
println nums // [1 5 3]
```

Type of a slice literal is infered automatically.

```go
a := [1, 2, 3]   // []int
b := [1, 2, 3.4] // []float64
c := ["Hi"]      // []string
d := ["Hi", 10]  // []any
d := []          // []any
```

And casting slice literals also works.

```go
a := []float64([1, 2, 3]) // []float64
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


### Maps

```go
a := {"Hello": 1, "xsw": 3}     // map[string]int
b := {"Hello": 1, "xsw": 3.4}   // map[string]float64
c := {"Hello": 1, "xsw": "Go+"} // map[string]any
d := {}                         // map[string]any
```

If a key is not found, a zero value is returned by default:

```go
a := {"Hello": 1, "xsw": 3}
c := {"Hello": 1, "xsw": "Go+"}
println a["bad_key"] // 0
println c["bad_key"] // <nil>
```

You can also check, if a key is present, and get its value.

```go
a := {"Hello": 1, "xsw": 3}
if v, ok := a["xsw"]; ok {
    println "its value is", v
}
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


## Module imports

For information about creating a module, see [Modules](#modules).

Modules can be imported using the `import` keyword:

```go
import "strings"

x := strings.NewReplacer("?", "!").Replace("Hello, world???")
println x // Hello, world!!!
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


### Module import aliasing

Any imported module name can be aliased:

```go
import strop "strings"

x := strop.NewReplacer("?", "!").Replace("Hello, world???")
println x // Hello, world!!!
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


## Statements & expressions


### If..else

In Go+, `if` statements are pretty straightforward and similar to most other languages.
Unlike other C-like languages,
there are no parentheses surrounding the condition and the braces are always required.

```go
a := 10
b := 20
if a < b {
    println "a < b"
} else if a > b {
    println "a > b"
} else {
    println "a == b"
}
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


### For loop

Go+ has only one looping keyword: `for`, with several forms.

#### `for`/`<-`

This is the most common form. You can use it with a slice, map, numeric range or custom iterators.

For information about creating a custom iterators, see [Custom iterators](#custom-iterators).

##### Slice `for`

The `for value <- arr` form is used for going through elements of a slice.

```go
numbers := [1, 3, 5, 7, 11, 13, 17]
sum := 0
for x <- numbers {
    sum += x
}
println sum // 57
```

If an index is required, an alternative form `for index, value <- arr` can be used.

```go
names := ["Sam", "Peter"]
for i, name <- names {
    println i, name
    // 0 Sam
    // 1 Peter
}
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


##### Map `for`

```go
m := {"one": 1, "two": 2}
for key, val <- m {
    println key, val
    // one 1
    // two 2
}
for key, _ <- m {
    println key
    // one
    // two
}
for val <- m {
    println val
    // 1
    // 2
}
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


##### Range `for`

You can use `range expression` (`start:end:step`) in for loop.

```go
for i <- :5 {
    println i
    // 0
    // 1
    // 2
    // 3
    // 4
}
for i <- 1:5 {
    println i
    // 1
    // 2
    // 3
    // 4
}
for i <- 1:5:2 {
    println i
    // 1
    // 3
}
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


##### `for`/`<-`/`if`

All loops of `for`/`<-` form can have an optional `if` condition.

```go
numbers := [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
for num <- numbers if num%3 == 0 {
    println num
    // 0
    // 3
    // 6
    // 9
}

for num <- :10 if num%3 == 0 {
    println num
    // 0
    // 3
    // 6
    // 9
}
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


#### Condition `for`

```go
sum := 0
i := 1
for i <= 100 {
    sum += i
    i++
}
println sum // 5050
```

This form of the loop is similar to `while` loops in other languages.
The loop will stop iterating once the boolean condition evaluates to false.
Again, there are no parentheses surrounding the condition, and the braces are always required.

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


#### C `for`

```go
for i := 0; i < 10; i += 2 {
    // Don't print 6
    if i == 6 {
        continue
    }
    println i
    // 0
    // 2
    // 4
    // 8
}
```

Finally, there's the traditional C style `for` loop. It's safer than the `while` form
because with the latter it's easy to forget to update the counter and get
stuck in an infinite loop.

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


#### Bare `for`

```go
for {
    // ...
}
```

The condition can be omitted, resulting in an infinite loop. You can use `break` or `return` to end the loop.

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


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

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


## Functions

```go
func add(x int, y int) int {
    return x + y
}

println add(2, 3) // 5
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


### Returning multiple values

```go
func foo() (int, int) {
    return 2, 3
}

a, b := foo()
println a // 2
println b // 3
c, _ := foo() // ignore values using `_`
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


### Variadic parameters

```go
func sum(a ...int) int {
    total := 0
    for x <- a {
        total += x
    }
    return total
}

println sum(2, 3, 5) // 10
```

Output parameters can have names.

```go
func sum(a ...int) (total int) {
    for x <- a {
        total += x
    }
    return // don't need return values if they are assigned
}

println sum(2, 3, 5) // 10
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


### Higher order functions

Functions can also be parameters.

```go
func square(x float64) float64 {
    return x*x
}

func abs(x float64) float64 {
    if x < 0 {
        return -x
    }
    return x
}

func transform(a []float64, f func(float64) float64) []float64 {
    return [f(x) for x <- a]
}

y := transform([1, 2, 3], square)
println y // [1 4 9]

z := transform([-3, 1, -5], abs)
println z // [3 1 5]
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


### Lambda expressions

You also can use `lambda expression` to define a anonymous function.

```go
func transform(a []float64, f func(float64) float64) []float64 {
    return [f(x) for x <- a]
}

y := transform([1, 2, 3], x => x*x)
println y // [1 4 9]

z := transform([-3, 1, -5], x => {
    if x < 0 {
        return -x
    }
    return x
})
println z // [3 1 5]
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


## Structs

### Custom iterators

#### For range of UDT

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

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


#### For range of UDT2

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

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


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

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


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

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


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

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


## Go/Go+ hybrid programming

This is an example to show how to mix Go/Go+ code in the same package.

In this example, we have a Go source file named `a.go`:

```go
package main

import "fmt"

func p(a interface{}) {
    sayMix()
    fmt.Println("Hello,", a)
}
```

And we have a Go+ source file named `b.gop`:

```go
func sayMix() {
    println "Mix Go and Go+"
}

p "world"
```

You can see that Go calls a Go+ function named `sayMix`, and Go+ calls a Go function named `p`. As you are used to in Go programming, this kind of circular reference is allowed.

Run `gop run .` to see the output of this example:

```
Mix Go and Go+
Hello, world
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


### Run Go+ in watch mode

The `gop` command can run in watch mode so that everytime a Go+ file is changed it is transpiled to a Go file:

```
gop watch [-gentest] [dir]
```

By default `gop watch` does not convert test files (normally ending with `_test.gop`). You can specify `-gentest` flag to force converting all Go+ files.

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


## Calling C from Go+

Here is [an example to show how Go+ interacts with C](https://github.com/goplus/gop/tree/main/demo/_llgo/hellollgo).

```go
import "c"

c.printf c"Hello, llgo!\n"
c.fprintf c.Stderr, c"Hi, %6.1f\n", 3.14
```

Here `import "c"` is used to import libc. In this example we call two C standard functions `printf` and `fprintf`, passing a C variable `stderr` and two C strings in the form of `c"xxx"` (a Go+ syntax to represent C-style strings).

To run this demo, you need to set the `GOP_GOCMD` environment variable first.

```sh
export GOP_GOCMD=llgo  # default is `go`
```

Then execute `gop run .` to see the output of this example:

```
Hello, llgo!
Hi,    3.1
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


## Data processing

### Rational numbers

We introduce rational numbers as primitive Go+ types. We use suffix `r` to denote rational literals. For example, `1r << 200` means a big int whose value is equal to 2<sup>200</sup>.

```go
a := 1r << 200
b := bigint(1 << 200)
```

By default, `1r` will have the type of `bigint`.

And `4/5r` means the rational constant `4/5`.
It will have the type of `bigrat`.

```go
a := 4/5r
b := a - 1/3r + 3 * 1/2r
println a, b // 4/5 59/30
```

Casting rational numbers works like other [primitive types](#primitive-types):

```go
a := 1r
b := bigrat(1r)
c := bigrat(1)
println a/3 // 0
println b/3 // 1/3
println c/3 // 1/3
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


### List comprehension

```go
a := [x*x for x <- [1, 3, 5, 7, 11]]
b := [x*x for x <- [1, 3, 5, 7, 11] if x > 3]
c := [i+v for i, v <- [1, 3, 5, 7, 11] if i%2 == 1]

arr := [1, 2, 3, 4, 5, 6]
d := [[a, b] for a <- arr if a < b for b <- arr if b > 2]

x := {x: i for i, x <- [1, 3, 5, 7, 11]}
y := {x: i for i, x <- [1, 3, 5, 7, 11] if i%2 == 1}
z := {v: k for k, v <- {1: "Hello", 3: "Hi", 5: "xsw", 7: "Go+"} if k > 3}
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


### Select data from a collection

```go
type student struct {
    name  string
    score int
}

students := [student{"Ken", 90}, student{"Jason", 80}, student{"Lily", 85}]

unknownScore, ok := {x.score for x <- students if x.name == "Unknown"}
jasonScore := {x.score for x <- students if x.name == "Jason"}

println unknownScore, ok // 0 false
println jasonScore // 80
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


### Check if data exists in a collection

```go
type student struct {
    name  string
    score int
}

students := [student{"Ken", 90}, student{"Jason", 80}, student{"Lily", 85}]

hasJason := {for x <- students if x.name == "Jason"} // is any student named Jason?
hasFailed := {for x <- students if x.score < 60}     // is any student failed?
```

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


## Unix shebang

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

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


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

<h5 align="right"><a href="#table-of-contents">⬆ back to toc</a></h5>


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
