Go+ Code Style
======

TODO


## Specification

### No package main

```go
package main

func f() {
}
```

will be converted into:

```go
func f() {
}
```


### No func main

```go
package main

func main() {
    a := 0
}
```

will be converted into:

```go
a := 0
```


### Replace fmt.Print to builtin

```go
import "fmt"

n, err := fmt.Println("Hello world")
```

will be converted into:

```go
n, err := println("Hello world")
```

Note:

* Convert `fmt.Print` => `print`
* Convert `fmt.Printf` => `printf`
* Convert `fmt.Println` => `println`
* ...


### Command style first

```go
import "fmt"

fmt.Println()
fmt.Println(fmt.Println("Hello world"))
```

will be converted into:

```go
println
println println("Hello world")
```

Note:

* Only the outermost function call statement is converted into command style. So `fmt.Println(fmt.Println("Hello world"))` is converted into `println println("Hello world")`, not `println println "Hello world"`.


### pkg.Fncall starting with lowercase

```go
import "math"

println math.Sin(math.Pi/3)
```

will be converted into:

```go
println math.sin(math.Pi/3)
```
