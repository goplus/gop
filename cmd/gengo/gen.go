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
	PkgFlagGmx
)

const (
	autoGenFile      = "gop_autogen.go"
	autoGenTestFile  = "gop_autogen_test.go"
	autoGen2TestFile = "gop_autogen2_test.go"
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

type Runner struct {
	errs  []*Error
	after func(p *Runner, dir string, pkgFlags int) error
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

func (p *Runner) GenGo(dir string, recursive bool, base *cl.Config) {
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
		if strings.HasPrefix(fname, "_") {
			continue
		}
		if fi.IsDir() {
			if recursive {
				pkgDir := path.Join(dir, fname)
				p.GenGo(pkgDir, true, base)
			}
			continue
		}
		ext := filepath.Ext(fname)
		if flag, ok := extPkgFlags[ext]; ok {
			modTime := fi.ModTime()
			switch flag {
			case PkgFlagGo:
				switch fname {
				case autoGenFile, autoGenTestFile, autoGen2TestFile:
					flag = PkgFlagGoGen
					if (pkgFlags&PkgFlagGoGen) == 0 || gogenTime.After(modTime) {
						gogenTime = modTime
					}
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
		if (pkgFlags & PkgFlagGo) != 0 { // a Go package
			// TODO: depency check
		} else if gopTime.After(gogenTime) { // update a Go+ package
			fmt.Printf("GenGoPkg %s\n", dir)
			pkgFlags |= PkgFlagGopModified
			p.GenGoPkg(dir, base)
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
		".gmx": PkgFlagGmx,
		".go":  PkgFlagGo,
	}
)

func (p *Runner) GenGoPkg(pkgDir string, base *cl.Config) (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch v := e.(type) {
			case error:
				err = p.addError(pkgDir, "recover", v)
			case string:
				err = p.addError(pkgDir, "recover", errors.New(v))
			default:
				panic(e)
			}
		}
	}()

	pkgDir, _ = filepath.Abs(pkgDir)

	conf := *base.Ensure()
	conf.Dir = pkgDir
	if conf.Fset == nil {
		conf.Fset = token.NewFileSet()
	}
	pkgs, err := parser.ParseDir(conf.Fset, pkgDir, nil, 0)
	if err != nil {
		return p.addError(pkgDir, "parse", err)
	}

	var pkgTest *ast.Package
	for name, pkg := range pkgs {
		if strings.HasSuffix(name, "_test") {
			if pkgTest != nil {
				panic("TODO: has multi test package?")
			}
			pkgTest = pkg
			continue
		}
		out, err := cl.NewPackage("", pkg, &conf)
		if err != nil {
			return p.addError(pkgDir, "compile", err)
		}
		err = saveGoFile(pkgDir, out)
		if err != nil {
			return p.addError(pkgDir, "save", err)
		}
	}
	if pkgTest != nil {
		out, err := cl.NewPackage("", pkgTest, &conf)
		if err != nil {
			return p.addError(pkgDir, "compile", err)
		}
		err = gox.WriteFile(filepath.Join(pkgDir, autoGen2TestFile), out, true)
		if err != nil {
			return p.addError(pkgDir, "save", err)
		}
	}
	return nil
}

func (p *Runner) addError(pkgDir string, stage string, err error) error {
	e := &Error{PkgDir: pkgDir, Stage: stage, Err: err}
	p.errs = append(p.errs, e)
	return e
}

func saveGoFile(dir string, pkg *gox.Package) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	err = gox.WriteFile(filepath.Join(dir, autoGenFile), pkg, false)
	if err != nil {
		return err
	}
	if pkg.HasTestingFile() {
		return gox.WriteFile(filepath.Join(dir, autoGenTestFile), pkg, true)
	}
	return nil
}

// -----------------------------------------------------------------------------
