
f, err = os.Open("file.ql")
if err != nil {
	fprintln(os.Stderr, err)
	return 1
}
defer println("exit!")
defer f.Close()

b = make([]byte, 8)
n, err = f.Read(b)
if err != nil {
	fprintln(os.Stderr, "Read failed:", err)
	return 2
}

println(string(b[:n]))
