package qlang_test

import (
	"reflect"
	"testing"

	"qlang.io/cl/qlang"
	"qlang.io/lib/math"
	qspec "qlang.io/spec"
)

// -----------------------------------------------------------------------------

const testCastFloatCode = `

x = sin(0)
`

func TestCastFloat(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testCastFloatCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 0.0 {
		t.Fatal("x != 0.0, x =", v)
	}
}

func init() {
	qlang.Import("", math.Exports)
}

// -----------------------------------------------------------------------------

const testTestByteSliceCode = `

a = ['a', 'b']
x = type(a)
`

func TestByteSlice(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testTestByteSliceCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != reflect.TypeOf([]byte(nil)) {
		t.Fatal("x != []byte, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testNilSliceCode = `

a = []
a = append(a, 1)
x = type(a)
y = a[0]
`

func TestNilSlice(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testNilSliceCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != reflect.TypeOf([]interface{}(nil)) {
		t.Fatal("x != []interface{}, x =", v)
	}
	if v, ok := lang.GetVar("y"); !ok || v != 1 {
		t.Fatal("y != 1, y =", v)
	}
}

// -----------------------------------------------------------------------------

type fooMode uint

func castFooMode(mode fooMode) int {
	return int(mode)
}

const testCastFooModeCode = `

y = castFooMode(1)
`

func TestCastFooMode(t *testing.T) {

	lang := qlang.New()
	qlang.Import("", map[string]interface{}{
		"castFooMode": castFooMode,
	})

	err := lang.SafeExec([]byte(testCastFooModeCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("y"); !ok || v != 1 {
		t.Fatal("y != 1, y =", v)
	}
}

// -----------------------------------------------------------------------------

type AImportType struct {
	X int
}

func (p *AImportType) GetX() int {
	return p.X
}

func init() {
	qlang.Import("fooMod", map[string]interface{}{
		"AImportType": qspec.StructOf((*AImportType)(nil)),
	})
}

const testImportTypeCode = `

a = new fooMod.AImportType
a.x = 1
y = a.getX()
`

func TestImportType(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testImportTypeCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("y"); !ok || v != 1 {
		t.Fatal("y != 1, y =", v)
	}
}

// -----------------------------------------------------------------------------

const testStructInitCode = `

a = &fooMod.AImportType{x: 1}
y = a.x
`

func TestStructInit(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testStructInitCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("y"); !ok || v != 1 {
		t.Fatal("y != 1, y =", v)
	}
}

// -----------------------------------------------------------------------------

const testMapInitCode = `

a = map[string]int16{"x": 1}
y = a.x
`

func TestMapInit(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testMapInitCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("y"); !ok || v != int16(1) {
		t.Fatal("y != 1, y =", v)
	}
}

// -----------------------------------------------------------------------------

const testMapInit2Code = `

y = map[string]int(nil)
`

func TestMapInit2(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testMapInit2Code), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("y"); !ok || v.(map[string]int) != nil {
		t.Fatal("y != nil, y =", v)
	}
}

// -----------------------------------------------------------------------------

const testMapInit3Code = `

a = map[string]var{"x": 1, "y": "hello"}
y = a.x
z = a.y
`

func TestMapInit3(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testMapInit3Code), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("y"); !ok || v != 1 {
		t.Fatal("y != 1, y =", v)
	}
	if v, ok := lang.GetVar("z"); !ok || v != "hello" {
		t.Fatal("z != hello, z =", v)
	}
}

// -----------------------------------------------------------------------------

const testMapInit32Code = `

a = map[string]int16{"x": 1, "hello": 2}
z = a.x
`

func TestMapInit32(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testMapInit32Code), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("z"); !ok || v != int16(1) {
		t.Fatal("z != 1, z =", v)
	}
}

// -----------------------------------------------------------------------------

const testMapIDCode = `

map = 1
y = map

t = class {
	fn map() {
		return "hello"
	}
}

obj = new t
z = obj.map()
`

func TestMapID(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(testMapIDCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("y"); !ok || v != 1 {
		t.Fatal("y != 1, y =", v)
	}
	if v, ok := lang.GetVar("z"); !ok || v != "hello" {
		t.Fatal("z != hello, z =", v)
	}
}

// -----------------------------------------------------------------------------
