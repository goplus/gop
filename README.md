The Q Language (Q语言)
========

[![Build Status](https://travis-ci.org/qiniu/qlang.png?branch=develop)](https://travis-ci.org/qiniu/qlang)

![logo](http://qiniutek.com/images/logo-2.png)

# 下载

### 源代码

```
go get -u -insecure qlang.io/qlang
```

或者在 src 目录执行如下命令：

```
mkdir qiniupkg.com
git clone https://github.com/qiniu/qlang.git qlang.io
git clone https://github.com/qiniu/text.git qiniupkg.com/text
```

# 语言特色

* 最大卖点：与 Go 语言有最好的互操作性。所有 Go 语言的社区资源可以直接为我所用。
* 有赖于 Go 语言的互操作性，这门语言不需要自己实现标准库。尽管年轻，但是这门语言已经具备商用的可行性。
* 微内核：语言的核心只有 1200 行代码。所有功能以可插拔的 module 方式提供。

预期的商业场景：

* 由于与 Go 语言的无缝配合，qlang 在嵌入式脚本领域有 lua、python、javascript 所不能比拟的优越性。比如：网络游戏中取代 lua 的位置。
* 作为编译原理的教学语言。由于 qlang 的 Compiler 代码极短，便于阅读和理解，非常方便教学实践之用。


# 快速入门

一个基础版本的 qlang 应该是这样的：

```go
import (
	"fmt"

	"qlang.io/qlang.v2/qlang"
	_ "qlang.io/qlang/builtin" // 导入 builtin 包
)

const scriptCode = `
	x = 1 + 2
`

func main() {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		// 错误处理
		return
	}

	err = lang.SafeExec([]byte(scriptCode), "")
	if err != nil {
		// 错误处理
		return
	}

	v, _ := lang.Var("x")
	fmt.Println("x:", v) // 输出 x: 3
}
```

这是一个最精简功能的 mini qlang。想要了解更多，可参考后文“定制 qlang”一节。实际项目中你也可以参考代码：

* [qlang/main.go](https://github.com/qiniu/qlang/blob/develop/app/qlang/main.go)

你也可以把 qlang 用于非嵌入式脚本领域，直接用 `qlang` 程序来执行 *.ql 的代码。如下：

```
qlang xxx.ql <arg1> <arg2> ... <argN>
```

为了方便学习和调试问题，我们还支持导出 qlang 编译的 “汇编指令”：

```
QLANG_DUMPCODE=1 qlang xxx.ql <arg1> <arg2> ... <argN>
```

如果 qlang 命令不带参数，则进入 qlang shell。更详细的内容，请参考《[Q语言手册](README_QL.md)》。

# 相关资源

## Go 语言包自动生成qlang library

* [qexport](app/qexport/README.md): 可为一个Go语言的标准库或者你自己写的Go Package自动导出相应的qlang library。

## Qlang Libs

* 准官方库: https://github.com/qlangio/libs
* 第三方库: https://github.com/qlang-libs

## IDE 插件

* liteide - https://github.com/visualfc/liteide
* sublime - https://github.com/damonchen/sublime-qlang
