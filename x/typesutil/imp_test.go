package typesutil

import (
	"go/token"
	"testing"
)

func TestNilImport(t *testing.T) {
	_, err := (&nilImporter{}).Import("fmt")
	if err == nil {
		t.Fatal("no error")
	}
	imp := newImporter(nil, nil, nil, token.NewFileSet())
	if imp.(*importer).imp == nil {
		t.Fatal("nilImporter")
	}
}
