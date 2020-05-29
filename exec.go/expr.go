/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

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

package exec

import (
	"go/ast"
	"go/token"
	"log"
	"reflect"
	"strconv"

	"github.com/qiniu/qlang/v6/exec.spec"
)

// -----------------------------------------------------------------------------

// Ident instr
func (p *Builder) Ident(name string) *Builder {
	p.code.Push(&ast.Ident{Name: name})
	return p
}

// Ident - ast.Ident
func Ident(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

// GoFuncIdent - ast.Ident or ast.SelectorExpr
func (p *Builder) GoFuncIdent(pkgPath, name string) ast.Expr {
	if pkgPath == "" {
		return Ident(name)
	}
	pkg := p.Import(pkgPath)
	return &ast.SelectorExpr{
		X:   Ident(pkg),
		Sel: Ident(name),
	}
}

// StringConst instr
func (p *Builder) StringConst(v string) *Builder {
	p.code.Push(StringConst(v))
	return p
}

// StringConst - ast.BasicLit
func StringConst(v string) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(v),
	}
}

// IntConst instr
func (p *Builder) IntConst(v int64) *Builder {
	p.code.Push(IntConst(v))
	return p
}

// IntConst - ast.BasicLit
func IntConst(v int64) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.FormatInt(v, 10),
	}
}

// UintConst instr
func (p *Builder) UintConst(v uint64) *Builder {
	p.code.Push(&ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.FormatUint(v, 10),
	})
	return p
}

// FloatConst instr
func (p *Builder) FloatConst(v float64) *Builder {
	p.code.Push(&ast.BasicLit{
		Kind:  token.FLOAT,
		Value: strconv.FormatFloat(v, 'g', -1, 64),
	})
	return p
}

// ImagConst instr
func (p *Builder) ImagConst(v float64) *Builder {
	p.code.Push(&ast.BasicLit{
		Kind:  token.IMAG,
		Value: strconv.FormatFloat(v, 'g', -1, 64) + "i",
	})
	return p
}

// ComplexConst instr
func (p *Builder) ComplexConst(v complex128) *Builder {
	r, i := real(v), imag(v)
	return p.FloatConst(r).ImagConst(i).BuiltinOp(exec.Float64, exec.OpAdd)
}

// Push instr
func (p *Builder) Push(val interface{}) *Builder {
	if val == nil {
		return p.Ident("nil")
	}
	v := reflect.ValueOf(val)
	kind := v.Kind()
	if kind == reflect.String {
		return p.StringConst(val.(string))
	}
	if kind >= reflect.Int && kind <= reflect.Int64 {
		p.IntConst(v.Int())
		if t := v.Type(); t != exec.TyInt {
			p.TypeCast(exec.TyInt, t)
		}
		return p
	}
	if kind >= reflect.Uint && kind <= reflect.Uintptr {
		p.UintConst(v.Uint())
		p.TypeCast(exec.TyInt, v.Type())
		return p
	}
	if kind >= reflect.Float32 && kind <= reflect.Float64 {
		p.FloatConst(v.Float())
		if t := v.Type(); t != exec.TyFloat64 {
			p.TypeCast(exec.TyFloat64, t)
		}
		return p
	}
	if kind >= reflect.Complex64 && kind <= reflect.Complex128 {
		p.ComplexConst(v.Complex())
		if t := v.Type(); t != exec.TyComplex128 {
			p.TypeCast(exec.TyComplex128, t)
		}
		return p
	}
	if kind == reflect.Bool {
		if val.(bool) {
			return p.Ident("true")
		}
		return p.Ident("false")
	}
	log.Panicln("Builder.Push: value type is unknown -", v.Type())
	return p
}

// UnaryOp instr
func (p *Builder) UnaryOp(tok token.Token) *Builder {
	x := p.code.Pop().(ast.Expr)
	p.code.Push(&ast.UnaryExpr{Op: tok, X: x})
	return p
}

// BinaryOp instr
func (p *Builder) BinaryOp(tok token.Token) *Builder {
	y := p.code.Pop().(ast.Expr)
	x := p.code.Pop().(ast.Expr)
	p.code.Push(&ast.BinaryExpr{Op: tok, X: x, Y: y})
	return p
}

// BuiltinOp instr
func (p *Builder) BuiltinOp(kind exec.Kind, op exec.Operator) *Builder {
	tok := opTokens[op]
	if tok == token.ILLEGAL {
		log.Panicln("BuiltinOp: unsupported op -", op)
	}
	oi := op.GetInfo()
	if oi.InSecond == 0 {
		return p.UnaryOp(tok)
	}
	return p.BinaryOp(tok)
}

var opTokens = [...]token.Token{
	exec.OpAdd:       token.ADD,
	exec.OpSub:       token.SUB,
	exec.OpMul:       token.MUL,
	exec.OpDiv:       token.QUO,
	exec.OpMod:       token.REM,
	exec.OpBitAnd:    token.AND,
	exec.OpBitOr:     token.OR,
	exec.OpBitXor:    token.XOR,
	exec.OpBitAndNot: token.AND_NOT,
	exec.OpBitSHL:    token.SHL,
	exec.OpBitSHR:    token.SHR,
	exec.OpLT:        token.LSS,
	exec.OpLE:        token.LEQ,
	exec.OpGT:        token.GTR,
	exec.OpGE:        token.GEQ,
	exec.OpEQ:        token.EQL,
	exec.OpEQNil:     token.ILLEGAL,
	exec.OpNE:        token.NEQ,
	exec.OpNENil:     token.ILLEGAL,
	exec.OpLAnd:      token.LAND,
	exec.OpLOr:       token.LOR,
	exec.OpNeg:       token.SUB,
	exec.OpNot:       token.NOT,
	exec.OpBitNot:    token.XOR,
}

// TypeCast instr
func (p *Builder) TypeCast(from, to reflect.Type) *Builder {
	t := Type(p, to)
	x := p.code.Pop().(ast.Expr)
	p.code.Push(&ast.CallExpr{
		Fun:  t,
		Args: []ast.Expr{x},
	})
	return p
}

// Call instr
func (p *Builder) Call(narg int, ellipsis bool) *Builder {
	fun := p.code.Pop().(ast.Expr)
	args := make([]ast.Expr, narg)
	for i := narg - 1; i >= 0; i-- {
		args[i] = p.code.Pop().(ast.Expr)
	}
	expr := &ast.CallExpr{Fun: fun, Args: args}
	if ellipsis {
		expr.Ellipsis++
	}
	p.code.Push(expr)
	return p
}

// CallGoFunc instr
func (p *Builder) CallGoFunc(fun exec.GoFuncAddr) *Builder {
	gfi := defaultImpl.GetGoFuncInfo(fun)
	pkgPath, name := gfi.Pkg.PkgPath(), gfi.Name
	fn := p.GoFuncIdent(pkgPath, name)
	p.code.Push(fn)
	arity := reflect.TypeOf(gfi.This).NumIn()
	return p.Call(arity, false)
}

// CallGoFuncv instr
func (p *Builder) CallGoFuncv(fun exec.GoFuncvAddr, arity int) *Builder {
	gfi := defaultImpl.GetGoFuncvInfo(fun)
	pkgPath, name := gfi.Pkg.PkgPath(), gfi.Name
	if pkgPath == "" {
		if alias, ok := builtinFnvs[name]; ok {
			pkgPath, name = alias[0], alias[1]
		}
	}
	fn := p.GoFuncIdent(pkgPath, name)
	p.code.Push(fn)
	ellipsis := arity == -1
	if ellipsis {
		arity = reflect.TypeOf(gfi.This).NumIn()
	}
	return p.Call(arity, ellipsis)
}

var builtinFnvs = map[string][2]string{
	"print":   {"fmt", "Print"},
	"printf":  {"fmt", "Printf"},
	"println": {"fmt", "Println"},
	"fprintf": {"fmt", "Fprintf"},
}

// -----------------------------------------------------------------------------
