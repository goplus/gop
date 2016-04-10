
sumi = 0
sumv = 0

for i, v = range [10, 20, 30, 40, 50, 60, 70, 80, 90] {
	sumi += i
	sumv += v
	if i == 4 {
		break
	}
}

println("sumi:", sumi, "sumv:", sumv)

