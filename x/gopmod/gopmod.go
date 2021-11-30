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

package gopmod

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/env"
	"github.com/goplus/gop/env/execgo"
)

// -----------------------------------------------------------------------------

type Source interface {
	Fingerp() [20]byte // source code fingerprint
	IsDirty(outFile string, temp bool) bool
	GenGo(outFile string) error
}

type Project struct {
	Source
	AutoGenFile   string // autogen file of output
	FriendlyFname string // friendly fname of source
}

type Context struct {
	GOPMOD string
	defctx bool
}

func New(dir string) *Context {
	defctx := false
	modfile, err := env.GOPMOD(dir)
	if err != nil {
		modfile = env.HOME() + "/.gop/cmd/go.mod"
		if _, err := os.Stat(modfile); os.IsNotExist(err) {
			genDefaultGopMod(modfile)
		}
		defctx = true
	}
	return &Context{GOPMOD: modfile, defctx: defctx}
}

func (p *Context) GoCommand(op string, src *Project, args ...string) *exec.Cmd {
	file := p.goFile(src)
	cmd := execgo.Command(op, file, args...)
	cmd.Dir, _ = filepath.Split(file)
	return cmd
}

func (p *Context) outFile(src *Project) (outFile string, temp bool) {
	if p.defctx {
		hash := src.Fingerp()
		fname := src.FriendlyFname
		if !strings.HasSuffix(fname, ".go") {
			fname += ".go"
		}
		dir, _ := filepath.Split(p.GOPMOD)
		outFile = dir + "/g" + base64.RawURLEncoding.EncodeToString(hash[:]) + fname
		temp = true
	} else {
		outFile = src.AutoGenFile
	}
	return
}

func (p *Context) goFile(src *Project) string {
	outFile, temp := p.outFile(src)
	if src.IsDirty(outFile, temp) {
		if temp {
			dir, _ := filepath.Split(outFile)
			os.Mkdir(dir, 0755)
		}
		if err := src.GenGo(outFile); err != nil {
			log.Panicln(err)
		}
	}
	return outFile
}

// -----------------------------------------------------------------------------

const (
	dummyGoFile = `package dummy

import (
	_ "github.com/goplus/gop"
)
`
	gomodFormat = `module goplus.org/userapp

go 1.16

require (
	github.com/goplus/gop %s
)

replace (
	github.com/goplus/gop => %s
)
`
)

func genGomodFile(modfile string) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, gomodFormat, execgo.GOPVERSION, execgo.GOPROOT)
	err := os.WriteFile(modfile, buf.Bytes(), 0666)
	if err != nil {
		log.Panicln(err)
	}
}

func genDummyProject(dir string) {
	err := os.WriteFile(dir+"/dummy.go", []byte(dummyGoFile), 0666)
	if err != nil {
		log.Panicln(err)
	}
}

func execCommand(dir, command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			os.Exit(e.ExitCode())
		default:
			log.Fatalln(err)
		}
	}
}

func genDefaultGopMod(modfile string) {
	dir, _ := filepath.Split(modfile)
	dummy := dir + "/dummy"
	os.MkdirAll(dummy, 0755)
	genGomodFile(modfile)
	genDummyProject(dummy)
	execCommand(dir, "go", "mod", "tidy")
}

// -----------------------------------------------------------------------------
