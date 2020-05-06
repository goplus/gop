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
	// ErrMainFuncNotFound error.
	ErrMainFuncNotFound = errors.New("main function not found")
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
	pkg     *pkgCtx
	imports map[string]string
}

func newFileCtx(pkg *pkgCtx) *fileCtx {
	return &fileCtx{pkg: pkg, imports: make(map[string]string)}
}

// -----------------------------------------------------------------------------

type iSymbol = interface{}

type blockCtx struct {
	*pkgCtx
	file   *fileCtx
	parent *blockCtx

	vlist []*exec.Var
	syms  map[string]iSymbol
}

func newBlockCtx(file *fileCtx, parent *blockCtx) *blockCtx {
	return &blockCtx{
		pkgCtx: file.pkg,
		file:   file,
		parent: parent,
		syms:   make(map[string]iSymbol),
	}
}

func (p *blockCtx) findVar(name string) (addr *exec.Var, err error) {
	for ; p != nil; p = p.parent {
		if v, ok := p.syms[name]; ok {
			if addr, ok = v.(*exec.Var); ok {
				return
			}
			return nil, syscall.EEXIST // exists, but not a variable
		}
	}
	return nil, syscall.ENOENT // not found
}

func (p *blockCtx) insertVar(name string, typ reflect.Type) *exec.Var {
	if _, ok := p.syms[name]; ok {
		log.Panicln("insertVar failed: symbol exists -", name)
	}
	idx := uint32(len(p.vlist))
	v := exec.NewVar(typ, name)
	v.SetAddr(p.out.NestDepth, idx)
	p.syms[name] = v
	p.vlist = append(p.vlist, v)
	return v
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
	ctx     *fileCtx
}

type typeDecl struct {
	Methods map[string]*methodDecl
	Alias   bool
}

type funcDecl struct {
	Type *funcTypeDecl
	Body *ast.BlockStmt
	ctx  *fileCtx
}

// A Package represents a qlang package.
type Package struct {
	types map[string]*typeDecl
	funcs map[string]*funcDecl
	vlist []*exec.Var
}

// NewPackage creates a qlang package instance.
func NewPackage(out *exec.Builder, pkg *ast.Package) (p *Package, err error) {
	p = &Package{
		types: make(map[string]*typeDecl),
		funcs: make(map[string]*funcDecl),
	}
	ctx := newPkgCtx(out)
	for _, f := range pkg.Files {
		p.loadFile(ctx, f)
	}
	if pkg.Name == "main" {
		entry, ok := p.funcs["main"]
		if !ok {
			return p, ErrMainFuncNotFound
		}
		block := newBlockCtx(entry.ctx, nil)
		p.compileBlockStmt(block, entry.Body)
		p.vlist = block.vlist
	}
	return
}

// GetGlobalVars returns the global variable list.
func (p *Package) GetGlobalVars() []*exec.Var {
	return p.vlist
}

func (p *Package) loadFile(pkg *pkgCtx, f *ast.File) {
	ctx := newFileCtx(pkg)
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
	log.Debug("import:", name, pkgPath)
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
		p.insertMethod(recv.Type, name, &methodDecl{
			Recv:    recv.Name,
			Pointer: recv.Pointer,
			Type:    &funcTypeDecl{X: d.Type},
			Body:    d.Body,
			ctx:     ctx,
		})
	} else if name == "init" {
		log.Panicln("loadFunc TODO: init")
	} else {
		p.insertFunc(name, &funcDecl{
			Type: &funcTypeDecl{X: d.Type},
			Body: d.Body,
			ctx:  ctx,
		})
	}
}

func (p *Package) insertFunc(name string, fun *funcDecl) {
	if _, ok := p.funcs[name]; ok {
		log.Panicln("insertFunc failed: func exists -", name)
	}
	p.funcs[name] = fun
}

func (p *Package) insertMethod(typeName, methodName string, method *methodDecl) {
	typ, ok := p.types[typeName]
	if !ok {
		typ = new(typeDecl)
		p.types[typeName] = typ
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
