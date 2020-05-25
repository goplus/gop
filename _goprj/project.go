/*
 Copyright 2020 Qiniu Cloud (七牛云)
 
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

package goprj

import "github.com/qiniu/x/log"

// -----------------------------------------------------------------------------

// Project represents a new Go project.
type Project struct {
	types      map[string]Type
	openedPkgs map[string]*Package // dir => Package
}

// NewProject creates a new Project.
func NewProject() *Project {
	return &Project{
		types:      make(map[string]Type),
		openedPkgs: make(map[string]*Package),
	}
}

// OpenPackage open a package by specified directory.
func (p *Project) OpenPackage(dir string) (pkg *Package, err error) {
	if pkg, ok := p.openedPkgs[dir]; ok {
		return pkg, nil
	}
	pkg, err = openPackage(dir, p)
	if err != nil {
		return
	}
	p.openedPkgs[dir] = pkg
	return
}

// FindVersionType lookups verPkgPath.name type.
func (p *Project) FindVersionType(verPkgPath string, name string) Type {
	log.Fatalln("FindVersionType: not impl")
	return nil
}

// UniqueType returns the unique instance of a type.
func (p *Project) UniqueType(t Type) Type {
	id := t.ID()
	if v, ok := p.types[id]; ok {
		return v
	}
	p.types[id] = t
	return t
}

// -----------------------------------------------------------------------------
