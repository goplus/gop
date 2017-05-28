package path

import (
	"path"
)

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name": "path",
	"base":  path.Base,
	"clean": path.Clean,
	"dir":   path.Dir,
	"ext":   path.Ext,
	"isAbs": path.IsAbs,
	"join":  path.Join,
	"match": path.Match,
	"split": path.Split,

	"Base":  path.Base,
	"Clean": path.Clean,
	"Dir":   path.Dir,
	"Ext":   path.Ext,
	"IsAbs": path.IsAbs,
	"Join":  path.Join,
	"Match": path.Match,
	"Split": path.Split,

	"ErrBadPattern": path.ErrBadPattern,
}

// -----------------------------------------------------------------------------
