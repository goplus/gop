package tpl

import (
	"testing"
)

func TestTokens(tt *testing.T) {

	for i, t := range tokens {
		if t == "" {
			continue
		}
		if i > 20 {
			n := len(t)
			if n < 3 {
				tt.Fatal("invalid token:", i, t)
			}
			if i < 0x80 {
				if t[0] != '\'' || t[n-1] != '\'' {
					tt.Fatal("invalid token:", i, t)
				}
				if n != 3 || t[1] != byte(i) {
					tt.Fatal("invalid token:", i, t)
				}
			} else {
				if t[0] != '"' || t[n-1] != '"' {
					tt.Fatal("invalid token:", i, t)
				}
			}
		}
	}
}

func TestToken(t *testing.T) {

	if operator_beg != 0x80 {
		t.Fatal("operator_beg != 0x80", int(operator_beg))
	}

	s := new(Scanner)
	if s.Ttol(LPAREN) != "'('" {
		t.Fatal("LPAREN")
	}
	if s.Ttol(IDENT) != "IDENT" {
		t.Fatal("IDENT")
	}
	if s.Ttol(ADD_ASSIGN) != "\"+=\"" {
		t.Fatal("ADD_ASSIGN")
	}

	if s.Ltot("\"++\"") != INC {
		t.Fatal("++")
	}
	if s.Ltot("')'") != RPAREN {
		t.Fatal(")")
	}
	if s.Ltot("COMMENT") != COMMENT {
		t.Fatal("COMMENT")
	}
}
