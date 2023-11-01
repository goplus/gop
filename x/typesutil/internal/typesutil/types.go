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

package typesutil

import (
	"go/token"
	"go/types"
	"unsafe"
)

// A Scope maintains a set of objects and links to its containing
// (parent) and contained (children) scopes. Objects may be inserted
// and looked up by name. The zero value for Scope is a ready-to-use
// empty scope.
type Scope struct {
	parent   *Scope
	children []*Scope
	number   int                     // parent.children[number-1] is this scope; 0 if there is no parent
	elems    map[string]types.Object // lazily allocated
	pos, end token.Pos               // scope extent; may be invalid
	comment  string                  // for debugging only
	isFunc   bool                    // set if this is a function scope (internal use only)
}

// ScopeDelete deletes an object from specified scope by its name.
func ScopeDelete(s *types.Scope, name string) types.Object {
	elems := (*Scope)(unsafe.Pointer(s)).elems
	if o, ok := elems[name]; ok {
		delete(elems, name)
		return o
	}
	return nil
}
