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

// Package golang implements a golang backend for Go+ to generate Go code.
package golang

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"path"
	"reflect"
	"sort"
	"strconv"

	"github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/exec/golang/internal/go/format"
	"github.com/goplus/gop/exec/golang/internal/go/printer"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// A Code represents generated go code.
type Code struct {
	fset *token.FileSet
	file *ast.File
}

// NewCode returns a new Code object.
func NewCode() *Code {
	return &Code{}
}

// Document returns the whole ast tree.
func (p *Code) Document() *ast.File {
	return p.file
}

// Format code.
func (p *Code) Format(dst io.Writer) error {
	return format.Node(dst, p.fset, p.Document())
}

// Bytes returns go source code.
func (p *Code) Bytes(buf []byte) ([]byte, error) {
	b := bytes.NewBuffer(buf)
	err := p.Format(b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Len returns code length.
func (p *Code) Len() int {
	panic("don't call me")
}

func (p *Code) String() string {
	b, _ := p.Bytes(nil)
	return string(b)
}

// -----------------------------------------------------------------------------

type callType int

const (
	callExpr callType = iota
	callByDefer
	callByGo
)

// Builder is a class that generates go code.
type Builder struct {
	lhs, rhs    exec.Stack
	out         *Code                    // golang code
	pkgName     string                   // package name
	types       map[reflect.Type]*GoType // type => gotype
	imports     map[string]string        // pkgPath => aliasName
	importPaths map[string]string        // aliasName => pkgPath
	gblScope    scopeCtx                 // global scope
	gblDecls    []ast.Decl               // global declarations
	fset        *token.FileSet           // fileset of Go+ code
	cfun        *FuncInfo                // current function
	cstmt       interface{}              // current statement
	reserveds   []interface{}
	comprehens  func()   // current comprehension
	identBase   int      // auo-increasement ident index
	*scopeCtx            // current block scope
	inDeferOrGo callType // in defer/go statement currently
}

// NewBuilder creates a new Code Builder instance.
func NewBuilder(pkgName string, code *Code, fset *token.FileSet) *Builder {
	if code == nil {
		code = NewCode()
	}
	p := &Builder{
		out:         code,
		gblDecls:    make([]ast.Decl, 0, 4),
		types:       make(map[reflect.Type]*GoType),
		imports:     make(map[string]string),
		importPaths: make(map[string]string),
		fset:        fset,
		pkgName:     pkgName,
	}
	p.scopeCtx = &p.gblScope // default scope is global
	p.lhs.Init()
	p.rhs.Init()
	return p
}

func (p *Builder) autoIdent() string {
	p.identBase++
	return "_gop_" + strconv.Itoa(p.identBase)
}

var (
	tyMainFunc  = reflect.TypeOf((*func())(nil)).Elem()
	unnamedVar  = Ident("_")
	gopRet      = Ident("_gop_ret")
	appendIdent = Ident("append")
	makeIdent   = Ident("make")
	newIdent    = Ident("new")
	nilIdent    = Ident("nil")
)

// Resolve resolves all unresolved labels/functions/consts/etc.
func (p *Builder) Resolve() *Code {
	decls := make([]ast.Decl, 0, 8)
	imports := p.resolveImports()
	if imports != nil {
		decls = append(decls, imports)
	}
	types := p.resolveTypes()
	if types != nil {
		decls = append(decls, types)
	}
	gblvars := p.gblScope.toGenDecl(p)
	if gblvars != nil {
		decls = append(decls, gblvars)
	}
	p.endBlockStmt(0)

	if len(p.gblScope.stmts) != 0 {
		fnName := "main"
		for _, decl := range p.gblDecls {
			if d, ok := decl.(*ast.FuncDecl); ok && d.Name.Name == "main" {
				fnName = "init"
				break
			}
		}
		body := &ast.BlockStmt{List: p.gblScope.stmts}
		fn := &ast.FuncDecl{
			Name: Ident(fnName),
			Type: FuncType(p, tyMainFunc),
			Body: body,
		}
		decls = append(decls, fn)
	}
	decls = append(decls, p.gblDecls...)
	p.out.fset = token.NewFileSet()
	p.out.file = &ast.File{
		Name:  Ident(p.pkgName),
		Decls: decls,
	}
	return p.out
}

func (p *Builder) resolveImports() *ast.GenDecl {
	n := len(p.imports)
	if n == 0 {
		return nil
	}
	specs := make([]ast.Spec, 0, n)

	// stable sort import path
	pkgs := make([]string, 0, len(p.imports))
	for k := range p.imports {
		pkgs = append(pkgs, k)
	}
	sort.Strings(pkgs)

	for _, pkg := range pkgs {
		name := p.imports[pkg]
		spec := &ast.ImportSpec{
			Path: StringConst(pkg),
		}
		if name != "" {
			spec.Name = Ident(name)
		}
		specs = append(specs, spec)
	}

	return &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: specs,
	}
}

func (p *Builder) resolveTypes() *ast.GenDecl {
	n := len(p.types)
	specs := make([]ast.Spec, 0, n)
	for _, t := range p.types {
		typ := t.Type()
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		spec := &ast.TypeSpec{
			Name: Ident(t.Name()),
			Type: Type(p, typ, true),
		}
		specs = append(specs, spec)
	}
	if len(specs) > 0 {
		return &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: specs,
		}
	}
	return nil
}

// Comment instr
func Comment(text string) *ast.CommentGroup {
	return &ast.CommentGroup{
		List: []*ast.Comment{
			{Text: text},
		},
	}
}

type stmtState struct {
	stmtOld interface{}
	rhsBase int
}

// StartStmt receives a `StartStmt` event.
func (p *Builder) StartStmt(stmt interface{}) interface{} {
	state := &stmtState{p.cstmt, p.rhs.Len()}
	p.cstmt = stmt
	return state
}

// EndStmt receives a `EndStmt` event.
func (p *Builder) EndStmt(stmt, start interface{}) *Builder {
	var node ast.Stmt
	var state = start.(*stmtState)
	defer func() { // restore parent statement
		p.cstmt = state.stmtOld
	}()
	if lhsLen := p.lhs.Len(); lhsLen > 0 { // assignment
		lhs := make([]ast.Expr, lhsLen)
		for i := 0; i < lhsLen; i++ {
			lhs[i] = p.lhs.Pop().(ast.Expr)
		}
		rhsLen := p.rhs.Len() - state.rhsBase
		rhs := make([]ast.Expr, rhsLen)
		for i, v := range p.rhs.GetArgs(rhsLen) {
			rhs[i] = v.(ast.Expr)
		}
		p.rhs.PopN(rhsLen)
		node = &ast.AssignStmt{Lhs: lhs, Tok: token.ASSIGN, Rhs: rhs}
	} else {
		if rhsLen := p.rhs.Len() - state.rhsBase; rhsLen != 1 {
			if rhsLen == 0 {
				return p
			}
			log.Panicln("EndStmt: comma expression? -", p.rhs.Len(), "stmt:", reflect.TypeOf(stmt))
		}
		var val = p.rhs.Pop()
		switch v := val.(type) {
		case ast.Expr:
			node = &ast.ExprStmt{X: v}
		case ast.Stmt:
			node = v
		default:
			log.Panicln("EndStmt: unexpected -", reflect.TypeOf(val))
		}
	}
	p.emitStmt(node)
	return p
}

func (p *Builder) emitStmt(node ast.Stmt) {
	if stmt := p.cstmt; stmt != nil {
		start := stmt.(ast.Node).Pos()
		pos := p.fset.Position(start)
		line := fmt.Sprintf("\n//line ./%s:%d", path.Base(pos.Filename), pos.Line)
		if node == nil {
			panic("node nil")
		}
		node = &printer.CommentedStmt{Comments: Comment(line), Stmt: node}
	}
	p.stmts = append(p.stmts, p.labeled(node, 0))
}

func (p *Builder) endBlockStmt(isEndFunc int) {
	if stmt := p.labeled(nil, isEndFunc); stmt != nil {
		p.stmts = append(p.stmts, stmt)
	}
}

func (p *Builder) labeled(stmt ast.Stmt, isEndFunc int) ast.Stmt {
	if p.labels != nil {
		if stmt == nil {
			stmt = endStmts[isEndFunc]
		}
		for _, l := range p.labels {
			stmt = &ast.LabeledStmt{
				Label: Ident(l.getName(p)),
				Stmt:  stmt,
			}
		}
		p.labels = nil
	}
	return stmt
}

var endStmts = [2]ast.Stmt{
	&ast.EmptyStmt{},
	&ast.ReturnStmt{},
}

// Import imports a package by pkgPath.
func (p *Builder) Import(pkgPath string) string {
	if name, ok := p.imports[pkgPath]; ok {
		return name
	}
	name := path.Base(pkgPath)
	if _, exists := p.importPaths[name]; exists {
		name = "q" + strconv.Itoa(len(p.imports)) + name
	}
	p.imports[pkgPath] = name
	p.importPaths[name] = pkgPath
	return name
}

// Reserve reserves an instruction.
func (p *Builder) Reserve() exec.Reserved {
	r := new(printer.ReservedExpr)
	idx := len(p.reserveds)
	p.reserveds = append(p.reserveds, r)
	p.rhs.Push(r)
	return exec.Reserved(idx)
}

// ReservedAsPush sets Reserved as Push(v)
func (p *Builder) ReservedAsPush(r exec.Reserved, v interface{}) {
	p.reserveds[r].(*printer.ReservedExpr).Expr = Const(p, v)
}

type reservedShift struct {
	expr *printer.ReservedExpr
	x    ast.Expr
	y    ast.Expr
}

// ReserveOpShift reserves an instruction.
func (p *Builder) ReserveOpShift() exec.Reserved {
	r := &reservedShift{}
	r.expr = new(printer.ReservedExpr)
	r.y = p.rhs.Pop().(ast.Expr)
	r.x = p.rhs.Pop().(ast.Expr)
	idx := len(p.reserveds)
	p.reserveds = append(p.reserveds, r)
	p.rhs.Push(r.expr)
	return exec.Reserved(idx)
}

// ReservedAsOpShift sets Reserved as OpLsh/OpRsh
func (p *Builder) ReservedAsOpShift(r exec.Reserved, kind exec.Kind, op exec.Operator) {
	tok := opTokens[op]
	rsh := p.reserveds[r].(*reservedShift)
	rsh.expr.Expr = &ast.BinaryExpr{Op: tok, X: rsh.x, Y: rsh.y}
}

// Pop instr
func (p *Builder) Pop(n int) *Builder {
	log.Panicln("todo")
	return p
}

// -----------------------------------------------------------------------------
