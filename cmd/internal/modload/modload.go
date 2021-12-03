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
	"log"
	"os"
	"path/filepath"

	gomodfile "golang.org/x/mod/modfile"
	"golang.org/x/mod/module"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/gengo"
	"github.com/goplus/gop/cmd/internal/search"
	"github.com/goplus/gop/env"
	"github.com/goplus/gop/x/mod/modfetch"
	"github.com/goplus/gop/x/mod/modfile"
)

var (
	modFile, classModFile *modfile.File

	initialized bool
	modRoot     string
	Target      module.Version
)

var ErrNoModRoot = errors.New("gop.mod file not found in current directory or any parent directory; see 'gop help modules'")

func findModuleRoot(dir string) (root string) {
	if dir == "" {
		panic("dir not set")
	}
	dir = filepath.Clean(dir)

	// Look for enclosing gop.mod or go.mod.
	for {
		if fi, err := os.Stat(filepath.Join(dir, "gop.mod")); err == nil && !fi.IsDir() {
			return dir
		}
		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
			return dir
		}
		d := filepath.Dir(dir)
		if d == dir {
			break
		}
		dir = d
	}
	return ""
}

func SetModRoot(dir string) {
	modRoot = dir
}

func getcwd() string {
	path, _ := os.Getwd()
	return path
}

var gopRoot = getcwd()

// HasModRoot reports whether a main module is present.
// HasModRoot may return false even if Enabled returns true: for example, 'get'
// does not require a main module.
func HasModRoot() bool {
	Init()
	return modRoot != ""
}

// GopModFilePath returns the effective path of the go.mod file. Normally, this
// "go.mod" in the directory returned by ModRoot, but the -modfile flag may
// change its location. ModFilePath calls base.Fatalf if there is no main
// module, even if -modfile is set.
func GopModFilePath() string {
	if !HasModRoot() {
		log.Fatalf("gop: %v", ErrNoModRoot)
	}
	return filepath.Join(modRoot, "gop.mod")
}

// GoModFilePath returns the effective path of the go.mod file. Normally, this
// "go.mod" in the directory returned by ModRoot, but the -modfile flag may
// change its location. ModFilePath calls base.Fatalf if there is no main
// module, even if -modfile is set.
func GoModFilePath() string {
	if !HasModRoot() {
		log.Fatalf("gop: %v", ErrNoModRoot)
	}
	return filepath.Join(modRoot, "go.mod")
}

// Init determines whether module mode is enabled, locates the root of the
// current module (if any), sets environment variables for Git subprocesses, and
// configures the cfg, codehost, load, modfetch, and search packages for use
// with modules.
func Init() {
	if initialized {
		return
	}
	initialized = true

	// Disable any prompting for passwords by Git.
	// Only has an effect for 2.3.0 or later, but avoiding
	// the prompt in earlier versions is just too hard.
	// If user has explicitly set GIT_TERMINAL_PROMPT=1, keep
	// prompting.
	// See golang.org/issue/9341 and golang.org/issue/12706.
	if os.Getenv("GIT_TERMINAL_PROMPT") == "" {
		os.Setenv("GIT_TERMINAL_PROMPT", "0")
	}

	// Disable any ssh connection pooling by Git.
	// If a Git subprocess forks a child into the background to cache a new connection,
	// that child keeps stdout/stderr open. After the Git subprocess exits,
	// os /exec expects to be able to read from the stdout/stderr pipe
	// until EOF to get all the data that the Git subprocess wrote before exiting.
	// The EOF doesn't come until the child exits too, because the child
	// is holding the write end of the pipe.
	// This is unfortunate, but it has come up at least twice
	// (see golang.org/issue/13453 and golang.org/issue/16104)
	// and confuses users when it does.
	// If the user has explicitly set GIT_SSH or GIT_SSH_COMMAND,
	// assume they know what they are doing and don't step on it.
	// But default to turning off ControlMaster.
	if os.Getenv("GIT_SSH") == "" && os.Getenv("GIT_SSH_COMMAND") == "" {
		os.Setenv("GIT_SSH_COMMAND", "ssh -o ControlMaster=no")
	}

	if modRoot != "" {
		// nothing to do
	} else {
		modRoot := findModuleRoot(gopRoot)
		if modRoot != "" {
			SetModRoot(modRoot)
		}
	}
}

func findModulePath(dir string) (string, error) {
	// Look for path in GOPATH.
	var badPathErr error
	for _, gpdir := range filepath.SplitList(getGoPath()) {
		if gpdir == "" {
			continue
		}
		if rel := search.InDir(dir, filepath.Join(gpdir, "src")); rel != "" && rel != "." {
			path := filepath.ToSlash(rel)
			return path, nil
		}
	}

	reason := "outside GOPATH, module path must be specified"
	if badPathErr != nil {
		// return a different error message if the module was in GOPATH, but
		// the module path determined above would be an invalid path.
		reason = fmt.Sprintf("bad module path inferred from directory in GOPATH: %v", badPathErr)
	}
	msg := `cannot determine module path for source directory %s (%s)

Example usage:
	'gop mod init example.com/m' to initialize a v0 or v1 module
	'gop mod init example.com/m/v2' to initialize a v2 module

Run 'gop help mod init' for more information.
`
	return "", fmt.Errorf(msg, dir, reason)
}

// CreateModFile initializes a new module by creating a go.mod file.
//
// If modPath is empty, CreateModFile will attempt to infer the path from the
// directory location within GOPATH.
//
// If a vendoring configuration file is present, CreateModFile will attempt to
// translate it to go.mod directives. The resulting build list may not be
// exactly the same as in the legacy configuration (for example, we can't get
// packages at multiple versions from the same module).
func CreateModFile(modPath string) {
	modRoot = gopRoot
	Init()
	modFilePath := GopModFilePath()
	if _, err := os.Stat(modFilePath); err == nil {
		log.Fatalf("gop: %s already exists", modFilePath)
	}

	if modPath == "" {
		var err error
		modPath, err = findModulePath(modRoot)
		if err != nil {
			log.Fatalf("gop: %v", err)
		}
	}

	fmt.Fprintf(os.Stderr, "gop: creating new gop.mod: module %s\n", modPath)
	modFile = new(modfile.File)
	modFile.AddModuleStmt(modPath)
	addGopStmt() // Add the gop directive before converted module requirements.
	WriteGopMod()
}

func Load() {
	LoadModFile()
	if modRoot == "" {
		return
	}
	SyncGoMod()
	if classModFile != nil && classModFile.Classfile != nil {
		gengo.RegisterPkgFlags(classModFile.Classfile.ProjExt, gengo.PkgFlagGmx)
		gengo.RegisterPkgFlags(classModFile.Classfile.WorkExt, gengo.PkgFlagSpx)
		cl.RegisterClassFileType(classModFile.Classfile.ProjExt,
			classModFile.Classfile.WorkExt, classModFile.Classfile.PkgPaths...)
	}
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

func fixGoVersion(fixed *bool) gomodfile.VersionFixer {
	return func(path, vers string) (resolved string, err error) {
		// do nothing
		return vers, nil
	}
}

// LoadModFile sets Target and, if there is a main module, parses the initial
// build list from its go.mod file.
//
// LoadModFile may make changes in memory, like adding a go directive and
// ensuring requirements are consistent. WriteGoMod should be called later to
// write changes out to disk or report errors in readonly mode.
//
// As a side-effect, LoadModFile may change cfg.BuildMod to "vendor" if
// -mod wasn't set explicitly and automatic vendoring should be enabled.
func LoadModFile() {
	Init()
	if modRoot == "" {
		return
	}
	// If gop.mod does not exist, then modroot does not exist,
	// and if go.mod exists then a copy of go.mod will be synchronized to gop.mod
	gopmod := GopModFilePath()
	gomod := GoModFilePath()
	if _, err := os.Stat(gopmod); os.IsNotExist(err) {
		if _, err := os.Stat(gomod); err == nil {
			SyncGopMod()
		}
		return
	}

	data, err := modfetch.Read(gopmod)
	if err != nil {
		log.Fatalf("gop: %v", err)
	}

	var fixed bool
	f, err := modfile.Parse(gopmod, data, fixVersion(&fixed))
	if err != nil {
		// Errors returned by modfile.Parse begin with file:line.
		log.Fatalf("gop: errors parsing gop.mod:\n%s\n", err)
	}
	modFile = f

	if f.Module == nil {
		// No module declaration. Must add module path.
		log.Fatalf("gop: no module declaration in gop.mod. To specify the module path:\n")
	}

	LoadClassFile()
}

// addGoStmt adds a gop directive to the gop.mod file if it does not already include one.
// The 'gop' version added, if any, is the latest version supported by this toolchain.
func addGopStmt() {
	if modFile.Gop != nil && modFile.Gop.Version != "" {
		return
	}
	version := env.MainVersion
	if !modfile.GopVersionRE.MatchString(version) {
		log.Fatalf("gop: unrecognized default version %q", version)
	}
	if err := modFile.AddGopStmt(version); err != nil {
		log.Fatalf("gop: internal error: %v", err)
	}
}

// WriteGopMod writes the current build list back to gop.mod.
func WriteGopMod() {
	// If we aren't in a module, we don't have anywhere to write a go.mod file.
	if modRoot == "" {
		return
	}
	addGopStmt()

	modFile.Cleanup()

	new, err := modFile.Format()
	if err != nil {
		log.Fatalf("gop: %v", err)
	}

	errNoChange := errors.New("no update needed")

	err = modfetch.Transform(GopModFilePath(), func(old []byte) ([]byte, error) {
		if bytes.Equal(old, new) {
			// The go.mod file is already equal to new, possibly as the result of some
			// other process.
			return nil, errNoChange
		}
		return new, nil
	})

	if err != nil && err != errNoChange {
		log.Fatalf("gop: updating gop.mod: %v", err)
	}
}

func getGoPath() string {
	return os.Getenv("GOPATH")
}

func SyncGoMod() {
	gomodPath := GoModFilePath()
	gomod := &gomodfile.File{}
	if _, err := os.Stat(gomodPath); err == nil {
		data, err := modfetch.Read(gomodPath)
		if err != nil {
			log.Fatalln(err)
		}
		var fixed bool
		gomod, err = gomodfile.Parse(gomodPath, data, fixGoVersion(&fixed))
		if err != nil {
			// Errors returned by modfile.Parse begin with file:line.
			log.Fatalf("gop: errors parsing gop.mod:\n%s\n", err)
		}
	}

	gomod.AddModuleStmt(modFile.Module.Mod.Path)
	if modFile.Go != nil {
		gomod.AddGoStmt(modFile.Go.Version)
	}

	for _, require := range modFile.Require {
		gomod.AddRequire(require.Mod.Path, require.Mod.Version)
	}

	for _, replace := range modFile.Replace {
		gomod.AddReplace(replace.Old.Path, replace.Old.Version, replace.New.Path, replace.New.Version)
	}

	for _, exclude := range modFile.Exclude {
		gomod.AddExclude(exclude.Mod.Path, exclude.Mod.Version)
	}

	for _, retract := range modFile.Retract {
		gomod.AddRetract(gomodfile.VersionInterval(retract.VersionInterval), retract.Rationale)
	}

	if classModFile != nil {
		for _, require := range classModFile.Require {
			gomod.AddRequire(require.Mod.Path, require.Mod.Version)
		}

		for _, replace := range classModFile.Replace {
			gomod.AddReplace(replace.Old.Path, replace.Old.Version, replace.New.Path, replace.New.Version)
		}

		for _, exclude := range classModFile.Exclude {
			gomod.AddExclude(exclude.Mod.Path, exclude.Mod.Version)
		}

		for _, retract := range classModFile.Retract {
			gomod.AddRetract(gomodfile.VersionInterval(retract.VersionInterval), retract.Rationale)
		}
	}

	gomod.Cleanup()

	new, err := gomod.Format()
	if err != nil {
		log.Fatalf("gop: %v", err)
	}

	errNoChange := errors.New("no update needed")
	err = modfetch.Transform(GoModFilePath(), func(old []byte) ([]byte, error) {
		if bytes.Equal(old, new) {
			// The go.mod file is already equal to new, possibly as the result of some
			// other process.
			return nil, errNoChange
		}
		return new, nil
	})

	if err != nil && err != errNoChange {
		log.Fatalf("gop: updating gop.mod: %v", err)
	}
}

func SyncGopMod() {
	gomodPath := GoModFilePath()
	gomod := &gomodfile.File{}
	if _, err := os.Stat(gomodPath); err == nil {
		data, err := modfetch.Read(gomodPath)
		if err != nil {
			log.Fatalln(err)
		}
		var fixed bool
		gomod, err = gomodfile.Parse(gomodPath, data, fixGoVersion(&fixed))
		if err != nil {
			// Errors returned by modfile.Parse begin with file:line.
			log.Fatalf("gop: errors parsing gop.mod:\n%s\n", err)
		}
	}

	if modFile == nil {
		modFile = &modfile.File{}
		modFile.AddModuleStmt(gomod.Module.Mod.Path)
	}
	if gomod.Go != nil {
		modFile.AddGoStmt(gomod.Go.Version)
	}

	for _, require := range gomod.Require {
		modFile.AddRequire(require.Mod.Path, require.Mod.Version)
	}

	for _, replace := range gomod.Replace {
		modFile.AddReplace(replace.Old.Path, replace.Old.Version, replace.New.Path, replace.New.Version)
	}

	for _, exclude := range gomod.Exclude {
		modFile.AddExclude(exclude.Mod.Path, exclude.Mod.Version)
	}

	for _, retract := range gomod.Retract {
		modFile.AddRetract(modfile.VersionInterval(retract.VersionInterval), retract.Rationale)
	}

	modFile.Cleanup()

	new, err := modFile.Format()
	if err != nil {
		log.Fatalf("gop: %v", err)
	}

	errNoChange := errors.New("no update needed")
	err = modfetch.Transform(GopModFilePath(), func(old []byte) ([]byte, error) {
		if bytes.Equal(old, new) {
			// The go.mod file is already equal to new, possibly as the result of some
			// other process.
			return nil, errNoChange
		}
		return new, nil
	})

	if err != nil && err != errNoChange {
		log.Fatalf("gop: updating gop.mod: %v", err)
	}
}
