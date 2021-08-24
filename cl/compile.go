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
	"go/constant"
	"go/types"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

const (
	gopPrefix = "Gop_"
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
	// If WorkingDir or TargetDir is empty, it is same as Dir.
	WorkingDir, TargetDir string

	// ModPath specifies module path.
	ModPath string

	// ModRootDir specifies root dir of this module.
	// If ModRootDir is empty, will lookup go.mod in all ancestor directories of Dir.
	// If you specify ModPath, you should specify ModRootDir at the same time.
	ModRootDir string

	// CacheFile specifies where cache data to write.
	CacheFile string

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

	// CacheLoadPkgs = true means to cache all loaded packages.
	CacheLoadPkgs bool

	// PersistLoadPkgs = true means to cache all loaded packages to disk.
	PersistLoadPkgs bool

	// NoFileLine = true means not to generate file line comments.
	NoFileLine bool

	// RelativePath = true means to generate file line comments with relative file path.
	RelativePath bool
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

type gmxSettings struct {
	Pkg    string
	Class  string
	This   string
	spx    *gox.PkgRef
	game   gox.Ref
	sprite gox.Ref
	titles []string
	params map[string]*types.Tuple
}

func newGmx(pkg *gox.Package, gmx *ast.File) *gmxSettings {
	p := &gmxSettings{Pkg: "github.com/goplus/spx", Class: "Game", This: "_gop_this"}
	getGameSettings(gmx, p)
	p.spx = pkg.Import(p.Pkg)
	p.game = spxRef(p.spx, "Gop_game", "Game")
	p.sprite = spxRef(p.spx, "Gop_sprite", "Sprite")
	p.titles = getStrings(p.spx, "Gop_title")
	p.params = getParams(p.spx, "Gop_params")
	return p
}

func spxRef(spx *gox.PkgRef, name, typ string) gox.Ref {
	if v := getStringConst(spx, name); v != "" {
		typ = v
	}
	return spx.Ref(typ)
}

// onMsg(, _gop_data interface{}); onCloned(_gop_data interface{})
func getParams(spx *gox.PkgRef, name string) map[string]*types.Tuple {
	ret := map[string]*types.Tuple{}
	if v := getStringConst(spx, name); v != "" {
		for _, part := range strings.Split(v, ";") {
			part = strings.TrimPrefix(part, " ")
			pos := strings.Index(part, "(")
			pos2 := strings.Index(part, ")")
			if pos > 0 && pos2 == len(part)-1 {
				key := part[:pos]
				params := strings.Split(part[pos+1:pos2], ",")
				tuple := make([]*types.Var, len(params))
				for i, param := range params {
					param = strings.TrimPrefix(param, " ")
					tuple[i] = getVar(param)
				}
				ret[key] = types.NewTuple(tuple...)
			}
		}
	}
	return ret
}

// _gop_data interface{}
func getVar(param string) *types.Var {
	if param == "" {
		return nil
	}
	idx := strings.Index(param, " ")
	if idx > 0 {
		name := param[:idx]
		if typ, ok := strTypes[param[idx+1:]]; ok {
			return types.NewParam(token.NoPos, nil, name, typ)
		}
	}
	panic("getVar failed: " + param)
}

var (
	strTypes = map[string]types.Type{
		"interface{}": gox.TyEmptyInterface,
		"string":      types.Typ[types.String],
		"int":         types.Typ[types.Int],
	}
)

// on,off
func getStrings(spx *gox.PkgRef, name string) []string {
	if v := getStringConst(spx, name); v != "" {
		return strings.Split(v, ",")
	}
	return nil
}

func getStringConst(spx *gox.PkgRef, name string) string {
	if o := spx.Ref(name); o != nil {
		if c, ok := o.(*types.Const); ok {
			return constant.StringVal(c.Val())
		}
	}
	return ""
}

func getStringVal(v ast.Expr, ret *string) {
	if lit, ok := v.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		if val, err := strconv.Unquote(lit.Value); err == nil {
			*ret = val
		}
	}
}

func getGameSettings(gmx *ast.File, ret *gmxSettings) {
	decls := gmx.Decls
	i, n := 0, len(decls)
	for i < n {
		g, ok := decls[i].(*ast.GenDecl)
		if !ok {
			break
		}
		if g.Tok != token.IMPORT {
			if g.Tok == token.CONST {
				rm := true
				for _, v := range g.Specs {
					spec := v.(*ast.ValueSpec)
					for i, name := range spec.Names {
						switch name.Name {
						case "GopGamePkg":
							getStringVal(spec.Values[i], &ret.Pkg)
						case "GopClass":
							getStringVal(spec.Values[i], &ret.Class)
						case "GopThis":
							getStringVal(spec.Values[i], &ret.This)
						default:
							rm = false
						}
					}
				}
				if rm {
					g.Specs = nil
				}
			}
			break
		}
		i++ // skip import, if any
	}
}

type pkgCtx struct {
	*nodeInterp
	*gmxSettings
	syms  map[string]loader
	inits []func()
	errs  []error
}

type blockCtx struct {
	*pkgCtx
	pkg          *gox.Package
	cb           *gox.CodeBuilder
	fset         *token.FileSet
	imports      map[string]*gox.PkgRef
	targetDir    string
	fileLine     bool
	relativePath bool
	fileType     int16
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
	if enableRecover {
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(error); ok {
					p.handleErr(err)
				} else {
					panic(e)
				}
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
		ModPath:         conf.ModPath,
		Env:             conf.Env,
		BuildFlags:      conf.BuildFlags,
		Fset:            conf.Fset,
		LoadPkgs:        conf.PkgsLoader.LoadPkgs,
		LoadNamed:       ctx.loadNamed,
		HandleErr:       ctx.handleErr,
		NodeInterpreter: interp,
		Prefix:          gopPrefix,
		ParseFile:       nil, // TODO
		NewBuiltin:      newBuiltinDefault,
	}
	p = gox.NewPackage(pkgPath, pkg.Name, confGox)
	for _, gmx := range pkg.Files {
		if gmx.FileType == ast.FileTypeGmx {
			ctx.gmxSettings = newGmx(p, gmx)
			break
		}
	}
	for fpath, f := range pkg.Files {
		testingFile := strings.HasSuffix(fpath, "_test.gop")
		preloadFile(p, ctx, f, targetDir, testingFile, conf)
	}
	for _, f := range pkg.Files {
		if f.FileType == ast.FileTypeGmx {
			loadFile(ctx, f)
			break
		}
	}
	for _, f := range pkg.Files {
		if f.FileType != ast.FileTypeGmx { // only one .gmx file
			loadFile(ctx, f)
		}
	}
	for _, load := range ctx.inits {
		load()
	}
	err = ctx.complete()
	return
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

func preloadFile(p *gox.Package, parent *pkgCtx, f *ast.File, targetDir string, testingFile bool, conf *Config) {
	syms := parent.syms
	fileLine := !conf.NoFileLine
	ctx := &blockCtx{
		pkg: p, pkgCtx: parent, cb: p.CB(), fset: p.Fset, targetDir: targetDir, fileType: f.FileType,
		fileLine: fileLine, relativePath: conf.RelativePath, imports: make(map[string]*gox.PkgRef),
	}
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			switch f.FileType {
			// TODO:
			}
			if d.Recv == nil {
				name := d.Name
				fn := func() {
					old := p.SetInTestingFile(testingFile)
					defer p.SetInTestingFile(old)
					loadFunc(ctx, nil, d)
				}
				if name.Name == "init" {
					if debugLoad {
						log.Println("==> Preload func init")
					}
					parent.inits = append(parent.inits, fn)
				} else {
					if debugLoad {
						log.Println("==> Preload func", name.Name)
					}
					initLoader(parent, syms, name.Pos(), name.Name, fn)
				}
			} else {
				if name, ok := getRecvTypeName(parent, d.Recv, true); ok {
					if debugLoad {
						log.Printf("==> Preload method %s.%s\n", name, d.Name.Name)
					}
					ld := getTypeLoader(syms, token.NoPos, name)
					ld.methods = append(ld.methods, func() {
						old := p.SetInTestingFile(testingFile)
						defer p.SetInTestingFile(old)
						doInitType(ld)
						recv := toRecv(ctx, d.Recv)
						loadFunc(ctx, recv, d)
					})
				}
			}
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				p.SetInTestingFile(testingFile)
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
					ld := getTypeLoader(syms, t.Name.Pos(), name)
					ld.typ = func() {
						old := p.SetInTestingFile(testingFile)
						defer p.SetInTestingFile(old)
						if t.Assign != token.NoPos { // alias type
							if debugLoad {
								log.Println("==> Load > AliasType", name)
							}
							ctx.pkg.AliasType(name, toType(ctx, t.Type))
							return
						}
						if debugLoad {
							log.Println("==> Load > NewType", name)
						}
						decl := ctx.pkg.NewType(name)
						ld.typInit = func() { // decycle
							if debugLoad {
								log.Println("==> Load > InitType", name)
							}
							decl.InitType(ctx.pkg, toType(ctx, t.Type))
						}
					}
				}
			case token.CONST:
				for _, spec := range d.Specs {
					vSpec := spec.(*ast.ValueSpec)
					if debugLoad {
						log.Println("==> Preload const", vSpec.Names)
					}
					setNamesLoader(parent, syms, vSpec.Names, func() {
						if v := vSpec; v != nil { // only init once
							old := p.SetInTestingFile(testingFile)
							defer p.SetInTestingFile(old)
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
					if debugLoad {
						log.Println("==> Preload var", vSpec.Names)
					}
					setNamesLoader(parent, syms, vSpec.Names, func() {
						if v := vSpec; v != nil { // only init once
							old := p.SetInTestingFile(testingFile)
							defer p.SetInTestingFile(old)
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
	if body := d.Body; body != nil {
		loadFuncBody(ctx, fn, body)
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
	"=":   "Gop_Assign",

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
	"^":  "Gop_Not",
	"!":  "Gop_LNot",
	"<-": "Gop_Recv",
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
		if name == "." {
			panic("TODO: not impl")
		}
		if name == "_" {
			pkg.MarkForceUsed()
			return
		}
	} else {
		name = path.Base(pkgPath)
	}
	ctx.imports[name] = pkg
}

func loadConsts(ctx *blockCtx, names []string, v *ast.ValueSpec) {
	var typ types.Type
	if v.Type != nil {
		typ = toType(ctx, v.Type)
	}
	if debugLoad {
		log.Println("==> Load const", typ, names)
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
	if debugLoad {
		log.Println("==> Load var", typ, names)
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
