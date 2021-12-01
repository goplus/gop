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
	"crypto/sha1"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

// -----------------------------------------------------------------------------

type gopFiles struct {
	files []string
}

func (p *Context) openFromGopFiles(files []string) (proj *Project, err error) {
	proj = &Project{
		Source:        &gopFiles{files: files},
		UseDefaultCtx: true,
	}
	return
}

func (p *gopFiles) Fingerp() [20]byte { // source code fingerprint
	files := make([]string, len(p.files))
	for i, file := range p.files {
		files[i], _ = filepath.Abs(file)
	}
	sort.Strings(files)
	var buf bytes.Buffer
	for _, file := range files {
		fi, err := os.Stat(file)
		if err != nil {
			log.Panicln(err)
		}
		buf.WriteString(file)
		buf.WriteByte('\t')
		buf.WriteString(fi.ModTime().UTC().String())
		buf.WriteByte('\n')
	}
	return sha1.Sum(buf.Bytes())
}

func (p *gopFiles) IsDirty(outFile string, temp bool) bool {
	if fi, err := os.Lstat(outFile); err == nil { // TODO: more strictly
		return fi.IsDir()
	}
	return true
}

const (
	parserMode = parser.ParseComments
)

func (p *gopFiles) GenGo(outFile, modFile string) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFiles(fset, p.files, parserMode)
	if err != nil {
		return err
	}
	if len(pkgs) != 1 {
		log.Panicln("TODO: mutli packages -", len(pkgs))
	}
	mainPkg, ok := pkgs["main"]
	if !ok {
		panic("TODO: main package not found")
	}

	srcDir, _ := filepath.Split(outFile)
	modDir, _ := filepath.Split(modFile)
	conf := &cl.Config{
		Dir: modDir, TargetDir: srcDir, Fset: fset, CacheLoadPkgs: true, PersistLoadPkgs: true}
	out, err := cl.NewPackage("", mainPkg, conf)
	if err != nil {
		return err
	}
	err = gox.WriteFile(outFile, out, false)
	if err != nil {
		return err
	}
	conf.PkgsLoader.Save()
	return nil
}

// -----------------------------------------------------------------------------
