a := make([]int, 2)
a = append(a, 1, 2, 3)
println(a)

b := make([]int, 0, 4)
c := [1, 2, 3]
b = append(b, c...)
println(b)
