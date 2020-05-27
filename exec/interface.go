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
type interfaceImpl struct {
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
	return nil
}

// GetGoFuncvType returns a Go function's type.
func (p *interfaceImpl) GetGoFuncvType(addr exec.GoFuncvAddr) reflect.Type {
	return nil
}

// -----------------------------------------------------------------------------
