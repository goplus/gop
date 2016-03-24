
f, err = os.open(__file__)
if err != nil {
	fprintln(os.stderr, err)
	return 1
}
defer println("exit!")
defer f.close()

b = slice("byte", 8)
n, err = f.read(b)
if err != nil {
	fprintln(os.stderr, "Read failed:", err)
	return 2
}

ret = sub(b, 0, n)
println(string(ret))

