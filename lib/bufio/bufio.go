package bufio

import (
	"bufio"
	"io"
)

// -----------------------------------------------------------------------------

func NewReader(rd io.Reader, size ...int) *bufio.Reader {

	if len(size) == 0 {
		return bufio.NewReader(rd)
	}
	return bufio.NewReaderSize(rd, size[0])
}

func NewWriter(w io.Writer, size ...int) *bufio.Writer {

	if len(size) == 0 {
		return bufio.NewWriter(w)
	}
	return bufio.NewWriterSize(w, size[0])
}

var Exports = map[string]interface{}{
	"reader":  NewReader,
	"writer":  NewWriter,
	"scanner": bufio.NewScanner,
}

// -----------------------------------------------------------------------------

