package tpl

import (
	"go/token"
	"testing"
)

type tokenTest struct {
	Pos     token.Pos
	Kind    uint
	Literal string
}

func TestScanner(t *testing.T) {

	var expected = []tokenTest{
		{3, IDENT, `term`},
		{8, '=', ``},
		{10, IDENT, `factor`},
		{17, '*', ``},
		{18, '(', ``},
		{19, CHAR, `'*'`},
		{23, IDENT, `factor`},
		{29, '/', ``},
		{30, IDENT, `mul`},
		{34, '|', ``},
		{36, CHAR, `'/'`},
		{40, IDENT, `factor`},
		{46, '/', ``},
		{47, IDENT, `div`},
		{50, ')', ``},
		{51, ';', "\n"},
		{53, IDENT, `expr`},
		{58, '=', ``},
		{60, IDENT, `term`},
		{65, '*', ``},
		{66, '(', ``},
		{67, CHAR, `'+'`},
		{71, IDENT, `term`},
		{75, '/', ``},
		{76, IDENT, `add`},
		{80, '|', ``},
		{82, CHAR, `'-'`},
		{86, IDENT, `term`},
		{90, '/', ``},
		{91, IDENT, `sub`},
		{94, ')', ``},
		{95, ';', "\n"},
		{97, IDENT, `factor`},
		{104, '=', ``},
		{107, IDENT, `FLOAT`},
		{112, '/', ``},
		{113, IDENT, `push`},
		{118, '|', ``},
		{121, CHAR, `'-'`},
		{125, IDENT, `factor`},
		{131, '/', ``},
		{132, IDENT, `neg`},
		{136, '|', ``},
		{139, CHAR, `'('`},
		{143, IDENT, `expr`},
		{148, CHAR, `')'`},
		{152, '|', ``},
		{155, '(', ``},
		{156, IDENT, `IDENT`},
		{162, CHAR, `'('`},
		{166, IDENT, `expr`},
		{171, '%', ``},
		{173, CHAR, `','`},
		{176, '/', ``},
		{177, IDENT, `arity`},
		{183, CHAR, `')'`},
		{186, ')', ``},
		{187, '/', ``},
		{188, IDENT, `call`},
		{193, '|', ``},
		{196, CHAR, `'+'`},
		{200, IDENT, `factor`},
		{206, ';', "\n"},
		{208, IDENT, `hello`},
		{214, ';', "\n"},
		{214, COMMENT, "#!/user/bin/env tpl 中文"},
		{241, IDENT, `hello`},
		{247, ';', "\n"},
		{247, COMMENT, "//!/user/bin/env tpl 中文"},
		{275, IDENT, `world`},
		{280, NOT, ``},
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
world!
`

	var s Scanner

	fset := token.NewFileSet()
	file := fset.AddFile("", -1, len(grammar))

	s.Init(file, []byte(grammar), nil /* no error handler */, ScanComments|InsertSemis)
	i := 0
	for {
		token := s.Scan()
		if token.Kind == EOF {
			break
		}
		expect := Token{expected[i].Kind, expected[i].Pos, expected[i].Literal}
		if token != expect {
			t.Fatal("Scan failed:", token, expect)
		}
		i++
	}
	if len(expected) != i {
		t.Fatalf("len(expected) != i: %d, %d\n", len(expected), i)
	}

	s.Init(file, []byte(grammar), nil, 0)
	i = 0
	for {
		token := s.Scan()
		if token.Kind == EOF {
			break
		}
		if expected[i].Kind == ';' {
			i++
		}
		if expected[i].Kind == COMMENT {
			i++
		}
		expect := Token{expected[i].Kind, expected[i].Pos, expected[i].Literal}
		if token != expect {
			t.Fatal("Scan failed:", token.Pos, s.Ttol(token.Kind), token.Literal, expect)
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
		token := s.Scan()
		if token.Kind == EOF {
			break
		}
		if expected[i].Kind == COMMENT {
			i++
		}
		expect := Token{expected[i].Kind, expected[i].Pos, expected[i].Literal}
		if token != expect {
			t.Fatal("Scan failed:", token.Pos, s.Ttol(token.Kind), token.Literal, expect)
		}
		i++
	}
	if len(expected) != i {
		t.Fatalf("len(expected) != i: %d, %d\n", len(expected), i)
	}
}
