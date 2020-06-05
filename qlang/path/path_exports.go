package path

import (
	"path"
)

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"base":  path.Base,
	"clean": path.Clean,
	"dir":   path.Dir,
	"ext":   path.Ext,
	"isAbs": path.IsAbs,
	"join":  path.Join,
	"match": path.Match,
	"split": path.Split,

	"ErrBadPattern": path.ErrBadPattern,
}

// -----------------------------------------------------------------------------
