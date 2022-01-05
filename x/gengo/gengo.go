/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
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

package gengo

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"go/types"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/goplus/gop/x/gopmod"
	"github.com/goplus/gop/x/mod/modfetch"
	"github.com/goplus/gox"
	"github.com/goplus/gox/packages"
)

const (
	pkgFlagIll = 1 << iota
	pkgFlagChanged
)

const (
	stateNormal = iota
	stateProcessing
	stateDone
	stateOccurErrors = -1
)

type pkgInfo struct {
	path    string
	runner  *Runner
	imports []string
	deps    []*pkgInfo
	hash    string
	nsource int
	flags   int
}

type Runner struct {
	runners map[string]*Runner
	mod     *gopmod.Module
	modRoot string
	modPath string
	modTime time.Time
	pkgs    map[string]*pkgInfo
	state   int
}

type Event interface {
	OnStart(pkgPath string)
	OnErr(stage string, err error)
	OnInfo(format string, args ...interface{})
	OnEnd()
}

type defaultEvent struct{}

func (p defaultEvent) OnStart(pkgPath string) {
	fmt.Fprintln(os.Stderr, pkgPath)
}

func (p defaultEvent) OnInfo(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func (p defaultEvent) OnErr(stage string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", stage, err)
}

func (p defaultEvent) OnEnd() {
}

type Config struct {
	Event
	Fset   *token.FileSet
	Loaded map[string]*types.Package
}

func (p *Runner) GenGoFlow() (pkgPaths []string, err error) {
	pkgPaths = make([]string, 0, len(p.pkgs))
	return p.genFlowPkgs(pkgPaths)
}

func (p *Runner) genFlowPkgs(pkgPaths []string) ([]string, error) {
	if p.state != stateNormal {
		if p.state != stateDone {
			panic("TODO: cycle or error?")
		}
		return pkgPaths, nil
	}
	p.state = stateProcessing
	for _, pkg := range p.pkgs {
		if pkg.flags == pkgFlagChanged {
			for _, dep := range pkg.deps {
				if r := dep.runner; r != p {
					pkgPaths, _ = dep.runner.genFlowPkgs(pkgPaths)
				}
			}
		}
	}
	for _, pkg := range p.pkgs {
		if pkg.flags == pkgFlagChanged {
			pkgPaths = p.genFlowPkg(pkgPaths, pkg)
		}
	}
	p.state = stateDone
	return pkgPaths, nil
}

func (p *Runner) genFlowPkg(pkgPaths []string, pkg *pkgInfo) []string {
	for _, dep := range pkg.deps {
		if r := dep.runner; r == p && dep.flags == pkgFlagChanged {
			pkgPaths = r.genFlowPkg(pkgPaths, dep)
		}
	}
	pkgPaths = append(pkgPaths, pkg.path)
	pkg.flags = 0
	return pkgPaths
}

func (p *Runner) GenGo(conf Config) bool {
	if conf.Loaded == nil {
		conf.Loaded = make(map[string]*types.Package)
	}
	if conf.Fset == nil {
		conf.Fset = token.NewFileSet()
	}
	if conf.Event == nil {
		conf.Event = defaultEvent{}
	}
	p.genGoPkgs(&conf)
	return p.state == stateDone
}

func (p *Runner) genGoPkgs(conf *Config) {
	if p.state != stateNormal {
		if p.state != stateDone {
			panic("TODO: cycle or error?")
		}
		return
	}
	p.state = stateProcessing
	imports, changed := make(map[string]none), false
	for _, pkg := range p.pkgs {
		if pkg.flags == pkgFlagChanged {
			changed = true
			for _, dep := range pkg.deps {
				if r := dep.runner; r != p {
					r.genGoPkgs(conf)
				}
			}
			for _, imp := range pkg.imports {
				imports[imp] = none{}
			}
		}
	}
	if !changed {
		p.state = stateDone
		return
	}
	imps := getKeys(imports)
	impConf := &packages.Config{
		ModRoot: p.modRoot,
		ModPath: p.modPath,
		Loaded:  conf.Loaded,
		Fset:    conf.Fset,
	}
	imp, _, err := packages.NewImporter(impConf, imps...)
	if err != nil {
		conf.OnErr("newImporter", err)
		p.state = stateOccurErrors
		return
	}
	defer imp.Close()
	conf.OnInfo("newImporter: %v\n", imps)
	for _, pkg := range p.pkgs {
		if pkg.flags == pkgFlagChanged {
			p.genGoPkg(pkg, imp, conf)
		}
	}
	if p.state == stateProcessing {
		p.state = stateDone
	}
}

func getKeys(v map[string]none) []string {
	keys := make([]string, 0, len(v))
	for key := range v {
		keys = append(keys, key)
	}
	return keys
}

const (
	autoGenFile      = "gop_autogen.go"
	autoGenTestFile  = "gop_autogen_test.go"
	autoGen2TestFile = "gop_autogen2_test.go"
)

func (p *Runner) genGoPkg(pkg *pkgInfo, imp types.Importer, conf *Config) {
	for _, dep := range pkg.deps {
		if r := dep.runner; r == p && dep.flags == pkgFlagChanged {
			r.genGoPkg(dep, imp, conf)
		}
	}
	if err := p.doGenGoPkg(pkg, imp, conf); err != nil {
		pkg.flags |= pkgFlagIll
		p.state = stateOccurErrors
		return
	}
	pkg.flags = 0
}

func (p *Runner) doGenGoPkg(pkgi *pkgInfo, imp types.Importer, conf *Config) (err error) {
	pkgPath := pkgi.path
	conf.OnStart(pkgPath)
	pkgDir := p.modRoot + pkgPath[len(p.modPath):]

	defer func() {
		if e := recover(); e != nil {
			switch v := e.(type) {
			case error:
				err = v
			case string:
				err = errors.New(v)
			default:
				panic(e)
			}
			conf.OnErr("genGoPkg", err)
		}
		conf.OnEnd()
	}()

	pkgs, err := parser.ParseDirEx(conf.Fset, pkgDir, parser.Config{
		IsClass: p.mod.IsClass,
		Mode:    parser.ParseComments,
	})
	if err != nil {
		conf.OnErr("parse", err)
		return
	}
	if len(pkgs) == 0 {
		panic("TODO: no source in directory")
	}

	var pkgTest *ast.Package
	clConf := &cl.Config{
		WorkingDir:  pkgDir,
		Fset:        conf.Fset,
		LookupClass: p.mod.LookupClass,
		Importer:    imp,
	}
	for name, pkg := range pkgs {
		if strings.HasSuffix(name, "_test") {
			if pkgTest != nil {
				panic("TODO: has multi test package?")
			}
			pkgTest = pkg
			continue
		}
		out, e := cl.NewPackage("", pkg, clConf)
		if e != nil {
			conf.OnErr("compile", e)
			return e
		}
		err = saveGoFile(pkgDir, out, pkgi)
		if err != nil {
			conf.OnErr("compile", err)
			return
		}
	}
	if pkgTest != nil {
		out, e := cl.NewPackage("", pkgTest, clConf)
		if e != nil {
			conf.OnErr("compile", e)
			return e
		}
		err = gox.WriteFile(filepath.Join(pkgDir, autoGen2TestFile), out, true)
		if err != nil {
			conf.OnErr("compile", err)
		}
	}
	return
}

func saveGoFile(dir string, pkg *gox.Package, pkgi *pkgInfo) error {
	os.MkdirAll(dir, 0755)
	err := func() (err error) {
		file := filepath.Join(dir, autoGenFile)
		f, err := os.Create(file)
		if err != nil {
			return
		}
		err = syscall.EFAULT
		defer func() {
			f.Close()
			if err != nil {
				os.Remove(file)
			}
		}()
		err = pkgi.writeCinfo(f)
		if err != nil {
			return
		}
		return gox.WriteTo(f, pkg, false)
	}()
	if err != nil {
		return err
	}
	if pkg.HasTestingFile() {
		return gox.WriteFile(filepath.Join(dir, autoGenTestFile), pkg, true)
	}
	return nil
}

// -----------------------------------------------------------------------------

type none = struct{}

func modTime(file string, errs *ErrorList) (mt time.Time) {
	fi, err := os.Stat(file)
	if err != nil {
		*errs = append(*errs, err)
		return
	}
	return fi.ModTime()
}

func New(mod *gopmod.Module, pattern ...string) (p *Runner, err error) {
	pkgPaths := make(map[string]struct{})
	err = mod.List(pkgPaths, pattern...)
	if err != nil {
		return
	}
	modFile := mod.Modfile()
	modRoot := filepath.Dir(modFile)

	var errs ErrorList
	p = &Runner{
		mod:     mod,
		modRoot: modRoot,
		modPath: mod.Path(),
		modTime: modTime(modFile, &errs),
		pkgs:    make(map[string]*pkgInfo),
	}
	p.runners = map[string]*Runner{modRoot: p}
	for pkgPath := range pkgPaths {
		p.genDeps(pkgPath, pkgPath, &errs)
	}
	if len(errs) > 0 {
		err = errs
	}
	return
}

func (p *Runner) genDeps(pkgPath, pkgPathBase string, errs *ErrorList) (pkg *pkgInfo) {
	switch p.mod.PkgType(pkgPath) {
	case gopmod.PkgtStandard:
		return nil
	case gopmod.PkgtLocal:
		pkgPath = filepath.ToSlash(filepath.Join(pkgPathBase, pkgPath))
	case gopmod.PkgtModule:
	case gopmod.PkgtExtern:
		return p.genExternDeps(pkgPath, errs)
	default:
		panic("TODO: invalid pkgPath")
	}
	if v, ok := p.pkgs[pkgPath]; ok {
		return v
	}
	pkg = &pkgInfo{
		path:   pkgPath,
		runner: p,
	}
	p.pkgs[pkgPath] = pkg
	dir := filepath.Join(p.modRoot, pkgPath[len(p.modPath):])

	ci, err := p.mod.ChangeInfo(dir)
	if err != nil {
		*errs, pkg.flags = append(*errs, err), pkgFlagIll
		return
	}
	if ci.GopFileNum == 0 {
		return
	}

	var buf bytes.Buffer
	imports, err := p.getImports(dir)
	if err != nil {
		errs.addError("import "+pkgPath, err)
		pkg.flags = pkgFlagIll
		return
	}
	imps := getSortedKeys(imports)
	for _, imp := range imps {
		if t := p.genDeps(imp, pkgPath, errs); t != nil {
			buf.WriteString(t.runner.modRoot)
			pkg.flags |= t.flags
			if t.flags == pkgFlagChanged {
				pkg.deps = append(pkg.deps, t)
			}
		} else {
			buf.WriteString(imp)
		}
		buf.WriteByte('\n')
	}
	pkg.imports = imps
	pkg.hash = hashOf(buf.Bytes())
	pkg.nsource = ci.SourceNum
	if p.sourceChanged(dir, ci.ModTime, pkg) {
		pkg.flags |= pkgFlagChanged
	}
	return
}

func (p *Runner) getImports(dir string) (ret map[string]struct{}, err error) {
	pkgs, err := p.mod.DirImports(dir)
	if err != nil {
		return
	}
	exists := false
	for name, imps := range pkgs {
		if !strings.HasSuffix(name, "_test") {
			if exists {
				return nil, errors.New("multiple packages in a directory")
			}
			ret, exists = imps, true
		}
	}
	if exists {
		return
	}
	return nil, syscall.ENOENT
}

func getSortedKeys(v map[string]none) []string {
	keys := make([]string, 0, len(v))
	for key := range v {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func hashOf(b []byte) string {
	h := md5.Sum(b)
	return hex.EncodeToString(h[:])
}

func (p *Runner) genExternDeps(pkgPath string, errs *ErrorList) *pkgInfo {
	_, modVer, ok := p.mod.LookupExternPkg(pkgPath)
	if !ok {
		log.Panicln("externPkg not found:", pkgPath)
	}
	modRoot, err := modfetch.ModCachePath(modVer)
	if err != nil {
		*errs = append(*errs, err)
		return &pkgInfo{path: pkgPath, flags: pkgFlagIll}
	}
	exr, ok := p.runners[modRoot]
	if !ok {
		mod, err := gopmod.LoadMod(modVer)
		if err != nil {
			*errs = append(*errs, err)
			return &pkgInfo{path: pkgPath, flags: pkgFlagIll}
		}
		exr = &Runner{
			mod:     mod,
			modRoot: modRoot,
			modPath: mod.Path(),
			modTime: modTime(mod.Modfile(), errs),
			pkgs:    make(map[string]*pkgInfo),
			runners: p.runners,
		}
		p.runners[modRoot] = exr
	}
	return exr.genDeps(pkgPath, pkgPath, errs)
}

func (p *Runner) sourceChanged(dir string, modt time.Time, pkg *pkgInfo) bool {
	f, err := os.Open(dir + "/gop_autogen.go")
	if err != nil {
		return true
	}
	fi, err := f.Stat()
	if err != nil || modt.After(fi.ModTime()) {
		return true
	}
	//cinfo $SourceNum $DepsHash
	var sourceNum int
	var hash string
	if _, err = fmt.Fscanf(f, "//cinfo %d %s\n", &sourceNum, &hash); err != nil {
		return true
	}
	return pkg.nsource != sourceNum || pkg.hash != hash
}

func (p pkgInfo) writeCinfo(w io.Writer) error {
	_, err := fmt.Fprintf(w, "//cinfo %d %s\n", p.nsource, p.hash)
	return err
}

// -----------------------------------------------------------------------------

type ErrorList []error

func (p *ErrorList) addError(stage string, err error) {
	*p = append(*p, fmt.Errorf("%s: %v", stage, err))
}

func (e ErrorList) Error() string {
	errStrs := make([]string, len(e))
	for i, err := range e {
		errStrs[i] = err.Error()
	}
	return strings.Join(errStrs, "\n")
}

// -----------------------------------------------------------------------------

func GenGo(conf Config, dontRun bool, dir string, pattern ...string) bool {
	if conf.Event == nil {
		conf.Event = defaultEvent{}
	}
	mod, err := gopmod.Load(dir)
	if err != nil {
		conf.OnErr("loadMod", err)
		return false
	}
	err = mod.RegisterClasses()
	if err != nil {
		conf.OnErr("registerClass", err)
		return false
	}
	p, err := New(mod, pattern...)
	if err != nil {
		conf.OnErr("genDeps", err)
		return false
	}
	if dontRun {
		pkgPaths, err := p.GenGoFlow()
		if err != nil {
			conf.OnErr("genGo", err)
			return false
		}
		for _, pkgPath := range pkgPaths {
			fmt.Fprintln(os.Stderr, pkgPath)
		}
		return true
	}
	return p.GenGo(conf)
}

// -----------------------------------------------------------------------------
