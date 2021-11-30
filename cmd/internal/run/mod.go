package run

import (
	"path/filepath"

	"github.com/goplus/gop/env"
	"github.com/qiniu/x/log"
)

func findGoModFile(dir string) (modfile string, noCacheFile bool, err error) {
	modfile, err = env.GOPMOD(dir)
	if err != nil {
		modfile = filepath.Join(env.GOPROOT(), "go.mod")
		return modfile, true, nil
	}
	return
}

func findGoModDir(dir string) (string, bool) {
	modfile, nocachefile, err := findGoModFile(dir)
	if err != nil {
		log.Fatalln("findGoModFile:", err)
	}
	return filepath.Dir(modfile), nocachefile
}
