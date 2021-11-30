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
)

// -----------------------------------------------------------------------------

func HOME() string {
	return os.Getenv(envHOME)
}

// -----------------------------------------------------------------------------

func GOPATH() string {
	s := os.Getenv("GOPATH")
	if s == "" {
		return defaultGOPATH()
	}
	return s
}

func defaultGOPATH() string {
	if home := HOME(); home != "" {
		def := filepath.Join(home, "go")
		if filepath.Clean(def) == filepath.Clean(runtime.GOROOT()) {
			// Don't set the default GOPATH to GOROOT,
			// as that will trigger warnings from the go tool.
			return ""
		}
		return def
	}
	return ""
}

// -----------------------------------------------------------------------------

func GOMODCACHE() string {
	val := os.Getenv("GOMODCACHE")
	if val == "" {
		return gopathJoin("pkg/mod")
	}
	return val
}

func gopathJoin(rel string) string {
	list := filepath.SplitList(GOPATH())
	if len(list) == 0 || list[0] == "" {
		return ""
	}
	return filepath.Join(list[0], rel)
}

// -----------------------------------------------------------------------------
