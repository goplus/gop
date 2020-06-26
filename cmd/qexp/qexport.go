package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/qiniu/goplus/cmd/qexp/gopkg"
	"github.com/qiniu/x/log"
)

var (
	exportFile string
)

func createExportFile(pkgDir string) (f io.WriteCloser, err error) {
	os.MkdirAll(pkgDir, 0777)
	exportFile = pkgDir + "/gomod_export.go"
	return os.Create(exportFile)
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: qexp <goPkgPath>\n")
		flag.PrintDefaults()
		return
	}
	pkgPath := flag.Arg(0)
	defer func() {
		if exportFile != "" {
			os.Remove(exportFile)
		}
	}()
	err := gopkg.Export(pkgPath, createExportFile)
	if err != nil {
		log.Panicln("export failed:", err)
	}
	exportFile = "" // don't remove file if success
}

// -----------------------------------------------------------------------------
