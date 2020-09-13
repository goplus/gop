func foo() []int {
	return make([]int, 10)
}

func foo1() map[int]int {
	return make(map[int]int, 10)
}

func foo2() chan int {
	return make(chan int, 10)
}
a := foo()
if a != nil {
    println("foo")
}

a1 := foo1()
if a1 != nil {
    println("foo1")
}
a2 := foo2()
if a2 != nil {
    println("foo2")
}