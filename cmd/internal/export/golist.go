/*
 Copyright 2020 The GoPlus Authors (goplus.org)

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

package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type jsonPackage struct {
	Dir           string      // directory containing package sources
	ImportPath    string      // import path of package in dir
	ImportComment string      // path in import comment on package statement
	Name          string      // package name
	Doc           string      // package documentation string
	Target        string      // install path
	Shlib         string      // the shared library that contains this package (only set when -linkshared)
	Goroot        bool        // is this package in the Go root?
	Standard      bool        // is this package part of the standard Go library?
	Stale         bool        // would 'go install' do anything for this package?
	StaleReason   string      // explanation for Stale==true
	Root          string      // Go root or Go path dir containing this package
	ConflictDir   string      // this directory shadows Dir in $GOPATH
	BinaryOnly    bool        // binary-only package (no longer supported)
	ForTest       string      // package is only for use in named test
	Export        string      // file containing export data (when using -export)
	Module        *jsonModule // info about package's containing module, if any (can be nil)
	Match         []string    // command-line patterns matching this package
	DepOnly       bool        // package is only a dependency, not explicitly listed

	// Source files
	GoFiles         []string // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles        []string // .go source files that import "C"
	CompiledGoFiles []string // .go files presented to compiler (when using -compiled)
	IgnoredGoFiles  []string // .go source files ignored due to build constraints
	CFiles          []string // .c source files
	CXXFiles        []string // .cc, .cxx and .cpp source files
	MFiles          []string // .m source files
	HFiles          []string // .h, .hh, .hpp and .hxx source files
	FFiles          []string // .f, .F, .for and .f90 Fortran source files
	SFiles          []string // .s source files
	SwigFiles       []string // .swig files
	SwigCXXFiles    []string // .swigcxx files
	SysoFiles       []string // .syso object files to add to archive
	TestGoFiles     []string // _test.go files in package
	XTestGoFiles    []string // _test.go files outside package

	// Cgo directives
	CgoCFLAGS    []string // cgo: flags for C compiler
	CgoCPPFLAGS  []string // cgo: flags for C preprocessor
	CgoCXXFLAGS  []string // cgo: flags for C++ compiler
	CgoFFLAGS    []string // cgo: flags for Fortran compiler
	CgoLDFLAGS   []string // cgo: flags for linker
	CgoPkgConfig []string // cgo: pkg-config names

	// Dependency information
	Imports      []string          // import paths used by this package
	ImportMap    map[string]string // map from source import to ImportPath (identity entries omitted)
	Deps         []string          // all (recursively) imported dependencies
	TestImports  []string          // imports from TestGoFiles
	XTestImports []string          // imports from XTestGoFiles

	// Error information
	Incomplete bool                // this package or a dependency has an error
	Error      *jsonPackageError   // error loading package
	DepsErrors []*jsonPackageError // errors loading dependencies
}

type jsonPackageError struct {
	ImportStack []string // shortest path from package named on command line to this one
	Pos         string   // position of error (if present, file:line:col)
	Err         string   // the error itself
}

type jsonModule struct {
	Path      string           // module path
	Version   string           // module version
	Versions  []string         // available module versions (with -versions)
	Replace   *jsonModule      // replaced by this module
	Time      *time.Time       // time version was created
	Update    *jsonModule      // available update, if any (with -u)
	Main      bool             // is this the main module?
	Indirect  bool             // is this module only an indirect dependency of main module?
	Dir       string           // directory holding files for this module, if any
	GoMod     string           // path to go.mod file used when loading this module, if any
	GoVersion string           // go version used in module
	Error     *jsonModuleError // error loading module
}

type jsonModuleError struct {
	Err string // the error itself
}

func checkGoPkgList(pkgPath string, workDir string, allpkgs bool) (pkgs []*jsonPackage, err error) {
	spkg := pkgPath
	if allpkgs && !strings.HasSuffix(pkgPath, "/...") {
		spkg += "/..."
	}
	cmd := exec.Command(gobin, "list", "-json", spkg)
	cmd.Dir = workDir
	cmd.Stderr = os.Stderr
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	for dec := json.NewDecoder(buf); dec.More(); {
		pkg := &jsonPackage{}
		if err := dec.Decode(&pkg); err != nil {
			return nil, fmt.Errorf("JSON decoding failed: %v", err)
		}
		pkgs = append(pkgs, pkg)
	}
	return
}

func checkGoPkg(srcDir string) (string, error) {
	cmd := exec.Command(gobin, "list")
	cmd.Dir = srcDir
	data, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
