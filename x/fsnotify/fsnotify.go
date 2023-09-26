/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
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

package fsnotify

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var (
	debugEvent bool
)

const (
	DbgFlagEvent = 1 << iota
	DbgFlagAll   = DbgFlagEvent
)

func SetDebug(dbgFlags int) {
	debugEvent = (dbgFlags & DbgFlagEvent) != 0
}

// -----------------------------------------------------------------------------------------

type FSChanged interface {
	FileChanged(name string)
	DirAdded(name string)
	EntryDeleted(name string, isDir bool)
}

type Ignore = func(name string, isDir bool) bool

// -----------------------------------------------------------------------------------------

type Watcher struct {
	w *fsnotify.Watcher
}

func New() Watcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panicln("[FATAL] fsnotify.NewWatcher:", err)
	}
	return Watcher{w}
}

func (p Watcher) Run(root string, fc FSChanged, ignore Ignore) error {
	go p.watchLoop(root, fc, ignore)
	return watchRecursive(p.w, root)
}

func (p Watcher) watchLoop(root string, fc FSChanged, ignore Ignore) {
	const (
		eventModify = fsnotify.Write | fsnotify.Create
	)
	for {
		select {
		case event, ok := <-p.w.Events:
			if !ok { // closed
				return
			}
			if debugEvent {
				log.Println("==> event:", event)
			}
			if (event.Op & fsnotify.Remove) != 0 {
				e := p.w.Remove(event.Name)
				name, err := filepath.Rel(root, event.Name)
				if err != nil {
					log.Println("[ERROR] fsnotify.EntryDeleted filepath.Rel:", err)
					continue
				}
				isDir := (e == nil || !errors.Is(e, fsnotify.ErrNonExistentWatch))
				name = filepath.ToSlash(name)
				if ignore != nil && ignore(name, isDir) {
					continue
				} else {
					fc.EntryDeleted(name, isDir)
				}
			} else if (event.Op & eventModify) != 0 {
				name, err := filepath.Rel(root, event.Name)
				if err != nil {
					log.Println("[ERROR] fsnotify.FileChanged filepath.Rel:", err)
					continue
				}
				isDir := isDir(event.Name)
				name = filepath.ToSlash(name)
				if ignore != nil && ignore(name, isDir) {
					continue
				} else if isDir {
					if (event.Op & fsnotify.Create) != 0 {
						watchRecursive(p.w, event.Name)
						fc.DirAdded(name)
					}
				} else {
					fc.FileChanged(name)
				}
			}
		case err, ok := <-p.w.Errors:
			if !ok {
				return
			}
			log.Println("[ERROR] fsnotify Errors:", err)
		}
	}
}

func (p *Watcher) Close() error {
	if w := p.w; w != nil {
		p.w = nil
		return w.Close()
	}
	return nil
}

func watchRecursive(w *fsnotify.Watcher, path string) error {
	err := filepath.WalkDir(path, func(walkPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if err = w.Add(walkPath); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func isDir(name string) bool {
	if fs, err := os.Lstat(name); err == nil {
		return fs.IsDir()
	}
	return false
}

// -----------------------------------------------------------------------------------------
