package cl

import (
	"testing"

	exec "github.com/qiniu/qlang/v6/exec/spec"
)

// -----------------------------------------------------------------------------

func TestValue(t *testing.T) {
	g := &goValue{t: exec.TyInt}
	_ = g.Value(0)

	nv := &nonValue{0}
	_ = nv.Kind()
	_ = nv.NumValues()
	_ = nv.Type()
	_ = nv.Value(0)

	f := new(qlFunc)
	_ = f.Kind()
	_ = f.Value(0)

	gf := new(goFunc)
	_ = gf.Kind()
	_ = gf.NumValues()
	_ = gf.Type()
	_ = gf.Value(0)

	c := new(constVal)
	_ = c.Value(0)
}

// -----------------------------------------------------------------------------
