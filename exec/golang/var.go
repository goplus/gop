/*
 Copyright 2020 The GoPlus Authors (goplus.org)

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

package golang

import (
	"go/ast"
	"go/token"
	"reflect"
	"strings"

	"github.com/goplus/gop/exec.spec"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// Var represents a variable.
type Var struct {
	typ   reflect.Type
	name  string
	where *scopeCtx
}

// NewVar creates a variable instance.
func NewVar(typ reflect.Type, name string) *Var {
	c := name[0]
	if c >= '0' && c <= '9' {
		name = "_ret_" + name
	}
	return &Var{typ: typ, name: name}
}

// Type returns variable's type.
func (p *Var) Type() reflect.Type {
	return p.typ
}

// Name returns variable's name.
func (p *Var) Name() string {
	return p.name
}

// IsUnnamedOut returns if variable unnamed or not.
func (p *Var) IsUnnamedOut() bool {
	return strings.HasPrefix(p.name, "_ret_")
}

func (p *Var) setScope(where *scopeCtx) {
	if p.where != nil {
		panic("Var.setScope: variable already defined")
	}
	p.where = where
}

// -----------------------------------------------------------------------------

type scopeCtx struct {
	parentCtx *scopeCtx
	vlist     []exec.Var
	stmts     []ast.Stmt
	labels    []*Label // labels of current statement
}

func (p *scopeCtx) addVar(vars ...exec.Var) {
	for _, v := range vars {
		v.(*Var).setScope(p)
	}
	p.vlist = append(p.vlist, vars...)
}

func (p *scopeCtx) toGenDecl(b *Builder) *ast.GenDecl {
	n := len(p.vlist)
	if n == 0 {
		return nil
	}
	specs := make([]ast.Spec, 0, n)
	for _, item := range p.vlist {
		v := item.(*Var)
		spec := &ast.ValueSpec{
			Names: []*ast.Ident{Ident(v.name)},
			Type:  Type(b, v.typ),
		}
		specs = append(specs, spec)
	}
	return &ast.GenDecl{
		Tok:   token.VAR,
		Specs: specs,
	}
}

func (p *scopeCtx) getStmts(b *Builder) []ast.Stmt {
	if decl := p.toGenDecl(b); decl != nil {
		p.stmts[0] = &ast.DeclStmt{Decl: decl}
		return p.stmts
	}
	return p.stmts[1:]
}

func (p *scopeCtx) initStmts() {
	p.stmts = make([]ast.Stmt, 1, 8)
}

// -----------------------------------------------------------------------------

// DefineVar defines variables.
func (p *Builder) DefineVar(vars ...exec.Var) *Builder {
	var vlist []exec.Var
	for _, v := range vars {
		if pkgPath := v.Type().PkgPath(); pkgPath != "" && pkgPath != p.pkgName {
			p.Import(pkgPath)
		} else if v.Name() == "_" {
			continue
		}
		pv := v.(*Var)
		if strings.HasPrefix(pv.name, "_") {
			pv.name = "_q" + pv.name
		}
		vlist = append(vlist, v)
	}
	p.addVar(vlist...)
	return p
}

// InCurrentCtx returns if a variable is in current context or not.
func (p *Builder) InCurrentCtx(v exec.Var) bool {
	return p.scopeCtx == v.(*Var).where
}

// Load instr
func (p *Builder) Load(idx int32) *Builder {
	p.rhs.Push(p.argIdent(idx))
	return p
}

// Addr instr
func (p *Builder) Addr(idx int32) *Builder {
	p.rhs.Push(p.argIdent(idx))
	return p
}

// Store instr
func (p *Builder) Store(idx int32) *Builder {
	p.lhs.Push(p.argIdent(idx))
	return p
}

func (p *Builder) argIdent(idx int32) *ast.Ident {
	i := len(p.cfun.in) + int(idx)
	if i == -1 {
		return Ident("_recv")
	}
	return Ident(toArg(i))
}

// LoadVar instr
func (p *Builder) LoadVar(v exec.Var) *Builder {
	p.rhs.Push(Ident(v.(*Var).name))
	return p
}

// StoreVar instr
func (p *Builder) StoreVar(v exec.Var) *Builder {
	p.lhs.Push(Ident(v.(*Var).name))
	return p
}

// AddrVar instr
func (p *Builder) AddrVar(v exec.Var) *Builder {
	p.rhs.Push(&ast.UnaryExpr{
		Op: token.AND,
		X:  Ident(v.(*Var).name),
	})
	return p
}

func (p *Builder) bigAddrOp(kind exec.Kind, op exec.AddrOperator) *Builder {
	if op == exec.OpAddrVal {
		return p
	}
	method := addropMethods[op]
	if method == "" {
		log.Panicln("bigAddrOp: unknown op -", op)
	}
	var expr ast.Expr
	var x = p.rhs.Pop()
	var val = p.rhs.Pop().(ast.Expr)
	if op == exec.OpInc || op == exec.OpDec {
		switch kind {
		case exec.BigInt:
			val = &ast.CallExpr{
				Fun:  p.GoSymIdent("math/big", "NewInt"),
				Args: []ast.Expr{&ast.Ident{Name: "1"}},
			}
		case exec.BigRat:
			val = &ast.CallExpr{
				Fun:  p.GoSymIdent("math/big", "NewRat"),
				Args: []ast.Expr{&ast.Ident{Name: "1"}, &ast.Ident{Name: "1"}},
			}
		case exec.BigFloat:
			val = &ast.CallExpr{
				Fun:  p.GoSymIdent("math/big", "NewFloat"),
				Args: []ast.Expr{&ast.Ident{Name: "1"}},
			}
		}
	}
	switch v := x.(type) {
	case *ast.UnaryExpr:
		if v.Op != token.AND {
			log.Panicln("bigAddrOp: unknown x expr -", reflect.TypeOf(x))
		}
		bigOp := &ast.SelectorExpr{X: v.X, Sel: Ident(method)}
		expr = &ast.CallExpr{Fun: bigOp, Args: []ast.Expr{v.X, val}}
	case *ast.Ident:
		bigOp := &ast.SelectorExpr{X: v, Sel: Ident(method)}
		expr = &ast.CallExpr{Fun: bigOp, Args: []ast.Expr{v, val}}
	default:
		log.Panicln("bigAddrOp: todo")
	}
	p.rhs.Push(&ast.ExprStmt{X: expr})
	return p
}

var addropMethods = [...]string{
	exec.OpAddAssign:    "Add",
	exec.OpSubAssign:    "Sub",
	exec.OpMulAssign:    "Mul",
	exec.OpQuoAssign:    "Quo",
	exec.OpModAssign:    "Mod",
	exec.OpAndAssign:    "And",
	exec.OpOrAssign:     "Or",
	exec.OpXorAssign:    "Xor",
	exec.OpAndNotAssign: "AndNot",
	exec.OpLshAssign:    "Lsh",
	exec.OpRshAssign:    "Rsh",
	exec.OpInc:          "Add",
	exec.OpDec:          "Sub",
}

// AddrOp instr
func (p *Builder) AddrOp(kind exec.Kind, op exec.AddrOperator) *Builder {
	if kind >= exec.BigInt {
		return p.bigAddrOp(kind, op)
	}
	if op == exec.OpAddrVal {
		p.rhs.Push(&ast.StarExpr{
			X: p.rhs.Pop().(ast.Expr),
		})
		return p
	}
	if op == exec.OpAssign {
		p.emitStmt(&ast.AssignStmt{
			Lhs: []ast.Expr{&ast.StarExpr{
				X: p.rhs.Pop().(ast.Expr),
			}},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{p.rhs.Pop().(ast.Expr)},
		})
		return p
	}
	var stmt ast.Stmt
	var x = p.rhs.Pop()
	var val = p.rhs.Pop().(ast.Expr)
	switch v := x.(type) {
	case *ast.UnaryExpr:
		if v.Op != token.AND {
			log.Panicln("AddrOp: unknown x expr -", reflect.TypeOf(x))
		}
		if op == exec.OpInc || op == exec.OpDec {
			stmt = &ast.IncDecStmt{X: v.X, TokPos: v.OpPos, Tok: addropTokens[op]}
		} else {
			stmt = &ast.AssignStmt{Lhs: []ast.Expr{v.X}, Tok: addropTokens[op], Rhs: []ast.Expr{val}}
		}
	default:
		log.Panicln("AddrOp: todo")
	}
	p.rhs.Push(stmt)
	return p
}

var addropTokens = [...]token.Token{
	exec.OpAddAssign:    token.ADD_ASSIGN,
	exec.OpSubAssign:    token.SUB_ASSIGN,
	exec.OpMulAssign:    token.MUL_ASSIGN,
	exec.OpQuoAssign:    token.QUO_ASSIGN,
	exec.OpModAssign:    token.REM_ASSIGN,
	exec.OpAndAssign:    token.AND_ASSIGN,
	exec.OpOrAssign:     token.OR_ASSIGN,
	exec.OpXorAssign:    token.XOR_ASSIGN,
	exec.OpAndNotAssign: token.AND_NOT_ASSIGN,
	exec.OpLshAssign:    token.SHL_ASSIGN,
	exec.OpRshAssign:    token.SHR_ASSIGN,
	exec.OpInc:          token.INC,
	exec.OpDec:          token.DEC,
}

// -----------------------------------------------------------------------------
