package meta

import (
	"bytes"
	"errors"
	"fmt"
	osexec "os/exec"
	"reflect"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"qlang.io/exec"
	qlang "qlang.io/spec"
)

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":   "qlang.io/lib/meta",
	"fnlist":  FnList,
	"fntable": FnTable,
	"pkgs":    GoPkgList,
	"dir":     Dir,
	"doc":     Doc,
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
		name := t.Field(i).Name
		if IsExported(name) {
			list = append(list, name)
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
	} else if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			ev := v.Index(i)
			list = append(list, fmt.Sprintf("%v", ev))
		}
	} else {
		switch e := i.(type) {
		case *exec.Class:
			for k := range e.Fns {
				list = append(list, k)
			}
		case *exec.Object:
			for k := range e.Cls.Fns {
				list = append(list, k)
			}
			for k := range e.Vars() {
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
				name := t.Method(i).Name
				if IsExported(name) {
					list = append(list, name)
				}
			}
		}
	}
	return
}

func findPackageName(i interface{}) (string, bool) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Map {
		for _, k := range v.MapKeys() {
			if k.Kind() == reflect.String && k.String() == "_name" {
				ev := v.MapIndex(k)
				if ev.Kind() == reflect.Interface {
					rv := ev.Elem()
					if rv.Kind() == reflect.String {
						return rv.String(), true
					}
				}
			}
		}
	}
	return "", false
}

// reflect.Value slice
type ReflectValues []reflect.Value

func (p ReflectValues) Len() int           { return len(p) }
func (p ReflectValues) Less(i, j int) bool { return p[i].String() < p[j].String() }
func (p ReflectValues) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Doc returns doc info of object
//
func Doc(i interface{}) string {
	if s, ok := i.(string); ok {
		return Docs(s)
	}
	var buf bytes.Buffer
	outf := func(format string, a ...interface{}) (err error) {
		_, err = buf.WriteString(fmt.Sprintf(format, a...))
		return
	}
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Map {
		mapKeys := v.MapKeys()
		sort.Sort(ReflectValues(mapKeys))
		pkgName, isPkg := findPackageName(i)
		if isPkg {
			outf("package %v", pkgName)
			for _, k := range mapKeys {
				if strings.HasPrefix(k.String(), "_") {
					continue
				}
				ev := v.MapIndex(k)
				if ev.Kind() == reflect.Interface {
					rv := ev.Elem()
					outf("\n%v\t%v ", k, rv.Type())
				}
			}
		} else {
			for _, k := range mapKeys {
				ev := v.MapIndex(k)
				outf("\n%v\t%v", k, ev)
			}
		}
	} else if v.Kind() == reflect.Slice {
		t := v.Type()
		outf("%v", t)
		for i := 0; i < v.Len(); i++ {
			ev := v.Index(i)
			if ev.Kind() == reflect.Interface {
				outf("\n%v\t%T", ev, ev.Interface())
			} else {
				outf("\n%v", ev)
			}
		}
	} else {
		switch e := i.(type) {
		case *exec.Class:
			outf("*exec.Class")
			for k, v := range e.Fns {
				outf("\n%v\t%T", k, v)
			}
		case *exec.Object:
			outf("*exec.Object")
			for k, v := range e.Cls.Fns {
				outf("\n%v\t%T", k, v)
			}
			for k, kv := range e.Vars() {
				outf("\n%v\t%T", k, kv)
			}
		default:
			t := v.Type()
			outf("%v", t)
			{
				t := v.Type()
				for t.Kind() == reflect.Ptr {
					t = t.Elem()
				}
				if t.Kind() == reflect.Struct {
					for i := 0; i < t.NumField(); i++ {
						field := t.Field(i)
						if IsExported(field.Name) {
							outf("\n%v\t%v", field.Name, field.Type)
						}
					}
				}
			}
			for i := 0; i < t.NumMethod(); i++ {
				m := t.Method(i)
				if IsExported(m.Name) {
					outf("\n%v\t%v", m.Name, m.Type)
				}
			}
		}
	}
	return buf.String()
}

func parserName(names string) (pkg string, name string) {
	if strings.Contains(names, ".") {
		ar := strings.Split(names, ".")
		pkg = ar[0]
		var sep string
		for i, v := range ar[1:] {
			if i == 1 {
				sep = "."
			}
			name += sep + strings.Title(v)
		}
	} else {
		pkg = names
	}
	return
}

func findFnTable(key string) (fn interface{}, gopkg string) {
	if i, ok := qlang.Fntable[key]; ok {
		fn = i
		if m, ok := i.(map[string]interface{}); ok {
			if name, ok := m["_name"]; ok {
				if s, ok := name.(string); ok {
					gopkg = s
				}
			}
		}
	}
	return
}

func godoc(pkg string, name string) (string, error) {
	bin, err := osexec.LookPath("godoc")
	if err != nil {
		return "", err
	}
	args := []string{pkg}
	if name != "" {
		args = append(args, name)
	}
	data, err := osexec.Command(bin, args...).Output()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Docs returns doc info of string, gopkg find in godoc
//
func Docs(names string) string {
	pkg, name := parserName(names)
	fn, gopkg := findFnTable(pkg)
	if fn == nil {
		return fmt.Sprintf("error find %q", names)
	}
	if gopkg != "" {
		info, err := godoc(gopkg, name)
		if err != nil {
			return fmt.Sprintf("godoc error %v", err)
		}
		return info
	}
	//check is string
	if s, ok := fn.(string); ok {
		return s
	}
	return Doc(fn)
}
