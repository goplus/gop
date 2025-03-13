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
package mod

import (
	"fmt"
	"log"
	"os"

	"github.com/goplus/gop/cmd/internal/base"
)

var Cmd = &base.Command{
	UsageLine: "gop mod",
	Short:     "Module maintenance",

	Commands: []*base.Command{
		cmdInit,
		cmdDownload,
		cmdTidy,
	},
}

func check(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func fatal(msg any) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
