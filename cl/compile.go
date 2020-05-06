package cl

import (
	"errors"
	"path"
	"reflect"
	"syscall"

	"github.com/qiniu/qlang/ast"
	"github.com/qiniu/qlang/ast/astutil"
	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/qlang/token"
	"github.com/qiniu/x/log"
)

var (
	// ErrNotFound error.
	ErrNotFound = syscall.ENOENT

	// ErrMainFuncNotFound error.
	ErrMainFuncNotFound = errors.New("main function not found")

	// ErrSymbolNotVariable error.
	ErrSymbolNotVariable = errors.New("symbol exists but not a variable")

	// ErrSymbolNotFunc error.
	ErrSymbolNotFunc = errors.New("symbol exists but not a func")

	// ErrSymbolNotType error.
	ErrSymbolNotType = errors.New("symbol exists but not a type")
)

// -----------------------------------------------------------------------------

type pkgCtx struct {
	infer   exec.Stack
	builtin *exec.GoPackage
	out     *exec.Builder
}

func newPkgCtx(out *exec.Builder) *pkgCtx {
	p := &pkgCtx{builtin: exec.Package(""), out: out}
	p.infer.Init()
	return p
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
// - pkgName => pkgPath
// - funcName => *funcDecl
// - typeName => *typeDecl
//
type iGblSymbol = interface{}

type blockCtx struct {
	*pkgCtx
	file   *fileCtx
	parent *blockCtx

	vlist []*exec.Var
	syms  map[string]iGblSymbol
}

func newGblBlockCtx(pkg *pkgCtx, parent *blockCtx) *blockCtx {
	return &blockCtx{
		pkgCtx: pkg,
		parent: parent,
		syms:   make(map[string]iGblSymbol),
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
	if ctx.parent == nil { // it's global blockCtx
		sym, ok = ctx.file.imports[name]
	}
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

func (p *blockCtx) findVar(name string) (addr *exec.Var, err error) {
	v, ok := p.find(name)
	if !ok {
		return nil, ErrNotFound
	}
	if addr, ok = v.(*exec.Var); ok {
		return
	}
	return nil, ErrSymbolNotVariable
}

func (p *blockCtx) insertVar(name string, typ reflect.Type) *exec.Var {
	if p.exists(name) {
		log.Panicln("insertVar failed: symbol exists -", name)
	}
	idx := uint32(len(p.vlist))
	v := exec.NewVar(typ, name)
	v.SetAddr(p.out.NestDepth, idx)
	p.syms[name] = v
	p.vlist = append(p.vlist, v)
	return v
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

type funcTypeDecl struct {
	X *ast.FuncType
}

type methodDecl struct {
	Recv    string // recv object name
	Pointer int
	Type    *funcTypeDecl
	Body    *ast.BlockStmt
	file    *fileCtx
}

type typeDecl struct {
	Methods map[string]*methodDecl
	Alias   bool
}

type funcDecl struct {
	Type *funcTypeDecl
	Body *ast.BlockStmt
	file *fileCtx
}

// A Package represents a qlang package.
type Package struct {
	vlist []*exec.Var
	syms  map[string]iGblSymbol
}

// NewPackage creates a qlang package instance.
func NewPackage(out *exec.Builder, pkg *ast.Package) (p *Package, err error) {
	p = &Package{}
	ctxPkg := newPkgCtx(out)
	ctx := newGblBlockCtx(ctxPkg, nil)
	for _, f := range pkg.Files {
		p.loadFile(ctx, f)
	}
	if pkg.Name == "main" {
		entry, err := ctx.findFunc("main")
		if err != nil {
			if err == ErrNotFound {
				err = ErrMainFuncNotFound
			}
			return p, err
		}
		ctx.file = entry.file
		p.compileBlockStmt(ctx, entry.Body)
	}
	p.vlist = ctx.vlist
	p.syms = ctx.syms
	return
}

// GetGlobalVars returns the global variable list.
func (p *Package) GetGlobalVars() []*exec.Var {
	return p.vlist
}

func (p *Package) loadFile(block *blockCtx, f *ast.File) {
	ctx := newFileCtx(block)
	block.file = ctx
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			p.loadFunc(ctx, d)
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				p.loadImports(ctx, d)
			case token.TYPE:
				p.loadTypes(d)
			case token.CONST:
				p.loadConsts(d)
			case token.VAR:
				p.loadVars(d)
			default:
				log.Panicln("tok:", d.Tok, "spec:", reflect.TypeOf(d.Specs).Elem())
			}
		default:
			log.Panicln("gopkg.Package.load: unknown decl -", reflect.TypeOf(decl))
		}
	}
}

func (p *Package) loadImports(ctx *fileCtx, d *ast.GenDecl) {
	for _, item := range d.Specs {
		p.loadImport(ctx, item.(*ast.ImportSpec))
	}
}

func (p *Package) loadImport(ctx *fileCtx, spec *ast.ImportSpec) {
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

func (p *Package) loadTypes(d *ast.GenDecl) {
	for _, item := range d.Specs {
		p.loadType(item.(*ast.TypeSpec))
	}
}

func (p *Package) loadType(spec *ast.TypeSpec) {
}

func (p *Package) loadConsts(d *ast.GenDecl) {
}

func (p *Package) loadVars(d *ast.GenDecl) {
	for _, item := range d.Specs {
		p.loadVar(item.(*ast.ValueSpec))
	}
}

func (p *Package) loadVar(spec *ast.ValueSpec) {
}

func (p *Package) loadFunc(ctx *fileCtx, d *ast.FuncDecl) {
	var name = d.Name.Name
	if d.Recv != nil {
		recv := astutil.ToRecv(d.Recv)
		ctx.insertMethod(recv.Type, name, &methodDecl{
			Recv:    recv.Name,
			Pointer: recv.Pointer,
			Type:    &funcTypeDecl{X: d.Type},
			Body:    d.Body,
			file:    ctx,
		})
	} else if name == "init" {
		log.Panicln("loadFunc TODO: init")
	} else {
		ctx.insertFunc(name, &funcDecl{
			Type: &funcTypeDecl{X: d.Type},
			Body: d.Body,
			file: ctx,
		})
	}
}

// -----------------------------------------------------------------------------
