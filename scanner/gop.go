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

package scanner

import (
	"go/scanner"
	"io"

	"github.com/goplus/gop/token"
	"github.com/qiniu/x/byteutil"
)

// -----------------------------------------------------------------------------

// Error is an alias of go/scanner.Error
type Error = scanner.Error

// ErrorList is an alias of go/scanner.ErrorList
type ErrorList = scanner.ErrorList

// PrintError is a utility function that prints a list of errors to w,
// one error per line, if the err parameter is an ErrorList. Otherwise
// it prints the err string.
func PrintError(w io.Writer, err error) {
	scanner.PrintError(w, err)
}

// -----------------------------------------------------------------------------

// New creates a scanner to tokenize the text src.
//
// Calls to Scan will invoke the error handler err if they encounter a
// syntax error and err is not nil. Also, for each error encountered,
// the Scanner field ErrorCount is incremented by one. The mode parameter
// determines how comments are handled.
func New(src string, err ErrorHandler, mode Mode) *Scanner {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	scanner := new(Scanner)
	scanner.Init(file, byteutil.Bytes(src), err, mode)
	return scanner
}

// -----------------------------------------------------------------------------
