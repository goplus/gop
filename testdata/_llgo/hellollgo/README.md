This is an example to show how Go+ interacts with C.

```go
import "c"

c.Printf C"Hello, llgo!\n"
c.Fprintf c.Stderr, C"Hi, %7.1f\n", 3.14
```

Here we use `import "c"` to import libc. It's an abbreviation for `import "github.com/goplus/llgo/c"`. It is equivalent to the following code:

```go
import "github.com/goplus/llgo/c"

c.Printf C"Hello, llgo!\n"
c.Fprintf c.Stderr, C"Hi, %7.1f\n", 3.14
```

In this example we call two C standard functions `printf` and `fprintf`, pass a C variable `stderr` and two C strings in the form of `C"xxx"`.

Run `gop run .` to see the output of this example:

```
Hello, llgo!
Hi,     3.1
```

### Give a Star! ⭐

If you like or are using Go+ to learn or start your projects, please give it a star. Thanks!
