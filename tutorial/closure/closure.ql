
test = fn() {
	g = 1
	return class {
		fn f() {
			return g
		}
	}
}

main {
	Foo = test()
	foo = new Foo
	g = 2
	println("foo.f:", foo.f())
}
