Q Language Manual (Q语言手册)
========

## 运算符

基本上 Go 语言的操作符都支持，且操作符优先级相同。包括：

* '+'、'-'、'*'、'/'、'%'、'='
* '^'、'&'、'&^'、'|'、'<<'、'>>'
* '+='、'-='、'*='、'/='、'%='、'++'、'--'
* '^='、'&='、'&^='、'|='、'<<='、'>>='
* '!'、'>='、'<='、'>'、'<'、'=='、'!='、'&&'、'||'
* '<-' (chan操作符)

## 类型

原理上支持所有 Go 语言中的类型。典型有：

* 基本类型：int、float (在 Go 语言里面是 float64)、string、byte、bool、var（在 Go 语言里面是 interface{}）。
* 复合类型：slice、map、chan
* 用户自定义：函数（闭包）、类成员函数、类


## 常量

* 布尔类型：true、false（由 builtin 模块支持）
* var 类型：nil（由 builtin 模块支持）
* 浮点类型：pi、e、phi （由 math 模块支持）


## 变量及初始化

### 基本类型

```go
a = 1 // 创建一个 int 类型变量，并初始化为 1
b = "hello" // string 类型
c = true // bool 类型
d = 1.0 // float 类型
e = 'h' // byte 类型
```

### string 类型

```go
a = "hello, world"
```

和 Go 语言类似，string 有如下内置的操作：

```go
a = "hello" + "world" // + 为字符串连接操作符
n = len(a) // 取 a 字符串的长度
b = a[1] // 取 a 字符串的某个字符，得到的 b 是 byte 类型
c = a[1:4] // 取子字符串
```

### slice 类型

```go
a = [1, 2, 3] // 创建一个 int slice，并初始化为 [1, 2, 3]
b = [1, 2.3, 5] // 创建一个 float slice
c = ["a", "b", "c"] // 创建一个 string slice
d = ["a", 1, 2.3] // 创建一个 var slice (等价于 Go 语言的 []interface{})
e = make([]int, len, cap) // 创建一个 int slice，并将长度设置为 len，容量设置为 cap
f = make([][]int, len, cap) // 创建一个 []int 的 slice，并将长度设置为 len，容量设置为 cap
g = []byte{1, 2, 3} // 创建一个 byte slice，并初始化为 [1, 2, 3]
h = []byte(nil) // 创建一个空 byte slice
```

和 Go 语言类似，slice 有如下内置的操作：

```go
a = append(a, 4, 5, 6) // 含义与 Go 语言完全一致
n = len(a) // 取 a 的元素个数
m = cap(a) // 取 slice a 的容量
ncopy = copy(a, b) // 复制 b 的内容到 a，复制的长度 ncopy = min(len(a), len(b))
b1 = b[2] // 取 b 这个 slice 中 index=2 的元素
b[2] = 888 // 设置 b 这个 slice 中 index=2 的元素值为 888
b[1], b[2], b[3] = 777, 888, 999 // 设置 b 这个 slice 中 index=1, 2, 3 的三个元素值
b2 = b[1:4] // 取子slice
```

特别地，在 qlang 中可以这样赋值：

```go
x, y, z = [1, 2, 3]
```

结果是 x = 1, y = 2, z = 3。

实际上，qlang 支持多返回值就是通过 slice 的多赋值完成：

* 对于那些返回了多个值的 Go 函数，在 qlang 会理解为返回 var slice，也就是 []interface{}。

举个例子：

```go
f, err = os.Open(fname)
```

这个例子，在 Go 里面返回的是 (*os.File, error)。但是 qlang 中是 var slice。


### map 类型

```go
a = {"a": 1, "b": 2, "c": 3} // 得到 map[string]int 类型的对象
b = {"a": 1, "b", 2.3, "c": 3} // 得到 map[string]float64 类型的对象
c = {1: "a", 2: "b", 3: "c"} // 得到 map[int]string 类型的对象
d = {"a": "hello", "b": 2.0, "c": true} // 得到 map[string]interface{} 类型的对象
e = make(map[string]int) // 创建一个空的 map[string]int 类型的对象
f = make(map[string]map[string]int) // 创建一个 map[string]map[string]int 类型的对象
g = map[string]int16{"a": 1, "b": 2} // 创建一个 map[string]int16 类型的对象
```

和 Go 语言类似，map 有如下内置的操作：

```go
n = len(a) // 取 a 的元素个数
x = a["b"] // 取 a map 中 key 为 "b" 的元素
x = a.b // 含义同上，但如果 "b" 元素不存在会 panic
a["e"], a["f"], a["g"] = 4, 5, 6 // 同 Go 语言
a.e, a.f, a.g = 4, 5, 6 // 含义同 a["e"], a["f"], a["g"] = 4, 5, 6
delete(a, "e") // 删除 a map 中的 "e" 元素
```

需要注意的是，a["b"] 的行为和 Go 语言中略有不同。在 Go 语言中，常见的范式是：

```go
x := map[string]int{"a": 1, "b": 2}
a, ok := x["a"] // 结果：a = 1, ok = true
if ok { // 判断a存在的逻辑
	...
}
c, ok2 := x["c"] // 结果：c = 0, ok2 = false
d := x["d"] // 结果：d = 0
```

而在 qlang 中是这样的：

```go
x = {"a": 1, "b": 2}
a = x["a"] // 结果：a = 1
if a != undefined { // 判断a存在的逻辑
	...
}
c = x["c"] // 结果：c = undefined，注意不是0，也不是nil
d = x["d"] // 结果：d = undefined，注意不是0，也不是nil
```

### chan 类型

```go
ch1 = make(chan bool, 2) // 得到 buffer = 2 的 chan bool
ch2 = make(chan int) // 得到 buffer = 0 的 chan int
ch3 = make(chan map[string]int) // 得到 buffer = 0 的 chan map[string]int
```

和 Go 语言类似，chan 有如下内置的操作：

```go
n = len(ch1) // 取得chan当前的元素个数
m = cap(ch1) // 取得chan的容量
ch1 <- true // 向chan发送一个值
v = <-ch1 // 从chan取出一个值
close(ch1) // 关闭chan，被关闭的chan是不能写，但是还可以读(直到已经写入的值全部被取完为止)
```

需要注意的是，在 chan 被关闭后，<-ch 取得 undefined 值。所以在 qlang 中应该这样：

```go
v = <-ch1
if v != undefined { // 判断chan没有被关闭的逻辑
	...
}
```

## 流程控制

### if 语句

```go
if booleanExpr1 {
	// ...
} elif booleanExpr2 {
	// ...
} elif booleanExpr3 {
	// ...
} else {
	// ...
}
```

### switch 语句

```go
switch expr {
case expr1:
	// ...
case expr2:
	// ...
default:
	// ...
}
```

或者：

```go
switch {
case booleanExpr1:
	// ...
case booleanExpr2:
	// ...
default:
	// ...
}
```

### for 语句

```go
for { // 无限循环，需要在中间 break 或 return 结束
	...
}

for booleanExpr { // 类似很多语言的 while 循环
	...
}

for initExpr; conditionExpr; stepExpr {
	...
}
```

典型例子：

```go
for i = 0; i < 10; i++ {
	...
}
```

另外我们也支持 for..range 语法：

```go
for range collectionExpr { // 其中 collectionExpr 可以是 slice, map 或 chan
	...
}

for index = range collectionExpr {
	...
}

for index, value = range collectionExpr {
	...
}
```

## 函数

### 函数和闭包

基本语法：

```go
funcName = fn(arg1, arg2, argN) {
	//...
	return expr
}
```

这就定义了一个名为 funcName 的函数。

本质上来说，函数只是和 1、"hello" 类似的一个值，只是值的类型是函数类型。

所有的用户自定义函数，在 Go 里面实际类型均为 func(args ...interface{}) interface{}。

你可以在一个函数中引用外层函数的变量。如：

```go
x = fn(a) {
	b = 1
	y = fn(t) {
		return b + t
	}
	return y(a)
}

println(x(3)) // 结果为 4
```

但是如果你直接修改外层变量会报错：

```go
x = fn(a) {
	b = 1
	y = fn(t) {
		b = t // 这里会抛出异常，因为不能确定你是想定义一个新的 b 变量，还是要修改外层 x 函数的 b 变量
	}
	y(a)
	return b
}
```

如果你想修改外层变量，需要先引用它，如下：

```go
x = fn(a) {
	b = 1
	y = fn(t) {
		b; b = t // 现在正常了，我们知道你要修改外层的 b 变量
	}
	y(a)
	return b
}

println(x(3)) // 输出 3
```

### 不定参数

和 Go 语言类似，qlang 也支持不定参数的函数。例如内置的 max、min 都是不定参数的：

```go
a = max(1.2, 3, 5, 6) // a 的值为 float 类型的 6
b = max(1, 3, 5, 6) // b 的值为 int 类型的 6
```

也可以自定义一个不定参数的函数，如：

```go
x = fn(fmt, args...) {
	printf(fmt, args...)
}
```

这样就得到了一个 x 函数，功能和内建的 printf 函数一模一样。


## 多赋值

形式上，qlang 和 Go 语言一样支持多赋值，也支持函数返回多个值：

```go
x, y, z = 1, 2, 3.5
a, b, c = fn() {
	return 1, 2, 3.5 // 返回的是 var slice
}()
```

另外我们也支持：

```go
x, y, z = [1, 2, 3.5]
a, b, c = fn() {
	return [1, 2, 3.5] // 返回的是 float slice
}()
```

需要注意的是，带上[]进行多赋值和不带[]进行多赋值在语义上有一点点不同。下面是例子：

```go
x1, y1, z1 = 1, 2, 3.5
x2, y2, z2 = [1, 2, 3.5]
println(type(x1), type(x2))
```

结果表明 x1 的类型为 int，而 x2 的类型是 float。


## defer

是的，qlang 也支持 defer。这在处理系统资源（如文件、锁等）释放场景非常有用。一个典型场景：

```go
f, err = os.Open(fname)
if err != nil {
	// 做些出错处理
	return
}
defer f.Close()

// 正常操作这个 f 文件
```

值得注意的是：

在一个细节上 qlang 的 defer 和 Go 语言处理并不一致，那就是 defer 表达式中的变量值。在 Go 语言中，所有 defer 引用的变量均在 defer 语句时刻固定下来（如上面的 f 变量），后面任何修改均不影响 defer 语句的行为。但 qlang 是会受到影响的。例如，假设你在 defer 之后，调用 f = nil 把 f 变量改为 nil，那么后面执行 f.Close() 时就会 panic。

### 匿名函数

所谓匿名函数，是指：

```go
fn {
	... // 一段复杂代码
}
```

它等价于：

```go
fn() {
	... // 一段复杂代码
}()
```

以前在 defer 要执行一段很复杂的代码段时，我们往往这样写：

```go
defer fn() {
	... // 一段复杂代码
}()
```

有了匿名函数，我们可以简写为：

```go
defer fn {
	... // 一段复杂代码
}
```


## 类

一个用户自定义类型的基本语法如下：

```go
Foo = class {
	fn SetAB(a, b) {
		this.a, this.b = a, b
	}
	fn GetA() {
		return this.a
	}
}
```

有了这个 class Foo，我们就可以创建 Foo 类型的 object 了：

```go
foo = new Foo
foo.SetAB(3, "hello")
a = foo.GetA()
println(a) // 输出 3
```

### 构造函数

在 qlang 中，构造函数只是一个名为 _init 的成员方法（method）：

```go
Foo = class {
	fn _init(a, b) {
		this.a, this.b = a, b
	}
}
```

有了这个 class Foo 后，我们 new Foo 时就必须携带2个构造参数了：

```go
foo = new Foo(3, "hello")
println(foo.a) // 输出 3
```

## goroutine

和 Go 语言一样，qlang 中通过 go 关键字启动一个新的 goroutine。如：

```go
go println("this is a goroutine")
```

一个比较复杂的例子：

```go
wg = sync.NewWaitGroup()
wg.Add(2)

go fn {
	defer wg.Done()
	println("in goroutine1")
}

go fn {
	defer wg.Done()
	println("in goroutine2")
}

wg.Wait()
```

这是一个经典的 goroutine 使用场景，把一个 task 分为 2 个子 task，交给 2 个 goroutine 执行。


## include

在 qlang 中，一个 .ql 文件可以通过 include 文法来将另一个 .ql 的内容包含进来。所谓包含，其实际的能力类似于将代码拷贝粘贴过来。例如，在某个目录下有 a.ql 和 b.ql 两个文件。

其中 a.ql 内容如下：


```go
println("in script A")

foo = fn() {
	println("in func foo:", a, b)
}
```

其中 b.ql 内容如下：

```go
a = 1
b = 2

include "a.ql"

println("in script B")
foo()
```

如果 include 语句的文件名不是以 .ql 为后缀，那么 qlang 会认为这是一个目录名，并为其补上 "/main.ql" 后缀。也就是说：

```go
include "foo/bar.v1"
```

等价于：

```go
include "foo/bar.v1/main.ql"
```

## 模块及 import

在 qlang 中，模块(module)是一个目录，该目录下要求有一个名为 main.ql 的文件。模块中的标识(ident)默认都是私有的。想要导出一个标识(ident)，需要用 export 语法。例如：

```go
a = 1
b = 2

println("in script A")

f = fn() {
	println("in func foo:", a, b)
}

export a, f
```

这个模块导出了两个标识(ident)：整型变量 a 和 函数 f。

要引用这个模块，我们需要用 import 文法：

```go
import "foo/bar.v1"
import "foo/bar.v1" as bar2

bar.a = 100 // 将 bar.a 值设置为 100
println(bar.a, bar2.a) // bar.a, bar2.a 的值现在都是 100

bar.f()
```

qlang 会在环境变量 `QLANG_PATH` 指示的目录列表中查找 `foo/bar.v1/main.ql` 文件。如个没有设置环境变量 `QLANG_PATH`，则会在 `~/qlang` 目录中查找。

将一个模块 import 多次并不会出现什么问题，事实上第二次导入不会发生什么，只是增加了一个别名。


### include vs. import

include 是拷贝粘贴，比较适合用于模块内的内容组织。比如一个模块比较复杂，全部写在 main.ql 文件中过于冗长，则可以用 include 语句分解到多个文件中。include 不会通过 `QLANG_PATH` 来找文件，它永远基于 `__dir__`(即 include 代码所在脚本的目录) 来定位文件。

import 是模块引用，适合用于作为业务分解的主要方式。import 基于 `QLANG_PATH` 这个环境变量搜寻被引用的模块，而不是基于 `__dir__`。


## 与 Go 语言的互操作性

qlang 是一个嵌入式语言，它的定位是作为 Go 语言应用的运行时嵌入脚本。

作为 Go 语言的伴生语言，它与 Go 语言有极佳的互操作性。任何 Go 语言的函数，可以几乎不做任何包装就可以直接在 qlang 中使用。

这太爽了！


### 定制 qlang

除了 qlang 语言的 import 支持外，qlang 的 Go 语言开发包也支持 Go package 编写 qlang 模块。

qlang 采用微内核设计，大部分你看到的功能，都通过 Go package 形式编写的 qlang 模块提供。你可以按需定制 qlang。

你可以自由定制你想要的 qlang 的样子。在没有引入任何模块的情况下，qlang 连最基本的 '+'、'-'、'*'、'/' 都做不了，因为提供这个能力的是 builtin 包。

在前面“[快速入门](README.md#快速入门)”给出的精简版本基础上，我们可以自由添加各种模块，如：

```go
import (
	"qlang.io/lib/math"
	"qlang.io/lib/strconv"
	"qlang.io/lib/strings"
	...
)

func main() {

	qlang.Import("math", math.Exports)
	qlang.Import("strconv", strconv.Exports)
	qlang.Import("strings", strings.Exports)

	...
}
```

这样，在 qlang 中就可以用 math.sin, strconv.itoa 等函数了。

如果你嫌 math.sin 太长，还可以将 math 模块作为 builtin 功能导入。这只需要略微修改下导入的文法：

```go
qlang.Import("", math.Exports) // 如此，你就可以直接用 sin 而不是 math.sin 了
```

### 制作 qlang 模块

制作 qlang 模块的成本极其低廉。我们打开 `qlang.io/lib/strings` 看看它是什么样的：

```go
package strings

import (
	"strings"
)

var Exports = map[string]interface{}{
	"Contains":  strings.Contains,
	"Index":     strings.Index,
	"IndexAny":  strings.IndexAny,
	"Join":      strings.Join,
	"Title":     strings.Title,
	"ToLower":   strings.ToLower,
	"ToTitle":   strings.ToTitle,
	"ToUpper":   strings.ToUpper,
	"Trim":      strings.Trim,

	"NewReader":   strings.NewReader,
	"NewReplacer": strings.NewReplacer,
	...
}
```

值得注意的一个细节是，我们几乎不需要对 Go 语言的 strings package 中的函数进行任何包装，你只需要把这个函数加入到导出表（Exports）即可。你也无需包装 Go 语言中的类，比如上面的我们导出了 strings.NewReplacer，但是我们不必去包装 strings.Replacer 类。这个类的所有功能可以直接使用。如：

```go
strings.NewReplacer("?", "!").Replace("hello, world???") // 得到 "hello, world!!!"
```

这是 qlang 最强大的地方，近乎免包装。甚至，你可以写一个自动的 Go package 转 qlang 模块的工具，找到 Go package 所有导出的全局函数，加入到 Exports 表即完成了该 Go package 的包装，几乎零成本。

### 导出 Go 结构体

假设我们有 Go 结构体如下：

```go
package foo

type Bar struct {
	X int
	Y string
}

func (p *Bar) GetX() int {
	return p.X
}
```

我们一样可以把它加入到 Export 表进行导出：

```go
import (
	qlang "qlang.io/spec"
)

var Exports = map[string]interface{}{
	"Bar": qlang.StructOf((*foo.Bar)(nil)),
}
```

这样你就可以在 qlang 里面使用这个结构体：

```go
bar = &foo.Bar{X: 1, Y: "hello"}
x = bar.GetX()
y = bar.Y
println(x, y)
```

## 反射

在任何时候，你都可以用 type 函数来查看一个变量的实际类型，结果在 Go 语言中是 reflect.Type。如：

```go
t1 = type(1) // 相当于调用 Go 语言中的 reflect.TypeOf
```

用 type 可以很好地研究 qlang 的内在实现。比如：

```go
t2 = type(fn() {})
```

我们得到了 *Function。这说明尽管用户自定义的函数原型多样，但是其 Go 类型是一致的。

我们也可以看看用户自定义的类型：

```go
Foo = class { fn f() {} }
t1 = type(Foo)
t2 = type(Foo.f)

foo = new Foo
t3 = type(foo)
t4 = type(foo.f)
```

可以看到，class Foo 的 Go 类型是 *Class，而 object foo 的 Go 类型是 *Object。而 Foo.f 和普通用户自定义函数一致，也是 *Function，但 foo.f 不一样，它是 *Method 类型。

# 附录

## 样例代码

### 求最大素数

输入 n，求 < n 的最大素数。用法：

* `qlang maxprime.ql <N>`

```go
primes = [2, 3]
n = 1
limit = 9

isPrime = fn(v) {
	for i = 0; i < n; i++ {
		if v % primes[i] == 0 {
			return false
		}
	}
	return true
}

listPrimes = fn(max) {

	v = 5
	for {
		for v < limit {
			if isPrime(v) {
				primes = append(primes, v)
				if v * v >= max {
					return
				}
			}
			v += 2
		}
		v += 2
		n; n++
		limit = primes[n] * primes[n]
	}
}

maxPrimeOf = fn(max) {

	if max % 2 == 0 {
		max--
	}

	listPrimes(max)
	n; n = len(primes)

	for {
		if isPrime(max) {
			return max
		}
		max -= 2
	}
}

// Usage: maxprime <Value>
//
if len(os.Args) < 2 {
	fprintln(os.Stderr, "Usage: maxprime <Value>")
	return
}

max, err = strconv.ParseInt(os.Args[1], 10, 64)
if err != nil {
	fprintln(os.Stderr, err)
	return 1
}
if max < 8 { // <8 的情况下，可直接建表，答案略
	return
}

max--
v = maxPrimeOf(max)
println(v)
```

### 计算器

实现一个支持四则运算及函数调用的计算器：

```go

grammar = `

term = factor *('*' factor/mul | '/' factor/quo | '%' factor/mod)

doc = term *('+' term/add | '-' term/sub)

factor =
	FLOAT/pushFloat |
	'-' factor/neg |
	'(' doc ')' |
	(IDENT '(' doc %= ','/ARITY ')')/call
`

fntable = nil

Stack = class {

	fn _init() {
		this.stk = []
	}

	fn clear() {
		this.stk = this.stk[:0]
	}

	fn pop() {
		n = len(this.stk)
		if n > 0 {
			v = this.stk[n-1]
			this.stk = this.stk[:n-1]
			return v, true
		}
		return nil, false
	}

	fn push(v) {
		this.stk = append(this.stk, v)
	}

	fn popArgs(arity) {
		n = len(this.stk)
		if n < arity {
			panic("Stack.popArgs: unexpected")
		}
		args = make([]var, arity)
		copy(args, this.stk[n-arity:])
		this.stk = this.stk[:n-arity]
		return args
	}
}

Calculator = class {

	fn _init() {
		this.stk = new Stack
	}

	fn grammar() {
		return grammar
	}

	fn stack() {
		return this.stk
	}

	fn fntable() {
		return fntable
	}

	fn ret() {
		v, _ = this.stk.pop()
		this.stk.clear()
		return v
	}

	fn call(name) {
		f = fntable[name]
		if f == undefined {
			panic("function not found: " + name)
		}
		arity, _ = this.stk.pop()
		args = this.stk.popArgs(arity)
		ret = f(args...)
		this.stk.push(ret)
	}
}

fntable = {
	"sin": sin,
	"cos": cos,
	"pow": pow,
	"max": max,
	"min": min,

	"$mul": fn(a, b) { return a*b },
	"$quo": fn(a, b) { return a/b },
	"$mod": fn(a, b) { return a%b },
	"$add": fn(a, b) { return a+b },
	"$sub": fn(a, b) { return a-b },
	"$neg": fn(a) { return -a },

	"$call": Calculator.call,
	"$pushFloat": Stack.push,
	"$ARITY": Stack.push,
}

main { // 使用main关键字将主程序括起来，是为了避免其中用的局部变量比如 err 对其他函数造成影响

	calc = new Calculator
	engine, err = interpreter(calc, nil)
	if err != nil {
		fprintln(os.Stderr, err)
		return 1
	}

	scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line = strings.Trim(scanner.Text(), " \t\r\n")
		if line != "" {
			err = engine.Eval(line)
			if err != nil {
				fprintln(os.Stderr, err)
			} else {
				printf("> %v\n\n", calc.ret())
			}
		}
	}
}
```
