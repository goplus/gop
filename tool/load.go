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

package tool

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/goplus/gogen"
	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/goplus/gop/x/gocmd"
	"github.com/goplus/gop/x/gopenv"
	"github.com/goplus/mod/env"
	"github.com/goplus/mod/gopmod"
	"github.com/qiniu/x/errors"
)

var (
	ErrNotFound      = gopmod.ErrNotFound
	ErrIgnoreNotated = errors.New("notated error ignored")
)

// NotFound returns if cause err is ErrNotFound or not
func NotFound(err error) bool {
	return gopmod.IsNotFound(err)
}

// IgnoreNotated returns if cause err is ErrIgnoreNotated or not.
func IgnoreNotated(err error) bool {
	return errors.Err(err) == ErrIgnoreNotated
}

// ErrorPos returns where the error occurs.
func ErrorPos(err error) token.Pos {
	switch v := err.(type) {
	case *gogen.CodeError:
		return v.Pos
	case *gogen.MatchError:
		return v.Pos()
	case *gogen.ImportError:
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

// Config represents a configuration for loading Go+ packages.
type Config struct {
	Gop      *env.Gop
	Fset     *token.FileSet
	Mod      *gopmod.Module
	Importer *Importer

	Filter func(fs.FileInfo) bool

	// If not nil, it is used for returning result of checks Go+ dependencies.
	// see https://pkg.go.dev/github.com/goplus/gogen#File.CheckGopDeps
	GopDeps *int

	// CacheFile specifies the file path of the cache.
	CacheFile string

	IgnoreNotatedError bool
	DontUpdateGoMod    bool
}

// ConfFlags represents configuration flags.
type ConfFlags int

const (
	ConfFlagIgnoreNotatedError ConfFlags = 1 << iota
	ConfFlagDontUpdateGoMod
	ConfFlagNoTestFiles
	ConfFlagNoCacheFile
)

// NewDefaultConf creates a dfault configuration for common cases.
func NewDefaultConf(dir string, flags ConfFlags, tags ...string) (conf *Config, err error) {
	mod, err := LoadMod(dir)
	if err != nil {
		return
	}
	gop := gopenv.Get()
	fset := token.NewFileSet()
	imp := NewImporter(mod, gop, fset)
	if len(tags) > 0 {
		imp.SetTags(strings.Join(tags, ","))
	}
	conf = &Config{
		Gop: gop, Fset: fset, Mod: mod, Importer: imp,
		IgnoreNotatedError: flags&ConfFlagIgnoreNotatedError != 0,
		DontUpdateGoMod:    flags&ConfFlagDontUpdateGoMod != 0,
	}
	if flags&ConfFlagNoCacheFile == 0 {
		conf.CacheFile = imp.CacheFile()
		imp.Cache().Load(conf.CacheFile)
	}
	if flags&ConfFlagNoTestFiles != 0 {
		conf.Filter = FilterNoTestFiles
	}
	return
}

func (conf *Config) NewGoCmdConf() *gocmd.Config {
	if cl := conf.Mod.Opt.Compiler; cl != nil {
		if os.Getenv("GOP_GOCMD") == "" {
			os.Setenv("GOP_GOCMD", cl.Name)
		}
	}
	return &gocmd.Config{
		Gop: conf.Gop,
	}
}

// UpdateCache updates the cache.
func (conf *Config) UpdateCache(verbose ...bool) {
	if conf.CacheFile != "" {
		c := conf.Importer.Cache()
		c.Save(conf.CacheFile)
		if verbose != nil && verbose[0] {
			fmt.Println("Times of calling go list:", c.ListTimes())
		}
	}
}

// LoadMod loads a Go+ module from a specified directory.
func LoadMod(dir string) (mod *gopmod.Module, err error) {
	mod, err = gopmod.Load(dir)
	if err != nil && !gopmod.IsNotFound(err) {
		err = errors.NewWith(err, `gopmod.Load(dir, 0)`, -2, "gopmod.Load", dir, 0)
		return
	}
	if mod == nil {
		mod = gopmod.Default
	}
	err = mod.ImportClasses()
	if err != nil {
		err = errors.NewWith(err, `mod.ImportClasses()`, -2, "(*gopmod.Module).ImportClasses", mod)
	}
	return
}

// FilterNoTestFiles filters to skip all testing files.
func FilterNoTestFiles(fi fs.FileInfo) bool {
	fname := fi.Name()
	suffix := ""
	switch path.Ext(fname) {
	case ".gox":
		suffix = "test.gox"
	case ".gop":
		suffix = "_test.gop"
	case ".go":
		suffix = "_test.go"
	default:
		return true
	}
	return !strings.HasSuffix(fname, suffix)
}

// -----------------------------------------------------------------------------

// LoadDir loads Go+ packages from a specified directory.
func LoadDir(dir string, conf *Config, genTestPkg bool, promptGenGo ...bool) (out, test *gogen.Package, err error) {
	if conf == nil {
		conf = new(Config)
	}

	mod := conf.Mod
	if mod == nil {
		if mod, err = LoadMod(dir); err != nil {
			err = errors.NewWith(err, `LoadMod(dir)`, -2, "tool.LoadMod", dir)
			return
		}
	}

	fset := conf.Fset
	if fset == nil {
		fset = token.NewFileSet()
	}
	pkgs, err := parser.ParseDirEx(fset, dir, parser.Config{
		ClassKind: mod.ClassKind,
		Filter:    conf.Filter,
		Mode:      parser.ParseComments | parser.SaveAbsFile,
	})
	if err != nil {
		return
	}
	if len(pkgs) == 0 {
		return nil, nil, ErrNotFound
	}

	gop := conf.Gop
	if gop == nil {
		gop = gopenv.Get()
	}
	imp := conf.Importer
	if imp == nil {
		imp = NewImporter(mod, gop, fset)
	}

	var pkgTest *ast.Package
	var clConf = &cl.Config{
		Fset:         fset,
		RelativeBase: relativeBaseOf(mod),
		Importer:     imp,
		LookupClass:  mod.LookupClass,
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
	if genTestPkg && pkgTest != nil {
		test, err = cl.NewPackage("", pkgTest, clConf)
	}
	afterLoad(mod, gop, out, test, conf)
	return
}

func afterLoad(mod *gopmod.Module, gop *env.Gop, out, test *gogen.Package, conf *Config) {
	if mod.Path() == gopMod { // nothing to do for Go+ itself
		return
	}
	updateMod := !conf.DontUpdateGoMod && mod.HasModfile()
	if updateMod || conf.GopDeps != nil {
		flags := checkGopDeps(out)
		if conf.GopDeps != nil { // for `gop run`
			*conf.GopDeps = flags
		}
		if updateMod {
			if test != nil {
				flags |= checkGopDeps(test)
			}
			if flags != 0 {
				mod.SaveWithGopMod(gop, flags)
			}
		}
	}
}

func checkGopDeps(pkg *gogen.Package) (flags int) {
	pkg.ForEachFile(func(fname string, file *gogen.File) {
		flags |= file.CheckGopDeps(pkg)
	})
	return
}

func relativeBaseOf(mod *gopmod.Module) string {
	if mod.HasModfile() {
		return mod.Root()
	}
	dir, _ := os.Getwd()
	return dir
}

// -----------------------------------------------------------------------------

// LoadFiles loads a Go+ package from specified files.
func LoadFiles(dir string, files []string, conf *Config) (out *gogen.Package, err error) {
	if conf == nil {
		conf = new(Config)
	}

	mod := conf.Mod
	if mod == nil {
		if mod, err = LoadMod(dir); err != nil {
			err = errors.NewWith(err, `LoadMod(dir)`, -2, "tool.LoadMod", dir)
			return
		}
	}

	fset := conf.Fset
	if fset == nil {
		fset = token.NewFileSet()
	}
	pkgs, err := parser.ParseEntries(fset, files, parser.Config{
		ClassKind: mod.ClassKind,
		Filter:    conf.Filter,
		Mode:      parser.ParseComments | parser.SaveAbsFile,
	})
	if err != nil {
		err = errors.NewWith(err, `parser.ParseFiles(fset, files, parser.ParseComments)`, -2, "parser.ParseFiles", fset, files, parser.ParseComments)
		return
	}
	if len(pkgs) != 1 {
		err = errors.NewWith(ErrMultiPackges, `len(pkgs) != 1`, -1, "!=", len(pkgs), 1)
		return
	}
	gop := conf.Gop
	if gop == nil {
		gop = gopenv.Get()
	}
	for _, pkg := range pkgs {
		imp := conf.Importer
		if imp == nil {
			imp = NewImporter(mod, gop, fset)
		}
		clConf := &cl.Config{
			Fset:         fset,
			RelativeBase: relativeBaseOf(mod),
			Importer:     imp,
			LookupClass:  mod.LookupClass,
		}
		out, err = cl.NewPackage("", pkg, clConf)
		if err != nil {
			if conf.IgnoreNotatedError {
				err = ignNotatedErrs(err, pkg, fset)
			}
		}
		break
	}
	afterLoad(mod, gop, out, nil, conf)
	return
}

// -----------------------------------------------------------------------------

var (
	ErrMultiPackges     = errors.New("multiple packages")
	ErrMultiTestPackges = errors.New("multiple test packages")
)

// -----------------------------------------------------------------------------

// GetFileClassType get gop module file classType.
func GetFileClassType(mod *gopmod.Module, file *ast.File, filename string) (classType string, isTest bool) {
	return cl.GetFileClassType(file, filename, mod.LookupClass)
}

// -----------------------------------------------------------------------------
