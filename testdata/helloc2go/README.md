This is an example to show how Go+ interacts with C.

```go
import "C"

C.printf C"Hello, c2go!\n"
C.fprintf C.stderr, C"Hi, %7.1f\n", 3.14
```

Here we use `import "C"` to import libc. It's an abbreviation for `import "C/github.com/goplus/libc"`. It is equivalent to the following code:

```go
import "C/github.com/goplus/libc"

C.printf C"Hello, c2go!\n"
C.fprintf C.stderr, C"Hi, %7.1f\n", 3.14
```

In this example we call two C standard functions `printf` and `fprintf`, pass a C variable `stderr` and two C string in the form of `C"xxx"`.

The output of this example is as follows:

```
Hello, c2go!
Hi,     3.1
```

