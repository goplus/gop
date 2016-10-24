package tpl

import (
	"strconv"
)

// -----------------------------------------------------------------------------
// Token

const (
	ILLEGAL uint = iota
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
	operator_beg uint = 0x80 + iota

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

	USER_TOKEN_BEGIN uint = 0xb0
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

	ADD: "'+'",
	SUB: "'-'",
	MUL: "'*'",
	QUO: "'/'",
	REM: "'%'",

	AND: "'&'",
	OR:  "'|'",
	XOR: "'^'",

	LT:     "'<'",
	GT:     "'>'",
	ASSIGN: "'='",
	NOT:    "'!'",

	LPAREN: "'('",
	LBRACK: "'['",
	LBRACE: "'{'",
	COMMA:  "','",
	PERIOD: "'.'",

	RPAREN:    "')'",
	RBRACK:    "']'",
	RBRACE:    "'}'",
	SEMICOLON: "';'",
	COLON:     "':'",
	QUESTION:  "'?'",
	TILDE:     "'~'",
	AT:        "'@'",

	SHL:     "\"<<\"",
	SHR:     "\">>\"",
	AND_NOT: "\"&^\"",

	ADD_ASSIGN: "\"+=\"",
	SUB_ASSIGN: "\"-=\"",
	MUL_ASSIGN: "\"*=\"",
	QUO_ASSIGN: "\"/=\"",
	REM_ASSIGN: "\"%=\"",

	AND_ASSIGN:     "\"&=\"",
	OR_ASSIGN:      "\"|=\"",
	XOR_ASSIGN:     "\"^=\"",
	SHL_ASSIGN:     "\"<<=\"",
	SHR_ASSIGN:     "\">>=\"",
	AND_NOT_ASSIGN: "\"&^=\"",

	LAND:  "\"&&\"",
	LOR:   "\"||\"",
	ARROW: "\"<-\"",
	INC:   "\"++\"",
	DEC:   "\"--\"",

	EQ:       "\"==\"",
	NE:       "\"!=\"",
	LE:       "\"<=\"",
	GE:       "\">=\"",
	DEFINE:   "\":=\"",
	ELLIPSIS: "\"...\"",
}

// -----------------------------------------------------------------------------

// token => literal
//
func (p *Scanner) Ttol(tok uint) (s string) {

	if tok < uint(len(tokens)) {
		s = tokens[tok]
	} else {
		panic("unknown token: " + strconv.Itoa(int(tok)))
	}
	return
}

var rtokens map[string]uint

func init() {

	rtokens = make(map[string]uint, 64)
	for i, s := range tokens {
		if s != "" {
			rtokens[s] = uint(i)
		}
	}
}

// literal => token
//
func (p *Scanner) Ltot(lit string) (tok uint) {

	return rtokens[lit]
}

// -----------------------------------------------------------------------------
