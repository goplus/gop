package meta

import (
	"errors"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"qlang.io/exec.v2"
	"qlang.io/qlang.spec.v1"
)

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"fnlist":  FnList,
	"fntable": FnTable,
	"pkglist": GoPkgList,
	"dir":     Dir,
}

// FnList returns qlang all function list
//
func FnList() (list []string) {
	for k, _ := range qlang.Fntable {
		if !strings.HasPrefix(k, "$") {
			list = append(list, k)
		}
	}
	return
}

// FnTable returns qlang all function table
//
func FnTable() map[string]interface{} {
	table := make(map[string]interface{})
	for k, v := range qlang.Fntable {
		if !strings.HasPrefix(k, "$") {
			table[k] = v
		}
	}
	return table
}

// GoPkgList returns qlang Go implemented module list
//
func GoPkgList() (list []string) {
	return qlang.GoModuleList()
}

func IsExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}

func ExporStructField(t reflect.Type) ([]string, error) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("type is not struct")
	}
	var list []string
	for i := 0; i < t.NumField(); i++ {
		if IsExported(t.Field(i).Name) {
			list = append(list, t.Field(i).Name)
		}
	}
	return list, nil
}

// Dir returns list object of strings.
// for a module: the module's attributes.
// for a go struct object: field and func list.
// for a qlang class: function list.
// for a qlang class object: function and vars list.
//
func Dir(i interface{}) (list []string) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Map {
		for _, k := range v.MapKeys() {
			list = append(list, k.String())
		}
	} else {
		switch e := i.(type) {
		case *exec.Class:
			for k, _ := range e.Fns {
				list = append(list, k)
			}
		case *exec.Object:
			for k, _ := range e.Cls.Fns {
				list = append(list, k)
			}
			for k, _ := range e.Vars() {
				list = append(list, k)
			}
		default:
			t := v.Type()
			// list struct field
			if field, err := ExporStructField(t); err == nil {
				list = append(list, field...)
			}
			// list type method
			for i := 0; i < t.NumMethod(); i++ {
				list = append(list, t.Method(i).Name)
			}
		}
	}
	return
}
