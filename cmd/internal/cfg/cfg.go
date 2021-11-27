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
package cfg

import (
	"os"
	"path/filepath"
)

var (
	GOMODCACHE = envOr("GOMODCACHE", gopathDir("pkg/mod"))
)

// envOr returns Getenv(key) if set, or else def.
func envOr(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}
	return val
}

// GetGoPath return the gopath
func GetGoPath() string {
	return os.Getenv("GOPATH")
}

func gopathDir(rel string) string {
	list := filepath.SplitList(GetGoPath())
	if len(list) == 0 || list[0] == "" {
		return ""
	}
	return filepath.Join(list[0], rel)
}
