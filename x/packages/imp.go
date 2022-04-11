package packages

import (
	"go/types"

	"github.com/goplus/gop/env"
)

type Config struct {
	// Dir specifies working directory. Use current directory if empty.
	Dir string
}

type Importer struct {
}

func NewImporter(conf *Config) (p *Importer, err error) {
	if conf == nil {
		conf = new(Config)
	}
	modfile, err := env.GOPMOD(conf.Dir)
	if err != nil {
		return
	}
	_ = modfile
	return
}

func (p *Importer) Import(pkgPath string) (pkg *types.Package, err error) {
	return
}
