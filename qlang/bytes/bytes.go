package bytes

import (
	"bytes"
)

// -----------------------------------------------------------------------------

func newBuffer(a ...interface{}) *bytes.Buffer {

	if len(a) == 0 {
		return new(bytes.Buffer)
	}
	switch v := a[0].(type) {
	case []byte:
		return bytes.NewBuffer(v)
	case string:
		return bytes.NewBufferString(v)
	}
	panic("bytes.buffer() - unsupported argument type")
}

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"buffer":   newBuffer,
	"reader":   bytes.NewReader,
	"contains": bytes.Contains,
	"index":    bytes.Index,
	"indexAny": bytes.IndexAny,
	"join":     bytes.Join,
	"title":    bytes.Title,
	"toLower":  bytes.ToLower,
	"toTitle":  bytes.ToTitle,
	"toUpper":  bytes.ToUpper,
	"trim":     bytes.Trim,
}

// -----------------------------------------------------------------------------
