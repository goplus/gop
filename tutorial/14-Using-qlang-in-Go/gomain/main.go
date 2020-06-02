package main

import (
	"fmt"

	"github.com/qiniu/qlang/v6/tutorial/14-Using-qlang-in-Go/qlfoo"
)

func main() {
	rmap := qlfoo.ReverseMap(map[string]int{"Hi": 1, "Hello": 2})
	fmt.Println(rmap)
}
