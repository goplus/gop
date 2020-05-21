import "fmt"

foo := func(prompt string) (n int, err error) {
    n, err = fmt.Println(prompt + x)
    return
}

x := "Hello, world!"
fmt.Println(foo("x: "))

printf := func(format string, args ...interface{}) (n int, err error) {
    n, err = fmt.Printf(format, args...)
}

fmt.Println(printf("Hello, %v\n", "qlang"))
