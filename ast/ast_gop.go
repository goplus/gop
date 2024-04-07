/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ast

import (
	"github.com/goplus/gop/token"
)

// -----------------------------------------------------------------------------

// OverloadFuncDecl node represents an overload function declaration:
//
// `func name = (overloadFuncs)`
// `func (T).nameOrOp = (overloadFuncs)`
//
// here overloadFunc represents
//
// `func(params) {...}`
// `funcName`
// `(*T).methodName`
type OverloadFuncDecl struct {
	Doc      *CommentGroup // associated documentation; or nil
	Func     token.Pos     // position of "func" keyword
	Recv     *FieldList    // receiver (methods); or nil (functions)
	Name     *Ident        // function/method name
	Assign   token.Pos     // position of token "="
	Lparen   token.Pos     // position of "("
	Funcs    []Expr        // overload functions. here `Expr` can be *FuncLit, *Ident or *SelectorExpr
	Rparen   token.Pos     // position of ")"
	Operator bool          // is operator or not
	IsClass  bool          // recv set by class
}

// Pos - position of first character belonging to the node.
func (p *OverloadFuncDecl) Pos() token.Pos {
	return p.Func
}

// End - position of first character immediately after the node.
func (p *OverloadFuncDecl) End() token.Pos {
	return p.Rparen + 1
}

func (*OverloadFuncDecl) declNode() {}

// -----------------------------------------------------------------------------

// A BasicLit node represents a literal of basic type.
type BasicLit struct {
	ValuePos token.Pos    // literal position
	Kind     token.Token  // token.INT, token.FLOAT, token.IMAG, token.RAT, token.CHAR, token.STRING, token.CSTRING
	Value    string       // literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 3r, 'a', '\x7f', "foo" or `\m\n\o`
	Extra    *StringLitEx // optional (only available when Kind == token.STRING)
}

type StringLitEx struct {
	Parts []any // can be (val string) or (xval Expr)
}

// NextPartPos - position of first character of next part.
// pos - position of this part (not including quote character).
func NextPartPos(pos token.Pos, part any) (nextPos token.Pos) {
	switch v := part.(type) {
	case string: // normal string literal or end with "$$"
		return pos + token.Pos(len(v))
	case Expr:
		return v.End()
	}
	panic("NextPartPos: unexpected parameters")
}

// -----------------------------------------------------------------------------

// A EnvExpr node represents a ${name} expression.
type EnvExpr struct {
	TokPos token.Pos // position of "$"
	Lbrace token.Pos // position of "{"
	Name   *Ident    // name
	Rbrace token.Pos // position of "}"
}

// Pos - position of first character belonging to the node.
func (p *EnvExpr) Pos() token.Pos {
	return p.TokPos
}

// End - position of first character immediately after the node.
func (p *EnvExpr) End() token.Pos {
	if p.Rbrace != token.NoPos {
		return p.Rbrace
	}
	return p.Name.End()
}

// HasBrace checks is this EnvExpr ${name} or $name.
func (p *EnvExpr) HasBrace() bool {
	return p.Rbrace != token.NoPos
}

func (*EnvExpr) exprNode() {}

// -----------------------------------------------------------------------------

// A SliceLit node represents a slice literal.
type SliceLit struct {
	Lbrack     token.Pos // position of "["
	Elts       []Expr    // list of slice elements; or nil
	Rbrack     token.Pos // position of "]"
	Incomplete bool      // true if (source) expressions are missing in the Elts list
}

// Pos - position of first character belonging to the node.
func (p *SliceLit) Pos() token.Pos {
	return p.Lbrack
}

// End - position of first character immediately after the node.
func (p *SliceLit) End() token.Pos {
	return p.Rbrack + 1
}

func (*SliceLit) exprNode() {}

// -----------------------------------------------------------------------------

// A MatrixLit node represents a matrix literal.
type MatrixLit struct {
	Lbrack     token.Pos // position of "["
	Elts       [][]Expr  // list of matrix elements
	Rbrack     token.Pos // position of "]"
	Incomplete bool      // true if (source) expressions are missing in the Elts list
}

// Pos - position of first character belonging to the node.
func (p *MatrixLit) Pos() token.Pos {
	return p.Lbrack
}

// End - position of first character immediately after the node.
func (p *MatrixLit) End() token.Pos {
	return p.Rbrack + 1
}

func (*MatrixLit) exprNode() {}

// -----------------------------------------------------------------------------

// A ElemEllipsis node represents a matrix row elements.
type ElemEllipsis struct {
	Elt      Expr      // ellipsis element
	Ellipsis token.Pos // position of "..."
}

// Pos - position of first character belonging to the node.
func (p *ElemEllipsis) Pos() token.Pos {
	return p.Elt.Pos()
}

// End - position of first character immediately after the node.
func (p *ElemEllipsis) End() token.Pos {
	return p.Ellipsis + 3
}

func (*ElemEllipsis) exprNode() {}

// -----------------------------------------------------------------------------

// ErrWrapExpr represents `expr!`, `expr?` or `expr?: defaultValue`.
type ErrWrapExpr struct {
	X       Expr
	Tok     token.Token // ! or ?
	TokPos  token.Pos
	Default Expr // can be nil
}

// Pos - position of first character belonging to the node.
func (p *ErrWrapExpr) Pos() token.Pos {
	return p.X.Pos()
}

// End - position of first character immediately after the node.
func (p *ErrWrapExpr) End() token.Pos {
	if p.Default != nil {
		return p.Default.End()
	}
	return p.TokPos + 1
}

func (*ErrWrapExpr) exprNode() {}

// -----------------------------------------------------------------------------

// LambdaExpr represents one of the following expressions:
//
//	`(x, y, ...) => exprOrExprTuple`
//	`x => exprOrExprTuple`
//	`=> exprOrExprTuple`
//
// here exprOrExprTuple represents
//
//	`expr`
//	`(expr1, expr2, ...)`
type LambdaExpr struct {
	First       token.Pos
	Lhs         []*Ident
	Rarrow      token.Pos
	Rhs         []Expr
	Last        token.Pos
	LhsHasParen bool
	RhsHasParen bool
}

// LambdaExpr2 represents one of the following expressions:
//
//	`(x, y, ...) => { ... }`
//	`x => { ... }`
//	`=> { ... }`
type LambdaExpr2 struct {
	First       token.Pos
	Lhs         []*Ident
	Rarrow      token.Pos
	Body        *BlockStmt
	LhsHasParen bool
}

// Pos - position of first character belonging to the node.
func (p *LambdaExpr) Pos() token.Pos {
	return p.First
}

// End - position of first character immediately after the node.
func (p *LambdaExpr) End() token.Pos {
	return p.Last
}

// Pos - position of first character belonging to the node.
func (p *LambdaExpr2) Pos() token.Pos {
	return p.First
}

// End - position of first character immediately after the node.
func (p *LambdaExpr2) End() token.Pos {
	return p.Body.End()
}

func (*LambdaExpr) exprNode()  {}
func (*LambdaExpr2) exprNode() {}

// -----------------------------------------------------------------------------

// ForPhrase represents `for k, v <- container if init; cond` phrase.
type ForPhrase struct {
	For        token.Pos // position of "for" keyword
	Key, Value *Ident    // Key may be nil
	TokPos     token.Pos // position of "<-" operator
	X          Expr      // value to range over
	IfPos      token.Pos // position of if or comma; or NoPos
	Init       Stmt      // initialization statement; or nil
	Cond       Expr      // value filter, can be nil
}

// Pos returns position of first character belonging to the node.
func (p *ForPhrase) Pos() token.Pos { return p.For }

// End returns position of first character immediately after the node.
func (p *ForPhrase) End() token.Pos {
	if p.Cond != nil {
		return p.Cond.End()
	}
	return p.X.End()
}

func (p *ForPhrase) exprNode() {}

// ComprehensionExpr represents one of the following expressions:
//
//	`[vexpr for k1, v1 <- container1, cond1 ...]` or
//	`{vexpr for k1, v1 <- container1, cond1 ...}` or
//	`{kexpr: vexpr for k1, v1 <- container1, cond1 ...}` or
//	`{for k1, v1 <- container1, cond1 ...}` or
type ComprehensionExpr struct {
	Lpos token.Pos   // position of "[" or "{"
	Tok  token.Token // token.LBRACK '[' or token.LBRACE '{'
	Elt  Expr        // *KeyValueExpr or Expr or nil
	Fors []*ForPhrase
	Rpos token.Pos // position of "]" or "}"
}

// Pos - position of first character belonging to the node.
func (p *ComprehensionExpr) Pos() token.Pos {
	return p.Lpos
}

// End - position of first character immediately after the node.
func (p *ComprehensionExpr) End() token.Pos {
	return p.Rpos + 1
}

func (*ComprehensionExpr) exprNode() {}

// -----------------------------------------------------------------------------

// A ForPhraseStmt represents a for statement with a for <- clause.
type ForPhraseStmt struct {
	*ForPhrase
	Body *BlockStmt
}

// Pos - position of first character belonging to the node.
func (p *ForPhraseStmt) Pos() token.Pos {
	return p.For
}

// End - position of first character immediately after the node.
func (p *ForPhraseStmt) End() token.Pos {
	return p.Body.End()
}

func (*ForPhraseStmt) stmtNode() {}

// -----------------------------------------------------------------------------

// A RangeExpr node represents a range expression.
type RangeExpr struct {
	First  Expr      // start of composite elements; or nil
	To     token.Pos // position of ":"
	Last   Expr      // end of composite elements
	Colon2 token.Pos // position of ":" or token.NoPos
	Expr3  Expr      // step (or max) of composite elements; or nil
}

// Pos - position of first character belonging to the node.
func (p *RangeExpr) Pos() token.Pos {
	if p.First != nil {
		return p.First.Pos()
	}
	return p.To
}

// End - position of first character immediately after the node.
func (p *RangeExpr) End() token.Pos {
	if p.Expr3 != nil {
		return p.Expr3.End()
	}
	if p.Colon2 != token.NoPos {
		return p.Colon2 + 1
	}
	if p.Last != nil {
		return p.Last.End()
	}
	return p.To + 1
}

func (*RangeExpr) exprNode() {}

// -----------------------------------------------------------------------------
