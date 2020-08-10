package cl

import (
	"testing"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/qiniu/x/ts"

	exec "github.com/goplus/gop/exec/bytecode"
)

// -----------------------------------------------------------------------------

var instrDeferDiscardsResult = map[string]goInstrInfo{
	"len":     {igoLen},
	"cap":     {igoCap},
	"make":    {igoMake},
	"new":     {igoNew},
	"complex": {igoComplex},
	"real":    {igoReal},
	"imag":    {igoImag},
}

func TestDeferFileNotFound(t *testing.T) {
	b := exec.NewBuilder(nil)
	pkg := &ast.Package{
		Files: map[string]*ast.File{},
	}
	pkgCtx := newPkgCtx(b.Interface(), pkg, token.NewFileSet())
	ctx := newGblBlockCtx(pkgCtx)
	ts.New(t).Call(func() {
		fn := &ast.Ident{Name: "len"}
		expr := &ast.CallExpr{Fun: fn}
		igoLen(ctx, expr, callByDefer)
	}).Panic(
		"pkgCtx.getCodeInfo failed: file not found - \n",
	)
}

func TestDeferDiscardsResult(t *testing.T) {
	for k, v := range instrDeferDiscardsResult {
		b := exec.NewBuilder(nil)
		file := &ast.File{
			Code: []byte(k + "()\n"),
		}
		pkg := &ast.Package{
			Files: map[string]*ast.File{
				"bar.gop": file,
			},
		}
		fset := token.NewFileSet()
		fset.AddFile("bar.gop", 1, 100)
		pkgCtx := newPkgCtx(b.Interface(), pkg, fset)
		ctx := newGblBlockCtx(pkgCtx)
		ts.New(t).Call(func() {
			fn := &ast.Ident{NamePos: 1, Name: k}
			expr := &ast.CallExpr{Fun: fn, Rparen: token.Pos(len(k) + 2)}
			v.instr(ctx, expr, callByDefer)
		}).Panic(
			"defer discards result of " + k + "()\n",
		)
	}
}

// -----------------------------------------------------------------------------
