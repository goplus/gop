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

func exportFunc(o *types.Func, prefix string) (err error) {
	fmt.Printf("%sfunc %s\n", prefix, o.Name())
	return nil
}

func exportTypeName(o *types.TypeName) (err error) {
	t := o.Type().(*types.Named)
	fmt.Printf("== type %s (%v) ==\n", o.Name(), t)
	n := t.NumMethods()
	for i := 0; i < n; i++ {
		m := t.Method(i)
		if !m.Exported() {
			continue
		}
		exportFunc(m, "  ")
	}
	return nil
}

func exportVar(o *types.Var) (err error) {
	fmt.Printf("==> var %s\n", o.Name())
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
		if !obj.Exported() {
			continue
		}
		switch o := obj.(type) {
		case *types.Var:
			err = exportVar(o)
		case *types.TypeName:
			err = exportTypeName(o)
		case *types.Func:
			err = exportFunc(o, "")
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
