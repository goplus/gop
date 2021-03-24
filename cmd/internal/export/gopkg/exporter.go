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
	ftyp string
	exec string
}

type exportedVar struct {
	name string
	addr string
}

type exportedConst struct {
	name string
	kind string
	val  string
}

type exportedType struct {
	name string
	kind string
}

// Exporter represents a go package exporter.
type Exporter struct {
	w            io.Writer
	pkg          *types.Package
	pkgDot       string
	execs        []string
	toTypes      []types.Type
	toSlices     []types.Type
	imports      map[string]string // pkgPath => pkg
	importPkgs   map[string]string // pkg => pkgPath
	exportFns    []exportedFunc
	exportFnvs   []exportedFunc
	exportedVars []exportedVar
	exportConsts []exportedConst
	exportTypes  []exportedType
}

// NewExporter creates a go package exporter.
func NewExporter(w io.Writer, pkg *types.Package) *Exporter {
	const gopPath = "github.com/goplus/gop"
	imports := map[string]string{gopPath: "gop"}
	importPkgs := map[string]string{"gop": gopPath}
	p := &Exporter{w: w, pkg: pkg, imports: imports, importPkgs: importPkgs}
	p.pkgDot = p.importPkg(pkg) + "."
	return p
}

// IsEmpty checks if there is nothing to exmport or not.
func (p *Exporter) IsEmpty() bool {
	return len(p.exportFns) == 0 && len(p.exportFnvs) == 0 &&
		len(p.exportedVars) == 0 && len(p.exportConsts) == 0
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
	reTyp, _ = regexp.Compile("[\\w\\.\\-_/]+")
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
		pos := strings.LastIndex(s, ".")
		if pos > 0 {
			pkg := s[:pos]
			if r, ok := p.imports[pkg]; ok {
				return r + s[pos:]
			}
		}
		return s
	})
}

func isInternalPkg(pkg string) bool {
	for _, a := range strings.Split(pkg, "/") {
		if a == "internal" {
			return true
		}
	}
	return false
}

// ExportFunc exports a go function/method.
func (p *Exporter) ExportFunc(fn *types.Func, iname string) error {
	tfn := fn.Type().(*types.Signature)
	isVariadic := tfn.Variadic()
	isMethod := tfn.Recv() != nil
	isIMethod := iname != ""
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
		if named, ok := typ.(*types.Named); ok {
			if !named.Obj().Exported() && named.String() != "error" {
				fmt.Println("ignore", fn)
				return nil
			}
			if pkg := named.Obj().Pkg(); pkg != nil && isInternalPkg(pkg.Path()) {
				fmt.Println("ignore", fn)
				return nil
			}
		}
		p.useType(typ)
		args[i] = p.typeCast("args["+strconv.Itoa(i+from)+"]", typ)
	}
	if isVariadic {
		var varg string
		if numIn == 0 {
			varg = "args"
			if from == 1 {
				varg = "args[1:]"
			}
		} else {
			varg = fmt.Sprintf("args[%d:]", numIn+from)
		}
		tyElem := tfn.Params().At(numIn).Type().(*types.Slice).Elem()
		p.useType(tyElem)
		args[numIn] = p.sliceCast(varg, tyElem) + "..."
	}
	name := fn.Name()
	exec := name
	if isMethod {
		recv := tfn.Recv()
		fullName := fn.FullName()
		exec = typeName(recv.Type()) + name
		name = withoutPkg(fullName)
		fnName = "args[0]." + withPkg(p.pkgDot, name)
	} else {
		fnName = p.pkgDot + name
	}
	var argsAssign string
	if arity != "0" {
		argsAssign = "	args := p.GetArgs(" + arity + ")\n"
	}
	if isIMethod {
		exec = "execi" + exec
	} else if isMethod {
		exec = "execm" + exec
	} else {
		exec = "exec" + exec
	}
	ftyp := name
	var skipExec bool
	if isIMethod && tfn.Recv().Pkg() == p.pkg &&
		typeName(tfn.Recv().Type()) != iname {
		skipExec = true
		name = "(" + iname + ")." + fn.Name()
	}
	if !skipExec {
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
	}
	exported := exportedFunc{name: name, ftyp: ftyp, exec: exec}
	if isVariadic {
		p.exportFnvs = append(p.exportFnvs, exported)
	} else {
		p.exportFns = append(p.exportFns, exported)
	}
	return nil
}

// ExportVar exports a go var.
func (p *Exporter) ExportVar(typ *types.Var) {
	name := typ.Name()
	pkg := p.importPkg(typ.Pkg())
	addr := pkg + "." + name
	p.exportedVars = append(p.exportedVars, exportedVar{name, addr})
}

// ExportType export a go type
func (p *Exporter) ExportType(typ *types.TypeName) {
	name := typ.Name()
	kind := fmt.Sprintf("%v.TypeOf((*%v)(nil)).Elem()", p.importPkg(reflectPkg), p.importPkg(typ.Pkg())+"."+name)
	p.exportTypes = append(p.exportTypes, exportedType{name, kind})
}

var (
	qspecPkg   = types.NewPackage("github.com/goplus/gop/exec.spec", "qspec")
	reflectPkg = types.NewPackage("reflect", "reflect")
)

// ExportConst exports a go const.
func (p *Exporter) ExportConst(typ *types.Const) error {
	kind, err := constKind(typ)
	if err != nil {
		return err
	}
	pkg := p.importPkg(qspecPkg)
	val := p.pkgDot + typ.Name()
	if typ.Val().Kind() == constant.Int && kind == "ConstUnboundInt" {
		value := typ.Val().String()
		_, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			if value[0] == '-' {
				kind = "Int64"
				val = "int64(" + val + ")"
			} else {
				kind = "Uint64"
				val = "uint64(" + val + ")"
			}
		}
	}
	var c exportedConst
	c.name = typ.Name()
	c.kind = pkg + "." + kind
	c.val = val
	p.exportConsts = append(p.exportConsts, c)
	return nil
}

func constKind(typ *types.Const) (string, error) {
	baisc, ok := typ.Type().Underlying().(*types.Basic)
	if !ok {
		return "", fmt.Errorf("unparse basic of const %v", typ)
	}
	switch baisc.Kind() {
	case types.UntypedBool:
		return "Bool", nil
	case types.UntypedInt:
		return "ConstUnboundInt", nil
	case types.UntypedRune:
		return "ConstBoundRune", nil
	case types.UntypedFloat:
		return "ConstUnboundFloat", nil
	case types.UntypedComplex:
		return "ConstUnboundComplex", nil
	case types.UntypedString:
		return "ConstBoundString", nil
	case types.UntypedNil:
		return "ConstUnboundPtr", nil
	case types.Byte:
		return "Uint8", nil
	case types.Rune:
		return "Uint32", nil
	}
	return strings.Title(baisc.Name()), nil
}

func withoutPkg(fullName string) string {
	pos := strings.Index(fullName, ")")
	if pos < 0 {
		return fullName
	}
	dot := strings.LastIndex(fullName[:pos], ".")
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
		name := withPkg(pkgDot, fn.ftyp)
		fmt.Fprintf(w, `		I.%s("%s", %s, %s),
`, tag, fn.name, name, fn.exec)
	}
	fmt.Fprintf(w, "	)\n")
}

func registerTypes(w io.Writer, types []exportedType) {
	if len(types) == 0 {
		return
	}
	fmt.Fprintf(w, `	I.RegisterTypes(
`)
	for _, v := range types {
		fmt.Fprintf(w, `		I.Type("%s", %s),
`, v.name, v.kind)
	}
	fmt.Fprintf(w, "	)\n")
}

func registerVars(w io.Writer, vars []exportedVar) {
	if len(vars) == 0 {
		return
	}
	fmt.Fprintf(w, `	I.RegisterVars(
`)
	for _, v := range vars {
		fmt.Fprintf(w, `		I.Var("%s", &%s),
`, v.name, v.addr)
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
	opkgs := make([]string, 0, len(p.importPkgs))
	for pkg, pkgPath := range p.importPkgs {
		if strings.ContainsAny(pkgPath, ".-_") {
			opkgs = append(opkgs, pkg)
		} else {
			pkgs = append(pkgs, pkg)
		}
	}
	sort.Strings(pkgs)
	sort.Strings(opkgs)
	pkg, pkgPath := p.pkg.Name(), p.pkg.Path()
	fmt.Fprintf(p.w, gopkgExportHeader, pkg, pkgPath, pkgPath, pkg)
	for _, pkg := range pkgs {
		pkgPath := p.importPkgs[pkg]
		fmt.Fprintf(p.w, `	%s "%s"
`, pkg, pkgPath)
	}
	fmt.Fprintf(p.w, "\n")
	for _, pkg := range opkgs {
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
	registerVars(p.w, p.exportedVars)
	registerTypes(p.w, p.exportTypes)
	registerConsts(p.w, p.exportConsts)
	fmt.Fprintf(p.w, gopkgInitExportFooter)
	return nil
}

// -----------------------------------------------------------------------------
