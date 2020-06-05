package strconv

import (
	"strconv"
)

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"itoa":      strconv.Itoa,
	"parseUint": strconv.ParseUint,
	"parseInt":  strconv.ParseInt,

	"unquoteChar": strconv.UnquoteChar,
	"unquote":     strconv.Unquote,
}

// -----------------------------------------------------------------------------
