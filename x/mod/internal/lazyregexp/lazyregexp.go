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

// Package lazyregexp is a thin wrapper over regexp, allowing the use of global
// regexp variables without forcing them to be compiled at init.
package lazyregexp

import (
	"os"
	"regexp"
	"strings"
	"sync"
)

// Regexp is a wrapper around regexp.Regexp, where the underlying regexp will be
// compiled the first time it is needed.
type Regexp struct {
	str  string
	once sync.Once
	rx   *regexp.Regexp
}

func (r *Regexp) re() *regexp.Regexp {
	r.once.Do(r.build)
	return r.rx
}

func (r *Regexp) build() {
	r.rx = regexp.MustCompile(r.str)
	r.str = ""
}

func (r *Regexp) FindSubmatch(s []byte) [][]byte {
	return r.re().FindSubmatch(s)
}

func (r *Regexp) FindStringSubmatch(s string) []string {
	return r.re().FindStringSubmatch(s)
}

func (r *Regexp) FindStringSubmatchIndex(s string) []int {
	return r.re().FindStringSubmatchIndex(s)
}

func (r *Regexp) ReplaceAllString(src, repl string) string {
	return r.re().ReplaceAllString(src, repl)
}

func (r *Regexp) FindString(s string) string {
	return r.re().FindString(s)
}

func (r *Regexp) FindAllString(s string, n int) []string {
	return r.re().FindAllString(s, n)
}

func (r *Regexp) MatchString(s string) bool {
	return r.re().MatchString(s)
}

func (r *Regexp) SubexpNames() []string {
	return r.re().SubexpNames()
}

var inTest = len(os.Args) > 0 && strings.HasSuffix(strings.TrimSuffix(os.Args[0], ".exe"), ".test")

// New creates a new lazy regexp, delaying the compiling work until it is first
// needed. If the code is being run as part of tests, the regexp compiling will
// happen immediately.
func New(str string) *Regexp {
	lr := &Regexp{str: str}
	if inTest {
		// In tests, always compile the regexps early.
		lr.re()
	}
	return lr
}
