package token

import (
	"go/token"
)

// -----------------------------------------------------------------------------

// Token is the set of lexical tokens of the qlang.
type Token = token.Token

// The list of tokens.
const (
	IDENT  = token.IDENT  // main
	INT    = token.INT    // 12345
	FLOAT  = token.FLOAT  // 123.45
	IMAG   = token.IMAG   // 123.45i
	CHAR   = token.CHAR   // 'a'
	STRING = token.STRING // "abc"

	ADD = token.ADD // +
	SUB = token.SUB // -
	MUL = token.MUL // *
	QUO = token.QUO // /
	REM = token.REM // %

	AND     = token.AND     // &
	OR      = token.OR      // |
	XOR     = token.XOR     // ^
	SHL     = token.SHL     // <<
	SHR     = token.SHR     // >>
	AND_NOT = token.AND_NOT // &^

	ADD_ASSIGN = token.ADD_ASSIGN // +=
	SUB_ASSIGN = token.SUB_ASSIGN // -=
	MUL_ASSIGN = token.MUL_ASSIGN // *=
	QUO_ASSIGN = token.QUO_ASSIGN // /=
	REM_ASSIGN = token.REM_ASSIGN // %=

	AND_ASSIGN     = token.AND_ASSIGN     // &=
	OR_ASSIGN      = token.OR_ASSIGN      // |=
	XOR_ASSIGN     = token.XOR_ASSIGN     // ^=
	SHL_ASSIGN     = token.SHL_ASSIGN     // <<=
	SHR_ASSIGN     = token.SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN = token.AND_NOT_ASSIGN // &^=

	LAND  = token.LAND  // &&
	LOR   = token.LOR   // ||
	ARROW = token.ARROW // <-
	INC   = token.INC   // ++
	DEC   = token.DEC   // --

	EQL    = token.EQL    // ==
	LSS    = token.LSS    // <
	GTR    = token.GTR    // >
	ASSIGN = token.ASSIGN // =
	NOT    = token.NOT    // !

	NEQ      = token.NEQ      // !=
	LEQ      = token.LEQ      // <=
	GEQ      = token.GEQ      // >=
	DEFINE   = token.DEFINE   // :=
	ELLIPSIS = token.ELLIPSIS // ...

	LPAREN = token.LPAREN // (
	LBRACK = token.LBRACK // [
	LBRACE = token.LBRACE // {
	COMMA  = token.COMMA  // ,
	PERIOD = token.PERIOD // .

	RPAREN    = token.RPAREN    // )
	RBRACK    = token.RBRACK    // ]
	RBRACE    = token.RBRACE    // }
	SEMICOLON = token.SEMICOLON // ;
	COLON     = token.COLON     // :

	// Keywords
	BREAK    = token.BREAK
	CASE     = token.CASE
	CHAN     = token.CHAN
	CONST    = token.CONST
	CONTINUE = token.CONTINUE

	DEFAULT     = token.DEFAULT
	DEFER       = token.DEFER
	ELSE        = token.ELSE
	FALLTHROUGH = token.FALLTHROUGH
	FOR         = token.FOR

	FUNC   = token.FUNC
	GO     = token.GO
	GOTO   = token.GOTO
	IF     = token.IF
	IMPORT = token.IMPORT

	INTERFACE = token.INTERFACE
	MAP       = token.MAP
	PACKAGE   = token.PACKAGE
	RANGE     = token.RANGE
	RETURN    = token.RETURN

	SELECT = token.SELECT
	STRUCT = token.STRUCT
	SWITCH = token.SWITCH
	TYPE   = token.TYPE
	VAR    = token.VAR
)

// -----------------------------------------------------------------------------

// Pos is a compact encoding of a source position within a file set.
// It can be converted into a Position for a more convenient, but much
// larger, representation.
//
// The Pos value for a given file is a number in the range [base, base+size],
// where base and size are specified when adding the file to the file set via
// AddFile.
//
// To create the Pos value for a specific source offset (measured in bytes),
// first add the respective file to the current file set using FileSet.AddFile
// and then call File.Pos(offset) for that file. Given a Pos value p
// for a specific file set fset, the corresponding Position value is
// obtained by calling fset.Position(p).
//
// Pos values can be compared directly with the usual comparison operators:
// If two Pos values p and q are in the same file, comparing p and q is
// equivalent to comparing the respective source file offsets. If p and q
// are in different files, p < q is true if the file implied by p was added
// to the respective file set before the file implied by q.
//
type Pos = token.Pos

// A File is a handle for a file belonging to a FileSet.
// A File has a name, size, and line offset table.
type File = token.File

// A FileSet represents a set of source files. Methods of file sets are synchronized;
// multiple goroutines may invoke them concurrently.
type FileSet = token.FileSet

// NewFileSet creates a new file set.
func NewFileSet() *FileSet {
	return token.NewFileSet()
}

// -----------------------------------------------------------------------------
