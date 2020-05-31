/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

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
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"path"
	"reflect"
	"strconv"

	"github.com/qiniu/qlang/v6/exec.spec"
	"github.com/qiniu/qlang/v6/exec/golang/internal/go/format"
	"github.com/qiniu/qlang/v6/exec/golang/internal/go/printer"
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

// Builder is a class that generates go code.
type Builder struct {
	lhs, rhs    exec.Stack
	out         *Code
	imports     map[string]string
	importPaths map[string]string
	gblvars     varManager
	stmts       []ast.Stmt
	fset        *token.FileSet
	reserveds   []*printer.ReservedExpr
	*varManager
}

// NewBuilder creates a new Code Builder instance.
func NewBuilder(code *Code, fset *token.FileSet) *Builder {
	if code == nil {
		code = NewCode()
	}
	p := &Builder{
		out:         code,
		imports:     make(map[string]string),
		importPaths: make(map[string]string),
		fset:        fset,
	}
	p.varManager = &p.gblvars // default scope is global
	p.lhs.Init()
	p.rhs.Init()
	return p
}

var (
	tyMainFunc = reflect.TypeOf((*func())(nil)).Elem()
)

// Resolve resolves all unresolved labels/functions/consts/etc.
func (p *Builder) Resolve() *Code {
	decls := make([]ast.Decl, 0, 8)
	imports := p.resolveImports()
	if imports != nil {
		decls = append(decls, imports)
	}
	gblvars := p.gblvars.toGenDecl(p)
	if gblvars != nil {
		decls = append(decls, gblvars)
	}
	if len(p.stmts) != 0 {
		body := &ast.BlockStmt{List: p.stmts}
		fn := &ast.FuncDecl{
			Name: Ident("main"),
			Type: FuncType(p, tyMainFunc),
			Body: body,
		}
		decls = append(decls, fn)
	}
	p.out.fset = token.NewFileSet()
	p.out.file = &ast.File{
		Name:  Ident("main"),
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
	for pkgPath, name := range p.imports {
		spec := &ast.ImportSpec{
			Path: StringConst(pkgPath),
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

// Comment instr
func Comment(text string) *ast.CommentGroup {
	return &ast.CommentGroup{
		List: []*ast.Comment{
			{Text: text},
		},
	}
}

// EndStmt recieves a `EndStmt` event.
func (p *Builder) EndStmt(stmt interface{}) *Builder {
	var node ast.Stmt
	if lhsLen := p.lhs.Len(); lhsLen > 0 { // assignment
		lhs := make([]ast.Expr, lhsLen)
		for i := 0; i < lhsLen; i++ {
			lhs[i] = p.lhs.Pop().(ast.Expr)
		}
		rhsLen := p.rhs.Len()
		rhs := make([]ast.Expr, rhsLen)
		for i, v := range p.rhs.GetArgs(rhsLen) {
			rhs[i] = v.(ast.Expr)
		}
		p.rhs.PopN(rhsLen)
		node = &ast.AssignStmt{Lhs: lhs, Tok: token.ASSIGN, Rhs: rhs}
	} else {
		if p.rhs.Len() != 1 {
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
	if stmt != nil {
		start := stmt.(ast.Node).Pos()
		pos := p.fset.Position(start)
		line := fmt.Sprintf("\n//line ./%s:%d", path.Base(pos.Filename), pos.Line)
		node = &printer.CommentedStmt{Comments: Comment(line), Stmt: node}
	}
	p.stmts = append(p.stmts, node)
	return p
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
	p.reserveds[r].Expr = Const(p, v)
}

// Pop instr
func (p *Builder) Pop(n int) *Builder {
	log.Panicln("todo")
	return p
}

// -----------------------------------------------------------------------------
