ch = make(chan bool, 2)
x = 1

go fn {
	defer ch <- true
	println("in goroutine1: x =", x)
	x; x++
}

go fn {
	defer ch <- true
	println("in goroutine2: x =", x)
	x; x++
}

<-ch
<-ch

// 结果并不一定是3，因为x++没有做到线程安全
//
println("in main routine: x =", x)
printf("chan cap=%d, len=%d\n", cap(ch), len(ch))

