package golang

import (
	"go/ast"
	"go/token"
	"reflect"

	"github.com/qiniu/qlang/v6/exec.spec"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// Var represents a variable.
type Var struct {
	typ   reflect.Type
	name  string
	where *varManager
}

// NewVar creates a variable instance.
func NewVar(typ reflect.Type, name string) *Var {
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
	c := p.name[0]
	return c >= '0' && c <= '9'
}

func (p *Var) setScope(where *varManager) {
	if p.where != nil {
		panic("Var.setScope: variable already defined")
	}
	p.where = where
}

// -----------------------------------------------------------------------------

type varManager struct {
	vlist []exec.Var
}

func newVarManager(vars ...exec.Var) *varManager {
	return &varManager{vlist: vars}
}

func (p *varManager) addVar(vars ...exec.Var) {
	for _, v := range vars {
		v.(*Var).setScope(p)
	}
	p.vlist = append(p.vlist, vars...)
}

func (p *varManager) toGenDecl(b *Builder) *ast.GenDecl {
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

// -----------------------------------------------------------------------------

// DefineVar defines variables.
func (p *Builder) DefineVar(vars ...exec.Var) *Builder {
	p.addVar(vars...)
	return p
}

// InCurrentCtx returns if a variable is in current context or not.
func (p *Builder) InCurrentCtx(v exec.Var) bool {
	return p.varManager == v.(*Var).where
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

// AddrOp instr
func (p *Builder) AddrOp(kind exec.Kind, op exec.AddrOperator) *Builder {
	var stmt ast.Stmt
	var x = p.rhs.Pop()
	var val = p.rhs.Pop().(ast.Expr)
	switch v := x.(type) {
	case *ast.UnaryExpr:
		stmt = &ast.AssignStmt{
			Lhs: []ast.Expr{v.X}, Tok: addropTokens[op], Rhs: []ast.Expr{val},
		}
	default:
		log.Panicln("AddrOp: todo")
	}
	p.rhs.Push(stmt)
	return p
}

var addropTokens = [...]token.Token{
	exec.OpAddAssign:       token.ADD_ASSIGN,
	exec.OpSubAssign:       token.SUB_ASSIGN,
	exec.OpMulAssign:       token.MUL_ASSIGN,
	exec.OpDivAssign:       token.QUO_ASSIGN,
	exec.OpModAssign:       token.REM_ASSIGN,
	exec.OpBitAndAssign:    token.AND_ASSIGN,
	exec.OpBitOrAssign:     token.OR_ASSIGN,
	exec.OpBitXorAssign:    token.XOR_ASSIGN,
	exec.OpBitAndNotAssign: token.AND_NOT_ASSIGN,
	exec.OpBitSHLAssign:    token.SHL_ASSIGN,
	exec.OpBitSHRAssign:    token.SHR_ASSIGN,
	exec.OpInc:             token.INC,
	exec.OpDec:             token.DEC,
}

// -----------------------------------------------------------------------------
