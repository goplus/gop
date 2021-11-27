package modfetch

import (
	"fmt"
	"path/filepath"

	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"

	"github.com/goplus/gop/cmd/internal/cfg"
)

func cacheDir(path string) (string, error) {
	enc, err := module.EscapePath(path)
	if err != nil {
		return "", err
	}
	return filepath.Join(cfg.GOMODCACHE, "cache/download", enc, "/@v"), nil
}

func CachePath(m module.Version, suffix string) (string, error) {
	dir, err := cacheDir(m.Path)
	if err != nil {
		return "", err
	}
	if !semver.IsValid(m.Version) {
		return "", fmt.Errorf("non-semver module version %q", m.Version)
	}
	if module.CanonicalVersion(m.Version) != m.Version {
		return "", fmt.Errorf("non-canonical module version %q", m.Version)
	}
	encVer, err := module.EscapeVersion(m.Version)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, encVer+"."+suffix), nil
}
