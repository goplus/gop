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

package strx

// -----------------------------------------------------------------------------

// Builder concatenates parts of a string together.
type Builder struct {
	b []byte
}

// NewBuilder creates a new Builder object.
func NewBuilder(ncap int) *Builder {
	return &Builder{b: make([]byte, 0, ncap)}
}

// Add appends a string part.
func (p *Builder) Add(s string) *Builder {
	p.b = append(p.b, s...)
	return p
}

// AddByte appends a bytes string part.
func (p *Builder) AddByte(s ...byte) *Builder {
	p.b = append(p.b, s...)
	return p
}

// Build concatenates parts of a string together and returns it.
func (p *Builder) Build() string {
	return String(p.b)
}

// -----------------------------------------------------------------------------
