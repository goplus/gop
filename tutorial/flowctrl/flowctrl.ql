weekday = "Friday"

switch weekday {
case "Monday":
	v = 1
case "Tuesday":
	v = 2
case "Wednesday":
	v = 3
case "Thursday":
	v = 4
case "Friday":
	v = 5
case "Saterday":
	v = 6
case "Sunday":
	v = 7
default:
	v = 0
}
println(weekday, "=>", v)

a = 3
b = 7
if a < b {
	min = a
} else {
	min = b
}
println("min", a, b, ":", min)

