// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Api computes the exported API of a set of Go packages.

// modify 2013-2016 visualfc

package goapi

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	apiVerbose      bool   = false
	ApiAllmethods   bool   = true
	ApiAlldecls     bool   = false
	apiShowpos      bool   = false
	ApiSeparate     string = ", "
	ApiImportParser bool   = true
	ApiDefaultCtx   bool   = true
	ApiCustomCtx    string = ""
)

//func init() {
//	Command.Flag.BoolVar(&apiVerbose, "v", false, "verbose debugging")
//	Command.Flag.BoolVar(&apiAllmethods, "e", true, "extract for all embedded methods")
//	Command.Flag.BoolVar(&apiAlldecls, "a", false, "extract for all declarations")
//	Command.Flag.BoolVar(&apiShowpos, "pos", false, "addition token position")
//	Command.Flag.StringVar(&apiSeparate, "sep", ", ", "setup separators")
//	Command.Flag.BoolVar(&apiImportParser, "dep", true, "parser package imports")
//	Command.Flag.BoolVar(&apiDefaultCtx, "default_ctx", true, "extract for default context")
//	Command.Flag.StringVar(&apiCustomCtx, "custom_ctx", "", "optional comma-separated list of <goos>-<goarch>[-cgo] to override default contexts.")
//	Command.Flag.StringVar(&apiLookupInfo, "cursor_info", "", "lookup cursor node info\"file.go:pos\"")
//	Command.Flag.BoolVar(&apiLookupStdin, "cursor_std", false, "cursor_info use stdin")
//	Command.Flag.StringVar(&apiOutput, "o", "", "output file")
//}

func LookupApi(pkgs ...string) ([]string, error) {
	if len(pkgs) == 0 {
		return nil, errors.New("empty pkg list")
	}
	var pkglist []string
	if pkgs[0] == "std" || pkgs[0] == "all" {
		out, err := exec.Command("go", "list", "-e", pkgs[0]).Output()
		if err != nil {
			log.Fatal(err)
		}
		pkglist = strings.Fields(string(out))
	} else {
		pkglist = pkgs
	}
	if ApiCustomCtx != "" {
		ApiDefaultCtx = false
		setCustomContexts()
	}

	var features []string
	w := NewWalker()
	w.sep = ApiSeparate

	if ApiDefaultCtx {
		w.context = &build.Default

		for _, pkg := range pkglist {
			w.wantedPkg[pkg] = true
		}

		for _, pkg := range pkglist {
			w.WalkPackage(pkg)
		}
		features = w.Features("")
	} else {
		for _, c := range contexts {
			c.Compiler = build.Default.Compiler
		}

		for _, pkg := range pkglist {
			w.wantedPkg[pkg] = true
		}

		var featureCtx = make(map[string]map[string]bool) // feature -> context name -> true
		for _, context := range contexts {
			w.context = context
			w.ctxName = contextName(w.context) + ":"

			for _, pkg := range pkglist {
				w.WalkPackage(pkg)
			}
		}

		for pkg, p := range w.packageMap {
			if w.wantedPkg[p.name] {
				pos := strings.Index(pkg, ":")
				if pos == -1 {
					continue
				}
				ctxName := pkg[:pos]
				for _, f := range p.Features() {
					if featureCtx[f] == nil {
						featureCtx[f] = make(map[string]bool)
					}
					featureCtx[f][ctxName] = true
				}
			}
		}

		for f, cmap := range featureCtx {
			if len(cmap) == len(contexts) {
				features = append(features, f)
				continue
			}
			comma := strings.Index(f, ",")
			for cname := range cmap {
				f2 := fmt.Sprintf("%s (%s)%s", f[:comma], cname, f[comma:])
				features = append(features, f2)
			}
		}
		sort.Strings(features)
	}

	return features, nil
}

type CursorInfo struct {
	pkg  string
	file string
	pos  token.Pos
	src  []byte
	std  bool
	info *TypeInfo
}

// contexts are the default contexts which are scanned, unless
// overridden by the -contexts flag.
var contexts = []*build.Context{
	{GOOS: "linux", GOARCH: "386", CgoEnabled: true},
	{GOOS: "linux", GOARCH: "386"},
	{GOOS: "linux", GOARCH: "amd64", CgoEnabled: true},
	{GOOS: "linux", GOARCH: "amd64"},
	{GOOS: "linux", GOARCH: "arm", CgoEnabled: true},
	{GOOS: "linux", GOARCH: "arm"},
	{GOOS: "darwin", GOARCH: "386", CgoEnabled: true},
	{GOOS: "darwin", GOARCH: "386"},
	{GOOS: "darwin", GOARCH: "amd64", CgoEnabled: true},
	{GOOS: "darwin", GOARCH: "amd64"},
	{GOOS: "windows", GOARCH: "amd64"},
	{GOOS: "windows", GOARCH: "386"},
	{GOOS: "freebsd", GOARCH: "386", CgoEnabled: true},
	{GOOS: "freebsd", GOARCH: "386"},
	{GOOS: "freebsd", GOARCH: "amd64", CgoEnabled: true},
	{GOOS: "freebsd", GOARCH: "amd64"},
	{GOOS: "freebsd", GOARCH: "arm", CgoEnabled: true},
	{GOOS: "freebsd", GOARCH: "arm"},
	{GOOS: "netbsd", GOARCH: "386", CgoEnabled: true},
	{GOOS: "netbsd", GOARCH: "386"},
	{GOOS: "netbsd", GOARCH: "amd64", CgoEnabled: true},
	{GOOS: "netbsd", GOARCH: "amd64"},
	{GOOS: "netbsd", GOARCH: "arm", CgoEnabled: true},
	{GOOS: "netbsd", GOARCH: "arm"},
	{GOOS: "openbsd", GOARCH: "386", CgoEnabled: true},
	{GOOS: "openbsd", GOARCH: "386"},
	{GOOS: "openbsd", GOARCH: "amd64", CgoEnabled: true},
	{GOOS: "openbsd", GOARCH: "amd64"},
}

func contextName(c *build.Context) string {
	s := c.GOOS + "-" + c.GOARCH
	if c.CgoEnabled {
		return s + "-cgo"
	}
	return s
}

func osArchName(c *build.Context) string {
	return c.GOOS + "-" + c.GOARCH
}

func parseContext(c string) *build.Context {
	parts := strings.Split(c, "-")
	if len(parts) < 2 {
		log.Fatalf("bad context: %q", c)
	}
	bc := &build.Context{
		GOOS:   parts[0],
		GOARCH: parts[1],
	}
	if len(parts) == 3 {
		if parts[2] == "cgo" {
			bc.CgoEnabled = true
		} else {
			log.Fatalf("bad context: %q", c)
		}
	}
	return bc
}

func setCustomContexts() {
	contexts = []*build.Context{}
	for _, c := range strings.Split(ApiCustomCtx, ",") {
		contexts = append(contexts, parseContext(c))
	}
}

func set(items []string) map[string]bool {
	s := make(map[string]bool)
	for _, v := range items {
		s[v] = true
	}
	return s
}

var spaceParensRx = regexp.MustCompile(` \(\S+?\)`)

func featureWithoutContext(f string) string {
	if !strings.Contains(f, "(") {
		return f
	}
	return spaceParensRx.ReplaceAllString(f, "")
}

func compareAPI(w io.Writer, features, required, optional, exception []string, allowNew bool) (ok bool) {
	ok = true

	optionalSet := set(optional)
	exceptionSet := set(exception)
	featureSet := set(features)

	sort.Strings(features)
	sort.Strings(required)

	take := func(sl *[]string) string {
		s := (*sl)[0]
		*sl = (*sl)[1:]
		return s
	}

	for len(required) > 0 || len(features) > 0 {
		switch {
		case len(features) == 0 || (len(required) > 0 && required[0] < features[0]):
			feature := take(&required)
			if exceptionSet[feature] {
				fmt.Fprintf(w, "~%s\n", feature)
			} else if featureSet[featureWithoutContext(feature)] {
				// okay.
			} else {
				fmt.Fprintf(w, "-%s\n", feature)
				ok = false // broke compatibility
			}
		case len(required) == 0 || (len(features) > 0 && required[0] > features[0]):
			newFeature := take(&features)
			if optionalSet[newFeature] {
				// Known added feature to the upcoming release.
				// Delete it from the map so we can detect any upcoming features
				// which were never seen.  (so we can clean up the nextFile)
				delete(optionalSet, newFeature)
			} else {
				fmt.Fprintf(w, "+%s\n", newFeature)
				if !allowNew {
					ok = false // we're in lock-down mode for next release
				}
			}
		default:
			take(&required)
			take(&features)
		}
	}

	// In next file, but not in API.
	var missing []string
	for feature := range optionalSet {
		missing = append(missing, feature)
	}
	sort.Strings(missing)
	for _, feature := range missing {
		fmt.Fprintf(w, "±%s\n", feature)
	}
	return
}

func fileFeatures(filename string) []string {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", filename, err)
	}
	text := strings.TrimSpace(string(bs))
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}

func isExtract(name string) bool {
	if ApiAlldecls {
		return true
	}
	return ast.IsExported(name)
}

// pkgSymbol represents a symbol in a package
type pkgSymbol struct {
	pkg    string // "net/http"
	symbol string // "RoundTripper"
}

//expression kind
type Kind int

const (
	KindBuiltin Kind = iota
	KindPackage
	KindImport
	KindVar
	KindConst
	KindInterface
	KindParam
	KindStruct
	KindMethod
	KindField
	KindType
	KindFunc
	KindChan
	KindArray
	KindMap
	KindSlice
	KindLabel
	KindBranch
)

func (k Kind) String() string {
	switch k {
	case KindBuiltin:
		return "builtin"
	case KindPackage:
		return "package"
	case KindImport:
		return "import"
	case KindVar:
		return "var"
	case KindConst:
		return "const"
	case KindParam:
		return "param"
	case KindInterface:
		return "interface"
	case KindStruct:
		return "struct"
	case KindMethod:
		return "method"
	case KindField:
		return "field"
	case KindType:
		return "type"
	case KindFunc:
		return "func"
	case KindChan:
		return "chan"
	case KindMap:
		return "map"
	case KindArray:
		return "array"
	case KindSlice:
		return "slice"
	case KindLabel:
		return "label"
	case KindBranch:
		return "branch"
	}
	return fmt.Sprint("unknown-kind")
}

//expression type
type TypeInfo struct {
	Kind Kind
	Name string
	Type string
	X    ast.Expr
	T    ast.Expr
}

type ExprType struct {
	X ast.Expr
	T string
}

type Package struct {
	dpkg             *doc.Package
	apkg             *ast.Package
	interfaceMethods map[string]([]typeMethod)
	interfaces       map[string]*ast.InterfaceType //interface
	structs          map[string]*ast.StructType    //struct
	types            map[string]ast.Expr           //type
	functions        map[string]typeMethod         //function
	consts           map[string]*ExprType          //const => type
	vars             map[string]*ExprType          //var => type
	name             string
	dir              string
	sep              string
	deps             []string
	features         map[string](token.Pos) // set
}

func NewPackage() *Package {
	return &Package{
		interfaceMethods: make(map[string]([]typeMethod)),
		interfaces:       make(map[string]*ast.InterfaceType),
		structs:          make(map[string]*ast.StructType),
		types:            make(map[string]ast.Expr),
		functions:        make(map[string]typeMethod),
		consts:           make(map[string]*ExprType),
		vars:             make(map[string]*ExprType),
		features:         make(map[string](token.Pos)),
		sep:              ", ",
	}
}

func (p *Package) Features() (fs []string) {
	for f, ps := range p.features {
		if apiShowpos {
			fs = append(fs, f+p.sep+strconv.Itoa(int(ps)))
		} else {
			fs = append(fs, f)
		}
	}
	sort.Strings(fs)
	return
}

func (p *Package) findType(name string) ast.Expr {
	for k, v := range p.interfaces {
		if k == name {
			return v
		}
	}
	for k, v := range p.structs {
		if k == name {
			return v
		}
	}
	for k, v := range p.types {
		if k == name {
			return v
		}
	}
	return nil
}

func funcRetType(ft *ast.FuncType, index int) ast.Expr {
	if ft.Results != nil {
		pos := 0
		for _, fi := range ft.Results.List {
			if fi.Names == nil {
				if pos == index {
					return fi.Type
				}
				pos++
			} else {
				for _ = range fi.Names {
					if pos == index {
						return fi.Type
					}
					pos++
				}
			}
		}
	}
	return nil
}

func findFunction(funcs []*doc.Func, name string) (*ast.Ident, *ast.FuncType) {
	for _, f := range funcs {
		if f.Name == name {
			return &ast.Ident{Name: name, NamePos: f.Decl.Pos()}, f.Decl.Type
		}
	}
	return nil, nil
}

func (p *Package) findSelectorType(name string) ast.Expr {
	if t, ok := p.vars[name]; ok {
		return &ast.Ident{
			NamePos: t.X.Pos(),
			Name:    t.T,
		}
	}
	if t, ok := p.consts[name]; ok {
		return &ast.Ident{
			NamePos: t.X.Pos(),
			Name:    t.T,
		}
	}
	if t, ok := p.functions[name]; ok {
		return t.ft
	}
	for k, v := range p.structs {
		if k == name {
			return &ast.Ident{
				NamePos: v.Pos(),
				Name:    name,
			}
		}
	}
	for k, v := range p.interfaces {
		if k == name {
			return &ast.Ident{
				NamePos: v.Pos(),
				Name:    name,
			}
		}
	}
	for k, v := range p.types {
		if k == name {
			return v
		}
	}
	return nil
}

func (p *Package) findCallFunc(name string) ast.Expr {
	if fn, ok := p.functions[name]; ok {
		return fn.ft
	}
	if s, ok := p.structs[name]; ok {
		return s
	}
	if t, ok := p.types[name]; ok {
		return t
	}
	if v, ok := p.vars[name]; ok {
		if strings.HasPrefix(v.T, "func(") {
			e, err := parser.ParseExpr(v.T + "{}")
			if err == nil {
				return e
			}
		}
	}
	return nil
}

func (p *Package) findCallType(name string, index int) ast.Expr {
	if fn, ok := p.functions[name]; ok {
		return funcRetType(fn.ft, index)
	}
	if s, ok := p.structs[name]; ok {
		return &ast.Ident{
			NamePos: s.Pos(),
			Name:    name,
		}
	}
	if t, ok := p.types[name]; ok {
		return &ast.Ident{
			NamePos: t.Pos(),
			Name:    name,
		}
	}
	return nil
}

func (p *Package) findMethod(typ, name string) (*ast.Ident, *ast.FuncType) {
	if t, ok := p.interfaces[typ]; ok && t.Methods != nil {
		for _, fd := range t.Methods.List {
			switch ft := fd.Type.(type) {
			case *ast.FuncType:
				for _, ident := range fd.Names {
					if ident.Name == name {
						return ident, ft
					}
				}
			}
		}
	}
	for k, v := range p.interfaceMethods {
		if k == typ {
			for _, m := range v {
				if m.name == name {
					return &ast.Ident{Name: name, NamePos: m.pos}, m.ft
				}
			}
		}
	}
	if p.dpkg == nil {
		return nil, nil
	}
	for _, t := range p.dpkg.Types {
		if t.Name == typ {
			return findFunction(t.Methods, name)
		}
	}
	return nil, nil
}

type Walker struct {
	context *build.Context
	fset    *token.FileSet
	scope   []string
	//	features        map[string](token.Pos) // set
	lastConstType   string
	curPackageName  string
	sep             string
	ctxName         string
	curPackage      *Package
	constDep        map[string]*ExprType // key's const identifier has type of future value const identifier
	packageState    map[string]loadState
	packageMap      map[string]*Package
	interfaces      map[pkgSymbol]*ast.InterfaceType
	selectorFullPkg map[string]string // "http" => "net/http", updated by imports
	wantedPkg       map[string]bool   // packages requested on the command line
	cursorInfo      *CursorInfo
	localvar        map[string]*ExprType
}

func NewWalker() *Walker {
	return &Walker{
		fset: token.NewFileSet(),
		//		features:        make(map[string]token.Pos),
		packageState:    make(map[string]loadState),
		interfaces:      make(map[pkgSymbol]*ast.InterfaceType),
		packageMap:      make(map[string]*Package),
		selectorFullPkg: make(map[string]string),
		wantedPkg:       make(map[string]bool),
		localvar:        make(map[string]*ExprType),
		sep:             ", ",
	}
}

// loadState is the state of a package's parsing.
type loadState int

const (
	notLoaded loadState = iota
	loading
	loaded
)

func (w *Walker) Features(ctx string) (fs []string) {
	for pkg, p := range w.packageMap {
		if w.wantedPkg[p.name] {
			if ctx == "" || strings.HasPrefix(pkg, ctx) {
				fs = append(fs, p.Features()...)
			}
		}
	}
	sort.Strings(fs)
	return
}

// fileDeps returns the imports in a file.
func fileDeps(f *ast.File) (pkgs []string) {
	for _, is := range f.Imports {
		fpkg, err := strconv.Unquote(is.Path.Value)
		if err != nil {
			log.Fatalf("error unquoting import string %q: %v", is.Path.Value, err)
		}
		if fpkg != "C" {
			pkgs = append(pkgs, fpkg)
		}
	}
	return
}

func (w *Walker) findPackage(pkg string) *Package {
	if full, ok := w.selectorFullPkg[pkg]; ok {
		if w.ctxName != "" {
			ctxName := w.ctxName + full
			for k, v := range w.packageMap {
				if k == ctxName {
					return v
				}
			}
		}
		for k, v := range w.packageMap {
			if k == full {
				return v
			}
		}
	}
	return nil
}

func (w *Walker) findPackageOSArch(pkg string) *Package {
	if full, ok := w.selectorFullPkg[pkg]; ok {
		ctxName := osArchName(w.context) + ":" + full
		for k, v := range w.packageMap {
			if k == ctxName {
				return v
			}
		}
	}
	return nil
}

// WalkPackage walks all files in package `name'.
// WalkPackage does nothing if the package has already been loaded.

func (w *Walker) WalkPackage(pkg string) {
	if build.IsLocalImport(pkg) {
		wd, err := os.Getwd()
		if err != nil {
			if apiVerbose {
				log.Println(err)
			}
			return
		}
		dir := filepath.Clean(filepath.Join(wd, pkg))
		bp, err := w.context.ImportDir(dir, 0)
		if err != nil {
			if apiVerbose {
				log.Println(err)
			}
			return
		}
		if w.wantedPkg[pkg] == true {
			w.wantedPkg[bp.Name] = true
			delete(w.wantedPkg, pkg)
		}
		if w.cursorInfo != nil && w.cursorInfo.pkg == pkg {
			w.cursorInfo.pkg = bp.Name
		}
		w.WalkPackageDir(bp.Name, bp.Dir, bp)
	} else if filepath.IsAbs(pkg) {
		bp, err := build.ImportDir(pkg, 0)
		if err != nil {
			if apiVerbose {
				log.Println(err)
			}
		}
		if w.wantedPkg[pkg] == true {
			w.wantedPkg[bp.Name] = true
			delete(w.wantedPkg, pkg)
		}
		if w.cursorInfo != nil && w.cursorInfo.pkg == pkg {
			w.cursorInfo.pkg = bp.Name
		}

		w.WalkPackageDir(bp.Name, bp.Dir, bp)
	} else {
		bp, err := build.Import(pkg, "", build.FindOnly)
		if err != nil {
			if apiVerbose {
				log.Println(err)
			}
			return
		}
		w.WalkPackageDir(pkg, bp.Dir, nil)
	}
}

func (w *Walker) WalkPackageDir(name string, dir string, bp *build.Package) {
	ctxName := w.ctxName + name
	curName := name
	switch w.packageState[ctxName] {
	case loading:
		log.Fatalf("import cycle loading package %q?", name)
		return
	case loaded:
		return
	}
	w.packageState[ctxName] = loading
	w.selectorFullPkg[name] = name

	defer func() {
		w.packageState[ctxName] = loaded
	}()

	sname := name[strings.LastIndexAny(name, ".-/\\")+1:]

	apkg := &ast.Package{
		Files: make(map[string]*ast.File),
	}
	if bp == nil {
		bp, _ = w.context.ImportDir(dir, 0)
	}
	if bp == nil {
		return
	}
	if w.ctxName != "" {
		isCgo := (len(bp.CgoFiles) > 0) && w.context.CgoEnabled
		if isCgo {
			curName = ctxName
		} else {
			isOSArch := false
			for _, file := range bp.GoFiles {
				if isOSArchFile(w.context, file) {
					isOSArch = true
					break
				}
			}
			var p *Package
			if isOSArch {
				curName = osArchName(w.context) + ":" + name
				p = w.findPackageOSArch(name)
			} else {
				curName = name
				p = w.findPackage(name)
			}
			if p != nil {
				if ApiImportParser {
					for _, dep := range p.deps {
						if _, ok := w.packageState[dep]; ok {
							continue
						}
						w.WalkPackage(dep)
					}
				}
				w.packageMap[ctxName] = p
				return
			}
		}
	}

	files := append(append([]string{}, bp.GoFiles...), bp.CgoFiles...)

	if w.cursorInfo != nil && w.cursorInfo.pkg == name {
		files = append(files, bp.TestGoFiles...)
		for _, v := range bp.XTestGoFiles {
			if v == w.cursorInfo.file {
				var xbp build.Package
				xbp.Name = name + "_test"
				xbp.GoFiles = append(xbp.GoFiles, bp.XTestGoFiles...)
				w.cursorInfo.pkg = xbp.Name
				w.WalkPackageDir(xbp.Name, dir, &xbp)
				break
			}
		}
	}

	if len(files) == 0 {
		if apiVerbose {
			log.Println("no Go source files in", bp.Dir)
		}
		return
	}
	var deps []string

	for _, file := range files {
		var src interface{} = nil
		if w.cursorInfo != nil &&
			w.cursorInfo.pkg == name &&
			w.cursorInfo.file == file &&
			w.cursorInfo.std {
			src = w.cursorInfo.src
		}
		f, err := parser.ParseFile(w.fset, filepath.Join(dir, file), src, 0)
		if err != nil {
			if apiVerbose {
				log.Printf("error parsing package %s, file %s: %v", name, file, err)
			}
		}

		if sname != f.Name.Name {
			continue
		}
		apkg.Files[file] = f
		if ApiImportParser {
			deps = fileDeps(f)
			for _, dep := range deps {
				if _, ok := w.packageState[dep]; ok {
					continue
				}
				w.WalkPackage(dep)
			}
		}
		if apiShowpos && w.wantedPkg[name] {
			tf := w.fset.File(f.Pos())
			if tf != nil {
				fmt.Printf("pos %s%s%s%s%d%s%d\n", name, w.sep, filepath.Join(dir, file), w.sep, tf.Base(), w.sep, tf.Size())
			}
		}
	}
	/* else {
		fdir, err := os.Open(dir)
		if err != nil {
			log.Fatalln(err)
		}
		infos, err := fdir.Readdir(-1)
		fdir.Close()
		if err != nil {
			log.Fatalln(err)
		}

		for _, info := range infos {
			if info.IsDir() {
				continue
			}
			file := info.Name()
			if strings.HasPrefix(file, "_") || strings.HasSuffix(file, "_test.go") {
				continue
			}
			if strings.HasSuffix(file, ".go") {
				f, err := parser.ParseFile(w.fset, filepath.Join(dir, file), nil, 0)
				if err != nil {
					if apiVerbose {
						log.Printf("error parsing package %s, file %s: %v", name, file, err)
					}
					continue
				}
				if f.Name.Name != sname {
					continue
				}

				apkg.Files[file] = f
				if apiImportParser {
					for _, dep := range fileDeps(f) {
						w.WalkPackage(dep)
					}
				}
				if apiShowpos && w.wantedPkg[name] {
					tf := w.fset.File(f.Pos())
					if tf != nil {
						fmt.Printf("pos %s%s%s%s%d:%d\n", name, w.sep, filepath.Join(dir, file), w.sep, tf.Base(), tf.Base()+tf.Size())
					}
				}
			}
		}
	}*/
	if curName != ctxName {
		w.packageState[curName] = loading

		defer func() {
			w.packageState[curName] = loaded
		}()
	}

	if apiVerbose {
		log.Printf("package %s => %s, %v", ctxName, curName, w.wantedPkg[curName])
	}
	pop := w.pushScope("pkg " + name)
	defer pop()

	w.curPackageName = curName
	w.constDep = map[string]*ExprType{}
	w.curPackage = NewPackage()
	w.curPackage.apkg = apkg
	w.curPackage.name = name
	w.curPackage.dir = dir
	w.curPackage.deps = deps
	w.curPackage.sep = w.sep
	w.packageMap[curName] = w.curPackage
	w.packageMap[ctxName] = w.curPackage

	for _, afile := range apkg.Files {
		w.recordTypes(afile)
	}

	// Register all function declarations first.
	for _, afile := range apkg.Files {
		for _, di := range afile.Decls {
			if d, ok := di.(*ast.FuncDecl); ok {
				if !w.isExtract(d.Name.Name) {
					continue
				}
				w.peekFuncDecl(d)
			}
		}
	}

	for _, afile := range apkg.Files {
		w.walkFile(afile)
	}

	w.resolveConstantDeps()

	if w.cursorInfo != nil && w.cursorInfo.pkg == name {
		for k, v := range apkg.Files {
			if k == w.cursorInfo.file {
				f := w.fset.File(v.Pos())
				if f == nil {
					log.Fatalf("error fset postion %v", v.Pos())
				}
				info, err := w.lookupFile(v, token.Pos(f.Base())+w.cursorInfo.pos-1)
				if err != nil {
					log.Fatalln("lookup error,", err)
				} else {
					if info != nil && info.Kind == KindImport {
						for _, is := range v.Imports {
							fpath, err := strconv.Unquote(is.Path.Value)
							if err == nil {
								if info.Name == path.Base(fpath) {
									info.T = is.Path
								}
							}
						}
					}
					w.cursorInfo.info = info
				}
				break
			}
		}
		return
	}

	// Now that we're done walking types, vars and consts
	// in the *ast.Package, use go/doc to do the rest
	// (functions and methods). This is done here because
	// go/doc is destructive.  We can't use the
	// *ast.Package after this.
	var mode doc.Mode
	if ApiAllmethods {
		mode |= doc.AllMethods
	}
	if ApiAlldecls && w.wantedPkg[w.ctxName] {
		mode |= doc.AllDecls
	}

	dpkg := doc.New(apkg, name, mode)
	w.curPackage.dpkg = dpkg

	if w.wantedPkg[name] != true {
		return
	}

	for _, t := range dpkg.Types {
		// Move funcs up to the top-level, not hiding in the Types.
		dpkg.Funcs = append(dpkg.Funcs, t.Funcs...)

		for _, m := range t.Methods {
			w.walkFuncDecl(m.Decl)
		}
	}

	for _, f := range dpkg.Funcs {
		w.walkFuncDecl(f.Decl)
	}
}

// pushScope enters a new scope (walking a package, type, node, etc)
// and returns a function that will leave the scope (with sanity checking
// for mismatched pushes & pops)
func (w *Walker) pushScope(name string) (popFunc func()) {
	w.scope = append(w.scope, name)
	return func() {
		if len(w.scope) == 0 {
			log.Fatalf("attempt to leave scope %q with empty scope list", name)
		}
		if w.scope[len(w.scope)-1] != name {
			log.Fatalf("attempt to leave scope %q, but scope is currently %#v", name, w.scope)
		}
		w.scope = w.scope[:len(w.scope)-1]
	}
}

func (w *Walker) recordTypes(file *ast.File) {
	cur := w.curPackage
	for _, di := range file.Decls {
		switch d := di.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.TYPE:
				for _, sp := range d.Specs {
					ts := sp.(*ast.TypeSpec)
					name := ts.Name.Name
					switch t := ts.Type.(type) {
					case *ast.InterfaceType:
						if isExtract(name) {
							w.noteInterface(name, t)
						}
						cur.interfaces[name] = t
					case *ast.StructType:
						cur.structs[name] = t
					default:
						cur.types[name] = ts.Type
					}
				}
			}
		}
	}
}

func inRange(node ast.Node, p token.Pos) bool {
	if node == nil {
		return false
	}
	return p >= node.Pos() && p <= node.End()
}

func (w *Walker) lookupLabel(body *ast.BlockStmt, name string) (*TypeInfo, error) {
	for _, stmt := range body.List {
		switch v := stmt.(type) {
		case *ast.BlockStmt:
			return w.lookupLabel(v, name)
		case *ast.LabeledStmt:
			return &TypeInfo{Kind: KindLabel, Name: v.Label.Name, Type: "branch", T: v.Label}, nil
		}
	}
	return nil, nil
}

func (w *Walker) lookupFile(file *ast.File, p token.Pos) (*TypeInfo, error) {
	if inRange(file.Name, p) {
		return &TypeInfo{Kind: KindPackage, X: file.Name, Name: file.Name.Name, Type: file.Name.Name, T: file.Name}, nil
	}
	for _, di := range file.Decls {
		switch d := di.(type) {
		case *ast.GenDecl:
			if inRange(d, p) {
				return w.lookupDecl(d, p, false)
			}
		case *ast.FuncDecl:
			if inRange(d, p) {
				info, err := w.lookupDecl(d, p, false)
				if info != nil && info.Kind == KindBranch {
					return w.lookupLabel(d.Body, info.Name)
				}
				return info, err
			}
			if d.Body != nil && inRange(d.Body, p) {
				return w.lookupStmt(d.Body, p)
			}
		default:
			return nil, fmt.Errorf("un parser decl %T", di)
		}
	}
	return nil, fmt.Errorf("un find cursor %v", w.fset.Position(p))
}

func (w *Walker) isExtract(name string) bool {
	if w.wantedPkg[w.curPackageName] || ApiAlldecls {
		return true
	}
	return ast.IsExported(name)
}

func (w *Walker) isType(typ string) *ExprType {
	pos := strings.Index(typ, ".")
	if pos != -1 {
		pkg := typ[:pos]
		typ = typ[pos+1:]
		if p := w.findPackage(pkg); p != nil {
			if t, ok := p.types[typ]; ok {
				if r := w.isType(typ); r != nil {
					return r
				}
				return &ExprType{X: t, T: w.pkgRetType(pkg, w.nodeString(t))}
			}
		}
		return nil
	}
	if t, ok := w.curPackage.types[typ]; ok {
		if r := w.isType(w.nodeString(t)); r != nil {
			return r
		}
		return &ExprType{X: t, T: w.nodeString(t)}
	}
	return nil
}

func (w *Walker) lookupStmt(vi ast.Stmt, p token.Pos) (*TypeInfo, error) {
	if vi == nil {
		return nil, nil
	}
	switch v := vi.(type) {
	case *ast.BadStmt:
		//
	case *ast.EmptyStmt:
		//
	case *ast.LabeledStmt:
		if inRange(v.Label, p) {
			return &TypeInfo{Kind: KindLabel, Name: v.Label.Name}, nil
		}
		return w.lookupStmt(v.Stmt, p)
		//
	case *ast.DeclStmt:
		return w.lookupDecl(v.Decl, p, true)
	case *ast.AssignStmt:
		if len(v.Lhs) == len(v.Rhs) {
			for i := 0; i < len(v.Lhs); i++ {
				switch lt := v.Lhs[i].(type) {
				case *ast.Ident:
					typ, err := w.varValueType(v.Rhs[i], 0)
					if err == nil && v.Tok == token.DEFINE {
						w.localvar[lt.Name] = &ExprType{T: typ, X: lt}
					} else if apiVerbose {
						log.Println(err)
					}
				}
				if inRange(v.Lhs[i], p) {
					return w.lookupExprInfo(v.Lhs[i], p)
				} else if inRange(v.Rhs[i], p) {
					return w.lookupExprInfo(v.Rhs[i], p)
				}
				if fl, ok := v.Rhs[i].(*ast.FuncLit); ok {
					if inRange(fl, p) {
						return w.lookupStmt(fl.Body, p)
					}
				}
			}
		} else if len(v.Rhs) == 1 {
			for i := 0; i < len(v.Lhs); i++ {
				switch lt := v.Lhs[i].(type) {
				case *ast.Ident:
					typ, err := w.varValueType(v.Rhs[0], i)
					if err == nil && v.Tok == token.DEFINE {
						w.localvar[lt.Name] = &ExprType{T: typ, X: lt}
					} else if apiVerbose {
						log.Println(err)
					}
				}
				if inRange(v.Lhs[i], p) {
					return w.lookupExprInfo(v.Lhs[i], p)
				} else if inRange(v.Rhs[0], p) {
					return w.lookupExprInfo(v.Rhs[0], p)
				}
				if fl, ok := v.Rhs[0].(*ast.FuncLit); ok {
					if inRange(fl, p) {
						return w.lookupStmt(fl.Body, p)
					}
				}
			}
		}
		return nil, nil
	case *ast.ExprStmt:
		return w.lookupExprInfo(v.X, p)
	case *ast.BlockStmt:
		for _, st := range v.List {
			if inRange(st, p) {
				return w.lookupStmt(st, p)
			}
			_, err := w.lookupStmt(st, p)
			if err != nil {
				log.Println(err)
			}
		}
	case *ast.IfStmt:
		if inRange(v.Init, p) {
			return w.lookupStmt(v.Init, p)
		} else {
			w.lookupStmt(v.Init, p)
		}
		if inRange(v.Cond, p) {
			return w.lookupExprInfo(v.Cond, p)
		} else if inRange(v.Body, p) {
			return w.lookupStmt(v.Body, p)
		} else if inRange(v.Else, p) {
			return w.lookupStmt(v.Else, p)
		}
	case *ast.SendStmt:
		if inRange(v.Chan, p) {
			return w.lookupExprInfo(v.Chan, p)
		} else if inRange(v.Value, p) {
			return w.lookupExprInfo(v.Value, p)
		}
	case *ast.IncDecStmt:
		return w.lookupExprInfo(v.X, p)
	case *ast.GoStmt:
		return w.lookupExprInfo(v.Call, p)
	case *ast.DeferStmt:
		return w.lookupExprInfo(v.Call, p)
	case *ast.ReturnStmt:
		for _, r := range v.Results {
			if inRange(r, p) {
				return w.lookupExprInfo(r, p)
			}
		}
	case *ast.BranchStmt:
		if inRange(v.Label, p) {
			return &TypeInfo{Kind: KindBranch, Name: v.Label.Name, Type: "label", T: v.Label}, nil
		}
		//
	case *ast.CaseClause:
		for _, r := range v.List {
			if inRange(r, p) {
				return w.lookupExprInfo(r, p)
			}
		}
		for _, body := range v.Body {
			if inRange(body, p) {
				return w.lookupStmt(body, p)
			} else {
				w.lookupStmt(body, p)
			}
		}
	case *ast.SwitchStmt:
		if inRange(v.Init, p) {
			return w.lookupStmt(v.Init, p)
		} else {
			w.lookupStmt(v.Init, p)
		}
		if inRange(v.Tag, p) {
			return w.lookupExprInfo(v.Tag, p)
		} else if inRange(v.Body, p) {
			return w.lookupStmt(v.Body, p)
		}
	case *ast.TypeSwitchStmt:
		if inRange(v.Assign, p) {
			return w.lookupStmt(v.Assign, p)
		} else {
			w.lookupStmt(v.Assign, p)
		}
		if inRange(v.Init, p) {
			return w.lookupStmt(v.Init, p)
		} else {
			w.lookupStmt(v.Init, p)
		}
		var vs string
		if as, ok := v.Assign.(*ast.AssignStmt); ok {
			if len(as.Lhs) == 1 {
				vs = w.nodeString(as.Lhs[0])
			}
		}
		if inRange(v.Body, p) {
			for _, s := range v.Body.List {
				if inRange(s, p) {
					switch cs := s.(type) {
					case *ast.CaseClause:
						for _, r := range cs.List {
							if inRange(r, p) {
								return w.lookupExprInfo(r, p)
							} else if vs != "" {
								typ, err := w.varValueType(r, 0)
								if err == nil {
									w.localvar[vs] = &ExprType{T: typ, X: r}
								}
							}
						}
						for _, body := range cs.Body {
							if inRange(body, p) {
								return w.lookupStmt(body, p)
							} else {
								w.lookupStmt(body, p)
							}
						}
					default:
						return w.lookupStmt(cs, p)
					}
				}
			}
		}
	case *ast.CommClause:
		if inRange(v.Comm, p) {
			return w.lookupStmt(v.Comm, p)
		}
		for _, body := range v.Body {
			if inRange(body, p) {
				return w.lookupStmt(body, p)
			}
		}
	case *ast.SelectStmt:
		if inRange(v.Body, p) {
			return w.lookupStmt(v.Body, p)
		}
	case *ast.ForStmt:
		if inRange(v.Init, p) {
			return w.lookupStmt(v.Init, p)
		} else {
			w.lookupStmt(v.Init, p)
		}
		if inRange(v.Cond, p) {
			return w.lookupExprInfo(v.Cond, p)
		} else if inRange(v.Body, p) {
			return w.lookupStmt(v.Body, p)
		} else if inRange(v.Post, p) {
			return w.lookupStmt(v.Post, p)
		}
	case *ast.RangeStmt:
		if inRange(v.X, p) {
			return w.lookupExprInfo(v.X, p)
		} else if inRange(v.Key, p) {
			return &TypeInfo{Kind: KindBuiltin, Name: w.nodeString(v.Key), Type: "int"}, nil
		} else if inRange(v.Value, p) {
			typ, err := w.lookupExprInfo(v.X, p)
			if typ != nil {
				typ.Name = w.nodeString(v.Value)
				return typ, err
			}
		} else {
			typ, err := w.varValueType(v.X, 0)
			//check is type
			if t := w.isType(typ); t != nil {
				typ = t.T
			}
			if err == nil {
				var kt, vt string
				if strings.HasPrefix(typ, "[]") {
					kt = "int"
					vt = typ[2:]
				} else if strings.HasPrefix(typ, "map[") {
					node, err := parser.ParseExpr(typ + "{}")
					if err == nil {
						if cl, ok := node.(*ast.CompositeLit); ok {
							if m, ok := cl.Type.(*ast.MapType); ok {
								kt = w.nodeString(w.namelessType(m.Key))
								vt = w.nodeString(w.namelessType(m.Value))
							}
						}
					}
				}
				if inRange(v.Key, p) {
					return &TypeInfo{Kind: KindVar, X: v.Key, Name: w.nodeString(v.Key), T: v.X, Type: kt}, nil
				} else if inRange(v.Value, p) {
					return &TypeInfo{Kind: KindVar, X: v.Value, Name: w.nodeString(v.Value), T: v.X, Type: vt}, nil
				}
				if key, ok := v.Key.(*ast.Ident); ok {
					w.localvar[key.Name] = &ExprType{T: kt, X: v.Key}
				}
				if value, ok := v.Value.(*ast.Ident); ok {
					w.localvar[value.Name] = &ExprType{T: vt, X: v.Value}
				}
			}
		}
		if inRange(v.Body, p) {
			return w.lookupStmt(v.Body, p)
		}
	}
	return nil, nil //fmt.Errorf("not lookup stmt %v %T", vi, vi)
}

func (w *Walker) lookupVar(vs *ast.ValueSpec, p token.Pos, local bool) (*TypeInfo, error) {
	if inRange(vs.Type, p) {
		return w.lookupExprInfo(vs.Type, p)
	}
	for _, v := range vs.Values {
		if inRange(v, p) {
			return w.lookupExprInfo(v, p)
		}
	}
	if vs.Type != nil {
		typ := w.nodeString(vs.Type)
		for _, ident := range vs.Names {
			if local {
				w.localvar[ident.Name] = &ExprType{T: typ, X: ident}
			}
			if inRange(ident, p) {
				return &TypeInfo{Kind: KindVar, X: ident, Name: ident.Name, T: vs.Type, Type: typ}, nil
			}
		}
	} else if len(vs.Names) == len(vs.Values) {
		for n, ident := range vs.Names {
			typ := ""
			if !local {
				if t, ok := w.curPackage.vars[ident.Name]; ok {
					typ = t.T
				}
			} else {
				typ, err := w.varValueType(vs.Values[n], n)
				if err != nil {
					if apiVerbose {
						log.Printf("unknown type of variable2 %q, type %T, error = %v, pos=%s",
							ident.Name, vs.Values[n], err, w.fset.Position(vs.Pos()))
					}
					typ = "unknown-type"
				}
				w.localvar[ident.Name] = &ExprType{T: typ, X: ident}
			}
			if inRange(ident, p) {
				return &TypeInfo{Kind: KindVar, X: ident, Name: ident.Name, T: ident, Type: typ}, nil
			}
		}
	} else if len(vs.Values) == 1 {
		for n, ident := range vs.Names {
			typ := ""
			if !local {
				if t, ok := w.curPackage.vars[ident.Name]; ok {
					typ = t.T
				}
			} else {
				typ, err := w.varValueType(vs.Values[0], n)
				if err != nil {
					if apiVerbose {
						log.Printf("unknown type of variable3 %q, type %T, error = %v, pos=%s",
							ident.Name, vs.Values[0], err, w.fset.Position(vs.Pos()))
					}
					typ = "unknown-type"
				}
				w.localvar[ident.Name] = &ExprType{T: typ, X: ident}
			}
			if inRange(ident, p) {
				return &TypeInfo{Kind: KindVar, X: ident, Name: ident.Name, T: ident, Type: typ}, nil
			}
		}
	}
	return nil, fmt.Errorf("not lookup var local:%v value:%v type:s%T", local, w.nodeString(vs), vs)
}

func (w *Walker) lookupConst(vs *ast.ValueSpec, p token.Pos, local bool) (*TypeInfo, error) {
	if inRange(vs.Type, p) {
		return w.lookupExprInfo(vs.Type, p)
	}
	for _, ident := range vs.Names {
		typ := ""
		if !local {
			if t, ok := w.curPackage.consts[ident.Name]; ok {
				typ = t.T
			}
		} else {
			litType := ""
			if vs.Type != nil {
				litType = w.nodeString(vs.Type)
			} else {
				litType = w.lastConstType
				if vs.Values != nil {
					if len(vs.Values) != 1 {
						if apiVerbose {
							log.Printf("const %q, values: %#v", ident.Name, vs.Values)
						}
						return nil, nil
					}
					var err error
					litType, err = w.constValueType(vs.Values[0])
					if err != nil {
						if apiVerbose {
							log.Printf("unknown kind in const %q (%T): %v", ident.Name, vs.Values[0], err)
						}
						litType = "unknown-type"
					}
				}
			}
			w.lastConstType = litType
			typ = litType
			w.localvar[ident.Name] = &ExprType{T: typ, X: ident}
		}
		if inRange(ident, p) {
			return &TypeInfo{Kind: KindConst, X: ident, Name: ident.Name, T: ident, Type: typ}, nil
		}
	}
	return nil, nil
}

func (w *Walker) lookupType(ts *ast.TypeSpec, p token.Pos, local bool) (*TypeInfo, error) {
	switch t := ts.Type.(type) {
	case *ast.StructType:
		if inRange(t.Fields, p) {
			for _, fd := range t.Fields.List {
				if inRange(fd.Type, p) {
					return w.lookupExprInfo(fd.Type, p)
				}
				for _, ident := range fd.Names {
					if inRange(ident, p) {
						return &TypeInfo{Kind: KindField, X: ident, Name: ts.Name.Name + "." + ident.Name, T: fd.Type, Type: w.nodeString(w.namelessType(fd.Type))}, nil
					}
				}
			}
		}
		return &TypeInfo{Kind: KindStruct, X: ts.Name, Name: ts.Name.Name, T: ts.Type, Type: "struct"}, nil
	case *ast.InterfaceType:
		if inRange(t.Methods, p) {
			for _, fd := range t.Methods.List {
				for _, ident := range fd.Names {
					if inRange(ident, p) {
						return &TypeInfo{Kind: KindMethod, X: ident, Name: ts.Name.Name + "." + ident.Name, T: ident, Type: w.nodeString(w.namelessType(fd.Type))}, nil
					}
				}
				if inRange(fd.Type, p) {
					return w.lookupExprInfo(fd.Type, p)
				}
			}
		}
		return &TypeInfo{Kind: KindInterface, X: ts.Name, Name: ts.Name.Name, T: ts.Type, Type: "interface"}, nil
	default:
		return &TypeInfo{Kind: KindType, X: ts.Name, Name: ts.Name.Name, T: ts.Type, Type: w.nodeString(w.namelessType(ts.Type))}, nil
	}
	return nil, nil
}

func (w *Walker) lookupDecl(di ast.Decl, p token.Pos, local bool) (*TypeInfo, error) {
	switch d := di.(type) {
	case *ast.GenDecl:
		switch d.Tok {
		case token.IMPORT:
			for _, sp := range d.Specs {
				is := sp.(*ast.ImportSpec)
				fpath, err := strconv.Unquote(is.Path.Value)
				if err != nil {
					return nil, err
				}
				name := path.Base(fpath)
				if is.Name != nil {
					name = is.Name.Name
				}
				if inRange(sp, p) {
					return &TypeInfo{Kind: KindImport, X: is.Name, Name: name, T: is.Name, Type: fpath}, nil
				}
			}
		case token.CONST:
			for _, sp := range d.Specs {
				if inRange(sp, p) {
					return w.lookupConst(sp.(*ast.ValueSpec), p, local)
				} else {
					w.lookupConst(sp.(*ast.ValueSpec), p, local)
				}
			}
			return nil, nil
		case token.TYPE:
			for _, sp := range d.Specs {
				if inRange(sp, p) {
					return w.lookupType(sp.(*ast.TypeSpec), p, local)
				} else {
					w.lookupType(sp.(*ast.TypeSpec), p, local)
				}
			}
		case token.VAR:
			for _, sp := range d.Specs {
				if inRange(sp, p) {
					return w.lookupVar(sp.(*ast.ValueSpec), p, local)
				} else {
					w.lookupVar(sp.(*ast.ValueSpec), p, local)
				}
			}
			return nil, nil
		default:
			return nil, fmt.Errorf("unknown token type %d %T in GenDecl", d.Tok, d)
		}
	case *ast.FuncDecl:
		if d.Type.Params != nil {
			for _, fd := range d.Type.Params.List {
				if inRange(fd, p) {
					return w.lookupExprInfo(fd.Type, p)
				}
				for _, ident := range fd.Names {
					if inRange(ident, p) {
						info, err := w.lookupExprInfo(fd.Type, p)
						if err == nil {
							return &TypeInfo{Kind: KindParam, X: ident, Name: ident.Name, T: info.T, Type: info.Type}, nil
						}
					}
					typ, err := w.varValueType(fd.Type, 0)
					if err == nil {
						w.localvar[ident.Name] = &ExprType{T: typ, X: ident}
					} else if apiVerbose {
						log.Println(err)
					}
				}
			}
		}
		if d.Type.Results != nil {
			for _, fd := range d.Type.Results.List {
				if inRange(fd, p) {
					return w.lookupExprInfo(fd.Type, p)
				}
				for _, ident := range fd.Names {
					typ, err := w.varValueType(fd.Type, 0)
					if err == nil {
						w.localvar[ident.Name] = &ExprType{T: typ, X: ident}
					}
				}
			}
		}
		if d.Recv != nil {
			for _, fd := range d.Recv.List {
				if inRange(fd, p) {
					return w.lookupExprInfo(fd.Type, p)
				}
				for _, ident := range fd.Names {
					w.localvar[ident.Name] = &ExprType{T: w.nodeString(fd.Type), X: ident}
				}
			}
		}
		if inRange(d.Body, p) {
			return w.lookupStmt(d.Body, p)
		}
		var fname = d.Name.Name
		kind := KindFunc
		if d.Recv != nil {
			recvTypeName, imp := baseTypeName(d.Recv.List[0].Type)
			if imp {
				return nil, nil
			}
			fname = recvTypeName + "." + d.Name.Name
			kind = KindMethod
		}
		return &TypeInfo{Kind: kind, X: d.Name, Name: fname, T: d.Type, Type: w.nodeString(w.namelessType(d.Type))}, nil
	default:
		return nil, fmt.Errorf("unhandled %T, %#v\n", di, di)
	}
	return nil, fmt.Errorf("not lookupDecl %v %T", w.nodeString(di), di)
}

func (w *Walker) lookupExprInfo(vi ast.Expr, p token.Pos) (*TypeInfo, error) {
	_, info, err := w.lookupExpr(vi, p)
	return info, err
}

// lookupExpr , return name,info,error
func (w *Walker) lookupExpr(vi ast.Expr, p token.Pos) (string, *TypeInfo, error) {
	if apiVerbose {
		log.Printf("lookup expr %v %T", w.nodeString(vi), vi)
	}
	switch v := vi.(type) {
	case *ast.BasicLit:
		litType, ok := varType[v.Kind]
		if !ok {
			return "", nil, fmt.Errorf("unknown basic literal kind %#v", v)
		}
		name := v.Value
		if len(name) >= 128 {
			name = name[:128] + "..."
		}
		return litType, &TypeInfo{Kind: KindBuiltin, X: v, Name: name, T: v, Type: litType}, nil
	case *ast.StarExpr:
		s, info, err := w.lookupExpr(v.X, p)
		if err != nil {
			return "", nil, err
		}
		return "*" + s, &TypeInfo{Kind: info.Kind, X: v, Name: "*" + info.Name, T: info.T, Type: "*" + info.Type}, err
	case *ast.InterfaceType:
		return "interface{}", &TypeInfo{Kind: KindInterface, X: v, Name: w.nodeString(v), T: v, Type: "interface{}"}, nil
	case *ast.Ellipsis:
		s, info, err := w.lookupExpr(v.Elt, p)
		if err != nil {
			return "", nil, err
		}
		return "[]" + s, &TypeInfo{Kind: KindArray, X: v.Elt, Name: "..." + s, T: info.T, Type: "[]" + info.Type}, nil
	case *ast.KeyValueExpr:
		if inRange(v.Key, p) {
			return w.lookupExpr(v.Key, p)
		} else if inRange(v.Value, p) {
			return w.lookupExpr(v.Value, p)
		}
	case *ast.CompositeLit:
		typ, err := w.varValueType(v.Type, 0)
		if err == nil {
			typ = strings.TrimLeft(typ, "*")
			if strings.HasPrefix(typ, "[]") {
				typ = strings.TrimLeft(typ[2:], "*")
			}
			pos := strings.Index(typ, ".")
			var pt *Package = w.curPackage
			var pkgdot string
			if pos != -1 {
				pkg := typ[:pos]
				typ = typ[pos+1:]
				pt = w.findPackage(pkg)
				if pt != nil {
					pkgdot = pkg + "."
				}
			}
			if pt != nil {
				if ss, ok := pt.structs[typ]; ok {
					for _, elt := range v.Elts {
						if inRange(elt, p) {
							if cl, ok := elt.(*ast.CompositeLit); ok {
								for _, elt := range cl.Elts {
									if inRange(elt, p) {
										if kv, ok := elt.(*ast.KeyValueExpr); ok {
											if inRange(kv.Key, p) {
												n, t := w.findStructField(ss, w.nodeString(kv.Key))
												if n != nil {
													return pkgdot + typ + "." + w.nodeString(kv.Key), &TypeInfo{Kind: KindField, X: kv.Key, Name: pkgdot + typ + "." + w.nodeString(kv.Key), T: n, Type: w.nodeString(w.namelessType(t))}, nil
												}
											} else if inRange(kv.Value, p) {
												return w.lookupExpr(kv.Value, p)
											}
										}
									}
								}
							}
							if kv, ok := elt.(*ast.KeyValueExpr); ok {
								if inRange(kv.Key, p) {
									n, t := w.findStructField(ss, w.nodeString(kv.Key))
									if n != nil {
										return typ + "." + w.nodeString(kv.Key), &TypeInfo{Kind: KindField, X: kv.Key, Name: typ + "." + w.nodeString(kv.Key), T: n, Type: w.nodeString(w.namelessType(t))}, nil
									}
								} else if inRange(kv.Value, p) {
									return w.lookupExpr(kv.Value, p)
								}
							}
						}
					}
				}
			}
		}
		for _, elt := range v.Elts {
			if inRange(elt, p) {
				return w.lookupExpr(elt, p)
			}
		}
		return w.lookupExpr(v.Type, p)
	case *ast.UnaryExpr:
		s, info, err := w.lookupExpr(v.X, p)
		return v.Op.String() + s, info, err
	case *ast.TypeAssertExpr:
		if inRange(v.X, p) {
			return w.lookupExpr(v.X, p)
		}
		return w.lookupExpr(v.Type, p)
	case *ast.BinaryExpr:
		if inRange(v.X, p) {
			return w.lookupExpr(v.X, p)
		} else if inRange(v.Y, p) {
			return w.lookupExpr(v.Y, p)
		}
		return "", nil, nil
	case *ast.CallExpr:
		for _, arg := range v.Args {
			if inRange(arg, p) {
				return w.lookupExpr(arg, p)
			}
		}
		switch ft := v.Fun.(type) {
		case *ast.Ident:
			if typ, ok := w.localvar[ft.Name]; ok {
				return ft.Name, &TypeInfo{Kind: KindVar, X: ft, Name: ft.Name, T: typ.X, Type: typ.T}, nil
			}
			if typ, ok := w.curPackage.vars[ft.Name]; ok {
				return ft.Name, &TypeInfo{Kind: KindVar, X: v, Name: ft.Name, T: typ.X, Type: typ.T}, nil
			}
			if typ, ok := w.curPackage.functions[ft.Name]; ok {
				return ft.Name, &TypeInfo{Kind: KindFunc, X: ft, Name: ft.Name, T: typ.ft, Type: typ.sig}, nil
			}
			if typ, ok := w.curPackage.interfaces[ft.Name]; ok {
				return ft.Name, &TypeInfo{Kind: KindInterface, X: ft, Name: ft.Name, T: typ, Type: w.nodeString(w.namelessType(typ))}, nil
			}
			if typ, ok := w.curPackage.interfaces[ft.Name]; ok {
				return ft.Name, &TypeInfo{Kind: KindInterface, X: ft, Name: ft.Name, T: typ, Type: w.nodeString(w.namelessType(typ))}, nil
			}
			if typ, ok := w.curPackage.structs[ft.Name]; ok {
				return ft.Name, &TypeInfo{Kind: KindStruct, X: ft, Name: ft.Name, T: typ, Type: w.nodeString(w.namelessType(typ))}, nil
			}
			if typ, ok := w.curPackage.types[ft.Name]; ok {
				return ft.Name, &TypeInfo{Kind: KindType, X: ft, Name: ft.Name, T: typ, Type: w.nodeString(w.namelessType(typ))}, nil
			}
			if isBuiltinType(ft.Name) {
				return ft.Name, &TypeInfo{Kind: KindBuiltin, X: ft, Name: ft.Name}, nil
			}
			return "", nil, fmt.Errorf("lookup unknown ident %v", v)
		case *ast.FuncLit:
			if inRange(ft.Body, p) {
				info, err := w.lookupStmt(ft.Body, p)
				if err == nil {
					return "", info, nil
				}
				return "", nil, err
			}
			return w.lookupExpr(ft.Type, p)
		case *ast.ParenExpr:
			return w.lookupExpr(ft.X, p)
		case *ast.SelectorExpr:
			switch st := ft.X.(type) {
			case *ast.Ident:
				if inRange(st, p) {
					return w.lookupExpr(st, p)
				}
				s, info, err := w.lookupExpr(st, p)
				if err != nil {
					return "", nil, err
				}
				typ := info.Type
				if typ == "" {
					typ = s
				}
				fname := typ + "." + ft.Sel.Name
				typ = strings.TrimLeft(typ, "*")
				if fn, ok := w.curPackage.functions[fname]; ok {
					return fname, &TypeInfo{Kind: KindMethod, X: st, Name: fname, T: fn.ft, Type: w.nodeString(w.namelessType(fn.ft))}, nil
				}
				info, e := w.lookupFunction(typ, ft.Sel.Name)
				if e != nil {
					return "", nil, e
				}
				return fname, info, nil
			case *ast.SelectorExpr:
				if inRange(st.X, p) {
					return w.lookupExpr(st.X, p)
				}
				if inRange(st, p) {
					return w.lookupExpr(st, p)
				}
				typ, err := w.varValueType(st, 0)
				if err != nil {
					return "", nil, err
				}
				/*
					typ = strings.TrimLeft(typ, "*")
					if t := w.curPackage.findType(typ); t != nil {
						if ss, ok := t.(*ast.StructType); ok {
							for _, fi := range ss.Fields.List {
								for _, n := range fi.Names {
									if n.Name == st.Sel.Name {
										//return fname, &TypeInfo{Kind: KindField, X: n, Name: fname, T: fi.Type, Type: w.nodeString(w.namelessType(fi.Type))}, nil
										typ = w.nodeString(w.namelessType(fi.Type))
									}
								}
							}
						}
					}
				*/
				info, e := w.lookupFunction(typ, ft.Sel.Name)
				if e != nil {
					return "", nil, e
				}
				return typ + "." + st.Sel.Name, info, nil
			case *ast.CallExpr:
				if inRange(st, p) {
					return w.lookupExpr(st, p)
				}
				if info, err := w.lookupExprInfo(st, p); err == nil {
					if fn, ok := info.X.(*ast.FuncType); ok {
						if fn.Results.NumFields() == 1 {
							info, err := w.lookupFunction(w.nodeString(fn.Results.List[0].Type), ft.Sel.Name)
							if err == nil {
								return info.Name, info, err
							}
							return "", nil, err
						}
					}
				}
				//w.lookupFunction(w.nodeString(info.X))
				typ, err := w.varValueType(st, 0)
				if err != nil {
					return "", nil, err
				}
				info, e := w.lookupFunction(typ, ft.Sel.Name)
				if e != nil {
					return "", nil, e
				}
				return typ + "." + ft.Sel.Name, info, nil
			case *ast.TypeAssertExpr:
				if inRange(st.X, p) {
					return w.lookupExpr(st.X, p)
				}
				typ := w.nodeString(w.namelessType(st.Type))
				info, e := w.lookupFunction(typ, ft.Sel.Name)
				if e != nil {
					return "", nil, e
				}
				return typ + "." + ft.Sel.Name, info, nil
			default:
				return "", nil, fmt.Errorf("not find select %v %T", v, st)
			}
		}
		return "", nil, fmt.Errorf("not find call %v %T", w.nodeString(v), v.Fun)
	case *ast.SelectorExpr:
		switch st := v.X.(type) {
		case *ast.Ident:
			if inRange(st, p) {
				return w.lookupExpr(st, p)
			}
			info, err := w.lookupSelector(st.Name, v.Sel.Name)
			if err != nil {
				return "", nil, err
			}
			return st.Name + "." + v.Sel.Name, info, nil
			//		case *ast.CallExpr:
			//			typ, err := w.varValueType(v.X, index)
			//			if err == nil {
			//				if strings.HasPrefix(typ, "*") {
			//					typ = typ[1:]
			//				}
			//				t := w.curPackage.findType(typ)
			//				if st, ok := t.(*ast.StructType); ok {
			//					for _, fi := range st.Fields.List {
			//						for _, n := range fi.Names {
			//							if n.Name == v.Sel.Name {
			//								return w.varValueType(fi.Type, index)
			//							}
			//						}
			//					}
			//				}
			//			}
		case *ast.SelectorExpr:
			if inRange(st.X, p) {
				return w.lookupExpr(st.X, p)
			}

			if inRange(st, p) {
				return w.lookupExpr(st, p)
			}

			typ, err := w.varValueType(st, 0)
			if err == nil {
				info, err := w.lookupSelector(typ, v.Sel.Name)
				if err != nil {
					return "", nil, err
				}
				return typ + v.Sel.Name, info, nil
			}
			//		case *ast.IndexExpr:
			//			typ, err := w.varValueType(st.X, 0)
			//			log.Println(typ, err)
			//			if err == nil {
			//				if strings.HasPrefix(typ, "[]") {
			//					return w.varSelectorType(typ[2:], v.Sel.Name)
			//				}
			//			}
		}
		return "", nil, fmt.Errorf("unknown lookup selector expr: %T %s.%s", v.X, w.nodeString(v.X), v.Sel)

		//		s, info, err := w.lookupExpr(v.X, p)
		//		if err != nil {
		//			return "", "", err
		//		}
		//		if strings.HasPrefix(s, "*") {
		//			s = s[1:]
		//		}
		//		if inRange(v.X, p) {
		//			return s, info, err
		//		}
		//		t := w.curPackage.findType(s)
		//		fname := s + "." + v.Sel.Name
		//		if st, ok := t.(*ast.StructType); ok {
		//			for _, fi := range st.Fields.List {
		//				for _, n := range fi.Names {
		//					if n.Name == v.Sel.Name {
		//						return fname, fmt.Sprintf("var,%s,%s,%s", fname, w.nodeString(w.namelessType(fi.Type)), w.fset.Position(n.Pos())), nil
		//					}
		//				}
		//			}
		//		}
		//		log.Println(">>", s)
		//		info, e := w.lookupSelector(s, v.Sel.Name)
		//		return fname, info, e
	case *ast.Ident:
		if typ, ok := w.localvar[v.Name]; ok {
			return typ.T, &TypeInfo{Kind: KindVar, X: v, Name: v.Name, T: typ.X, Type: typ.T}, nil
		}
		if typ, ok := w.curPackage.interfaces[v.Name]; ok {
			return v.Name, &TypeInfo{Kind: KindInterface, X: v, Name: v.Name, T: typ, Type: "interface"}, nil
		}
		if typ, ok := w.curPackage.structs[v.Name]; ok {
			return v.Name, &TypeInfo{Kind: KindStruct, X: v, Name: v.Name, T: typ, Type: "struct"}, nil
		}
		if typ, ok := w.curPackage.types[v.Name]; ok {
			return v.Name, &TypeInfo{Kind: KindType, X: v, Name: v.Name, T: typ, Type: v.Name}, nil
		}
		if typ, ok := w.curPackage.vars[v.Name]; ok {
			return v.Name, &TypeInfo{Kind: KindVar, X: v, Name: v.Name, T: typ.X, Type: typ.T}, nil
		}
		if typ, ok := w.curPackage.consts[v.Name]; ok {
			return v.Name, &TypeInfo{Kind: KindConst, X: v, Name: v.Name, T: typ.X, Type: typ.T}, nil
		}
		if typ, ok := w.curPackage.functions[v.Name]; ok {
			return v.Name, &TypeInfo{Kind: KindFunc, X: typ.ft, Name: v.Name, T: typ.ft, Type: typ.sig}, nil
		}
		if p := w.findPackage(v.Name); p != nil {
			return v.Name, &TypeInfo{Kind: KindImport, X: v, Name: v.Name, Type: p.name}, nil
		}
		if isBuiltinType(v.Name) {
			return v.Name, &TypeInfo{Kind: KindBuiltin, Name: v.Name}, nil
		}
		return "", nil, fmt.Errorf("lookup unknown ident %v", v)
		//return v.Name, &TypeInfo{Kind: KindVar, X: v, Name: v.Name, T: v, Type: v.Name}, nil
	case *ast.IndexExpr:
		if inRange(v.Index, p) {
			return w.lookupExpr(v.Index, p)
		}
		return w.lookupExpr(v.X, p)
	case *ast.ParenExpr:
		return w.lookupExpr(v.X, p)
	case *ast.FuncLit:
		if inRange(v.Type, p) {
			return w.lookupExpr(v.Type, p)
		} else {
			w.lookupExpr(v.Type, p)
		}
		typ, err := w.varValueType(v.Type, 0)
		if err != nil {
			return "", nil, err
		}
		info, e := w.lookupStmt(v.Body, p)
		if e != nil {
			return "", nil, err
		}
		return typ, info, nil
	case *ast.FuncType:
		if v.Params != nil {
			for _, fd := range v.Params.List {
				if inRange(fd, p) {
					return w.lookupExpr(fd.Type, p)
				}
				for _, ident := range fd.Names {
					typ, err := w.varValueType(fd.Type, 0)
					if err == nil {
						w.localvar[ident.Name] = &ExprType{T: typ, X: ident}
					}
				}
			}
		}
		if v.Results != nil {
			for _, fd := range v.Results.List {
				if inRange(fd, p) {
					return w.lookupExpr(fd.Type, p)
				}
				for _, ident := range fd.Names {
					typ, err := w.varValueType(fd.Type, 0)
					if err == nil {
						w.localvar[ident.Name] = &ExprType{T: typ, X: ident}
					}
				}
			}
		}
		return "", nil, nil
	case *ast.ArrayType:
		s, info, err := w.lookupExpr(v.Elt, p)
		if err != nil {
			return "", nil, err
		}
		return "[]" + s, &TypeInfo{Kind: KindArray, Name: "[]" + info.Name, Type: "[]" + info.Type, T: info.T}, nil
	case *ast.SliceExpr:
		if inRange(v.High, p) {
			return w.lookupExpr(v.High, p)
		} else if inRange(v.Low, p) {
			return w.lookupExpr(v.Low, p)
		}
		return w.lookupExpr(v.X, p)
	case *ast.MapType:
		if inRange(v.Key, p) {
			return w.lookupExpr(v.Key, p)
		} else if inRange(v.Value, p) {
			return w.lookupExpr(v.Value, p)
		}
		typ, err := w.varValueType(v, 0)
		if err != nil {
			return "", nil, err
		}
		return typ, &TypeInfo{Kind: KindMap, X: v, Name: w.nodeString(v), T: v, Type: typ}, nil
	case *ast.ChanType:
		if inRange(v.Value, p) {
			return w.lookupExpr(v.Value, p)
		}
		typ, err := w.varValueType(v, 0)
		if err != nil {
			return "", nil, err
		}
		return typ, &TypeInfo{Kind: KindChan, X: v, Name: w.nodeString(v), T: v, Type: typ}, nil
	default:
		return "", nil, fmt.Errorf("not lookupExpr %v %T", w.nodeString(v), v)
	}
	return "", nil, fmt.Errorf("not lookupExpr %v %T", w.nodeString(vi), vi)
}

func (w *Walker) walkFile(file *ast.File) {
	// Not entering a scope here; file boundaries aren't interesting.
	for _, di := range file.Decls {
		switch d := di.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				for _, sp := range d.Specs {
					is := sp.(*ast.ImportSpec)
					fpath, err := strconv.Unquote(is.Path.Value)
					if err != nil {
						log.Fatal(err)
					}
					//name := path.Base(fpath)
					name := fpath
					if i := strings.LastIndexAny(name, ".-/\\"); i > 0 {
						name = name[i+1:]
					}
					if is.Name != nil {
						name = is.Name.Name
					}
					w.selectorFullPkg[name] = fpath
				}
			case token.CONST:
				for _, sp := range d.Specs {
					w.walkConst(sp.(*ast.ValueSpec))
				}
			case token.TYPE:
				for _, sp := range d.Specs {
					w.walkTypeSpec(sp.(*ast.TypeSpec))
				}
			case token.VAR:
				for _, sp := range d.Specs {
					w.walkVar(sp.(*ast.ValueSpec))
				}
			default:
				log.Fatalf("unknown token type %d in GenDecl", d.Tok)
			}
		case *ast.FuncDecl:
			// Ignore. Handled in subsequent pass, by go/doc.
		default:
			log.Printf("unhandled %T, %#v\n", di, di)
			printer.Fprint(os.Stderr, w.fset, di)
			os.Stderr.Write([]byte("\n"))
		}
	}
}

var constType = map[token.Token]string{
	token.INT:    "ideal-int",
	token.FLOAT:  "ideal-float",
	token.STRING: "ideal-string",
	token.CHAR:   "ideal-char",
	token.IMAG:   "ideal-imag",
}

var varType = map[token.Token]string{
	token.INT:    "int",
	token.FLOAT:  "float64",
	token.STRING: "string",
	token.CHAR:   "rune",
	token.IMAG:   "complex128",
}

var builtinTypes = []string{
	"bool", "byte", "complex64", "complex128", "error", "float32", "float64",
	"int", "int8", "int16", "int32", "int64", "rune", "string",
	"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
}

func isBuiltinType(typ string) bool {
	for _, v := range builtinTypes {
		if v == typ {
			return true
		}
	}
	return false
}

func constTypePriority(typ string) int {
	switch typ {
	case "complex128":
		return 100
	case "ideal-imag":
		return 99
	case "complex64":
		return 98
	case "float64":
		return 97
	case "ideal-float":
		return 96
	case "float32":
		return 95
	case "int64":
		return 94
	case "int", "uint", "uintptr":
		return 93
	case "ideal-int":
		return 92
	case "int16", "uint16", "int8", "uint8", "byte":
		return 91
	case "ideal-char":
		return 90
	}
	return 101
}

func (w *Walker) constRealType(typ string) string {
	pos := strings.Index(typ, ".")
	if pos >= 0 {
		pkg := typ[:pos]
		if pkg == "C" {
			return "int"
		}
		typ = typ[pos+1:]
		if p := w.findPackage(pkg); p != nil {
			ret := p.findType(typ)
			if ret != nil {
				return w.nodeString(w.namelessType(ret))
			}
		}
	} else {
		ret := w.curPackage.findType(typ)
		if ret != nil {
			return w.nodeString(w.namelessType(ret))
		}
	}
	return typ
}

func (w *Walker) constValueType(vi interface{}) (string, error) {
	switch v := vi.(type) {
	case *ast.BasicLit:
		litType, ok := constType[v.Kind]
		if !ok {
			return "", fmt.Errorf("unknown basic literal kind %#v", v)
		}
		return litType, nil
	case *ast.UnaryExpr:
		return w.constValueType(v.X)
	case *ast.SelectorExpr:
		lhs := w.nodeString(v.X)
		rhs := w.nodeString(v.Sel)
		//if CGO
		if lhs == "C" {
			return lhs + "." + rhs, nil
		}
		if p := w.findPackage(lhs); p != nil {
			if ret, ok := p.consts[rhs]; ok {
				return w.pkgRetType(p.name, ret.T), nil
			}
		}
		return "", fmt.Errorf("unknown constant reference to %s.%s", lhs, rhs)
	case *ast.Ident:
		if v.Name == "iota" {
			return "ideal-int", nil // hack.
		}
		if v.Name == "false" || v.Name == "true" {
			return "bool", nil
		}
		if t, ok := w.curPackage.consts[v.Name]; ok {
			return t.T, nil
		}
		return constDepPrefix + v.Name, nil
	case *ast.BinaryExpr:
		//== > < ! != >= <=
		if v.Op == token.EQL || v.Op == token.LSS || v.Op == token.GTR || v.Op == token.NOT ||
			v.Op == token.NEQ || v.Op == token.LEQ || v.Op == token.GEQ {
			return "bool", nil
		}
		left, err := w.constValueType(v.X)
		if err != nil {
			return "", err
		}
		if v.Op == token.SHL || v.Op == token.SHR {
			return left, err
		}
		right, err := w.constValueType(v.Y)
		if err != nil {
			return "", err
		}
		//const left != right , one or two is ideal-
		if left != right {
			if strings.HasPrefix(left, constDepPrefix) && strings.HasPrefix(right, constDepPrefix) {
				// Just pick one.
				// e.g. text/scanner GoTokens const-dependency:ScanIdents, const-dependency:ScanFloats
				return left, nil
			}
			lp := constTypePriority(w.constRealType(left))
			rp := constTypePriority(w.constRealType(right))
			if lp >= rp {
				return left, nil
			} else {
				return right, nil
			}
			return "", fmt.Errorf("in BinaryExpr, unhandled type mismatch; left=%q, right=%q", left, right)
		}
		return left, nil
	case *ast.CallExpr:
		// Not a call, but a type conversion.
		typ := w.nodeString(v.Fun)
		switch typ {
		case "complex":
			return "complex128", nil
		case "real", "imag":
			return "float64", nil
		}
		return typ, nil
	case *ast.ParenExpr:
		return w.constValueType(v.X)
	}
	return "", fmt.Errorf("unknown const value type %T", vi)
}

func (w *Walker) pkgRetType(pkg, ret string) string {
	pkg = pkg[strings.LastIndex(pkg, "/")+1:]
	if strings.HasPrefix(ret, "[]") {
		return "[]" + w.pkgRetType(pkg, ret[2:])
	}
	if strings.HasPrefix(ret, "*") {
		return "*" + w.pkgRetType(pkg, ret[1:])
	}
	if ast.IsExported(ret) {
		return pkg + "." + ret
	}
	return ret
}

func (w *Walker) findStructFieldType(st ast.Expr, name string) ast.Expr {
	_, expr := w.findStructField(st, name)
	return expr
}

func (w *Walker) findStructFieldFunction(st ast.Expr, name string) (*TypeInfo, error) {
	if s, ok := st.(*ast.StructType); ok {
		for _, fi := range s.Fields.List {
			typ := fi.Type
			if fi.Names == nil {
				switch v := typ.(type) {
				case *ast.Ident:
					if t := w.curPackage.findType(v.Name); t != nil {
						return w.lookupFunction(v.Name, name)
					}
				case *ast.SelectorExpr:
					pt := w.nodeString(typ)
					pos := strings.Index(pt, ".")
					if pos != -1 {
						if p := w.findPackage(pt[:pos]); p != nil {
							if t := p.findType(pt[pos+1:]); t != nil {
								return w.lookupFunction(pt, name)
							}
						}
					}
				case *ast.StarExpr:
					return w.findStructFieldFunction(v.X, name)
				default:
					if apiVerbose {
						log.Printf("unable to handle embedded %T", typ)
					}
				}
			}
		}
	}
	return nil, nil
}

func (w *Walker) findStructField(st ast.Expr, name string) (*ast.Ident, ast.Expr) {
	if s, ok := st.(*ast.StructType); ok {
		for _, fi := range s.Fields.List {
			typ := fi.Type
			for _, n := range fi.Names {
				if n.Name == name {
					return n, fi.Type
				}
			}
			if fi.Names == nil {
				switch v := typ.(type) {
				case *ast.Ident:
					if t := w.curPackage.findType(v.Name); t != nil {
						if v.Name == name {
							return v, v
						}
						id, expr := w.findStructField(t, name)
						if id != nil {
							return id, expr
						}
					}
				case *ast.StarExpr:
					switch vv := v.X.(type) {
					case *ast.Ident:
						if t := w.curPackage.findType(vv.Name); t != nil {
							if vv.Name == name {
								return vv, v.X
							}
							id, expr := w.findStructField(t, name)
							if id != nil {
								return id, expr
							}
						}
					case *ast.SelectorExpr:
						pt := w.nodeString(typ)
						pos := strings.Index(pt, ".")
						if pos != -1 {
							if p := w.findPackage(pt[:pos]); p != nil {
								if t := p.findType(pt[pos+1:]); t != nil {
									return w.findStructField(t, name)
								}
							}
						}
					default:
						if apiVerbose {
							log.Printf("unable to handle embedded starexpr before %T", typ)
						}
					}
				case *ast.SelectorExpr:
					pt := w.nodeString(typ)
					pos := strings.Index(pt, ".")
					if pos != -1 {
						if p := w.findPackage(pt[:pos]); p != nil {
							if t := p.findType(pt[pos+1:]); t != nil {
								return w.findStructField(t, name)
							}
						}
					}
				default:
					if apiVerbose {
						log.Printf("unable to handle embedded %T", typ)
					}
				}
			}
		}
	}
	return nil, nil
}

func (w *Walker) lookupFunction(name, sel string) (*TypeInfo, error) {
	name = strings.TrimLeft(name, "*")
	if p := w.findPackage(name); p != nil {
		fn := p.findCallFunc(sel)
		if fn != nil {
			return &TypeInfo{Kind: KindFunc, X: fn, Name: name + "." + sel, T: fn, Type: w.nodeString(w.namelessType(fn))}, nil
		}
	}
	pos := strings.Index(name, ".")
	if pos != -1 {
		pkg := name[:pos]
		typ := name[pos+1:]
		if p := w.findPackage(pkg); p != nil {
			if ident, fn := p.findMethod(typ, sel); fn != nil {
				return &TypeInfo{Kind: KindMethod, X: fn, Name: name + "." + sel, T: ident, Type: w.nodeString(w.namelessType(fn))}, nil
			}
		}
		return nil, fmt.Errorf("not lookup pkg type function pkg: %s, %s. %s. %s", name, pkg, typ, sel)
	}

	//find local var.func()
	if ns, nt, n := w.resolveName(name); n >= 0 {
		var vt string
		if nt != nil {
			vt = w.nodeString(w.namelessType(nt))
		} else if ns != nil {
			typ, err := w.varValueType(ns, n)
			if err == nil {
				vt = typ
			}
		} else {
			typ := w.curPackage.findSelectorType(name)
			if typ != nil {
				vt = w.nodeString(w.namelessType(typ))
			}
		}
		if strings.HasPrefix(vt, "*") {
			vt = vt[1:]
		}
		if vt == "error" && sel == "Error" {
			return &TypeInfo{Kind: KindBuiltin, Name: "error.Error", Type: "()string"}, nil
		}
		if fn, ok := w.curPackage.functions[vt+"."+sel]; ok {
			return &TypeInfo{Kind: KindMethod, X: fn.ft, Name: name + "." + sel, T: fn.ft, Type: w.nodeString(w.namelessType(fn))}, nil
		}
	}
	if typ, ok := w.curPackage.structs[name]; ok {
		if fn, ok := w.curPackage.functions[name+"."+sel]; ok {
			return &TypeInfo{Kind: KindMethod, X: fn.ft, Name: name + "." + sel, T: fn.ft, Type: w.nodeString(w.namelessType(fn.ft))}, nil
		}
		if info, err := w.findStructFieldFunction(typ, sel); err == nil {
			return info, nil
		}
		// struct field is type function
		if ft := w.findStructFieldType(typ, sel); ft != nil {
			typ, err := w.varValueType(ft, 0)
			if err != nil {
				typ = w.nodeString(ft)
			}
			return &TypeInfo{Kind: KindField, X: ft, Name: name + "." + sel, T: ft, Type: typ}, nil
		}
	}

	if ident, fn := w.curPackage.findMethod(name, sel); ident != nil && fn != nil {
		return &TypeInfo{Kind: KindMethod, X: fn, Name: name + "." + sel, T: ident, Type: w.nodeString(w.namelessType(fn))}, nil
	}

	if p := w.findPackage(name); p != nil {
		fn := p.findCallFunc(sel)
		if fn != nil {
			return &TypeInfo{Kind: KindFunc, X: fn, Name: name + "." + sel, T: fn, Type: w.nodeString(w.namelessType(fn))}, nil
		}
		return nil, fmt.Errorf("not find pkg func0 %v.%v", p.name, sel)
	}
	return nil, fmt.Errorf("not lookup func %v.%v", name, sel)
}

func (w *Walker) varFunctionType(name, sel string, index int) (string, error) {
	name = strings.TrimLeft(name, "*")
	pos := strings.Index(name, ".")
	if pos != -1 {
		pkg := name[:pos]
		typ := name[pos+1:]

		if p := w.findPackage(pkg); p != nil {
			_, fn := p.findMethod(typ, sel)
			if fn != nil {
				ret := funcRetType(fn, index)
				if ret != nil {
					return w.pkgRetType(p.name, w.nodeString(w.namelessType(ret))), nil
				}
			}
		}
		return "", fmt.Errorf("unknown pkg type function pkg: %s.%s.%s", pkg, typ, sel)
	}
	//find local var
	if v, ok := w.localvar[name]; ok {
		vt := v.T
		if strings.HasPrefix(vt, "*") {
			vt = vt[1:]
		}
		if vt == "error" && sel == "Error" {
			return "string", nil
		}
		typ, err := w.varFunctionType(vt, sel, 0)
		if err == nil {
			return typ, nil
		}
	}
	//find global var.func()
	if ns, nt, n := w.resolveName(name); n >= 0 {
		var vt string
		if nt != nil {
			vt = w.nodeString(w.namelessType(nt))
		} else if ns != nil {
			typ, err := w.varValueType(ns, n)
			if err == nil {
				vt = typ
			}
		} else {
			typ := w.curPackage.findSelectorType(name)
			if typ != nil {
				vt = w.nodeString(w.namelessType(typ))
			}
		}
		if strings.HasPrefix(vt, "*") {
			vt = vt[1:]
		}
		if vt == "error" && sel == "Error" {
			return "string", nil
		}
		if fn, ok := w.curPackage.functions[vt+"."+sel]; ok {
			return w.nodeString(w.namelessType(funcRetType(fn.ft, index))), nil
		}
	}
	if typ, ok := w.curPackage.structs[name]; ok {
		if ft := w.findStructFieldType(typ, sel); ft != nil {
			return w.varValueType(ft, index)
		}
	}
	//find pkg.func()
	if p := w.findPackage(name); p != nil {
		typ := p.findCallType(sel, index)
		if typ != nil {
			return w.pkgRetType(p.name, w.nodeString(w.namelessType(typ))), nil
		}
		//log.Println("->", p.functions)
		return "", fmt.Errorf("not find pkg func1 %v . %v", p.name, sel)
	}
	return "", fmt.Errorf("not find func %v.%v", name, sel)
}

func (w *Walker) lookupSelector(name string, sel string) (*TypeInfo, error) {
	name = strings.TrimLeft(name, "*")
	pos := strings.Index(name, ".")
	if pos != -1 {
		pkg := name[:pos]
		typ := name[pos+1:]
		if p := w.findPackage(pkg); p != nil {
			t := p.findType(typ)
			if t != nil {
				typ := w.findStructFieldType(t, sel)
				if typ != nil {
					return &TypeInfo{Kind: KindField, X: typ, Name: name + "." + sel, T: typ, Type: w.pkgRetType(p.name, w.nodeString(w.namelessType(typ)))}, nil
				}
			}
		}
		return nil, fmt.Errorf("lookup unknown pkg type selector pkg: %s.%s %s", pkg, typ, sel)
	}

	if lv, ok := w.localvar[name]; ok {
		return w.lookupSelector(lv.T, sel)
	}

	vs, vt, n := w.resolveName(name)
	if n >= 0 {
		var typ string
		if vt != nil {
			typ = w.nodeString(w.namelessType(vt))
		} else {
			typ, _ = w.varValueType(vs, n)
		}
		if strings.HasPrefix(typ, "*") {
			typ = typ[1:]
		}
		//typ is type, find real type
		for k, v := range w.curPackage.types {
			if k == typ {
				typ = w.nodeString(w.namelessType(v))
			}
		}
		pos := strings.Index(typ, ".")
		if pos == -1 {
			t := w.curPackage.findType(typ)
			if t != nil {
				typ := w.findStructFieldType(t, sel)
				if typ != nil {
					return &TypeInfo{Kind: KindField, X: typ, Name: name + "." + sel, T: typ, Type: w.nodeString(w.namelessType(typ))}, nil
				}
			}
		} else {
			name := typ[:pos]
			typ = typ[pos+1:]
			if p := w.findPackage(name); p != nil {
				t := p.findType(typ)
				if t != nil {
					typ := w.findStructFieldType(t, sel)
					if typ != nil {
						return &TypeInfo{Kind: KindField, X: typ, Name: name + "." + sel, T: typ, Type: w.nodeString(w.namelessType(typ))}, nil
					}
				}
			}
		}
	}
	if p := w.findPackage(name); p != nil {
		typ := p.findSelectorType(sel)
		if typ != nil {
			return &TypeInfo{Kind: KindType, X: typ, Name: name + "." + sel, T: typ, Type: w.pkgRetType(p.name, w.nodeString(w.namelessType(typ)))}, nil
		}
	}
	t := w.curPackage.findType(name)
	if t != nil {
		typ := w.findStructFieldType(t, sel)
		if typ != nil {
			return &TypeInfo{Kind: KindField, X: typ, Name: name + "." + sel, T: typ, Type: w.nodeString(w.namelessType(typ))}, nil
		}
	}
	if t, ok := w.curPackage.types[name]; ok {
		return w.lookupSelector(w.nodeString(t), sel)
	}
	return nil, fmt.Errorf("unknown selector expr ident: %s.%s", name, sel)
}

func (w *Walker) varSelectorType(name string, sel string) (string, error) {
	name = strings.TrimLeft(name, "*")
	pos := strings.Index(name, ".")
	if pos != -1 {
		pkg := name[:pos]
		typ := name[pos+1:]
		if p := w.findPackage(pkg); p != nil {
			t := p.findType(typ)
			if t != nil {
				typ := w.findStructFieldType(t, sel)
				if typ != nil {
					return w.pkgRetType(pkg, w.nodeString(w.namelessType(typ))), nil
				}
			}
		}
		return "", fmt.Errorf("unknown pkg type selector pkg: %s.%s.%s", pkg, typ, sel)
	}
	//check local
	if lv, ok := w.localvar[name]; ok {
		return w.varSelectorType(lv.T, sel)
	}
	//check struct
	if t := w.curPackage.findType(name); t != nil {
		typ := w.findStructFieldType(t, sel)
		if typ != nil {
			return w.nodeString(w.namelessType(typ)), nil
		}
	}
	//check var
	vs, vt, n := w.resolveName(name)
	if n >= 0 {
		var typ string
		if vt != nil {
			typ = w.nodeString(w.namelessType(vt))
		} else {
			typ, _ = w.varValueType(vs, n)
		}
		if strings.HasPrefix(typ, "*") {
			typ = typ[1:]
		}
		//typ is type, find real type
		for k, v := range w.curPackage.types {
			if k == typ {
				typ = w.nodeString(w.namelessType(v))
			}
		}
		pos := strings.Index(typ, ".")
		if pos == -1 {
			t := w.curPackage.findType(typ)
			if t != nil {
				typ := w.findStructFieldType(t, sel)
				if typ != nil {
					return w.nodeString(w.namelessType(typ)), nil
				}
			}
		} else {
			name := typ[:pos]
			typ = typ[pos+1:]
			if p := w.findPackage(name); p != nil {
				t := p.findType(typ)
				if t != nil {
					typ := w.findStructFieldType(t, sel)
					if typ != nil {
						return w.nodeString(w.namelessType(typ)), nil
					}
				}
			}
		}
	}

	if p := w.findPackage(name); p != nil {
		typ := p.findSelectorType(sel)
		if typ != nil {
			return w.pkgRetType(p.name, w.nodeString(w.namelessType(typ))), nil
		}
	}
	return "", fmt.Errorf("unknown var selector expr ident: %s.%s", name, sel)
}

func (w *Walker) varValueType(vi ast.Expr, index int) (string, error) {
	if vi == nil {
		return "", nil
	}
	switch v := vi.(type) {
	case *ast.BasicLit:
		litType, ok := varType[v.Kind]
		if !ok {
			return "", fmt.Errorf("unknown basic literal kind %#v", v)
		}
		return litType, nil
	case *ast.CompositeLit:
		return w.nodeString(v.Type), nil
	case *ast.FuncLit:
		return w.nodeString(w.namelessType(v.Type)), nil
	case *ast.InterfaceType:
		return w.nodeString(v), nil
	case *ast.Ellipsis:
		typ, err := w.varValueType(v.Elt, index)
		if err != nil {
			return "", err
		}
		return "[]" + typ, nil
	case *ast.StarExpr:
		typ, err := w.varValueType(v.X, index)
		if err != nil {
			return "", err
		}
		return "*" + typ, err
	case *ast.UnaryExpr:
		if v.Op == token.AND {
			typ, err := w.varValueType(v.X, index)
			return "*" + typ, err
		}
		return "", fmt.Errorf("unknown unary expr: %#v", v)
	case *ast.SelectorExpr:
		switch st := v.X.(type) {
		case *ast.Ident:
			return w.varSelectorType(st.Name, v.Sel.Name)
		case *ast.CallExpr:
			typ, err := w.varValueType(v.X, index)
			if err == nil {
				if strings.HasPrefix(typ, "*") {
					typ = typ[1:]
				}
				t := w.curPackage.findType(typ)
				if st, ok := t.(*ast.StructType); ok {
					for _, fi := range st.Fields.List {
						for _, n := range fi.Names {
							if n.Name == v.Sel.Name {
								return w.varValueType(fi.Type, index)
							}
						}
					}
				}
			}
		case *ast.SelectorExpr:
			typ, err := w.varValueType(v.X, index)
			if err == nil {
				return w.varSelectorType(typ, v.Sel.Name)
			}
		case *ast.IndexExpr:
			typ, err := w.varValueType(st.X, index)
			if err == nil {
				if strings.HasPrefix(typ, "[]") {
					return w.varSelectorType(typ[2:], v.Sel.Name)
				}
			}
		case *ast.CompositeLit:
			typ, err := w.varValueType(st.Type, 0)
			if err == nil {
				//log.Println(typ, v.Sel.Name)
				t, err := w.varSelectorType(typ, v.Sel.Name)
				if err == nil {
					return t, nil
				}
			}
		}
		return "", fmt.Errorf("var unknown selector expr: %T %s.%s", v.X, w.nodeString(v.X), v.Sel)
	case *ast.Ident:
		if v.Name == "true" || v.Name == "false" {
			return "bool", nil
		}
		if isBuiltinType(v.Name) {
			return v.Name, nil
		}
		if lv, ok := w.localvar[v.Name]; ok {
			return lv.T, nil
		}
		vt := w.curPackage.findType(v.Name)
		if vt != nil {
			if _, ok := vt.(*ast.StructType); ok {
				return v.Name, nil
			}
			return w.nodeString(vt), nil
		}
		vs, _, n := w.resolveName(v.Name)
		if n >= 0 {
			return w.varValueType(vs, n)
		}
		return "", fmt.Errorf("unresolved identifier: %q", v.Name)
	case *ast.BinaryExpr:
		//== > < ! != >= <=
		if v.Op == token.EQL || v.Op == token.LSS || v.Op == token.GTR || v.Op == token.NOT ||
			v.Op == token.NEQ || v.Op == token.LEQ || v.Op == token.GEQ {
			return "bool", nil
		}
		left, err := w.varValueType(v.X, index)
		if err != nil {
			return "", err
		}
		right, err := w.varValueType(v.Y, index)
		if err != nil {
			return "", err
		}
		if left != right {
			return "", fmt.Errorf("in BinaryExpr, unhandled type mismatch; left=%q, right=%q", left, right)
		}
		return left, nil
	case *ast.ParenExpr:
		return w.varValueType(v.X, index)
	case *ast.CallExpr:
		switch ft := v.Fun.(type) {
		case *ast.ArrayType:
			return w.nodeString(v.Fun), nil
		case *ast.Ident:
			switch ft.Name {
			case "make":
				return w.nodeString(w.namelessType(v.Args[0])), nil
			case "new":
				return "*" + w.nodeString(w.namelessType(v.Args[0])), nil
			case "append":
				return w.varValueType(v.Args[0], 0)
			case "recover":
				return "interface{}", nil
			case "len", "cap", "copy":
				return "int", nil
			case "complex":
				return "complex128", nil
			case "real":
				return "float64", nil
			case "imag":
				return "float64", nil
			}
			if isBuiltinType(ft.Name) {
				return ft.Name, nil
			}
			typ := w.curPackage.findCallType(ft.Name, index)
			if typ != nil {
				return w.nodeString(w.namelessType(typ)), nil
			}
			//if local var type
			if fn, ok := w.localvar[ft.Name]; ok {
				typ := fn.T
				if strings.HasPrefix(typ, "func(") {
					expr, err := parser.ParseExpr(typ + "{}")
					if err == nil {
						if fl, ok := expr.(*ast.FuncLit); ok {
							retType := funcRetType(fl.Type, index)
							if retType != nil {
								return w.nodeString(w.namelessType(retType)), nil
							}
						}
					}
				}
			}
			//if var is func() type
			vs, _, n := w.resolveName(ft.Name)
			if n >= 0 {
				if vs != nil {
					typ, err := w.varValueType(vs, n)
					if err == nil {
						if strings.HasPrefix(typ, "func(") {
							expr, err := parser.ParseExpr(typ + "{}")
							if err == nil {
								if fl, ok := expr.(*ast.FuncLit); ok {
									retType := funcRetType(fl.Type, index)
									if retType != nil {
										return w.nodeString(w.namelessType(retType)), nil
									}
								}
							}
						}
					}
				}
			}
			return "", fmt.Errorf("unknown funcion %s %s", w.curPackageName, ft.Name)
		case *ast.SelectorExpr:
			typ, err := w.varValueType(ft.X, index)
			if err == nil {
				if strings.HasPrefix(typ, "*") {
					typ = typ[1:]
				}
				retType := w.curPackage.findCallType(typ+"."+ft.Sel.Name, index)
				if retType != nil {
					return w.nodeString(w.namelessType(retType)), nil
				}
			}
			switch st := ft.X.(type) {
			case *ast.Ident:
				return w.varFunctionType(st.Name, ft.Sel.Name, index)
			case *ast.CallExpr:
				typ, err := w.varValueType(st, 0)
				if err != nil {
					return "", err
				}
				return w.varFunctionType(typ, ft.Sel.Name, index)
			case *ast.SelectorExpr:
				typ, err := w.varValueType(st, index)
				if err == nil {
					return w.varFunctionType(typ, ft.Sel.Name, index)
				}
			case *ast.IndexExpr:
				typ, err := w.varValueType(st.X, index)
				if err == nil {
					if strings.HasPrefix(typ, "[]") {
						return w.varFunctionType(typ[2:], ft.Sel.Name, index)
					}
				}
			case *ast.TypeAssertExpr:
				typ := w.nodeString(w.namelessType(st.Type))
				typ = strings.TrimLeft(typ, "*")
				return w.varFunctionType(typ, ft.Sel.Name, index)
			}
			return "", fmt.Errorf("unknown var function selector %v %T", w.nodeString(ft.X), ft.X)
		case *ast.FuncLit:
			retType := funcRetType(ft.Type, index)
			if retType != nil {
				return w.nodeString(w.namelessType(retType)), nil
			}
		case *ast.CallExpr:
			typ, err := w.varValueType(v.Fun, 0)
			if err == nil && strings.HasPrefix(typ, "func(") {
				expr, err := parser.ParseExpr(typ + "{}")
				if err == nil {
					if fl, ok := expr.(*ast.FuncLit); ok {
						retType := funcRetType(fl.Type, index)
						if retType != nil {
							return w.nodeString(w.namelessType(retType)), nil
						}
					}
				}
			}
		}
		return "", fmt.Errorf("not a known function %T %v", v.Fun, w.nodeString(v.Fun))
	case *ast.MapType:
		return fmt.Sprintf("map[%s](%s)", w.nodeString(w.namelessType(v.Key)), w.nodeString(w.namelessType(v.Value))), nil
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", w.nodeString(w.namelessType(v.Elt))), nil
	case *ast.FuncType:
		return w.nodeString(w.namelessType(v)), nil
	case *ast.IndexExpr:
		typ, err := w.varValueType(v.X, index)
		typ = strings.TrimLeft(typ, "*")
		if err == nil {
			if index == 0 {
				return typ, nil
			} else if index == 1 {
				return "bool", nil
			}
			if strings.HasPrefix(typ, "[]") {
				return typ[2:], nil
			} else if strings.HasPrefix(typ, "map[") {
				node, err := parser.ParseExpr(typ + "{}")
				if err == nil {
					if cl, ok := node.(*ast.CompositeLit); ok {
						if m, ok := cl.Type.(*ast.MapType); ok {
							return w.nodeString(w.namelessType(m.Value)), nil
						}
					}
				}
			}
		}
		return "", fmt.Errorf("unknown index %v %v %v %v", typ, v.X, index, err)
	case *ast.SliceExpr:
		return w.varValueType(v.X, index)
	case *ast.ChanType:
		typ, err := w.varValueType(v.Value, index)
		if err == nil {
			if v.Dir == ast.RECV {
				return "<-chan " + typ, nil
			} else if v.Dir == ast.SEND {
				return "chan<- " + typ, nil
			}
			return "chan " + typ, nil
		}
	case *ast.TypeAssertExpr:
		if index == 1 {
			return "bool", nil
		}
		return w.nodeString(w.namelessType(v.Type)), nil
	default:
		return "", fmt.Errorf("unknown value type %v %T", w.nodeString(vi), vi)
	}
	//panic("unreachable")
	return "", fmt.Errorf("unreachable value type %v %T", vi, vi)
}

// resolveName finds a top-level node named name and returns the node
// v and its type t, if known.
func (w *Walker) resolveName(name string) (v ast.Expr, t interface{}, n int) {
	for _, file := range w.curPackage.apkg.Files {
		for _, di := range file.Decls {
			switch d := di.(type) {
			case *ast.GenDecl:
				switch d.Tok {
				case token.VAR:
					for _, sp := range d.Specs {
						vs := sp.(*ast.ValueSpec)
						for i, vname := range vs.Names {
							if vname.Name == name {
								if len(vs.Values) == 1 {
									return vs.Values[0], vs.Type, i
								}
								return nil, vs.Type, i
							}
						}
					}
				}
			}
		}
	}
	return nil, nil, -1
}

// constDepPrefix is a magic prefix that is used by constValueType
// and walkConst to signal that a type isn't known yet. These are
// resolved at the end of walking of a package's files.
const constDepPrefix = "const-dependency:"

func (w *Walker) walkConst(vs *ast.ValueSpec) {
	for _, ident := range vs.Names {
		if !w.isExtract(ident.Name) {
			continue
		}
		litType := ""
		if vs.Type != nil {
			litType = w.nodeString(vs.Type)
		} else {
			litType = w.lastConstType
			if vs.Values != nil {
				if len(vs.Values) != 1 {
					log.Fatalf("const %q, values: %#v", ident.Name, vs.Values)
				}
				var err error
				litType, err = w.constValueType(vs.Values[0])
				if err != nil {
					if apiVerbose {
						log.Printf("unknown kind in const %q (%T): %v", ident.Name, vs.Values[0], err)
					}
					litType = "unknown-type"
				}
			}
		}
		if strings.HasPrefix(litType, constDepPrefix) {
			dep := litType[len(constDepPrefix):]
			w.constDep[ident.Name] = &ExprType{T: dep, X: ident}
			continue
		}
		if litType == "" {
			if apiVerbose {
				log.Printf("unknown kind in const %q", ident.Name)
			}
			continue
		}
		w.lastConstType = litType

		w.curPackage.consts[ident.Name] = &ExprType{T: litType, X: ident}

		if isExtract(ident.Name) {
			w.emitFeature(fmt.Sprintf("const %s %s", ident, litType), ident.Pos())
		}
	}
}

func (w *Walker) resolveConstantDeps() {
	var findConstType func(string) string
	findConstType = func(ident string) string {
		if dep, ok := w.constDep[ident]; ok {
			return findConstType(dep.T)
		}
		if t, ok := w.curPackage.consts[ident]; ok {
			return t.T
		}
		return ""
	}
	for ident, info := range w.constDep {
		if !isExtract(ident) {
			continue
		}
		t := findConstType(ident)
		if t == "" {
			if apiVerbose {
				log.Printf("failed to resolve constant %q", ident)
			}
			continue
		}
		w.curPackage.consts[ident] = &ExprType{T: t, X: info.X}
		w.emitFeature(fmt.Sprintf("const %s %s", ident, t), info.X.Pos())
	}
}

func (w *Walker) walkVar(vs *ast.ValueSpec) {
	if vs.Type != nil {
		typ := w.nodeString(vs.Type)
		for _, ident := range vs.Names {
			w.curPackage.vars[ident.Name] = &ExprType{T: typ, X: ident}
			if isExtract(ident.Name) {
				w.emitFeature(fmt.Sprintf("var %s %s", ident, typ), ident.Pos())
			}
		}
	} else if len(vs.Names) == len(vs.Values) {
		for n, ident := range vs.Names {
			if !w.isExtract(ident.Name) {
				continue
			}
			typ, err := w.varValueType(vs.Values[n], n)
			if err != nil {
				if apiVerbose {
					log.Printf("unknown type of variable0 %q, type %T, error = %v, pos=%s",
						ident.Name, vs.Values[n], err, w.fset.Position(vs.Pos()))
				}
				typ = "unknown-type"
			}
			w.curPackage.vars[ident.Name] = &ExprType{T: typ, X: ident}
			if isExtract(ident.Name) {
				w.emitFeature(fmt.Sprintf("var %s %s", ident, typ), ident.Pos())
			}
		}
	} else if len(vs.Values) == 1 {
		for n, ident := range vs.Names {
			if !w.isExtract(ident.Name) {
				continue
			}
			typ, err := w.varValueType(vs.Values[0], n)
			if err != nil {
				if apiVerbose {
					log.Printf("unknown type of variable1 %q, type %T, error = %v, pos=%s",
						ident.Name, vs.Values[0], err, w.fset.Position(vs.Pos()))
				}
				typ = "unknown-type"
			}
			w.curPackage.vars[ident.Name] = &ExprType{T: typ, X: ident}
			if isExtract(ident.Name) {
				w.emitFeature(fmt.Sprintf("var %s %s", ident, typ), ident.Pos())
			}
		}
	}
}

func (w *Walker) nodeString(node interface{}) string {
	if node == nil {
		return ""
	}
	var b bytes.Buffer
	printer.Fprint(&b, w.fset, node)
	return b.String()
}

func (w *Walker) nodeDebug(node interface{}) string {
	if node == nil {
		return ""
	}
	var b bytes.Buffer
	ast.Fprint(&b, w.fset, node, nil)
	return b.String()
}

func (w *Walker) noteInterface(name string, it *ast.InterfaceType) {
	w.interfaces[pkgSymbol{w.curPackageName, name}] = it
}

func (w *Walker) walkTypeSpec(ts *ast.TypeSpec) {
	name := ts.Name.Name
	if !isExtract(name) {
		return
	}
	switch t := ts.Type.(type) {
	case *ast.StructType:
		w.walkStructType(name, t)
	case *ast.InterfaceType:
		w.walkInterfaceType(name, t)
	default:
		w.emitFeature(fmt.Sprintf("type %s %s", name, w.nodeString(ts.Type)), t.Pos()-token.Pos(len(name)+1))
	}
}

func (w *Walker) walkStructType(name string, t *ast.StructType) {
	typeStruct := fmt.Sprintf("type %s struct", name)
	w.emitFeature(typeStruct, t.Pos()-token.Pos(len(name)+1))
	pop := w.pushScope(typeStruct)
	defer pop()
	for _, f := range t.Fields.List {
		typ := f.Type
		for _, name := range f.Names {
			if isExtract(name.Name) {
				w.emitFeature(fmt.Sprintf("%s %s", name, w.nodeString(w.namelessType(typ))), name.Pos())
			}
		}
		if f.Names == nil {
			switch v := typ.(type) {
			case *ast.Ident:
				if isExtract(v.Name) {
					w.emitFeature(fmt.Sprintf("embedded %s", v.Name), v.Pos())
				}
			case *ast.StarExpr:
				switch vv := v.X.(type) {
				case *ast.Ident:
					if isExtract(vv.Name) {
						w.emitFeature(fmt.Sprintf("embedded *%s", vv.Name), vv.Pos())
					}
				case *ast.SelectorExpr:
					w.emitFeature(fmt.Sprintf("embedded %s", w.nodeString(typ)), v.Pos())
				default:
					log.Fatalf("unable to handle embedded starexpr before %T", typ)
				}
			case *ast.SelectorExpr:
				w.emitFeature(fmt.Sprintf("embedded %s", w.nodeString(typ)), v.Pos())
			default:
				if apiVerbose {
					log.Printf("unable to handle embedded %T", typ)
				}
			}
		}
	}
}

// typeMethod is a method of an interface.
type typeMethod struct {
	name string // "Read"
	sig  string // "([]byte) (int, error)", from funcSigString
	ft   *ast.FuncType
	pos  token.Pos
	recv ast.Expr
}

// interfaceMethods returns the expanded list of exported methods for an interface.
// The boolean complete reports whether the list contains all methods (that is, the
// interface has no unexported methods).
// pkg is the complete package name ("net/http")
// iname is the interface name.
func (w *Walker) interfaceMethods(pkg, iname string) (methods []typeMethod, complete bool) {
	t, ok := w.interfaces[pkgSymbol{pkg, iname}]
	if !ok {
		if apiVerbose {
			log.Printf("failed to find interface %s.%s", pkg, iname)
		}
		return
	}

	complete = true
	for _, f := range t.Methods.List {
		typ := f.Type
		switch tv := typ.(type) {
		case *ast.FuncType:
			for _, mname := range f.Names {
				if isExtract(mname.Name) {
					ft := typ.(*ast.FuncType)
					methods = append(methods, typeMethod{
						name: mname.Name,
						sig:  w.funcSigString(ft),
						ft:   ft,
						pos:  f.Pos(),
					})
				} else {
					complete = false
				}
			}
		case *ast.Ident:
			embedded := typ.(*ast.Ident).Name
			if embedded == "error" {
				methods = append(methods, typeMethod{
					name: "Error",
					sig:  "() string",
					ft: &ast.FuncType{
						Params: nil,
						Results: &ast.FieldList{
							List: []*ast.Field{
								&ast.Field{
									Type: &ast.Ident{
										Name: "string",
									},
								},
							},
						},
					},
					pos: f.Pos(),
				})
				continue
			}
			if !isExtract(embedded) {
				log.Fatalf("unexported embedded interface %q in exported interface %s.%s; confused",
					embedded, pkg, iname)
			}
			m, c := w.interfaceMethods(pkg, embedded)
			methods = append(methods, m...)
			complete = complete && c
		case *ast.SelectorExpr:
			lhs := w.nodeString(tv.X)
			rhs := w.nodeString(tv.Sel)
			fpkg, ok := w.selectorFullPkg[lhs]
			if !ok {
				log.Fatalf("can't resolve selector %q in interface %s.%s", lhs, pkg, iname)
			}
			m, c := w.interfaceMethods(fpkg, rhs)
			methods = append(methods, m...)
			complete = complete && c
		default:
			log.Fatalf("unknown type %T in interface field", typ)
		}
	}
	return
}

func (w *Walker) walkInterfaceType(name string, t *ast.InterfaceType) {
	methNames := []string{}
	pop := w.pushScope("type " + name + " interface")
	methods, complete := w.interfaceMethods(w.curPackageName, name)
	w.packageMap[w.curPackageName].interfaceMethods[name] = methods
	for _, m := range methods {
		methNames = append(methNames, m.name)
		w.emitFeature(fmt.Sprintf("%s%s", m.name, m.sig), m.pos)
	}
	if !complete {
		// The method set has unexported methods, so all the
		// implementations are provided by the same package,
		// so the method set can be extended. Instead of recording
		// the full set of names (below), record only that there were
		// unexported methods. (If the interface shrinks, we will notice
		// because a method signature emitted during the last loop,
		// will disappear.)
		w.emitFeature("unexported methods", 0)
	}
	pop()

	if !complete {
		return
	}

	sort.Strings(methNames)
	if len(methNames) == 0 {
		w.emitFeature(fmt.Sprintf("type %s interface {}", name), t.Pos()-token.Pos(len(name)+1))
	} else {
		w.emitFeature(fmt.Sprintf("type %s interface { %s }", name, strings.Join(methNames, ", ")), t.Pos()-token.Pos(len(name)+1))
	}
}

func baseTypeName(x ast.Expr) (name string, imported bool) {
	switch t := x.(type) {
	case *ast.Ident:
		return t.Name, false
	case *ast.SelectorExpr:
		if _, ok := t.X.(*ast.Ident); ok {
			// only possible for qualified type names;
			// assume type is imported
			return t.Sel.Name, true
		}
	case *ast.StarExpr:
		return baseTypeName(t.X)
	}
	return
}

func (w *Walker) peekFuncDecl(f *ast.FuncDecl) {
	var fname = f.Name.Name
	var recv ast.Expr
	if f.Recv != nil {
		recvTypeName, imp := baseTypeName(f.Recv.List[0].Type)
		if imp {
			return
		}
		fname = recvTypeName + "." + f.Name.Name
		recv = f.Recv.List[0].Type
	}
	// Record return type for later use.
	//if f.Type.Results != nil && len(f.Type.Results.List) >= 1 {
	// record all function
	w.curPackage.functions[fname] = typeMethod{
		name: fname,
		sig:  w.funcSigString(f.Type),
		ft:   f.Type,
		pos:  f.Pos(),
		recv: recv,
	}
	//}
}

func (w *Walker) walkFuncDecl(f *ast.FuncDecl) {
	if !w.isExtract(f.Name.Name) {
		return
	}
	if f.Recv != nil {
		// Method.
		recvType := w.nodeString(f.Recv.List[0].Type)
		keep := isExtract(recvType) ||
			(strings.HasPrefix(recvType, "*") &&
				isExtract(recvType[1:]))
		if !keep {
			return
		}
		w.emitFeature(fmt.Sprintf("method (%s) %s%s", recvType, f.Name.Name, w.funcSigString(f.Type)), f.Name.Pos())
		return
	}
	// Else, a function
	w.emitFeature(fmt.Sprintf("func %s%s", f.Name.Name, w.funcSigString(f.Type)), f.Name.Pos())
}

func (w *Walker) funcSigString(ft *ast.FuncType) string {
	var b bytes.Buffer
	writeField := func(b *bytes.Buffer, f *ast.Field) {
		if n := len(f.Names); n > 1 {
			for i := 0; i < n; i++ {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(w.nodeString(w.namelessType(f.Type)))
			}
		} else {
			b.WriteString(w.nodeString(w.namelessType(f.Type)))
		}
	}
	b.WriteByte('(')
	if ft.Params != nil {
		for i, f := range ft.Params.List {
			if i > 0 {
				b.WriteString(", ")
			}
			writeField(&b, f)
		}
	}
	b.WriteByte(')')
	if ft.Results != nil {
		nr := 0
		for _, f := range ft.Results.List {
			if n := len(f.Names); n > 1 {
				nr += n
			} else {
				nr++
			}
		}
		if nr > 0 {
			b.WriteByte(' ')
			if nr > 1 {
				b.WriteByte('(')
			}
			for i, f := range ft.Results.List {
				if i > 0 {
					b.WriteString(", ")
				}
				writeField(&b, f)
			}
			if nr > 1 {
				b.WriteByte(')')
			}
		}
	}
	return b.String()
}

// namelessType returns a type node that lacks any variable names.
func (w *Walker) namelessType(t interface{}) interface{} {
	ft, ok := t.(*ast.FuncType)
	if !ok {
		return t
	}
	return &ast.FuncType{
		Params:  w.namelessFieldList(ft.Params),
		Results: w.namelessFieldList(ft.Results),
	}
}

// namelessFieldList returns a deep clone of fl, with the cloned fields
// lacking names.
func (w *Walker) namelessFieldList(fl *ast.FieldList) *ast.FieldList {
	fl2 := &ast.FieldList{}
	if fl != nil {
		for _, f := range fl.List {
			n := len(f.Names)
			if n >= 1 {
				for i := 0; i < n; i++ {
					fl2.List = append(fl2.List, w.namelessField(f))
				}
			} else {
				fl2.List = append(fl2.List, w.namelessField(f))
			}
		}
	}
	return fl2
}

// namelessField clones f, but not preserving the names of fields.
// (comments and tags are also ignored)
func (w *Walker) namelessField(f *ast.Field) *ast.Field {
	return &ast.Field{
		Type: f.Type,
	}
}

func (w *Walker) emitFeature(feature string, pos token.Pos) {
	if !w.wantedPkg[w.curPackage.name] {
		return
	}
	more := strings.Index(feature, "\n")
	if more != -1 {
		if len(feature) <= 1024 {
			feature = strings.Replace(feature, "\n", " ", 1)
			feature = strings.Replace(feature, "\n", ";", -1)
			feature = strings.Replace(feature, "\t", " ", -1)
		} else {
			feature = feature[:more] + " ...more"
			if apiVerbose {
				log.Printf("feature contains newlines: %v, %s", feature, w.fset.Position(pos))
			}
		}
	}
	f := strings.Join(w.scope, w.sep) + w.sep + feature

	if _, dup := w.curPackage.features[f]; dup {
		return
	}
	w.curPackage.features[f] = pos
}

func strListContains(l []string, s string) bool {
	for _, v := range l {
		if v == s {
			return true
		}
	}
	return false
}

const goosList = "darwin freebsd linux netbsd openbsd plan9 windows "
const goarchList = "386 amd64 arm "

// goodOSArchFile returns false if the name contains a $GOOS or $GOARCH
// suffix which does not match the current system.
// The recognized name formats are:
//
//     name_$(GOOS).*
//     name_$(GOARCH).*
//     name_$(GOOS)_$(GOARCH).*
//     name_$(GOOS)_test.*
//     name_$(GOARCH)_test.*
//     name_$(GOOS)_$(GOARCH)_test.*
//
func isOSArchFile(ctxt *build.Context, name string) bool {
	if dot := strings.Index(name, "."); dot != -1 {
		name = name[:dot]
	}
	l := strings.Split(name, "_")
	if n := len(l); n > 0 && l[n-1] == "test" {
		l = l[:n-1]
	}
	n := len(l)
	if n >= 2 && knownOS[l[n-2]] && knownArch[l[n-1]] {
		return l[n-2] == ctxt.GOOS && l[n-1] == ctxt.GOARCH
	}
	if n >= 1 && knownOS[l[n-1]] {
		return l[n-1] == ctxt.GOOS
	}
	if n >= 1 && knownArch[l[n-1]] {
		return l[n-1] == ctxt.GOARCH
	}
	return false
}

var knownOS = make(map[string]bool)
var knownArch = make(map[string]bool)

func init() {
	for _, v := range strings.Fields(goosList) {
		knownOS[v] = true
	}
	for _, v := range strings.Fields(goarchList) {
		knownArch[v] = true
	}
}
