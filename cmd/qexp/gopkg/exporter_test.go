package gopkg

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------

type goFunc struct {
	Name string
	This interface{}
}

const expected = `package strings

import (
	fmt "fmt"
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

func Test(t *testing.T) {
	exports := []goFunc{
		{"NewReplacer", strings.NewReplacer},
		{"(*Replacer).Replace", (*strings.Replacer).Replace},
	}
	b := bytes.NewBuffer(nil)
	e := NewExporter(b, "strings", "strings")
	e.importPkg("", "fmt")
	for _, gof := range exports {
		e.ExportFunc(gof.Name, gof.This)
	}
	e.Close()
	if real := b.String(); real != expected {
		fmt.Println(real)
		t.Fatal("Test failed")
	}
}

// -----------------------------------------------------------------------------
