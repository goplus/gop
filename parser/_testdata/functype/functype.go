type T struct {
	*T
	A int `json:"a"`
}

func bar(v chan bool) (int, <-chan error) {
	v <- true
	<-v
	return 0, (<-chan error)(nil)
}

func foo(f func([]byte, *string, ...T) chan<- int) (v int, err error) {
	return
}
