package search

import (
	"path/filepath"
	"strings"
)

// InDir checks whether path is in the file tree rooted at dir.
// If so, InDir returns an equivalent path relative to dir.
// If not, InDir returns an empty string.
// InDir makes some effort to succeed even in the presence of symbolic links.
// TODO(rsc): Replace internal/test.inDir with a call to this function for Go 1.12.
func InDir(path, dir string) string {
	if rel := inDirLex(path, dir); rel != "" {
		return rel
	}
	xpath, err := filepath.EvalSymlinks(path)
	if err != nil || xpath == path {
		xpath = ""
	} else {
		if rel := inDirLex(xpath, dir); rel != "" {
			return rel
		}
	}

	xdir, err := filepath.EvalSymlinks(dir)
	if err == nil && xdir != dir {
		if rel := inDirLex(path, xdir); rel != "" {
			return rel
		}
		if xpath != "" {
			if rel := inDirLex(xpath, xdir); rel != "" {
				return rel
			}
		}
	}
	return ""
}

// inDirLex is like inDir but only checks the lexical form of the file names.
// It does not consider symbolic links.
// TODO(rsc): This is a copy of str.HasFilePathPrefix, modified to
// return the suffix. Most uses of str.HasFilePathPrefix should probably
// be calling InDir instead.
func inDirLex(path, dir string) string {
	pv := strings.ToUpper(filepath.VolumeName(path))
	dv := strings.ToUpper(filepath.VolumeName(dir))
	path = path[len(pv):]
	dir = dir[len(dv):]
	switch {
	default:
		return ""
	case pv != dv:
		return ""
	case len(path) == len(dir):
		if path == dir {
			return "."
		}
		return ""
	case dir == "":
		return path
	case len(path) > len(dir):
		if dir[len(dir)-1] == filepath.Separator {
			if path[:len(dir)] == dir {
				return path[len(dir):]
			}
			return ""
		}
		if path[len(dir)] == filepath.Separator && path[:len(dir)] == dir {
			if len(path) == len(dir)+1 {
				return "."
			}
			return path[len(dir)+1:]
		}
		return ""
	}
}
