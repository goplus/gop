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
	"strings"
)

const (
	MainVersion = "1.0"
)

// buildVersion is the GoPlus tree's version string at build time.
// This is set by the linker.
var (
	buildVersion string
)

func init() {
	if buildVersion != "" && !strings.HasPrefix(buildVersion, "v"+MainVersion+".") {
		panic("Invalid buildVersion: " + buildVersion)
	}
}

// Version returns the GoPlus tree's version string.
// It is either the commit hash and date at the time of the build or,
// when possible, a release tag like "v1.0.0-rc1".
func Version() string {
	if buildVersion == "" {
		return "v" + MainVersion + ".x"
	}
	return buildVersion
}
