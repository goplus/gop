//go:build ignore
// +build ignore

package make_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

const (
	inWindows = (runtime.GOOS == "windows")
)

var script = "all.bash"
var gopRoot = ""
var gopBinFiles = []string{"gop", "gopfmt"}
var installer = "cmd/make.go"
var versionFile = "VERSION"
var mainVersionFile = "env/version.go"

func checkPathExist(path string, isDir bool) bool {
	stat, err := os.Stat(path)
	isExists := !os.IsNotExist(err)
	if isDir {
		return isExists && stat.IsDir()
	}
	return isExists && !stat.IsDir()
}

func trimRight(s string) string {
	return strings.TrimRight(s, " \t\r\n")
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
	installer = filepath.Join(gopRoot, installer)
	versionFile = filepath.Join(gopRoot, versionFile)
	mainVersionFile = filepath.Join(gopRoot, mainVersionFile)

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

func getBranch() string {
	gitCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	b, _ := gitCmd.CombinedOutput()
	return trimRight(string(b))
}

func TestTagFlagInGitRepo(t *testing.T) {
	os.Chdir(gopRoot)

	// Setup
	tag := "v1.0.90"
	autoVersion := "v1.0.91"
	releaseBranch := "v1.0"
	currentBranch := getBranch()

	// Teardown
	defer func() {
		gitCmd := exec.Command("git", "checkout", currentBranch)
		gitCmd.CombinedOutput()
		gitCmd = exec.Command("git", "branch", "-d", releaseBranch)
		gitCmd.CombinedOutput()
		gitCmd = exec.Command("git", "tag", "-d", tag)
		gitCmd.CombinedOutput()
		gitCmd = exec.Command("git", "tag", "-d", autoVersion)
		gitCmd.CombinedOutput()
		if checkPathExist(versionFile, false) {
			os.Remove(versionFile)
		}
	}()

	t.Run("failed tag operation", func(t *testing.T) {
		cmd := exec.Command("go", "run", installer, "--tag", tag)
		if _, err := cmd.CombinedOutput(); err == nil {
			t.Fatal("Failed: Release on a non-release branch should be failed.")
		}
	})

	t.Run("successfull tag operation with specified tag", func(t *testing.T) {
		gitCmd := exec.Command("git", "checkout", "-b", releaseBranch)
		if output, err := gitCmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed: %v, output: %s\n", err, output)
		}

		cmd := exec.Command("go", "run", installer, "--tag", tag)
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed: Release tag: %s on branch: %s should not be failed", tag, releaseBranch)
		}

		if !checkPathExist(versionFile, false) {
			t.Fatal("Failed: a VERSION file not found.")
		}

		if data, _ := os.ReadFile(versionFile); string(data) != tag {
			t.Fatalf("Failed: content of VERSION file: '%s' not match tag: %s", data, tag)
		}

		gitCmd = exec.Command("git", "tag")
		if allTags, _ := gitCmd.CombinedOutput(); !strings.Contains(string(allTags), tag) {
			t.Fatalf("Failed: %s tag not found in tag list of this git repo.\n", tag)
		}
	})

	t.Run("successfull tag operation with auto generated tag", func(t *testing.T) {
		cmd := exec.Command("go", "run", installer, "--tag", "auto")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed: %v, out: %s\n", err, out)
		}

		if data, _ := os.ReadFile(versionFile); string(data) != autoVersion {
			t.Fatalf("Failed: content of VERSION file: '%s' not match tag: %s", data, autoVersion)
		}

		gitCmd := exec.Command("git", "tag")
		if allTags, _ := gitCmd.CombinedOutput(); !strings.Contains(string(allTags), autoVersion) {
			t.Fatalf("Failed: %s tag not found in tag list of this git repo.\n", autoVersion)
		}
	})

	t.Run("empty tag", func(t *testing.T) {
		cmd := exec.Command("go", "run", installer, "--tag", "")
		if out, err := cmd.CombinedOutput(); err != nil || !strings.Contains(string(out), "Usage") {
			t.Fatalf("Failed: %v, out: %s\n", err, out)
		}
	})
}

func TestTagFlagInNonGitRepo(t *testing.T) {
	os.Chdir(gopRoot)

	// Setup
	gitDir := filepath.Join(gopRoot, ".git")
	gitBackupDir := filepath.Join(gopRoot, ".gitBackup")

	// Rename .git dir
	if checkPathExist(gitDir, true) {
		os.Rename(gitDir, gitBackupDir)
	}

	// Teardown
	defer func() {
		if checkPathExist(gitBackupDir, true) {
			os.Rename(gitBackupDir, gitDir)
		}
	}()

	t.Run("specify new tag", func(t *testing.T) {
		cmd := exec.Command("go", "run", installer, "--tag", "v1.0.98")
		output, err := cmd.CombinedOutput()
		if err == nil || !strings.Contains(string(output), "Error") {
			t.Fatal("Failed: release on a non-git repo should be failed.")
		}
	})

	t.Run("specify auto tag", func(t *testing.T) {
		cmd := exec.Command("go", "run", installer, "--tag", "auto")
		output, err := cmd.CombinedOutput()
		if err == nil || !strings.Contains(string(output), "Error") {
			t.Fatal("Failed: release on a non-git repo should be failed.")
		}
	})
}

func TestInstallInNonGitRepo(t *testing.T) {
	os.Chdir(gopRoot)

	// Setup
	gitDir := filepath.Join(gopRoot, ".git")
	gitBackupDir := filepath.Join(gopRoot, ".gitBackup")

	// Rename .git dir
	if checkPathExist(gitDir, true) {
		os.Rename(gitDir, gitBackupDir)
	}

	// Teardown
	defer func() {
		if checkPathExist(gitBackupDir, true) {
			os.Rename(gitBackupDir, gitDir)
		}
		if checkPathExist(versionFile, false) {
			os.Remove(versionFile)
		}
	}()

	installCmd := func() *exec.Cmd {
		return exec.Command("go", "run", installer, "--install")
	}

	t.Run("failed build operation", func(t *testing.T) {
		cmd := installCmd()
		output, err := cmd.CombinedOutput()
		if err == nil || !strings.Contains(string(output), "Error") {
			t.Fatal("Failed: build Go+ in a non-git repo and no VERSION file should be failed.")
		}
	})

	t.Run("install with VERSION file", func(t *testing.T) {
		version := "v1.0.98"
		// Create VERSION file
		if err := os.WriteFile(versionFile, []byte(version), 0666); err != nil {
			t.Fatal(err)
		}

		cmd := installCmd()
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed: %v, output: %s\n", err, output)
		}

		cmd = exec.Command(filepath.Join(gopRoot, "bin", gopBinFiles[0]), "version")
		output, err := cmd.CombinedOutput()
		if err != nil || !strings.Contains(string(output), version) {
			t.Fatalf("Failed: %v, output: %s\n", err, output)
		}
	})
}

func TestReleaseMinorVersion(t *testing.T) {
	os.Chdir(gopRoot)

	// Setup
	nextVersion := "v1.98.1"
	smallNextVersion := "v1.98.2"
	bigNextVersion := "v1.98.98"
	mainVersion := "1.98"
	releaseBranch := "v" + mainVersion
	currentBranch := getBranch()
	// Modify gop/env.MainVersion to be test value
	mainVersionContent, _ := os.ReadFile(mainVersionFile)
	re := regexp.MustCompile(`(?sm)(^const \(.*?MainVersion = ")\d+.\d+(".*?\))`)
	newContent := re.ReplaceAll(mainVersionContent, []byte("${1}"+mainVersion+"${2}"))
	if err := os.WriteFile(mainVersionFile, newContent, 0666); err != nil {
		t.Fatalf("Error: %v, Modify %s file failed.\n", err, mainVersionFile)
	}

	// Teardown
	defer func() {
		gitCmd := exec.Command("git", "checkout", currentBranch)
		gitCmd.CombinedOutput()
		gitCmd = exec.Command("git", "branch", "-d", releaseBranch)
		gitCmd.CombinedOutput()
		gitCmd = exec.Command("git", "tag", "-d", nextVersion)
		gitCmd.CombinedOutput()
		gitCmd = exec.Command("git", "tag", "-d", smallNextVersion)
		gitCmd.CombinedOutput()
		gitCmd = exec.Command("git", "tag", "-d", bigNextVersion)
		gitCmd.CombinedOutput()
		if checkPathExist(versionFile, false) {
			os.Remove(versionFile)
		}
		if err := os.WriteFile(mainVersionFile, mainVersionContent, 0666); err != nil {
			t.Fatalf("Error: %v, Restore %s file failed.\n", err, mainVersionFile)
		}
	}()

	t.Run("auto generates new version", func(t *testing.T) {
		gitCmd := exec.Command("git", "checkout", "-b", releaseBranch)
		if output, err := gitCmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed: %v, output: %s\n", err, output)
		}

		cmd := exec.Command("go", "run", installer, "--tag", "auto")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Log(string(out))
			t.Fatalf("Failed: %v, auto release tag on branch: %s should not be failed", err, releaseBranch)
		}

		if !checkPathExist(versionFile, false) {
			t.Fatal("Failed: a VERSION file not found.")
		}

		if data, _ := os.ReadFile(versionFile); string(data) != nextVersion {
			t.Fatalf("Failed: content of VERSION file: '%s' not match tag: %s", data, nextVersion)
		}

		gitCmd = exec.Command("git", "tag")
		if allTags, _ := gitCmd.CombinedOutput(); !strings.Contains(string(allTags), nextVersion) {
			t.Fatalf("Failed: %s tag not found in tag list of this git repo.\n", nextVersion)
		}
	})

	t.Run("another auto generate operation", func(t *testing.T) {
		cmd := exec.Command("go", "run", installer, "--tag", "auto")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Log(string(out))
			t.Fatalf("Failed: %v, auto release tag on branch: %s should not be failed", err, releaseBranch)
		}

		if !checkPathExist(versionFile, false) {
			t.Fatal("Failed: a VERSION file not found.")
		}

		if data, _ := os.ReadFile(versionFile); string(data) != smallNextVersion {
			t.Fatalf("Failed: content of VERSION file: '%s' not match tag: %s", data, smallNextVersion)
		}

		gitCmd := exec.Command("git", "tag")
		if allTags, _ := gitCmd.CombinedOutput(); !strings.Contains(string(allTags), smallNextVersion) {
			t.Fatalf("Failed: %s tag not found in tag list of this git repo.\n", smallNextVersion)
		}
	})

	t.Run("manual specify new version", func(t *testing.T) {
		cmd := exec.Command("go", "run", installer, "--tag", bigNextVersion)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Log(string(out))
			t.Fatalf("Failed: Release tag: %s on branch: %s should not be failed", bigNextVersion, releaseBranch)
		}

		if !checkPathExist(versionFile, false) {
			t.Fatal("Failed: a VERSION file not found.")
		}

		if data, _ := os.ReadFile(versionFile); string(data) != bigNextVersion {
			t.Fatalf("Failed: content of VERSION file: '%s' not match tag: %s", data, bigNextVersion)
		}

		gitCmd := exec.Command("git", "tag")
		if allTags, _ := gitCmd.CombinedOutput(); !strings.Contains(string(allTags), bigNextVersion) {
			t.Fatalf("Failed: %s tag not found in tag list of this git repo.\n", bigNextVersion)
		}
	})
}

func TestHandleMultiFlags(t *testing.T) {
	os.Chdir(gopRoot)

	// Setup
	nextVersion := "v1.0.988"
	releaseBranch := "v1.0"
	currentBranch := getBranch()

	// Teardown
	defer func() {
		gitCmd := exec.Command("git", "checkout", currentBranch)
		gitCmd.CombinedOutput()
		gitCmd = exec.Command("git", "branch", "-d", releaseBranch)
		gitCmd.CombinedOutput()
		gitCmd = exec.Command("git", "tag", "-d", nextVersion)
		gitCmd.CombinedOutput()
		if checkPathExist(versionFile, false) {
			os.Remove(versionFile)
		}
	}()

	gitCmd := exec.Command("git", "checkout", "-b", releaseBranch)
	if output, err := gitCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed: %v, output: %s\n", err, output)
	}

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
		if !strings.Contains(string(allTags), nextVersion) {
			t.Fatalf("Failed: %s tag not found in this git repo.\n", nextVersion)
		}
	}

	if !checkPathExist(versionFile, false) {
		t.Fatal("Failed: a VERSION file not found.")
	}

	if data, _ := os.ReadFile(versionFile); string(data) != nextVersion {
		t.Fatalf("Failed: content of VERSION file: '%s' not match tag: %s", data, nextVersion)
	}
}
