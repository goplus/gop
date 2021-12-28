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

package modfetch

import (
	"bytes"
	"os/exec"
	"strings"

	"golang.org/x/mod/module"
)

/*
import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

// downloadDir returns the directory to which m should have been downloaded.
// An error will be returned if the module path or version cannot be escaped.
// An error satisfying errors.Is(err, fs.ErrNotExist) will be returned
// along with the directory if the directory does not exist or if the directory
// is not completely populated.
func downloadDir(m module.Version) (string, error) {
	enc, err := module.EscapePath(m.Path)
	if err != nil {
		return "", err
	}
	if !semver.IsValid(m.Version) {
		return "", fmt.Errorf("non-semver module version %q", m.Version)
	}
	if module.CanonicalVersion(m.Version) != m.Version {
		return "", fmt.Errorf("non-canonical module version %q", m.Version)
	}
	encVer, err := module.EscapeVersion(m.Version)
	if err != nil {
		return "", err
	}

	// Check whether the directory itself exists.
	dir := filepath.Join(GOMODCACHE, enc+"@"+encVer)
	if fi, err := os.Stat(dir); os.IsNotExist(err) {
		return dir, err
	} else if err != nil {
		return dir, &DownloadDirPartialError{dir, err}
	} else if !fi.IsDir() {
		return dir, &DownloadDirPartialError{dir, errors.New("not a directory")}
	}

	// Check if a .partial file exists. This is created at the beginning of
	// a download and removed after the zip is extracted.
	partialPath, err := CachePath(m, "partial")
	if err != nil {
		return dir, err
	}
	if _, err := os.Stat(partialPath); err == nil {
		return dir, &DownloadDirPartialError{dir, errors.New("not completely extracted")}
	} else if !os.IsNotExist(err) {
		return dir, err
	}

	// Check if a .ziphash file exists. It should be created before the
	// zip is extracted, but if it was deleted (by another program?), we need
	// to re-calculate it.
	ziphashPath, err := CachePath(m, "ziphash")
	if err != nil {
		return dir, err
	}
	if _, err := os.Stat(ziphashPath); os.IsNotExist(err) {
		return dir, &DownloadDirPartialError{dir, errors.New("ziphash file is missing")}
	} else if err != nil {
		return dir, err
	}
	return dir, nil
}

func Download(m module.Version) (dir string, err error) {
	dir, err = downloadDir(m)
	if err == nil {
		// The directory has already been completely extracted (no .partial file exists).
		return dir, nil
	} else if dir == "" || !errors.Is(err, fs.ErrNotExist) {
		return "", err
	}

	DownloadArgs(dir, m.String())
	return dir, nil
}

func DownloadArgs(dir string, args ...string) {
	runCmd(dir, "go", append([]string{"mod", "download"}, args...)...)
}

func TidyArgs(dir string, args ...string) {
	runCmd(dir, "go", append([]string{"mod", "tidy"}, args...)...)
}

func InitArgs(dir string, args ...string) {
	runCmd(dir, "go", append([]string{"mod", "init"}, args...)...)
}

// runCmd executes a command tool.
func runCmd(dir string, name string, args ...string) {
	cmd := exec.Command(name, args...)
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
			log.Fatalln("RunGoCmd failed:", err)
		}
	}
}

// DownloadDirPartialError is returned by DownloadDir if a module directory
// exists but was not completely populated.
//
// DownloadDirPartialError is equivalent to fs.ErrNotExist.
type DownloadDirPartialError struct {
	Dir string
	Err error
}

func (e *DownloadDirPartialError) Error() string     { return fmt.Sprintf("%s: %v", e.Dir, e.Err) }
func (e *DownloadDirPartialError) Is(err error) bool { return err == fs.ErrNotExist }
*/

// -----------------------------------------------------------------------------

func Get(modPath string, action func(act string, mod module.Version)) (mod module.Version, err error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("go", "get", modPath)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	mod = getResult(stderr.String(), action)
	if err != nil {
		err = &ExecCmdError{Err: err, Stderr: stderr.Bytes()}
	}
	return
}

func getResult(data string, action func(act string, mod module.Version)) (mod module.Version) {
	// go: downloading github.com/xushiwei/foogop v0.1.0
	const downloading = "go: downloading "
	if strings.HasPrefix(data, downloading) {
		mod = moduleAct("download", data[len(downloading):], action, &data)
	}
	// go get: added github.com/xushiwei/foogop v0.1.0
	const added = "go get: added "
	if strings.HasPrefix(data, added) {
		mod = moduleAct("add", data[len(added):], action, &data)
	}
	return
}

func moduleAct(
	act, data string, action func(act string, mod module.Version), next *string) (mod module.Version) {
	if pos := strings.IndexByte(data, '\n'); pos > 0 {
		line := data[:pos]
		*next = data[pos+1:]
		if pos = strings.IndexByte(line, ' '); pos > 0 {
			mod.Path, mod.Version = line[:pos], line[pos+1:]
			if action != nil {
				action(act, mod)
			}
		}
	}
	return
}

// -----------------------------------------------------------------------------

type ExecCmdError struct {
	Err    error
	Stderr []byte
}

func (p *ExecCmdError) Error() string {
	if e := p.Stderr; e != nil {
		return string(e)
	}
	return p.Err.Error()
}

// -----------------------------------------------------------------------------
