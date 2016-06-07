stringset.v1
======

StringSet is a set of string, implemented via `map[string]struct{}` for minimal memory consumption. Here is an example:

```go
import (
    "stringset.v1"
)

set := stringset.New("a", "b")
set.Add("c")

if set.Has("d") {
    println("set has item `"d"`")
}
```
