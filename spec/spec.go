package spec

import (
	"errors"
	"reflect"
)

// -----------------------------------------------------------------------------

var (
	Undefined interface{} = errors.New("undefined")
)

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
	Set       = dummySet
	MapFrom   = dummyN
	SliceFrom = dummyN
	SubSlice  = dummy3
	EQ        = dummy2
	Get       = dummy2
	Add       = dummy2
	Sub       = dummy2
	Mul       = dummy2
	Quo       = dummy2
	Mod       = dummy2
	Inc       = dummy1
	Dec       = dummy1
)

var Fntable = map[string]interface{}{
	"$neg": dummy1,
	"$mul": dummy2,
	"$quo": dummy2,
	"$mod": dummy2,
	"$add": dummy2,
	"$sub": dummy2,

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

var SafeMode bool

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
	DumpStack = false
	AutoCall  = make(map[reflect.Type]bool)
)

func SetAutoCall(t reflect.Type) {
	AutoCall[t] = true
}

func SetDumpStack(dump bool) {
	DumpStack = dump
}

// -----------------------------------------------------------------------------
