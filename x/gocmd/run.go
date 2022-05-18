/*
 * Copyright (c) 2022 The GoPlus Authors (goplus.org). All rights reserved.
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

package gocmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RunConfig struct {
	Gop   *GopEnv
	Flags []string
	Run   func(cmd *exec.Cmd) error
}

// -----------------------------------------------------------------------------

func RunDir(dir string, conf *RunConfig) (err error) {
	fis, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var files []string
	for _, fi := range fis {
		if !fi.IsDir() {
			if fname := fi.Name(); filterRunFname(fname) {
				files = append(files, filepath.Join(dir, fname))
			}
		}
	}
	return RunFiles(files, conf)
}

func filterRunFname(fname string) bool {
	return strings.HasSuffix(fname, ".go") &&
		!(strings.HasSuffix(fname, "_test.go") || strings.HasPrefix(fname, "_"))
}

// -----------------------------------------------------------------------------

func RunFiles(files []string, conf *RunConfig) (err error) {
	if conf == nil {
		conf = new(RunConfig)
	}
	exargs := make([]string, 1, 16)
	exargs[0] = "run"
	exargs = appendLdflags(exargs, conf.Gop)
	exargs = append(exargs, conf.Flags...)
	exargs = append(exargs, files...)
	cmd := exec.Command("go", exargs...)
	run := conf.Run
	if run == nil {
		run = (*exec.Cmd).Run
	}
	return run(cmd)
}

// -----------------------------------------------------------------------------
