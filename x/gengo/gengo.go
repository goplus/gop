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
	"fmt"
	"io"
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
		gens:    make(map[string]none),
		imports: make(map[string]none),
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

type Runner struct {
	runners map[string]*Runner
	deps    map[string]*Runner
	mod     *gopmod.Module
	modRoot string
	modPath string
	modTime time.Time
	imports map[string]none
	gens    map[string]none
}

func (p *Runner) gen(pkgPath, pkgPathBase string, errs *ErrorList) {
	switch p.mod.PkgType(pkgPath) {
	case gopmod.PkgtStandard:
		return
	case gopmod.PkgtLocal:
		pkgPath = filepath.ToSlash(filepath.Join(pkgPathBase, pkgPath))
	case gopmod.PkgtModule:
	case gopmod.PkgtExtern:
		p.genExtern(pkgPath, errs)
		return
	default:
		panic("TODO: invalid pkgPath")
	}
	dir := filepath.Join(p.modRoot, pkgPath[len(p.modPath):])
	if p.notChanged(dir) {
		return
	}
	imps, err := p.mod.Imports(dir, false)
	if err != nil {
		*errs = append(*errs, err)
		return
	}
	for _, imp := range imps {
		p.gen(imp, pkgPath, errs)
		p.imports[imp] = none{}
	}
	p.gens[pkgPath] = none{}
}

func (p *Runner) genExtern(pkgPath string, errs *ErrorList) {
	modVer, _, ok := p.mod.LookupExternPkg(pkgPath)
	if !ok {
		panic("TODO: externPkg not found")
	}
	modRoot, err := modfetch.ModCachePath(modVer)
	if err != nil {
		*errs = append(*errs, err)
		return
	}
	exr, ok := p.runners[modVer.Path]
	if !ok {
		mod, err := gopmod.Load(modRoot)
		if err != nil {
			*errs = append(*errs, err)
			return
		}
		exr = &Runner{
			mod:     mod,
			modRoot: modRoot,
			modPath: modVer.Path,
			modTime: modTime(mod.Modfile(), errs),
			gens:    make(map[string]none),
			imports: make(map[string]none),
			runners: p.runners,
		}
		p.runners[modVer.Path] = exr
	}
	exr.gen(pkgPath, pkgPath, errs)
	if p.deps == nil {
		p.deps = make(map[string]*Runner)
	}
	p.deps[modVer.Path] = exr
}

func (p *Runner) notChanged(dir string) bool {
	ci, err := p.mod.ChangeInfo(dir)
	if err != nil {
		return false
	}
	f, err := os.Open(dir + "/gop_autogen.go")
	if err != nil {
		return false
	}
	fi, err := f.Stat()
	if err != nil || ci.ModTime.After(fi.ModTime()) {
		return false
	}
	//cinfo $SourceNum $ModuleModTime
	const cinfo = "//cinfo "
	const maxBufSize = 128
	var buf [maxBufSize]byte
	n, _ := io.ReadFull(f, buf[:])
	if pos := bytes.IndexByte(buf[:n], '\n'); pos > 0 {
		n = pos
	}
	s := string(buf[:n])
	if strings.HasPrefix(s, cinfo) {
		var sourceNum int
		var modTime int64
		if _, err = fmt.Sscan(s[len(cinfo):], &sourceNum, &modTime); err == nil {
			return sourceNum == ci.SourceNum && modTime == ci.ModTime.UnixNano()
		}
	}
	return false
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
