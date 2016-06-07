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
	flagCustomContext            string
	flagExportPath               string
)

const help = `Export go packages to qlang modules.

Usage:
  qexport [-contexts=""] [-defctx=false] [-convnew=true] [-skiperrimpl=true] [-outpath="./qlang"] packages

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
	flag.StringVar(&flagExportPath, "outpath", "./qlang", "optional set export root path")
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

var (
	uint64_const_keys = []string{
		"crc64.ECMA",
		"crc64.ISO",
		"math.MaxUint64",
	}
)

func isUint64Const(key string) bool {
	for _, k := range uint64_const_keys {
		if key == k {
			return true
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
		}
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

	// go ver map
	verMap := make(map[string][]string)
	outfv := func(ver string, k, v string) {
		verMap[ver] = append(verMap[ver], fmt.Sprintf("Exports[%q]=%s", k, v))
	}

	var hasTypeExport bool
	verHasTypeExport := make(map[string]bool)

	//write exports
	outf(`// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name": "%s",	
`, pkg)

	//const
	if keys, _ := p.FilterCommon(Const); len(keys) > 0 {
		outf("\n")
		for _, v := range keys {
			name := v
			fn := pkgName + "." + v
			if isUint64Const(fn) {
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
			name := toLowerCaseStyle(v)
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
				if ast.IsExported(f.Name) {
					if strings.HasPrefix(f.Name, "New"+v) {
						funcsNew = append(funcsNew, f.Name)
					} else {
						funcsOther = append(funcsOther, f.Name)
					}
				}
			}

			for _, f := range funcsNew {
				name := toLowerCaseStyle(f)
				if flagRenameNewTypeFunc && len(funcsNew) == 1 {
					name = toLowerCaseStyle(v)
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
				name := toLowerCaseStyle(f)
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
			// exported funcs
			var funcsNew []string
			var funcsOther []string
			for _, f := range dt.Funcs {
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

			//export type, qlang.NewType(reflect.TypeOf((*http.Client)(nil)).Elem())
			//export type, qlang.StructOf((*strings.Reader)(nil))
			if ast.IsExported(v) {
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
				name := toLowerCaseStyle(f)
				if flagRenameNewTypeFunc && len(funcsNew) == 1 {
					name = toLowerCaseStyle(v)
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
				name := toLowerCaseStyle(f)
				fn := pkgName + "." + f
				if vers, ok := checkVer(f); ok {
					outfv(vers, name, fn)
				} else {
					outf("\t%q:\t%s,\n", name, fn)
				}
			}
		}
	}

	// end exports
	outf("}")

	var head bytes.Buffer
	outHeadf := func(format string, a ...interface{}) (err error) {
		_, err = head.WriteString(fmt.Sprintf(format, a...))
		return
	}

	//write package
	outHeadf("package %s\n", pkgName)

	if strings.Count(buf.String(), ",") > 1 {
		//write imports
		outHeadf("import (\n")
		outHeadf("\t%q\n", pkg)
		if hasTypeExport {
			outHeadf("\n\t\"qlang.io/qlang.spec.v1\"\n")
		}
		outHeadf(")\n\n")
	}

	// format
	data, err := format.Source(append(head.Bytes(), buf.Bytes()...))
	if err != nil {
		return err
	}

	// write file
	root := filepath.Join(outpath, pkg)
	err = os.MkdirAll(root, 0777)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(root, pkgName+".go"))
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(data)

	// write version
	for ver, lines := range verMap {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("// +build %s\n\n", ver))
		buf.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
		if verHasTypeExport[ver] {
			buf.WriteString("import (\n")
			buf.WriteString(fmt.Sprintf("\t%q\n\n", bp.ImportPath))
			buf.WriteString(fmt.Sprintf("\t%q\n", "qlang.io/qlang.spec.v1"))
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

// convert to lower case style, Name => name, NAME => NAME
func toLowerCaseStyle(s string) string {
	if len(s) <= 1 {
		return s
	}
	if unicode.IsLower(rune(s[1])) {
		return strings.ToLower(s[0:1]) + s[1:]
	}
	return s
}
