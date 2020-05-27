package exec

import (
	"reflect"
)

// -----------------------------------------------------------------------------

// Var represents a variable.
type Var interface {
	// Name returns the variable name.
	Name() string

	// Type returns the variable type.
	Type() reflect.Type

	// IsUnnamedOut returns if variable unnamed or not.
	IsUnnamedOut() bool
}

// Label represents a label.
type Label interface {
	// Name returns the label name.
	Name() string
}

// Reserved represents a reserved instruction position.
type Reserved int

// InvalidReserved is an invalid reserved position.
const InvalidReserved Reserved = -1

// Push instr
func (p Reserved) Push(b Builder, val interface{}) {
	b.ReservedPush(p, val)
}

// ForPhrase represents a for range phrase.
type ForPhrase interface {
}

// Comprehension represents a list/map comprehension.
type Comprehension interface {
}

// FuncInfo represents a qlang function information.
type FuncInfo interface {
	// Name returns the function name.
	Name() string

	// Type returns type of this function.
	Type() reflect.Type

	// Args sets argument types of a qlang function.
	Args(in ...reflect.Type) FuncInfo

	// Vargs sets argument types of a variadic qlang function.
	Vargs(in ...reflect.Type) FuncInfo

	// Return sets return types of a qlang function.
	Return(out ...Var) FuncInfo

	// NumOut returns a function type's output parameter count.
	// It panics if the type's Kind is not Func.
	NumOut() int

	// Out returns the type of a function type's i'th output parameter.
	// It panics if i is not in the range [0, NumOut()).
	Out(i int) Var

	// IsVariadic returns if this function is variadic or not.
	IsVariadic() bool

	// IsUnnamedOut returns if function results unnamed or not.
	IsUnnamedOut() bool
}

// SymbolKind represents symbol kind.
type SymbolKind uint32

const (
	// SymbolVar - variable
	SymbolVar SymbolKind = 0
	// SymbolFunc - function
	SymbolFunc SymbolKind = 1
	// SymbolFuncv - variadic function
	SymbolFuncv SymbolKind = 2
)

// GoFuncAddr represents a Go function address.
type GoFuncAddr uint32

// GoFuncvAddr represents a variadic Go function address.
type GoFuncvAddr uint32

// GoVarAddr represents a variadic Go variable address.
type GoVarAddr uint32

// GoPackage represents a Go package.
type GoPackage interface {
	PkgPath() string

	// Find lookups a symbol by specified its name.
	Find(name string) (addr uint32, kind SymbolKind, ok bool)

	// FindType lookups a Go type by name.
	FindType(name string) (typ reflect.Type, ok bool)

	// FindConst lookups a Go constant by name.
	FindConst(name string) (ci *GoConstInfo, ok bool)
}

// Builder represents a executing byte code generator.
type Builder interface {
	// Push instr
	Push(val interface{}) Builder

	// Pop instr
	Pop(n int) Builder

	// BuiltinOp instr
	BuiltinOp(kind Kind, op Operator) Builder

	// Label defines a label to jmp here.
	Label(l Label) Builder

	// Jmp instr
	Jmp(l Label) Builder

	// JmpIf instr
	JmpIf(zeroOrOne uint32, l Label) Builder

	// CaseNE instr
	CaseNE(l Label, arity int) Builder

	// Default instr
	Default() Builder

	// ForPhrase instr
	ForPhrase(f ForPhrase, key, val Var, hasExecCtx ...bool) Builder

	// FilterForPhrase instr
	FilterForPhrase(f ForPhrase) Builder

	// EndForPhrase instr
	EndForPhrase(f ForPhrase) Builder

	// ListComprehension instr
	ListComprehension(c Comprehension) Builder

	// MapComprehension instr
	MapComprehension(c Comprehension) Builder

	// EndComprehension instr
	EndComprehension(c Comprehension) Builder

	// Closure instr
	Closure(fun FuncInfo) Builder

	// GoClosure instr
	GoClosure(fun FuncInfo) Builder

	// CallClosure instr
	CallClosure(arity int) Builder

	// CallGoClosure instr
	CallGoClosure(arity int) Builder

	// CallFunc instr
	CallFunc(fun FuncInfo) Builder

	// CallFuncv instr
	CallFuncv(fun FuncInfo, arity int) Builder

	// CallGoFunc instr
	CallGoFunc(fun GoFuncAddr) Builder

	// CallGoFuncv instr
	CallGoFuncv(fun GoFuncvAddr, arity int) Builder

	// DefineFunc instr
	DefineFunc(fun FuncInfo) Builder

	// Return instr
	Return(n int32) Builder

	// Load instr
	Load(idx int32) Builder

	// Store instr
	Store(idx int32) Builder

	// EndFunc instr
	EndFunc(fun FuncInfo) Builder

	// DefineVar defines variables.
	DefineVar(vars ...Var) Builder

	// InCurrentCtx returns if a variable is in current context or not.
	InCurrentCtx(v Var) bool

	// LoadVar instr
	LoadVar(v Var) Builder

	// StoreVar instr
	StoreVar(v Var) Builder

	// AddrVar instr
	AddrVar(v Var) Builder

	// AddrOp instr
	AddrOp(kind Kind, op AddrOperator) Builder

	// Append instr
	Append(typ reflect.Type, arity int) Builder

	// MakeArray instr
	MakeArray(typ reflect.Type, arity int) Builder

	// MakeMap instr
	MakeMap(typ reflect.Type, arity int) Builder

	// Make instr
	Make(typ reflect.Type, arity int) Builder

	// MapIndex instr
	MapIndex() Builder

	// SetMapIndex instr
	SetMapIndex() Builder

	// Index instr
	Index(idx int) Builder

	// SetIndex instr
	SetIndex(idx int) Builder

	// Slice instr
	Slice(i, j int) Builder

	// Slice3 instr
	Slice3(i, j, k int) Builder

	// TypeCast instr
	TypeCast(from, to reflect.Type) Builder

	// Zero instr
	Zero(typ reflect.Type) Builder

	// Reserve reserves an instruction.
	Reserve() Reserved

	// ReservedPush sets Reserved as Push(v)
	ReservedPush(r Reserved, v interface{})

	// GlobalInterface returns the global Interface.
	GlobalInterface() Interface
}

// Interface represents all global functions of a executing byte code generator.
type Interface interface {
	// NewVar creates a variable instance.
	NewVar(typ reflect.Type, name string) Var

	// NewLabel creates a label object.
	NewLabel(name string) Label

	// NewForPhrase creates a new ForPhrase instance.
	NewForPhrase(in reflect.Type) ForPhrase

	// NewComprehension creates a new Comprehension instance.
	NewComprehension(out reflect.Type) Comprehension

	// NewFunc create a qlang function.
	NewFunc(name string, nestDepth uint32) FuncInfo

	// FindGoPackage lookups a Go package by pkgPath. It returns nil if not found.
	FindGoPackage(pkgPath string) GoPackage

	// GetGoFuncType returns a Go function's type.
	GetGoFuncType(addr GoFuncAddr) reflect.Type

	// GetGoFuncvType returns a Go function's type.
	GetGoFuncvType(addr GoFuncvAddr) reflect.Type
}

// -----------------------------------------------------------------------------
