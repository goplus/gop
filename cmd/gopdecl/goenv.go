package main

import (
	"os"
	"path/filepath"
	"runtime"
)

func initGoModCache() string {
	val := os.Getenv("GOMODCACHE")
	if val == "" {
		val = filepath.Join(runtime.GOROOT(), "pkg/mod")
		os.Setenv("GOMODCACHE", val)
	}
	return val
}
