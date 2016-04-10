import sys

primes = [2, 3]
n = 1
limit = 9

def isPrime(v):
	i = 0
	while i < n:
		if v % primes[i] == 0:
			return False
		i += 1
	return True

def listPrimes(max):
	global n, limit
	v = 5
	while True:
		while v < limit:
			if isPrime(v):
				primes.append(v)
				if v * v >= max:
					return
			v += 2
		v += 2
		n += 1
		limit = primes[n] * primes[n]

def maxPrimeOf(max):
	global n
	if max % 2 == 0:
		max -= 1
	listPrimes(max)
	n = len(primes)
	while True:
		if isPrime(max):
			return max
		max -= 2

if len(sys.argv) < 2:
	print 'Usage: maxprime <Value>'
	sys.exit(1)

max = int(sys.argv[1])
if max < 8:
	sys.exit(2)

max -= 1
v = maxPrimeOf(max)
print v

