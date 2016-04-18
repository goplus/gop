import "internal/a"
import "internal/a" as g

b = 3

g.b = 4
println("in script C:", a.b, g.b)
a.foo(b)

println("\nnow, we will panic")
println(a.a) // 因为 module a 并没有导出变量 a
