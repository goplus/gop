sum := 0
for x <- [1, 3, 5, 7, 11, 13, 17], x > 3 {
    sum += x
}
println("sum(5,7,11,13,17):", sum)

fns := make([]func() int, 3)
for i, x <- [3, 15, 777] {
    v := x
    fns[i] = func() int {
        return v
    }
}
println("values:", fns[0](), fns[1](), fns[2]())
