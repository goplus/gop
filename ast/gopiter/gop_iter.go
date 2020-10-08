package gopiter

import (
	"log"
	"reflect"
)

type sliceIter struct {
	slice    reflect.Value
	sliceLen int
	cursor   int
}

func newSliceIter(val reflect.Value) *sliceIter {
	return &sliceIter{
		slice:    val,
		sliceLen: val.Len(),
		cursor:   0,
	}
}

func (p *sliceIter) Next() bool {
	p.cursor++
	return p.cursor-1 < p.sliceLen
}

func (p *sliceIter) Key() reflect.Value {
	return reflect.ValueOf(p.cursor - 1)
}

func (p *sliceIter) Value() reflect.Value {
	return p.slice.Index(p.cursor - 1)
}

type Iterator interface {
	Next() bool
	Key() reflect.Value
	Value() reflect.Value
}

func NewIter(container interface{}) (iter Iterator) {
	data := reflect.ValueOf(container)
	switch data.Kind() {
	case reflect.Map:
		iter = data.MapRange()
	case reflect.Slice, reflect.Array, reflect.String:
		iter = newSliceIter(reflect.Indirect(data))
	case reflect.Ptr:
		if v, ok := container.(Iterator); ok {
			iter = v
		} else {
			log.Panicln("iter must be Iterator")
		}
	default:
		log.Panicln("Iterator not support")
	}
	return
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
