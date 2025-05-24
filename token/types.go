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

package token

import (
	"go/token"

	"github.com/goplus/xgo/token/internal/tokenutil"
)

// Pos is a compact encoding of a source position within a file set.
// It can be converted into a Position for a more convenient, but much
// larger, representation.
//
// The Pos value for a given file is a number in the range [base, base+size],
// where base and size are specified when adding the file to the file set via
// AddFile.
//
// To create the Pos value for a specific source offset (measured in bytes),
// first add the respective file to the current file set using FileSet.AddFile
// and then call File.Pos(offset) for that file. Given a Pos value p
// for a specific file set fset, the corresponding Position value is
// obtained by calling fset.Position(p).
//
// Pos values can be compared directly with the usual comparison operators:
// If two Pos values p and q are in the same file, comparing p and q is
// equivalent to comparing the respective source file offsets. If p and q
// are in different files, p < q is true if the file implied by p was added
// to the respective file set before the file implied by q.
type Pos = token.Pos

const (
	// NoPos - The zero value for Pos is NoPos; there is no file and line
	// information associated with it, and NoPos.IsValid() is false. NoPos
	// is always smaller than any other Pos value. The corresponding
	// Position value for NoPos is the zero value for Position.
	NoPos = token.NoPos
)

// Position describes an arbitrary source position
// including the file, line, and column location.
// A Position is valid if the line number is > 0.
type Position = token.Position

// A File is a handle for a file belonging to a FileSet.
// A File has a name, size, and line offset table.
type File = token.File

// Lines returns the effective line offset table of the form described by SetLines.
// Callers must not mutate the result.
func Lines(f *File) []int {
	return tokenutil.Lines(f)
}

// A FileSet represents a set of source files. Methods of file sets are
// synchronized; multiple goroutines may invoke them concurrently.
type FileSet = token.FileSet

// NewFileSet creates a new file set.
func NewFileSet() *FileSet {
	return token.NewFileSet()
}
