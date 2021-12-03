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
// gop mod tidy
package mod

import (
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/cmd/internal/modload"
	"github.com/goplus/gop/x/mod/modfetch"
)

var cmdTidy = &base.Command{
	UsageLine: "gop mod tidy [-e] [-v]",
	Short:     "add missing and remove unused modules",
}

func init() {
	cmdTidy.Run = runTidy
}

func runTidy(cmd *base.Command, args []string) {
	modload.LoadModFile()
	modload.SyncGoMod()
	modfetch.TidyArgs(".", args...)
	modload.SyncGopMod()
}
