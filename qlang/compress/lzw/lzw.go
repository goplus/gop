package lzw

import (
	"compress/lzw"
)

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"reader": lzw.NewReader,
	"writer": lzw.NewWriter,
	"LSB":    lzw.LSB,
	"MSB":    lzw.MSB,
}

// -----------------------------------------------------------------------------
