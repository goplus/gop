package strings

import (
	"strings"
)

type caller interface {
	Call(...interface{}) interface{}
}

func fieldsFunc(s string, f caller) []string {
	return strings.FieldsFunc(s, func(c rune) bool {
		v := f.Call(c)
		if b, ok := v.(bool); ok {
			return b
		}
		panic("strings.fieldsFunc: filter function must return bool")
	})
}

func indexFunc(s string, f caller) int {
	return strings.IndexFunc(s, func(c rune) bool {
		v := f.Call(c)
		if b, ok := v.(bool); ok {
			return b
		}
		panic("strings.indexFunc: filter function must return bool")
	})
}

func lastIndexFunc(s string, f caller) int {
	return strings.LastIndexFunc(s, func(c rune) bool {
		v := f.Call(c)
		if b, ok := v.(bool); ok {
			return b
		}
		panic("strings.lastIndexFunc: filter function must return bool")
	})
}

func trimFunc(s string, f caller) string {
	return strings.TrimFunc(s, func(c rune) bool {
		v := f.Call(c)
		if b, ok := v.(bool); ok {
			return b
		}
		panic("strings.trimFunc: filter function must return bool")
	})
}

func trimLeftFunc(s string, f caller) string {
	return strings.TrimLeftFunc(s, func(c rune) bool {
		v := f.Call(c)
		if b, ok := v.(bool); ok {
			return b
		}
		panic("strings.trimLeftFunc: filter function must return bool")
	})
}

func trimRightFunc(s string, f caller) string {
	return strings.TrimRightFunc(s, func(c rune) bool {
		v := f.Call(c)
		if b, ok := v.(bool); ok {
			return b
		}
		panic("strings.trimRightFunc: filter function must return bool")
	})
}

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"contains":       strings.Contains,
	"containsAny":    strings.ContainsAny,
	"containsRune":   strings.ContainsRune,
	"count":          strings.Count,
	"equalFold":      strings.EqualFold,
	"fields":         strings.Fields,
	"fieldsFunc":     fieldsFunc,
	"hasPrefix":      strings.HasPrefix,
	"hasSuffix":      strings.HasSuffix,
	"index":          strings.Index,
	"indexAny":       strings.IndexAny,
	"indexByte":      strings.IndexByte,
	"indexFunc":      indexFunc,
	"indexRune":      strings.IndexRune,
	"join":           strings.Join,
	"lastIndex":      strings.LastIndex,
	"lastIndexAny":   strings.LastIndexAny,
	"lastIndexFunc":  lastIndexFunc,
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
	"trimFunc":       trimFunc,
	"trimLeft":       strings.TrimLeft,
	"trimLeftFunc":   trimLeftFunc,
	"trimPrefix":     strings.TrimPrefix,
	"trimRight":      strings.TrimRight,
	"trimRightFunc":  trimRightFunc,
	"trimSpace":      strings.TrimSpace,
	"trimSuffix":     strings.TrimSuffix,
	"reader":         strings.NewReader,
	"replacer":       strings.NewReplacer,
}

// -----------------------------------------------------------------------------
