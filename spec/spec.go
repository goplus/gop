package spec

// -----------------------------------------------------------------------------

type undefinedType int

var undefinedBytes = []byte("\"```undefined```\"")

func (p undefinedType) Error() string {
	return "undefined"
}

func (p undefinedType) MarshalJSON() ([]byte, error) {
	return undefinedBytes, nil
}

var (
	// Undefined is `undefined` in qlang.
	Undefined interface{} = undefinedType(0)
)

// -----------------------------------------------------------------------------
