package errors

import "errors"

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"new": errors.New,
}

// -----------------------------------------------------------------------------
