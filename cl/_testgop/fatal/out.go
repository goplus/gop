package main

import (
	"github.com/goplus/gop/builtin/osx"
	"github.com/qiniu/x/stringutil"
	"os"
)

func main() {
	f, err := os.Open("hello.txt")
	if err != nil {
		osx.Fatal(stringutil.Concat("open file failed: ", err.Error()))
	}
	f.Close()
}
