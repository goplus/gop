// Copyright 2011-2016 visualfc <visualfc@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkgwalk

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/visualfc/gotools/pkgutil"
	"github.com/visualfc/gotools/stdlib"
)

var (
	typesVerbose     bool
	typesAllowBinary bool
)

var (
	typeParserComment bool = true
	typesFindUse      bool
	typesFindDef      bool
	typesFindUseAll   bool
	typesFindInfo     bool
	typesFindDoc      bool
)

type ObjKind int

const (
	ObjNone ObjKind = iota
	ObjPackage
	ObjPackageX
	ObjPkgName //import
	ObjTypeName
	ObjInterface
	ObjStruct
	ObjConst
	ObjVar
	ObjField
	ObjFunc
	ObjMethod
	ObjLabel
	ObjBuiltin
	ObjNil
	ObjImplicit
	ObjUnknown
	ObjComment
)

var ObjKindName = []string{"none", "package", "packagex", "import",
	"type", "interface", "struct",
	"const", "var", "field",
	"func", "method",
	"label", "builtin", "nil", "implicit", "unknown"}

func (k ObjKind) String() string {
	if k >= 0 && int(k) < len(ObjKindName) {
		return ObjKindName[k]
	}
	return "unkwnown"
}

var builtinInfoMap = map[string]string{
	"append":   "func append(slice []Type, elems ...Type) []Type",
	"copy":     "func copy(dst, src []Type) int",
	"delete":   "func delete(m map[Type]Type1, key Type)",
	"len":      "func len(v Type) int",
	"cap":      "func cap(v Type) int",
	"make":     "func make(Type, size IntegerType) Type",
	"new":      "func new(Type) *Type",
	"complex":  "func complex(r, i FloatType) ComplexType",
	"real":     "func real(c ComplexType) FloatType",
	"imag":     "func imag(c ComplexType) FloatType",
	"close":    "func close(c chan<- Type)",
	"panic":    "func panic(v interface{})",
	"recover":  "func recover() interface{}",
	"print":    "func print(args ...Type)",
	"println":  "func println(args ...Type)",
	"error":    "type error interface {Error() string}",
	"Sizeof":   "func unsafe.Sizeof(any) uintptr",
	"Offsetof": "func unsafe.Offsetof(any) uintptr",
	"Alignof":  "func unsafe.Alignof(any) uintptr",
}

func builtinInfo(id string) string {
	if info, ok := builtinInfoMap[id]; ok {
		return "builtin " + info
	}
	return "builtin " + id
}

func simpleObjInfo(obj types.Object) string {
	s := obj.String()
	pkg := obj.Pkg()
	if pkg != nil {
		s = strings.Replace(s, pkg.Path(), pkg.Name(), -1)
		s = simpleType(s)
		if pkg.Name() == "main" {
			s = strings.Replace(s, "main.", "", -1)
		}
	}
	return s
}

func simpleType(src string) string {
	re, _ := regexp.Compile("[\\w\\./]+")
	return re.ReplaceAllStringFunc(src, func(s string) string {
		r := s
		if i := strings.LastIndex(s, "/"); i != -1 {
			r = s[i+1:]
		}
		if strings.Count(r, ".") > 1 {
			r = r[strings.Index(r, ".")+1:]
		}
		return r
	})
}

func DefaultPkgConfig() *PkgConfig {
	conf := &PkgConfig{IgnoreFuncBodies: true, AllowBinary: true, WithTestFiles: true}
	conf.IgnoreFuncBodies = false
	conf.Info = &types.Info{
		Uses:       make(map[*ast.Ident]types.Object),
		Defs:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Scopes:     make(map[ast.Node]*types.Scope),
		Implicits:  make(map[ast.Node]types.Object),
	}
	conf.XInfo = &types.Info{
		Uses:       make(map[*ast.Ident]types.Object),
		Defs:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Implicits:  make(map[ast.Node]types.Object),
	}
	return conf
}

func ImportPackage(w *PkgWalker, pkgPath string, conf *PkgConfig) (*types.Package, error) {
	return w.Import("", pkgPath, conf)
}

func LookupObjList(w *PkgWalker, pkg *types.Package, conf *PkgConfig) []*ObjList {
	var objList []*ObjList
	for ident, obj := range conf.Info.Defs {
		if ident.Name == "_" {
			continue
		}
		kind, err := ParserObjectKind(ident, obj, conf)
		// ObjUndefined, v := switch xxx.(type)
		if err != nil && kind == ObjUnknown {
			for _, v := range conf.Info.Implicits {
				if v.Pos() == ident.Pos() {
					kind = ObjImplicit
					break
				}
			}
		}

		//		if kind == ObjUnknown {
		//			log.Println("warning, %v, %v %v\n", err, w.fset.Position(ident.Pos()), kind)
		//			continue
		//		}

		//create new object list
		ol := &ObjList{Ident: ident, Kind: kind, Obj: obj}
		ol.List = append(ol.List, int(ident.Pos()))

		switch kind {
		case ObjPackage, ObjPackageX, ObjPkgName:
		case ObjConst, ObjVar:
		case ObjFunc:
			p := w.FileSet.Position(ident.Pos())
			if strings.HasSuffix(filepath.Base(p.Filename), "_test.go") &&
				(strings.HasPrefix(ident.Name, "Test") ||
					strings.HasPrefix(ident.Name, "Benchmark") ||
					strings.HasPrefix(ident.Name, "Example")) {
				ol.IsTestFunc = true
			}
		case ObjInterface, ObjStruct, ObjTypeName:
		case ObjField:
			if v, ok := obj.(*types.Var); ok {
				if v.Anonymous() {
					ol.IsAnonymous = true
					ol.List = make([]int, 0)
				}
			}
		case ObjMethod:
		case ObjImplicit:
		case ObjLabel:
		default:
			log.Fatalf("unprocess ident name:%v, obj:%v, kind:%v, %v", ident.Name, obj, kind, w.FileSet.Position(ident.Pos()))
		}
		//find all uses
		if obj == nil {
			for k1, v1 := range conf.Info.Uses {
				if v1.Pos() == ident.Pos() {
					ol.List = append(ol.List, int(k1.Pos()))
				}
			}
		} else {
			for k1, v1 := range conf.Info.Uses {
				if v1 == obj {
					ol.List = append(ol.List, int(k1.Pos()))
				}
			}
		}
		objList = append(objList, ol)
	}
	return objList
}

type FileCursor struct {
	pkg       string
	fileName  string
	fileDir   string
	cursorPos int
	pos       token.Pos
	src       interface{}
	xtest     bool
}

type PkgConfig struct {
	IgnoreFuncBodies bool
	AllowBinary      bool
	WithTestFiles    bool
	Cursor           *FileCursor
	Pkg              *types.Package
	XPkg             *types.Package
	Info             *types.Info
	XInfo            *types.Info
	Files            []*ast.File
	TestFiles        []*ast.File
	XTestFiles       []*ast.File
	Dir              string
}

func NewPkgWalker(context *build.Context) *PkgWalker {
	return &PkgWalker{
		Context:         context,
		FileSet:         token.NewFileSet(),
		parsedFileCache: map[string]*ast.File{},
		importingName:   map[string]bool{},
		imported:        map[string]*types.Package{"unsafe": types.Unsafe},
		gcimported:      importer.Default(),
	}
}

type PkgWalker struct {
	FileSet         *token.FileSet
	Context         *build.Context
	current         *types.Package
	importingName   map[string]bool
	parsedFileCache map[string]*ast.File
	imported        map[string]*types.Package // packages already imported
	gcimported      types.Importer
	//importing       types.Package
}

func findDecl(fset *token.FileSet, file *ast.File, pos token.Pos) ast.Decl {
	for _, decl := range file.Decls {
		if pos >= decl.Pos() && pos <= decl.End() {
			return decl
		}
	}
	return nil
}

func (w *PkgWalker) FindDeclForPos(p token.Pos) (*ast.File, ast.Decl) {
	pos := w.FileSet.Position(p)
	if !pos.IsValid() {
		return nil, nil
	}
	dir, file := filepath.Split(pos.Filename)
	f, err := w.parseFile(dir, file, nil)
	if err != nil {
		return nil, nil
	}
	return f, findDecl(w.FileSet, f, p)
}

func importName(is *ast.ImportSpec) string {
	fpath, _ := strconv.Unquote(is.Path.Value)

	fbase := fpath
	pos := strings.LastIndexAny(fpath, "./-\\")
	if pos != -1 {
		fbase = fpath[pos+1:]
	}

	var fname string
	if is.Name != nil {
		fname = is.Name.Name
	} else {
		fname = fbase
	}
	return fname
}

func (w *PkgWalker) FindImportName(ast *ast.File, path string) string {
	for _, im := range ast.Imports {
		ipath, _ := strconv.Unquote(im.Path.Value)
		if ipath == path {
			return importName(im)
		}
	}
	return ""
}

func (w *PkgWalker) FindImportEndPos(ast *ast.File) token.Pos {
	size := len(ast.Imports)
	if size == 0 {
		return token.NoPos
	}
	return ast.Imports[size-1].End()
}

func contains(list []string, s string) bool {
	for _, t := range list {
		if t == s {
			return true
		}
	}
	return false
}

func (w *PkgWalker) isBinaryPkg(pkg string) bool {
	return stdlib.IsStdPkg(pkg)
}

func (w *PkgWalker) importPath(path string, mode build.ImportMode) (*build.Package, error) {
	if filepath.IsAbs(path) {
		return w.Context.ImportDir(path, 0)
	}
	if stdlib.IsStdPkg(path) {
		return stdlib.ImportStdPkg(w.Context, path, mode)
	}
	return w.Context.Import(path, "", mode)
}

func (w *PkgWalker) Import(parentDir string, name string, conf *PkgConfig) (pkg *types.Package, err error) {
	defer func() {
		err := recover()
		if err != nil && typesVerbose {
			log.Println(err)
		}
	}()

	if parentDir != "" {
		if strings.HasPrefix(name, ".") {
			name = filepath.Join(parentDir, name)
		} else if pkgutil.IsVendorExperiment() {
			parentPkg := pkgutil.ImportDir(parentDir)
			name, _ = pkgutil.VendoredImportPath(parentPkg, name)
		}
	}

	pkg = w.imported[name]
	if pkg != nil {
		//		if pkg == &w.importing {
		//			return nil, fmt.Errorf("cycle importing package %q", name)
		//		}
		return pkg, nil
	}

	if typesVerbose {
		log.Println("parser pkg", parentDir, name)
	}

	bp, err := w.importPath(name, 0)
	if err != nil {
		return nil, err
	}

	checkName := name

	if bp.ImportPath == "." {
		checkName = bp.Name
	} else {
		checkName = bp.ImportPath
	}

	if w.importingName[checkName] {
		return nil, fmt.Errorf("cycle importing package %q", name)
	}

	w.importingName[checkName] = true

	//	if err != nil {
	//		return nil, err
	//		//if _, nogo := err.(*build.NoGoError); nogo {
	//		//	return
	//		//}
	//		//return
	//		//log.Fatalf("pkg %q, dir %q: ScanDir: %v", name, info.Dir, err)
	//	}

	filenames := append(append([]string{}, bp.GoFiles...), bp.CgoFiles...)
	if conf.WithTestFiles {
		filenames = append(filenames, bp.TestGoFiles...)
	}

	if name == "runtime" {
		n := fmt.Sprintf("zgoos_%s.go", w.Context.GOOS)
		if !contains(filenames, n) {
			filenames = append(filenames, n)
		}

		n = fmt.Sprintf("zgoarch_%s.go", w.Context.GOARCH)
		if !contains(filenames, n) {
			filenames = append(filenames, n)
		}
	}

	parserFiles := func(filenames []string, cursor *FileCursor, xtest bool) (files []*ast.File) {
		for _, file := range filenames {
			var f *ast.File
			if cursor != nil && cursor.fileName == file {
				f, err = w.parseFile(bp.Dir, file, cursor.src)
				cursor.pos = token.Pos(w.FileSet.File(f.Pos()).Base()) + token.Pos(cursor.cursorPos)
				cursor.fileDir = bp.Dir
				cursor.xtest = xtest
			} else {
				f, err = w.parseFile(bp.Dir, file, nil)
			}
			if err != nil && typesVerbose {
				log.Printf("error parsing package %s: %s\n", name, err)
			}
			files = append(files, f)
		}
		return
	}
	files := parserFiles(filenames, conf.Cursor, false)
	xfiles := parserFiles(bp.XTestGoFiles, conf.Cursor, true)

	conf.Files = files
	conf.XTestFiles = xfiles
	conf.Dir = bp.Dir

	typesConf := types.Config{
		IgnoreFuncBodies: conf.IgnoreFuncBodies,
		FakeImportC:      true,
		Importer:         &Importer{w, conf, bp.Dir},
		Error: func(err error) {
			if typesVerbose {
				log.Println(err)
			}
		},
	}
	if pkg == nil {
		pkg, err = typesConf.Check(checkName, w.FileSet, files, conf.Info)
		conf.Pkg = pkg
	}
	w.importingName[checkName] = false
	w.imported[name] = pkg

	if len(xfiles) > 0 {
		xpkg, _ := typesConf.Check(checkName+"_test", w.FileSet, xfiles, conf.XInfo)
		w.imported[checkName+"_test"] = xpkg
		conf.XPkg = xpkg
	}
	return
}

type Importer struct {
	w    *PkgWalker
	conf *PkgConfig
	dir  string
}

func (im *Importer) Import(name string) (pkg *types.Package, err error) {
	if im.conf.AllowBinary && im.w.isBinaryPkg(name) {
		pkg, err = im.w.gcimported.Import(name)
		if pkg != nil && pkg.Complete() {
			return
		}
		//		pkg = im.w.gcimporter[name]
		//		if pkg != nil && pkg.Complete() {
		//			return
		//		}
		//		pkg, err = importer.Default().Import(name)
		//		if pkg != nil && pkg.Complete() {
		//			im.w.gcimporter[name] = pkg
		//			return
		//		}
	}
	return im.w.Import(im.dir, name, &PkgConfig{IgnoreFuncBodies: true, AllowBinary: true, WithTestFiles: false})
}

func (w *PkgWalker) parseFile(dir, file string, src interface{}) (*ast.File, error) {
	filename := filepath.Join(dir, file)
	f, _ := w.parsedFileCache[filename]
	if f != nil {
		return f, nil
	}

	var err error

	// generate missing context-dependent files.
	if w.Context != nil && file == fmt.Sprintf("zgoos_%s.go", w.Context.GOOS) {
		src := fmt.Sprintf("package runtime; const theGoos = `%s`", w.Context.GOOS)
		f, err = parser.ParseFile(w.FileSet, filename, src, 0)
		if err != nil {
			log.Fatalf("incorrect generated file: %s", err)
		}
	}

	if w.Context != nil && file == fmt.Sprintf("zgoarch_%s.go", w.Context.GOARCH) {
		src := fmt.Sprintf("package runtime; const theGoarch = `%s`", w.Context.GOARCH)
		f, err = parser.ParseFile(w.FileSet, filename, src, 0)
		if err != nil {
			log.Fatalf("incorrect generated file: %s", err)
		}
	}

	if f == nil {
		flag := parser.AllErrors
		if typeParserComment {
			flag |= parser.ParseComments
		}
		f, err = parser.ParseFile(w.FileSet, filename, src, flag)
		if err != nil {
			return f, err
		}
	}

	w.parsedFileCache[filename] = f
	return f, nil
}

func ParserObjectKind(ident *ast.Ident, obj types.Object, conf *PkgConfig) (ObjKind, error) {
	if obj != nil {
		return parserObjKind(obj)
	}
	for _, f := range conf.Files {
		if f.Name == ident {
			return ObjPackage, nil
		}
	}
	for _, f := range conf.XTestFiles {
		if f.Name == ident {
			return ObjPackageX, nil
		}
	}
	return ObjUnknown, fmt.Errorf("unknown ident %v", ident)
}

type ObjList struct {
	Kind        ObjKind
	Ident       *ast.Ident
	Obj         types.Object
	List        []int
	IsTestFunc  bool
	IsAnonymous bool
}

func NewObjList(kind ObjKind, ident *ast.Ident, obj types.Object) *ObjList {
	v := &ObjList{Kind: kind, Ident: ident, Obj: obj}
	v.List = append(v.List, int(ident.Pos()))
	return v
}

func (o *ObjList) String() string {
	return fmt.Sprintf("%v %v %v", o.Kind, o.Ident.Name, o.List)
}

type Rename struct {
	Old string
	New string
}

func (w *PkgWalker) LookupMethodRecvObject(fn *types.Func) types.Object {
	sig, ok := fn.Type().(*types.Signature)
	if !ok {
		return nil
	}
	typ := sig.Recv().Type()
	if pt, ok := typ.(*types.Pointer); ok {
		typ = pt.Elem()
	}
	if named, ok := typ.(*types.Named); ok {
		return named.Obj()
	}
	return nil
}

func (w *PkgWalker) LookupIdentComment(value token.Pos) *ast.CommentGroup {
	pos := w.FileSet.Position(value)
	file := w.parsedFileCache[pos.Filename]
	if file == nil {
		return nil
	}
	line := pos.Line
	var group *ast.CommentGroup
	for _, v := range file.Comments {
		lastLine := w.FileSet.Position(v.End()).Line
		if lastLine == line || lastLine == line-1 {
			group = v
		} else if lastLine > line {
			break
		}
	}
	return group
}

func (w *PkgWalker) LookupCursor(pkg *types.Package, conf *PkgConfig, cursor *FileCursor) {
	is := w.CheckIsImport(cursor)
	if is != nil {
		if cursor.xtest {
			w.LookupImport(conf.XPkg, conf.XInfo, cursor, is)
		} else {
			w.LookupImport(conf.Pkg, conf.Info, cursor, is)
		}
	} else {
		w.LookupObjects(conf, cursor)
	}
}

func (w *PkgWalker) LookupImport(pkg *types.Package, pkgInfo *types.Info, cursor *FileCursor, is *ast.ImportSpec) {
	fpath, err := strconv.Unquote(is.Path.Value)
	if err != nil {
		return
	}

	if typesFindDef {
		fmt.Println(w.FileSet.Position(is.Pos()))
	}

	fbase := fpath
	pos := strings.LastIndexAny(fpath, "./-\\")
	if pos != -1 {
		fbase = fpath[pos+1:]
	}

	var fname string
	if is.Name != nil {
		fname = is.Name.Name
	} else {
		fname = fbase
	}

	if typesFindInfo {
		if fname == fpath {
			fmt.Printf("package %s\n", fname)
		} else {
			fmt.Printf("package %s (\"%s\")\n", fname, fpath)
		}
	}

	if !typesFindUse {
		return
	}

	fid := pkg.Path() + "." + fname

	var usages []int
	for id, obj := range pkgInfo.Uses {
		if obj != nil && obj.Id() == fid { //!= nil && cursorObj.Pos() == obj.Pos() {
			if _, ok := obj.(*types.PkgName); ok {
				usages = append(usages, int(id.Pos()))
			}
		}
	}
	(sort.IntSlice(usages)).Sort()
	for _, pos := range usages {
		fmt.Println(w.FileSet.Position(token.Pos(pos)))
	}
}

func testObjKind(obj types.Object, kind ObjKind) bool {
	k, err := parserObjKind(obj)
	if err != nil {
		return false
	}
	return k == kind
}

func parserObjKind(obj types.Object) (ObjKind, error) {
	var kind ObjKind
	switch t := obj.(type) {
	case *types.PkgName:
		kind = ObjPkgName
	case *types.Const:
		kind = ObjConst
	case *types.TypeName:
		kind = ObjTypeName
		switch t.Type().Underlying().(type) {
		case *types.Interface:
			kind = ObjInterface
		case *types.Struct:
			kind = ObjStruct
		}
	case *types.Var:
		kind = ObjVar
		if t.IsField() {
			kind = ObjField
		}
	case *types.Func:
		kind = ObjFunc
		if sig, ok := t.Type().(*types.Signature); ok {
			if sig.Recv() != nil {
				kind = ObjMethod
			}
		}
	case *types.Label:
		kind = ObjLabel
	case *types.Builtin:
		kind = ObjBuiltin
	case *types.Nil:
		kind = ObjNil
	default:
		return ObjNone, fmt.Errorf("unknown obj type %T", obj)
	}
	return kind, nil
}

func (w *PkgWalker) LookupStructFromField(info *types.Info, cursorPkg *types.Package, cursorObj types.Object, cursorPos token.Pos) types.Object {
	if info == nil {
		conf := &PkgConfig{
			IgnoreFuncBodies: true,
			AllowBinary:      true,
			WithTestFiles:    true,
			Info: &types.Info{
				Defs: make(map[*ast.Ident]types.Object),
			},
		}
		w.imported[cursorPkg.Path()] = nil
		pkg, _ := w.Import("", cursorPkg.Path(), conf)
		if pkg != nil {
			info = conf.Info
		}
	}
	for _, obj := range info.Defs {
		if obj == nil {
			continue
		}
		if _, ok := obj.(*types.TypeName); ok {
			if t, ok := obj.Type().Underlying().(*types.Struct); ok {
				for i := 0; i < t.NumFields(); i++ {
					if t.Field(i).Pos() == cursorPos {
						return obj
					}
				}
			}
		}
	}
	return nil
}

func (w *PkgWalker) lookupNamedField(named *types.Named, name string) *types.Named {
	if istruct, ok := named.Underlying().(*types.Struct); ok {
		for i := 0; i < istruct.NumFields(); i++ {
			field := istruct.Field(i)
			if field.Anonymous() {
				fieldType := orgType(field.Type())
				if typ, ok := fieldType.(*types.Named); ok {
					if na := w.lookupNamedField(typ, name); na != nil {
						return na
					}
				}
			} else {
				if field.Name() == name {
					return named
				}
			}
		}
	}
	return nil
}

func (w *PkgWalker) lookupNamedMethod(named *types.Named, name string) (types.Object, *types.Named) {
	if iface, ok := named.Underlying().(*types.Interface); ok {
		for i := 0; i < iface.NumMethods(); i++ {
			fn := iface.Method(i)
			if fn.Name() == name {
				return fn, named
			}
		}
		for i := 0; i < iface.NumEmbeddeds(); i++ {
			if obj, na := w.lookupNamedMethod(iface.Embedded(i), name); obj != nil {
				return obj, na
			}
		}
		return nil, nil
	}
	if istruct, ok := named.Underlying().(*types.Struct); ok {
		for i := 0; i < named.NumMethods(); i++ {
			fn := named.Method(i)
			if fn.Name() == name {
				return fn, named
			}
		}
		for i := 0; i < istruct.NumFields(); i++ {
			field := istruct.Field(i)
			if !field.Anonymous() {
				continue
			}
			if typ, ok := field.Type().(*types.Named); ok {
				if obj, na := w.lookupNamedMethod(typ, name); obj != nil {
					return obj, na
				}
			}
		}
	}
	return nil, nil
}

func IsSamePkg(a, b *types.Package) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Path() == b.Path()
}

func IsSameObject(a, b types.Object) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	var apath string
	var bpath string
	if a.Pkg() != nil {
		apath = a.Pkg().Path()
	}
	if b.Pkg() != nil {
		bpath = b.Pkg().Path()
	}
	if apath != bpath {
		return false
	}
	if a.Id() != b.Id() {
		return false
	}
	return a.String() == b.String()
}

func orgType(typ types.Type) types.Type {
	if pt, ok := typ.(*types.Pointer); ok {
		return pt.Elem()
	}
	return typ
}

func (w *PkgWalker) LookupObjects(conf *PkgConfig, cursor *FileCursor) {
	var cursorObj types.Object
	var cursorSelection *types.Selection
	var cursorObjIsDef bool
	//lookup defs

	var pkg *types.Package
	var pkgInfo *types.Info
	if cursor.xtest {
		pkgInfo = conf.XInfo
		pkg = conf.XPkg
	} else {
		pkgInfo = conf.Info
		pkg = conf.Pkg
	}

	_ = cursorObjIsDef
	if cursorObj == nil {
		for sel, obj := range pkgInfo.Selections {
			if cursor.pos >= sel.Sel.Pos() && cursor.pos <= sel.Sel.End() {
				cursorObj = obj.Obj()
				cursorSelection = obj
				break
			}
		}
	}
	if cursorObj == nil {
		for id, obj := range pkgInfo.Defs {
			if cursor.pos >= id.Pos() && cursor.pos <= id.End() {
				cursorObj = obj
				cursorObjIsDef = true
				break
			}
		}
	}
	_ = cursorSelection
	if cursorObj == nil {
		for id, obj := range pkgInfo.Uses {
			if cursor.pos >= id.Pos() && cursor.pos <= id.End() {
				cursorObj = obj
				break
			}
		}
	}
	if cursorObj == nil {
		return
	}
	kind, err := parserObjKind(cursorObj)
	if err != nil {
		log.Fatalln(err)
	}

	if kind == ObjField {
		if cursorObj.(*types.Var).Anonymous() {
			typ := orgType(cursorObj.Type())
			if named, ok := typ.(*types.Named); ok {
				cursorObj = named.Obj()
			}
		}
	}
	cursorPkg := cursorObj.Pkg()
	cursorPos := cursorObj.Pos()
	//var fieldTypeInfo *types.Info
	var fieldTypeObj types.Object
	//	if cursorPkg == pkg {
	//		fieldTypeInfo = pkgInfo
	//	}
	cursorIsInterfaceMethod := false
	var cursorInterfaceTypeName string

	if kind == ObjMethod && cursorSelection != nil && cursorSelection.Recv() != nil {
		sig := cursorObj.(*types.Func).Type().Underlying().(*types.Signature)
		if _, ok := sig.Recv().Type().Underlying().(*types.Interface); ok {
			if named, ok := cursorSelection.Recv().(*types.Named); ok {
				obj, typ := w.lookupNamedMethod(named, cursorObj.Name())
				if obj != nil {
					cursorObj = obj
				}
				if typ != nil {
					cursorPkg = typ.Obj().Pkg()
					cursorInterfaceTypeName = typ.Obj().Name()
				}
				cursorIsInterfaceMethod = true
			}
		}
	} else if kind == ObjField && cursorSelection != nil {
		if recv := cursorSelection.Recv(); recv != nil {
			typ := orgType(recv)
			if typ != nil {
				if name, ok := typ.(*types.Named); ok {
					fieldTypeObj = name.Obj()
					na := w.lookupNamedField(name, cursorObj.Name())
					if na != nil {
						fieldTypeObj = na.Obj()
					}
				}
			}
		}
	}
	if cursorPkg != nil && cursorPkg != pkg &&
		kind != ObjPkgName && w.isBinaryPkg(cursorPkg.Path()) {
		conf := &PkgConfig{
			IgnoreFuncBodies: true,
			AllowBinary:      true,
			WithTestFiles:    true,
			Info: &types.Info{
				Defs: make(map[*ast.Ident]types.Object),
			},
		}
		pkg, _ := w.Import("", cursorPkg.Path(), conf)
		if pkg != nil {
			if cursorIsInterfaceMethod {
				for _, obj := range conf.Info.Defs {
					if obj == nil {
						continue
					}
					if fn, ok := obj.(*types.Func); ok {
						if fn.Name() == cursorObj.Name() {
							if sig, ok := fn.Type().Underlying().(*types.Signature); ok {
								if named, ok := sig.Recv().Type().(*types.Named); ok {
									if named.Obj() != nil && named.Obj().Name() == cursorInterfaceTypeName {
										cursorPos = obj.Pos()
										break
									}
								}
							}
						}
					}
				}
			} else if kind == ObjField && fieldTypeObj != nil {
				for _, obj := range conf.Info.Defs {
					if obj == nil {
						continue
					}
					if _, ok := obj.(*types.TypeName); ok {
						if IsSameObject(fieldTypeObj, obj) {
							if t, ok := obj.Type().Underlying().(*types.Struct); ok {
								for i := 0; i < t.NumFields(); i++ {
									if t.Field(i).Id() == cursorObj.Id() {
										cursorPos = t.Field(i).Pos()
										break
									}
								}
							}
							break
						}
					}
				}
			} else {
				for k, v := range conf.Info.Defs {
					if k != nil && v != nil && IsSameObject(v, cursorObj) {
						cursorPos = k.Pos()
						break
					}
				}
			}
		}
		//		if kind == ObjField || cursorIsInterfaceMethod {
		//			fieldTypeInfo = conf.Info
		//		}
	}
	//	if kind == ObjField {
	//		fieldTypeObj = w.LookupStructFromField(fieldTypeInfo, cursorPkg, cursorObj, cursorPos)
	//	}
	if typesFindDef {
		fmt.Println(w.FileSet.Position(cursorPos))
	}
	if typesFindInfo {
		if kind == ObjField && fieldTypeObj != nil {
			typeName := fieldTypeObj.Name()
			if fieldTypeObj.Pkg() != nil && fieldTypeObj.Pkg() != pkg {
				typeName = fieldTypeObj.Pkg().Name() + "." + fieldTypeObj.Name()
			}
			fmt.Println(typeName, simpleObjInfo(cursorObj))
		} else if kind == ObjBuiltin {
			fmt.Println(builtinInfo(cursorObj.Name()))
		} else if kind == ObjPkgName {
			fmt.Println(cursorObj.String())
		} else if cursorIsInterfaceMethod {
			fmt.Println(strings.Replace(simpleObjInfo(cursorObj), "(interface)", cursorPkg.Name()+"."+cursorInterfaceTypeName, 1))
		} else {
			fmt.Println(simpleObjInfo(cursorObj))
		}
	}

	if typesFindDoc && typesFindDef {
		pos := w.FileSet.Position(cursorPos)
		file := w.parsedFileCache[pos.Filename]
		if file != nil {
			line := pos.Line
			var group *ast.CommentGroup
			for _, v := range file.Comments {
				lastLine := w.FileSet.Position(v.End()).Line
				if lastLine == line || lastLine == line-1 {
					group = v
				} else if lastLine > line {
					break
				}
			}
			if group != nil {
				fmt.Println(group.Text())
			}
		}
	}
	if !typesFindUse {
		return
	}

	var usages []int
	if kind == ObjPkgName {
		for id, obj := range pkgInfo.Uses {
			if obj != nil && obj.Id() == cursorObj.Id() { //!= nil && cursorObj.Pos() == obj.Pos() {
				if _, ok := obj.(*types.PkgName); ok {
					usages = append(usages, int(id.Pos()))
				}
			}
		}
	} else {
		//		for id, obj := range pkgInfo.Defs {
		//			if obj == cursorObj { //!= nil && cursorObj.Pos() == obj.Pos() {
		//				usages = append(usages, int(id.Pos()))
		//			}
		//		}
		for id, obj := range pkgInfo.Uses {
			if obj == cursorObj { //!= nil && cursorObj.Pos() == obj.Pos() {
				usages = append(usages, int(id.Pos()))
			}
		}
	}
	var pkg_path string
	var xpkg_path string
	if conf.Pkg != nil {
		pkg_path = conf.Pkg.Path()
	}
	if conf.XPkg != nil {
		xpkg_path = conf.XPkg.Path()
	}

	if cursorPkg != nil &&
		(cursorPkg.Path() == pkg_path || cursorPkg.Path() == xpkg_path) &&
		kind != ObjPkgName {
		usages = append(usages, int(cursorPos))
	}

	(sort.IntSlice(usages)).Sort()
	for _, pos := range usages {
		fmt.Println(w.FileSet.Position(token.Pos(pos)))
	}
	//check look for current pkg.object on pkg_test
	if typesFindUseAll || IsSamePkg(cursorPkg, conf.Pkg) {
		var addInfo *types.Info
		if conf.Cursor.xtest {
			addInfo = conf.Info
		} else {
			addInfo = conf.XInfo
		}
		if addInfo != nil && cursorPkg != nil {
			var usages []int
			//		for id, obj := range addInfo.Defs {
			//			if id != nil && obj != nil && obj.Id() == cursorObj.Id() {
			//				usages = append(usages, int(id.Pos()))
			//			}
			//		}
			for k, v := range addInfo.Uses {
				if k != nil && v != nil && IsSameObject(v, cursorObj) {
					usages = append(usages, int(k.Pos()))
				}
			}
			(sort.IntSlice(usages)).Sort()
			for _, pos := range usages {
				fmt.Println(w.FileSet.Position(token.Pos(pos)))
			}
		}
	}
	if !typesFindUseAll {
		return
	}

	if cursorPkg == nil {
		return
	}

	var find_def_pkg string
	var uses_paths []string
	if cursorPkg.Path() != pkg_path && cursorPkg.Path() != xpkg_path {
		find_def_pkg = cursorPkg.Path()
		uses_paths = append(uses_paths, cursorPkg.Path())
	}

	cursorPkgPath := cursorObj.Pkg().Path()
	if pkgutil.IsVendorExperiment() {
		cursorPkgPath = pkgutil.VendorPathToImportPath(cursorPkgPath)
	}

	ForEachPackage(&build.Default, func(importPath string, err error) {
		if err != nil {
			return
		}
		if importPath == conf.Pkg.Path() {
			return
		}
		bp, err := w.importPath(importPath, 0)
		if err != nil {
			return
		}
		find := false
		if bp.ImportPath == cursorPkg.Path() {
			find = true
		} else {
			for _, v := range bp.Imports {
				if v == cursorPkgPath {
					find = true
					break
				}
			}
		}
		if find {
			for _, v := range uses_paths {
				if v == bp.ImportPath {
					return
				}
			}
			uses_paths = append(uses_paths, bp.ImportPath)
		}
	})

	w.imported = make(map[string]*types.Package)
	for _, v := range uses_paths {
		conf := &PkgConfig{
			IgnoreFuncBodies: false,
			AllowBinary:      true,
			WithTestFiles:    true,
			Info: &types.Info{
				Uses: make(map[*ast.Ident]types.Object),
			},
			XInfo: &types.Info{
				Uses: make(map[*ast.Ident]types.Object),
			},
		}
		w.imported[v] = nil
		var usages []int
		vpkg, _ := w.Import("", v, conf)
		if vpkg != nil && vpkg != pkg {
			if conf.Info != nil {
				for k, v := range conf.Info.Uses {
					if k != nil && v != nil && IsSameObject(v, cursorObj) {
						usages = append(usages, int(k.Pos()))
					}
				}
			}
			if conf.XInfo != nil {
				for k, v := range conf.XInfo.Uses {
					if k != nil && v != nil && IsSameObject(v, cursorObj) {
						usages = append(usages, int(k.Pos()))
					}
				}
			}
		}
		if v == find_def_pkg {
			usages = append(usages, int(cursorPos))
		}
		(sort.IntSlice(usages)).Sort()
		for _, pos := range usages {
			fmt.Println(w.FileSet.Position(token.Pos(pos)))
		}
	}
}

func (w *PkgWalker) CheckIsImport(cursor *FileCursor) *ast.ImportSpec {
	if cursor.fileDir == "" {
		return nil
	}
	file, _ := w.parseFile(cursor.fileDir, cursor.fileName, cursor.src)
	if file == nil {
		return nil
	}
	for _, is := range file.Imports {
		if inRange(is, cursor.pos) {
			return is
		}
	}
	return nil
}

func inRange(node ast.Node, p token.Pos) bool {
	if node == nil {
		return false
	}
	return p >= node.Pos() && p <= node.End()
}

func (w *PkgWalker) nodeString(node interface{}) string {
	if node == nil {
		return ""
	}
	var b bytes.Buffer
	printer.Fprint(&b, w.FileSet, node)
	return b.String()
}
