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

package golang

import (
	"go/ast"
	"go/token"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/qiniu/qlang/v6/exec.spec"
)

// -----------------------------------------------------------------------------

// Ident instr
func (p *Builder) Ident(name string) *Builder {
	p.rhs.Push(&ast.Ident{Name: name})
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
	if strings.HasPrefix(name, "(") { // eg. name = "(*Replacer).Replace"
		pos := strings.LastIndexByte(name, '*')
		if pos == -1 {
			pos = 1
		} else {
			pos++
		}
		name = name[:pos] + pkg + "." + name[pos:]
		return Ident(name)
	}
	return &ast.SelectorExpr{
		X:   Ident(pkg),
		Sel: Ident(name),
	}
}

// StringConst - ast.BasicLit
func StringConst(v string) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(v),
	}
}

// IntConst - ast.BasicLit
func IntConst(v int64) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.FormatInt(v, 10),
	}
}

// UintConst instr
func UintConst(v uint64) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.FormatUint(v, 10),
	}
}

// FloatConst instr
func FloatConst(v float64) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.FLOAT,
		Value: strconv.FormatFloat(v, 'g', -1, 64),
	}
}

// ImagConst instr
func ImagConst(v float64) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.IMAG,
		Value: strconv.FormatFloat(v, 'g', -1, 64) + "i",
	}
}

// ComplexConst instr
func ComplexConst(v complex128) ast.Expr {
	r, i := real(v), imag(v)
	x, y := FloatConst(r), ImagConst(i)
	return &ast.ParenExpr{X: BinaryOp(token.ADD, x, y)}
}

// Const instr
func Const(p *Builder, val interface{}) ast.Expr {
	if val == nil {
		return Ident("nil")
	}
	v := reflect.ValueOf(val)
	kind := v.Kind()
	if kind == reflect.String {
		return StringConst(val.(string))
	}
	if kind >= reflect.Int && kind <= reflect.Int64 {
		var expr ast.Expr = IntConst(v.Int())
		if t := v.Type(); t != exec.TyInt {
			expr = TypeCast(p, expr, exec.TyInt, t)
		}
		return expr
	}
	if kind >= reflect.Uint && kind <= reflect.Uintptr {
		var expr ast.Expr = UintConst(v.Uint())
		return TypeCast(p, expr, exec.TyInt, v.Type())
	}
	if kind >= reflect.Float32 && kind <= reflect.Float64 {
		var expr ast.Expr = FloatConst(v.Float())
		if t := v.Type(); t != exec.TyFloat64 {
			expr = TypeCast(p, expr, exec.TyFloat64, t)
		}
		return expr
	}
	if kind >= reflect.Complex64 && kind <= reflect.Complex128 {
		var expr ast.Expr = ComplexConst(v.Complex())
		if t := v.Type(); t != exec.TyComplex128 {
			expr = TypeCast(p, expr, exec.TyComplex128, t)
		}
		return expr
	}
	if kind == reflect.Bool {
		if val.(bool) {
			return Ident("true")
		}
		return Ident("false")
	}
	log.Panicln("Const: value type is unknown -", v.Type())
	return nil
}

// Push instr
func (p *Builder) Push(val interface{}) *Builder {
	p.rhs.Push(Const(p, val))
	return p
}

// UnaryOp instr
func (p *Builder) UnaryOp(tok token.Token) *Builder {
	x := p.rhs.Pop().(ast.Expr)
	p.rhs.Push(&ast.UnaryExpr{Op: tok, X: x})
	return p
}

// BinaryOp instr
func (p *Builder) BinaryOp(tok token.Token) *Builder {
	y := p.rhs.Pop().(ast.Expr)
	x := p.rhs.Pop().(ast.Expr)
	p.rhs.Push(&ast.BinaryExpr{Op: tok, X: x, Y: y})
	return p
}

// BinaryOp instr
func BinaryOp(tok token.Token, x, y ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{Op: tok, X: x, Y: y}
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
	x := p.rhs.Pop().(ast.Expr)
	TypeCast(p, x, from, to)
	return p
}

// TypeCast instr
func TypeCast(p *Builder, x ast.Expr, from, to reflect.Type) *ast.CallExpr {
	t := Type(p, to)
	return &ast.CallExpr{
		Fun:  t,
		Args: []ast.Expr{x},
	}
}

// Call instr
func (p *Builder) Call(narg int, ellipsis bool, args ...ast.Expr) *Builder {
	fun := p.rhs.Pop().(ast.Expr)
	for _, item := range p.rhs.GetArgs(narg) {
		args = append(args, item.(ast.Expr))
	}
	p.rhs.PopN(narg)
	expr := &ast.CallExpr{Fun: fun, Args: args}
	if ellipsis {
		expr.Ellipsis++
	}
	p.rhs.Push(expr)
	return p
}

// CallGoFunc instr
func (p *Builder) CallGoFunc(fun exec.GoFuncAddr) *Builder {
	gfi := defaultImpl.GetGoFuncInfo(fun)
	pkgPath, name := gfi.Pkg.PkgPath(), gfi.Name
	fn := p.GoFuncIdent(pkgPath, name)
	p.rhs.Push(fn)
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
	p.rhs.Push(fn)
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

// Append instr
func (p *Builder) Append(typ reflect.Type, arity int) *Builder {
	p.rhs.Push(Ident("append"))
	var ellipsis bool
	if arity == -1 {
		ellipsis = true
		arity = 2
	}
	p.Call(arity, ellipsis)
	return p
}

// Make instr
func (p *Builder) Make(typ reflect.Type, arity int) *Builder {
	p.rhs.Push(Ident("make"))
	p.Call(arity, false, Type(p, typ))
	return p
}

// -----------------------------------------------------------------------------
