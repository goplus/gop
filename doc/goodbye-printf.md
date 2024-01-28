Goodbye printf
=====

For professional programmers, `printf` is a very familiar function, and it can be found in basically every language. However, `printf` is one of the most difficult functions for beginners to master.

Unfortunately, formatting a piece of information and displaying it to the user is a very common operation, so one has to remember their usage. While finding its documentation over the Internet can somewhat ease the burden of using it every time, it's far from a pleasant experience.

The most primitive way to format information is to use `string concat`:

```go
age := 10
println "age = " + age.string
```

And the most classic way of formatting information is to use `printf`:

```go
age := 10
printf "age = %d\n", age
```

Here `%d` means to format an integer value and `\n` means a newline.

To simplify format information in most cases, Go+ introduces `${expr}` expressions in string literals. For above example, you can replace `age.string` to `"${age}"`:

```go
age := 10
println "age = ${age}"
```

Here is a more complex example of `${expr}`:

```go
host := "foo.com"
page := 0
limit := 20
println "https://${host}/items?page=${page+1}&limit=${limit}" // https://foo.com/items?page=1&limit=20
println "$$" // $
```

This is a bit like how you feel at the `*nix` command line, right? To be more like it, we introduced a new builtin `echo` as an alias for `println`:

```go
age := 10
echo "age = ${age}"
```

