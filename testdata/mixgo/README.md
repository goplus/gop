This is an example to show how to mix Go/Go+ programming in the same package.

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

The output of this example is as follows:

```
Mix Go and Go+
Hello, world
```
