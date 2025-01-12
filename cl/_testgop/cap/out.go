package main

import (
	"fmt"
	"github.com/goplus/gop/builtin/stringslice"
)

func main() {
	a := []string{"hello", "world", "123"}
	fmt.Println(stringslice.Capitalize(a))
}
