package eqlang

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"unicode"
)

var (
	// ErrEndRequired is returned if eql script doesn't end by `%>`.
	ErrEndRequired = errors.New("eql script requires `%>` to end")

	// ErrEndOfString is returned if string doesn't end.
	ErrEndOfString = errors.New("string doesn't end")
)

// -----------------------------------------------------------------------------

// Parse parses eql source into qlang code.
//
func Parse(source string) (code []byte, err error) {

	var b bytes.Buffer
	for {
		pos := strings.Index(source, "<%")
		if pos < 0 {
			err = parseText(&b, source)
			break
		}
		if pos > 0 {
			err = parseText(&b, source[:pos])
			if err != nil {
				return
			}
		}
		source, err = parseEql(&b, source[pos+2:])
		if err != nil {
			return
		}
	}
	code = b.Bytes()
	return
}

func parseText(b *bytes.Buffer, source string) (err error) {

	b.WriteString("print(eql.Subst(`")
	for {
		pos := strings.IndexByte(source, '`')
		if pos < 0 {
			b.WriteString(source)
			break
		}
		b.WriteString(source[:pos])
		pos2 := pos + 1
		for ; pos2 < len(source); pos2++ {
			if source[pos2] != '`' {
				break
			}
		}
		b.WriteString("` + \"")
		b.WriteString(source[pos:pos2])
		b.WriteString("\" + `")
		source = source[pos2:]
	}
	b.WriteString("`)); ")
	return
}

func parseEql(b *bytes.Buffer, source string) (ret string, err error) {

	fexpr := strings.HasPrefix(source, "=")
	if fexpr {
		b.WriteString("print(")
		source = source[1:]
	}
	for {
		pos := strings.IndexAny(source, "%\"`")
		if pos < 0 {
			return "", ErrEndRequired
		}
		if c := source[pos]; c == '%' {
			if strings.HasPrefix(source[pos+1:], ">") {
				ret = source[pos+2:]
				b.WriteString(source[:pos])
				if fexpr {
					b.WriteString("); ")
				} else if strings.HasPrefix(ret, "\n") {
					ret = ret[1:]
					b.WriteString("\n")
				} else {
					b.WriteString("; ")
				}
				return
			}
			b.WriteString(source[:pos+1])
			source = source[pos+1:]
		} else {
			n := findEnd(source[pos+1:], c)
			if n < 0 {
				return "", ErrEndOfString
			}
			n += pos + 2
			b.WriteString(source[:n])
			source = source[n:]
		}
	}
}

func findEnd(line string, c byte) int {

	for i := 0; i < len(line); i++ {
		switch line[i] {
		case c:
			return i
		case '\\':
			i++
		}
	}
	return -1
}

// -----------------------------------------------------------------------------

// Variables represent how to get value of a variable.
//
type Variables interface {
	GetVar(name string) (v interface{}, ok bool)
}

type mapVars map[string]interface{}
type mapStrings map[string]string

func (vars mapVars) GetVar(name string) (v interface{}, ok bool) {
	v, ok = vars[name]
	return
}

func (vars mapStrings) GetVar(name string) (v interface{}, ok bool) {
	v, ok = vars[name]
	return
}

// Subst substs variables in text.
//
func Subst(text string, lang interface{}) string {

	var vars Variables
	switch v := lang.(type) {
	case map[string]interface{}:
		vars = mapVars(v)
	case map[string]string:
		vars = mapStrings(v)
	case Variables:
		vars = v
	default:
		panic(fmt.Sprintf("eql.Subst: unsupported lang type `%v`", reflect.TypeOf(lang)))
	}
	return subst(text, vars)
}

func subst(text string, vars Variables) string {

	var b bytes.Buffer
	for {
		pos := strings.IndexByte(text, '$')
		if pos < 0 || pos+1 >= len(text) {
			if b.Len() == 0 {
				return text
			}
			b.WriteString(text)
			break
		}
		switch c := text[pos+1]; {
		case c == '$':
			b.WriteString(text[:pos+1])
			text = text[pos+2:]
		case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z'):
			b.WriteString(text[:pos])
			pos1 := pos + 2
			n := strings.IndexFunc(text[pos1:], func(c rune) bool {
				return !unicode.IsLetter(c) && !unicode.IsDigit(c)
			})
			if n < 0 {
				n = len(text) - pos1
			}
			pos2 := pos1 + n
			key := text[pos+1 : pos2]
			val, ok := vars.GetVar(key)
			if !ok {
				panic("variable not found: " + key)
			}
			b.WriteString(fmt.Sprint(val))
			text = text[pos2:]
		default:
			b.WriteString(text[:pos+1])
			text = text[pos+1:]
		}
	}
	return b.String()
}

// -----------------------------------------------------------------------------

// Input decodes a json string.
//
func Input(input string) (ret map[string]interface{}, err error) {

	err = json.NewDecoder(strings.NewReader(input)).Decode(&ret)
	return
}

// InputFile decodes a json file.
//
func InputFile(input string) (ret map[string]interface{}, err error) {

	in := os.Stdin
	if input != "-" {
		f, err1 := os.Open(input)
		if err1 != nil {
			return nil, err1
		}
		defer f.Close()
		in = f
	}
	err = json.NewDecoder(in).Decode(&ret)
	return
}

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"new":   New,
	"parse": Parse,
	"subst": Subst,

	"input":     Input,
	"inputFile": InputFile,

	"New":   New,
	"Parse": Parse,
	"Subst": Subst,

	"Input":     Input,
	"InputFile": InputFile,

	"ErrEndRequired": ErrEndRequired,
	"ErrEndOfString": ErrEndOfString,
}

// -----------------------------------------------------------------------------
