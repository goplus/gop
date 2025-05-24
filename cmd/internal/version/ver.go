/*
 * Copyright (c) 2021 The XGo Authors (xgo.dev). All rights reserved.
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

package version

import (
	"fmt"
	"runtime"

	"github.com/goplus/xgo/cmd/internal/base"
	"github.com/goplus/xgo/env"
)

// -----------------------------------------------------------------------------

// Cmd - gop version
var Cmd = &base.Command{
	UsageLine: "gop version [-v]",
	Short:     "Print XGo version",
}

var (
	flag = &Cmd.Flag
	_    = flag.Bool("v", false, "print verbose information.")
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	fmt.Printf("gop %s %s/%s\n", env.Version(), runtime.GOOS, runtime.GOARCH)
}

// -----------------------------------------------------------------------------
