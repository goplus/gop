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

package build

import (
	"fmt"
)

// The value of variables come form `go build -ldflags '-X "build.Date=xxxxx" -X "build.CommitID=xxxx"' `
var (
	// Date build time
	Date string
	// Branch current git branch
	Branch string
	// Commit git commit id
	Commit string
	// Default GOPROOT
	GopRoot string
)

// Build information
func Build() string {
	//if Date == "" {
	//	return ""
	//}
	return fmt.Sprintf("%s(%s) %s", Branch, Commit, Date)
}
