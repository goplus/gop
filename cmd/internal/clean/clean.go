/*
 * Copyright (c) 2021-2021 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package clean

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/qiniu/x/log"

	"github.com/goplus/gop/cmd/internal/base"
)

const (
	autoGenFile      = "gop_autogen.go"
	autoGenTestFile  = "gop_autogen_test.go"
	autoGen2TestFile = "gop_autogen2_test.go"
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
			pkgDir := filepath.Join(dir, fname)
			if fname == ".gop" {
				removeGopDir(pkgDir)
			} else {
				cleanAGFiles(pkgDir)
			}
			continue
		}
	}
	autogens := []string{autoGenFile, autoGenTestFile, autoGen2TestFile}
	for _, autogen := range autogens {
		file := filepath.Join(dir, autogen)
		if _, err = os.Stat(file); err == nil {
			fmt.Printf("Cleaning %s ...\n", file)
			os.Remove(file)
		}
	}
}

func removeGopDir(dir string) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, fi := range fis {
		fname := fi.Name()
		if strings.HasSuffix(fname, ".gop.go") {
			genfile := filepath.Join(dir, fname)
			fmt.Printf("Cleaning %s ...\n", genfile)
			os.Remove(genfile)
		}
	}
	os.Remove(dir)
}

// -----------------------------------------------------------------------------

// Cmd - gop clean
var Cmd = &base.Command{
	UsageLine: "gop clean [-v] <gopSrcDir>",
	Short:     "Clean all Go+ auto generated files",
}

var (
	flag = &Cmd.Flag
	_    = flag.Bool("v", false, "print verbose information.")
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(_ *base.Command, args []string) {
	err := flag.Parse(args)
	if err != nil {
		log.Fatalln("parse input arguments failed:", err)
	}
	var dir string
	if flag.NArg() == 0 {
		dir = "."
	} else {
		dir = flag.Arg(0)
	}
	cleanAGFiles(dir)
}

// -----------------------------------------------------------------------------
