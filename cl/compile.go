/*
 Copyright 2021 The GoPlus Authors (goplus.org)

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

// Package cl compiles Go+ syntax trees (ast).
package cl

import (
	"context"
	"go/types"
	"log"
	"path"
	"reflect"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

// -----------------------------------------------------------------------------

// Config of loading Go+ packages.
type Config struct {
	// Context specifies the context for the load operation.
	// If the context is cancelled, the loader may stop early
	// and return an ErrCancelled error.
	// If Context is nil, the load cannot be cancelled.
	Context context.Context

	// Logf is the logger for the config.
	// If the user provides a logger, debug logging is enabled.
	// If the GOPACKAGESDEBUG environment variable is set to true,
	// but the logger is nil, default to log.Printf.
	Logf func(format string, args ...interface{})

	// Dir is the directory in which to run the build system's query tool
	// that provides information about the packages.
	// If Dir is empty, the tool is run in the current directory.
	Dir string

	// Env is the environment to use when invoking the build system's query tool.
	// If Env is nil, the current environment is used.
	// As in os/exec's Cmd, only the last value in the slice for
	// each environment key is used. To specify the setting of only
	// a few variables, append to the current environment, as in:
	//
	//	opt.Env = append(os.Environ(), "GOOS=plan9", "GOARCH=386")
	//
	Env []string

	// BuildFlags is a list of command-line flags to be passed through to
	// the build system's query tool.
	BuildFlags []string

	// LoadPkgs is called to load all import packages.
	LoadPkgs gox.LoadPkgsFunc
}

// NewPackage creates a Go+ package instance.
func NewPackage(
	pkgPath string, pkg *ast.Package, fset *token.FileSet, conf *Config) (p *gox.Package, err error) {
	if conf == nil {
		conf = &Config{}
	}
	confGox := &gox.Config{
		Context:    conf.Context,
		Logf:       conf.Logf,
		Dir:        conf.Dir,
		Env:        conf.Env,
		BuildFlags: conf.BuildFlags,
		Fset:       fset,
		ParseFile:  nil, // TODO
		LoadPkgs:   conf.LoadPkgs,
		Prefix:     nil,
		Contracts:  nil,
		NewBuiltin: newBuiltinDefault,
	}
	p = gox.NewPackage(pkgPath, pkg.Name, confGox)
	for _, f := range pkg.Files {
		loadFile(p, f)
	}
	return
}

type blockCtx struct {
	pkg     *gox.Package
	fns     []*clFuncBody
	cb      *gox.CodeBuilder
	imports map[string]*gox.PkgRef
}

func (p *blockCtx) complete() {
	for _, fn := range p.fns {
		fn.load(p)
	}
	p.fns = nil
}

type clFuncBody struct {
	fn   *gox.Func
	body *ast.BlockStmt
}

func (p *clFuncBody) load(ctx *blockCtx) {
	loadFuncBody(ctx, p.fn, p.body)
}

func loadFuncBody(ctx *blockCtx, fn *gox.Func, body *ast.BlockStmt) {
	cb := fn.BodyStart(ctx.pkg)
	compileStmts(ctx, body.List)
	cb.End()
}

func loadFile(p *gox.Package, f *ast.File) {
	ctx := &blockCtx{pkg: p, cb: p.CB(), imports: make(map[string]*gox.PkgRef)}
	last := len(f.Decls) - 1
	for i, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			loadFunc(ctx, d, f.NoEntrypoint && i == last)
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				loadImports(ctx, d)
			case token.TYPE:
				loadTypes(ctx, d)
			case token.CONST:
				loadConsts(ctx, d)
			case token.VAR:
				loadVars(ctx, d)
			default:
				log.Panicln("TODO - tok:", d.Tok, "spec:", reflect.TypeOf(d.Specs).Elem())
			}
		default:
			log.Panicln("TODO - gopkg.Package.load: unknown decl -", reflect.TypeOf(decl))
		}
	}
	ctx.complete()
}

func loadFunc(ctx *blockCtx, d *ast.FuncDecl, isUnnamed bool) {
	pkg := ctx.pkg
	name := d.Name.Name
	if d.Recv != nil {
		log.Panicln("loadFunc TODO: method")
	} else if name == "init" {
		log.Panicln("loadFunc TODO: init")
	} else {
		sig := toFuncType(ctx, d.Type)
		fn := pkg.NewFuncWith(name, sig)
		if body := d.Body; body != nil {
			ctx.fns = append(ctx.fns, &clFuncBody{fn: fn, body: body})
		}
	}
}

func insertParams(scope *types.Scope, params *types.Tuple) {
	if params != nil {
		for i, n := 0, params.Len(); i < n; i++ {
			v := params.At(i)
			if name := v.Name(); name != "" {
				scope.Insert(v)
			}
		}
	}
}

func loadImports(ctx *blockCtx, d *ast.GenDecl) {
	for _, item := range d.Specs {
		loadImport(ctx, item.(*ast.ImportSpec))
	}
}

func loadImport(ctx *blockCtx, spec *ast.ImportSpec) {
	pkgPath := toString(spec.Path)
	pkg := ctx.pkg.Import(pkgPath)
	var name string
	if spec.Name != nil {
		name = spec.Name.Name
	} else {
		name = path.Base(pkgPath) // TODO: pkg.Name()
	}
	ctx.imports[name] = pkg
}

func loadTypes(ctx *blockCtx, d *ast.GenDecl) {
	panic("TODO: loadTypes")
}

func loadConsts(ctx *blockCtx, d *ast.GenDecl) {
	panic("TODO: loadConsts")
}

func loadVars(ctx *blockCtx, d *ast.GenDecl) {
	for _, spec := range d.Specs {
		var v = spec.(*ast.ValueSpec)
		var typ types.Type
		names := makeNames(v.Names)
		if v.Type != nil {
			typ = toType(ctx, v.Type)
		}
		pkg, cb := ctx.pkg, ctx.cb
		varDecl := pkg.NewVar(typ, names...)
		if v.Values != nil {
			varDecl.InitStart(pkg)
			for _, val := range v.Values {
				compileExpr(ctx, val)
			}
			cb.EndInit(len(v.Values))
		}
	}
}

func makeNames(vals []*ast.Ident) []string {
	names := make([]string, len(vals))
	for i, v := range vals {
		names[i] = v.Name
	}
	return names
}

// -----------------------------------------------------------------------------

/*
import (
	"errors"
	"path"
	"reflect"
	"syscall"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/ast/astutil"
	"github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/token"
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
func NewPackage(pkg *ast.Package, fset *token.FileSet, act PkgAct) (p *gox.Package, err error) {
	if pkg == nil {
		return nil, ErrNotFound
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
			ctx.checkLabels()
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
			loadFunc(ctx, d, f.NoEntrypoint && i == last)
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				loadImports(file, d)
			case token.TYPE:
				loadTypes(ctx, d)
			case token.CONST:
				loadConsts(ctx, d)
			case token.VAR:
				compileStmt(ctx, &ast.DeclStmt{decl})
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
	if ctx.exists(spec.Name.Name) {
		log.Panicln("loadType failed: symbol exists -", spec.Name.Name)
	}
	t := toType(ctx, spec.Type).(reflect.Type)

	ctx.out.DefineType(t, spec.Name.Name)

	tDecl := &typeDecl{
		Type: t,
	}
	ctx.syms[spec.Name.Name] = tDecl
	ctx.types[t] = tDecl
}

func loadConsts(ctx *blockCtx, d *ast.GenDecl) {
}

func loadVars(ctx *blockCtx, d *ast.GenDecl, stmt ast.Stmt) {
	for _, item := range d.Specs {
		spec := item.(*ast.ValueSpec)
		for i := 0; i < len(spec.Names); i++ {
			start := ctx.out.StartStmt(stmt)
			name := spec.Names[i].Name
			if len(spec.Values) > i {
				loadVar(ctx, name, spec.Type, spec.Values[i])
			} else {
				loadVar(ctx, name, spec.Type, nil)
			}
			ctx.out.EndStmt(stmt, start)
		}
	}
}

func loadVar(ctx *blockCtx, name string, typ ast.Expr, value ast.Expr) {
	var t reflect.Type
	if typ != nil {
		t = toType(ctx, typ).(reflect.Type)
	}
	if value != nil {
		expr := compileExpr(ctx, value)
		in := ctx.infer.Get(-1)
		if t == nil {
			t = boundType(in.(iValue))
		}
		expr()
		addr := ctx.insertVar(name, t)
		checkType(addr.getType(), in, ctx.out)
		ctx.infer.PopN(1)
		ctx.out.StoreVar(addr.v)
	} else {
		ctx.insertVar(name, t)
	}
}

func loadFunc(ctx *blockCtx, d *ast.FuncDecl, isUnnamed bool) {
	var name = d.Name.Name
	if d.Recv != nil {
		recv := astutil.ToRecv(d.Recv)
		funCtx := newExecBlockCtx(ctx)
		funCtx.noExecCtx = isUnnamed
		funCtx.funcCtx = newFuncCtx(nil)
		ctx.insertMethod(recv, name, d, funCtx)
	} else if name == "init" {
		log.Panicln("loadFunc TODO: init")
	} else {
		funCtx := newExecBlockCtx(ctx)
		funCtx.noExecCtx = isUnnamed
		funCtx.funcCtx = newFuncCtx(nil)
		ctx.insertFunc(name, newFuncDecl(name, nil, d.Type, d.Body, funCtx))
	}
}
*/
// -----------------------------------------------------------------------------
