package exec

import (
	"go/ast"
	"reflect"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// MapType instr
func MapType(p *Builder, typ reflect.Type) *ast.MapType {
	key := Type(p, typ.Key())
	val := Type(p, typ.Elem())
	return &ast.MapType{Key: key, Value: val}
}

// ArrayType instr
func ArrayType(p *Builder, typ reflect.Type) *ast.ArrayType {
	var len ast.Expr
	var kind = typ.Kind()
	if kind == reflect.Array {
		len = IntConst(int64(typ.Len()))
	}
	return &ast.ArrayType{Len: len, Elt: Type(p, typ.Elem())}
}

// PtrType instr
func PtrType(p *Builder, typElem reflect.Type) ast.Expr {
	elt := Type(p, typElem)
	return &ast.StarExpr{X: elt}
}

// Type instr
func Type(p *Builder, typ reflect.Type) ast.Expr {
	pkgPath, name := typ.PkgPath(), typ.Name()
	log.Debug(typ, "-", "pkgPath:", pkgPath, "name:", name)
	if name != "" {
		if pkgPath != "" {
			pkg := p.Import(pkgPath)
			return &ast.SelectorExpr{X: Ident(pkg), Sel: Ident(name)}
		}
		return Ident(name)
	}
	kind := typ.Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		return ArrayType(p, typ)
	case reflect.Map:
		return MapType(p, typ)
	case reflect.Ptr:
		return PtrType(p, typ.Elem())
	case reflect.Func:
	case reflect.Chan:
	case reflect.Interface:
	case reflect.Struct:
	}
	log.Panicln("Type: unknown type -", typ)
	return nil
}

// -----------------------------------------------------------------------------
