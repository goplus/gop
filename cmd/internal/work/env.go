package work

import (
	"os"
	"path/filepath"
)

func GOPROOT() string {
	if gopbin := os.Getenv("GOPROOT"); gopbin != "" {
		return gopbin
	}
	return filepath.Join(os.Getenv("HOME"), "gop")
}

func GOPBIN() string {
	if gopbin := os.Getenv("GOPBIN"); gopbin != "" {
		return gopbin
	}
	return filepath.Join(os.Getenv("HOME"), "gop", "bin")
}
