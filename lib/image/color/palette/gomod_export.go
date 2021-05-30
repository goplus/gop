// Package palette provide Go+ "image/color/palette" package, as "image/color/palette" package in Go.
package palette

import (
	palette "image/color/palette"

	gop "github.com/goplus/gop"
)

// I is a Go package instance.
var I = gop.NewGoPackage("image/color/palette")

func init() {
	I.RegisterVars(
		I.Var("Plan9", &palette.Plan9),
		I.Var("WebSafe", &palette.WebSafe),
	)
}
