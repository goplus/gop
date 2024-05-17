package hello

import "github.com/goplus/llgo/c"

func Main() {
	c.Printf(c.Str("Hello world\n"))
}
