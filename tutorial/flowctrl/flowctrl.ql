weekday = "Friday"

v = switch weekday {
case "Monday":
	1
case "Tuesday":
	2
case "Wednesday":
	3
case "Thursday":
	4
case "Friday":
	5
case "Saterday":
	6
case "Sunday":
	7
default:
	0
}
println(weekday, "=>", v)

a = 3
b = 7
min = if a < b {
	a
} else {
	b
}
println("min", a, b, ":", min)

