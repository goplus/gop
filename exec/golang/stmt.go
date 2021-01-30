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

	"github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/exec/golang/internal/go/printer"
	"github.com/qiniu/x/errors"
)

// ----------------------------------------------------------------------------

// Label represents a label.
type Label struct {
	name      string
	wrapIfErr bool // for wrapIfErr instr
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
	if l.wrapIfErr {
		defval := p.rhs.Pop().(ast.Expr)
		reservedExpr := p.rhs.Pop().(*printer.ReservedExpr)
		reservedExpr.Expr = defval
		return p
	}
	// make sure all labels in golang code  will be used
	p.Jmp(l)
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
func (p *Builder) JmpIf(jc exec.JmpCondFlag, l *Label) *Builder {
	if jc.IsNotPop() {
		return p
	}
	cond := p.rhs.Pop().(ast.Expr)
	switch jc {
	case exec.JcFalse:
		cond = &ast.UnaryExpr{Op: token.NOT, X: cond}
	case exec.JcNil:
		cond = &ast.BinaryExpr{Op: token.EQL, X: cond, Y: nilIdent}
	case exec.JcNotNil:
		cond = &ast.BinaryExpr{Op: token.NEQ, X: cond, Y: nilIdent}
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
	var cond *ast.BinaryExpr
	for i := 1; i <= arity; i++ {
		if cond == nil {
			cond = &ast.BinaryExpr{X: x, Op: token.NEQ, Y: args[i].(ast.Expr)}
		} else {
			cond = &ast.BinaryExpr{
				X:     cond,
				OpPos: 0,
				Op:    token.LAND,
				Y:     &ast.BinaryExpr{X: x, Op: token.NEQ, Y: args[i].(ast.Expr)},
			}
		}
	}
	p.emitStmt(GotoIf(p, cond, l))
	p.rhs.PopN(arity)
	return p
}

// Default instr
func (p *Builder) Default() *Builder {
	p.rhs.PopN(1)
	return p
}

// WrapIfErr instr
func (p *Builder) WrapIfErr(nret int, l *Label) *Builder {
	args := p.wrapCall(nret)
	reservedExpr := &printer.ReservedExpr{}
	assignVal := &ast.AssignStmt{
		Lhs: []ast.Expr{args[0]},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{reservedExpr},
	}
	ifStmt := &ast.IfStmt{
		Cond: &ast.BinaryExpr{X: args[nret-1], Op: token.NEQ, Y: nilIdent},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{assignVal},
		},
	}
	p.emitStmt(ifStmt)
	p.rhs.Push(args[0])
	p.rhs.Push(reservedExpr)
	l.wrapIfErr = true
	return p
}

func (p *Builder) wrapCall(nret int) []ast.Expr {
	x := p.rhs.Pop().(ast.Expr)
	args := make([]ast.Expr, nret)
	for i := 0; i < nret; i++ {
		args[i] = Ident(p.autoIdent())
	}
	assignStmt := &ast.AssignStmt{
		Lhs: args,
		Tok: token.DEFINE,
		Rhs: []ast.Expr{x},
	}
	p.emitStmt(assignStmt)
	return args
}

// ErrWrap instr
func (p *Builder) ErrWrap(nret int, retErr exec.Var, frame *errors.Frame, narg int) *Builder {
	args := p.wrapCall(nret)
	frameFun := p.GoSymIdent("github.com/qiniu/x/errors", "NewFrame")
	frameArgs := make([]ast.Expr, narg+6)
	frameArgs[0] = args[nret-1]
	frameArgs[1] = StringConst(frame.Code)
	frameArgs[2] = StringConst(frame.File)
	frameArgs[3] = IntConst(int64(frame.Line))
	frameArgs[4] = StringConst(frame.Pkg)
	frameArgs[5] = StringConst(frame.Func)
	for i := 0; i < narg; i++ {
		frameArgs[6+i] = Ident(toArg(i, p.cfun.nestDepth))
	}
	frameErr := &ast.CallExpr{Fun: frameFun, Args: frameArgs}
	var body *ast.BlockStmt
	if retErr != nil {
		outErr := Ident(retErr.Name())
		assignErr := &ast.AssignStmt{
			Lhs: []ast.Expr{outErr},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{frameErr},
		}
		returnStmt := &ast.ReturnStmt{}
		body = &ast.BlockStmt{
			List: []ast.Stmt{assignErr, returnStmt},
		}
	} else {
		panicArgs := []ast.Expr{frameErr}
		panicStmt := &ast.ExprStmt{X: &ast.CallExpr{Fun: Ident("panic"), Args: panicArgs}}
		body = &ast.BlockStmt{
			List: []ast.Stmt{panicStmt},
		}
	}
	ifStmt := &ast.IfStmt{
		Cond: &ast.BinaryExpr{X: args[nret-1], Op: token.NEQ, Y: nilIdent},
		Body: body,
	}
	p.emitStmt(ifStmt)
	for _, arg := range args[:nret-1] {
		p.rhs.Push(arg)
	}
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
	p.endBlockStmt(0)
	body := &ast.BlockStmt{List: f.getStmts(p)}
	if f.Cond != nil {
		body = &ast.BlockStmt{List: []ast.Stmt{
			&ast.IfStmt{Cond: f.Cond, Body: body},
		}}
	}
	p.rhs.Push(&ast.RangeStmt{
		Key: func() ast.Expr {
			if f.Key == nil && f.Value == nil {
				return nil
			}
			return toVarExpr(f.Key, unnamedVar)
		}(),
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
			Fun:  appendIdent,
			Args: []ast.Expr{gopRet, x},
		}
		assign := &ast.AssignStmt{
			Lhs: []ast.Expr{gopRet},
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
			Lhs: []ast.Expr{&ast.IndexExpr{X: gopRet, Index: key}},
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
		Names: []*ast.Ident{gopRet},
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
	makeExpr := &ast.CallExpr{Fun: makeIdent, Args: makeArgs}
	stmt := p.rhs.Pop().(ast.Stmt)
	stmtInit := &ast.AssignStmt{
		Lhs: []ast.Expr{gopRet},
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

// Defer instr
func (p *Builder) Defer() *Builder {
	p.inDeferOrGo = callByDefer
	return p
}

// Send instr
func (p *Builder) Send() *Builder {
	args := p.rhs.GetArgs(2)
	p.rhs.PopN(2)
	stmt := &ast.SendStmt{
		Chan:  args[0].(ast.Expr),
		Value: args[1].(ast.Expr),
	}
	p.emitStmt(stmt)
	return p
}

// Recv instr
func (p *Builder) Recv() *Builder {
	expr := &ast.UnaryExpr{
		Op: token.ARROW,
		X:  p.rhs.Get(-1).(ast.Expr),
	}
	p.rhs.Ret(1, expr)
	return p
}

// Go instr
func (p *Builder) Go() *Builder {
	p.inDeferOrGo = callByGo
	return p
}

// DefineBlock starts a new block.
func (p *Builder) DefineBlock() *Builder {
	p.scopeCtx = &scopeCtx{parentCtx: p.scopeCtx}
	p.initStmts()
	return p
}

// EndBlock ends a block.
func (p *Builder) EndBlock() *Builder {
	p.endBlockStmt(0)
	blockStmt := &ast.BlockStmt{List: p.getStmts(p)}
	p.scopeCtx = p.parentCtx
	p.stmts = append(p.stmts, p.labeled(blockStmt, 0))
	return p
}

// ----------------------------------------------------------------------------
