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

package formatutil

import (
	"bytes"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/goplus/gop/scanner"
	"github.com/goplus/gop/token"
)

func TestSeekAfter(t *testing.T) {
	if seekAfter(nil, 0, 0) != nil {
		t.Fatal("seekAfter failed")
	}
}

func TestTokOf(t *testing.T) {
	words := []aWord{{tok: token.COMMENT}}
	if tok, _ := tokOf(words); token.COMMENT != tok {
		t.Fatal("tokOf failed:", tok)
	}
	if startWith(words, token.VAR) {
		t.Fatal("startWith failed")
	}
}

func doSplitStmts(src []byte) (ret []string) {
	fset := token.NewFileSet()
	base := fset.Base()
	f := fset.AddFile("", base, len(src))

	var s scanner.Scanner
	s.Init(f, src, nil, scanner.ScanComments)
	stmts := splitStmts(&s)

	ret = make([]string, len(stmts))
	for i, stmt := range stmts {
		ret[i] = stmtKind(stmt)
	}
	return
}

func stmtKind(s aStmt) string {
	tok := s.tok
	if tok == token.FUNC {
		if isFuncDecl(s.words[s.at+1:]) {
			return "FUNC"
		}
		return "FNCALL"
	}
	return tok.String()
}

func testFrom(t *testing.T, pkgDir, sel string, doIt func(in []byte) ([]byte, error)) {
	if sel != "" && !strings.Contains(pkgDir, sel) {
		return
	}
	t.Helper()
	log.Println("Parsing", pkgDir)
	in, err := os.ReadFile(pkgDir + "/in.data")
	if err != nil {
		t.Fatal("Parsing", pkgDir, "-", err)
	}
	out, err := os.ReadFile(pkgDir + "/out.expect")
	if err != nil {
		t.Fatal("Parsing", pkgDir, "-", err)
	}
	ret, err := doIt(in)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(ret, out) {
		t.Fatal("Parsing", pkgDir, "- failed:\n"+string(ret))
	}
}

func testFromDir(t *testing.T, sel, relDir string, doIt func(in []byte) ([]byte, error)) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = path.Join(dir, relDir)
	fis, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		name := fi.Name()
		if strings.HasPrefix(name, "_") {
			continue
		}
		t.Run(name, func(t *testing.T) {
			testFrom(t, dir+"/"+name, sel, doIt)
		})
	}
}

func TestSplitStmts(t *testing.T) {
	testFromDir(t, "", "./_testdata/splitstmts", func(in []byte) ([]byte, error) {
		ret := strings.Join(doSplitStmts(in), "\n") + "\n"
		return []byte(ret), nil
	})
}

func TestRearrangeFuncs(t *testing.T) {
	if runtime.GOOS == "windows" { // skip temporarily
		return
	}
	testFromDir(t, "", "./_testdata/rearrange", func(in []byte) ([]byte, error) {
		return RearrangeFuncs(in)
	})
}

func TestFormat(t *testing.T) {
	testFromDir(t, "", "./_testdata/format", func(in []byte) ([]byte, error) {
		return SourceEx(in, false, "")
	})
}
