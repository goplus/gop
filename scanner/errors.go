// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner

import (
	"go/scanner"
	"io"
)

// Error is an alias of go/scanner.Error
//
type Error = scanner.Error

// ErrorList is an alias of go/scanner.ErrorList
//
type ErrorList = scanner.ErrorList

// PrintError is a utility function that prints a list of errors to w,
// one error per line, if the err parameter is an ErrorList. Otherwise
// it prints the err string.
//
func PrintError(w io.Writer, err error) {
	scanner.PrintError(w, err)
}
