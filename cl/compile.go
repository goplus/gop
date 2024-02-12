/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package cl compiles Go+ syntax trees (ast).
package cl

import (
	"fmt"
	"go/types"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/ast/fromgo"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
	"github.com/goplus/gox/cpackages"
	"github.com/goplus/mod/modfile"
	"github.com/qiniu/x/errors"
)

type dbgFlags int

const (
	DbgFlagLoad dbgFlags = 1 << iota
	DbgFlagLookup
	FlagNoMarkAutogen
	DbgFlagAll = DbgFlagLoad | DbgFlagLookup
)

var (
	enableRecover = true
)

var (
	debugLoad     bool
	debugLookup   bool
	noMarkAutogen bool // add const _ = true
)

func SetDisableRecover(disableRecover bool) {
	enableRecover = !disableRecover
}

func SetDebug(flags dbgFlags) {
	debugLoad = (flags & DbgFlagLoad) != 0
	debugLookup = (flags & DbgFlagLookup) != 0
	noMarkAutogen = (flags & FlagNoMarkAutogen) != 0
}

// -----------------------------------------------------------------------------

// Recorder represents a compiling event recorder.
type Recorder interface {
	// Type maps expressions to their types, and for constant
	// expressions, also their values. Invalid expressions are
	// omitted.
	//
	// For (possibly parenthesized) identifiers denoting built-in
	// functions, the recorded signatures are call-site specific:
	// if the call result is not a constant, the recorded type is
	// an argument-specific signature. Otherwise, the recorded type
	// is invalid.
	//
	// The Types map does not record the type of every identifier,
	// only those that appear where an arbitrary expression is
	// permitted. For instance, the identifier f in a selector
	// expression x.f is found only in the Selections map, the
	// identifier z in a variable declaration 'var z int' is found
	// only in the Defs map, and identifiers denoting packages in
	// qualified identifiers are collected in the Uses map.
	Type(ast.Expr, types.TypeAndValue)

	// Instantiate maps identifiers denoting generic types or functions to their
	// type arguments and instantiated type.
	//
	// For example, Instantiate will map the identifier for 'T' in the type
	// instantiation T[int, string] to the type arguments [int, string] and
	// resulting instantiated *Named type. Given a generic function
	// func F[A any](A), Instances will map the identifier for 'F' in the call
	// expression F(int(1)) to the inferred type arguments [int], and resulting
	// instantiated *Signature.
	//
	// Invariant: Instantiating Uses[id].Type() with Instances[id].TypeArgs
	// results in an equivalent of Instances[id].Type.
	Instantiate(*ast.Ident, types.Instance)

	// Def maps identifiers to the objects they define (including
	// package names, dots "." of dot-imports, and blank "_" identifiers).
	// For identifiers that do not denote objects (e.g., the package name
	// in package clauses, or symbolic variables t in t := x.(type) of
	// type switch headers), the corresponding objects are nil.
	//
	// For an embedded field, Def maps the field *Var it defines.
	//
	// Invariant: Defs[id] == nil || Defs[id].Pos() == id.Pos()
	Def(id *ast.Ident, obj types.Object)

	// Use maps identifiers to the objects they denote.
	//
	// For an embedded field, Use maps the *TypeName it denotes.
	//
	// Invariant: Uses[id].Pos() != id.Pos()
	Use(id *ast.Ident, obj types.Object)

	// Implicit maps nodes to their implicitly declared objects, if any.
	// The following node and object types may appear:
	//
	//     node               declared object
	//
	//     *ast.ImportSpec    *PkgName for imports without renames
	//     *ast.CaseClause    type-specific *Var for each type switch case clause (incl. default)
	//     *ast.Field         anonymous parameter *Var (incl. unnamed results)
	//
	Implicit(node ast.Node, obj types.Object)

	// Select maps selector expressions (excluding qualified identifiers)
	// to their corresponding selections.
	Select(*ast.SelectorExpr, *types.Selection)

	// Scope maps ast.Nodes to the scopes they define. Package scopes are not
	// associated with a specific node but with all files belonging to a package.
	// Thus, the package scope can be found in the type-checked Package object.
	// Scopes nest, with the Universe scope being the outermost scope, enclosing
	// the package scope, which contains (one or more) files scopes, which enclose
	// function scopes which in turn enclose statement and function literal scopes.
	// Note that even though package-level functions are declared in the package
	// scope, the function scopes are embedded in the file scope of the file
	// containing the function declaration.
	//
	// The following node types may appear in Scopes:
	//
	//     *ast.File
	//     *ast.FuncType
	//     *ast.TypeSpec
	//     *ast.BlockStmt
	//     *ast.IfStmt
	//     *ast.SwitchStmt
	//     *ast.TypeSwitchStmt
	//     *ast.CaseClause
	//     *ast.CommClause
	//     *ast.ForStmt
	//     *ast.RangeStmt
	//
	Scope(ast.Node, *types.Scope)
}

// -----------------------------------------------------------------------------

type Project = modfile.Project
type Class = modfile.Class

// Config of loading Go+ packages.
type Config struct {
	// Types provides type information for the package (optional).
	Types *types.Package

	// Fset provides source position information for syntax trees and types (required).
	Fset *token.FileSet

	// RelativeBase is the root directory of relative path.
	RelativeBase string

	// C2goBase specifies base of standard c2go packages (optional).
	// Default is github.com/goplus/.
	C2goBase string

	// LookupPub lookups the c2go package pubfile named c2go.a.pub (required).
	// See gop/x/c2go.LookupPub.
	LookupPub func(pkgPath string) (pubfile string, err error)

	// LookupClass lookups a class by specified file extension (required).
	// See (*github.com/goplus/mod/gopmod.Module).LookupClass.
	LookupClass func(ext string) (c *Project, ok bool)

	// An Importer resolves import paths to Packages (optional).
	Importer types.Importer

	// A Recorder records existing objects including constants, variables and
	// types etc (optional).
	Recorder Recorder

	// NoFileLine = true means not to generate file line comments.
	NoFileLine bool

	// NoAutoGenMain = true means not to auto generate main func is no entry.
	NoAutoGenMain bool

	// NoSkipConstant = true means to disable optimization of skipping constants.
	NoSkipConstant bool

	// Outline = true means to skip compiling function bodies.
	Outline bool
}

type nodeInterp struct {
	fset       *token.FileSet
	files      map[string]*ast.File
	relBaseDir string
}

func (p *nodeInterp) Position(start token.Pos) (pos token.Position) {
	pos = p.fset.Position(start)
	pos.Filename = relFile(p.relBaseDir, pos.Filename)
	return
}

func (p *nodeInterp) Caller(node ast.Node) string {
	if expr, ok := node.(*ast.CallExpr); ok {
		return p.LoadExpr(expr.Fun)
	}
	return "the function call"
}

func (p *nodeInterp) LoadExpr(node ast.Node) string {
	start := node.Pos()
	pos := p.fset.Position(start)
	f := p.files[pos.Filename]
	n := int(node.End() - start)
	return string(f.Code[pos.Offset : pos.Offset+n])
}

type loader interface {
	load()
	pos() token.Pos
}

type baseLoader struct {
	fn    func()
	start token.Pos
}

func initLoader(ctx *pkgCtx, syms map[string]loader, start token.Pos, name string, fn func(), genBody bool) bool {
	if name == "_" {
		if genBody {
			ctx.inits = append(ctx.inits, fn)
		}
		return false
	}
	if old, ok := syms[name]; ok {
		oldpos := ctx.Position(old.pos())
		ctx.handleErrorf(
			start, "%s redeclared in this block\n\tprevious declaration at %v", name, oldpos)
		return false
	}
	syms[name] = &baseLoader{start: start, fn: fn}
	return true
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

func getTypeLoader(ctx *pkgCtx, syms map[string]loader, start token.Pos, name string) *typeLoader {
	t, ok := syms[name]
	if ok {
		if start != token.NoPos {
			ld := t.(*typeLoader)
			if ld.start == token.NoPos {
				ld.start = start
			} else {
				ctx.handleErrorf(
					start, "%s redeclared in this block\n\tprevious declaration at %v",
					name, ctx.Position(ld.pos()))
			}
			return ld
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
	*nodeInterp
	projs    map[string]*gmxProject // .gmx => project
	classes  map[*ast.File]gmxClass
	fset     *token.FileSet
	cpkgs    *cpackages.Importer
	syms     map[string]loader
	lbinames []any // names that should load before initGopPkg (can be string/func or *ast.Ident/type)
	inits    []func()
	tylds    []*typeLoader
	errs     errors.List

	generics map[string]bool // generic type record
	idents   []*ast.Ident    // toType ident recored
	inInst   int             // toType in generic instance
}

type pkgImp struct {
	gox.PkgRef
	pkgName *types.PkgName
}

type blockCtx struct {
	*pkgCtx
	proj       *gmxProject
	pkg        *gox.Package
	cb         *gox.CodeBuilder
	imports    map[string]pkgImp
	autoimps   map[string]pkgImp
	lookups    []gox.PkgRef
	clookups   []cpackages.PkgRef
	tlookup    *typeParamLookup
	c2goBase   string // default is `github.com/goplus/`
	relBaseDir string
	classRecv  *ast.FieldList // available when gmxSettings != nil
	fileScope  *types.Scope   // only valid when isGopFile
	rec        *typesRecorder
	fileLine   bool
	isClass    bool
	isGopFile  bool // is Go+ file or not
}

func (bc *blockCtx) recorder() *typesRecorder {
	if bc.isGopFile {
		return bc.rec
	}
	return nil
}

func (bc *blockCtx) findImport(name string) (pi pkgImp, ok bool) {
	pi, ok = bc.imports[name]
	if !ok && bc.autoimps != nil {
		pi, ok = bc.autoimps[name]
	}
	return
}

func (p *pkgCtx) newCodeError(pos token.Pos, msg string) error {
	return &gox.CodeError{Fset: p.nodeInterp, Pos: pos, Msg: msg}
}

func (p *pkgCtx) newCodeErrorf(pos token.Pos, format string, args ...interface{}) error {
	return &gox.CodeError{Fset: p.nodeInterp, Pos: pos, Msg: fmt.Sprintf(format, args...)}
}

func (p *pkgCtx) handleErrorf(pos token.Pos, format string, args ...interface{}) {
	p.handleErr(p.newCodeErrorf(pos, format, args...))
}

func (p *pkgCtx) handleErr(err error) {
	p.errs = append(p.errs, err)
}

func (p *pkgCtx) loadNamed(at *gox.Package, t *types.Named) {
	o := t.Obj()
	if o.Pkg() == at.Types {
		p.loadType(o.Name())
	}
}

func (p *pkgCtx) complete() error {
	return p.errs.ToError()
}

func (p *pkgCtx) loadType(name string) {
	if sym, ok := p.syms[name]; ok {
		if ld, ok := sym.(*typeLoader); ok {
			ld.load()
		}
	}
}

func (p *pkgCtx) loadSymbol(name string) bool {
	if enableRecover {
		defer func() {
			if e := recover(); e != nil {
				p.handleRecover(e)
			}
		}()
	}
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

func (p *pkgCtx) handleRecover(e interface{}) {
	err, ok := e.(error)
	if !ok {
		if msg, ok := e.(string); ok {
			err = errors.New(msg)
		} else {
			panic(e)
		}
	}
	p.handleErr(err)
}

const (
	defaultGoFile  = ""
	skippingGoFile = "_skip"
	testingGoFile  = "_test"
)

const (
	ioxPkgPath = "github.com/goplus/gop/builtin/iox"
)

// NewPackage creates a Go+ package instance.
func NewPackage(pkgPath string, pkg *ast.Package, conf *Config) (p *gox.Package, err error) {
	relBaseDir := conf.RelativeBase
	fset := conf.Fset
	files := pkg.Files
	interp := &nodeInterp{
		fset: fset, files: files, relBaseDir: relBaseDir,
	}
	ctx := &pkgCtx{
		fset:       fset,
		nodeInterp: interp,
		projs:      make(map[string]*gmxProject),
		classes:    make(map[*ast.File]gmxClass),
		syms:       make(map[string]loader),
		generics:   make(map[string]bool),
	}
	confGox := &gox.Config{
		Types:           conf.Types,
		Fset:            fset,
		Importer:        conf.Importer,
		LoadNamed:       ctx.loadNamed,
		HandleErr:       ctx.handleErr,
		NodeInterpreter: interp,
		NewBuiltin:      newBuiltinDefault,
		DefaultGoFile:   defaultGoFile,
		NoSkipConstant:  conf.NoSkipConstant,
		PkgPathIox:      ioxPkgPath,
		DbgPositioner:   interp,
	}
	var rec *typesRecorder
	if conf.Recorder != nil {
		rec = newTypeRecord(conf.Recorder)
		confGox.Recorder = &goxRecorder{rec: rec}
	}
	if enableRecover {
		defer func() {
			if e := recover(); e != nil {
				ctx.handleRecover(e)
				err = ctx.errs.ToError()
			}
		}()
	}
	p = gox.NewPackage(pkgPath, pkg.Name, confGox)

	if !noMarkAutogen {
		p.CB().NewConstStart(nil, "_").Val(true).EndInit(1)
	}

	ctx.cpkgs = cpackages.NewImporter(&cpackages.Config{
		Pkg: p, LookupPub: conf.LookupPub,
	})

	// sort files
	type File struct {
		*ast.File
		path string
	}
	sfiles := make([]*File, 0, len(files))
	for fpath, f := range files {
		sfiles = append(sfiles, &File{f, fpath})
	}
	sort.Slice(sfiles, func(i, j int) bool {
		return sfiles[i].path < sfiles[j].path
	})

	for _, f := range sfiles {
		gmx := f.File
		if gmx.IsClass && !gmx.IsNormalGox {
			if debugLoad {
				log.Println("==> File", f.path, "normalGox:", gmx.IsNormalGox)
			}
			loadClass(ctx, p, f.path, gmx, conf)
		}
	}

	for _, f := range sfiles {
		fileLine := !conf.NoFileLine
		fileScope := types.NewScope(p.Types.Scope(), f.Pos(), f.End(), f.path)
		ctx := &blockCtx{
			pkg: p, pkgCtx: ctx, cb: p.CB(), relBaseDir: relBaseDir, fileScope: fileScope,
			fileLine: fileLine, isClass: f.IsClass, rec: rec,
			c2goBase: c2goBase(conf.C2goBase), imports: make(map[string]pkgImp), isGopFile: true,
		}
		if rec := ctx.rec; rec != nil {
			rec.Scope(f, fileScope)
		}
		preloadGopFile(p, ctx, f.path, f.File, conf)
	}

	gopSyms := make(map[string]bool) // TODO: remove this map
	for name := range ctx.syms {
		gopSyms[name] = true
	}

	gofiles := make([]*ast.File, 0, len(pkg.GoFiles))
	for fpath, gof := range pkg.GoFiles {
		f := fromgo.ASTFile(gof, 0)
		gofiles = append(gofiles, f)
		ctx := &blockCtx{
			pkg: p, pkgCtx: ctx, cb: p.CB(), relBaseDir: relBaseDir,
			imports: make(map[string]pkgImp),
		}
		preloadFile(p, ctx, fpath, f, skippingGoFile, false)
	}

	initGopPkg(ctx, p, gopSyms)

	// genMain = true if it is main package and no main func
	var genMain bool
	var gen func()
	if pkg.Name == "main" {
		_, hasMain := ctx.syms["main"]
		genMain = !hasMain
	}

	for _, f := range sfiles {
		if f.IsProj {
			loadFile(ctx, f.File)
		}
	}
	if genMain { // make classfile main func if need
		gen = gmxMainFunc(p, ctx, conf.NoAutoGenMain)
	}

	for _, f := range sfiles {
		if !f.IsProj {
			loadFile(ctx, f.File)
		}
	}
	if conf.Outline {
		for _, f := range gofiles {
			loadFile(ctx, f)
		}
	}
	for _, ld := range ctx.tylds {
		ld.load()
	}
	for _, load := range ctx.inits {
		load()
	}
	err = ctx.complete()

	if gen != nil { // generate classfile main func
		gen()
	} else if genMain && !conf.NoAutoGenMain { // generate empty main func
		old, _ := p.SetCurFile(defaultGoFile, false)
		p.NewFunc(nil, "main", nil, nil, false).BodyStart(p).End()
		p.RestoreCurFile(old)
	}
	return
}

func isOverloadFunc(name string) bool {
	n := len(name)
	return n > 3 && name[n-3:n-1] == "__"
}

func initGopPkg(ctx *pkgCtx, pkg *gox.Package, gopSyms map[string]bool) {
	for name, f := range ctx.syms {
		if gopSyms[name] {
			continue
		}
		if _, ok := f.(*typeLoader); ok {
			ctx.loadType(name)
		} else if isOverloadFunc(name) {
			ctx.loadSymbol(name)
		}
	}
	for _, lbi := range ctx.lbinames {
		if name, ok := lbi.(string); ok {
			ctx.loadSymbol(name)
		} else {
			ctx.loadType(lbi.(*ast.Ident).Name)
		}
	}
	gox.InitThisGopPkg(pkg.Types)
}

func getEntrypoint(f *ast.File) string {
	switch {
	case f.IsProj:
		return "MainEntry"
	case f.IsClass:
		return "Main"
	case f.Name.Name != "main":
		return "init"
	default:
		return "main"
	}
}

func loadFile(ctx *pkgCtx, f *ast.File) {
	for _, decl := range f.Decls {
		switch d := decl.(type) {
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
		case *ast.FuncDecl:
			if d.Recv == nil {
				name := d.Name.Name
				if name != "init" {
					ctx.loadSymbol(name)
				}
			} else {
				if name, ok := getRecvTypeName(ctx, d.Recv, false); ok {
					getTypeLoader(ctx, ctx.syms, token.NoPos, name).load()
				}
			}
		}
	}
}

// gen testingGoFile for:
//
//	*_test.gop
//	*test.gox
func genGoFile(file string, goxTestFile bool) string {
	if goxTestFile || strings.HasSuffix(file, "_test.gop") {
		return testingGoFile
	}
	return defaultGoFile
}

func preloadGopFile(p *gox.Package, ctx *blockCtx, file string, f *ast.File, conf *Config) {
	var proj *gmxProject
	var classType string
	var testType string
	var baseTypeName string
	var baseType types.Type
	var spxClass bool
	var goxTestFile bool
	var parent = ctx.pkgCtx
	if f.IsClass {
		if f.IsNormalGox {
			classType, _ = ClassNameAndExt(file)
			if classType == "main" {
				classType = "_main"
			}
		} else {
			c := parent.classes[f]
			proj, ctx.proj = c.proj, c.proj
			ctx.autoimps = proj.autoimps
			classType = c.tname
			if isGoxTestFile(c.ext) { // test classfile
				testType = c.tname
				goxTestFile, proj.isTest = true, true
				if !f.IsProj {
					classType = casePrefix + testNameSuffix(testType)
				}
			}
			if f.IsProj {
				o := proj.game
				baseTypeName, baseType = o.Name(), o.Type()
				if proj.gameIsPtr {
					baseType = types.NewPointer(baseType)
				}
			} else {
				o := proj.sprite[c.ext]
				baseTypeName, baseType, spxClass = o.Name(), o.Type(), true
			}
		}
	}
	goFile := genGoFile(file, goxTestFile)
	if classType != "" {
		if debugLoad {
			log.Println("==> Preload type", classType)
		}
		if proj != nil {
			ctx.lookups = make([]gox.PkgRef, len(proj.pkgPaths))
			for i, pkgPath := range proj.pkgPaths {
				ctx.lookups[i] = p.Import(pkgPath)
			}
		}
		syms := parent.syms
		pos := f.Pos()
		specs := getFields(f)
		ld := getTypeLoader(parent, syms, pos, classType)
		ld.typ = func() {
			if debugLoad {
				log.Println("==> Load > NewType", classType)
			}
			old, _ := p.SetCurFile(goFile, true)
			defer p.RestoreCurFile(old)

			decl := p.NewTypeDefs().NewType(classType)
			ld.typInit = func() { // decycle
				if debugLoad {
					log.Println("==> Load > InitType", classType)
				}
				old, _ := p.SetCurFile(goFile, true)
				defer p.RestoreCurFile(old)

				pkg := p.Types
				var flds []*types.Var
				var tags []string
				chk := newCheckRedecl()
				if baseTypeName != "" {
					flds = append(flds, types.NewField(pos, pkg, baseTypeName, baseType, true))
					tags = append(tags, "")
					chk.chkRedecl(ctx, baseTypeName, pos)
				}
				if spxClass && proj.gameClass != "" {
					typ := toType(ctx, &ast.StarExpr{X: &ast.Ident{Name: proj.gameClass}})
					name := getTypeName(typ)
					if !chk.chkRedecl(ctx, name, pos) {
						fld := types.NewField(pos, pkg, name, typ, true)
						flds = append(flds, fld)
						tags = append(tags, "")
					}
				}
				rec := ctx.recorder()
				for _, v := range specs {
					spec := v.(*ast.ValueSpec)
					typ := toType(ctx, spec.Type)
					tag := toFieldTag(spec.Tag)
					if len(spec.Names) == 0 {
						name := parseTypeEmbedName(spec.Type)
						if chk.chkRedecl(ctx, name.Name, spec.Type.Pos()) {
							continue
						}
						fld := types.NewField(spec.Type.Pos(), pkg, name.Name, typ, true)
						if rec != nil {
							rec.Def(name, fld)
						}
						flds = append(flds, fld)
						tags = append(tags, tag)
					} else {
						for _, name := range spec.Names {
							if chk.chkRedecl(ctx, name.Name, name.Pos()) {
								continue
							}
							fld := types.NewField(name.Pos(), pkg, name.Name, typ, false)
							if rec != nil {
								rec.Def(name, fld)
							}
							flds = append(flds, fld)
							tags = append(tags, tag)
						}
					}
				}
				decl.InitType(p, types.NewStruct(flds, tags))
			}
			parent.tylds = append(parent.tylds, ld)
		}

		// bugfix: see TestGopxNoFunc
		parent.lbinames = append(parent.lbinames, classType)

		ctx.classRecv = &ast.FieldList{List: []*ast.Field{{
			Names: []*ast.Ident{
				{Name: "this"},
			},
			Type: &ast.StarExpr{
				X: &ast.Ident{Name: classType},
			},
		}}}
	}
	if d := f.ShadowEntry; d != nil {
		d.Name.Name = getEntrypoint(f)
	} else if f.IsProj && !conf.NoAutoGenMain && f.Name.Name == "main" {
		var entry = getEntrypoint(f)
		var hasEntry bool
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Name.Name == entry {
					hasEntry = true
				}
			}
		}
		if !hasEntry {
			f.Decls = append(f.Decls, &ast.FuncDecl{
				Name: &ast.Ident{
					Name: entry,
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
				},
				Body:   &ast.BlockStmt{},
				Shadow: true,
			})
		}
	}
	preloadFile(p, ctx, file, f, goFile, !conf.Outline)
	if goxTestFile {
		parent.inits = append(parent.inits, func() {
			old, _ := p.SetCurFile(testingGoFile, true)
			gmxTestFunc(p, testType, f.IsProj)
			p.RestoreCurFile(old)
		})
	}
}

func parseTypeEmbedName(typ ast.Expr) *ast.Ident {
retry:
	switch t := typ.(type) {
	case *ast.Ident:
		return t
	case *ast.SelectorExpr:
		return t.Sel
	case *ast.StarExpr:
		typ = t.X
		goto retry
	}
	panic("TODO: parseTypeEmbedName unexpected")
}

func preloadFile(p *gox.Package, ctx *blockCtx, file string, f *ast.File, goFile string, genFnBody bool) {
	parent := ctx.pkgCtx
	syms := parent.syms
	old, _ := p.SetCurFile(goFile, true)
	defer p.RestoreCurFile(old)
	var skipClassFields bool
	if f.IsClass {
		skipClassFields = true
	}

	preloadFuncDecl := func(d *ast.FuncDecl) {
		if ctx.classRecv != nil { // in class file (.spx/.gmx)
			if d.Recv == nil {
				d.Recv = ctx.classRecv
				d.IsClass = true
			}
		}
		if d.Recv == nil {
			name := d.Name
			fn := func() {
				old, _ := p.SetCurFile(goFile, true)
				defer p.RestoreCurFile(old)
				loadFunc(ctx, nil, d, genFnBody)
			}
			fname := name.Name
			if fname == "init" {
				if genFnBody {
					if debugLoad {
						log.Println("==> Preload func init")
					}
					parent.inits = append(parent.inits, fn)
				}
			} else {
				if debugLoad {
					log.Println("==> Preload func", fname)
				}
				if initLoader(parent, syms, name.Pos(), fname, fn, genFnBody) {
					if strings.HasPrefix(fname, "Gopx_") { // Gopx_xxx func
						ctx.lbinames = append(ctx.lbinames, fname)
					}
				}
			}
		} else {
			if name, ok := getRecvTypeName(parent, d.Recv, true); ok {
				if debugLoad {
					log.Printf("==> Preload method %s.%s\n", name, d.Name.Name)
				}
				ld := getTypeLoader(parent, syms, token.NoPos, name)
				fn := func() {
					old, _ := p.SetCurFile(goFile, true)
					defer p.RestoreCurFile(old)
					doInitType(ld)
					recv := toRecv(ctx, d.Recv)
					loadFunc(ctx, recv, d, genFnBody)
				}
				ld.methods = append(ld.methods, fn)
			}
		}
	}

	preloadConst := func(d *ast.GenDecl) {
		pkg := ctx.pkg
		cdecl := pkg.NewConstDefs(pkg.Types.Scope())
		for _, spec := range d.Specs {
			vSpec := spec.(*ast.ValueSpec)
			if debugLoad {
				log.Println("==> Preload const", vSpec.Names)
			}
			setNamesLoader(parent, syms, vSpec.Names, func() {
				if c := cdecl; c != nil {
					cdecl = nil
					loadConstSpecs(ctx, c, d.Specs)
					for _, s := range d.Specs {
						v := s.(*ast.ValueSpec)
						removeNames(syms, v.Names)
					}
				}
			})
		}
	}

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				for _, item := range d.Specs {
					loadImport(ctx, item.(*ast.ImportSpec))
				}
			case token.TYPE:
				for _, spec := range d.Specs {
					t := spec.(*ast.TypeSpec)
					tName := t.Name
					name := tName.Name
					if debugLoad {
						log.Println("==> Preload type", name)
					}
					pos := tName.Pos()
					ld := getTypeLoader(parent, syms, pos, name)
					defs := ctx.pkg.NewTypeDefs()
					if goFile != skippingGoFile { // is Go+ file
						ld.typ = func() {
							old, _ := p.SetCurFile(goFile, true)
							defer p.RestoreCurFile(old)
							if t.Assign != token.NoPos { // alias type
								if debugLoad {
									log.Println("==> Load > AliasType", name)
								}
								defs.AliasType(name, toType(ctx, t.Type), tName)
								return
							}
							if debugLoad {
								log.Println("==> Load > NewType", name)
							}
							decl := defs.NewType(name, tName)
							if t.Doc != nil {
								defs.SetComments(t.Doc)
							} else if d.Doc != nil {
								defs.SetComments(d.Doc)
							}
							ld.typInit = func() { // decycle
								if debugLoad {
									log.Println("==> Load > InitType", name)
								}
								decl.InitType(ctx.pkg, toType(ctx, t.Type))
								if rec := ctx.recorder(); rec != nil {
									rec.Def(tName, decl.Type().Obj())
								}
							}
						}
					} else {
						ctx.generics[name] = true
						ld.typ = func() {
							pkg := ctx.pkg.Types
							if t.Assign != token.NoPos { // alias type
								if debugLoad {
									log.Println("==> Load > AliasType", name)
								}
								aliasType(pkg, t.Pos(), name, toType(ctx, t.Type))
								return
							}
							if debugLoad {
								log.Println("==> Load > NewType", name)
							}
							named := newType(pkg, t.Pos(), name)

							ld.typInit = func() { // decycle
								if debugLoad {
									log.Println("==> Load > InitType", name)
								}
								initType(ctx, named, t)
							}
						}
					}
				}
			case token.CONST:
				preloadConst(d)
			case token.VAR:
				if skipClassFields {
					skipClassFields = false
					continue
				}
				for _, spec := range d.Specs {
					vSpec := spec.(*ast.ValueSpec)
					if debugLoad {
						log.Println("==> Preload var", vSpec.Names)
					}
					setNamesLoader(parent, syms, vSpec.Names, func() {
						if v := vSpec; v != nil { // only init once
							vSpec = nil
							old, _ := p.SetCurFile(goFile, true)
							defer p.RestoreCurFile(old)
							loadVars(ctx, v, d.Doc, true)
							removeNames(syms, v.Names)
						}
					})
				}
			default:
				log.Panicln("TODO - tok:", d.Tok, "spec:", reflect.TypeOf(d.Specs).Elem())
			}

		case *ast.FuncDecl:
			preloadFuncDecl(d)

		case *ast.OverloadFuncDecl:
			var recv *ast.Ident
			if d.Recv != nil {
				otyp, ok := d.Recv.List[0].Type.(*ast.Ident)
				if !ok {
					log.Panicln("TODO - cl.preloadFile OverloadFuncDecl: invalid recv")
				}
				recv = otyp
			}
			onames := make([]string, 0, 4)
			exov := false
			name := d.Name
			for idx, fn := range d.Funcs {
				switch expr := fn.(type) {
				case *ast.Ident:
					checkOverloadFunc(d)
					onames = append(onames, expr.Name)
					ctx.lbinames = append(ctx.lbinames, expr.Name)
					exov = true
				case *ast.SelectorExpr:
					checkOverloadMethod(d)
					checkOverloadMethodRecvType(recv, expr.X)
					onames = append(onames, "."+expr.Sel.Name)
					ctx.lbinames = append(ctx.lbinames, recv)
					exov = true
				case *ast.FuncLit:
					checkOverloadFunc(d)
					name1 := overloadFuncName(name.Name, idx)
					ctx.lbinames = append(ctx.lbinames, name1)
					preloadFuncDecl(&ast.FuncDecl{
						Doc:  d.Doc,
						Name: &ast.Ident{NamePos: name.NamePos, Name: name1},
						Type: expr.Type,
						Body: expr.Body,
					})
				default:
					log.Panicf("TODO - cl.preloadFile OverloadFuncDecl: unknown func - %T\n", expr)
				}
			}
			if exov { // need Gopo_xxx
				oname := overloadName(recv, name.Name, d.Operator)
				oval := strings.Join(onames, ",")
				preloadConst(&ast.GenDecl{
					Doc: d.Doc,
					Tok: token.CONST,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names:  []*ast.Ident{{Name: oname}},
							Values: []ast.Expr{stringLit(oval)},
						},
					},
				})
				ctx.lbinames = append(ctx.lbinames, oname)
			}

		default:
			log.Panicf("TODO - cl.preloadFile: unknown decl - %T\n", decl)
		}
	}
}

func checkOverloadFunc(d *ast.OverloadFuncDecl) {
	if d.Recv != nil && !d.Operator {
		log.Panicln("TODO - cl.preloadFile OverloadFuncDecl: checkOverloadFunc")
	}
}

func checkOverloadMethod(d *ast.OverloadFuncDecl) {
	if d.Recv == nil {
		log.Panicln("TODO - cl.preloadFile OverloadFuncDecl: checkOverloadMethod")
	}
}

func checkOverloadMethodRecvType(ot *ast.Ident, recv ast.Expr) {
	rtyp, _ := getRecvType(recv)
	rt, ok := rtyp.(*ast.Ident)
	if !ok || ot.Name != rt.Name {
		log.Panicln("TODO - checkOverloadMethodRecvType:", recv)
	}
}

const (
	indexTable = "0123456789abcdefghijklmnopqrstuvwxyz"
)

func overloadFuncName(name string, idx int) string {
	return name + "__" + indexTable[idx:idx+1]
}

func overloadName(recv *ast.Ident, name string, isOp bool) string {
	if isOp {
		if oname, ok := binaryGopNames[name]; ok {
			name = oname
		} else {
			log.Panicln("TODO - can't overload operator", name)
		}
	}
	sep := "_"
	if strings.ContainsRune(name, '_') || (recv != nil && strings.ContainsRune(recv.Name, '_')) {
		sep = "__"
	}
	typ := ""
	if recv != nil {
		typ = recv.Name + sep
	}
	return "Gopo" + sep + typ + name
}

func stringLit(val string) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(val)}
}

func newType(pkg *types.Package, pos token.Pos, name string) *types.Named {
	typName := types.NewTypeName(pos, pkg, name, nil)
	if old := pkg.Scope().Insert(typName); old != nil {
		log.Panicf("%s redeclared in this block\n\tprevious declaration at %v\n", name, "<TODO>")
	}
	return types.NewNamed(typName, nil, nil)
}

func aliasType(pkg *types.Package, pos token.Pos, name string, typ types.Type) {
	o := types.NewTypeName(pos, pkg, name, typ)
	pkg.Scope().Insert(o)
}

func loadFunc(ctx *blockCtx, recv *types.Var, d *ast.FuncDecl, genBody bool) {
	name := d.Name.Name
	if debugLoad {
		if recv == nil {
			log.Println("==> Load func", name)
		} else {
			log.Printf("==> Load method %v.%s\n", recv.Type(), name)
		}
	}
	pkg := ctx.pkg
	if d.Operator {
		if recv != nil { // binary op
			if v, ok := binaryGopNames[name]; ok {
				name = v
			}
		} else { // unary op
			if v, ok := unaryGopNames[name]; ok {
				name = v
				at := pkg.Types
				arg1 := d.Type.Params.List[0]
				typ := toType(ctx, arg1.Type)
				recv = types.NewParam(arg1.Pos(), at, arg1.Names[0].Name, typ)
				d.Type.Params.List = nil
			}
		}
	}
	sig := toFuncType(ctx, d.Type, recv, d)
	fn, err := pkg.NewFuncWith(d.Name.Pos(), name, sig, func() token.Pos {
		return d.Recv.List[0].Type.Pos()
	})
	if err != nil {
		ctx.handleErr(err)
		return
	}
	commentFunc(ctx, fn, d)
	if rec := ctx.recorder(); rec != nil {
		rec.Def(d.Name, fn.Func)
		if recv == nil {
			ctx.fileScope.Insert(fn.Func)
		}
	}
	if genBody {
		if body := d.Body; body != nil {
			if recv != nil {
				file := pkg.CurFile()
				ctx.inits = append(ctx.inits, func() { // interface issue: #795
					old := pkg.RestoreCurFile(file)
					loadFuncBody(ctx, fn, body, d)
					pkg.RestoreCurFile(old)
				})
			} else {
				loadFuncBody(ctx, fn, body, d)
			}
		}
	}
}

var binaryGopNames = map[string]string{
	"+": "Gop_Add",
	"-": "Gop_Sub",
	"*": "Gop_Mul",
	"/": "Gop_Quo",
	"%": "Gop_Rem",

	"&":  "Gop_And",
	"|":  "Gop_Or",
	"^":  "Gop_Xor",
	"<<": "Gop_Lsh",
	">>": "Gop_Rsh",
	"&^": "Gop_AndNot",

	"+=": "Gop_AddAssign",
	"-=": "Gop_SubAssign",
	"*=": "Gop_MulAssign",
	"/=": "Gop_QuoAssign",
	"%=": "Gop_RemAssign",

	"&=":  "Gop_AndAssign",
	"|=":  "Gop_OrAssign",
	"^=":  "Gop_XorAssign",
	"<<=": "Gop_LshAssign",
	">>=": "Gop_RshAssign",
	"&^=": "Gop_AndNotAssign",

	"==": "Gop_EQ",
	"!=": "Gop_NE",
	"<=": "Gop_LE",
	"<":  "Gop_LT",
	">=": "Gop_GE",
	">":  "Gop_GT",

	"&&": "Gop_LAnd",
	"||": "Gop_LOr",

	"<-": "Gop_Send",
}

var unaryGopNames = map[string]string{
	"++": "Gop_Inc",
	"--": "Gop_Dec",
	"-":  "Gop_Neg",
	"+":  "Gop_Dup",
	"^":  "Gop_Not",
	"!":  "Gop_LNot",
	"<-": "Gop_Recv",
}

func loadFuncBody(ctx *blockCtx, fn *gox.Func, body *ast.BlockStmt, src ast.Node) {
	cb := fn.BodyStart(ctx.pkg, body)
	compileStmts(ctx, body.List)
	cb.End(src)
}

func simplifyGopPackage(pkgPath string) string {
	if strings.HasPrefix(pkgPath, "gop/") {
		return "github.com/goplus/" + pkgPath
	}
	return pkgPath
}

func loadImport(ctx *blockCtx, spec *ast.ImportSpec) {
	if enableRecover {
		defer func() {
			if e := recover(); e != nil {
				ctx.handleRecover(e)
			}
		}()
	}
	var pkg gox.PkgRef
	var pkgPath = toString(spec.Path)
	if realPath, kind := checkC2go(pkgPath); kind != c2goInvalid {
		if kind == c2goStandard {
			realPath = ctx.c2goBase + realPath
		}
		if pkg = loadC2goPkg(ctx, realPath, spec.Path); pkg.Types == nil {
			return
		}
	} else {
		pkg = ctx.pkg.Import(simplifyGopPackage(pkgPath), spec)
	}

	var pos token.Pos
	var name string
	var specName = spec.Name
	if specName != nil {
		name, pos = specName.Name, specName.Pos()
	} else {
		name, pos = pkg.Types.Name(), spec.Path.Pos()
	}
	pkgName := types.NewPkgName(pos, ctx.pkg.Types, name, pkg.Types)
	if rec := ctx.recorder(); rec != nil {
		if specName != nil {
			rec.Def(specName, pkgName)
		} else {
			rec.Implicit(spec, pkgName)
		}
	}

	if specName != nil {
		if name == "." {
			ctx.lookups = append(ctx.lookups, pkg)
			return
		}
		if name == "_" {
			pkg.MarkForceUsed(ctx.pkg)
			return
		}
	}
	ctx.imports[name] = pkgImp{pkg, pkgName}
}

func loadConstSpecs(ctx *blockCtx, cdecl *gox.ConstDefs, specs []ast.Spec) {
	for iotav, spec := range specs {
		vSpec := spec.(*ast.ValueSpec)
		loadConsts(ctx, cdecl, vSpec, iotav)
	}
}

func loadConsts(ctx *blockCtx, cdecl *gox.ConstDefs, v *ast.ValueSpec, iotav int) {
	vNames := v.Names
	names := makeNames(vNames)
	if v.Values == nil {
		if debugLoad {
			log.Println("==> Load const", names)
		}
		cdecl.Next(iotav, v.Pos(), names...)
		defNames(ctx, vNames, nil)
		return
	}
	var typ types.Type
	if v.Type != nil {
		typ = toType(ctx, v.Type)
	}
	if debugLoad {
		log.Println("==> Load const", names, typ)
	}
	fn := func(cb *gox.CodeBuilder) int {
		for _, val := range v.Values {
			compileExpr(ctx, val)
		}
		return len(v.Values)
	}
	cdecl.New(fn, iotav, v.Pos(), typ, names...)
	defNames(ctx, v.Names, nil)
}

func loadVars(ctx *blockCtx, v *ast.ValueSpec, doc *ast.CommentGroup, global bool) {
	var typ types.Type
	if v.Type != nil {
		typ = toType(ctx, v.Type)
	}
	names := makeNames(v.Names)
	if debugLoad {
		log.Println("==> Load var", typ, names)
	}
	var scope *types.Scope
	if global {
		scope = ctx.pkg.Types.Scope()
	} else {
		scope = ctx.cb.Scope()
	}
	varDefs := ctx.pkg.NewVarDefs(scope).SetComments(doc)
	varDecl := varDefs.New(v.Names[0].Pos(), typ, names...)
	if nv := len(v.Values); nv > 0 {
		cb := varDecl.InitStart(ctx.pkg)
		if enableRecover {
			defer func() {
				if e := recover(); e != nil {
					cb.ResetInit()
					panic(e)
				}
			}()
		}
		if nv == 1 && len(names) == 2 {
			compileExpr(ctx, v.Values[0], clCallWithTwoValue)
		} else {
			for _, val := range v.Values {
				switch e := val.(type) {
				case *ast.LambdaExpr, *ast.LambdaExpr2:
					if len(v.Values) == 1 {
						sig := checkLambdaFuncType(ctx, e, typ, clLambaAssign, v.Names[0])
						compileLambda(ctx, e, sig)
					} else {
						panic(ctx.newCodeErrorf(e.Pos(), "lambda unsupport multiple assignment"))
					}
				case *ast.SliceLit:
					compileSliceLit(ctx, e, typ)
				case *ast.CompositeLit:
					compileCompositeLit(ctx, e, typ, false)
				default:
					compileExpr(ctx, val)
				}
			}
		}
		cb.EndInit(nv)
	}
	defNames(ctx, v.Names, scope)
}

func defNames(ctx *blockCtx, names []*ast.Ident, scope *types.Scope) {
	if rec := ctx.recorder(); rec != nil {
		if scope == nil {
			scope = ctx.cb.Scope()
		}
		for _, name := range names {
			if o := scope.Lookup(name.Name); o != nil {
				rec.Def(name, o)
			}
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

func removeNames(syms map[string]loader, names []*ast.Ident) {
	for _, name := range names {
		delete(syms, name.Name)
	}
}

func setNamesLoader(ctx *pkgCtx, syms map[string]loader, names []*ast.Ident, load func()) {
	for _, name := range names {
		initLoader(ctx, syms, name.Pos(), name.Name, load, true)
	}
}

// -----------------------------------------------------------------------------
