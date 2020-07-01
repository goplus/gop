package gopkg

import (
	"bytes"
	"fmt"
	"go/importer"
	"go/types"
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

func execmReplacerReplace(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*strings.Replacer).Replace(args[1].(string))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("strings")

func init() {
	I.RegisterFuncs(
		I.Func("(*Replacer).Replace", (*strings.Replacer).Replace, execmReplacerReplace),
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

func TestFixPkgString(t *testing.T) {
	pkg, _ := importer.Default().Import("go/types")
	b := bytes.NewBuffer(nil)
	e := NewExporter(b, pkg)
	e.imports["io"] = "io1"
	e.imports["go/types"] = "types1"
	if v := e.fixPkgString("*go/types.Interface"); v != "*types1.Interface" {
		t.Fatal(v)
	}
	if v := e.fixPkgString("[]*go/types.Package"); v != "[]*types1.Package" {
		t.Fatal(v)
	}
	if v := e.fixPkgString("io.Writer"); v != "io1.Writer" {
		t.Fatal(v)
	}
}

// -----------------------------------------------------------------------------

func TestExportStrings(t *testing.T) {
	err := Export("strings", ioutil.Discard)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

func TestExportStrconv(t *testing.T) {
	err := Export("strconv", ioutil.Discard)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

func TestExportReflect(t *testing.T) {
	err := Export("reflect", ioutil.Discard)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

func TestExportIo(t *testing.T) {
	err := Export("io", os.Stdout)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

func TestExportGoTypes(t *testing.T) {
	err := Export("go/types", ioutil.Discard)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

func TestExportMath(t *testing.T) {
	err := Export("math", ioutil.Discard)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

func TestImport(t *testing.T) {
	pkg, err := Import("go/types")
	if err != nil {
		t.Fatal("TestImport failed:", err)
	}
	if pkg.Path() != "go/types" {
		t.Fatal(pkg.Path())
	}
}

func TestBadImport(t *testing.T) {
	_, err := Import("go/badtypes")
	if err == nil {
		t.Fatal("TestBadImport failed")
	}
}

func TestBadExport(t *testing.T) {
	err := Export("go/badtypes", ioutil.Discard)
	if err == nil {
		t.Fatal("TestBadExport failed")
	}
}

func TestExportPackage(t *testing.T) {
	pkg, err := Import("go/build")
	if err != nil {
		t.Fatal("TestPackage Import failed:", err)
	}
	err = ExportPackage(pkg, ioutil.Discard)
	if err != nil {
		t.Fatal("TestPackage ExportPackage failed:", err)
	}
}

// -----------------------------------------------------------------------------
