package cl

import (
	"log"
	"reflect"

	"github.com/qiniu/qlang/ast"
)

type iType = reflect.Type

func toTypeEx(ctx *fileCtx, typ ast.Expr) (t iType, variadic bool) {
	if t = toType(ctx, typ); t != nil {
		return
	}
	if v, ok := typ.(*ast.Ellipsis); ok {
		elem := toType(ctx, v.Elt)
		return reflect.SliceOf(elem), true
	}
	log.Panicln("toType: unknown -", reflect.TypeOf(typ))
	return nil, false
}

func toType(ctx *fileCtx, typ ast.Expr) iType {
	switch v := typ.(type) {
	case *ast.Ident:
		return toIdentType(ctx, v.Name)
	case *ast.SelectorExpr:
		return toExternalType(ctx, v)
	case *ast.StarExpr:
		elem := toType(ctx, v.X)
		return reflect.PtrTo(elem)
	case *ast.ArrayType:
		return toArrayType(ctx, v)
	case *ast.FuncType:
		return toFuncType(ctx, v)
	case *ast.InterfaceType:
		return toInterfaceType(ctx, v)
	case *ast.StructType:
		return toStructType(ctx, v)
	case *ast.MapType:
		key := toType(ctx, v.Key)
		elem := toType(ctx, v.Value)
		return reflect.MapOf(key, elem)
	case *ast.ChanType:
		val := toType(ctx, v.Value)
		return reflect.ChanOf(toChanDir(v.Dir), val)
	}
	return nil
}

func toChanDir(dir ast.ChanDir) (ret reflect.ChanDir) {
	if (dir & ast.SEND) != 0 {
		ret |= reflect.SendDir
	}
	if (dir & ast.RECV) != 0 {
		ret |= reflect.RecvDir
	}
	return
}

func toFuncType(ctx *fileCtx, t *ast.FuncType) iType {
	in, variadic := toTypesEx(ctx, t.Params)
	out := toTypes(ctx, t.Results)
	return reflect.FuncOf(in, out, variadic)
}

func toTypes(ctx *fileCtx, fields *ast.FieldList) (types []iType) {
	if fields == nil {
		return
	}
	for _, field := range fields.List {
		n := len(field.Names)
		if n == 0 {
			n = 1
		}
		typ := toType(ctx, field.Type)
		if typ == nil {
			log.Panicln("toType: unknown -", reflect.TypeOf(field.Type))
		}
		for i := 0; i < n; i++ {
			types = append(types, typ)
		}
	}
	return
}

func toTypesEx(ctx *fileCtx, fields *ast.FieldList) ([]iType, bool) {
	var types []iType
	last := len(fields.List) - 1
	for i := 0; i <= last; i++ {
		field := fields.List[i]
		n := len(field.Names)
		if n == 0 {
			n = 1
		}
		typ, variadic := toTypeEx(ctx, field.Type)
		for i := 0; i < n; i++ {
			types = append(types, typ)
		}
		if variadic {
			if i != last {
				log.Panicln("toTypes failed: the variadic type isn't last argument?")
			}
			return types, true
		}
	}
	return types, false
}

func toStructType(ctx *fileCtx, v *ast.StructType) iType {
	panic("toStructType: todo")
}

func toInterfaceType(ctx *fileCtx, v *ast.InterfaceType) iType {
	panic("toInterfaceType: todo")
}

func toExternalType(ctx *fileCtx, v *ast.SelectorExpr) iType {
	panic("toExternalType: todo")
}

func toIdentType(ctx *fileCtx, ident string) iType {
	panic("toIdentType: todo")
}

func toArrayType(ctx *fileCtx, v *ast.ArrayType) iType {
	panic("toArrayType: todo")
}

// -----------------------------------------------------------------------------

type typeDecl struct {
	Methods map[string]*methodDecl
	Alias   bool
}

type funcTypeDecl struct {
	X *ast.FuncType
}

type methodDecl struct {
	Recv    string // recv object name
	Pointer int
	Type    *funcTypeDecl
	Body    *ast.BlockStmt
	file    *fileCtx
}

type funcDecl struct {
	typ  *funcTypeDecl
	body *ast.BlockStmt
	file *fileCtx
	t    iType
}

func (p *funcDecl) typeOf() iType {
	if p.t == nil {
		p.t = toType(p.file, p.typ.X)
	}
	return p.t
}

// -----------------------------------------------------------------------------
