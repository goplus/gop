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
	goast "go/ast"
	"go/types"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/token"
)

type Project = cl.Project

type Config struct {
	// Types provides type information for the package.
	Types *types.Package

	// Fset provides source position information for syntax trees and types.
	// If Fset is nil, Load will use a new fileset, but preserve Fset's value.
	Fset *token.FileSet

	// WorkingDir is the directory in which to run gop compiler.
	WorkingDir string

	// C2goBase specifies base of standard c2go packages.
	// Default is github.com/goplus/.
	C2goBase string

	// LookupPub lookups the c2go package pubfile (named c2go.a.pub).
	LookupPub func(pkgPath string) (pubfile string, err error)

	// LookupClass lookups a class by specified file extension.
	LookupClass func(ext string) (c *Project, ok bool)
}

// A Checker maintains the state of the type checker.
// It must be created with NewChecker.
type Checker struct {
	conf    *types.Config
	opts    *Config
	goInfo  *types.Info
	gopInfo *Info
}

// NewChecker returns a new Checker instance for a given package.
// Package files may be added incrementally via checker.Files.
func NewChecker(conf *types.Config, opts *Config, goInfo *types.Info, gopInfo *Info) *Checker {
	return &Checker{conf, opts, goInfo, gopInfo}
}

// Files checks the provided files as part of the checker's package.
func (p *Checker) Files(goFiles []*goast.File, gopFiles []*ast.File) error {
	opts := p.opts
	pkgTypes := opts.Types
	fset := opts.Fset
	conf := p.conf
	if len(gopFiles) == 0 {
		checker := types.NewChecker(conf, fset, pkgTypes, p.goInfo)
		return checker.Files(goFiles)
	}
	gofs := make(map[string]*goast.File)
	gopfs := make(map[string]*ast.File)
	for _, goFile := range goFiles {
		f := fset.File(goFile.Pos())
		gofs[f.Name()] = goFile
	}
	for _, gopFile := range gopFiles {
		f := fset.File(gopFile.Pos())
		gopfs[f.Name()] = gopFile
	}
	pkg := &ast.Package{
		Name:    pkgTypes.Name(),
		Files:   gopfs,
		GoFiles: gofs,
	}
	_, err := cl.NewPackage(pkgTypes.Path(), pkg, &cl.Config{
		Types:          pkgTypes,
		Fset:           fset,
		WorkingDir:     opts.WorkingDir,
		C2goBase:       opts.C2goBase,
		LookupPub:      opts.LookupPub,
		LookupClass:    opts.LookupClass,
		Importer:       conf.Importer,
		Recorder:       gopRecorder{p.gopInfo},
		NoFileLine:     true,
		NoAutoGenMain:  true,
		NoSkipConstant: true,
	})
	return err
}
