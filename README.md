Q Language - A script language for Go
========

[![Build Status](https://travis-ci.org/qiniu/qlang.png?branch=develop)](https://travis-ci.org/qiniu/qlang)
[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/qlang)](https://goreportcard.com/report/github.com/qiniu/qlang)

![logo](http://qiniutek.com/images/logo-2.png)

## 语言特色

* 与 Go 语言有最好的互操作性。可不进行任何包装即可直接使用 Go 语言的函数、类及其成员变量和方法。
* 有赖于 Go 语言的互操作性，这门语言直接拥有了一套非常完整且您十分熟悉的标准库，无额外学习成本。
* 与 Go 十分相近的语法，降低您的学习成本。
* 支持 Go 绝大部分语言特性，包括：for range, string, slice, map, chan, goroutine, defer, etc。
* 微内核：语言的核心只有大约 1200 行代码。所有功能以可插拔的 module 方式提供。

预期的商业场景：

* 由于与 Go 语言的无缝配合，qlang 在嵌入式脚本领域有 lua、python、javascript 所不能比拟的优越性。比如：网络游戏中取代 lua 的位置。
* 作为编译原理的教学语言。由于 qlang 的 Compiler 代码极短，便于阅读和理解，非常方便教学实践之用。


## 快速入门

### 在您的 Go 代码中整合 qlang

```go
import (
	"fmt"

	"qlang.io/cl/qlang"
	_ "qlang.io/lib/builtin" // 导入 builtin 包
)

const scriptCode = `x = 1 + 2`

func main() {

	ql := qlang.New()
	err := ql.SafeExec([]byte(scriptCode), "")
	if err != nil {
		// 错误处理
		return
	}

	fmt.Println("x:", ql.Var("x")) // 输出 x: 3
}
```

这是一个最精简功能的 mini qlang。想要了解更多，可参考“[定制 qlang](README_QL.md#定制-qlang)”相关内容。实际项目中你也可以参考代码：

* [qlang/main.go](cmd/qlang/main.go)

### 非嵌入式场景下使用 qlang

尽管我们认为 qlang 的优势领域是在与 Go 配合的嵌入式场景，但您也可以把 qlang 语言用于非嵌入式脚本领域。

您可以直接通过 `qlang` 命令行程序来执行 *.ql 的代码。如下：

```
qlang xxx.ql <arg1> <arg2> ... <argN>
```

为了方便学习 qlang 工作机理，我们支持导出 qlang 编译的 “汇编指令”：

```
QLANG_DUMPCODE=1 qlang xxx.ql <arg1> <arg2> ... <argN>
```

在 Unix 系的操作系统下，您可以在 xxx.ql 文件开头加上：

```
#!/usr/bin/env qlang
```

并给 xxx.ql 文件加上可执行权限，即可直接运行 xxx.ql 脚本。

### 使用 qlang shell

命令行下输入 `qlang` 命令（不带参数），直接进入 qlang shell。

您同样也可以设置 QLANG_DUMPCODE 环境变量来学习 qlang 的工作机理：

```
QLANG_DUMPCODE=1 qlang
```

### 学习 qlang 语言特性

* [Q 语言手册](README_QL.md): 这里有语言特性的详细介绍。
* [Qlang Tutorials](tutorial): 这里是一些 qlang 的样例代码，供您学习 qlang 时参考。


## 下载

### 发行版本

* https://github.com/qiniu/qlang/releases

### 最新版本

```
go get -u -insecure qlang.io/qlang
```

或者在 src 目录执行如下命令：

```
git clone https://github.com/qiniu/qlang.git qlang.io
```

## 社区资源

### Embedded qlang (eql)

* [eql](cmd/eql): 全称 [embedded qlang](cmd/eql)，是类似 erubis/erb 的东西。结合 go generate 可很方便地让 Go 支持模板（不是 html template，是指泛型）。

### 为 Go package 导出 qlang module

* [qexport](cmd/qexport): 可为 Go 语言的标准库或者你自己写的 Go Package 自动导出相应的 qlang module。

### Qlang Modules

* 准官方库: https://github.com/qlangio/libs
* 第三方库: https://github.com/qlang-libs

### IDE 插件

* liteide - https://github.com/visualfc/liteide
* sublime - https://github.com/damonchen/sublime-qlang
* vim - https://github.com/simon-xia/vim-qlang
* atom - https://github.com/bingohuang/atom-language-q
