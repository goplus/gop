package main

import (
	"strings"
)

// -----------------------------------------------------------------------------

type goPkgExporter struct {
}

// -----------------------------------------------------------------------------

type goFunc struct {
	Name string
	This interface{}
}

var exports = []goFunc{
	{"NewReplacer", strings.NewReplacer},
	{"(*Replacer).Replace", (*strings.Replacer).Replace},
}

func exportGof(gof goFunc) {

}

// -----------------------------------------------------------------------------

func main() {
	for _, gof := range exports {
		exportGof(gof)
	}
}

// -----------------------------------------------------------------------------
