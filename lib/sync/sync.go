package sync

import (
	"sync"
)

// -----------------------------------------------------------------------------

func newMutex() *sync.Mutex {
	return new(sync.Mutex)
}

func newRWMutex() *sync.RWMutex {
	return new(sync.RWMutex)
}

func newWaitGroup() *sync.WaitGroup {
	return new(sync.WaitGroup)
}

// -----------------------------------------------------------------------------

var Exports = map[string]interface{}{
	"cond":      sync.NewCond,
	"mutex":     newMutex,
	"waitGroup": newWaitGroup,
}

// -----------------------------------------------------------------------------

