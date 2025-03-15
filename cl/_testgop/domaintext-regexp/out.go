package main

import (
	"fmt"
	"github.com/goplus/gop/tpl/encoding/regexp"
)

func main() {
	re, err := regexp.New(`^[a-z]+\[[0-9]+\]$`)
	fmt.Println(re.MatchString("adam[23]"))
	_ = err
}
