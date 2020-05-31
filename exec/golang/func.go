package golang

import (
	"go/ast"
	"go/token"
	"reflect"
	"strconv"

	"github.com/qiniu/qlang/v6/exec.spec"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

const (
	nVariadicInvalid      = 0
	nVariadicFixedArgs    = 1
	nVariadicVariadicArgs = 2
)

// FuncInfo represents a qlang function information.
type FuncInfo struct {
	name   string
	t      reflect.Type
	in     []reflect.Type
	numOut int
	stmts  []ast.Stmt
	varManager
	blockCtx
	nVariadic uint16
}

// NewFunc create a qlang function.
func NewFunc(name string, nestDepth uint32) *FuncInfo {
	return &FuncInfo{name: name}
}

func (p *FuncInfo) getName(b *Builder) string {
	if p.name == "" {
		p.name = b.autoIdent()
	}
	return p.name
}

// Name returns the function name.
func (p *FuncInfo) Name() string {
	return p.name
}

// Type returns type of this function.
func (p *FuncInfo) Type() reflect.Type {
	if p.t == nil {
		out := make([]reflect.Type, p.numOut)
		for i := 0; i < p.numOut; i++ {
			out[i] = p.vlist[i].(*Var).typ
		}
		p.t = reflect.FuncOf(p.in, out, p.IsVariadic())
	}
	return p.t
}

// NumOut returns a function type's output parameter count.
// It panics if the type's Kind is not Func.
func (p *FuncInfo) NumOut() int {
	return p.numOut
}

// Out returns the type of a function type's i'th output parameter.
// It panics if i is not in the range [0, NumOut()).
func (p *FuncInfo) Out(i int) exec.Var {
	if i >= p.numOut {
		log.Panicln("FuncInfo.Out: out of range -", i, "func:", p.name)
	}
	return p.vlist[i]
}

// Args sets argument types of a qlang function.
func (p *FuncInfo) Args(in ...reflect.Type) exec.FuncInfo {
	p.in = in
	p.setVariadic(nVariadicFixedArgs)
	return p
}

// Vargs sets argument types of a variadic qlang function.
func (p *FuncInfo) Vargs(in ...reflect.Type) exec.FuncInfo {
	if in[len(in)-1].Kind() != reflect.Slice {
		log.Panicln("Vargs failed: last argument must be a slice.")
	}
	p.in = in
	p.setVariadic(nVariadicVariadicArgs)
	return p
}

// Return sets return types of a qlang function.
func (p *FuncInfo) Return(out ...exec.Var) exec.FuncInfo {
	if p.vlist != nil {
		log.Panicln("don't call DefineVar before calling Return.")
	}
	p.addVar(out...)
	p.numOut = len(out)
	return p
}

// IsUnnamedOut returns if function results unnamed or not.
func (p *FuncInfo) IsUnnamedOut() bool {
	if p.numOut > 0 {
		return p.vlist[0].IsUnnamedOut()
	}
	return false
}

// IsVariadic returns if this function is variadic or not.
func (p *FuncInfo) IsVariadic() bool {
	if p.nVariadic == 0 {
		log.Panicln("FuncInfo is unintialized.")
	}
	return p.nVariadic == nVariadicVariadicArgs
}

func (p *FuncInfo) setVariadic(nVariadic uint16) {
	if p.nVariadic == 0 {
		p.nVariadic = nVariadic
	} else if p.nVariadic != nVariadic {
		log.Panicln("setVariadic failed: unmatched -", p.name)
	}
}

// -----------------------------------------------------------------------------

// Closure instr
func (p *Builder) Closure(fun *FuncInfo) *Builder {
	p.rhs.Push(Ident(fun.getName(p)))
	return p
}

// CallFunc instr
func (p *Builder) CallFunc(fun *FuncInfo) *Builder {
	p.rhs.Push(Ident(fun.getName(p)))
	return p.Call(len(fun.in), false)
}

// CallFuncv instr
func (p *Builder) CallFuncv(fun *FuncInfo, arity int) *Builder {
	var ellipsis bool
	if arity == -1 {
		arity = len(fun.in)
		ellipsis = true
	}
	p.rhs.Push(Ident(fun.getName(p)))
	return p.Call(arity, ellipsis)
}

// DefineFunc instr
func (p *Builder) DefineFunc(fun exec.FuncInfo) *Builder {
	f := fun.(*FuncInfo)
	f.saveEnv(p)
	p.varManager = &f.varManager
	p.stmts = &f.stmts
	p.cfun = f
	return p
}

// Return instr
func (p *Builder) Return(n int32) *Builder {
	var results []ast.Expr
	if n > 0 {
		arity := int(n)
		args := p.rhs.GetArgs(arity)
		results = make([]ast.Expr, n)
		for i, arg := range args {
			results[i] = arg.(ast.Expr)
		}
		p.rhs.PopN(arity)
	}
	p.rhs.Push(&ast.ReturnStmt{Results: results})
	return p
}

// EndFunc instr
func (p *Builder) EndFunc(fun *FuncInfo) *Builder {
	p.endBlockStmt()
	body := &ast.BlockStmt{List: fun.stmts}
	fn := &ast.FuncDecl{
		Name: Ident(fun.getName(p)),
		Type: toFuncType(p, fun),
		Body: body,
	}
	p.gbldecls = append(p.gbldecls, fn)
	p.cfun = nil
	fun.restoreEnv(p)
	return p
}

func toFuncType(p *Builder, typ *FuncInfo) *ast.FuncType {
	numIn, numOut := len(typ.in), typ.NumOut()
	variadic := typ.IsVariadic()
	var opening token.Pos
	var params, results []*ast.Field
	if numIn > 0 {
		params = make([]*ast.Field, numIn)
		if variadic {
			numIn--
		}
		for i := 0; i < numIn; i++ {
			params[i] = Field(p, toArg(i), typ.in[i], "", false)
		}
		if variadic {
			params[numIn] = Field(p, toArg(numIn), typ.in[numIn], "", true)
		}
	}
	if numOut > 0 {
		results = make([]*ast.Field, numOut)
		for i := 0; i < numOut; i++ {
			out := typ.Out(i).(*Var)
			results[i] = Field(p, out.name, out.typ, "", false)
		}
		opening++
	}
	return &ast.FuncType{
		Params:  &ast.FieldList{Opening: 1, List: params, Closing: 1},
		Results: &ast.FieldList{Opening: opening, List: results, Closing: opening},
	}
}

func toArg(i int) string {
	return "_arg_" + strconv.Itoa(i)
}
