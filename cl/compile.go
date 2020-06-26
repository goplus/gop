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

// Package cl compiles Go+ syntax trees (ast) into a backend code.
// For now the supported backends are `bytecode` and `golang`.
package cl

import (
	"errors"
	"fmt"
	"path"
	"reflect"
	"syscall"

	"github.com/qiniu/goplus/ast"
	"github.com/qiniu/goplus/ast/astutil"
	"github.com/qiniu/goplus/exec.spec"
	"github.com/qiniu/goplus/token"
	"github.com/qiniu/x/log"
)

var (
	// ErrNotFound error
	ErrNotFound = syscall.ENOENT

	// ErrNotAMainPackage error
	ErrNotAMainPackage = errors.New("not a main package")

	// ErrMainFuncNotFound error
	ErrMainFuncNotFound = errors.New("main function not found")

	// ErrSymbolNotVariable error
	ErrSymbolNotVariable = errors.New("symbol exists but not a variable")

	// ErrSymbolNotFunc error
	ErrSymbolNotFunc = errors.New("symbol exists but not a func")

	// ErrSymbolNotType error
	ErrSymbolNotType = errors.New("symbol exists but not a type")
)

var (
	// CallBuiltinOp calls BuiltinOp
	CallBuiltinOp func(kind exec.Kind, op exec.Operator, data ...interface{}) interface{}
)

// -----------------------------------------------------------------------------

// CompileError represents a compiling time error.
type CompileError struct {
	At  ast.Node
	Err error
}

func (p *CompileError) Error() string {
	return p.Err.Error()
}

func newError(at ast.Node, format string, params ...interface{}) *CompileError {
	err := fmt.Errorf(format, params...)
	return &CompileError{at, err}
}

func logError(ctx *blockCtx, at ast.Node, format string, params ...interface{}) {
	err := newError(at, format, params...)
	log.Error(err)
}

func logPanic(ctx *blockCtx, at ast.Node, format string, params ...interface{}) {
	err := newError(at, format, params...)
	log.Panicln(err)
}

func logNonIntegerIdxPanic(ctx *blockCtx, v ast.Node, kind reflect.Kind) {
	logPanic(ctx, v, `non-integer %v index %v`, kind, ctx.code(v))
}

func logIllTypeMapIndexPanic(ctx *blockCtx, v ast.Node, t, typIdx reflect.Type) {
	logPanic(ctx, v, `cannot use %v (type %v) as type %v in map index`, ctx.code(v), t, typIdx)
}

// -----------------------------------------------------------------------------

type pkgCtx struct {
	exec.Package
	infer   exec.Stack
	builtin exec.GoPackage
	out     exec.Builder
	usedfns []*funcDecl
	pkg     *ast.Package
	fset    *token.FileSet
}

func newPkgCtx(out exec.Builder, pkg *ast.Package, fset *token.FileSet) *pkgCtx {
	pkgOut := out.GetPackage()
	builtin := pkgOut.FindGoPackage("")
	p := &pkgCtx{Package: pkgOut, builtin: builtin, out: out, pkg: pkg, fset: fset}
	p.infer.Init()
	return p
}

func (p *pkgCtx) code(v ast.Node) string {
	_, code := p.getCodeInfo(v)
	return code
}

func (p *pkgCtx) getCodeInfo(v ast.Node) (token.Position, string) {
	start, end := v.Pos(), v.End()
	pos := p.fset.Position(start)
	if f, ok := p.pkg.Files[pos.Filename]; ok {
		return pos, string(f.Code[pos.Offset : pos.Offset+int(end-start)])
	}
	log.Panicln("pkgCtx.getCodeInfo failed: file not found -", pos.Filename)
	return pos, ""
}

func (p *pkgCtx) use(f *funcDecl) {
	if f.used {
		return
	}
	p.usedfns = append(p.usedfns, f)
	f.used = true
}

func (p *pkgCtx) resolveFuncs() {
	for {
		n := len(p.usedfns)
		if n == 0 {
			break
		}
		f := p.usedfns[n-1]
		p.usedfns = p.usedfns[:n-1]
		f.Compile()
	}
}

type fileCtx struct {
	*blockCtx // it's global blockCtx
	imports   map[string]string
}

func newFileCtx(block *blockCtx) *fileCtx {
	return &fileCtx{blockCtx: block, imports: make(map[string]string)}
}

// -----------------------------------------------------------------------------

// - varName => *exec.Var
// - stkVarName => *stackVar
// - pkgName => pkgPath
// - funcName => *funcDecl
// - typeName => *typeDecl
//
type iSymbol = interface{}

type iVar interface {
	inCurrentCtx(ctx *blockCtx) bool
	getType() reflect.Type
}

type execVar struct {
	v exec.Var
}

func (p *execVar) inCurrentCtx(ctx *blockCtx) bool {
	return ctx.out.InCurrentCtx(p.v)
}

func (p *execVar) getType() reflect.Type {
	return p.v.Type()
}

type stackVar struct {
	typ   reflect.Type
	index int32
}

func (p *stackVar) inCurrentCtx(ctx *blockCtx) bool {
	return true
}

func (p *stackVar) getType() reflect.Type {
	return p.typ
}

// -----------------------------------------------------------------------------

type blockCtx struct {
	*pkgCtx
	file      *fileCtx
	parent    *blockCtx
	fun       exec.FuncInfo
	syms      map[string]iSymbol
	noExecCtx bool
	checkFlag bool
}

// function block ctx
func newExecBlockCtx(parent *blockCtx) *blockCtx {
	return &blockCtx{
		pkgCtx:    parent.pkgCtx,
		file:      parent.file,
		parent:    parent,
		syms:      make(map[string]iSymbol),
		noExecCtx: false,
	}
}

// normal block ctx, eg. if/switch/for/etc.
func newNormBlockCtx(parent *blockCtx) *blockCtx {
	return newNormBlockCtxEx(parent, true)
}

func newNormBlockCtxEx(parent *blockCtx, noExecCtx bool) *blockCtx {
	return &blockCtx{
		pkgCtx:    parent.pkgCtx,
		file:      parent.file,
		parent:    parent,
		fun:       parent.fun,
		syms:      make(map[string]iSymbol),
		noExecCtx: noExecCtx,
	}
}

// global block ctx
func newGblBlockCtx(pkg *pkgCtx) *blockCtx {
	return &blockCtx{
		pkgCtx:    pkg,
		parent:    nil,
		syms:      make(map[string]iSymbol),
		noExecCtx: true,
	}
}

func (p *blockCtx) getNestDepth() (nestDepth uint32) {
	for {
		if !p.noExecCtx {
			nestDepth++
		}
		if p = p.parent; p == nil {
			return
		}
	}
}

func (p *blockCtx) exists(name string) (ok bool) {
	if _, ok = p.syms[name]; ok {
		return
	}
	if p.parent == nil { // it's global blockCtx
		_, ok = p.file.imports[name]
	}
	return
}

func (p *blockCtx) find(name string) (sym interface{}, ok bool) {
	ctx := p
	for ; p != nil; p = p.parent {
		if sym, ok = p.syms[name]; ok {
			return
		}
	}
	sym, ok = ctx.file.imports[name]
	return
}

func (p *blockCtx) findType(name string) (decl *typeDecl, err error) {
	v, ok := p.find(name)
	if !ok {
		return nil, ErrNotFound
	}
	if decl, ok = v.(*typeDecl); ok {
		return
	}
	return nil, ErrSymbolNotType
}

func (p *blockCtx) findFunc(name string) (addr *funcDecl, err error) {
	v, ok := p.find(name)
	if !ok {
		return nil, ErrNotFound
	}
	if addr, ok = v.(*funcDecl); ok {
		return
	}
	return nil, ErrSymbolNotFunc
}

func (p *blockCtx) findVar(name string) (addr iVar, err error) {
	v, ok := p.find(name)
	if !ok {
		return nil, ErrNotFound
	}
	if addr, ok = v.(iVar); ok {
		return
	}
	return nil, ErrSymbolNotVariable
}

func (p *blockCtx) insertFuncVars(in []reflect.Type, args []string, rets []exec.Var) {
	n := len(args)
	if n > 0 {
		for i := n - 1; i >= 0; i-- {
			name := args[i]
			if name == "" { // unnamed argument
				continue
			}
			if p.exists(name) {
				log.Panicln("insertStkVars failed: symbol exists -", name)
			}
			p.syms[name] = &stackVar{index: int32(i - n), typ: in[i]}
		}
	}
	for _, ret := range rets {
		if ret.IsUnnamedOut() {
			continue
		}
		p.syms[ret.Name()] = &execVar{ret}
	}
}

func (p *blockCtx) insertVar(name string, typ reflect.Type, inferOnly ...bool) *execVar {
	if p.exists(name) {
		log.Panicln("insertVar failed: symbol exists -", name)
	}
	v := p.NewVar(typ, name)
	if inferOnly == nil {
		p.out.DefineVar(v)
	}
	ev := &execVar{v}
	p.syms[name] = ev
	return ev
}

func (p *blockCtx) insertFunc(name string, fun *funcDecl) {
	if p.exists(name) {
		log.Panicln("insertFunc failed: symbol exists -", name)
	}
	p.syms[name] = fun
}

func (p *blockCtx) insertMethod(typeName, methodName string, method *methodDecl) {
	if p.parent != nil {
		log.Panicln("insertMethod failed: unexpected - non global method declaration?")
	}
	typ, err := p.findType(typeName)
	if err == ErrNotFound {
		typ = new(typeDecl)
		p.syms[typeName] = typ
	} else if err != nil {
		log.Panicln("insertMethod failed:", err)
	} else if typ.Alias {
		log.Panicln("insertMethod failed: alias?")
	}
	if typ.Methods == nil {
		typ.Methods = map[string]*methodDecl{methodName: method}
	} else {
		if _, ok := typ.Methods[methodName]; ok {
			log.Panicln("insertMethod failed: method exists -", typeName, methodName)
		}
		typ.Methods[methodName] = method
	}
}

// -----------------------------------------------------------------------------

// A Package represents a Go+ package.
type Package struct {
	syms map[string]iSymbol
}

// PkgAct represents a package compiling action.
type PkgAct int

const (
	// PkgActNone - do nothing
	PkgActNone PkgAct = iota
	// PkgActClMain - compile main function
	PkgActClMain
	// PkgActClAll - compile all things
	PkgActClAll
)

// NewPackage creates a Go+ package instance.
func NewPackage(out exec.Builder, pkg *ast.Package, fset *token.FileSet, act PkgAct) (p *Package, err error) {
	if pkg == nil {
		return nil, ErrNotFound
	}
	if CallBuiltinOp == nil {
		log.Panicln("NewPackage failed: variable CallBuiltinOp is uninitialized")
	}
	p = &Package{}
	ctxPkg := newPkgCtx(out, pkg, fset)
	ctx := newGblBlockCtx(ctxPkg)
	for _, f := range pkg.Files {
		loadFile(ctx, f)
	}
	switch act {
	case PkgActClAll:
		for _, sym := range ctx.syms {
			if f, ok := sym.(*funcDecl); ok && f.fi != nil {
				ctxPkg.use(f)
			}
		}
		if pkg.Name != "main" {
			break
		}
		fallthrough
	case PkgActClMain:
		if pkg.Name != "main" {
			return nil, ErrNotAMainPackage
		}
		entry, err := ctx.findFunc("main")
		if err != nil {
			if err == ErrNotFound {
				err = ErrMainFuncNotFound
			}
			return p, err
		}
		if entry.ctx.noExecCtx {
			ctx.file = entry.ctx.file
			compileBlockStmtWithout(ctx, entry.body)
		} else {
			out.CallFunc(entry.Get(), 0)
			ctxPkg.use(entry)
		}
		out.Return(-1)
	}
	ctxPkg.resolveFuncs()
	p.syms = ctx.syms
	return
}

// SymKind represents a symbol kind.
type SymKind uint

const (
	// SymInvalid - invalid symbol kind
	SymInvalid SymKind = iota
	// SymVar - symbol is a variable
	SymVar
	// SymFunc - symbol is a function
	SymFunc
	// SymType - symbol is a type
	SymType
)

// Find lookups a symbol and returns it's kind and the object instance.
func (p *Package) Find(name string) (kind SymKind, v interface{}, ok bool) {
	if v, ok = p.syms[name]; !ok {
		return
	}
	switch v.(type) {
	case *exec.Var:
		kind = SymVar
	case *funcDecl:
		kind = SymFunc
	case *typeDecl:
		kind = SymType
	default:
		log.Panicln("Package.Find: unknown symbol type -", reflect.TypeOf(v))
	}
	return
}

func loadFile(ctx *blockCtx, f *ast.File) {
	file := newFileCtx(ctx)
	last := len(f.Decls) - 1
	ctx.file = file
	for i, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			loadFunc(ctx, d, f.HasUnnamed && i == last)
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				loadImports(file, d)
			case token.TYPE:
				loadTypes(ctx, d)
			case token.CONST:
				loadConsts(ctx, d)
			case token.VAR:
				loadVars(ctx, d)
			default:
				log.Panicln("tok:", d.Tok, "spec:", reflect.TypeOf(d.Specs).Elem())
			}
		default:
			log.Panicln("gopkg.Package.load: unknown decl -", reflect.TypeOf(decl))
		}
	}
}

func loadImports(ctx *fileCtx, d *ast.GenDecl) {
	for _, item := range d.Specs {
		loadImport(ctx, item.(*ast.ImportSpec))
	}
}

func loadImport(ctx *fileCtx, spec *ast.ImportSpec) {
	var pkgPath = astutil.ToString(spec.Path)
	var name string
	if spec.Name != nil {
		name = spec.Name.Name
		switch name {
		case "_", ".":
			panic("not impl")
		}
	} else {
		name = path.Base(pkgPath)
	}
	ctx.imports[name] = pkgPath
}

func loadTypes(ctx *blockCtx, d *ast.GenDecl) {
	for _, item := range d.Specs {
		loadType(ctx, item.(*ast.TypeSpec))
	}
}

func loadType(ctx *blockCtx, spec *ast.TypeSpec) {
}

func loadConsts(ctx *blockCtx, d *ast.GenDecl) {
}

func loadVars(ctx *blockCtx, d *ast.GenDecl) {
	for _, item := range d.Specs {
		loadVar(ctx, item.(*ast.ValueSpec))
	}
}

func loadVar(ctx *blockCtx, spec *ast.ValueSpec) {
}

func loadFunc(ctx *blockCtx, d *ast.FuncDecl, isUnnamed bool) {
	var name = d.Name.Name
	if d.Recv != nil {
		recv := astutil.ToRecv(d.Recv)
		ctx.insertMethod(recv.Type, name, &methodDecl{
			recv:    recv.Name,
			pointer: recv.Pointer,
			typ:     d.Type,
			body:    d.Body,
			file:    ctx.file,
		})
	} else if name == "init" {
		log.Panicln("loadFunc TODO: init")
	} else {
		funCtx := newExecBlockCtx(ctx)
		funCtx.noExecCtx = isUnnamed
		ctx.insertFunc(name, newFuncDecl(name, d.Type, d.Body, funCtx))
	}
}

// -----------------------------------------------------------------------------
