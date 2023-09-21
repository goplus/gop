package fsx

import (
	"io/fs"
	"os"
	"path/filepath"
)

// -----------------------------------------------------------------------------

// FileSystem represents a file system.
type FileSystem interface {
	ReadDir(dirname string) ([]fs.DirEntry, error)
	ReadFile(filename string) ([]byte, error)
	Join(elem ...string) string
}

// -----------------------------------------------------------------------------

type localFS struct{}

func (p localFS) ReadDir(dirname string) ([]fs.DirEntry, error) {
	return os.ReadDir(dirname)
}

func (p localFS) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (p localFS) Join(elem ...string) string {
	return filepath.Join(elem...)
}

var Local FileSystem = localFS{}

// -----------------------------------------------------------------------------
