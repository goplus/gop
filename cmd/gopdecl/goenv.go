package main

import (
	"os"
	"path/filepath"
	"runtime"
)

func initGoEnv() {
	val := os.Getenv("GOMODCACHE")
	if val == "" {
		os.Setenv("GOMODCACHE", filepath.Join(runtime.GOROOT(), "pkg/mod"))
	}
}
