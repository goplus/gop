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

// Package exec defines the specification of a Go+ backend.
package exec

import (
	"reflect"

	"github.com/qiniu/x/errors"
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
	b.ReservedAsPush(p, val)
}

// ForPhrase represents a for range phrase.
type ForPhrase interface {
}

// Comprehension represents a list/map comprehension.
type Comprehension interface {
}

// FuncInfo represents a Go+ function information.
type FuncInfo interface {
	// Name returns the function name.
	Name() string

	// Type returns type of this function.
	Type() reflect.Type

	// Args sets argument types of a Go+ function.
	Args(in ...reflect.Type) FuncInfo

	// Vargs sets argument types of a variadic Go+ function.
	Vargs(in ...reflect.Type) FuncInfo

	// Return sets return types of a Go+ function.
	Return(out ...Var) FuncInfo

	// NumIn returns a function's input parameter count.
	NumIn() int

	// NumOut returns a function's output parameter count.
	NumOut() int

	// Out returns the type of a function type's i'th output parameter.
	// It panics if i is not in the range [0, NumOut()).
	Out(i int) Var

	// IsVariadic returns if this function is variadic or not.
	IsVariadic() bool

	// IsUnnamedOut returns if function results unnamed or not.
	IsUnnamedOut() bool
}

// JmpCondFlag represents condition of Jmp intruction.
type JmpCondFlag uint32

const (
	// JcFalse - JmpIfFalse
	JcFalse JmpCondFlag = 0
	// JcTrue - JmpIfTrue
	JcTrue JmpCondFlag = 1
	// JcNil - JmpIfNil
	JcNil JmpCondFlag = 2
	// JcNotNil - JmpIfNotNil
	JcNotNil JmpCondFlag = 3
	// JcNotPopMask - jump but not pop
	JcNotPopMask JmpCondFlag = 4
)

// IsNotPop returns to pop condition value or not
func (v JmpCondFlag) IsNotPop() bool {
	return v&JcNotPopMask == JcNotPopMask
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

	// FindFunc lookups a Go function by name.
	FindFunc(name string) (addr GoFuncAddr, ok bool)

	// FindFuncv lookups a Go function by name.
	FindFuncv(name string) (addr GoFuncvAddr, ok bool)

	// FindType lookups a Go type by name.
	FindType(name string) (typ reflect.Type, ok bool)

	// FindConst lookups a Go constant by name.
	FindConst(name string) (ci *GoConstInfo, ok bool)
}

// GoSymInfo represents a Go symbol (function or variable) information.
type GoSymInfo struct {
	Pkg  GoPackage
	Name string
	This interface{} // address of this symbol
}

// GoFuncInfo represents a Go function information.
type GoFuncInfo = GoSymInfo

// GoVarInfo represents a Go variable information.
type GoVarInfo = GoSymInfo

// A Code represents generated instructions to execute.
type Code interface {
	// Len returns code length.
	Len() int
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
	JmpIf(cond JmpCondFlag, l Label) Builder

	// CaseNE instr
	CaseNE(l Label, arity int) Builder

	// Default instr
	Default() Builder

	// WrapIfErr instr
	WrapIfErr(nret int, l Label) Builder

	// ErrWrap instr
	ErrWrap(nret int, retErr Var, frame *errors.Frame, narg int) Builder

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
	CallClosure(nexpr, arity int, ellipsis bool) Builder

	// CallGoClosure instr
	CallGoClosure(nexpr, arity int, ellipsis bool) Builder

	// CallFunc instr
	CallFunc(fun FuncInfo, nexpr int) Builder

	// CallFuncv instr
	CallFuncv(fun FuncInfo, nexpr, arity int) Builder

	// CallGoFunc instr
	CallGoFunc(fun GoFuncAddr, nexpr int) Builder

	// CallGoFuncv instr
	CallGoFuncv(fun GoFuncvAddr, nexpr, arity int) Builder

	// GoBuiltin instr
	GoBuiltin(typ reflect.Type, op GoBuiltin) Builder

	// Defer instr
	Defer() Builder

	// Go instr
	Go() Builder

	// DefineFunc instr
	DefineFunc(fun FuncInfo) Builder

	// DefineType name string,reflect.Typeinstr
	DefineType(typ reflect.Type, name string) Builder

	// Return instr
	Return(n int32) Builder

	// Load instr
	Load(idx int32) Builder

	// Addr instr
	Addr(idx int32) Builder

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

	// LoadGoVar instr
	LoadGoVar(addr GoVarAddr) Builder

	// StoreGoVar instr
	StoreGoVar(addr GoVarAddr) Builder

	// AddrGoVar instr
	AddrGoVar(addr GoVarAddr) Builder

	// LoadField instr
	LoadField(typ reflect.Type, index []int) Builder

	// StoreField instr
	StoreField(typ reflect.Type, index []int) Builder

	// AddrField instr
	AddrField(typ reflect.Type, index []int) Builder

	// AddrOp instr
	AddrOp(kind Kind, op AddrOperator) Builder

	// MakeArray instr
	MakeArray(typ reflect.Type, arity int) Builder

	// MakeMap instr
	MakeMap(typ reflect.Type, arity int) Builder

	// Make instr
	Make(typ reflect.Type, arity int) Builder

	// Val instr
	Struct(typ reflect.Type, arity int) Builder

	// Append instr
	Append(typ reflect.Type, arity int) Builder

	// MapIndex instr
	MapIndex(twoValue bool) Builder

	// SetMapIndex instr
	SetMapIndex() Builder

	// Index instr
	Index(idx int) Builder

	// AddrIndex instr
	AddrIndex(idx int) Builder

	// SetIndex instr
	SetIndex(idx int) Builder

	// Slice instr
	Slice(i, j int) Builder

	// Slice3 instr
	Slice3(i, j, k int) Builder

	// TypeCast instr
	TypeCast(from, to reflect.Type) Builder

	// TypeAssert instr
	TypeAssert(from, to reflect.Type, twoValue bool) Builder

	// Zero instr
	Zero(typ reflect.Type) Builder

	// New instr
	New(typ reflect.Type) Builder

	// StartStmt emit a `StartStmt` event.
	StartStmt(stmt interface{}) interface{}

	// EndStmt emit a `EndStmt` event.
	EndStmt(stmt, start interface{}) Builder

	// Reserve reserves an instruction.
	Reserve() Reserved

	// ReservedAsPush sets Reserved as Push(v)
	ReservedAsPush(r Reserved, v interface{})

	// ReserveOpLsh reserves an instruction.
	ReserveOpLsh() Reserved

	// ReservedAsOpLsh sets Reserved as OpLsh
	ReservedAsOpLsh(r Reserved, kind Kind, op Operator)

	// GetPackage returns the Go+ package that the Builder works for.
	GetPackage() Package

	// Resolve resolves all unresolved labels/functions/consts/etc.
	Resolve() Code

	// DefineBlock instr
	DefineBlock() Builder

	// EndBlock instr
	EndBlock() Builder

	// Send instr
	Send() Builder

	// Recv instr
	Recv() Builder
}

// Package represents a Go+ package.
type Package interface {
	// NewVar creates a variable instance.
	NewVar(typ reflect.Type, name string) Var

	// NewLabel creates a label object.
	NewLabel(name string) Label

	// NewForPhrase creates a new ForPhrase instance.
	NewForPhrase(in reflect.Type) ForPhrase

	// NewComprehension creates a new Comprehension instance.
	NewComprehension(out reflect.Type) Comprehension

	// NewFunc creates a Go+ function.
	NewFunc(name string, nestDepth uint32, funcType ...int) FuncInfo

	// FindGoPackage lookups a Go package by pkgPath. It returns nil if not found.
	FindGoPackage(pkgPath string) GoPackage

	// GetGoFuncType returns a Go function's type.
	GetGoFuncType(addr GoFuncAddr) reflect.Type

	// GetGoFuncvType returns a Go function's type.
	GetGoFuncvType(addr GoFuncvAddr) reflect.Type

	// GetGoFuncInfo returns a Go function's information.
	GetGoFuncInfo(addr GoFuncAddr) *GoFuncInfo

	// GetGoFuncvInfo returns a Go function's information.
	GetGoFuncvInfo(addr GoFuncvAddr) *GoFuncInfo

	// GetGoVarInfo returns a Go variable's information.
	GetGoVarInfo(addr GoVarAddr) *GoVarInfo
}

// -----------------------------------------------------------------------------
