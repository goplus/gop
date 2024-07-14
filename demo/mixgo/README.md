This is an example to show how to mix Go/Go+ code in the same package.

In this example, we have a Go source file named `a.go`:

```go
package main

import "fmt"

func p(a interface{}) {
	sayMix()
	fmt.Println("Hello,", a)
}
```

And we have a Go+ source file named `b.gop`:

```go
func sayMix() {
	println "Mix Go and Go+"
}

p "world"
```

You can see that Go calls a Go+ function named `sayMix`, and Go+ calls a Go function named `p`. As you are used to in Go programming, this kind of circular reference is allowed.

Run `gop run .` to see the output of this example:

```
Mix Go and Go+
Hello, world
```

### Give a Star! ⭐

If you like or are using Go+ to learn or start your projects, please give it a star. Thanks!
