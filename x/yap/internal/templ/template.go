/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
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

package templ

import (
	"strings"
)

func Translate(text string) string {
	offs := make([]int, 0, 16)
	base := 0
	for {
		pos := strings.Index(text[base:], "{{")
		if pos < 0 {
			break
		}
		begin := base + pos + 2 // script begin
		n := strings.Index(text[begin:], "}}")
		if n < 0 {
			n = len(text) - begin // script length
		}
		base = begin + n
		code := text[begin:base]
		nonBlank := false
		for i := 0; i < n; i++ {
			c := code[i]
			if !isSpace(c) {
				nonBlank = true
			} else if c == '\n' && nonBlank {
				off := begin + i
				if i, nonBlank = findScript(code, i+1, n); !nonBlank {
					break // script not found
				}
				offs = append(offs, off) // insert }}{{
			}
		}
	}
	n := len(offs)
	if n == 0 {
		return text
	}
	var b strings.Builder
	b.Grow(len(text) + n*4)
	base = 0
	for i := 0; i < n; i++ {
		off := offs[i]
		b.WriteString(text[base:off])
		b.WriteString("}}{{")
		base = off
	}
	b.WriteString(text[base:])
	return b.String()
}

func isSpace(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\v', '\f', '\r', 0x85, 0xA0:
		return true
	}
	return false
}

func findScript(code string, i, n int) (int, bool) {
	for i < n {
		if !isSpace(code[i]) {
			return i, true
		}
		i++
	}
	return -1, false
}
