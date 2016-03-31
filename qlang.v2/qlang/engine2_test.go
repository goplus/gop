package qlang

import (
	"testing"

	_ "qlang.io/qlang/builtin"
)

// -----------------------------------------------------------------------------

const scriptA = `

println("in script A")

foo = fn() {
	println("in func foo")
}
`

const scriptB = `

include "a.ql"

println("in script B")
foo()
`

func TestInclude(t *testing.T) {

	lang, _ := New(InsertSemis)

	lang.Include(func(file string) int {
		end, err := lang.Cl([]byte(scriptA), "a.ql", true)
		if err != nil {
			panic(err)
		}
		return end
	})

	err := lang.SafeExec([]byte(scriptB), "b.ql")
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------

