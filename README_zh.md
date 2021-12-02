## 适用于工程、STEM 教育和数据科学的 Go+ 语言

[![Build Status](https://github.com/goplus/gop/actions/workflows/go.yml/badge.svg)](https://github.com/goplus/gop/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/gop)](https://goreportcard.com/report/github.com/goplus/gop)
[![Coverage Status](https://codecov.io/gh/goplus/gop/branch/main/graph/badge.svg)](https://codecov.io/gh/goplus/gop)
[![GitHub release](https://img.shields.io/github/v/tag/goplus/gop.svg?label=release)](https://github.com/goplus/gop/releases)
[![Tutorials](https://img.shields.io/badge/tutorial-Go+-blue.svg)](https://tutorial.goplus.org/)
[![Playground](https://img.shields.io/badge/playground-Go+-blue.svg)](https://play.goplus.org/)
[![VSCode](https://img.shields.io/badge/vscode-Go+-teal.svg)](https://github.com/gopcode/vscode-goplus)
[![Readme](https://img.shields.io/badge/README-%E4%B8%AD%E6%96%87-teal.svg)](https://github.com/goplus/gop/blob/main/README.md)

-   **对于工程领域而言** ：意味着大家将会使用一门简单、易学、好用的编程语言。
-   **对于 STEM 教育而言** ：掌握一门可能会在将来工作中用到的工程语言将大有裨益。
-   **对于数据科学家而言** ：意味着你们将和工程师使用同一种语言进行合作和交流。

## 如何安装 Go+?

当前，我们推荐下载源代码来安装 Go+。

```bash
git clone https://github.com/goplus/gop.git
cd gop

# On linux run:
./all.bash
# On Windows run:
all.bat
```

实际上， `all.bash` 和 `all.bat` 内部都会执行 `go run cmd/make.go`

## 代码风格（很重要）

-   [https://tutorial.goplus.org/hello-world](https://tutorial.goplus.org/hello-world)

## Go+ 概要

关于 Go+ 的主要印象是什么？

-   一种静态类型语言。
-   完全兼容 [Go 语言](https://github.com/golang/go) 。
-   类似脚本语言的风格，代码比 Go 更易读。

例如，以下是合法的 Go+ 源代码：

```coffee
println [1, 2, 3.4]
```

我们如何在 Go 语言中做到这一点？

```go
package main

import "fmt"

func main() {
    fmt.Println([]float64{1, 2, 3.4})
}
```

当然，我们不仅仅做少打字的事情。

例如，我们支持 [列表推导](https://en.wikipedia.org/wiki/List_comprehension) ，这使得数据处理更容易。

```go
println [x*x for x <- 1:6:2] // output: [1 9 25]

mapData := {"Hi": 1, "Hello": 2, "Go+": 3}
reversedMap := {v: k for k, v <- mapData}
println reversedMap // output: map[1:Hi 2:Hello 3:Go+]
```

我们会保持 Go+ 的简单性。这就是为什么我们称它为 Go+，而不是 Go++。

少即是指数级的多。

Go 是如此，Go+ 亦是如此。

## 教程

-   [https://tutorial.goplus.org/](https://tutorial.goplus.org/)

## Playground

基于 Docker 的 Go+ Playground（推荐）：

-   [https://play.goplus.org/](https://play.goplus.org/)

基于 GopherJS 的 Go+ Playground（目前仅在 v0.7.x 中可用）：

-   [https://jsplay.goplus.org/](https://jsplay.goplus.org/)

## 与 Go 的兼容性

将支持所有 Go 特性（包括部分支持 `cgo` ，见 [下文](#bytecode-vs-go-code) ）。

**所有 Go 包（即便这些软件包使用 `cgo` ）都可以由 Go+ 导入。**

```coffee
import (
    "fmt"
    "strings"
)

x := strings.NewReplacer("?", "!").Replace("hello, world???")
fmt.Println "x:", x
```

**并且所有 Go+ 包也可以导入到 Go 程序中。您需要做的只是使用 `gop` 命令而不是 `go` 。**

首先，让我们创建一个名为 `14-Using-goplus-in-Go` 的目录。

然后在里面编写一个名为 [foo](https://github.com/goplus/tutorial/tree/main/14-Using-goplus-in-Go/foo) 的 Go+ 包：

```go
package foo

func ReverseMap(m map[string]int) map[int]string {
    return {v: k for k, v <- m}
}
```

然后在 Go 包 [14-Using-goplus-in-Go/gomain](https://github.com/goplus/tutorial/tree/main/14-Using-goplus-in-Go/gomain) 中使用它：

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

如何构建这个例子？您可以使用：

```bash
gop install -v ./...
```

访问 [github.com/goplus/tutorial/14-Using-goplus-in-Go](https://github.com/goplus/tutorial/tree/main/14-Using-goplus-in-Go) 来获取源代码。

## 字节码与 Go 代码

Go+ 支持字节码后端和 Go 代码生成。

当我们使用 `gop` 命令时，它会生成 Go 代码，将 Go+ 包转换为 Go 包。

```bash
gop run     # Run a Go+ program
gop install # Build Go+ files and install target to GOBIN
gop build   # Build Go+ files
gop test    # Test Go+ packages
gop fmt     # Format Go+ packages
gop clean   # Clean all Go+ auto generated files
gop go      # Convert Go+ packages into Go packages
```

当我们使用 [`igop`](https://github.com/goplus/igop) 命令时，它会执行生成的的字节码。

```bash
igop  # Run a Go+ program
```

在字节码模式下，Go+ 不支持 `cgo` 。但是，在 Go 代码生成模式下，Go+ 完全支持 `cgo` 。

## Go+ 的特性

### 有理数：bigint、bigrat、bigfloat

我们将有理数作为原生 Go+ 类型引入。我们使用后缀 `r` 来表示有理文字。例如， (1r << 200) 表示一个 big int ，其值等于 2 200 。 4/5r 表示有理常数 4/5。

```go
var a bigint = 1r << 65  // bigint, large than int64
var b bigrat = 4/5r      // bigrat
c := b - 1/3r + 3 * 1/2r // bigrat
println a, b, c

var x *big.Int = 1r << 65 // (1r << 65) is untyped bigint, and can be assigned to *big.Int
var y *big.Rat = 4/5r
println x, y
```

### 映射字面量

```go
x := {"Hello": 1, "xsw": 3.4} // map[string]float64
y := {"Hello": 1, "xsw": "Go+"} // map[string]interface{}
z := {"Hello": 1, "xsw": 3} // map[string]int
empty := {} // map[string]interface{}
```

### 切片字面量

```go
x := [1, 3.4] // []float64
y := [1] // []int
z := [1+2i, "xsw"] // []interface{}
a := [1, 3.4, 3+4i] // []complex128
b := [5+6i] // []complex128
c := ["xsw", 3] // []interface{}
empty := [] // []interface{}
```

### Lambda 表达式

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

### 推导结构类型

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

这里 `foo {Dir: "/foo/bar", Level: 1}` 等价于 `foo (&Config {Dir: "/foo/bar", Level: 1})` 。但是，您不能将 `foo (&Config {"/foo/bar", 1})` 替换为 `foo {"/foo/bar", 1}` ，因为将 `{"/foo/bar", 1}` 视为结构字面量会让人摸不着头脑。

您还可以在 return 语句中省略结构类型。例如：

```go
type Result struct {
    Text string
}

func foo() *Result {
    return {Text: "Hi, Go+"} // return &Result{Text: "Hi, Go+"}
}
```

### 列表推导

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

### 从集合中选择数据

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

### 检查数据是否存在于集合中

```go
type student struct {
    name  string
    score int
}

students := [student{"Ken", 90}, student{"Jason", 80}, student{"Lily", 85}]

hasJason := {for x <- students, x.name == "Jason"} // is any student named Jason?
hasFailed := {for x <- students, x.score < 60}     // is any student failed?
```

### For 循环

```go
sum := 0
for x <- [1, 3, 5, 7, 11, 13, 17], x > 3 {
    sum += x
}
```

### Range 表达式（ `start:end:step` ）

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

**注意：对于 udt.Gop\_Enum（回调）的范围，无法使用 break/continue 或 return 语句。**

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

### 重载运算符

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

### 错误处理

我们在 Go+ 中重新发明了错误处理规范。我们称之为 `ErrWrap expressions` ：

```go
expr! // panic if err
expr? // return if err
expr?:defval // use defval if err
```

如何使用它们？以下是一个例子：

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

这个例子的输出是：

```
add("100", "23"): 123
add("10", "abc"): 0 strconv.Atoi: parsing "abc": invalid syntax

===> errors stack:
main.add("10", "abc")
    /Users/xsw/tutorial/15-ErrWrap/err_wrap.gop:6 strconv.Atoi(y)?

addSafe("10", "abc"): 10
```

与相应的 Go 代码相比，它更清晰、更具可读性。

最有趣的是，返回错误包含完整的错误堆栈。当我们遇到错误时，很容易定位根本原因是什么。

这些 `ErrWrap expressions` 如何工作的？有关更多信息，请参阅 [错误处理。](https://github.com/goplus/gop/wiki/Error-Handling)

### Auto property

让我们看一个用 Go+ 编写的例子：

```go
import "gop/ast/goptest"

doc := goptest.New(`... Go+ code ...`)!

println doc.Any().FuncDecl().Name()
```

在许多语言中，有一个名为 `property` 的概念，它具有 `get` 和 `set` 方法。

假设我们有 `get property` ，上面的例子将是：

```go
import "gop/ast/goptest"

doc := goptest.New(`... Go+ code ...`)!

println doc.any.funcDecl.name
```

在 Go+ 中，我们引入了一个名为 `auto property` 的概念。它是一个 `get property` ，但会自动实现。如果我们有一个名为 `Bar()` 的方法，那么我们将同时拥有一个名为 `bar` 的 `get property` 。

### Unix shebang

您现在可以将 Go+ 程序用作 shell 脚本。例如：

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

访问 [20-Unix-Shebang/shebang](https://github.com/goplus/tutorial/blob/main/20-Unix-Shebang/shebang) 来获取源代码。

### Go 特性

将支持所有 Go 特性（包括部分支持 `cgo`）。在字节码模式下，Go+ 不支持 `cgo` 。然而，在 Go 代码生成模式下，Go+ 完全支持 `cgo` 。

## IDE 插件

-   vscode： [https://github.com/gopcode/vscode-goplus](https://github.com/gopcode/vscode-goplus)

## 贡献

Go+ 项目欢迎所有贡献者。我们感谢您的帮助！

这是 [Go+ 贡献者列表](https://github.com/goplus/gop/wiki/Goplus-Contributors) 。我们为每位贡献者颁发一个电子邮件帐户 ([XXX@goplus.org)。](mailto:XXX@goplus.org) 我们建议您使用此电子邮件帐户提交代码：

```bash
git config --global user.email XXX@goplus.org
```

如果您这样做了，请记住将您的 `XXX@goplus.org` 电子邮件添加到 [https://github.com/settings/emails](https://github.com/settings/emails) 。

如何成为一个 `Go+ 贡献者`？您必须满足以下条件之一：

-   至少一个实现完整功能的拉取请求。
-   至少三个功能增强的拉取请求。
-   至少十个任意类型的 issues 拉取请求。

您可以从哪里开始？

-   [![Issues](https://img.shields.io/badge/ISSUEs-Go+-blue.svg)](https://github.com/goplus/gop/issues)
-   [![Issues](https://img.shields.io/badge/ISSUEs-NumGo+-blue.svg)](https://github.com/numgoplus/ng/issues)
-   [![Issues](https://img.shields.io/badge/ISSUEs-PandasGo+-blue.svg)](https://github.com/goplus/pandas/issues)
-   [![Issues](https://img.shields.io/badge/ISSUEs-vscode%20Go+-blue.svg)](https://github.com/gopcode/vscode-goplus/issues)
-   [![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/goplus/gop)](https://www.tickgit.com/browse?repo=github.com/goplus/gop)
