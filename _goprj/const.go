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

package goprj

import (
	"go/ast"
	"go/token"
	"reflect"
	"unsafe"

	"github.com/qiniu/goplus/ast/astutil"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// ToLen converts ast.Expr to a Len.
func (p *fileLoader) ToLen(e ast.Expr) int64 {
	if e != nil {
		_, val := p.ToConst(e, -1)
		if _, ok := val.(*UninferedType); ok {
			log.Fatal("ToLen:", reflect.TypeOf(e))
		}
		if n, ok := val.(int64); ok {
			return n
		}
		log.Debug("ToLen:", reflect.TypeOf(val))
		return reflect.ValueOf(val).Int()
	}
	return 0
}

// ToConst infers constant value from a ast.Expr.
func (p *fileLoader) ToConst(expr ast.Expr, i int64) (typ Type, val interface{}) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		kind, n := astutil.ToConst(v)
		if astutil.IsConstBound(kind) {
			return AtomType(kind), n
		}
		return Unbound, n
	case *ast.Ident:
		if v.Name == "iota" {
			return Unbound, i
		}
	case *ast.SelectorExpr:
	case *ast.BinaryExpr:
		tx, x := p.ToConst(v.X, i)
		ty, y := p.ToConst(v.Y, i)
		tret, ret := binaryOp(v.Op, tx.(AtomType), ty.(AtomType), x, y)
		assertNotOverflow(tret, ret)
		return tret, ret
	case *ast.UnaryExpr:
		tx, x := p.ToConst(v.X, i)
		tret, ret := unaryOp(v.Op, tx.(AtomType), x)
		assertNotOverflow(tret, ret)
		return tret, ret
	case *ast.Ellipsis:
		return Unbound, -1
	case *ast.CallExpr:
		switch fun := v.Fun.(type) {
		case *ast.SelectorExpr:
			switch recv := fun.X.(type) {
			case *ast.Ident:
				if recv.Name == "unsafe" {
					switch fun.Sel.Name {
					case "Sizeof":
						t, _ := p.ToConst(v.Args[0], -1)
						return Uintptr, uint64(t.Sizeof(p.prj))
					}
				}
				log.Fatalln("ToConst CallExpr/SelectorExpr: unknown -", recv.Name, fun.Sel.Name)
			}
		case *ast.Ident:
			return p.callBuiltin(fun.Name, v.Args)
		default:
			log.Fatalln("ToConst CallExpr: unknown -", reflect.TypeOf(fun))
		}
	case *ast.CompositeLit:
		if v.Type != nil {
			return p.ToType(v.Type), nil
		}
	case *ast.ParenExpr:
		return p.ToConst(v.X, i)
	}
	log.Fatalln("ToConst: unknown -", reflect.TypeOf(expr), "-", expr)
	return &UninferedType{Expr: expr}, expr
}

// -----------------------------------------------------------------------------

func (p *fileLoader) callBuiltin(fun string, args []ast.Expr) (Type, interface{}) {
	if bfn, ok := builtinFuns[fun]; ok {
		return bfn(p, fun, args)
	}
	log.Fatalln("callBuiltin: unknown builtin func -", fun)
	return nil, nil
}

var builtinFuns map[string]func(*fileLoader, string, []ast.Expr) (Type, interface{})

var _builtinFuns = map[string]func(*fileLoader, string, []ast.Expr) (Type, interface{}){
	"bool":       atomTypeCast,
	"int":        atomTypeCast,
	"int16":      atomTypeCast,
	"int32":      atomTypeCast,
	"int64":      atomTypeCast,
	"uint":       atomTypeCast,
	"uint8":      atomTypeCast,
	"uint16":     atomTypeCast,
	"uint32":     atomTypeCast,
	"uint64":     atomTypeCast,
	"uintptr":    atomTypeCast,
	"float32":    atomTypeCast,
	"float64":    atomTypeCast,
	"complex64":  atomTypeCast,
	"complex128": atomTypeCast,
	"string":     atomTypeCast,
	"rune":       atomTypeCast,
	"byte":       atomTypeCast,
	"complex":    toComplex,
}

func init() {
	builtinFuns = _builtinFuns
}

func toComplex(p *fileLoader, fun string, args []ast.Expr) (Type, interface{}) {
	_, x := atomTypeCast(p, "float64", args[:1])
	_, y := atomTypeCast(p, "float64", args[1:])
	return Complex128, complex(x.(float64), y.(float64))
}

func atomTypeCast(p *fileLoader, fun string, args []ast.Expr) (Type, interface{}) {
	_, x := p.ToConst(args[0], -1)
	switch vx := x.(type) {
	case int64:
		bti := intCastTypes[fun]
		if (bti.from & (1 << reflect.Int64)) != 0 {
			return bti.typ, vx << bti.shbits >> bti.shbits
		}
		if (bti.from & (1 << reflect.Int)) != 0 {
			switch bti.typ {
			case String:
				return String, string(vx)
			case Float64:
				return Float64, float64(vx)
			case Float32:
				return Float32, float64(float32(vx))
			}
			return bti.typ, uint64(vx) << bti.shbits >> bti.shbits
		}
	case uint64:
		bti := intCastTypes[fun]
		if (bti.from & (1 << reflect.Uint64)) != 0 {
			return bti.typ, vx << bti.shbits >> bti.shbits
		}
		if (bti.from & (1 << reflect.Uint)) != 0 {
			switch bti.typ {
			case String:
				return String, string(vx)
			case Float64:
				return Float64, float64(vx)
			case Float32:
				return Float32, float64(float32(vx))
			}
			return bti.typ, int64(vx) << bti.shbits >> bti.shbits
		}
	case string:
		if fun == "string" {
			if _, ok := x.(string); ok {
				return String, x
			}
		}
	case float64:
		bti := intCastTypes[fun]
		if (bti.from & (1 << reflect.Float64)) != 0 {
			return bti.typ, x
		}
		if (bti.from & (1 << reflect.Float32)) != 0 {
			if (bti.from & (1 << reflect.Int64)) != 0 {
				return bti.typ, int64(vx) << bti.shbits >> bti.shbits
			}
			if (bti.from & (1 << reflect.Uint64)) != 0 {
				return bti.typ, uint64(vx) << bti.shbits >> bti.shbits
			}
		}
	case bool:
		if fun == "bool" {
			if _, ok := x.(bool); ok {
				return Bool, x
			}
		}
	case complex128:
		if _, ok := x.(complex128); ok {
			switch fun {
			case "complex64":
				return Complex64, x
			case "complex128":
				return Complex128, x
			}
		}
	}
	log.Fatalln("atomTypeCast:", fun, x)
	return nil, nil
}

var intCastTypes = map[string]bTI{
	"int":     {Int, (1 << reflect.Int64) | (1 << reflect.Uint), 64 - unsafe.Sizeof(int(0))*8},
	"int16":   {Int16, (1 << reflect.Int64) | (1 << reflect.Uint) | (1 << reflect.Float32), 64 - 16},
	"int32":   {Int32, (1 << reflect.Int64) | (1 << reflect.Uint) | (1 << reflect.Float32), 64 - 32},
	"int64":   {Int64, (1 << reflect.Int64) | (1 << reflect.Uint) | (1 << reflect.Float32), 0},
	"uint":    {Uint, (1 << reflect.Int) | (1 << reflect.Uint64) | (1 << reflect.Float32), 64 - unsafe.Sizeof(uint(0))*8},
	"uint8":   {Uint8, (1 << reflect.Int) | (1 << reflect.Uint64) | (1 << reflect.Float32), 64 - 8},
	"uint16":  {Uint16, (1 << reflect.Int) | (1 << reflect.Uint64) | (1 << reflect.Float32), 64 - 16},
	"uint32":  {Uint32, (1 << reflect.Int) | (1 << reflect.Uint64) | (1 << reflect.Float32), 64 - 32},
	"uint64":  {Uint64, (1 << reflect.Int) | (1 << reflect.Uint64) | (1 << reflect.Float32), 0},
	"uintptr": {Uintptr, (1 << reflect.Int) | (1 << reflect.Uint64) | (1 << reflect.Float32), 64 - unsafe.Sizeof(uintptr(0))*8},
	"float32": {Float32, (1 << reflect.Int) | (1 << reflect.Uint) | (1 << reflect.Float64), 0},
	"float64": {Float64, (1 << reflect.Int) | (1 << reflect.Uint) | (1 << reflect.Float64), 0},
	"string":  {String, (1 << reflect.Int) | (1 << reflect.Uint) | (1 << reflect.String), 0},
	"rune":    {Rune, (1 << reflect.Int64) | (1 << reflect.Uint) | (1 << reflect.Float32), 64 - unsafe.Sizeof(rune(0))*8},
	"byte":    {Byte, (1 << reflect.Int) | (1 << reflect.Uint64) | (1 << reflect.Float32), 64 - 1},
}

type bTI struct {
	typ    AtomType
	from   int
	shbits uintptr
}

// -----------------------------------------------------------------------------

func assertNotOverflow(tx AtomType, x interface{}) {
	// TODO:
	// log.Fatalln("assertNotOverflow:", tx, x)
}

func checkType(tx, ty AtomType) (AtomType, bool) {
	if tx == ty {
		return tx, true
	}
	if tx == Unbound {
		return ty, true
	}
	if ty == Unbound {
		return tx, true
	}
	return Unbound, false
}

func assertUnbound(t AtomType) {
	if t != Unbound {
		log.Fatalln("assertUnbound: type -", t)
	}
}

func checkValue(tx, ty AtomType, x, y interface{}) (nx, ny interface{}) {
	switch vx := x.(type) {
	case int64:
		switch vy := y.(type) {
		case int64:
			return x, y
		case uint64:
			assertUnbound(ty)
			return x, int64(vy)
		case float64:
			assertUnbound(tx)
			return float64(vx), y
		case complex128:
			assertUnbound(tx)
			return complex(float64(vx), 0), y
		default:
			log.Fatalln("checkValue failed: <int> op <unnkown> -", reflect.TypeOf(y))
		}
	case string:
		if _, ok := y.(string); ok {
			return x, y
		}
		log.Fatalln("checkValue failed: <string> op <unnkown> -", reflect.TypeOf(y))
	case uint64:
		switch y.(type) {
		case int64:
			assertUnbound(tx)
			return int64(vx), y
		case uint64:
			return x, y
		case float64:
			assertUnbound(tx)
			return float64(vx), y
		case complex128:
			assertUnbound(tx)
			return complex(float64(vx), 0), y
		default:
			log.Fatalln("checkValue failed: <uint> op <unnkown> -", reflect.TypeOf(y))
		}
	case float64:
		switch vy := y.(type) {
		case int64:
			assertUnbound(ty)
			return x, float64(vy)
		case uint64:
			assertUnbound(ty)
			return x, float64(vy)
		case float64:
			return x, y
		case complex128:
			assertUnbound(tx)
			return complex(vx, 0), y
		default:
			log.Fatalln("checkValue failed: <float64> op <unnkown> -", reflect.TypeOf(y))
		}
	case complex128:
		switch vy := y.(type) {
		case int64:
			assertUnbound(ty)
			return x, complex(float64(vy), 0)
		case uint64:
			assertUnbound(ty)
			return x, complex(float64(vy), 0)
		case float64:
			assertUnbound(ty)
			return x, complex(vy, 0)
		case complex128:
			return x, y
		default:
			log.Fatalln("checkValue failed: <complex128> op <unnkown> -", reflect.TypeOf(y))
		}
	}
	return nil, false
}

func binaryOp(op token.Token, tx, ty AtomType, x, y interface{}) (AtomType, interface{}) {
	switch op {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM, // + - * / %
		token.EQL, token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ, // == < > != <= >=
		token.AND, token.OR, token.XOR, token.AND_NOT: // & | ^ &^
		if t, ok := checkType(tx, ty); ok {
			x, y = checkValue(tx, ty, x, y)
			switch op {
			case token.ADD:
				switch vx := x.(type) {
				case int64:
					return t, vx + y.(int64)
				case string:
					return t, vx + y.(string)
				case uint64:
					return t, vx + y.(uint64)
				case float64:
					return t, vx + y.(float64)
				case complex128:
					return t, vx + y.(complex128)
				default:
					log.Fatalln("binaryOp + failed: unknown -", reflect.TypeOf(x))
				}
			case token.SUB:
				switch vx := x.(type) {
				case int64:
					return t, vx - y.(int64)
				case uint64:
					return t, vx - y.(uint64)
				case float64:
					return t, vx - y.(float64)
				case complex128:
					return t, vx - y.(complex128)
				default:
					log.Fatalln("binaryOp - failed: unknown -", reflect.TypeOf(x))
				}
			case token.MUL:
				switch vx := x.(type) {
				case int64:
					return t, vx * y.(int64)
				case uint64:
					return t, vx * y.(uint64)
				case float64:
					return t, vx * y.(float64)
				case complex128:
					return t, vx * y.(complex128)
				default:
					log.Fatalln("binaryOp * failed: unknown -", reflect.TypeOf(x))
				}
			case token.QUO:
				switch vx := x.(type) {
				case int64:
					return t, vx / y.(int64)
				case uint64:
					return t, vx / y.(uint64)
				case float64:
					return t, vx / y.(float64)
				case complex128:
					return t, vx / y.(complex128)
				default:
					log.Fatalln("binaryOp * failed: unknown -", reflect.TypeOf(x))
				}
			case token.REM:
				switch vx := x.(type) {
				case int64:
					return t, vx % y.(int64)
				case uint64:
					return t, vx % y.(uint64)
				default:
					log.Fatalln("binaryOp % failed: unknown -", reflect.TypeOf(x))
				}
			case token.EQL: // ==
				return Bool, x == y
			case token.LSS: // <
				switch vx := x.(type) {
				case int64:
					return Bool, vx < y.(int64)
				case string:
					return Bool, vx < y.(string)
				case uint64:
					return Bool, vx < y.(uint64)
				case float64:
					return Bool, vx < y.(float64)
				default:
					log.Fatalln("binaryOp < failed: unknown -", reflect.TypeOf(x))
				}
			case token.GTR: // >
				switch vx := x.(type) {
				case int64:
					return Bool, vx > y.(int64)
				case string:
					return Bool, vx > y.(string)
				case uint64:
					return Bool, vx > y.(uint64)
				case float64:
					return Bool, vx > y.(float64)
				default:
					log.Fatalln("binaryOp > failed: unknown -", reflect.TypeOf(x))
				}
			case token.NEQ: // !=
				return Bool, x != y
			case token.LEQ: // <=
				switch vx := x.(type) {
				case int64:
					return Bool, vx <= y.(int64)
				case string:
					return Bool, vx <= y.(string)
				case uint64:
					return Bool, vx <= y.(uint64)
				case float64:
					return Bool, vx <= y.(float64)
				default:
					log.Fatalln("binaryOp <= failed: unknown -", reflect.TypeOf(x))
				}
			case token.GEQ: // >=
				switch vx := x.(type) {
				case int64:
					return Bool, vx >= y.(int64)
				case string:
					return Bool, vx >= y.(string)
				case uint64:
					return Bool, vx >= y.(uint64)
				case float64:
					return Bool, vx >= y.(float64)
				default:
					log.Fatalln("binaryOp >= failed: unknown -", reflect.TypeOf(x))
				}
			case token.AND: // &
				switch vx := x.(type) {
				case int64:
					return t, vx & y.(int64)
				case uint64:
					return t, vx & y.(uint64)
				default:
					log.Fatalln("binaryOp & failed: unknown -", reflect.TypeOf(x))
				}
			case token.OR: // |
				switch vx := x.(type) {
				case int64:
					return t, vx | y.(int64)
				case uint64:
					return t, vx | y.(uint64)
				default:
					log.Fatalln("binaryOp | failed: unknown -", reflect.TypeOf(x))
				}
			case token.XOR: //  ^
				switch vx := x.(type) {
				case int64:
					return t, vx ^ y.(int64)
				case uint64:
					return t, vx ^ y.(uint64)
				default:
					log.Fatalln("binaryOp ^ failed: unknown -", reflect.TypeOf(x))
				}
			case token.AND_NOT: // &^
				switch vx := x.(type) {
				case int64:
					return t, vx &^ y.(int64)
				case uint64:
					return t, vx &^ y.(uint64)
				default:
					log.Fatalln("binaryOp &^ failed: unknown -", reflect.TypeOf(x))
				}
			}
		}
		log.Fatalln("binaryOp: checkType failed -", tx, op, ty)
	case token.LAND, token.LOR: // && ||
		vx, ok1 := x.(bool)
		vy, ok2 := y.(bool)
		if ok1 && ok2 {
			if op == token.LAND {
				return Bool, vx && vy
			}
			return Bool, vx || vy
		}
		log.Fatalln("binaryOp: bool expression required -", tx, op, ty)
	case token.SHL: // <<
		switch vx := x.(type) {
		case int64:
			switch vy := y.(type) {
			case int64:
				return tx, vx << vy
			case uint64:
				return tx, vx << vy
			default:
				log.Fatalln("binaryOp failed: int << unknown -", ty)
			}
		case uint64:
			switch vy := y.(type) {
			case int64:
				return tx, vx << vy
			case uint64:
				return tx, vx << vy
			default:
				log.Fatalln("binaryOp failed: int << unknown -", ty)
			}
		default:
			log.Fatalln("binaryOp << failed: unknown -", tx, op, ty)
		}
	case token.SHR: // >>
		switch vx := x.(type) {
		case int64:
			switch vy := y.(type) {
			case int64:
				return tx, vx >> vy
			case uint64:
				return tx, vx >> vy
			default:
				log.Fatalln("binaryOp failed: int >> unknown -", ty)
			}
		case uint64:
			switch vy := y.(type) {
			case int64:
				return tx, vx >> vy
			case uint64:
				return tx, vx >> vy
			default:
				log.Fatalln("binaryOp failed: int >> unknown -", ty)
			}
		default:
			log.Fatalln("binaryOp >> failed: unknown -", tx, op, ty)
		}
	}
	log.Fatalln("binaryOp failed: unknown -", tx, op, ty)
	return Unbound, nil
}

func unaryOp(op token.Token, tx AtomType, x interface{}) (AtomType, interface{}) {
	switch op {
	case token.SUB: // -
		switch vx := x.(type) {
		case int64:
			return tx, -vx
		case uint64:
			return tx, -vx
		default:
			log.Fatalln("unaryOp - failed:", tx, x)
		}
	case token.XOR: // ^
		switch vx := x.(type) {
		case int64:
			return tx, ^vx
		case uint64:
			return tx, ^vx
		default:
			log.Fatalln("unaryOp ^ failed:", tx, x)
		}
	case token.NOT:
		if vx, ok := x.(bool); ok {
			return Bool, !vx
		}
		log.Fatalln("unaryOp ! failed:", tx, x)
	}
	log.Fatalln("unaryOp failed:", op, tx, x)
	return Unbound, nil
}

// -----------------------------------------------------------------------------
