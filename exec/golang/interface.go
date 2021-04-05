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
	"reflect"

	"github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/exec/bytecode"
	"github.com/qiniu/x/errors"
)

var qexecImpl = bytecode.NewPackage(nil)

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

// NewFunc create a Go+ function.
func (p *interfaceImpl) NewFunc(name string, nestDepth uint32, funcType ...int) exec.FuncInfo {
	if nestDepth == 0 {
		return nil
	}
	return NewFunc(name, nestDepth, funcType[0])
}

// FindGoPackage lookups a Go package by pkgPath. It returns nil if not found.
func (p *interfaceImpl) FindGoPackage(pkgPath string) exec.GoPackage {
	return qexecImpl.FindGoPackage(pkgPath)
}

// GetGoFuncType returns a Go function's type.
func (p *interfaceImpl) GetGoFuncType(addr exec.GoFuncAddr) reflect.Type {
	return qexecImpl.GetGoFuncType(addr)
}

// GetGoFuncvType returns a Go function's type.
func (p *interfaceImpl) GetGoFuncvType(addr exec.GoFuncvAddr) reflect.Type {
	return qexecImpl.GetGoFuncvType(addr)
}

// GetGoFuncInfo returns a Go function's information.
func (p *interfaceImpl) GetGoFuncInfo(addr exec.GoFuncAddr) *exec.GoFuncInfo {
	return qexecImpl.GetGoFuncInfo(addr)
}

// GetGoFuncvInfo returns a Go function's information.
func (p *interfaceImpl) GetGoFuncvInfo(addr exec.GoFuncvAddr) *exec.GoFuncInfo {
	return qexecImpl.GetGoFuncvInfo(addr)
}

// GetGoVarInfo returns a Go variable's information.
func (p *interfaceImpl) GetGoVarInfo(addr exec.GoVarAddr) *exec.GoVarInfo {
	return qexecImpl.GetGoVarInfo(addr)
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

// Defer instr
func (p *iBuilder) Defer() exec.Builder {
	((*Builder)(p)).Defer()
	return p
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

// Go instr
func (p *iBuilder) Go() exec.Builder {
	((*Builder)(p)).Go()
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
	((*Builder)(p)).Closure(fun.(*FuncInfo))
	return p
}

// GoClosure instr
func (p *iBuilder) GoClosure(fun exec.FuncInfo) exec.Builder {
	((*Builder)(p)).Closure(fun.(*FuncInfo))
	return p
}

// CallClosure instr
func (p *iBuilder) CallClosure(nexpr, arity int, ellipsis bool) exec.Builder {
	((*Builder)(p)).Call(nexpr, ellipsis)
	return p
}

// CallGoClosure instr
func (p *iBuilder) CallGoClosure(nexpr, arity int, ellipsis bool) exec.Builder {
	((*Builder)(p)).Call(nexpr, ellipsis)
	return p
}

// CallFunc instr
func (p *iBuilder) CallFunc(fun exec.FuncInfo, nexpr int) exec.Builder {
	((*Builder)(p)).CallFunc(fun.(*FuncInfo), nexpr)
	return p
}

// CallFuncv instr
func (p *iBuilder) CallFuncv(fun exec.FuncInfo, nexpr, arity int) exec.Builder {
	((*Builder)(p)).CallFuncv(fun.(*FuncInfo), nexpr, arity)
	return p
}

// CallGoFunc instr
func (p *iBuilder) CallGoFunc(fun exec.GoFuncAddr, nexpr int) exec.Builder {
	((*Builder)(p)).CallGoFunc(fun, nexpr)
	return p
}

// CallGoFuncv instr
func (p *iBuilder) CallGoFuncv(fun exec.GoFuncvAddr, nexpr, arity int) exec.Builder {
	((*Builder)(p)).CallGoFuncv(fun, nexpr, arity)
	return p
}

// DefineFunc instr
func (p *iBuilder) DefineFunc(fun exec.FuncInfo) exec.Builder {
	((*Builder)(p)).DefineFunc(fun)
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
	((*Builder)(p)).EndFunc(fun.(*FuncInfo))
	return p
}

// DefineVar defines variables.
func (p *iBuilder) DefineVar(vars ...exec.Var) exec.Builder {
	((*Builder)(p)).DefineVar(vars...)
	return p
}

// DefineType defines variables.
func (p *iBuilder) DefineType(typ reflect.Type, name string) exec.Builder {
	((*Builder)(p)).DefineType(typ, name)
	return p
}

// InCurrentCtx returns if a variable is in current context or not.
func (p *iBuilder) InCurrentCtx(v exec.Var) bool {
	return ((*Builder)(p)).InCurrentCtx(v.(*Var))
}

// LoadVar instr
func (p *iBuilder) LoadVar(v exec.Var) exec.Builder {
	((*Builder)(p)).LoadVar(v)
	return p
}

// StoreVar instr
func (p *iBuilder) StoreVar(v exec.Var) exec.Builder {
	((*Builder)(p)).StoreVar(v)
	return p
}

// AddrVar instr
func (p *iBuilder) AddrVar(v exec.Var) exec.Builder {
	((*Builder)(p)).AddrVar(v.(*Var))
	return p
}

// LoadGoVar instr
func (p *iBuilder) LoadGoVar(addr exec.GoVarAddr) exec.Builder {
	((*Builder)(p)).LoadGoVar(addr)
	return p
}

// StoreGoVar instr
func (p *iBuilder) StoreGoVar(addr exec.GoVarAddr) exec.Builder {
	((*Builder)(p)).StoreGoVar(addr)
	return p
}

// AddrGoVar instr
func (p *iBuilder) AddrGoVar(addr exec.GoVarAddr) exec.Builder {
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

// New instr
func (p *iBuilder) New(typ reflect.Type) exec.Builder {
	((*Builder)(p)).New(typ)
	return p
}

// Zero instr
func (p *iBuilder) Zero(typ reflect.Type) exec.Builder {
	((*Builder)(p)).Zero(typ)
	return p
}

// StartStmt receives a `StartStmt` event.
func (p *iBuilder) StartStmt(stmt interface{}) interface{} {
	return ((*Builder)(p)).StartStmt(stmt)
}

// EndStmt receives a `EndStmt` event.
func (p *iBuilder) EndStmt(stmt, start interface{}) exec.Builder {
	((*Builder)(p)).EndStmt(stmt, start)
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
	return ((*Builder)(p)).ReserveOpLsh()
}

// ReservedAsOpLsh sets Reserved as GoBuiltin
func (p *iBuilder) ReservedAsOpLsh(r exec.Reserved, kind exec.Kind, op exec.Operator) {
	((*Builder)(p)).ReservedAsOpLsh(r, kind, op)
}

// GetPackage returns the Go+ package that the Builder works for.
func (p *iBuilder) GetPackage() exec.Package {
	return defaultImpl
}

// Resolve resolves all unresolved labels/functions/consts/etc.
func (p *iBuilder) Resolve() exec.Code {
	return ((*Builder)(p)).Resolve()
}

// DefineBlock
func (p *iBuilder) DefineBlock() exec.Builder {
	((*Builder)(p)).DefineBlock()
	return p
}

// EndBlock
func (p *iBuilder) EndBlock() exec.Builder {
	((*Builder)(p)).EndBlock()
	return p
}

// -----------------------------------------------------------------------------
