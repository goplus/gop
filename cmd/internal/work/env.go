package work

import (
	"os"
	"path/filepath"
)

func GOPBIN() string {
	if gopbin := os.Getenv("GOPBIN"); gopbin != "" {
		return gopbin
	}
	return filepath.Join(os.Getenv("HOME"), "gop", "bin")
}
