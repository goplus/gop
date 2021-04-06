/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package ast

import (
	"github.com/goplus/gop/token"
)

// ----------------------------------------------------------------------------
/*
// A Scope maintains the set of named language entities declared
// in the scope and a link to the immediately surrounding (outer)
// scope.
type Scope = ast.Scope

// A Package node represents a set of source files collectively building a Go+ package.
type Package = ast.Package

// A File node represents a Go+ source file.
type File = ast.File

// Expr - All expression nodes implement the Expr interface.
type Expr = ast.Expr

// Stmt - All statement nodes implement the Stmt interface.
type Stmt = ast.Stmt

// ----------------------------------------------------------------------------
// Declarations

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

// ----------------------------------------------------------------------------
// Comments

// A CommentGroup represents a sequence of comments
// with no other tokens and no empty lines between.
type CommentGroup = ast.CommentGroup

// ----------------------------------------------------------------------------
// Statements

// A statement is represented by a tree consisting of one
// or more of the following concrete statement nodes.
type (
	// A BadStmt node is a placeholder for statements containing
	// syntax errors for which no correct statement nodes can be
	// created.
	BadStmt = ast.BadStmt

	// A DeclStmt node represents a declaration in a statement list.
	DeclStmt = ast.DeclStmt

	// An EmptyStmt node represents an empty statement.
	// The "position" of the empty statement is the position
	// of the immediately following (explicit or implicit) semicolon.
	EmptyStmt = ast.EmptyStmt

	// A LabeledStmt node represents a labeled statement.
	LabeledStmt = ast.LabeledStmt

	// An ExprStmt node represents a (stand-alone) expression
	// in a statement list.
	ExprStmt = ast.ExprStmt

	// A SendStmt node represents a send statement.
	SendStmt = ast.SendStmt

	// An IncDecStmt node represents an increment or decrement statement.
	IncDecStmt = ast.IncDecStmt

	// An AssignStmt node represents an assignment or
	// a short variable declaration.
	AssignStmt = ast.AssignStmt

	// A GoStmt node represents a go statement.
	GoStmt = ast.GoStmt

	// A DeferStmt node represents a defer statement.
	DeferStmt = ast.DeferStmt

	// A ReturnStmt node represents a return statement.
	ReturnStmt = ast.ReturnStmt

	// A BranchStmt node represents a break, continue, goto,
	// or fallthrough statement.
	BranchStmt = ast.BranchStmt

	// A BlockStmt node represents a braced statement list.
	BlockStmt = ast.BlockStmt

	// An IfStmt node represents an if statement.
	IfStmt = ast.IfStmt

	// A CaseClause represents a case of an expression or type switch statement.
	CaseClause = ast.CaseClause

	// A SwitchStmt node represents an expression switch statement.
	SwitchStmt = ast.SwitchStmt

	// A TypeSwitchStmt node represents a type switch statement.
	TypeSwitchStmt = ast.TypeSwitchStmt

	// A CommClause node represents a case of a select statement.
	CommClause = ast.CommClause

	// A SelectStmt node represents a select statement.
	SelectStmt = ast.SelectStmt

	// A ForStmt represents a for statement.
	ForStmt = ast.ForStmt

	// A RangeStmt represents a for statement with a range clause.
	RangeStmt = ast.RangeStmt
)

// ----------------------------------------------------------------------------
// Expressions and types

// A Field represents a Field declaration list in a struct type,
// a method list in an interface type, or a parameter/result declaration
// in a signature.
// Field.Names is nil for unnamed parameters (parameter lists which only contain types)
// and embedded struct fields. In the latter case, the field name is the type name.
//
type Field = ast.Field

// A FieldList represents a list of Fields, enclosed by parentheses or braces.
type FieldList = ast.FieldList

// An expression is represented by a tree consisting of one
// or more of the following concrete expression nodes.
//
type (
	// A BadExpr node is a placeholder for expressions containing
	// syntax errors for which no correct expression nodes can be
	// created.
	BadExpr = ast.BadExpr

	// An Ident node represents an identifier.
	Ident = ast.Ident

	// An Ellipsis node stands for the "..." type in a
	// parameter list or the "..." length in an array type.
	Ellipsis = ast.Ellipsis

	// A BasicLit node represents a literal of basic type.
	BasicLit = ast.BasicLit

	// A FuncLit node represents a function literal.
	FuncLit = ast.FuncLit

	// A CompositeLit node represents a composite literal.
	CompositeLit = ast.CompositeLit

	// A ParenExpr node represents a parenthesized expression.
	ParenExpr = ast.ParenExpr

	// A SelectorExpr node represents an expression followed by a selector.
	SelectorExpr = ast.SelectorExpr

	// An IndexExpr node represents an expression followed by an index.
	IndexExpr = ast.IndexExpr

	// A SliceExpr node represents an expression followed by slice indices.
	SliceExpr = ast.SliceExpr

	// A TypeAssertExpr node represents an expression followed by a
	// type assertion.
	TypeAssertExpr = ast.TypeAssertExpr

	// A CallExpr node represents an expression followed by an argument list.
	CallExpr = ast.CallExpr

	// A StarExpr node represents an expression of the form "*" Expression.
	// Semantically it could be a unary "*" expression, or a pointer type.
	StarExpr = ast.StarExpr

	// A UnaryExpr node represents a unary expression.
	// Unary "*" expressions are represented via StarExpr nodes.
	UnaryExpr = ast.UnaryExpr

	// A BinaryExpr node represents a binary expression.
	BinaryExpr = ast.BinaryExpr

	// A KeyValueExpr node represents (key : value) pairs
	// in composite literals.
	KeyValueExpr = ast.KeyValueExpr
)

// ChanDir - the direction of a channel type is indicated by a bit
// mask including one or both of the following constants.
type ChanDir = ast.ChanDir

const (
	// SEND flag
	SEND = ast.SEND
	// RECV flag
	RECV = ast.RECV
)

// A type is represented by a tree consisting of one
// or more of the following type-specific expression
// nodes.
type (
	// An ArrayType node represents an array or slice type.
	ArrayType = ast.ArrayType

	// A StructType node represents a struct type.
	StructType = ast.StructType

	// Pointer types are represented via StarExpr nodes.

	// A FuncType node represents a function type.
	FuncType = ast.FuncType

	// An InterfaceType node represents an interface type.
	InterfaceType = ast.InterfaceType

	// A MapType node represents a map type.
	MapType = ast.MapType

	// A ChanType node represents a channel type.
	ChanType = ast.ChanType
)
*/
// -----------------------------------------------------------------------------

// A SliceLit node represents a slice literal.
type SliceLit struct {
	Lbrack     token.Pos // position of "["
	Elts       []Expr    // list of composite elements; or nil
	Rbrack     token.Pos // position of "]"
	Incomplete bool      // true if (source) expressions are missing in the Elts list
}

// Pos - position of first character belonging to the node
func (p *SliceLit) Pos() token.Pos {
	return p.Lbrack
}

// End - position of first character immediately after the node
func (p *SliceLit) End() token.Pos {
	return p.Rbrack + 1
}

func (*SliceLit) exprNode() {}

// -----------------------------------------------------------------------------

// TernaryExpr represents `cond ? expr1 : expr2`
type TernaryExpr struct {
	Cond     Expr
	Question token.Pos
	X        Expr
	Colon    token.Pos
	Y        Expr
}

// Pos - position of first character belonging to the node
func (p *TernaryExpr) Pos() token.Pos {
	return p.Cond.Pos()
}

// End - position of first character immediately after the node
func (p *TernaryExpr) End() token.Pos {
	return p.Y.End()
}

func (*TernaryExpr) exprNode() {}

// -----------------------------------------------------------------------------

// ErrWrapExpr represents `expr!`, `expr?` or `expr? defaultValue`
type ErrWrapExpr struct {
	X       Expr
	Tok     token.Token // ! or ?
	TokPos  token.Pos
	Default Expr // can be nil
}

// Pos - position of first character belonging to the node
func (p *ErrWrapExpr) Pos() token.Pos {
	return p.X.Pos()
}

// End - position of first character immediately after the node
func (p *ErrWrapExpr) End() token.Pos {
	if p.Default != nil {
		return p.Default.End()
	}
	return p.TokPos + 1
}

func (*ErrWrapExpr) exprNode() {}

// -----------------------------------------------------------------------------

// ForPhrase represents `for k, v <- listOrMap`
type ForPhrase struct {
	For        token.Pos // position of "for" keyword
	Key, Value *Ident    // Key may be nil
	TokPos     token.Pos // position of "<-" operator
	X          Expr      // value to range over, must be list or map
	Cond       Expr      // value filter, can be nil
}

func (p *ForPhrase) Pos() token.Pos {
	return p.For
}

func (p *ForPhrase) End() token.Pos {
	return p.X.End()
}

func (p *ForPhrase) exprNode() {}

// ListComprehensionExpr represents `[expr for k1, v1 <- listOrMap1, cond1 ...]`
type ListComprehensionExpr struct {
	Lbrack token.Pos // position of "["
	Elt    Expr
	Fors   []ForPhrase
	Rbrack token.Pos // position of "]"
}

// Pos - position of first character belonging to the node
func (p *ListComprehensionExpr) Pos() token.Pos {
	return p.Lbrack
}

// End - position of first character immediately after the node
func (p *ListComprehensionExpr) End() token.Pos {
	return p.Rbrack + 1
}

func (*ListComprehensionExpr) exprNode() {}

// -----------------------------------------------------------------------------

// MapComprehensionExpr represents `{kexpr: vexpr for k1, v1 <- listOrMap1, cond1 ...}`
type MapComprehensionExpr struct {
	Lbrace token.Pos // position of "{"
	Elt    *KeyValueExpr
	Fors   []ForPhrase
	Rbrace token.Pos // position of "}"
}

// Pos - position of first character belonging to the node
func (p *MapComprehensionExpr) Pos() token.Pos {
	return p.Lbrace
}

// End - position of first character immediately after the node
func (p *MapComprehensionExpr) End() token.Pos {
	return p.Rbrace + 1
}

func (*MapComprehensionExpr) exprNode() {}

// -----------------------------------------------------------------------------

// A ForPhraseStmt represents a for statement with a for <- clause.
type ForPhraseStmt struct {
	ForPhrase
	Body *BlockStmt
}

// Pos - position of first character belonging to the node
func (p *ForPhraseStmt) Pos() token.Pos {
	return p.For
}

// End - position of first character immediately after the node
func (p *ForPhraseStmt) End() token.Pos {
	return p.Body.End()
}

func (*ForPhraseStmt) stmtNode() {}

// -----------------------------------------------------------------------------

// A TwoValueIndexExpr node represents a two-value assignment expression (v, ok := m["key"])
type TwoValueIndexExpr struct {
	*IndexExpr
}

type TwoValueTypeAssertExpr struct {
	*TypeAssertExpr
}

// -----------------------------------------------------------------------------
