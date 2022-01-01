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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/goplus/gop/x/gopmod"
	"github.com/goplus/gop/x/mod/modfetch"
	"github.com/goplus/gox/packages"
)

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

	var errs ErrorList
	p = &Runner{
		mod:     mod,
		modRoot: modRoot,
		modPath: modPath,
		modTime: modTime(modFile, &errs),
		imports: make(map[string]none),
		pkgs:    make(map[string]*pkgInfo),
	}
	p.runners = map[string]*Runner{modPath: p}
	for _, pkgPath := range pkgPaths {
		p.gen(pkgPath, pkgPath, &errs)
	}
	if len(errs) > 0 {
		err = errs
	}
	return
}

const (
	pkgFlagIll = 1 << iota
	pkgFlagChanged
)

type pkgInfo struct {
	path    string
	runner  *Runner
	imports []string
	flags   int
}

type Runner struct {
	runners map[string]*Runner
	mod     *gopmod.Module
	modRoot string
	modPath string
	modTime time.Time
	imports map[string]none
	pkgs    map[string]*pkgInfo
}

func (p *Runner) gen(pkgPath, pkgPathBase string, errs *ErrorList) (pkg *pkgInfo) {
	switch p.mod.PkgType(pkgPath) {
	case gopmod.PkgtStandard:
		return nil
	case gopmod.PkgtLocal:
		pkgPath = filepath.ToSlash(filepath.Join(pkgPathBase, pkgPath))
	case gopmod.PkgtModule:
	case gopmod.PkgtExtern:
		return p.genExtern(pkgPath, errs)
	default:
		panic("TODO: invalid pkgPath")
	}
	pkg = &pkgInfo{
		path:   pkgPath,
		runner: p,
	}
	dir := filepath.Join(p.modRoot, pkgPath[len(p.modPath):])
	imps, err := p.mod.Imports(dir, false)
	if err != nil {
		*errs, pkg.flags = append(*errs, err), pkgFlagIll
		return
	}
	pkg.imports = imps
	if p.sourceChanged(dir) {
		pkg.flags = pkgFlagChanged
	}
	for _, imp := range imps {
		p.imports[imp] = none{}
		if t := p.gen(imp, pkgPath, errs); t != nil {
			pkg.flags |= t.flags
		}
	}
	return
}

func (p *Runner) genExtern(pkgPath string, errs *ErrorList) *pkgInfo {
	modVer, _, ok := p.mod.LookupExternPkg(pkgPath)
	if !ok {
		panic("TODO: externPkg not found")
	}
	modRoot, err := modfetch.ModCachePath(modVer)
	if err != nil {
		*errs = append(*errs, err)
		return &pkgInfo{path: pkgPath, flags: pkgFlagIll}
	}
	exr, ok := p.runners[modVer.Path]
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
			imports: make(map[string]none),
			runners: p.runners,
		}
		p.runners[modVer.Path] = exr
	}
	return exr.gen(pkgPath, pkgPath, errs)
}

func (p *Runner) sourceChanged(dir string) bool {
	ci, err := p.mod.ChangeInfo(dir)
	if err != nil {
		return true
	}
	f, err := os.Open(dir + "/gop_autogen.go")
	if err != nil {
		return true
	}
	fi, err := f.Stat()
	if err != nil || ci.ModTime.After(fi.ModTime()) {
		return true
	}
	//cinfo $SourceNum
	var sourceNum int
	if _, err = fmt.Fscanf(f, "//cinfo %d\n", &sourceNum); err != nil {
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
