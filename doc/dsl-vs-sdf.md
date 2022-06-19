Go+ ClassFile: DSL vs. SDF
=====

Go+ doesn't support DSL (Domain Specific Language), but it's SDF (Specific Domain Friendly).

The Go+ philosophy of Domain is:

```
Don't define a language for specific domain.
Abstract domain knowledge for it.
```

Go+ introduces `ClassFile` to abstract domain knowledge.

NOTE: `ClassFile` is still in beta and may be subject to change.

What is `ClassFile`?

TODO


## DEMO: Go+ DevOps Tools

This is a demo project to implement a new `ClassFile` work: [github.com/xushiwei/gsh](https://github.com/xushiwei/gsh).

Before we dive into implementation of `github.com/xushiwei/gsh`, let's take a look how to use it.

First, let's create a new module named [gshexample](https://github.com/xushiwei/gshexample):

```sh
mkdir gshexample
cd gshexample
gop mod init gshexample
```

Second, get `github.com/xushiwei/gsh` into `gshexample`:

```sh
gop get github.com/xushiwei/gsh
```

Third, create a Go+ source file named `./example.gsh` and write the following code:

```coffee
mkdir "testgsh"
```

Fourth, run this project:

```sh
gop mod tidy
gop run .
```

It's strange to you that the file extension of Go+ source is not `.gop` but `.gsh`.

But the Go+ compiler can recognize `.gsh` files as Go+ source because `github.com/xushiwei/gsh` register `.gsh` as its `ClassFile` file extension.

We can change `./example.gsh` more complicated:

```coffee
type file struct {
	name  string
	fsize int
}

mkdir! "testgsh"

mkdir "testgsh2"
lastErr!

mkdir "testgsh3"
if lastErr != nil {
	panic lastErr
}

capout => { ls }
println output.fields

capout => { ls "-l" }
files := [file{flds[8], flds[4].int!} for e <- output.split("\n") if flds := e.fields; flds.len > 2]
println files

rmdir "testgsh", "testgsh2", "testgsh3"
```

### Check last error

If we want to ensure `mkdir` successfully, there are three ways:

The simplest way is:

```coffee
mkdir! "testsh"  # will panic if mkdir failed
```

The second way is:

```coffee
mkdir "testsh"
lastErr!
```

Yes, `github.com/xushiwei/gsh` provides `lastErr` to check last error.

The third way is:

```coffee
mkdir "testsh"
if lastErr != nil {
    panic lastErr
}
```

This is the most familiar way to Go developers.

### Capture output of commands

And, `github.com/xushiwei/gsh` provides a way to capture output of commands:

```coffee
capout => {
    ...
}
```

Similar to `lastErr`, the captured output result is saved to `output`.

For an example:

```coffee
capout => { ls "-l" }
println output
```

Here is a possible output:

```s
total 72
-rw-r--r--  1 xushiwei  staff  11357 Jun 19 00:20 LICENSE
-rw-r--r--  1 xushiwei  staff    127 Jun 19 10:00 README.md
-rw-r--r--  1 xushiwei  staff    365 Jun 19 00:25 example.gsh
-rw-r--r--  1 xushiwei  staff    126 Jun 19 09:33 go.mod
-rw-r--r--  1 xushiwei  staff    165 Jun 19 09:33 go.sum
-rw-r--r--  1 xushiwei  staff    110 Jun 19 09:33 gop.mod
-rw-r--r--  1 xushiwei  staff   1938 Jun 19 10:00 gop_autogen.go
```

We can use [Go+ powerful built-in data processing capabilities](doc/docs.md#data-processing) to process captured `output`:

```coffee
type file struct {
	name  string
	fsize int
}

files := [file{flds[8], flds[4].int!} for e <- output.split("\n") if flds := e.fields; flds.len > 2]
```

In this example, we split `output` by `"\n"`, and for each entry `e`, split it by spaces (`e.fields`) and save into `flds`. Condition `flds.len > 2` is to remove special line of output:

```s
total 72
```

At last, pick file name and size of all selected entries and save into `files`.


### Implementation of gsh

To most developers, they just use `ClassFile`. They don't need to implement a `ClassFile` project.

We provide a `ClassFile` template project: [github.com/goplus/classfile-project-template](https://github.com/goplus/classfile-project-template).

TODO
