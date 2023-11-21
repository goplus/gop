/*
 * Copyright (c) 2022 The GoPlus Authors (goplus.org). All rights reserved.
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

package cl_test

import (
	"path/filepath"
	"strings"
	"testing"
)

func lookupPub(pkgPath string) (pubfile string, err error) {
	relPath := strings.TrimPrefix(pkgPath, "github.com/goplus/gop/")
	return filepath.Join(gopRootDir, relPath, "c2go.a.pub"), nil
}

func TestHelloC2go(t *testing.T) {
	gopClTestC(t, `
import "C"

C.printf C"Hello, world!\n"
`, `package main

import (
	"github.com/goplus/gop/cl/internal/libc"
	"unsafe"
)

func main() {
	libc.Printf((*int8)(unsafe.Pointer(&[15]int8{'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', '!', '\n', '\x00'})))
}
`)
}

func TestHelloC2go2(t *testing.T) {
	gopClTestC(t, `
import "C/github.com/goplus/gop/cl/internal/libc"

C.printf C"Hello, world!\n"
`, `package main

import (
	"github.com/goplus/gop/cl/internal/libc"
	"unsafe"
)

func main() {
	libc.Printf((*int8)(unsafe.Pointer(&[15]int8{'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', '!', '\n', '\x00'})))
}
`)
}

func TestHelloC2go3(t *testing.T) {
	gopClTestC(t, `
import "C/libc"

C.printf C"Hello, world!\n"
`, `package main

import (
	"github.com/goplus/gop/cl/internal/libc"
	"unsafe"
)

func main() {
	libc.Printf((*int8)(unsafe.Pointer(&[15]int8{'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', '!', '\n', '\x00'})))
}
`)
}

func TestErrHelloC2go(t *testing.T) {
	codeErrorTestC(t, `./bar.gop:7:3: confliction: printf declared both in "github.com/goplus/gop/cl/internal/libc" and "github.com/goplus/gop/cl/internal/libc"`, `
import (
	"C"
	"C/libc"
)

C.printf C"Hello, world!\n"
`)
}
