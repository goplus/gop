package main

import (
	"io"
	"log"
	"testing"
)

type ift interface {
	io.Closer
	f(int) string
	g()
}

type impl struct {
	a T
}

func Bar(t *testing.T) int {
	log.Println("Hello")
	t.Log("Hello",
		"world")
	return 0
}
