package tpl

import (
	"bytes"
	"fmt"
	"go/token"
	"path/filepath"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// A Token is a lexical unit returned by Scan.
//
type Token struct {
	Kind    uint
	Pos     token.Pos
	Literal string
}

// An ScanErrorHandler may be provided to Scanner.Init. If a syntax error is
// encountered and a handler was installed, the handler is called with a
// position and an error message. The position points to the beginning of
// the offending token.
//
type ScanErrorHandler func(pos token.Position, msg string)

// A Scanner holds the scanner's internal state while processing
// a given text.  It can be allocated as part of another data
// structure but must be initialized via Init before use.
//
type Scanner struct {
	// immutable state
	file *token.File      // source file handle
	dir  string           // directory portion of file.Name()
	src  []byte           // source
	err  ScanErrorHandler // error reporting; or nil
	mode ScanMode         // scanning mode

	// scanning state
	ch         rune // current character
	offset     int  // character offset
	rdOffset   int  // reading offset (position after current character)
	lineOffset int  // current line offset
	insertSemi bool // insert a semicolon before next newline

	// public state - ok to modify
	ErrorCount int // number of errors encountered
}

const bom = 0xFEFF // byte order mark, only permitted as very first character

// Read the next Unicode char into s.ch.
// s.ch < 0 means end-of-file.
//
func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.file.AddLine(s.offset)
		}
		r, w := rune(s.src[s.rdOffset]), 1
		switch {
		case r == 0:
			s.error(s.offset, "illegal character NUL")
		case r >= 0x80:
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				s.error(s.offset, "illegal UTF-8 encoding")
			} else if r == bom && s.offset > 0 {
				s.error(s.offset, "illegal byte order mark")
			}
		}
		s.rdOffset += w
		s.ch = r
	} else {
		s.offset = len(s.src)
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.file.AddLine(s.offset)
		}
		s.ch = -1 // eof
	}
}

// A ScanMode value is a set of flags (or 0).
// They control scanner behavior.
//
type ScanMode uint

const (
	// ScanComments means returning comments as COMMENT tokens
	ScanComments ScanMode = 1 << iota

	// InsertSemis means automatically insert semicolons
	InsertSemis
)

// Init prepares the scanner s to tokenize the text src by setting the
// scanner at the beginning of src. The scanner uses the file set file
// for position information and it adds line information for each line.
// It is ok to re-use the same file when re-scanning the same file as
// line information which is already present is ignored. Init causes a
// panic if the file size does not match the src size.
//
// Calls to Scan will invoke the error handler err if they encounter a
// syntax error and err is not nil. Also, for each error encountered,
// the Scanner field ErrorCount is incremented by one. The mode parameter
// determines how comments are handled.
//
// Note that Init may call err if there is an error in the first character
// of the file.
//
func (s *Scanner) Init(file *token.File, src []byte, err ScanErrorHandler, mode ScanMode) {
	// Explicitly initialize all fields since a scanner may be reused.
	if file.Size() != len(src) {
		panic(fmt.Sprintf("file size (%d) does not match src len (%d)", file.Size(), len(src)))
	}
	s.file = file
	s.dir, _ = filepath.Split(file.Name())
	s.src = src
	s.err = err
	s.mode = mode

	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.lineOffset = 0
	s.insertSemi = false
	s.ErrorCount = 0

	s.next()
	if s.ch == bom {
		s.next() // ignore BOM at file beginning
	}
}

func (s *Scanner) error(offs int, msg string) {
	if s.err != nil {
		s.err(s.file.Position(s.file.Pos(offs)), msg)
	}
	s.ErrorCount++
}

var prefix = []byte("//line ")

func (s *Scanner) interpretLineComment(text []byte) {
	if bytes.HasPrefix(text, prefix) {
		// get filename and line number, if any
		if i := bytes.LastIndex(text, []byte{':'}); i > 0 {
			if line, err := strconv.Atoi(string(text[i+1:])); err == nil && line > 0 {
				// valid //line filename:line comment
				filename := string(bytes.TrimSpace(text[len(prefix):i]))
				if filename != "" {
					filename = filepath.Clean(filename)
					if !filepath.IsAbs(filename) {
						// make filename relative to current directory
						filename = filepath.Join(s.dir, filename)
					}
				}
				// update scanner position
				s.file.AddLineInfo(
					s.lineOffset+len(text)+1, filename, line) // +len(text)+1 since comment applies to next line
			}
		}
	}
}

func (s *Scanner) scanComment() string {
	// initial '/' already consumed; s.ch == '/' || s.ch == '*'
	offs := s.offset - 1 // position of initial '/'
	hasCR := false

	if s.ch == '/' {
		//-style comment
		s.next()
		for s.ch != '\n' && s.ch >= 0 {
			if s.ch == '\r' {
				hasCR = true
			}
			s.next()
		}
		if offs == s.lineOffset {
			// comment starts at the beginning of the current line
			s.interpretLineComment(s.src[offs:s.offset])
		}
		goto exit
	}

	/*-style comment */
	s.next()
	for s.ch >= 0 {
		ch := s.ch
		if ch == '\r' {
			hasCR = true
		}
		s.next()
		if ch == '*' && s.ch == '/' {
			s.next()
			goto exit
		}
	}

	s.error(offs, "comment not terminated")

exit:
	lit := s.src[offs:s.offset]
	if hasCR {
		lit = stripCR(lit)
	}

	return string(lit)
}

func (s *Scanner) findLineEnd() bool {
	// initial '/' already consumed

	defer func(offs int) {
		// reset scanner state to where it was upon calling findLineEnd
		s.ch = '/'
		s.offset = offs
		s.rdOffset = offs + 1
		s.next() // consume initial '/' again
	}(s.offset - 1)

	// read ahead until a newline, EOF, or non-comment token is found
	for s.ch == '/' || s.ch == '*' {
		if s.ch == '/' {
			//-style comment always contains a newline
			return true
		}
		/*-style comment: look for newline */
		s.next()
		for s.ch >= 0 {
			ch := s.ch
			if ch == '\n' {
				return true
			}
			s.next()
			if ch == '*' && s.ch == '/' {
				s.next()
				break
			}
		}
		s.skipWhitespace() // s.insertSemi is set
		if s.ch < 0 || s.ch == '\n' {
			return true
		}
		if s.ch != '/' {
			// non-comment token
			return false
		}
		s.next() // consume '/'
	}

	return false
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= 0x80 && unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' || ch >= 0x80 && unicode.IsDigit(ch)
}

func (s *Scanner) scanIdentifier() string {
	offs := s.offset
	for isLetter(s.ch) || isDigit(s.ch) {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	}
	return 16 // larger than any legal digit val
}

func (s *Scanner) scanMantissa(base int) {
	for digitVal(s.ch) < base {
		s.next()
	}
}

func (s *Scanner) scanNumber(seenDecimalPoint bool) (uint, string) {
	// digitVal(s.ch) < 10
	offs := s.offset
	tok := INT

	if seenDecimalPoint {
		offs--
		tok = FLOAT
		s.scanMantissa(10)
		goto exponent
	}

	if s.ch == '0' {
		// int or float
		offs := s.offset
		s.next()
		if s.ch == 'x' || s.ch == 'X' {
			// hexadecimal int
			s.next()
			s.scanMantissa(16)
			if s.offset-offs <= 2 {
				// only scanned "0x" or "0X"
				s.error(offs, "illegal hexadecimal number")
			}
		} else {
			// octal int or float
			seenDecimalDigit := false
			s.scanMantissa(8)
			if s.ch == '8' || s.ch == '9' {
				// illegal octal int or float
				seenDecimalDigit = true
				s.scanMantissa(10)
			}
			if s.ch == '.' || s.ch == 'e' || s.ch == 'E' || s.ch == 'i' {
				goto fraction
			}
			// octal int
			if seenDecimalDigit {
				s.error(offs, "illegal octal number")
			}
		}
		goto exit
	}

	// decimal int or float
	s.scanMantissa(10)

fraction:
	if s.ch == '.' {
		tok = FLOAT
		s.next()
		s.scanMantissa(10)
	}

exponent:
	if s.ch == 'e' || s.ch == 'E' {
		tok = FLOAT
		s.next()
		if s.ch == '-' || s.ch == '+' {
			s.next()
		}
		s.scanMantissa(10)
	}

	if s.ch == 'i' {
		tok = IMAG
		s.next()
	}

exit:
	return tok, string(s.src[offs:s.offset])
}

// scanEscape parses an escape sequence where rune is the accepted
// escaped quote. In case of a syntax error, it stops at the offending
// character (without consuming it) and returns false. Otherwise
// it returns true.
func (s *Scanner) scanEscape(quote rune) bool {
	offs := s.offset

	var n int
	var base, max uint32
	switch s.ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		s.next()
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		s.next()
		n, base, max = 2, 16, 255
	case 'u':
		s.next()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		s.next()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		msg := "unknown escape sequence"
		if s.ch < 0 {
			msg = "escape sequence not terminated"
		}
		s.error(offs, msg)
		return false
	}

	var x uint32
	for n > 0 {
		d := uint32(digitVal(s.ch))
		if d >= base {
			msg := fmt.Sprintf("illegal character %#U in escape sequence", s.ch)
			if s.ch < 0 {
				msg = "escape sequence not terminated"
			}
			s.error(s.offset, msg)
			return false
		}
		x = x*base + d
		s.next()
		n--
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		s.error(offs, "escape sequence is invalid Unicode code point")
		return false
	}

	return true
}

func (s *Scanner) scanRune() string {
	// '\'' opening already consumed
	offs := s.offset - 1

	valid := true
	n := 0
	for {
		ch := s.ch
		if ch == '\n' || ch < 0 {
			// only report error if we don't have one already
			if valid {
				s.error(offs, "rune literal not terminated")
				valid = false
			}
			break
		}
		s.next()
		if ch == '\'' {
			break
		}
		n++
		if ch == '\\' {
			if !s.scanEscape('\'') {
				valid = false
			}
			// continue to read to closing quote
		}
	}

	if valid && n != 1 {
		s.error(offs, "illegal rune literal")
	}

	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanSharpComment() string {
	// '#' opening already consumed
	offs := s.offset - 1

	for {
		ch := s.ch
		if ch == '\n' || ch < 0 {
			break
		}
		s.next()
	}

	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanString() string {
	// '"' opening already consumed
	offs := s.offset - 1

	for {
		ch := s.ch
		if ch == '\n' || ch < 0 {
			s.error(offs, "string literal not terminated")
			break
		}
		s.next()
		if ch == '"' {
			break
		}
		if ch == '\\' {
			s.scanEscape('"')
		}
	}

	return string(s.src[offs:s.offset])
}

func stripCR(b []byte) []byte {
	c := make([]byte, len(b))
	i := 0
	for _, ch := range b {
		if ch != '\r' {
			c[i] = ch
			i++
		}
	}
	return c[:i]
}

func (s *Scanner) scanRawString() string {
	// '`' opening already consumed
	offs := s.offset - 1

	hasCR := false
	for {
		ch := s.ch
		if ch < 0 {
			s.error(offs, "raw string literal not terminated")
			break
		}
		s.next()
		if ch == '`' {
			break
		}
		if ch == '\r' {
			hasCR = true
		}
	}

	lit := s.src[offs:s.offset]
	if hasCR {
		lit = stripCR(lit)
	}

	return string(lit)
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' && !s.insertSemi || s.ch == '\r' {
		s.next()
	}
}

// Helper functions for scanning multi-byte tokens such as >> += >>= .
// Different routines recognize different length tok_i based on matches
// of ch_i. If a token ends in '=', the result is tok1 or tok3
// respectively. Otherwise, the result is tok0 if there was no other
// matching character, or tok2 if the matching character was ch2.

func (s *Scanner) switch2(tok0, tok1 uint) uint {
	if s.ch == '=' {
		s.next()
		return tok1
	}
	return tok0
}

func (s *Scanner) switch3(tok0, tok1 uint, ch2 rune, tok2 uint) uint {
	if s.ch == '=' {
		s.next()
		return tok1
	}
	if s.ch == ch2 {
		s.next()
		return tok2
	}
	return tok0
}

func (s *Scanner) switch4(tok0, tok1 uint, ch2 rune, tok2, tok3 uint) uint {
	if s.ch == '=' {
		s.next()
		return tok1
	}
	if s.ch == ch2 {
		s.next()
		if s.ch == '=' {
			s.next()
			return tok3
		}
		return tok2
	}
	return tok0
}

// Scan scans the next token and returns the token position, the token,
// and its literal string if applicable. The source end is indicated by
// token.EOF.
//
// If the returned token is a literal (IDENT, INT, FLOAT, IMAG, CHAR, STRING)
// or COMMENT, the literal string has the corresponding value.
//
// If the returned token is SEMICOLON, the corresponding
// literal string is ";" if the semicolon was present in the source,
// and "\n" if the semicolon was inserted because of a newline or
// at EOF.
//
// If the returned token is ILLEGAL, the literal string is the offending
// character.
//
// In all other cases, Scan returns an empty literal string.
//
// For more tolerant parsing, Scan will return a valid token if
// possible even if a syntax error was encountered. Thus, even
// if the resulting token sequence contains no illegal tokens,
// a client may not assume that no error occurred. Instead it
// must check the scanner's ErrorCount or the number of calls
// of the error handler, if there was one installed.
//
// Scan adds line information to the file added to the file
// set with Init. Token positions are relative to that file
// and thus relative to the file set.
//
func (s *Scanner) Scan() (t Token) {
scanAgain:
	s.skipWhitespace()

	// current token start
	t.Pos = s.file.Pos(s.offset)

	// determine token value
	insertSemi := false
	switch ch := s.ch; {
	case isLetter(ch):
		t.Literal = s.scanIdentifier()
		insertSemi = true
		t.Kind = IDENT
	case '0' <= ch && ch <= '9':
		insertSemi = true
		t.Kind, t.Literal = s.scanNumber(false)
	default:
		s.next() // always make progress
		switch ch {
		case -1:
			if s.insertSemi {
				s.insertSemi = false // EOF consumed
				t.Kind, t.Literal = SEMICOLON, "\n"
				return
			}
			t.Kind = EOF
		case '\n':
			// we only reach here if s.insertSemi was
			// set in the first place and exited early
			// from s.skipWhitespace()
			s.insertSemi = false // newline consumed
			t.Kind, t.Literal = SEMICOLON, "\n"
			return
		case '"':
			insertSemi = true
			t.Kind = STRING
			t.Literal = s.scanString()
		case '\'':
			insertSemi = true
			t.Kind = CHAR
			t.Literal = s.scanRune()
		case '`':
			insertSemi = true
			t.Kind = STRING
			t.Literal = s.scanRawString()
		case ':':
			t.Kind = s.switch2(COLON, DEFINE)
		case '.':
			if '0' <= s.ch && s.ch <= '9' {
				insertSemi = true
				t.Kind, t.Literal = s.scanNumber(true)
			} else if s.ch == '.' {
				s.next()
				if s.ch == '.' {
					s.next()
					t.Kind = ELLIPSIS
				}
			} else {
				t.Kind = PERIOD
			}
		case ',':
			t.Kind = COMMA
		case ';':
			t.Kind = SEMICOLON
			t.Literal = ";"
		case '(':
			t.Kind = LPAREN
		case ')':
			insertSemi = true
			t.Kind = RPAREN
		case '[':
			t.Kind = LBRACK
		case ']':
			insertSemi = true
			t.Kind = RBRACK
		case '{':
			t.Kind = LBRACE
		case '}':
			insertSemi = true
			t.Kind = RBRACE
		case '+':
			t.Kind = s.switch3(ADD, ADD_ASSIGN, '+', INC)
			if t.Kind == INC {
				insertSemi = true
			}
		case '-':
			t.Kind = s.switch3(SUB, SUB_ASSIGN, '-', DEC)
			if t.Kind == DEC {
				insertSemi = true
			}
		case '*':
			t.Kind = s.switch2(MUL, MUL_ASSIGN)
		case '/':
			if s.ch == '/' || s.ch == '*' {
				// comment
				if s.insertSemi && s.findLineEnd() {
					// reset position to the beginning of the comment
					s.ch = '/'
					s.offset = s.file.Offset(t.Pos)
					s.rdOffset = s.offset + 1
					s.insertSemi = false // newline consumed
					t.Kind, t.Literal = SEMICOLON, "\n"
					return
				}
				comment := s.scanComment()
				if s.mode&ScanComments == 0 {
					// skip comment
					s.insertSemi = false // newline consumed
					goto scanAgain
				}
				t.Kind = COMMENT
				t.Literal = comment
			} else {
				t.Kind = s.switch2(QUO, QUO_ASSIGN)
			}
		case '#':
			if s.insertSemi {
				s.ch = '#'
				s.offset = s.file.Offset(t.Pos)
				s.rdOffset = s.offset + 1
				s.insertSemi = false
				t.Kind, t.Literal = SEMICOLON, "\n"
				return
			}
			comment := s.scanSharpComment()
			if s.mode&ScanComments == 0 { // skip comment
				goto scanAgain
			}
			t.Kind = COMMENT
			t.Literal = comment
		case '%':
			t.Kind = s.switch2(REM, REM_ASSIGN)
		case '^':
			t.Kind = s.switch2(XOR, XOR_ASSIGN)
		case '<':
			if s.ch == '-' {
				s.next()
				t.Kind = ARROW
			} else {
				t.Kind = s.switch4(LT, LE, '<', SHL, SHL_ASSIGN)
			}
		case '>':
			t.Kind = s.switch4(GT, GE, '>', SHR, SHR_ASSIGN)
		case '=':
			t.Kind = s.switch2(ASSIGN, EQ)
		case '!':
			t.Kind = s.switch2(NOT, NE)
		case '&':
			if s.ch == '^' {
				s.next()
				t.Kind = s.switch2(AND_NOT, AND_NOT_ASSIGN)
			} else {
				t.Kind = s.switch3(AND, AND_ASSIGN, '&', LAND)
			}
		case '|':
			t.Kind = s.switch3(OR, OR_ASSIGN, '|', LOR)
		case '?':
			t.Kind = QUESTION
		case '~':
			t.Kind = TILDE
		case '@':
			t.Kind = AT
		default:
			// next reports unexpected BOMs - don't repeat
			if ch != bom {
				s.error(s.file.Offset(t.Pos), fmt.Sprintf("illegal character %#U", ch))
			}
			insertSemi = s.insertSemi // preserve insertSemi info
			t.Kind = ILLEGAL
			t.Literal = string(ch)
		}
	}
	if s.mode&InsertSemis != 0 {
		s.insertSemi = insertSemi
	}

	return
}

// Source returns the scanning source.
//
func (s *Scanner) Source() TokenSource {

	return TokenSource{s.file, s.src}
}
