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

package tpl

import (
	"testing"
)

func TestTok(t *testing.T) {
	if IDENT.String() != "IDENT" {
		t.Fatal("IDENT")
	}
	if Tok(' ').String() != "token(32)" {
		t.Fatal("token(32)")
	}
}

/*
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
*/
