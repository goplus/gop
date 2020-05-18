x := 0
if t := false; t {
    x = 3
} else {
    x = 5
}
println("x:", x)

x = 0
switch s := "Hello"; s {
default:
    x = 7
case "world", "hi":
    x = 5
case "xsw":
    x = 3
}
println("x:", x)
