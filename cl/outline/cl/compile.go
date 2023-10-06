/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
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

package cl

import (
	"fmt"
	"go/types"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/ast/fromgo"
	"github.com/goplus/gop/parser"
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

func SetDebug(flags dbgFlags) {
	debugLoad = (flags & DbgFlagLoad) != 0
	debugLookup = (flags & DbgFlagLookup) != 0
}

// -----------------------------------------------------------------------------

type Project = modfile.Project

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
	LookupClass func(ext string) (c *Project, ok bool)

	// An Importer resolves import paths to Packages.
	Importer types.Importer

	// NoFileLine = true means not to generate file line comments.
	NoFileLine bool

	// RelativePath = true means to generate file line comments with relative file path.
	RelativePath bool

	// NoAutoGenMain = true means not to auto generate main func is no entry.
	NoAutoGenMain bool

	// NoSkipConstant = true means to disable optimization of skipping constants.
	NoSkipConstant bool
}

// -----------------------------------------------------------------------------

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
	if f == nil { // not found
		return
	}
	n := int(node.End() - start)
	pos.Filename = relFile(p.workingDir, pos.Filename)
	if pos.Offset+n < 0 {
		log.Println("LoadExpr:", node, pos.Filename, pos.Line, pos.Offset, node.Pos(), node.End(), n)
	}
	src = string(f.Code[pos.Offset : pos.Offset+n])
	return
}

// -----------------------------------------------------------------------------

type lazyKind int

const (
	lazyAll lazyKind = iota
	lazyAliasType
	lazyInitType
	lazyVarOrConst
)

type lazy struct {
	pos  token.Pos
	load func()
	kind lazyKind
}

func (f *lazy) resolve(kind lazyKind) bool {
	if load := f.load; load != nil && (kind == lazyAll || f.kind == kind) {
		f.load = nil
		load()
		return true
	}
	return false
}

type lazyobjs map[string]*lazy

func (p lazyobjs) insert(ctx *pkgCtx, name string, o *lazy) {
	if old, ok := p[name]; ok {
		npos := ctx.Position(o.pos)
		ctx.handleCodeErrorf(&npos,
			"%v redeclared\n\t%v other declaration of %v", name, ctx.Position(old.pos), name)
	}
	p[name] = o
}

func (p lazyobjs) load(ctx *pkgCtx, name string, kind lazyKind) bool {
	if enableRecover {
		defer func() {
			if e := recover(); e != nil {
				ctx.handleRecover(e)
			}
		}()
	}
	if f, ok := p[name]; ok {
		return f.resolve(kind)
	}
	return false
}

// -----------------------------------------------------------------------------

type pkgCtx struct {
	*nodeInterp
	*gmxSettings
	cpkgs *cpackages.Importer
	lazys lazyobjs // delay: (1) alias types; (2) init named type; (4) global constants/variables
	funcs []func() // delay: (3) func prototypes
	bodys []func() // delay: (5) generate function body
	errs  errors.List
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

func (p *pkgCtx) insertObject(name string, pos token.Pos, kind lazyKind, load func()) {
	p.lazys.insert(p, name, &lazy{pos, load, kind})
}

func (p *pkgCtx) insertNames(pos token.Pos, names []*ast.Ident, load func()) {
	lazy := &lazy{pos, load, lazyVarOrConst}
	for _, name := range names {
		p.lazys.insert(p, name.Name, lazy)
	}
}

func (p *pkgCtx) loadObject(name string) bool {
	return p.lazys.load(p, name, lazyAll)
}

// -----------------------------------------------------------------------------

type blockCtx struct {
	*pkgCtx
	pkg          *gox.Package
	cb           *gox.CodeBuilder
	imports      map[string]*gox.PkgRef
	lookups      []*gox.PkgRef
	clookups     []*cpackages.PkgRef
	tlookup      *typeParamLookup
	c2goBase     string // default is `github.com/goplus/`
	targetDir    string
	classRecv    *ast.FieldList // available when gmxSettings != nil
	isClass      bool
	fileLine     bool
	relativePath bool
}

func (p *blockCtx) findImport(name string) (pr *gox.PkgRef, ok bool) {
	pr, ok = p.imports[name]
	return
}

// -----------------------------------------------------------------------------

// Context saves package level states during compiling.
type Context = pkgCtx

func (p *Context) Complete(pkg *gox.Package, autoGenMain bool) error {
	for _, genBody := range p.bodys {
		genBody()
	}
	pkgTypes := pkg.Types
	if autoGenMain && pkgTypes.Name() == "main" {
		if obj := pkgTypes.Scope().Lookup("main"); obj == nil {
			old, _ := pkg.SetCurFile(defaultGoFile, false)
			pkg.NewFunc(nil, "main", nil, nil, false).BodyStart(pkg).End()
			pkg.RestoreCurFile(old)
		}
	}
	return p.errs.ToError()
}

func (p *Context) genOutline() error {
	for _, f := range p.lazys {
		f.resolve(lazyAliasType)
	}
	for _, f := range p.lazys {
		f.resolve(lazyInitType)
	}
	for _, fnProto := range p.funcs {
		fnProto()
	}
	for _, f := range p.lazys {
		f.resolve(lazyVarOrConst)
	}
	return p.errs.ToError()
}

// -----------------------------------------------------------------------------

// NewPackage creates a Go+ package instance.
func NewPackage(pkgPath string, pkg *ast.Package, conf *Config) (p *gox.Package, err error) {
	p, ctx, err := NewOutline(pkgPath, pkg, conf)
	if err != nil {
		return
	}
	err = ctx.Complete(p, !conf.NoAutoGenMain)
	return
}

// -----------------------------------------------------------------------------

const (
	defaultGoFile  = ""
	skippingGoFile = "_skip"
	testingGoFile  = "_test"
)

const (
	ioxPkgPath = "github.com/goplus/gop/builtin/iox"
)

// NewOutline creates a Go+ package outline.
func NewOutline(pkgPath string, pkg *ast.Package, conf *Config) (p *gox.Package, ctx *Context, err error) {
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
	ctx = &pkgCtx{
		nodeInterp: interp,
		lazys:      make(lazyobjs),
	}
	confGox := &gox.Config{
		Fset:            fset,
		Importer:        conf.Importer,
		HandleErr:       ctx.handleErr,
		NodeInterpreter: interp,
		NewBuiltin:      newBuiltinDefault,
		DefaultGoFile:   defaultGoFile,
		NoSkipConstant:  conf.NoSkipConstant,
		PkgPathIox:      ioxPkgPath,
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
	for file, gmx := range files {
		if gmx.IsProj {
			ctx.gmxSettings = newGmx(ctx, p, file, gmx, conf)
			break
		}
	}
	if ctx.gmxSettings == nil {
		for file, gmx := range files {
			if gmx.IsClass && !gmx.IsNormalGox {
				ctx.gmxSettings = newGmx(ctx, p, file, gmx, conf)
				break
			}
		}
	}

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

	// fset = p.Fset // conf.Fset maybe is nil
	for _, f := range sfiles {
		fileLine := !conf.NoFileLine
		ctx := &blockCtx{
			pkg: p, pkgCtx: ctx, cb: p.CB(), targetDir: targetDir,
			fileLine: fileLine, relativePath: conf.RelativePath, isClass: f.IsClass,
			c2goBase: c2goBase(conf.C2goBase), imports: make(map[string]*gox.PkgRef),
		}
		preloadGopFile(p, ctx, f.path, f.File, conf)
	}

	for fpath, gof := range pkg.GoFiles {
		f := fromgo.ASTFile(gof, 0)
		ctx := &blockCtx{
			pkg: p, pkgCtx: ctx, cb: p.CB(), targetDir: targetDir,
			imports: make(map[string]*gox.PkgRef),
		}
		preloadFile(p, ctx, fpath, f, false)
	}
	err = ctx.genOutline()
	return
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

func genGoFile(file string, gopFile bool) string {
	if gopFile {
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
		classType = getDefaultClass(file)
		if parent.gmxSettings != nil {
			ext := parser.ClassFileExt(file)
			o, ok := parent.sprite[ext]
			if ok {
				baseTypeName, baseType, spxClass = o.Name(), o.Type(), true
			}
		}
	}
	if classType != "" {
		ctx.classRecv = &ast.FieldList{List: []*ast.Field{{
			Names: []*ast.Ident{
				{Name: "this"},
			},
			Type: &ast.StarExpr{
				X: &ast.Ident{Name: classType},
			},
		}}}
		if parent.gmxSettings != nil {
			ctx.lookups = make([]*gox.PkgRef, len(parent.pkgPaths))
			for i, pkgPath := range parent.pkgPaths {
				ctx.lookups[i] = p.Import(pkgPath)
			}
		}
		if debugLoad {
			log.Println("==> Load > NewType", classType)
		}
		pos := f.Pos()
		specs := getFields(f)
		decl := p.NewTypeDefs().NewType(classType, pos)
		parent.insertObject(classType, pos, lazyInitType, func() {
			if debugLoad {
				log.Println("==> Load > InitType", classType)
			}
			pkg := p.Types
			var flds []*types.Var
			var tags []string
			chk := newCheckRedecl()
			if len(baseTypeName) != 0 {
				flds = append(flds, types.NewField(pos, pkg, baseTypeName, baseType, true))
				tags = append(tags, "")
				chk.chkRedecl(ctx, baseTypeName, pos)
			}
			if spxClass && parent.gmxSettings != nil && parent.gameClass != "" {
				typ := toType(ctx, &ast.StarExpr{X: &ast.Ident{Name: parent.gameClass}})
				name := getTypeName(typ)
				if !chk.chkRedecl(ctx, name, pos) {
					fld := types.NewField(pos, pkg, name, typ, true)
					flds = append(flds, fld)
					tags = append(tags, "")
				}
			}
			for _, v := range specs {
				spec := v.(*ast.ValueSpec)
				var embed bool
				if spec.Names == nil {
					embed = true
					v := parseTypeEmbedName(spec.Type)
					spec.Names = []*ast.Ident{v}
				}
				typ := toType(ctx, spec.Type)
				for _, name := range spec.Names {
					if chk.chkRedecl(ctx, name.Name, name.Pos()) {
						continue
					}
					flds = append(flds, types.NewField(name.Pos(), pkg, name.Name, typ, embed))
					tags = append(tags, toFieldTag(spec.Tag))
				}
			}
			decl.InitType(p, types.NewStruct(flds, tags))
		})
	}
	// check class project no MainEntry and auto added
	if f.IsProj && !conf.NoAutoGenMain && !f.NoEntrypoint() && f.Name.Name == "main" {
		entry := getEntrypoint(f)
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
	if d := f.ShadowEntry; d != nil {
		d.Name.Name = getEntrypoint(f)
	}
	preloadFile(p, ctx, file, f, true)
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
	return nil
}

func preloadFile(p *gox.Package, ctx *blockCtx, file string, f *ast.File, gopFile bool) {
	parent := ctx.pkgCtx
	goFile := genGoFile(file, gopFile)
	old, _ := p.SetCurFile(goFile, true)
	defer p.RestoreCurFile(old)
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			dRecv := d.Recv
			if dRecv == nil {
				dRecv = ctx.classRecv // maybe in class file (.spx/.gmx)
			}
			parent.funcs = append(parent.funcs, func() {
				old, _ := p.SetCurFile(goFile, true)
				defer p.RestoreCurFile(old)
				var recv *types.Var
				if dRecv != nil {
					recv = toRecv(ctx, dRecv)
				}
				preloadFunc(ctx, recv, d)
			})
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				for _, item := range d.Specs {
					loadImport(ctx, item.(*ast.ImportSpec))
				}
			case token.TYPE:
				tdecl := ctx.pkg.NewTypeDefs()
				for _, spec := range d.Specs {
					t := spec.(*ast.TypeSpec)
					name, pos := t.Name.Name, t.Pos()
					if t.Assign != token.NoPos { // alias type
						parent.insertObject(name, pos, lazyAliasType, func() {
							if debugLoad {
								log.Println("==> Load > AliasType", name)
							}
							tdecl.AliasType(name, toType(ctx, t.Type), pos)
						})
						continue
					}
					if debugLoad {
						log.Println("==> Load > NewType", name)
					}
					decl := tdecl.NewType(name, pos)
					if t.Doc != nil {
						decl.SetComments(t.Doc)
					} else if d.Doc != nil {
						decl.SetComments(d.Doc)
					}
					parent.insertObject(name, pos, lazyInitType, func() { // decycle
						if debugLoad {
							log.Println("==> Load > InitType", name)
						}
						decl.InitType(ctx.pkg, toType(ctx, t.Type))
					})
				}
			case token.CONST:
				pkg := ctx.pkg
				cdecl := pkg.NewConstDefs(pkg.Types.Scope())
				for iotav, spec := range d.Specs {
					vSpec := spec.(*ast.ValueSpec)
					if debugLoad {
						log.Println("==> Preload const", vSpec.Names)
					}
					parent.insertNames(vSpec.Pos(), vSpec.Names, func() {
						loadConsts(ctx, cdecl, vSpec, iotav)
					})
				}
			case token.VAR:
				pkg := ctx.pkg
				vdecl := pkg.NewVarDefs(pkg.Types.Scope())
				for _, spec := range d.Specs {
					vSpec := spec.(*ast.ValueSpec)
					if debugLoad {
						log.Println("==> Preload var", vSpec.Names)
					}
					parent.insertNames(vSpec.Pos(), vSpec.Names, func() {
						loadVars(ctx, vdecl, vSpec)
					})
				}
			default:
				log.Panicln("TODO - tok:", d.Tok, "spec:", reflect.TypeOf(d.Specs).Elem())
			}
		default:
			log.Panicln("TODO - gopkg.Package.load: unknown decl -", reflect.TypeOf(decl))
		}
	}
	if f.IsProj {
		parent.bodys = append(parent.bodys, func() {
			old, _ := p.SetCurFile(goFile, true)
			defer p.RestoreCurFile(old)
			gmxMainFunc(p, parent)
		})
	}
}

func preloadFunc(ctx *blockCtx, recv *types.Var, d *ast.FuncDecl) {
	name := d.Name.Name
	if debugLoad {
		if recv == nil {
			log.Println("==> New func", name)
		} else {
			log.Printf("==> New method %v.%s\n", recv.Type(), name)
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
	sig := toFuncType(ctx, d.Type, recv, d)
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
		ctx.bodys = append(ctx.bodys, func() { // interface issue: #795
			loadFuncBody(ctx, fn, body, d)
		})
	}
}

func loadFuncBody(ctx *blockCtx, fn *gox.Func, body *ast.BlockStmt, src ast.Node) {
	cb := fn.BodyStart(ctx.pkg)
	compileStmts(ctx, body.List)
	cb.End(src)
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

func loadConstSpecs(ctx *blockCtx, decl *gox.ConstDefs, specs []ast.Spec) {
	for iotav, spec := range specs {
		loadConsts(ctx, decl, spec.(*ast.ValueSpec), iotav)
	}
}

func loadConsts(ctx *blockCtx, decl *gox.ConstDefs, v *ast.ValueSpec, iotav int) {
	names := makeNames(v.Names)
	if v.Values == nil {
		if debugLoad {
			log.Println("==> Load const", names)
		}
		decl.Next(iotav, v.Pos(), names...)
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
	decl.New(fn, iotav, v.Pos(), typ, names...)
}

func loadVarSpecs(ctx *blockCtx, decl *gox.VarDefs, specs []ast.Spec) {
	for _, spec := range specs {
		loadVars(ctx, decl, spec.(*ast.ValueSpec))
	}
}

func loadVars(ctx *blockCtx, decl *gox.VarDefs, v *ast.ValueSpec) {
	var typ types.Type
	if v.Type != nil {
		typ = toType(ctx, v.Type)
	}
	names := makeNames(v.Names)
	if debugLoad {
		log.Println("==> Load var", typ, names)
	}
	varDecl := decl.New(v.Names[0].Pos(), typ, names...)
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
}

func makeNames(vals []*ast.Ident) []string {
	names := make([]string, len(vals))
	for i, v := range vals {
		names[i] = v.Name
	}
	return names
}

// -----------------------------------------------------------------------------
