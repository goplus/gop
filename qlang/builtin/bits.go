package builtin

// -----------------------------------------------------------------------------

// Lshr returns a << b
//
func Lshr(a, b interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		switch b1 := b.(type) {
		case int:
			return a1 << uint(b1)
		}
	}
	return panicUnsupportedOp2("<<", a, b)
}

// Rshr returns a >> b
//
func Rshr(a, b interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		switch b1 := b.(type) {
		case int:
			return a1 >> uint(b1)
		}
	}
	return panicUnsupportedOp2(">>", a, b)
}

// Xor returns a ^ b
//
func Xor(a, b interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		switch b1 := b.(type) {
		case int:
			return a1 ^ b1
		}
	}
	return panicUnsupportedOp2("^", a, b)
}

// BitAnd returns a & b
//
func BitAnd(a, b interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		switch b1 := b.(type) {
		case int:
			return a1 & b1
		}
	}
	return panicUnsupportedOp2("&", a, b)
}

// BitOr returns a | b
//
func BitOr(a, b interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		switch b1 := b.(type) {
		case int:
			return a1 | b1
		}
	}
	return panicUnsupportedOp2("|", a, b)
}

// BitNot returns ^a
//
func BitNot(a interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		return ^a1
	}
	return panicUnsupportedOp1("^", a)
}

// AndNot returns a &^ b
//
func AndNot(a, b interface{}) interface{} {

	switch a1 := a.(type) {
	case int:
		switch b1 := b.(type) {
		case int:
			return a1 &^ b1
		}
	}
	return panicUnsupportedOp2("&^", a, b)
}

// -----------------------------------------------------------------------------
