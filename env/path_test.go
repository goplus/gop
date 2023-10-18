/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
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

package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goplus/mod"
)

func findGoModFile(dir string) (modfile string, noCacheFile bool, err error) {
	modfile, err = mod.GOPMOD(dir, 0)
	if err != nil {
		gopRoot, err := findGopRoot()
		if err == nil {
			modfile = filepath.Join(gopRoot, "go.mod")
			return modfile, true, nil
		}
	}
	return
}

// Common testing directory structure:
// testing_root/
//
//	  src/
//		   subdir/
//	  valid_goproot/
//		   go.mod
//		   go.sum
//	    cmd/gop/
func makeTestDir(t *testing.T) (root string, src string, gopRoot string) {
	root, _ = filepath.EvalSymlinks(t.TempDir())
	src = filepath.Join(root, "src")
	gopRoot = filepath.Join(root, "valid_goproot")
	makeValidGopRoot(gopRoot)
	os.Mkdir(src, 0755)
	return
}

func makeValidGopRoot(root string) {
	os.Mkdir(root, 0755)
	os.MkdirAll(filepath.Join(root, "cmd/gop"), 0755)
	os.WriteFile(filepath.Join(root, "go.mod"), []byte(""), 0644)
	os.WriteFile(filepath.Join(root, "go.sum"), []byte(""), 0644)
}

func writeDummyFile(path string) {
	os.WriteFile(path, []byte(""), 0644)
}

func cleanup() {
	os.Setenv("GOPROOT", "")
	os.Setenv(envHOME, "")
	defaultGopRoot = ""
}

func TestBasic(t *testing.T) {
	defaultGopRoot = ".."
	if GOPROOT() == "" {
		t.Fatal("TestBasic failed")
	}
	defaultGopRoot = ""
}

func TestFindGoModFileInGoModDir(t *testing.T) {
	cleanup()

	t.Run("the src/ is a valid mod dir", func(tt *testing.T) {
		tt.Cleanup(cleanup)
		_, src, _ := makeTestDir(tt)
		subdir := filepath.Join(src, "subdir")
		writeDummyFile(filepath.Join(src, "go.mod"))
		os.Mkdir(subdir, 0755)

		{
			modfile, noCacheFile, err := findGoModFile(src)

			if err != nil || modfile != filepath.Join(src, "go.mod") || noCacheFile {
				tt.Fatal("got:", modfile, noCacheFile, err)
			}
		}

		{
			// Should found go.mod in parent dir
			modfile, noCacheFile, err := findGoModFile(subdir)

			if err != nil || modfile != filepath.Join(src, "go.mod") || noCacheFile {
				tt.Fatal("got:", modfile, noCacheFile, err)
			}
		}
	})

	t.Run("the src/ is not a valid mod dir", func(tt *testing.T) {
		tt.Cleanup(cleanup)
		_, src, _ := makeTestDir(tt)

		modfile, noCacheFile, err := findGoModFile(src)

		if err == nil {
			tt.Fatal("should not found the mod file, but got:", modfile, noCacheFile)
		}
	})
}

func TestFindGoModFileInGopRoot(t *testing.T) {
	originDir, _ := os.Getwd()
	origiExecutable := executable
	home := filepath.Join(os.TempDir(), "test_home")
	os.Mkdir(home, 0755)
	t.Cleanup(func() {
		os.Chdir(originDir)
		os.RemoveAll(home)
		executable = origiExecutable
	})

	bin := filepath.Join(home, "bin")

	// Don't find go.mod in gop source dir when testing
	os.Chdir(home)

	cleanupAll := func() {
		cleanup()
		executable = func() (string, error) {
			return filepath.Join(bin, "run"), nil
		}
	}
	cleanupAll()

	t.Run("without gop root", func(tt *testing.T) {
		tt.Cleanup(cleanupAll)
		root, _, _ := makeTestDir(tt)

		modfile, noCacheFile, err := findGoModFile(root)

		if err == nil || noCacheFile || modfile != "" {
			tt.Fatal("should not found go.mod without gop root, got:", modfile, noCacheFile, err)
		}
	})

	t.Run("set GOPROOT to a valid goproot path", func(tt *testing.T) {
		tt.Cleanup(cleanupAll)
		_, src, gopRoot := makeTestDir(tt)

		os.Setenv("GOPROOT", gopRoot)
		modfile, noCacheFile, err := findGoModFile(src)

		if err != nil || modfile != filepath.Join(gopRoot, "go.mod") || !noCacheFile {
			tt.Fatal("should found mod file in GOPROOT, got:", modfile, noCacheFile, err)
		}
	})

	t.Run("set GOPROOT to an invalid goproot path", func(tt *testing.T) {
		tt.Cleanup(cleanupAll)
		root, src, _ := makeTestDir(tt)
		invalidGopRoot := filepath.Join(root, "invalid_goproot")

		defer func() {
			r := recover()
			if r == nil {
				tt.Fatal("should panic, but not")
			}
		}()

		os.Setenv("GOPROOT", invalidGopRoot)
		findGoModFile(src)
	})

	t.Run("set defaultGopRoot to a valid goproot path", func(tt *testing.T) {
		tt.Cleanup(cleanupAll)
		_, src, gopRoot := makeTestDir(tt)

		defaultGopRoot = gopRoot
		modfile, noCacheFile, err := findGoModFile(src)

		if err != nil || modfile != filepath.Join(gopRoot, "go.mod") || !noCacheFile {
			tt.Fatal("should found go.mod in the dir of defaultGopRoot, got:", modfile, noCacheFile, err)
		}
	})

	t.Run("set defaultGopRoot to an invalid path", func(tt *testing.T) {
		tt.Cleanup(cleanupAll)
		root, src, _ := makeTestDir(tt)
		invalidGopRoot := filepath.Join(root, "invalid_goproot")

		defaultGopRoot = invalidGopRoot
		{
			modfile, noCacheFile, err := findGoModFile(src)

			if err == nil || noCacheFile || modfile != "" {
				tt.Fatal("should not found go.mod when defaultGopRoot isn't exists, got:", modfile, noCacheFile, err)
			}
		}

		{
			os.Mkdir(invalidGopRoot, 0755)
			modfile, noCacheFile, err := findGoModFile(src)

			if err == nil || noCacheFile || modfile != "" {
				tt.Fatal("should not found go.mod when defaultGopRoot isn't an valid gop root dir, got:", modfile, noCacheFile, err)
			}
		}
	})

	t.Run("use $HOME/gop or $HOME/goplus", func(tt *testing.T) {
		tt.Cleanup(cleanupAll)
		root, src, _ := makeTestDir(tt)
		home := filepath.Join(root, "home")
		os.Mkdir(home, 0755)
		os.Setenv(envHOME, home)

		{
			gopRoot := filepath.Join(home, "goplus")
			makeValidGopRoot(gopRoot)

			modfile, noCacheFile, err := findGoModFile(src)

			if err != nil || !noCacheFile || modfile != filepath.Join(gopRoot, "go.mod") {
				tt.Fatal("should found go.mod in $HOME/goplus, but got:", modfile, noCacheFile, err)
			}
		}

		{
			gopRoot := filepath.Join(home, "gop")
			makeValidGopRoot(gopRoot)

			modfile, noCacheFile, err := findGoModFile(src)

			if err != nil || !noCacheFile || modfile != filepath.Join(gopRoot, "go.mod") {
				tt.Fatal("should found go.mod in $HOME/gop, but got:", modfile, noCacheFile, err)
			}
		}
	})

	t.Run("check if parent dir of the executable is valid gop root", func(tt *testing.T) {
		tt.Cleanup(cleanupAll)
		_, src, gopRoot := makeTestDir(tt)
		bin := filepath.Join(gopRoot, "bin")
		exePath := filepath.Join(bin, "run")
		os.Mkdir(bin, 0755)
		writeDummyFile(exePath)

		// Mock executable location
		executable = func() (string, error) {
			return exePath, nil
		}

		modfile, noCacheFile, err := findGoModFile(src)

		if err != nil || !noCacheFile || modfile != filepath.Join(gopRoot, "go.mod") {
			tt.Fatal("should found go.mod in gopRoot, but got:", modfile, noCacheFile, err)
		}
	})

	t.Run("test gop root priority", func(tt *testing.T) {
		tt.Cleanup(cleanupAll)
		root, src, _ := makeTestDir(tt)

		tt.Run("without gop root", func(tt *testing.T) {
			modfile, noCacheFile, err := findGoModFile(src)

			if err == nil || noCacheFile || modfile != "" {
				tt.Fatal("should not found go.mod without gop root, got:", modfile, noCacheFile, err)
			}
		})

		tt.Run("set HOME but hasn't $HOME/gop/ and $HOME/goplus/", func(tt *testing.T) {
			os.Setenv(envHOME, root)

			modfile, noCacheFile, err := findGoModFile(src)

			if err == nil || noCacheFile || modfile != "" {
				tt.Fatal("should not found go.mod without gop root, got:", modfile, noCacheFile, err)
			}
		})

		tt.Run("set HOME, and has valid $HOME/goplus/", func(tt *testing.T) {
			gopRoot := filepath.Join(root, "goplus")
			makeValidGopRoot(gopRoot)

			modfile, noCacheFile, err := findGoModFile(src)

			if err != nil || !noCacheFile || modfile != filepath.Join(gopRoot, "go.mod") {
				tt.Fatal("should found go.mod in $HOME/goplus, but got:", modfile, noCacheFile, err)
			}
		})

		tt.Run("set HOME, and has valid $HOME/gop/", func(tt *testing.T) {
			gopRoot := filepath.Join(root, "gop")
			makeValidGopRoot(gopRoot)

			modfile, noCacheFile, err := findGoModFile(src)

			if err != nil || !noCacheFile || modfile != filepath.Join(gopRoot, "go.mod") {
				tt.Fatal("should found go.mod in $HOME/gop, but got:", modfile, noCacheFile, err)
			}
		})

		tt.Run("set defaultGopRoot to an invalid gop root dir", func(tt *testing.T) {
			gopRoot := filepath.Join(root, "gop")

			defaultGopRoot = filepath.Join(root, "invalid_goproot")
			modfile, noCacheFile, err := findGoModFile(src)

			if err != nil || !noCacheFile || modfile != filepath.Join(gopRoot, "go.mod") {
				tt.Fatal("should found go.mod in $HOME/gop, but got:", modfile, noCacheFile, err)
			}
		})

		tt.Run("set defaultGopRoot to a valid gop root dir", func(tt *testing.T) {
			newGopRoot := filepath.Join(root, "new_gop_root")
			makeValidGopRoot(newGopRoot)

			defaultGopRoot = newGopRoot
			modfile, noCacheFile, err := findGoModFile(src)

			if err != nil || !noCacheFile || modfile != filepath.Join(newGopRoot, "go.mod") {
				tt.Fatal("should found go.mod in new_gop_root/, but got:", modfile, noCacheFile, err)
			}
		})

		tt.Run("the executable's parent dir is a valid gop root dir", func(tt *testing.T) {
			newGopRoot2 := filepath.Join(root, "new_gop_root2")
			makeValidGopRoot(newGopRoot2)
			bin := filepath.Join(newGopRoot2, "bin")
			exePath := filepath.Join(bin, "run")
			os.Mkdir(bin, 0755)
			writeDummyFile(exePath)
			// Mock executable location
			executable = func() (string, error) {
				return exePath, nil
			}

			modfile, noCacheFile, err := findGoModFile(src)

			if err != nil || !noCacheFile || modfile != filepath.Join(newGopRoot2, "go.mod") {
				tt.Fatal("should found go.mod in new_gop_root2/, but got:", modfile, noCacheFile, err)
			}
		})

		tt.Run("set GOPROOT to an invalid gop root dir", func(tt *testing.T) {
			newGopRoot3 := filepath.Join(root, "new_gop_root3")

			defer func() {
				r := recover()
				if r == nil {
					tt.Fatal("should panic, but not")
				}
			}()

			os.Setenv("GOPROOT", newGopRoot3)
			findGoModFile(src)
		})

		tt.Run("set GOPROOT to a valid gop root dir", func(tt *testing.T) {
			newGopRoot4 := filepath.Join(root, "new_gop_root4")
			makeValidGopRoot(newGopRoot4)

			os.Setenv("GOPROOT", newGopRoot4)
			modfile, noCacheFile, err := findGoModFile(src)

			if err != nil || !noCacheFile || modfile != filepath.Join(newGopRoot4, "go.mod") {
				tt.Fatal("should found go.mod in new_gop_root3/, but got:", modfile, noCacheFile, err)
			}
		})
	})
}
