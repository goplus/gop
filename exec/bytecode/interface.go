package bytecode

import (
	"reflect"

	"github.com/goplus/gop/exec.spec"
	"github.com/qiniu/x/errors"
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

// NumIn returns a function's input parameter count.
func (p *iFuncInfo) NumIn() int {
	return ((*FuncInfo)(p)).NumIn()
}

// NumOut returns a function's output parameter count.
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
	p.vlist = nil
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
func (p *iBuilder) JmpIf(cond exec.JmpCondFlag, l exec.Label) exec.Builder {
	((*Builder)(p)).JmpIf(cond, l.(*Label))
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

// Defer instr
func (p *iBuilder) Defer() exec.Builder {
	((*Builder)(p)).Defer()
	return p
}

// Go instr
func (p *iBuilder) Go() exec.Builder {
	return ((*Builder)(p)).Go()
}

// Send instr
func (p *iBuilder) Send() exec.Builder {
	((*Builder)(p)).Send()
	return p
}

// Recv instr
func (p *iBuilder) Recv() exec.Builder {
	((*Builder)(p)).Recv()
	return p
}

// WrapIfErr instr
func (p *iBuilder) WrapIfErr(nret int, l exec.Label) exec.Builder {
	((*Builder)(p)).WrapIfErr(nret, l.(*Label))
	return p
}

// ErrWrap instr
func (p *iBuilder) ErrWrap(nret int, retErr exec.Var, frame *errors.Frame, narg int) exec.Builder {
	((*Builder)(p)).ErrWrap(nret, retErr, frame, narg)
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

// Addr instr
func (p *iBuilder) Addr(idx int32) exec.Builder {
	((*Builder)(p)).Addr(idx)
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

// DefineType name string,reflect.Type instr
func (p *iBuilder) DefineType(typ reflect.Type, name string) exec.Builder {
	// Nothing todo in bytecode
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

// LoadField instr
func (p *iBuilder) LoadField(typ reflect.Type, index []int) exec.Builder {
	((*Builder)(p)).LoadField(typ, index)
	return p
}

// AddrField instr
func (p *iBuilder) AddrField(typ reflect.Type, index []int) exec.Builder {
	((*Builder)(p)).AddrField(typ, index)
	return p
}

// StoreField instr
func (p *iBuilder) StoreField(typ reflect.Type, index []int) exec.Builder {
	((*Builder)(p)).StoreField(typ, index)
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
func (p *iBuilder) MapIndex(twoValue bool) exec.Builder {
	((*Builder)(p)).MapIndex(twoValue)
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

// AddrIndex instr
func (p *iBuilder) AddrIndex(idx int) exec.Builder {
	((*Builder)(p)).AddrIndex(idx)
	return p
}

// SetIndex instr
func (p *iBuilder) SetIndex(idx int) exec.Builder {
	((*Builder)(p)).SetIndex(idx)
	return p
}

// Struct instr
func (p *iBuilder) Struct(typ reflect.Type, arity int) exec.Builder {
	((*Builder)(p)).Struct(typ, arity)
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

// TypeAssert instr
func (p *iBuilder) TypeAssert(from, to reflect.Type, twoValue bool) exec.Builder {
	((*Builder)(p)).TypeAssert(from, to, twoValue)
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

// New instr
func (p *iBuilder) New(typ reflect.Type) exec.Builder {
	((*Builder)(p)).New(typ)
	return p
}

// StartStmt receives a `StartStmt` event.
func (p *iBuilder) StartStmt(stmt interface{}) interface{} {
	return nil
}

// EndStmt receives a `EndStmt` event.
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

// ReserveOpLsh reserves an instruction.
func (p *iBuilder) ReserveOpLsh() exec.Reserved {
	return ((*Builder)(p)).Reserve()
}

// ReservedAsOpLsh sets Reserved as GoBuiltin
func (p *iBuilder) ReservedAsOpLsh(r exec.Reserved, kind exec.Kind, op exec.Operator) {
	((*Builder)(p)).ReservedAsOpLsh(r, kind, op)
}

// GetPackage returns the Go+ package that the Builder works for.
func (p *iBuilder) GetPackage() exec.Package {
	return NewPackage(p.code)
}

// Resolve resolves all unresolved labels/functions/consts/etc.
func (p *iBuilder) Resolve() exec.Code {
	return ((*Builder)(p)).Resolve()
}

// DefineBlock
func (p *iBuilder) DefineBlock() exec.Builder {
	return p
}

// EndBlock
func (p *iBuilder) EndBlock() exec.Builder {
	return p
}

func (p *iBuilder) MethodOf(typ reflect.Type, infos []*exec.MethodInfo) exec.Builder {
	((*Builder)(p)).MethodOf(typ, infos)
	return p
}

// -----------------------------------------------------------------------------
