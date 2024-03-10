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


## Command Style Code

Different from the function call style of most languages, Go+ recommends command style code:

```coffee
println "Hello world"
```

To emphasize our preference for command style, we introduce `echo` as an alias for `println`:

```coffee
echo "Hello world"
```

For more discussion on coding style, see https://tutorial.goplus.org/hello-world.


## Go+ Classfiles

```
One language can change the whole world.
Go+ is a "DSL" for all domains.
```

Rob Pike once said that if he could only introduce one feature to Go, he would choose `interface` instead of `goroutine`. `classfile` is as important to Go+ as `interface` is to Go.

In the design philosophy of Go+, we do not recommend `DSL` (Domain Specific Language). But `SDF` (Specific Domain Friendliness) is very important. The Go+ philosophy about `SDF` is:

```
Don't define a language for specific domain.
Abstract domain knowledge for it.
```

Go+ introduces `classfile` to abstract domain knowledge.

* [What's Classfile?](doc/classfile.md#whats-classfile)
* [Dive into Go+ Classfiles](doc/classfile.md)

Sound a bit abstract? Let's see some Go+ classfiles.

* Unit Test: [classfile: Unit Test](https://github.com/goplus/gop/blob/main/doc/classfile.md#classfile-unit-test)
* DevOps: [gsh: Go+ DevOps Tools](https://github.com/qiniu/x/tree/main/gsh)
* Web Programming: [yap: Yet Another HTTP Web Framework](https://github.com/goplus/yap)
* Web Programming: [yaptest: HTTP Test Framework](https://github.com/goplus/yap#yaptest-http-test-framework)
* Web Programming: [ydb: Database Framework](https://github.com/goplus/yap#ydb-database-framework)
* STEM Education: [spx: A Go+ 2D Game Engine](https://github.com/goplus/spx)


### yap: Yet Another HTTP Web Framework

This classfile has the file suffix `.yap`.

Create a file named [get.yap](https://github.com/goplus/yap/blob/main/demo/classfile2_hello/get.yap) with the following content:

```go
html `<html><body>Hello, YAP!</body></html>`
```

Execute the following commands:

```sh
gop mod init hello
gop get github.com/goplus/yap@latest
gop mod tidy
gop run .
```

A simplest web program is running now. At this time, if you visit http://localhost:8080, you will get:

```
Hello, YAP!
```

YAP uses filenames to define routes. `get.yap`'s route is `get "/"` (GET homepage), and `get_p_#id.yap`'s route is `get "/p/:id"` (In fact, the filename can also be `get_p_:id.yap`, but it is not recommended because `:` is not allowed to exist in filenames under Windows).

Let's create a file named [get_p_#id.yap](https://github.com/goplus/yap/blob/main/demo/classfile2_hello/get_p_%23id.yap) with the following content:

```coffee
json {
	"id": ${id},
}
```

Execute `gop run .` and visit http://localhost:8080/p/123, you will get:

```
{"id": "123"}
```

See [yap: Yet Another HTTP Web Framework](https://github.com/goplus/yap) for more details.


### spx: A Go+ 2D Game Engine

![Screen Shot1](https://github.com/goplus/spx/blob/main/tutorial/01-Weather/1.jpg) ![Screen Shot2](https://github.com/goplus/spx/blob/main/tutorial/01-Weather/2.jpg)

Through this example you can learn how to implement dialogues between multiple actors.

Here are some codes in [Kai.spx](https://github.com/goplus/spx/blob/main/tutorial/01-Weather/Kai.spx):

```coffee
onStart => {
	say "Where do you come from?", 2
	broadcast "1"
}

onMsg "2", => {
	say "What's the climate like in your country?", 3
	broadcast "3"
}
```

We call `onStart` and `onMsg` to listen events. `onStart` is called when the program is started. And `onMsg` is called when someone calls `broadcast` to broadcast a message.

When the program starts, Kai says `Where do you come from?`, and then broadcasts the message `1`. Who will recieve this message? Let's see codes in [Jaime.spx](https://github.com/goplus/spx/blob/main/tutorial/01-Weather/Jaime.spx):

```coffee
onMsg "1", => {
	say "I come from England.", 2
	broadcast "2"
}
```

Yes, Jaime recieves the message `1` and says `I come from England.`. Then he broadcasts the message `2`. Kai recieves it and says `What's the climate like in your country?`.

The following procedures are very similar. In this way you can implement dialogues between multiple actors.

See [spx: A Go+ 2D Game Engine](https://github.com/goplus/spx) for more details.


### gsh: Go+ DevOps Tools

Yes, now you can write `shell script` in Go+. It supports all shell commands.

Let's create a file named [example.gsh](https://github.com/qiniu/x/blob/main/gsh/demo/hello/example.gsh) and write the following code:

```coffee
mkdir "testgsh"
```

Don't need a `go.mod` file, just enter `gop run ./example.gsh` directly to run.

See [gsh: Go+ DevOps Tools](https://github.com/qiniu/x/tree/main/gsh) for more details.


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

### Web Programming

* [yap: Yet Another HTTP Web Framework](https://github.com/goplus/yap)
* [yaptest: HTTP Test Framework](https://github.com/goplus/yap/tree/main/ytest)
* [ydb: Database Framework](https://github.com/goplus/yap#ydb-database-framework)

### DevOps Tools

* [gsh: Go+ DevOps Tools](https://github.com/qiniu/x/tree/main/gsh)

### Data Processing

* [hdq: HTML DOM Query Language for Go+](https://github.com/goplus/hdq)


## IDE Plugins

* vscode: [Go/Go+ for Visual Studio Code](https://marketplace.visualstudio.com/items?itemName=goplus.gop)


## Contributing

The Go+ project welcomes all contributors. We appreciate your help!

For more details, see [Contributing & compiler design](doc/contributing.md).


## Give a Star! ⭐

If you like or are using Go+ to learn or start your projects, please give it a star. Thanks!
