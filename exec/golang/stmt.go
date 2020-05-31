package golang

import (
	"go/ast"
	"go/token"
	"reflect"
)

// ----------------------------------------------------------------------------

// Label represents a label.
type Label struct {
	name string
}

// NewLabel creates a label object.
func NewLabel(name string) *Label {
	return &Label{name: name}
}

func (p *Label) getName(b *Builder) string {
	if p.name == "" {
		p.name = b.autoIdent()
	}
	return p.name
}

// Name returns the label name.
func (p *Label) Name() string {
	return p.name
}

// Label defines a label to jmp here.
func (p *Builder) Label(l *Label) *Builder {
	p.labels = append(p.labels, l)
	return p
}

// Jmp instr
func (p *Builder) Jmp(l *Label) *Builder {
	p.emitStmt(Goto(p, l))
	return p
}

// Goto instr
func Goto(p *Builder, l *Label) *ast.BranchStmt {
	return &ast.BranchStmt{
		Tok:   token.GOTO,
		Label: Ident(l.getName(p)),
	}
}

// GotoIf instr
func GotoIf(p *Builder, cond ast.Expr, l *Label) *ast.IfStmt {
	return &ast.IfStmt{
		Cond: cond,
		Body: &ast.BlockStmt{
			List: []ast.Stmt{Goto(p, l)},
		},
	}
}

// JmpIf instr
func (p *Builder) JmpIf(zeroOrOne uint32, l *Label) *Builder {
	cond := p.rhs.Pop().(ast.Expr)
	if zeroOrOne == 0 {
		cond = &ast.UnaryExpr{Op: token.NOT, X: cond}
	}
	p.emitStmt(GotoIf(p, cond, l))
	return p
}

// CaseNE instr
func (p *Builder) CaseNE(l *Label, arity int) *Builder {
	args := p.rhs.GetArgs(arity + 1)
	arg0 := args[0]
	x, ok := arg0.(*ast.Ident)
	if !ok {
		x = Ident(p.autoIdent())
		p.emitStmt(&ast.AssignStmt{
			Lhs: []ast.Expr{x},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{arg0.(ast.Expr)},
		})
		args[0] = x
	}
	for i := 1; i <= arity; i++ {
		p.emitStmt(GotoIf(p, &ast.BinaryExpr{X: x, Op: token.NEQ, Y: args[i].(ast.Expr)}, l))
	}
	p.rhs.PopN(arity)
	return p
}

// Default instr
func (p *Builder) Default() *Builder {
	p.rhs.PopN(1)
	return p
}

// ----------------------------------------------------------------------------

type blockCtx struct {
	parent *scopeCtx
}

func (p *blockCtx) saveEnv(b *Builder) {
	if p.parent != nil {
		panic("blockCtx.setParent: already defined")
	}
	p.parent = b.scopeCtx
}

func (p *blockCtx) restoreEnv(b *Builder) {
	b.scopeCtx = p.parent
}

// ForPhrase represents a for range phrase.
type ForPhrase struct {
	Key, Value *Var // Key, Value may be nil
	X, Cond    ast.Expr
	TypeIn     reflect.Type
	scopeCtx
	blockCtx
}

// NewForPhrase creates a new ForPhrase instance.
func NewForPhrase(in reflect.Type) *ForPhrase {
	return &ForPhrase{TypeIn: in}
}

// Comprehension represents a list/map comprehension.
type Comprehension struct {
	TypeOut reflect.Type
	old     func()
}

// NewComprehension creates a new Comprehension instance.
func NewComprehension(out reflect.Type) *Comprehension {
	return &Comprehension{TypeOut: out}
}

// ForPhrase instr
func (p *Builder) ForPhrase(f *ForPhrase, key, val *Var, hasExecCtx ...bool) *Builder {
	f.initStmts()
	f.saveEnv(p)
	p.scopeCtx = &f.scopeCtx
	f.Key, f.Value = key, val
	f.X = p.rhs.Pop().(ast.Expr)
	return p
}

// FilterForPhrase instr
func (p *Builder) FilterForPhrase(f *ForPhrase) *Builder {
	f.Cond = p.rhs.Pop().(ast.Expr)
	return p
}

// EndForPhrase instr
func (p *Builder) EndForPhrase(f *ForPhrase) *Builder {
	if p.comprehens != nil {
		p.comprehens()
	}
	body := &ast.BlockStmt{List: f.getStmts(p)}
	if f.Cond != nil {
		body = &ast.BlockStmt{List: []ast.Stmt{
			&ast.IfStmt{Cond: f.Cond, Body: body},
		}}
	}
	p.rhs.Push(&ast.RangeStmt{
		Key:   toVarExpr(f.Key, unnamedVar),
		Value: toVarExpr(f.Value, nil),
		Tok:   token.DEFINE,
		X:     f.X,
		Body:  body,
	})
	f.restoreEnv(p)
	return p
}

func toVarExpr(v *Var, defval ast.Expr) ast.Expr {
	if v == nil {
		return defval
	}
	return Ident(v.name)
}

// ListComprehension instr
func (p *Builder) ListComprehension(c *Comprehension) *Builder {
	c.old = p.comprehens
	p.comprehens = func() {
		v := p.rhs.Pop()
		if stmt, ok := v.(*ast.RangeStmt); ok {
			p.emitStmt(stmt)
			return
		}
		x := v.(ast.Expr)
		appendExpr := &ast.CallExpr{
			Fun:  appendIden,
			Args: []ast.Expr{qlangRet, x},
		}
		assign := &ast.AssignStmt{
			Lhs: []ast.Expr{qlangRet},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{appendExpr},
		}
		p.emitStmt(assign)
	}
	return p
}

// MapComprehension instr
func (p *Builder) MapComprehension(c *Comprehension) *Builder {
	c.old = p.comprehens
	p.comprehens = func() {
		v := p.rhs.Pop()
		if stmt, ok := v.(*ast.RangeStmt); ok {
			p.emitStmt(stmt)
			return
		}
		val := v.(ast.Expr)
		key := p.rhs.Pop().(ast.Expr)
		assign := &ast.AssignStmt{
			Lhs: []ast.Expr{&ast.IndexExpr{X: qlangRet, Index: key}},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{val},
		}
		p.emitStmt(assign)
	}
	return p
}

// EndComprehension instr
func (p *Builder) EndComprehension(c *Comprehension) *Builder {
	typOut := Type(p, c.TypeOut)
	fldOut := &ast.Field{
		Names: []*ast.Ident{qlangRet},
		Type:  typOut,
	}
	typFun := &ast.FuncType{
		Params:  &ast.FieldList{Opening: 1, Closing: 1},
		Results: &ast.FieldList{Opening: 1, Closing: 1, List: []*ast.Field{fldOut}},
	}
	makeArgs := []ast.Expr{typOut}
	if c.TypeOut.Kind() == reflect.Slice {
		makeArgs = append(makeArgs, IntConst(0), IntConst(4))
	}
	makeExpr := &ast.CallExpr{Fun: makeIden, Args: makeArgs}
	stmt := p.rhs.Pop().(ast.Stmt)
	stmtInit := &ast.AssignStmt{
		Lhs: []ast.Expr{qlangRet},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{makeExpr},
	}
	stmtReturn := &ast.ReturnStmt{}
	p.rhs.Push(&ast.CallExpr{
		Fun: &ast.FuncLit{
			Type: typFun,
			Body: &ast.BlockStmt{List: []ast.Stmt{stmtInit, stmt, stmtReturn}},
		},
	})
	p.comprehens = c.old
	return p
}

// ----------------------------------------------------------------------------
