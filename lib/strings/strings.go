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

	"Contains":       strings.Contains,
	"ContainsAny":    strings.ContainsAny,
	"ContainsRune":   strings.ContainsRune,
	"Count":          strings.Count,
	"EqualFold":      strings.EqualFold,
	"Fields":         strings.Fields,
	"FieldsFunc":     strings.FieldsFunc,
	"HasPrefix":      strings.HasPrefix,
	"HasSuffix":      strings.HasSuffix,
	"Index":          strings.Index,
	"IndexAny":       strings.IndexAny,
	"IndexByte":      strings.IndexByte,
	"IndexFunc":      strings.IndexFunc,
	"IndexRune":      strings.IndexRune,
	"Join":           strings.Join,
	"LastIndex":      strings.LastIndex,
	"LastIndexAny":   strings.LastIndexAny,
	"LastIndexFunc":  strings.LastIndexFunc,
	"Map":            strings.Map,
	"Repeat":         strings.Repeat,
	"Replace":        strings.Replace,
	"Split":          strings.Split,
	"SplitAfter":     strings.SplitAfter,
	"SplitAfterN":    strings.SplitAfterN,
	"SplitN":         strings.SplitN,
	"Title":          strings.Title,
	"ToLower":        strings.ToLower,
	"ToLowerSpecial": strings.ToLowerSpecial,
	"ToTitle":        strings.ToTitle,
	"ToTitleSpecial": strings.ToTitleSpecial,
	"ToUpper":        strings.ToUpper,
	"Trim":           strings.Trim,
	"TrimFunc":       strings.TrimFunc,
	"TrimLeft":       strings.TrimLeft,
	"TrimLeftFunc":   strings.TrimLeftFunc,
	"TrimPrefix":     strings.TrimPrefix,
	"TrimRight":      strings.TrimRight,
	"TrimRightFunc":  strings.TrimRightFunc,
	"TrimSpace":      strings.TrimSpace,
	"TrimSuffix":     strings.TrimSuffix,

	"reader":   strings.NewReader,
	"replacer": strings.NewReplacer,

	"NewReader":   strings.NewReader,
	"NewReplacer": strings.NewReplacer,

	"Reader":   tyReader,
	"Replacer": tyReplacer,
}

// -----------------------------------------------------------------------------
