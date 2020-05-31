package golang

import (
	"go/ast"
	"go/token"
	"reflect"

	"github.com/qiniu/qlang/v6/exec.spec"
	"github.com/qiniu/x/log"
)

// ----------------------------------------------------------------------------

// Return instr
func (p *Builder) Return(n int32) *Builder {
	var results []ast.Expr
	if n > 0 {
		log.Panicln("todo")
	}
	p.rhs.Push(&ast.ReturnStmt{Results: results})
	return p
}

// Label defines a label to jmp here.
func (p *Builder) Label(l exec.Label) *Builder {
	log.Panicln("todo")
	return p
}

// Jmp instr
func (p *Builder) Jmp(l exec.Label) *Builder {
	log.Panicln("todo")
	return p
}

// JmpIf instr
func (p *Builder) JmpIf(zeroOrOne uint32, l exec.Label) *Builder {
	log.Panicln("todo")
	return p
}

// CaseNE instr
func (p *Builder) CaseNE(l exec.Label, arity int) *Builder {
	log.Panicln("todo")
	return p
}

// Default instr
func (p *Builder) Default() *Builder {
	log.Panicln("todo")
	return p
}

// ----------------------------------------------------------------------------

type blockCtx struct {
	parent *varManager
	old    *[]ast.Stmt
}

func (p *blockCtx) saveEnv(b *Builder) {
	if p.parent != nil {
		panic("blockCtx.setParent: already defined")
	}
	p.parent = b.varManager
	p.old = b.stmts
}

func (p *blockCtx) restoreEnv(b *Builder) {
	b.varManager = p.parent
	b.stmts = p.old
}

// ForPhrase represents a for range phrase.
type ForPhrase struct {
	Key, Value *Var // Key, Value may be nil
	X, Cond    ast.Expr
	TypeIn     reflect.Type
	stmts      []ast.Stmt
	varManager
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
	f.saveEnv(p)
	p.varManager = &f.varManager
	p.stmts = &f.stmts
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
	body := &ast.BlockStmt{List: f.stmts}
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
			*p.stmts = append(*p.stmts, stmt)
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
		*p.stmts = append(*p.stmts, assign)
	}
	return p
}

// MapComprehension instr
func (p *Builder) MapComprehension(c *Comprehension) *Builder {
	c.old = p.comprehens
	p.comprehens = func() {
		v := p.rhs.Pop()
		if stmt, ok := v.(*ast.RangeStmt); ok {
			*p.stmts = append(*p.stmts, stmt)
			return
		}
		val := v.(ast.Expr)
		key := p.rhs.Pop().(ast.Expr)
		assign := &ast.AssignStmt{
			Lhs: []ast.Expr{&ast.IndexExpr{X: qlangRet, Index: key}},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{val},
		}
		*p.stmts = append(*p.stmts, assign)
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
