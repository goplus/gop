package gopiter

import (
	"reflect"
	"testing"
)

func Test_Slice(t *testing.T) {
	c := []int{1, 2, 3, 4, 5}
	iter := NewIter(c)
	cnt := 0
	keySum := 0
	valSum := 0
	k := 0
	v := 0
	for iter.Next() {
		Key(iter, &k)
		Value(iter, &v)
		keySum += k
		valSum += v
		cnt++
	}
	if cnt != 5 || keySum != 10 || valSum != 15 {
		t.Fail()
	}
}

func Test_Map(t *testing.T) {
	c := map[int]int{0: 1, 1: 2, 2: 3, 3: 4, 4: 5}
	iter := NewIter(c)
	cnt := 0
	keySum := 0
	valSum := 0
	k := 0
	v := 0
	for iter.Next() {
		Key(iter, &k)
		Value(iter, &v)
		keySum += k
		valSum += v
		cnt++
	}
	if cnt != 5 || keySum != 10 || valSum != 15 {
		t.Fail()
	}
}

func Test_String(t *testing.T) {
	c := "hello"
	iter := NewIter(c)
	bs := []byte{'h', 'e', 'l', 'l', 'o'}
	var k int
	var v byte
	var idx int
	for iter.Next() {
		Key(iter, &k)
		Value(iter, &v)
		if k != idx {
			t.Fail()
		}
		if v != bs[idx] {
			t.Fail()
		}
		idx++
	}
}
func Test_Struct(t *testing.T) {
	c := &testIter{cap: 5}
	iter := NewIter(c)
	sum := 0
	var k int
	var v int
	for iter.Next() {
		Key(iter, &k)
		Value(iter, &v)
		sum += v
	}
	// 1*0+2*1+3*2+4*3+5*4
	if sum != 40 {
		t.Fail()
	}
}

func Test_NotSupport(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fail()
			}
		}()
		NewIter(testIter{})
	}()
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fail()
			}
		}()
		NewIter(&struct {
		}{})
	}()
}

type testIter struct {
	cap  int
	size int
}

func (t *testIter) Next() bool {
	return t.size < t.cap
}

func (t *testIter) Key() reflect.Value {
	t.size++
	return reflect.ValueOf(t.size)
}

func (t *testIter) Value() reflect.Value {
	return reflect.ValueOf(t.size * (t.size - 1))
}
