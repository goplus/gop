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

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/qiniu/goplus/cmd/qexp/gopkg"
)

var (
	flagExportDir string
)

func init() {
	flag.StringVar(&flagExportDir, "outdir", "", "optional set export lib path, default goplus/lib path")
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: qexp [packages]\n")
		flag.PrintDefaults()
		return
	}
	var libDir string
	if flagExportDir != "" {
		libDir = flagExportDir
	} else {
		root, err := goplusRoot()
		if err != nil {
			fmt.Fprintln(os.Stderr, "find goplus root failed:", err)
			os.Exit(-1)
		}
		libDir = filepath.Join(root, "lib")
	}

	for _, pkgPath := range flag.Args() {
		err := exportPkg(pkgPath, libDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "export pkg %q failed, %v\n", pkgPath, err)
		} else {
			fmt.Fprintf(os.Stdout, "export pkg %q success\n", pkgPath)
		}
	}
}

func exportPkg(pkgPath string, libDir string) error {
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

func goplusRoot() (root string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for {
		modfile := filepath.Join(dir, "go.mod")
		if hasFile(modfile) {
			if isGoplus(modfile) {
				return dir, nil
			}
			return "", errors.New("current directory is not under goplus root")
		}
		next := filepath.Dir(dir)
		if dir == next {
			return "", errors.New("go.mod not found, please run under goplus root")
		}
		dir = next
	}
}

func isGoplus(modfile string) bool {
	b, err := ioutil.ReadFile(modfile)
	return err == nil && bytes.HasPrefix(b, goplusPrefix)
}

var (
	goplusPrefix = []byte("module github.com/qiniu/goplus")
)

func hasFile(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode().IsRegular()
}

// -----------------------------------------------------------------------------
