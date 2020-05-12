import (
    "fmt"
    "strings"
)

func foo(x string) string {
    return strings.NewReplacer("?", "!").Replace(x)
}

fmt.Println(foo("Hello, world???"))
