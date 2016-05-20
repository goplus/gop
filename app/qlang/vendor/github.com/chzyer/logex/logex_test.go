package logex

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func (s *S) hello() {
	s.Warn("warn in hello")
}

func log(buf io.Writer, s string) {
	Logger{Logger: NewGoLog(buf)}.Output(2, s)
}

func test(buf io.Writer) {
	log(buf, "aa")
	Info("b")
	NewLoggerEx(buf).Info("c")
	Error("ec")
	s.hello()
	Struct(&s, 1, "", false)
}

// -------------------------------------------------------------------

type S struct {
	Logger
}

var s = S{}

func TestLogex(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := NewLoggerEx(buf)
	logger.depth = 1
	goLogStd = logger.Logger
	SetStd(logger)

	test(buf)
	ret := buf.String()

	println("--------\n", ret)

	except := []string{
		".test:logex_test.go:19]aa",
		".test:logex_test.go:20][INFO] b",
		".test:logex_test.go:21][INFO] c",
		".test:logex_test.go:22][ERROR] ec",
		".(*S).hello:logex_test.go:11][WARN] warn in hello",
	}

	for _, e := range except {
		idx := strings.Index(ret, e)
		if idx < 0 {
			t.Fatal("except", e, "not found")
		}
		ret = ret[idx+len(e):]
	}
}
