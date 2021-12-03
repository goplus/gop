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
	"fmt"
	"path/filepath"

	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"

	"github.com/goplus/gop/env"
)

var (
	GOMODCACHE = env.GOMODCACHE()
)

func cacheDir(path string) (string, error) {
	enc, err := module.EscapePath(path)
	if err != nil {
		return "", err
	}
	return filepath.Join(GOMODCACHE, "cache/download", enc, "/@v"), nil
}

func CachePath(m module.Version, suffix string) (string, error) {
	dir, err := cacheDir(m.Path)
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
	return filepath.Join(dir, encVer+"."+suffix), nil
}
