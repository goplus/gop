package extractor_test

import (
	"fmt"
	"reflect"
	"testing"

	"qiniupkg.com/text/tpl.v1/extractor"
)

// -----------------------------------------------------------------------------

const grammar1 = `

doc = (INT/items_IntArray | STRING/items_StringArray) % ','

`

func TestExtract1(t *testing.T) {

	e, err := extractor.New(grammar1, nil)
	if err != nil {
		t.Fatal("extractor.New failed:", err)
	}

	doc, err := e.SafeExtract(`1, 2, "hello, world!", 3,4`)
	if err != nil {
		t.Fatal("Extract failed:", err)
	}
	exp := map[string]interface{}{
		"items": []interface{}{
			int64(1), int64(2), "hello, world!", int64(3), int64(4),
		},
	}
	if !reflect.DeepEqual(doc, exp) {
		t.Fatal("doc:", doc, "exp:", exp)
	}
	fmt.Println("doc:", doc)
}

// -----------------------------------------------------------------------------

const grammar2 = `

doc = (INT/items_Ints | STRING/items_TextArray | CHAR/items_Chars) % ','

`

func TestExtract2(t *testing.T) {

	e, err := extractor.New(grammar2, nil)
	if err != nil {
		t.Fatal("extractor.New failed:", err)
	}

	doc, err := e.SafeExtract(`1, 2, "hello, world!", '\\', 3,4`)
	if err != nil {
		t.Fatal("Extract failed:", err)
	}
	exp := map[string]interface{}{
		"items": []interface{}{
			int64(1), int64(2), `"hello, world!"`, byte('\\'), int64(3), int64(4),
		},
	}
	if !reflect.DeepEqual(doc, exp) {
		t.Fatal("doc:", doc, "exp:", exp)
	}
	fmt.Println("doc:", doc)
}

// -----------------------------------------------------------------------------
