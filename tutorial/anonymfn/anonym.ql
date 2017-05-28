defer fn {
	println("x:", x)
}

defer fn {
	x; x = 2
}

x = 1