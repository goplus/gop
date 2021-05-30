// Package draw provide Go+ "image/draw" package, as "image/draw" package in Go.
package draw

import (
	image "image"
	color "image/color"
	draw "image/draw"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func toType0(v interface{}) draw.Image {
	if v == nil {
		return nil
	}
	return v.(draw.Image)
}

func toType1(v interface{}) image.Image {
	if v == nil {
		return nil
	}
	return v.(image.Image)
}

func execDraw(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	draw.Draw(toType0(args[0]), args[1].(image.Rectangle), toType1(args[2]), args[3].(image.Point), args[4].(draw.Op))
	p.PopN(5)
}

func execDrawMask(_ int, p *gop.Context) {
	args := p.GetArgs(7)
	draw.DrawMask(toType0(args[0]), args[1].(image.Rectangle), toType1(args[2]), args[3].(image.Point), toType1(args[4]), args[5].(image.Point), args[6].(draw.Op))
	p.PopN(7)
}

func execiDrawerDraw(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	args[0].(draw.Drawer).Draw(toType0(args[1]), args[2].(image.Rectangle), toType1(args[3]), args[4].(image.Point))
	p.PopN(5)
}

func execiImageAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(draw.Image).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execiImageBounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(draw.Image).Bounds()
	p.Ret(1, ret0)
}

func execiImageColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(draw.Image).ColorModel()
	p.Ret(1, ret0)
}

func toType2(v interface{}) color.Color {
	if v == nil {
		return nil
	}
	return v.(color.Color)
}

func execiImageSet(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(draw.Image).Set(args[1].(int), args[2].(int), toType2(args[3]))
	p.PopN(4)
}

func execmOpDraw(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	args[0].(draw.Op).Draw(toType0(args[1]), args[2].(image.Rectangle), toType1(args[3]), args[4].(image.Point))
	p.PopN(5)
}

func execiQuantizerQuantize(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(draw.Quantizer).Quantize(args[1].(color.Palette), toType1(args[2]))
	p.Ret(3, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("image/draw")

func init() {
	I.RegisterFuncs(
		I.Func("Draw", draw.Draw, execDraw),
		I.Func("DrawMask", draw.DrawMask, execDrawMask),
		I.Func("(Drawer).Draw", (draw.Drawer).Draw, execiDrawerDraw),
		I.Func("(Image).At", (draw.Image).At, execiImageAt),
		I.Func("(Image).Bounds", (draw.Image).Bounds, execiImageBounds),
		I.Func("(Image).ColorModel", (draw.Image).ColorModel, execiImageColorModel),
		I.Func("(Image).Set", (draw.Image).Set, execiImageSet),
		I.Func("(Op).Draw", (draw.Op).Draw, execmOpDraw),
		I.Func("(Quantizer).Quantize", (draw.Quantizer).Quantize, execiQuantizerQuantize),
	)
	I.RegisterVars(
		I.Var("FloydSteinberg", &draw.FloydSteinberg),
	)
	I.RegisterTypes(
		I.Type("Drawer", reflect.TypeOf((*draw.Drawer)(nil)).Elem()),
		I.Type("Image", reflect.TypeOf((*draw.Image)(nil)).Elem()),
		I.Type("Op", reflect.TypeOf((*draw.Op)(nil)).Elem()),
		I.Type("Quantizer", reflect.TypeOf((*draw.Quantizer)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("Over", qspec.Int, draw.Over),
		I.Const("Src", qspec.Int, draw.Src),
	)
}
