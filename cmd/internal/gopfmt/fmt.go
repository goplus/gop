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

// Package gopfmt implements the ``gop fmt'' command.
package gopfmt

import (
	"github.com/goplus/gop/cmd/internal/base"
)

/*
import (
	"fmt"
	"go/printer"
	"go/scanner"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/format"
)

const (
	tabWidth    = 8
	printerMode = printer.UseSpaces | printer.TabIndent
)

var (
	// 0-do nothing,1-formatted,2-error
	exitCode = 0
)

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
	backupName := f.Name()
	if chmodSupported {
		err = f.Chmod(perm)
		if err != nil {
			f.Close()
			os.Remove(backupName)
			return backupName, err
		}
	}

	// write data to backup file
	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}

	return backupName, err
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

	if *write && string(src) != string(res) {
		exitCode = 1
		// make a temporary backup before overwriting original
		backupName, err := backupFile(filename+".", src, perm)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filename, res, perm)
		if err != nil {
			os.Rename(backupName, filename)
			return err
		}
		err = os.Remove(backupName)
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

func qfmtMain(args []string) {
	flag.Parse(args)
	narg := flag.NArg()
	if narg == 0 {
		if *write {
			fmt.Fprintln(os.Stderr, "error: cannot use -w with standard input")
			exitCode = 2
			return
		}
		if err := processFile("<standard input>", os.Stdin, os.Stdout, true); err != nil {
			report(err)
		}
		return
	}
	for i := 0; i < narg; i++ {
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
*/
// -----------------------------------------------------------------------------

// Cmd - gop go
var Cmd = &base.Command{
	UsageLine: "gop fmt [-n] path ...",
	Short:     "Format Go+ packages",
}

var (
	flag        = &Cmd.Flag
	flagNotExec = flag.Bool("n", false, "prints commands that would be executed.")
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	// call qfmtMain in a separate function
	// so that it can use defer and have them
	// run before the exit.
	/*
		qfmtMain(args)
		os.Exit(exitCode)
	*/
	panic("TODO: gop fmt not impl")
}

// -----------------------------------------------------------------------------
