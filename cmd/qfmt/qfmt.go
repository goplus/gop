/*
 Copyright 2020 damonchen (netubu@gmail.com)

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

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	//"bytes"
	"flag"
	"fmt"
	"go/printer"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/qiniu/goplus/format"
	"github.com/qiniu/goplus/scanner"
)

var (
	// main operation modes
	write = flag.Bool("w", false, "write result to (source) file instead of stdout")
)

const (
	tabWidth    = 8
	printerMode = printer.UseSpaces | printer.TabIndent
)

var (
	exitCode = 0
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: qfmt [flags] [path ...]\n")
	flag.PrintDefaults()
}

func isGopFile(f os.FileInfo) bool {
	// ignore non-Gop files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".gop")
}

func report(err error) {
	scanner.PrintError(os.Stderr, err)
	exitCode = 2
}

func visitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isGopFile(f) {
		err = processFile(path, nil, os.Stdout, false)
	}
	// Don't complain if a file was deleted in the meantime (i.e.
	// the directory changed concurrently while running gofmt).
	if err != nil && !os.IsNotExist(err) {
		report(err)
	}
	return nil
}

const chmodSupported = runtime.GOOS != "windows"

// backupFile writes data to a new file named filename<number> with permissions perm,
// with <number randomly chosen such that the file name is unique. backupFile returns
// the chosen file name.
func backupFile(filename string, data []byte, perm os.FileMode) (string, error) {
	// create backup file
	f, err := ioutil.TempFile(filepath.Dir(filename), filepath.Base(filename))
	if err != nil {
		return "", err
	}
	bakname := f.Name()
	if chmodSupported {
		err = f.Chmod(perm)
		if err != nil {
			f.Close()
			os.Remove(bakname)
			return bakname, err
		}
	}

	// write data to backup file
	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}

	return bakname, err
}

// If in == nil, the source is the contents of the file with the given filename.
func processFile(filename string, in io.Reader, out io.Writer, stdin bool) error {
	var perm os.FileMode = 0644
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			return err
		}
		in = f
		perm = fi.Mode().Perm()
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	res, err := format.Source(src)
	if err != nil {
		return err
	}

	if *write {
		// make a temporary backup before overwriting original
		bakname, err := backupFile(filename+".", src, perm)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filename, res, perm)
		if err != nil {
			os.Rename(bakname, filename)
			return err
		}
		err = os.Remove(bakname)
		if err != nil {
			return err
		}
	}

	if !*write {
		_, err = out.Write(res)
	}

	return err
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func _main() {
	flag.Usage = usage
	flag.Parse()

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			report(err)
		case dir.IsDir():
			walkDir(path)
		default:
			if err := processFile(path, nil, os.Stdout, false); err != nil {
				report(err)
			}
		}
	}
}

func main() {
	_main()
	os.Exit(exitCode)
}
