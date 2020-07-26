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

// Package export implements the ``gop download'' command.
package download

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
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
	UsageLine: "gop download [-outdir <outRootDir>] [packages]",
	Short:     "Download Go packages for Go+ programs",
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
		log.Fatalln("not found go bin in PATH", err)
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
		err := checkPkg(pkgPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "export pkg %q failed, %v\n", pkgPath, err)
		} else {
			fmt.Fprintf(os.Stdout, "export pkg %q success\n", pkgPath)
		}
	}
}

type ModInfo struct {
	Path     string
	Version  string
	Info     string
	GoMod    string
	Zip      string
	Dir      string
	Sum      string
	GoModSum string
}

func downMod(pkgPath string) error {
	cmd := exec.Command(gobin, "mod", "download", "-json", pkgPath)
	fmt.Println("go download mod", pkgPath)
	data, err := cmd.Output()
	if err != nil {
		return err
	}
	var info ModInfo
	err = json.Unmarshal(data, &info)
	if err != nil {
		return err
	}
	return exportDir(info)
}

func exportDir(info ModInfo) error {
	cmd := exec.Command(gobin, "list", "-f", "{{.ImportPath}} {{.Dir}}", "./...")
	cmd.Dir = info.Dir
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(data))
		return err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		ar := strings.Split(line, " ")
		if len(ar) == 2 {
			err := exportMod(ar[0], ar[1])
			if err == nil {
				fmt.Println("export mod", ar[1])
			} else if err != gopkg.ErrIgnore {
				fmt.Printf("export %v error,%v", ar[1], err)
			}
		}
	}

	return nil
}

func checkPkg(pkgPath string) error {
	if strings.Contains(pkgPath, "@") {
		return downMod(pkgPath)
	}
	return exportPkg(pkgPath)
}

func exportMod(pkgPath string, srcDir string) error {
	pkg, err := gopkg.ImportSource(pkgPath, srcDir)
	if err != nil {
		return err
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
	return err
}

func exportPkg(pkgPath string) error {
	pkg, err := gopkg.Import(pkgPath)
	if err != nil {
		return err
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
	return err
}

// -----------------------------------------------------------------------------
