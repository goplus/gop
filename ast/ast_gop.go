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
/*
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
*/
// -----------------------------------------------------------------------------

// ErrWrapExpr represents `expr!`, `expr?` or `expr?: defaultValue`
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

// LambdaExpr represents
//    `(x, y, ...) => exprOrExprTuple`
//    `x => exprOrExprTuple`
//    `=> exprOrExprTuple`
// here exprOrExprTuple represents
//    `expr`
//    `(expr1, expr2, ...)`
type LambdaExpr struct {
	First       token.Pos
	Lhs         []*Ident
	Rarrow      token.Pos
	Rhs         []Expr
	Last        token.Pos
	LhsHasParen bool
	RhsHasParen bool
}

// LambdaExpr2 represents
//    `(x, y, ...) => { ... }`
//    `x => { ... }`
//    `=> { ... }`
type LambdaExpr2 struct {
	First       token.Pos
	Lhs         []*Ident
	Rarrow      token.Pos
	Body        *BlockStmt
	LhsHasParen bool
}

func (p *LambdaExpr) Pos() token.Pos {
	return p.First
}

func (p *LambdaExpr) End() token.Pos {
	return p.Last
}

func (p *LambdaExpr2) Pos() token.Pos {
	return p.First
}

func (p *LambdaExpr2) End() token.Pos {
	return p.Body.End()
}

func (*LambdaExpr) exprNode()  {}
func (*LambdaExpr2) exprNode() {}

// -----------------------------------------------------------------------------

// ForPhrase represents `for k, v <- container, cond`
type ForPhrase struct {
	For        token.Pos // position of "for" keyword
	Key, Value *Ident    // Key may be nil
	TokPos     token.Pos // position of "<-" operator
	X          Expr      // value to range over
	Init       Stmt      // initialization statement; or nil
	Cond       Expr      // value filter, can be nil
}

// Pos returns position of first character belonging to the node.
func (p *ForPhrase) Pos() token.Pos { return p.For }

// End returns position of first character immediately after the node.
func (p *ForPhrase) End() token.Pos { return p.X.End() }

func (p *ForPhrase) exprNode() {}

// ComprehensionExpr represents
//    `[vexpr for k1, v1 <- container1, cond1 ...]` or
//    `{vexpr for k1, v1 <- container1, cond1 ...}` or
//    `{kexpr: vexpr for k1, v1 <- container1, cond1 ...}` or
//    `{for k1, v1 <- container1, cond1 ...}` or
type ComprehensionExpr struct {
	Lpos token.Pos   // position of "[" or "{"
	Tok  token.Token // token.LBRACK '[' or token.LBRACE '{'
	Elt  Expr        // *KeyValueExpr or Expr or nil
	Fors []*ForPhrase
	Rpos token.Pos // position of "]" or "}"
}

// Pos - position of first character belonging to the node
func (p *ComprehensionExpr) Pos() token.Pos {
	return p.Lpos
}

// End - position of first character immediately after the node
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
