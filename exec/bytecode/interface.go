package bytecode

import (
	"reflect"

	exec "github.com/qiniu/goplus/exec.spec"
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

// Args sets argument types of a Go+ function.
func (p *iFuncInfo) Args(in ...reflect.Type) exec.FuncInfo {
	((*FuncInfo)(p)).Args(in...)
	return p
}

// Vargs sets argument types of a variadic Go+ function.
func (p *iFuncInfo) Vargs(in ...reflect.Type) exec.FuncInfo {
	((*FuncInfo)(p)).Vargs(in...)
	return p
}

// Return sets return types of a Go+ function.
func (p *iFuncInfo) Return(out ...exec.Var) exec.FuncInfo {
	if p.vlist != nil {
		log.Panicln("don't call DefineVar before calling Return.")
	}
	p.addVars(out...)
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

// GlobalInterface returns the global Interface.
func GlobalInterface() exec.Interface {
	return defaultImpl
}

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

// NewFunc create a Go+ function.
func (p *interfaceImpl) NewFunc(name string, nestDepth uint32) exec.FuncInfo {
	return (*iFuncInfo)(NewFunc(name, nestDepth))
}

// FindGoPackage lookups a Go package by pkgPath. It returns nil if not found.
func (p *interfaceImpl) FindGoPackage(pkgPath string) exec.GoPackage {
	return FindGoPackage(pkgPath)
}

// GetGoFuncType returns a Go function's type.
func (p *interfaceImpl) GetGoFuncType(addr exec.GoFuncAddr) reflect.Type {
	return reflect.TypeOf(gofuns[addr].This)
}

// GetGoFuncvType returns a Go function's type.
func (p *interfaceImpl) GetGoFuncvType(addr exec.GoFuncvAddr) reflect.Type {
	return reflect.TypeOf(gofunvs[addr].This)
}

// GetGoFuncInfo returns a Go function's information.
func (p *interfaceImpl) GetGoFuncInfo(addr exec.GoFuncAddr) *exec.GoFuncInfo {
	gfi := &gofuns[addr]
	return &exec.GoFuncInfo{Pkg: gfi.Pkg, Name: gfi.Name, This: gfi.This}
}

// GetGoFuncvInfo returns a Go function's information.
func (p *interfaceImpl) GetGoFuncvInfo(addr exec.GoFuncvAddr) *exec.GoFuncInfo {
	gfi := &gofunvs[addr]
	return &exec.GoFuncInfo{Pkg: gfi.Pkg, Name: gfi.Name, This: gfi.This}
}

// GetGoVarInfo returns a Go variable's information.
func (p *interfaceImpl) GetGoVarInfo(addr exec.GoVarAddr) *exec.GoVarInfo {
	gvi := &govars[addr]
	return &exec.GoVarInfo{Pkg: gvi.Pkg, Name: gvi.Name, This: gvi.Addr}
}

// -----------------------------------------------------------------------------

type iBuilder Builder

// Interface converts *Builder to exec.Builder interface.
func (p *Builder) Interface() exec.Builder {
	return (*iBuilder)(p)
}

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
	((*Builder)(p)).ForPhrase(f.(*ForPhrase), toVar(key), toVar(val), hasExecCtx...)
	return p
}

func toVar(v exec.Var) *Var {
	if v == nil {
		return nil
	}
	return v.(*Var)
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
func (p *iBuilder) CallClosure(nexpr, arity int, ellipsis bool) exec.Builder {
	if ellipsis {
		arity = -1
	}
	((*Builder)(p)).CallClosure(arity)
	return p
}

// CallGoClosure instr
func (p *iBuilder) CallGoClosure(nexpr, arity int, ellipsis bool) exec.Builder {
	if ellipsis {
		arity = -1
	}
	((*Builder)(p)).CallGoClosure(arity)
	return p
}

// CallFunc instr
func (p *iBuilder) CallFunc(fun exec.FuncInfo, nexpr int) exec.Builder {
	((*Builder)(p)).CallFunc((*FuncInfo)(fun.(*iFuncInfo)))
	return p
}

// CallFuncv instr
func (p *iBuilder) CallFuncv(fun exec.FuncInfo, nexpr, arity int) exec.Builder {
	((*Builder)(p)).CallFuncv((*FuncInfo)(fun.(*iFuncInfo)), arity)
	return p
}

// CallGoFunc instr
func (p *iBuilder) CallGoFunc(fun exec.GoFuncAddr, nexpr int) exec.Builder {
	((*Builder)(p)).CallGoFunc(fun)
	return p
}

// CallGoFuncv instr
func (p *iBuilder) CallGoFuncv(fun exec.GoFuncvAddr, nexpr, arity int) exec.Builder {
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
	((*Builder)(p)).DefineVars(vars...)
	return p
}

// InCurrentCtx returns if a variable is in current context or not.
func (p *iBuilder) InCurrentCtx(v exec.Var) bool {
	return ((*Builder)(p)).InCurrentCtx(v.(*Var))
}

// LoadVar instr
func (p *iBuilder) LoadVar(v exec.Var) exec.Builder {
	((*Builder)(p)).LoadVar(v.(*Var))
	return p
}

// StoreVar instr
func (p *iBuilder) StoreVar(v exec.Var) exec.Builder {
	((*Builder)(p)).StoreVar(v.(*Var))
	return p
}

// AddrVar instr
func (p *iBuilder) AddrVar(v exec.Var) exec.Builder {
	((*Builder)(p)).AddrVar(v.(*Var))
	return p
}

// LoadGoVar instr
func (p *iBuilder) LoadGoVar(addr GoVarAddr) exec.Builder {
	((*Builder)(p)).LoadGoVar(addr)
	return p
}

// StoreGoVar instr
func (p *iBuilder) StoreGoVar(addr GoVarAddr) exec.Builder {
	((*Builder)(p)).StoreGoVar(addr)
	return p
}

// AddrGoVar instr
func (p *iBuilder) AddrGoVar(addr GoVarAddr) exec.Builder {
	((*Builder)(p)).AddrGoVar(addr)
	return p
}

// AddrOp instr
func (p *iBuilder) AddrOp(kind exec.Kind, op exec.AddrOperator) exec.Builder {
	((*Builder)(p)).AddrOp(kind, op)
	return p
}

// Append instr
func (p *iBuilder) Append(typ reflect.Type, arity int) exec.Builder {
	((*Builder)(p)).Append(typ, arity)
	return p
}

// MakeArray instr
func (p *iBuilder) MakeArray(typ reflect.Type, arity int) exec.Builder {
	((*Builder)(p)).MakeArray(typ, arity)
	return p
}

// MakeMap instr
func (p *iBuilder) MakeMap(typ reflect.Type, arity int) exec.Builder {
	((*Builder)(p)).MakeMap(typ, arity)
	return p
}

// Make instr
func (p *iBuilder) Make(typ reflect.Type, arity int) exec.Builder {
	((*Builder)(p)).Make(typ, arity)
	return p
}

// MapIndex instr
func (p *iBuilder) MapIndex() exec.Builder {
	((*Builder)(p)).MapIndex()
	return p
}

// SetMapIndex instr
func (p *iBuilder) SetMapIndex() exec.Builder {
	((*Builder)(p)).SetMapIndex()
	return p
}

// Index instr
func (p *iBuilder) Index(idx int) exec.Builder {
	((*Builder)(p)).Index(idx)
	return p
}

// SetIndex instr
func (p *iBuilder) SetIndex(idx int) exec.Builder {
	((*Builder)(p)).SetIndex(idx)
	return p
}

// Slice instr
func (p *iBuilder) Slice(i, j int) exec.Builder {
	((*Builder)(p)).Slice(i, j)
	return p
}

// Slice3 instr
func (p *iBuilder) Slice3(i, j, k int) exec.Builder {
	((*Builder)(p)).Slice3(i, j, k)
	return p
}

// TypeCast instr
func (p *iBuilder) TypeCast(from, to reflect.Type) exec.Builder {
	((*Builder)(p)).TypeCast(from, to)
	return p
}

// GoBuiltin instr
func (p *iBuilder) GoBuiltin(typ reflect.Type, op exec.GoBuiltin) exec.Builder {
	((*Builder)(p)).GoBuiltin(typ, op)
	return p
}

// Zero instr
func (p *iBuilder) Zero(typ reflect.Type) exec.Builder {
	((*Builder)(p)).Zero(typ)
	return p
}

// StartStmt recieves a `StartStmt` event.
func (p *iBuilder) StartStmt(stmt interface{}) interface{} {
	return nil
}

// EndStmt recieves a `EndStmt` event.
func (p *iBuilder) EndStmt(stmt, start interface{}) exec.Builder {
	return p
}

// Reserve reserves an instruction.
func (p *iBuilder) Reserve() exec.Reserved {
	return ((*Builder)(p)).Reserve()
}

// ReservedAsPush sets Reserved as Push(v)
func (p *iBuilder) ReservedAsPush(r exec.Reserved, v interface{}) {
	((*Builder)(p)).ReservedAsPush(r, v)
}

// GlobalInterface returns the global Interface.
func (p *iBuilder) GlobalInterface() exec.Interface {
	return defaultImpl
}

// Resolve resolves all unresolved labels/functions/consts/etc.
func (p *iBuilder) Resolve() exec.Code {
	return ((*Builder)(p)).Resolve()
}

// -----------------------------------------------------------------------------
