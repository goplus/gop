/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

// Package gop provides all interfaces for implementing a Go+ package from
// existed Go packages.
package gop

import (
	exec "github.com/goplus/gop/exec/bytecode"
)

// -----------------------------------------------------------------------------

// A Context represents the context of an executor.
type Context = exec.Context

// -----------------------------------------------------------------------------

// GoPackage represents a Go package.
type GoPackage = exec.GoPackage

// NewGoPackage creates a new builtin Go Package.
func NewGoPackage(pkgPath string) *GoPackage {
	return exec.NewGoPackage(pkgPath)
}

// NewGoPackageEx creates a new builtin Go Package.
func NewGoPackageEx(pkgPath string, name string) *GoPackage {
	return exec.NewGoPackageEx(pkgPath, name)
}

// ToBools converts []interface{} into []bool.
func ToBools(args []interface{}) []bool {
	ret := make([]bool, len(args))
	for i, arg := range args {
		ret[i] = arg.(bool)
	}
	return ret
}

// ToStrings converts []interface{} into []string.
func ToStrings(args []interface{}) []string {
	ret := make([]string, len(args))
	for i, arg := range args {
		ret[i] = arg.(string)
	}
	return ret
}

// ToError converts a value into error.
func ToError(v interface{}) error {
	if v == nil {
		return nil
	}
	return v.(error)
}

// -----------------------------------------------------------------------------
