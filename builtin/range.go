/*
 Copyright 2021 The GoPlus Authors (goplus.org)
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package builtin

func NewRange__0(low, high, step int) *IntRange {
	return &IntRange{Low: low, High: high, Step: step}
}

func NewRange__1(low, high, step float64) *FloatRange {
	return &FloatRange{Low: low, High: high, Step: step}
}

type IntRange struct {
	Low, High, Step int
}

func (p *IntRange) Gop_Enum() *intRangeIter {
	return &intRangeIter{i: p.Low, high: p.High, step: p.Step}
}

type FloatRange struct {
	Low, High, Step float64
}

func (p *FloatRange) Gop_Enum() *floatRangeIter {
	return &floatRangeIter{i: p.Low, high: p.High, step: p.Step}
}

type intRangeIter struct {
	i, high, step int
}

func (p *intRangeIter) Next() (val int, ok bool) {
	if p.i < p.high {
		val, ok = p.i, true
		p.i += p.step
	}
	return
}

type floatRangeIter struct {
	i, high, step float64
}

func (p *floatRangeIter) Next() (val float64, ok bool) {
	if p.i < p.high {
		val, ok = p.i, true
		p.i += p.step
	}
	return
}
