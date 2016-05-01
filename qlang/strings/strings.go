package strings

import (
	"strings"
)

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"contains": strings.Contains,
	"index":    strings.Index,
	"indexAny": strings.IndexAny,
	"join":     strings.Join,
	"title":    strings.Title,
	"toLower":  strings.ToLower,
	"toTitle":  strings.ToTitle,
	"toUpper":  strings.ToUpper,
	"trim":     strings.Trim,
	"reader":   strings.NewReader,
	"replacer": strings.NewReplacer,
}

// -----------------------------------------------------------------------------
