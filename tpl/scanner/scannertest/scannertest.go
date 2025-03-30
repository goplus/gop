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

package scannertest

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/goplus/gop/tpl/scanner"
	"github.com/goplus/gop/tpl/token"

	goscanner "go/scanner"
	gotoken "go/token"

	gopscanner "github.com/goplus/gop/scanner"
	goptoken "github.com/goplus/gop/token"
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

// Scan scans the input and writes the tokens to the writer.
func Scan(w io.Writer, in []byte) {
	var s scanner.Scanner
	fset := gotoken.NewFileSet()
	f := fset.AddFile("", -1, len(in))
	s.Init(f, in, nil, scanner.ScanComments)
	for {
		t := s.Scan()
		if t.Tok == token.EOF {
			break
		}
		fmt.Fprintln(w, t.Pos, t.Tok, t.Lit)
	}
}

// -----------------------------------------------------------------------------

// GopScan scans the input using the Go+ standard library scanner and
// writes the tokens to the writer.
func GopScan(w io.Writer, in []byte) {
	var s gopscanner.Scanner
	fset := gotoken.NewFileSet()
	f := fset.AddFile("", -1, len(in))
	s.Init(f, in, nil, gopscanner.ScanComments)
	for {
		pos, tok, lit := s.Scan()
		if tok == goptoken.EOF {
			break
		}
		fmt.Fprintln(w, pos, tok, lit)
	}
}

// -----------------------------------------------------------------------------

// GoScan scans the input using the Go standard library scanner and
// writes the tokens to the writer.
func GoScan(w io.Writer, in []byte) {
	var s goscanner.Scanner
	fset := gotoken.NewFileSet()
	f := fset.AddFile("", -1, len(in))
	s.Init(f, in, nil, goscanner.ScanComments)
	for {
		pos, tok, lit := s.Scan()
		if tok == gotoken.EOF {
			break
		}
		fmt.Fprintln(w, pos, tok, lit)
	}
}

// -----------------------------------------------------------------------------
