<div align="center">
<p></p>
<p>
    <img width="80" src="https://goplus.org/favicon.svg">
</p>
<h1>The Go+ Programming Language</h1>

[goplus.org](https://goplus.org) | [Docs](doc/docs.md) | [Go+ vs. Go](doc/goplus-vs-go.md) | [Tutorials](https://tutorial.goplus.org/) | [Playground](https://play.goplus.org) | [iGo+ Playground](https://repl.goplus.org/) | [Contributing & compiler design](doc/contributing.md)

</div>

<div align="center">
<!--
[![VSCode](https://img.shields.io/badge/vscode-Go+-teal.svg)](https://github.com/gopcode/vscode-goplus)
[![Discord](https://img.shields.io/discord/983646982100897802?label=Discord&logo=discord&logoColor=white)](https://discord.gg/mYjWCJDcAr)
[![Interpreter](https://img.shields.io/badge/interpreter-iGo+-seagreen.svg)](https://github.com/goplus/igop)
-->

[![Build Status](https://github.com/goplus/gop/actions/workflows/go.yml/badge.svg)](https://github.com/goplus/gop/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/gop)](https://goreportcard.com/report/github.com/goplus/gop)
[![Coverage Status](https://codecov.io/gh/goplus/gop/branch/main/graph/badge.svg)](https://codecov.io/gh/goplus/gop)
[![GitHub release](https://img.shields.io/github/v/tag/goplus/gop.svg?label=release)](https://github.com/goplus/gop/releases)
[![Discord](https://img.shields.io/badge/Discord-online-success.svg?logo=discord&logoColor=white)](https://discord.com/invite/mYjWCJDcAr)

</div>

Our vision is to **enable everyone to create production-level applications**.

#### Easy to learn

* Simple and easy to understand
* Smaller syntax set than Python in best practices

#### Ready for large projects

* Derived from Go and easy to build large projects from its good engineering foundation

The Go+ programming language is designed for engineering, STEM education, and data science.

* **For engineering**: working in the simplest language that can be mastered by children.
* **For STEM education**: studying an engineering language that can be used for work in the future.
* **For data science**: communicating with engineers in the same language.

For more details, see [Quick Start](doc/docs.md).


## Go+ Classfiles

Rob Pike once said that if he could only introduce one feature to Go, he would choose `interface` instead of `goroutine`. `classfile` is as important to Go+ as `interface` is to Go.

In the design philosophy of Go+, we do not recommend `DSL` (Domain Specific Language). But `SDF` (Specific Domain Friendliness) is very important. The Go+ philosophy about `SDF` is:

```
Don't define a language for specific domain.
Abstract domain knowledge for it.
```

Go+ introduces `classfile` to abstract domain knowledge.

Sound a bit abstract? Let's take web programming as an example. First let us initialize a hello project:

```sh
gop mod init hello
```

Then we have it reference a classfile called `yap` as the HTTP Web Framework:

```sh
gop get github.com/goplus/yap@latest
```

We can use it to implement a static file server:

```coffee
static "/foo", FS("public")
static "/"    # Equivalent to static "/", FS("static")

run ":8080"
```

We can also add the ability to handle dynamic GET/POST requests:

```coffee
static "/foo", FS("public")
static "/"    # Equivalent to static "/", FS("static")

get "/p/:id", ctx => {
	ctx.json {
		"id": ctx.param("id"),
	}
}

run ":8080"
```

Save this code to `hello_yap.gox` file and execute:

```sh
mkdir -p yap/static yap/public    # Static resources can be placed in these directories
gop mod tidy
gop run .
```

A simplest web program is running now. At this time, if you visit http://localhost:8080/p/123, you will get:

```
{"id":"123"}
```

Why is `yap` so easy to use? How does it do it? Click here to learn more about the [Go+ Classfiles](doc/classfile.md) mechanism and [YAP HTTP Web Framework](https://github.com/goplus/yap).


## Key Features of Go+

* A static typed language.
* The simplest engineering language that can be mastered by children (script-like style).
* Performance: as fast as Go (Go+'s main backend compiles to human-readable Go).
* Fully compatible with [Go](https://github.com/golang/go) and can mix Go/Go+ code in the same package (see [Go/Go+ hybrid programming](doc/docs.md#gogo-hybrid-programming)).
* No DSL (Domain Specific Language) support, but SDF ([Specific Domain Friendliness](doc/classfile.md)).
* Support Go code generation (main backend) and [bytecode backend](https://github.com/goplus/igop) (REPL: see [iGo+](https://repl.goplus.org/)).
* [Simplest way to interaction with C](doc/docs.md#calling-c-from-go) (cgo is supported but not recommended).
* [Powerful built-in data processing capabilities](doc/docs.md#data-processing).


## How to install

### on Windows

```sh
winget install goplus.gop
```

### on Debian/Ubuntu

```sh
sudo bash -c ' echo "deb [trusted=yes] https://pkgs.goplus.org/apt/ /" > /etc/apt/sources.list.d/goplus.list'
sudo apt update
sudo apt install gop
```

### on RedHat/CentOS/Fedora

```sh
sudo bash -c 'echo -e "[goplus]\nname=Go+ Repo\nbaseurl=https://pkgs.goplus.org/yum/\nenabled=1\ngpgcheck=0" > /etc/yum.repos.d/goplus.repo'
sudo yum install gop
```

### on macOS/Linux (Homebrew)

Install via [brew](https://brew.sh/)

```sh
$ brew install goplus
```

### from source code

For now, we suggest you install Go+ from source code.

Note: Requires go1.18 or later

```bash
git clone https://github.com/goplus/gop.git
cd gop

# On mac/linux run:
./all.bash
# On Windows run:
all.bat
```

## Go+ Applications

### 2D Games powered by Go+

* [A Go+ 2D Game Engine for STEM education](https://github.com/goplus/spx)
* [Aircraft War](https://github.com/goplus/AircraftWar)
* [Flappy Bird](https://github.com/goplus/FlappyCalf)
* [Maze Play](https://github.com/goplus/MazePlay)
* [BetaGo](https://github.com/xushiwei/BetaGo)
* [Gobang](https://github.com/xushiwei/Gobang)
* [Dinosaur](https://github.com/xushiwei/Dinosaur)

### HTTP Web Framework

* [yap: Yet Another Go/Go+ HTTP Web Framework](https://github.com/goplus/yap)

### HTTP Test

* [yaptest: HTTP Test Framework](https://github.com/goplus/yap/tree/main/ytest)

### DevOps tools

* [Go+ DevOps Tools](https://github.com/goplus/gop/blob/main/doc/dsl-vs-sdf.md#demo-go-devops-tools)

### Data processing

* [HTML DOM Query Language for Go+](https://github.com/goplus/hdq)


## IDE Plugins

* vscode: https://github.com/goplus/vscode-gop


## Contributing

The Go+ project welcomes all contributors. We appreciate your help!

For more details, see [Contributing & compiler design](doc/contributing.md).


## Give a Star! ⭐

If you like or are using Go+ to learn or start your projects, please give it a star. Thanks!
