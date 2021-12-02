/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
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

// Package gopfmt implements the ``gop fmt'' command.
package gopfmt

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/qiniu/x/log"

	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/format"

	xformat "github.com/goplus/gop/x/format"
)

// Cmd - gop go
var Cmd = &base.Command{
	UsageLine: "gop fmt [-n --smart] path ...",
	Short:     "Format Go+ packages",
}

var (
	flag        = &Cmd.Flag
	flagNotExec = flag.Bool("n", false, "prints commands that would be executed.")
	flagSmart   = flag.Bool("smart", false, "convert Go code style into Go+ style.")
)

func init() {
	Cmd.Run = runCmd
}

var (
	procCnt    = 0
	walkSubDir = false
	extGops    = map[string]struct{}{
		".go":  {},
		".gop": {},
		".spx": {},
		".gmx": {},
	}
	rootDir = ""
)

func gopfmt(path string, smart bool) (err error) {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	var target []byte
	if smart {
		target, err = xformat.Source(src)
	} else {
		target, err = format.Source(src)
	}
	if err != nil {
		return
	}
	if bytes.Equal(src, target) {
		return
	}
	fmt.Println(path)
	return writeFileWithBackup(path, target)
}

func writeFileWithBackup(path string, target []byte) (err error) {
	dir, file := filepath.Split(path)
	f, err := ioutil.TempFile(dir, file)
	if err != nil {
		return
	}
	tmpfile := f.Name()
	_, err = f.Write(target)
	f.Close()
	if err != nil {
		return
	}
	err = os.Remove(path)
	if err != nil {
		return
	}
	return os.Rename(tmpfile, path)
}

func walk(path string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else if d.IsDir() {
		if !walkSubDir && path != rootDir {
			return filepath.SkipDir
		}
	} else {
		ext := filepath.Ext(path)
		if _, ok := extGops[ext]; ok {
			procCnt++
			if *flagNotExec {
				fmt.Println("gop fmt", path)
			} else {
				err = gopfmt(path, *flagSmart && ext != ".go")
			}
		}
	}
	return err
}

func runCmd(cmd *base.Command, args []string) {
	err := flag.Parse(args)
	if err != nil {
		log.Fatalln("parse input arguments failed:", err)
	}
	narg := flag.NArg()
	if narg < 1 {
		cmd.Usage(os.Stderr)
	}
	for i := 0; i < narg; i++ {
		path := flag.Arg(i)
		walkSubDir = strings.HasSuffix(path, "/...")
		if walkSubDir {
			path = path[:len(path)-4]
		}
		procCnt = 0
		rootDir = path
		filepath.WalkDir(path, walk)
		if procCnt == 0 {
			fmt.Println("no Go+ files in", path)
		}
	}
}

// -----------------------------------------------------------------------------
