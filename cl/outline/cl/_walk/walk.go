package walk

import (
	"go/types"
	"reflect"
	"strconv"

	"github.com/goplus/gop/ast"
)

// -----------------------------------------------------------------------------

type Visitor interface {
}

// -----------------------------------------------------------------------------

type Context struct {
	pkg *types.Package
	ctx *types.Scope
	vi  Visitor
}

func New(pkg *types.Package, visitor Visitor) *Context {
	return &Context{pkg: pkg, vi: visitor}
}

func (p *Context) File(f *ast.File) error {
	p.ctx = p.pkg.Scope()
	for _, decl := range f.Decls {
		if e := p.decl(f, decl); e != nil {
			return e
		}
	}
	return nil
}

func (p *Context) decl(f *ast.File, decl ast.Decl) error {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return p.funcDecl(f, d)
	case *ast.GenDecl:
		return p.genDecl(d)
	default:
		panic("Unknown ast.decl: " + reflect.TypeOf(decl).String())
	}
}

func (p *Context) genDecl(f *ast.GenDecl) error {
	return nil
}

func (p *Context) funcDecl(f *ast.File, d *ast.FuncDecl) error {
	body := p.block(d.Body.List)
	if f.ShadowEntry == d {
		return body
	}
	flds := Block{
		"NAME": d.Name.Name,
	}
	p.funcType(flds, d.Type)
	return Block{
		"type":   "goplus_func",
		"fields": flds,
		"inputs": Block{
			"BODY": Block{
				"block": body,
			},
		},
	}
}

func (p *Context) funcType(ret Block, fnt *ast.FuncType) {
	p.fieldList(ret, "ARG", fnt.Params)
	p.fieldList(ret, "RET", fnt.Results)
	// _, variadic := params[n-1].Type.(*ast.Ellipsis)
}

func (p *Context) fieldList(ret Block, namePrefix string, flds *ast.FieldList) {
	if flds == nil {
		return
	}
	i := 0
	for _, param := range flds.List {
		last := len(param.Names) - 1
		for j, name := range param.Names {
			idx := strconv.Itoa(i)
			i++
			arg := namePrefix + idx
			argt := arg + "_TYPE"
			typ := ""
			if j == last {
				typ = p.toType(param.Type)
			}
			ret[arg] = name.Name
			ret[argt] = typ
		}
	}
}

// -----------------------------------------------------------------------------

func (p *Context) block(stmts []ast.Stmt) (ret Block) {
	var prev Block
	for _, stmt := range stmts {
		b := p.stmt(stmt)
		if prev != nil {
			prev.setNext(b)
		} else {
			ret = b
		}
		prev = b
	}
	return
}

func (p *Context) stmt(stmt ast.Stmt) (ret Block) {
	switch v := stmt.(type) {
	case *ast.IfStmt:
		return p.ifStmt(v)
	case *ast.BlockStmt:
		return p.block(v.List)
	}
	return Block{
		"type": "goplus_stmt",
		"inputs": Block{
			"CODE": p.toCode(stmt),
		},
	}
}

func (p *Context) ifStmt(v *ast.IfStmt) (ret Block) {
	cond := p.expr(v.Cond)
	body := p.block(v.Body.List)
	inputs := Block{
		"COND": cond,
		"BODY": body,
	}
	if v.Init != nil {
		inputs["INIT"] = p.stmt(v.Init)
	}
	if v.Else != nil {
		inputs["ELSE"] = p.stmt(v.Else)
	}
	return Block{
		"type":   "goplus_if",
		"inputs": inputs,
	}
}

// -----------------------------------------------------------------------------

func (p *Context) expr(expr ast.Expr) (ret Block) {
	switch v := expr.(type) {
	case *ast.BinaryExpr:
		return p.binaryExpr(v)
	case *ast.UnaryExpr:
		return p.unaryExpr(v)
	}
	return Block{
		"type": "goplus_expr",
		"inputs": Block{
			"CODE": p.toCode(expr),
		},
	}
}

func (p *Context) unaryExpr(v *ast.UnaryExpr) (ret Block) {
	return
}

func (p *Context) binaryExpr(v *ast.BinaryExpr) (ret Block) {
	return
}

// -----------------------------------------------------------------------------

func (p *Context) toType(t ast.Expr) string {
	return p.str(t)
}

func (p *Context) toCode(v ast.Node) (ret Block) {
	return Block{
		"block": Block{
			"type": "text",
			"fields": Block{
				"TEXT": p.str(v),
			},
		},
	}
}

// -----------------------------------------------------------------------------
