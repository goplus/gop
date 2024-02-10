/*
 * Copyright (c) 2024 The GoPlus Authors (goplus.org). All rights reserved.
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

package test

const (
	GopPackage = true
)

type basetype interface {
	string | int | bool | float64
}

type Case struct {
	CaseT
}

func Gopt_Case_Equal__0[T basetype](t CaseT, a, b T) bool {
	return a == b
}

func Gopt_Case_Equal__1(t CaseT, a, b any) bool {
	return true
}

func Gopt_Case_Match__0[T basetype](t CaseT, got, expected T, name ...string) {
}

func Gopt_Case_Match__1(t CaseT, got, expected any, name ...string) {
}
