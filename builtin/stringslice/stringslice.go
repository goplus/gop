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

package stringslice

import (
	"strings"

	"github.com/qiniu/x/stringutil"
)

// Capitalize capitalizes the first letter of each string in the slice.
func Capitalize(a []string) []string {
	r := make([]string, len(a))
	for i, v := range a {
		r[i] = stringutil.Capitalize(v)
	}
	return r
}

// ToTitle title-cases all strings in the slice.
func ToTitle(a []string) []string {
	r := make([]string, len(a))
	for i, v := range a {
		r[i] = strings.ToTitle(v)
	}
	return r
}

// ToUpper upper-cases all strings in the slice.
func ToUpper(a []string) []string {
	r := make([]string, len(a))
	for i, v := range a {
		r[i] = strings.ToUpper(v)
	}
	return r
}

// ToLower lower-cases all strings in the slice.
func ToLower(a []string) []string {
	r := make([]string, len(a))
	for i, v := range a {
		r[i] = strings.ToLower(v)
	}
	return r
}

// Repeat repeats each string in the slice count times.
func Repeat(a []string, count int) []string {
	r := make([]string, len(a))
	for i, v := range a {
		r[i] = strings.Repeat(v, count)
	}
	return r
}

// Replace replaces all occurrences of old in each string in the slice with new.
func Replace(a []string, old, new string, n int) []string {
	r := make([]string, len(a))
	for i, v := range a {
		r[i] = strings.Replace(v, old, new, n)
	}
	return r
}

// ReplaceAll replaces all occurrences of old in each string in the slice with new.
func ReplaceAll(a []string, old, new string) []string {
	r := make([]string, len(a))
	for i, v := range a {
		r[i] = strings.ReplaceAll(v, old, new)
	}
	return r
}

// Trim removes leading and trailing white space from each string in the slice.
func Trim(a []string, cutset string) []string {
	r := make([]string, len(a))
	for i, v := range a {
		r[i] = strings.Trim(v, cutset)
	}
	return r
}

// TrimSpace removes leading and trailing white space from each string in the slice.
func TrimSpace(a []string) []string {
	r := make([]string, len(a))
	for i, v := range a {
		r[i] = strings.TrimSpace(v)
	}
	return r
}

// TrimLeft removes leading white space from each string in the slice.
func TrimLeft(a []string, cutset string) []string {
	r := make([]string, len(a))
	for i, v := range a {
		r[i] = strings.TrimLeft(v, cutset)
	}
	return r
}

// TrimRight removes trailing white space from each string in the slice.
func TrimRight(a []string, cutset string) []string {
	r := make([]string, len(a))
	for i, v := range a {
		r[i] = strings.TrimRight(v, cutset)
	}
	return r
}

// TrimPrefix removes leading prefix from each string in the slice.
func TrimPrefix(a []string, prefix string) []string {
	r := make([]string, len(a))
	for i, v := range a {
		r[i] = strings.TrimPrefix(v, prefix)
	}
	return r
}

// TrimSuffix removes trailing suffix from each string in the slice.
func TrimSuffix(a []string, suffix string) []string {
	r := make([]string, len(a))
	for i, v := range a {
		r[i] = strings.TrimSuffix(v, suffix)
	}
	return r
}
