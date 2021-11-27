package modload

import (
	"context"
	"log"
	"path/filepath"

	"golang.org/x/mod/module"

	"github.com/goplus/gop/cmd/internal/modfetch"
	"github.com/goplus/gop/x/mod/modfile"
)

func LoadClassFile() {
	if ModFile.Register == nil {
		return
	}

	var dir string
	var err error
	var claassMod module.Version

	for _, require := range ModFile.Require {
		if require.Mod.Path == ModFile.Register.ClassfileMod {
			claassMod = require.Mod
		}
	}
	for _, replace := range ModFile.Replace {
		if replace.Old.Path == ModFile.Register.ClassfileMod {
			claassMod = replace.New
		}
	}

	if claassMod.Version != "" {
		dir, err = modfetch.Download(claassMod)
		if err != nil {
			log.Fatalf("gop: download classsfile module error %v", err)
		}
	} else {
		dir = claassMod.Path
	}

	if dir == "" {
		log.Fatalf("gop: can't find classfile path in require statment")
	}

	gopmod := filepath.Join(dir, "gop.mod")
	data, err := modfetch.Read(gopmod)
	if err != nil {
		log.Fatalf("gop: %v", err)
	}

	var fixed bool
	f, err := modfile.Parse(gopmod, data, fixVersion(context.Background(), &fixed))
	if err != nil {
		// Errors returned by modfile.Parse begin with file:line.
		log.Fatalf("go: errors parsing go.mod:\n%s\n", err)
	}
	ClassModFile = f
}
