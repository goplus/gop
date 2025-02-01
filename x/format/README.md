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
n, err := echo("Hello world")
```

Note:

* Convert `fmt.Errorf => errorf`
* Convert `fmt.Fprint => fprint`
* Convert `fmt.Fprintf => fprintf`
* Convert `fmt.Fprintln => fprintln`
* Convert `fmt.Print` => `print`
* Convert `fmt.Printf` => `printf`
* Convert `fmt.Println` => `echo`
* Convert `fmt.Sprint => sprint`
* Convert `fmt.Sprintf => sprintf`
* Convert `fmt.Sprintln => sprintln`

### Command style first

```go
import "fmt"

fmt.Println()
fmt.Println(fmt.Println("Hello world"))
```

will be converted into:

```go
echo
echo echo("Hello world")
```

Note:

* Only the outermost function call statement is converted into command style. So `fmt.Println(fmt.Println("Hello world"))` is converted into `echo echo("Hello world")`, not `echo echo "Hello world"`.


### pkg.Fncall starting with lowercase

```go
import "math"

echo math.Sin(math.Pi/3)
```

will be converted into:

```go
echo math.sin(math.Pi/3)
```

### Funclit of params convert to lambda in fncall (skip named results)

```go
echo(demo(func(n int) int {
	return n+100
}))

echo(demo(func(n int) (v int) {
	return n+100
}))

onStart(func() {
	echo("start")
})
```

will be converted into:
```
echo demo(n => n + 100)

echo demo(func(n int) (v int) {
	return n + 100
})

onStart => {
	echo "start"
}
```
