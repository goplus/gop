package typesutil

import (
	"errors"
	"testing"

	"github.com/goplus/gogen"
)

func TestConvErr(t *testing.T) {
	e := errors.New("foo")
	if ret, ok := convErr(nil, &gogen.ImportError{Err: e}); !ok || ret.Msg != "foo" {
		t.Fatal("convErr:", ret, ok)
	}
	if _, ok := convErr(nil, e); ok {
		t.Fatal("convErr: ok?")
	}
}
