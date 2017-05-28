
runtime.GOMAXPROCS(8)

wg = sync.NewWaitGroup()
wg.Add(2)

x = 1

go fn {
	defer wg.Done()
	println("in goroutine1: x =", x)
	x++
}

go fn {
	defer wg.Done()
	println("in goroutine2: x =", x)
	x++
}

wg.Wait()

// 结果并不一定是3，因为x++没有做到线程安全
//
println("in main routine: x =", x)
