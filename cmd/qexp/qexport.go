package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/qiniu/goplus/cmd/qexp/gopkg"
	"github.com/qiniu/x/log"
)

func createExportFile(pkgDir string) (f io.WriteCloser, err error) {
	os.MkdirAll(pkgDir, 0777)
	return os.Create(pkgDir + "/gomod_export.go")
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: qexp <goPkgPath>\n")
		flag.PrintDefaults()
		return
	}
	pkgPath := flag.Arg(0)
	err := gopkg.Export(pkgPath, createExportFile)
	if err != nil {
		log.Panicln("export failed:", err)
		os.Exit(1)
	}
}

// -----------------------------------------------------------------------------
