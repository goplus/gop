package json

import (
	"encoding/json"
	"strings"
	"syscall"
)

// -----------------------------------------------------------------------------

func Pretty(v interface{}) string {

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

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

var Exports = map[string]interface{}{
	"decoder":       json.NewDecoder,
	"encoder":       json.NewEncoder,
	"marshal":       json.Marshal,
	"marshalIndent": json.MarshalIndent,
	"pretty":        Pretty,
	"unmarshal":     Unmarshal,
	"compact":       json.Compact,
	"indent":        json.Indent,
	"htmlEscape":    json.HTMLEscape,
}

// -----------------------------------------------------------------------------
