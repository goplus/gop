package interpreter

import (
	"testing"
)

// -----------------------------------------------------------------------------

func TestParseInt(t *testing.T) {

	if v, err := ParseInt("0"); err != nil || v != 0 {
		t.Fatal(`ParseInt("0")`, v, err)
	}

	if v, err := ParseInt("012"); err != nil || v != 012 {
		t.Fatal(`ParseInt("012")`, v, err)
	}

	if v, err := ParseInt("0x12"); err != nil || v != 18 {
		t.Fatal(`ParseInt("0x12")`, v, err)
	}
}

func TestParseFloat(t *testing.T) {

	if v, err := ParseFloat("0"); err != nil || v != 0 {
		t.Fatal(`ParseFloat("0")`, v, err)
	}

	if v, err := ParseFloat("012"); err != nil || v != 012 {
		t.Fatal(`ParseFloat("012")`, v, err)
	}

	if v, err := ParseFloat("012.34"); err != nil || v != 012.34 {
		t.Fatal(`ParseFloat("012.34")`, v, err)
	}

	if v, err := ParseFloat("0x12"); err != nil || v != 0x12 {
		t.Fatal(`ParseFloat("0x12")`, v, err)
	}
}

// -----------------------------------------------------------------------------

