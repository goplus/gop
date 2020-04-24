package token

import (
	"go/token"
)

// -----------------------------------------------------------------------------

// Token is the set of lexical tokens of the qlang.
type Token = token.Token

// -----------------------------------------------------------------------------

// A FileSet represents a set of source files. Methods of file sets are synchronized;
// multiple goroutines may invoke them concurrently.
type FileSet = token.FileSet

// NewFileSet creates a new file set.
func NewFileSet() *FileSet {
	return token.NewFileSet()
}

// -----------------------------------------------------------------------------
