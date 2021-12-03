Go+ Code Style
======

TODO

## Go+ Style Specification

### No package main

```go
package main

func f() {
}
```

will be converted into:

```
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

```
a := 0
```
