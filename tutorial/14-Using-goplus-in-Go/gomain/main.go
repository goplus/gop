package main

import (
	"fmt"

	"github.com/goplus/gop/tutorial/14-Using-goplus-in-Go/foo"
)

func main() {
	rmap := foo.ReverseMap(map[string]int{"Hi": 1, "Hello": 2})
	fmt.Println(rmap)
}
