package goapi

import "testing"

func TestApi(t *testing.T) {
	ApiDefaultCtx = false
	s, err := LookupApi("os")
	if err != nil {
		t.Fatalf("error %v", err)
	}
	if len(s) == 0 {
		t.FailNow()
	}
}
