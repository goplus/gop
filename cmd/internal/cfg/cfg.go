package cfg

import (
	"os"
	"path/filepath"
)

var (
	GOMODCACHE = envOr("GOMODCACHE", gopathDir("pkg/mod"))
)

// envOr returns Getenv(key) if set, or else def.
func envOr(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}
	return val
}

// GetGoPath return the gopath
func GetGoPath() string {
	return os.Getenv("GOPATH")
}

func gopathDir(rel string) string {
	list := filepath.SplitList(GetGoPath())
	if len(list) == 0 || list[0] == "" {
		return ""
	}
	return filepath.Join(list[0], rel)
}
