package main

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/lib/py"
	"github.com/goplus/lib/py/math"
)

func main() {
	x := math.Sqrt(py.Float(2))
	c.Printf(c.Str("sqrt(2) = %f\n"), x.Float64())
}
