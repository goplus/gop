package exec

import (
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"path"
	"reflect"
	"strconv"

	"github.com/qiniu/qlang/v6/exec.spec"
)

// -----------------------------------------------------------------------------

// A Code represents generated go code.
type Code struct {
	fset *token.FileSet
}

// NewCode returns a new Code object.
func NewCode() *Code {
	return &Code{}
}

// Document returns the whole ast tree.
func (p *Code) Document() *ast.File {
	return &ast.File{}
}

// Format code.
func (p *Code) Format(dst io.Writer) error {
	return format.Node(dst, p.fset, p.Document())
}

// -----------------------------------------------------------------------------

// Builder is a class that generates go code.
type Builder struct {
	code        exec.Stack
	out         *Code
	imports     map[string]string
	importPaths map[string]string
}

// NewBuilder creates a new Code Builder instance.
func NewBuilder(code *Code) *Builder {
	if code == nil {
		code = NewCode()
	}
	p := &Builder{
		out:         code,
		imports:     make(map[string]string),
		importPaths: make(map[string]string),
	}
	p.code.Init()
	return p
}

// Import imports a package by pkgPath
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

// CallGoFunc instr
func (p *Builder) CallGoFunc(fun exec.GoFuncAddr) *Builder {
	return p
}

// CallGoFuncv instr
func (p *Builder) CallGoFuncv(fun exec.GoFuncvAddr, arity int) *Builder {
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
	return nil
}

// Resolve resolves all unresolved labels/functions/consts/etc.
func (p *Builder) Resolve() *Code {
	return p.out
}

// -----------------------------------------------------------------------------
