package json

import (
	"encoding/json"
	"strings"
	"syscall"
)

// -----------------------------------------------------------------------------

// Pretty prints a value in pretty mode.
//
func Pretty(v interface{}) string {

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

// Unmarshal unmarshals a []byte or string
//
func Unmarshal(b interface{}) (v interface{}, err error) {

	switch in := b.(type) {
	case []byte:
		err = json.Unmarshal(in, &v)
	case string:
		err = json.NewDecoder(strings.NewReader(in)).Decode(&v)
	default:
		err = syscall.EINVAL
	}
	return
}

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":         "encoding/json",
	"decoder":       json.NewDecoder,
	"encoder":       json.NewEncoder,
	"marshal":       json.Marshal,
	"marshalIndent": json.MarshalIndent,
	"pretty":        Pretty,
	"unmarshal":     Unmarshal,
	"compact":       json.Compact,
	"indent":        json.Indent,
	"htmlEscape":    json.HTMLEscape,

	"NewDecoder":    json.NewDecoder,
	"NewEncoder":    json.NewEncoder,
	"Marshal":       json.Marshal,
	"MarshalIndent": json.MarshalIndent,
	"Pretty":        Pretty,
	"Unmarshal":     Unmarshal,
	"Compact":       json.Compact,
	"Indent":        json.Indent,
	"HTMLEscape":    json.HTMLEscape,
}

// -----------------------------------------------------------------------------
