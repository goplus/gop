package main

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

func main() {
	fmt.Println(html.Parse(strings.NewReader(`<html><body><h1>hello</h1></body></html>`)))
}
