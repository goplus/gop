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
	"errors"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/env"
	"golang.org/x/mod/module"
)

// -----------------------------------------------------------------------------

var (
	GOMODCACHE = env.GOMODCACHE()
)

var (
	ErrNoNeedToDownload = errors.New("no need to download")
)

func DownloadCachePath(mod module.Version) (string, error) {
	if mod.Version == "" {
		return mod.Path, ErrNoNeedToDownload
	}
	encPath, err := module.EscapePath(mod.Path)
	if err != nil {
		return "", err
	}
	return filepath.Join(GOMODCACHE, "cache/download", encPath, "@v", mod.Version+".zip"), nil
}

func ModCachePath(mod module.Version) (string, error) {
	if mod.Version == "" {
		return mod.Path, nil
	}
	encPath, err := module.EscapePath(mod.Path)
	if err != nil {
		return "", err
	}
	return filepath.Join(GOMODCACHE, encPath+"@"+mod.Version), nil
}

func InModCachePath(path string) bool {
	if strings.HasPrefix(path, GOMODCACHE) {
		name := path[len(GOMODCACHE):]
		return name == "" || name[0] == '/' || name[0] == '\\'
	}
	return false
}

// -----------------------------------------------------------------------------
