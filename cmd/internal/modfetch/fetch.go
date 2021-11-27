package modfetch

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"

	"github.com/goplus/gop/cmd/internal/cfg"
)

// DownloadDir returns the directory to which m should have been downloaded.
// An error will be returned if the module path or version cannot be escaped.
// An error satisfying errors.Is(err, fs.ErrNotExist) will be returned
// along with the directory if the directory does not exist or if the directory
// is not completely populated.
func DownloadDir(m module.Version) (string, error) {
	enc, err := module.EscapePath(m.Path)
	if err != nil {
		return "", err
	}
	if !semver.IsValid(m.Version) {
		return "", fmt.Errorf("non-semver module version %q", m.Version)
	}
	if module.CanonicalVersion(m.Version) != m.Version {
		return "", fmt.Errorf("non-canonical module version %q", m.Version)
	}
	encVer, err := module.EscapeVersion(m.Version)
	if err != nil {
		return "", err
	}

	// Check whether the directory itself exists.
	dir := filepath.Join(cfg.GOMODCACHE, enc+"@"+encVer)
	if fi, err := os.Stat(dir); os.IsNotExist(err) {
		return dir, err
	} else if err != nil {
		return dir, &DownloadDirPartialError{dir, err}
	} else if !fi.IsDir() {
		return dir, &DownloadDirPartialError{dir, errors.New("not a directory")}
	}

	// Check if a .partial file exists. This is created at the beginning of
	// a download and removed after the zip is extracted.
	partialPath, err := CachePath(m, "partial")
	if err != nil {
		return dir, err
	}
	if _, err := os.Stat(partialPath); err == nil {
		return dir, &DownloadDirPartialError{dir, errors.New("not completely extracted")}
	} else if !os.IsNotExist(err) {
		return dir, err
	}

	// Check if a .ziphash file exists. It should be created before the
	// zip is extracted, but if it was deleted (by another program?), we need
	// to re-calculate it.
	ziphashPath, err := CachePath(m, "ziphash")
	if err != nil {
		return dir, err
	}
	if _, err := os.Stat(ziphashPath); os.IsNotExist(err) {
		return dir, &DownloadDirPartialError{dir, errors.New("ziphash file is missing")}
	} else if err != nil {
		return dir, err
	}
	return dir, nil
}

func Download(m module.Version) (dir string, err error) {
	dir, err = DownloadDir(m)
	if err == nil {
		// The directory has already been completely extracted (no .partial file exists).
		return dir, nil
	} else if dir == "" || !errors.Is(err, fs.ErrNotExist) {
		return "", err
	}

	DownloadArgs(dir, m.String())
	return dir, nil
}

func DownloadArgs(dir string, args ...string) {
	RunGoCmd(dir, "mod", append([]string{"download"}, args...)...)
}

func TidyArgs(dir string, args ...string) {
	RunGoCmd(dir, "mod", append([]string{"tidy"}, args...)...)
}

func InitArgs(dir string, args ...string) {
	RunGoCmd(dir, "mod", append([]string{"init"}, args...)...)
}

// RunGoCmd executes `go` command tools.
func RunGoCmd(dir string, op string, args ...string) {
	cmd := exec.Command("go", append([]string{op}, args...)...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			os.Exit(e.ExitCode())
		default:
			log.Fatalln("RunGoCmd failed:", err)
		}
	}
}

// DownloadDirPartialError is returned by DownloadDir if a module directory
// exists but was not completely populated.
//
// DownloadDirPartialError is equivalent to fs.ErrNotExist.
type DownloadDirPartialError struct {
	Dir string
	Err error
}

func (e *DownloadDirPartialError) Error() string     { return fmt.Sprintf("%s: %v", e.Dir, e.Err) }
func (e *DownloadDirPartialError) Is(err error) bool { return err == fs.ErrNotExist }
