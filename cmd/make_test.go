//go:build ignore
// +build ignore

package make_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	inWindows     = (runtime.GOOS == "windows")
	installerFile = "cmd/make.go"
)

var script = "all.bash"
var gopRoot = ""
var gopBinFiles = []string{"gop", "gopfmt"}
var installer = ""

func checkPathExist(path string, isDir bool) bool {
	stat, err := os.Stat(path)
	isExists := !os.IsNotExist(err)
	if isDir {
		return isExists && stat.IsDir()
	}
	return isExists && !stat.IsDir()
}

func detectGoBinPath() string {
	goBin, ok := os.LookupEnv("GOBIN")
	if ok {
		return goBin
	}

	goPath, ok := os.LookupEnv("GOPATH")
	if ok {
		list := filepath.SplitList(goPath)
		if len(list) > 0 {
			// Put in first directory of $GOPATH.
			return filepath.Join(list[0], "bin")
		}
	}

	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, "go", "bin")
}

func init() {
	pwd, _ := os.Getwd()
	gopRoot = filepath.Join(pwd, "..")
	installer = filepath.Join(gopRoot, installerFile)
	if inWindows {
		script = "all.bat"
		for index, file := range gopBinFiles {
			file += ".exe"
			gopBinFiles[index] = file
		}
	}
}

func cleanGopRunCacheFiles(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	runCacheDir := filepath.Join(homeDir, ".gop", "run")
	files := []string{"go.mod", "go.sum"}
	for _, file := range files {
		fullPath := filepath.Join(runCacheDir, file)
		if checkPathExist(fullPath, false) {
			t.Fatalf("Failed: %s found in %s directory\n", file, runCacheDir)
		}
	}
}

func TestAllScript(t *testing.T) {
	os.Chdir(gopRoot)
	cmd := exec.Command(filepath.Join(gopRoot, script))
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed: %v:\nOut: %s\n", err, output)
	}

	goBinPath := detectGoBinPath()

	for _, file := range gopBinFiles {
		if !checkPathExist(filepath.Join(gopRoot, "bin", file), false) {
			t.Fatalf("Failed: %s not found in ./bin directory\n", file)
		}
		if !checkPathExist(filepath.Join(goBinPath, file), false) {
			t.Fatalf("Failed: %s not found in %s/bin directory\n", file, goBinPath)
		}
	}

	cleanGopRunCacheFiles(t)

	cmd = exec.Command(filepath.Join(gopRoot, "bin", gopBinFiles[0]), "version")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed: %v:\nOut: %s\n", err, output)
	}
}

func TestHandleMultiFlags(t *testing.T) {
	os.Chdir(gopRoot)
	nextVersion := "v1.2.3"
	cmd := exec.Command("go", "run", installer, "--install", "--test", "--uninstall", "--tag", nextVersion)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed: %v:\nOut: %s\n", err, output)
	}

	goBinPath := detectGoBinPath()

	// Uninstall will be the last action to run, test if all build assets being cleared.
	for _, file := range gopBinFiles {
		if checkPathExist(filepath.Join(gopRoot, "bin", file), false) {
			t.Fatalf("Failed: %s found in ./bin directory\n", file)
		}
		if checkPathExist(filepath.Join(goBinPath, file), false) {
			t.Fatalf("Failed: %s found in %s/bin directory\n", file, goBinPath)
		}
	}

	// Test if git tag list contains nextVersion
	cmd = exec.Command("git", "tag")
	if allTags, err := cmd.CombinedOutput(); err != nil {
		if !strings.Contains(allTags, nextVersion) {
			t.Fatalf("Failed: %s tag not found in this git repo.\n", nextVersion)
		}
	}
}

func TestTagFlagOnGitRepo(t *testing.T) {
	os.Chdir(gopRoot)

	// Specify next tag
	cmd := exec.Command("go", "run", installer, "--tag", "v1.0.98")

	// Specify tag: auto
	cmd := exec.Command("go", "run", installer, "--tag", "auto")

	// Specify bad tag
	cmd := exec.Command("go", "run", installer, "--tag", "1.23.45")

	// Specify empty tag
	cmd := exec.Command("go", "run", installer, "--tag", "")

	// Without specified tag
	cmd := exec.Command("go", "run", installer, "--tag")
}

func TestTagFlagOnNonGitRepo(t *testing.T) {
	os.Chdir(gopRoot)

	// Specify next tag
	cmd := exec.Command("go", "run", installer, "--tag", "v1.0.98")

	// Specify tag: auto
	cmd := exec.Command("go", "run", installer, "--tag", "auto")

	// Specify bad tag
	cmd := exec.Command("go", "run", installer, "--tag", "1.23.45")

	// Specify empty tag
	cmd := exec.Command("go", "run", installer, "--tag", "")

	// Without specified tag
	cmd := exec.Command("go", "run", installer, "--tag")
}
