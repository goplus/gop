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

package scanner

import (
	"testing"

	"github.com/goplus/gop/tpl/token"
	"github.com/goplus/gop/tpl/types"
)

type Token = types.Token

type tokenTest struct {
	Pos     token.Pos
	Kind    token.Token
	Literal string
}

func TestScanner(t *testing.T) {
	var expected = []tokenTest{
		{3, token.IDENT, `term`},
		{8, '=', ``},
		{10, token.IDENT, `factor`},
		{17, '*', ``},
		{18, '(', ``},
		{19, token.CHAR, `'*'`},
		{23, token.IDENT, `factor`},
		{29, '/', ``},
		{30, token.IDENT, `mul`},
		{34, '|', ``},
		{36, token.CHAR, `'/'`},
		{40, token.IDENT, `factor`},
		{46, '/', ``},
		{47, token.IDENT, `div`},
		{50, ')', ``},
		{51, ';', "\n"},
		{53, token.IDENT, `expr`},
		{58, '=', ``},
		{60, token.IDENT, `term`},
		{65, '*', ``},
		{66, '(', ``},
		{67, token.CHAR, `'+'`},
		{71, token.IDENT, `term`},
		{75, '/', ``},
		{76, token.IDENT, `add`},
		{80, '|', ``},
		{82, token.CHAR, `'-'`},
		{86, token.IDENT, `term`},
		{90, '/', ``},
		{91, token.IDENT, `sub`},
		{94, ')', ``},
		{95, ';', "\n"},
		{97, token.IDENT, `factor`},
		{104, '=', ``},
		{107, token.IDENT, `FLOAT`},
		{112, '/', ``},
		{113, token.IDENT, `push`},
		{118, '|', ``},
		{121, token.CHAR, `'-'`},
		{125, token.IDENT, `factor`},
		{131, '/', ``},
		{132, token.IDENT, `neg`},
		{136, '|', ``},
		{139, token.CHAR, `'('`},
		{143, token.IDENT, `expr`},
		{148, token.CHAR, `')'`},
		{152, '|', ``},
		{155, '(', ``},
		{156, token.IDENT, `IDENT`},
		{162, token.CHAR, `'('`},
		{166, token.IDENT, `expr`},
		{171, '%', ``},
		{173, token.CHAR, `','`},
		{176, '/', ``},
		{177, token.IDENT, `arity`},
		{183, token.CHAR, `')'`},
		{186, ')', ``},
		{187, '/', ``},
		{188, token.IDENT, `call`},
		{193, '|', ``},
		{196, token.CHAR, `'+'`},
		{200, token.IDENT, `factor`},
		{206, ';', "\n"},
		{208, token.IDENT, `hello`},
		{214, ';', "\n"},
		{214, token.COMMENT, "#!/user/bin/env tpl 中文"},
		{241, token.IDENT, `hello`},
		{247, ';', "\n"},
		{247, token.COMMENT, "//!/user/bin/env tpl 中文"},
		{275, token.IDENT, `world`},
		{281, token.NE, ``},
	}

	const grammar = `

term = factor *('*' factor/mul | '/' factor/div)

expr = term *('+' term/add | '-' term/sub)

factor =
	FLOAT/push |
	'-' factor/neg |
	'(' expr ')' |
	(IDENT '(' expr % ','/arity ')')/call |
	'+' factor

hello #!/user/bin/env tpl 中文
hello //!/user/bin/env tpl 中文
world !=
`
	var s Scanner

	fset := token.NewFileSet()
	file := fset.AddFile("", -1, len(grammar))

	s.Init(file, []byte(grammar), nil /* no error handler */, ScanComments|InsertSemis)
	i := 0
	for {
		c := s.Scan()
		if c.Tok == token.EOF {
			break
		}
		expect := Token{Tok: expected[i].Kind, Pos: expected[i].Pos, Lit: expected[i].Literal}
		if c != expect {
			t.Fatal("Scan failed:", c, expect)
		}
		i++
	}
	if len(expected) != i {
		t.Fatalf("len(expected) != i: %d, %d\n", len(expected), i)
	}

	s.Init(file, []byte(grammar), nil, 0)
	i = 0
	for {
		c := s.Scan()
		if c.Tok == token.EOF {
			break
		}
		if expected[i].Kind == ';' {
			i++
		}
		if expected[i].Kind == token.COMMENT {
			i++
		}
		expect := Token{Tok: expected[i].Kind, Pos: expected[i].Pos, Lit: expected[i].Literal}
		if c != expect {
			t.Fatal("Scan failed:", c.Pos, c.Lit, expect)
		}
		switch c.Tok.Len() {
		case 0, 1, 2, 3:
		default:
			t.Fatal("TokenLen failed:", c.Tok.Len())
		}
		i++
	}
	if i < len(expected) && expected[i].Kind == ';' {
		i++
	}
	if len(expected) != i {
		t.Fatalf("len(expected) != i: %d, %d\n", len(expected), i)
	}

	s.Init(file, []byte(grammar), nil, InsertSemis)
	i = 0
	for {
		c := s.Scan()
		if c.Tok == token.EOF {
			break
		}
		if expected[i].Kind == token.COMMENT {
			i++
		}
		expect := Token{Tok: expected[i].Kind, Pos: expected[i].Pos, Lit: expected[i].Literal}
		if c != expect {
			t.Fatal("Scan failed:", c.Pos, c.Lit, expect)
		}
		i++
	}
	if len(expected) != i {
		t.Fatalf("len(expected) != i: %d, %d\n", len(expected), i)
	}
}
