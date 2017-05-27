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

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":     "sync",
	"cond":      sync.NewCond,
	"mutex":     newMutex,
	"waitGroup": newWaitGroup,

	"NewCond":      sync.NewCond,
	"NewMutex":     newMutex,
	"NewWaitGroup": newWaitGroup,
}

// -----------------------------------------------------------------------------
