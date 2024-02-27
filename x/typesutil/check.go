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
	"path/filepath"
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/token"
	"github.com/goplus/gop/x/c2go"
	"github.com/goplus/gop/x/typesutil/internal/typesutil"
	"github.com/goplus/gox"
	"github.com/goplus/mod/gopmod"
	"github.com/qiniu/x/errors"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

type dbgFlags int

const (
	DbgFlagVerbose dbgFlags = 1 << iota
	DbgFlagPrintError
	DbgFlagDisableRecover
	DbgFlagDefault = DbgFlagVerbose | DbgFlagPrintError
	DbgFlagAll     = DbgFlagDefault | DbgFlagDisableRecover
)

var (
	debugVerbose  bool
	debugPrintErr bool
)

func SetDebug(flags dbgFlags) {
	debugVerbose = (flags & DbgFlagVerbose) != 0
	debugPrintErr = (flags & DbgFlagPrintError) != 0
	if (flags & DbgFlagDisableRecover) != 0 {
		cl.SetDisableRecover(true)
	}
}

// -----------------------------------------------------------------------------

type Project = cl.Project

type Config struct {
	// Types provides type information for the package (required).
	Types *types.Package

	// Fset provides source position information for syntax trees and types (required).
	Fset *token.FileSet

	// WorkingDir is the directory in which to run gop compiler (optional).
	// If WorkingDir is not set, os.Getwd() is used.
	WorkingDir string

	// C2goBase specifies base of standard c2go packages (optional).
	// Default is github.com/goplus/.
	C2goBase string

	// Mod represents a Go+ module (optional).
	Mod *gopmod.Module

	// If IgnoreFuncBodies is set, skip compiling function bodies (optional).
	IgnoreFuncBodies bool
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
func (p *Checker) Files(goFiles []*goast.File, gopFiles []*ast.File) (err error) {
	opts := p.opts
	pkgTypes := opts.Types
	fset := opts.Fset
	conf := p.conf
	if len(gopFiles) == 0 {
		checker := types.NewChecker(conf, fset, pkgTypes, p.goInfo)
		return checker.Files(goFiles)
	}
	files := make([]*goast.File, 0, len(goFiles))
	gofs := make(map[string]*goast.File)
	gopfs := make(map[string]*ast.File)
	for _, goFile := range goFiles {
		f := fset.File(goFile.Pos())
		if f == nil {
			continue
		}
		file := f.Name()
		fname := filepath.Base(file)
		if strings.HasPrefix(fname, "gop_autogen") {
			continue
		}
		gofs[file] = goFile
		files = append(files, goFile)
	}
	for _, gopFile := range gopFiles {
		f := fset.File(gopFile.Pos())
		if f == nil {
			continue
		}
		gopfs[f.Name()] = gopFile
	}
	if debugVerbose {
		log.Println("typesutil.Check:", pkgTypes.Path(), "gopFiles =", len(gopfs), "goFiles =", len(gofs))
	}
	pkg := &ast.Package{
		Name:    pkgTypes.Name(),
		Files:   gopfs,
		GoFiles: gofs,
	}
	mod := opts.Mod
	if mod == nil {
		mod = gopmod.Default
	}
	_, err = cl.NewPackage(pkgTypes.Path(), pkg, &cl.Config{
		Types:          pkgTypes,
		Fset:           fset,
		C2goBase:       opts.C2goBase,
		LookupPub:      c2go.LookupPub(mod),
		LookupClass:    mod.LookupClass,
		Importer:       conf.Importer,
		Recorder:       NewRecorder(p.gopInfo),
		NoFileLine:     true,
		NoAutoGenMain:  true,
		NoSkipConstant: true,
		Outline:        opts.IgnoreFuncBodies,
	})
	if err != nil {
		if onErr := conf.Error; onErr != nil {
			if list, ok := err.(errors.List); ok {
				for _, e := range list {
					if ce, ok := convErr(fset, e); ok {
						onErr(ce)
					}
				}
			} else if ce, ok := convErr(fset, err); ok {
				onErr(ce)
			}
		}
		if debugPrintErr {
			log.Println("typesutil.Check err:", err)
			log.SingleStack()
		}
		// don't return even if err != nil
	}
	if len(files) > 0 && p.goInfo != nil {
		scope := pkgTypes.Scope()
		objMap := DeleteObjects(scope, files)
		checker := types.NewChecker(conf, fset, pkgTypes, p.goInfo)
		err = checker.Files(files)
		// TODO: how to process error?
		CorrectTypesInfo(scope, objMap, p.gopInfo.Uses)
	}
	return
}

type astIdent interface {
	comparable
	ast.Node
}

type objMapT = map[types.Object]types.Object

// CorrectTypesInfo corrects types info to avoid there are two instances for the same Go object.
func CorrectTypesInfo[Ident astIdent](scope *types.Scope, objMap objMapT, uses map[Ident]types.Object) {
	for o := range objMap {
		objMap[o] = scope.Lookup(o.Name())
	}
	for id, old := range uses {
		if new := objMap[old]; new != nil {
			uses[id] = new
		}
	}
}

// DeleteObjects deletes all objects defined in Go files and returns deleted objects.
func DeleteObjects(scope *types.Scope, files []*goast.File) objMapT {
	objMap := make(objMapT)
	for _, f := range files {
		for _, decl := range f.Decls {
			switch v := decl.(type) {
			case *goast.GenDecl:
				for _, spec := range v.Specs {
					switch v := spec.(type) {
					case *goast.ValueSpec:
						for _, name := range v.Names {
							scopeDelete(objMap, scope, name.Name)
						}
					case *goast.TypeSpec:
						scopeDelete(objMap, scope, v.Name.Name)
					}
				}
			case *goast.FuncDecl:
				if v.Recv == nil {
					scopeDelete(objMap, scope, v.Name.Name)
				}
			}
		}
	}
	return objMap
}

func convErr(fset *token.FileSet, e error) (ret types.Error, ok bool) {
	switch v := e.(type) {
	case *gox.CodeError:
		ret.Pos, ret.Msg = v.Pos, v.Msg
		typesutil.SetErrorGo116(&ret, 0, v.Pos, v.Pos)
	case *gox.MatchError:
		end := token.NoPos
		if v.Src != nil {
			ret.Pos, end = v.Src.Pos(), v.Src.End()
		}
		ret.Msg = v.Message("")
		typesutil.SetErrorGo116(&ret, 0, ret.Pos, end)
	case *gox.ImportError:
		ret.Pos, ret.Msg = v.Pos, v.Err.Error()
		typesutil.SetErrorGo116(&ret, 0, v.Pos, v.Pos)
	default:
		return
	}
	ret.Fset, ok = fset, true
	return
}

func scopeDelete(objMap map[types.Object]types.Object, scope *types.Scope, name string) {
	if o := typesutil.ScopeDelete(scope, name); o != nil {
		objMap[o] = nil
	}
}
