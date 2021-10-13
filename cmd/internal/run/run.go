/*
 Copyright 2021 The GoPlus Authors (goplus.org)

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

// Package run implements the ``gop run'' command.
package run

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goplus/gox"
	"github.com/qiniu/x/log"
	"golang.org/x/tools/go/packages"
	"gop/ast"
	"gop/cl"
	"gop/cmd/gengo"
	"gop/cmd/internal/base"
	"gop/parser"
	"gop/scanner"
	"gop/token"
)

// -----------------------------------------------------------------------------

// Cmd - gop run
var Cmd = &base.Command{
	UsageLine: "gop run [-asm -quiet -debug -nr -gop -prof] <gopSrcDir|gopSrcFile>",
	Short:     "Run a Go+ program",
}

var (
	flag        = &Cmd.Flag
	flagAsm     = flag.Bool("asm", false, "generates `asm` code of Go+ bytecode backend")
	flagVerbose = flag.Bool("v", false, "print verbose information")
	flagQuiet   = flag.Bool("quiet", false, "don't generate any compiling stage log")
	flagDebug   = flag.Bool("debug", false, "set log level to debug")
	flagNorun   = flag.Bool("nr", false, "don't run if no change")
	flagGop     = flag.Bool("gop", false, "parse a .go file as a .gop file")
	flagProf    = flag.Bool("prof", false, "do profile and generate profile report")
)

func init() {
	Cmd.Run = runCmd
}

func saveGoFile(gofile string, pkg *gox.Package) error {
	dir := filepath.Dir(gofile)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	return gox.WriteFile(gofile, pkg, false)
}

func findGoModFile(dir string) (modfile string, noCacheFile bool, err error) {
	modfile, err = cl.FindGoModFile(dir)
	if err != nil {
		home := os.Getenv("HOME")
		modfile = home + "/gop/go.mod"
		if fi, e := os.Lstat(modfile); e == nil && !fi.IsDir() {
			return modfile, true, nil
		}
		modfile = home + "/goplus/go.mod"
		if fi, e := os.Lstat(modfile); e == nil && !fi.IsDir() {
			return modfile, true, nil
		}
	}
	return
}

func findGoModDir(dir string) (string, bool) {
	modfile, nocachefile, err := findGoModFile(dir)
	if err != nil {
		log.Fatalln("findGoModFile:", err)
	}
	return filepath.Dir(modfile), nocachefile
}

func runCmd(cmd *base.Command, args []string) {
	flag.Parse(args)
	if flag.NArg() < 1 {
		cmd.Usage(os.Stderr)
	}
	args = flag.Args()[1:]

	if *flagQuiet {
		log.SetOutputLevel(0x7000)
	} else if *flagDebug {
		log.SetOutputLevel(log.Ldebug)
		gox.SetDebug(gox.DbgFlagAll)
		cl.SetDebug(cl.DbgFlagAll)
	}
	if *flagVerbose {
		gox.SetDebug(gox.DbgFlagAll &^ gox.DbgFlagComments)
		cl.SetDebug(cl.DbgFlagAll)
		cl.SetDisableRecover(true)
	} else if *flagAsm {
		gox.SetDebug(gox.DbgFlagInstruction)
	}
	if *flagProf {
		panic("TODO: profile not impl")
	}

	fset := token.NewFileSet()
	src, _ := filepath.Abs(flag.Arg(0))
	fi, err := os.Stat(src)
	if err != nil {
		log.Fatalln("input arg check failed:", err)
	}
	isDir := fi.IsDir()

	var isDirty bool
	var srcDir, file, gofile string
	var pkgs map[string]*ast.Package
	if isDir {
		srcDir = src
		gofile = src + "/gop_autogen.go"
		isDirty = true // TODO: check if code changed
		if isDirty {
			pkgs, err = parser.ParseDir(fset, src, nil, 0)
		} else if *flagNorun {
			return
		}
	} else {
		srcDir, file = filepath.Split(src)
		isGo := filepath.Ext(file) == ".go"
		if isGo {
			hash := sha1.Sum([]byte(src))
			dir := os.Getenv("HOME") + "/.gop/run"
			os.MkdirAll(dir, 0777)
			gofile = dir + "/g" + base64.RawURLEncoding.EncodeToString(hash[:]) + file
		} else if hasMultiFiles(srcDir, ".gop") {
			gofile = filepath.Join(srcDir, "gop_autogen_"+file+".go")
		} else {
			gofile = srcDir + "/gop_autogen.go"
		}
		isDirty = fileIsDirty(fi, gofile)
		if isDirty {
			if isGo {
				fmt.Println("==> GenGo to", gofile)
			}
			if *flagGop {
				pkgs, err = parser.Parse(fset, src, nil, 0)
			} else {
				pkgs, err = parser.Parse(fset, src, nil, 0) // TODO: only to check dependencies
			}
		} else if *flagNorun {
			return
		}
	}
	if err != nil {
		scanner.PrintError(os.Stderr, err)
		os.Exit(10)
	}

	if isDirty {
		mainPkg, ok := pkgs["main"]
		if !ok {
			if len(pkgs) == 0 && isDir { // not a Go+ package, try runGoPkg
				runGoPkg(src, args, true)
				return
			}
			fmt.Fprintln(os.Stderr, "TODO: not a main package")
			os.Exit(12)
		}

		modDir, noCacheFile := findGoModDir(srcDir)
		conf := &cl.Config{
			Dir: modDir, TargetDir: srcDir, Fset: fset, CacheLoadPkgs: true, PersistLoadPkgs: !noCacheFile}
		out, err := cl.NewPackage("", mainPkg, conf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(11)
		}
		err = saveGoFile(gofile, out)
		if err != nil {
			log.Fatalln("saveGoFile failed:", err)
		}
		conf.PkgsLoader.Save()
	}

	goRun(gofile, args)
	if *flagProf {
		panic("TODO: profile not impl")
	}
}

func fileIsDirty(fi os.FileInfo, gofile string) bool {
	fiDest, err := os.Stat(gofile)
	if err != nil {
		return true
	}
	return fi.ModTime().After(fiDest.ModTime())
}

func goRun(file string, args []string) {
	goArgs := make([]string, len(args)+2)
	goArgs[0] = "run"
	goArgs[1] = file
	copy(goArgs[2:], args)
	cmd := exec.Command("go", goArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			os.Exit(e.ExitCode())
		default:
			log.Fatalln("go run failed:", err)
		}
	}
}

func runGoPkg(src string, args []string, doRun bool) {
	modfile, noCacheFile, err := findGoModFile(src)
	if err != nil {
		log.Fatalln("findGoModFile:", err)
	}
	base := filepath.Dir(modfile)
	rel, _ := filepath.Rel(base, src)
	modPath, err := cl.GetModulePath(modfile)
	if err != nil {
		log.Fatalln("GetModulePath:", err)
	}
	pkgPath := filepath.Join(modPath, rel)
	const (
		loadTypes = packages.NeedImports | packages.NeedDeps | packages.NeedTypes
		loadModes = loadTypes | packages.NeedName | packages.NeedModule
	)
	baseConf := &cl.Config{
		Fset:            token.NewFileSet(),
		GenGoPkg:        new(gengo.Runner).GenGoPkg,
		CacheLoadPkgs:   true,
		PersistLoadPkgs: !noCacheFile,
		NoFileLine:      true,
	}
	loadConf := &packages.Config{Mode: loadModes, Fset: baseConf.Fset}
	pkgs, err := baseConf.Ensure().PkgsLoader.Load(loadConf, pkgPath)
	if err != nil || len(pkgs) == 0 {
		log.Fatalln("PkgsLoader.Load failed:", err)
	}
	if pkgs[0].Name != "main" {
		fmt.Fprintln(os.Stderr, "TODO: not a main package")
		os.Exit(12)
	}
	baseConf.PkgsLoader.Save()
	if doRun {
		goRun(src+"/.", args)
	}
}

func hasMultiFiles(srcDir string, ext string) bool {
	var has bool
	if f, err := os.Open(srcDir); err == nil {
		defer f.Close()
		fis, _ := f.ReadDir(-1)
		for _, fi := range fis {
			if !fi.IsDir() && filepath.Ext(fi.Name()) == ext {
				if has {
					return true
				}
				has = true
			}
		}
	}
	return false
}

// -----------------------------------------------------------------------------
