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
)

// -----------------------------------------------------------------------------

type Source interface {
	Fingerp() [20]byte // source code fingerprint
	IsDirty(outFile string, defctx bool) bool
	GenGo(outFile, modFile string) error
}

type Project struct {
	Source
	AutoGenFile   string // autogen file of output
	FriendlyFname string // friendly fname of source
	BuildArgs     []string
	ExecArgs      []string
	FlagNRINC     bool // do not run if not changed
	FlagRTOE      bool // remove tempfile on error
}

type Context struct {
	modfile string
	dir     string
	defctx  bool
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
	return &Context{modfile: modfile, dir: dir, defctx: defctx}
}

func (p *Context) GoCommand(op string, src *Project) GoCmd {
	out := p.out(src)
	if src.IsDirty(out.goFile, p.defctx) {
		if p.defctx {
			dir, _ := filepath.Split(out.goFile)
			os.Mkdir(dir, 0755)
		}
		if err := src.GenGo(out.goFile, p.modfile); err != nil {
			log.Panicln(err)
		}
	} else if src.FlagNRINC { // do not run if not changed
		return GoCmd{}
	}
	return goCommand(p.dir, op, &out)
}

type goTarget struct {
	goFile  string
	outFile string
	proj    *Project
	defctx  bool
}

func (p *Context) out(src *Project) (ret goTarget) {
	hash := src.Fingerp()
	fname := src.FriendlyFname
	if !strings.HasSuffix(fname, ".go") {
		fname += ".go"
	}
	dir, _ := filepath.Split(p.modfile)
	ret.outFile = dir + "g" + base64.RawURLEncoding.EncodeToString(hash[:])
	ret.proj = src
	ret.defctx = p.defctx
	if ret.defctx || src.AutoGenFile == "" {
		ret.goFile = ret.outFile + fname
	} else {
		ret.goFile = src.AutoGenFile
	}
	return
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
	fmt.Fprintf(&buf, gomodFormat, GOPVERSION, GOPROOT)
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
	err := runCommand(dir, command, args...)
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			os.Exit(e.ExitCode())
		default:
			log.Fatalln(err)
		}
	}
}

func runCommand(dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Dir = dir
	return cmd.Run()
}

func genDefaultGopMod(modfile string) {
	dir, _ := filepath.Split(modfile)
	dummy := dir + "dummy"
	os.MkdirAll(dummy, 0755)
	genGomodFile(modfile)
	genDummyProject(dummy)
	execCommand(dir, "go", "mod", "tidy")
}

// -----------------------------------------------------------------------------
