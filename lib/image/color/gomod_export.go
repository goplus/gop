// Package color provide Go+ "image/color" package, as "image/color" package in Go.
package color

import (
	color "image/color"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execmAlphaRGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2, ret3 := args[0].(color.Alpha).RGBA()
	p.Ret(1, ret0, ret1, ret2, ret3)
}

func execmAlpha16RGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2, ret3 := args[0].(color.Alpha16).RGBA()
	p.Ret(1, ret0, ret1, ret2, ret3)
}

func execmCMYKRGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2, ret3 := args[0].(color.CMYK).RGBA()
	p.Ret(1, ret0, ret1, ret2, ret3)
}

func execCMYKToRGB(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1, ret2 := color.CMYKToRGB(args[0].(uint8), args[1].(uint8), args[2].(uint8), args[3].(uint8))
	p.Ret(4, ret0, ret1, ret2)
}

func execiColorRGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2, ret3 := args[0].(color.Color).RGBA()
	p.Ret(1, ret0, ret1, ret2, ret3)
}

func execmGrayRGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2, ret3 := args[0].(color.Gray).RGBA()
	p.Ret(1, ret0, ret1, ret2, ret3)
}

func execmGray16RGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2, ret3 := args[0].(color.Gray16).RGBA()
	p.Ret(1, ret0, ret1, ret2, ret3)
}

func toType0(v interface{}) color.Color {
	if v == nil {
		return nil
	}
	return v.(color.Color)
}

func execiModelConvert(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(color.Model).Convert(toType0(args[1]))
	p.Ret(2, ret0)
}

func execModelFunc(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := color.ModelFunc(args[0].(func(color.Color) color.Color))
	p.Ret(1, ret0)
}

func execmNRGBARGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2, ret3 := args[0].(color.NRGBA).RGBA()
	p.Ret(1, ret0, ret1, ret2, ret3)
}

func execmNRGBA64RGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2, ret3 := args[0].(color.NRGBA64).RGBA()
	p.Ret(1, ret0, ret1, ret2, ret3)
}

func execmNYCbCrARGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2, ret3 := args[0].(color.NYCbCrA).RGBA()
	p.Ret(1, ret0, ret1, ret2, ret3)
}

func execmPaletteConvert(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(color.Palette).Convert(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmPaletteIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(color.Palette).Index(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmRGBARGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2, ret3 := args[0].(color.RGBA).RGBA()
	p.Ret(1, ret0, ret1, ret2, ret3)
}

func execmRGBA64RGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2, ret3 := args[0].(color.RGBA64).RGBA()
	p.Ret(1, ret0, ret1, ret2, ret3)
}

func execRGBToCMYK(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1, ret2, ret3 := color.RGBToCMYK(args[0].(uint8), args[1].(uint8), args[2].(uint8))
	p.Ret(3, ret0, ret1, ret2, ret3)
}

func execRGBToYCbCr(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1, ret2 := color.RGBToYCbCr(args[0].(uint8), args[1].(uint8), args[2].(uint8))
	p.Ret(3, ret0, ret1, ret2)
}

func execmYCbCrRGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2, ret3 := args[0].(color.YCbCr).RGBA()
	p.Ret(1, ret0, ret1, ret2, ret3)
}

func execYCbCrToRGB(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1, ret2 := color.YCbCrToRGB(args[0].(uint8), args[1].(uint8), args[2].(uint8))
	p.Ret(3, ret0, ret1, ret2)
}

// I is a Go package instance.
var I = gop.NewGoPackage("image/color")

func init() {
	I.RegisterFuncs(
		I.Func("(Alpha).RGBA", (color.Alpha).RGBA, execmAlphaRGBA),
		I.Func("(Alpha16).RGBA", (color.Alpha16).RGBA, execmAlpha16RGBA),
		I.Func("(CMYK).RGBA", (color.CMYK).RGBA, execmCMYKRGBA),
		I.Func("CMYKToRGB", color.CMYKToRGB, execCMYKToRGB),
		I.Func("(Color).RGBA", (color.Color).RGBA, execiColorRGBA),
		I.Func("(Gray).RGBA", (color.Gray).RGBA, execmGrayRGBA),
		I.Func("(Gray16).RGBA", (color.Gray16).RGBA, execmGray16RGBA),
		I.Func("(Model).Convert", (color.Model).Convert, execiModelConvert),
		I.Func("ModelFunc", color.ModelFunc, execModelFunc),
		I.Func("(NRGBA).RGBA", (color.NRGBA).RGBA, execmNRGBARGBA),
		I.Func("(NRGBA64).RGBA", (color.NRGBA64).RGBA, execmNRGBA64RGBA),
		I.Func("(NYCbCrA).RGBA", (color.NYCbCrA).RGBA, execmNYCbCrARGBA),
		I.Func("(Palette).Convert", (color.Palette).Convert, execmPaletteConvert),
		I.Func("(Palette).Index", (color.Palette).Index, execmPaletteIndex),
		I.Func("(RGBA).RGBA", (color.RGBA).RGBA, execmRGBARGBA),
		I.Func("(RGBA64).RGBA", (color.RGBA64).RGBA, execmRGBA64RGBA),
		I.Func("RGBToCMYK", color.RGBToCMYK, execRGBToCMYK),
		I.Func("RGBToYCbCr", color.RGBToYCbCr, execRGBToYCbCr),
		I.Func("(YCbCr).RGBA", (color.YCbCr).RGBA, execmYCbCrRGBA),
		I.Func("YCbCrToRGB", color.YCbCrToRGB, execYCbCrToRGB),
	)
	I.RegisterVars(
		I.Var("Alpha16Model", &color.Alpha16Model),
		I.Var("AlphaModel", &color.AlphaModel),
		I.Var("Black", &color.Black),
		I.Var("CMYKModel", &color.CMYKModel),
		I.Var("Gray16Model", &color.Gray16Model),
		I.Var("GrayModel", &color.GrayModel),
		I.Var("NRGBA64Model", &color.NRGBA64Model),
		I.Var("NRGBAModel", &color.NRGBAModel),
		I.Var("NYCbCrAModel", &color.NYCbCrAModel),
		I.Var("Opaque", &color.Opaque),
		I.Var("RGBA64Model", &color.RGBA64Model),
		I.Var("RGBAModel", &color.RGBAModel),
		I.Var("Transparent", &color.Transparent),
		I.Var("White", &color.White),
		I.Var("YCbCrModel", &color.YCbCrModel),
	)
	I.RegisterTypes(
		I.Type("Alpha", reflect.TypeOf((*color.Alpha)(nil)).Elem()),
		I.Type("Alpha16", reflect.TypeOf((*color.Alpha16)(nil)).Elem()),
		I.Type("CMYK", reflect.TypeOf((*color.CMYK)(nil)).Elem()),
		I.Type("Color", reflect.TypeOf((*color.Color)(nil)).Elem()),
		I.Type("Gray", reflect.TypeOf((*color.Gray)(nil)).Elem()),
		I.Type("Gray16", reflect.TypeOf((*color.Gray16)(nil)).Elem()),
		I.Type("Model", reflect.TypeOf((*color.Model)(nil)).Elem()),
		I.Type("NRGBA", reflect.TypeOf((*color.NRGBA)(nil)).Elem()),
		I.Type("NRGBA64", reflect.TypeOf((*color.NRGBA64)(nil)).Elem()),
		I.Type("NYCbCrA", reflect.TypeOf((*color.NYCbCrA)(nil)).Elem()),
		I.Type("Palette", reflect.TypeOf((*color.Palette)(nil)).Elem()),
		I.Type("RGBA", reflect.TypeOf((*color.RGBA)(nil)).Elem()),
		I.Type("RGBA64", reflect.TypeOf((*color.RGBA64)(nil)).Elem()),
		I.Type("YCbCr", reflect.TypeOf((*color.YCbCr)(nil)).Elem()),
	)
}
