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

package modfile

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// A File is the parsed, interpreted form of a gop.mod file.
type File struct {
	modfile.File
	Gop       *Gop
	Classfile *Classfile
	Register  []*Register
}

// A Module is the module statement.
type Module = modfile.Module

// A Go is the go statement.
type Go = modfile.Go

// A Require is a single require statement.
type Require = modfile.Require

// An Exclude is a single exclude statement.
type Exclude = modfile.Exclude

// A Replace is a single replace statement.
type Replace = modfile.Replace

// A Retract is a single retract statement.
type Retract = modfile.Retract

// A Gop is the gop statement.
type Gop = modfile.Go

// A Classfile is the classfile statement.
type Classfile struct {
	ProjExt  string   // ".gmx"
	WorkExt  string   // ".spx"
	PkgPaths []string // package paths of classfile
	Syntax   *Line
}

// A Register is the register statement.
type Register struct {
	ClassfileMod string // module path of classfile
	Syntax       *Line
}

// A VersionInterval represents a range of versions with upper and lower bounds.
// Intervals are closed: both bounds are included. When Low is equal to High,
// the interval may refer to a single version ('v1.2.3') or an interval
// ('[v1.2.3, v1.2.3]'); both have the same representation.
type VersionInterval = modfile.VersionInterval

type VersionFixer = modfile.VersionFixer

// Parse parses and returns a gop.mod file.
//
// file is the name of the file, used in positions and errors.
//
// data is the content of the file.
//
// fix is an optional function that canonicalizes module versions.
// If fix is nil, all module versions must be canonical (module.CanonicalVersion
// must return the same string).
func Parse(file string, data []byte, fix VersionFixer) (*File, error) {
	return parseToFile(file, data, fix, true)
}

// ParseLax is like Parse but ignores unknown statements.
// It is used when parsing gop.mod files other than the main module,
// under the theory that most statement types we add in the future will
// only apply in the main module, like exclude and replace,
// and so we get better gradual deployments if old go commands
// simply ignore those statements when found in gop.mod files
// in dependencies.
func ParseLax(file string, data []byte, fix VersionFixer) (*File, error) {
	return parseToFile(file, data, fix, false)
}

func parseToFile(file string, data []byte, fix VersionFixer, strict bool) (parsed *File, err error) {
	f, err := modfile.ParseLax(file, data, fix)
	if err != nil {
		return
	}
	parsed = &File{File: *f}

	var errs ErrorList
	var fs = f.Syntax
	for _, x := range fs.Stmt {
		switch x := x.(type) {
		case *Line:
			parsed.parseVerb(&errs, x.Token[0], x, x.Token[1:], strict)
		case *LineBlock:
			verb := x.Token[0]
			for _, line := range x.Line {
				parsed.parseVerb(&errs, verb, line, line.Token, strict)
			}
		}
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return
}

func (f *File) parseVerb(errs *ErrorList, verb string, line *Line, args []string, strict bool) {
	wrapModPathError := func(modPath string, err error) {
		*errs = append(*errs, Error{
			Filename: f.Syntax.Name,
			Pos:      line.Start,
			ModPath:  modPath,
			Verb:     verb,
			Err:      err,
		})
	}
	wrapError := func(err error) {
		*errs = append(*errs, Error{
			Filename: f.Syntax.Name,
			Pos:      line.Start,
			Err:      err,
		})
	}
	errorf := func(format string, args ...interface{}) {
		wrapError(fmt.Errorf(format, args...))
	}
	switch verb {
	case "require", "exclude", "module", "go", "retract":
	case "replace":
		arrow := 2
		if len(args) >= 2 && args[1] == "=>" {
			arrow = 1
		}
		if len(args) < arrow+2 || len(args) > arrow+3 || args[arrow] != "=>" {
			errorf("usage: %s module/path [v1.2.3] => other/module v1.4\n\t or %s module/path [v1.2.3] => ../local/directory", verb, verb)
			return
		}
		s, err := parseString(&args[0])
		if err != nil {
			errorf("invalid quoted string: %v", err)
			return
		}
		pathMajor, err := modulePathMajor(s)
		if err != nil {
			wrapModPathError(s, err)
			return
		}
		var v string
		if arrow == 2 {
			v, err = parseVersion(verb, s, &args[1])
			if err != nil {
				wrapError(err)
				return
			}
			if err := module.CheckPathMajor(v, pathMajor); err != nil {
				wrapModPathError(s, err)
				return
			}
		}
		ns, err := parseString(&args[arrow+1])
		if err != nil {
			errorf("invalid quoted string: %v", err)
			return
		}
		nv := ""
		if len(args) == arrow+2 {
			if !IsDirectoryPath(ns) {
				errorf("replacement module without version must be directory path (rooted or starting with ./ or ../)")
				return
			}
			if filepath.Separator == '/' && strings.Contains(ns, `\`) {
				errorf("replacement directory appears to be Windows path (on a non-windows system)")
				return
			}
		}
		if len(args) == arrow+3 {
			nv, err = parseVersion(verb, ns, &args[arrow+2])
			if err != nil {
				wrapError(err)
				return
			}
			if IsDirectoryPath(ns) {
				errorf("replacement module directory path %q cannot have version", ns)
				return
			}
		}
		f.Replace = append(f.Replace, &Replace{
			Old:    module.Version{Path: s, Version: v},
			New:    module.Version{Path: ns, Version: nv},
			Syntax: line,
		})
	case "gop":
		if f.Gop != nil {
			errorf("repeated go statement")
			return
		}
		if len(args) != 1 {
			errorf("go directive expects exactly one argument")
			return
		} else if !modfile.GoVersionRE.MatchString(args[0]) {
			errorf("invalid gop version '%s': must match format 1.23", args[0])
			return
		}
		f.Gop = &Gop{Syntax: line}
		f.Gop.Version = args[0]
	case "register":
		if len(args) != 1 {
			errorf("register directive expects exactly one argument")
			return
		}
		s, err := parseString(&args[0])
		if err != nil {
			errorf("invalid quoted string: %v", err)
			return
		}
		err = module.CheckPath(s)
		if err != nil {
			wrapError(err)
			return
		}
		f.Register = append(f.Register, &Register{
			ClassfileMod: s,
			Syntax:       line,
		})
	case "classfile":
		if f.Classfile != nil {
			errorf("repeated classfile statement")
			return
		}
		if len(args) < 3 {
			errorf("usage: classfile projExt workExt [classFilePkgPath ...]")
			return
		}
		projExt, err := parseExt(&args[0])
		if err != nil {
			wrapError(err)
			return
		}
		workExt, err := parseExt(&args[1])
		if err != nil {
			wrapError(err)
			return
		}
		pkgPaths, err := parseStrings(args[2:])
		if err != nil {
			errorf("invalid quoted string: %v", err)
			return
		}
		f.Classfile = &Classfile{
			ProjExt: projExt, WorkExt: workExt, PkgPaths: pkgPaths, Syntax: line,
		}
	default:
		if strict {
			errorf("unknown directive: %s", verb)
		}
	}
}

func parseVersion(verb string, path string, s *string) (string, error) {
	t, err := parseString(s)
	if err != nil {
		return "", &Error{
			Verb:    verb,
			ModPath: path,
			Err: &module.InvalidVersionError{
				Version: *s,
				Err:     err,
			},
		}
	}
	cv := module.CanonicalVersion(t)
	if cv == "" {
		return "", &Error{
			Verb:    verb,
			ModPath: path,
			Err: &module.InvalidVersionError{
				Version: t,
				Err:     errors.New("must be of the form v1.2.3"),
			},
		}
	}
	t = cv
	*s = t
	return *s, nil
}

func modulePathMajor(path string) (string, error) {
	_, major, ok := module.SplitPathVersion(path)
	if !ok {
		return "", fmt.Errorf("invalid module path")
	}
	return major, nil
}

// IsDirectoryPath reports whether the given path should be interpreted
// as a directory path. Just like on the go command line, relative paths
// and rooted paths are directory paths; the rest are module paths.
func IsDirectoryPath(ns string) bool {
	return modfile.IsDirectoryPath(ns)
}

// MustQuote reports whether s must be quoted in order to appear as
// a single token in a gop.mod line.
func MustQuote(s string) bool {
	return modfile.MustQuote(s)
}

// AutoQuote returns s or, if quoting is required for s to appear in a gop.mod,
// the quotation of s.
func AutoQuote(s string) string {
	return modfile.AutoQuote(s)
}

func parseString(s *string) (string, error) {
	t := *s
	if strings.HasPrefix(t, `"`) {
		var err error
		if t, err = strconv.Unquote(t); err != nil {
			return "", err
		}
	} else if strings.ContainsAny(t, "\"'`") {
		// Other quotes are reserved both for possible future expansion
		// and to avoid confusion. For example if someone types 'x'
		// we want that to be a syntax error and not a literal x in literal quotation marks.
		return "", fmt.Errorf("unquoted string cannot contain quote")
	}
	*s = AutoQuote(t)
	return t, nil
}

func parseStrings(args []string) (arr []string, err error) {
	arr = make([]string, len(args))
	for i := range args {
		if arr[i], err = parseString(&args[i]); err != nil {
			return
		}
	}
	return
}

func parseExt(s *string) (t string, err error) {
	t, err = parseString(s)
	if err != nil {
		goto failed
	}
	if len(t) > 1 && t[0] == '.' || t == "" {
		return
	}
	err = errors.New("invalid ext format")
failed:
	return "", &InvalidExtError{
		Ext: *s,
		Err: err,
	}
}

type InvalidExtError struct {
	Ext string
	Err error
}

func (e *InvalidExtError) Error() string {
	return fmt.Sprintf("ext %s invalid: %s", e.Ext, e.Err)
}

func (e *InvalidExtError) Unwrap() error { return e.Err }

type ErrorList = modfile.ErrorList
type Error = modfile.Error

// -----------------------------------------------------------------------------

const (
	directiveInvalid = iota
	directiveModule
	directiveGo
	directiveGop
	directiveClassfile
)

const (
	directiveLineBlock = 0x80 + iota
	directiveRegister
	directiveRequire
	directiveExclude
	directiveReplace
	directiveRetract
)

var directiveWeights = map[string]int{
	"module":    directiveModule,
	"go":        directiveGo,
	"gop":       directiveGop,
	"classfile": directiveClassfile,
	"register":  directiveRegister,
	"require":   directiveRequire,
	"exclude":   directiveExclude,
	"replace":   directiveReplace,
	"retract":   directiveRetract,
}

func getWeight(e Expr) int {
	if line, ok := e.(*Line); ok {
		return directiveWeights[line.Token[0]]
	}
	if w, ok := directiveWeights[e.(*LineBlock).Token[0]]; ok {
		return w
	}
	return directiveLineBlock
}

func updateLine(line *Line, tokens ...string) {
	if line.InBlock {
		tokens = tokens[1:]
	}
	line.Token = tokens
}

func addLine(x *FileSyntax, tokens ...string) *Line {
	new := &Line{Token: tokens}
	w := directiveWeights[tokens[0]]
	for i, e := range x.Stmt {
		w2 := getWeight(e)
		if w <= w2 {
			x.Stmt = append(x.Stmt, nil)
			copy(x.Stmt[i+1:], x.Stmt[i:])
			x.Stmt[i] = new
			return new
		}
	}
	x.Stmt = append(x.Stmt, new)
	return new
}

func (f *File) AddGopStmt(version string) error {
	if !modfile.GoVersionRE.MatchString(version) {
		return fmt.Errorf("invalid language version string %q", version)
	}
	if f.Gop == nil {
		f.Gop = &Gop{
			Version: version,
			Syntax:  addLine(f.Syntax, "gop", version),
		}
	} else {
		f.Gop.Version = version
		updateLine(f.Gop.Syntax, "gop", version)
	}
	return nil
}

func (f *File) AddRegister(modPath string) {
	for _, r := range f.Register {
		if r.ClassfileMod == modPath {
			return
		}
	}
	f.AddNewRegister(modPath)
}

func (f *File) AddNewRegister(modPath string) {
	line := addLine(f.Syntax, "register", AutoQuote(modPath))
	r := &Register{
		ClassfileMod: modPath,
		Syntax:       line,
	}
	f.Register = append(f.Register, r)
}

// -----------------------------------------------------------------------------

// markRemoved modifies line so that it (and its end-of-line comment, if any)
// will be dropped by (*FileSyntax).Cleanup.
func markRemoved(line *Line) {
	line.Token = nil
	line.Comments.Suffix = nil
}

func (f *File) DropAllRequire() {
	for _, r := range f.Require {
		markRemoved(r.Syntax)
	}
	f.Require = nil
}

func (f *File) DropAllReplace() {
	for _, r := range f.Replace {
		markRemoved(r.Syntax)
	}
	f.Replace = nil
}

// -----------------------------------------------------------------------------
