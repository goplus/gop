package main

import "github.com/goplus/gop/cl/internal/huh"

func main() {
	form := huh.New("1", 2, "<form id=\"test\">\n</form>\n")
	form.Run()
}
