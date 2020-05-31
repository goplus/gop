x := []float64{1, 3.4, 5}
y := map[string]float64{"Hello": 1, "xsw": 3.4}

a := [...]float64{1, 3.4, 5}
b := [...]float64{1, 3: 3.4, 5}
c := []float64{2: 1.2, 3, 6: 4.5}

x[1], y["xsw"] = 1.7, 2.8
println("x:", x, "y:", y)
println(`x[1]:`, x[1], `y["xsw"]:`, y["xsw"], `a[1]:`, a[1])

i := uint16(4)
b[uint32(4)], c[i] = 123, 1.7
println("a:", a, "b:", b, "c:", c)

arr := [...]float64{1, 2}
title := "Hello,world!" + "2020-05-27"
println(title[:len(title)-len("2006-01-02")], len(arr), arr[1:])
