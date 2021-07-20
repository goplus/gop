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

const (
	gopPrefix = "Gop_"
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

func getUnderlying(ctx *blockCtx, typ types.Type) types.Type {
	if t := typ.Underlying(); t != nil {
		return t
	}
	if t, ok := typ.(*types.Named); ok {
		return doLoadUnderlying(ctx.pkgCtx, t)
	}
	panic("TODO: getUnderlying: not named type - " + typ.String())
}

func doLoadUnderlying(ctx *pkgCtx, typ *types.Named) types.Type {
	panic("TODO: doLoadUnderlying")
	/*
		ld := ctx.syms[typ.Obj().Name()].(*typeLoader)
		doInitType(ld)
		return typ.Underlying()
	*/
}

// NewPackage creates a Go+ package instance.
func NewPackage(
	pkgPath string, pkg *ast.Package, fset *token.FileSet, conf *Config) (p *gox.Package, err error) {
	if conf == nil {
		conf = &Config{}
	}
	ctx := &pkgCtx{syms: make(map[string]loader)}
	loadUnderlying := func(at *gox.Package, typ *types.Named) types.Type {
		return doLoadUnderlying(ctx, typ)
	}
	loadPkgs := conf.LoadPkgs
	if loadPkgs == nil {
		loadPkgs = LoadGopPkgs
	}
	confGox := &gox.Config{
		Context:        conf.Context,
		Logf:           conf.Logf,
		Dir:            conf.Dir,
		Env:            conf.Env,
		BuildFlags:     conf.BuildFlags,
		Fset:           fset,
		ParseFile:      nil, // TODO
		LoadPkgs:       loadPkgs,
		LoadUnderlying: loadUnderlying,
		Prefix:         gopPrefix,
		Contracts:      nil,
		NewBuiltin:     newBuiltinDefault,
	}
	p = gox.NewPackage(pkgPath, pkg.Name, confGox)
	for _, f := range pkg.Files {
		loadFile(p, ctx, f)
	}
	for _, f := range pkg.Files {
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				ctx.loadSymbol(d.Name.Name)
			case *ast.GenDecl:
				switch d.Tok {
				case token.TYPE:
					for _, spec := range d.Specs {
						ctx.loadType(spec.(*ast.TypeSpec).Name.Name)
					}
				case token.CONST, token.VAR:
					for _, spec := range d.Specs {
						for _, name := range spec.(*ast.ValueSpec).Names {
							ctx.loadSymbol(name.Name)
						}
					}
				}
			}
		}
	}
	ctx.complete()
	return
}

type loader interface {
	load()
}

type loaderFunc func()

func (fn loaderFunc) load() {
	fn()
}

type typeLoader struct {
	typ, typInit func()
	methods      []func()
}

func getTypeLoader(syms map[string]loader, name string) *typeLoader {
	t, ok := syms[name]
	if !ok {
		t = &typeLoader{}
		syms[name] = t
	}
	return t.(*typeLoader)
}

func (p *typeLoader) load() {
	doNewType(p)
	doInitType(p)
	doInitMethods(p)
}

func doNewType(ld *typeLoader) {
	if typ := ld.typ; typ != nil {
		ld.typ = nil
		typ()
	}
}

func doInitType(ld *typeLoader) {
	if typInit := ld.typInit; typInit != nil {
		ld.typInit = nil
		typInit()
	}
}

func doInitMethods(ld *typeLoader) {
	if methods := ld.methods; methods != nil {
		ld.methods = nil
		for _, method := range methods {
			method()
		}
	}
}

type pkgCtx struct {
	syms map[string]loader
}

type blockCtx struct {
	*pkgCtx
	pkg     *gox.Package
	cb      *gox.CodeBuilder
	imports map[string]*gox.PkgRef
}

func (p *pkgCtx) complete() {
}

func (p *pkgCtx) loadType(name string) {
	if sym, ok := p.syms[name]; ok {
		if ld, ok := sym.(*typeLoader); ok {
			ld.load()
		}
	}
}

func (p *pkgCtx) loadSymbol(name string) bool {
	if f, ok := p.syms[name]; ok {
		if ld, ok := f.(*typeLoader); ok {
			doNewType(ld) // create this type, but don't init
			return true
		}
		delete(p.syms, name)
		f.load()
		return true
	}
	return false
}

func loadFile(p *gox.Package, parent *pkgCtx, f *ast.File) {
	syms := parent.syms
	ctx := &blockCtx{
		pkg: p, pkgCtx: parent, cb: p.CB(), imports: make(map[string]*gox.PkgRef),
	}
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil {
				name := d.Name.Name
				if name == "init" {
					log.Panicln("loadFunc TODO: init")
				}
				syms[name] = loaderFunc(func() {
					loadFunc(ctx, nil, d)
				})
			} else {
				name := getRecvTypeName(d.Recv)
				ld := getTypeLoader(syms, name)
				ld.methods = append(ld.methods, func() {
					doInitType(ld)
					recv := toRecv(ctx, d.Recv)
					loadFunc(ctx, recv, d)
				})
			}
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				for _, item := range d.Specs {
					loadImport(ctx, item.(*ast.ImportSpec))
				}
			case token.TYPE:
				for _, spec := range d.Specs {
					t := spec.(*ast.TypeSpec)
					name := t.Name.Name
					log.Println("TYPE:", name) // TODO: remove
					ld := getTypeLoader(syms, name)
					ld.typ = func() {
						if t.Assign != token.NoPos { // alias type
							ctx.pkg.AliasType(name, toType(ctx, t.Type))
							return
						}
						decl := ctx.pkg.NewType(name)
						ld.typInit = func() { // decycle
							decl.InitType(ctx.pkg, toType(ctx, t.Type))
						}
					}
				}
			case token.CONST:
				for _, spec := range d.Specs {
					v := spec.(*ast.ValueSpec)
					names := makeNames(v.Names)
					setNamesLoader(syms, names, func() {
						loadConsts(ctx, names, v)
						removeNames(syms, names)
					})
				}
			case token.VAR:
				for _, spec := range d.Specs {
					v := spec.(*ast.ValueSpec)
					names := makeNames(v.Names)
					setNamesLoader(syms, names, func() {
						loadVars(ctx, names, v)
						removeNames(syms, names)
					})
				}
			default:
				log.Panicln("TODO - tok:", d.Tok, "spec:", reflect.TypeOf(d.Specs).Elem())
			}
		default:
			log.Panicln("TODO - gopkg.Package.load: unknown decl -", reflect.TypeOf(decl))
		}
	}
}

func loadType(ctx *blockCtx, t *ast.TypeSpec) {
	pkg := ctx.pkg
	if t.Assign != token.NoPos { // alias type
		pkg.AliasType(t.Name.Name, toType(ctx, t.Type))
	} else {
		pkg.NewType(t.Name.Name).InitType(pkg, toType(ctx, t.Type))
	}
}

func loadFunc(ctx *blockCtx, recv *types.Var, d *ast.FuncDecl) {
	sig := toFuncType(ctx, d.Type, recv)
	fn := ctx.pkg.NewFuncWith(d.Name.Name, sig)
	log.Println("loadFunc:", fn.Name(), fn.String(), fn.FullName()) // TODO: remove
	if body := d.Body; body != nil {
		loadFuncBody(ctx, fn, body)
	}
}

func loadFuncBody(ctx *blockCtx, fn *gox.Func, body *ast.BlockStmt) {
	cb := fn.BodyStart(ctx.pkg)
	compileStmts(ctx, body.List)
	cb.End()
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

func loadConsts(ctx *blockCtx, names []string, v *ast.ValueSpec) {
	var typ types.Type
	if v.Type != nil {
		typ = toType(ctx, v.Type)
	}
	cb := ctx.pkg.NewConstStart(typ, names...)
	for _, val := range v.Values {
		compileExpr(ctx, val)
	}
	cb.EndInit(len(v.Values))
}

func loadVars(ctx *blockCtx, names []string, v *ast.ValueSpec) {
	var typ types.Type
	if v.Type != nil {
		typ = toType(ctx, v.Type)
	}
	varDecl := ctx.pkg.NewVar(typ, names...)
	if v.Values != nil {
		cb := varDecl.InitStart(ctx.pkg)
		for _, val := range v.Values {
			compileExpr(ctx, val)
		}
		cb.EndInit(len(v.Values))
	}
}

func makeNames(vals []*ast.Ident) []string {
	names := make([]string, len(vals))
	for i, v := range vals {
		names[i] = v.Name
	}
	return names
}

func removeNames(syms map[string]loader, names []string) {
	for _, name := range names {
		delete(syms, name)
	}
}

func setNamesLoader(syms map[string]loader, names []string, load func()) {
	for _, name := range names {
		syms[name] = loaderFunc(load)
	}
}

// -----------------------------------------------------------------------------

/*
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

*/
// -----------------------------------------------------------------------------
