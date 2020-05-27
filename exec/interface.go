package exec

import (
	"reflect"

	"github.com/qiniu/qlang/v6/exec.spec"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

type iFuncInfo FuncInfo

// Name returns the function name.
func (p *iFuncInfo) Name() string {
	return p.name
}

// Type returns type of this function.
func (p *iFuncInfo) Type() reflect.Type {
	return ((*FuncInfo)(p)).Type()
}

// NumOut returns a function type's output parameter count.
// It panics if the type's Kind is not Func.
func (p *iFuncInfo) NumOut() int {
	return p.numOut
}

// Out returns the type of a function type's i'th output parameter.
// It panics if i is not in the range [0, NumOut()).
func (p *iFuncInfo) Out(i int) exec.Var {
	return ((*FuncInfo)(p)).Out(i)
}

// Args sets argument types of a qlang function.
func (p *iFuncInfo) Args(in ...reflect.Type) exec.FuncInfo {
	((*FuncInfo)(p)).Args(in...)
	return p
}

// Vargs sets argument types of a variadic qlang function.
func (p *iFuncInfo) Vargs(in ...reflect.Type) exec.FuncInfo {
	((*FuncInfo)(p)).Vargs(in...)
	return p
}

// Return sets return types of a qlang function.
func (p *iFuncInfo) Return(out ...exec.Var) exec.FuncInfo {
	if p.vlist != nil {
		log.Panicln("don't call DefineVar before calling Return.")
	}
	for _, v := range out {
		p.addVars(v.(*Var))
	}
	p.numOut = len(out)
	return p
}

// IsUnnamedOut returns if function results unnamed or not.
func (p *iFuncInfo) IsUnnamedOut() bool {
	return ((*FuncInfo)(p)).IsUnnamedOut()
}

// IsVariadic returns if this function is variadic or not.
func (p *iFuncInfo) IsVariadic() bool {
	return ((*FuncInfo)(p)).IsVariadic()
}

// -----------------------------------------------------------------------------

// Interface represents all global functions of a executing byte code generator.
type interfaceImpl struct{}

var defaultImpl = &interfaceImpl{}

// NewVar creates a variable instance.
func (p *interfaceImpl) NewVar(typ reflect.Type, name string) exec.Var {
	return NewVar(typ, name)
}

// NewLabel creates a label object.
func (p *interfaceImpl) NewLabel(name string) exec.Label {
	return NewLabel(name)
}

// NewForPhrase creates a new ForPhrase instance.
func (p *interfaceImpl) NewForPhrase(in reflect.Type) exec.ForPhrase {
	return NewForPhrase(in)
}

// NewComprehension creates a new Comprehension instance.
func (p *interfaceImpl) NewComprehension(out reflect.Type) exec.Comprehension {
	return NewComprehension(out)
}

// NewFunc create a qlang function.
func (p *interfaceImpl) NewFunc(name string, nestDepth uint32) exec.FuncInfo {
	return (*iFuncInfo)(NewFunc(name, nestDepth))
}

// FindGoPackage lookups a Go package by pkgPath. It returns nil if not found.
func (p *interfaceImpl) FindGoPackage(pkgPath string) exec.GoPackage {
	return FindGoPackage(pkgPath)
}

// GetGoFuncType returns a Go function's type.
func (p *interfaceImpl) GetGoFuncType(addr exec.GoFuncAddr) reflect.Type {
	panic("todo")
}

// GetGoFuncvType returns a Go function's type.
func (p *interfaceImpl) GetGoFuncvType(addr exec.GoFuncvAddr) reflect.Type {
	panic("todo")
}

// -----------------------------------------------------------------------------

type iBuilder Builder

// Push instr
func (p *iBuilder) Push(val interface{}) exec.Builder {
	((*Builder)(p)).Push(val)
	return p
}

// Pop instr
func (p *iBuilder) Pop(n int) exec.Builder {
	((*Builder)(p)).Pop(n)
	return p
}

// BuiltinOp instr
func (p *iBuilder) BuiltinOp(kind exec.Kind, op exec.Operator) exec.Builder {
	((*Builder)(p)).BuiltinOp(kind, op)
	return p
}

// Label defines a label to jmp here.
func (p *iBuilder) Label(l exec.Label) exec.Builder {
	((*Builder)(p)).Label(l.(*Label))
	return p
}

// Jmp instr
func (p *iBuilder) Jmp(l exec.Label) exec.Builder {
	((*Builder)(p)).Jmp(l.(*Label))
	return p
}

// JmpIf instr
func (p *iBuilder) JmpIf(zeroOrOne uint32, l exec.Label) exec.Builder {
	((*Builder)(p)).JmpIf(zeroOrOne, l.(*Label))
	return p
}

// CaseNE instr
func (p *iBuilder) CaseNE(l exec.Label, arity int) exec.Builder {
	((*Builder)(p)).CaseNE(l.(*Label), arity)
	return p
}

// Default instr
func (p *iBuilder) Default() exec.Builder {
	((*Builder)(p)).Default()
	return p
}

// ForPhrase instr
func (p *iBuilder) ForPhrase(f exec.ForPhrase, key, val exec.Var, hasExecCtx ...bool) exec.Builder {
	((*Builder)(p)).ForPhrase(f.(*ForPhrase), key.(*Var), val.(*Var), hasExecCtx...)
	return p
}

// FilterForPhrase instr
func (p *iBuilder) FilterForPhrase(f exec.ForPhrase) exec.Builder {
	((*Builder)(p)).FilterForPhrase(f.(*ForPhrase))
	return p
}

// EndForPhrase instr
func (p *iBuilder) EndForPhrase(f exec.ForPhrase) exec.Builder {
	((*Builder)(p)).EndForPhrase(f.(*ForPhrase))
	return p
}

// ListComprehension instr
func (p *iBuilder) ListComprehension(c exec.Comprehension) exec.Builder {
	((*Builder)(p)).ListComprehension(c.(*Comprehension))
	return p
}

// MapComprehension instr
func (p *iBuilder) MapComprehension(c exec.Comprehension) exec.Builder {
	((*Builder)(p)).MapComprehension(c.(*Comprehension))
	return p
}

// EndComprehension instr
func (p *iBuilder) EndComprehension(c exec.Comprehension) exec.Builder {
	((*Builder)(p)).EndComprehension(c.(*Comprehension))
	return p
}

// Closure instr
func (p *iBuilder) Closure(fun exec.FuncInfo) exec.Builder {
	((*Builder)(p)).Closure((*FuncInfo)(fun.(*iFuncInfo)))
	return p
}

// GoClosure instr
func (p *iBuilder) GoClosure(fun exec.FuncInfo) exec.Builder {
	((*Builder)(p)).GoClosure((*FuncInfo)(fun.(*iFuncInfo)))
	return p
}

// CallClosure instr
func (p *iBuilder) CallClosure(arity int) exec.Builder {
	((*Builder)(p)).CallClosure(arity)
	return p
}

// CallGoClosure instr
func (p *iBuilder) CallGoClosure(arity int) exec.Builder {
	((*Builder)(p)).CallGoClosure(arity)
	return p
}

// CallFunc instr
func (p *iBuilder) CallFunc(fun exec.FuncInfo) exec.Builder {
	((*Builder)(p)).CallFunc((*FuncInfo)(fun.(*iFuncInfo)))
	return p
}

// CallFuncv instr
func (p *iBuilder) CallFuncv(fun exec.FuncInfo, arity int) exec.Builder {
	((*Builder)(p)).CallFuncv((*FuncInfo)(fun.(*iFuncInfo)), arity)
	return p
}

// CallGoFunc instr
func (p *iBuilder) CallGoFunc(fun exec.GoFuncAddr) exec.Builder {
	((*Builder)(p)).CallGoFunc(fun)
	return p
}

// CallGoFuncv instr
func (p *iBuilder) CallGoFuncv(fun exec.GoFuncvAddr, arity int) exec.Builder {
	((*Builder)(p)).CallGoFuncv(fun, arity)
	return p
}

// DefineFunc instr
func (p *iBuilder) DefineFunc(fun exec.FuncInfo) exec.Builder {
	((*Builder)(p)).DefineFunc((*FuncInfo)(fun.(*iFuncInfo)))
	return p
}

// Return instr
func (p *iBuilder) Return(n int32) exec.Builder {
	((*Builder)(p)).Return(n)
	return p
}

// Load instr
func (p *iBuilder) Load(idx int32) exec.Builder {
	((*Builder)(p)).Load(idx)
	return p
}

// Store instr
func (p *iBuilder) Store(idx int32) exec.Builder {
	((*Builder)(p)).Store(idx)
	return p
}

// EndFunc instr
func (p *iBuilder) EndFunc(fun exec.FuncInfo) exec.Builder {
	((*Builder)(p)).EndFunc((*FuncInfo)(fun.(*iFuncInfo)))
	return p
}

// DefineVar defines variables.
func (p *iBuilder) DefineVar(vars ...exec.Var) exec.Builder {
	panic("todo")
}

// InCurrentCtx returns if a variable is in current context or not.
func (p *iBuilder) InCurrentCtx(v exec.Var) bool {
	panic("todo")
}

// LoadVar instr
func (p *iBuilder) LoadVar(v exec.Var) exec.Builder {
	panic("todo")
}

// StoreVar instr
func (p *iBuilder) StoreVar(v exec.Var) exec.Builder {
	panic("todo")
}

// AddrVar instr
func (p *iBuilder) AddrVar(v exec.Var) exec.Builder {
	panic("todo")
}

// AddrOp instr
func (p *iBuilder) AddrOp(kind exec.Kind, op exec.AddrOperator) exec.Builder {
	panic("todo")
}

// Append instr
func (p *iBuilder) Append(typ reflect.Type, arity int) exec.Builder {
	panic("todo")
}

// MakeArray instr
func (p *iBuilder) MakeArray(typ reflect.Type, arity int) exec.Builder {
	panic("todo")
}

// MakeMap instr
func (p *iBuilder) MakeMap(typ reflect.Type, arity int) exec.Builder {
	panic("todo")
}

// Make instr
func (p *iBuilder) Make(typ reflect.Type, arity int) exec.Builder {
	panic("todo")
}

// MapIndex instr
func (p *iBuilder) MapIndex() exec.Builder {
	panic("todo")
}

// SetMapIndex instr
func (p *iBuilder) SetMapIndex() exec.Builder {
	panic("todo")
}

// Index instr
func (p *iBuilder) Index(idx int) exec.Builder {
	panic("todo")
}

// SetIndex instr
func (p *iBuilder) SetIndex(idx int) exec.Builder {
	panic("todo")
}

// Slice instr
func (p *iBuilder) Slice(i, j int) exec.Builder {
	panic("todo")
}

// Slice3 instr
func (p *iBuilder) Slice3(i, j, k int) exec.Builder {
	panic("todo")
}

// TypeCast instr
func (p *iBuilder) TypeCast(from, to reflect.Type) exec.Builder {
	panic("todo")
}

// Zero instr
func (p *iBuilder) Zero(typ reflect.Type) exec.Builder {
	panic("todo")
}

// Reserve reserves an instruction.
func (p *iBuilder) Reserve() exec.Reserved {
	panic("todo")
}

// ReservedPush sets Reserved as Push(v)
func (p *iBuilder) ReservedPush(r exec.Reserved, v interface{}) {
	panic("todo")
}

// GlobalInterface returns the global Interface.
func (p *iBuilder) GlobalInterface() exec.Interface {
	panic("todo")
}

// -----------------------------------------------------------------------------
