package astutil

import (
	"go/ast"
	"go/token"
	"strconv"
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
