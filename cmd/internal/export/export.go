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
	"encoding/json"
	"fmt"
	"go/format"
	"go/types"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/cmd/internal/export/gopkg"
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
		pkg, err := LookupPkg(pkgPath)
		if err != nil {
			fmt.Printf("lookup %q source error: %v\n", pkgPath, err)
			continue
		}
		fmt.Printf("source %v\n", pkg.Dir)
		if exporAll {
			exportDir(pkg.ImportPath, pkg.Dir)
			continue
		}
		if pkg.Goroot {
			err = exportPkg(pkg.ImportPath, "")
		} else {
			err = exportPkg(pkg.ImportPath, pkg.Dir)
		}
		if err != nil {
			fmt.Printf("export %q error: %v\n", pkg.ImportPath, err)
		} else {
			fmt.Printf("export %q success\n", pkg.ImportPath)
		}
	}
}

type Package struct {
	Dir        string // directory containing package sources
	ImportPath string // import path of package in dir
	Goroot     bool   // is this package in the Go root?
}

func LookupPkg(pkgPath string) (*Package, error) {
	// check go list
	if !strings.Contains(pkgPath, "@") {
		pkg, err := checkGoList(pkgPath)
		if err == nil {
			return pkg, nil
		}
	}
	// check go mod
	srcDir, err := gopkg.LookupMod(pkgPath)
	if err == nil {
		pkgPath, _, _ = gopkg.ParsePkgVer(pkgPath)
		rpkg, _ := checkGoPkg(srcDir)
		if rpkg != pkgPath {
			return nil, fmt.Errorf("unsupport replace module %q", rpkg)
		}
		return &Package{ImportPath: pkgPath, Dir: srcDir}, nil
	}
	// check download mod
	pkg, mod, sub := gopkg.ParsePkgVer(pkgPath)
	info, err := downloadMod(mod)
	if err != nil {
		return nil, fmt.Errorf("download %q failed, %v", mod, err)
	}
	dir := filepath.Join(info.Dir, sub)
	rpkg, _ := checkGoPkg(info.Dir)
	if rpkg != info.Path {
		return nil, fmt.Errorf("unsupport replace module %q", rpkg)
	}
	return &Package{ImportPath: pkg, Dir: dir}, nil
}

type ModInfo struct {
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

func downloadMod(pkgPath string) (*ModInfo, error) {
	fmt.Println("exec: go download mod", pkgPath)
	cmd := exec.Command(gobin, "mod", "download", "-json", pkgPath)
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var info ModInfo
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func exportDir(pkgPath string, dir string) error {
	cmd := exec.Command(gobin, "list", "-f", "{{.ImportPath}} {{.Dir}}", "./...")
	cmd.Dir = dir
	data, err := cmd.Output()
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		ar := strings.Split(line, " ")
		if len(ar) == 2 {
			pkg, dir := ar[0], ar[1]
			err := exportPkg(pkg, dir)
			if strings.Contains(pkg, "internal") || strings.Contains(pkg, "vendor") {
				continue
			}
			if err == nil {
				fmt.Printf("export %q success\n", pkg)
			} else if err != gopkg.ErrIgnore {
				fmt.Printf("export %v error,%v\n", pkg, err)
			}
		}
	}
	return nil
}

func checkGoList(pkgPath string) (*Package, error) {
	cmd := exec.Command(gobin, "list", "-json", pkgPath)
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var pkg Package
	err = json.Unmarshal(data, &pkg)
	if err != nil {
		return nil, err
	}
	return &pkg, nil
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

func exportPkg(pkgPath string, srcDir string) (err error) {
	var pkg *types.Package
	if srcDir == "" {
		pkg, err = gopkg.Import(pkgPath)
	} else {
		pkg, err = gopkg.ImportSource(pkgPath, srcDir)
	}
	if err != nil {
		return fmt.Errorf("import %q failed: %v", pkgPath, err)
	}
	if pkg.Name() == "main" {
		return gopkg.ErrIgnore
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
