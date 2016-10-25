tpl
===

TPL 全称是 `Text Processing Language`（文本处理语言）。整个库的组成如下：

* `qiniupkg.com/text/tpl.v1`: 基础的 TPL 文本处理引擎，外加一个 TPL 文法编译器。
* `qiniupkg.com/text/tpl.v1/interpreter`: 在 TPL 基础上实现的解释器引擎。典型使用场景是实现一个计算器（Calculator）、解释型的语言执行器（直接产生执行结果）或翻译器（翻译成另一种语言，比如字节码）。
* `qiniupkg.com/text/tpl.v1/{number, cmplx, rat}`: 在解释器框架基础上实现的3个计算器，分别对应三种类型的运算：`float64`（浮点数）、`complex128`（复数）、`*big.Rat`（有理数）。
* `qiniupkg.com/text/tpl.v1/exmples`: TPL 库的一些样例。


## 文本处理引擎

TPL 把文本处理分为了两个阶段。

第一个阶段是词法分析，由 Tokener（也叫 Scanner）完成。Tokener/Scanner 将文本流转换为 Token 序列。简单来说，就是一个 `text []byte => tokens []Token` 的过程。尽管世上语言多样，但是词法非常接近，所以这块 TPL 只是抽象了一个 Tokener 接口，方便用户自定义。TPL 也内置了一个与 Go 语言词法完全类似的 Scanner（做了一个非常细微的调整，就是增加了 '?'、'~'、'@' 等操作符）。

第二个阶段是语法分析。由于 Go 语言不支持操作符重载，语法表示我们用文本来表达，通过 TPL 文法编译器编译成语法 DOM 树。当前我们支持的 TPL 文法如下：

### 基本token

基本token的表示方式又分为以下几类：

* TOKEN。通常以全大些的字母表达。如：IDENT、INT、FLOAT、IMAG、STRING、CHAR、COMMENT 等。
* 单字符操作符。如：'+'、'-'、'*'、'/' 等等。这些其实也可以用上面 TOKEN 形式表达，只是不那么直观。如 '+' 可以用 `ADD`，'-' 可以用 `SUB` 等等。
* 多字符操作符。如: "++"、"--"、"+="、"<=" 等等。这些其实也可以用上面 TOKEN 形式表达，只是不那么直观。"++" 可以用 `INC`，"+=" 可以用 `ADD_ASSIGN` 等等。
* 关键字。如："if"、"return"、"class" 等等。TPL 内建的 AutoKwScanner 可以自动识别语言的关键字。其原理是，只要在 TPL 文法中出现了满足 `'"' IDENT '"'` 形式的规则，那么就认为里面 IDENT 就是一个关键字。

### 特殊token

* `EOF`：即 tpl 中的 `GrEOF` 规则。不匹配任何内容，但是可用于判断当前是否已经匹配完整个 token 序列。
* `1`：即 tpl 中的 `GrTrue` 规则。不匹配任何内容，无论如何永远匹配成功。


### 复合规则

* `*G`: 反复匹配规则 G，直到无法成功匹配为止。
* `+G`: 反复匹配规则 G，直到无法成功匹配为止。要求至少匹配成功1次。
* `?G`: 匹配规则 G 1次或0次。
* `@G`: 要求其后的文本满足规则 G。如果要匹配的文本满足规则 G，则匹配成功，但是不真匹配任何内容（相当于只是 peek 下后续的 token 序列，要匹配的文本的当前匹配位置不前进）。例如：`@'{'` 这样一段规则表示的含义是要求要匹配的后续文本应该以 { 开头。
* `~G`: 要求其后的文本不满足规则 G。如果要匹配的文本满足规则 G，则匹配失败；如果不满足 G，则匹配成功，但不匹配任何内容（相当于只是 peek 下后续的 token 序列，要匹配的文本的当前匹配位置不前进）。例如: `~'{'` 这样一段规则表示的含义是要求要匹配的后续文本不能以 { 开头。
* `G1 G2 ... Gn`: 要求要匹配的文本满足规则序列 G1 G2 ... Gn。例如：`"fn" '(' ')'` 表示要匹配的文本去除所有空白和注释后是 `fn()` 这样的文本。
* `G1 G2! ... Gn`: 可以在 G1 G2 ... Gn 中间的任何地方插入 ! 操纵符。它表示的意思是不可回退，或者说快速失败（fail fast）。例如 `G1 G2! ... Gn` 表达的含义是，如果 G1 G2 匹配成功了，那么后续的 G3 ... Gn 必须匹配成功，如果不成功则直接结束整个匹配报告失败。善用 ! 操作符可以改善错误提示信息，因为通常这时候提示的错误更为准确。
* `G1 | G2 | ... | Gn` 要求要匹配的文本满足规则 G1 G2 ... Gn 中的任意一个。
* `G1 % G2`: 列表运算。从规则上来说等价于 `G1 *(G2 G1)`。比如 `IDENT % ','` 可以匹配 `abc`, `abc, defg, fhi` 这样的文本。
* `G1 %= G2`: 可选列表运算。从规则上来说等价于 `?(G1 % G2)`。和列表运算唯一差别是允许匹配的内容为空。
* `(G)`: 从规则上来说就是 G，只是改变运算的次序。详细见下“运算符优先级”。

运算符按优先次序，从高到低排序如下：

* `(G)`
* `*G`, `+G`, `?G`, `@G`, `~G`: 这些操作遵循右结合律，在最右边的运算优先。如：`~+G` 表示 `~(+G)`
* `G1 % G2`, `G1 %= G2`
* `G1 G2! ... Gn`
* `G1 | G2 | ... | Gn`

### 动作和标记

动作(action) 是指规则匹配成功的情况下执行的代码。如下：

```go
tpl.Action(G, func(tokens []tpl.Token, g tpl.Grammar) {
	...
})
```

这里的 func 就是动作(action)，整个表达式得到一个带动作的规则，使用上和一般的规则无异。

但在 TPL 文法毕竟无法内嵌 Go 语言代码，我们并无法直接内嵌一个动作(action)进去，取而代之的是标记(mark)。

标记的文法是这样的：

```
G/mark
```

这里 G 代表一个规则，而 mark 是一个合法的符号(IDENT)。在 TPL 文法中遇到标记(mark)时，TPL 编译器会产生一个回调(Marker)，交给业务方来处理这个标记。如下：

```go
Marker := func(g tpl.Grammar, mark string) tpl.Grammar {
	...
}
```

Maker 概念上并不是动作。但是你可以在 Maker 回调中生成相应的动作，并且通常你都在这样做。所以一般我们在沟通惯例上，会简单把标记(mark)和执行动作(action)等同起来。

通过定制 Marker，业务方就可以完成自己期望的业务逻辑。我们也有一些内建的 Marker 实现，进一步简化大家的文本处理过程。例如，我们接下来要介绍的“解释器(interpreter)”。

### 使用TPL

使用 TPL 的范式如下：

```go
import (
	"qiniupkg.com/text/tpl.v1"
)

// 定义要处理的文本内容对应的TPL文法
//
const grammar = `
	...
`

func eval(text []byte, fname string) (..., err error) {

	defer func() {
		if e := recover(); e != nil {
			// 错误处理
			return
		}
	}()

	marker := func(g tpl.Grammar, mark string) tpl.Grammar {
		...
	}
	compiler := &tpl.Compiler{
		Grammar: grammar,
		Marker:  marker,
	}
	m, err := compiler.Cl()
	if err != nil {
		// 错误处理
		return
	}
	err = m.MatchExactly(texr, fname)
	if err != nil {
		// 错误处理
		return
	}

	// ...
	return
}
```

## 解释器(interpreter)

解析器引擎(interpreter engine)在文本处理的模型上做出了更多的假设。它假设业务方实现如下接口：

```go
type Stack interface {
	PopArgs(arity int) (args []reflect.Value, ok bool)
	PushRet(ret []reflect.Value) error
}

type Interface interface {
	Grammar() string
	Fntable() map[string]interface{}
	Stack() Stack
}
```

也就是说：

* 要有一个栈(stack)；
* 要有一个函数表(fntable)；

在 解释器(interpreter) 的 TPL 文法中，标记(mark)分为如下这些情况：

### _mute

`_mute` 是一个内建的标记(mark)。顾名思义，它有禁止发言(禁止执行动作)的意思。展开来说：

* 在第一次遇到 `_mute` 时，会禁用后续普通 mark 的执行，但 '_' 开头的 mark 不受影响。
* 后续如果再遇到 `_mute` 时，mute 引用计数++，当 mute 计数 > 1 时，所有非内建的 mark 都会禁止执行（包括 '_' 开头的）。

### _unmute

`_unmute` 是一个内建的标记(mark)。它是 _mute 的反操作，执行 mute 引用计数-- 的行为。当 mute 计数减少到 0 时，所有 mark 回到正常执行状态。

### _tr

_tr 是一个内建的标记(mark)。它是一个调试用的标记，它在规则匹配成功时，将规则所匹配的文本打印出来。

### 用户定义标记

解释器(interpreter) 在遇到用户定义标记时，会查找 fntable 得到对应的函数。例如，假设用户自定义标记叫 `add`，那么我们会到 fntable 中查找名为 `$add` 的函数。如果找到，则：

先通过反射查看函数的第一个参数。这分为 3 种情况：

1) 第一个参数是 interpreter.Stack 类型。会自动传入 interpreter 的 stack 实例。

2) 第一个参数是 interpreter.Interface 类型。会自动传入用户实现的 interpreter 实例。

3) 第一个参数是其他类型。

对于情形1和2，我们要求函数的参数只能是1个或2个，返回值要么没有，要么error类型。对于函数只有1个参数的情形，我们传入 stack 或 interpreter 实例然后调用之；对于函数有2个参数的情形，我们依据参数的类型分为如下几种情况：

* 参数为 interface{} 类型。这表示这个函数希望传入规则匹配到的 tokens 序列（tokens []tpl.Token）。
* 参数为 interpreter.Engine 类型。这表示这个函数希望传入解释器引擎(interpreter engine)。
* 参数为其他类型，这时也有两种情形。一种是 mark 以大写字母开头（或者以 _ + 大写字母开头开头），则表示函数接受的是 Grammar.Len()，通常是当规则为 *G、+G、?G、G1 % G2、G1 %= G2 时告知你准确的成功匹配次数，此时函数第二个参数必须是 int 类型。另一种情况是小写开头（或者以 _ + 小写字母开头），表示函数希望传入规则匹配到的 tokens 序列的第一个，也就是 tokens[0] 对应的值（可能会依据类型的不同进行自动的类型转换，比如如果参数为 float64，那么我们会自动调用 strconv.ParseFloat 完成转换）。

对于情形3，我们自动根据所需的参数个数，从 stack 中弹出参数列表（通过调用 PopArgs），然后调用该函数，把返回的结果压回 stack（通过调用 PushRet）。

### 使用解释器

使用 interpreter 的范式如下：

第一步，先实作一个 interpreter 包（假设叫 foo，我们的样例 qiniupkg.com/text/tpl.v1/{number, cmplx, rat} 都属于这一类）：

```go
package foo

import (
	"reflect"
	"qiniupkg.com/text/tpl.v1/interpreter.util"
)

// 定义要处理的文本内容对应的TPL文法
//
const grammar = `
	...
`

// ----------------------------------------------------

type Stack struct {

}

func (p *Stack) PopArgs(arity int) (args []reflect.Value, ok bool) {
	...
}

func (p *Stack) PushRet(ret []reflect.Value) error {
	...
}

// ----------------------------------------------------

type Calculator struct {
	stk *Stack
}

func (p *Calculator) Grammar() string {
	return grammar
}

func (p *Calculator) Stack() interpreter.Stack {
	return p.stk
}

func (p *Calculator) Fntable() map[string]interface{} {
	return fntable
}

func init() {
	fntable = Fntable
}

var fntable map[string]interface{}

// ----------------------------------------------------

var Fntable = map[string]interface{}{
	...
}
```

然后引用这个 interpreter 实作：

```go
import (
	"foo"

	"qiniupkg.com/text/tpl.v1/interpreter"
)

func eval(text []byte, fname string) (..., err error) {

	defer func() {
		if e := recover(); e != nil {
			// 错误处理
			return
		}
	}()

	calc := foo.New()
	engine, err := interpreter.New(calc, nil)
	if err != nil {
		// 错误处理
		return
	}

	err = engine.MatchExtractly(text, fname)
	if err != nil {
		// 错误处理
		return
	}

	// ...
	return
}
```

