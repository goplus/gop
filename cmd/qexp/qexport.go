package main

import (
	"flag"
	"fmt"
	"go/importer"
	"go/types"
	"os"
	"reflect"
	"strings"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

func filterGotest(fi os.FileInfo) bool {
	name := fi.Name()
	if strings.HasSuffix(name, "_test.go") {
		return false
	}
	if strings.HasPrefix(name, "_") {
		return false
	}
	return true
}

func exportFunc(o *types.Func) (err error) {
	fmt.Printf("func %s\n", o.Name())
	return nil
}

func export(pkgPath string) (err error) {
	pkg, err := importer.Default().Import(pkgPath)
	if err != nil {
		return
	}
	gbl := pkg.Scope()
	names := gbl.Names()
	for _, name := range names {
		obj := gbl.Lookup(name)
		switch o := obj.(type) {
		case *types.Var:
		case *types.TypeName:
		case *types.Func:
			if o.Exported() {
				err = exportFunc(o)
			}
		default:
			log.Panicln("export failed: unknown type -", reflect.TypeOf(o))
		}
		if err != nil {
			return
		}
	}
	return nil
}

// -----------------------------------------------------------------------------

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: qexp <goPkgPath>\n")
		flag.PrintDefaults()
		return
	}
	pkgPath := flag.Arg(0)
	err := export(pkgPath)
	if err != nil {
		log.Panicln("export failed:", err)
		os.Exit(1)
	}
}

// -----------------------------------------------------------------------------
