y := [x * x for x <- [1, 3, 5, 7, 11]]
println(y)

z := [i + v for i, v <- [1, 3, 5, 7, 11]]
println(z)

println([k + "," + s for k, s <- {"Hello": "xsw", "Hi": "qlang"}])
