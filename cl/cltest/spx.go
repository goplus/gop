/*
 * Copyright (c) 2024 The GoPlus Authors (goplus.org). All rights reserved.
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

package cltest

import (
	"bytes"
	"os"
	"testing"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/parser/fsx/memfs"
	"github.com/goplus/gop/scanner"
	"github.com/goplus/mod/modfile"
)

func spxParserConf() parser.Config {
	return parser.Config{
		ClassKind: func(fname string) (isProj bool, ok bool) {
			ext := modfile.ClassExt(fname)
			c, ok := LookupClass(ext)
			if ok {
				isProj = c.IsProj(ext, fname)
			}
			return
		},
	}
}

func Spx(t *testing.T, gmx, spxcode, expected string) {
	SpxEx(t, gmx, spxcode, expected, "index.tgmx", "bar.tspx")
}

func SpxEx(t *testing.T, gmx, spxcode, expected, gmxfile, spxfile string) {
	SpxWithConf(t, "gopSpxTest", Conf, gmx, spxcode, expected, gmxfile, spxfile, "")
}

func SpxEx2(t *testing.T, gmx, spxcode, expected, gmxfile, spxfile, resultFile string) {
	SpxWithConf(t, "gopSpxTest", Conf, gmx, spxcode, expected, gmxfile, spxfile, resultFile)
}

func SpxWithConf(t *testing.T, name string, conf *cl.Config, gmx, spxcode, expected, gmxfile, spxfile, resultFile string) {
	t.Run(name, func(t *testing.T) {
		cl.SetDisableRecover(true)
		defer cl.SetDisableRecover(false)

		fs := memfs.TwoFiles("/foo", spxfile, spxcode, gmxfile, gmx)
		if gmxfile == "" {
			fs = memfs.SingleFile("/foo", spxfile, spxcode)
		}
		pkgs, err := parser.ParseFSDir(Conf.Fset, fs, "/foo", spxParserConf())
		if err != nil {
			scanner.PrintError(os.Stderr, err)
			t.Fatal("ParseFSDir:", err)
		}
		bar := pkgs["main"]
		pkg, err := cl.NewPackage("", bar, conf)
		if err != nil {
			t.Fatal("NewPackage:", err)
		}
		var b bytes.Buffer
		err = pkg.WriteTo(&b, resultFile)
		if err != nil {
			t.Fatal("gogen.WriteTo failed:", err)
		}
		result := b.String()
		if result != expected {
			t.Fatalf("\nResult:\n%s\nExpected:\n%s\n", result, expected)
		}
	})
}

func SpxErrorEx(t *testing.T, msg, gmx, spxcode, gmxfile, spxfile string) {
	fs := memfs.TwoFiles("/foo", spxfile, spxcode, gmxfile, gmx)
	pkgs, err := parser.ParseFSDir(Conf.Fset, fs, "/foo", spxParserConf())
	if err != nil {
		scanner.PrintError(os.Stderr, err)
		t.Fatal("ParseFSDir:", err)
	}
	conf := *Conf
	conf.RelativeBase = "/foo"
	conf.Recorder = nil
	conf.NoFileLine = false
	bar := pkgs["main"]
	_, err = cl.NewPackage("", bar, &conf)
	if err == nil {
		t.Fatal("no error?")
	}
	if ret := err.Error(); ret != msg {
		t.Fatalf("\nError: \"%s\"\nExpected: \"%s\"\n", ret, msg)
	}
}
