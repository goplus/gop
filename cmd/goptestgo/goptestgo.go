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

package main

import (
	"bytes"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// RunGopCmd executes `gop` command tools.
func RunGopCmd(dir string, op string, args ...string) {
	opwargs := make([]string, len(args)+1)
	opwargs[0] = op
	copy(opwargs[1:], args)
	cmd := exec.Command("gop", opwargs...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			os.Exit(e.ExitCode())
		default:
			log.Fatalln("RunGopCmd failed:", err)
		}
	}
}

var (
	goRunPrefix = []byte("// run\n")
)

func gopTestRunGo(dir string) {
	filepath.Walk(dir, func(file string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return nil
		}
		ext := filepath.Ext(fi.Name())
		if ext != ".go" {
			return nil
		}
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Panicln(err)
		}
		if !bytes.HasPrefix(data, goRunPrefix) {
			return nil
		}
		log.Println("==> Gop run", file)
		RunGopCmd("", "run", "-gop", file)
		return nil
	})
}

// goptestgo: run all $GOROOT/test/*.go
func main() {
	dir := filepath.Join(build.Default.GOROOT, "test")
	gopTestRunGo(dir)
}
