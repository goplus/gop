package gopkg

import (
	"bytes"
	"fmt"
	"go/format"
	"go/importer"
	"go/types"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

type goFunc struct {
	Name string
	This interface{}
}

const expected = `// Package strings provide Go+ "strings" package, as "strings" package in Go.
package strings

import (
	strings "strings"

	gop "github.com/goplus/gop"
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
	e.ExportFunc(newReplacer, "")
	e.ExportFunc(getMethod(replacer, "Replace"), "")
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
	e.imports["go/types-pkg1"] = "pkg1"
	e.imports["go/types_pkg2"] = "pkg2"
	e.imports["golang.org/x/crypto/chacha20"] = "chacha20"
	if v := e.fixPkgString("*go/types.Interface"); v != "*types1.Interface" {
		t.Fatal(v)
	}
	if v := e.fixPkgString("[]*go/types.Package"); v != "[]*types1.Package" {
		t.Fatal(v)
	}
	if v := e.fixPkgString("*go/types-pkg1.Interface"); v != "*pkg1.Interface" {
		t.Fatal(v)
	}
	if v := e.fixPkgString("*go/types_pkg2.Interface"); v != "*pkg2.Interface" {
		t.Fatal(v)
	}
	if v := e.fixPkgString("io.Writer"); v != "io1.Writer" {
		t.Fatal(v)
	}
	if v := e.fixPkgString("*golang.org/x/crypto/chacha20.Cipher"); v != "*chacha20.Cipher" {
		t.Fatal(v)
	}
}

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
	err := Export("io", ioutil.Discard)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

func TestExportIoUtil(t *testing.T) {
	err := Export("io/ioutil", ioutil.Discard)
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

func TestExportGopQ(t *testing.T) {
	err := Export("github.com/goplus/gop/ast/gopq", ioutil.Discard)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

func TestExportX(t *testing.T) {
	err := Export("github.com/qiniu/x/log", ioutil.Discard)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

func TestExportNotFound(t *testing.T) {
	err := Export("github.com/qiniu", ioutil.Discard)
	if err != ErrInvalidPkgPath {
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

func TestLooupMod(t *testing.T) {
	dir, err := LookupMod("github.com/qiniu/x/log")
	if err != nil {
		t.Fatal("LookupMod failed:", err)
	}
	i := strings.Index(dir, "github.com")
	dir1, err := LookupMod(filepath.ToSlash(dir[i:]))
	if err != nil {
		t.Fatal("LookupMod failed:", err)
	}
	if dir != dir1 {
		t.Fatal("LookupMod failed:", dir, dir1)
	}
}

func TestIgnorePkg(t *testing.T) {
	if !isInternalPkg("golang.org/x/tools/internal/span") {
		t.Fatal("is internal pkg")
	}
	if isInternalPkg("golang.org/x/tools") {
		t.Fatal("not internal pkg")
	}
}

func TestWithoutPkg(t *testing.T) {
	if withoutPkg("github.com/qiniu/x/log") != "github.com/qiniu/x/log" {
		t.Fatal()
	}
}

func TestBadLooupMod(t *testing.T) {
	dir, err := LookupMod("github.com/qiniu/x/log2")
	if err == nil {
		t.Fatal("LookupMod bad:", dir)
	}
	dir, err = LookupMod("github.com/qiniu/x/doc.go")
	if err == nil {
		t.Fatal("LookupMod bad:", dir)
	}
	dir, err = LookupMod("unicode")
	if err == nil {
		t.Fatal("LookupMod bad:", dir)
	}
}

func TestExportVersion(t *testing.T) {
	pkgPath := "github.com/qiniu/x/log"
	dir, err := LookupMod(pkgPath)
	if err != nil {
		t.Fatal("LookupMod failed:", err)
	}
	pkgPath = filepath.ToSlash(dir[len(getModRoot())+1:])
	err = Export(pkgPath, ioutil.Discard)
	if err != nil {
		t.Fatal("TestExport failed:", err)
	}
}

func TestParsePkgVer(t *testing.T) {
	pkg, mod, sub := ParsePkgVer("github.com/qiniu/x@v1.11.5/log")
	if pkg != "github.com/qiniu/x/log" {
		t.Fatal("ParsePkgVer failed:", pkg)
	}
	if mod != "github.com/qiniu/x@v1.11.5" {
		t.Fatal("ParsePkgVer failed:", mod)
	}
	if sub != "log" {
		t.Fatal("ParsePkgVer failed:", sub)
	}
	pkg, mod, sub = ParsePkgVer("github.com/qiniu/x@v1.11.5")
	if pkg != "github.com/qiniu/x" {
		t.Fatal("ParsePkgVer failed:", pkg)
	}
	if mod != "github.com/qiniu/x@v1.11.5" {
		t.Fatal("ParsePkgVer failed:", mod)
	}
	if sub != "" {
		t.Fatal("ParsePkgVer failed:", sub)
	}
	pkg, mod, sub = ParsePkgVer("github.com/qiniu/x")
	if pkg != "github.com/qiniu/x" {
		t.Fatal("ParsePkgVer failed:", pkg)
	}
	if mod != "" {
		t.Fatal("ParsePkgVer failed:", mod)
	}
	if sub != "" {
		t.Fatal("ParsePkgVer failed:", sub)
	}
}

func TestImportSource(t *testing.T) {
	pkgPath := "github.com/qiniu/x/log"
	dir, err := LookupMod(pkgPath)
	if err != nil {
		t.Fatal("LookupMod failed:", err)
	}
	pkg, err := ImportSource(pkgPath, dir)
	if err != nil {
		t.Fatal("ImportSource failed:", err)
	}
	err = ExportPackage(pkg, ioutil.Discard)
	if err != nil {
		t.Fatal("ExportPackage failed:", err)
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

func TestExportStd(t *testing.T) {
	out, err := exec.Command("go", "list", "-e", "std").Output()
	if err != nil {
		log.Fatal(err)
	}
	pkgs := strings.Fields(string(out))
	for _, pkg := range pkgs {
		if pkg == "unsafe" {
			continue
		}
		if strings.Contains(pkg, "internal") {
			continue
		}
		if strings.HasPrefix(pkg, "vendor/") {
			continue
		}
		p, err := Import(pkg)
		if err != nil {
			t.Fatal("import error:", pkg, err)
		}
		var buf bytes.Buffer
		err = ExportPackage(p, &buf)
		if err != nil && err != ErrIgnore {
			t.Fatal("export error:", pkg, err)
		}
		_, err = format.Source(buf.Bytes())
		if err != nil {
			t.Log(buf.String())
			t.Fatal("format error:", pkg, err)
		}
	}
}

// -----------------------------------------------------------------------------
