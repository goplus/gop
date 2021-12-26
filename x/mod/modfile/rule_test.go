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
package modfile

import (
	"bytes"
	"syscall"
	"testing"
)

// -----------------------------------------------------------------------------

var addParseExtTests = []struct {
	desc    string
	ext     string
	want    string
	wantErr string
}{
	{
		"spx ok",
		".spx",
		".spx",
		"",
	},
	{
		"no ext",
		"",
		"",
		"",
	},
	{
		"not a ext",
		"gmx",
		"",
		"ext gmx invalid: invalid ext format",
	},
}

func TestParseExt(t *testing.T) {
	if (&InvalidExtError{Err: syscall.EINVAL}).Unwrap() != syscall.EINVAL {
		t.Fatal("InvalidExtError.Unwrap failed")
	}
	for _, tt := range addParseExtTests {
		t.Run(tt.desc, func(t *testing.T) {
			ext, err := parseExt(&tt.ext)
			if err != nil {
				if err.Error() != tt.wantErr {
					t.Fatalf("wanterr: %s, but got: %s", tt.wantErr, err)
				}
			}
			if ext != tt.want {
				t.Fatalf("want: %s, but got: %s", tt.want, ext)
			}
		})
	}
}

func TestIsDirectoryPath(t *testing.T) {
	if !IsDirectoryPath("./...") {
		t.Fatal("IsDirectoryPath failed")
	}
}

func TestFormat(t *testing.T) {
	if b := Format(&FileSyntax{}); len(b) != 0 {
		t.Fatal("Format failed:", b)
	}
}

func TestMustQuote(t *testing.T) {
	if !MustQuote("") {
		t.Fatal("MustQuote failed")
	}
}

// -----------------------------------------------------------------------------

var addGopTests = []struct {
	desc    string
	in      string
	version string
	out     string
}{
	{
		`module_only`,
		`module m
		`,
		`1.14`,
		`module m
		gop 1.14
		`,
	},
	{
		`module_before_require`,
		`module m
		require x.y/a v1.2.3
		`,
		`1.14`,
		`module m
		gop 1.14
		require x.y/a v1.2.3
		`,
	},
	{
		`require_before_module`,
		`require x.y/a v1.2.3
		module example.com/inverted
		`,
		`1.14`,
		`gop 1.14
		require x.y/a v1.2.3
		module example.com/inverted
		`,
	},
	{
		`require_only`,
		`require x.y/a v1.2.3
		`,
		`1.14`,
		`gop 1.14
		require x.y/a v1.2.3
		`,
	},
	{
		`replace gop`,
		`require x.y/a v1.2.3
		gop 1.10
		`,
		`1.14`,
		`require x.y/a v1.2.3
		gop 1.14
		`,
	},
}

func TestAddGop(t *testing.T) {
	for _, tt := range addGopTests {
		t.Run(tt.desc, func(t *testing.T) {
			testEdit(t, tt.in, tt.out, true, func(f *File) error {
				return f.AddGopStmt(tt.version)
			})
		})
	}
}

func TestAddGopErr(t *testing.T) {
	if new(File).AddGopStmt("1.x") == nil {
		t.Fatal("AddGoStmt failed")
	}
}

func testEdit(t *testing.T, in, want string, strict bool, transform func(f *File) error) *File {
	t.Helper()
	parse := Parse
	if !strict {
		parse = ParseLax
	}
	f, err := parse("in", []byte(in), nil)
	if err != nil {
		t.Fatal(err)
	}
	g, err := parse("out", []byte(want), nil)
	if err != nil {
		t.Fatal(err)
	}
	golden, err := g.Format()
	if err != nil {
		t.Fatal(err)
	}

	if err := transform(f); err != nil {
		t.Fatal(err)
	}
	out, err := f.Format()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, golden) {
		t.Errorf("have:\n%s\nwant:\n%s", out, golden)
	}

	return f
}

// -----------------------------------------------------------------------------
