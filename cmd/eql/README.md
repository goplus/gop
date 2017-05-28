Embedded qlang (eql)
========

eql 全称 embedded qlang，是类似 erubis/erb 的东西。结合 go generate 可很方便地让 Go 支持模板（不是 html template，是指语言特性中的泛型）。

## eql 程序

命令行：

```bash
eql <templatefile> [-i -j <jsoninput> -o <outputfile> --key1=value1 --key2=value2 ...]
eql <templatedir> [-i -j <jsoninput> -o <outputdir> --key1=value1 --key2=value2 ...]
```

其中

* `<templatefile>`: 要解析的 eql 文件，也就是模板文件。
* `<templatedir>`: 要解析的 template package，也就是整个目录是一个模板。
* `-i`: 以 stdin 作为 json input。
* `-j <jsoninput>`: 输入 json input。
* `-o <outputfile>`: 要生成的渲染后的文件。如果没有指定则为 stdout。
* `-o <outputdir>`: 要生成的渲染后的目标目录。如果没有指定则为对 `<template_dir>` 进行渲染后的值。
* `--key1=val1 --key2=val2 ...`: 重载 json input 涉及到的模板变量的值。

假设我们有 a.eql，内容如下：

```
Hello, $a!
```

则我们可以有这样一些使用 `eql` 程序的方式：

```bash
eql -j '{"a":"qlang"}' a.eql
```

或者：

```bash
echo '{"a":"qlang"}' | eql -i a.eql
```

或者：

```bash
# 这种方式看起来最为简洁，只是无法输入 string 类型之外的数据
eql --a=qlang a.eql
```

## 样例

### 单文件模板

* [example.eql](example.eql): 展示 eql 语法的样例，下文将详细介绍“模板文法”。

### 目录模板

* [$set.v1](example/$set.v1): 以集合类为例展示如何构建一个 template package。下文将详细介绍“目录模板规则”。


## 模板文法

### 插入 qlang 代码

```go
<%
    // 在此插入 qlang 代码
%>
```

### 输出 qlang 表达式

```go
<%= qlang_expr %>
```

你可以理解为这只是插入 qlang 代码的一种简写手法。它等价于：

```go
<% print(qlang_expr) %>
```

### 输出一个变量

```
$var
```

你可以理解为这只是插入 qlang 代码的一种简写手法。它等价于：

```go
<% print(var) %>
```

特别地，我们用 `$$` 表示普通字符 `$`。也就是说：

```
$$
```

等价于：

```go
<% print('$') %>
```

## 目录模板规则

目录模板的根目录名本身允许是模板，但是只有在用户没有指定 `-o <output_dir>` 的情况下有效。例如上面例子中的 `$set.v1` 包含模板变量 `$set`。

目录模板的渲染过程（实例化过程）比较简单，就是对模板目录进行遍历，并执行如下规则：

* 如果遇到子目录，创建出对应子目录，并递归本过程。
* 如果是 `.eql` 后缀的文件，进行模板渲染并去掉 `.eql` 后缀进行存储。例如 `README.md.eql` 渲染后保存为 `README.md`。
* 如果是其他后缀的文件，则进行简单复制。


## eql 函数

### eql.Imports()

这个函数解析 imports 变量（通常是用户通过命令行 `--imports=package1,package2,...` 传入的）并将 `,` 分隔风格改为符合 Go 语言中 import (...) 风格的字符串。

如果 imports 变量不存在，则该函数返回 "" (空字符串)。

例如，我们假设传入的 imports 变量为 "bufio,bytes"，则 `eql.Imports()` 为如下风格的字符串：

```go
	"bufio"
	"bytes"
```

### eql.Var("varname", defaultval)

这个函数主要解决 qlang 中目前还没有支持判断一个变量是否存在的缺陷。其语义看起来是这样的：

```go
if varname != undefined {
	return varname
} else {
	return defaultval
}
```

当然目前因为在 qlang 中如果 varname 不存在就会直接报错，所以以上代码仅仅是表达 `eql.Var("varname", defaultval)` 的逻辑语义。

## 用 eql 实现 Go 对泛型的支持

我们举例说明。假设我们现在实现了一个 Go 的模板类，文件名为 `example.eql`，内容如下：

```go
package eql_test

import (
	<%= eql.Imports() %>
	"encoding/binary"
)

// -----------------------------------------------------------------------------

type $module string

func (p $module) write(out $Writer, b []byte) {

	_, err := out.Write(b)
	if err != nil {
		panic(err)
	}
}

<% if Writer == "*bytes.Buffer" { %>
func (p $module) flush(out $Writer) {
}
<% } else { %>
func (p $module) flush(out $Writer) {

	err := out.Flush()
	if err != nil {
		panic(err)
	}
}
<% } %>

// -----------------------------------------------------------------------------
```

这个模板里面，有 3 个模板变量：

* `imports`: 需要额外引入的 package 列表，用 `,` 分隔。
* `module`: 模板类的类名。
* `Writer`: 模板类的用到的参数类型。

有了这个模板，我们就可以用如下命令生成具体的类：

```
eql example.eql -o example_bytes.go --imports=bytes --module=modbytes --Writer="*bytes.Buffer"
```

这会生成 example_bytes.go 文件，内容如下：

```go
package eql_test

import (
	"bytes"
	"encoding/binary"
)

// -----------------------------------------------------------------------------

type modbytes string

func (p modbytes) write(out *bytes.Buffer, b []byte) {

	_, err := out.Write(b)
	if err != nil {
		panic(err)
	}
}

func (p modbytes) flush(out *bytes.Buffer) {
}

// -----------------------------------------------------------------------------
```

再试试换一个 Writer：

```
eql example.eql -o example_bufio.go --imports=bufio --module=modbufio --Writer="*bufio.Writer"
```

我们得到 example_bufio.go，内容如下：

```go
package eql_test

import (
	"bufio"
	"encoding/binary"
)

// -----------------------------------------------------------------------------

type modbufio string

func (p modbufio) write(out *bufio.Writer, b []byte) {

	_, err := out.Write(b)
	if err != nil {
		panic(err)
	}
}

func (p modbufio) flush(out *bufio.Writer) {

	err := out.Flush()
	if err != nil {
		panic(err)
	}
}

// -----------------------------------------------------------------------------
```

### 结合 go generate

结合 go generate 工具，我们就可以很好地支持 Go 泛型了。

例如假设我们在 foo.go 里面引用了 Writer = `*bufio.Writer` 版本的实现，则只需要在 foo.go 文件中插入以下代码：

```go
//go:generate eql example.eql -o example_bufio.go --imports=bufio --module=module --Writer=*bufio.Writer
```

如此，你只需要在 foo.go 所在的目录执行 go generate 就可以生成 example_bufio.go 文件了。


## 以库的方式使用 eql

### 初始化

```go
ql = qlang.New()
eql = eqlang.New(ql)

ql.ResetVars(map[string]var{"name": "qlang"}) // 重置 ql 对象的变量集合
ql.SetVar("a", 123)
```

### eql.Subst("template text")

替换 $varname 为对应的值。例如：

```go
output = eql.Subst("Hello, $name!")
println(output)
```

得到的结果是：

```
Hello, qlang!
```

### eql.Execute("template script", fname, output)

执行 eql 模板。

* fname:  为模板文件名。
* output: 为输出文件，为 "" 表示输出到 stdout。

例如：

```go
eql.Execute("Hello, <%= strings.ToUpper(name) %>!\n", "", "")
```

输出为：

```
Hello, QLANG!
```

### eql.ExecuteFile(source, output)

执行 eql 模板文件。

* source: 为 eql 模板文件。
* output: 为输出文件，为 "" 表示输出到 stdout。

### eql.ExecuteDir(vars, source, output)

执行 eql 模板目录。

* vars:   是模板依赖的变量集。
* source: 为 eql 模板目录。
* output: 为输出目录，为 "" 表示输出的目录名是以 source 为模板的渲染结果。
