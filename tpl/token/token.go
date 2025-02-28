/*
 * Copyright (c) 2025 The GoPlus Authors (goplus.org). All rights reserved.
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

package token

import (
	"strconv"
)

// -----------------------------------------------------------------------------

type Token uint

func (tok Token) String() (s string) {
	if tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return
}

// Len returns
// 1) len of is token literal, if token is an operator.
// 2) 0 for else.
func (tok Token) Len() int {
	if tok > ' ' && tok <= Token(len(tokens)) {
		return len(tokens[tok])
	}
	return 0
}

// -----------------------------------------------------------------------------

const (
	ILLEGAL Token = iota
	EOF
	COMMENT

	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	IMAG   // 123.45i
	CHAR   // 'a'
	STRING // "abc"

	ADD = '+'
	SUB = '-'
	MUL = '*'
	QUO = '/'
	REM = '%'

	AND = '&'
	OR  = '|'
	XOR = '^'

	LT     = '<'
	GT     = '>'
	ASSIGN = '='
	NOT    = '!'

	LPAREN = '('
	LBRACK = '['
	LBRACE = '{'
	COMMA  = ','
	PERIOD = '.'

	RPAREN    = ')'
	RBRACK    = ']'
	RBRACE    = '}'
	SEMICOLON = ';'
	COLON     = ':'
	QUESTION  = '?'
	TILDE     = '~'
	AT        = '@'
)

const (
	operator_beg Token = 0x80 + iota

	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND  // &&
	LOR   // ||
	ARROW // <-
	INC   // ++
	DEC   // --

	EQ       // ==
	NE       // !=
	LE       // <=
	GE       // >=
	DEFINE   // :=
	ELLIPSIS // ...

	_ = operator_beg
)

// -----------------------------------------------------------------------------

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	IMAG:   "IMAG",
	CHAR:   "CHAR",
	STRING: "STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND: "&",
	OR:  "|",
	XOR: "^",

	LT:     "<",
	GT:     ">",
	ASSIGN: "=",
	NOT:    "!",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",
	QUESTION:  "?",
	TILDE:     "~",
	AT:        "@",

	SHL:     "<<",
	SHR:     ">>",
	AND_NOT: "&^",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	AND_ASSIGN:     "&=",
	OR_ASSIGN:      "|=",
	XOR_ASSIGN:     "^=",
	SHL_ASSIGN:     "<<=",
	SHR_ASSIGN:     ">>=",
	AND_NOT_ASSIGN: "&^=",

	LAND:  "&&",
	LOR:   "||",
	ARROW: "<-",
	INC:   "++",
	DEC:   "--",

	EQ:       "==",
	NE:       "!=",
	LE:       "<=",
	GE:       ">=",
	DEFINE:   ":=",
	ELLIPSIS: "...",
}

// ForEach iterates all tokens.
func ForEach(f func(tok Token, lit string)) {
	for i, s := range tokens {
		if s != "" {
			f(Token(i), s)
		}
	}
}

// -----------------------------------------------------------------------------
