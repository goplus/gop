package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/qiniu/goplus/cmd/qexp/gopkg"
	"github.com/qiniu/x/log"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: qexp <goPkgPath>\n")
		flag.PrintDefaults()
		return
	}
	pkgPath := flag.Arg(0)
	err := gopkg.Export(pkgPath, nil)
	if err != nil {
		log.Panicln("export failed:", err)
		os.Exit(1)
	}
}

// -----------------------------------------------------------------------------
