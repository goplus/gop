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
	"go/build"
	"os"
	"path/filepath"
)

var (
	GOPATH = build.Default.GOPATH
)

func GOMODCACHE() string {
	val := os.Getenv("GOMODCACHE")
	if val == "" {
		val = gopathJoin("pkg/mod")
	}
	return val
}

func gopathJoin(rel string) string {
	list := filepath.SplitList(GOPATH)
	if len(list) == 0 || list[0] == "" {
		return ""
	}
	return filepath.Join(list[0], rel)
}
