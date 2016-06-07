uint32set.v1
======

Uint32Set is a set of uint32, implemented via `map[uint32]struct{}` for minimal memory consumption. Here is an example:

```go
import (
    "uint32set.v1"
)

set := uint32set.New(1, 2)
set.Add(3)

if set.Has(4) {
    println("set has item `4`")
}
```
