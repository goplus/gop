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
	"time"

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
		Source: &gopFiles{files: files},
	}
	if len(files) == 1 {
		file := files[0]
		srcDir, fname := filepath.Split(file)
		var autogen string
		if hasMultiFiles(srcDir, ".gop") {
			autogen = "gop_autogen_" + fname + ".go"
		} else {
			autogen = "gop_autogen.go"
		}
		proj.FriendlyFname = fname
		proj.AutoGenFile = srcDir + autogen
	}
	return
}

func hasMultiFiles(srcDir string, ext string) bool {
	var has bool
	if f, err := os.Open(srcDir); err == nil {
		defer f.Close()
		fis, _ := f.ReadDir(-1)
		for _, fi := range fis {
			if !fi.IsDir() && filepath.Ext(fi.Name()) == ext {
				if has {
					return true
				}
				has = true
			}
		}
	}
	return false
}

func (p *gopFiles) Fingerp() (*Fingerp, error) { // source code fingerprint
	var buf bytes.Buffer
	var lastModTime time.Time
	for _, file := range p.files {
		absfile, err := filepath.Abs(file)
		if err != nil {
			return nil, err
		}
		if buf.Len() >= 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(absfile)
		fi, err := os.Stat(absfile)
		if err != nil {
			return nil, err
		}
		modTime := fi.ModTime()
		if modTime.After(lastModTime) {
			lastModTime = modTime
		}
	}
	hash := sha1.Sum(buf.Bytes())
	return &Fingerp{Hash: hash, ModTime: lastModTime}, nil
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
