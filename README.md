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

Our vision is to **enable everyone to become a builder of the digital world**.

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


## Key Features of Go+

* Approaching natural language expression and intuitive (see [How Go+ simplifies Go's expressions](#how-go-simplifies-gos-expressions)).
* Smallest but Turing-complete syntax set in best practices (see [The Go+ Mini Specification](doc/spec-mini.md)).
* Fully compatible with [Go](https://github.com/golang/go) and can mix Go/Go+ code in the same package (see [The Go+ Full Specification](doc/spec.md) and [Go/Go+ Hybrid Programming](doc/docs.md#gogo-hybrid-programming)).
* Integrating with the C ecosystem including Python and providing limitless possibilities based on [LLGo](https://github.com/goplus/llgo) (see [Importing C/C++ and Python libraries](#importing-cc-and-python-libraries)).
* Does not support DSL (Domain-Specific Languages), but supports SDF (Specific Domain Friendliness) (see [Go+ Classfiles](#go-classfiles)).


## How Go+ simplifies Go's expressions

Different from the function call style of most languages, Go+ recommends command style code:

```coffee
println "Hello world"
```

To emphasize our preference for command style, we introduce `echo` as an alias for `println`:

```coffee
echo "Hello world"
```

For more discussion on coding style, see https://tutorial.goplus.org/hello-world.

Code style is just the first step. We have made many efforts to make the code more intuitive and closer to natural language expression. These include:

| Go code | Go+ code | Note |
| ---- | ---- | ---- |
| package main<br><br>import "fmt"<br><br>func main() {<br>&nbsp;&nbsp;&nbsp;&nbsp;fmt.Println("Hi")<br>} | import "fmt"<br><br>fmt.Println("Hi")<br> | Program structure: Go+ allows omitting `package main` and `func main` |
| fmt.Println("Hi") | echo("Hi") | More builtin functions: It simplifies the expression of the most common tasks |
| fmt.Println("Hi") | echo "Hi" | Command-line style code: It reduces the number of parentheses in the code as much as possible, making it closer to natural language |
| a := []int{1, 2, 3} | a := [1, 2, 3] | List literals |
| a = append(a, 4)<br>a = append(a, 5, 6, 7) | a <- 4<br>a <- 5, 6, 7 | Append values to a list |
| a := map[string]int{<br>&nbsp;&nbsp;&nbsp;&nbsp;"Monday": 1,<br>&nbsp;&nbsp;&nbsp;&nbsp;"Tuesday": 2,<br>} | a := {<br>&nbsp;&nbsp;&nbsp;&nbsp;"Monday": 1,<br>&nbsp;&nbsp;&nbsp;&nbsp;"Tuesday": 2,<br>} | Mapping literals |
| OnStart(func() {<br>&nbsp;&nbsp;&nbsp;&nbsp;...<br>}) | onStart => {<br>&nbsp;&nbsp;&nbsp;&nbsp;...<br>} | Lambda expressions |
| type Rect struct {<br>&nbsp;&nbsp;&nbsp;&nbsp;Width&nbsp; float64<br>&nbsp;&nbsp;&nbsp;&nbsp;Height float64<br>}<br><br>func (this *Rect) Area() float64 { <br>&nbsp;&nbsp;&nbsp;&nbsp;return this.Width * this.Height<br>} | var (<br>&nbsp;&nbsp;&nbsp;&nbsp;Width&nbsp; float64<br>&nbsp;&nbsp;&nbsp;&nbsp;Height float64<br>)<br><br>func Area() float64 { <br>&nbsp;&nbsp;&nbsp;&nbsp;return Width * Height<br>} | [Go+ Classfiles](doc/classfile.md): We can express OOP with global variables and functions. |

For more details, see [The Go+ Mini Specification](doc/spec-mini.md).


## Importing C/C++ and Python libraries

Go+ can choose different Go compilers as its underlying support. Currently known supported Go compilers include:

* [go](https://go.dev/) (The official Go compiler supported by Google)
* [llgo](https://github.com/goplus/llgo) (The Go compiler supported by the Go+ team)
* [tinygo](https://tinygo.org/) (A Go compiler for small places)

Currently, Go+ defaults to using [go](https://go.dev/) as its underlying support, but in the future, it will be [llgo](https://github.com/goplus/llgo).

LLGo is a Go compiler based on [LLVM](https://llvm.org/) in order to better integrate Go with the C ecosystem including Python. It aims to expand the boundaries of Go/Go+, providing limitless possibilities such as:

* Game development
* AI and data science
* WebAssembly
* Embedded development
* ...

If you wish to use [llgo](https://github.com/goplus/llgo), specify the `-llgo` flag when initializing a Go+ module:

```sh
gop mod init -llgo YourModulePath
```

This will generate a `go.mod` file with the following contents (It may vary slightly depending on the versions of local Go+ and LLGo):

```go
module YourModulePath

go 1.21 // llgo 1.0

require github.com/goplus/llgo v0.9.9
```

Based on LLGo, Go+ can import libraries written in C/C++ and Python.

Here is an example (see [chello](demo/_llgo/chello/hello.gop)) of printing `Hello world` using C's `printf`:

```go
import "c"

c.printf c"Hello world\n"
```

Here, `c"Hello world\n"` is a syntax supported by Go+, representing a null-terminated C-style string.

To run this example, you can:

```sh
cd YourModulePath  # set work directory to your module
gop mod tidy       # for generating go.sum file
gop run .
```

And here is an example (see [pyhello](demo/_llgo/pyhello/hello.gop)) of printing `Hello world` using Python's `print`:

```go
import "py/std"

std.print py"Hello world"
```

Here, `py"Hello world"` is a syntax supported by Go+, representing a Python string.

Here are more examples of Go+ calling C/C++ and Python libraries:

* [pytensor](demo/_llgo/pytensor/tensor.gop): a simple demo using [py/torch](https://pkg.go.dev/github.com/goplus/llgo/py/torch)
* [tetris](demo/_llgo/tetris/tetris.gop): a tetris game based on [c/raylib](https://pkg.go.dev/github.com/goplus/llgo/c/raylib)
* [sqlitedemo](demo/_llgo/sqlitedemo/sqlitedemo.gop): a demo using [c/sqlite](https://pkg.go.dev/github.com/goplus/llgo/c/sqlite)

To find out more about LLGo/Go+'s support for C/C++ and Python in detail, please refer to homepage of [llgo](https://github.com/goplus/llgo).


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
* Web Programming: [yaptest: HTTP Test Framework](https://github.com/goplus/yap/tree/main/ytest)
* Web Programming: [ydb: Database Framework](https://github.com/goplus/yap/tree/main/ydb)
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


## How to install

Note: Requires go1.19 or later

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
