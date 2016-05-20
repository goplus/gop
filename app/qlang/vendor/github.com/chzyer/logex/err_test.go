package logex

import (
	"os"
	"strings"
	"testing"
)

func b() error {
	_, err := os.Open("dflkjasldfkas")
	return Trace(err)
}

func a() error {
	return Trace(b())
}

func TestError(t *testing.T) {
	te := TraceError(a())
	errInfo := te.StackError()
	prefixes := []string{"logex%2ev1", "logex"}
	for _, p := range prefixes {
		if strings.Contains(errInfo, p+".b:11") &&
			strings.Contains(errInfo, p+".a:15") &&
			strings.Contains(errInfo, p+".TestError:19") {
			return
		}
	}

	t.Error("fail", te.StackError())
}
