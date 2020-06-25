/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package gopkg

import (
	"fmt"
	"io"
	"log"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/qiniu/goplus/exec.spec"
)

// -----------------------------------------------------------------------------

type exportedFunc struct {
	name string
	exec string
}

// Exporter represents a go package exporter.
type Exporter struct {
	w          io.Writer
	pkgDot     string
	pkgPath    string
	pkgExport  string
	execs      []string
	imports    map[string]string // pkgPath => pkg
	importPkgs map[string]string // pkg => pkgPath
	exportFns  []exportedFunc
	exportFnvs []exportedFunc
}

// NewExporter creates a go package exporter.
func NewExporter(w io.Writer, pkg, pkgPath string) *Exporter {
	const gopPath = "github.com/qiniu/goplus/gop"
	imports := map[string]string{gopPath: "gop"}
	importPkgs := map[string]string{"gop": gopPath}
	p := &Exporter{w: w, pkgExport: pkg, pkgPath: pkgPath, imports: imports, importPkgs: importPkgs}
	p.pkgDot = p.importPkg(pkg, pkgPath) + "."
	return p
}

func (p *Exporter) importPkg(pkg, pkgPath string) string {
	if name, ok := p.imports[pkgPath]; ok {
		return name
	}
	if pkg == "" {
		pkg = path.Base(pkgPath)
		pos := strings.IndexAny(pkg, ".-")
		if pos >= 0 {
			pkg = pkg[:pos]
		}
		if pkg == "" {
			pkg = "p"
		}
	}
	n := len(pkg)
	idx := 1
	for {
		if _, ok := p.importPkgs[pkg]; !ok {
			break
		}
		pkg = pkg[:n] + strconv.Itoa(idx)
		idx++
	}
	p.imports[pkgPath] = pkg
	p.importPkgs[pkg] = pkgPath
	return pkg
}

// ExportFunc exports a go function/method.
func (p *Exporter) ExportFunc(name string, fn interface{}) {
	tfn := reflect.TypeOf(fn)
	isVariadic := tfn.IsVariadic()
	isMethod := isMethod(name)
	numIn := tfn.NumIn()
	numOut := tfn.NumOut()
	args := make([]string, numIn)
	var arityName, arity, fnName, retAssign, retReturn string
	if isVariadic {
		arityName, arity = "arity", "arity"
		numIn--
	} else {
		arityName, arity = "_", strconv.Itoa(numIn)
	}
	if numOut > 0 {
		retOut := make([]string, numOut)
		for i := 0; i < numOut; i++ {
			retOut[i] = "ret" + strconv.Itoa(i)
		}
		retAssign = strings.Join(retOut, ", ")
		retReturn = arity + ", " + retAssign
		retAssign += " := "
	} else {
		retReturn = arity
	}
	for i := 0; i < numIn; i++ {
		t := tfn.In(i)
		args[i] = fmt.Sprintf("args[%d].(%s)", i, t.String())
	}
	if isVariadic {
		var varg string
		if numIn == 0 {
			varg = "args"
		} else {
			varg = fmt.Sprintf("args[%d:]", numIn)
		}
		switch tyElem := tfn.In(numIn).Elem(); tyElem {
		case exec.TyEmptyInterface:
		case exec.TyString:
			varg = "gop.ToStrings(" + varg + ")"
		default:
			log.Panicf("ExportFunc: unsupported type - ...%v\n", tyElem)
		}
		args[numIn] = varg + "..."
	}
	if isMethod {
		pos := strings.Index(name, ".")
		fnName = args[0] + "." + name[pos+1:]
		args = args[1:]
	} else {
		fnName = p.pkgDot + name
	}
	exec := strings.Map(skipNsymch, name)
	repl := strings.NewReplacer(
		"$name", exec,
		"$ariName", arityName,
		"$arity", arity,
		"$args", strings.Join(args, ", "),
		"$retAssign", retAssign,
		"$retReturn", retReturn,
		"$fn", fnName,
	)
	p.execs = append(p.execs, repl.Replace(`
func exec$name($ariName int, p *gop.Context) {
	args := p.GetArgs($arity)
	$retAssign$fn($args)
	p.Ret($retReturn)
}
`))
	exported := exportedFunc{name: name, exec: exec}
	if isVariadic {
		p.exportFnvs = append(p.exportFnvs, exported)
	} else {
		p.exportFns = append(p.exportFns, exported)
	}
}

func isMethod(name string) bool {
	return strings.HasPrefix(name, "(")
}

func skipNsymch(c rune) rune {
	if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_' || c >= 0x80 {
		return c
	}
	return -1
}

func exportFns(w io.Writer, pkgDot string, fns []exportedFunc, tag string) {
	if len(fns) == 0 {
		return
	}
	fmt.Fprintf(w, `	I.Register%ss(
`, tag)
	for _, fn := range fns {
		name := fn.name
		if isMethod(name) {
			name = name[:2] + pkgDot + name[2:]
		} else {
			name = pkgDot + name
		}
		fmt.Fprintf(w, `		I.%s("%s", %s, exec%s),
`, tag, fn.name, name, fn.exec)
	}
	fmt.Fprintf(w, "	)\n")
}

const gopkgInitExportHeader = `
// I is a Go package instance.
var I = gop.NewGoPackage("%s")

func init() {
`

const gopkgInitExportFooter = `}
`

const gopkgExportHeader = `package %s

import (
`

const gopkgExportFooter = `)
`

// Close finishes go package export.
func (p *Exporter) Close() error {
	pkgs := make([]string, 0, len(p.importPkgs))
	for pkg := range p.importPkgs {
		pkgs = append(pkgs, pkg)
	}
	sort.Strings(pkgs)
	fmt.Fprintf(p.w, gopkgExportHeader, p.pkgExport)
	for _, pkg := range pkgs {
		pkgPath := p.importPkgs[pkg]
		fmt.Fprintf(p.w, `	%s "%s"
`, pkg, pkgPath)
	}
	fmt.Fprintf(p.w, gopkgExportFooter)
	for _, exec := range p.execs {
		io.WriteString(p.w, exec)
	}
	fmt.Fprintf(p.w, gopkgInitExportHeader, p.pkgPath)
	pkgDot := p.pkgDot
	exportFns(p.w, pkgDot, p.exportFns, "Func")
	exportFns(p.w, pkgDot, p.exportFnvs, "Funcv")
	fmt.Fprintf(p.w, gopkgInitExportFooter)
	return nil
}

// -----------------------------------------------------------------------------
