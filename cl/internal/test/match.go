/*
 * Copyright (c) 2024 The XGo Authors (xgo.dev). All rights reserved.
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

const (
	Gopo_Gopt_Case_Match = "Gopt_Case_MatchTBase,Gopt_Case_MatchAny"
)

func Gopt_Case_MatchTBase[T basetype](t CaseT, got, expected T, name ...string) {
}

func Gopt_Case_MatchAny(t CaseT, got, expected any, name ...string) {
}
