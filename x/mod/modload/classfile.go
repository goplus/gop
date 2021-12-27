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

/*
import (
	"log"
	"path/filepath"

	"golang.org/x/mod/module"

	"github.com/goplus/gop/x/mod/modfetch"
	"github.com/goplus/gop/x/mod/modfile"
)

func LoadClassFile() {
	if modFile.Register == nil {
		return
	}

	var dir string
	var err error
	var claassMod module.Version

	for _, require := range modFile.Require {
		if require.Mod.Path == modFile.Register[0].ClassfileMod {
			claassMod = require.Mod
		}
	}
	for _, replace := range modFile.Replace {
		if replace.Old.Path == modFile.Register[0].ClassfileMod {
			claassMod = replace.New
		}
	}

	if claassMod.Version != "" {
		dir, err = modfetch.Download(claassMod)
		if err != nil {
			log.Fatalf("gop: download classsfile module error %v", err)
		}
	} else {
		dir = claassMod.Path
	}

	if dir == "" {
		log.Fatalf("gop: can't find classfile path in require statment")
	}

	gopmod := filepath.Join(modRoot, dir, "gop.mod")
	data, err := modfetch.Read(gopmod)
	if err != nil {
		log.Fatalf("gop: %v", err)
	}

	var fixed bool
	f, err := modfile.Parse(gopmod, data, fixVersion(&fixed))
	if err != nil {
		// Errors returned by modfile.Parse begin with file:line.
		log.Fatalf("go: errors parsing go.mod:\n%s\n", err)
	}
	classModFile = f
}
*/

/*
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
*/

/*
// addGoStmt adds a gop directive to the gop.mod file if it does not already include one.
// The 'gop' version added, if any, is the latest version supported by this toolchain.
func addGopStmt() {
	if modFile.Gop != nil && modFile.Gop.Version != "" {
		return
	}
	version := env.MainVersion
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
*/

/*
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
*/
