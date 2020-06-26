package gopkg

import (
	"bytes"
	"fmt"
	"go/importer"
	"go/types"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

// -----------------------------------------------------------------------------

type goFunc struct {
	Name string
	This interface{}
}

const expected = `// Package strings provide Go+ "strings" package, as "strings" package in Go.
package strings

import (
	gop "github.com/qiniu/goplus/gop"
	strings "strings"
)

func execNewReplacer(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := strings.NewReplacer(gop.ToStrings(args)...)
	p.Ret(arity, ret0)
}

func execReplacerReplace(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*strings.Replacer).Replace(args[1].(string))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("strings")

func init() {
	I.RegisterFuncs(
		I.Func("(*Replacer).Replace", (*strings.Replacer).Replace, execReplacerReplace),
	)
	I.RegisterFuncvs(
		I.Funcv("NewReplacer", strings.NewReplacer, execNewReplacer),
	)
}
`

func getMethod(o *types.Named, name string) *types.Func {
	n := o.NumMethods()
	for i := 0; i < n; i++ {
		m := o.Method(i)
		if m.Name() == name {
			return m
		}
	}
	return nil
}

func TestBasic(t *testing.T) {
	pkg, err := importer.Default().Import("strings")
	if err != nil {
		t.Fatal("Import failed:", err)
	}
	gbl := pkg.Scope()
	newReplacer := gbl.Lookup("NewReplacer").(*types.Func)
	replacer := gbl.Lookup("Replacer").(*types.TypeName).Type().(*types.Named)
	b := bytes.NewBuffer(nil)
	e := NewExporter(b, pkg)
	e.ExportFunc(newReplacer)
	e.ExportFunc(getMethod(replacer, "Replace"))
	e.Close()
	if real := b.String(); real != expected {
		fmt.Println(real)
		t.Fatal("Test failed")
	}
}

// -----------------------------------------------------------------------------

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func nilExportFile(pkgDir string) (f io.WriteCloser, err error) {
	return &nopCloser{ioutil.Discard}, nil
}

func stdoutExportFile(pkgDir string) (f io.WriteCloser, err error) {
	return &nopCloser{os.Stdout}, nil
}

func TestExportStrings(t *testing.T) {
	err := Export("strings", nilExportFile)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

func TestExportStrconv(t *testing.T) {
	err := Export("strconv", nilExportFile)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

func TestExportReflect(t *testing.T) {
	err := Export("reflect", nilExportFile)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

func TestExportIo(t *testing.T) {
	err := Export("io", stdoutExportFile)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

func TestExportGoTypes(t *testing.T) {
	err := Export("go/types", nilExportFile)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

// -----------------------------------------------------------------------------
