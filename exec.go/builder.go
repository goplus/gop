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

package exec

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"log"
	"path"
	"reflect"
	"strconv"

	"github.com/qiniu/qlang/v6/exec.go/format"
	"github.com/qiniu/qlang/v6/exec.go/printer"
	"github.com/qiniu/qlang/v6/exec.spec"

	qexec "github.com/qiniu/qlang/v6/exec"
)

var defaultImpl = qexec.GlobalInterface()

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

func (p *Code) String() string {
	b, _ := p.Bytes(nil)
	return string(b)
}

// -----------------------------------------------------------------------------

// Builder is a class that generates go code.
type Builder struct {
	code        exec.Stack
	out         *Code
	imports     map[string]string
	importPaths map[string]string
	stmts       []ast.Stmt
	fset        *token.FileSet
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
	p.code.Init()
	return p
}

var (
	// TyMainFunc type
	TyMainFunc = reflect.TypeOf((*func())(nil)).Elem()
)

// Resolve resolves all unresolved labels/functions/consts/etc.
func (p *Builder) Resolve() *Code {
	decls := make([]ast.Decl, 0, 8)
	imports := p.resolveImports()
	if imports != nil {
		decls = append(decls, imports)
	}
	if len(p.stmts) != 0 {
		body := &ast.BlockStmt{List: p.stmts}
		fn := &ast.FuncDecl{
			Name: Ident("main"),
			Type: FuncType(p, TyMainFunc),
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
	var val = p.code.Pop()
	switch v := val.(type) {
	case ast.Expr:
		node = &ast.ExprStmt{X: v}
	case ast.Stmt:
		node = v
	default:
		log.Panicln("EndStmt: unexpected -", reflect.TypeOf(val))
	}
	if stmt != nil {
		start := stmt.(ast.Node).Pos()
		pos := p.fset.Position(start)
		line := fmt.Sprintf("\n//line %s:%d:%d", pos.Filename, pos.Line, pos.Column)
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

// Pop instr
func (p *Builder) Pop(n int) *Builder {
	return p
}

// Label defines a label to jmp here.
func (p *Builder) Label(l exec.Label) *Builder {
	return p
}

// Jmp instr
func (p *Builder) Jmp(l exec.Label) *Builder {
	return p
}

// JmpIf instr
func (p *Builder) JmpIf(zeroOrOne uint32, l exec.Label) *Builder {
	return p
}

// CaseNE instr
func (p *Builder) CaseNE(l exec.Label, arity int) *Builder {
	return p
}

// Default instr
func (p *Builder) Default() *Builder {
	return p
}

// ForPhrase instr
func (p *Builder) ForPhrase(f exec.ForPhrase, key, val exec.Var, hasExecCtx ...bool) *Builder {
	return p
}

// FilterForPhrase instr
func (p *Builder) FilterForPhrase(f exec.ForPhrase) *Builder {
	return p
}

// EndForPhrase instr
func (p *Builder) EndForPhrase(f exec.ForPhrase) *Builder {
	return p
}

// ListComprehension instr
func (p *Builder) ListComprehension(c exec.Comprehension) *Builder {
	return p
}

// MapComprehension instr
func (p *Builder) MapComprehension(c exec.Comprehension) *Builder {
	return p
}

// EndComprehension instr
func (p *Builder) EndComprehension(c exec.Comprehension) *Builder {
	return p
}

// Closure instr
func (p *Builder) Closure(fun exec.FuncInfo) *Builder {
	return p
}

// GoClosure instr
func (p *Builder) GoClosure(fun exec.FuncInfo) *Builder {
	return p
}

// CallClosure instr
func (p *Builder) CallClosure(arity int) *Builder {
	return p
}

// CallGoClosure instr
func (p *Builder) CallGoClosure(arity int) *Builder {
	return p
}

// CallFunc instr
func (p *Builder) CallFunc(fun exec.FuncInfo) *Builder {
	return p
}

// CallFuncv instr
func (p *Builder) CallFuncv(fun exec.FuncInfo, arity int) *Builder {
	return p
}

// DefineFunc instr
func (p *Builder) DefineFunc(fun exec.FuncInfo) *Builder {
	return p
}

// Return instr
func (p *Builder) Return(n int32) *Builder {
	return p
}

// Load instr
func (p *Builder) Load(idx int32) *Builder {
	return p
}

// Store instr
func (p *Builder) Store(idx int32) *Builder {
	return p
}

// EndFunc instr
func (p *Builder) EndFunc(fun exec.FuncInfo) *Builder {
	return p
}

// DefineVar defines variables.
func (p *Builder) DefineVar(vars ...exec.Var) *Builder {
	return p
}

// InCurrentCtx returns if a variable is in current context or not.
func (p *Builder) InCurrentCtx(v exec.Var) bool {
	return false
}

// LoadVar instr
func (p *Builder) LoadVar(v exec.Var) *Builder {
	return p
}

// StoreVar instr
func (p *Builder) StoreVar(v exec.Var) *Builder {
	return p
}

// AddrVar instr
func (p *Builder) AddrVar(v exec.Var) *Builder {
	return p
}

// AddrOp instr
func (p *Builder) AddrOp(kind exec.Kind, op exec.AddrOperator) *Builder {
	return p
}

// Append instr
func (p *Builder) Append(typ reflect.Type, arity int) *Builder {
	return p
}

// MakeArray instr
func (p *Builder) MakeArray(typ reflect.Type, arity int) *Builder {
	return p
}

// MakeMap instr
func (p *Builder) MakeMap(typ reflect.Type, arity int) *Builder {
	return p
}

// Make instr
func (p *Builder) Make(typ reflect.Type, arity int) *Builder {
	return p
}

// MapIndex instr
func (p *Builder) MapIndex() *Builder {
	return p
}

// SetMapIndex instr
func (p *Builder) SetMapIndex() *Builder {
	return p
}

// Index instr
func (p *Builder) Index(idx int) *Builder {
	return p
}

// SetIndex instr
func (p *Builder) SetIndex(idx int) *Builder {
	return p
}

// Slice instr
func (p *Builder) Slice(i, j int) *Builder {
	return p
}

// Slice3 instr
func (p *Builder) Slice3(i, j, k int) *Builder {
	return p
}

// Zero instr
func (p *Builder) Zero(typ reflect.Type) *Builder {
	return p
}

// Reserve reserves an instruction.
func (p *Builder) Reserve() exec.Reserved {
	return exec.InvalidReserved
}

// ReservedAsPush sets Reserved as Push(v)
func (p *Builder) ReservedAsPush(r exec.Reserved, v interface{}) {
}

// GlobalInterface returns the global Interface.
func (p *Builder) GlobalInterface() exec.Interface {
	return defaultImpl
}

// -----------------------------------------------------------------------------
