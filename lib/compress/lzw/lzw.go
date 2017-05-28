package lzw

import (
	"compress/lzw"
)

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":  "compress/lzw",
	"reader": lzw.NewReader,
	"writer": lzw.NewWriter,
	"LSB":    lzw.LSB,
	"MSB":    lzw.MSB,

	"NewReader": lzw.NewReader,
	"NewWriter": lzw.NewWriter,
}

// -----------------------------------------------------------------------------
