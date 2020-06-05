package builtin

import (
	"fmt"
	"reflect"
	"strings"
)

// -----------------------------------------------------------------------------

// Inc returns a+1
//
func Inc(a interface{}) interface{} {

	switch v := a.(type) {
	case int:
		return v + 1
	case uint:
		return v + 1
	case int64:
		return v + 1
	case uint64:
		return v + 1
	case int32:
		return v + 1
	case uint32:
		return v + 1
	case uint8:
		return v + 1
	case int8:
		return v + 1
	case uint16:
		return v + 1
	case int16:
		return v + 1
	}
	return panicUnsupportedOp1("++", a)
}

// Dec returns a-1
//
func Dec(a interface{}) interface{} {

	switch v := a.(type) {
	case int:
		return v - 1
	case uint:
		return v - 1
	case int64:
		return v - 1
	case uint64:
		return v - 1
	case int32:
		return v - 1
	case uint32:
		return v - 1
	case uint8:
		return v - 1
	case int8:
		return v - 1
	case uint16:
		return v - 1
	case int16:
		return v - 1
	}
	return panicUnsupportedOp1("--", a)
}

// Neg returns -a
//
func Neg(a interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		return -a1
	case float64:
		return -a1
	}
	return panicUnsupportedOp1("-", a)
}

// Float returns float64(a)
//
func Float(a interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		return float64(a1)
	case float64:
		return a1
	}
	return panicUnsupportedFn("float", a)
}

// Int returns int(a)
//
func Int(a interface{}) interface{} {

	switch a1 := a.(type) {
	case float64:
		return int(a1)
	case int:
		return a1
	}
	return panicUnsupportedFn("int", a)
}

// Int8 returns int8(a)
//
func Int8(a interface{}) interface{} {

	switch a1 := a.(type) {
	case float64:
		return int8(a1)
	case int:
		return int8(a1)
	}
	return panicUnsupportedFn("int8", a)
}

// Int16 returns int16(a)
//
func Int16(a interface{}) interface{} {

	switch a1 := a.(type) {
	case float64:
		return int16(a1)
	case int:
		return int16(a1)
	}
	return panicUnsupportedFn("int16", a)
}

// Int32 returns int32(a)
//
func Int32(a interface{}) interface{} {

	switch a1 := a.(type) {
	case float64:
		return int32(a1)
	case int:
		return int32(a1)
	}
	return panicUnsupportedFn("int32", a)
}

// Int64 returns int64(a)
//
func Int64(a interface{}) interface{} {

	switch a1 := a.(type) {
	case float64:
		return int64(a1)
	case int:
		return int64(a1)
	}
	return panicUnsupportedFn("int64", a)
}

// Uint16 returns uint16(a)
//
func Uint16(a interface{}) interface{} {

	switch a1 := a.(type) {
	case float64:
		return uint16(a1)
	case int:
		return uint16(a1)
	}
	return panicUnsupportedFn("uint16", a)
}

// Uint32 returns uint32(a)
//
func Uint32(a interface{}) interface{} {

	switch a1 := a.(type) {
	case float64:
		return uint32(a1)
	case int:
		return uint32(a1)
	}
	return panicUnsupportedFn("uint32", a)
}

// Uint64 returns uint64(a)
//
func Uint64(a interface{}) interface{} {

	switch a1 := a.(type) {
	case float64:
		return uint64(a1)
	case int:
		return uint64(a1)
	}
	return panicUnsupportedFn("uint64", a)
}

// Uint returns uint(a)
//
func Uint(a interface{}) interface{} {

	switch a1 := a.(type) {
	case float64:
		return uint(a1)
	case int:
		return uint(a1)
	}
	return panicUnsupportedFn("uint", a)
}

// Byte returns byte(a)
//
func Byte(a interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		return byte(a1)
	case float64:
		return byte(a1)
	}
	return panicUnsupportedFn("byte", a)
}

// String returns string(a)
//
func String(a interface{}) interface{} {

	switch a1 := a.(type) {
	case []byte:
		return string(a1)
	case int:
		return string(a1)
	}
	return panicUnsupportedFn("string", a)
}

// Mul returns a*b
//
func Mul(a, b interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		switch b1 := b.(type) {
		case int:
			return a1 * b1
		case float64:
			return float64(a1) * b1
		}
	case float64:
		switch b1 := b.(type) {
		case int:
			return a1 * float64(b1)
		case float64:
			return a1 * b1
		}
	}
	return panicUnsupportedOp2("*", a, b)
}

// Quo returns a/b
//
func Quo(a, b interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		switch b1 := b.(type) {
		case int:
			return a1 / b1
		case float64:
			return float64(a1) / b1
		}
	case float64:
		switch b1 := b.(type) {
		case int:
			return a1 / float64(b1)
		case float64:
			return a1 / b1
		}
	}
	return panicUnsupportedOp2("/", a, b)
}

// Mod returns a%b
//
func Mod(a, b interface{}) interface{} {

	if a1, ok := a.(int); ok {
		if b1, ok := b.(int); ok {
			return a1 % b1
		}
	}
	return panicUnsupportedOp2("%", a, b)
}

// Add returns a+b
//
func Add(a, b interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		switch b1 := b.(type) {
		case int:
			return a1 + b1
		case float64:
			return float64(a1) + b1
		}
	case float64:
		switch b1 := b.(type) {
		case int:
			return a1 + float64(b1)
		case float64:
			return a1 + b1
		}
	case string:
		if b1, ok := b.(string); ok {
			return a1 + b1
		}
	case uint:
		switch b1 := b.(type) {
		case int:
			return a1 + uint(b1)
		}
	case uint64:
		switch b1 := b.(type) {
		case int:
			return a1 + uint64(b1)
		}
	case int64:
		switch b1 := b.(type) {
		case int:
			return a1 + int64(b1)
		}
	case uint32:
		switch b1 := b.(type) {
		case int:
			return a1 + uint32(b1)
		}
	case int32:
		switch b1 := b.(type) {
		case int:
			return a1 + int32(b1)
		}
	case uint16:
		switch b1 := b.(type) {
		case int:
			return a1 + uint16(b1)
		}
	case int16:
		switch b1 := b.(type) {
		case int:
			return a1 + int16(b1)
		}
	case uint8:
		switch b1 := b.(type) {
		case int:
			return a1 + uint8(b1)
		}
	case int8:
		switch b1 := b.(type) {
		case int:
			return a1 + int8(b1)
		}
	}
	return panicUnsupportedOp2("+", a, b)
}

// Sub returns a-b
//
func Sub(a, b interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		switch b1 := b.(type) {
		case int:
			return a1 - b1
		case float64:
			return float64(a1) - b1
		}
	case float64:
		switch b1 := b.(type) {
		case int:
			return a1 - float64(b1)
		case float64:
			return a1 - b1
		}
	case uint:
		switch b1 := b.(type) {
		case int:
			return a1 - uint(b1)
		}
	case uint64:
		switch b1 := b.(type) {
		case int:
			return a1 - uint64(b1)
		}
	case int64:
		switch b1 := b.(type) {
		case int:
			return a1 - int64(b1)
		}
	case uint32:
		switch b1 := b.(type) {
		case int:
			return a1 - uint32(b1)
		}
	case int32:
		switch b1 := b.(type) {
		case int:
			return a1 - int32(b1)
		}
	case uint16:
		switch b1 := b.(type) {
		case int:
			return a1 - uint16(b1)
		}
	case int16:
		switch b1 := b.(type) {
		case int:
			return a1 - int16(b1)
		}
	case uint8:
		switch b1 := b.(type) {
		case int:
			return a1 - uint8(b1)
		}
	case int8:
		switch b1 := b.(type) {
		case int:
			return a1 - int8(b1)
		}
	}
	return panicUnsupportedOp2("-", a, b)
}

// Max returns max(a1, a2, ...)
//
func Max(args ...interface{}) (max interface{}) {

	if len(args) == 0 {
		return 0
	}

	switch kindOfArgs(args) {
	case reflect.Int:
		return maxInt(args)
	case reflect.Float64:
		return maxFloat(args)
	}
	return panicUnsupportedFn("max", args)
}

// Min returns min(a1, a2, ...)
//
func Min(args ...interface{}) (min interface{}) {

	if len(args) == 0 {
		return 0
	}

	switch kindOfArgs(args) {
	case reflect.Int:
		return minInt(args)
	case reflect.Float64:
		return minFloat(args)
	}
	return panicUnsupportedFn("min", args)
}

func kindOfArgs(args []interface{}) reflect.Kind {

	kind := kindOf(args[0])
	for i := 1; i < len(args); i++ {
		if t := kindOf(args[i]); t != kind {
			if kind == reflect.Float64 || kind == reflect.Int {
				if t == reflect.Int {
					continue
				}
				if t == reflect.Float64 {
					kind = reflect.Float64
					continue
				}
			}
			return reflect.Invalid
		}
	}
	return kind
}

func kindOf2Args(args []interface{}, idx int) reflect.Kind {

	kind := kindOf(args[idx])
	for i := 2; i < len(args); i += 2 {
		if t := kindOf(args[i+idx]); t != kind {
			if kind == reflect.Float64 || kind == reflect.Int {
				if t == reflect.Int {
					continue
				}
				if t == reflect.Float64 {
					kind = reflect.Float64
					continue
				}
			}
			return reflect.Invalid
		}
	}
	return kind
}

func maxFloat(args []interface{}) (max float64) {

	max = asFloat(args[0])
	for i := 1; i < len(args); i++ {
		if t := asFloat(args[i]); t > max {
			max = t
		}
	}
	return
}

func minFloat(args []interface{}) (min float64) {

	min = asFloat(args[0])
	for i := 1; i < len(args); i++ {
		if t := asFloat(args[i]); t < min {
			min = t
		}
	}
	return
}

func maxInt(args []interface{}) (max int) {

	max = args[0].(int)
	for i := 1; i < len(args); i++ {
		if t := args[i].(int); t > max {
			max = t
		}
	}
	return
}

func minInt(args []interface{}) (min int) {

	min = args[0].(int)
	for i := 1; i < len(args); i++ {
		if t := args[i].(int); t < min {
			min = t
		}
	}
	return
}

func asFloat(a interface{}) float64 {

	switch v := a.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	}
	panic("unreachable")
}

func asInt(a interface{}) int {

	switch v := a.(type) {
	case int:
		return v
	}
	panic(fmt.Sprintf("param `%v` not a integer", a))
}

func kindOf(a interface{}) reflect.Kind {

	return reflect.ValueOf(a).Kind()
}

func panicUnsupportedFn(fn string, args ...interface{}) interface{} {

	targs := make([]string, len(args))
	for i, a := range args {
		targs[i] = typeString(a)
	}
	panic("unsupported function: " + fn + "(" + strings.Join(targs, ",") + ")")
}

func panicUnsupportedOp1(op string, a interface{}) interface{} {

	ta := typeString(a)
	panic("unsupported operator: " + op + ta)
}

func panicUnsupportedOp2(op string, a, b interface{}) interface{} {

	ta := typeString(a)
	tb := typeString(b)
	panic("unsupported operator: " + ta + op + tb)
}

func typeString(a interface{}) string {

	if a == nil {
		return "nil"
	}
	return reflect.TypeOf(a).String()
}

// -----------------------------------------------------------------------------
