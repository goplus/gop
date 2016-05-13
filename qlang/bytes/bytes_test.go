package bytes

import (
	"testing"
)

func TestCastNil(t *testing.T) {

	if from(nil) != nil {
		t.Fatal("from(nil) != nil")
	}
}
