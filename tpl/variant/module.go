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

package variant

import (
	"github.com/goplus/gop/tpl"
)

// -----------------------------------------------------------------------------

// Module represents a module.
type Module struct {
	Name string
	objs map[string]any
}

// NewModule creates a new module.
func NewModule(name string) *Module {
	mod := &Module{
		Name: name,
		objs: make(map[string]any),
	}
	universe.Insert(name, mod)
	return mod
}

// Lookup looks up an object in the module.
func (p *Module) Lookup(name string) (v any, ok bool) {
	v, ok = p.objs[name]
	return
}

// Insert inserts an object into the module.
func (p *Module) Insert(name string, v any) {
	objs := p.objs
	if _, ok := objs[name]; ok {
		panic("object exists: " + name)
	}
	objs[name] = v
}

// Merge merges a module into the current module.
func (p *Module) Merge(mod *Module) {
	for name, obj := range mod.objs {
		p.Insert(name, obj)
	}
}

// -----------------------------------------------------------------------------

var (
	universe = &Module{
		objs: map[string]any{},
	}
)

// Call calls a function.
func Call(needList bool, name string, arglist any) any {
	if o, ok := universe.objs[name]; ok {
		if fn, ok := o.(func(...any) any); ok {
			var args []any
			if arglist != nil {
				args = arglist.([]any)
				if needList {
					args = tpl.List(args)
				}
			}
			return fn(args...)
		}
		panic("not a function: " + name)
	}
	panic("function not found: " + name)
}

// CallObject calls a function object.
func CallObject(needList bool, fn any, arglist any) any {
	if fn, ok := Eval(fn).(func(...any) any); ok {
		var args []any
		if arglist != nil {
			args = arglist.([]any)
			if needList {
				args = tpl.List(args)
			}
		}
		return fn(args...)
	}
	panic("call of non function")
}

// -----------------------------------------------------------------------------

// InitUniverse initializes the universe module with the specified modules.
func InitUniverse(names ...string) {
	for _, name := range names {
		if o, ok := universe.objs[name]; ok {
			if mod, ok := o.(*Module); ok {
				universe.Merge(mod)
			} else {
				panic("not a module: " + name)
			}
		} else {
			panic("module not found: " + name)
		}
	}
}

// -----------------------------------------------------------------------------
