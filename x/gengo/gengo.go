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
	"fmt"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/goplus/gop/x/gopmod"
	"github.com/goplus/gop/x/mod/modfetch"
	"github.com/goplus/gox/packages"
)

const (
	pkgFlagIll = 1 << iota
	pkgFlagChanged
)

type pkgInfo struct {
	path    string
	runner  *Runner
	imports []string
	deps    []*pkgInfo
	hash    string
	flags   int
}

type Runner struct {
	runners map[string]*Runner
	mod     *gopmod.Module
	modRoot string
	modPath string
	modTime time.Time
	pkgs    map[string]*pkgInfo
}

func (p *Runner) GenGo(fset *token.FileSet, loaded map[string]*types.Package) (err error) {
	if loaded == nil {
		loaded = make(map[string]*types.Package)
	}
	if fset == nil {
		fset = token.NewFileSet()
	}
	var errs ErrorList
	p.genGoPkgs(fset, loaded, &errs)
	if len(errs) > 0 {
		err = errs
	}
	return
}

func (p *Runner) genGoPkgs(fset *token.FileSet, loaded map[string]*types.Package, errs *ErrorList) {
	imports := make(map[string]none)
	for path, pkg := range p.pkgs {
		if _, ok := loaded[path]; ok { // loaded
			continue
		}
		if pkg.flags == pkgFlagChanged {
			for _, imp := range pkg.imports {
				imports[imp] = none{}
			}
		}
	}
	imps := getKeys(imports)
	conf := &packages.Config{
		ModRoot:       p.modRoot,
		ModPath:       p.modPath,
		SupportedExts: p.mod.SupportedExts(),
		Loaded:        loaded,
		Fset:          fset,
	}
	imp, _, err := packages.NewImporter(conf, imps...)
	if err != nil {
		return
	}
	_ = imp
	for path, pkg := range p.pkgs {
		if _, ok := loaded[path]; ok { // loaded
			continue
		}
		if pkg.flags == pkgFlagChanged {
			p.genGoPkg(path)
		}
	}
}

func getKeys(v map[string]none) []string {
	keys := make([]string, 0, len(v))
	for key := range v {
		keys = append(keys, key)
	}
	return keys
}

func (p *Runner) genGoPkg(pkgPath string) {
	// cl.NewPackage(pkgPath)
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
	modFile := mod.Modfile()
	modRoot, _ := filepath.Split(modFile)
	modPath := mod.Path()
	conf := &packages.Config{
		ModRoot:       modRoot,
		ModPath:       modPath,
		SupportedExts: mod.SupportedExts(),
	}
	pkgPaths, err := packages.List(conf, pattern...)
	if err != nil {
		return
	}
	modRoot = strings.TrimSuffix(modRoot, string(os.PathSeparator))

	var errs ErrorList
	p = &Runner{
		mod:     mod,
		modRoot: modRoot,
		modPath: modPath,
		modTime: modTime(modFile, &errs),
		pkgs:    make(map[string]*pkgInfo),
	}
	p.runners = map[string]*Runner{modRoot: p}
	for _, pkgPath := range pkgPaths {
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
	var depsHash string
	imps, err := p.mod.Imports(dir, false)
	if err != nil {
		*errs, pkg.flags = append(*errs, err), pkgFlagIll
		return
	}
	sort.Strings(imps)
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
	if p.sourceChanged(dir, ci, &depsHash) || pkg.hash != depsHash {
		pkg.flags |= pkgFlagChanged
	}
	return
}

func hashOf(b []byte) string {
	h := md5.Sum(b)
	return hex.EncodeToString(h[:])
}

func (p *Runner) genExternDeps(pkgPath string, errs *ErrorList) *pkgInfo {
	_, modVer, ok := p.mod.LookupExternPkg(pkgPath)
	if !ok {
		panic("TODO: externPkg not found")
	}
	modRoot, err := modfetch.ModCachePath(modVer)
	if err != nil {
		*errs = append(*errs, err)
		return &pkgInfo{path: pkgPath, flags: pkgFlagIll}
	}
	exr, ok := p.runners[modRoot]
	if !ok {
		mod, err := gopmod.Load(modRoot)
		if err != nil {
			*errs = append(*errs, err)
			return &pkgInfo{path: pkgPath, flags: pkgFlagIll}
		}
		exr = &Runner{
			mod:     mod,
			modRoot: modRoot,
			modPath: modVer.Path,
			modTime: modTime(mod.Modfile(), errs),
			pkgs:    make(map[string]*pkgInfo),
			runners: p.runners,
		}
		p.runners[modRoot] = exr
	}
	return exr.genDeps(pkgPath, pkgPath, errs)
}

func (p *Runner) sourceChanged(dir string, ci gopmod.ChangeInfo, depsHash *string) bool {
	f, err := os.Open(dir + "/gop_autogen.go")
	if err != nil {
		return true
	}
	fi, err := f.Stat()
	if err != nil || ci.ModTime.After(fi.ModTime()) {
		return true
	}
	//cinfo $SourceNum $DepsHash
	var sourceNum int
	if _, err = fmt.Fscanf(f, "//cinfo %d %s\n", &sourceNum, depsHash); err != nil {
		return true
	}
	return sourceNum != ci.SourceNum
}

// -----------------------------------------------------------------------------

type ErrorList []error

func (e ErrorList) Error() string {
	errStrs := make([]string, len(e))
	for i, err := range e {
		errStrs[i] = err.Error()
	}
	return strings.Join(errStrs, "\n")
}

// -----------------------------------------------------------------------------
