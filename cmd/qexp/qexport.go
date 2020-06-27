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
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/qiniu/goplus/cmd/qexp/gopkg"
)

var (
	exportFile string
)

func createExportFile(pkgDir string) (f io.WriteCloser, err error) {
	os.MkdirAll(pkgDir, 0777)
	exportFile = pkgDir + "/gomod_export.go"
	return os.Create(exportFile)
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: qexp <goPkgPath>\n")
		flag.PrintDefaults()
		return
	}
	pkgPath := flag.Arg(0)
	defer func() {
		if exportFile != "" {
			os.Remove(exportFile)
		}
	}()
	err := gopkg.Export(pkgPath, createExportFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "export failed:", err)
		os.Exit(1)
	}
	exportFile = "" // don't remove file if success
}

// -----------------------------------------------------------------------------
