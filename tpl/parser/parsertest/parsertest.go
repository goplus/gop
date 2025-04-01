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

package parsertest

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/goplus/gop/tpl/ast"
	"github.com/goplus/gop/tpl/token"
)

// -----------------------------------------------------------------------------

func Diff(t *testing.T, outfile string, dst, src []byte) bool {
	line := 1
	offs := 0 // line offset
	for i := 0; i < len(dst) && i < len(src); i++ {
		d := dst[i]
		s := src[i]
		if d != s {
			os.WriteFile(outfile, dst, 0644)
			t.Errorf("dst:%d: %s\n", line, dst[offs:])
			t.Errorf("src:%d: %s\n", line, src[offs:])
			return true
		}
		if s == '\n' {
			line++
			offs = i + 1
		}
	}
	if len(dst) != len(src) {
		os.WriteFile(outfile, dst, 0644)
		t.Errorf("len(dst) = %d, len(src) = %d\ndst = %q\nsrc = %q", len(dst), len(src), dst, src)
		return true
	}
	return false
}

// -----------------------------------------------------------------------------

var (
	tyNode   = reflect.TypeOf((*ast.Node)(nil)).Elem()
	tyString = reflect.TypeOf("")
	tyToken  = reflect.TypeOf(token.Token(0))
)

// FprintNode prints a tpl AST node.
func FprintNode(w io.Writer, lead string, v any, prefix, indent string) {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Slice:
		n := val.Len()
		if n > 0 && lead != "" {
			io.WriteString(w, lead)
		}
		for i := 0; i < n; i++ {
			FprintNode(w, "", val.Index(i).Interface(), prefix, indent)
		}
	case reflect.Ptr:
		t := val.Type()
		if val.IsNil() {
			return
		}
		if t.Implements(tyNode) {
			if lead != "" {
				io.WriteString(w, lead)
			}
			elem, tyElem := val.Elem(), t.Elem()
			fmt.Fprintf(w, "%s%v:\n", prefix, tyElem)
			n := elem.NumField()
			prefix += indent
			for i := 0; i < n; i++ {
				sf := tyElem.Field(i)
				if sf.Name == "RetProc" { // skip RetProc field, see gop/tpl/ast.Rule
					continue
				}
				sfv := elem.Field(i).Interface()
				switch sf.Type {
				case tyString, tyToken:
					fmt.Fprintf(w, "%s%v: %v\n", prefix, sf.Name, sfv)
				default:
					FprintNode(w, fmt.Sprintf("%s%v:\n", prefix, sf.Name), sfv, prefix+indent, indent)
				}
			}
		} else {
			log.Panicln("FprintNode unexpected type:", t)
		}
	case reflect.Int, reflect.Bool, reflect.Invalid:
		// skip
	default:
		log.Panicln("FprintNode unexpected kind:", val.Kind(), "type:", val.Type())
	}
}

// Fprint prints a tpl ast.File node.
func Fprint(w io.Writer, f *ast.File) {
	FprintNode(w, "", f.Decls, "", "  ")
}

// Expect asserts a tpl AST equals output or not.
func Expect(t *testing.T, outfile string, f *ast.File, expected []byte) {
	b := bytes.NewBuffer(nil)
	Fprint(b, f)
	if Diff(t, outfile, b.Bytes(), []byte(expected)) {
		t.Fatal("tpl.Parser: unexpect result")
	}
}

// -----------------------------------------------------------------------------
