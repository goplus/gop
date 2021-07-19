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

package gengo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

const (
	PkgFlagGo = 1 << iota
	PkgFlagGoPlus
	PkgFlagGoGen
	PkgFlagGopModified
	PkgFlagSpx
)

const (
	autoGenFileName = "gop_autogen.go"
)

type Error struct {
	PkgDir string
	Stage  string
	Err    error
}

func (p *Error) Error() string { // TODO: more friendly
	return p.Err.Error()
}

// -----------------------------------------------------------------------------

type Config = cl.Config

type Runner struct {
	conf  *Config
	fset  *token.FileSet
	errs  []*Error
	after func(p *Runner, dir string, pkgFlags int) error
}

func NewRunner(fset *token.FileSet, conf *Config) *Runner {
	if fset == nil {
		fset = token.NewFileSet()
	}
	return &Runner{fset: fset, conf: conf}
}

func (p *Runner) SetAfter(after func(p *Runner, dir string, flags int) error) {
	p.after = after
}

func (p *Runner) Errors() []*Error {
	return p.errs
}

func (p *Runner) ResetErrors() []*Error {
	errs := p.errs
	p.errs = nil
	return errs
}

func (p *Runner) GenGo(dir string, recursive bool) {
	if strings.HasPrefix(dir, "_") {
		return
	}
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		p.addError(dir, "readDir", err)
		return
	}
	var gopTime time.Time
	var gogenTime time.Time
	var pkgFlags int
	for _, fi := range fis {
		fname := fi.Name()
		if fi.IsDir() {
			if recursive {
				pkgDir := path.Join(dir, fname)
				p.GenGo(pkgDir, true)
			}
			continue
		}
		ext := filepath.Ext(fname)
		if flag, ok := extPkgFlags[ext]; ok {
			modTime := fi.ModTime()
			switch flag {
			case PkgFlagGo:
				if fname == autoGenFileName {
					flag = PkgFlagGoGen
					gogenTime = modTime
				}
			default:
				if modTime.After(gopTime) {
					gopTime = modTime
				}
			}
			pkgFlags |= flag
		}
	}
	if pkgFlags != 0 {
		if pkgFlags == PkgFlagGo { // a pure Go package
			// TODO: depency check
		} else if gopTime.After(gogenTime) { // update a Go+ package
			fmt.Printf("GenGoPkg %s ...\n", dir)
			pkgFlags |= PkgFlagGopModified
			p.GenGoPkg(dir)
		}
		if p.after != nil {
			if err = p.after(p, dir, pkgFlags); err != nil {
				p.addError(dir, "after", err)
			}
		}
	}
}

var (
	extPkgFlags = map[string]int{
		".gop": PkgFlagGoPlus,
		".spx": PkgFlagSpx,
		".go":  PkgFlagGo,
	}
)

func (p *Runner) GenGoPkg(pkgDir string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch v := e.(type) {
			case string:
				err = p.addError(pkgDir, "recover", errors.New(v))
			case error:
				err = p.addError(pkgDir, "recover", v)
			default:
				panic(e)
			}
		}
	}()

	pkgDir, _ = filepath.Abs(pkgDir)
	pkgs, err := parser.ParseDir(p.fset, pkgDir, nil, 0)
	if err != nil {
		return p.addError(pkgDir, "parse", err)
	}
	if len(pkgs) != 1 {
		err = fmt.Errorf("too many packages (%d) in the same directory", len(pkgs))
		return p.addError(pkgDir, "parse", err)
	}

	pkg := getPkg(pkgs)
	out, err := cl.NewPackage("", pkg, p.fset, p.conf)
	if err != nil {
		return p.addError(pkgDir, "compile", err)
	}
	err = saveGoFile(pkgDir, out)
	if err != nil {
		return p.addError(pkgDir, "save", err)
	}
	return nil
}

func (p *Runner) addError(pkgDir string, stage string, err error) error {
	e := &Error{PkgDir: pkgDir, Stage: stage, Err: err}
	p.errs = append(p.errs, e)
	return e
}

func saveGoFile(dir string, pkg *gox.Package) error {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	return gox.WriteFile(filepath.Join(dir, autoGenFileName), pkg)
}

func getPkg(pkgs map[string]*ast.Package) *ast.Package {
	for _, pkg := range pkgs {
		return pkg
	}
	return nil
}

// -----------------------------------------------------------------------------
