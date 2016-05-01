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

// Import imports a qlang module implemented by Go.
//
func Import(mod string, table map[string]interface{}) {

	if SafeMode {
		if v, ok := table["_initSafe"]; ok {
			if fn, ok := v.(func(map[string]interface{}, func(...interface{}) interface{})); ok {
				fn(table, dummyN)
			} else {
				panic("invalid prototype of _initSafe")
			}
		}
	}

	if mod != "" {
		if _, ok := Fntable[mod]; ok {
			panic("module to import exists already: " + mod)
		}
		Fntable[mod] = table
		return
	}

	for name, fn := range table {
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
