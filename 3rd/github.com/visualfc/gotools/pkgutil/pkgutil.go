package pkgutil

import (
	"go/build"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var go15VendorExperiment = os.Getenv("GO15VENDOREXPERIMENT") == "1"

func IsVendorExperiment() bool {
	return go15VendorExperiment
}

// matchPattern(pattern)(name) reports whether
// name matches pattern.  Pattern is a limited glob
// pattern in which '...' means 'any string' and there
// is no other special syntax.
func matchPattern(pattern string) func(name string) bool {
	re := regexp.QuoteMeta(pattern)
	re = strings.Replace(re, `\.\.\.`, `.*`, -1)
	// Special case: foo/... matches foo too.
	if strings.HasSuffix(re, `/.*`) {
		re = re[:len(re)-len(`/.*`)] + `(/.*)?`
	}
	reg := regexp.MustCompile(`^` + re + `$`)
	return func(name string) bool {
		return reg.MatchString(name)
	}
}

// hasPathPrefix reports whether the path s begins with the
// elements in prefix.
func hasPathPrefix(s, prefix string) bool {
	switch {
	default:
		return false
	case len(s) == len(prefix):
		return s == prefix
	case len(s) > len(prefix):
		if prefix != "" && prefix[len(prefix)-1] == '/' {
			return strings.HasPrefix(s, prefix)
		}
		return s[len(prefix)] == '/' && s[:len(prefix)] == prefix
	}
}

// hasFilePathPrefix reports whether the filesystem path s begins with the
// elements in prefix.
func hasFilePathPrefix(s, prefix string) bool {
	sv := strings.ToUpper(filepath.VolumeName(s))
	pv := strings.ToUpper(filepath.VolumeName(prefix))
	s = s[len(sv):]
	prefix = prefix[len(pv):]
	switch {
	default:
		return false
	case sv != pv:
		return false
	case len(s) == len(prefix):
		return s == prefix
	case len(s) > len(prefix):
		if prefix != "" && prefix[len(prefix)-1] == filepath.Separator {
			return strings.HasPrefix(s, prefix)
		}
		return s[len(prefix)] == filepath.Separator && s[:len(prefix)] == prefix
	}
}

// treeCanMatchPattern(pattern)(name) reports whether
// name or children of name can possibly match pattern.
// Pattern is the same limited glob accepted by matchPattern.
func treeCanMatchPattern(pattern string) func(name string) bool {
	wildCard := false
	if i := strings.Index(pattern, "..."); i >= 0 {
		wildCard = true
		pattern = pattern[:i]
	}
	return func(name string) bool {
		return len(name) <= len(pattern) && hasPathPrefix(pattern, name) ||
			wildCard && strings.HasPrefix(name, pattern)
	}
}

var isDirCache = map[string]bool{}

func isDir(path string) bool {
	result, ok := isDirCache[path]
	if ok {
		return result
	}

	fi, err := os.Stat(path)
	result = err == nil && fi.IsDir()
	isDirCache[path] = result
	return result
}

type Package struct {
	Root       string
	Dir        string
	ImportPath string
}

func ImportFile(fileName string) *Package {
	return ImportDir(filepath.Dir(fileName))
}

func ImportDir(dir string) *Package {
	pkg, err := build.ImportDir(dir, build.FindOnly)
	if err != nil {
		return &Package{"", dir, ""}
	}
	return &Package{pkg.Root, pkg.Dir, pkg.ImportPath}
}

// vendoredImportPath returns the expansion of path when it appears in parent.
// If parent is x/y/z, then path might expand to x/y/z/vendor/path, x/y/vendor/path,
// x/vendor/path, vendor/path, or else stay x/y/z if none of those exist.
// vendoredImportPath returns the expanded path or, if no expansion is found, the original.
// If no expansion is found, vendoredImportPath also returns a list of vendor directories
// it searched along the way, to help prepare a useful error message should path turn
// out not to exist.
func VendoredImportPath(parent *Package, path string) (found string, searched []string) {
	if parent == nil || parent.Root == "" || !go15VendorExperiment {
		return path, nil
	}
	dir := filepath.Clean(parent.Dir)
	root := filepath.Join(parent.Root, "src")
	if !hasFilePathPrefix(dir, root) || len(dir) <= len(root) || dir[len(root)] != filepath.Separator {
		//log.Fatalf("invalid vendoredImportPath: dir=%q root=%q separator=%q", dir, root, string(filepath.Separator))
		return path, nil
	}
	vpath := "vendor/" + path
	for i := len(dir); i >= len(root); i-- {
		if i < len(dir) && dir[i] != filepath.Separator {
			continue
		}
		// Note: checking for the vendor directory before checking
		// for the vendor/path directory helps us hit the
		// isDir cache more often. It also helps us prepare a more useful
		// list of places we looked, to report when an import is not found.
		if !isDir(filepath.Join(dir[:i], "vendor")) {
			continue
		}
		targ := filepath.Join(dir[:i], vpath)
		if isDir(targ) {
			// We started with parent's dir c:\gopath\src\foo\bar\baz\quux\xyzzy.
			// We know the import path for parent's dir.
			// We chopped off some number of path elements and
			// added vendor\path to produce c:\gopath\src\foo\bar\baz\vendor\path.
			// Now we want to know the import path for that directory.
			// Construct it by chopping the same number of path elements
			// (actually the same number of bytes) from parent's import path
			// and then append /vendor/path.
			chopped := len(dir) - i
			if chopped == len(parent.ImportPath)+1 {
				// We walked up from c:\gopath\src\foo\bar
				// and found c:\gopath\src\vendor\path.
				// We chopped \foo\bar (length 8) but the import path is "foo/bar" (length 7).
				// Use "vendor/path" without any prefix.
				return vpath, nil
			}
			return parent.ImportPath[:len(parent.ImportPath)-chopped] + "/" + vpath, nil
		}
		// Note the existence of a vendor directory in case path is not found anywhere.
		searched = append(searched, targ)
	}
	return path, searched
}

// findVendor looks for the last non-terminating "vendor" path element in the given import path.
// If there isn't one, findVendor returns ok=false.
// Otherwise, findInternal returns ok=true and the index of the "vendor".
//
// Note that terminating "vendor" elements don't count: "x/vendor" is its own package,
// not the vendored copy of an import "" (the empty import path).
// This will allow people to have packages or commands named vendor.
// This may help reduce breakage, or it may just be confusing. We'll see.
func findVendor(path string) (index int, ok bool) {
	// Two cases, depending on internal at start of string or not.
	// The order matters: we must return the index of the final element,
	// because the final one is where the effective import path starts.
	switch {
	case strings.Contains(path, "/vendor/"):
		return strings.LastIndex(path, "/vendor/") + 1, true
	case strings.HasPrefix(path, "vendor/"):
		return 0, true
	}
	return 0, false
}

func VendorPathToImportPath(path string) string {
	if i, ok := findVendor(path); ok {
		return path[i+len("vendor/"):]
	}
	return path
}
