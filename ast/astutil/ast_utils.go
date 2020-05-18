package astutil

import (
	"go/token"
	"strconv"

	"github.com/qiniu/qlang/v6/ast"
	"github.com/qiniu/qlang/v6/ast/spec"
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

// -----------------------------------------------------------------------------

// A ConstKind represents the specific kind of type that a Type represents.
// The zero Kind is not a valid kind.
type ConstKind = spec.ConstKind

const (
	// ConstBoundRune - bound type: rune
	ConstBoundRune = spec.ConstBoundRune
	// ConstBoundString - bound type: string
	ConstBoundString = spec.ConstBoundString
	// ConstUnboundInt - unbound int type
	ConstUnboundInt = spec.ConstUnboundInt
	// ConstUnboundFloat - unbound float type
	ConstUnboundFloat = spec.ConstUnboundFloat
	// ConstUnboundComplex - unbound complex type
	ConstUnboundComplex = spec.ConstUnboundComplex
)

// ToConst converts a ast.BasicLit to constant value.
func ToConst(v *ast.BasicLit) (ConstKind, interface{}) {
	switch v.Kind {
	case token.INT:
		n, err := strconv.ParseInt(v.Value, 0, 0)
		if err != nil {
			n2, err2 := strconv.ParseUint(v.Value, 0, 0)
			if err2 != nil {
				log.Fatalln("ToConst: strconv.ParseInt failed:", err2)
			}
			return ConstUnboundInt, n2
		}
		return ConstUnboundInt, n
	case token.CHAR, token.STRING:
		n, err := strconv.Unquote(v.Value)
		if err != nil {
			log.Fatalln("ToConst: strconv.Unquote failed:", err)
		}
		if v.Kind == token.CHAR {
			for _, c := range n {
				return ConstBoundRune, int64(c)
			}
			panic("not here")
		}
		return ConstBoundString, n
	case token.FLOAT:
		n, err := strconv.ParseFloat(v.Value, 64)
		if err != nil {
			log.Fatalln("ToConst: strconv.ParseFloat failed:", err)
		}
		return ConstUnboundFloat, n
	case token.IMAG: // 123.45i
		val := v.Value
		n, err := strconv.ParseFloat(val[:len(val)-1], 64)
		if err != nil {
			log.Fatalln("ToConst: strconv.ParseFloat failed:", err)
		}
		return ConstUnboundComplex, complex(0, n)
	}
	log.Fatalln("ToConst: unknown -", v)
	return 0, nil
}

// -----------------------------------------------------------------------------

// RecvInfo represents recv information of a method.
type RecvInfo struct {
	Name    string
	Type    string
	Pointer int
}

// ToRecv converts a ast.FieldList to recv information.
func ToRecv(recv *ast.FieldList) (ret RecvInfo) {
	fields := recv.List
	if len(fields) != 1 {
		panic("ToRecv: multi recv object?")
	}
	field := fields[0]
	if field.Names != nil {
		ret.Name = field.Names[0].Name
	}
	t := field.Type
retry:
	switch v := t.(type) {
	case *ast.Ident: // T
		ret.Type = v.Name
	case *ast.StarExpr: // *T
		ret.Pointer++
		t = v.X
		goto retry
	default:
		panic("ToRecv: recv can only be *T or T")
	}
	return
}

// -----------------------------------------------------------------------------
