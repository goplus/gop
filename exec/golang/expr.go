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
	"math/big"
	"strconv"
	"strings"

	"github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/lib/builtin"
	"github.com/goplus/gop/reflect"
	"github.com/qiniu/x/log"
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

// GoSymIdent - ast.Ident or ast.SelectorExpr
func (p *Builder) GoSymIdent(pkgPath, name string) ast.Expr {
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

// NewGopkgType instr
func NewGopkgType(p *Builder, pkgPath, typName string) ast.Expr {
	typ := p.GoSymIdent(pkgPath, typName)
	args := []ast.Expr{typ}
	return &ast.CallExpr{Fun: newIdent, Args: args}
}

func valBySetString(p *Builder, typ reflect.Type, x ast.Expr, args ...ast.Expr) ast.Expr {
	setString := &ast.SelectorExpr{X: x, Sel: Ident("SetString")}
	setStringCall := &ast.CallExpr{Fun: setString, Args: args}
	stmt := &ast.AssignStmt{
		Lhs: []ast.Expr{gopRet, unnamedVar},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{setStringCall},
	}
	fldOut := &ast.Field{
		Names: []*ast.Ident{gopRet},
		Type:  Type(p, typ),
	}
	typFun := &ast.FuncType{
		Params:  &ast.FieldList{Opening: 1, Closing: 1},
		Results: &ast.FieldList{Opening: 1, Closing: 1, List: []*ast.Field{fldOut}},
	}
	stmtReturn := &ast.ReturnStmt{}
	return &ast.CallExpr{
		Fun: &ast.FuncLit{
			Type: typFun,
			Body: &ast.BlockStmt{List: []ast.Stmt{stmt, stmtReturn}},
		},
	}
}

// BigIntConst instr
func BigIntConst(p *Builder, v *big.Int) ast.Expr {
	if v.IsInt64() {
		newInt := p.GoSymIdent("math/big", "NewInt")
		args := []ast.Expr{IntConst(v.Int64())}
		return &ast.CallExpr{Fun: newInt, Args: args}
	}
	bigInt := NewGopkgType(p, "math/big", "Int")
	return valBySetString(p, exec.TyBigInt, bigInt, StringConst(v.String()), IntConst(10))
}

// BigRatConst instr
func BigRatConst(p *Builder, v *big.Rat) ast.Expr {
	a, b := v.Num(), v.Denom()
	if a.IsInt64() && b.IsInt64() {
		newRat := p.GoSymIdent("math/big", "NewRat")
		args := []ast.Expr{IntConst(a.Int64()), IntConst(b.Int64())}
		return &ast.CallExpr{Fun: newRat, Args: args}
	}
	rat := NewGopkgType(p, "math/big", "Rat")
	setFrac := &ast.SelectorExpr{X: rat, Sel: Ident("SetFrac")}
	args := []ast.Expr{BigIntConst(p, a), BigIntConst(p, b)}
	return &ast.CallExpr{Fun: setFrac, Args: args}
}

// BigFloatConst instr
func BigFloatConst(p *Builder, v *big.Float) ast.Expr {
	val, acc := v.Float64()
	if acc == big.Exact {
		newFloat := p.GoSymIdent("math/big", "NewFloat")
		args := []ast.Expr{FloatConst(val)}
		return &ast.CallExpr{Fun: newFloat, Args: args}
	}
	prec := v.Prec()
	sval := v.Text('g', int(prec))
	bigFlt := NewGopkgType(p, "math/big", "Float")
	setPrec := &ast.SelectorExpr{X: bigFlt, Sel: Ident("SetPrec")}
	setPrecCall := &ast.CallExpr{Fun: setPrec, Args: []ast.Expr{IntConst(int64(prec))}}
	return valBySetString(p, exec.TyBigFloat, setPrecCall, StringConst(sval))
}

// Const instr
func Const(p *Builder, val interface{}) ast.Expr {
	if val == nil {
		return nilIdent
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
	switch kind {
	case reflect.Bool:
		if val.(bool) {
			return Ident("true")
		}
		return Ident("false")
	case reflect.Ptr:
		switch v.Type() {
		case exec.TyBigRat:
			return BigRatConst(p, val.(*big.Rat))
		case exec.TyBigInt:
			return BigIntConst(p, val.(*big.Int))
		case exec.TyBigFloat:
			return BigFloatConst(p, val.(*big.Float))
		}
	}
	log.Panicln("Const: value type is unknown -", v.Type())
	return nil
}

// Push instr
func (p *Builder) Push(val interface{}) *Builder {
	p.rhs.Push(Const(p, val))
	return p
}

func (p *Builder) bigBuiltinOp(kind exec.Kind, op exec.Operator) *Builder {
	val := p.rhs.Pop().(ast.Expr)
	if op >= exec.OpLT && op <= exec.OpNE {
		x := p.rhs.Pop().(ast.Expr)
		bigOp := &ast.SelectorExpr{X: x, Sel: Ident("Cmp")}
		bigOpCall := &ast.CallExpr{Fun: bigOp, Args: []ast.Expr{val}}
		p.rhs.Push(&ast.BinaryExpr{X: bigOpCall, Y: IntConst(0), Op: opTokens[op]})
		return p
	}
	method := opOpMethods[op]
	if method == "" {
		log.Panicln("bigBuiltinOp: unsupported op -", op)
	}
	t := exec.TypeFromKind(kind).Elem()
	typ := NewGopkgType(p, t.PkgPath(), t.Name())
	bigOp := &ast.SelectorExpr{X: typ, Sel: Ident(method)}
	oi := op.GetInfo()
	if oi.InSecond == 0 {
		p.rhs.Push(&ast.CallExpr{Fun: bigOp, Args: []ast.Expr{val}})
		return p
	}
	x := p.rhs.Pop().(ast.Expr)
	p.rhs.Push(&ast.CallExpr{Fun: bigOp, Args: []ast.Expr{x, val}})
	return p
}

var opOpMethods = [...]string{
	exec.OpAdd:    "Add",
	exec.OpSub:    "Sub",
	exec.OpMul:    "Mul",
	exec.OpQuo:    "Quo",
	exec.OpMod:    "Mod",
	exec.OpAnd:    "And",
	exec.OpOr:     "Or",
	exec.OpXor:    "Xor",
	exec.OpAndNot: "AndNot",
	exec.OpLsh:    "Lsh",
	exec.OpRsh:    "Rsh",
	exec.OpNeg:    "Neg",
	exec.OpBitNot: "Not",
}

// BinaryOp instr
func BinaryOp(tok token.Token, x, y ast.Expr) ast.Expr {
	return &ast.BinaryExpr{Op: tok, X: x, Y: y}
}

// BuiltinOp instr
func (p *Builder) BuiltinOp(kind exec.Kind, op exec.Operator) *Builder {
	if kind >= exec.BigInt {
		return p.bigBuiltinOp(kind, op)
	}
	tok := opTokens[op]
	if tok == token.ILLEGAL {
		log.Panicln("BuiltinOp: unsupported op -", op)
	}
	oi := op.GetInfo()
	val := p.rhs.Pop().(ast.Expr)
	if oi.InSecond == 0 {
		p.rhs.Push(&ast.UnaryExpr{Op: tok, X: val})
		return p
	}
	x := p.rhs.Pop().(ast.Expr)
	p.rhs.Push(&ast.BinaryExpr{Op: tok, X: x, Y: val})
	return p
}

var opTokens = [...]token.Token{
	exec.OpAdd:    token.ADD,
	exec.OpSub:    token.SUB,
	exec.OpMul:    token.MUL,
	exec.OpQuo:    token.QUO,
	exec.OpMod:    token.REM,
	exec.OpAnd:    token.AND,
	exec.OpOr:     token.OR,
	exec.OpXor:    token.XOR,
	exec.OpAndNot: token.AND_NOT,
	exec.OpLsh:    token.SHL,
	exec.OpRsh:    token.SHR,
	exec.OpLT:     token.LSS,
	exec.OpLE:     token.LEQ,
	exec.OpGT:     token.GTR,
	exec.OpGE:     token.GEQ,
	exec.OpEQ:     token.EQL,
	exec.OpEQNil:  token.ILLEGAL,
	exec.OpNE:     token.NEQ,
	exec.OpNENil:  token.ILLEGAL,
	exec.OpLAnd:   token.LAND,
	exec.OpLOr:    token.LOR,
	exec.OpLNot:   token.NOT,
	exec.OpNeg:    token.SUB,
	exec.OpBitNot: token.XOR,
}

// TypeCast instr
func (p *Builder) TypeCast(from, to reflect.Type) *Builder {
	x := p.rhs.Pop().(ast.Expr)
	p.rhs.Push(TypeCast(p, x, from, to))
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
	if ct := p.inDeferOrGo; ct == callExpr {
		p.rhs.Push(expr)
	} else {
		p.inDeferOrGo = callExpr
		if ct == callByDefer {
			p.rhs.Push(&ast.DeferStmt{
				Call: expr,
			})
		} else {
			p.rhs.Push(&ast.GoStmt{
				Call: expr,
			})
		}
	}
	return p
}

// CallGoFunc instr
func (p *Builder) CallGoFunc(fun exec.GoFuncAddr, nexpr int) *Builder {
	gfi := defaultImpl.GetGoFuncInfo(fun)
	pkgPath, name := gfi.Pkg.PkgPath(), gfi.Name
	if pkgPath == "" {
		if alias, ok := builtin.FuncGoInfo(name); ok {
			pkgPath, name = alias[0], alias[1]
		}
	}
	fn := p.GoSymIdent(pkgPath, name)
	p.rhs.Push(fn)
	return p.Call(nexpr, false)
}

// CallGoFuncv instr
func (p *Builder) CallGoFuncv(fun exec.GoFuncvAddr, nexpr, arity int) *Builder {
	gfi := defaultImpl.GetGoFuncvInfo(fun)
	pkgPath, name := gfi.Pkg.PkgPath(), gfi.Name
	if pkgPath == "" {
		if alias, ok := builtin.FuncGoInfo(name); ok {
			pkgPath, name = alias[0], alias[1]
		}
	}
	fn := p.GoSymIdent(pkgPath, name)
	p.rhs.Push(fn)
	return p.Call(nexpr, arity == -1)
}

// LoadGoVar instr
func (p *Builder) LoadGoVar(addr exec.GoVarAddr) *Builder {
	gvi := defaultImpl.GetGoVarInfo(addr)
	p.rhs.Push(p.GoSymIdent(gvi.Pkg.PkgPath(), gvi.Name))
	return p
}

// StoreGoVar instr
func (p *Builder) StoreGoVar(addr exec.GoVarAddr) *Builder {
	gvi := defaultImpl.GetGoVarInfo(addr)
	p.lhs.Push(p.GoSymIdent(gvi.Pkg.PkgPath(), gvi.Name))
	return p
}

// AddrGoVar instr
func (p *Builder) AddrGoVar(addr exec.GoVarAddr) *Builder {
	gvi := defaultImpl.GetGoVarInfo(addr)
	p.rhs.Push(&ast.UnaryExpr{
		Op: token.AND,
		X:  p.GoSymIdent(gvi.Pkg.PkgPath(), gvi.Name),
	})
	return p
}

func (p *Builder) fieldExpr(typ reflect.Type, index []int) ast.Expr {
	expr := &ast.SelectorExpr{}
	expr.X = p.rhs.Pop().(ast.Expr)
	if unary, ok := expr.X.(*ast.UnaryExpr); ok {
		expr.X = unary.X
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	for i := 0; i < len(index); i++ {
		sf := typ.FieldByIndex(index[:i+1])
		if sf.Anonymous {
			continue
		}
		if expr.Sel != nil {
			expr.X = &ast.SelectorExpr{X: expr.X, Sel: expr.Sel}
		}
		expr.Sel = Ident(sf.Name)
	}

	return expr
}

// LoadField instr
func (p *Builder) LoadField(typ reflect.Type, index []int) *Builder {
	expr := p.fieldExpr(typ, index)
	p.rhs.Push(expr)
	return p
}

// AddrField instr
func (p *Builder) AddrField(typ reflect.Type, index []int) *Builder {
	expr := p.fieldExpr(typ, index)
	p.rhs.Push(&ast.UnaryExpr{
		Op: token.AND,
		X:  expr,
	})
	return p
}

// StoreField instr
func (p *Builder) StoreField(typ reflect.Type, index []int) *Builder {
	expr := p.fieldExpr(typ, index)
	p.lhs.Push(expr)
	return p
}

// Append instr
func (p *Builder) Append(typ reflect.Type, arity int) *Builder {
	p.rhs.Push(appendIdent)
	var ellipsis bool
	if arity == -1 {
		ellipsis = true
		arity = 2
	}
	p.Call(arity, ellipsis)
	return p
}

// New instr
func (p *Builder) New(typ reflect.Type) *Builder {
	p.rhs.Push(newIdent)
	p.Call(0, false, Type(p, typ))
	return p
}

// Make instr
func (p *Builder) Make(typ reflect.Type, arity int) *Builder {
	p.rhs.Push(makeIdent)
	p.Call(arity, false, Type(p, typ))
	return p
}

// MakeArray instr
func (p *Builder) MakeArray(typ reflect.Type, arity int) *Builder {
	typExpr := Type(p, typ)
	elts := make([]ast.Expr, arity)
	for i, v := range p.rhs.GetArgs(arity) {
		elts[i] = v.(ast.Expr)
	}
	var xExpr ast.Expr = &ast.CompositeLit{
		Type: typExpr,
		Elts: elts,
	}
	p.rhs.Ret(arity, xExpr)
	return p
}

// MakeMap instr
func (p *Builder) MakeMap(typ reflect.Type, arity int) *Builder {
	typExpr := Type(p, typ)
	elts := make([]ast.Expr, arity)
	args := p.rhs.GetArgs(arity << 1)
	for i := 0; i < arity; i++ {
		elts[i] = &ast.KeyValueExpr{
			Key:   args[i<<1].(ast.Expr),
			Value: args[(i<<1)+1].(ast.Expr),
		}
	}
	p.rhs.Ret(arity<<1, &ast.CompositeLit{
		Type: typExpr,
		Elts: elts,
	})
	return p
}

// MapIndex instr
func (p *Builder) MapIndex() *Builder {
	p.rhs.Push(Index(p))
	return p
}

// SetMapIndex instr
func (p *Builder) SetMapIndex() *Builder {
	p.lhs.Push(Index(p))
	return p
}

// Index instr
func (p *Builder) Index(idx int) *Builder {
	p.rhs.Push(IndexWith(p, idx))
	return p
}

// AddrIndex instr
func (p *Builder) AddrIndex(idx int) *Builder {
	p.rhs.Push(&ast.UnaryExpr{
		Op: token.AND,
		X:  IndexWith(p, idx),
	})
	return p
}

// SetIndex instr
func (p *Builder) SetIndex(idx int) *Builder {
	p.lhs.Push(IndexWith(p, idx))
	return p
}

// Index instr
func Index(p *Builder) *ast.IndexExpr {
	idx := p.rhs.Pop().(ast.Expr)
	x := p.rhs.Pop().(ast.Expr)
	return &ast.IndexExpr{
		X:     x,
		Index: idx,
	}
}

// IndexWith instr
func IndexWith(p *Builder, idx int) *ast.IndexExpr {
	if idx == -1 {
		return Index(p)
	}
	x := p.rhs.Pop().(ast.Expr)
	return &ast.IndexExpr{
		X:     x,
		Index: IntConst(int64(idx)),
	}
}

// Slice instr
func (p *Builder) Slice(i, j int) *Builder {
	jExpr := SliceIndex(p, j)
	iExpr := SliceIndex(p, i)
	p.rhs.Push(&ast.SliceExpr{
		X:    p.rhs.Pop().(ast.Expr),
		Low:  iExpr,
		High: jExpr,
	})
	return p
}

// Slice3 instr
func (p *Builder) Slice3(i, j, k int) *Builder {
	kExpr := SliceIndex(p, k)
	jExpr := SliceIndex(p, j)
	iExpr := SliceIndex(p, i)
	p.rhs.Push(&ast.SliceExpr{
		X:      p.rhs.Pop().(ast.Expr),
		Low:    iExpr,
		High:   jExpr,
		Max:    kExpr,
		Slice3: true,
	})
	return p
}

// SliceIndex instr
func SliceIndex(p *Builder, idx int) ast.Expr {
	if idx == exec.SliceDefaultIndex {
		return nil
	}
	if idx == -1 {
		return p.rhs.Pop().(ast.Expr)
	}
	return IntConst(int64(idx))
}

// Zero instr
func (p *Builder) Zero(typ reflect.Type) *Builder {
	p.rhs.Push(Zero(p, typ))
	return p
}

// Zero instr
func Zero(p *Builder, typ reflect.Type) ast.Expr {
	kind := typ.Kind()
	if kind == reflect.String {
		return StringConst("")
	}
	if kind >= reflect.Int && kind <= reflect.Complex128 {
		return IntConst(0)
	}
	if kind == reflect.Bool {
		return Ident("false")
	}
	log.Panicln("Zero: unknown -", typ)
	return nil
}

// GoBuiltin instr
func (p *Builder) GoBuiltin(typ reflect.Type, op exec.GoBuiltin) *Builder {
	arity := goBuiltinArities[op]
	p.rhs.Push(Ident(op.String()))
	return p.Call(arity, false)
}

var goBuiltinArities = [...]int{
	exec.GobLen:     1,
	exec.GobCap:     1,
	exec.GobCopy:    2,
	exec.GobDelete:  2,
	exec.GobComplex: 2,
	exec.GobReal:    1,
	exec.GobImag:    1,
	exec.GobClose:   1,
}

// -----------------------------------------------------------------------------

// Struct instr
func (p *Builder) Struct(typ reflect.Type, arity int) *Builder {
	var ptr bool
	if typ.Kind() == reflect.Ptr {
		ptr = true
		typ = typ.Elem()
	}
	typExpr := Type(p, typ)
	elts := make([]ast.Expr, arity)
	args := p.rhs.GetArgs(arity << 1)
	for i := 0; i < arity; i++ {
		elts[i] = &ast.KeyValueExpr{
			Key:   toField(args[i<<1].(ast.Expr), typ),
			Value: args[(i<<1)+1].(ast.Expr),
		}
	}

	var ret ast.Expr

	ret = &ast.CompositeLit{
		Type: typExpr,
		Elts: elts,
	}
	if ptr {
		ret = &ast.UnaryExpr{
			Op: token.AND,
			X:  ret,
		}
	}
	p.rhs.Ret(arity<<1, ret)
	return p
}

func toField(expr ast.Expr, typ reflect.Type) *ast.Ident {
	if blit, ok := expr.(*ast.BasicLit); ok {
		i, err := strconv.Atoi(blit.Value)
		if err == nil {
			field := typ.Field(i)
			return Ident(field.Name)
		}
	}
	log.Panicln("toField expression must be arrayType")
	return nil
}
