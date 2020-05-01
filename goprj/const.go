package goprj

import (
	"go/ast"
	"go/token"
	"reflect"
	"strconv"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// ToString converts a ast.BasicLit to string value.
func ToString(l *ast.BasicLit) string {
	if l.Kind == token.STRING {
		s, err := strconv.Unquote(l.Value)
		if err == nil {
			return s
		}
	}
	panic("ToString: convert ast.BasicLit to string failed")
}

// ToLen converts ast.Expr to a Len.
func (p *fileLoader) ToLen(e ast.Expr) int {
	if e != nil {
		_, val := p.ToConst(e, -1)
		if _, ok := val.(*UninferedType); ok {
			log.Fatal("ToLen:", reflect.TypeOf(e))
		}
		if n, ok := val.(int); ok {
			return n
		}
		log.Debug("ToLen:", reflect.TypeOf(val))
		return int(reflect.ValueOf(val).Int())
	}
	return 0
}

// ToConst infers constant value from a ast.Expr.
func (p *fileLoader) ToConst(expr ast.Expr, i int) (typ Type, val interface{}) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		switch v.Kind {
		case token.INT:
			n, err := strconv.ParseInt(v.Value, 0, 0)
			if err != nil {
				n2, err2 := strconv.ParseUint(v.Value, 0, 0)
				if err2 != nil {
					log.Fatalln("ToConst: strconv.ParseInt failed:", err2)
				}
				return Unbound, uint(n2)
			}
			return Unbound, int(n)
		case token.CHAR, token.STRING:
			n, err := strconv.Unquote(v.Value)
			if err != nil {
				log.Fatalln("ToConst: strconv.Unquote failed:", err)
			}
			if v.Kind == token.CHAR {
				for _, c := range n {
					return Rune, int(c)
				}
				panic("not here")
			}
			return String, n
		case token.FLOAT:
			n, err := strconv.ParseFloat(v.Value, 64)
			if err != nil {
				log.Fatalln("ToConst: strconv.ParseFloat failed:", err)
			}
			return Unbound, n
		case token.IMAG: // 123.45i
			val := v.Value
			n, err := strconv.ParseFloat(val[:len(val)-1], 64)
			if err != nil {
				log.Fatalln("ToConst: strconv.ParseFloat failed:", err)
			}
			return Unbound, complex(0, n)
		default:
			log.Fatalln("ToConst: unknown -", expr)
		}
	case *ast.Ident:
		if v.Name == "iota" {
			return Unbound, i
		}
	case *ast.SelectorExpr:
	case *ast.BinaryExpr:
		tx, x := p.ToConst(v.X, i)
		ty, y := p.ToConst(v.Y, i)
		return binaryOp(v.Op, tx, ty, x, y)
	case *ast.Ellipsis:
		return Unbound, -1
	}
	log.Fatalln("ToConst: unknown -", reflect.TypeOf(expr), "-", expr)
	return &UninferedType{Expr: expr}, expr
}

func checkType(tx, ty Type) (Type, bool) {
	if tx == ty {
		return tx, true
	}
	if tx == Unbound {
		return ty, true
	}
	if ty == Unbound {
		return tx, true
	}
	return nil, false
}

func assertUnbound(t Type) {
	if t != Unbound {
		log.Fatalln("assertUnbound: type -", t)
	}
}

func checkValue(tx, ty Type, x, y interface{}) (nx, ny interface{}) {
	switch vx := x.(type) {
	case int:
		switch vy := y.(type) {
		case int:
			return x, y
		case uint:
			assertUnbound(ty)
			return x, int(vy)
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
	case uint:
		switch y.(type) {
		case int:
			assertUnbound(tx)
			return int(vx), y
		case uint:
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
		case int:
			assertUnbound(ty)
			return x, float64(vy)
		case uint:
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
		case int:
			assertUnbound(ty)
			return x, complex(float64(vy), 0)
		case uint:
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

func binaryOp(op token.Token, tx, ty Type, x, y interface{}) (Type, interface{}) {
	switch op {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM, // + - * / %
		token.EQL, token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ, // == < > != <= >=
		token.AND, token.OR, token.XOR, token.AND_NOT: // & | ^ &^
		if t, ok := checkType(tx, ty); ok {
			x, y = checkValue(tx, ty, x, y)
			switch op {
			case token.ADD:
				switch vx := x.(type) {
				case int:
					return t, vx + y.(int)
				case string:
					return t, vx + y.(string)
				case uint:
					return t, vx + y.(uint)
				case float64:
					return t, vx + y.(float64)
				case complex128:
					return t, vx + y.(complex128)
				default:
					log.Fatalln("binaryOp + failed: unknown -", reflect.TypeOf(x))
				}
			case token.SUB:
				switch vx := x.(type) {
				case int:
					return t, vx - y.(int)
				case uint:
					return t, vx - y.(uint)
				case float64:
					return t, vx - y.(float64)
				case complex128:
					return t, vx - y.(complex128)
				default:
					log.Fatalln("binaryOp - failed: unknown -", reflect.TypeOf(x))
				}
			case token.MUL:
				switch vx := x.(type) {
				case int:
					return t, vx * y.(int)
				case uint:
					return t, vx * y.(uint)
				case float64:
					return t, vx * y.(float64)
				case complex128:
					return t, vx * y.(complex128)
				default:
					log.Fatalln("binaryOp * failed: unknown -", reflect.TypeOf(x))
				}
			case token.QUO:
				switch vx := x.(type) {
				case int:
					return t, vx / y.(int)
				case uint:
					return t, vx / y.(uint)
				case float64:
					return t, vx / y.(float64)
				case complex128:
					return t, vx / y.(complex128)
				default:
					log.Fatalln("binaryOp * failed: unknown -", reflect.TypeOf(x))
				}
			case token.REM:
				switch vx := x.(type) {
				case int:
					return t, vx % y.(int)
				case uint:
					return t, vx % y.(uint)
				default:
					log.Fatalln("binaryOp % failed: unknown -", reflect.TypeOf(x))
				}
			case token.EQL: // ==
				return Bool, x == y
			case token.LSS: // <
				switch vx := x.(type) {
				case int:
					return Bool, vx < y.(int)
				case string:
					return Bool, vx < y.(string)
				case uint:
					return Bool, vx < y.(uint)
				case float64:
					return Bool, vx < y.(float64)
				default:
					log.Fatalln("binaryOp < failed: unknown -", reflect.TypeOf(x))
				}
			case token.GTR: // >
				switch vx := x.(type) {
				case int:
					return Bool, vx > y.(int)
				case string:
					return Bool, vx > y.(string)
				case uint:
					return Bool, vx > y.(uint)
				case float64:
					return Bool, vx > y.(float64)
				default:
					log.Fatalln("binaryOp > failed: unknown -", reflect.TypeOf(x))
				}
			case token.NEQ: // !=
				return Bool, x != y
			case token.LEQ: // <=
				switch vx := x.(type) {
				case int:
					return Bool, vx <= y.(int)
				case string:
					return Bool, vx <= y.(string)
				case uint:
					return Bool, vx <= y.(uint)
				case float64:
					return Bool, vx <= y.(float64)
				default:
					log.Fatalln("binaryOp <= failed: unknown -", reflect.TypeOf(x))
				}
			case token.GEQ: // >=
				switch vx := x.(type) {
				case int:
					return Bool, vx >= y.(int)
				case string:
					return Bool, vx >= y.(string)
				case uint:
					return Bool, vx >= y.(uint)
				case float64:
					return Bool, vx >= y.(float64)
				default:
					log.Fatalln("binaryOp >= failed: unknown -", reflect.TypeOf(x))
				}
			case token.AND: // &
				switch vx := x.(type) {
				case int:
					return t, vx & y.(int)
				case uint:
					return t, vx & y.(uint)
				default:
					log.Fatalln("binaryOp & failed: unknown -", reflect.TypeOf(x))
				}
			case token.OR: // |
				switch vx := x.(type) {
				case int:
					return t, vx | y.(int)
				case uint:
					return t, vx | y.(uint)
				default:
					log.Fatalln("binaryOp | failed: unknown -", reflect.TypeOf(x))
				}
			case token.XOR: //  ^
				switch vx := x.(type) {
				case int:
					return t, vx ^ y.(int)
				case uint:
					return t, vx ^ y.(uint)
				default:
					log.Fatalln("binaryOp ^ failed: unknown -", reflect.TypeOf(x))
				}
			case token.AND_NOT: // &^
				switch vx := x.(type) {
				case int:
					return t, vx &^ y.(int)
				case uint:
					return t, vx &^ y.(uint)
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
		case int:
			switch vy := y.(type) {
			case int:
				return tx, vx << vy
			case uint:
				return tx, vx << vy
			default:
				log.Fatalln("binaryOp failed: int << unknown -", ty)
			}
		case uint:
			switch vy := y.(type) {
			case int:
				return tx, vx << vy
			case uint:
				return tx, vx << vy
			default:
				log.Fatalln("binaryOp failed: int << unknown -", ty)
			}
		default:
			log.Fatalln("binaryOp << failed: unknown -", tx, op, ty)
		}
	case token.SHR: // >>
		switch vx := x.(type) {
		case int:
			switch vy := y.(type) {
			case int:
				return tx, vx >> vy
			case uint:
				return tx, vx >> vy
			default:
				log.Fatalln("binaryOp failed: int >> unknown -", ty)
			}
		case uint:
			switch vy := y.(type) {
			case int:
				return tx, vx >> vy
			case uint:
				return tx, vx >> vy
			default:
				log.Fatalln("binaryOp failed: int >> unknown -", ty)
			}
		default:
			log.Fatalln("binaryOp >> failed: unknown -", tx, op, ty)
		}
	}
	log.Fatalln("binaryOp failed: unknown -", tx, op, ty)
	return nil, nil
}

// -----------------------------------------------------------------------------
