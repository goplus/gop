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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goplus/gop/env"
)

// -----------------------------------------------------------------------------

const (
	ldFlagVersion   = "-X \"github.com/goplus/gop/env.buildVersion=%s\""
	ldFlagBuildDate = "-X \"github.com/goplus/gop/env.buildDate=%s\""
	ldFlagBuildRev  = "-X \"github.com/goplus/gop/env.buildCommit=%s\""
	ldFlagGopRoot   = "-X \"github.com/goplus/gop/env.defaultGopRoot=%s\""
)

const (
	ldFlagAll = ldFlagVersion + " " + ldFlagBuildDate + " " + ldFlagBuildRev + " " + ldFlagGopRoot
)

var (
	GOPVERSION   = env.Version()
	GOPBUILDDATE = env.BuildDate()
	GOPBUILDREV  = env.BuildRevision()
	GOPROOT      = env.GOPROOT()
)

func LoadFlags() string {
	return fmt.Sprintf(ldFlagAll, GOPVERSION, GOPBUILDDATE, GOPBUILDREV, GOPROOT)
}

// -----------------------------------------------------------------------------

type GoCmd struct {
	*exec.Cmd
	after func(error) error
}

func (p GoCmd) IsValid() bool {
	return p.Cmd != nil
}

func (p GoCmd) Run() error {
	err := p.Cmd.Run()
	if p.after != nil {
		return p.after(err)
	}
	return err
}

func goCommand(dir, op string, t *goTarget) (ret GoCmd) {
	proj := t.proj
	exargs := make([]string, 1, len(proj.BuildArgs)+len(proj.ExecArgs)+6)
	exargs[0] = op                             // 1
	exargs = append(exargs, proj.BuildArgs...) // len(proj.BuildArgs)
	exargs = appendLdflags(exargs, op)         // 2
	if op == "run" && t.defctx {               // 2
		afterDir, goFile, outFile := dir, t.goFile, t.outFile
		dir, _ = filepath.Split(goFile)
		exargs[0] = "build"
		exargs = append(exargs, "-o", outFile, goFile)
		ret.after = func(e error) error {
			if e == nil {
				e = runCommand(afterDir, outFile, proj.ExecArgs...)
				os.Remove(outFile)
			}
			if e != nil && t.proj.FlagRTOE { // remove tempfile on error
				os.Remove(goFile)
			}
			return e
		}
	} else {
		exargs = append(exargs, t.goFile)         // 1
		exargs = append(exargs, proj.ExecArgs...) // len(proj.ExecArgs)
	}
	ret.Cmd = exec.Command("go", exargs...)
	ret.Cmd.Dir = dir
	return
}

func appendLdflags(exargs []string, op string) []string {
	for _, v := range opsWithLdflags {
		if op == v {
			return append(exargs, "-ldflags", LoadFlags())
		}
	}
	return exargs
}

var (
	opsWithLdflags = []string{"run", "install", "build", "test"}
)

// -----------------------------------------------------------------------------
