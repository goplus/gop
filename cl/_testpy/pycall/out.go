package main

import (
	"github.com/goplus/llgo/c"
	"github.com/goplus/llgo/py"
	"github.com/goplus/llgo/py/math"
)

func main() {
	x := math.Sqrt(py.Float(2))
	c.Printf(c.Str("sqrt(2) = %f\n"), x.Float64())
}
