
primes = [2, 3]
n = 1
limit = 9

isPrime = fn(v) {
	for i = 0; i < n; i++ {
		if v % primes[i] == 0 {
			return false
		}
	}
	return true
}

listPrimes = fn(max) {

	v = 5
	for {
		for v < limit {
			if isPrime(v) {
				primes = append(primes, v)
				if v * v >= max {
					return
				}
			}
			v += 2
		}
		v += 2
		n; n++
		limit = primes[n] * primes[n]
	}
}

maxPrimeOf = fn(max) {

	if max % 2 == 0 {
		max--
	}

	listPrimes(max)
	n; n = len(primes)

	for {
		if isPrime(max) {
			return max
		}
		max -= 2
	}
}

// Usage: maxprime <Value>
//
if len(os.args) < 2 {
	fprintln(os.stderr, "Usage: maxprime <Value>")
	return
}

max, err = strconv.parseInt(os.args[1], 10, 64)
if err != nil {
	fprintln(os.stderr, err)
	return 1
}
if max < 8 { // <8 的情况下，可直接建表，答案略
	return
}

max--
v = maxPrimeOf(max)
println(v)

