mutex = sync.mutex()

n = 10
x = 1

wg = sync.waitGroup()
wg.add(n)

for i = 0; i < n; i++ {
	go fn {
		defer wg.done()

		mutex.lock()
		defer mutex.unlock()
		x; x++
	}
}

wg.wait()

println("x:", x)

