import (
    "fmt"
    "strings"
)

func foo(x string) string {
    return strings.NewReplacer("?", "!").Replace(x)
}

func printf(format string, args ...interface{}) (n int, err error) {
    n, err = fmt.Printf(format, args...)
    return
}

fmt.Println(foo("Hello, world???"))
fmt.Println(printf("Hello, %v\n", "qlang"))
