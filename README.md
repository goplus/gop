The Q Language (Q语言)
========

# 下载

### 源代码

```
go get -u qlang.io/qlang
```

### 二进制(qlang.v1)

* linux/amd64: [qnlang-linux-amd64.tar.gz](http://open.qiniudn.com/qnlang-linux-amd64.tar.gz)
* Mac OS: [qnlang-darwin-amd64.tar.gz](http://open.qiniudn.com/qnlang-darwin-amd64.tar.gz)

### 二进制(qlang.v2)

* 性能比qlang.v2更优。暂未提供二进制下载包，请通过源代码自行编译。


# 语言特色

* 最大卖点：与 Go 语言有最好的互操作性。所有 Go 语言的社区资源可以直接为我所用。
* 有赖 Go 语言的互操作性，这门语言不需要自己实现标准库。尽管年轻，但是这门语言已经具备商用的可行性。
* 微内核：语言的核心只有 1200 行代码。所有功能以可插拔的 module 方式提供。

预期的商业场景：

* 由于与 Go 语言的无缝配合，qlang 在嵌入式脚本领域有 lua、python、javascript 所不能比拟的优越性。
* 比如：网络游戏中取代 lua 的位置。


## 样例

### 最大素数

输入 n，求 < n 的最大素数：

* [maxprime.ql](https://github.com/qiniu/qlang/blob/develop/app/qlang/maxprime.ql)

用法：

```
qlang maxprime.ql <N>
```

### 计算器

实现一个支持四则运算及函数调用的计算器：

* [calc.ql](https://github.com/qiniu/qlang/blob/develop/app/qlang/calc.ql)

用法：

```
qlang calc.ql
```

### qlang自举

qlang的自举（用qlang实现一个qlang）：

* [qlang.ql](https://github.com/qiniu/qlang/blob/develop/app/qlang/qlang.ql)

交互模式跑 qlang 版本的 qlang（可以认为是上面计算器的增强版本）：

```
qlang.v1 qlang.ql  #目前暂时只能用qlang.v1版本完成自举
```

当然你还可以用 qlang 版本的 qlang 来跑最大素数问题：

```
qlang.v1 qlang.ql maxprime.ql <N>
```

# 使用说明

## 运算符

基本上除了位运算：'&'、'|'、'>>'、'<<' 和 chan 操作 '<-' 之外，Go 语言的操作符都支持。包括：

* '+'、'-'、'*'、'/'、'%'、'='
* '+='、'-='、'*='、'/='、'%='、'++'、'--'
* '!'、'>='、'<='、'>'、'<'、'=='、'!='、'&&'、'||'


## 类型

原理上支持所有 Go 语言中的类型。典型有：

* 基本类型：int、float (在 Go 语言里面是 float64)、string、byte、bool、var（在 Go 语言里面是 interface{}）。
* 复合类型：slice、map
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
e = slice("int", len, cap) // 创建一个 int slice，并将长度设置为 len，容量设置为 cap
f = slice(type(e), len, cap) // 创建一个 int slice 的 slice，也就是 Go 语言里面的 [][]int
```

和 Go 语言类似，slice 有如下内置的操作：

```go
a = append(a, 4, 5, 6) // 含义与 Go 语言完全一致
n = len(a) // 取 a 的元素个数
b1 = b[2] // 取 b 这个 slice 的第二个元素
set(b, 2, 888) // 设置 b 这个 slice 的第二个元素的值为 888。在 Go 语言里面是 b[2] = 888
set(b, 1, 777, 2, 888, 3, 999) // Go 里面是：b[1], b[2], b[3] = 777, 888, 999
b2 = b[1:4] // Go 里面是： b2 = b[1:4]
```

特别地，在 qlang 中可以这样赋值：

```go
x, y, z = [1, 2, 3]
```

结果是 x = 1, y = 2, z = 3。这是 qlang 和 Go 语言的基础设计不同导致的：

* qlang 不支持多返回值。对于那些返回了多个值的 Go 函数，在 qlang 会理解为返回 var slice，也就是 []interface{}。

举个例子：

```go
f, err = os.open(fname)
```

这个例子，在 Go 里面返回的是 (*os.File, error)。但是 qlang 中是 var slice。

### map 类型

```go
a = {"a": 1, "b": 2, "c": 3} // 得到 map[string]int 类型的对象
b = {"a": 1, "b", 2.3, "c": 3} // 得到 map[string]float64 类型的对象
c = {1: "a", 2: "b", 3: "c"} // 得到 map[int]string 类型的对象
d = {"a": "hello", "b": 2.0, "c": true} // 得到 map[string]interface{} 类型的对象
e = map("string:int") // 创建一个空的 map[string]int 类型的对象
f = map(mapOf("string", type(e))) // 创建一个 map[string]map[string]int 类型的对象
```

和 Go 语言类似，map 有如下内置的操作：

```go
n = len(a) // 取 a 的元素个数
x = a["b"] // 取 a map 中 key 为 "b" 的元素
x = a.b // 含义同上
set(a, "e", 4, "f", 5, "g", 6) // 在 Go 语言里面是 a["e"], b["f"], b["g"] = 4, 5, 6
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

## 类型转换

### 自动类型转换

大部分情况下，我们不会自动进行类型转换。一些例外是：

* 如果一个函数接受的是 float，但是传入的是 int，会进行自动类型转换。

### 强制类型转换

```go
a = int('a') // 强制将 byte 类型转为 int 类型
b = float(b) // 强制将 int 类型转为 float 类型
c = string('a') // 强制将 byte 类型转为 string 类型
```

## 流程控制

### if语句

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

需要注意的是，if 语句是有值的。比如：

```go
x = if a < b { a } else { b } // x 取 a 和 b 两者中的小值。即 x = min(a, b)
```

如果你希望使用 if 表达式的值，建议不要写成多行：

```go
x = if a < b {
	a
} else {
	b
}
```

这段代码不会如你所愿工作。原因是编译器会在行末自动插入 ';'，所以相当于是：

```go
x = if a < b { a; } else { b; }
```

结果是 `x = nil`。

### switch语句

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

### for语句

除了不支持 for range 文法，也不支持中途 continue、break（但是支持 return）。其他和 Go 语言完全类似：

```go
for { // 无限循环，需要在中间 return，或者 os.exit(code)，否则不能退出
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

### 私有变量

在 qlang 中，我们引入了私有变量的概念。

私有变量以 _ 开头来标识。私有变量只能在本函数内使用，并且互不干扰（没有名字冲突）。如：

```go
x = fn(a) {
	_b = 1
	y = fn(t) {
		_b = t
	}
	y(a)
	return _b
}

println(x(3)) // 输出 1，因为 y(a) 这个调用并不影响 x 函数中的 _b 变量
```

### 可变参数

如 Go 语言类似，qlang 也支持可变参数的函数。例如内置的 max、min 都是可变参数的：

```go
a = max(1.2, 3, 5, 6) // a 的值为 float 类型的 6
b = max(1, 3, 5, 6) // b 的值为 int 类型的 6
```

也可以自定义一个可变参数的函数，如：

```go
x = fn(fmt, args...) {
	printf(fmt, args...)
}
```

这样就得到了一个 x 函数，功能和内建的 printf 函数一模一样。


## defer

是的，qlang 也支持 defer。这在处理系统资源（如文件、锁等）释放场景非常有用。一个典型场景：

```go
f, err = os.open(fname)
if err != nil {
	// 做些出错处理
	return
}
defer f.close()

// 正常操作这个 f 文件
```

值得注意的是：

在一个细节上 qlang 的 defer 和 Go 语言处理并不一致，那就是 defer 表达式中的变量值。在 Go 语言中，所有 defer 引用的变量均在 defer 语句时刻固定下来（如上面的 f 变量），后面任何修改均不影响 defer 语句的行为。但 qlang 是会受到影响的。例如，假设你在 defer 之后，调用 f = nil 把 f 变量改为 nil，那么后面执行 f.close() 时就会 panic。


## 类

一个用户自定义类型的基本语法如下：

```go
Foo = class {
	fn setAB(a, b) {
		set(this, "a", a, "b", b)
	}
	fn getA() {
		return this.a
	}
}
```

有了这个 class Foo，我们就可以创建 Foo 类型的 object 了：

```go
foo = new Foo
foo.setAB(3, "hello")
a = foo.getA()
println(a) // 输出 3
```

### 构造函数

在 qlang 中，构造函数只是一个名为 _init 的成员方法（method）：

```go
Foo = class {
	fn _init(a, b) {
		set(this, "a", a, "b", b)
	}
}
```

有了这个 class Foo 后，我们 new Foo 时就必须携带2个构造参数了：

```go
foo = new Foo(3, "hello")
println(foo.a) // 输出 3
```


## 模块及模块import

在 qlang 中的确有模块的概念，但是并不像大部分脚本语言那样，有模块的 import 语法。

要理解这一点，最关键是理解 qlang 的核心理念：

* qlang 是一个嵌入式语言，它的定位是作为 Go 语言应用的运行时嵌入脚本。
* 我们鼓励依据应用按需求对 qlang 进行裁剪。这也是 qlang 采用微内核设计的原因。

看起来 qlang 将自己的使用范围收得非常窄，但考虑到 Go 语言在未来的霸主地位，qlang 有着广阔的应用空间。

而在 Go 语言应用的嵌入脚本领域而言，qlang 有无可比拟的核心优势：

* 任何 Go 语言的函数，都可以不做任何包装直接在 qlang 中使用

这太爽了！


### 定制qlang

尽管 qlang 本身没有 import，但是 qlang 的 Go 语言包支持 import。通过 import 你可以自由定制你想要的 qlang 的样子。在没有引入任何模块的情况下，qlang 连最基本的 '+'、'-'、'*'、'/' 都做不了，因为提供这个能力的是 builtin 包。

所以一个最基础版本的 qlang 应该是这样的：

```go
import (
	"fmt"

	"qlang.io/qlang.v2/qlang"
	_ "qlang.io/qlang/builtin" // 导入 builtin 包
)

func main() {

	lang, err := qlang.New(nil) // 参数 nil 也可以改为 qlang.InsertSemis
	if err != nil {
		// 错误处理
		return
	}

	err = lang.SafeEval(`1 + 2`)
	if err != nil {
		// 错误处理
		return
	}

	v, _ := lang.Ret()
	fmt.Println(v) // 输出 3
}
```

在大部分正式的场合，qlang.New 传入 qlang.InsertSemis 会更多一些。它表示在各行的末尾智能的插入 ';'。但是在本例中我们希望执行的是一个表达式，而不是语句，所以传入 nil 参数更为合适。如果我们改为传入 qlang.InsertSemis，那么得到的结果就不再是 3，而是 nil。因为表达式 `1+2` 的结果是 3，但是表达式 `1+2;` 的结果是 nil。

有了这个基础版本以后，我们可以自由添加各种模块，如：

```go
import (
	"qlang.io/qlang/math"
	"qlang.io/qlang/strconv"
	"qlang.io/qlang/strings"
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

### 制作qlang模块

制作qlang的模块成本极其低廉。我们打开 `qlang.io/qlang/strings` 看看它是什么样的：

```go
package strings

import (
	"strings"
)

var Exports = map[string]interface{}{
	"contains":  strings.Contains,
	"index":     strings.Index,
	"indexAny":  strings.IndexAny,
	"join":      strings.Join,
	"title":     strings.Title,
	"toLower":   strings.ToLower,
	"toTitle":   strings.ToTitle,
	"toUpper":   strings.ToUpper,
	"trim":      strings.Trim,
	"reader":    strings.NewReader,
	"replacer":	 strings.NewReplacer,
	...
}
```

值得注意的一个细节是，我们并没有对 Go 语言的 strings package 进行包装。比如上面的我们导出了 strings.NewReplacer，但是我们不必去包装 strings.Replacer 类。这个类的所有功能可以直接使用。如：

```go
strings.replacer("?", "!").replace("hello, world???") // 得到 "hello, world!!!"
```

这是 qlang 最强大的地方，近乎免包装。甚至，你可以写一个自动的 Go package 转 qlang 模块的工具，找到 Go package 所有导出的全局函数，加入到 Exports 表即完成了该 Go package 的包装，几乎零成本。


## 反射

在任何时候，你都可以用 type 函数来查看一个变量的实际类型，结果在 Go 语言中是 reflect.Type。如：

```go
t1 = type(1) // 相当于调用 Go 语言中的 reflect.TypeOf
```

用 type 可以很好地研究 qlang 的内在实现。比如：

```go
t2 = type(fn() {})
```

我们得到了 *qlang.Function。这说明尽管用户自定义的函数原型多样，但是其 Go 类型是一致的。

我们也可以看看用户自定义的类型：

```go
Foo = class { fn f() {} }
t1 = type(Foo)
t2 = type(Foo.f)

foo = new Foo
t3 = type(foo)
t4 = type(foo.f)
```

可以看到，class Foo 的 Go 类型是 *qlang.Class，而 object foo 的 Go 类型是 *qlang.Object。而 Foo.f 和普通用户自定义函数一致，也是 *qlang.Function，但 foo.f 不一样，它是 *qlang.thisDref 类型。

# 附录

## 样例代码

### 求最大素数

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
if len(os.args) < 2 {
	fprintln(os.stderr, "Usage: maxprime <Value>")
	return
}

max, err = strconv.parseInt(os.args[1], 10, 64)
if err != nil {
	fprintln(os.stderr, err)
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
		set(this, "stk", [])
	}

	fn clear() {
		set(this, "stk", this.stk[:0])
	}

	fn pop() {
		n = len(this.stk)
		if n > 0 {
			v = this.stk[n-1]
			set(this, "stk", this.stk[:n-1])
			return [v, true]
		}
		return [nil, false]
	}

	fn push(v) {
		set(this, "stk", append(this.stk, v))
	}

	fn popArgs(arity) {
		n = len(this.stk)
		if n < arity {
			panic("Stack.popArgs: unexpected")
		}
		args = sliceFrom(this.stk[n-arity:n]...)
		set(this, "stk", this.stk[:n-arity])
		return args
	}
}

Calculator = class {

	fn _init() {
		set(this, "stk", new Stack)
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
		if f == nil {
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

calc = new Calculator
engine, err = interpreter(calc, nil)
if err != nil {
	fprintln(os.stderr, err)
	return 1
}

scanner = bufio.scanner(os.stdin)
for scanner.scan() {
	line = strings.trim(scanner.text(), " \t\r\n")
	if line != "" {
		err = engine.eval(line)
		if err != nil {
			fprintln(os.stderr, err)
		} else {
			printf("> %v\n\n", calc.ret())
		}
	}
}
```

