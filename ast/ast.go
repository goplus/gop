package ast

import (
	"go/ast"
)

// -----------------------------------------------------------------------------

// A Package node represents a set of source files collectively building a qlang package.
type Package = ast.Package

// A File node represents a qlang source file.
type File = ast.File

// -----------------------------------------------------------------------------
