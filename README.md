Q Language - A script language for Go
========

[![LICENSE](https://img.shields.io/github/license/qiniu/qlang.svg)](https://github.com/qiniu/qlang/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/qiniu/qlang.png?branch=master)](https://travis-ci.org/qiniu/qlang)
[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/qlang)](https://goreportcard.com/report/github.com/qiniu/qlang)
[![GitHub release](https://img.shields.io/github/v/tag/qiniu/qlang.svg?label=release)](https://github.com/qiniu/qlang/releases)
[![Coverage Status](https://codecov.io/gh/qiniu/qlang/branch/master/graph/badge.svg)](https://codecov.io/gh/qiniu/qlang)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://godoc.org/github.com/qiniu/qlang)

[![Qiniu Logo](http://open.qiniudn.com/logo.png)](http://www.qiniu.com/)

## Mission & vision

最新的 qlang v6 版本在能力上较以前的版本有极大的调整。其核心变化为：

- 完全推翻重来，从动态类型转向静态类型！
- 完全兼容 Go 语言文法。
- 在 Go 语言兼容基础上，保留当初 qlang 动态类型版本的重要特性。比如：

```
a := [1, 2, 3.4]
```

为什么人们需要 qlang？它的使命与愿景是什么？

一句话：qlang 希望能够将 Go 语言带到数据科技（DT）的世界。

对于服务端编程的最佳实践而言，Go 语言非常优雅。所以 Go 制霸了云计算了云计算领域。

但是当前 Go 是有舒适区的，从数据科技的角度，Go 语言显得很笨拙。有没有办法让 Go 在数据科技领域变得一样优雅？

这就是 qlang v6 的由来。它兼容 Go，扩展 Go，以 Go++ 的姿态出现，让数据科技享受 Go 的简洁之美。

qlang 将支持生成 Go 代码，方便 Go 语言编译并与其他代码集成。

关于新版本的详细规划，参考：

* https://github.com/qiniu/qlang/issues/176

代码样例：

* https://github.com/qiniu/qlang/tree/v6.x/tutorial

## Old versions

当前 qlang v6 还在快速迭代中。在正式场合建议使用正式 release 的版本：

* https://github.com/qiniu/qlang/releases

最近的老版本代码可以从 qlang v1.5 分支获得：

* https://github.com/qiniu/qlang/tree/v1.5

## Supported features

### Variable & operator

```go
x := 123.1 - 3i
y, z := 1, 123
s := "Hello"

println(s + " complex")
println(x - 1, y * z)
```

### Condition

```go
x := 0
if t := false; t {
    x = 3
} else {
    x = 5
}

x = 0
switch s := "Hello"; s {
default:
    x = 7
case "world", "hi":
    x = 5
case "xsw":
    x = 3
}

v := "Hello"
switch {
case v == "xsw":
    x = 3
case v == "Hello", v == "world":
    x = 5
default:
    x = 7
}
```

### Import go package

```go
import (
    "fmt"
    "strings"
)

x := strings.NewReplacer("?", "!").Replace("hello, world???")
fmt.Println("x:", x)
```

### Func & closure

```go
import (
    "fmt"
    "strings"
)

func foo(x string) string {
    return strings.NewReplacer("?", "!").Replace(x)
}

func printf(format string, args ...interface{}) (n int, err error) {
    n, err = fmt.Printf(format, args...)
    return
}

x := "qlang"
fooVar := func(prompt string) (n int, err error) {
    n, err = fmt.Println(prompt + x)
    return
}

printfVar := func(format string, args ...interface{}) (n int, err error) {
    n, err = fmt.Printf(format, args...)
    return
}
```

### Map, array & slice

```go
x := []float64{1, 3.4, 5}
y := map[string]float64{"Hello": 1, "xsw": 3.4}

a := [...]float64{1, 3.4, 5}
b := [...]float64{1, 3: 3.4, 5}
c := []float64{2: 1.2, 3, 6: 4.5}
```

### Map literal

```go
x := {"Hello": 1, "xsw": 3.4} // map[string]float64
y := {"Hello": 1, "xsw": "qlang"} // map[string]interface{}
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
a := [x * x for x <- [1, 3, 5, 7, 11]]
b := [x * x for x <- [1, 3, 5, 7, 11], x > 3]
c := [i + v for i, v <- [1, 3, 5, 7, 11], i%2 == 1]
d := [k + "," + s for k, s <- {"Hello": "xsw", "Hi": "qlang"}]

arr := [1, 2, 3, 4, 5, 6]
e := [[a, b] for a <- arr, a < b for b <- arr, b > 2]

x := {x: i for i, x <- [1, 3, 5, 7, 11]}
y := {x: i for i, x <- [1, 3, 5, 7, 11], i%2 == 1}
z := {v: k for k, v <- {1: "Hello", 3: "Hi", 5: "xsw", 7: "qlang"}, k > 3}
```

### For loop

```go
sum := 0
for x <- [1, 3, 5, 7, 11, 13, 17], x > 3 {
    sum += x
}
```
