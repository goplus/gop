import (
    "fmt"
    "strings"
)

func foo(x string) string {
    return strings.NewReplacer("?", "!").Replace(x)
}

func bar(n int, err error) {
	fmt.Println(n, err)
}

bar(fmt.Println(foo("Hello, world???")))
