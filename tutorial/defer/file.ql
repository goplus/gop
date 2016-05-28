
f, err = os.open("file.ql")
if err != nil {
	fprintln(os.stderr, err)
	return 1
}
defer println("exit!")
defer f.close()

b = make([]byte, 8)
n, err = f.read(b)
if err != nil {
	fprintln(os.stderr, "Read failed:", err)
	return 2
}

println(string(b[:n]))
