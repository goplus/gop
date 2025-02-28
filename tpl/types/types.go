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

package types

import (
	"github.com/goplus/gop/tpl/token"
)

// A Token is a lexical unit returned by Scan.
type Token struct {
	Tok token.Token
	Pos token.Pos
	Lit string
}

// End returns end position of this token.
func (p *Token) End() token.Pos {
	n := len(p.Lit)
	if n == 0 {
		n = p.Tok.Len()
	}
	return p.Pos + token.Pos(n)
}
