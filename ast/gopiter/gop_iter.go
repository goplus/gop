package gopiter

import (
	"log"
	"reflect"
)

type rangeIterType int

const (
	arrayIterType rangeIterType = iota
	mapIterType
	iterHolderType
)

type rangeIter struct {
	typ         rangeIterType
	array       reflect.Value
	arrayLen    int
	arrayCursor int
	mapIter     *reflect.MapIter
	holderIter  Iterator
}

func newIter(d interface{}) *rangeIter {
	f := &rangeIter{}
	data := reflect.ValueOf(d)
	switch data.Kind() {
	case reflect.Map:
		f.typ = mapIterType
		f.mapIter = data.MapRange()
	case reflect.Slice, reflect.Array, reflect.String:
		f.typ = arrayIterType
		f.array = reflect.Indirect(data)
		f.arrayLen = f.array.Len()
	case reflect.Struct:
		if v, ok := d.(Iterator); ok {
			f.typ = iterHolderType
			f.holderIter = v
		} else {
			log.Panicln("iter must be Iterator")
		}
	}
	return f
}

func (p *rangeIter) Next() bool {
	switch p.typ {
	case mapIterType:
		return p.mapIter.Next()
	default:
		p.arrayCursor++
		return p.arrayCursor-1 < p.arrayLen
	}
}

func (p *rangeIter) Key() reflect.Value {
	switch p.typ {
	case mapIterType:
		return p.mapIter.Key()
	default:
		return reflect.ValueOf(p.arrayCursor - 1)
	}
}

func (p *rangeIter) Value() reflect.Value {
	switch p.typ {
	case mapIterType:
		return p.mapIter.Value()
	default:
		return p.array.Index(p.arrayCursor - 1)
	}
}

type Iterator interface {
	Next() bool
	Key() reflect.Value
	Value() reflect.Value
}

func NewIter(container interface{}) Iterator {
	return newIter(container)
}

func Key(i Iterator, key interface{}) {
	reflect.ValueOf(key).Elem().Set(i.Key())
}

func Value(i Iterator, value interface{}) {
	reflect.ValueOf(value).Elem().Set(i.Value())
}
func Next(i Iterator) bool {
	return i.Next()
}
