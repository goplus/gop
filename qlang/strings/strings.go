package strings

import (
	"strings"

	qlang "qlang.io/spec"
)

// -----------------------------------------------------------------------------

var (
	tyReader   = qlang.StructOf((*strings.Reader)(nil))
	tyReplacer = qlang.StructOf((*strings.Replacer)(nil))
)

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":          "strings",
	"contains":       strings.Contains,
	"containsAny":    strings.ContainsAny,
	"containsRune":   strings.ContainsRune,
	"count":          strings.Count,
	"equalFold":      strings.EqualFold,
	"fields":         strings.Fields,
	"fieldsFunc":     strings.FieldsFunc,
	"hasPrefix":      strings.HasPrefix,
	"hasSuffix":      strings.HasSuffix,
	"index":          strings.Index,
	"indexAny":       strings.IndexAny,
	"indexByte":      strings.IndexByte,
	"indexFunc":      strings.IndexFunc,
	"indexRune":      strings.IndexRune,
	"join":           strings.Join,
	"lastIndex":      strings.LastIndex,
	"lastIndexAny":   strings.LastIndexAny,
	"lastIndexFunc":  strings.LastIndexFunc,
	"map":            strings.Map,
	"repeat":         strings.Repeat,
	"replace":        strings.Replace,
	"split":          strings.Split,
	"splitAfter":     strings.SplitAfter,
	"splitAfterN":    strings.SplitAfterN,
	"splitN":         strings.SplitN,
	"title":          strings.Title,
	"toLower":        strings.ToLower,
	"toLowerSpecial": strings.ToLowerSpecial,
	"toTitle":        strings.ToTitle,
	"toTitleSpecial": strings.ToTitleSpecial,
	"toUpper":        strings.ToUpper,
	"trim":           strings.Trim,
	"trimFunc":       strings.TrimFunc,
	"trimLeft":       strings.TrimLeft,
	"trimLeftFunc":   strings.TrimLeftFunc,
	"trimPrefix":     strings.TrimPrefix,
	"trimRight":      strings.TrimRight,
	"trimRightFunc":  strings.TrimRightFunc,
	"trimSpace":      strings.TrimSpace,
	"trimSuffix":     strings.TrimSuffix,

	"reader":   strings.NewReader,
	"replacer": strings.NewReplacer,

	"Reader":   tyReader,
	"Replacer": tyReplacer,
}

// -----------------------------------------------------------------------------
