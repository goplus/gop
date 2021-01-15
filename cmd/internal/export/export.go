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

// Package export implements the ``gop export'' command.
package export

import (
	"bytes"
	"fmt"
	"go/build"
	"go/format"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/cmd/internal/export/gopkg"
	"github.com/goplus/gop/cmd/internal/export/srcimporter"
)

var (
	flag          = &Cmd.Flag
	flagExportDir string
)

func init() {
	flag.StringVar(&flagExportDir, "outdir", "", "optional set export lib path, default is $GoPlusRoot/lib path")
}

// -----------------------------------------------------------------------------

// Cmd - gop go
var Cmd = &base.Command{
	UsageLine: "gop export [-outdir <outRootDir>] [packages]",
	Short:     "Export Go packages for Go+ programs",
}

func init() {
	Cmd.Run = runCmd
}

var (
	libDir string
	gobin  string
)

func init() {
	var err error
	gobin, err = exec.LookPath("go")
	if err != nil {
		panic("not found go bin in PATH")
	}
}

func runCmd(cmd *base.Command, args []string) {
	flag.Parse(args)
	if flag.NArg() < 1 {
		cmd.Usage(os.Stderr)
		return
	}
	if flagExportDir != "" {
		libDir = flagExportDir
	} else {
		root, err := gopkg.GoPlusRoot()
		if err != nil {
			fmt.Fprintln(os.Stderr, "find goplus root failed:", err)
			os.Exit(-1)
		}
		libDir = filepath.Join(root, "lib")
	}

	for _, pkgPath := range flag.Args() {
		var exporAll bool
		if strings.HasSuffix(pkgPath, "/...") {
			pkgPath = pkgPath[:len(pkgPath)-4]
			exporAll = true
		}
		pkgs, err := LookupPkgList(pkgPath, libDir, exporAll)
		if err != nil {
			fmt.Fprintf(os.Stderr, "go: lookup pkg %q error: %v\n", pkgPath, err)
			continue
		}
		var goMod string
		var exportList []string
		for _, pkg := range pkgs {
			if pkg.Name == "main" || isIgnorePkg(pkg.ImportPath) {
				continue
			}
			if pkg.Module != nil && pkg.Module.GoMod != goMod {
				goMod = pkg.Module.GoMod
				fmt.Fprintf(os.Stderr, "go: found module in %v\n", pkg.Module.Dir)
			}
			err := exportPkg(pkg.ImportPath, pkg.Dir, pkg.Goroot)
			if err == nil {
				exportList = append(exportList, pkg.ImportPath)
				fmt.Fprintf(os.Stdout, "export %q success\n", pkg.ImportPath)
			} else if err != gopkg.ErrIgnore {
				fmt.Fprintf(os.Stderr, "export %q error: %v\n", pkg.ImportPath, err)
			}
		}
		if len(exportList) == 0 {
			fmt.Fprintf(os.Stderr, "export %q error: empty exported package\n", pkgPath)
		}
	}
}

func isIgnorePkg(pkg string) bool {
	for _, a := range strings.Split(pkg, "/") {
		if a == "vendor" || a == "internal" {
			return true
		}
	}
	return false
}

func LookupPkgList(pkgPath string, workDir string, exporAll bool) ([]*jsonPackage, error) {
	pkg, path, _ := gopkg.ParsePkgVer(pkgPath)
	if strings.Contains(path, "@") {
		cmd := exec.Command(gobin, "get", path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
	}
	return checkGoPkgList(pkg, workDir, exporAll)
}

type jsonModInfo struct {
	Path     string // module path
	Version  string // module version
	Error    string // error loading module
	Info     string // absolute path to cached .info file
	GoMod    string // absolute path to cached .mod file
	Zip      string // absolute path to cached .zip file
	Dir      string // absolute path to cached source root directory
	Sum      string // checksum for path, version (as in go.sum)
	GoModSum string // checksum for go.mod (as in go.sum)
}

func exportPkg(pkgPath string, srcDir string, goRoot bool) (err error) {
	var pkg *types.Package
	if goRoot {
		pkg, err = gopkg.Import(pkgPath)
	} else {
		imp := srcimporter.New(&build.Default, token.NewFileSet(), make(map[string]*types.Package))
		pkg, err = imp.ImportFrom(pkgPath, srcDir, 0)
	}
	if err != nil {
		return fmt.Errorf("import %q failed: %v", pkgPath, err)
	}
	var buf bytes.Buffer
	err = gopkg.ExportPackage(pkg, &buf)
	if err != nil {
		return err
	}
	data, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println(buf.String())
		return err
	}
	dir := filepath.Join(libDir, pkg.Path())
	os.MkdirAll(dir, 0777)
	err = ioutil.WriteFile(filepath.Join(dir, "gomod_export.go"), data, 0666)
	return
}

// -----------------------------------------------------------------------------
