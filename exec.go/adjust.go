package exec

import (
	"go/ast"
	"go/token"
	"log"
	"reflect"
)

// -----------------------------------------------------------------------------

func adjustFile(f *ast.File, cmap ast.CommentMap) {
	pos := adjustIdent(f.Name, 1)
	adjustDecls(f.Decls, cmap, pos)
}

func adjustDecls(decls []ast.Decl, cmap ast.CommentMap, pos token.Pos) token.Pos {
	for _, decl := range decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			pos = adjustGenDecl(d, cmap, pos)
		case *ast.FuncDecl:
			pos = adjustFieldList(d.Recv, cmap, pos)
			pos = adjustIdent(d.Name, pos)
			pos = adjustFuncType(d.Type, cmap, pos)
			pos = adjustBlockStmt(d.Body, cmap, pos)
		default:
			log.Panicln("adjustDecls: todo -", reflect.TypeOf(decl))
		}
	}
	return pos
}

func adjustFieldList(flds *ast.FieldList, cmap ast.CommentMap, pos token.Pos) token.Pos {
	if flds == nil {
		return pos
	}
	if flds.Opening != 0 {
		flds.Opening = pos
		pos++
	}
	for _, fld := range flds.List {
		pos = adjustField(fld, cmap, pos)
	}
	if flds.Closing != 0 {
		flds.Closing = pos
		pos++
	}
	return pos
}

func adjustField(fld *ast.Field, cmap ast.CommentMap, pos token.Pos) token.Pos {
	for _, name := range fld.Names {
		pos = adjustIdent(name, pos)
	}
	pos = adjustType(fld.Type, cmap, pos)
	pos = adjustBasicLit(fld.Tag, pos)
	return pos
}

func adjustGenDecl(d *ast.GenDecl, cmap ast.CommentMap, pos token.Pos) token.Pos {
	d.TokPos = pos
	pos++
	d.Lparen = pos
	pos++
	for _, spec := range d.Specs {
		pos = adjustSpec(spec, cmap, pos)
	}
	d.Rparen = pos
	pos++
	return pos
}

func adjustSpec(spec ast.Spec, cmap ast.CommentMap, pos token.Pos) token.Pos {
	switch s := spec.(type) {
	case *ast.ImportSpec:
		pos = adjustIdent(s.Name, pos)
		pos = adjustBasicLit(s.Path, pos)
		s.EndPos = pos
	default:
		log.Panicln("adjustSpec: todo -", reflect.TypeOf(spec))
	}
	pos++
	return pos
}

func adjustBlockStmt(body *ast.BlockStmt, cmap ast.CommentMap, pos token.Pos) token.Pos {
	if body == nil {
		return pos
	}
	body.Lbrace = pos
	pos++
	for _, stmt := range body.List {
		if comments, ok := cmap[stmt]; ok {
			pos = adjustComments(comments, pos)
		}
		pos = adjustStmt(stmt, cmap, pos)
	}
	body.Rbrace = pos
	pos++
	return pos
}

func adjustComments(comments []*ast.CommentGroup, pos token.Pos) token.Pos {
	for _, cg := range comments {
		for _, c := range cg.List {
			c.Slash = pos
			pos += token.Pos(len(c.Text))
		}
	}
	return pos
}

func adjustStmt(stmt ast.Stmt, cmap ast.CommentMap, pos token.Pos) token.Pos {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		pos = adjustExpr(s.X, cmap, pos)
	case *ast.DeclStmt:
		pos = adjustGenDecl(s.Decl.(*ast.GenDecl), cmap, pos)
	default:
		log.Panicln("adjustStmt: todo -", reflect.TypeOf(stmt))
	}
	return pos
}

func adjustExpr(expr ast.Expr, cmap ast.CommentMap, pos token.Pos) token.Pos {
	switch e := expr.(type) {
	case *ast.Ident:
		pos = adjustIdent(e, pos)
	case *ast.BasicLit:
		pos = adjustBasicLit(e, pos)
	case *ast.CallExpr:
		pos = adjustExpr(e.Fun, cmap, pos)
		e.Lparen = pos
		pos++
		pos = adjustExprs(e.Args, cmap, pos)
		if e.Ellipsis != token.NoPos {
			e.Ellipsis = pos
			pos++
		}
		e.Rparen = pos
		pos++
	case *ast.BinaryExpr:
		pos = adjustExpr(e.X, cmap, pos)
		e.OpPos = pos
		pos++
		pos = adjustExpr(e.Y, cmap, pos)
	case *ast.SelectorExpr:
		pos = adjustExpr(e.X, cmap, pos)
		pos++
		e.Sel.NamePos = pos
		pos++
	case *ast.Ellipsis:
		e.Ellipsis = pos
		pos++
	default:
		log.Panicln("adjustExpr: todo -", reflect.TypeOf(expr))
	}
	return pos
}

func adjustBasicLit(lit *ast.BasicLit, pos token.Pos) token.Pos {
	if lit != nil {
		lit.ValuePos = pos
		pos++
	}
	return pos
}

func adjustIdent(ident *ast.Ident, pos token.Pos) token.Pos {
	if ident != nil {
		ident.NamePos = pos
		pos += token.Pos(len(ident.Name))
	}
	return pos
}

func adjustExprs(exprs []ast.Expr, cmap ast.CommentMap, pos token.Pos) token.Pos {
	for _, expr := range exprs {
		pos = adjustExpr(expr, cmap, pos)
	}
	return pos
}

// -----------------------------------------------------------------------------

func adjustType(expr ast.Expr, cmap ast.CommentMap, pos token.Pos) token.Pos {
	log.Panicln("adjustType: todo -", reflect.TypeOf(expr))
	return pos
}

func adjustFuncType(t *ast.FuncType, cmap ast.CommentMap, pos token.Pos) token.Pos {
	t.Func = pos
	pos++
	pos = adjustFieldList(t.Params, cmap, pos)
	pos = adjustFieldList(t.Results, cmap, pos)
	return pos
}

// -----------------------------------------------------------------------------
