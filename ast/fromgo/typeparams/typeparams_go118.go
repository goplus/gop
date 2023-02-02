//go:build go1.18
// +build go1.18

package typeparams

import "go/ast"

// IndexListExpr is an alias for ast.IndexListExpr.
type IndexListExpr = ast.IndexListExpr

// ForFuncType returns n.TypeParams.
func ForFuncType(n *ast.FuncType) *ast.FieldList {
	if n == nil {
		return nil
	}
	return n.TypeParams
}

// ForTypeSpec returns n.TypeParams.
func ForTypeSpec(n *ast.TypeSpec) *ast.FieldList {
	if n == nil {
		return nil
	}
	return n.TypeParams
}
