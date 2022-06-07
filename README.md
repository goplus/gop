<div align="center">
<p></p>
<p>
    <img width="80" src="https://goplus.org/favicon.svg">
</p>
<h1>The Go+ Programming Language</h1>

[goplus.org](https://goplus.org) | [Docs](doc/docs.md) | [Go+ vs. Go](doc/goplus-vs-go.md) | [Playground](https://play.goplus.org) | [Tutorials](https://tutorial.goplus.org/) | [Contributing & compiler design](doc/contributing.md)

</div>

<div align="center">
<!--
[![VSCode](https://img.shields.io/badge/vscode-Go+-teal.svg)](https://github.com/gopcode/vscode-goplus)
[![Discord](https://img.shields.io/discord/983646982100897802?label=Discord&logo=discord&logoColor=white)](https://discord.gg/mYjWCJDcAr)
-->

[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/gop)](https://goreportcard.com/report/github.com/goplus/gop)
[![Coverage Status](https://codecov.io/gh/goplus/gop/branch/main/graph/badge.svg)](https://codecov.io/gh/goplus/gop)
[![GitHub release](https://img.shields.io/github/v/tag/goplus/gop.svg?label=release)](https://github.com/goplus/gop/releases)
[![Discord](https://img.shields.io/badge/Discord-online-success.svg?logo=discord&logoColor=white)](https://discord.gg/mYjWCJDcAr)
[![Interpreter](https://img.shields.io/badge/interpreter-iGo+-seagreen.svg)](https://github.com/goplus/igop)

</div>

The Go+ programming language is designed for engineering, STEM education, and data science.

* **For engineering**: working in the simplest language that can be mastered by children.
* **For STEM education**: studying an engineering language that can be used for work in the future.
* **For data science**: communicating with engineers in the same language.

For more detials, see [Quick Start](doc/docs.md).

## Key Features of Go+

* A static typed language.
* The simplest engineering language that can be mastered by children (script-like style).
* Performance: as fast as Go (Go+'s main backend compiles to human-readable Go).
* Fully compatible with [Go](https://github.com/golang/go) and can mix Go/Go+ code in the same package (see [Go/Go+ hybrid programming](doc/docs.md#gogo-hybrid-programming)).
* REPL: see [iGo+](https://github.com/goplus/igop).
* [Simplest way to interaction with C](doc/docs.md#calling-c-from-go) (cgo is supported but not recommended).
* [Powerful built-in data processing capabilities](doc/docs.md#data-processing).

## Installing Go+

### from source code

For now, we suggest you install Go+ from source code.

Note: Requires go1.16 or later

```bash
git clone https://github.com/goplus/gop.git
cd gop

# On mac/linux run:
./all.bash
# On Windows run:
all.bat
```

### on macOS

```sh
brew install goplus
```
