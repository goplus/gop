// Package image provide Go+ "image" package, as "image" package in Go.
package image

import (
	image "image"
	color "image/color"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmAlphaColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Alpha).ColorModel()
	p.Ret(1, ret0)
}

func execmAlphaBounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Alpha).Bounds()
	p.Ret(1, ret0)
}

func execmAlphaAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Alpha).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmAlphaAlphaAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Alpha).AlphaAt(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmAlphaPixOffset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Alpha).PixOffset(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func toType0(v interface{}) color.Color {
	if v == nil {
		return nil
	}
	return v.(color.Color)
}

func execmAlphaSet(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.Alpha).Set(args[1].(int), args[2].(int), toType0(args[3]))
	p.PopN(4)
}

func execmAlphaSetAlpha(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.Alpha).SetAlpha(args[1].(int), args[2].(int), args[3].(color.Alpha))
	p.PopN(4)
}

func execmAlphaSubImage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*image.Alpha).SubImage(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmAlphaOpaque(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Alpha).Opaque()
	p.Ret(1, ret0)
}

func execmAlpha16ColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Alpha16).ColorModel()
	p.Ret(1, ret0)
}

func execmAlpha16Bounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Alpha16).Bounds()
	p.Ret(1, ret0)
}

func execmAlpha16At(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Alpha16).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmAlpha16Alpha16At(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Alpha16).Alpha16At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmAlpha16PixOffset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Alpha16).PixOffset(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmAlpha16Set(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.Alpha16).Set(args[1].(int), args[2].(int), toType0(args[3]))
	p.PopN(4)
}

func execmAlpha16SetAlpha16(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.Alpha16).SetAlpha16(args[1].(int), args[2].(int), args[3].(color.Alpha16))
	p.PopN(4)
}

func execmAlpha16SubImage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*image.Alpha16).SubImage(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmAlpha16Opaque(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Alpha16).Opaque()
	p.Ret(1, ret0)
}

func execmCMYKColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.CMYK).ColorModel()
	p.Ret(1, ret0)
}

func execmCMYKBounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.CMYK).Bounds()
	p.Ret(1, ret0)
}

func execmCMYKAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.CMYK).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmCMYKCMYKAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.CMYK).CMYKAt(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmCMYKPixOffset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.CMYK).PixOffset(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmCMYKSet(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.CMYK).Set(args[1].(int), args[2].(int), toType0(args[3]))
	p.PopN(4)
}

func execmCMYKSetCMYK(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.CMYK).SetCMYK(args[1].(int), args[2].(int), args[3].(color.CMYK))
	p.PopN(4)
}

func execmCMYKSubImage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*image.CMYK).SubImage(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmCMYKOpaque(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.CMYK).Opaque()
	p.Ret(1, ret0)
}

func toType1(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execDecode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := image.Decode(toType1(args[0]))
	p.Ret(1, ret0, ret1, ret2)
}

func execDecodeConfig(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := image.DecodeConfig(toType1(args[0]))
	p.Ret(1, ret0, ret1, ret2)
}

func execmGrayColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Gray).ColorModel()
	p.Ret(1, ret0)
}

func execmGrayBounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Gray).Bounds()
	p.Ret(1, ret0)
}

func execmGrayAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Gray).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmGrayGrayAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Gray).GrayAt(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmGrayPixOffset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Gray).PixOffset(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmGraySet(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.Gray).Set(args[1].(int), args[2].(int), toType0(args[3]))
	p.PopN(4)
}

func execmGraySetGray(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.Gray).SetGray(args[1].(int), args[2].(int), args[3].(color.Gray))
	p.PopN(4)
}

func execmGraySubImage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*image.Gray).SubImage(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmGrayOpaque(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Gray).Opaque()
	p.Ret(1, ret0)
}

func execmGray16ColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Gray16).ColorModel()
	p.Ret(1, ret0)
}

func execmGray16Bounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Gray16).Bounds()
	p.Ret(1, ret0)
}

func execmGray16At(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Gray16).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmGray16Gray16At(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Gray16).Gray16At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmGray16PixOffset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Gray16).PixOffset(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmGray16Set(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.Gray16).Set(args[1].(int), args[2].(int), toType0(args[3]))
	p.PopN(4)
}

func execmGray16SetGray16(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.Gray16).SetGray16(args[1].(int), args[2].(int), args[3].(color.Gray16))
	p.PopN(4)
}

func execmGray16SubImage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*image.Gray16).SubImage(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmGray16Opaque(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Gray16).Opaque()
	p.Ret(1, ret0)
}

func execiImageAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(image.Image).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execiImageBounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(image.Image).Bounds()
	p.Ret(1, ret0)
}

func execiImageColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(image.Image).ColorModel()
	p.Ret(1, ret0)
}

func execmNRGBAColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.NRGBA).ColorModel()
	p.Ret(1, ret0)
}

func execmNRGBABounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.NRGBA).Bounds()
	p.Ret(1, ret0)
}

func execmNRGBAAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.NRGBA).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmNRGBANRGBAAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.NRGBA).NRGBAAt(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmNRGBAPixOffset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.NRGBA).PixOffset(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmNRGBASet(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.NRGBA).Set(args[1].(int), args[2].(int), toType0(args[3]))
	p.PopN(4)
}

func execmNRGBASetNRGBA(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.NRGBA).SetNRGBA(args[1].(int), args[2].(int), args[3].(color.NRGBA))
	p.PopN(4)
}

func execmNRGBASubImage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*image.NRGBA).SubImage(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmNRGBAOpaque(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.NRGBA).Opaque()
	p.Ret(1, ret0)
}

func execmNRGBA64ColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.NRGBA64).ColorModel()
	p.Ret(1, ret0)
}

func execmNRGBA64Bounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.NRGBA64).Bounds()
	p.Ret(1, ret0)
}

func execmNRGBA64At(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.NRGBA64).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmNRGBA64NRGBA64At(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.NRGBA64).NRGBA64At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmNRGBA64PixOffset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.NRGBA64).PixOffset(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmNRGBA64Set(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.NRGBA64).Set(args[1].(int), args[2].(int), toType0(args[3]))
	p.PopN(4)
}

func execmNRGBA64SetNRGBA64(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.NRGBA64).SetNRGBA64(args[1].(int), args[2].(int), args[3].(color.NRGBA64))
	p.PopN(4)
}

func execmNRGBA64SubImage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*image.NRGBA64).SubImage(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmNRGBA64Opaque(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.NRGBA64).Opaque()
	p.Ret(1, ret0)
}

func execmNYCbCrAColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.NYCbCrA).ColorModel()
	p.Ret(1, ret0)
}

func execmNYCbCrAAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.NYCbCrA).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmNYCbCrANYCbCrAAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.NYCbCrA).NYCbCrAAt(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmNYCbCrAAOffset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.NYCbCrA).AOffset(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmNYCbCrASubImage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*image.NYCbCrA).SubImage(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmNYCbCrAOpaque(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.NYCbCrA).Opaque()
	p.Ret(1, ret0)
}

func execNewAlpha(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := image.NewAlpha(args[0].(image.Rectangle))
	p.Ret(1, ret0)
}

func execNewAlpha16(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := image.NewAlpha16(args[0].(image.Rectangle))
	p.Ret(1, ret0)
}

func execNewCMYK(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := image.NewCMYK(args[0].(image.Rectangle))
	p.Ret(1, ret0)
}

func execNewGray(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := image.NewGray(args[0].(image.Rectangle))
	p.Ret(1, ret0)
}

func execNewGray16(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := image.NewGray16(args[0].(image.Rectangle))
	p.Ret(1, ret0)
}

func execNewNRGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := image.NewNRGBA(args[0].(image.Rectangle))
	p.Ret(1, ret0)
}

func execNewNRGBA64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := image.NewNRGBA64(args[0].(image.Rectangle))
	p.Ret(1, ret0)
}

func execNewNYCbCrA(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := image.NewNYCbCrA(args[0].(image.Rectangle), args[1].(image.YCbCrSubsampleRatio))
	p.Ret(2, ret0)
}

func execNewPaletted(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := image.NewPaletted(args[0].(image.Rectangle), args[1].(color.Palette))
	p.Ret(2, ret0)
}

func execNewRGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := image.NewRGBA(args[0].(image.Rectangle))
	p.Ret(1, ret0)
}

func execNewRGBA64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := image.NewRGBA64(args[0].(image.Rectangle))
	p.Ret(1, ret0)
}

func execNewUniform(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := image.NewUniform(toType0(args[0]))
	p.Ret(1, ret0)
}

func execNewYCbCr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := image.NewYCbCr(args[0].(image.Rectangle), args[1].(image.YCbCrSubsampleRatio))
	p.Ret(2, ret0)
}

func execmPalettedColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Paletted).ColorModel()
	p.Ret(1, ret0)
}

func execmPalettedBounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Paletted).Bounds()
	p.Ret(1, ret0)
}

func execmPalettedAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Paletted).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmPalettedPixOffset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Paletted).PixOffset(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmPalettedSet(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.Paletted).Set(args[1].(int), args[2].(int), toType0(args[3]))
	p.PopN(4)
}

func execmPalettedColorIndexAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Paletted).ColorIndexAt(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmPalettedSetColorIndex(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.Paletted).SetColorIndex(args[1].(int), args[2].(int), args[3].(uint8))
	p.PopN(4)
}

func execmPalettedSubImage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*image.Paletted).SubImage(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmPalettedOpaque(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Paletted).Opaque()
	p.Ret(1, ret0)
}

func execiPalettedImageColorIndexAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(image.PalettedImage).ColorIndexAt(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmPointString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(image.Point).String()
	p.Ret(1, ret0)
}

func execmPointAdd(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Point).Add(args[1].(image.Point))
	p.Ret(2, ret0)
}

func execmPointSub(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Point).Sub(args[1].(image.Point))
	p.Ret(2, ret0)
}

func execmPointMul(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Point).Mul(args[1].(int))
	p.Ret(2, ret0)
}

func execmPointDiv(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Point).Div(args[1].(int))
	p.Ret(2, ret0)
}

func execmPointIn(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Point).In(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmPointMod(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Point).Mod(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmPointEq(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Point).Eq(args[1].(image.Point))
	p.Ret(2, ret0)
}

func execPt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := image.Pt(args[0].(int), args[1].(int))
	p.Ret(2, ret0)
}

func execmRGBAColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.RGBA).ColorModel()
	p.Ret(1, ret0)
}

func execmRGBABounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.RGBA).Bounds()
	p.Ret(1, ret0)
}

func execmRGBAAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.RGBA).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmRGBARGBAAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.RGBA).RGBAAt(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmRGBAPixOffset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.RGBA).PixOffset(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmRGBASet(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.RGBA).Set(args[1].(int), args[2].(int), toType0(args[3]))
	p.PopN(4)
}

func execmRGBASetRGBA(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.RGBA).SetRGBA(args[1].(int), args[2].(int), args[3].(color.RGBA))
	p.PopN(4)
}

func execmRGBASubImage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*image.RGBA).SubImage(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmRGBAOpaque(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.RGBA).Opaque()
	p.Ret(1, ret0)
}

func execmRGBA64ColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.RGBA64).ColorModel()
	p.Ret(1, ret0)
}

func execmRGBA64Bounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.RGBA64).Bounds()
	p.Ret(1, ret0)
}

func execmRGBA64At(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.RGBA64).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmRGBA64RGBA64At(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.RGBA64).RGBA64At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmRGBA64PixOffset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.RGBA64).PixOffset(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmRGBA64Set(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.RGBA64).Set(args[1].(int), args[2].(int), toType0(args[3]))
	p.PopN(4)
}

func execmRGBA64SetRGBA64(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*image.RGBA64).SetRGBA64(args[1].(int), args[2].(int), args[3].(color.RGBA64))
	p.PopN(4)
}

func execmRGBA64SubImage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*image.RGBA64).SubImage(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmRGBA64Opaque(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.RGBA64).Opaque()
	p.Ret(1, ret0)
}

func execRect(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := image.Rect(args[0].(int), args[1].(int), args[2].(int), args[3].(int))
	p.Ret(4, ret0)
}

func execmRectangleString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(image.Rectangle).String()
	p.Ret(1, ret0)
}

func execmRectangleDx(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(image.Rectangle).Dx()
	p.Ret(1, ret0)
}

func execmRectangleDy(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(image.Rectangle).Dy()
	p.Ret(1, ret0)
}

func execmRectangleSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(image.Rectangle).Size()
	p.Ret(1, ret0)
}

func execmRectangleAdd(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Rectangle).Add(args[1].(image.Point))
	p.Ret(2, ret0)
}

func execmRectangleSub(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Rectangle).Sub(args[1].(image.Point))
	p.Ret(2, ret0)
}

func execmRectangleInset(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Rectangle).Inset(args[1].(int))
	p.Ret(2, ret0)
}

func execmRectangleIntersect(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Rectangle).Intersect(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmRectangleUnion(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Rectangle).Union(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmRectangleEmpty(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(image.Rectangle).Empty()
	p.Ret(1, ret0)
}

func execmRectangleEq(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Rectangle).Eq(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmRectangleOverlaps(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Rectangle).Overlaps(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmRectangleIn(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(image.Rectangle).In(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmRectangleCanon(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(image.Rectangle).Canon()
	p.Ret(1, ret0)
}

func execmRectangleAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(image.Rectangle).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmRectangleBounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(image.Rectangle).Bounds()
	p.Ret(1, ret0)
}

func execmRectangleColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(image.Rectangle).ColorModel()
	p.Ret(1, ret0)
}

func execRegisterFormat(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	image.RegisterFormat(args[0].(string), args[1].(string), args[2].(func(io.Reader) (image.Image, error)), args[3].(func(io.Reader) (image.Config, error)))
	p.PopN(4)
}

func execmUniformRGBA(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2, ret3 := args[0].(*image.Uniform).RGBA()
	p.Ret(1, ret0, ret1, ret2, ret3)
}

func execmUniformColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Uniform).ColorModel()
	p.Ret(1, ret0)
}

func execmUniformConvert(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*image.Uniform).Convert(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmUniformBounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Uniform).Bounds()
	p.Ret(1, ret0)
}

func execmUniformAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.Uniform).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmUniformOpaque(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.Uniform).Opaque()
	p.Ret(1, ret0)
}

func execmYCbCrColorModel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.YCbCr).ColorModel()
	p.Ret(1, ret0)
}

func execmYCbCrBounds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.YCbCr).Bounds()
	p.Ret(1, ret0)
}

func execmYCbCrAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.YCbCr).At(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmYCbCrYCbCrAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.YCbCr).YCbCrAt(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmYCbCrYOffset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.YCbCr).YOffset(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmYCbCrCOffset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*image.YCbCr).COffset(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmYCbCrSubImage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*image.YCbCr).SubImage(args[1].(image.Rectangle))
	p.Ret(2, ret0)
}

func execmYCbCrOpaque(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*image.YCbCr).Opaque()
	p.Ret(1, ret0)
}

func execmYCbCrSubsampleRatioString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(image.YCbCrSubsampleRatio).String()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("image")

func init() {
	I.RegisterFuncs(
		I.Func("(*Alpha).ColorModel", (*image.Alpha).ColorModel, execmAlphaColorModel),
		I.Func("(*Alpha).Bounds", (*image.Alpha).Bounds, execmAlphaBounds),
		I.Func("(*Alpha).At", (*image.Alpha).At, execmAlphaAt),
		I.Func("(*Alpha).AlphaAt", (*image.Alpha).AlphaAt, execmAlphaAlphaAt),
		I.Func("(*Alpha).PixOffset", (*image.Alpha).PixOffset, execmAlphaPixOffset),
		I.Func("(*Alpha).Set", (*image.Alpha).Set, execmAlphaSet),
		I.Func("(*Alpha).SetAlpha", (*image.Alpha).SetAlpha, execmAlphaSetAlpha),
		I.Func("(*Alpha).SubImage", (*image.Alpha).SubImage, execmAlphaSubImage),
		I.Func("(*Alpha).Opaque", (*image.Alpha).Opaque, execmAlphaOpaque),
		I.Func("(*Alpha16).ColorModel", (*image.Alpha16).ColorModel, execmAlpha16ColorModel),
		I.Func("(*Alpha16).Bounds", (*image.Alpha16).Bounds, execmAlpha16Bounds),
		I.Func("(*Alpha16).At", (*image.Alpha16).At, execmAlpha16At),
		I.Func("(*Alpha16).Alpha16At", (*image.Alpha16).Alpha16At, execmAlpha16Alpha16At),
		I.Func("(*Alpha16).PixOffset", (*image.Alpha16).PixOffset, execmAlpha16PixOffset),
		I.Func("(*Alpha16).Set", (*image.Alpha16).Set, execmAlpha16Set),
		I.Func("(*Alpha16).SetAlpha16", (*image.Alpha16).SetAlpha16, execmAlpha16SetAlpha16),
		I.Func("(*Alpha16).SubImage", (*image.Alpha16).SubImage, execmAlpha16SubImage),
		I.Func("(*Alpha16).Opaque", (*image.Alpha16).Opaque, execmAlpha16Opaque),
		I.Func("(*CMYK).ColorModel", (*image.CMYK).ColorModel, execmCMYKColorModel),
		I.Func("(*CMYK).Bounds", (*image.CMYK).Bounds, execmCMYKBounds),
		I.Func("(*CMYK).At", (*image.CMYK).At, execmCMYKAt),
		I.Func("(*CMYK).CMYKAt", (*image.CMYK).CMYKAt, execmCMYKCMYKAt),
		I.Func("(*CMYK).PixOffset", (*image.CMYK).PixOffset, execmCMYKPixOffset),
		I.Func("(*CMYK).Set", (*image.CMYK).Set, execmCMYKSet),
		I.Func("(*CMYK).SetCMYK", (*image.CMYK).SetCMYK, execmCMYKSetCMYK),
		I.Func("(*CMYK).SubImage", (*image.CMYK).SubImage, execmCMYKSubImage),
		I.Func("(*CMYK).Opaque", (*image.CMYK).Opaque, execmCMYKOpaque),
		I.Func("Decode", image.Decode, execDecode),
		I.Func("DecodeConfig", image.DecodeConfig, execDecodeConfig),
		I.Func("(*Gray).ColorModel", (*image.Gray).ColorModel, execmGrayColorModel),
		I.Func("(*Gray).Bounds", (*image.Gray).Bounds, execmGrayBounds),
		I.Func("(*Gray).At", (*image.Gray).At, execmGrayAt),
		I.Func("(*Gray).GrayAt", (*image.Gray).GrayAt, execmGrayGrayAt),
		I.Func("(*Gray).PixOffset", (*image.Gray).PixOffset, execmGrayPixOffset),
		I.Func("(*Gray).Set", (*image.Gray).Set, execmGraySet),
		I.Func("(*Gray).SetGray", (*image.Gray).SetGray, execmGraySetGray),
		I.Func("(*Gray).SubImage", (*image.Gray).SubImage, execmGraySubImage),
		I.Func("(*Gray).Opaque", (*image.Gray).Opaque, execmGrayOpaque),
		I.Func("(*Gray16).ColorModel", (*image.Gray16).ColorModel, execmGray16ColorModel),
		I.Func("(*Gray16).Bounds", (*image.Gray16).Bounds, execmGray16Bounds),
		I.Func("(*Gray16).At", (*image.Gray16).At, execmGray16At),
		I.Func("(*Gray16).Gray16At", (*image.Gray16).Gray16At, execmGray16Gray16At),
		I.Func("(*Gray16).PixOffset", (*image.Gray16).PixOffset, execmGray16PixOffset),
		I.Func("(*Gray16).Set", (*image.Gray16).Set, execmGray16Set),
		I.Func("(*Gray16).SetGray16", (*image.Gray16).SetGray16, execmGray16SetGray16),
		I.Func("(*Gray16).SubImage", (*image.Gray16).SubImage, execmGray16SubImage),
		I.Func("(*Gray16).Opaque", (*image.Gray16).Opaque, execmGray16Opaque),
		I.Func("(Image).At", (image.Image).At, execiImageAt),
		I.Func("(Image).Bounds", (image.Image).Bounds, execiImageBounds),
		I.Func("(Image).ColorModel", (image.Image).ColorModel, execiImageColorModel),
		I.Func("(*NRGBA).ColorModel", (*image.NRGBA).ColorModel, execmNRGBAColorModel),
		I.Func("(*NRGBA).Bounds", (*image.NRGBA).Bounds, execmNRGBABounds),
		I.Func("(*NRGBA).At", (*image.NRGBA).At, execmNRGBAAt),
		I.Func("(*NRGBA).NRGBAAt", (*image.NRGBA).NRGBAAt, execmNRGBANRGBAAt),
		I.Func("(*NRGBA).PixOffset", (*image.NRGBA).PixOffset, execmNRGBAPixOffset),
		I.Func("(*NRGBA).Set", (*image.NRGBA).Set, execmNRGBASet),
		I.Func("(*NRGBA).SetNRGBA", (*image.NRGBA).SetNRGBA, execmNRGBASetNRGBA),
		I.Func("(*NRGBA).SubImage", (*image.NRGBA).SubImage, execmNRGBASubImage),
		I.Func("(*NRGBA).Opaque", (*image.NRGBA).Opaque, execmNRGBAOpaque),
		I.Func("(*NRGBA64).ColorModel", (*image.NRGBA64).ColorModel, execmNRGBA64ColorModel),
		I.Func("(*NRGBA64).Bounds", (*image.NRGBA64).Bounds, execmNRGBA64Bounds),
		I.Func("(*NRGBA64).At", (*image.NRGBA64).At, execmNRGBA64At),
		I.Func("(*NRGBA64).NRGBA64At", (*image.NRGBA64).NRGBA64At, execmNRGBA64NRGBA64At),
		I.Func("(*NRGBA64).PixOffset", (*image.NRGBA64).PixOffset, execmNRGBA64PixOffset),
		I.Func("(*NRGBA64).Set", (*image.NRGBA64).Set, execmNRGBA64Set),
		I.Func("(*NRGBA64).SetNRGBA64", (*image.NRGBA64).SetNRGBA64, execmNRGBA64SetNRGBA64),
		I.Func("(*NRGBA64).SubImage", (*image.NRGBA64).SubImage, execmNRGBA64SubImage),
		I.Func("(*NRGBA64).Opaque", (*image.NRGBA64).Opaque, execmNRGBA64Opaque),
		I.Func("(*NYCbCrA).ColorModel", (*image.NYCbCrA).ColorModel, execmNYCbCrAColorModel),
		I.Func("(*NYCbCrA).At", (*image.NYCbCrA).At, execmNYCbCrAAt),
		I.Func("(*NYCbCrA).NYCbCrAAt", (*image.NYCbCrA).NYCbCrAAt, execmNYCbCrANYCbCrAAt),
		I.Func("(*NYCbCrA).AOffset", (*image.NYCbCrA).AOffset, execmNYCbCrAAOffset),
		I.Func("(*NYCbCrA).SubImage", (*image.NYCbCrA).SubImage, execmNYCbCrASubImage),
		I.Func("(*NYCbCrA).Opaque", (*image.NYCbCrA).Opaque, execmNYCbCrAOpaque),
		I.Func("NewAlpha", image.NewAlpha, execNewAlpha),
		I.Func("NewAlpha16", image.NewAlpha16, execNewAlpha16),
		I.Func("NewCMYK", image.NewCMYK, execNewCMYK),
		I.Func("NewGray", image.NewGray, execNewGray),
		I.Func("NewGray16", image.NewGray16, execNewGray16),
		I.Func("NewNRGBA", image.NewNRGBA, execNewNRGBA),
		I.Func("NewNRGBA64", image.NewNRGBA64, execNewNRGBA64),
		I.Func("NewNYCbCrA", image.NewNYCbCrA, execNewNYCbCrA),
		I.Func("NewPaletted", image.NewPaletted, execNewPaletted),
		I.Func("NewRGBA", image.NewRGBA, execNewRGBA),
		I.Func("NewRGBA64", image.NewRGBA64, execNewRGBA64),
		I.Func("NewUniform", image.NewUniform, execNewUniform),
		I.Func("NewYCbCr", image.NewYCbCr, execNewYCbCr),
		I.Func("(*Paletted).ColorModel", (*image.Paletted).ColorModel, execmPalettedColorModel),
		I.Func("(*Paletted).Bounds", (*image.Paletted).Bounds, execmPalettedBounds),
		I.Func("(*Paletted).At", (*image.Paletted).At, execmPalettedAt),
		I.Func("(*Paletted).PixOffset", (*image.Paletted).PixOffset, execmPalettedPixOffset),
		I.Func("(*Paletted).Set", (*image.Paletted).Set, execmPalettedSet),
		I.Func("(*Paletted).ColorIndexAt", (*image.Paletted).ColorIndexAt, execmPalettedColorIndexAt),
		I.Func("(*Paletted).SetColorIndex", (*image.Paletted).SetColorIndex, execmPalettedSetColorIndex),
		I.Func("(*Paletted).SubImage", (*image.Paletted).SubImage, execmPalettedSubImage),
		I.Func("(*Paletted).Opaque", (*image.Paletted).Opaque, execmPalettedOpaque),
		I.Func("(PalettedImage).At", (image.Image).At, execiImageAt),
		I.Func("(PalettedImage).Bounds", (image.Image).Bounds, execiImageBounds),
		I.Func("(PalettedImage).ColorIndexAt", (image.PalettedImage).ColorIndexAt, execiPalettedImageColorIndexAt),
		I.Func("(PalettedImage).ColorModel", (image.Image).ColorModel, execiImageColorModel),
		I.Func("(Point).String", (image.Point).String, execmPointString),
		I.Func("(Point).Add", (image.Point).Add, execmPointAdd),
		I.Func("(Point).Sub", (image.Point).Sub, execmPointSub),
		I.Func("(Point).Mul", (image.Point).Mul, execmPointMul),
		I.Func("(Point).Div", (image.Point).Div, execmPointDiv),
		I.Func("(Point).In", (image.Point).In, execmPointIn),
		I.Func("(Point).Mod", (image.Point).Mod, execmPointMod),
		I.Func("(Point).Eq", (image.Point).Eq, execmPointEq),
		I.Func("Pt", image.Pt, execPt),
		I.Func("(*RGBA).ColorModel", (*image.RGBA).ColorModel, execmRGBAColorModel),
		I.Func("(*RGBA).Bounds", (*image.RGBA).Bounds, execmRGBABounds),
		I.Func("(*RGBA).At", (*image.RGBA).At, execmRGBAAt),
		I.Func("(*RGBA).RGBAAt", (*image.RGBA).RGBAAt, execmRGBARGBAAt),
		I.Func("(*RGBA).PixOffset", (*image.RGBA).PixOffset, execmRGBAPixOffset),
		I.Func("(*RGBA).Set", (*image.RGBA).Set, execmRGBASet),
		I.Func("(*RGBA).SetRGBA", (*image.RGBA).SetRGBA, execmRGBASetRGBA),
		I.Func("(*RGBA).SubImage", (*image.RGBA).SubImage, execmRGBASubImage),
		I.Func("(*RGBA).Opaque", (*image.RGBA).Opaque, execmRGBAOpaque),
		I.Func("(*RGBA64).ColorModel", (*image.RGBA64).ColorModel, execmRGBA64ColorModel),
		I.Func("(*RGBA64).Bounds", (*image.RGBA64).Bounds, execmRGBA64Bounds),
		I.Func("(*RGBA64).At", (*image.RGBA64).At, execmRGBA64At),
		I.Func("(*RGBA64).RGBA64At", (*image.RGBA64).RGBA64At, execmRGBA64RGBA64At),
		I.Func("(*RGBA64).PixOffset", (*image.RGBA64).PixOffset, execmRGBA64PixOffset),
		I.Func("(*RGBA64).Set", (*image.RGBA64).Set, execmRGBA64Set),
		I.Func("(*RGBA64).SetRGBA64", (*image.RGBA64).SetRGBA64, execmRGBA64SetRGBA64),
		I.Func("(*RGBA64).SubImage", (*image.RGBA64).SubImage, execmRGBA64SubImage),
		I.Func("(*RGBA64).Opaque", (*image.RGBA64).Opaque, execmRGBA64Opaque),
		I.Func("Rect", image.Rect, execRect),
		I.Func("(Rectangle).String", (image.Rectangle).String, execmRectangleString),
		I.Func("(Rectangle).Dx", (image.Rectangle).Dx, execmRectangleDx),
		I.Func("(Rectangle).Dy", (image.Rectangle).Dy, execmRectangleDy),
		I.Func("(Rectangle).Size", (image.Rectangle).Size, execmRectangleSize),
		I.Func("(Rectangle).Add", (image.Rectangle).Add, execmRectangleAdd),
		I.Func("(Rectangle).Sub", (image.Rectangle).Sub, execmRectangleSub),
		I.Func("(Rectangle).Inset", (image.Rectangle).Inset, execmRectangleInset),
		I.Func("(Rectangle).Intersect", (image.Rectangle).Intersect, execmRectangleIntersect),
		I.Func("(Rectangle).Union", (image.Rectangle).Union, execmRectangleUnion),
		I.Func("(Rectangle).Empty", (image.Rectangle).Empty, execmRectangleEmpty),
		I.Func("(Rectangle).Eq", (image.Rectangle).Eq, execmRectangleEq),
		I.Func("(Rectangle).Overlaps", (image.Rectangle).Overlaps, execmRectangleOverlaps),
		I.Func("(Rectangle).In", (image.Rectangle).In, execmRectangleIn),
		I.Func("(Rectangle).Canon", (image.Rectangle).Canon, execmRectangleCanon),
		I.Func("(Rectangle).At", (image.Rectangle).At, execmRectangleAt),
		I.Func("(Rectangle).Bounds", (image.Rectangle).Bounds, execmRectangleBounds),
		I.Func("(Rectangle).ColorModel", (image.Rectangle).ColorModel, execmRectangleColorModel),
		I.Func("RegisterFormat", image.RegisterFormat, execRegisterFormat),
		I.Func("(*Uniform).RGBA", (*image.Uniform).RGBA, execmUniformRGBA),
		I.Func("(*Uniform).ColorModel", (*image.Uniform).ColorModel, execmUniformColorModel),
		I.Func("(*Uniform).Convert", (*image.Uniform).Convert, execmUniformConvert),
		I.Func("(*Uniform).Bounds", (*image.Uniform).Bounds, execmUniformBounds),
		I.Func("(*Uniform).At", (*image.Uniform).At, execmUniformAt),
		I.Func("(*Uniform).Opaque", (*image.Uniform).Opaque, execmUniformOpaque),
		I.Func("(*YCbCr).ColorModel", (*image.YCbCr).ColorModel, execmYCbCrColorModel),
		I.Func("(*YCbCr).Bounds", (*image.YCbCr).Bounds, execmYCbCrBounds),
		I.Func("(*YCbCr).At", (*image.YCbCr).At, execmYCbCrAt),
		I.Func("(*YCbCr).YCbCrAt", (*image.YCbCr).YCbCrAt, execmYCbCrYCbCrAt),
		I.Func("(*YCbCr).YOffset", (*image.YCbCr).YOffset, execmYCbCrYOffset),
		I.Func("(*YCbCr).COffset", (*image.YCbCr).COffset, execmYCbCrCOffset),
		I.Func("(*YCbCr).SubImage", (*image.YCbCr).SubImage, execmYCbCrSubImage),
		I.Func("(*YCbCr).Opaque", (*image.YCbCr).Opaque, execmYCbCrOpaque),
		I.Func("(YCbCrSubsampleRatio).String", (image.YCbCrSubsampleRatio).String, execmYCbCrSubsampleRatioString),
	)
	I.RegisterVars(
		I.Var("Black", &image.Black),
		I.Var("ErrFormat", &image.ErrFormat),
		I.Var("Opaque", &image.Opaque),
		I.Var("Transparent", &image.Transparent),
		I.Var("White", &image.White),
		I.Var("ZP", &image.ZP),
		I.Var("ZR", &image.ZR),
	)
	I.RegisterTypes(
		I.Type("Alpha", reflect.TypeOf((*image.Alpha)(nil)).Elem()),
		I.Type("Alpha16", reflect.TypeOf((*image.Alpha16)(nil)).Elem()),
		I.Type("CMYK", reflect.TypeOf((*image.CMYK)(nil)).Elem()),
		I.Type("Config", reflect.TypeOf((*image.Config)(nil)).Elem()),
		I.Type("Gray", reflect.TypeOf((*image.Gray)(nil)).Elem()),
		I.Type("Gray16", reflect.TypeOf((*image.Gray16)(nil)).Elem()),
		I.Type("Image", reflect.TypeOf((*image.Image)(nil)).Elem()),
		I.Type("NRGBA", reflect.TypeOf((*image.NRGBA)(nil)).Elem()),
		I.Type("NRGBA64", reflect.TypeOf((*image.NRGBA64)(nil)).Elem()),
		I.Type("NYCbCrA", reflect.TypeOf((*image.NYCbCrA)(nil)).Elem()),
		I.Type("Paletted", reflect.TypeOf((*image.Paletted)(nil)).Elem()),
		I.Type("PalettedImage", reflect.TypeOf((*image.PalettedImage)(nil)).Elem()),
		I.Type("Point", reflect.TypeOf((*image.Point)(nil)).Elem()),
		I.Type("RGBA", reflect.TypeOf((*image.RGBA)(nil)).Elem()),
		I.Type("RGBA64", reflect.TypeOf((*image.RGBA64)(nil)).Elem()),
		I.Type("Rectangle", reflect.TypeOf((*image.Rectangle)(nil)).Elem()),
		I.Type("Uniform", reflect.TypeOf((*image.Uniform)(nil)).Elem()),
		I.Type("YCbCr", reflect.TypeOf((*image.YCbCr)(nil)).Elem()),
		I.Type("YCbCrSubsampleRatio", reflect.TypeOf((*image.YCbCrSubsampleRatio)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("YCbCrSubsampleRatio410", qspec.Int, image.YCbCrSubsampleRatio410),
		I.Const("YCbCrSubsampleRatio411", qspec.Int, image.YCbCrSubsampleRatio411),
		I.Const("YCbCrSubsampleRatio420", qspec.Int, image.YCbCrSubsampleRatio420),
		I.Const("YCbCrSubsampleRatio422", qspec.Int, image.YCbCrSubsampleRatio422),
		I.Const("YCbCrSubsampleRatio440", qspec.Int, image.YCbCrSubsampleRatio440),
		I.Const("YCbCrSubsampleRatio444", qspec.Int, image.YCbCrSubsampleRatio444),
	)
}
