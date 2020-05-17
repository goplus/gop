package strconv

import (
	"strconv"
)

// -----------------------------------------------------------------------------

var Exports = map[string]interface{}{
	"itoa":      strconv.Itoa,
	"parseUint": strconv.ParseUint,
	"parseInt":  strconv.ParseInt,

	"unquoteChar": strconv.UnquoteChar,
	"unquote":     strconv.Unquote,
}

// -----------------------------------------------------------------------------
