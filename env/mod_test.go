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

package env

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const (
	sep = string(os.PathSeparator)
)

func TestGopMod(t *testing.T) {
	file, err := GOPMOD("")
	if err != nil {
		t.Fatal("GOPMOD failed:", err)
	}
	if !strings.HasSuffix(file, "gop"+sep+"go.mod") {
		t.Fatal("GOPMOD failed:", file)
	}
	file, err = GOPMOD("../cl/internal/spx")
	if err != nil {
		t.Fatal("GOPMOD spx failed:", err)
	}
	if !strings.HasSuffix(file, "cl"+sep+"internal"+sep+"gop.mod") {
		t.Fatal("GOPMOD spx failed:", file)
	}
}

func TestGOPATH(t *testing.T) {
	os.Setenv("GOPATH", "")
	gopath := GOPATH()
	home := HOME()
	expect := filepath.Join(home, "go")
	if expect == filepath.Clean(runtime.GOROOT()) {
		if gopath != "" {
			t.Fatal("TestGOPATH failed:", gopath)
		}
	} else if gopath != expect {
		t.Fatal("TestGOPATH failed:", gopath)
	}
	const ugopath = "/abc/work"
	os.Setenv("GOPATH", ugopath)
	gopath = GOPATH()
	if gopath != ugopath {
		t.Fatal("TestGOPATH failed:", gopath)
	}

	os.Setenv("GOPATH", "")
	os.Setenv(envHOME, "")
	defer os.Setenv(envHOME, home)
	if gopath := GOPATH(); gopath != "" {
		t.Fatal("TestGOPATH failed:", gopath)
	}
}

func TestGOMODCACHE(t *testing.T) {
	gopath := GOPATH()
	os.Setenv("GOMODCACHE", "")
	modcache := GOMODCACHE()
	if filepath.Join(gopath, "pkg/mod") != modcache {
		t.Fatal("TestGOMODCACHE failed:", modcache)
	}
	const umodcache = "/abc/mod"
	os.Setenv("GOMODCACHE", umodcache)
	modcache = GOMODCACHE()
	if modcache != umodcache {
		t.Fatal("TestGOMODCACHE (umodcache) failed:", modcache)
	}
	/* TODO:
	os.Setenv("GOMODCACHE", "")
	os.Setenv("GOPATH", "")
	home := HOME()
	defer os.Setenv(envHOME, home)
	os.Setenv(envHOME, strings.TrimRight(runtime.GOROOT(), "go"))
	modcache = GOMODCACHE()
	if modcache != "" {
		t.Fatal("TestGOMODCACHE (GOROOT) failed:", modcache)
	}
	*/
}
