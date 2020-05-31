b := make([]float64, 5)
c := [2.4]
b[uint32(4)], c[0] = 123, 1.7
println(b, c)
