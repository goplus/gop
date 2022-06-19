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
	"os"
	"reflect"
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/ast/fromgo"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
	"github.com/goplus/gox/cpackages"
	"github.com/goplus/mod/modfile"
	"github.com/qiniu/x/errors"
)

const (
	DbgFlagLoad = 1 << iota
	DbgFlagLookup
	DbgFlagAll = DbgFlagLoad | DbgFlagLookup
)

var (
	enableRecover = true
)

var (
	debugLoad   bool
	debugLookup bool
)

func SetDisableRecover(disableRecover bool) {
	enableRecover = !disableRecover
}

func SetDebug(flags int) {
	debugLoad = (flags & DbgFlagLoad) != 0
	debugLookup = (flags & DbgFlagLookup) != 0
}

// -----------------------------------------------------------------------------

type Class = modfile.Classfile

// Config of loading Go+ packages.
type Config struct {
	// Fset provides source position information for syntax trees and types.
	// If Fset is nil, Load will use a new fileset, but preserve Fset's value.
	Fset *token.FileSet

	// WorkingDir is the directory in which to run gop compiler.
	WorkingDir string

	// TargetDir is the directory in which to generate Go files.
	TargetDir string

	// C2goBase specifies base of standard c2go packages.
	// Default is github.com/goplus/.
	C2goBase string

	// LookupPub lookups the c2go package pubfile (named c2go.a.pub).
	LookupPub func(pkgPath string) (pubfile string, err error)

	// LookupClass lookups a class by specified file extension.
	LookupClass func(ext string) (c *Class, ok bool)

	// An Importer resolves import paths to Packages.
	Importer types.Importer

	// NoFileLine = true means not to generate file line comments.
	NoFileLine bool

	// RelativePath = true means to generate file line comments with relative file path.
	RelativePath bool

	// NoAutoGenMain = true means not to auto generate main func is no entry.
	NoAutoGenMain bool

	// NoSkipConstant = true means disable optimization of skip constants
	NoSkipConstant bool
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
	if pos.Offset+n < 0 {
		log.Println("LoadExpr:", node, pos.Filename, pos.Line, pos.Offset, node.Pos(), node.End(), n)
	}
	src = string(f.Code[pos.Offset : pos.Offset+n])
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

func initLoader(ctx *pkgCtx, syms map[string]loader, start token.Pos, name string, fn func(), genCode bool) {
	if name == "_" {
		if genCode {
			ctx.inits = append(ctx.inits, fn)
		}
		return
	}
	if old, ok := syms[name]; ok {
		var pos token.Position
		if start != token.NoPos {
			pos = ctx.Position(start)
		}
		oldpos := ctx.Position(old.pos())
		ctx.handleCodeErrorf(
			&pos, "%s redeclared in this block\n\tprevious declaration at %v", name, oldpos)
		return
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

func getTypeLoader(ctx *pkgCtx, syms map[string]loader, start token.Pos, name string) *typeLoader {
	t, ok := syms[name]
	if ok {
		if start != token.NoPos {
			ld := t.(*typeLoader)
			if ld.start == token.NoPos {
				ld.start = start
			} else {
				pos := ctx.Position(start)
				ctx.handleCodeErrorf(&pos, "%s redeclared in this block\n\tprevious declaration at %v",
					name, ctx.Position(ld.start))
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
	*gmxSettings
	cpkgs *cpackages.Importer
	syms  map[string]loader
	inits []func()
	tylds []*typeLoader
	errs  errors.List
}

type blockCtx struct {
	*pkgCtx
	pkg          *gox.Package
	cb           *gox.CodeBuilder
	fset         *token.FileSet
	imports      map[string]*gox.PkgRef
	lookups      []*gox.PkgRef
	clookups     []*cpackages.PkgRef
	c2goBase     string // default is `github.com/goplus/`
	targetDir    string
	classRecv    *ast.FieldList // available when gmxSettings != nil
	fileLine     bool
	relativePath bool
	isClass      bool
}

func (bc *blockCtx) findImport(name string) (pr *gox.PkgRef, ok bool) {
	pr, ok = bc.imports[name]
	return
}

func newCodeErrorf(pos *token.Position, format string, args ...interface{}) *gox.CodeError {
	return &gox.CodeError{Pos: pos, Msg: fmt.Sprintf(format, args...)}
}

func (p *pkgCtx) newCodeError(start token.Pos, msg string) error {
	pos := p.Position(start)
	return &gox.CodeError{Pos: &pos, Msg: msg}
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

// NewPackage creates a Go+ package instance.
func NewPackage(pkgPath string, pkg *ast.Package, conf *Config) (p *gox.Package, err error) {
	workingDir := conf.WorkingDir
	if workingDir == "" {
		workingDir, _ = os.Getwd()
	}
	targetDir := conf.TargetDir
	if targetDir == "" {
		targetDir = workingDir
	}
	fset := conf.Fset
	files := pkg.Files
	interp := &nodeInterp{
		fset: fset, files: files, workingDir: workingDir,
	}
	ctx := &pkgCtx{
		syms: make(map[string]loader), nodeInterp: interp,
	}
	confGox := &gox.Config{
		Fset:            fset,
		Importer:        conf.Importer,
		LoadNamed:       ctx.loadNamed,
		HandleErr:       ctx.handleErr,
		NodeInterpreter: interp,
		NewBuiltin:      newBuiltinDefault,
		DefaultGoFile:   defaultGoFile,
		NoSkipConstant:  conf.NoSkipConstant,
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
	ctx.cpkgs = cpackages.NewImporter(&cpackages.Config{
		Pkg: p, LookupPub: conf.LookupPub,
	})
	for file, gmx := range files {
		if gmx.IsProj {
			ctx.gmxSettings = newGmx(ctx, p, file, conf)
			break
		}
	}
	for fpath, f := range files {
		fileLine := !conf.NoFileLine
		ctx := &blockCtx{
			pkg: p, pkgCtx: ctx, cb: p.CB(), fset: p.Fset, targetDir: targetDir,
			fileLine: fileLine, relativePath: conf.RelativePath, isClass: f.IsClass,
			c2goBase: c2goBase(conf.C2goBase), imports: make(map[string]*gox.PkgRef),
		}
		preloadGopFile(p, ctx, fpath, f, conf)
	}
	for fpath, gof := range pkg.GoFiles {
		f := fromgo.ASTFile(gof, 0)
		ctx := &blockCtx{
			pkg: p, pkgCtx: ctx, cb: p.CB(), fset: p.Fset, targetDir: targetDir,
			imports: make(map[string]*gox.PkgRef),
		}
		preloadFile(p, ctx, fpath, f, false)
	}
	for _, f := range files {
		if f.IsProj {
			loadFile(ctx, f)
			gmxMainFunc(p, ctx)
			break
		}
	}
	for _, f := range files {
		if !f.IsProj { // only one .gmx file
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

	if !conf.NoAutoGenMain && pkg.Name == "main" {
		if obj := p.Types.Scope().Lookup("main"); obj == nil {
			old, _ := p.SetCurFile(defaultGoFile, false)
			p.NewFunc(nil, "main", nil, nil, false).BodyStart(p).End()
			p.RestoreCurFile(old)
		}
	}
	return
}

func hasMethod(o types.Object, name string) bool {
	if obj, ok := o.(*types.TypeName); ok {
		if t, ok := obj.Type().(*types.Named); ok {
			for i, n := 0, t.NumMethods(); i < n; i++ {
				if t.Method(i).Name() == name {
					return true
				}
			}
		}
	}
	return false
}

func getEntrypoint(f *ast.File, isMod bool) string {
	switch {
	case f.IsProj:
		return "MainEntry"
	case f.IsClass:
		return "Main"
	case isMod:
		return "init"
	default:
		return "main"
	}
}

func loadFile(ctx *pkgCtx, f *ast.File) {
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
					getTypeLoader(ctx, ctx.syms, token.NoPos, name).load()
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

func getGoFile(file string, genCode bool) string {
	if genCode {
		if strings.HasSuffix(file, "_test.gop") {
			return testingGoFile
		}
		return defaultGoFile
	}
	return skippingGoFile
}

func preloadGopFile(p *gox.Package, ctx *blockCtx, file string, f *ast.File, conf *Config) {
	var parent = ctx.pkgCtx
	var classType string
	var baseTypeName string
	var baseType types.Type
	var spxClass bool
	switch {
	case f.IsProj:
		classType = parent.gameClass
		o := parent.game
		baseTypeName, baseType = o.Name(), o.Type()
		if parent.gameIsPtr {
			baseType = types.NewPointer(baseType)
		}
	case f.IsClass:
		if parent.gmxSettings != nil {
			classType = getDefaultClass(file)
			o := parent.sprite
			baseTypeName, baseType, spxClass = o.Name(), o.Type(), true
		}
		// TODO: panic
	}
	if classType != "" {
		if debugLoad {
			log.Println("==> Preload type", classType)
		}
		ctx.lookups = make([]*gox.PkgRef, len(parent.pkgPaths))
		for i, pkgPath := range parent.pkgPaths {
			ctx.lookups[i] = p.Import(pkgPath)
		}
		syms := parent.syms
		pos := f.Pos()
		specs := getFields(ctx, f)
		ld := getTypeLoader(parent, syms, pos, classType)
		ld.typ = func() {
			if debugLoad {
				log.Println("==> Load > NewType", classType)
			}
			decl := p.NewType(classType)
			ld.typInit = func() { // decycle
				if debugLoad {
					log.Println("==> Load > InitType", classType)
				}
				pkg := p.Types
				flds := make([]*types.Var, 1, 2)
				flds[0] = types.NewField(pos, pkg, baseTypeName, baseType, true)
				if spxClass {
					typ := toType(ctx, &ast.StarExpr{X: &ast.Ident{Name: parent.gameClass}})
					fld := types.NewField(pos, pkg, getTypeName(typ), typ, true)
					flds = append(flds, fld)
				}
				for _, v := range specs {
					spec := v.(*ast.ValueSpec)
					if spec.Type == nil {
						pos := ctx.Position(v.Pos())
						ctx.handleCodeErrorf(&pos, "missing field type in class file")
						continue
					}
					if len(spec.Values) > 0 {
						pos := ctx.Position(v.Pos())
						ctx.handleCodeErrorf(&pos, "cannot assign value to field in class file")
					}
					typ := toType(ctx, spec.Type)
					for _, name := range spec.Names {
						flds = append(flds, types.NewField(name.Pos(), pkg, name.Name, typ, false))
					}
				}
				decl.InitType(p, types.NewStruct(flds, nil))
			}
			parent.tylds = append(parent.tylds, ld)
		}
		ctx.classRecv = &ast.FieldList{List: []*ast.Field{{
			Names: []*ast.Ident{
				{Name: "this"},
			},
			Type: &ast.StarExpr{
				X: &ast.Ident{Name: classType},
			},
		}}}
	}
	// check class project no MainEntry and auto added
	if f.IsProj && !conf.NoAutoGenMain && !f.NoEntrypoint && f.Name.Name == "main" {
		entry := getEntrypoint(f, false)
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
				Name: ast.NewIdent(entry),
				Type: &ast.FuncType{Params: &ast.FieldList{}},
				Body: &ast.BlockStmt{},
			})
		}
	}
	preloadFile(p, ctx, file, f, true)
}

func preloadFile(p *gox.Package, ctx *blockCtx, file string, f *ast.File, genCode bool) {
	parent := ctx.pkgCtx
	syms := parent.syms
	goFile := getGoFile(file, genCode)
	old, _ := p.SetCurFile(goFile, true)
	defer p.RestoreCurFile(old)
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if f.NoEntrypoint && d.Name.Name == "main" {
				d.Name.Name = getEntrypoint(f, f.Name.Name != "main")
			}
			if ctx.classRecv != nil { // in class file (.spx/.gmx)
				if d.Recv == nil {
					d.Recv = ctx.classRecv
				}
			}
			if d.Recv == nil {
				var name = d.Name
				var fn func()
				if genCode {
					fn = func() {
						old, _ := p.SetCurFile(goFile, true)
						defer p.RestoreCurFile(old)
						loadFunc(ctx, nil, d)
					}
				} else {
					fn = func() {
						declFunc(ctx, nil, d)
					}
				}
				if name.Name == "init" {
					if genCode {
						if debugLoad {
							log.Println("==> Preload func init")
						}
						parent.inits = append(parent.inits, fn)
					}
				} else {
					if debugLoad {
						log.Println("==> Preload func", name.Name)
					}
					initLoader(parent, syms, name.Pos(), name.Name, fn, genCode)
				}
			} else {
				if name, ok := getRecvTypeName(parent, d.Recv, true); ok {
					if debugLoad {
						log.Printf("==> Preload method %s.%s\n", name, d.Name.Name)
					}
					var ld = getTypeLoader(parent, syms, token.NoPos, name)
					var fn func()
					if genCode {
						fn = func() {
							old, _ := p.SetCurFile(goFile, true)
							defer p.RestoreCurFile(old)
							doInitType(ld)
							recv := toRecv(ctx, d.Recv)
							loadFunc(ctx, recv, d)
						}
					} else {
						fn = func() {
							doInitType(ld)
							recv := toRecv(ctx, d.Recv)
							declFunc(ctx, recv, d)
						}
					}
					ld.methods = append(ld.methods, fn)
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
					if debugLoad {
						log.Println("==> Preload type", name)
					}
					ld := getTypeLoader(parent, syms, t.Name.Pos(), name)
					if genCode {
						ld.typ = func() {
							old, _ := p.SetCurFile(goFile, true)
							defer p.RestoreCurFile(old)
							if t.Assign != token.NoPos { // alias type
								if debugLoad {
									log.Println("==> Load > AliasType", name)
								}
								ctx.pkg.AliasType(name, toType(ctx, t.Type), t.Pos())
								return
							}
							if debugLoad {
								log.Println("==> Load > NewType", name)
							}
							decl := ctx.pkg.NewType(name)
							if t.Doc != nil {
								decl.SetComments(t.Doc)
							} else if d.Doc != nil {
								decl.SetComments(d.Doc)
							}
							ld.typInit = func() { // decycle
								if debugLoad {
									log.Println("==> Load > InitType", name)
								}
								decl.InitType(ctx.pkg, toType(ctx, t.Type))
							}
						}
					} else {
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
								initType(ctx, named, toType(ctx, t.Type))
							}
						}
					}
				}
			case token.CONST:
				pkg := ctx.pkg
				cdecl := pkg.NewConstDecl(pkg.Types.Scope())
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
			case token.VAR:
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
							loadVars(ctx, v, true)
							removeNames(syms, v.Names)
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

func newType(pkg *types.Package, pos token.Pos, name string) *types.Named {
	typName := types.NewTypeName(pos, pkg, name, nil)
	if old := pkg.Scope().Insert(typName); old != nil {
		log.Panicf("%s redeclared in this block\n\tprevious declaration at %v\n", name, "<TODO>")
	}
	return types.NewNamed(typName, nil, nil)
}

func initType(ctx *blockCtx, named *types.Named, typ types.Type) {
	if named, ok := typ.(*types.Named); ok {
		typ = getUnderlying(ctx, named)
	}
	named.SetUnderlying(typ)
}

func aliasType(pkg *types.Package, pos token.Pos, name string, typ types.Type) {
	o := types.NewTypeName(pos, pkg, name, typ)
	pkg.Scope().Insert(o)
}

func declFunc(ctx *blockCtx, recv *types.Var, d *ast.FuncDecl) {
	name := d.Name.Name
	if debugLoad {
		if recv == nil {
			log.Println("==> Load func", name)
		} else {
			log.Printf("==> Load method %v.%s\n", recv.Type(), name)
		}
	}
	if name == "_" {
		return
	}
	pkg := ctx.pkg.Types
	sig := toFuncType(ctx, d.Type, recv)
	fn := types.NewFunc(d.Pos(), pkg, name, sig)
	if recv != nil {
		typ := recv.Type()
		switch t := typ.(type) {
		case *types.Named:
			t.AddMethod(fn)
			return
		case *types.Pointer:
			if tt, ok := t.Elem().(*types.Named); ok {
				tt.AddMethod(fn)
				return
			}
		}
		log.Panicf("invalid receiver type %v (%v is not a defined type)\n", typ, typ)
	} else {
		pkg.Scope().Insert(fn)
	}
}

func loadFunc(ctx *blockCtx, recv *types.Var, d *ast.FuncDecl) {
	name := d.Name.Name
	if debugLoad {
		if recv == nil {
			log.Println("==> Load func", name)
		} else {
			log.Printf("==> Load method %v.%s\n", recv.Type(), name)
		}
	}
	if d.Operator {
		if recv != nil { // binary op
			if v, ok := binaryGopNames[name]; ok {
				name = v
			}
		} else { // unary op
			if v, ok := unaryGopNames[name]; ok {
				name = v
				at := ctx.pkg.Types
				arg1 := d.Type.Params.List[0]
				typ := toType(ctx, arg1.Type)
				recv = types.NewParam(arg1.Pos(), at, arg1.Names[0].Name, typ)
				d.Type.Params.List = nil
			}
		}
	}
	sig := toFuncType(ctx, d.Type, recv)
	fn, err := ctx.pkg.NewFuncWith(d.Pos(), name, sig, func() token.Pos {
		return d.Recv.List[0].Type.Pos()
	})
	if err != nil {
		ctx.handleErr(err)
		return
	}
	if d.Doc != nil {
		fn.SetComments(d.Doc)
	}
	if body := d.Body; body != nil {
		if recv != nil {
			ctx.inits = append(ctx.inits, func() { // interface issue: #795
				loadFuncBody(ctx, fn, body)
			})
		} else {
			loadFuncBody(ctx, fn, body)
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

func loadFuncBody(ctx *blockCtx, fn *gox.Func, body *ast.BlockStmt) {
	cb := fn.BodyStart(ctx.pkg)
	compileStmts(ctx, body.List)
	cb.End()
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
	var pkg *gox.PkgRef
	var pkgPath = toString(spec.Path)
	if realPath, kind := checkC2go(pkgPath); kind != c2goInvalid {
		if kind == c2goStandard {
			realPath = ctx.c2goBase + realPath
		}
		if pkg = loadC2goPkg(ctx, realPath, spec.Path); pkg == nil {
			return
		}
	} else {
		pkg = ctx.pkg.Import(simplifyGopPackage(pkgPath), spec)
	}
	var name string
	if spec.Name != nil {
		name = spec.Name.Name
		if name == "." {
			ctx.lookups = append(ctx.lookups, pkg)
			return
		}
		if name == "_" {
			pkg.MarkForceUsed()
			return
		}
	} else {
		name = pkg.Types.Name()
	}
	ctx.imports[name] = pkg
}

func loadConstSpecs(ctx *blockCtx, cdecl *gox.ConstDecl, specs []ast.Spec) {
	for iotav, spec := range specs {
		vSpec := spec.(*ast.ValueSpec)
		loadConsts(ctx, cdecl, vSpec, iotav)
	}
}

func loadConsts(ctx *blockCtx, cdecl *gox.ConstDecl, v *ast.ValueSpec, iotav int) {
	names := makeNames(v.Names)
	if v.Values == nil {
		if debugLoad {
			log.Println("==> Load const", names)
		}
		cdecl.Next(iotav, v.Pos(), names...)
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
}

func loadVars(ctx *blockCtx, v *ast.ValueSpec, global bool) {
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
	varDecl := ctx.pkg.NewVarEx(scope, v.Names[0].Pos(), typ, names...)
	if nv := len(v.Values); nv > 0 {
		cb := varDecl.InitStart(ctx.pkg)
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
