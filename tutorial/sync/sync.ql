mutex = sync.NewMutex()

n = 10
x = 1

wg = sync.NewWaitGroup()
wg.Add(n)

for i = 0; i < n; i++ {
	go fn {
		defer wg.Done()

		mutex.Lock()
		defer mutex.Unlock()
		x; x++
	}
}

wg.Wait()

println("x:", x)
