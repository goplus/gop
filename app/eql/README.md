Embedded qlang (eql)
========

eql 全称 embedded qlang，是类似 erubis/erb 的东西。结合 go generate 可很方便地让 Go 支持模板（不是 html template，是指语言特性中的泛型）。

## eql 程序

命令行：

```
eql <template_file> [-o <output_file>] [--key1=val1 --key2=val2 ...]
eql <template_dir> [-o <output_dir>] [--key1=val1 --key2=val2 ...]
```

其中

* `<template_file>`: 要解析的 eql 文件，也就是模板文件。
* `<template_dir>`: 要解析的 template package，也就是整个目录是一个模板。
* `<output_file>`: 要生成的渲染后的文件。如果没有指定则为 stdout。
* `<output_dir>`: 要生成的渲染后的目标目录。如果没有指定则为对 `<template_dir>` 进行渲染后的值。
* `--key1=val1 --key2=val2 ...`: 渲染涉及到的模板变量的值。

## 语法

### 插入 qlang 代码

```
<%
    // 在此插入 qlang 代码
%>
```

### 输出 qlang 表达式

```
<%= qlang_expr %>
```

你可以理解为这只是插入 qlang 代码的一种简写手法。它等价于：

```
<% print(qlang_expr) %>
```

### 输出一个变量

```
$var
```

你可以理解为这只是插入 qlang 代码的一种简写手法。它等价于：

```
<% print(var) %>
```

特别地，我们用 `$$` 表示普通字符 `$`。也就是说：

```
$$
```

等价于：

```
<% print('$') %>
```

## 用 eql 实现 Go 对泛型的支持

我们举例说明。假设我们现在实现了一个 Go 的模板类，文件名为 `example.eql`，内容如下：

```go
package eql_test

import (
	<%= eql.imports(imports) %>
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

### 样例: template package

* [$set.v1](example/$set.v1): 以集合类为例展示如何构建一个 template package。
