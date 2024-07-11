package imager

import (
	"image"
	"image/color"
	"image/draw"
)

type additiveDrawer struct{}

func (additiveDrawer) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	sBounds := src.Bounds()
	sw, sh := src.Bounds().Dx(), src.Bounds().Dy()
	for y := 0; y < r.Dy() && (sp.Y+y) < sh; y++ {
		for x := 0; x < r.Dx() && (sp.X+x) < sw; x++ {
			dc := dst.At(r.Min.X+x, r.Min.Y+y)
			dr, dg, db, da := dc.RGBA()
			if da == 0 {
				continue
			}
			sc := src.At(sBounds.Min.X+sp.X+x, sBounds.Min.Y+sp.Y+y)
			sr, sg, sb, _ := sc.RGBA()
			dc = color.NRGBA{
				R: uint8(min(255, (dr+sr)>>8)),
				G: uint8(min(255, (dg+sg)>>8)),
				B: uint8(min(255, (db+sb)>>8)),
				A: uint8(min(255, da>>8)),
			}
			dst.Set(r.Min.X+x, r.Min.Y+y, dc)
		}
	}
}

type alphaImage struct {
	src   image.Image
	alpha uint8
}

func (img alphaImage) ColorModel() color.Model {
	return img.src.ColorModel()
}

func (img alphaImage) Bounds() image.Rectangle {
	return img.src.Bounds()
}

func (img alphaImage) At(x, y int) color.Color {
	sc := img.src.At(x, y)
	switch sc := sc.(type) {
	case *color.RGBA:
		return color.NRGBA{
			R: sc.R,
			G: sc.G,
			B: sc.B,
			A: min(sc.A, img.alpha),
		}
	default:
		sr, sg, sb, sa := sc.RGBA()
		return color.NRGBA{
			R: uint8(sr >> 8),
			G: uint8(sg >> 8),
			B: uint8(sb >> 8),
			A: min(uint8(sa>>8), img.alpha),
		}
	}
}
