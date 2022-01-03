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

package modload

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	gomodfile "golang.org/x/mod/modfile"
	"golang.org/x/mod/module"

	"github.com/goplus/gop/env"
	"github.com/goplus/gop/x/mod/modfile"
)

var (
	ErrNoModDecl = errors.New("no module declaration in gop.mod (or go.mod)")
	ErrNoModRoot = errors.New("gop.mod or go.mod file not found in current directory or any parent directory")
)

type Module struct {
	*modfile.File
}

// Modfile returns absolute path of the module file (gop.mod or go.mod).
func (p Module) Modfile() string {
	return p.Syntax.Name
}

// Path returns the module path.
func (p Module) Path() string {
	if mod := p.Module; mod != nil {
		return mod.Mod.Path
	}
	return ""
}

// DepMods returns all depended modules.
// If a depended module path is replace to be a local path, it will be canonical to an absolute path.
func (p Module) DepMods() map[string]module.Version {
	vers := make(map[string]module.Version)
	for _, r := range p.Require {
		if r.Mod.Path != "" {
			vers[r.Mod.Path] = r.Mod
		}
	}
	for _, r := range p.Replace {
		if r.Old.Path != "" {
			real := r.New
			if real.Version == "" {
				if strings.HasPrefix(real.Path, ".") {
					dir, _ := filepath.Split(p.Modfile())
					real.Path = dir + real.Path
				}
				if a, err := filepath.Abs(real.Path); err == nil {
					real.Path = a
				}
			}
			vers[r.Old.Path] = real
		}
	}
	return vers
}

// Create creates a new module in `dir`.
// You should call `Save` manually to save this module.
func Create(dir string, modPath, goVer, gopVer string) (p Module, err error) {
	gopmod, err := filepath.Abs(filepath.Join(dir, "gop.mod"))
	if err != nil {
		return
	}
	if _, err := os.Stat(gopmod); err == nil {
		return Module{}, fmt.Errorf("gop: %s already exists", gopmod)
	}
	mod := new(modfile.File)
	mod.AddModuleStmt(modPath)
	mod.AddGoStmt(goVer)
	mod.AddGopStmt(gopVer)
	mod.Syntax.Name = gopmod
	return Module{File: mod}, nil
}

// fixVersion returns a modfile.VersionFixer implemented using the Query function.
//
// It resolves commit hashes and branch names to versions,
// canonicalizes versions that appeared in early vgo drafts,
// and does nothing for versions that already appear to be canonical.
//
// The VersionFixer sets 'fixed' if it ever returns a non-canonical version.
func fixVersion(fixed *bool) modfile.VersionFixer {
	return func(path, vers string) (resolved string, err error) {
		// do nothing
		return vers, nil
	}
}

// Load loads a module from `dir`.
func Load(dir string) (p Module, err error) {
	gopmod, err := env.GOPMOD(dir)
	if err != nil {
		return
	}

	data, err := os.ReadFile(gopmod)
	if err != nil {
		return
	}

	var fixed bool
	f, err := modfile.Parse(gopmod, data, fixVersion(&fixed))
	if err != nil {
		// Errors returned by modfile.Parse begin with file:line.
		return
	}
	if f.Module == nil {
		// No module declaration. Must add module path.
		return Module{}, ErrNoModDecl
	}
	return Module{File: f}, nil
}

func loadGoMod(gomod string) (f *gomodfile.File, err error) {
	data, err := os.ReadFile(gomod)
	if err != nil {
		return
	}
	var fixed bool
	return gomodfile.Parse(gomod, data, fixVersion(&fixed))
}

// -----------------------------------------------------------------------------

// Save saves all changes of this module.
func (p Module) Save() error {
	modfile := p.Modfile()
	data, err := p.Format()
	if err == nil {
		err = os.WriteFile(modfile, data, 0644)
	}
	return err
}

const (
	gopMod = "github.com/goplus/gop"
)

func (p Module) Tidy() (err error) {
	gopmod := p.Modfile()
	dir, file := filepath.Split(gopmod)
	gomod := dir + "go.mod"
	isGopMod := file != "go.mod"
	if isGopMod {
		if err = p.saveGoMod(gomod); err != nil {
			return
		}
	}
	if err = tidy(dir); err != nil {
		return
	}
	if isGopMod {
		gof, loadErr := loadGoMod(gomod)
		if loadErr != nil {
			return loadErr
		}
		p.DropAllRequire()
		p.DropAllReplace()
		for _, r := range gof.Require {
			switch r.Mod.Path {
			case gopMod, "":
				continue
			}
			p.AddRequire(r.Mod.Path, r.Mod.Version)
		}
		for _, r := range gof.Replace {
			switch r.Old.Path {
			case gopMod, "":
				continue
			}
			p.AddReplace(r.Old.Path, r.Old.Version, r.New.Path, r.New.Version)
		}
		err = p.Save()
	}
	return
}

func tidy(modRoot string) (err error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = modRoot
	err = cmd.Run()
	if err != nil {
		err = &ExecCmdError{Err: err, Stderr: stderr.Bytes()}
	}
	return
}

// UpdateGoMod updates the go.mod file.
func (p Module) UpdateGoMod(checkChanged bool) error {
	gopmod := p.Modfile()
	dir, file := filepath.Split(gopmod)
	if file == "go.mod" {
		return nil
	}
	gomod := dir + "go.mod"
	if checkChanged && notChanged(gomod, gopmod) {
		return nil
	}
	return p.saveGoMod(gomod)
}

func (p Module) saveGoMod(gomod string) error {
	gof := p.convToGoMod()
	data, err := gof.Format()
	if err == nil {
		err = os.WriteFile(gomod, data, 0644)
	}
	return err
}

func (p Module) convToGoMod() *gomodfile.File {
	copy := p.File.File
	copy.Syntax = cloneGoFileSyntax(copy.Syntax)
	addRequireIfNotExist(&copy, gopMod, env.Version())
	addReplaceIfNotExist(&copy, gopMod, "", env.GOPROOT(), "")
	return &copy
}

func addRequireIfNotExist(f *gomodfile.File, path, vers string) {
	for _, r := range f.Require {
		if r.Mod.Path == path {
			return
		}
	}
	f.AddNewRequire(path, vers, false)
}

func addReplaceIfNotExist(f *gomodfile.File, oldPath, oldVers, newPath, newVers string) {
	for _, r := range f.Replace {
		if r.Old.Path == oldPath && (oldVers == "" || r.Old.Version == oldVers) {
			return
		}
	}
	f.AddReplace(oldPath, oldVers, newPath, newVers)
}

func notChanged(target, src string) bool {
	fiTarget, err := os.Stat(target)
	if err != nil {
		return false
	}
	fiSrc, err := os.Stat(src)
	if err != nil {
		return false
	}
	return fiTarget.ModTime().After(fiSrc.ModTime())
}

// -----------------------------------------------------------------------------

func cloneGoFileSyntax(syn *modfile.FileSyntax) *modfile.FileSyntax {
	stmt := make([]modfile.Expr, 0, len(syn.Stmt))
	for _, e := range syn.Stmt {
		if isGopExpr(e) {
			continue
		}
		stmt = append(stmt, cloneExpr(e))
	}
	return &modfile.FileSyntax{
		Name:     syn.Name,
		Comments: syn.Comments,
		Stmt:     stmt,
	}
}

func cloneExpr(e modfile.Expr) modfile.Expr {
	if v, ok := e.(*modfile.LineBlock); ok {
		copy := *v
		return &copy
	}
	return e
}

func isGopExpr(e modfile.Expr) bool {
	switch verb := getVerb(e); verb {
	case "gop", "register", "classfile":
		return true
	}
	return false
}

func getVerb(e modfile.Expr) string {
	if line, ok := e.(*modfile.Line); ok {
		return line.Token[0]
	}
	return e.(*modfile.LineBlock).Token[0]
}

// -----------------------------------------------------------------------------

type ExecCmdError struct {
	Err    error
	Stderr []byte
}

func (p *ExecCmdError) Error() string {
	if e := p.Stderr; e != nil {
		return string(e)
	}
	return p.Err.Error()
}

// -----------------------------------------------------------------------------
