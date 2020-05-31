x := [1, 3.4] // []float64
println("x:", x)

y := [1] // []int
println("y:", y)

z := [1+2i, "xsw"] // []interface{}
println("z:", z)

println([1, 3.4, 3+4i]) // []complex128
println([5+6i])         // []complex128
println(["xsw", 3])     // []interface{}

println("empty slice:", []) // []interface{}
