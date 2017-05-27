package bufio

import (
	"bufio"
	"io"
)

// -----------------------------------------------------------------------------

// NewReader returns a new Reader with size (optional) or using default size.
//
func NewReader(rd io.Reader, size ...int) *bufio.Reader {

	if len(size) == 0 {
		return bufio.NewReader(rd)
	}
	return bufio.NewReaderSize(rd, size[0])
}

// NewWriter returns a new Writer with size (optional) or using default size.
//
func NewWriter(w io.Writer, size ...int) *bufio.Writer {

	if len(size) == 0 {
		return bufio.NewWriter(w)
	}
	return bufio.NewWriterSize(w, size[0])
}

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":   "bufio",
	"reader":  NewReader,
	"writer":  NewWriter,
	"scanner": bufio.NewScanner,

	"NewReader":  NewReader,
	"NewWriter":  NewWriter,
	"NewScanner": bufio.NewScanner,
}

// -----------------------------------------------------------------------------
