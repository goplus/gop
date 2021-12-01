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
	"crypto/sha1"
	"os"
	"path/filepath"
)

// -----------------------------------------------------------------------------

type goFile struct {
	file string
	code []byte
}

func openFromGoFile(file string) (proj *Project, err error) {
	code, err := os.ReadFile(file)
	if err != nil {
		return
	}
	proj = &Project{
		Source:        &goFile{file: file, code: code},
		AutoGenFile:   file,
		FriendlyFname: filepath.Base(file),
	}
	return
}

func (p *goFile) Fingerp() [20]byte { // source code fingerprint
	return sha1.Sum(p.code)
}

func (p *goFile) IsDirty(outFile string, temp bool) bool {
	return true
}

func (p *goFile) GenGo(outFile, modFile string) error {
	if p.file != outFile {
		return os.WriteFile(outFile, p.code, 0666)
	}
	return nil
}

// -----------------------------------------------------------------------------
