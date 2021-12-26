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
	"testing"

	"golang.org/x/mod/modfile"
)

// -----------------------------------------------------------------------------

const gopmodSpx = `
module spx

go 1.17
gop 1.1

classfile .gmx .spx github.com/goplus/spx math

require (
	github.com/ajstarks/svgo v0.0.0-20210927141636-6d70534b1098
)
`

func TestGoModCompat(t *testing.T) {
	const (
		gopmod = gopmodSpx
	)
	f, err := modfile.ParseLax("go.mod", []byte(gopmod), nil)
	if err != nil || len(f.Syntax.Stmt) != 5 {
		t.Fatal("modfile.ParseLax failed:", f, err)
	}

	gop := f.Syntax.Stmt[2].(*modfile.Line)
	if len(gop.Token) != 2 || gop.Token[0] != "gop" || gop.Token[1] != "1.1" {
		t.Fatal("modfile.ParseLax gop:", gop)
	}

	require := f.Syntax.Stmt[4].(*modfile.LineBlock)
	if len(require.Token) != 1 || require.Token[0] != "require" {
		t.Fatal("modfile.ParseLax require:", require)
	}
	if len(require.Line) != 1 {
		t.Fatal("modfile.ParseLax require.Line:", require.Line)
	}
}

// -----------------------------------------------------------------------------

func TestParse1(t *testing.T) {
	const (
		gopmod = gopmodSpx
	)
	f, err := ParseLax("github.com/goplus/gop/gop.mod", []byte(gopmod), nil)
	if err != nil {
		t.Error(err)
		return
	}
	if f.Gop.Version != "1.1" {
		t.Errorf("gop version expected be 1.1, but %s got", f.Gop.Version)
	}

	if f.Classfile.ProjExt != ".gmx" {
		t.Errorf("classfile exts expected be .gmx, but %s got", f.Classfile.ProjExt)
	}
	if f.Classfile.WorkExt != ".spx" {
		t.Errorf("classfile exts expected be .spx, but %s got", f.Classfile.WorkExt)
	}

	if len(f.Classfile.PkgPaths) != 2 {
		t.Errorf("classfile pkgpaths length expected be 2, but %d got", len(f.Classfile.PkgPaths))
	}

	if f.Classfile.PkgPaths[0] != "github.com/goplus/spx" {
		t.Errorf("classfile path expected be github.com/goplus/spx, but %s got", f.Classfile.PkgPaths[0])
	}
	if f.Classfile.PkgPaths[1] != "math" {
		t.Errorf("classfile path expected be math, but %s got", f.Classfile.PkgPaths[1])
	}
}

const gopmodUserProj = `
module moduleUserProj

go 1.17
gop 1.1

register github.com/goplus/spx

require (
    github.com/goplus/spx v1.0
)
`

func TestParse2(t *testing.T) {
	const (
		gopmod = gopmodUserProj
	)
	f, err := Parse("github.com/goplus/gop/gop.mod", []byte(gopmod), nil)
	if err != nil || len(f.Register) != 1 {
		t.Fatal("Parse:", f, err)
		return
	}
	if f.Register[0].ClassfileMod != "github.com/goplus/spx" {
		t.Fatal("Parse => Register:", f.Register)
	}
}

func TestParseErr(t *testing.T) {
	doTestParseErr(t, `gop.mod:3: repeated go statement`, `
gop 1.1
gop 1.2
`)
	doTestParseErr(t, `gop.mod:2: go directive expects exactly one argument`, `
gop 1.1 1.2
`)
	doTestParseErr(t, `gop.mod:2: invalid gop version '1.x': must match format 1.23`, `
gop 1.x
`)
	doTestParseErr(t, `gop.mod:2: register directive expects exactly one argument`, `
register 1 2 3
`)
	doTestParseErr(t, `gop.mod:2: invalid quoted string: invalid syntax`, `
register "\?"
`)
	doTestParseErr(t, `gop.mod:2: malformed module path "-": leading dash`, `
register -
`)
	doTestParseErr(t, `gop.mod:3: repeated classfile statement`, `
classfile .gmx .spx github.com/goplus/spx math
classfile .gmx .spx github.com/goplus/spx math
`)
	doTestParseErr(t, `gop.mod:2: usage: classfile projExt workExt [classFilePkgPath ...]`, `
classfile .gmx .spx
`)
	doTestParseErr(t, `gop.mod:2: ext . invalid: invalid ext format`, `
classfile .gmx . math
`)
	doTestParseErr(t, `gop.mod:2: ext "\?" invalid: invalid syntax`, `
classfile "\?" .spx math
`)
	doTestParseErr(t, `gop.mod:2: invalid quoted string: invalid syntax`, `
classfile .123 .spx "\?"
`)
	doTestParseErr(t, `gop.mod:2: unknown directive: unknown`, `
unknown .spx
`)
	doTestParseErr(t, `gop.mod:2: invalid go version '1.x': must match format 1.23`, `
go 1.x
`)
}

func doTestParseErr(t *testing.T, errMsg string, gopmod string) {
	t.Run(errMsg, func(t *testing.T) {
		_, err := Parse("gop.mod", []byte(gopmod), nil)
		if err == nil {
			t.Fatal("Parse: no error?")
			return
		}
		if err.Error() != errMsg {
			t.Error("Parse got:", err, "\nExpected:", errMsg)
		}
	})
}

// -----------------------------------------------------------------------------
