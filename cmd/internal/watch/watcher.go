package watch

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/howeyc/fsnotify"
	"github.com/qiniu/x/log"
)

var fileTypesToWatch = []string{
	".gop",
}

type Watcher struct {
	dir          string
	eventHandler func(string)
	watch        *fsnotify.Watcher
}

func NewWatcher(dir string, handler func(string)) (*Watcher, error) {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		dir:          dir,
		eventHandler: handler,
		watch:        watch,
	}, nil
}

func (w Watcher) Start() error {
	go func() {
		for {
			select {
			case e := <-w.watch.Event:
				if w.shallFireEvent(e) {
					w.handleEvent(e)
				}
			case err := <-w.watch.Error:
				log.Error("failed to watch file changes:", err)
			}
		}
	}()

	return filepath.Walk(w.dir, w.processWalked)
}

func (w Watcher) handleEvent(e *fsnotify.FileEvent) {
	if e.IsCreate() {
		w.handleCreateEvent(e)
		return
	}

	if w.isFileConcerned(e.Name) {
		w.eventHandler(e.Name)
	}
}

func (w Watcher) handleCreateEvent(e *fsnotify.FileEvent) {
	info, err := os.Stat(e.Name)
	if err != nil {
		log.Errorf("failed to stat %q:", err)
		return
	}

	if info.IsDir() {
		if err := filepath.Walk(e.Name, w.processWalked); err != nil {
			log.Errorf("failed to watch %q: %v", err)
		}
	} else if w.isFileConcerned(e.Name) {
		if err := w.watch.Watch(e.Name); err != nil {
			log.Errorf("failed to watch %q: %v", err)
		}
		w.eventHandler(e.Name)
	}
}

func (w Watcher) isFileConcerned(file string) bool {
	ext := filepath.Ext(file)

	for _, tp := range fileTypesToWatch {
		if tp == ext {
			return true
		}
	}

	return false
}

func (w Watcher) processWalked(path string, info fs.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() || w.isFileConcerned(path) {
		return w.watch.Watch(path)
	}

	return nil
}

func (w Watcher) shallFireEvent(e *fsnotify.FileEvent) bool {
	// why we put IsAttrib() at the first,
	// because modify on vim fires the attrib event and modify event.
	if e.IsAttrib() {
		return false
	}
	if e.IsCreate() {
		return true
	}
	if e.IsDelete() {
		return true
	}
	if e.IsModify() {
		return true
	}

	return false
}
