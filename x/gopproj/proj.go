package gopproj

import (
	"errors"
	"path/filepath"
	"strings"
	"syscall"
)

// -----------------------------------------------------------------------------

type Proj = interface {
	projObj()
}

type FilesProj struct {
	Files []string
}

type PkgPathProj struct {
	Path string
}

type DirProj struct {
	Dir string
}

func (p *FilesProj) projObj()   {}
func (p *PkgPathProj) projObj() {}
func (p *DirProj) projObj()     {}

// -----------------------------------------------------------------------------

func ParseOne(args ...string) (proj Proj, next []string, err error) {
	if len(args) == 0 {
		return nil, nil, syscall.ENOENT
	}
	arg := args[0]
	if isFile(arg) {
		n := 1
		for n < len(args) && isFile(args[n]) {
			n++
		}
		return &FilesProj{Files: args[:n]}, args[n:], nil
	}
	if strings.HasPrefix(arg, ".") || strings.HasPrefix(arg, "/") {
		return &DirProj{Dir: arg}, args[1:], nil
	}
	return &PkgPathProj{Path: arg}, args[1:], nil
}

func isFile(fname string) bool {
	switch filepath.Ext(fname) {
	case ".gop", ".go":
		return true
	}
	return false
}

// -----------------------------------------------------------------------------

func ParseAll(args ...string) (projs []Proj, err error) {
	var hasFiles, hasNotFiles bool
	for {
		proj, next, e := ParseOne(args...)
		if e != nil {
			if hasFiles && hasNotFiles {
				return nil, ErrMixedFilesProj
			}
			return
		}
		if _, ok := proj.(*FilesProj); ok {
			hasFiles = true
		} else {
			hasNotFiles = true
		}
		projs = append(projs, proj)
		args = next
	}
}

var (
	ErrMixedFilesProj = errors.New("mixed files project")
)

// -----------------------------------------------------------------------------
