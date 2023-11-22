/*
 * Copyright (c) 2022 The GoPlus Authors (goplus.org). All rights reserved.
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

package gop

import (
	"fmt"
	"go/types"
	"io/fs"
	"os"
	"strings"
	"syscall"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/goplus/gop/x/c2go"
	"github.com/goplus/gop/x/gopenv"
	"github.com/goplus/gox"
	"github.com/goplus/mod/env"
	"github.com/goplus/mod/gopmod"
	"github.com/qiniu/x/errors"
)

var (
	ErrNotFound      = syscall.ENOENT
	ErrIgnoreNotated = errors.New("notated error ignored")
)

// NotFound returns if cause err is ErrNotFound or not
func NotFound(err error) bool {
	return errors.Err(err) == ErrNotFound
}

// IgnoreNotated returns if cause err is ErrIgnoreNotated or not.
func IgnoreNotated(err error) bool {
	return errors.Err(err) == ErrIgnoreNotated
}

// ErrorPos returns where the error occurs.
func ErrorPos(err error) token.Pos {
	switch v := err.(type) {
	case *gox.CodeError:
		return v.Pos
	case *gox.MatchError:
		return v.Pos()
	case *gox.ImportError:
		return v.Pos
	}
	return token.NoPos
}

func ignNotatedErrs(err error, pkg *ast.Package, fset *token.FileSet) error {
	switch v := err.(type) {
	case errors.List:
		var ret errors.List
		for _, e := range v {
			if isNotatedErr(e, pkg, fset) {
				continue
			}
			ret = append(ret, e)
		}
		if len(ret) == 0 {
			return ErrIgnoreNotated
		}
		return ret
	default:
		if isNotatedErr(err, pkg, fset) {
			return ErrIgnoreNotated
		}
		return err
	}
}

func isNotatedErr(err error, pkg *ast.Package, fset *token.FileSet) (notatedErr bool) {
	pos := ErrorPos(err)
	f := fset.File(pos)
	if f == nil {
		return
	}
	gopf, ok := pkg.Files[f.Name()]
	if !ok {
		return
	}
	lines := token.Lines(f)
	i := f.Line(pos) - 1 // base 0
	start := lines[i]
	var end int
	if i+1 < len(lines) {
		end = lines[i+1]
	} else {
		end = f.Size()
	}
	text := string(gopf.Code[start:end])
	commentOff := strings.Index(text, "//")
	if commentOff < 0 {
		return
	}
	return strings.Contains(text[commentOff+2:], "compile error:")
}

// -----------------------------------------------------------------------------

type Config struct {
	Gop      *env.Gop
	Fset     *token.FileSet
	Filter   func(fs.FileInfo) bool
	Importer types.Importer

	DontUpdateGoMod     bool
	DontCheckModChanged bool
	IgnoreNotatedError  bool
}

func LoadMod(dir string, gop *env.Gop, conf *Config) (mod *gopmod.Module, err error) {
	mod, err = gopmod.Load(dir, 0)
	if err != nil && !NotFound(err) {
		err = errors.NewWith(err, `gopmod.Load(dir, 0)`, -2, "gopmod.Load", dir, 0)
		return
	}
	if mod != nil {
		err = mod.ImportClasses()
		if err != nil {
			err = errors.NewWith(err, `mod.RegisterClasses()`, -2, "(*gopmod.Module).RegisterClasses", mod)
			return
		}
		if !conf.DontUpdateGoMod {
			err = mod.UpdateGoMod(gop, !conf.DontCheckModChanged)
			if err != nil {
				err = errors.NewWith(err, `mod.UpdateGoMod(gop, !conf.DontCheckModChanged)`, -2, "(*gopmod.Module).UpdateGoMod", mod, gop, !conf.DontCheckModChanged)
			}
		}
		return
	}
	return new(gopmod.Module), nil
}

// -----------------------------------------------------------------------------

func LoadDir(dir string, conf *Config, genTestPkg bool, promptGenGo ...bool) (out, test *gox.Package, err error) {
	if conf == nil {
		conf = new(Config)
	}
	fset := conf.Fset
	if fset == nil {
		fset = token.NewFileSet()
	}
	gop := conf.Gop
	if gop == nil {
		gop = gopenv.Get()
	}
	mod, err := LoadMod(dir, gop, conf)
	if err != nil {
		return
	}

	pkgs, err := parser.ParseDirEx(fset, dir, parser.Config{
		ClassKind: mod.ClassKind,
		Filter:    conf.Filter,
		Mode:      parser.ParseComments,
	})
	if err != nil {
		return
	}
	if len(pkgs) == 0 {
		return nil, nil, ErrNotFound
	}

	imp := conf.Importer
	if imp == nil {
		imp = NewImporter(mod, gop, fset)
	}

	var pkgTest *ast.Package
	var clConf = &cl.Config{
		WorkingDir:   dir,
		Fset:         fset,
		Importer:     imp,
		LookupClass:  mod.LookupClass,
		LookupPub:    c2go.LookupPub(mod),
		FileLineRoot: mod.Root(),
	}
	for name, pkg := range pkgs {
		if strings.HasSuffix(name, "_test") {
			if pkgTest != nil {
				return nil, nil, ErrMultiTestPackges
			}
			pkgTest = pkg
			continue
		}
		if out != nil {
			return nil, nil, ErrMultiPackges
		}
		if len(pkg.Files) == 0 { // no Go+ source files
			continue
		}
		if promptGenGo != nil && promptGenGo[0] {
			fmt.Fprintln(os.Stderr, "GenGo", dir, "...")
		}
		out, err = cl.NewPackage("", pkg, clConf)
		if err != nil {
			if conf.IgnoreNotatedError {
				err = ignNotatedErrs(err, pkg, fset)
			}
			return
		}
	}
	if out == nil {
		return nil, nil, ErrNotFound
	}
	if pkgTest != nil && genTestPkg {
		test, err = cl.NewPackage("", pkgTest, clConf)
	}
	return
}

// -----------------------------------------------------------------------------

func LoadFiles(files []string, conf *Config) (out *gox.Package, err error) {
	if conf == nil {
		conf = new(Config)
	}
	fset := conf.Fset
	if fset == nil {
		fset = token.NewFileSet()
	}
	gop := conf.Gop
	if gop == nil {
		gop = gopenv.Get()
	}
	mod, err := LoadMod("", gop, conf)
	if err != nil {
		err = errors.NewWith(err, `loadMod("", gop, conf)`, -2, "gop.loadMod", "", gop, conf)
		return
	}

	pkgs, err := parser.ParseFiles(fset, files, parser.ParseComments)
	if err != nil {
		err = errors.NewWith(err, `parser.ParseFiles(fset, files, parser.ParseComments)`, -2, "parser.ParseFiles", fset, files, parser.ParseComments)
		return
	}
	if len(pkgs) != 1 {
		err = errors.NewWith(ErrMultiPackges, `len(pkgs) != 1`, -1, "!=", len(pkgs), 1)
		return
	}
	for _, pkg := range pkgs {
		imp := conf.Importer
		if imp == nil {
			imp = NewImporter(mod, gop, fset)
		}
		clConf := &cl.Config{
			Fset:        fset,
			Importer:    imp,
			LookupClass: mod.LookupClass,
			LookupPub:   c2go.LookupPub(mod),
		}
		out, err = cl.NewPackage("", pkg, clConf)
		if err != nil {
			if conf.IgnoreNotatedError {
				err = ignNotatedErrs(err, pkg, fset)
			}
		}
		break
	}
	return
}

// -----------------------------------------------------------------------------

var (
	ErrMultiPackges     = errors.New("multiple packages")
	ErrMultiTestPackges = errors.New("multiple test packages")
)

// -----------------------------------------------------------------------------
