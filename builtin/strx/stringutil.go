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

package strx

import (
	"unicode"
	"unicode/utf8"

	"github.com/qiniu/x/stringutil"
)

// Capitalize returns a copy of the string str with the first letter mapped to
// its upper case.
func Capitalize(str string) string {
	c, nc := utf8.DecodeRuneInString(str)
	if c == utf8.RuneError || unicode.IsUpper(c) {
		return str
	}
	ret := make([]byte, len(str))
	nr := utf8.EncodeRune(ret, unicode.ToUpper(c))
	ret = append(ret[:nr], str[nc:]...)
	return stringutil.String(ret)
}
