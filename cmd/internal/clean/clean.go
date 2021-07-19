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

package clean

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/cmd/internal/base"
)

const (
	autoGenFileName = "gop_autogen.go"
)

// -----------------------------------------------------------------------------

func cleanAGFiles(dir string) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, fi := range fis {
		fname := fi.Name()
		if strings.HasPrefix(fname, "_") {
			continue
		}
		if fi.IsDir() {
			pkgDir := path.Join(dir, fname)
			cleanAGFiles(pkgDir)
			continue
		}
	}
	file := filepath.Join(dir, autoGenFileName)
	if _, err = os.Stat(file); err == nil {
		fmt.Printf("Cleaning %s ...\n", file)
		os.Remove(file)
	}
}

// -----------------------------------------------------------------------------

// Cmd - gop build
var Cmd = &base.Command{
	UsageLine: "gop clean [-v] <gopSrcDir>",
	Short:     "Clean all Go+ auto generated files",
}

var (
	flag        = &Cmd.Flag
	flagVerbose = flag.Bool("v", false, "print verbose information.")
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	flag.Parse(args)
	if flag.NArg() < 1 {
		cmd.Usage(os.Stderr)
		return
	}
	dir := flag.Arg(0)
	cleanAGFiles(dir)
}

// -----------------------------------------------------------------------------
