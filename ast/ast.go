package ast

import (
	"go/ast"
)

// -----------------------------------------------------------------------------

// A Package node represents a set of source files collectively building a qlang package.
type Package = ast.Package

// A File node represents a qlang source file.
type File = ast.File

// A Spec node represents a single (non-parenthesized) import,
// constant, type, or variable declaration.
type (
	// The Spec type stands for any of *ImportSpec, *ValueSpec, and *TypeSpec.
	Spec = ast.Spec

	// An ImportSpec node represents a single package import.
	ImportSpec = ast.ImportSpec

	// A ValueSpec node represents a constant or variable declaration
	// (ConstSpec or VarSpec production).
	ValueSpec = ast.ValueSpec

	// A TypeSpec node represents a type declaration (TypeSpec production).
	TypeSpec = ast.TypeSpec
)

// A declaration is represented by one of the following declaration nodes.
type (
	// A GenDecl node (generic declaration node) represents an import,
	// constant, type or variable declaration. A valid Lparen position
	// (Lparen.IsValid()) indicates a parenthesized declaration.
	//
	// Relationship between Tok value and Specs element type:
	//
	//	token.IMPORT  *ImportSpec
	//	token.CONST   *ValueSpec
	//	token.TYPE    *TypeSpec
	//	token.VAR     *ValueSpec
	//
	GenDecl = ast.GenDecl

	// A FuncDecl node represents a function declaration.
	FuncDecl = ast.FuncDecl
)

// -----------------------------------------------------------------------------
