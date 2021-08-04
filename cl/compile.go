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
	"fmt"
	"go/types"
	"log"
	"os"
	"path"
	"reflect"
	"strings"

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

	// WorkingDir is the directory in which to run gop compiler.
	// TargetDir is the directory in which to generate Go files.
	// If SrcDir or TargetDir is empty, it is same as Dir.
	WorkingDir, TargetDir string

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

	// Fset provides source position information for syntax trees and types.
	// If Fset is nil, Load will use a new fileset, but preserve Fset's value.
	Fset *token.FileSet

	// GenGoPkg is called to convert a Go+ package into Go.
	GenGoPkg func(pkgDir string, base *Config) error

	// PkgsLoader is the Go+ packages loader (will be set if it is nil).
	PkgsLoader *PkgsLoader

	// CacheLoadPkgs set to cache all loaded packages or not.
	CacheLoadPkgs bool

	// NoFileLine = true means not generate file line comments.
	NoFileLine bool
}

func (conf *Config) Ensure() *Config {
	if conf == nil {
		conf = &Config{Fset: token.NewFileSet()}
	}
	if conf.PkgsLoader == nil {
		initPkgsLoader(conf)
	}
	return conf
}

type nodeInterp struct {
	fset       *token.FileSet
	files      map[string]*ast.File
	workingDir string
}

func (p *nodeInterp) Position(start token.Pos) token.Position {
	pos := p.fset.Position(start)
	pos.Filename = relFile(p.workingDir, pos.Filename)
	return pos
}

func (p *nodeInterp) Caller(node ast.Node) string {
	if expr, ok := node.(*ast.CallExpr); ok {
		node = expr.Fun
		start := node.Pos()
		pos := p.fset.Position(start)
		f := p.files[pos.Filename]
		n := int(node.End() - start)
		return string(f.Code[pos.Offset : pos.Offset+n])
	}
	return "the function call"
}

func (p *nodeInterp) LoadExpr(node ast.Node) (src string, pos token.Position) {
	start := node.Pos()
	pos = p.fset.Position(start)
	f := p.files[pos.Filename]
	n := int(node.End() - start)
	pos.Filename = relFile(p.workingDir, pos.Filename)
	src = string(f.Code[pos.Offset : pos.Offset+n])
	return
}

// NewPackage creates a Go+ package instance.
func NewPackage(pkgPath string, pkg *ast.Package, conf *Config) (p *gox.Package, err error) {
	conf = conf.Ensure()
	dir := conf.Dir
	if dir == "" {
		dir, _ = os.Getwd()
	}
	workingDir := conf.WorkingDir
	if workingDir == "" {
		workingDir, _ = os.Getwd()
	}
	targetDir := conf.TargetDir
	if targetDir == "" {
		targetDir = dir
	}
	interp := &nodeInterp{fset: conf.Fset, files: pkg.Files, workingDir: workingDir}
	ctx := &pkgCtx{syms: make(map[string]loader), nodeInterp: interp}
	confGox := &gox.Config{
		Context:         conf.Context,
		Logf:            conf.Logf,
		Dir:             dir,
		Env:             conf.Env,
		BuildFlags:      conf.BuildFlags,
		Fset:            conf.Fset,
		LoadPkgs:        conf.PkgsLoader.LoadPkgs,
		HandleErr:       ctx.handleErr,
		NodeInterpreter: interp,
		Prefix:          gopPrefix,
		ParseFile:       nil, // TODO
		NewBuiltin:      newBuiltinDefault,
	}
	fileLine := !conf.NoFileLine
	p = gox.NewPackage(pkgPath, pkg.Name, confGox)
	for _, f := range pkg.Files {
		loadFile(p, ctx, f, targetDir, fileLine)
	}
	for _, f := range pkg.Files {
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Recv == nil {
					name := d.Name.Name
					if name != "init" {
						ctx.loadSymbol(name)
					}
				} else {
					if name, ok := getRecvTypeName(ctx, d.Recv, false); ok {
						getTypeLoader(ctx.syms, token.NoPos, name).load()
					}
				}
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
	for _, load := range ctx.inits {
		load()
	}
	err = ctx.complete()
	return
}

type loader interface {
	load()
	pos() token.Pos
}

type baseLoader struct {
	fn    func()
	start token.Pos
}

func initLoader(ctx *pkgCtx, syms map[string]loader, start token.Pos, name string, fn func()) {
	if old, ok := syms[name]; ok && start != token.NoPos {
		pos := ctx.Position(start)
		oldpos := ctx.Position(old.pos())
		ctx.handleCodeErrorf(
			&pos, "%s redeclared in this block\n\tprevious declaration at %v", name, oldpos)
	}
	syms[name] = &baseLoader{start: start, fn: fn}
}

func (p *baseLoader) load() {
	p.fn()
}

func (p *baseLoader) pos() token.Pos {
	return p.start
}

type typeLoader struct {
	typ, typInit func()
	methods      []func()
	start        token.Pos
}

func getTypeLoader(syms map[string]loader, start token.Pos, name string) *typeLoader {
	t, ok := syms[name]
	if ok {
		if start != token.NoPos {
			panic("TODO: redefine")
		}
	} else {
		t = &typeLoader{start: start}
		syms[name] = t
	}
	return t.(*typeLoader)
}

func (p *typeLoader) pos() token.Pos {
	return p.start
}

func (p *typeLoader) load() {
	doNewType(p)
	doInitType(p)
	doInitMethods(p)
}

func doNewType(ld *typeLoader) bool {
	if typ := ld.typ; typ != nil {
		ld.typ = nil
		typ()
		return true
	}
	return false
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

type Errors struct {
	Errs []error
}

func (p *Errors) Error() string {
	msgs := make([]string, len(p.Errs))
	for i, err := range p.Errs {
		msgs[i] = err.Error()
	}
	return strings.Join(msgs, "\n")
}

type pkgCtx struct {
	*nodeInterp
	syms  map[string]loader
	inits []func()
	errs  []error
}

type blockCtx struct {
	*pkgCtx
	pkg       *gox.Package
	cb        *gox.CodeBuilder
	fset      *token.FileSet
	targetDir string
	fileLine  bool
	imports   map[string]*gox.PkgRef
}

func newCodeErrorf(pos *token.Position, format string, args ...interface{}) *gox.CodeError {
	return &gox.CodeError{Pos: pos, Msg: fmt.Sprintf(format, args...)}
}

func (p *pkgCtx) newCodeErrorf(start token.Pos, format string, args ...interface{}) error {
	pos := p.Position(start)
	return newCodeErrorf(&pos, format, args...)
}

func (p *pkgCtx) handleCodeErrorf(pos *token.Position, format string, args ...interface{}) {
	p.handleErr(newCodeErrorf(pos, format, args...))
}

func (p *pkgCtx) handleErr(err error) {
	p.errs = append(p.errs, err)
}

func (p *pkgCtx) complete() error {
	if p.errs != nil {
		return &Errors{Errs: p.errs}
	}
	return nil
}

func (p *pkgCtx) loadType(name string) {
	if sym, ok := p.syms[name]; ok {
		if ld, ok := sym.(*typeLoader); ok {
			ld.load()
		}
	}
}

func (p *pkgCtx) loadSymbol(name string) bool {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				p.handleErr(err)
			} else {
				panic(e)
			}
		}
	}()

	if f, ok := p.syms[name]; ok {
		if ld, ok := f.(*typeLoader); ok {
			return doNewType(ld) // create this type, but don't init
		}
		delete(p.syms, name)
		f.load()
		return true
	}
	return false
}

func loadFile(p *gox.Package, parent *pkgCtx, f *ast.File, targetDir string, fileLine bool) {
	syms := parent.syms
	ctx := &blockCtx{
		pkg: p, pkgCtx: parent, cb: p.CB(), fset: p.Fset, targetDir: targetDir, fileLine: fileLine,
		imports: make(map[string]*gox.PkgRef),
	}
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil {
				name := d.Name
				fn := func() {
					loadFunc(ctx, nil, d)
				}
				if name.Name == "init" {
					parent.inits = append(parent.inits, fn)
				} else {
					initLoader(parent, syms, name.Pos(), name.Name, fn)
				}
			} else {
				if name, ok := getRecvTypeName(parent, d.Recv, true); ok {
					ld := getTypeLoader(syms, token.NoPos, name)
					ld.methods = append(ld.methods, func() {
						doInitType(ld)
						recv := toRecv(ctx, d.Recv)
						loadFunc(ctx, recv, d)
					})
				}
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
					ld := getTypeLoader(syms, t.Name.Pos(), name)
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
					vSpec := spec.(*ast.ValueSpec)
					setNamesLoader(parent, syms, vSpec.Names, func() {
						if v := vSpec; v != nil { // only init once
							vSpec = nil
							names := makeNames(v.Names)
							loadConsts(ctx, names, v)
							removeNames(syms, names)
						}
					})
				}
			case token.VAR:
				for _, spec := range d.Specs {
					vSpec := spec.(*ast.ValueSpec)
					setNamesLoader(parent, syms, vSpec.Names, func() {
						if v := vSpec; v != nil { // only init once
							vSpec = nil
							names := makeNames(v.Names)
							loadVars(ctx, names, v)
							removeNames(syms, names)
						}
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
	fn, err := ctx.pkg.NewFuncWith(d.Pos(), d.Name.Name, sig, func() token.Pos {
		return d.Recv.List[0].Type.Pos()
	})
	if err != nil {
		ctx.handleErr(err)
		return
	}
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
	cb := ctx.pkg.NewConstStart(v.Names[0].Pos(), typ, names...)
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
	varDecl := ctx.pkg.NewVar(v.Names[0].Pos(), typ, names...)
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

func setNamesLoader(ctx *pkgCtx, syms map[string]loader, names []*ast.Ident, load func()) {
	for _, name := range names {
		initLoader(ctx, syms, name.Pos(), name.Name, load)
	}
}

// -----------------------------------------------------------------------------
