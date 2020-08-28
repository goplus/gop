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
	"go/ast"
	"go/token"
	"strconv"

	"github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/exec/golang/internal/go/printer"
	"github.com/goplus/gop/reflect"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

const (
	nVariadicInvalid      = 0
	nVariadicFixedArgs    = 1
	nVariadicVariadicArgs = 2
)

// FuncInfo represents a Go+ function information.
type FuncInfo struct {
	isMethod int
	name     string
	closure  *printer.ReservedExpr // only when name="" (closure)
	recv     reflect.Type
	t        reflect.Type
	in       []reflect.Type
	out      []exec.Var
	scopeCtx
	nVariadic uint16
}

// NewFunc create a Go+ function.
func NewFunc(name string, nestDepth uint32, isMethod int) *FuncInfo {
	if name != "" {
		return &FuncInfo{name: name, isMethod: isMethod}
	}
	return &FuncInfo{closure: &printer.ReservedExpr{}}
}

func (p *FuncInfo) getFuncExpr(recv ast.Expr) ast.Expr {
	if p.isMethod == 1 {
		if p.name != "" {
			if unary, ok := recv.(*ast.UnaryExpr); ok {
				recv = unary.X
			}
			return &ast.SelectorExpr{
				X:   recv,
				Sel: Ident(p.name),
			}
		}
	}
	if p.name != "" {
		return Ident(p.name)
	}
	return p.closure
}

// Name returns the function name.
func (p *FuncInfo) Name() string {
	return p.name
}

// Type returns type of this function.
func (p *FuncInfo) Type() reflect.Type {
	if p.t == nil {
		out := make([]reflect.Type, len(p.out))
		for i, v := range p.out {
			out[i] = v.(*Var).typ
		}
		in := make([]reflect.Type, 0, p.NumIn())
		if p.isMethod == 1 {
			in = append(in, p.recv)
		}
		in = append(in, p.in...)
		p.t = reflect.FuncOf(in, out, p.IsVariadic())
	}
	return p.t
}

// NumIn returns a function's input parameter count.
func (p *FuncInfo) NumIn() int {
	return len(p.in) + p.isMethod
}

// NumOut returns a function's output parameter count.
func (p *FuncInfo) NumOut() int {
	return len(p.out)
}

// Out returns the type of a function type's i'th output parameter.
// It panics if i is not in the range [0, NumOut()).
func (p *FuncInfo) Out(i int) exec.Var {
	return p.out[i]
}

// Args sets argument types of a Go+ function.
func (p *FuncInfo) Args(in ...reflect.Type) exec.FuncInfo {
	if p.isMethod == 1 {
		p.recv = in[0]
		in = in[1:]
	}
	p.in = in
	p.setVariadic(nVariadicFixedArgs)
	return p
}

// Vargs sets argument types of a variadic Go+ function.
func (p *FuncInfo) Vargs(in ...reflect.Type) exec.FuncInfo {
	if in[len(in)-1].Kind() != reflect.Slice {
		log.Panicln("Vargs failed: last argument must be a slice.")
	}
	if p.isMethod == 1 {
		p.recv = in[0]
		in = in[1:]
	}
	p.in = in
	p.setVariadic(nVariadicVariadicArgs)
	return p
}

// Return sets return types of a Go+ function.
func (p *FuncInfo) Return(out ...exec.Var) exec.FuncInfo {
	p.out = out
	return p
}

// IsUnnamedOut returns if function results unnamed or not.
func (p *FuncInfo) IsUnnamedOut() bool {
	if len(p.out) > 0 {
		return p.out[0].IsUnnamedOut()
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
	p.rhs.Push(fun.getFuncExpr(nil))
	return p
}

// CallFunc instr
func (p *Builder) CallFunc(fun *FuncInfo, nexpr int) *Builder {
	if fun.isMethod == 1 {
		args := p.rhs.GetArgs(len(fun.in) + 1)
		recv := args[0]
		args = args[1:]
		p.rhs.Ret(len(fun.in)+1, args...)
		p.rhs.Push(fun.getFuncExpr(recv.(ast.Expr)))
	} else {
		p.rhs.Push(fun.getFuncExpr(nil))
	}
	return p.Call(nexpr, false)
}

// CallFuncv instr
func (p *Builder) CallFuncv(fun *FuncInfo, nexpr, arity int) *Builder {
	if fun.isMethod == 1 {
		args := p.rhs.GetArgs(arity)
		recv := args[0]
		args = args[1:]
		p.rhs.Ret(arity, args...)
		p.rhs.Push(fun.getFuncExpr(recv.(ast.Expr)))
	} else {
		p.rhs.Push(fun.getFuncExpr(nil))
	}
	return p.Call(nexpr, arity == -1)
}

// DefineFunc instr
func (p *Builder) DefineFunc(fun exec.FuncInfo) *Builder {
	f := fun.(*FuncInfo)
	f.initStmts()
	p.scopeCtx = &f.scopeCtx
	p.cfun = f
	return p
}

// Return instr
func (p *Builder) Return(n int32) *Builder {
	var results []ast.Expr
	var stmt ast.Stmt
	switch n {
	case exec.BreakAsReturn:
		stmt = &ast.BranchStmt{
			Tok: token.BREAK,
		}
	case exec.ContinueAsReturn:
		stmt = &ast.BranchStmt{
			Tok: token.CONTINUE,
		}
	default:
		if n > 0 {
			arity := int(n)
			args := p.rhs.GetArgs(arity)
			results = make([]ast.Expr, n)
			for i, arg := range args {
				results[i] = arg.(ast.Expr)
			}
			p.rhs.PopN(arity)
		}
		stmt = &ast.ReturnStmt{Results: results}
	}
	p.rhs.Push(stmt)
	return p
}

// EndFunc instr
func (p *Builder) EndFunc(fun *FuncInfo) *Builder {
	p.endBlockStmt(1)
	body := &ast.BlockStmt{List: fun.getStmts(p)}
	name := fun.name
	if name != "" {
		fn := &ast.FuncDecl{
			Name: Ident(name),
			Type: toFuncType(p, fun),
			Body: body,
		}
		if fun.isMethod == 1 {
			params := make([]*ast.Field, 1)
			params[0] = Field(p, "recv", fun.recv, "", false)
			fn.Recv = &ast.FieldList{Opening: 1, List: params, Closing: 1}
		}
		p.gblDecls = append(p.gblDecls, fn)
	} else {
		fun.closure.Expr = &ast.FuncLit{
			Type: toFuncType(p, fun),
			Body: body,
		}
	}
	p.cfun = nil
	p.scopeCtx = &p.gblScope
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
