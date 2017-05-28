package strconv

import (
	"strconv"
)

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":      "strconv",
	"itoa":       strconv.Itoa,
	"parseUint":  strconv.ParseUint,
	"parseInt":   strconv.ParseInt,
	"parseFloat": strconv.ParseFloat,

	"Itoa":       strconv.Itoa,
	"ParseUint":  strconv.ParseUint,
	"ParseInt":   strconv.ParseInt,
	"ParseFloat": strconv.ParseFloat,

	"unquoteChar": strconv.UnquoteChar,
	"unquote":     strconv.Unquote,

	"UnquoteChar": strconv.UnquoteChar,
	"Unquote":     strconv.Unquote,
}

// -----------------------------------------------------------------------------
