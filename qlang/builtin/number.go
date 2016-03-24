package builtin

import (
	"fmt"
	"reflect"
	"strings"
)

// -----------------------------------------------------------------------------

func Inc(a interface{}) interface{} {

	switch v := a.(type) {
	case int:
		return v+1
	case byte:
		return v+1
	}
	return panicUnsupportedOp1("++", a)
}

func Dec(a interface{}) interface{} {

	switch v := a.(type) {
	case int:
		return v-1
	case byte:
		return v-1
	}
	return panicUnsupportedOp1("--", a)
}

func Neg(a interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		return -a1
	case float64:
		return -a1
	}
	return panicUnsupportedOp1("-", a)
}

func Float(a interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		return float64(a1)
	case float64:
		return a1
	case byte:
		return float64(a1)
	}
	return panicUnsupportedFn("float", a)
}

func Int(a interface{}) interface{} {

	switch a1 := a.(type) {
	case float64:
		return int(a1)
	case int:
		return a1
	case byte:
		return int(a1)
	}
	return panicUnsupportedFn("int", a)
}

func Byte(a interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		return byte(a1)
	case byte:
		return a1
	case float64:
		return byte(a1)
	}
	return panicUnsupportedFn("byte", a)
}

func String(a interface{}) interface{} {

	switch a1 := a.(type) {
	case []byte:
		return string(a1)
	case byte:
		return string(a1)
	case int:
		return string(a1)
	}
	return panicUnsupportedFn("string", a)
}

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

func Mod(a, b interface{}) interface{} {

	if a1, ok := a.(int); ok {
		if b1, ok := b.(int); ok {
			return a1 % b1
		}
	}
	return panicUnsupportedOp2("%", a, b)
}

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
	}
	return panicUnsupportedOp2("+", a, b)
}

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
	}
	return panicUnsupportedOp2("-", a, b)
}

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
	return nil
}

func panicUnsupportedOp1(op string, a interface{}) interface{} {

	ta := typeString(a)
	panic("unsupported operator: " + op + ta)
	return nil
}

func panicUnsupportedOp2(op string, a, b interface{}) interface{} {

	ta := typeString(a)
	tb := typeString(b)
	panic("unsupported operator: " + ta + op + tb)
	return nil
}

func typeString(a interface{}) string {

	if a == nil {
		return "nil"
	}
	return reflect.TypeOf(a).String()
}

// -----------------------------------------------------------------------------
