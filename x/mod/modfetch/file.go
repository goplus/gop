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
package modfetch

import (
	"io"
	"os"
)

// Edit creates the named file with mode 0666 (before umask),
// but does not truncate existing contents.
//
// If Edit succeeds, methods on the returned File can be used for I/O.
// The associated file descriptor has mode O_RDWR and the file is write-locked.
func Edit(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)
}

// Transform invokes t with the result of reading the named file, with its lock
// still held.
//
// If t returns a nil error, Transform then writes the returned contents back to
// the file, making a best effort to preserve existing contents on error.
//
// t must not modify the slice passed to it.
func Transform(name string, t func([]byte) ([]byte, error)) (err error) {
	f, err := Edit(name)
	if err != nil {
		return err
	}
	defer f.Close()

	old, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	new, err := t(old)
	if err != nil {
		return err
	}

	if len(new) > len(old) {
		// The overall file size is increasing, so write the tail first: if we're
		// about to run out of space on the disk, we would rather detect that
		// failure before we have overwritten the original contents.
		if _, err := f.WriteAt(new[len(old):], int64(len(old))); err != nil {
			// Make a best effort to remove the incomplete tail.
			f.Truncate(int64(len(old)))
			return err
		}
	}

	// We're about to overwrite the old contents. In case of failure, make a best
	// effort to roll back before we close the file.
	defer func() {
		if err != nil {
			if _, err := f.WriteAt(old, 0); err == nil {
				f.Truncate(int64(len(old)))
			}
		}
	}()

	if len(new) >= len(old) {
		if _, err := f.WriteAt(new[:len(old)], 0); err != nil {
			return err
		}
	} else {
		if _, err := f.WriteAt(new, 0); err != nil {
			return err
		}
		// The overall file size is decreasing, so shrink the file to its final size
		// after writing. We do this after writing (instead of before) so that if
		// the write fails, enough filesystem space will likely still be reserved
		// to contain the previous contents.
		if err := f.Truncate(int64(len(new))); err != nil {
			return err
		}
	}

	return nil
}

// Read opens the named file and returns its contents.
func Read(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}
