package qlang

import (
	"errors"
	"reflect"
)

// -----------------------------------------------------------------------------

var (
	// Undefined is `undefined` in qlang.
	Undefined = interface{}(errors.New("undefined"))
)

// A Chan represents chan class in qlang.
//
type Chan struct {
	Data reflect.Value
}

// A DataIndex represents a compound data and its index. such as map[index], slice[index], object.member.
//
type DataIndex struct {
	Data  interface{}
	Index interface{}
}

// A Var represents a qlang variable.
//
type Var struct {
	Name string
}

// -----------------------------------------------------------------------------

type Module struct {
	Exports map[string]interface{}
}

func (p Module) Disable(fnNames ...string) {

	for _, fnName := range fnNames {
		p.Exports[fnName] = dummyN
	}
}

func GoModuleName(table map[string]interface{}) (name string, ok bool) {

	name, ok = table["_name"].(string)
	return
}

func GetInitSafe(table map[string]interface{}) (initSafe func(mod Module), ok bool) {

	vinitSafe, ok := table["_initSafe"]
	if !ok {
		return
	}
	initSafe, ok = vinitSafe.(func(mod Module))
	if !ok {
		panic("invalid prototype of initSafe: must be `func initSafe(mod qlang.Module)`")
	}
	return
}

// -----------------------------------------------------------------------------

func dummy1(a interface{}) interface{} {
	panic("function not implemented")
}

func dummy2(a, b interface{}) interface{} {
	panic("function not implemented")
}

func dummy3(a, i, j interface{}) interface{} {
	panic("function not implemented")
}

func dummyN(a ...interface{}) interface{} {
	panic("function not implemented")
}

func dummySetIndex(m, index, val interface{}) {
	panic("function not implemented")
}

func dummySet(m interface{}, args ...interface{}) {
	panic("function not implemented")
}

func dummyChanIn(ch, val interface{}, try bool) interface{} {
	panic("operator ch<-value not implemented")
}

func dummyChanOut(ch interface{}, try bool) interface{} {
	panic("operator <-ch not implemented")
}

var (
	fnDummy1  = reflect.ValueOf(dummy1).Pointer()
	fnDummy2  = reflect.ValueOf(dummy2).Pointer()
	ChanIn    = dummyChanIn
	ChanOut   = dummyChanOut
	GetEx     = dummy2
	SetEx     = dummySet
	SetIndex  = dummySetIndex
	MapFrom   = dummyN
	SliceFrom = dummyN
	SubSlice  = dummy3
	Xor       = dummy2
	Lshr      = dummy2
	Rshr      = dummy2
	BitAnd    = dummy2
	BitOr     = dummy2
	AndNot    = dummy2
	EQ        = dummy2
	GetVar    = dummy2
	Get       = dummy2
	Add       = dummy2
	Sub       = dummy2
	Mul       = dummy2
	Quo       = dummy2
	Mod       = dummy2
	Inc       = dummy1
	Dec       = dummy1
)

// Fntable is function table required by tpl.Interpreter engine.
//
var Fntable = map[string]interface{}{
	"$neg": dummy1,
	"$mul": dummy2,
	"$quo": dummy2,
	"$mod": dummy2,
	"$add": dummy2,
	"$sub": dummy2,

	"$xor":    dummy2,
	"$lshr":   dummy2,
	"$rshr":   dummy2,
	"$bitand": dummy2,
	"$bitor":  dummy2,
	"$bitnot": dummy1,
	"$andnot": dummy2,

	"$lt":  dummy2,
	"$gt":  dummy2,
	"$le":  dummy2,
	"$ge":  dummy2,
	"$eq":  dummy2,
	"$ne":  dummy2,
	"$and": dummy2,
	"$or":  dummy2,
	"$not": dummy1,
}

// SafeMode is the init mode of qlang.
//
var SafeMode bool

var goMods []string

// GoModuleList returns qlang modules implemented by Go.
//
func GoModuleList() []string {
	return goMods
}

// Import imports a qlang module implemented by Go.
//
func Import(mod string, table map[string]interface{}) {

	if SafeMode {
		if initSafe, ok := GetInitSafe(table); ok {
			initSafe(Module{Exports: table})
		}
	}

	if mod != "" {
		if _, ok := Fntable[mod]; ok {
			panic("module to import exists already: " + mod)
		}
		Fntable[mod] = table
		goMods = append(goMods, mod)
		return
	}

	for name, fn := range table {
		if name == "_name" {
			continue
		}
		if fn, ok := Fntable[name]; ok {
			if v := reflect.ValueOf(fn); v.Kind() == reflect.Func {
				p := v.Pointer()
				if p != fnDummy1 && p != fnDummy2 {
					panic("symbol to import exists already: " + name)
				}
			} else {
				panic("symbol to import exists already: " + name)
			}
		}
		Fntable[name] = fn
	}
}

// -----------------------------------------------------------------------------

var (
	// DumpStack indicates to dump stack when error.
	DumpStack = false

	// AutoCall is reserved for internal use.
	AutoCall = make(map[reflect.Type]bool)
)

// SetAutoCall is reserved for internal use.
//
func SetAutoCall(t reflect.Type) {
	AutoCall[t] = true
}

// SetDumpStack set to dump stack or not.
//
func SetDumpStack(dump bool) {
	DumpStack = dump
}

// -----------------------------------------------------------------------------
