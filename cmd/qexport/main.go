package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/doc"
	"go/format"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

var (
	flagDefaultContext           bool
	flagRenameNewTypeFunc        bool
	flagSkipErrorImplementStruct bool
	flagQlangLowerCaseStyle      bool
	flagCustomContext            string
	flagExportPath               string
	flagUpdatePath               string
)

const help = `Export go packages to qlang modules.

Usage:
  qexport [option] packages

The packages for go package list or std for golang all standard packages.
`

func usage() {
	fmt.Fprintln(os.Stderr, help)
	flag.PrintDefaults()
}

func init() {
	flag.StringVar(&flagCustomContext, "contexts", "", "optional comma-separated list of <goos>-<goarch>[-cgo] to override default contexts.")
	flag.BoolVar(&flagDefaultContext, "defctx", false, "optional use default context for build, default use all contexts.")
	flag.BoolVar(&flagRenameNewTypeFunc, "convnew", true, "optional convert NewType func to type func")
	flag.BoolVar(&flagSkipErrorImplementStruct, "skiperrimpl", true, "optional skip error interface implement struct.")
	flag.BoolVar(&flagQlangLowerCaseStyle, "lowercase", true, "optional use qlang lower case style.")
	flag.StringVar(&flagExportPath, "outpath", "./qlang", "optional set export root path")
	flag.StringVar(&flagUpdatePath, "updatepath", "", "option set update qlang package root")
}

var (
	ac *ApiCheck
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		usage()
		return
	}

	if flagCustomContext != "" {
		flagDefaultContext = false
		setCustomContexts(flagCustomContext)
	}

	//load ApiCheck
	ac = NewApiCheck()
	err := ac.LoadBase("go1", "go1.1", "go1.2", "go1.3", "go1.4")
	if err != nil {
		log.Println(err)
	}
	err = ac.LoadApi("go1.5", "go1.6", "go1.7")
	if err != nil {
		log.Println(err)
	}

	var outpath string
	if filepath.IsAbs(flagExportPath) {
		outpath = flagExportPath
	} else {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
		outpath = filepath.Join(dir, flagExportPath)
	}

	var pkgs []string
	if args[0] == "std" {
		out, err := exec.Command("go", "list", "-e", args[0]).Output()
		if err != nil {
			log.Fatal(err)
		}
		pkgs = strings.Fields(string(out))
	} else {
		pkgs = args
	}
	var exportd []string
	for _, pkg := range pkgs {
		err := export(pkg, outpath, true)
		if err != nil {
			log.Printf("warning skip pkg %q, error %v.\n", pkg, err)
		} else {
			exportd = append(exportd, pkg)
		}
	}
	for _, pkg := range exportd {
		log.Printf("export pkg %q success.\n", pkg)
	}
}

func checkStructHasUnexportField(decl *ast.GenDecl) bool {
	if len(decl.Specs) > 0 {
		if ts, ok := decl.Specs[0].(*ast.TypeSpec); ok {
			if st, ok := ts.Type.(*ast.StructType); ok {
				if st.Fields != nil {
					for _, f := range st.Fields.List {
						for _, n := range f.Names {
							if !ast.IsExported(n.Name) {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func export(pkg string, outpath string, skipOSArch bool) error {
	p, err := NewPackage(pkg, flagDefaultContext)
	if err != nil {
		return err
	}

	p.Parser()

	bp := p.BuildPackage()
	if bp == nil {
		return errors.New("not find build")
	}

	pkgName := bp.Name

	if bp.Name == "main" {
		return errors.New("skip main pkg")
	}

	if pkg == "unsafe" {
		return errors.New("skip unsafe pkg")
	}

	if p.CommonCount() == 0 {
		return errors.New("empty common exports")
	}

	//skip internal
	for _, path := range strings.Split(bp.ImportPath, "/") {
		if path == "internal" {
			return errors.New("skip internal pkg")
		} else if path == "vendor" {
			return errors.New("skip vendor pkg")
		}
	}

	var upinfo *UpdateInfo
	if flagUpdatePath != "" {
		upinfo, _ = CheckUpdateInfo(bp.Name, pkg)
	}

	if upinfo != nil {
		log.Printf("update pkg %q from %q\n", pkg, upinfo.updataPkg)
	}

	fnskip := func(key string) bool {
		return false
	}

	//check fnskip
	if upinfo != nil {
		fnskip = func(key string) bool {
			if _, ok := upinfo.usesMap[key]; ok {
				return true
			}
			return false
		}
		p.fnskip = fnskip
	}

	var buf bytes.Buffer
	outf := func(format string, a ...interface{}) (err error) {
		_, err = buf.WriteString(fmt.Sprintf(format, a...))
		return
	}

	checkVer := func(key string) (string, bool) {
		vers := ac.FincApis(bp.ImportPath + "." + key)
		if len(vers) > 0 {
			return vers[0], true
		}
		return "", false
	}
	checkConst := func(key string) KeyType {
		return ac.CheckConstType(bp.ImportPath + "." + key)
	}

	// go ver map
	verMap := make(map[string][]string)
	outfv := func(ver string, k, v string) {
		verMap[ver] = append(verMap[ver], fmt.Sprintf("Exports[%q]=%s", k, v))
	}

	var hasTypeExport bool
	verHasTypeExport := make(map[string]bool)

	//const
	if keys, _ := p.FilterCommon(Const); len(keys) > 0 {
		outf("\n")
		for _, v := range keys {
			name := v
			fn := pkgName + "." + v
			typ := checkConst(v)
			if typ == ConstInt64 {
				fn = "int64(" + fn + ")"
			} else if typ == ConstUnit64 {
				fn = "uint64(" + fn + ")"
			}
			if vers, ok := checkVer(v); ok {
				outfv(vers, name, fn)
			} else {
				outf("\t%q:\t%s,\n", name, fn)
			}
		}
	}

	//vars
	if keys, m := p.FilterCommon(Var); len(keys) > 0 {
		outf("\n")
		skeys, _ := p.FilterCommon(Struct)
		for _, v := range keys {
			mv := m[v].(*doc.Value)
			var isStructVar bool
			if typ := p.simpleValueDeclType(mv.Decl); typ != "" {
				for _, k := range skeys {
					if typ == k {
						isStructVar = true
						log.Printf("warning convert struct var to ref %s %s (%s)\n", bp.ImportPath, v, typ)
					}
				}
			}
			name := v
			fn := pkgName + "." + v
			if isStructVar {
				fn = "&" + fn
			}
			if vers, ok := checkVer(v); ok {
				outfv(vers, name, fn)
			} else {
				outf("\t%q:\t%s,\n", name, fn)
			}
		}
	}

	//funcs
	if keys, _ := p.FilterCommon(Func); len(keys) > 0 {
		outf("\n")
		for _, v := range keys {
			name := toQlangName(v)
			fn := pkgName + "." + v
			if vers, ok := checkVer(v); ok {
				outfv(vers, name, fn)
			} else {
				outf("\t%q:\t%s,\n", name, fn)
			}
		}
	}

	//interface
	if keys, m := p.FilterCommon(Interface); len(keys) > 0 {
		outf("\n")
		for _, v := range keys {
			t, ok := m[v]
			if !ok {
				continue
			}
			dt, ok := t.(*doc.Type)
			if !ok {
				continue
			}

			// exported funcs
			var funcsNew []string
			var funcsOther []string
			for _, f := range dt.Funcs {
				if fnskip(f.Name) {
					continue
				}
				if ast.IsExported(f.Name) {
					if strings.HasPrefix(f.Name, "New"+v) {
						funcsNew = append(funcsNew, f.Name)
					} else {
						funcsOther = append(funcsOther, f.Name)
					}
				}
			}

			for _, f := range funcsNew {
				name := toQlangName(f)
				if flagRenameNewTypeFunc && len(funcsNew) == 1 {
					name = toQlangName(v)
					if ast.IsExported(name) {
						name = strings.ToLower(name)
						log.Printf("waring convert %s to %s", bp.ImportPath+"."+f, name)
					}
				}
				fn := pkgName + "." + f
				if vers, ok := checkVer(f); ok {
					outfv(vers, name, fn)
				} else {
					outf("\t%q:\t%s,\n", name, fn)
				}
			}

			for _, f := range funcsOther {
				name := toQlangName(f)
				fn := pkgName + "." + f
				if vers, ok := checkVer(f); ok {
					outfv(vers, name, fn)
				} else {
					outf("\t%q:\t%s,\n", name, fn)
				}
			}
		}
	}

	//structs
	if keys, m := p.FilterCommon(Struct); len(keys) > 0 {
		outf("\n")
		for _, v := range keys {
			t, ok := m[v]
			if !ok {
				continue
			}
			dt, ok := t.(*doc.Type)
			if !ok {
				continue
			}
			//fmt.Printf("%v\n", dt.Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List)
			// exported funcs
			var funcsNew []string
			var funcsOther []string
			for _, f := range dt.Funcs {
				if fnskip(f.Name) {
					continue
				}
				if ast.IsExported(f.Name) {
					if strings.HasPrefix(f.Name, "New"+v) {
						funcsNew = append(funcsNew, f.Name)
					} else {
						funcsOther = append(funcsOther, f.Name)
					}
				}
			}
			// check error interface implement struct
			if flagSkipErrorImplementStruct && strings.HasSuffix(v, "Error") {
				check := func(name string) bool {
					for _, f := range dt.Methods {
						if f.Name == name {
							return true
						}
					}
					return false
				}
				if check("Error") {
					log.Printf("warning skip struct %s, is error interface{} implement\n", bp.ImportPath+"."+v)
					continue
				}
			}

			//export type, qlang.StructOf((*strings.Reader)(nil))
			if ast.IsExported(v) && !checkStructHasUnexportField(dt.Decl) {
				name := v
				fn := fmt.Sprintf("qlang.StructOf((*%s.%s)(nil))", pkgName, v)
				if vers, ok := checkVer(v); ok {
					verHasTypeExport[vers] = true
					outfv(vers, name, fn)
				} else {
					hasTypeExport = true
					outf("\t%q:\t%s,\n", name, fn)
				}
			}

			for _, f := range funcsNew {
				name := toQlangName(f)
				if flagRenameNewTypeFunc && len(funcsNew) == 1 {
					name = toQlangName(v)
					//NewRGBA => rgba
					if ast.IsExported(name) {
						name = strings.ToLower(name)
						log.Printf("waring convert %s to %s", bp.ImportPath+"."+f, name)
					}
				}
				fn := pkgName + "." + f
				if vers, ok := checkVer(f); ok {
					outfv(vers, name, fn)
				} else {
					outf("\t%q:\t%s,\n", name, fn)
				}
			}

			for _, f := range funcsOther {
				name := toQlangName(f)
				fn := pkgName + "." + f
				if vers, ok := checkVer(f); ok {
					outfv(vers, name, fn)
				} else {
					outf("\t%q:\t%s,\n", name, fn)
				}
			}
		}
	}
	// write root dir
	root := filepath.Join(outpath, pkg)
	os.MkdirAll(root, 0777)

	verUpdateMap := make(map[string]*InsertInfo)

	//check update ver map
	if upinfo != nil {
		for _, ii := range upinfo.initsInfo {
			for _, group := range ii.file.Comments {
				for _, comment := range group.List {
					text := comment.Text
					if strings.HasPrefix(text, "//") {
						text = strings.TrimSpace(text[2:])
						if strings.HasPrefix(text, "+build ") {
							tag := strings.TrimSpace(text[6:])
							verUpdateMap[tag] = ii
						}
					}
				}
			}
		}
	}

	//write exports
	if upinfo != nil && upinfo.exportsInfo != nil {
		upinfo.exportsInfo.UpdateFile(root, buf.Bytes(), hasTypeExport)
	} else {
		var head bytes.Buffer
		outHeadf := func(format string, a ...interface{}) (err error) {
			_, err = head.WriteString(fmt.Sprintf(format, a...))
			return
		}

		//write package
		outHeadf("package %s\n", pkgName)

		//check pkgName used
		if strings.Index(buf.String(), pkgName+".") > 0 {
			//write imports
			outHeadf("import (\n")
			outHeadf("\t%q\n", pkg)
			//check qlang used
			if hasTypeExport {
				outHeadf("\n\tqlang \"qlang.io/spec\"\n")
			}
			outHeadf(")\n\n")
		}

		var source []byte
		source = append(source, head.Bytes()...)
		exportsDecl := fmt.Sprintf(`// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name": "%s",	
`, pkg)
		source = append(source, []byte(exportsDecl)...)
		source = append(source, buf.Bytes()...)
		source = append(source, '}')

		// format
		data, err := format.Source(source)
		if err != nil {
			return err
		}

		file, err := os.Create(filepath.Join(root, pkgName+".go"))
		if err != nil {
			return err
		}
		file.Write(data)
		file.Close()
	}
	// write version data
	for ver, lines := range verMap {
		if ii, ok := verUpdateMap[ver]; ok {
			data := strings.Join(lines, "\n\t")
			ii.UpdateFile(root, []byte(data), verHasTypeExport[ver])
			continue
		}
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("// +build %s\n\n", ver))
		buf.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
		if verHasTypeExport[ver] {
			buf.WriteString("import (\n")
			buf.WriteString(fmt.Sprintf("\t%q\n\n", bp.ImportPath))
			buf.WriteString(fmt.Sprintf("\tqlang %q\n", "qlang.io/spec"))
			buf.WriteString(")\n")
		} else {
			buf.WriteString(fmt.Sprintf("import %q\n", bp.ImportPath))
		}
		buf.WriteString("func init() {\n\t")
		buf.WriteString(strings.Join(lines, "\n\t"))
		buf.WriteString("\n}")
		data, err := format.Source(buf.Bytes())
		if err != nil {
			return err
		}
		file, err := os.Create(filepath.Join(root, pkgName+"-"+strings.Replace(ver, ".", "", 1)+".go"))
		if err != nil {
			return err
		}
		file.Write(data)
		file.Close()
	}
	return nil
}

// convert to qlang name, default use lower case style, Name => name, NAME => NAME
func toQlangName(s string) string {
	if !flagQlangLowerCaseStyle {
		return s
	}
	if len(s) <= 1 {
		return s
	}
	if unicode.IsLower(rune(s[1])) {
		return strings.ToLower(s[0:1]) + s[1:]
	}
	return s
}
