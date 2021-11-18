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

package builtin

// -----------------------------------------------------------------------------

type IntRange struct {
	Start, End, Step int
}

func NewRange__0(start, end, step int) *IntRange {
	return &IntRange{Start: start, End: end, Step: step}
}

func (p *IntRange) Gop_Enum() *intRangeIter {
	step := p.Step
	n := p.End - p.Start + step
	if step > 0 {
		n = (n - 1) / step
	} else {
		n = (n + 1) / step
	}
	return &intRangeIter{n: n, val: p.Start, step: p.Step}
}

// -----------------------------------------------------------------------------

type intRangeIter struct {
	n, val, step int
}

func (p *intRangeIter) Next() (val int, ok bool) {
	if p.n > 0 {
		val, ok = p.val, true
		p.val += p.step
		p.n--
	}
	return
}

// -----------------------------------------------------------------------------
