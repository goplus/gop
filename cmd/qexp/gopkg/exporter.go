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
	"go/constant"
	"go/types"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// -----------------------------------------------------------------------------

type exportedFunc struct {
	name string
	exec string
}

type exportedConst struct {
	name string
	kind string
	val  string
}

// Exporter represents a go package exporter.
type Exporter struct {
	w             io.Writer
	pkg           *types.Package
	pkgDot        string
	execs         []string
	toTypes       []types.Type
	toSlices      []types.Type
	imports       map[string]string // pkgPath => pkg
	importPkgs    map[string]string // pkg => pkgPath
	exportFns     []exportedFunc
	exportFnvs    []exportedFunc
	exportConsts  []exportedConst
	hasReflectPkg bool
	hasSpecPkg    bool
}

// NewExporter creates a go package exporter.
func NewExporter(w io.Writer, pkg *types.Package) *Exporter {
	const gopPath = "github.com/qiniu/goplus/gop"
	const specPath = "github.com/qiniu/goplus/exec.spec"
	const reflectPath = "reflect"
	imports := map[string]string{gopPath: "gop", specPath: "spec", reflectPath: "reflect"}
	importPkgs := map[string]string{"gop": gopPath, "spec": specPath, "reflect": reflectPath}
	p := &Exporter{w: w, pkg: pkg, imports: imports, importPkgs: importPkgs}
	p.pkgDot = p.importPkg(pkg) + "."
	return p
}

func (p *Exporter) importPkg(pkgObj *types.Package) string {
	pkgPath := pkgObj.Path()
	if name, ok := p.imports[pkgPath]; ok {
		return name
	}
	pkg := pkgObj.Name()
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

func (p *Exporter) useType(typ types.Type) {
	switch t := typ.(type) {
	case *types.Basic:
		if t.Kind() == types.UnsafePointer {
			p.imports["unsafe"] = "unsafe"
			p.importPkgs["unsafe"] = "unsafe"
		}
	case *types.Pointer:
		p.useType(t.Elem())
	case *types.Slice:
		p.useType(t.Elem())
	case *types.Map:
		p.useType(t.Key())
		p.useType(t.Elem())
	case *types.Chan:
		p.useType(t.Elem())
	case *types.Array:
		p.useType(t.Elem())
	case *types.Struct:
		n := t.NumFields()
		for i := 0; i < n; i++ {
			p.useType(t.Field(i).Type())
		}
	case *types.Signature:
		p.useType(t.Params())
		p.useType(t.Results())
	case *types.Tuple:
		n := t.Len()
		for i := 0; i < n; i++ {
			p.useType(t.At(i).Type())
		}
	case *types.Named:
		if pkg := t.Obj().Pkg(); pkg != nil {
			p.importPkg(pkg)
		}
	case *types.Interface:
		n := t.NumMethods()
		for i := 0; i < n; i++ {
			m := t.Method(i)
			p.useType(m.Type())
		}
	default:
		panic("not here")
	}
}

func (p *Exporter) toType(typ types.Type) string {
	for i, t := range p.toTypes {
		if types.Identical(typ, t) {
			return toTypeName(i)
		}
	}
	idx := toTypeName(len(p.toTypes))

	var typStr = p.typeString(typ)
	p.execs = append(p.execs, fmt.Sprintf(`
func %s(v interface{}) %s {
	if v == nil {
		return nil
	}
	return v.(%s)
}
`, idx, typStr, typStr))
	p.toTypes = append(p.toTypes, typ)
	return idx
}

func toTypeName(i int) string {
	return "toType" + strconv.Itoa(i)
}

func toSliceName(i int) string {
	return "toSlice" + strconv.Itoa(i)
}

func (p *Exporter) toSlice(tyElem types.Type) string {
	for i, t := range p.toSlices {
		if types.Identical(tyElem, t) {
			return toSliceName(i)
		}
	}
	idx := toSliceName(len(p.toSlices))
	typCast := p.typeCast("arg", tyElem)
	typStr := p.typeString(tyElem)
	p.execs = append(p.execs, fmt.Sprintf(`
func %s(args []interface{}) []%v {
	ret := make([]%v, len(args))
	for i, arg := range args {
		ret[i] = %s
	}
	return ret
}
`, idx, typStr, typStr, typCast))
	p.toSlices = append(p.toSlices, tyElem)
	return idx
}

func (p *Exporter) sliceCast(varg string, tyElem types.Type) string {
	if e, ok := tyElem.(*types.Basic); ok {
		uName := strings.Title(e.Name())
		varg = "gop.To" + uName + "s(" + varg + ")"
	} else {
		tyElemIntf, isInterface := tyElem.Underlying().(*types.Interface)
		if !(isInterface && tyElemIntf.Empty()) { // is not empty interface
			varg = p.toSlice(tyElem) + "(" + varg + ")"
		}
	}
	return varg
}

func (p *Exporter) typeCast(varg string, typ types.Type) string {
	typIntf, isInterface := typ.Underlying().(*types.Interface)
	if isInterface {
		if typIntf.Empty() {
			return varg
		}
		return p.toType(typ) + "(" + varg + ")"
	}
	typStr := p.typeString(typ)
	return varg + ".(" + typStr + ")"
}

var (
	reTyp, _ = regexp.Compile("[\\w\\./]+")
)

func (p *Exporter) typeString(typ types.Type) string {
	typStr := typ.String()
	if !strings.Contains(typStr, ".") {
		return typStr
	}
	return p.fixPkgString(typ.String())
}

func (p *Exporter) fixPkgString(typ string) string {
	return reTyp.ReplaceAllStringFunc(typ, func(s string) string {
		pos := strings.Index(s, ".")
		if pos > 0 {
			pkg := s[:pos]
			if r, ok := p.imports[pkg]; ok {
				return r + s[pos:]
			}
		}
		return s
	})
}

// ExportFunc exports a go function/method.
func (p *Exporter) ExportFunc(fn *types.Func) error {
	tfn := fn.Type().(*types.Signature)
	isVariadic := tfn.Variadic()
	isMethod := tfn.Recv() != nil
	numIn := tfn.Params().Len()
	numOut := tfn.Results().Len()
	args := make([]string, numIn)
	from := 0
	if isMethod {
		from = 1
	}
	var arityName, arity, fnName, retAssign, retReturn string
	if isVariadic {
		arityName, arity = "arity", "arity"
		numIn--
	} else {
		arityName, arity = "_", strconv.Itoa(numIn+from)
	}
	if numOut > 0 {
		retOut := make([]string, numOut)
		for i := 0; i < numOut; i++ {
			retOut[i] = "ret" + strconv.Itoa(i)
		}
		retAssign = strings.Join(retOut, ", ")
		retReturn = "Ret(" + arity + ", " + retAssign + ")"
		retAssign += " := "
	} else {
		retReturn = "PopN(" + arity + ")"
	}
	for i := 0; i < numIn; i++ {
		typ := tfn.Params().At(i).Type()
		p.useType(typ)
		args[i] = p.typeCast("args["+strconv.Itoa(i+from)+"]", typ)
	}
	if isVariadic {
		var varg string
		if numIn == 0 {
			varg = "args"
		} else {
			varg = fmt.Sprintf("args[%d:]", numIn)
		}
		tyElem := tfn.Params().At(numIn).Type().(*types.Slice).Elem()
		p.useType(tyElem)
		args[numIn] = p.sliceCast(varg, tyElem) + "..."
	}
	name := fn.Name()
	exec := name
	if isMethod {
		fullName := fn.FullName()
		exec = typeName(tfn.Recv().Type()) + name
		name = withoutPkg(fullName)
		fnName = "args[0]." + withPkg(p.pkgDot, name)
	} else {
		fnName = p.pkgDot + name
	}
	var argsAssign string
	if arity != "0" {
		argsAssign = "	args := p.GetArgs(" + arity + ")\n"
	}
	if isMethod {
		exec = "execm" + exec
	} else {
		exec = "exec" + exec
	}
	repl := strings.NewReplacer(
		"$execFunc", exec,
		"$ariName", arityName,
		"$args", strings.Join(args, ", "),
		"$argInit", argsAssign,
		"$retAssign", retAssign,
		"$retReturn", retReturn,
		"$fn", fnName,
	)
	p.execs = append(p.execs, repl.Replace(`
func $execFunc($ariName int, p *gop.Context) {
$argInit	$retAssign$fn($args)
	p.$retReturn
}
`))
	exported := exportedFunc{name: name, exec: exec}
	if isVariadic {
		p.exportFnvs = append(p.exportFnvs, exported)
	} else {
		p.exportFns = append(p.exportFns, exported)
	}
	return nil
}

// ExportConst exports a go consts.
func (p *Exporter) ExportConst(typ *types.Const) error {
	kind, err := constKind(typ, "spec")
	if err != nil {
		return err
	}
	if strings.HasPrefix(kind, "reflect.") {
		p.hasReflectPkg = true
	}
	if strings.HasPrefix(kind, "spec.") {
		p.hasSpecPkg = true
	}
	fullName := p.pkgDot + typ.Name()
	var c exportedConst
	c.name = typ.Name()
	c.kind = kind
	c.val = fullName
	if kind == "spec.ConstUnboundInt" && typ.Val().Kind() == constant.Int {
		value := typ.Val().String()
		_, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			if value[0] == '-' {
				c.kind = "spec.Int64"
				c.val = "int64(" + fullName + ")"
			} else {
				c.kind = "spec.Uint64"
				c.val = "uint64(" + fullName + ")"
			}
		}
	}
	p.exportConsts = append(p.exportConsts, c)
	return nil
}

func constKind(typ *types.Const, pkg string) (string, error) {
	baisc, ok := typ.Type().Underlying().(*types.Basic)
	if !ok {
		return "", fmt.Errorf("unparse basic of const %v", typ)
	}
	switch baisc.Kind() {
	case types.UntypedBool:
		return "reflect.Bool", nil
	case types.UntypedInt:
		return pkg + ".ConstUnboundInt", nil
	case types.UntypedRune:
		return pkg + ".ConstBoundRune", nil
	case types.UntypedFloat:
		return pkg + ".ConstUnboundFloat", nil
	case types.UntypedComplex:
		return pkg + ".ConstUnboundComplex", nil
	case types.UntypedString:
		return pkg + ".ConstBoundString", nil
	case types.UntypedNil:
		return pkg + ".ConstUnboundPtr", nil
	case types.Byte:
		return "reflect.Uint8", nil
	case types.Rune:
		return "reflect.Uint32", nil
	}
	return "reflect." + strings.Title(baisc.Name()), nil
}

func withoutPkg(fullName string) string {
	pos := strings.Index(fullName, ")")
	if pos < 0 {
		return fullName
	}
	dot := strings.Index(fullName[:pos], ".")
	if dot < 0 {
		return fullName
	}
	start := strings.IndexFunc(fullName[:dot], func(c rune) bool {
		return c != '(' && c != '*'
	})
	if start < 0 {
		return fullName
	}
	return fullName[:start] + fullName[dot+1:]
}

func typeName(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Pointer:
		return typeName(t.Elem())
	case *types.Named:
		return t.Obj().Name()
	}
	panic("not here")
}

func isMethod(name string) bool {
	return strings.HasPrefix(name, "(")
}

func withPkg(pkgDot, name string) string {
	if isMethod(name) {
		n := len(name) - len(strings.TrimLeft(name[1:], "*"))
		return name[:n] + pkgDot + name[n:]
	}
	return pkgDot + name
}

func registerFns(w io.Writer, pkgDot string, fns []exportedFunc, tag string) {
	if len(fns) == 0 {
		return
	}
	fmt.Fprintf(w, `	I.Register%ss(
`, tag)
	for _, fn := range fns {
		name := withPkg(pkgDot, fn.name)
		fmt.Fprintf(w, `		I.%s("%s", %s, %s),
`, tag, fn.name, name, fn.exec)
	}
	fmt.Fprintf(w, "	)\n")
}

func registerConsts(w io.Writer, consts []exportedConst) {
	if len(consts) == 0 {
		return
	}
	fmt.Fprintf(w, `	I.RegisterConsts(
`)
	for _, c := range consts {
		fmt.Fprintf(w, `		I.Const("%s", %s, %s),
`, c.name, c.kind, c.val)
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

const gopkgExportHeader = `// Package %s provide Go+ "%s" package, as "%s" package in Go.
package %s

import (
`

const gopkgExportFooter = `)
`

// Close finishes go package export.
func (p *Exporter) Close() error {
	pkgs := make([]string, 0, len(p.importPkgs))
	for pkg := range p.importPkgs {
		if pkg == "reflect" && !p.hasReflectPkg {
			continue
		}
		if pkg == "spec" && !p.hasSpecPkg {
			continue
		}
		pkgs = append(pkgs, pkg)
	}
	sort.Strings(pkgs)
	pkg, pkgPath := p.pkg.Name(), p.pkg.Path()
	fmt.Fprintf(p.w, gopkgExportHeader, pkg, pkgPath, pkgPath, pkg)
	for _, pkg := range pkgs {
		pkgPath := p.importPkgs[pkg]
		fmt.Fprintf(p.w, `	%s "%s"
`, pkg, pkgPath)
	}
	fmt.Fprintf(p.w, gopkgExportFooter)
	for _, exec := range p.execs {
		io.WriteString(p.w, exec)
	}
	fmt.Fprintf(p.w, gopkgInitExportHeader, pkgPath)
	pkgDot := p.pkgDot
	registerFns(p.w, pkgDot, p.exportFns, "Func")
	registerFns(p.w, pkgDot, p.exportFnvs, "Funcv")
	registerConsts(p.w, p.exportConsts)
	fmt.Fprintf(p.w, gopkgInitExportFooter)
	return nil
}

// -----------------------------------------------------------------------------
