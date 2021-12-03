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


### Replace fmt.Print to builtin (NotImpl)

```go
import "fmt"

var n, err = fmt.Println("Hello world")
```

will be converted into:

```go
var n, err = println("Hello world")
```

Note:

* Convert `fmt.Print` => `print`
* Convert `fmt.Printf` => `printf`
* Convert `fmt.Println` => `println`
* ...


### Command style first (NotImpl)

```go
import "fmt"
import "os/exec"

cmd := exec.Command("go", "run", "hello.go")
cmd.Run()
fmt.Println("Hello world")
```

will be converted into:

```go
import "os/exec"

cmd := exec.command("go", "run", "hello.go")
cmd.run
println "Hello world"
```

Note:

* We prefer making a function call starting with lowercase, eg. `exec.Command` => `exec.command`, `cmd.Run` => `cmd.run`.
* Only function call statements can be converted into command style.

